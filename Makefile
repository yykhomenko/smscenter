build: ## Build version
	go build -v ./cmd/smscenter
	go build -v ./cmd/smsclient

test: ## Run all tests
	go test -v -race -timeout 30s ./...

run: ## Run version
	go run ./cmd/smscenter

install: ## Install version
	make build
	make test
	go install -v ./cmd/smscenter
	go install -v ./cmd/smsclient

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
