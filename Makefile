BINARY=aads
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build install clean

build:
	go build -ldflags "-X github.com/SaadBelfqih/apple-ads-cli/cmd.Version=$(VERSION)" -o $(BINARY) .

install:
	go install -ldflags "-X github.com/SaadBelfqih/apple-ads-cli/cmd.Version=$(VERSION)" .

clean:
	rm -f $(BINARY)
