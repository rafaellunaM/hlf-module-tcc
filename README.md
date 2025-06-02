## Sequence

* create-ca.go
* register-user-peers-cas.go
* register-user-orderes-cas.go
* deploy-peers.go
* deploy-orderer.go
* channel-ca-register.go
* channel-ca-enroll.go
* create-generic-wallet.go
* create-main-channel.go

### Depends on: https://github.com/rafaellunaM/IaC-tcc/blob/main/terraform/modules/deployments/code/hlf-config.json

# secret
kubectl delete secrets wallet && \
go run create-generic-wallet.go 

# consulta
kubectl get secrete |grep wallet
watch kubectl get fabricmainchannels 
watch kubectl get pods

# delete
rm -rf /tmp/*
rm *.yaml
kubectl delete fabricorderernodes.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricchaincode.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricmainchannels --all-namespaces --all
kubectl delete fabricfollowerchannels --all-namespaces --all
kubectl delete secrets wallet
watch kubectl get pods

# Enroll
kubectl hlf ca enroll --name=ord-ca --namespace=default \
    --user=admin --secret=adminpw --mspid OrdererMSP \
    --ca-name tlsca  --output orderermsp.yaml
    
kubectl hlf ca enroll --name=ord-ca --namespace=default \
    --user=admin --secret=adminpw --mspid OrdererMSP \
    --ca-name ca  --output orderermspsign.yaml

kubectl hlf ca enroll --name=org1-ca --namespace=default \
    --user=admin --secret=adminpw --mspid Org1MSP \
    --ca-name tlsca  --output org1msp-tlsca.yaml
kubectl hlf ca enroll --name=org1-ca --namespace=default \
    --user=admin --secret=adminpw --mspid Org1MSP \
    --ca-name ca  --output org1msp.yaml

kubectl hlf identity create --name org1-admin --namespace default \
    --ca-name org1-ca --ca-namespace default \
    --ca ca --mspid Org1MSP --enroll-id admin --enroll-secret adminpw

kubectl create secret generic wallet --namespace=default \
        --from-file=org1msp.yaml=$PWD/org1msp.yaml \
        --from-file=orderermsp.yaml=$PWD/orderermsp.yaml \
        --from-file=orderermspsign.yaml=$PWD/orderermspsign.yaml

# Pem script and main-channel manual
´´´
set -e
cat orderermsp.yaml | grep -A 100 "pem: |" | sed 's/.*pem: |//' | sed '/^[[:space:]]*$/d' | sed 's/^[[:space:]]*//' > /tmp/orderer-cert.pem


kubectl hlf channelcrd main create \
  --name demo \
  --channel-name demo \
  --secret-name wallet \
  --admin-orderer-orgs OrdererMSP \
  --orderer-orgs OrdererMSP \
  --identities "OrdererMSP;orderermsp.yaml" \
  --identities "OrdererMSP-sign;orderermspsign.yaml" \
  --admin-peer-orgs Org1MSP \
  --peer-orgs Org1MSP \
  --identities "Org1MSP;org1msp.yaml" \
  --secret-ns default \
  --consenters "orderer0-ord.localho.st:443" \
  --consenter-certificates /tmp/orderer-cert.pem
´´´
