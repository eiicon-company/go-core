package dsn

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
)

type (
	// FileDSN ...
	FileDSN struct {
		Folder string
	}
)

// Join returns ...
func (dsn *FileDSN) Join(filename string) string {
	return filepath.Join(dsn.Folder, filename)
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
