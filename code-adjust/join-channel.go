package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)


func getIndentedCert(resource, jsonPath string) (string, error) {
	cmd := exec.Command("kubectl", "get", resource, "-o", "jsonpath="+jsonPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for i := range lines {
		lines[i] = "        " + lines[i] // 8 espaços
	}
	return strings.Join(lines, "\n"), nil
}

func main() {

	namespace := "default"
	ordererNode := "ord-node1"
	mspID := "Org1MSP"
	secretName := "wallet"
	secretKey := "org1msp.yaml"
	channelName := "demo"
	peerHost := "peer0-org1.localho.st"
	peerName := "org1-peer0"

	fmt.Println("📄 Coletando certificado TLS do orderer...")

	ordererTLS, err := getIndentedCert("fabricorderernodes/"+ordererNode, "{.status.tlsCert}")
	if err != nil {
		fmt.Printf(" Erro ao pegar TLS cert do orderer: %v\n", err)
		os.Exit(1)
	}


	yaml := fmt.Sprintf(`
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricFollowerChannel
metadata:
  name: demo-org1msp
spec:
  anchorPeers:
    - host: %s
      port: 443
  hlfIdentity:
    secretKey: %s
    secretName: %s
    secretNamespace: %s
  mspId: %s
  name: %s
  externalPeersToJoin: []
  orderers:
    - certificate: |
%s
      url: grpcs://%s.%s:7050
  peersToJoin:
    - name: %s
      namespace: %s
`, peerHost, secretKey, secretName, namespace, mspID, channelName, ordererTLS, ordererNode, namespace, peerName, namespace)

	fmt.Println("📤 Aplicando recurso FabricFollowerChannel...")

	applyCmd := exec.Command("kubectl", "apply", "-f", "-")
	applyCmd.Stdin = bytes.NewBufferString(yaml)
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Run(); err != nil {
		fmt.Printf(" Erro ao aplicar FabricFollowerChannel: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("FabricFollowerChannel criado com sucesso.")
}
