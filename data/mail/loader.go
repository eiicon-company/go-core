package mail

import (
	"go.uber.org/dig"

	"github.com/eiicon-company/go-core/util/logger"
)

// Inject injects dependencies
func Inject(di *dig.Container) {
	// Injects
	var deps = []interface{}{
		newMail,
	}

	for _, dep := range deps {
		if err := di.Provide(dep); err != nil {
			logger.Panicf("failed to process go-core mail injection: %s", err)
		}
	}
}
