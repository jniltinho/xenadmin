## ADDED Requirements

### Requirement: CLI Cobra no mesmo binário
O binário SHALL usar Cobra com um arquivo por grupo de comandos em `cmd/` (padrão `ad-manager`/`ova2xva`): sem
argumentos abre a janela; `version`, `server`, `vm`, `host` e `sr` rodam sem janela com as mesmas regras da UI, e
`--json` muda a saída para JSON. A senha MUST vir de `XCP_PASS` ou de um prompt sem eco, nunca de uma flag.

#### Scenario: Listar VMs em JSON
- **WHEN** o usuário roda `xcp-ng-center-go vm list --json` com `XCP_HOST`, `XCP_USER` e `XCP_PASS` definidos
- **THEN** a saída é um array JSON com nome, uuid, estado e host de cada VM, e nenhuma senha aparece na saída

#### Scenario: Console do Windows
- **WHEN** `xcp-ng-center-go.exe version` ou `mcp` roda num terminal ou com stdin/stdout redirecionados no Windows
- **THEN** a saída chega ao terminal ou ao pipe (o exe anexa o console do processo pai ou usa os handles herdados)

#### Scenario: Flags sobre o ambiente
- **WHEN** o usuário roda `xcp-ng-center-go vm list --host 10.0.0.5 --user admin` com `XCP_HOST` apontando para outro servidor
- **THEN** o comando usa `10.0.0.5` e `admin`, e a senha continua vindo de `XCP_PASS` ou do prompt

#### Scenario: Senha em flag
- **WHEN** o usuário tenta passar a senha como flag
- **THEN** a flag não existe e o comando mostra o uso das variáveis de ambiente

### Requirement: Servidor MCP equivalente ao do app C#
`xcp-ng-center-go mcp [--allow-write]` SHALL servir MCP por stdio com `github.com/modelcontextprotocol/go-sdk`, expondo as
mesmas ferramentas e regras de `openspec/changes/add-mcp-server`: somente leitura por padrão, escrita só com `--allow-write` e `confirm: true`,
`operation_id` para operações longas, textos dos servidores tratados como dados e o sanitizador comum de segredos do
núcleo.

#### Scenario: Servidores definidos pelo operador
- **WHEN** o MCP inicia com `--host`/`--user` ou `XCP_HOST`/`XCP_USER`, senha em `XCP_PASS` e certificado já confiado ou `--trust-fingerprint`/`XCP_TRUST_FINGERPRINT`
- **THEN** só esses servidores podem ser conectados, e nenhuma ferramenta aceita servidor, usuário ou senha do XAPI

#### Scenario: Escrita sem confirmação
- **WHEN** o agente chama `vm_power` com `--allow-write` mas sem `confirm: true`
- **THEN** a ferramenta não altera nada e pede a confirmação

#### Scenario: Somente leitura
- **WHEN** o servidor roda sem `--allow-write`
- **THEN** `tools/list` não inclui ferramentas de escrita

#### Scenario: Paridade de ferramentas
- **WHEN** o teste compara os nomes e esquemas de ferramentas do servidor Go com a lista do servidor C#
- **THEN** os nomes são os mesmos e as diferenças de esquema estão documentadas

### Requirement: Um binário para janela, CLI e MCP
O mesmo binário e o mesmo pacote de cada plataforma SHALL abrir a janela sem subcomando e rodar `mcp` e os comandos da
CLI sem abrir janela nem conectar a um display.

#### Scenario: MCP sem display
- **WHEN** `xcp-ng-center-go mcp` roda num Linux sem `DISPLAY` nem `WAYLAND_DISPLAY`
- **THEN** o servidor MCP inicia e responde a `initialize` sem tentar abrir janela

#### Scenario: AppImage como servidor MCP
- **WHEN** o AppImage roda com o subcomando `mcp` numa máquina sem GTK4/WebKitGTK instalados
- **THEN** o servidor MCP inicia com as bibliotecas que o AppImage traz
