package email

import (
	"abstractions/db/models"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
)

// fakeStore is an in-memory SuppressionStore for exercising the send-failure logic
// without a database or a live SES connection.
type fakeStore struct {
	suppressed map[string]string
	softCount  int
	softErr    error
}

func newFakeStore() *fakeStore {
	return &fakeStore{suppressed: map[string]string{}}
}

func (f *fakeStore) IsSuppressed(ctx context.Context, email string) (bool, error) {
	_, ok := f.suppressed[email]
	return ok, nil
}
func (f *fakeStore) Suppress(ctx context.Context, email string, reason string) error {
	f.suppressed[email] = reason
	return nil
}
func (f *fakeStore) LogSend(ctx context.Context, entry EmailLogEntry) error { return nil }
func (f *fakeStore) CountSoftFailures(ctx context.Context, email string) (int, error) {
	return f.softCount, f.softErr
}

func TestClassifySESError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"message rejected is hard", awserr.New(ses.ErrCodeMessageRejected, "rejected", nil), models.EmailErrorTypeHard},
		{"unverified from domain is system", awserr.New(ses.ErrCodeMailFromDomainNotVerifiedException, "unverified", nil), models.EmailErrorTypeSystem},
		{"missing config set is system", awserr.New(ses.ErrCodeConfigurationSetDoesNotExistException, "missing", nil), models.EmailErrorTypeSystem},
		{"unrecognized aws code is soft", awserr.New(ses.ErrCodeLimitExceededException, "throttled", nil), models.EmailErrorTypeSoft},
		{"generic error is soft", errors.New("dial tcp: timeout"), models.EmailErrorTypeSoft},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifySESError(tc.err); got != tc.want {
				t.Fatalf("classifySESError = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHandleSendFailure_HardSuppressesImmediately(t *testing.T) {
	fake := newFakeStore()
	SetSuppressionStore(fake)
	defer SetSuppressionStore(nil)

	handleSendFailure(context.Background(), "bounce@example.com", models.EmailErrorTypeHard)

	if fake.suppressed["bounce@example.com"] != models.SuppressionReasonHardBounce {
		t.Fatalf("expected hard bounce to suppress immediately, got %q", fake.suppressed["bounce@example.com"])
	}
}

func TestHandleSendFailure_SystemNeverSuppresses(t *testing.T) {
	fake := newFakeStore()
	fake.softCount = 100 // even with a high count, system errors must not suppress
	SetSuppressionStore(fake)
	defer SetSuppressionStore(nil)

	handleSendFailure(context.Background(), "ops@example.com", models.EmailErrorTypeSystem)

	if _, ok := fake.suppressed["ops@example.com"]; ok {
		t.Fatalf("system error must not suppress the address")
	}
}

func TestHandleSendFailure_SoftSuppressesAtThreshold(t *testing.T) {
	// Below the threshold: no suppression.
	below := newFakeStore()
	below.softCount = FailureThreshold - 1
	SetSuppressionStore(below)
	handleSendFailure(context.Background(), "soft@example.com", models.EmailErrorTypeSoft)
	if _, ok := below.suppressed["soft@example.com"]; ok {
		t.Fatalf("soft failure below threshold must not suppress")
	}
	SetSuppressionStore(nil)

	// At the threshold: suppress.
	at := newFakeStore()
	at.softCount = FailureThreshold
	SetSuppressionStore(at)
	defer SetSuppressionStore(nil)
	handleSendFailure(context.Background(), "soft@example.com", models.EmailErrorTypeSoft)
	if at.suppressed["soft@example.com"] != models.SuppressionReasonThresholdExceeded {
		t.Fatalf("soft failure at threshold must suppress with threshold_exceeded, got %q", at.suppressed["soft@example.com"])
	}
}
