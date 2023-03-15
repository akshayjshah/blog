.DEFAULT_GOAL := build
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: build
build:
	@go run .

.PHONY: run
run: build ## Run on :8080
	@wrangler pages dev --port 8080 --live-reload --compatibility-date=2023-03-02 dist

.PHONY: clean
clean: ## Delete intermediate build artifacts
	rm -rf dist/*.html dist/recipes/*.html

.PHONY: upgrade
upgrade: ## Upgrade Go dependencies
	go get -u -t ./...
	go mod tidy -v

.PHONY: lint
lint: ## Lint project
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest
