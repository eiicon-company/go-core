package dsn

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGCS(t *testing.T) {
	t.Helper()

	f, err := GCS("redis://127.0.0.1:6379/4")
	require.Error(t, err)

	f, err = GCS("gs://bucket/path/data.flac")
	if err != nil {
		require.Contains(t, err.Error(), "could not find default credentials")
	}

	t.Logf("GCS: %#+v", f)
}

func TestGCSString(t *testing.T) {
	t.Helper()

	f := GCSDSN{
		Bucket: "data-bucket",
		Key:    "/path/data.flac",
	}

	require.Equal(t, "gs://data-bucket/path/filename.jpg", f.String("filename.jpg"))

	t.Logf("GCS.String: %s", f.String("filename.jpg"))
}

func TestGCSPublicURL(t *testing.T) {
	t.Helper()

	f := GCSDSN{
		Bucket: "data-bucket",
		Key:    "/path/data.flac",
	}

	f.PublicURL, _ = url.Parse("https://example.com")

	expected := fmt.Sprintf("%s/%s", "https://example.com", "filename.jpg")
	require.Equal(t, expected, f.URL("filename.jpg"))

	t.Logf("GCS.URL: %s", f.URL("filename.jpg"))
}
