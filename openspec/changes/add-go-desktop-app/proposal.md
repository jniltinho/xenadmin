## Why

O app desktop multiplataforma em C# (Avalonia, `XenAdmin.Desktop`) depende do XenModel do WinForms, do .NET e de um
runtime de UI que não roda no navegador. Uma versão em Go, com a interface em Vue + TailwindCSS dentro do Wails v3,
gera binários nativos menores para Linux, macOS e Windows, reaproveita o SDK oficial do XAPI em Go (XenServer SDK
26.16.1, o mesmo modelo do `sdk/xenapi` do ova2xva) e deixa o terreno pronto para uma versão web: a mesma interface e
o mesmo núcleo servidos por HTTP. O objetivo é paridade funcional e de UI/UX com o app desktop em C# (que por sua vez
segue o XCP-ng Center v25.04 do Windows), com um layout mais moderno e sem cantos arredondados.

## What Changes

- Nova pasta `desktop-go/` com um projeto Go independente do .NET, no padrão dos projetos `ad-manager` e `ova2xva`:
  `main.go` → `cmd.Execute()` (Cobra), `internal/` com um pacote por responsabilidade, `sdk/xenapi/` (SDK oficial
  vendorizado com `replace xenapi => ./sdk/xenapi`), `frontend/` (Vue 3 + TailwindCSS 4 + Vite), `desktop/` (módulo
  Wails v3.1.0), `web/` (embed da SPA para o futuro modo web), `docs/`, `packaging/`, `openspec` do repositório.
- Núcleo sem UI (`internal/`): conexão e sessão XAPI com TLS confiado na primeira vez e certificados de membros do pool
  confiados pelo coordenador; cache de objetos alimentado por `event.from`; regras portadas de
  `XenAdmin.Desktop.Core` (VM, host, pool, SR, migração, notificações, performance, New VM); operações longas com
  acompanhamento de `task`; proxy do console VNC (RFB) para WebSocket.
- Camada de serviço neutra de transporte (`internal/app`): os mesmos métodos são expostos como serviços do Wails no
  desktop, como ferramentas MCP e, depois, como API HTTP + eventos no modo web.
- CLI Cobra no mesmo binário: `xcp-ng-center-go` (abre a janela), `xcp-ng-center-go mcp [--allow-write]` (servidor MCP por
  stdio com `github.com/modelcontextprotocol/go-sdk`, as mesmas regras de segurança do MCP em C#), `version` e comandos
  de leitura/operação (`server`, `vm`, `host`, `sr`) para scripts.
- Interface Vue com paridade de telas do app C#: árvore de recursos, abas (General, Memory, Storage, Networking,
  Console, Performance, Snapshots, Search), menus e barra de ferramentas, Notifications (Alerts, Events), Properties,
  wizards (New VM, New SR NFS/SMB/iSCSI/ISO, New Pool, Copy VM, Delete VM, Migrate), modo manutenção.
  Tema moderno (claro e escuro), tipografia e espaçamento atuais, **sem `border-radius` em nenhum componente**.
- Recursos removidos do XenServer 9 e os que não existem no XCP-ng (Health Check, licenciamento Citrix, CDN de
  updates, guests PV, SSL legado, FCoE por software, entre outros listados no design) não são portados.
- Empacotamento: Linux (AppImage e tarball), macOS universal (Apple Silicon + Intel, `.app` em zip) e Windows x64
  (exe em zip), mais o binário `xcp-ng-center-go-cli` sem interface (CLI e MCP) nas três plataformas; builds em
  container, exceto o macOS, que precisa do SDK da Apple num runner macOS.

## Capabilities

### New Capabilities
- `go-app-core`: núcleo Go sem UI (sessão XAPI, TLS, cache por eventos, regras portadas, operações, console RFB).
- `go-app-desktop`: app Wails v3 + Vue + TailwindCSS com paridade de telas e fluxos do app C# e tema sem cantos
  arredondados.
- `go-app-cli-mcp`: CLI Cobra e servidor MCP no mesmo binário, com as regras de segurança do MCP em C#.
- `go-app-web-ready`: separação de transporte que permite servir a mesma SPA e o mesmo núcleo por HTTP no futuro.
- `go-app-packaging`: builds e pacotes para Linux, macOS universal e Windows x64, e o fluxo de release.

### Modified Capabilities
(nenhuma; o app C# e o WinForms seguem como estão)

## Impact

- **Código**: nova árvore `desktop-go/` (Go 1.27, Wails v3.1.0, Vue 3, TailwindCSS 4); nenhuma mudança no código C#.
- **Dependências**: `github.com/wailsapp/wails/v3`, `github.com/spf13/cobra`, `github.com/modelcontextprotocol/go-sdk`,
  `github.com/coder/websocket`, `github.com/pelletier/go-toml/v2`, o SDK `xenapi` 26.16.1 vendorizado (BSD-2) e, no
  frontend, Vue, vue-router, TailwindCSS, `@wailsio/runtime` e noVNC (MPL-2.0) para o console.
- **Build/CI**: novo workflow `desktop-go.yml` e alvos `make go-*`; releases com tag própria `go-v*`.
- **Docs**: `desktop-go/README.md`, `desktop-go/docs/`, seção no README raiz, AGENTS.md e CHANGELOG próprio.
- **Fora de escopo nesta mudança**: o servidor web multiusuário em si (só a preparação), atualizações de hosts, HA/DR,
  vApps, VM groups, GPU/vGPU, PVS, WLB, Import/Export (XVA e OVA via ova2xva), assinatura/notarização dos binários.
