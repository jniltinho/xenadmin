## Why

O XCP-ng Center só roda no Windows (WinForms + .NET Framework 4.8.1), e quem administra XCP-ng a partir de estações Linux depende de Xen Orchestra (web) ou da CLI `xe`. Uma prova de conceito nesta branch mostrou que o núcleo (XenAPI/XenModel) já compila para .NET 10 e fala com um host XCP-ng real (API 2.21) a partir de um container Linux, então um app desktop nativo em C# é viável reaproveitando a maior parte da lógica existente.

## What Changes

- Os projetos de núcleo `CommandLib`, `XenCenterLib`, `XenOvfApi` e `XenModel` passam a ter multi-target `net481;net10.0`; o alvo `net481` (Windows) continua idêntico em comportamento.
- Troca de `System.Web.Script.Serialization.JavaScriptSerializer` por `Newtonsoft.Json` (dependência já existente) nos dois pontos onde era usado, para funcionar nos dois alvos.
- Novo app desktop `XenAdmin.Linux` (Avalonia, .NET 10) com um MVP: conectar a servidores, árvore de recursos com atualização ao vivo, ciclo de vida de VMs e console gráfico (VNC).
- Build, testes unitários e testes de integração executados em container (`mcr.microsoft.com/dotnet/sdk:10.0`), via `tools/linux/container.sh`.
- Credenciais de teste lidas de `.env` local (ignorado pelo git); `.env.example` versionado.
- Fora de escopo nesta change: wizards (New VM, Import/Export OVF, New SR), HA/WLB, atualizações/patching, RDP, relatórios, plugins, e qualquer alteração na UI WinForms existente.

## Capabilities

### New Capabilities
- `desktop-server-connection`: adicionar/remover servidores, login, confiança em certificado, persistência da lista de servidores e de segredos no Linux.
- `desktop-resource-tree`: visualização hierárquica de pools, hosts, VMs e SRs com estado atualizado por eventos do XAPI.
- `desktop-vm-lifecycle`: operações de energia de VMs (start, shutdown, reboot, suspend, resume, variantes forçadas) com progresso e erros.
- `desktop-vm-console`: console gráfico de VM via VNC tunelado pelo XAPI, com teclado e mouse.
- `desktop-linux-packaging`: build reproduzível em container e artefato executável para Linux x64.

### Modified Capabilities
(nenhuma — não há specs existentes em `openspec/specs/`)

## Impact

- **Código**: `*.csproj` dos 4 projetos de núcleo; `XenModel/XenAPI-Extensions/Host.cs`; `XenModel/Actions/Updates/TokenManager.cs`; novos `XenAdmin.Linux/`, `tests/XenModel.Linux.Tests/`, `tools/linux/`.
- **Dependências novas (somente `net10.0`)**: `System.Security.Cryptography.ProtectedData`, `System.Management`, `System.Drawing.Common` (para compilar; caminhos Windows-only não são usados no Linux), Avalonia 11, MSTest.
- **Windows**: o `XenAdmin` WinForms continua `net481` e consome o alvo `net481` do núcleo; sem mudança de comportamento esperada (verificado: `XenModel -f net481` compila).
- **Riscos conhecidos**: ~490 avisos CA1416 (APIs Windows-only) no núcleo em `net10.0`; `CookComputing.XmlRPCv2` (XCM) e `SharpZipLib 0.85` são pacotes .NET Framework e não são exercitados pelo MVP — ver `design.md`.
