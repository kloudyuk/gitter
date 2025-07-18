# Makefile for gitter
.PHONY: build test cover lint clean help

# Variables
BINARY_NAME=gitter

# Default target
.DEFAULT_GOAL := help

## Build the application
build:
	go build -o $(BINARY_NAME) .

## Run tests
test:
	go test -v -race ./...

## Run tests with coverage
cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Run linter
lint:
	golangci-lint run

## Clean build artifacts
clean:
	rm -f $(BINARY_NAME) coverage.out coverage.html *.log

## Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the application"
	@echo "  test   - Run tests"
	@echo "  cover  - Run tests with coverage report"
	@echo "  lint   - Run golangci-lint"
	@echo "  clean  - Clean build artifacts"
