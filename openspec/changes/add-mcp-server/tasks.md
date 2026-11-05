## 1. Protocolo e bootstrap

- [x] 1.1 `McpServer`: JSON-RPC 2.0 por linha, negociação de `protocolVersion`, notificações sem resposta, IDs string/número, códigos -32700/-32601/-32602, `isError`, pedidos em paralelo com respostas serializadas
- [x] 1.2 `McpHost`: bootstrap headless (config, `InvokeHelper`, Folders, TLS do app com membros do pool)
- [x] 1.3 Modo `xcp-ng-center mcp [--allow-write]` em `Program.Main`, sem Avalonia; logs em stderr
- [x] 1.4 Testes unitários do protocolo com streams em memória (handshake, notificações, erros, read-only)

## 2. Conexão

- [x] 2.1 `list_servers`, `connect` (só hosts de `XCP_HOST`), `disconnect`; credenciais do ambiente; TLS confiado fora das ferramentas (`XCP_TRUST_FINGERPRINT` ou configuração do app); recusa de HTTP

## 3. Ferramentas

- [x] 3.1 Leitura: `list_objects`, `get_object`, `list_snapshots`, `list_alerts`, `get_performance`, `iscsi_probe`, `operation_status`
- [x] 3.2 VMs: `vm_power`, `vm_snapshot`, `snapshot_revert`, `snapshot_delete`, `vm_migrate`
- [x] 3.3 Hosts: `host_power`, `host_maintenance`
- [x] 3.4 Storage: `sr_create` (nfs, smb, iscsi com probe e reattach/format explícitos, nfs_iso, smb_iso)
- [x] 3.5 Alertas: `dismiss_alerts`
- [x] 3.6 Escrita só com `--allow-write` e `confirm: true`; operações longas com `operation_id`

## 4. Validação

- [x] 4.1 Teste no lab KVM pelo processo real (stdio, container na rede do lab) com TLS real: recusa sem fingerprint, conexão com `XCP_TRUST_FINGERPRINT`, leituras, paginação, guardas (start em VM ligada, sem `confirm`, uuid de VM em `snapshot_delete`, LUN em uso), snapshot criado e excluído
- [x] 4.2 Revisão Codex (quatro passadas; GO)
- [x] 4.3 README (Claude Code / Codex), CHANGELOG, AGENTS.md
