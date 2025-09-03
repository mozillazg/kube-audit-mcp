FROM golang:1.24.7-bookworm@sha256:b37bdf7efa519791172a520fe8d2bd2d0eead166bdf5d299358abe4ba3e24e1b as builder

WORKDIR /app
COPY . .
RUN make build

FROM busybox:latest@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b

COPY --from=builder /app/kube-audit-mcp /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/kube-audit-mcp"]
