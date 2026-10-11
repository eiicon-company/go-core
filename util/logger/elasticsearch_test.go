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

// afterFuncContext records the callbacks context.AfterFunc registers on it and which were stopped.
// Done must be non-nil, or AfterFunc registers nothing.
type afterFuncContext struct {
	context.Context
	done       chan struct{}
	registered int
	stopped    int
}

func (c *afterFuncContext) Done() <-chan struct{} { return c.done }

func (c *afterFuncContext) AfterFunc(func()) func() bool {
	c.registered++
	return func() bool { c.stopped++; return true }
}

func TestSentryErrorLoggerStopsCancelWatchWhenDone(t *testing.T) {
	failure := errors.New("connection refused")
	tests := []struct {
		name    string
		last    error
		reports int
	}{
		{name: "succeeded on retry", last: nil, reports: 0},
		{name: "reported after the last attempt", last: failure, reports: 1},
		{name: "ended by its context", last: context.Canceled, reports: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: a long-lived context, as a server or batch job passes.
			ctx := &afterFuncContext{Context: context.Background(), done: make(chan struct{})}
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://es.invalid/", nil)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			l := NewSentryErrorLogger(2)
			reports := 0
			l.Report = func(string) { reports++ }
			_ = l.LogRoundTrip(req, nil, failure, time.Now(), 0)

			// Act
			_ = l.LogRoundTrip(req, nil, tt.last, time.Now(), 0)

			// Assert: the watch on the context must not outlive the request.
			if ctx.registered != 1 || ctx.stopped != 1 {
				t.Errorf("Miss match value: want 1 watch registered and stopped, got registered=%d stopped=%d",
					ctx.registered, ctx.stopped)
			}
			if _, ok := l.failures.Load(req); ok {
				t.Errorf("Miss match value: want the failure count dropped, still present")
			}
			if reports != tt.reports {
				t.Errorf("Miss match value: want %d reports, got %d", tt.reports, reports)
			}
		})
	}
}
