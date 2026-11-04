## Purpose

Garante que o app desktop Linux e os testes do núcleo sejam construídos e testados de forma reproduzível em container, gerando um artefato executável sem exigir o .NET instalado na máquina do usuário.

## ADDED Requirements

### Requirement: Build e testes em container
O projeto SHALL fornecer um comando único que compila e executa os testes dentro da imagem oficial do .NET SDK, sem exigir SDK instalado no host, e MUST deixar os arquivos gerados com o dono do usuário que executou.

#### Scenario: Testes sem credenciais
- **WHEN** o desenvolvedor executa os testes sem `XCP_HOST` definido
- **THEN** os testes unitários rodam e os testes de integração são reportados como ignorados

#### Scenario: Testes de integração
- **WHEN** `.env` define `XCP_HOST`, `XCP_USER` e `XCP_PASS`
- **THEN** os testes de integração autenticam no host e listam hosts e VMs

### Requirement: Segredos fora do repositório
Credenciais de teste MUST ficar apenas em `.env` local, que MUST estar no `.gitignore`; o repositório SHALL conter apenas `.env.example` com valores fictícios.

#### Scenario: Verificar o repositório
- **WHEN** alguém executa `git check-ignore .env`
- **THEN** o arquivo é reportado como ignorado e nenhuma senha real aparece nos arquivos versionados

### Requirement: Artefato executável para Linux x64
O build de release SHALL publicar o app como self-contained para `linux-x64` e empacotá-lo como `.tar.gz` e AppImage, executáveis em Ubuntu 22.04+ e Fedora 39+ sem .NET instalado.

#### Scenario: Executar em máquina limpa
- **WHEN** o usuário extrai o `.tar.gz` (ou executa o AppImage) em um Ubuntu 22.04 sem .NET
- **THEN** o app abre a janela principal

### Requirement: Alvo Windows preservado
O build MUST continuar compilando os projetos de núcleo para `net481`, de modo que o app WinForms existente não mude de comportamento.

#### Scenario: Compilar o núcleo para Windows
- **WHEN** o núcleo é compilado com alvo `net481`
- **THEN** a compilação termina sem erros
