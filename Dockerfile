FROM golang:1.24.6-bookworm@sha256:2679c15c940573aded505b2f2fbbd4e718b5172327aae3ab9f43a10a5c700dfc as builder

WORKDIR /app
COPY . .
RUN make build

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/kube-audit-mcp /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/kube-audit-mcp"]
