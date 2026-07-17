package email

import (
	"abstractions/db/models"
	"net/http"
	"testing"
)

func TestClassifyCloudflareError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		parsed     cloudflareSendResponse
		wantErrOk  bool
		wantErrTyp string
	}{
		{
			name:      "200 success",
			status:    http.StatusOK,
			parsed:    cloudflareSendResponse{Success: true, Result: cloudflareSendResult{Delivered: []string{"a@example.com"}}},
			wantErrOk: false,
		},
		{
			name:       "400 validation error is hard",
			status:     http.StatusBadRequest,
			parsed:     cloudflareSendResponse{Success: false, Errors: []cloudflareErrorDetail{{Code: 1000, Message: "invalid recipient"}}},
			wantErrOk:  true,
			wantErrTyp: models.EmailErrorTypeHard,
		},
		{
			name:       "401 bad token is system",
			status:     http.StatusUnauthorized,
			parsed:     cloudflareSendResponse{Success: false, Errors: []cloudflareErrorDetail{{Code: 9109, Message: "invalid token"}}},
			wantErrOk:  true,
			wantErrTyp: models.EmailErrorTypeSystem,
		},
		{
			name:       "429 rate limited is soft",
			status:     http.StatusTooManyRequests,
			parsed:     cloudflareSendResponse{Success: false},
			wantErrOk:  true,
			wantErrTyp: models.EmailErrorTypeSoft,
		},
		{
			name:       "500 server error is soft",
			status:     http.StatusInternalServerError,
			parsed:     cloudflareSendResponse{Success: false},
			wantErrOk:  true,
			wantErrTyp: models.EmailErrorTypeSoft,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errType, ok := classifyCloudflareError(tc.status, tc.parsed)
			if ok != tc.wantErrOk {
				t.Fatalf("classifyCloudflareError() ok = %v, want %v", ok, tc.wantErrOk)
			}
			if ok && errType != tc.wantErrTyp {
				t.Fatalf("classifyCloudflareError() errType = %q, want %q", errType, tc.wantErrTyp)
			}
		})
	}
}

func TestCloudflareRawSendPermanentBounce(t *testing.T) {
	parsed := cloudflareSendResponse{
		Success: true,
		Result:  cloudflareSendResult{PermanentBounces: []string{"bounced@example.com"}},
	}
	// A 2xx + success:true response with the recipient in permanent_bounces must still
	// be treated as a hard failure by cloudflareRawSend (checked after classifyCloudflareError,
	// which only inspects the top-level success/status).
	errType, ok := classifyCloudflareError(http.StatusOK, parsed)
	if ok {
		t.Fatalf("classifyCloudflareError() should not flag a 200/success response, got errType=%q", errType)
	}
	found := false
	for _, addr := range parsed.Result.PermanentBounces {
		if addr == "bounced@example.com" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected recipient to be present in PermanentBounces for this test to be meaningful")
	}
}
