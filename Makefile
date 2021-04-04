.PHONY: run
run: ## runs the application
	go run main.go

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
generate-from-sql: ## generate Go structs from SQL files.
	sqlc generate -f "./actors/postgres/sqlc.yaml"

.PHONY: install-generate-api
install-generate-api: ## install code generator for API files.
	go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

.PHONY: generate-api
generate-api: ## generate api.gen.go file based on doc/api/spec.yaml
	oapi-codegen -generate=types,server,spec -package http docs/http/spec.yaml > actors/http/generated.go
## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.

.DEFAULT_GOAL := help

help: ## prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'