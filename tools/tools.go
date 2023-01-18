//go:build tools

package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/sqs/goreturns"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
)
