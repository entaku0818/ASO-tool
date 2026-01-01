#!/bin/bash
set -e

echo "=== Running pre-commit checks ==="

# Backend checks
echo ""
echo "--- Backend ---"
cd "$(dirname "$0")/../backend"

echo "Running Go tests..."
go test ./...

echo "Running golangci-lint..."
golangci-lint run

# Frontend checks
echo ""
echo "--- Frontend ---"
cd "$(dirname "$0")/../frontend"

echo "Running type check..."
npm run type-check

echo "Running lint..."
npm run lint

echo "Running tests..."
npm test

echo ""
echo "=== All checks passed! ==="
