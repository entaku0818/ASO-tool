.PHONY: test build run docker-build docker-run lint clean

# Run tests in Docker (Go 1.23)
test:
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go test ./... -v

# Run tests with coverage
test-coverage:
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go test ./... -coverprofile=coverage.out
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go tool cover -html=coverage.out -o coverage.html

# Build binary
build:
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go build -o api ./cmd/api

# Run locally (requires local Go 1.23)
run:
	go run ./cmd/api

# Build Docker image
docker-build:
	docker build -t aso-api .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -e DATABASE_URL=$(DATABASE_URL) aso-api

# Lint code
lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run

# Clean build artifacts
clean:
	rm -f api coverage.out coverage.html

# Download dependencies
deps:
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go mod download

# Tidy dependencies
tidy:
	docker run --rm -v $(PWD):/app -w /app golang:1.23-alpine go mod tidy
