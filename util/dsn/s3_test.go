package dsn

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestS3(t *testing.T) {
	t.Helper()

	f, err := S3("redis://127.0.0.1:6379/4")
	require.Error(t, err)

	f, err = S3("s3://bucket/data.flac")
	require.NoError(t, err)

	t.Logf("S3: %#+v", f)
}

func TestS3String(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/path/data.flac")
	require.NoError(t, err)

	require.Equal(t, "s3://data-bucket/path/filename.jpg", f.String("filename.jpg"))

	t.Logf("S3.String: %s", f.String("filename.jpg"))
}

func TestS3URL(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/path/data.flac")
	require.NoError(t, err)

	name := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", "data-bucket", os.Getenv("AWS_REGION"), "path", "filename.jpg")
	require.Equal(t, name, f.URL("filename.jpg"))

	t.Logf("S3.URL: %s", f.URL("filename.jpg"))
}

func TestS3PublicURL(t *testing.T) {
	t.Helper()

	f, err := S3("s3://data-bucket/data.flac?url=https://example.com")
	require.NoError(t, err)

	expected := fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg")
	require.Equal(t, expected, f.URL("filename.jpg"))

	t.Logf("S3.URL: %s", f.URL("filename.jpg"))
}
