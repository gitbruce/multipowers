.PHONY: test test-smoke test-unit test-integration test-e2e test-live test-coverage test-all test-plugin-name clean-tests sync help

# Default: smoke + unit (fast feedback)
test: test-smoke test-unit

# Validate plugin name (critical - prevents command prefix breakage)
test-plugin-name:
	@go run ./cmd/mp-devx --action suite --suite smoke

# Run all tests
test-all: test-smoke test-unit test-integration test-e2e

# Smoke tests (pre-commit, <30s)
test-smoke: test-plugin-name
	@echo "Running smoke tests..."
	@go run ./cmd/mp-devx --action suite --suite smoke

# Unit tests (1-2min)
test-unit:
	@echo "Running unit tests..."
	@go run ./cmd/mp-devx --action suite --suite unit

# Integration tests (5-10min)
test-integration:
	@echo "Running integration tests..."
	@go run ./cmd/mp-devx --action suite --suite integration

# E2E tests (15-30min)
test-e2e:
	@echo "Running E2E tests..."
	@go run ./cmd/mp-devx --action suite --suite e2e

# Live tests - real Claude Code sessions (2-5min per test, uses API)
test-live:
	@echo "Running live tests (real Claude Code sessions)..."
	@echo "WARNING: This makes real API calls"
	@go run ./cmd/mp-devx --action suite --suite live

# Performance tests
test-performance:
	@echo "Running performance tests..."
	@go run ./cmd/mp-devx --action suite --suite performance

# Regression tests
test-regression:
	@echo "Running regression tests..."
	@go run ./cmd/mp-devx --action suite --suite regression

# Coverage report
test-coverage:
	@echo "Generating coverage report..."
	@go run ./cmd/mp-devx --action suite --suite coverage

# Verbose mode for debugging
test-verbose:
	@go run ./cmd/mp-devx --action suite --suite all

# Clean test artifacts
clean-tests:
	@echo "Cleaning test artifacts..."
	@rm -rf tests/tmp/
	@rm -f test-results*.xml
	@rm -f coverage*.xml
	@rm -f /tmp/test_*.log
	@echo "Test artifacts cleaned"

# Sync fork main with upstream and rebase custom branch
sync:
	@echo "sync moved to go runtime; use project docs for upstream sync workflow"

# Help
help:
	@echo "Claude Octopus Test Suite"
	@echo ""
	@echo "Usage:"
	@echo "  make test              - Run smoke + unit tests (default)"
	@echo "  make test-all          - Run all test categories"
	@echo "  make test-smoke        - Run smoke tests (<30s)"
	@echo "  make test-unit         - Run unit tests (1-2min)"
	@echo "  make test-integration  - Run integration tests (5-10min)"
	@echo "  make test-e2e          - Run E2E tests (15-30min)"
	@echo "  make test-live         - Run live tests (real Claude sessions)"
	@echo "  make test-performance  - Run performance tests"
	@echo "  make test-regression   - Run regression tests"
	@echo "  make test-coverage     - Generate coverage report"
	@echo "  make test-verbose      - Run all tests with verbose output"
	@echo "  make clean-tests       - Clean test artifacts"
	@echo "  make sync              - Sync main from upstream and rebase custom branch"
	@echo "  make help              - Show this help message"
	@echo ""
	@echo "For more details, see tests/README.md"


.PHONY: build-go test-go lint-go check-go-lines

build-go:
	@mkdir -p .claude-plugin/bin
	@go build -o .claude-plugin/bin/mp ./cmd/mp
	@go build -o .claude-plugin/bin/mp-devx ./cmd/mp-devx

test-go:
	@go test ./...

check-go-lines:
	@go run ./scripts/check-go-file-length.go

lint-go:
	@go vet ./...
