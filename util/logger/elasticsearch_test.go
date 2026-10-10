package logger

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestSentryErrorLoggerForgetsCanceledRequests(t *testing.T) {
	// Arrange: one failed attempt, then the caller gives up during the retry backoff.
	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://es.invalid/", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	l := NewSentryErrorLogger(3)
	l.Report = func(string) { t.Errorf("Miss match value: want no report for a canceled request") }
	_ = l.LogRoundTrip(req, nil, errors.New("connection refused"), time.Now(), 0)

	// Act
	cancel()

	// Assert: the transport never calls the logger again, so the count must go on its own.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, ok := l.failures.Load(req); !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("Miss match value: want the failure count dropped after cancel, still present")
}
