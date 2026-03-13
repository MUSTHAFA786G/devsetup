BINARY      := devsetup
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
COMMIT_SHA  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -s -w \
  -X 'github.com/devsetup/devsetup/cmd/devsetup.Version=$(VERSION)' \
  -X 'github.com/devsetup/devsetup/cmd/devsetup.BuildDate=$(BUILD_DATE)' \
  -X 'github.com/devsetup/devsetup/cmd/devsetup.CommitSHA=$(COMMIT_SHA)'

.PHONY: all build install test lint clean cross help

## all: Build for current platform
all: build

## build: Compile the binary
build:
	@echo ">> Building $(BINARY) $(VERSION)"
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .
	@echo ">> Binary: ./$(BINARY)"

## install: Install to GOPATH/bin
install:
	go install -ldflags="$(LDFLAGS)" .

## test: Run all tests
test:
	go test -v -race ./...

## test-cover: Run tests + generate HTML coverage report
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo ">> coverage.html"

## lint: Run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## clean: Remove build artefacts
clean:
	rm -f $(BINARY) coverage.out coverage.html
	rm -rf dist/

## cross: Build for Linux / macOS / Windows
cross:
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64  go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64   .
	GOOS=linux   GOARCH=arm64  go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY)-linux-arm64   .
	GOOS=darwin  GOARCH=amd64  go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY)-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64  go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64  .
	GOOS=windows GOARCH=amd64  go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe .
	@echo ">> Binaries in dist/"

## help: Show available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
