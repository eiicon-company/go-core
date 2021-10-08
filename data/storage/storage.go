// Package storage will be used here https://github.com/google/go-cloud someday
package storage

import (
	"context"
	"net/url"
	"regexp"
	"time"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
)

var (
	gzipPtn = regexp.MustCompile(".gz$") // gzipPtn uses gzip file determination.
)

type (
	// Storage provides interface for writes some of kinda data.
	Storage interface {
		Write(ctx context.Context, filename string, data []byte) error
		Read(ctx context.Context, filename string) ([]byte, error)
		Delete(ctx context.Context, filename string) error
		Merge(ctx context.Context, filename string, data []byte) error
		Files(ctx context.Context, ptn string) ([]string, error)
		URL(ctx context.Context, filename string) string
		String(ctx context.Context, filename string) string
		PresignedUploadURL(ctx context.Context, filename string, expire time.Duration) (string, error)
		PresignedDownloadURL(ctx context.Context, filename string, expire time.Duration) (string, error)
	}
)

func newStorage(env util.Environment) Storage {
	return SelectStorage(env.EnvString("FURI"))
}

// SelectStorage can choose storage connection
func SelectStorage(fURI string) Storage {
	fu, _ := url.Parse(fURI)
	switch fu.Scheme {
	default: // file://<bucket_name>/<file_path_inside_bucket>.
		file, err := dsn.File(fURI)
		if err != nil {
			logger.Panicf("failed to parse file uri <%s>: %s", fURI, err)
		}

		msg := "A storage folder is chosen filesystems to <%s> Public URL: <%s>"
		logger.Infof(msg, file.Folder, file.PublicURL)

		return &fileStorage{dsn: file}

	case "s3": // s3://<bucket_name>/<file_path_inside_bucket>.
		s3, err := dsn.S3(fURI)
		if err != nil {
			logger.Panicf("failed to parse s3 uri <%s>: %s", fURI, err)
		}

		msg := "a storage folder is chosen s3 by <%s> Public URL: <%s>"
		logger.Infof(msg, fURI, s3.PublicURL)

		return &s3Storage{dsn: s3}

	case "gs": // gs://<bucket_name>/<file_path_inside_bucket>.
		gcs, err := dsn.GCS(fURI)
		if err != nil {
			logger.Panicf("failed to parse gcs uri <%s>: %s", fURI, err)
		}

		msg := "a storage folder is chosen gcs by <%s> Public URL: <%s>"
		logger.Infof(msg, fURI, gcs.PublicURL)

		return &gcsStorage{dsn: gcs}
	}
}
