VERSION := $(shell grep 'VERSION' version.go | awk '{ print $$4 }' | tr -d '"')
FILES ?= $(shell find . -type f -name '*.go' ! -path "./vendor/*")

.PHONY: help test version clean vendor vet fmt build release

default: help

help: ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

test: ## run unit tests
	@go test -v ./...

version: ## print the version of the project
	@echo ${VERSION}

clean: ## remove build related files
	@go clean
	@rm -f ./out/*

vendor: ## copy go dependencies to vendor directory
	@go mod tidy
	@go mod vendor

vet: ## run go vet on the source files
	@go vet ./...

fmt: vet ## format the project source files
	@go run mvdan.cc/gofumpt -w .
	@go run golang.org/x/tools/cmd/goimports -w $(FILES)
	@go run github.com/google/addlicense -c "c-fraser" -l apache -y "2022" $(FILES)

build: test ## build the CLI application
	@go build -a -o ./out/extendz ./cmd/cli

release: fmt vet test ## release a version of the project
	@git tag -a "v${VERSION}"
	@git push origin "v${VERSION}"
	@go run github.com/goreleaser/goreleaser release
