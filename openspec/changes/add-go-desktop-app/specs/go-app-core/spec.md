## ADDED Requirements

### Requirement: Núcleo sem UI e sem transporte
O núcleo em `desktop-go/internal/` SHALL conter toda a lógica de conexão, cache, regras e operações, e MUST NOT
importar Wails, `net/http` de servidor, Cobra ou o pacote MCP. Cada regra portada MUST citar no comentário do pacote ou
da função o arquivo C# que espelha (`XenAdmin.Desktop.Core/*.cs` ou o WinForms de origem).

#### Scenario: Dependências do núcleo
- **WHEN** o teste de arquitetura lista os imports dos pacotes `internal/xapi`, `internal/model`, `internal/rules` e `internal/ops`
- **THEN** nenhum deles importa `github.com/wailsapp/wails`, `github.com/spf13/cobra`, `github.com/modelcontextprotocol/go-sdk` nem `internal/server`

### Requirement: SDK XAPI oficial vendorizado
O app SHALL usar o XenServer SDK para Go 26.16.1 vendorizado em `desktop-go/sdk/xenapi` (`replace xenapi => ./sdk/xenapi`).
Alterações locais MUST ser mínimas, estar listadas em `sdk/xenapi/PATCHES.md` e se limitar a: aceitar um
`*http.Client` do chamador (para a política de TLS do app), devolver erros em vez de `log.Fatal`, devolver erros XAPI
tipados (código e parâmetros) e exportar os decodificadores de registros para os snapshots de `event.from`.

#### Scenario: Atualização do SDK
- **WHEN** um SDK mais novo é copiado por cima de `sdk/xenapi` e os patches de `PATCHES.md` são reaplicados
- **THEN** `go build ./...` e os testes do núcleo passam sem outras mudanças

#### Scenario: Erro de URL inválida
- **WHEN** o núcleo cria uma sessão com uma URL que não é válida
- **THEN** recebe um erro e o processo continua rodando

#### Scenario: Erro XAPI tipado
- **WHEN** o servidor responde `HOST_NOT_ENOUGH_FREE_MEMORY` com parâmetros
- **THEN** o núcleo recebe o código e os parâmetros separados e mostra o texto de `Messages` correspondente

### Requirement: Conexão, login e coordenador do pool
O núcleo SHALL conectar por HTTPS (porta 443 padrão, IPv4, IPv6 puro ou `[addr]:porta`) com
`session.login_with_password`, e MUST recusar as portas HTTP 80 e 8080. O cliente HTTP MUST NOT seguir
redirecionamentos automaticamente: um redirecionamento para HTTP ou para um host não confiado MUST falhar sem reenviar
a senha. Um login em membro de pool que responder `HOST_IS_SLAVE` MUST ser refeito no coordenador indicado só depois
de o certificado do coordenador ser confiado.

#### Scenario: Membro do pool
- **WHEN** o usuário conecta a um membro do pool
- **THEN** o núcleo refaz o login no coordenador e a conexão passa a apontar para ele

#### Scenario: Redirecionamento para HTTP
- **WHEN** o servidor responde 302, 307 ou 308 apontando para `http://` ou para outro host
- **THEN** a chamada falha e a senha não é enviada ao destino

#### Scenario: Porta HTTP
- **WHEN** o endereço informado é `host:80` ou `host:8080`
- **THEN** a conexão é recusada antes de enviar a senha

### Requirement: Variantes por versão do XAPI
Funções com variantes por versão (`Host.Evacuate3`...) SHALL ser escolhidas pela versão do XAPI do servidor
(`XAPIVersion` e `APIVersion` da sessão); um servidor de versão desconhecida MUST usar a variante mais antiga
compatível.

#### Scenario: Servidores 8.2 e 8.3
- **WHEN** os testes rodam contra fixtures gravadas de XCP-ng 8.2 e 8.3
- **THEN** a escolha das variantes e a leitura dos registros funcionam nas duas

### Requirement: Confiança de certificados
O núcleo SHALL usar a mesma política do app C# (`CertificateTrust`): o primeiro certificado de um servidor é mostrado
com a impressão digital SHA-256 e confiado só após aceite do usuário; um certificado diferente do salvo MUST bloquear a
conexão até novo aceite; certificados dos membros de um pool SHALL ser confiados quando o coordenador autenticado os
informa. Nenhum caminho MUST usar `InsecureSkipVerify` sem verificação própria da impressão digital.

#### Scenario: Primeiro acesso
- **WHEN** o usuário conecta a um servidor cujo certificado não está salvo
- **THEN** a conexão pausa, a UI mostra o assunto, o emissor e o SHA-256, e só continua se o usuário aceitar

#### Scenario: Certificado trocado
- **WHEN** o certificado apresentado difere do salvo para aquele servidor
- **THEN** a conexão é bloqueada e a UI avisa que o certificado mudou

### Requirement: Cache ao vivo por eventos
O núcleo SHALL manter por conexão um cache dos objetos (pool, host, VM, SR, VDI, VBD, VIF, network, PIF, PBD, message,
task, metrics) carregado com `event.from` e atualizado continuamente, decodificando os snapshots com os
decodificadores do SDK, e SHALL publicar mudanças a assinantes sem bloquear o laço de eventos. Ao perder a conexão ou
receber `EVENTS_LOST` MUST reconectar, recarregar tudo e publicar só as diferenças.

#### Scenario: VM ligada em outro cliente
- **WHEN** uma VM muda de `Halted` para `Running` por outro cliente
- **THEN** o cache atualiza o objeto e os assinantes recebem a mudança em poucos segundos

#### Scenario: Eventos perdidos
- **WHEN** `event.from` responde `EVENTS_LOST` ou a conexão cai e volta
- **THEN** o cache recarrega os registros e os assinantes recebem só o que mudou, incluindo objetos removidos

#### Scenario: Valor de enum desconhecido
- **WHEN** o servidor (XCP-ng 8.3 ou mais novo) devolve um valor de enum que o SDK não conhece
- **THEN** o objeto é carregado com o valor `Unrecognized` e o cache continua funcionando

### Requirement: Regras portadas do app C#
As regras de habilitação, confirmação e validação SHALL ser portadas de `XenAdmin.Desktop.Core` (`VmOperations`,
`VmRules`, `HostRules`, `PoolRules`, `MigrateRules`, `SrRules`, `StorageRules`, `NewVmDefaults`, `NotificationRules`,
`PerformanceRules`, `ConsoleRules`) com os mesmos textos de `Messages` em inglês e testes de tabela equivalentes aos
do C#.

#### Scenario: Operação não permitida
- **WHEN** uma VM não lista `clean_shutdown` em `allowed_operations`
- **THEN** a regra de Shut Down devolve desabilitado e o motivo

#### Scenario: Confirmação
- **WHEN** a UI pede Shut Down, Reboot ou Suspend de uma VM
- **THEN** a regra devolve o texto de confirmação do app C#; Start e Resume não pedem confirmação

### Requirement: Operações longas
Operações que geram `task` no XAPI SHALL ser iniciadas de forma assíncrona e acompanhadas (progresso, sucesso, erro
traduzido do código XAPI, cancelamento quando o XAPI permitir), e MUST ficar visíveis na lista de Events.

#### Scenario: Migração
- **WHEN** o usuário migra uma VM entre hosts do pool
- **THEN** a operação aparece em Events com progresso e termina com sucesso ou com a mensagem do erro XAPI

### Requirement: Console VNC por WebSocket
O núcleo SHALL abrir o console RFB da VM numa conexão TLS própria (mesma política de certificados do app) com
`CONNECT` em `console.location` e o session id, preservando os bytes já lidos após a resposta 200, e SHALL entregá-lo
como WebSocket binário a um cliente noVNC. O túnel MUST fechar quando o cliente fecha ou a VM desliga.

#### Scenario: Abrir console
- **WHEN** o usuário abre a aba Console de uma VM ligada
- **THEN** o noVNC conecta pelo endpoint com o token e mostra a tela da VM

### Requirement: Endpoint do console protegido
O endpoint WebSocket do console MUST aceitar só origens do webview da plataforma, e cada token (256 bits, validade de
30 s) MUST valer para uma VM de um `Workspace`, ser consumido de forma atômica uma única vez e ser revogado quando a
conexão ao servidor fecha.

#### Scenario: Token reutilizado
- **WHEN** um segundo cliente tenta usar um token já usado ou expirado, ou dois clientes usam o mesmo token ao mesmo tempo
- **THEN** o endpoint aceita no máximo um e recusa os outros

#### Scenario: Origem hostil
- **WHEN** uma conexão chega sem `Origin` ou com uma origem que não é a do webview
- **THEN** o endpoint recusa antes de abrir o túnel

#### Scenario: Webviews reais
- **WHEN** o teste de fumaça de cada plataforma abre o console de uma VM do lab
- **THEN** o noVNC conecta no WebKitGTK, no WKWebView e no WebView2

### Requirement: Segredos nunca expostos
Senhas, `session_id`, cookies e tokens do console MUST NOT aparecer em logs, mensagens de erro, respostas da CLI, do
MCP ou dos serviços do Wails, nem em URLs registradas; o núcleo SHALL passar todo texto que sai por um sanitizador
comum. A única exceção SHALL ser a resposta de abertura do console ao frontend do desktop, que entrega o token ao
noVNC.

#### Scenario: Segredo canário
- **WHEN** os testes conectam com uma senha e um session id canários e forçam erros de login, de HTTP e de console
- **THEN** nenhum log, erro ou resposta de qualquer transporte contém os canários

#### Scenario: Token entregue só ao console
- **WHEN** o frontend do desktop abre o console de uma VM
- **THEN** só a resposta de abertura do console contém o token, e ele não aparece em logs, erros, CLI nem MCP

### Requirement: Operações compostas portadas das Actions
Operações de vários passos (New VM, criação e reattach de SR, evacuação, Copy VM, Create Pool) SHALL seguir a
orquestração das Actions do XenModel que o app C# usa, incluindo a limpeza em falha parcial. Uma operação cujo
resultado remoto é incerto (timeout) MUST NOT ser repetida automaticamente, e operações sobre o mesmo LUN MUST ser
serializadas.

#### Scenario: Falha parcial no New VM
- **WHEN** a criação do disco falha depois de a VM ser criada
- **THEN** o núcleo desfaz o que a Action do C# desfaz e reporta o erro original

#### Scenario: Timeout
- **WHEN** uma operação de escrita estoura o timeout sem resposta do servidor
- **THEN** o núcleo reporta resultado incerto, consulta o estado real e não repete a chamada

#### Scenario: Mesmo LUN
- **WHEN** duas criações de SR iSCSI pedem o mesmo LUN ao mesmo tempo
- **THEN** a segunda falha com LUN em uso antes de chamar o XAPI
