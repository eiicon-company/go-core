package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryErrorLogger satisfies elastictransport.Logger and reports failed round trips to Sentry.
type SentryErrorLogger struct{}

// LogRoundTrip reports transport failures only; successful and rejected requests stay silent.
func (a *SentryErrorLogger) LogRoundTrip(req *http.Request, _ *http.Response, err error, _ time.Time, _ time.Duration) error {
	if err == nil {
		return nil
	}

	url := ""
	if req != nil && req.URL != nil {
		url = req.URL.Host
	}
	errdeps(sentry.CaptureMessage, 4, fmt.Sprintf("elastic: %s is dead: %s", url, err))
	return nil
}

// RequestBodyEnabled reports that request bodies are not needed.
func (a *SentryErrorLogger) RequestBodyEnabled() bool { return false }

// ResponseBodyEnabled reports that response bodies are not needed.
func (a *SentryErrorLogger) ResponseBodyEnabled() bool { return false }
