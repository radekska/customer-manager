SHELL := /bin/bash
.PHONY: help

help: ## Show help
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

start: ## Start the application
	dev/start.sh

tests: ## Run unit tests
	dev/test.sh

coverage: tests ## Run tests with coverage
	@go tool cover -html=/tmp/customer-manager.out

static-analysis: ## Run static analysis configured by .golangci.yml
	@golangci-lint run --config .golangci.yml --out-format=colored-line-number

fixer: ## Adjust line length to 120 and fix static analysis errors
	@golines -w --max-len=120 server/ repositories/ database/ cmd/
	@golangci-lint run --config .golangci.yml --fix --out-format=colored-line-number

api-docs: ## Generate API docs
	swag init
