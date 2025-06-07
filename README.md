# Overview
This repository is a module of this IaC: https://github.com/rafaellunaM/IaC-tcc, but be used alone.
The objective is automation HLF deployment and save a hlf infrastructure status, because it's set from json file.

# Depends on: 
* Custom config: https://github.com/rafaellunaM/IaC-tcc/blob/main/terraform/modules/deployments/code/hlf-config.json
* Bevel operator in the cluster
* Go Lang

# Get started
* Step-1: go run main.go
* Step-2: Choice file config
  + option-1: custom config, but default values for hlf-config.json is: 1 orderer node, 1 peer node and 1 CA
  + option-2: 4 orderers nodes, 1 peer node and 1 CA
  + option-3: 1 orderer node, 4 peers nodes and 1 CA
* Step-3: Choice"all" to create HLF of way automatic
* Step-4: Choice "12" for show resourcers, wait they are running status

# Interface

 ██   ██ ██      ███████      █████  ██    ██ ████████  ██████  
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ███████ ██      █████       ███████ ██    ██    ██    ██    ██ 
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ██   ██ ███████ ██          ██   ██  ██████     ██     ██████  

           Hyperledger Fabric Automated Deployment Tool
               Version 1.0 - Created by Rafael Luna


==================================================
SELEÇÃO DE ARQUIVO DE CONFIGURAÇÃO
==================================================
1. hlf-config.json (configuração padrão)
2. templates/4-orderers.json (4 orderers)
3. templates/4-peers.json (4 peers)
==================================================
Escolha o arquivo de configuração (1-3): 3
Usando: templates/4-peers.json

============================================================
HYPERLEDGER FABRIC DEPLOYMENT CLI
Configuração atual: templates/4-peers.json
============================================================
 1. Create   CAs         - Criar Certificate Authorities
 2. Register Orderers    - Registrar orderers
 3. Register Peers       - Registrar peers
 4. Deploy   Peers       - Fazer deploy dos peers
 5. Deploy   Orderers    - Fazer deploy dos orderers
 6. Register Channels    - Registrar channels
 7. Enroll   Channels    - Fazer enroll dos channels
 8. Create   Wallet      - Criar wallet
 9. Execute  PEM Script  - Extrair certificado PEM
10. Create   Main Channel - Criar canal principal
11. Delete   recursos    - Delete todos os recursos HLF e secret
12. Mostrar  recursos    - Mostra todos recursos para o HLF
13. Change   Config      - Alterar arquivo de configuração
------------------------------------------------------------
Opções especiais:
   all      - Executar todos os passos em sequência
   1-5      - Executar passos de 1 a 5
   1,3,5    - Executar passos 1, 3 e 5
   q        - Sair
============================================================
