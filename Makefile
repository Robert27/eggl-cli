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
GOPATH := $(shell go env GOPATH)
GORELEASER := $(GOPATH)/bin/goreleaser

.PHONY: build install release-build release-snapshot test fmt check hooks clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

install:
	go install -ldflags "$(LDFLAGS)" .

release-build:
	@mkdir -p $(dir $(OUTPUT))
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(OUTPUT) .

release-snapshot: $(GORELEASER)
	$(GORELEASER) release --snapshot --clean

$(GORELEASER):
	go install github.com/goreleaser/goreleaser/v2@latest

test:
	go test -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | tail -1

fmt:
	gofmt -w .

fmt-check:
	@unformatted="$$(gofmt -l .)"; \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		echo "Run: make fmt"; \
		exit 1; \
	fi

check: fmt-check test

hooks:
	git config core.hooksPath .githooks

clean:
	rm -rf bin/
