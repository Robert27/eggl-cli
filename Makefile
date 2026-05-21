MODULE := github.com/Robert27/eggl-cli
BINARY := eggl
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(MODULE)/cmd.version=$(VERSION) \
           -X $(MODULE)/cmd.commit=$(COMMIT) \
           -X $(MODULE)/cmd.date=$(DATE)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUTPUT ?= bin/$(BINARY)

.PHONY: build install release-build test clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

install:
	go install -ldflags "$(LDFLAGS)" .

release-build:
	@mkdir -p $(dir $(OUTPUT))
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) .

test:
	go test ./...

clean:
	rm -rf bin/
