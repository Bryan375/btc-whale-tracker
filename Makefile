.PHONY: lint fmt check test run

fmt:
	@echo "🔧 Formatting code..."
	@gofmt -s -w .
	@goimports -w -local github.com/Bryan375/btc-whale-tracker .

lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

check: fmt lint
	@echo "✅ All checks passed!"

run:
	@go run cmd/main.go
