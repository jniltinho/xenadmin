# Development

- Go 1.27; everything runs in `golang:1.27.3` through the Makefile (`make check`, `make cli`).
- Tests use a fake JSON-RPC server over TLS (`internal/xapi/client_test.go`) and an in-memory MCP client
  (`internal/mcp/server_test.go`); none needs a real server.
- Against a server: `XCP_HOST`, `XCP_USER`, `XCP_PASS` in the environment (never in files); the lab servers and rules
  are in the repository's [AGENTS.md](../../AGENTS.md) (read-only `.236`, view-only `.114`).
- A new SDK: `tools/sdk-update.sh XenServer-SDK-<version>.zip` replaces `sdk/xenapi` and re-applies the
  `XCP-NG PATCH` changes listed in `sdk/xenapi/PATCHES.md`.
