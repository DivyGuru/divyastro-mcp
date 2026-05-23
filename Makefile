# Makefile for divyastro-mcp — the Model Context Protocol server
# that exposes DivyAstroAPI endpoints as tools for AI clients.

BINARY  := divyastro-mcp
VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -trimpath -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -s -w"

.PHONY: build build-all test fmt vet tidy clean

## Build for the current platform.
## Override version: make build VERSION=0.2.1
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BINARY) ./cmd/divyastro-mcp/
	@echo "Built: $(BINARY) ($(shell du -sh $(BINARY) 2>/dev/null | cut -f1))"

## Cross-compile for all 5 release platforms into dist/.
## Pure Go (CGO disabled) — no C cross-toolchain needed.
## Usage: make build-all VERSION=0.2.1
build-all:
	@mkdir -p dist
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64    ./cmd/divyastro-mcp/
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64    ./cmd/divyastro-mcp/
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64     ./cmd/divyastro-mcp/
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY)-linux-arm64     ./cmd/divyastro-mcp/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/divyastro-mcp/
	@cd dist && shasum -a 256 $(BINARY)-* > $(BINARY)-checksums.txt
	@echo "Release artifacts in dist/ (version=$(VERSION)):"
	@ls -lh dist/$(BINARY)-*

test:
	go test ./... -count=1

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf dist/ $(BINARY)
