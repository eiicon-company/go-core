package storage

import (
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
		Write(filename string, data []byte) error
		Read(filename string) ([]byte, error)
		Merge(filename string, data []byte) error
		Files(ptn string) ([]string, error)
		URL(filename string) string
	}
)

func newStorage(env util.Environment) Storage {
	fURI, fURL := env.EnvString("FURI"), env.EnvString("FURL")

	fu, _ := url.Parse(fURI)
	switch fu.Scheme {
	default:
		file, err := dsn.File(fURI)
		if err != nil {
			msg := "failed to parse file uri <%s>: %s"
			logger.Panicf(msg, fURI, err)
		}

		file.PublicURL = fURL

		msg := "A storage folder is chosen filesystems to <%s> Public URL: <%v>"
		logger.Infof(msg, file.Folder, fURL)

		return &fileStorage{dsn: file}

	case "s3":
		s3, err := dsn.S3(fURI)
		if err != nil {
			msg := "failed to parse s3 uri <%s>: %s"
			logger.Panicf(msg, fURI, err)
		}

		s3.PublicURL = fURL

		msg := "a storage folder is chosen s3 by <%s> Public URL: <%v>"
		logger.Infof(msg, fURI, fURL)

		return &s3Storage{dsn: s3}

		// case "gcs": TODO: gs://<bucket_name>/<file_path_inside_bucket>.
		//
		//
		//
		//
	}
}
