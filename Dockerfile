FROM golang:1.25.1-bookworm@sha256:2960a1db140a9a6dd42b15831ec6f8da0c880df98930411194cf11875d433021 as builder

WORKDIR /app
COPY . .
RUN make build

FROM gcr.io/distroless/static-debian12@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c

COPY --from=builder /app/kube-audit-mcp /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/kube-audit-mcp"]
