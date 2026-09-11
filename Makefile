.PHONY: help build test lint fmt vet coverage clean check-all

help:
	@echo "Radic Development Commands"
	@echo "=========================="
	@echo "  make build       - Build the application"
	@echo "  make test        - Run tests with coverage"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make fmt         - Format code"
	@echo "  make vet         - Run go vet"
	@echo "  make coverage    - Show test coverage"
	@echo "  make check-all   - Run all checks (fmt, vet, lint, test)"
	@echo "  make clean       - Clean build artifacts"

build:
	@echo "🔨 Building..."
	go build -v ./cmd/...

test:
	@echo "✅ Testing (race detection enabled)..."
	go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "🔍 Linting with golangci-lint..."
	golangci-lint run --config .golangci.yml

fmt:
	@echo "📝 Formatting code..."
	gofmt -s -w .

vet:
	@echo "🔎 Running go vet..."
	go vet ./...

coverage: test
	@echo "📊 Generating coverage..."
	go tool cover -func=coverage.out | tail -1

check-all: fmt vet lint test
	@echo "✅ All checks passed!"

clean:
	@echo "🧹 Cleaning..."
	rm -f coverage.out
	go clean -cache -testcache
