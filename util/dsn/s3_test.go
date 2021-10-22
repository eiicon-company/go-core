package dsn

import (
	"fmt"
	"os"
	"testing"
)

func TestS3(t *testing.T) {
	t.Helper()

	f, err := S3("redis://127.0.0.1:6379/4")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = S3("s3://bucket/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	t.Logf("S3: %#+v", f)
}

func TestS3String(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/path/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	if f.String("filename.jpg") != "s3://data-bucket/path/filename.jpg" {
		t.Fatalf("Miss match value: %v", f.String("filename.jpg"))
	}

	t.Logf("S3.String: %s", f.String("filename.jpg"))
}

func TestS3URL(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/path/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	name := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", "data-bucket", os.Getenv("AWS_REGION"), "path", "filename.jpg")
	if name != f.URL("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.URL("filename.jpg"))
	}

	t.Logf("S3.URL: %s", f.URL("filename.jpg"))
}

func TestS3PublicURL(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/data.flac?url=https://example.com")
	if err != nil {
		t.Fatalf("S3.URL: %v", f)
	}

	if fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg") != f.URL("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.URL("filename.jpg"))
	}

	t.Logf("S3.URL: %s", f.URL("filename.jpg"))
}
