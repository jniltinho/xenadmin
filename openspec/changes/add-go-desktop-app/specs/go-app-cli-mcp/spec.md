## ADDED Requirements

### Requirement: CLI Cobra no mesmo binário
O binário SHALL usar Cobra com um arquivo por grupo de comandos em `cmd/` (padrão `ad-manager`/`ova2xva`): sem
argumentos abre a janela; `version` mostra versão, commit e data (via `internal/buildinfo`); `server`, `vm`, `host` e
`sr` listam e operam objetos com as mesmas regras da UI; `--json` muda a saída para JSON. As credenciais MUST vir do
ambiente (`XCP_HOST`, `XCP_USER`, `XCP_PASS`) ou de um prompt sem eco, nunca de flags.

#### Scenario: Listar VMs em JSON
- **WHEN** o usuário roda `xcp-ng-center-go vm list --json` com `XCP_HOST`, `XCP_USER` e `XCP_PASS` definidos
- **THEN** a saída é um array JSON com nome, uuid, estado e host de cada VM, e nenhuma senha aparece na saída

#### Scenario: Console do Windows
- **WHEN** `xcp-ng-center-go.exe version` ou `mcp` roda num terminal ou com stdin/stdout redirecionados no Windows
- **THEN** a saída chega ao terminal ou ao pipe (o exe GUI anexa o console do processo pai), e o `xcp-ng-center-go-cli.exe` sem interface faz o mesmo sem anexar nada

#### Scenario: Senha em flag
- **WHEN** o usuário tenta passar a senha como flag
- **THEN** a flag não existe e o comando mostra o uso das variáveis de ambiente

### Requirement: Servidor MCP equivalente ao do app C#
`xcp-ng-center-go mcp [--allow-write]` SHALL servir MCP por stdio com `github.com/modelcontextprotocol/go-sdk`, expondo as
mesmas ferramentas e regras de `openspec/changes/add-mcp-server`: servidores e credenciais só do ambiente, certificado
já confiado ou `XCP_TRUST_FINGERPRINT`, somente leitura por padrão, escrita só com `--allow-write` e `confirm: true`,
`operation_id` para operações longas, textos dos servidores tratados como dados e o sanitizador comum de segredos do
núcleo.

#### Scenario: Escrita sem confirmação
- **WHEN** o agente chama `vm_power` com `--allow-write` mas sem `confirm: true`
- **THEN** a ferramenta não altera nada e pede a confirmação

#### Scenario: Somente leitura
- **WHEN** o servidor roda sem `--allow-write`
- **THEN** `tools/list` não inclui ferramentas de escrita

#### Scenario: Paridade de ferramentas
- **WHEN** o teste compara os nomes e esquemas de ferramentas do servidor Go com a lista do servidor C#
- **THEN** os nomes são os mesmos e as diferenças de esquema estão documentadas

### Requirement: Binário sem interface para MCP e CLI
O binário `xcp-ng-center-go-cli` (módulo raiz, sem CGO) SHALL oferecer os mesmos comandos, exceto abrir a janela, e a
configuração recomendada nos clientes MCP SHALL usá-lo.

#### Scenario: Servidor sem interface gráfica
- **WHEN** `xcp-ng-center-go-cli mcp` roda num Linux sem GTK nem WebKit
- **THEN** o servidor MCP inicia e responde a `initialize`
