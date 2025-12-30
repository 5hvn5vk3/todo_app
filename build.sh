#!/bin/bash
# Render build script

set -euo pipefail

# Build the Go application
go build -o todo_app .
