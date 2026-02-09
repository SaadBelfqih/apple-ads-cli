BINARY := aads
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Install destination for `make install`.
GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: build install clean

build:
	go build -ldflags "-X github.com/SaadBelfqih/apple-ads-cli/cmd.Version=$(VERSION)" -o $(BINARY) .

install: build
	install -m 0755 $(BINARY) $(GOBIN)/$(BINARY)

clean:
	rm -f $(BINARY)
