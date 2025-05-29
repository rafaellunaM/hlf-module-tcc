package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"hlf"
)

func getIndentedCert(resourceType, resourceName, jsonPath string) (string, error) {
	cmd := exec.Command("kubectl", "get", resourceType, resourceName, "-o", "jsonpath="+jsonPath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	for i := range lines {
		lines[i] = "        " + lines[i]
	}
	return strings.Join(lines, "\n"), nil
}

func buildPeerOrganizationsYAML(peers []Peer, cas []CA) string {
	var result strings.Builder

	mspToPeer := make(map[string]Peer)
	for _, peer := range peers {
		mspToPeer[peer.Mspid] = peer
	}
	
	for _, peer := range mspToPeer {
		var caName string
		for _, ca := range cas {
			if ca.Name == peer.CAName {
				caName = ca.Name
				break
			}
		}
		
		result.WriteString(fmt.Sprintf(`    - mspID: %s
      caName: "%s"
      caNamespace: "default"
`, peer.Mspid, caName))
	}
	
	return result.String()
}

func buildIdentitiesYAML(channels []Channel) string {
	var result strings.Builder
	
	for _, ch := range channels {
		result.WriteString(fmt.Sprintf(`    %s:
      secretKey: %s
      secretName: wallet
      secretNamespace: default
`, ch.MspID, ch.FileOutput))

		if ch.MspID == "OrdererMSP" {
			result.WriteString(fmt.Sprintf(`    %s-tls:
      secretKey: %s
      secretName: wallet
      secretNamespace: default
    %s-sign:
      secretKey: %s
      secretName: wallet
      secretNamespace: default
`, ch.MspID, ch.FileOutput, ch.MspID, ch.FileOutputTls))
		}
	}
	
	return result.String()
}

func buildOrdererOrganizationsYAML(orderers []Orderer, cas []CA) string {
	var result strings.Builder

	mspToOrderer := make(map[string][]Orderer)
	for _, orderer := range orderers {
		mspToOrderer[orderer.Mspid] = append(mspToOrderer[orderer.Mspid], orderer)
	}
	
	for mspID, ordererList := range mspToOrderer {
		var caName string
		for _, ca := range cas {
			if ca.MspID == mspID {
				caName = ca.Name
				break
			}
		}
		
		result.WriteString(fmt.Sprintf(`    - caName: "%s"
      caNamespace: "default"
      externalOrderersToJoin:
`, caName))
		for _, orderer := range ordererList {
			result.WriteString(fmt.Sprintf(`        - host: %s.default
          port: 7053
`, orderer.Name))
		}
		
		result.WriteString(fmt.Sprintf(`      mspID: %s
      ordererEndpoints:
`, mspID))
		for _, orderer := range ordererList {
			result.WriteString(fmt.Sprintf(`        - %s:443
`, orderer.Hosts))
		}
		
		result.WriteString(`      orderersToJoin: []
`)
	}
	
	return result.String()
}

func buildOrderersYAML(orderers []Orderer, tlsCerts []string) string {
	var result strings.Builder
	
	for i, orderer := range orderers {
		result.WriteString(fmt.Sprintf(`    - host: %s
      port: 443
      tlsCert: |-
%s
`, orderer.Hosts, tlsCerts[i]))
	}
	
	return result.String()
}

func buildAdminOrgsYAML(channels []Channel, orgType string) string {
	var result strings.Builder
	seenMSPs := make(map[string]bool)
	
	for _, ch := range channels {
		if !seenMSPs[ch.MspID] {
			if (orgType == "orderer" && ch.MspID == "OrdererMSP") ||
			   (orgType == "peer" && ch.MspID != "OrdererMSP") {
				result.WriteString(fmt.Sprintf(`    - mspID: %s
`, ch.MspID))
				seenMSPs[ch.MspID] = true
			}
		}
	}
	
	return result.String()
}

func main() {
	configFile := "output.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	
	raw, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro lendo %s: %v\n", configFile, err)
		os.Exit(1)
	}

	var config fabric.Config
	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro no unmarshal: %v\n", err)
		os.Exit(1)
	}

	channelName := "demo"

	fmt.Println("📄 Coletando certificados TLS dos orderers...")

	var tlsCerts []string
	var errors []error
	
	for _, orderer := range config.Orderers {
		cert, err := getIndentedCert("fabricorderernodes", orderer.Name, "{.status.tlsCert}")
		tlsCerts = append(tlsCerts, cert)
		errors = append(errors, err)
	}

	hasErrors := false
	for i, err := range errors {
		if err != nil {
			if !hasErrors {
				fmt.Printf("❌ Erros ao pegar TLS certs dos orderers:\n")
				hasErrors = true
			}
			fmt.Printf("  %s: %v\n", config.Orderers[i].Name, err)
		}
	}
	
	if hasErrors {
		os.Exit(1)
	}
	yaml := fmt.Sprintf(`apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricMainChannel
metadata:
  name: %s
spec:
  name: %s
  adminOrdererOrganizations:
%s  adminPeerOrganizations:
%s  channelConfig:
    application:
      acls: null
      capabilities:
        - V2_0
        - V2_5
      policies: null
    capabilities:
      - V2_0
    orderer:
      batchSize:
        absoluteMaxBytes: 1048576
        maxMessageCount: 10
        preferredMaxBytes: 524288
      batchTimeout: 2s
      capabilities:
        - V2_0
      etcdRaft:
        options:
          electionTick: 10
          heartbeatTick: 1
          maxInflightBlocks: 5
          snapshotIntervalSize: 16777216
          tickInterval: 500ms
      ordererType: etcdraft
      policies: null
      state: STATE_NORMAL
    policies: null
  externalOrdererOrganizations: []
  externalPeerOrganizations: []
  peerOrganizations:
%s  identities:
%s  ordererOrganizations:
%s  orderers:
%s`,
		channelName,
		channelName,
		buildAdminOrgsYAML(config.Channels, "orderer"),
		buildAdminOrgsYAML(config.Channels, "peer"),
		buildPeerOrganizationsYAML(config.Peers, config.CAs),
		buildIdentitiesYAML(config.Channels),
		buildOrdererOrganizationsYAML(config.Orderers, config.CAs),
		buildOrderersYAML(config.Orderers, tlsCerts))

	fmt.Printf("📤 Aplicando recurso FabricMainChannel '%s'...\n", channelName)

	applyCmd := exec.Command("kubectl", "apply", "-f", "-")
	applyCmd.Stdin = bytes.NewBufferString(yaml)
	applyCmd.Stdout = os.Stdout
	applyCmd.Stderr = os.Stderr

	if err := applyCmd.Run(); err != nil {
		fmt.Printf("❌ Erro ao aplicar canal: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Canal '%s' criado com sucesso.\n", channelName)
}