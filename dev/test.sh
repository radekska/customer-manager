#!/bin/bash

# Run tests a root
docker-compose -f ./docker/docker-compose.yml exec -T -u 0 customer-manager go test -coverprofile=/tmp/customer-manager.out ./...
