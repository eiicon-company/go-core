package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eiicon-company/go-core/util/logger"
)

// esStubEnv is the minimal Environment for ES connection tests.
type esStubEnv struct {
	esurl string
}

func (e esStubEnv) IsProd() bool   { return false }
func (e esStubEnv) IsDev() bool    { return false }
func (e esStubEnv) IsLocal() bool  { return true }
func (e esStubEnv) IsDebug() bool  { return false }
func (e esStubEnv) IsSentry() bool { return false }
func (e esStubEnv) EnvString(prop string) string {
	if prop == "ESURL" {
		return e.esurl
	}
	return ""
}
func (e esStubEnv) EnvInt(string) int { return 0 }

// newESStub serves the root info endpoint the client probes on connect. The
// product header is what the official client uses to tell Elasticsearch apart
// from any other HTTP server.
func newESStub(t *testing.T, withProductHeader bool) (*httptest.Server, *int32) {
	t.Helper()

	var infoCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		atomic.AddInt32(&infoCalls, 1)
		if withProductHeader {
			w.Header().Set("X-Elastic-Product", "Elasticsearch")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"stub","cluster_name":"stub","cluster_uuid":"u","version":{"number":"9.4.1","build_flavor":"default","build_type":"docker","build_hash":"h","build_date":"2026-01-01T00:00:00.000Z","build_snapshot":false,"lucene_version":"10.4.0","minimum_wire_compatibility_version":"8.19.0","minimum_index_compatibility_version":"8.0.0"},"tagline":"You Know, for Search"}`))
	}))
	t.Cleanup(srv.Close)

	return srv, &infoCalls
}

func TestESConnEstablishesConnection(t *testing.T) {
	// Arrange
	srv, infoCalls := newESStub(t, true)

	// Act
	es, err := ESConn(esStubEnv{esurl: srv.URL})

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if es == nil {
		t.Fatalf("Miss match value: want a client, got nil")
	}
	if got := atomic.LoadInt32(infoCalls); got != 1 {
		t.Errorf("Miss match value: want 1 info probe, got %d", got)
	}
}

func TestESBulkConnEstablishesConnection(t *testing.T) {
	// Arrange
	srv, infoCalls := newESStub(t, true)

	// Act
	es, err := ESBulkConn(esStubEnv{esurl: srv.URL})

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if es == nil {
		t.Fatalf("Miss match value: want a client, got nil")
	}
	if got := atomic.LoadInt32(infoCalls); got != 1 {
		t.Errorf("Miss match value: want 1 info probe, got %d", got)
	}
}

func TestESConnRejectsNonElasticsearchServer(t *testing.T) {
	// Arrange
	srv, _ := newESStub(t, false)

	// Act
	es, err := ESConn(esStubEnv{esurl: srv.URL})

	// Assert
	if err == nil {
		t.Fatalf("Miss match value: want an error, got nil")
	}
	if es != nil {
		t.Errorf("Miss match value: want nil client, got %v", es)
	}
	if !strings.Contains(err.Error(), "error got es version") {
		t.Errorf("Miss match value: want version probe error, got %s", err)
	}
}

func TestESConnDoesNotRetryOnServerErrorStatus(t *testing.T) {
	// Arrange: the previous client retried transport failures only, never a response status.
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	// Act
	_, err := ESConn(esStubEnv{esurl: srv.URL})

	// Assert
	if err == nil {
		t.Fatalf("Miss match value: want an error, got nil")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("Miss match value: want 1 request, got %d", got)
	}
}

func TestRetryTransportErrorSkipsCancellationsAndTemporaryErrors(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "canceled context is returned to the caller", err: context.Canceled, want: false},
		{name: "expired deadline is returned to the caller", err: context.DeadlineExceeded, want: false},
		{name: "temporary error is returned to the caller", err: &transportError{temporary: true}, want: false},
		{name: "non-temporary timeout retries", err: &transportError{timeout: true}, want: true},
		{name: "wrapped temporary error is returned to the caller", err: fmt.Errorf("wrapped: %w", &transportError{temporary: true}), want: false},
		{name: "other transport failures retry", err: io.ErrUnexpectedEOF, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			got := retryTransportError(nil, tt.err)

			// Assert
			if got != tt.want {
				t.Errorf("Miss match value: want %t, got %t", tt.want, got)
			}
		})
	}
}

type transportError struct {
	temporary bool
	timeout   bool
}

func (e *transportError) Error() string   { return "stub transport error" }
func (e *transportError) Temporary() bool { return e.temporary }
func (e *transportError) Timeout() bool   { return e.timeout }

func TestDeadlineTransportBoundsTheBodyRead(t *testing.T) {
	// Arrange: headers arrive at once, then the body stalls until the test ends.
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-release
	}))
	t.Cleanup(func() { close(release); srv.Close() })

	rt := &deadlineTransport{base: http.DefaultTransport, timeout: 100 * time.Millisecond}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Act
	res, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)

	// Assert
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Miss match value: want a deadline error, got %v", err)
	}
}

func TestDeadlineTransportReadsNormalBody(t *testing.T) {
	// Arrange
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "complete body")
	}))
	t.Cleanup(srv.Close)
	rt := &deadlineTransport{base: http.DefaultTransport, timeout: time.Second}
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Act
	res, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	// Assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if string(body) != "complete body" {
		t.Errorf("Miss match value: want complete body, got %s", body)
	}
}

func TestCancelOnCloseReleasesDeadlineOnReadError(t *testing.T) {
	for _, readErr := range []error{io.EOF, io.ErrUnexpectedEOF} {
		t.Run(readErr.Error(), func(t *testing.T) {
			// Arrange
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			body := &cancelOnClose{ReadCloser: errorBody{err: readErr}, cancel: cancel}

			// Act
			_, err := body.Read(make([]byte, 1))

			// Assert
			if !errors.Is(err, readErr) || !errors.Is(ctx.Err(), context.Canceled) {
				t.Errorf("Miss match value: want %v and canceled context, got %v and %v", readErr, err, ctx.Err())
			}
		})
	}
}

type errorBody struct{ err error }

func (b errorBody) Read([]byte) (int, error) { return 0, b.err }
func (b errorBody) Close() error             { return nil }

type fakeESLogger struct {
	requestBody  bool
	responseBody bool
	calls        int
	err          error
}

func (l *fakeESLogger) LogRoundTrip(*http.Request, *http.Response, error, time.Time, time.Duration) error {
	l.calls++
	return l.err
}
func (l *fakeESLogger) RequestBodyEnabled() bool  { return l.requestBody }
func (l *fakeESLogger) ResponseBodyEnabled() bool { return l.responseBody }

func TestMultiLoggerCallsEveryLogger(t *testing.T) {
	// Arrange
	firstErr := errors.New("first logger failed")
	first := &fakeESLogger{err: firstErr}
	second := &fakeESLogger{requestBody: true, err: errors.New("second logger failed")}
	third := &fakeESLogger{responseBody: true}
	multi := logger.MultiLogger{first, second, third}

	// Act
	err := multi.LogRoundTrip(nil, nil, nil, time.Time{}, 0)

	// Assert
	if !errors.Is(err, firstErr) {
		t.Errorf("Miss match value: want first error, got %v", err)
	}
	for _, l := range []*fakeESLogger{first, second, third} {
		if l.calls != 1 {
			t.Errorf("Miss match value: want 1 logger call, got %d", l.calls)
		}
	}
	if !multi.RequestBodyEnabled() || !multi.ResponseBodyEnabled() {
		t.Errorf("Miss match value: want request and response bodies enabled")
	}
	if (logger.MultiLogger{first}).RequestBodyEnabled() || (logger.MultiLogger{first}).ResponseBodyEnabled() {
		t.Errorf("Miss match value: want bodies disabled when no logger requests them")
	}
}

func TestSentryErrorLoggerCountsFailedAttempts(t *testing.T) {
	// Arrange
	l := logger.NewSentryErrorLogger(3)
	var reports []string
	l.Report = func(message string) { reports = append(reports, message) }
	req, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	failure := errors.New("connection failed")

	// Act / Assert
	logAttempt := func(err error, wantReports int) {
		t.Helper()
		if logErr := l.LogRoundTrip(req, nil, err, time.Time{}, 0); logErr != nil {
			t.Fatalf("Unexpected error: %s", logErr)
		}
		if len(reports) != wantReports {
			t.Fatalf("Miss match value: want %d reports, got %d", wantReports, len(reports))
		}
	}
	logAttempt(failure, 0)
	logAttempt(failure, 0)
	logAttempt(failure, 1)
	if !strings.Contains(reports[0], "localhost is dead: connection failed") {
		t.Errorf("Miss match value: want host and error in report, got %s", reports[0])
	}
	logAttempt(failure, 1)
	logAttempt(failure, 1)
	logAttempt(nil, 1)
	logAttempt(failure, 1)
	logAttempt(failure, 1)
	logAttempt(failure, 2)
	for _, ignored := range []error{context.Canceled, context.DeadlineExceeded, &transportError{temporary: true}, fmt.Errorf("wrapped: %w", context.Canceled)} {
		logAttempt(failure, 2)
		logAttempt(failure, 2)
		logAttempt(ignored, 2)
		logAttempt(failure, 2)
		logAttempt(nil, 2)
	}
}

func TestESConnFailsWhenUnreachable(t *testing.T) {
	// Arrange: a server that is closed before the client dials it.
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	// Act
	es, err := ESConn(esStubEnv{esurl: url})

	// Assert
	if err == nil {
		t.Fatalf("Miss match value: want an error, got nil")
	}
	if es != nil {
		t.Errorf("Miss match value: want nil client, got %v", es)
	}
}
