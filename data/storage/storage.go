package storage

import (
	"context"
	"net/url"
	"regexp"

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
	}
)

func newStorage(env util.Environment) Storage {
	fURI := env.EnvString("FURI")

	fu, _ := url.Parse(fURI)
	switch fu.Scheme {
	default:
		file, err := dsn.File(fURI)
		if err != nil {
			msg := "failed to parse file uri <%s>: %s"
			logger.Panicf(msg, fURI, err)
		}

		msg := "A storage folder is chosen filesystems to <%s> Public URL: <%s>"
		logger.Infof(msg, file.Folder, file.PublicURL)

		return &fileStorage{dsn: file}

	case "s3":
		s3, err := dsn.S3(fURI)
		if err != nil {
			msg := "failed to parse s3 uri <%s>: %s"
			logger.Panicf(msg, fURI, err)
		}

		msg := "a storage folder is chosen s3 by <%s> Public URL: <%s>"
		logger.Infof(msg, fURI, s3.PublicURL)

		return &s3Storage{dsn: s3}

	case "gs": // gs://<bucket_name>/<file_path_inside_bucket>.
		gcs, err := dsn.GCS(fURI)
		if err != nil {
			msg := "failed to parse gcs uri <%s>: %s"
			logger.Panicf(msg, fURI, err)
		}

		msg := "a storage folder is chosen gcs by <%s> Public URL: <%s>"
		logger.Infof(msg, fURI, gcs.PublicURL)

		return &gcsStorage{dsn: gcs}
	}
}
