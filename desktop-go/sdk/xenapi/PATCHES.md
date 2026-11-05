# Local patches to the XenServer Go SDK

`sdk/xenapi` is the Go SDK of the XenServer SDK **26.16.1** (`XenServerGo/src` of `XenServer-SDK-26.16.01.zip`, BSD
2-Clause, Cloud Software Group), vendored as its own module (`go/xenapi`) and used by the app through a `replace`.
Every change is marked `XCP-NG PATCH` in the code and kept as one diff from the pristine SDK in
`sdk/patches/xcp-ng.patch`: `tools/sdk-update.sh` applies it to a new SDK in a staging folder and replaces
`sdk/xenapi` only when it applies cleanly. After changing a patch, regenerate that file.

| Patch | Where | Why |
| --- | --- | --- |
| `ClientOpts.HTTPClient` | `jsonrpc_client.go` | The SDK otherwise builds its own client with `InsecureSkipVerify` when no CA is given. The app passes its client: TLS 1.2+, its certificate trust (TOFU, pins), no redirects, timeouts. |
| No `log.Fatal` | `jsonrpc_client.go`, `session.go` | A bad option ended the whole process. `newJSONRPCClient` returns the error; `NewSessionE` exposes it; `NewSession` keeps its signature and fails the calls instead. |
| `APIError` | `jsonrpc_client.go` | XAPI failures were flattened into a string. `*APIError` keeps the code (`HOST_IS_SLAVE`, `SESSION_AUTHENTICATION_FAILED`...) and its parameters for `errors.As`. |
| `Session.Ref()` | `session.go` | The session reference for the HTTP calls of the console and of import/export. |

Not yet: `DecodeRecord` (exported record decoding for `event.from`), planned with the model cache (task 2.4).
