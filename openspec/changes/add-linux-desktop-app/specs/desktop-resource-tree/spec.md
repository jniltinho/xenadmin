## Purpose

Mostra ao usuário do app desktop Linux a hierarquia de pools, hosts, VMs e storage dos servidores conectados, sempre refletindo o estado atual do XAPI.

## ADDED Requirements

### Requirement: Hierarquia de recursos
O app SHALL exibir, para cada servidor conectado, uma árvore com pool (quando existir), hosts, VMs (excluindo templates, snapshots e o domínio de controle) e SRs, cada nó com nome e ícone de estado.

#### Scenario: Pool com hosts e VMs
- **WHEN** um pool com dois hosts e VMs residentes é conectado
- **THEN** a árvore mostra o pool, os dois hosts sob ele e cada VM em execução sob o host em que reside

#### Scenario: VM desligada
- **WHEN** uma VM está `Halted` e não tem host residente
- **THEN** a VM aparece sob o pool (ou sob o host único, se não houver pool), com ícone de desligada

### Requirement: Atualização ao vivo por eventos
O app SHALL acompanhar o XAPI pelo mecanismo de eventos e SHALL refletir criação, remoção e mudança de estado de objetos sem ação do usuário, em até 5 segundos após o evento.

#### Scenario: VM ligada por outra ferramenta
- **WHEN** uma VM é iniciada via `xe` ou Xen Orchestra
- **THEN** o ícone e a posição da VM na árvore são atualizados sem recarregar

#### Scenario: Conexão perdida
- **WHEN** a conexão com o servidor cai
- **THEN** o nó do servidor é marcado como desconectado e o app tenta reconectar automaticamente, sem intervenção do usuário

### Requirement: Painel de propriedades
O app SHALL mostrar, para o nó selecionado, um painel somente leitura com as propriedades principais (nome, descrição, UUID, estado de energia, memória, vCPUs, host residente para VMs; versão e uso de memória para hosts; tamanho e uso para SRs).

#### Scenario: Selecionar uma VM
- **WHEN** o usuário seleciona uma VM na árvore
- **THEN** o painel mostra nome, UUID, estado de energia, memória, vCPUs e host residente atuais

### Requirement: Busca por nome
O app SHALL filtrar a árvore por texto, mantendo visíveis os ancestrais dos nós que casam.

#### Scenario: Filtrar VMs
- **WHEN** o usuário digita parte do nome de uma VM no campo de busca
- **THEN** a árvore mostra apenas os nós cujo nome contém o texto (sem diferenciar maiúsculas) e seus ancestrais
