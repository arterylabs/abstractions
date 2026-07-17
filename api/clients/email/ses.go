package email

import (

	//go get -u github.com/aws/aws-sdk-go

	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/utils/email_utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var service *ses.SES

// ErrEmailSuppressed is returned when a send is skipped because the address is on the
// suppression list. ErrEmailBlockedDomain is returned for a disposable/blocked domain.
var (
	ErrEmailSuppressed    = errors.New("email: recipient is suppressed")
	ErrEmailBlockedDomain = errors.New("email: recipient domain is blocked")
)

func CreateSESService() *ses.SES {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.Environment.AWSRegion)},
	)
	sess.Config.Credentials = credentials.NewStaticCredentials(config.Environment.AWSAccessKey, config.Environment.AWSSecretKey, "")
	return ses.New(sess)
}

func GetSESService() *ses.SES {
	if service == nil {
		service = CreateSESService()
	}
	return service
}

type AWSEmailParams struct {
	Email    string
	Subject  string
	BodyHTML string
	BodyText string
}

var CharSet = "UTF-8"

// rawSender performs the actual provider network call for a single, already-resolved
// recipient. It returns the provider's message ID on success, or an EmailErrorType
// (models.EmailErrorType*) describing the failure.
type rawSender func(ctx context.Context, recipient string, params AWSEmailParams) (messageID string, errType string, err error)

// sendEmail is the provider-agnostic send shell shared by every provider: the
// @test.arterylabs.com skip, domain blocklist check, suppression check, dev-email
// override, logging, and post-failure suppression handling. Providers only need to
// supply a rawSender that performs the network call.
func sendEmail(ctx context.Context, params AWSEmailParams, raw rawSender) error {
	// the @test.arterylabs.com domain is used for testing purposes and should not actually send emails
	if strings.HasSuffix(params.Email, "@test.arterylabs.com") {
		return nil
	}

	// Pre-send checks are evaluated against the intended recipient (params.Email),
	// not the development override address.
	if email_utils.IsEmailDomainBlocked(params.Email) {
		logSend(ctx, EmailLogEntry{
			Email:        params.Email,
			Subject:      params.Subject,
			Status:       models.EmailLogStatusSuppressed,
			ErrorType:    models.EmailErrorTypeBlockedDomain,
			ErrorMessage: "recipient domain is on the blocklist",
		})
		logger.Info(ctx, "email: skipped send to blocked domain: "+params.Email)
		return ErrEmailBlockedDomain
	}

	if store != nil {
		suppressed, err := store.IsSuppressed(ctx, params.Email)
		if err != nil {
			logger.Error(ctx, "email: failed to check suppression list for "+params.Email, err)
		} else if suppressed {
			logSend(ctx, EmailLogEntry{
				Email:     params.Email,
				Subject:   params.Subject,
				Status:    models.EmailLogStatusSuppressed,
				ErrorType: models.EmailErrorTypeSuppressed,
			})
			logger.Info(ctx, "email: skipped send to suppressed address: "+params.Email)
			return ErrEmailSuppressed
		}
	}

	recipient := params.Email
	if config.Environment.IsDevelopment() {
		recipient = config.Environment.DevelopmentEmail
	}

	messageID, errType, err := raw(ctx, recipient, params)
	if err != nil {
		logSend(ctx, EmailLogEntry{
			Email:        params.Email,
			Subject:      params.Subject,
			Status:       models.EmailLogStatusFailed,
			ErrorType:    errType,
			ErrorMessage: err.Error(),
		})
		handleSendFailure(ctx, params.Email, errType)
		logger.Error(ctx, "email: send failed for "+params.Email+" ("+errType+")", err)
		return err
	}

	logSend(ctx, EmailLogEntry{
		Email:     params.Email,
		Subject:   params.Subject,
		MessageID: messageID,
		Status:    models.EmailLogStatusSent,
	})
	logger.Info(ctx, "Email Sent to address: "+recipient)
	return nil
}

// SendEmailUsingSES sends a single email via SES. Before sending it checks the domain
// blocklist and the suppression list; after sending it records the outcome in the email
// log and, on failure, updates the suppression list. All logging/suppression is best-effort
// (fails open) when no store has been registered.
func SendEmailUsingSES(ctx context.Context, params AWSEmailParams) error {
	return sendEmail(ctx, params, sesRawSend)
}

// sesRawSend performs the actual SES SendEmail call for the resolved recipient.
func sesRawSend(ctx context.Context, recipient string, params AWSEmailParams) (string, string, error) {
	svc := GetSESService()

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(params.BodyHTML),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(params.BodyText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(params.Subject),
			},
		},
		Source: aws.String(fmt.Sprintf("%s<%s>", SenderDetails.Name, SenderDetails.Email)),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		return "", classifySESError(err), err
	}

	messageID := ""
	if result != nil && result.MessageId != nil {
		messageID = *result.MessageId
	}
	return messageID, "", nil
}

// classifySESError maps an SES send error to an error type that drives suppression.
// Hard errors are the recipient's fault (invalid address); system errors are our own
// misconfiguration and must not count against the address; everything else is transient.
func classifySESError(err error) string {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			return models.EmailErrorTypeHard
		case ses.ErrCodeMailFromDomainNotVerifiedException, ses.ErrCodeConfigurationSetDoesNotExistException:
			return models.EmailErrorTypeSystem
		}
	}
	return models.EmailErrorTypeSoft
}

// handleSendFailure updates the suppression list based on the failure type: hard bounces
// suppress immediately, soft failures suppress once they cross FailureThreshold, and system
// errors never count against the address.
func handleSendFailure(ctx context.Context, email string, errType string) {
	if store == nil {
		return
	}
	switch errType {
	case models.EmailErrorTypeHard:
		suppress(ctx, email, models.SuppressionReasonHardBounce)
	case models.EmailErrorTypeSoft:
		count, err := store.CountSoftFailures(ctx, email)
		if err != nil {
			logger.Error(ctx, "email: failed to count soft failures for "+email, err)
			return
		}
		if count >= FailureThreshold {
			suppress(ctx, email, models.SuppressionReasonThresholdExceeded)
		}
	}
}

func suppress(ctx context.Context, email string, reason string) {
	if store == nil {
		return
	}
	if err := store.Suppress(ctx, email, reason); err != nil {
		logger.Error(ctx, "email: failed to suppress "+email, err)
	}
}

func logSend(ctx context.Context, entry EmailLogEntry) {
	if store == nil {
		return
	}
	if err := store.LogSend(ctx, entry); err != nil {
		logger.Error(ctx, "email: failed to write email log for "+entry.Email, err)
	}
}
