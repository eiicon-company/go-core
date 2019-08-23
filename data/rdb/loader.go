package rdb

import (
	"context"

	"github.com/facebookgo/inject"

	"github.com/eiicon-company/go-utils/util"
	"github.com/eiicon-company/go-utils/util/logger"
)

// Inject injects dependencies
func Inject(ctx context.Context, env util.Environment, g *inject.Graph, rt interface{}) {
	// inject
	err := g.Provide(
		&inject.Object{Value: &rdb{}},
	)
	if err != nil {
		logger.Panicf("[PANIC] Failed to process injection: %s", err)
	}
}
