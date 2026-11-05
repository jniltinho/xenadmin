## Context

O app desktop C# (`XenAdmin.Desktop`, Avalonia, release `linux-v0.2.0`) já tem a lógica separada da UI
(`XenAdmin.Desktop.Core`) e um servidor MCP, mas depende do XenModel do WinForms e do .NET. Esta mudança cria, em
`desktop-go/`, uma segunda implementação em Go com a mesma lógica, a interface em Vue + TailwindCSS dentro do Wails
v3.1.0, e preparada para virar uma versão web.

Referências usadas:
- **Padrão de projeto do autor**: `ad-manager` (`main.go` → `cmd.Execute()`, Cobra com um arquivo por grupo e `init()`
  com `rootCmd.AddCommand`, `internal/<pkg>` com godoc, `frontend/` Vue 3 + Tailwind 4 com build `--mode desktop` e
  `--mode web`, `web/embed.go`, `desktop/` como módulo Go separado com Wails v3, MCP com
  `modelcontextprotocol/go-sdk`, config TOML, Makefile com alvos em container, Keep a Changelog) e `ova2xva` (SDK
  `xenapi` vendorizado com `replace`, upload `/import` acompanhado por `Task`).
- **XenServer SDK para Go 26.16.1** (`XenServer-SDK-26.16.01.zip`, `XenServerGo/src`, BSD-2, API 2.23). O
  `sdk/xenapi` do ova2xva é a mesma geração em versão anterior (API 2.21), sem patches.
- **Wails v3.1.0** (tag de 2026-11-01): Services com bindings gerados por análise estática, eventos tipados, várias
  janelas, bandeja, menus, diálogos, notificações, instância única, Taskfile, AppImage/deb/rpm, `.app` universal com
  `lipo`, NSIS/MSIX, e build `-tags server` (o mesmo app como servidor HTTP, sem janela e sem autenticação embutida).
  Linux usa GTK4 + WebKitGTK 6.0 por padrão.
- **XenServer 9 removed features** e a tabela do XenServer 8, para não portar o que saiu.
- **XCP-ng 8.3**: XAPI 26.1.x (próximo do SDK 26.16), sem licenciamento Citrix, CDN ou Health Check; atualizações por
  `host.call_plugin` (`updater.py`).

## Goals / Non-Goals

- Goals:
  - Paridade funcional e de UI/UX com o app C# 0.2.0, comparando cada tela com ele e com o XCP-ng Center v25.04.
  - Layout moderno, claro/escuro, sem cantos arredondados.
  - Binário único para desktop + CLI + MCP em Linux x64, macOS universal e Windows x64.
  - Núcleo testável sem UI e sem rede (servidor XAPI falso), mais testes no lab KVM.
  - Estrutura que permita servir a mesma SPA e o mesmo núcleo por HTTP sem reescrever regras.
- Non-Goals (nesta mudança):
  - Servidor web multiusuário (só a preparação; mudança própria depois).
  - Atualizações de hosts, HA/DR, vApps, VM groups, GPU/vGPU, PVS, WLB, Conversion Manager e o resto da 8.4g do app C#.
    Entram em mudanças seguintes, na ordem em que o app C# as tiver.
  - Assinatura e notarização (macOS Developer ID, Authenticode).
  - Import/Export de VMs (XVA e OVA via ova2xva).
  - Substituir o app C#: os dois convivem até a paridade.

## Decisions

### D1. Layout de `desktop-go/`

```
desktop-go/
  go.mod                 module xcpng-center (Go 1.27), Wails v3.1.0, replace xenapi => ./sdk/xenapi
  main.go                cmd.Execute(): um binário para janela, CLI e MCP
  cmd/                   Cobra: root.go (sem subcomando abre a janela), gui.go, version.go, mcp.go, server.go,
                         vm.go, host.go, sr.go (+ _test.go)
  internal/
    buildinfo/           versão, commit e data por -ldflags
    config/              settings.toml em os.UserConfigDir()/xcp-ng-center-go (0600; recusa arquivo legível por outros)
    trust/               confiança de certificados (porta de CertificateTrust)
    xapi/                cliente HTTP com TLS do app, sessão, login, coordenador, ServerAddress, endpoints HTTP do XAPI
    model/               cache por conexão alimentado por event.from, assinaturas, Workspace
    rules/               regras portadas de XenAdmin.Desktop.Core (um arquivo por classe C#, testes de tabela)
    ops/                 operações longas: Async* + task, progresso, erros XAPI traduzidos, Events
    console/             túnel RFB (HTTP CONNECT em console.location) e ponte WebSocket com token de uso único
    rrd/                 rrd_updates e archives para Performance
    messages/            textos de XenModel/Messages.resx usados pelas regras (gerado por tools/resx2go)
    app/                 casos de uso neutros de transporte (D4)
    mcp/                 ferramentas MCP sobre internal/app (github.com/modelcontextprotocol/go-sdk)
    gui/                 app Wails: janela, serviços finos sobre internal/app, nativos (SSH, RDP, janela, bandeja)
  sdk/xenapi/            XenServer SDK 26.16.1 + PATCHES.md (D3)
  frontend/              Vue 3 + TypeScript + TailwindCSS 4 + Vite (D6); dist embutido por internal/gui
  web/                   embed.go (//go:embed all:dist) reservado para o modo web
  build/                 Taskfile e configs do Wails (ícones, Info.plist, NSIS, AppImage), layout do wails3 init
  packaging/             .desktop, ícones, Info.plist, launcher
  docs/                  architecture.md, development.md, parity.md, mcp.md, cli.md
  tools/resx2go/         gera internal/messages a partir dos .resx
  Makefile, README.md, CHANGELOG.md, AGENTS.md
```

Um binário e um pacote por plataforma, como o app C# (`xcp-ng-center mcp`): `xcp-ng-center-go` sem subcomando abre a
janela; `mcp`, `vm`, `host`, `sr`, `server` e `version` rodam sem janela, e só `cmd/gui.go` importa `internal/gui`
(Wails). São dois módulos Go (raiz e `sdk/xenapi`); `make check` e `make test` rodam nos dois, porque `go test ./...`
não entra em módulos aninhados. O núcleo e a CLI são testados sem CGO (`go test ./internal/... ./cmd/...` com a tag
`nogui`, que troca `cmd/gui.go` por um stub).

Diferente do ad-manager, cujo `desktop/` é um módulo separado porque só faz proxy para o painel web remoto, aqui o
desktop roda o núcleo localmente, então janela, CLI e MCP ficam no mesmo módulo e no mesmo binário.

Flags (Cobra, persistentes no root): `--config`, `--json`, `--host` (repetível; sobrepõe `XCP_HOST`), `--user`
(sobrepõe `XCP_USER`) e `--trust-fingerprint` (sobrepõe `XCP_TRUST_FINGERPRINT`); `mcp` acrescenta `--allow-write`
(padrão somente leitura). A senha nunca é flag: vem de `XCP_PASS` ou, na CLI interativa, de um prompt sem eco.

O binário liga GTK4/WebKitGTK no Linux mesmo nos modos sem janela, então MCP e CLI precisam dessas bibliotecas; numa
máquina sem elas, o próprio AppImage roda `xcp-ng-center-go-<ver>-x86_64.AppImage mcp` com as bibliotecas que traz
(`--appimage-extract-and-run` quando não houver FUSE). No Windows o exe é `-H windowsgui` e, com subcomando, chama
`AttachConsole(ATTACH_PARENT_PROCESS)` para escrever no terminal; com stdin/stdout redirecionados (MCP) usa os handles
herdados. Testes cobrem os dois casos. Durante o preview o nome termina em `-go` para conviver
com o `xcp-ng-center` do app C# (mesmas configs de MCP, PATH, AppImage); o nome final é decidido na paridade.

### D2. Núcleo portado, não reinventado
Cada arquivo de `internal/rules` espelha uma classe de `XenAdmin.Desktop.Core` (`VmOperations`, `VmRules`, `HostRules`,
`PoolRules`, `MigrateRules`, `SrRules`, `StorageRules`, `NewVmDefaults`, `NotificationRules`, `PerformanceRules`,
`ConsoleRules`, `Details`, `OverviewColumns`) e cita o arquivo de origem e o do WinForms/XenModel que ele segue. Os
testes de tabela do C# (`tests/XenAdmin.Desktop.Tests`) são portados caso a caso. Os textos vêm de `Messages.resx`
por geração (`tools/resx2go` com lista de chaves usadas), não copiados à mão.

### D3. SDK XAPI vendorizado com dois patches
O SDK 26.16.1 é copiado de `XenServerGo/src` para `sdk/xenapi` (módulo `xenapi`, como no ova2xva). O cliente gerado:
desliga a verificação TLS quando não recebe um arquivo de CA (`InsecureSkipVerify: true`), não aceita `http.Client`
nem `context.Context`, fixa cipher suites ECDHE-RSA e chama `log.Fatal` em erros de URL/certificado. Patches, listados
em `sdk/xenapi/PATCHES.md` e reaplicáveis por script a cada SDK novo:
1. `ClientOpts.HTTPClient *http.Client`: quando presente, é usado como está (TLS do app, timeouts, sem redirects).
2. Erros de construção devolvidos (`NewSessionE` ou equivalente) em vez de `log.Fatal`.
3. Erros XAPI tipados: o cliente achata `error.data` numa string; o patch devolve `*APIError{Code, Params}` com
   `errors.As`, para as regras traduzirem pelo `Messages` como o C# faz com `Failure`.
4. Decodificadores de registro exportados (`DecodeRecord(class, snapshot)`): `EventRecord.Snapshot` chega como mapa
   bruto, e o cache precisa dos mesmos conversores de enum/número/data que o `get_all_records` usa.

O TLS do app (`internal/xapi`) aceita TLS 1.2+ com as suítes padrão do Go e verifica a impressão digital em
`VerifyConnection` (D5). Sem `context.Context` no SDK, cada chamada usa o timeout do cliente; `event.from` usa um
cliente próprio com timeout maior que o long-poll. Endpoints HTTP fora do SDK (`/rrd_updates`, `/host_rrd`, `/vm_rrd` e,
depois, `/import` e `/export`) usam o mesmo `http.Client` com a sessão. O `/console` usa uma conexão TLS dedicada que
compartilha só a política de TLS (D8).

Valores de enum desconhecidos viram `Unrecognized` e campos ausentes são ignorados, o que mantém a compatibilidade com
XCP-ng 8.2/8.3. Métodos com versão (`Host.Evacuate3`...) são escolhidos por uma tabela de capacidades a partir de
`session.XAPIVersion` e `APIVersion`; versão desconhecida (`APIVersionUnknown`) usa a variante mais antiga. Fixtures
gravadas de XCP-ng 8.2 e 8.3 (respostas JSON-RPC reais do lab, sem segredos) alimentam os testes.

### D4. Camada de serviço neutra de transporte (`internal/app`) e `Workspace`
`internal/app` expõe os casos de uso (conectar, árvore, detalhes, operações de VM/host/SR, wizards, alertas,
performance, console) com `context.Context`, DTOs JSON e erros tipados (`code`, `message`, `details`). O estado do
usuário fica num `Workspace` (conexões, cache, confiança, operações) obtido do contexto. Adaptadores finos, sem regras:
- **Desktop**: serviços Wails em `internal/gui` (bindings TS gerados por `wails3 generate bindings`); mudanças do
  cache e das operações chegam como eventos Wails tipados (`xcp:objects`, `xcp:operation`, `xcp:connection`), com
  coalescência a cada 100 ms para não inundar o webview.
- **MCP**: `internal/mcp` registra as ferramentas sobre os mesmos casos de uso.
- **CLI**: `cmd/` chama os mesmos casos de uso, saída tabela ou `--json`.
- **Web (futuro)**: o mesmo `internal/app` atrás de HTTP + WebSocket/SSE (D9).

### D5. Conexão e confiança de certificados
Porta direta de `ServerAddress` (IPv4, IPv6 puro, `[v6]:porta`; recusa 80/8080) e de `CertificateTrust`: TOFU com
diálogo (assunto, emissor, validade, SHA-256), bloqueio quando o certificado muda, membros do pool confiados pelo
coordenador autenticado (`host.get_server_certificate` / `certificate` records), `XCP_TRUST_FINGERPRINT` para MCP/CLI.
`HOST_IS_SLAVE` refaz o login no coordenador só depois de o certificado dele passar pela mesma política. O cliente
HTTP usa `CheckRedirect` que recusa todo redirecionamento (o XAPI não precisa deles); testes com 302/307/308 para
`http://` e para outro host garantem que a senha não sai (a mesma regressão corrigida no C# 0.2.0). As senhas ficam só
em memória (como o app C# hoje); cofre do sistema é mudança futura.

Segredos (senha, `session_id`, cookies, tokens do console) passam por um sanitizador único em `internal/xapi` antes de
qualquer log, erro ou resposta (CLI, MCP, serviços Wails); URLs registradas perdem a query. A única exceção é a
resposta de `OpenConsole` ao frontend do desktop, que precisa entregar a URL loopback com o token ao noVNC; o token
nunca sai por CLI ou MCP nem aparece em logs. Testes com valores
canário forçam erros em todos os caminhos e procuram os canários na saída.

### D6. Frontend: Vue 3 + TypeScript + TailwindCSS 4 + Vite
Mesmo stack do `ad-manager/frontend` (Vue 3.5, vue-router, Tailwind 4 com `@tailwindcss/vite`, Vite; Vitest e
Playwright). O ad-manager usa Vite+ (`vp`); aqui fica o Vite padrão até o Vite+ ser necessário, para reduzir o
número de ferramentas. Estrutura:

```
frontend/src/
  api/          index.ts (interface), wails.ts (bindings + eventos), http.ts (fetch + WebSocket/SSE, modo web)
  stores/       estado reativo (Pinia ou composables) a partir dos eventos de objetos
  components/   Tree, Tabs, DataGrid (virtualizado), MenuBar, Toolbar, ContextMenu, Dialog, Wizard, PropertyPages,
                StatusBar, Splitter, Field, Button
  views/        General, Memory, Storage, Networking, Console, Performance, Snapshots, Search, Notifications
  wizards/      NewVm, NewSr, NewPool, CopyVm, Migrate
  dialogs/      AddServer, CertificateTrust, Properties, Confirm, Evacuate
  desktop/      código só do desktop (fica fora do build web)
  theme/        tokens, claro/escuro, ícones
```

Sem biblioteca de componentes pronta (elas trazem cantos arredondados e um visual web); componentes próprios, poucos e
testados. A barra de menus é desenhada na janela (como no app C#, e funciona no modo web); no macOS o menu nativo
também tem os itens padrão do app. Todos os ícones são os mesmos do app desktop C#: os PNG de `XenAdmin/Images`
(BSD-2) e o logo/ícone de `branding-xcp-ng/Images`, copiados no build para `frontend/public/images/` (sem duplicar no
git) e escolhidos pelas mesmas regras de `IconNames` (objeto, estado, alertas, barra de ferramentas e menus), para que
árvore, abas e botões fiquem iguais ao app C#; nenhum outro conjunto de ícones é usado. Gráficos de Performance com uPlot. Console com noVNC (MPL-2.0, sem modificações, arquivo de licença no
pacote).

### D7. Tema: moderno e sem cantos arredondados
- Tokens em `@theme`: todas as `--radius-*` = `0`; regra base `*, ::before, ::after { border-radius: 0 }` para cobrir
  componentes de terceiros (noVNC, uPlot); foco por `outline` quadrado de 2 px.
- Densidade de desktop (linhas de 24–28 px na árvore e nas grades), fonte do sistema (`system-ui`), cor de destaque
  única, superfícies planas com bordas de 1 px em vez de sombras grandes; ícones de 16 px do app C# (os PNG de 32 px
  onde o app C# os usa), renderizados sem suavização extra para não borrar.
- Claro e escuro pelo `prefers-color-scheme` com escolha manual em Settings; contraste WCAG AA.
- Teste Playwright percorre as telas e diálogos e falha se algum `getComputedStyle(el).borderRadius` não for `0px`.
- Os cantos da própria janela pertencem ao sistema (macOS e Windows 11 arredondam a moldura); janelas sem moldura não
  são usadas para não perder o comportamento nativo.

### D8. Console VNC
`internal/console` abre uma conexão TLS própria (`tls.Dial` com o mesmo `tls.Config` e verificação de impressão
digital do app; o `http.Client` não entrega o fluxo bidirecional de um `CONNECT`), escreve
`CONNECT <path>?session_id=… HTTP/1.0`, lê a resposta com um `bufio.Reader` e repassa ao WebSocket também os bytes já
bufferizados depois do `200`; deadlines de conexão e de inatividade. O proxy de assets do Wails não carrega WebSocket (lição do ad-manager), então no desktop o núcleo
abre um listener só em `127.0.0.1` porta aleatória, aceita apenas `/console/{token}` com token aleatório de 256 bits,
uso único (consumo atômico) e validade de 30 s, ligado à VM e ao `Workspace` e revogado quando a conexão ao servidor
fecha, e aceita só o `Origin` do webview de cada plataforma (`wails://localhost` no Linux/macOS e
`http://wails.localhost` no Windows, a URL base do asset server do Wails v3.1.0; o `Origin` efetivamente enviado por
cada webview real é conferido no teste de fumaça), recusando `Origin` ausente. noVNC conecta nesse endereço. No modo web o mesmo
handler fica atrás do servidor autenticado. Alternativa a investigar: os streams do Wails v3 como canal do noVNC, que
evitam a porta local; adotada se o spike mostrar latência equivalente.

### D8b. Segurança do webview
Como no ad-manager: CSP restrita servida pelo middleware de assets (`default-src 'self'`, sem `unsafe-eval`,
`connect-src` só o app e `ws://127.0.0.1:<porta do console>`), bloqueio de navegação para fora do app (links externos
por `Browser.OpenURL` depois de validados), mapa de permissões do webview negando câmera, microfone, geolocalização e
notificações, e textos vindos do XAPI sempre interpolados como texto (lint que proíbe `v-html`). O webview é um
cliente privilegiado (chama os serviços que operam o pool), então qualquer conteúdo do servidor é dado, nunca markup.

### D8c. Operações compostas
As regras dizem se algo pode ser feito; a orquestração vem das Actions do XenModel usadas pelo app C#
(`CreateVMAction`, `SrCreateAction`/`SrReattachAction`/`SrIntroduceAction`, `EvacuateHostAction`, `VMCopyAction`,
`CreatePoolAction`...), portadas passo a passo com a limpeza em falha parcial. Escrita com resultado incerto (timeout)
nunca é repetida automaticamente: o núcleo marca a operação como incerta e relê o estado. Criações sobre o mesmo LUN
são serializadas por uma reserva no `Workspace` (como o MCP C# faz). Testes do servidor falso injetam falhas em cada
passo.

### D9. Caminho para a versão web
Duas opções, decididas numa mudança própria depois de um spike nesta:
1. **Build `-tags server` do Wails v3**: os mesmos serviços por HTTP e eventos por WebSocket. Falta autenticação e
   estado por usuário; o spike verifica se um middleware de assets consegue autenticar e passar o `Workspace` do
   usuário até o serviço pelo contexto da chamada.
2. **Servidor próprio** (`internal/server`, `net/http` ou Echo v5 como no ad-manager) com o `frontend/src/api/http.ts`.

Nos dois casos: sessão por usuário no servidor, cookie `HttpOnly`/`Secure`/`SameSite=Strict`, credenciais XAPI só no
servidor, console pelo mesmo handler WebSocket, `web/dist` embutido. Esta mudança só garante as costuras (D4, o
adaptador `api/`, `Workspace` por contexto, `make go-web-check`).

### D10. Recursos não portados (removidos ou inexistentes no XCP-ng)
Do XenServer 8/9 e do XCP-ng: Health Check e uploads ao Citrix Insight Services; licenciamento Citrix, LAS, Standard
Edition e o gerenciador de licenças; atualizações por CDN (`pool.sync_updates`, repositórios); VM Recovery Mode;
PuTTY embutido (usa o `ssh` do sistema); SSL legado, TLS 1.0/1.1; criação de guests PV (VMs PV existentes aparecem,
mas sem opções só de PV); VSS e snapshots quiesced; Container Management; Measured Boot; Transfer VM; vSwitch
Controller e redes privadas entre servidores; ligar hosts por DRAC/iLO; AMD MxGPU e Intel GVT-g; Linux bridge (só
vSwitch); FCoE por software (o tipo FCoE não entra no New Storage); hosts 7.x ou mais antigos (mínimo XCP-ng 8.2).
BIOS para VMs continua disponível (XCP-ng 8.3 suporta), mas o New VM prefere UEFI quando o template permitir, como o
XenServer 9 recomenda.

### D11. Testes
- **Unidade (Go)**: regras com testes de tabela portados do C#; `internal/xapi`, `model`, `ops` e `console` contra um
  servidor JSON-RPC falso (`httptest`) que grava as chamadas; teste de arquitetura dos imports (núcleo sem Wails,
  Cobra, MCP); `-race`.
- **Frontend**: Vitest para stores e componentes; Playwright contra a SPA com a API falsa para capturas, teste de
  cantos e acessibilidade básica.
- **Lab KVM** (as regras do AGENTS.md): conexão, TLS, cache, operações de VM e host, storage NFS/SMB/iSCSI/ISO,
  pool, migração e manutenção, com `XCP_LAB_HOSTS`, `XCP_ALLOW_WRITE`; nunca nos servidores reais, salvo leitura no
  `.236`.
- **Paridade**: o app Go roda no container Xvfb ao lado do app C#, e cada tela é capturada nos dois (e no XCP-ng
  Center v25.04 do lab Windows quando a tela vier de lá), em `_output/parity/` (nunca commitado). Além das telas, a
  lista de comportamentos do CHANGELOG 0.2.0 (duplo clique, aba por objeto, RDP, reattach/introduce, console que
  sobrevive à troca de aba, full copy, SRs ocultos) tem um teste ou roteiro de lab cada;
  `docs/parity.md` lista tela por tela o estado.
- **Fumaça**: AppImage sob Xvfb com `XCP_SMOKE=1` (abre, carrega a SPA, sai 0) nas distribuições suportadas; `.app`
  e `.exe` validados como hoje no Mac e no Windows do lab, incluindo o console noVNC nos três webviews reais (a API
  falsa do Playwright não prova o WebSocket nem as origens).

### D12. Build, CI e release
- `desktop-go/Makefile` com alvos em container (Go 1.27, Node 24, GTK4/WebKitGTK 6.0): `test`, `check` (gofmt, vet,
  staticcheck, `vp`/`vite` check), `frontend`, `linux`, `appimage`, `windows`, `macos`, `dist`, `web-check`, `smoke`;
  o `Makefile` raiz ganha `go-*` que delegam a eles.
- Linux: binário CGO com GTK4/WebKitGTK 6.0, compilado em Ubuntu 24.04 (base fixa do container); AppImage do
  `wails3`. Empacotar o WebKit não remove a dependência da glibc da base, então o suporte é Ubuntu 24.04+, Debian 13+,
  Fedora 42+ (fumaça em containers dessas três). Ubuntu 22.04, Debian 12 e RHEL 9 não são suportados nesta mudança
  (o caminho GTK3 do Wails está saindo).
- macOS: exceção à regra de builds em container. Dois builds (arm64, amd64) unidos com `lipo` num `.app` (mínimo
  macOS 14), feitos no runner `macos-15` como no ad-manager, porque o Wails precisa de CGO com o SDK da Apple; o
  backend experimental sem CGO (`-tags purego`) do 3.1 não é usado até sair do experimental. Zip com assinatura
  ad-hoc.
- Windows: `CGO_ENABLED=0`, `-H windowsgui`, ícone/manifesto, exe em zip; NSIS fica para depois.
- Workflow `desktop-go.yml`: testes e checks em PR; tag `go-v*` publica pre-release com
  `xcp-ng-center-go-<ver>-x86_64.AppImage`, `-linux-x64.tar.gz`, `_macos_universal.zip`, `-win-x64.zip` (os três
  com o mesmo binário para janela, CLI e MCP) e `SHA256SUMS`, notas de `desktop-go/CHANGELOG.md` no estilo do rustXterm; skill `create-release` ganha a seção Go.

## Risks / Trade-offs

- **SDK sem context e com TLS inseguro por padrão** → patches mínimos e documentados (D3); teste que falha se uma
  conexão for aberta sem o `HTTPClient` do app.
- **Wails v3 é recente** (3.0.0 em 2026-10-26, 3.1.0 em 2026-11-01) → versão fixa em `go.mod`, atualizações só com
  CHANGELOG lido e testes de fumaça nos três sistemas; o núcleo não depende do Wails.
- **WebKitGTK no Linux** é mais lento que o WebView2/WKWebView em grades grandes → DataGrid virtualizado, eventos
  coalescidos, medição com 1000 VMs no servidor falso.
- **Porta local do console** → só loopback, token de uso único, `Origin` conferido, fecha com a aba; spike com
  streams do Wails para eliminar a porta.
- **Dois apps para manter** (C# e Go) → as regras são portadas com citação e testes de tabela iguais; mudanças de regra
  no C# entram nos dois na mesma mudança OpenSpec até a decisão de qual segue.
- **Paridade visual × layout moderno** → paridade é de conteúdo, ordem, rótulos e fluxo; o visual (cores, densidade,
  cantos) é o novo tema. `docs/parity.md` registra cada desvio intencional.

## Migration Plan

Marcos, cada um com testes (unidade + lab), revisão Codex, commit e push (regra do AGENTS.md):
1. Esqueleto: módulos, SDK com patches, Makefile/containers, CI, `version`.
2. Núcleo de conexão: `xapi`, `trust`, `model`, `config`; CLI `server`/`vm list`.
3. Shell do desktop: Wails, SPA com árvore, abas General e barra de status, tema e teste de cantos.
4. Operações: VM lifecycle, confirmações, Events, Notifications/Alerts, host power e manutenção.
5. Console noVNC, SSH e RDP.
6. Storage, Snapshots, Networking, Memory, Search, Performance; Properties.
7. Wizards: New VM, New SR (NFS, SMB, iSCSI, ISO), New Pool, Copy, Migrate, Delete.
8. MCP com paridade de ferramentas; CLI completa.
9. Empacotamento e release `go-v0.1.0` (pre-release); spike do modo web documentado.

Rollback: a pasta é independente; nada do app C# muda, e o release `linux-v*` continua.

## Open Questions

- Nome final do binário e do produto quando o app Go alcançar a paridade (substitui o C# ou convive?).
- Import/Export (fora desta mudança): XVA direto por `/import` e `/export`; OVA pelo ova2xva, que tem a conversão em
  `internal/` (não importável por outro módulo): promover para `pkg/` no repositório do ova2xva ou chamar o binário.
  Decidido na mudança própria.
- Pinia ou só composables para o estado da SPA (decidir no marco 3 conforme o tamanho real).
- Vite+ (`vp`) como no ad-manager ou Vite padrão (D6 começa com Vite padrão).
