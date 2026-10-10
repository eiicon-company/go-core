package logger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/getsentry/sentry-go"
)

// MultiLogger fans one round trip out to every logger. The transport takes a single logger, so
// without this a debug trace would displace the Sentry report the previous client always kept.
type MultiLogger []elastictransport.Logger

// LogRoundTrip passes the round trip to every logger and returns the first failure.
func (m MultiLogger) LogRoundTrip(req *http.Request, res *http.Response, err error, start time.Time, dur time.Duration) error {
	var first error
	for _, l := range m {
		if lerr := l.LogRoundTrip(req, res, err, start, dur); lerr != nil && first == nil {
			first = lerr
		}
	}
	return first
}

// RequestBodyEnabled reports whether any logger wants request bodies.
func (m MultiLogger) RequestBodyEnabled() bool {
	for _, l := range m {
		if l.RequestBodyEnabled() {
			return true
		}
	}
	return false
}

// ResponseBodyEnabled reports whether any logger wants response bodies.
func (m MultiLogger) ResponseBodyEnabled() bool {
	for _, l := range m {
		if l.ResponseBodyEnabled() {
			return true
		}
	}
	return false
}

// SentryErrorLogger reports exhausted requests rather than each failed retry.
type SentryErrorLogger struct {
	attempts int
	failures sync.Map
	Report   func(string)
}

func NewSentryErrorLogger(attempts int) *SentryErrorLogger {
	return &SentryErrorLogger{
		attempts: attempts,
		Report:   func(message string) { errdeps(sentry.CaptureMessage, 4, message) },
	}
}

// LogRoundTrip suppresses errors the previous client returned without reporting.
func (a *SentryErrorLogger) LogRoundTrip(req *http.Request, _ *http.Response, err error, _ time.Time, _ time.Duration) error {
	var temporary interface{ Temporary() bool }
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
		(errors.As(err, &temporary) && temporary.Temporary()) {
		a.failures.Delete(req)
		return nil
	}

	count, _ := a.failures.LoadOrStore(req, &atomic.Int64{})
	if count.(*atomic.Int64).Add(1) != int64(a.attempts) {
		return nil
	}
	a.failures.Delete(req)

	url := ""
	if req != nil && req.URL != nil {
		url = req.URL.Host
	}
	a.Report(fmt.Sprintf("elastic: %s is dead: %s", url, err))
	return nil
}

// RequestBodyEnabled reports that request bodies are not needed.
func (a *SentryErrorLogger) RequestBodyEnabled() bool { return false }

// ResponseBodyEnabled reports that response bodies are not needed.
func (a *SentryErrorLogger) ResponseBodyEnabled() bool { return false }
