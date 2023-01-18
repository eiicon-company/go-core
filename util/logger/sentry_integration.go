package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryDevNullTransport output to like dev null
type SentryDevNullTransport struct{}

// Configure ...
func (t *SentryDevNullTransport) Configure(options sentry.ClientOptions) {
	dsn, _ := sentry.NewDsn(options.Dsn)
	fmt.Println("[FakeSentry] Stores Endpoint:", dsn.StoreAPIURL())
	fmt.Println("[FakeSentry] Headers:", dsn.RequestHeaders())
}

// SendEvent ...
func (t *SentryDevNullTransport) SendEvent(event *sentry.Event) {
	b, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("[FakeSentry] log failed: %+v", err)
		return
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		fmt.Printf("[FakeSentry] log failed: %+v", err)
		return
	}

	fmt.Println("[FakeSentry] SentEvent", out.String())
}

// Flush ...
func (t *SentryDevNullTransport) Flush(timeout time.Duration) bool {
	return true
}

// SentryLoggerIntegration filters no need stacktrace frames
type SentryLoggerIntegration struct{}

// Name ...
func (it *SentryLoggerIntegration) Name() string {
	return "SentryLogger"
}

// SetupOnce ...
func (it *SentryLoggerIntegration) SetupOnce(client *sentry.Client) {
	client.AddEventProcessor(it.processor)
}

func (it *SentryLoggerIntegration) processor(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	for _, thread := range event.Threads {
		if thread.Stacktrace == nil {
			continue
		}

		it.filterFrames(thread.Stacktrace)
	}

	for _, exc := range event.Exception {
		if exc.Stacktrace == nil {
			continue
		}

		it.filterFrames(exc.Stacktrace)
	}

	return event
}

func (it *SentryLoggerIntegration) filterFrames(trace *sentry.Stacktrace) {
	frames := trace.Frames[:0]
	for _, frame := range trace.Frames {
		if frame.Module == "github.com/eiicon-company/go-core/util/logger" {
			continue
		}

		frames = append(frames, frame)
	}

	trace.Frames = frames
}
