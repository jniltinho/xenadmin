## Purpose

Dá ao usuário do app desktop Linux acesso ao console gráfico das VMs, tunelado pela conexão XAPI, sem depender de clientes externos.

## ADDED Requirements

### Requirement: Abrir console gráfico
O app SHALL abrir o console RFB (VNC) de uma VM em execução usando o console de protocolo `rfb` da VM, tunelado por HTTP CONNECT na sessão XAPI autenticada, sem abrir portas adicionais no host.

#### Scenario: VM em execução com console
- **WHEN** o usuário abre a aba Console de uma VM `Running`
- **THEN** o app mostra a tela da VM em até 5 segundos

#### Scenario: VM desligada
- **WHEN** a VM não está `Running`
- **THEN** a aba Console mostra que o console está indisponível e reconecta sozinha quando a VM iniciar

### Requirement: Entrada de teclado e mouse
O app SHALL encaminhar teclado e mouse à VM enquanto o console tiver foco, e SHALL oferecer uma ação para enviar Ctrl+Alt+Del.

#### Scenario: Enviar Ctrl+Alt+Del
- **WHEN** o usuário aciona "Enviar Ctrl+Alt+Del"
- **THEN** a VM recebe a sequência de teclas Ctrl, Alt e Delete

#### Scenario: Digitação com layout não-US
- **WHEN** o usuário digita caracteres acentuados com layout ABNT2
- **THEN** a VM recebe os keysyms correspondentes aos caracteres digitados

### Requirement: Redimensionamento da tela
O app SHALL escalar a imagem do console para caber na área disponível, preservando a proporção, com opção de tamanho real.

#### Scenario: Janela menor que a resolução da VM
- **WHEN** a área do console é menor que a resolução da VM
- **THEN** a imagem é reduzida proporcionalmente e as coordenadas do mouse são convertidas para a resolução da VM
