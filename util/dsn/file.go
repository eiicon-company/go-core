package dsn

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

type (
	// FileDSN file://./storage/data.flac
	FileDSN struct {
		Folder string
		Route  string
	}
)

// FileRoute Http URL
var FileRoute = "http://localhost:8000"

// Join returns ...
func (dsn *FileDSN) Join(filename string) string {
	return filepath.Join(dsn.Folder, filename)
}

func (dsn *FileDSN) String(filename string) string {
	return fmt.Sprintf("file://%s", dsn.Join(filename))
}

// URL returns https URL
func (dsn *FileDSN) URL(filename string) string {
	route := dsn.Route
	if route == "" {
		route = FileRoute
	}
	return filepath.Join(route, filename)
}

// File ...
func File(uri string) (*FileDSN, error) {
	if uri == "" {
		return nil, ef("invalid file dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "invalid file dsn")
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
		return nil, errors.Wrap(err, "invalid file dsn")
	}

	return &FileDSN{Folder: filepath.Dir(abs)}, nil
}
