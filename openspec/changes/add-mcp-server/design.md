## Context

O app desktop já separa a lógica (Core, sem Avalonia) da interface. O XenModel roda sem UI quando `InvokeHelper` recebe
um sincronizador de thread própria (como nos testes de integração). O servidor MCP usa esse caminho, mas com a mesma
política de TLS do app (os testes aceitam qualquer certificado; o MCP não).

## Goals / Non-Goals

- Goals: operar pools XCP-ng por agentes com as validações do app; um único binário; nenhuma dependência nova; nenhum
  segredo gravado localmente nem devolvido ao agente.
- Non-Goals: console, wizards de VM, import/export, HA/WLB, atualizações; transporte HTTP (só stdio); cancelamento de
  operações (o XAPI não desfaz uma migração ou um modo manutenção pela metade).

## Decisions

### D1. Mesmo binário, modo `mcp`
`Program.Main` desvia para o servidor quando o primeiro argumento é `mcp`, antes do Avalonia. stdout é só do protocolo;
logs vão para stderr. Opções: `--allow-write` (habilita as ferramentas de escrita; padrão é somente leitura).

### D2. Bootstrap headless explícito (`McpHost`)
Antes de atender: `PlatformServices`, `XenAdminConfigManager.Provider`, `InvokeHelper` com um sincronizador de thread
própria, `Folders.InitFolders` nessa thread, e a mesma política TLS do app: `CertificateTrust` com a configuração do app
em `ServicePointManager.ServerCertificateValidationCallback` e em `XenAPI.HTTP.ServerCertificateValidator`, e os
certificados dos membros do pool confiados pelo coordenador autenticado (`VouchedFingerprints`). As ferramentas nunca
esperam Actions na thread de eventos.

### D3. Protocolo
JSON-RPC 2.0, uma mensagem por linha. `initialize` negocia `protocolVersion`: se o cliente pedir `2025-06-18` ou
`2024-11-05`, o servidor responde a mesma; senão `2025-06-18`. `2025-03-26` não é anunciada porque exigiria batches
(arrays de pedidos); um array recebido é respondido com `-32600`. Depois: `notifications/initialized`,
`ping`, `tools/list`, `tools/call`. IDs string ou número são devolvidos como vieram; notificações (sem id) nunca têm
resposta, e as desconhecidas são ignoradas. JSON inválido → `-32700`; método desconhecido → `-32601`; parâmetros
inválidos ou ferramenta desconhecida → `-32602`. Falha de uma ferramenta → resultado com `isError: true`. Resultados:
`content: [{type: "text", text: <JSON>}]` (compatível com as três versões). Pedidos são atendidos em paralelo e as
respostas serializadas, então `ping` responde enquanto uma operação longa roda.

### D4. Credenciais e servidores vinculados pelo operador
As ferramentas não recebem senha nem endereço livre. O operador define `XCP_HOST` (um ou mais, separados por vírgula),
`XCP_USER` e `XCP_PASS` no ambiente do processo; `connect` só aceita um desses hosts (ou o coordenador para o qual o
XAPI redirecionar um deles). A senha fica só em memória e é removida de qualquer texto devolvido ou logado.

### D5. TLS fora das ferramentas
O certificado precisa já estar fixado (pin) na configuração do app (o usuário conectou pela interface e conferiu o
fingerprint) ou vir em `XCP_TRUST_FINGERPRINT` (definido pelo operador). No modo MCP, `CertificateTrust.RequirePin`
desliga o atalho do app que aceita e fixa sozinho um certificado válido pela cadeia do sistema; continuam aceitos só os
fixados e os dos membros do pool atestados pelo coordenador já autenticado. Certificado desconhecido ou alterado →
`connect` falha mostrando o fingerprint e o que o operador deve fazer; nenhuma ferramenta confia em certificado.
Portas 80/8080 são recusadas.

### D6. Escrita
Desligada por padrão; com `--allow-write`, cada ferramenta de escrita ainda exige `confirm: true`, e ações que destroem
dados exigem uma confirmação específica (ex.: `format_lun_scsi_id` igual ao SCSIid do LUN). A aprovação humana fica no
cliente MCP (ex.: Claude Code pede permissão por ferramenta); o README orienta mantê-la.

### D7. Dados não confiáveis
Nomes, descrições, tags e textos de alertas vêm dos servidores e podem conter texto malicioso. As descrições das
ferramentas são fixas e avisam disso; os resultados trazem esses campos como dados JSON, nunca em mensagens do servidor
MCP, e o texto longo é truncado.

### D8. Ferramentas
Objetos são indicados por `uuid` (obrigatório para escrita; leitura também aceita nome exato, ambíguo é erro com os
candidatos). Listas são paginadas (`offset`, `limit` até 500, padrão 100), em ordem estável (tipo, nome, uuid), e
devolvem `total` e `next_offset` (ausente na última página).
- Leitura: `list_servers`, `connect`, `disconnect`, `list_objects`, `get_object`, `list_snapshots`, `list_alerts`,
  `get_performance`, `iscsi_probe` (IQNs, LUNs e a SR já existente em cada LUN), `operation_status`.
- Escrita: `vm_power`, `vm_snapshot`, `snapshot_revert`, `snapshot_delete`, `vm_migrate`, `host_power`,
  `host_maintenance`, `sr_create`, `dismiss_alerts`.
- `snapshot_revert` e `snapshot_delete` só aceitam uuid de snapshot (`is_a_snapshot`), guarda de
  `DeleteSnapshotCommand` extraída para o Core: o uuid de uma VM comum é recusado antes de qualquer chamada.
- `sr_create` para iSCSI exige `scsi_id` e `mode`: `reattach` (usa a SR achada no probe) ou `format` (só com
  `format_lun_scsi_id` igual ao `scsi_id`); NFS/SMB criam uma SR nova (o share continua com o que tinha) e as
  bibliotecas de ISO só são criadas.

### D9. Operações longas
Toda escrita roda como Action do XenModel e recebe um `operation_id`, devolvido na hora por padrão
(`wait_seconds` 0; até 50 se o agente pedir, abaixo do tempo limite comum de 60 s dos clientes). `operation_status`
(com `wait_seconds` opcional até 50) informa estado, progresso, resultado ou erro. Expirar a espera não desfaz nada.

### D10. Segredos remotos e senhas de storage
Senhas de SMB e CHAP passadas a `sr_create` são enviadas ao XAPI, que as guarda como secrets do pool (como no app);
nunca são gravadas localmente nem devolvidas.

## Risks / Trade-offs

- Agente sob prompt injection: mitigado por escrita desligada por padrão, `confirm`, confirmações específicas para
  formatar, aprovação do cliente e servidores fixados pelo operador.
- Operações longas sem cancelamento: o estado fica consultável por `operation_status`.
