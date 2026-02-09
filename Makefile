BINARY := aads
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
PKG := github.com/SaadBelfqih/apple-ads-cli/cmd
LDFLAGS := -X $(PKG).Version=$(VERSION) -X $(PKG).Commit=$(COMMIT) -X $(PKG).Date=$(DATE)

# Install destination for `make install`.
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: build install clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(GOBIN)/$(BINARY)

clean:
	rm -f $(BINARY)
