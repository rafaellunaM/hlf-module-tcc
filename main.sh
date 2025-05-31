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

  # --identities "Org1MSP-tls;org1msp-tlsca.yaml" \