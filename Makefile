.PHONY: test

help:
	@echo "Available commands:"
	@echo "  make test - Run all tests"

test:
	go test ./...
