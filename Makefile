.DEFAULT_GOAL := help

SHELL := /bin/bash
TOOL_BIN_DIR  ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT_VERSION := 1.60.3

install-golangci-lint:  ## Install golangci-lint
	@rm -f $(TOOL_BIN_DIR)/golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOL_BIN_DIR) v$(GOLANGCI_LINT_VERSION)


gomodule:  ## Tidy up Golang dependencies, see https://github.com/golang/go/wiki/Modules
	@go mod tidy


gomodule-upgradable:  ## List to view available minor and patch upgrades only for the direct dependencies
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null


gomodule-upgrade:  ## Upgrade to use the latest minor or patch releases (and add -t to also upgrade test dependencies)
	go get -u ./...


gomodule-upgrade-patch:  ## Upgrade to use the latest patch releases (and add -t to also upgrade test dependencies)
	go get -u=patch ./...


test:  ## Test to all of directories
	AWS_REGION=ap-northeast-1 AWS_ACCESS_KEY_ID=1 AWS_SECRET_ACCESS_KEY=2 go test -mod=mod -cover -race ./...


linter:  ## Golang completely all of style checking
	@test -f $(TOOL_BIN_DIR)/golangci-lint || make install-golangci-lint
	@if [ "`golangci-lint run -c .golangci.yml --timeout 10m0s | tee /dev/stderr`" ]; then \
			echo "^ linter errors!" && echo && exit 1; \
	fi


golint:  ## run golint to all of gofiles
	@go get -v golang.org/x/lint/golint 2> /dev/null
	@go install golang.org/x/lint/golint
	@if [ "`golint ./... | tee /dev/stderr`" ]; then \
		echo "^ golint errors!" && echo && exit 1; \
	fi


misspell:  ## Check misspelling to files except go files
	@go get -v github.com/client9/misspell/cmd/misspell 2> /dev/null
	@go install github.com/client9/misspell/cmd/misspell
	@if [ "`find . -type f | xargs misspell -error | tee /dev/stderr`" ]; then \
		echo "^ misspell errors!" && echo && exit 1; \
	fi


format:  ## Run go formater
	@go install golang.org/x/tools/cmd/goimports 2> /dev/null
	@go install github.com/sqs/goreturns 2> /dev/null
	@make format-target target="data/" &
	@make format-target target="util/"


format-target:  ## Run go formater: ${target}
	goimports -w ${target}
	goreturns -w ${target}

circleci-validate:  ## Validate ./circleci/config.yml
	circleci config validate


help:  ## Show all of tasks
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


