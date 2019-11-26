package dsn

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"golang.org/x/xerrors"
)

type (
	// FileDSN file://./storage/data.flac
	FileDSN struct {
		Folder    string
		PublicURL *url.URL
	}
)

// filePublicURL Http URL
var filePublicURL = "http://localhost:8000"

// Join returns ...
func (dsn *FileDSN) Join(filename string) string {
	return filepath.Join(dsn.Folder, filename)
}

func (dsn *FileDSN) String(filename string) string {
	return fmt.Sprintf("file://%s", dsn.Join(filename))
}

// URL returns https URL
func (dsn *FileDSN) URL(filename string) string {
	u, _ := url.Parse(filePublicURL)
	if dsn.PublicURL != nil {
		u, _ = url.Parse(dsn.PublicURL.String())
	}

	u.Path = path.Join(u.Path, filename)
	return u.String()
}

// File ...
func File(uri string) (*FileDSN, error) {
	if uri == "" {
		return nil, ef("invalid file dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, xerrors.Errorf("invalid file dsn: %w", err)
	}
	if u.Scheme != "file" {
		return nil, ef("invalid file scheme: %s", u.Scheme)
	}

	if u.Host != "" && u.Host != "." && u.Host != ".." {
		msg := "invalid file path prefix. that must be dotslach(./) or slash(/) char set."
		return nil, ef(msg)
	}
	if u.Path == "" {
		return nil, ef("invalid file path is blank")
	}

	folder := fmt.Sprintf("%s%s", u.Host, u.Path)
	abs, err := filepath.Abs(folder)
	if err != nil {
		return nil, xerrors.Errorf("invalid file dsn: %w", err)
	}

	pubURL, err := url.Parse(u.Query().Get("url"))
	if err != nil {
		return nil, xerrors.Errorf("invalid url='' queryString: %w", err)
	}

	if pubURL.Scheme == "" || pubURL.Host == "" {
		return &FileDSN{Folder: filepath.Dir(abs)}, nil
	}

	return &FileDSN{Folder: filepath.Dir(abs), PublicURL: pubURL}, nil
}
