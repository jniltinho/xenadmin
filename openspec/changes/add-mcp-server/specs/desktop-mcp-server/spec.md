## Purpose

Permite que agentes de IA operem pools XCP-ng pelo protocolo MCP usando o próprio app desktop, com as mesmas regras,
confirmações e cuidados com credenciais da interface gráfica.

## ADDED Requirements

### Requirement: Servidor MCP por stdio
O binário `xcp-ng-center` SHALL aceitar o modo `mcp`, que MUST falar MCP (JSON-RPC 2.0, uma mensagem por linha) em
stdin/stdout sem abrir janela, e MUST escrever logs apenas em stderr.

#### Scenario: Inicialização
- **WHEN** um cliente envia `initialize` pedindo a versão `2024-11-05`, depois `notifications/initialized` e `tools/list`
- **THEN** o servidor responde com `protocolVersion` `2024-11-05`, seu nome e versão e a capacidade `tools`, não responde à notificação e lista as ferramentas com seus esquemas

#### Scenario: Batch
- **WHEN** o cliente envia um array de pedidos
- **THEN** o servidor responde `-32600` (batches não são suportados nas versões anunciadas)

#### Scenario: Notificação desconhecida
- **WHEN** o cliente envia uma notificação que o servidor não conhece
- **THEN** o servidor não responde e continua atendendo

#### Scenario: Ferramenta desconhecida
- **WHEN** o cliente chama uma ferramenta inexistente
- **THEN** o servidor responde com erro JSON-RPC `-32602` e continua atendendo

### Requirement: Credenciais e certificados
As credenciais do XAPI MUST vir do ambiente definido pelo operador (`XCP_HOST`, `XCP_USER`, `XCP_PASS`) e MUST ser
usadas só com esses servidores; nenhuma ferramenta aceita credenciais do XAPI (só `sr_create` recebe as senhas do share
SMB ou do CHAP, que vão ao pool como secrets). O servidor MUST manter senhas apenas em memória e MUST NOT
devolvê-las nem registrá-las. Certificados SHALL ser confiados pela mesma política do app, fora das ferramentas, e
portas HTTP (80, 8080) MUST ser recusadas.

#### Scenario: Servidor não configurado
- **WHEN** o agente pede `connect` a um host que não está em `XCP_HOST`
- **THEN** a ferramenta falha e nenhuma credencial é enviada

#### Scenario: Certificado desconhecido
- **WHEN** `connect` encontra um certificado que nem a configuração do app nem `XCP_TRUST_FINGERPRINT` confiam
- **THEN** a ferramenta falha mostrando o fingerprint SHA-256 e como o operador pode confiá-lo, mesmo que o certificado seja válido pela cadeia do sistema

#### Scenario: Certificado alterado
- **WHEN** o fingerprint difere do confiado antes
- **THEN** a conexão é recusada

### Requirement: Operações com as regras do app
Nomes, descrições e textos de alertas vindos dos servidores MUST ser devolvidos apenas como dados. As ferramentas de VMs,
hosts, storage e alertas SHALL validar com as mesmas regras do app e executar as mesmas Actions
do XenModel, devolvendo o resultado ou o motivo da recusa.

#### Scenario: Operação não permitida
- **WHEN** o agente pede `vm_power` com `start` para uma VM que já está ligada
- **THEN** a ferramenta falha explicando que a operação não está permitida, sem chamar o XAPI

#### Scenario: Migração impossível
- **WHEN** o agente pede `vm_migrate` para um host que não pode receber a VM
- **THEN** a ferramenta falha com o motivo mostrado pelo app no menu Migrate

### Requirement: Confirmação e modo somente leitura
O servidor MUST iniciar em modo somente leitura; ferramentas de escrita SHALL existir só com `--allow-write` e MUST
exigir `confirm: true`. Formatar um LUN MUST exigir a confirmação do seu SCSIid.

#### Scenario: Sem confirmação
- **WHEN** o agente chama `vm_power` sem `confirm: true`
- **THEN** a ferramenta falha pedindo confirmação e nada é executado

#### Scenario: Somente leitura
- **WHEN** o servidor roda sem `--allow-write`
- **THEN** `tools/list` mostra apenas ferramentas de leitura e chamadas de escrita são recusadas

#### Scenario: Excluir snapshot com uuid de VM
- **WHEN** o agente chama `snapshot_delete` com o uuid de uma VM que não é snapshot
- **THEN** a ferramenta falha sem chamar o XAPI

#### Scenario: LUN com SR existente
- **WHEN** `sr_create` iSCSI encontra uma SR no LUN e o agente não pede `reattach` nem confirma o SCSIid para formatar
- **THEN** nada é criado nem formatado

### Requirement: Operações longas
Cada operação de escrita SHALL receber um `operation_id`, devolvido imediatamente por padrão; `operation_status` SHALL
informar estado, progresso, resultado ou erro.

#### Scenario: Migração demorada
- **WHEN** o agente chama `vm_migrate`
- **THEN** a resposta traz o `operation_id` na hora e `operation_status` mostra quando a migração termina

### Requirement: Listas paginadas
As listas SHALL ser paginadas com ordem estável e SHALL informar o total e o próximo deslocamento.

#### Scenario: Segunda página
- **WHEN** há mais objetos que o `limit` pedido
- **THEN** a resposta traz `next_offset`, e a chamada com esse `offset` devolve os seguintes sem repetir nenhum

