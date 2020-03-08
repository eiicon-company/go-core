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
		Delete(filename string) error
		Merge(filename string, data []byte) error
		Files(ptn string) ([]string, error)
		URL(filename string) string
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

		// case "gcs": // gs://<bucket_name>/<file_path_inside_bucket>.
		// 	s3, err := dsn.S3(fURI)
		// 	if err != nil {
		// 		msg := "failed to parse s3 uri <%s>: %s"
		// 		logger.Panicf(msg, fURI, err)
		// 	}
		//
		// 	msg := "a storage folder is chosen s3 by <%s> Public URL: <%s>"
		// 	logger.Infof(msg, fURI, s3.PublicURL)
		//
		// 	return &s3Storage{dsn: s3}
	}
}
