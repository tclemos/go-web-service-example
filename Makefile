.PHONY: run
run: ## runs the application
	go run main.go

.PHONY: test
test: ## runs short tests
	go test ./... -short

.PHONY: test-all
test-all: ## runs all tests
	go test ./...

.PHONY: generate-sql
generate-from-sql: ## generate Go structs from SQL files.
	sqlc generate -f "./actors/postgres/sqlc.yaml"

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.

.DEFAULT_GOAL := help

help: ## prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'