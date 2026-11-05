## ADDED Requirements

### Requirement: Builds em container
`make go-test`, `make go-linux` e `make go-windows` SHALL rodar em container (Go 1.27, Node 24, GTK4/WebKitGTK 6.0
do Wails v3), sem instalar nada no host, como os alvos atuais do app C#, e MUST rodar os checks em cada módulo Go
(raiz, `desktop/` e `sdk/xenapi`). O build macOS é a exceção: SHALL rodar num runner macOS, porque o Wails precisa de
CGO com o SDK da Apple.

#### Scenario: Build limpo
- **WHEN** `make go-test go-linux` roda numa máquina só com Docker
- **THEN** os testes passam e o binário Linux é gerado em `desktop-go/dist/`

### Requirement: Artefatos por plataforma
Cada release SHALL publicar: `xcp-ng-center-go-<ver>-x86_64.AppImage` e `-linux-x64.tar.gz`,
`xcp-ng-center-go-<ver>_macos_universal.zip` (um `.app` com binário universal arm64 + x86_64 via `lipo`, mínimo
macOS 14) e `xcp-ng-center-go-<ver>-win-x64.zip`, cada um com o binário `xcp-ng-center-go-cli` sem interface, e
`SHA256SUMS`. O AppImage MUST ser gerado sobre uma base fixa (Ubuntu 24.04) e o README MUST listar as distribuições
suportadas (glibc e WebKitGTK 6.0 dessa base ou mais novas).

#### Scenario: App macOS universal
- **WHEN** o release roda o job macOS
- **THEN** `lipo -archs` do binário do `.app` mostra `x86_64 arm64`

#### Scenario: Distribuições suportadas
- **WHEN** o CI roda o AppImage sob Xvfb em containers Ubuntu 24.04, Debian 13 e Fedora 42
- **THEN** a janela abre nas três; numa distribuição mais antiga o README já avisa que não é suportada

#### Scenario: Teste de fumaça
- **WHEN** o CI abre o AppImage sob Xvfb com `XCP_SMOKE=1`
- **THEN** a janela carrega a SPA e o processo sai com código 0

### Requirement: Release separado
Os releases do app Go SHALL usar tags `go-v*` e `desktop-go/CHANGELOG.md`, sem mudar os releases `linux-v*` do app C#,
e MUST ser pre-releases até a paridade com o app C#.

#### Scenario: Tag do app Go
- **WHEN** a tag `go-v0.1.0` é enviada
- **THEN** só o workflow do app Go publica um pre-release com os artefatos acima
