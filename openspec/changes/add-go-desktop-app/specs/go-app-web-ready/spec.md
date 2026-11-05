## ADDED Requirements

### Requirement: Camada de serviço neutra de transporte
Os casos de uso SHALL ficar em `internal/app` como métodos Go com `context.Context`, entradas e saídas em structs
serializáveis em JSON e erros tipados; os serviços do Wails, as ferramentas MCP e a CLI MUST ser adaptadores finos
sobre essa camada, sem regras próprias.

#### Scenario: Mesma regra em todos os transportes
- **WHEN** a mesma operação inválida é pedida pela UI, pela CLI e pelo MCP
- **THEN** os três recebem o mesmo erro tipado com o mesmo texto

### Requirement: Frontend com adaptador de API
A SPA SHALL acessar o backend só por `frontend/src/api/` com uma implementação Wails (bindings gerados e eventos) e
uma implementação HTTP (fetch, SSE ou WebSocket), escolhidas no build (`--mode desktop` ou `--mode web`). Código só do
desktop MUST ficar em `frontend/src/desktop/` e MUST NOT entrar no build web.

#### Scenario: Build web sem código desktop
- **WHEN** o CI roda `make go-web-check`
- **THEN** o build web não contém `@wailsio/runtime` nem arquivos de `src/desktop/`

### Requirement: Estado por Workspace
O estado do usuário (conexões, sessões XAPI, certificados aceitos, operações) SHALL ficar num `Workspace` do núcleo,
e a camada `internal/app` MUST receber o `Workspace` pelo `context.Context` de cada chamada, nunca por variável
global; o desktop usa um único `Workspace`.

#### Scenario: Dois workspaces isolados
- **WHEN** um teste chama a camada `internal/app` com dois `Workspace` diferentes conectados a servidores diferentes
- **THEN** cada chamada vê só as conexões, operações e certificados do seu `Workspace`

### Requirement: Servidor web em mudança própria
O servidor web multiusuário MUST ficar para uma mudança OpenSpec própria; esta mudança SHALL só documentar em
`docs/web.md` o spike do build `-tags server` do Wails v3 com autenticação por sessão no servidor.

#### Scenario: Comando ainda não disponível
- **WHEN** o usuário roda `xcp-ng-center-go serve` nesta mudança
- **THEN** o comando não existe e nenhuma porta é aberta
