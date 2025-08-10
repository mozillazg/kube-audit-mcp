


.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o kube-audit-mcp main.go


.PHONY: test
test:
	go test -v -cover ./...

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
