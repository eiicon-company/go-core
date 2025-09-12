package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type (
	// SentryErrorLogger is satisfied of elastic.Logger
	SentryErrorLogger struct{}
	// SentryInfoLogger is satisfied of elastic.Logger
	SentryInfoLogger struct{}
)

// Printf prints out message as error
func (a *SentryErrorLogger) Printf(format string, v ...any) {
	errdeps(sentry.CaptureMessage, 4, fmt.Sprintf(format, v...))
}


// Printf prints out message as info
func (a *SentryInfoLogger) Printf(format string, v ...any) {
	infodeps(sentry.CaptureMessage, 4, fmt.Sprintf(format, v...))
}

func (a *SentryErrorLogger) LogRoundTrip( req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	var l func(format string, args ...any)
	
	if err != nil || res != nil && res.StatusCode >= 500 {
		l = Errorf
	}

	if l != nil {
		l("%s, method: %s, status_code: %d", req.URL.String(), req.Method, res.StatusCode)
	}

	return nil
}

func (l *SentryErrorLogger) RequestBodyEnabled() bool { return true }
func (l *SentryErrorLogger) ResponseBodyEnabled() bool { return true }
