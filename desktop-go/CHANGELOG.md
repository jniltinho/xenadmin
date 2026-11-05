# Changelog — XCP-ng Center (Go)

Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions are the `go-v*` tags.

## [Unreleased]

### Added
- One binary (`xcp-ng-center-go`) with Cobra: `server add|list`, `pool|host|vm|template|sr|network list` (table or
  `--json`), `version`, and `mcp`.
- Connection core: server addresses (IPv6; 80/8080 refused), TLS 1.2+ without redirects, trust on first use with
  SHA-256 pins and `--trust-fingerprint`, the pool coordinator redirect, error messages without secrets.
- MCP server on the official Go SDK: `list_servers`, `connect`, `disconnect`, `list_objects` (read-only), with the
  C# app's tool names.
- The XenServer Go SDK 26.16.1 vendored in `sdk/xenapi` with documented patches (the app's HTTP client, no
  `log.Fatal`, typed `APIError`, `Session.Ref`).
