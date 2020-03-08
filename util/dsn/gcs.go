package dsn

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

type (
	// GCSDSN gs://data-bucket/path/
	// 			  gs://data-bucket/path/?url=https://exampl.ecom:80
	GCSDSN struct {
		Sess   *google.Credentials
		Bucket string
		Key    string

		PublicURL *url.URL
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

// URL returns https URL
//
// TODO: Get no auth or authed or private or public URL
//
func (dsn *GCSDSN) URL(filename string) string {
	if dsn.PublicURL != nil {
		u, _ := url.Parse(dsn.PublicURL.String())
		u.Path = path.Join(filepath.Dir(u.Path), filename)
		return u.String()
	}

	conf, err := google.JWTConfigFromJSON(dsn.Sess.JSON, storage.ScopeReadOnly)
	if err != nil {
		return ""
	}

	// https://cloud.google.com/storage/docs/access-control/signed-urls
	opts := &storage.SignedURLOptions{
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Method:         http.MethodGet,
		Expires:        time.Now().Add(3 * time.Hour),
	}

	uri, err := storage.SignedURL(dsn.Bucket, dsn.Key, opts)
	if err != nil {
		return ""
	}

	u, _ := url.Parse(uri) // TODO: Auth URL: Public or Private URL
	u.Path = path.Join(filepath.Dir(u.Path), filename)
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

// GCS ...
func GCS(uri string) (*GCSDSN, error) {
	if uri == "" {
		return nil, ef("invalid gcs dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, xerrors.Errorf("invalid gcs dsn: %w", err)
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

	sess, err := gcpSession()
	if err != nil {
		msg := "invalid s3 environment variables: %w"
		return nil, xerrors.Errorf(msg, err)
	}

	pubURL, err := url.Parse(u.Query().Get("url"))
	if err != nil {
		return nil, xerrors.Errorf("invalid url='' queryString: %w", err)
	}

	dsn := &GCSDSN{
		Sess:   sess,
		Bucket: u.Host,
		Key:    u.Path,
	}

	if pubURL.Scheme != "" && pubURL.Host != "" {
		dsn.PublicURL = pubURL
	}

	return dsn, nil
}

func gcpSession() (*google.Credentials, error) {
	ctx := context.Background()

	// https://github.com/golang/oauth2/blob/master/google/default.go#L61
	//
	creds, err := google.FindDefaultCredentials(ctx, storage.ScopeFullControl)
	if err != nil {
		return nil, xerrors.Errorf("gcp session failed creation: %w", err)
	}

	return creds, nil
}
