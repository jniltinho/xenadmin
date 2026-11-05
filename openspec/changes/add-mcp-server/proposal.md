## Why

Operações repetitivas de XCP-ng (listar o pool, ligar/desligar VMs, snapshots, migração, modo manutenção, storage, alertas)
hoje exigem o app gráfico ou a CLI `xe` no dom0. Agentes de IA (Claude Code, Codex e outros clientes MCP) poderiam fazer
isso com segurança se o app expusesse as mesmas regras que já validam cada operação na interface. O núcleo do app
desktop (`XenAdmin.Desktop.Core`) já não depende de UI, então um servidor MCP no mesmo binário reaproveita tudo.

## What Changes

- Novo modo `xcp-ng-center mcp`: servidor MCP (Model Context Protocol) por stdio, sem abrir janela, no mesmo binário
  das versões Linux, macOS e Windows.
- Ferramentas: conectar/desconectar servidores, listar e detalhar pools, hosts, VMs e SRs, operações de energia de VMs,
  snapshot, migração entre hosts do pool, reboot/shutdown e modo manutenção de hosts, criação de SR (NFS, SMB, iSCSI e
  bibliotecas de ISO), alertas (listar e dispensar) e dados de desempenho.
- As ferramentas usam as mesmas regras do app (`VmRules`, `HostRules`, `MigrateRules`, `SrRules`,
  `NotificationRules`, `PerformanceRules`) e as mesmas Actions do XenModel.
- Segurança: servidores e credenciais definidos pelo operador no ambiente (nunca pelas ferramentas), senha só em
  memória, certificado confiado fora das ferramentas, modo somente leitura por padrão (`--allow-write` habilita a
  escrita, que ainda exige `confirm: true`), confirmação específica para formatar LUNs.
- Fora de escopo: console VNC, wizards de VM/import/export, HA/WLB, atualizações.

## Capabilities

### New Capabilities
- `desktop-mcp-server`: servidor MCP por stdio que expõe as operações do app a agentes, com as mesmas regras e
  confirmações.

### Modified Capabilities
(nenhuma)

## Impact

- **Código**: novos `XenAdmin.Desktop.Core/Mcp/` (protocolo e ferramentas) e um desvio em `XenAdmin.Desktop/Program.cs`
  para o modo `mcp`; testes em `tests/XenAdmin.Desktop.Tests/`.
- **Dependências**: nenhuma nova (JSON-RPC com o `Newtonsoft.Json` já usado pelo XenModel).
- **Docs**: README do app (configuração em clientes MCP), CHANGELOG e AGENTS.md.
