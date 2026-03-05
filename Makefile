.PHONY: all build-go test test-smoke test-unit test-integration test-e2e test-coverage test-all test-plugin-name clean-tests parity help check-go-lines lint-go

# Default target
all: build-go test-smoke test-unit

# Build production and devx binaries
build-go:
	@echo "Building binaries..."
	@mkdir -p .claude-plugin/bin
	@go build -o .claude-plugin/bin/mp ./cmd/mp
	@go build -o .claude-plugin/bin/mp-devx ./cmd/mp-devx
	@echo "Binaries built in .claude-plugin/bin/"

# Run consistency parity check
parity: build-go
	@./.claude-plugin/bin/mp-devx -action parity

# Default test: smoke + unit (fast feedback)
test: test-smoke test-unit

# Smoke tests (critical path, <30s)
test-smoke: build-go
	@echo "Running smoke tests..."
	@./.claude-plugin/bin/mp-devx --action suite --suite smoke

# Unit tests (full internal logic)
test-unit:
	@echo "Running unit tests..."
	@go test -v ./internal/...

# Integration tests
test-integration: build-go
	@echo "Running integration tests..."
	@./.claude-plugin/bin/mp-devx --action suite --suite integration

# Coverage report
test-coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out

# Linting
lint-go:
	@echo "Running go vet..."
	@go vet ./...

check-go-lines:
	@go run ./scripts/check-go-file-length.go

# Clean test artifacts
clean-tests:
	@echo "Cleaning test artifacts..."
	@rm -rf tests/tmp/
	@rm -f coverage.out
	@echo "Test artifacts cleaned"

# Help
help:
	@echo "Multipowers Go Native Test & Build Suite"
	@echo ""
	@echo "Usage:"
	@echo "  make build-go          - Build mp and mp-devx binaries"
	@echo "  make parity            - Run plugin namespace parity check"
	@echo "  make test              - Run smoke + unit tests (default)"
	@echo "  make test-smoke        - Run critical smoke tests"
	@echo "  make test-unit         - Run Go unit tests"
	@echo "  make test-coverage     - Generate and display coverage"
	@echo "  make clean-tests       - Remove test artifacts"
	@echo "  make help              - Show this help message"
