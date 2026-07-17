package email

import (
	"abstractions/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"abstractions/db/models"
)

const cloudflareSendURLFormat = "https://api.cloudflare.com/client/v4/accounts/%s/email/sending/send"

var cloudflareHTTPClient = &http.Client{Timeout: 15 * time.Second}

// cloudflareSendRequest is the request body for POST /accounts/{account_id}/email/sending/send.
type cloudflareSendRequest struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	HTML    string `json:"html,omitempty"`
	Text    string `json:"text,omitempty"`
}

// cloudflareSendResponse is the response body for a send request, success or failure.
type cloudflareSendResponse struct {
	Success bool                    `json:"success"`
	Errors  []cloudflareErrorDetail `json:"errors"`
	Result  cloudflareSendResult    `json:"result"`
}

type cloudflareErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cloudflareSendResult struct {
	Delivered        []string `json:"delivered"`
	PermanentBounces []string `json:"permanent_bounces"`
	Queued           []string `json:"queued"`
}

// SendEmailUsingCloudflare sends a single email via the Cloudflare Email Service REST
// API. It shares the same pre/post-send handling (blocklist, suppression, logging) as
// SendEmailUsingSES via sendEmail.
func SendEmailUsingCloudflare(ctx context.Context, params AWSEmailParams) error {
	return sendEmail(ctx, params, cloudflareRawSend)
}

// cloudflareRawSend performs the actual Cloudflare send request for the resolved recipient.
func cloudflareRawSend(ctx context.Context, recipient string, params AWSEmailParams) (string, string, error) {
	reqBody := cloudflareSendRequest{
		To:      recipient,
		From:    fmt.Sprintf("%s <%s>", SenderDetails.Name, SenderDetails.Email),
		Subject: params.Subject,
		HTML:    params.BodyHTML,
		Text:    params.BodyText,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", models.EmailErrorTypeSystem, fmt.Errorf("cloudflare: marshal request: %w", err)
	}

	url := fmt.Sprintf(cloudflareSendURLFormat, config.Environment.CloudflareAccountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", models.EmailErrorTypeSystem, fmt.Errorf("cloudflare: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Environment.CloudflareAPIToken)

	resp, err := cloudflareHTTPClient.Do(req)
	if err != nil {
		return "", models.EmailErrorTypeSoft, fmt.Errorf("cloudflare: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", models.EmailErrorTypeSoft, fmt.Errorf("cloudflare: read response: %w", err)
	}

	var parsed cloudflareSendResponse
	// A malformed body is a system problem (our parsing), not the recipient's fault.
	if jsonErr := json.Unmarshal(respBody, &parsed); jsonErr != nil {
		return "", models.EmailErrorTypeSystem, fmt.Errorf("cloudflare: decode response (status %d): %w", resp.StatusCode, jsonErr)
	}

	if errType, ok := classifyCloudflareError(resp.StatusCode, parsed); ok {
		return "", errType, cloudflareError(resp.StatusCode, parsed)
	}

	for _, addr := range parsed.Result.PermanentBounces {
		if addr == recipient {
			return "", models.EmailErrorTypeHard, fmt.Errorf("cloudflare: recipient permanently bounced: %s", recipient)
		}
	}

	return "", "", nil
}

// classifyCloudflareError maps a Cloudflare response to an EmailErrorType. The second
// return value is false when the send was successful and no classification is needed.
//
//   - HTTP 400 (validation): the recipient/request was rejected -> hard failure.
//   - HTTP 401 (bad token): our own misconfiguration -> system, never counts against the address.
//   - HTTP 429 / 500 / anything else non-2xx: transient -> soft failure.
func classifyCloudflareError(status int, parsed cloudflareSendResponse) (string, bool) {
	if status >= 200 && status < 300 && parsed.Success {
		return "", false
	}
	switch status {
	case http.StatusBadRequest:
		return models.EmailErrorTypeHard, true
	case http.StatusUnauthorized:
		return models.EmailErrorTypeSystem, true
	default:
		return models.EmailErrorTypeSoft, true
	}
}

func cloudflareError(status int, parsed cloudflareSendResponse) error {
	if len(parsed.Errors) > 0 {
		return fmt.Errorf("cloudflare: send failed (status %d): %s", status, parsed.Errors[0].Message)
	}
	return fmt.Errorf("cloudflare: send failed (status %d)", status)
}
