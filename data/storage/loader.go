package storage

import (
	"context"
	"net/url"

	"github.com/facebookgo/inject"

	"github.com/eiicon-company/go-utils/util"
	"github.com/eiicon-company/go-utils/util/dsn"
	"github.com/eiicon-company/go-utils/util/logger"
)

// Inject injects dependencies
func Inject(ctx context.Context, env util.Environment, g *inject.Graph, rt interface{}) {
	// inject
	var objects []*inject.Object

	fu, _ := url.Parse(env.EnvString("FURI"))
	switch fu.Scheme {
	default:
		file, err := dsn.File(env.EnvString("FURI"))
		if err != nil {
			msg := "[PANIC] failed to parse file uri <%s>: %s"
			logger.Panicf(msg, env.EnvString("FURI"), err)
		}

		msg := "[INFO] a storage folder is chosen filesystems by <%s>"
		logger.Printf(msg, env.EnvString("FURI"))

		objects = []*inject.Object{
			{Value: &fileStorage{dsn: file}},
		}

	case "s3":
		s3, err := dsn.S3(env.EnvString("FURI"))
		if err != nil {
			msg := "[PANIC] failed to parse s3 uri <%s>: %s"
			logger.Panicf(msg, env.EnvString("FURI"), err)
		}

		msg := "[INFO] a storage folder is chosen s3 by <%s>"
		logger.Printf(msg, env.EnvString("FURI"))

		objects = []*inject.Object{
			{Value: &s3Storage{dsn: s3}},
		}
		// case "gcs": TODO: gs://<bucket_name>/<file_path_inside_bucket>.
		//
		//
		//
		//
	}

	if err := g.Provide(objects...); err != nil {
		logger.Panicf("[PANIC] failed to process injection: %s", err)
	}
}
