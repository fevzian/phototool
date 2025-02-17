# Simple Makefile for a Go project

# Build the application
all: build test

build: clean
	@mkdir build
	@echo "Building..."
	@if [ "$(shell go env GOOS)" = "windows" ]; then \
		go build -o ./build/phototool.exe cmd/main.go; \
	else \
		go build -o ./build/phototool cmd/main.go; \
	fi
	@echo "Build complete"

# Run the application
run:
	@go run cmd/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf ./build

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

.PHONY: all build run test clean watch
