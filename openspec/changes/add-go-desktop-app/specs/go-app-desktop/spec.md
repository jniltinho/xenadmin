## ADDED Requirements

### Requirement: App desktop Wails v3 com Vue e TailwindCSS
O app desktop SHALL ser um binário Wails v3.1.0 (`desktop-go/desktop`) com a SPA Vue 3 + TailwindCSS 4 embutida, rodando
em Linux (WebKitGTK), macOS (WKWebView) e Windows (WebView2). A SPA MUST falar com o núcleo só pela interface `api` do
frontend, nunca chamando o runtime do Wails diretamente nas telas.

#### Scenario: Abrir o app
- **WHEN** o usuário executa `xcp-ng-center-go` sem argumentos
- **THEN** abre a janela principal com a árvore, as abas e a barra de status, e a lista de servidores salvos é carregada

#### Scenario: Instância única
- **WHEN** o usuário abre o app uma segunda vez
- **THEN** a janela já aberta vem para a frente e nenhuma segunda instância conecta aos servidores

### Requirement: Paridade de telas com o app C#
A SPA SHALL oferecer todas as telas e fluxos do app desktop C# (`XenAdmin.Desktop` 0.2.0) listados em
`desktop-go/docs/parity.md`, com os mesmos rótulos, ordem de campos, validações e confirmações. Cada tela MUST ser
comparada com a do app C# e do XCP-ng Center v25.04 antes de ser dada como pronta.

#### Scenario: Telas cobertas
- **WHEN** a paridade é conferida no fim da mudança
- **THEN** existem: árvore de recursos com ícones de estado; abas General, Memory, Storage, Networking, Console, Performance, Snapshots e Search; menus, barra de ferramentas e menus de contexto; Notifications (Alerts com filtros e Dismiss, Events); Properties (General, CPU, Boot Options, Home Server); wizards New VM, New Storage Repository (NFS, SMB, iSCSI, NFS ISO, SMB ISO), Create New Pool, Copy VM, Delete VM e Migrate to Server; Shut Down, Reboot e Enter/Exit Maintenance Mode de hosts

#### Scenario: Checklist de paridade
- **WHEN** uma tela é marcada como pronta em `tasks.md`
- **THEN** existe um par de capturas (Go e C#/Windows) em `_output/parity/` e as diferenças funcionais estão resolvidas ou registradas

#### Scenario: Comportamentos do 0.2.0
- **WHEN** a paridade é conferida no fim da mudança
- **THEN** funcionam como no app C#: duplo clique conecta um servidor ou liga uma VM; a aba é lembrada por objeto; Turn on / Switch to Remote Desktop; New SR reanexa ou introduz um SR NFS/SMB existente; sair e voltar à aba Console não derruba a sessão; Full copy nunca vira fast clone; SRs ocultos nunca são o padrão de discos novos

#### Scenario: Comando desabilitado
- **WHEN** um comando não é permitido pelo objeto selecionado
- **THEN** o item de menu e o botão aparecem desabilitados, como no app C#

### Requirement: Layout moderno sem cantos arredondados
A SPA SHALL usar um tema próprio em TailwindCSS (claro e escuro, densidade de app desktop, fonte do sistema) em que
nenhum elemento tem cantos arredondados: o tema MUST definir todos os tokens `--radius-*` como `0` e um teste MUST
falhar se algum componente renderizado tiver `border-radius` diferente de `0`.

#### Scenario: Verificação de cantos
- **WHEN** o teste visual percorre as telas principais e os diálogos
- **THEN** `getComputedStyle(el).borderRadius` é `0px` em todos os elementos

#### Scenario: Tema escuro
- **WHEN** o sistema está em modo escuro e o usuário não escolheu tema
- **THEN** a SPA abre no tema escuro com contraste AA nos textos

### Requirement: Segurança do webview
A janela SHALL carregar só a SPA embutida com uma CSP restrita (`default-src 'self'`, sem `unsafe-eval`, `connect-src`
limitado ao app e ao endpoint loopback do console), MUST bloquear navegação e janelas para fora do app (links externos
abrem no navegador do sistema), MUST negar câmera, microfone, geolocalização e notificações do webview, e MUST
renderizar nomes, descrições e textos de alertas vindos dos servidores sempre como texto (nunca `v-html`).

#### Scenario: Nome de VM com HTML
- **WHEN** uma VM se chama `<img src=x onerror=alert(1)>`
- **THEN** a árvore e as abas mostram o texto literal e nenhum script roda

#### Scenario: Navegação externa
- **WHEN** a SPA tenta navegar para `https://example.com`
- **THEN** a janela não sai do app e o endereço abre no navegador do sistema somente se veio de um link permitido

### Requirement: Integração nativa mínima
Serviços nativos do Wails SHALL cobrir apenas o que a SPA não faz sozinha: abrir o terminal SSH e o cliente RDP do sistema (as mesmas regras do app C#), estado da janela, bandeja e notificações
do sistema.

#### Scenario: SSH em VM Linux
- **WHEN** o usuário escolhe Launch SSH Console numa VM Linux com IP conhecido
- **THEN** o terminal do sistema abre com `ssh` para o IP da VM

### Requirement: Interface não bloqueante
Nenhuma chamada XAPI SHALL rodar na thread de UI do webview; a SPA MUST mostrar progresso e continuar responsiva
durante operações longas e reconexões.

#### Scenario: Servidor lento
- **WHEN** um servidor demora a responder ao login
- **THEN** a árvore mostra "Connecting..." para ele e o resto da interface continua utilizável
