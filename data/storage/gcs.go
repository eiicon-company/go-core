// Package storage gonna be implementation
// that stream io processing for memory performance.
package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
	"github.com/gobwas/glob"
)

// gcsStorage provides implementation gcs resource interface.
type gcsStorage struct {
	Env util.Environment
	dsn *dsn.GCSDSN
}

// Write will create file into the gcs.
func (adp *gcsStorage) Write(ctx context.Context, filename string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return xerrors.Errorf("[F] gcs write client failed: %w", err)
	}

	wc := client.Bucket(adp.dsn.Bucket).
		Object(strings.TrimLeft(adp.dsn.Join(filename), "/")).
		NewWriter(ctx)

	var reader io.Reader = bytes.NewReader(data)

	if gzipPtn.MatchString(filename) {
		var writer *io.PipeWriter

		reader, writer = io.Pipe()
		go func() {
			gz := gzip.NewWriter(writer)
			if _, err := io.Copy(gz, bytes.NewReader(data)); err != nil {
				logger.ErrorfWithContext(ctx, "[F] gcs write gzip failed: %s", err)
			}

			gz.Close()
			writer.Close()
		}()
	}

	if _, err := io.Copy(wc, reader); err != nil {
		return xerrors.Errorf("[F] gcs write failed: %w", err)
	}
	if err := wc.Close(); err != nil {
		return xerrors.Errorf("[F] gcs write close failed: %w", err)
	}

	return nil
}

// Read returns file data from the gcs
func (adp *gcsStorage) Read(ctx context.Context, filename string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("[F] gcs read client failed: %w", err)
	}

	rc, err := client.Bucket(adp.dsn.Bucket).
		Object(strings.TrimLeft(adp.dsn.Join(filename), "/")).
		NewReader(ctx)

	if err != nil {
		return nil, xerrors.Errorf("[F] gcs read reader failed: %w", err)
	}
	defer rc.Close()

	var reader io.ReadCloser = rc
	defer reader.Close()

	if gzipPtn.MatchString(filename) {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, xerrors.Errorf("[F] gcs read gzip failed: %w", err)
		}
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, xerrors.Errorf("[F] gcs read failed: %w", err)
	}

	return data, nil
}

// Delete will delete file from the file systems.
func (adp *gcsStorage) Delete(ctx context.Context, filename string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return xerrors.Errorf("[F] gcs delete client failed: %w", err)
	}

	o := client.Bucket(adp.dsn.Bucket).
		Object(strings.TrimLeft(adp.dsn.Join(filename), "/"))
	if err := o.Delete(ctx); err != nil {
		return xerrors.Errorf("[F] gcs delete failed: %w", err)
	}

	return nil
}

// Merge will merge file into the gcs
func (adp *gcsStorage) Merge(ctx context.Context, filename string, data []byte) error {
	entire, _ := adp.Read(ctx, filename)
	entire = append(entire, data...)

	return adp.Write(ctx, filename, entire)
}

// Files returns filename list which is traversing with glob from gcs storage.
func (adp *gcsStorage) Files(ctx context.Context, ptn string) ([]string, error) {
	base := strings.TrimLeft(adp.dsn.Join(ptn), "/")

	g, err := glob.Compile(base)
	if err != nil {
		return nil, xerrors.Errorf("[F] gcs files pattern arg failed: %w", err)
	}
	prefix := strings.TrimSuffix(base, filepath.Base(base)) // XXX: Prefix setter is so fuzzy now.
	if !strings.Contains(prefix, "/") {
		prefix = ""
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("[F] gcs files client failed: %w", err)
	}

	files := []string{}
	it := client.Bucket(adp.dsn.Bucket).Objects(ctx, &storage.Query{Prefix: prefix})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("[F] gcs files failed: %w", err)
		}

		if g.Match(attrs.Name) {
			files = append(files, attrs.Name)
		}
	}

	return files, nil
}

// URL returns Public URL
func (adp *gcsStorage) URL(_ context.Context, filename string) string {
	return adp.dsn.URL(filename)
}

// String returns a URI
func (adp *gcsStorage) String(_ context.Context, filename string) string {
	return adp.dsn.String(filename)
}

// PresignedUploadURL returns a presigned upload URI
func (adp *gcsStorage) PresignedUploadURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", fmt.Errorf("PresignedUploadURL not supported yet")
}

// PresignedDownloadURL returns a presigned download URI
func (adp *gcsStorage) PresignedDownloadURL(_ context.Context, _ string, _ time.Duration) (string, error) {
	return "", fmt.Errorf("PresignedDownloadURL not supported yet")
}
