.DEFAULT_GOAL := help

SHELL := /bin/bash



gomodule:  ## Tidy up Golang dependencies, see https://github.com/golang/go/wiki/Modules
	@go mod tidy


gomodule-upgradable:  ## List to view available minor and patch upgrades only for the direct dependencies
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null


gomodule-upgrade:  ## Upgrade to use the latest minor or patch releases (and add -t to also upgrade test dependencies)
	go get -u ./...


gomodule-upgrade-patch:  ## Upgrade to use the latest patch releases (and add -t to also upgrade test dependencies)
	go get -u=patch ./...


help:  ## Show all of tasks
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


