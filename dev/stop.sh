#!/bin/bash

# Stop the application
docker-compose -f ./docker/docker-compose.yml down --remove-orphans -t 0
