package util

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
)

// captureTransport records what sentry would have sent over the wire.
type captureTransport struct {
	events []*sentry.Event
}

func (t *captureTransport) Configure(sentry.ClientOptions)        {}
func (t *captureTransport) SendEvent(event *sentry.Event)         { t.events = append(t.events, event) }
func (t *captureTransport) Flush(time.Duration) bool              { return true }
func (t *captureTransport) FlushWithContext(context.Context) bool { return true }
func (t *captureTransport) Close()                                {}

// emittedSpans runs fn inside a fully sampled transaction and returns the spans
// it produced, as sentry would have received them.
func emittedSpans(t *testing.T, fn func(ctx context.Context)) []*sentry.Span {
	t.Helper()

	transport := &captureTransport{}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://user@example.com/1",
		Transport:        transport,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		t.Fatalf("Unable to init sentry: %s", err)
	}

	tx := sentry.StartTransaction(context.Background(), "conn_test")
	fn(tx.Context())
	tx.Finish()

	if len(transport.events) != 1 {
		t.Fatalf("Miss match value: want 1 event, got %d", len(transport.events))
	}

	return transport.events[0].Spans
}

// TestEmitSlowSpanDatesSpanByMeasuredElapsed guards the span duration reported to
// sentry: it must equal the measured elapsed time, not twice it.
func TestEmitSlowSpanDatesSpanByMeasuredElapsed(t *testing.T) {
	// Arrange
	const since = 120 * time.Millisecond
	startTime := time.Now().Add(-since)

	// Act
	spans := emittedSpans(t, func(ctx context.Context) {
		emitSlowSpan(ctx, "db.sql.query.slow", startTime, since, 50*time.Millisecond, "SELECT 1", nil)
	})

	// Assert
	if len(spans) != 1 {
		t.Fatalf("Miss match value: want 1 span, got %d", len(spans))
	}

	span := spans[0]
	if got := span.EndTime.Sub(span.StartTime); got != since {
		t.Errorf("Miss match value: want %s, got %s", since, got)
	}
	if !span.StartTime.Equal(startTime) {
		t.Errorf("Miss match value: want %s, got %s", startTime, span.StartTime)
	}
	if span.EndTime.After(time.Now()) {
		t.Errorf("Miss match value: end time %s is in the future", span.EndTime)
	}
}

func TestEmitSlowSpanDescribesStatement(t *testing.T) {
	// Arrange
	const query = "SELECT * FROM users WHERE id = ?"

	// Act
	spans := emittedSpans(t, func(ctx context.Context) {
		emitSlowSpan(ctx, "db.sql.exec.slow", time.Now().Add(-time.Second), time.Second, 50*time.Millisecond, query, nil)
	})

	// Assert
	if len(spans) != 1 {
		t.Fatalf("Miss match value: want 1 span, got %d", len(spans))
	}
	if spans[0].Op != "db.sql.exec.slow" {
		t.Errorf("Miss match value: want db.sql.exec.slow, got %s", spans[0].Op)
	}
	if spans[0].Description != query {
		t.Errorf("Miss match value: want %s, got %s", query, spans[0].Description)
	}
}

func TestEmitSlowSpanPeriodGate(t *testing.T) {
	period := 50 * time.Millisecond

	tests := []struct {
		name  string
		since time.Duration
		want  int
	}{
		{name: "faster than period is not reported", since: 10 * time.Millisecond, want: 0},
		{name: "exactly period is not reported", since: period, want: 0},
		{name: "slower than period is reported", since: period + time.Millisecond, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			spans := emittedSpans(t, func(ctx context.Context) {
				emitSlowSpan(ctx, "db.sql.query.slow", time.Now().Add(-tt.since), tt.since, period, "SELECT 1", nil)
			})

			// Assert
			if len(spans) != tt.want {
				t.Errorf("Miss match value: want %d spans, got %d", tt.want, len(spans))
			}
		})
	}
}

func TestEmitSlowSpanBoundArgs(t *testing.T) {
	// namedArgs builds args that are keyed by Name; ordinalArgs fall back to Ordinal.
	namedArgs := []driver.NamedValue{
		{Name: "user_id", Ordinal: 1, Value: int64(7)},
		{Ordinal: 2, Value: "abc"},
	}

	// The loop keeps indexes 0..50 and breaks at 51, so 60 args collapse to 51.
	manyArgs := make([]driver.NamedValue, 0, 60)
	for i := 1; i <= 60; i++ {
		manyArgs = append(manyArgs, driver.NamedValue{Ordinal: i, Value: i})
	}

	tests := []struct {
		name string
		args []driver.NamedValue
		want map[string]interface{}
		size int
	}{
		{
			name: "named args key by name, unnamed fall back to ordinal",
			args: namedArgs,
			want: map[string]interface{}{"user_id": "7", "2": "abc"},
			size: 2,
		},
		{
			name: "no args yields empty data",
			args: nil,
			size: 0,
		},
		{
			name: "args beyond the cap are dropped",
			args: manyArgs,
			size: 51,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			spans := emittedSpans(t, func(ctx context.Context) {
				emitSlowSpan(ctx, "db.sql.query.slow", time.Now().Add(-time.Second), time.Second, 50*time.Millisecond, "SELECT 1", tt.args)
			})

			// Assert
			if len(spans) != 1 {
				t.Fatalf("Miss match value: want 1 span, got %d", len(spans))
			}

			data := spans[0].Data
			if len(data) != tt.size {
				t.Errorf("Miss match value: want %d args, got %d", tt.size, len(data))
			}
			for k, want := range tt.want {
				if got := data[k]; got != want {
					t.Errorf("Miss match value: want %s=%v, got %v", k, want, got)
				}
			}
		})
	}
}
