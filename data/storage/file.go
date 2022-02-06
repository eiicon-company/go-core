// Package storage gonna be implementation
// that stream io processing for memory performance.
//
package storage

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
	"golang.org/x/xerrors"
)

// fileStorage provides implementation file object interface.
type fileStorage struct {
	Env util.Environment
	dsn *dsn.FileDSN
}

// Write will create file into the file systems.
func (adp *fileStorage) Write(ctx context.Context, filename string, data []byte) error {
	filename = adp.dsn.Join(filename)
	folder := filepath.Dir(filename)

	fi, err := os.Stat(folder)
	if err != nil {
		_ = os.MkdirAll(folder, 0755)
	} else if !fi.IsDir() {
		return fmt.Errorf("[F] %s should be a directory", folder)
	}

	file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	fi, err = file.Stat()
	if err != nil {
		return fmt.Errorf("[F] %s file not exists", filename)
	} else if !fi.Mode().IsRegular() {
		return fmt.Errorf("[F] %s should be a file", filename)
	}

	if gzipPtn.MatchString(filename) {
		adp.gzip(file, data)
	} else {
		adp.plain(file, data)
	}

	return nil
}

// Read returns file data from the file systems.
func (adp *fileStorage) Read(ctx context.Context, filename string) ([]byte, error) {
	var reader io.ReadCloser

	reader, err := os.Open(adp.dsn.Join(filename))
	if err != nil {
		return nil, xerrors.Errorf("[F] file read failed: %w", err)
	}

	defer reader.Close()

	if gzipPtn.MatchString(filename) {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, xerrors.Errorf("[F] gzip read failed: %w", err)
		}
	}

	return ioutil.ReadAll(reader)
}

// Delete will delete file from the file systems.
func (adp *fileStorage) Delete(ctx context.Context, filename string) error {
	path := adp.dsn.Join(filename)
	return os.Remove(path)
}

// Merge will merge file into the file systems.
func (adp *fileStorage) Merge(ctx context.Context, filename string, data []byte) error {
	entire, _ := adp.Read(ctx, filename)
	entire = append(entire, data...)

	return adp.Write(ctx, filename, entire)
}

// Files returns filename list which is traversing with glob from filesystem.
func (adp *fileStorage) Files(ctx context.Context, ptn string) ([]string, error) {
	matches, err := filepath.Glob(adp.dsn.Join(ptn))
	if err != nil {
		logger.Printf("Failed to retrieve list files %s", err)
		return []string{}, err
	}

	files := []string{}
	for _, file := range matches {
		files = append(files, fmt.Sprintf("file://%s", file))
	}

	return files, nil
}

// URL returns a Public URL
func (adp *fileStorage) URL(ctx context.Context, filename string) string {
	return adp.dsn.URL(filename)
}

// String returns a URI
func (adp *fileStorage) String(ctx context.Context, filename string) string {
	return adp.dsn.String(filename)
}

// PresignedUploadURL returns a presigned upload URI
func (adp *fileStorage) PresignedUploadURL(ctx context.Context, filename string, _ time.Duration) (string, error) {
	return adp.URL(ctx, filename), nil
}

// PresignedDownloadURL returns a presigned download URI
func (adp *fileStorage) PresignedDownloadURL(ctx context.Context, filename string, _ time.Duration) (string, error) {
	return adp.URL(ctx, filename), nil
}

// gzip will create sitemap file as a gzip.
func (adp *fileStorage) gzip(file io.Writer, data []byte) {
	gz := gzip.NewWriter(file)
	defer gz.Close()
	if _, err := gz.Write(data); err != nil {
		logger.E("file gzip: %s", err)
	}

}

// plain will create uncompressed file.
func (adp *fileStorage) plain(file io.WriteCloser, data []byte) {
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		logger.E("file plain: %s", err)
	}
}
