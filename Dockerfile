FROM golang:1.24.6-bookworm@sha256:2679c15c940573aded505b2f2fbbd4e718b5172327aae3ab9f43a10a5c700dfc as builder

WORKDIR /app
COPY . .
RUN make build

FROM busybox:latest@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74

COPY --from=builder /app/kube-audit-mcp /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/kube-audit-mcp"]
