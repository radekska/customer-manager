#!/bin/bash

set -e

golangci-lint run --config .golangci.yml --out-format=colored-line-number

# Run tests a root
docker-compose -f ./docker/docker-compose.yml exec -T -u 0 customer-manager go test -coverprofile=/tmp/customer-manager.out ./...
