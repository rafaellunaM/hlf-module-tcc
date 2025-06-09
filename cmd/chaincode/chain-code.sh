#!/bin/bash

# Hyperledger Fabric Chaincode Installation Script

kubectl hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP

kubectl hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP \
        --ca-name ca  --output peer-org1.yaml

kubectl hlf utils adduser --userPath=peer-org1.yaml --config=org1.yaml --username=admin --mspid=Org1MSP

# Remove existing chaincode files
rm -f code.tar.gz chaincode.tgz

export CHAINCODE_NAME=asset
export CHAINCODE_LABEL=asset

# Create metadata.json
cat << METADATA-EOF > "metadata.json"
{
    "type": "ccaas",
    "label": "${CHAINCODE_LABEL}"
}
METADATA-EOF

# Create connection.json for chaincode as a service
cat > "connection.json" <<CONN_EOF
{
  "address": "${CHAINCODE_NAME}:7052",
  "dial_timeout": "10s",
  "tls_required": false
}
CONN_EOF

# Package chaincode
tar cfz code.tar.gz connection.json
tar cfz chaincode.tgz metadata.json code.tar.gz

export PACKAGE_ID=$(kubectl hlf chaincode calculatepackageid --path=chaincode.tgz --language=node --label=$CHAINCODE_LABEL)

echo "Instalando chaincode..."

kubectl hlf chaincode install --path=./chaincode.tgz \
    --config=org1.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=org1-peer0.default

kubectl hlf externalchaincode sync --image=kfsoftware/chaincode-external:latest \
    --name=$CHAINCODE_NAME \
    --namespace=default \
    --package-id=$PACKAGE_ID \
    --tls-required=false \
    --replicas=1

kubectl hlf chaincode queryinstalled --config=org1.yaml --user=admin --peer=org1-peer0.default

export SEQUENCE=1
export VERSION="1.0"

kubectl hlf chaincode approveformyorg --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --package-id=$PACKAGE_ID \
    --version "$VERSION" --sequence "$SEQUENCE" --name=asset \
    --policy="OR('Org1MSP.member')" --channel=demo

kubectl hlf chaincode commit --config=org1.yaml --user=admin --mspid=Org1MSP \
    --version "$VERSION" --sequence "$SEQUENCE" --name=asset \
    --policy="OR('Org1MSP.member')" --channel=demo
