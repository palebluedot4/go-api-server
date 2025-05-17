#!/bin/sh

set -e

echo "Running pre-commit hook..."

echo "Formatting Go files..."
go fmt ./...

echo "Tidying go modules..."
go mod tidy

if
    command -v golangci-lint &
    >/dev/null
then
    echo "Running linter (golangci-lint)..."
    golangci-lint run --fast ./...
else
    echo "Linter (golangci-lint) not found, skipping lint check"
fi

echo "Running go build..."
go build -a -trimpath -buildvcs=false -ldflags="-s -w" -installsuffix cgo -o main cmd/api/main.go

echo "Pre-commit checks passed"
exit 0
