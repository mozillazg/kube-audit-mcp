

# Get version from git tag, fallback to "dev" if no tags
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Binary name
BINARY_NAME := kube-audit-mcp

# Go build flags
LDFLAGS := -s -w
LDFLAGS += -X github.com/mozillazg/kube-audit-mcp/pkg/cli.version=$(VERSION)
LDFLAGS += -X github.com/mozillazg/kube-audit-mcp/pkg/cli.commit=$(COMMIT)


.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" \
		-o $(BINARY_NAME) main.go


.PHONY: test
test:
	go test -v -cover ./...


.PHONY: lint
lint: vendor
	go vet ./...
	gofmt -s -w ./...


.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
