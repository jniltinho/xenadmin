# XCP-ng Center (Go). Everything runs in a container (only docker and make needed), as the C# app's Makefile.
# The desktop window (Wails, task 1.3) is not wired yet: the binary is the CLI and the MCP server.
GO_IMAGE ?= golang:1.27.3
VERSION  ?= $(shell git describe --tags --match 'go-v*' 2>/dev/null | sed 's/^go-v//' || echo dev)
LDFLAGS  := -s -w -X xcpng-center/internal/buildinfo.Version=$(VERSION) -X xcpng-center/internal/buildinfo.Commit=$(shell git rev-parse --short HEAD 2>/dev/null)
RUN      := docker run --rm -v $(CURDIR):/src -w /src -e GOFLAGS=-buildvcs=false -e GOCACHE=/src/.cache/go-build -e GOMODCACHE=/src/.cache/mod $(GO_IMAGE)

.PHONY: test check cli dist clean

test:            ## unit tests of the app and the vendored SDK
	$(RUN) sh -c 'go test ./... && cd sdk/xenapi && go test ./...'

check:           ## gofmt, go vet (both modules) and the tests
	$(RUN) sh -c 'test -z "$$(gofmt -l cmd internal main.go)" || { gofmt -l cmd internal main.go; exit 1; }; go vet ./... && (cd sdk/xenapi && go vet ./...) && go test ./...'

cli:             ## static CLI/MCP binaries (no desktop window) for Linux, macOS and Windows in dist/
	$(RUN) sh -c 'for t in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
	  os=$${t%/*}; arch=$${t#*/}; ext=; [ $$os = windows ] && ext=.exe; \
	  CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -tags nogui -trimpath -ldflags "$(LDFLAGS)" -o dist/xcp-ng-center-go-$(VERSION)-$$os-$$arch$$ext . || exit 1; done'

dist: check cli

clean:
	rm -rf dist .cache
