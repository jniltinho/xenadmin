## Purpose

Permite que o usuário do app desktop Linux conecte-se a hosts e pools XCP-ng/XenServer, mantenha a lista de servidores entre sessões e confie explicitamente nos certificados TLS desses servidores.

## ADDED Requirements

### Requirement: Adicionar servidor e autenticar
O app SHALL permitir informar endereço (host ou host:porta), usuário e senha, e SHALL abrir uma sessão XAPI via HTTPS usando login por senha. Em caso de falha, o app MUST exibir a mensagem de erro retornada pelo servidor sem gravar a senha.

#### Scenario: Login bem-sucedido
- **WHEN** o usuário informa endereço, usuário e senha válidos e confirma
- **THEN** o app abre a sessão, mostra o servidor como conectado e exibe a versão da API negociada

#### Scenario: Credenciais inválidas
- **WHEN** o servidor responde com `SESSION_AUTHENTICATION_FAILED`
- **THEN** o app mostra o erro de autenticação, mantém o diálogo aberto e não persiste a senha

#### Scenario: Host inalcançável
- **WHEN** a conexão TCP/TLS falha ou excede o tempo limite
- **THEN** o app mostra o erro de rede e o servidor permanece desconectado

### Requirement: Login em membro de pool redireciona ao coordenador
O app SHALL reconectar automaticamente ao coordenador quando o servidor informado for um membro de pool que responde `HOST_IS_SLAVE`.

#### Scenario: Usuário informa um membro do pool
- **WHEN** o login retorna `HOST_IS_SLAVE` com o endereço do coordenador
- **THEN** o app autentica no coordenador com as mesmas credenciais e mostra o pool conectado

### Requirement: Confiança em certificado TLS
O app MUST validar o certificado de toda conexão TLS com o servidor, incluindo a API e os túneis de console. Uma impressão digital SHA-256 gravada tem precedência sobre a confiança do sistema. Sem impressão gravada, certificado não confiável pelo sistema exige confirmação do usuário.

#### Scenario: Primeiro contato com certificado autoassinado
- **WHEN** o servidor apresenta um certificado não confiável ainda não visto
- **THEN** o app mostra a impressão digital e só conecta se o usuário aceitar, gravando a impressão digital aceita

#### Scenario: Certificado alterado
- **WHEN** o servidor apresenta impressão digital diferente da gravada, mesmo que o certificado seja confiável pelo sistema
- **THEN** o app recusa a conexão e mostra um alerta com a impressão antiga e a nova, só conectando se o usuário aceitar explicitamente a nova

#### Scenario: Túnel de console usa a mesma confiança
- **WHEN** o console de uma VM é aberto em um host do pool cujo certificado não confere com o gravado para esse host
- **THEN** o túnel é recusado antes de enviar o identificador da sessão

### Requirement: Persistência da lista de servidores
O app SHALL salvar endereço, usuário e impressão digital aceita de cada servidor em arquivo de configuração no diretório de configuração do usuário (XDG) e SHALL restaurar a lista ao iniciar.

#### Scenario: Reabrir o app
- **WHEN** o usuário fecha e reabre o app
- **THEN** os servidores adicionados anteriormente aparecem na lista, desconectados, prontos para reconectar

### Requirement: Armazenamento opcional de senha
O app SHALL oferecer salvar a senha somente quando um serviço de segredos do desktop (Secret Service) estiver disponível, e MUST NOT gravar a senha em texto puro em disco.

#### Scenario: Secret Service disponível
- **WHEN** o usuário marca "lembrar senha" em um desktop com Secret Service
- **THEN** a senha é guardada no chaveiro do usuário e usada no próximo login automático

#### Scenario: Secret Service indisponível
- **WHEN** não há Secret Service no ambiente
- **THEN** a opção "lembrar senha" fica desabilitada com explicação, e a senha é pedida a cada conexão

### Requirement: Desconectar
O app SHALL encerrar a sessão XAPI (logout) ao desconectar um servidor ou ao fechar o app.

#### Scenario: Usuário desconecta
- **WHEN** o usuário escolhe desconectar um servidor conectado
- **THEN** o app faz logout da sessão e remove os objetos desse servidor da árvore
