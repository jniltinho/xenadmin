## Purpose

Permite que o usuário do app desktop Linux controle o estado de energia das VMs, com as mesmas regras de permissão e operações permitidas que o XAPI aplica.

## ADDED Requirements

### Requirement: Operações de energia
O app SHALL oferecer para a VM selecionada as operações start, clean shutdown, clean reboot, suspend, resume, force shutdown e force reboot, habilitando somente as que constam em `allowed_operations` da VM.

#### Scenario: VM desligada
- **WHEN** o usuário seleciona uma VM `Halted`
- **THEN** apenas "Iniciar" está habilitado entre as operações de energia

#### Scenario: Iniciar VM
- **WHEN** o usuário escolhe "Iniciar" em uma VM `Halted`
- **THEN** o app dispara a operação assíncrona e a VM passa a `Running` na árvore quando a tarefa termina

### Requirement: Confirmação de operações destrutivas
O app MUST pedir confirmação antes de force shutdown e force reboot.

#### Scenario: Forçar desligamento
- **WHEN** o usuário escolhe "Forçar desligamento"
- **THEN** o app mostra um diálogo de confirmação e só executa se o usuário confirmar

### Requirement: Progresso e erros de tarefas
O app SHALL mostrar progresso das tarefas XAPI disparadas e SHALL exibir a mensagem de falha retornada pelo servidor quando a tarefa falhar.

#### Scenario: Falha por falta de memória
- **WHEN** o start falha com `HOST_NOT_ENOUGH_FREE_MEMORY`
- **THEN** o app mostra o erro descrito em linguagem legível e a VM permanece `Halted`

#### Scenario: Permissão insuficiente
- **WHEN** o usuário autenticado não tem papel RBAC para a operação
- **THEN** o app mostra a falha `RBAC_PERMISSION_DENIED` sem travar a interface

### Requirement: VM protegida por HA não é alterada implicitamente
O app MUST NOT alterar a configuração de HA de uma VM como efeito colateral de start ou resume. Se a VM estiver protegida por HA e não puder ser iniciada mantendo a proteção, o app SHALL abortar e explicar o motivo.

#### Scenario: Start de VM protegida e não ágil
- **WHEN** o usuário inicia uma VM com HA `restart` que não é ágil no pool
- **THEN** o app não inicia a VM, não remove a proteção HA e mostra uma mensagem orientando o ajuste no pool

### Requirement: Interface não bloqueante
O app MUST executar chamadas XAPI fora da thread de interface, mantendo a janela responsiva durante operações longas.

#### Scenario: Várias operações simultâneas
- **WHEN** o usuário inicia duas VMs em sequência rápida
- **THEN** ambas as tarefas aparecem em andamento e a árvore continua navegável
