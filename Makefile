.PHONY: lint fmt check test run install-hooks

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

install-hooks:
	@echo "📦 Installing git hooks..."
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo '#!/bin/sh\nmake test' > .git/hooks/pre-push
	@chmod +x .git/hooks/pre-push
	@echo "✅ Git hooks installed!"
