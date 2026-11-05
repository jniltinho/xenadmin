## 1. Núcleo portável (prova de conceito já feita nesta branch)

- [x] 1.1 Multi-target `net481;net10.0` em CommandLib, XenCenterLib, XenOvfApi e XenModel, com referências .NET Framework condicionadas a `net481`
- [x] 1.2 Substituir `JavaScriptSerializer` por `Newtonsoft.Json` em `Host.cs` e `TokenManager.cs`
- [x] 1.3 Projeto `tests/XenAdmin.Desktop.Tests` (MSTest) com testes unitários e de integração (login + listar hosts/VMs)
- [x] 1.4 `tools/linux/container.sh build|test` usando `dotnet/sdk:10.0`; `.env` ignorado e `.env.example` versionado
- [x] 1.5 Verificar que `XenModel -f net481` compila no container (reference assemblies)
- [x] 1.6 Adicionar o build `net481` do núcleo ao `container.sh` (`container.sh build-win`) para checagem contínua
- [x] 1.7 Corrigir `container.sh` para não anular valores do `.env` com `-e VAR` vazio

## 2. Abstrações de plataforma

- [x] 2.1 `IXenAdminConfigProvider` Linux com JSON em `$XDG_CONFIG_HOME/xcp-ng-center/` (servidores, usuários, impressões digitais)
- [x] 2.2 Adaptador `ISynchronizeInvoke` sobre `Dispatcher.UIThread` para `InvokeHelper.Initialize`
- [x] 2.3 Validador de certificado para todas as sessões XAPI via `ServicePointManager.ServerCertificateValidationCallback` (fallback que o `JsonRpc` já usa), cobrindo reconexão, `HOST_IS_SLAVE` e sessões derivadas sem mudar o `XenConnection`
- [x] 2.4 Política TOFU (pin tem precedência sobre a cadeia do sistema) com testes unitários: primeiro contato, pin confere, pin diverge com cadeia confiável
- [x] 2.5 Gancho opcional `XenAPI.HTTP.ServerCertificateValidator` (nulo = comportamento atual do Windows); no Linux aplica o TOFU a túneis de console e downloads RRD, e aceita hosts membros cujo certificado o XAPI autenticado informa (`host.certificates`)
- [x] 2.4b Callback TLS não bloqueante: certificado rejeitado fica pendente, a UI pergunta e reconecta (bloquear no handshake causava `408 Request Timeout` do XAPI)
- [ ] 2.6 Armazenamento de senha via Secret Service com detecção em runtime e teste do caminho "indisponível" (MVP atual: senha nunca é gravada)

## 3. App Avalonia: conexão e árvore

- [x] 3.1 Criar `XenAdmin.Desktop` (Avalonia 11.3, `net10.0`), com `CA1416` como erro; fora do `XenAdmin.sln` (msbuild.exe compila só `net481`)
- [x] 3.2 Diálogo Adicionar Servidor, login, `HOST_IS_SLAVE`, erros de rede/autenticação
- [x] 3.3 Árvore pool/host/VM/SR a partir do cache do `XenConnection`, com atualização por eventos e reconexão
- [x] 3.4 Painel de propriedades somente leitura e busca por nome
- [x] 3.5 Teste de integração: conectar via `XenConnection` (não `Session` direta), aguardar cache e verificar hosts/VMs, incluindo queda e reconexão

## 4. Ciclo de vida de VM

- [x] 4.1 Comandos start/shutdown/reboot/suspend/resume/force-* habilitados por `allowed_operations`, reusando as Actions do XenModel com callbacks Linux não nulos
- [x] 4.1b Veto de HA na `VMStartAbstractAction` via propriedade `AllowHaUnprotect` (padrão `true` = Windows inalterado; Linux usa `false` e a Action falha sem `HAUnprotectVMAction`); callbacks nulos tolerados; teste unitário
- [ ] 4.1c Pré-check de HA na UI Linux apenas para desabilitar o comando cedo
- [x] 4.2 Confirmação para operações forçadas; painel de tarefas com progresso e erros (`Failure` legível)
- [~] 4.3 Matriz de integração (gate de release, `XCP_ALLOW_WRITE=1`): feito com VM descartável sem disco — start, force reboot, force shutdown (sucesso) e suspend sem guest tools (falha esperada). Pendente: clean shutdown/reboot e suspend/resume com uma VM com guest tools
- [ ] 4.4 Testes unitários: regras de habilitação por `allowed_operations`, confirmação de operações forçadas, e exibição de `RBAC_PERMISSION_DENIED`/`HOST_NOT_ENOUGH_FREE_MEMORY` com `Failure` simulada

## 5. Console VNC (marco 2)

- [ ] 5.0 Teste de regressão do console no Windows (manual documentado ou automatizado) antes de mexer no WinForms
- [x] 5.1a (atalho) Compilar `VNCStream`/`MyStream`/`VNCException`/`IVNCGraphicsClient` sem alterações no app Linux, com shims portáveis (`Bitmap`/`Color`/`PixelFormat`/`Win32`) no namespace `DotNetVnc`
- [ ] 5.1 Extrair `VNCStream`/`MyStream`/`VNCException` para `XenAdmin.Vnc` multi-target; trocar `Win32.QueryPerformance*` por `Stopwatch`
- [ ] 5.1b Novo contrato de framebuffer por buffer de pixels no lugar de `Bitmap` em `IVNCGraphicsClient`, substituindo `GraphicsUtils` (GDI)
- [ ] 5.2 Adaptador WinForms (`Bitmap`) para o novo contrato e WinForms consumindo `XenAdmin.Vnc`; repetir o teste 5.0
- [x] 5.3 Renderizador Avalonia (`WriteableBitmap`) com escala proporcional e conversão de coordenadas do mouse
- [~] 5.4 Teclado por keysym X11 e Ctrl+Alt+Del (+ Ctrl+Alt+Insert) feitos; falta teste com layout ABNT2 numa VM descartável
- [~] 5.5 Validado manualmente (somente visualização) no .114 com túnel CONNECT + TOFU; falta teste automatizado

## 6. Empacotamento

- [x] 6.1 `dotnet publish -r linux-x64 --self-contained` em container gerando `.tar.gz`
- [x] 6.2 AppImage via `appimagetool` em container, com `.desktop` e ícone
- [x] 6.3 Smoke test: executar o artefato em container Ubuntu 22.04 limpo com `xvfb-run` e verificar que a janela abre
- [x] 6.4 Workflow `.github/workflows/linux.yml` (via `make`): testes, build `net481`, AppImage, smoke e release em tags `linux-v*`
- [x] 6.5 README: build em container, `.env`, limitações conhecidas (Wayland, XCM, export)

## 7. Paridade com o app Windows (revisão Codex: visual ~85%, funcional ~40%)

- [x] 7.1 Layout do WinForms: menus, toolbar com ícones originais (cinza quando desabilitados), painel de navegação, faixa de título, cabeçalho teal, seções PDSection, grids com botões
- [x] 7.2 Árvore: pool/host/VM/SR via `Home()`, SR de tools oculto, templates de usuário, `+/−`, expansão por nó, nome do pool para servidor desconectado
- [x] 7.3 Eventos: duplo clique (conectar/iniciar/retomar), menu de contexto com os comandos e textos originais, diálogos "Connecting to…" e "Connect to Server", confirmações de operações forçadas, aba lembrada por objeto
- [x] 7.4 Abas na ordem original; as não portadas aparecem desabilitadas
- [x] 7.5 Diálogo Properties: páginas General (nome, descrição, pasta via "Change Folder", tags via "Edit Tags", IQN do host), CPU (vCPUs, topologia, prioridade) e Boot Options (autostart, ordem de boot HVM, dispositivo/parâmetros PV), salvando como o original (`SaveChangesAction` no clone + actions das páginas em `MultipleAction`, objeto travado); teste de integração de escrita no 110
- [ ] 7.5b Páginas restantes do Properties (Custom Fields, Start Options/HA, Alerts, USB, Advanced, Multipathing, Log Destination) — hoje listadas desabilitadas
- [x] 7.5c Controles com o visual do WinForms: botões 75×23 (borda azul no padrão), TextBox/ComboBox 23 px, checkbox/radio de 13 px; "Connect to Server" com grupo "User login credentials"
- [x] 7.5d `Folders.InitFolders()` na inicialização (mantém `IXenObject.Path`; sem isso pastas apareciam como <None>)
- [x] 7.6a Aba Snapshots (lista, Take Snapshot/Revert/Delete/Properties, detalhes) e aba Search (overview pool→hosts→VMs→templates com métricas ao vivo do `MetricUpdater`); testes: snapshot→revert→delete no 110 e métricas RRD no pool .236
- [x] 7.6b Aba Performance: gráficos padrão de host e VM (CPU, memória usada, rede, disco) sobre os arquivos RRD (`rrd_updates`), Zoom de 10 minutos a 1 ano, eventos de ciclo de vida; teste somente leitura no pool .236
- [ ] 7.6 Performance: designer de gráficos (New/Edit/Delete, Move Up/Down) e navegador de período; abas NICs, GPU, USB, console do host; Snapshots em árvore e com screenshot; Search editável
- [x] 7.7a Aba Storage da VM: criar drive de DVD, trocar/ejetar ISO, Add (New Disk), Attach Disk, Activate/Deactivate, Detach, Delete, Properties do disco, com as regras de habilitação dos comandos originais; teste de escrita no 110
- [x] 7.7b Ações de host: Shut Down, Reboot (confirmações com coordenador e HCI), Enter/Exit Maintenance Mode (salva e restaura as VMs evacuadas, não reduz a tolerância de HA sem consentimento); teste no lab KVM com evacuação, restauração e reboot
- [ ] 7.7 Diálogo de evacuação (destinos por VM), Move de disco
- [x] 7.8a Wizard New VM (Template, Name, Installation Media, Home Server, CPU & Memory, Storage, Networking, Finish → `CreateVMAction`), Copy VM (fast clone / full copy) e Delete VM (discos e snapshots); teste no 110: criar (Debian 11, 512 MB, disco 64 MB) → iniciar → reboot → desligar → clone → copy → apagar tudo
- [x] 7.8b Create New Pool (coordenador, membros com as razões de `PoolJoinRules`, avisos de CPU/licença/AD, `CreatePoolAction`); teste no lab: forma o pool, confere e ejeta o membro
- [x] 7.8c Lab de 2 hosts XCP-ng 8.3 aninhados no VirtualBox (`tools/vagrant/xcp-ng-lab`), instalação automática por answerfile
- [ ] 7.8d New Storage Repository (NFS e SMB: Scan, criar, reanexar/introduzir) e Migrate to Server (intra-pool), com testes no lab
- [x] 7.8e Console: scancodes para HVM (QEMU extended key event), keysyms com fallback US para PV; logout da sessão do console removido
- [x] 7.8f Console: "Launch SSH Console" (terminal do sistema com OpenSSH) e "Turn on / Switch to Remote Desktop" (cliente RDP do sistema: FreeRDP no terminal ou Remmina; mstsc no Windows; app RDP no macOS), com as regras do VNCTabView; reconexão do console corrigida
- [ ] 7.8 Wizards: Import/Export, demais tipos de SR (iSCSI, HBA, ISO libraries)
- [x] 7.9a Notifications > Alerts (severidade, texto expansível, filtros por severidade e servidor, Dismiss/Dismiss All, contador no botão e no status bar) e Events
- [ ] 7.9b Objects, Organization Views e Saved Searches; filtro por data e Export dos alertas
- [ ] 7.10 Detalhes visuais pendentes: abas inclinadas do WinForms, recuo e seleção da árvore só no rótulo, menus Back/Forward com histórico

## 8. Modularização para outros sistemas

- [x] 8.1 Renomear `XenAdmin.Linux` → `XenAdmin.Desktop` e `tests/XenModel.Linux.Tests` → `tests/XenAdmin.Desktop.Tests`
- [x] 8.2 `XenAdmin.Desktop.Core` sem UI: árvore, detalhes, regras dos comandos, defaults do New VM, certificados e configuração; UI só consulta
- [x] 8.3 Seam `IPlatformServices` (diretório de configuração por SO, cofre de senhas) com testes por plataforma
- [x] 8.4a macOS universal: `make macos` publica `osx-arm64` e `osx-x64` num único `XCP-ng Center.app` (launcher escolhe a CPU, Info.plist, ícone) e gera `xcp-ng-center-<ver>_macos_universal.zip`; job no workflow e asset do release; testado no macOS 15 (x86_64, KVM)
- [x] 8.4b Página Home Server do Properties (AffinityPicker: sem servidor com storage compartilhado, ou um servidor vivo; `VM.set_affinity`), testada no pool KVM
- [x] 8.4c Windows x64: `make windows` publica o app Avalonia como `xcp-ng-center.exe` único (`-win-x64.zip`), no workflow e no release; testado no Windows 11 (dockur/KVM) conectando ao pool do lab
- [x] 8.4e New Storage: iSCSI (LVMoISCSI: CHAP, IQN/LUN, probe → Format/Reattach) e bibliotecas NFS ISO / SMB ISO; target iSCSI (tgt) no lab; testes de integração iSCSI e NFS ISO no pool KVM
- [x] 8.4f QA de paridade (agente; relatório em `_output/qa/report.md`): corrigidos New VM sem Add de disco, confirmações de Shut Down/Reboot/Suspend e de modo manutenção, ações de SR (Set as Default, Detach, Forget, Destroy), faixa de DVD do console, seleção na ordem de boot, Esc nos diálogos, layout de Events, IP duplicado, separadores, SSH para Alpine
- [x] 8.4g1 Ações de rede da VM (Add Interface/Properties/Remove/Activate, coluna Limit), testadas no .110
- [x] 8.4g2 Aba Storage do SR (Rescan, Add, Delete, Properties; colunas Storage Volume e CBT), testada no .110
- [x] 8.4g3 Aba Storage de host e pool (New SR, Reclaim freed space, Properties) e grade selecionável comum (`SelectableGrid`)
- [x] 8.4g4 Aba Networking de host e pool e wizard New Network (External/VLAN e Single-Server Private), VLAN testada no lab num host sozinho e no pool; Bonded, SR-IOV e Delete Bond pendentes
- [x] 8.4g5 Aba NICs do host (lista e Rescan); Create/Delete Bond pendentes
- [x] 8.4g6 Menus: Restart Toolstack, Remove From Pool, Take a Snapshot, Convert to Template (testados no .110 e no lab)
- [x] 8.4g7 Alertas: filtro por data (presets e intervalo personalizado) e Export CSV
- [x] 8.4g8 Move/Migrate Virtual Disk (abas Storage de VM e SR), testado no lab (local → NFS)
- [x] 8.4g9 Menu Templates (Create VM From Selection, Copy, Delete Template, Properties); Export to File pendente
- [x] 8.4g10 Bonded network (wizard e Create Bond), Delete Bond, testados no pool do lab (eth1+eth2)
- [x] 8.4g11 New VM: ícones de SO, categorias e ordem dos templates
- [x] 8.4g12 Console: Undock/Redock e Fullscreen (testado no lab com VM Alpine)
- [x] 8.4g13 Aba Memory: barras de VM e host e Memory Settings (fixa, intervalo, avançado)
- [x] 8.4g14 Menu Pool: Reconnect As, Disconnect, Change Server Password, Rotate Pool Secret (testado no lab)
- [x] 8.4g15 Menu Pool: Add Server, Remove Server, Make into standalone server (testados no lab)
- [x] 8.4g16 HA: wizard Configure HA, prioridades com HA ligado e Disable (testado no pool do lab com heartbeat NFS)
- [x] 8.4g17 Import/Export XVA: VM > Export (verificação opcional) e File > Import em pool/servidor e SR (testado no .110)
- [ ] 8.4g Restantes do QA: botões das abas Networking/Storage de host, pool e SR, menus Pool/Server/VM/Templates/Help completos (HA, DR, Restart Toolstack, Export, Convert to Template, vApp, vTPM), abas NICs/GPU/USB e Console do host, barras da aba Memory, Search editável, Tree View de snapshots, Objects/Organization Views/Saved Searches, HBA/FCoE/Local no New Storage, Undock/Fullscreen do console, filtro de data e Export de alertas, ícones de SO no New VM
- [ ] 8.4d Assinatura e notarização do app macOS (precisa de Mac com Developer ID), assinatura Authenticode no Windows; FreeBSD — fora do escopo atual
