BINARY := golinks
PKG := ./...

.PHONY: build run dev test vet fmt tidy clean docker-build help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary
	go build -o $(BINARY) .

run: ## Run the server
	go run . serve

dev: ## Run with live reload (air)
	air

test: ## Run tests
	go test $(PKG)

vet: ## Run go vet
	go vet $(PKG)

fmt: ## Format code
	gofmt -w .

tidy: ## Tidy go modules
	go mod tidy

clean: ## Remove build artifacts
	rm -rf $(BINARY) tmp/

docker-build: ## Build Docker image
	docker build . -t $(BINARY):latest
