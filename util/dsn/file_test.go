package dsn

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	t.Helper()

	f, err := File("redis://127.0.0.1:6379/4")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = File("file://./storage/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	t.Logf("File: %#+v", f)
}

func TestFileDotORSlash(t *testing.T) {
	t.Helper()

	f, err := File("file://storage/data.flac")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = File("file://.storage/data.flac")
	if err == nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = File("file://./storage/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	f, err = File("file:///storage/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	t.Logf("File: %#+v", f)
}

func TestFileString(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	abs, _ := filepath.Abs("./")
	if fmt.Sprintf("file://%s/%s/%s", abs, "storage", "filename.jpg") != f.String("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.String("filename.jpg"))
	}

	t.Logf("File.String: %s", f.String("filename.jpg"))
}

func TestFileURL(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac")
	if err != nil {
		t.Fatalf("Unknown Scheme: file=%#+v err=%v", f, err)
	}

	if fmt.Sprintf("%s/%s", filePublicURL, "filename.jpg") != f.URL("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.URL("filename.jpg"))
	}

	t.Logf("File.URL: %s", f.URL("filename.jpg"))
}

func TestFilePublicURL(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac?url=https://example.com")
	if err != nil {
		t.Fatalf("File.URL: %v", f)
	}

	if fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg") != f.URL("filename.jpg") {
		t.Fatalf("Miss match value: %v", f.URL("filename.jpg"))
	}

	t.Logf("File.URL: %s", f.URL("filename.jpg"))
}
