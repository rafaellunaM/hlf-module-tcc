## Sequence

* create-ca.go
* register-user-peers-cas.go
* register-user-orderes-cas.go
* deploy-peers.go
* deploy-orderer.go
* channel-ca-register.go

* channel-ca-enroll.go _ERRO 
* create-generic-wallet.go
* create-main-channel.go

### Depends on: https://github.com/rafaellunaM/IaC-tcc/blob/main/terraform/modules/deployments/code/hlf-config.json


# HLF deploy
go run create-cas.go && \ 
kubectl wait --timeout=60s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all && \
go run register-user-peers-cas.go && \
go run register-user-orderes-cas.go && \
go run deploy-peers.go && \
go run deploy-orderer.go && \
go run channel-ca-register.go && \
go run channel-ca-enroll.go && \
kubectl delete secrets wallet && \
go run create-generic-wallet.go 

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
watch kubectl get pods

# Enroll
kubectl hlf ca enroll --name=ord-ca --namespace=default \
    --user=admin --secret=adminpw --mspId OrdererMSP \
    --ca-name tlsca  --output orderermsp.yaml
    
kubectl hlf ca enroll --name=ord-ca --namespace=default \
    --user=admin --secret=adminpw --mspId OrdererMSP \
    --ca-name ca  --output orderermspsign.yaml

kubectl hlf ca enroll --name=org1-ca --namespace=default \
    --user=admin --secret=adminpw --mspId Org1MSP \
    --ca-name tlsca  --output org1msp-tlsca.yaml
kubectl hlf ca enroll --name=org1-ca --namespace=default \
    --user=admin --secret=adminpw --mspId Org1MSP \
    --ca-name ca  --output org1msp.yaml

kubectl hlf identity create --name org1-admin --namespace default \
    --ca-name org1-ca --ca-namespace default \
    --ca ca --mspId Org1MSP --enroll-id admin --enroll-secret adminpw

kubectl create secret generic wallet --namespace=default \
        --from-file=org1msp.yaml=$PWD/org1msp.yaml \
        --from-file=orderermsp.yaml=$PWD/orderermsp.yaml \
        --from-file=orderermspsign.yaml=$PWD/orderermspsign.yaml
