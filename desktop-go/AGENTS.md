# Agent guide — desktop-go

The Go version of XCP-ng Center (proposal: `openspec/changes/add-go-desktop-app`, validate with
`openspec validate add-go-desktop-app --strict`). The repository's root `AGENTS.md` rules apply: credentials only in
the environment, the lab servers' limits, a Codex review before each commit.

- `make check` before every commit; `make cli` for binaries. Containers only.
- The core (`internal/xapi`, `trust`, `config`, `app`) never imports Cobra, MCP or Wails (`internal/archtest`).
- Port behaviour from the C# app and cite the file it comes from; keep the MCP tool names of `XcpTools.cs`.
- `sdk/xenapi` changes only through `XCP-NG PATCH` blocks documented in `PATCHES.md`.
