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
	}
)

func newStorage(env util.Environment) Storage {
	fu, _ := url.Parse(env.EnvString("FURI"))

	switch fu.Scheme {
	default:
		file, err := dsn.File(env.EnvString("FURI"))
		if err != nil {
			msg := "[PANIC] failed to parse file uri <%s>: %s"
			logger.Panicf(msg, env.EnvString("FURI"), err)
		}

		msg := "[INFO] a storage folder is chosen filesystems to <%s>"
		logger.Printf(msg, file.Folder)

		return &fileStorage{dsn: file}

	case "s3":
		s3, err := dsn.S3(env.EnvString("FURI"))
		if err != nil {
			msg := "[PANIC] failed to parse s3 uri <%s>: %s"
			logger.Panicf(msg, env.EnvString("FURI"), err)
		}

		msg := "[INFO] a storage folder is chosen s3 by <%s>"
		logger.Printf(msg, env.EnvString("FURI"))

		return &s3Storage{dsn: s3}

		// case "gcs": TODO: gs://<bucket_name>/<file_path_inside_bucket>.
		//
		//
		//
		//
	}
}
