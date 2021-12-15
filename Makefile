build: ## Build a version
	go build -v ./cmd/smscenter

test: ## Run all the tests
	go test -v -race -timeout 30s ./...

install: ## Install a version
	make build
	make test
	go install -v ./cmd/smscenter

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
