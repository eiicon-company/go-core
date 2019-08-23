package dsn

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

type (
	// GCSDSN ...
	GCSDSN struct {
		Bucket string
		Key    string
	}
)

// Join returns file joined string that discards
// key's basename and then combine filename.
func (dsn *GCSDSN) Join(filename string) string {
	return filepath.Join(filepath.Dir(dsn.Key), filename)
}

func (dsn *GCSDSN) String(filename string) string {
	return fmt.Sprintf("gs://%s%s", dsn.Bucket, dsn.Join(filename))
}

// GCS ...
func GCS(uri string) (*GCSDSN, error) {
	if uri == "" {
		return nil, ef("invalid gcs dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "invalid gcs dsn")
	}
	if u.Scheme != "gs" {
		return nil, ef("invalid gs scheme: %s", u.Scheme)
	}
	if u.Host == "" {
		return nil, ef("invalid gcs bucket is blank")
	}
	if u.Path == "" {
		return nil, ef("invalid gcs key is blank")
	}

	dsn := &GCSDSN{
		Bucket: u.Host,
		Key:    u.Path,
	}
	return dsn, nil
}
