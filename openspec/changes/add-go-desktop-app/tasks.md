## 1. Esqueleto e ferramentas

- [x] 1.1 `desktop-go/` com `go.mod` (`xcpng-center`, Go 1.27), `main.go` → `cmd.Execute()`, `cmd/root.go` (`--json`, `--config`), `cmd/version.go`, `internal/buildinfo`
- [~] 1.2 (feito: SDK 26.16.1, `PATCHES.md`, `tools/sdk-update.sh`, patches `HTTPClient`, sem `log.Fatal`, `APIError`, `Session.Ref`; falta `DecodeRecord` e o teste dedicado de que o SDK só usa o cliente do app) `sdk/xenapi` com o SDK 26.16.1, `PATCHES.md`, script `tools/sdk-update.sh` e os patches `HTTPClient`, erros sem `log.Fatal`, `APIError` tipado e `DecodeRecord`; teste que falha se o SDK abrir conexão sem o cliente do app
- [x] 1.3 `internal/gui` com Wails v3.1.0 e `cmd/gui.go` (root sem subcomando abre a janela vazia carregando a SPA); tag `nogui` para testar núcleo e CLI sem CGO
- [~] 1.4 (feito: Vue 3 + TS + Tailwind 4 + Vite → `internal/gui/dist`, `--mode web` → `web/dist`, `src/api/workspace.ts`; falta a interface comum e `http.ts`) `frontend/` Vue 3 + TS + Tailwind 4 + Vite com `--mode desktop` (→ `frontend/dist`) e `--mode web` (→ `web/dist`); `src/api/` com interface, `wails.ts` e `http.ts` vazio
- [~] 1.5 (feito: `test`, `check`, `cli`, `frontend` (node:24), `linux` (Ubuntu 24.04 + GTK4/WebKitGTK 6.0, `build/linux.Dockerfile`), alvos `go-*` na raiz; `windows` (WebView2, sem CGO) e `macos` (`.app` arm64 + amd64 com `-tags purego`, sem teste num Mac); falta `appimage`, `web-check`, `smoke`) Makefile (`test`, `check`, `frontend`, `linux`, `windows`, `macos`, `appimage`, `dist`, `web-check`, `smoke`) em container Ubuntu 24.04 com Go 1.27, Node 24 e GTK4/WebKitGTK 6.0, rodando os checks nos dois módulos (raiz e `sdk/xenapi`); alvos `go-*` no Makefile raiz
- [~] 1.6 (feito: README, `docs/architecture.md`, `docs/development.md`, `AGENTS.md`, `CHANGELOG.md`; falta o workflow, que pede o escopo `workflow` no token do GitHub) Workflow `desktop-go.yml` (PR: check, testes, web-check, fumaça Xvfb, build Windows); README, `docs/architecture.md`, `docs/development.md`, `AGENTS.md` e `CHANGELOG.md` do `desktop-go/`
- [x] 1.7 Teste de arquitetura: `internal/{xapi,model,rules,ops,console,app}` sem Wails, Cobra, MCP ou `internal/server`

## 2. Núcleo de conexão

- [~] 2.1 (feito: endereço, TLS 1.2+, sem redirects, login, `HOST_IS_SLAVE`, logout, com testes; falta a tabela de variantes por versão) `internal/xapi`: `ServerAddress` (porta de `ServerAddress.cs`, recusa 80/8080, IPv6), cliente HTTP com TLS 1.2+ e sem redirects (testes 302/307/308), login, `HOST_IS_SLAVE` → coordenador após confiar no certificado, logout, tabela de variantes por versão com fixtures 8.2/8.3
- [ ] 2.1b Sanitizador de segredos (senha, `session_id`, cookies, tokens) em logs, erros e respostas; testes com canários em todos os transportes e teste de que o token do console só sai na resposta de abertura ao frontend
- [x] 2.2 `internal/trust`: porta de `CertificateTrust` (TOFU, certificado trocado, membros pelo coordenador, `XCP_TRUST_FINGERPRINT`), store em `config`
- [x] 2.3 `internal/config`: `settings.toml` em `os.UserConfigDir()/xcp-ng-center-go`, permissões 0600, lista de servidores, preferências de UI
- [ ] 2.4 `internal/model`: cache por conexão com `event.from` e `DecodeRecord`, `Workspace`, assinaturas, reconexão e recarga em `EVENTS_LOST` com diffs (inclusive remoções); teste com 1000 VMs no servidor falso
- [ ] 2.5 Servidor JSON-RPC falso (`internal/xapitest`) para testes do núcleo
- [ ] 2.6 `internal/app`: casos de uso de conexão e árvore com `Workspace` pelo contexto; teste de dois workspaces isolados
- [x] 2.7 CLI `server add/list/connect`, `vm list` (tabela e `--json`); servidor e usuário do ambiente ou das flags, senha de `XCP_PASS` ou prompt sem eco
- [ ] 2.8 Lab KVM: conectar ao pool, TLS real (recusa sem aceite, aceite, troca de certificado), membro → coordenador

## 3. Shell do desktop e tema

- [~] 3.1 (feito: tokens e regra base sem cantos em `style.css`; falta claro/escuro, densidade e o teste Playwright) Tema Tailwind (tokens, claro/escuro, densidade, `--radius-*: 0`, regra base) e teste Playwright de cantos
- [ ] 3.1b Full HD, QHD e 4K: tokens em `rem`, zoom 80/90/100/110/125/150% (`Window.SetZoom`, Ctrl+=/-/0), janela inicial em 80% da área de trabalho com estado restaurado, ícones `_32`/`_24` em `devicePixelRatio` ≥ 1,5; capturas Playwright em 1920×1080@1x, 2560×1440@1x/1,25x e 3840×2160@2x
- [ ] 3.2 Componentes base: Splitter, Tree, Tabs, DataGrid virtualizado, MenuBar, Toolbar, ContextMenu, Dialog, Wizard, StatusBar, Field, Button
- [~] 3.2b (feito: cópia dos ícones da árvore no `make frontend`; falta `IconNames` completo) Ícones: cópia de `XenAdmin/Images` e `branding-xcp-ng/Images` no build e porta de `IconNames` (objeto, estado, toolbar, menus) com teste de tabela
- [~] 3.3 (feito: árvore pool → hosts → VMs → SRs com ícones do C#, aba General, Add New Server, diálogo de certificado, barra de status, testado no .110; faltam menus e toolbar completos) Janela principal: árvore (ícones do app C#), abas, menus e barra de ferramentas, Add Server e diálogo de certificado, barra de status
- [ ] 3.4 Abas General de pool, host, VM e SR (porta de `Details` e `OverviewColumns`)
- [ ] 3.4b Segurança do webview: CSP, bloqueio de navegação externa, permissões negadas, lint sem `v-html`, teste com nome de VM contendo HTML
- [ ] 3.5 Serviços Wails e eventos `xcp:*` coalescidos; instância única; estado da janela; menu nativo no macOS
- [ ] 3.6 Paridade: capturas Go × C# das telas do marco em `_output/parity/`; `docs/parity.md`

## 4. Operações

- [~] 4.1 (feito: `rules/vm.go` e `rules/host.go` com testes de tabela, `tools/resx2go` e `internal/messages` gerado de Messages.resx e FriendlyNames.resx; falta `PoolRules`) `internal/rules`: `VmOperations`, `VmRules`, `HostRules`, `PoolRules` com testes de tabela do C#; `internal/messages` gerado por `tools/resx2go`
- [~] 4.2 (feito: tasks assíncronas com polling até o fim (`ops.wait`) e erros XAPI tipados; falta progresso, cancelamento e o resultado incerto) `internal/ops`: Async + task, progresso, erros XAPI traduzidos, cancelamento quando permitido, resultado incerto sem repetição
- [ ] 4.2b Orquestração portada das Actions (Create VM, SR create/reattach/introduce, Evacuate, Copy, Create Pool) com testes de falha em cada passo e reserva de LUN
- [x] 4.3 VM: Start, Shut Down, Reboot, Suspend, Resume e forçados, com as confirmações do C#; menus e toolbar
- [x] 4.4 Host: Shut Down, Reboot, Enter/Exit Maintenance Mode com o diálogo de evacuação e a volta das VMs
- [~] 4.5 (feito: Alerts com os textos de `NotificationRules` (login falho, nomes amigáveis, prioridade), Dismiss, Dismiss Selected e Dismiss All com as confirmações do C#, painel Infrastructure/Notifications; faltam os filtros, o contador na barra e Events) Notifications: Alerts (filtros por severidade e servidor, Dismiss, Dismiss All/Selected, contador) e Events
- [~] 4.6 (feito: pool lab-1 + lab-2 com SR NFS e VM Alpine: Enter Maintenance Mode migrou a VM, Exit com Restore a trouxe de volta, Reboot desligou a VM e reiniciou o host, MCP `host_maintenance`; faltam alertas no lab) Lab KVM: ciclo de vida de VM, manutenção com VMs, alertas

## 5. Console

- [x] 5.1 (origens tiradas do assetserver do Wails 3.1: `wails://localhost` no Linux/macOS, `http://wails.localhost` no Windows; conferida no WebKitGTK) `internal/console`: `tls.Dial` com a política do app, `CONNECT` com bytes bufferizados preservados e deadlines; listener loopback com token de uso único (atômico, ligado à VM e ao `Workspace`, revogado no disconnect) e origens por plataforma (`wails://localhost`, `http://wails.localhost`) conferidas nos webviews reais; testes com origem ausente/hostil e consumo concorrente
- [~] 5.2 (feito: noVNC 1.5 com escala, Send Ctrl+Alt+Del, Reconnect, console do host pelo dom0; faltam scancodes ABNT2 e a faixa de DVD) Aba Console com noVNC (escala, Ctrl+Alt+Del, teclado ABNT2/scancodes), faixa de DVD; spike com streams do Wails
- [ ] 5.3 Launch SSH Console e Remote Desktop (porta de `ConsoleRules` e do lançador do app C#)
- [~] 5.4 (feito: VM Alpine no membro do pool e console do host no WebKitGTK, login e símbolos `|@#$%&*()_+{}:<>?~^[]=;` corretos; faltam WKWebView e WebView2) Lab KVM: console de VM Alpine e teclado com símbolos no WebKitGTK, WKWebView (Mac do lab) e WebView2 (Windows do lab)

## 6. Abas e Properties

- [~] 6.1 (feito: tabela de discos do `Details.Storage`, Activate/Deactivate, drive de DVD (criar, ISO, Eject) com o aviso de reinício do C#, testado no lab KVM; faltam novo, anexar, mover, desanexar, excluir disco e Properties) Storage da VM (novo, anexar, ativar, desativar, desanexar, excluir disco; DVD e ISO) — `StorageRules`
- [~] 6.2 (feito: Networking (tabela do `Details.Networking`) e Memory (valores), sem edição ainda; Snapshots com Take Snapshot (disco ou disco e memória, com o motivo quando não disponível), Revert To com "tirar snapshot antes" e volta do estado de energia, Delete só de snapshots com os discos próprios, testado no lab KVM e pelo MCP; faltam Networking, Memory, Search) Snapshots (tirar, reverter, excluir), Networking, Memory, Search
- [ ] 6.3 Performance (`internal/rrd`, `PerformanceRules`, uPlot, zoom de 10 min a 1 ano, eventos de ciclo de vida)
- [ ] 6.4 Properties: General, CPU, Boot Options, Home Server
- [ ] 6.5 SR: Set as Default, Detach, Forget, Destroy com as regras e confirmações do C#
- [ ] 6.6 Paridade e lab das telas do marco, incluindo os comportamentos do 0.2.0 (duplo clique, aba por objeto, console que sobrevive à troca de aba, full copy, SRs ocultos)

## 7. Wizards

- [ ] 7.1 New VM (`NewVmDefaults`: templates, instalação por ISO, discos com Add/Delete, rede, CPU/memória, UEFI preferido)
- [ ] 7.2 New Storage Repository: NFS, SMB, iSCSI (Scan, IQN, LUN, Format/Reattach), NFS ISO, SMB ISO (`SrRules`)
- [ ] 7.3 Create New Pool (`PoolRules`), Copy VM (fast clone/full copy), Delete VM, Migrate to Server (`MigrateRules`)
- [ ] 7.4 Lab KVM: instalação do Alpine por ISO com `xe-guest-utilities`, iSCSI no target do lab, NFS ISO, pool, migração

## 8. MCP e CLI

- [~] 8.1 (feito: `list_servers`, `connect`, `disconnect`, `list_objects`, `list_alerts`; com `--allow-write` e `confirm: true`: `vm_power`, `host_power`, `host_maintenance`, `dismiss_alerts`, testados no .110 e no lab KVM; faltam `get_object`, snapshots, `vm_migrate`, `sr_create`, `get_performance` e `operation_id`) `internal/mcp` com `modelcontextprotocol/go-sdk`: mesmas ferramentas do MCP C#, somente leitura por padrão, `--allow-write` + `confirm: true`, `operation_id`, redação de senhas, textos do servidor como dados
- [ ] 8.2 Teste de paridade de nomes e esquemas com o MCP C#; teste do processo real por stdio no lab
- [ ] 8.3 CLI `vm`, `host` e `sr` com as operações e confirmações (`--yes` explícito para escrita)
- [ ] 8.3b Flags globais (`--host`, `--user`, `--trust-fingerprint`, `--json`, `--config`), `mcp` sem display no Linux e pelo AppImage, `AttachConsole` no exe do Windows; testes de terminal e de stdin/stdout redirecionados
- [ ] 8.4 `docs/mcp.md` e `docs/cli.md`

## 9. Empacotamento, release e web

- [ ] 9.1 Linux: binário, tarball e AppImage (base Ubuntu 24.04) com fumaça Xvfb em Ubuntu 24.04, Debian 13 e Fedora 42; distribuições suportadas no README
- [ ] 9.2 macOS universal (`lipo`, `.app`, mínimo 14) no runner `macos-15`; validação no Mac do lab
- [ ] 9.3 Windows x64 (exe em zip); validação no Windows do lab
- [ ] 9.4 Workflow de release `go-v*` (pre-release, `SHA256SUMS`, notas no estilo rustXterm) e seção Go na skill `create-release`
- [ ] 9.5 Spike do modo web (build `-tags server` do Wails com middleware de autenticação e `Workspace` por usuário) documentado em `docs/web.md`, sem publicar o comando `serve`
- [ ] 9.6 README raiz, AGENTS.md raiz e revisão Codex final
