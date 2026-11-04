## Context

Motivação em `proposal.md` (Why). Estado atual medido nesta branch:

| Projeto | Arquivos .cs | Dependência de WinForms | Situação em `net10.0` |
|---|---|---|---|
| CommandLib | 5 | 0 | compila |
| XenCenterLib | 30 | 0 (DPAPI, `user32` P/Invoke em `Win32.cs`) | compila; DPAPI e P/Invoke falham em runtime no Linux |
| XenOvfApi | 23 | 0 (`System.Management`) | compila; WMI indisponível no Linux |
| XenModel (XenAPI, Actions, cache) | 466 | 0 | compila; login + `get_all_records` testados contra host real |
| XenAdmin (UI) | 1134 | 752 arquivos | não portável |

O XenAPI já usa JSON-RPC via `HttpWebRequest`, com callback de validação de certificado por sessão — funciona igual no .NET 10. O protocolo VNC (`XenAdmin/VNC/VNCStream.cs`, 1434 linhas) não depende de WinForms; só o renderizador (`IVNCGraphicsClient`) e o mapeamento de teclas são de UI.

Restrições: manter o WinForms/Windows intacto; .NET 8 sai de suporte em 10/11/2026, então o alvo é .NET 10 (LTS até nov/2028).

## Goals / Non-Goals

**Goals:**
- Um único núcleo C# compartilhado entre o app Windows (`net481`) e o app Linux (`net10.0`), sem fork de código.
- MVP Linux cobrindo as specs `desktop-*` desta change.
- Toda compilação e teste reproduzível em container.

**Non-Goals:**
- Portar telas WinForms existentes (o app Linux é UI nova).
- Migrar o app Windows para .NET 10 ou Avalonia.
- macOS (Avalonia permitiria, mas não será testado nem empacotado).

## Decisions

### D1. UI em Avalonia 11 (C#/XAML), não WinForms/Mono, MAUI ou GTK
- **Avalonia**: maduro no Linux (X11/Wayland via XWayland), MVVM, XAML parecido com WPF, controle `TreeView`/`DataGrid`, `WriteableBitmap` para o framebuffer VNC. Era também a intenção do branch `avalonia` do upstream (que só tem um README).
- WinForms no .NET 10 é Windows-only; WinForms do Mono está sem manutenção e falha com controles customizados/owner-draw que o XenAdmin usa muito.
- MAUI não tem alvo Linux. GTK# (GirCore) seria viável, mas com menos ferramental e MVVM menos natural.

### D2. Multi-target do núcleo (`net481;net10.0`) em vez de copiar para um novo projeto
- Evita divergência: correções no XenAPI/XenModel valem para os dois apps.
- Referências `<Reference>` do .NET Framework ficam condicionadas a `net481`; pacotes de compatibilidade (`ProtectedData`, `System.Management`, `System.Drawing.Common`) entram só em `net10.0` para compilar, mas os caminhos que dependem deles não são chamados pelo app Linux.
- Alternativa descartada: `netstandard2.0` — exigiria mais polyfills e não permite APIs do .NET 10.

### D3. Camada de UI fina sobre `IXenConnection`/cache do XenModel
- O app Linux usa `XenConnection` (login, `HOST_IS_SLAVE`, eventos `event.from`, cache) e as `Actions` existentes (`VMStartAction`, `VMShutdownAction`, …) que já tratam tarefas assíncronas, progresso e RBAC.
- O XenModel já tem os pontos de extensão certos: `InvokeHelper.Initialize(ISynchronizeInvoke)` recebe um adaptador sobre `Dispatcher.UIThread`, e `XenAdminConfigManager.Provider` (`IXenAdminConfigProvider`, hoje `WinformsXenAdminConfigProvider`) ganha uma implementação Linux com JSON em `$XDG_CONFIG_HOME/xcp-ng-center/`.
- Alternativa descartada: chamar XenAPI direto da UI — duplicaria regras já existentes nas Actions.

### D4. Certificados: TOFU com impressão SHA-256
- **Implementado**: o `JsonRpc` do XenAPI cai em `ServicePointManager.ServerCertificateValidationCallback` quando a sessão não tem callback próprio (o mesmo ponto que o Windows usa em `Program.cs:209`); o app Linux registra ali o validador TOFU, cobrindo login, reconexão e redirecionamento ao coordenador sem alterar o `XenConnection`. O callback **nunca bloqueia**: bloquear no handshake esperando o usuário fez o XAPI responder `408 Request Timeout` no teste de UI; o certificado rejeitado fica pendente, a UI pergunta e reconecta se aceito.
- Ordem: (1) impressão digital salva para o host, se existir, tem que conferir — mesmo com cadeia confiável; (2) sem impressão, aceita se a cadeia for confiável e grava; (3) senão pergunta ao usuário. Segue o XenCenter (`Settings.KnownServers`), persistido em JSON.
- **Ponto de injeção (achado da revisão)**: `XenConnection` cria e autentica a `Session` sem callback (`XenConnection.cs:444-450`). Adicionar um validador por conexão (ex.: `IXenConnection.CertificateValidator`) atribuído antes de `login_with_password` e propagado nas reconexões, no redirecionamento `HOST_IS_SLAVE` e nas sessões derivadas.
- **Túnel CONNECT**: `XenAPI/HTTP.cs` usa `ValidateServerCertificate` que retorna sempre `true` (`HTTP.cs:279-286`, `:482-486`) e envia o `session_id` depois do handshake. O validador da conexão MUST ser passado também para `HTTP.CONNECT`. Isso corrige a mesma falha no app Windows; a mudança de comportamento no Windows deve ser revisada à parte (ou ficar atrás de um parâmetro opcional com padrão atual).

### D5. Segredos via Secret Service (libsecret), sem fallback em disco
- DPAPI (`ProtectedData`) não existe no Linux. Usa o Secret Service por D-Bus (via `secret-tool`/libsecret, detectado em runtime). Sem ele, a senha não é salva (spec `desktop-server-connection`).

### D6. Console: reaproveitar `VNCStream` movendo-o para biblioteca sem UI
- Mover `VNCStream`, `MyStream` e `VNCException` para um projeto `XenAdmin.Vnc` multi-target. Não é só extração: `VNCStream` chama `Win32.QueryPerformanceFrequency/Counter` no construtor e no processamento de framebuffer (`VNCStream.cs:167`, `:1244`) → trocar por `Stopwatch`; `IVNCGraphicsClient` expõe `Bitmap` e `GraphicsUtils` usa GDI → novo contrato baseado em buffer de pixels (BGRA + retângulo), com adaptadores WinForms (`Bitmap`) e Avalonia (`WriteableBitmap`). `System.Drawing.Common` não é suportado fora do Windows desde o .NET 7.
- Por isso o console é o **marco 2** desta change, com teste de regressão do console no Windows antes de trocar o WinForms para o novo contrato.
- Túnel: `HTTPHelper.CONNECT` existente no XenModel (`/console?ref=…&session_id=…`).
- Teclado: enviar keysyms X11 derivados do texto/tecla do Avalonia (Linux já fala keysym nativamente), em vez do `KeyMap` de scancodes Windows.
- Alternativa descartada: abrir `vncviewer`/`virt-viewer` externo — o XAPI exige o túnel autenticado, o que exigiria um proxy local.

### D6b. Actions de VM com callbacks Linux e política de HA
- `VMStartAction` invoca callbacks de diagnóstico e de aviso de HA sem checar `null` (`VMStartAction.cs:84-95`) e, para VM protegida não ágil, executa `HAUnprotectVMAction` incondicionalmente. O app Linux MUST fornecer callbacks reais (diagnóstico → mostra a `Failure`).
- Veto na própria Action (o pré-check sozinho tem corrida, pois a Action revalida a agilidade em `VMStartAction.cs:79-87`): propriedade `AllowHaUnprotect` (padrão `true`, Windows inalterado, sem mexer nos 12+ chamadores WinForms). O Linux usa `false` e a Action relança a `Failure` do `assert_agile` sem chamar `HAUnprotectVMAction`.

### D6c. msbuild.exe compila só `net481`
- Os 4 projetos de núcleo usam `TargetFrameworks` condicionado a `$(MSBuildRuntimeType)`: Visual Studio/msbuild.exe (CI Windows existente) continuam compilando apenas `net481`; o `dotnet` CLI compila `net481;net10.0`. `XenAdmin.Desktop` fica fora do `XenAdmin.sln` pelo mesmo motivo.

### D7. Build/teste em container
- `tools/linux/container.sh build|test|build-win` roda `mcr.microsoft.com/dotnet/sdk:10.0` como root e faz `chown` no final (rodar como UID do host quebra o mutex de migração do NuGet). Cache NuGet em `_output/nuget` (já ignorado).
- Testes MSTest; integração marcada com `TestCategory=Integration` e `Assert.Inconclusive` sem `XCP_HOST`.
- Publicação: `make dist|appimage|smoke` usa a imagem `tools/linux/Dockerfile` (SDK 10 + appimagetool 1.9.1 + runtime 20251108 fixados); `dotnet publish -r linux-x64 --self-contained -p:PublishSingleFile=true`. Smoke em `ubuntu:22.04` limpo com Xvfb.

### D8. Módulos separados por camada, prontos para macOS, FreeBSD e Windows sem WinForms
- `XenAdmin.Desktop.Core` (`net10.0`, sem Avalonia, `CA1416` como erro): árvore, detalhes, regras de habilitação dos comandos (`VmOperations`, `VmRules`, `StorageRules`), defaults do wizard New VM, colunas do Search, nomes de ícones, confiança de certificados e configuração. Testável sem UI.
- `XenAdmin.Desktop` (Avalonia, assembly `xcp-ng-center`): só janelas, diálogos e o console; não decide regras, consulta o Core.
- Seam de SO: o contrato `IPlatformServices` fica em `XenAdmin.Desktop.Core/Platform` e os adaptadores em `XenAdmin.Desktop/Platform` (escolhidos no startup): diretório de configuração (XDG no Linux/FreeBSD, `~/Library/Application Support` no macOS, `%APPDATA%\xcp-ng-center` no Windows, separado do app WinForms) e cofre de senhas (`ISecretStore`; hoje `NoSecretStore`, a D5 entra aqui por SO). Escolha via `OperatingSystem.IsX()`; testes injetam a plataforma.
- Empacotamento por SO fica fora do código: `XenAdmin.Desktop/packaging/` (AppImage hoje); macOS e Windows adicionam RID (`osx-arm64`, `win-x64`) e seu pacote (`.app`/`.dmg`, MSIX/zip) sem tocar nas camadas; FreeBSD precisa também do runtime do port da comunidade e de nativos SkiaSharp/HarfBuzz compilados para FreeBSD (os pacotes NuGet não trazem).
- Alternativa descartada: um projeto de UI por SO — duplicaria telas que o Avalonia já desenha igual nas três plataformas.

## Risks / Trade-offs

- [~490 avisos CA1416 no núcleo `net10.0`] → caminhos Windows-only podem ser alcançados por Actions não previstas. Compilar não prova portabilidade. Mitigação: o app Linux só usa uma lista explícita de Actions, e **cada uma** tem teste de integração obrigatório via `XenConnection` + cache (sucesso e falha/RBAC); `CA1416` como erro no `XenAdmin.Desktop` só cobre o código do app, não o corpo das bibliotecas.
- [`CookComputing.XmlRPCv2` e `SharpZipLib 0.85` são pacotes .NET Framework (NU1701)] → XCM/conversão e alguns exports não funcionam no Linux. Mitigação: fora do MVP; trocar SharpZipLib por versão 1.x quando export entrar.
- [Mudança `JavaScriptSerializer` → `Newtonsoft.Json` afeta o Windows] → semântica equivalente para os dois usos (dicionário de strings e DTO com propriedades simples). Mitigação: teste unitário de desserialização do `FileServiceToken`.
- [Wayland] → Avalonia roda via XWayland; captura de teclado do console pode ter limitações em Wayland nativo. Mitigação: documentar; testar em GNOME/KDE.
- [`log4net` 2.0.15 tem aviso de vulnerabilidade (NU1902)] → atualizar para 3.x numa change separada, para não misturar com o port.
- [Testes de integração mexem em host real] → testes de leitura rodam sempre que houver `XCP_HOST`; testes de ciclo de vida exigem uma VM descartável em `XCP_TEST_VM` e ficam *inconclusive* sem ela, mas são **gate obrigatório** antes de release (CI do lab falha se estiverem inconclusive).
- [RBAC não testável no lab] → o host de teste só tem `root` local (RBAC exige autenticação externa). Mitigação: falhas `RBAC_PERMISSION_DENIED` e `HOST_NOT_ENOUGH_FREE_MEMORY` cobertas por testes unitários com `Failure` simulada na camada de Actions/UI.

## Migration Plan

Aditivo: o app Windows não muda e continua sendo o padrão. O app Linux é distribuído como artefato separado. Rollback = remover `XenAdmin.Desktop/` e voltar `TargetFramework` para `net481` nos 4 projetos (o diff do núcleo é pequeno).

## Open Questions

- Nome do produto no Linux (`xcp-ng-center` vs. outro) e ícone/branding — não muda a arquitetura.
- Distribuir também via Flatpak além de AppImage/tar.gz — pode entrar depois sem mudar as specs.
