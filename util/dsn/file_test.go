package dsn

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	t.Helper()

	f, err := File("redis://127.0.0.1:6379/4")
	require.Error(t, err)

	f, err = File("file://./storage/data.flac")
	require.NoError(t, err)

	t.Logf("File: %#+v", f)
}

func TestFileDotORSlash(t *testing.T) {
	t.Helper()

	f, err := File("file://storage/data.flac")
	require.Error(t, err)

	f, err = File("file://.storage/data.flac")
	require.Error(t, err)

	f, err = File("file://./storage/data.flac")
	require.NoError(t, err)

	f, err = File("file:///storage/data.flac")
	require.NoError(t, err)

	t.Logf("File: %#+v", f)
}

func TestFileString(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac")
	require.NoError(t, err)

	abs, _ := filepath.Abs("./")
	expected := fmt.Sprintf("file://%s/%s/%s", abs, "storage", "filename.jpg")
	require.Equal(t, expected, f.String("filename.jpg"))

	t.Logf("File.String: %s", f.String("filename.jpg"))
}

func TestFileURL(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac")
	require.NoError(t, err)

	expected := fmt.Sprintf("%s/%s", filePublicURL, "filename.jpg")
	require.Equal(t, expected, f.URL("filename.jpg"))

	t.Logf("File.URL: %s", f.URL("filename.jpg"))
}

func TestFilePublicURL(t *testing.T) {
	t.Helper()

	f, err := File("file://./storage/data.flac?url=https://example.com")
	require.NoError(t, err)

	expected := fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg")
	require.Equal(t, expected, f.URL("filename.jpg"))

	t.Logf("File.URL: %s", f.URL("filename.jpg"))
}
