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

# HLF Auto - Interface

```
 ██   ██ ██      ███████      █████  ██    ██ ████████  ██████  
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ███████ ██      █████       ███████ ██    ██    ██    ██    ██ 
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ██   ██ ███████ ██          ██   ██  ██████     ██     ██████  
```

**Hyperledger Fabric Automated Deployment Tool**  
*Version 1.0 - Created by Rafael Luna*

---

## Seleção de Arquivo de Configuração

| Opção | Arquivo | Descrição |
|-------|---------|-----------|
| 1 | `hlf-config.json` | Configuração padrão |
| 2 | `templates/4-orderers.json` | Template com 4 orderers |
| 3 | `templates/4-peers.json` | Template com 4 peers |

**Exemplo de uso:**
```
Escolha o arquivo de configuração (1-3): 3
Usando: templates/4-peers.json
```

---

## Menu Principal - Hyperledger Fabric Deployment CLI

### Comandos Disponíveis

| # | Comando | Descrição |
|---|---------|-----------|
| 1 | **Create CAs** | Criar Certificate Authorities |
| 2 | **Register Orderers** | Registrar orderers |
| 3 | **Register Peers** | Registrar peers |
| 4 | **Deploy Peers** | Fazer deploy dos peers |
| 5 | **Deploy Orderers** | Fazer deploy dos orderers |
| 6 | **Register Channels** | Registrar channels |
| 7 | **Enroll Channels** | Fazer enroll dos channels |
| 8 | **Create Wallet** | Criar wallet |
| 9 | **Execute PEM Script** | Extrair certificado PEM |
| 10 | **Create Main Channel** | Criar canal principal |
| 11 | **Delete recursos** | Deletar todos os recursos HLF e secrets |
| 12 | **Mostrar recursos** | Mostrar todos recursos para o HLF |
| 13 | **Change Config** | Alterar arquivo de configuração |

### Opções Especiais

| Comando | Descrição |
|---------|-----------|
| `all` | Executar todos os passos em sequência |
| `1-5` | Executar passos de 1 a 5 |
| `1,3,5` | Executar passos específicos (1, 3 e 5) |
| `q` | Sair da aplicação |



# Query for test
* Query for default and four orderes templates 
```
kubectl hlf chaincode invoke --config=org1.yaml \
   --user=admin --peer=org1-peer0.default \
   --chaincode=asset --channel=demo \
   --fcn=initLedger -a '[]'
```

```
kubectl hlf chaincode query --config=org1.yaml \
   --user=admin --peer=org1-peer0.default \
   --chaincode=asset --channel=demo \
   --fcn=GetAllAssets -a '[]'
```

* Query for four peers template

```
kubectl hlf chaincode invoke --config=org1.yaml \
    --user=admin --peer=org1-peer1.default \
    --chaincode=asset --channel=demo \
    --fcn=initLedger -a '[]'

kubectl hlf chaincode query --config=org1.yaml \
    --user=admin --peer=org1-peer1.default \
    --chaincode=asset --channel=demo \
    --fcn=GetAllAssets -a '[]'
```