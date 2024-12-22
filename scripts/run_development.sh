#!/bin/bash

# Get the project root directory (assuming script is in project/scripts/)
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"
go run "$PROJECT_ROOT/cmd/user-service/main.go"
