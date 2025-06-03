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

# Secret
kubectl delete secrets wallet && \
go run create-generic-wallet.go 

# Consulta
kubectl get secrete |grep wallet
watch kubectl get fabricmainchannels 
watch kubectl get pods

# Delete
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
