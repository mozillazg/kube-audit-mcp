FROM golang:1.24.6-bookworm@sha256:2679c15c940573aded505b2f2fbbd4e718b5172327aae3ab9f43a10a5c700dfc as builder

WORKDIR /app
COPY . .
RUN make build

FROM busybox:latest@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b

COPY --from=builder /app/kube-audit-mcp /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/kube-audit-mcp"]
