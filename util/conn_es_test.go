package util

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
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
