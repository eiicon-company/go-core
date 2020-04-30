package dsn

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestGCS(t *testing.T) {
	t.Helper()

	f, err := GCS("redis://127.0.0.1:6379/4")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = GCS("gs://bucket/path/data.flac")
	if err != nil && !strings.Contains(err.Error(), "could not find default credentials") {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	t.Logf("GCS: %#+v", f)
}

func TestGCSString(t *testing.T) {
	t.Helper()

	f := GCSDSN{
		Bucket: "data-bucket",
		Key:    "/path/data.flac",
	}

	if "gs://data-bucket/path/filename.jpg" != f.String("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.String("filename.jpg"))
	}

	t.Logf("GCS.String: %s", f.String("filename.jpg"))
}

func TestGCSPublicURL(t *testing.T) {
	t.Helper()

	f := GCSDSN{
		Bucket: "data-bucket",
		Key:    "/path/data.flac",
	}

	f.PublicURL, _ = url.Parse("https://example.com")

	if fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg") != f.URL("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.URL("filename.jpg"))
	}

	t.Logf("GCS.URL: %s", f.URL("filename.jpg"))
}
