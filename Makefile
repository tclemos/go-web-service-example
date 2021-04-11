.PHONY: run
run: ## runs the application
	go run main.go

.PHONY: build
build: ## builds the application
	go build ./...

.PHONY: test
test: ## runs short tests
	go test ./... -short

.PHONY: test-all
test-all: ## runs all tests
	go test ./...

.PHONY: install-generate-sql
install-generate-sql: ## install code generator for SQL files.
	go get github.com/kyleconroy/sqlc/cmd/sqlc

.PHONY: generate-sql
generate-sql: ## generate persistence Go files on adapters/postgres/db based on adapters/postgres/migrations and adapters/postgres/queries .sql files
	sqlc generate -f "./adapters/postgres/sqlc.yaml"

.PHONY: install-generate-api
install-generate-api: ## install code generator for API files.
	go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

.PHONY: generate-api
generate-api: ## generate api Go http server and client files based on docs/api/spec.yaml
	oapi-codegen -generate=types,server -package api docs/api/spec.yaml > adapters/http/api/api.gen.go && \
	oapi-codegen -generate=types,client -package thingshttpclient docs/api/spec.yaml > adapters/http/client/things.go

.PHONY: generate
generate: generate-api generate-sql ## run all generate commands

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.

.DEFAULT_GOAL := help

help: ## prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'