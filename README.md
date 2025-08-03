# kube-audit-mcp

kube-audit-mcp is an Model Context Protocol (MCP) server that gives AI agents, assistants,
and chatbots the ability to query Kubernetes Audit Logs.


## Installation

1. First, download the latest release from this repository.

2. Then, configure the backend of Kubernetes Audit Logs.

  1. You can get an sample config via `kube-audit-mcp sample-config`
  2. xxx


## Transport Options

### STDIO Transport (Default)

The default transport mode uses standard input/output for communication.
This is the standard MCP transport used by most clients like Claude Desktop.

```
# Run with default stdio transport
kube-audit-mcp mcp --config /path/to/config.yaml

# Or explicitly specify stdio
kube-audit-mcp mcp --config /path/to/config.yaml --transport stdio
```

## MCP Clients

Theoretically, any MCP client should work with kube-audit-mcp. Three examples are given below.



## Configurations
