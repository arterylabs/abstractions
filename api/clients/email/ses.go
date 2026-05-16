package email

import (

	//go get -u github.com/aws/aws-sdk-go

	"abstractions/config"
	"abstractions/logger"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var service *ses.SES

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

func SendEmailUsingSES(ctx context.Context, params AWSEmailParams) {
	// the @test.arterylabs.com domain is used for testing purposes and should not actually send emails
	if strings.HasSuffix(params.Email, "@test.arterylabs.com") {
		return
	}

	svc := GetSESService()

	email := params.Email
	if config.Environment.IsDevelopment() {
		email = config.Environment.DevelopmentEmail
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(email),
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

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	logger.Info(ctx, "Email Sent to address: "+email)
	logger.Info(ctx, result)
}
