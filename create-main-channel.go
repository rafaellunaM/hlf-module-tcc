package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"hlf/internal/fabric"
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

func buildPeerOrganizationsYAML(peers []fabric.Peer, cas []fabric.CA) string {
	var result strings.Builder

	mspToPeer := make(map[string]fabric.Peer)
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

func buildIdentitiesYAML(channels []fabric.Channel) string {
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

func buildOrdererOrganizationsYAML(orderers []fabric.Orderer, cas []fabric.CA) string {
	var result strings.Builder

	mspToOrderers := make(map[string][]fabric.Orderer)
	for _, orderer := range orderers {
		mspToOrderers[orderer.Mspid] = append(mspToOrderers[orderer.Mspid], orderer)
	}
	
	for mspID, ordererList := range mspToOrderers {
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

func buildOrderersYAML(orderers []fabric.Orderer, tlsCerts []string) string {
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

func buildAdminOrgsYAML(channels []fabric.Channel, orgType string) string {
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

	var jsonStructure map[string]interface{}
	if err := json.Unmarshal(raw, &jsonStructure); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro no unmarshal para debug: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("🔍 Debug - Estrutura do JSON:\n")
	for key, value := range jsonStructure {
		switch v := value.(type) {
		case []interface{}:
			fmt.Printf("  %s: array com %d elementos\n", key, len(v))
		case interface{}:
			fmt.Printf("  %s: %T\n", key, v)
		}
	}
	
	var config fabric.Config
	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro no unmarshal: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🔍 Debug - Dados carregados:\n")
	fmt.Printf("  CAs: %d\n", len(config.CAs))
	fmt.Printf("  Peers: %d\n", len(config.Peers))
	fmt.Printf("  Orderers: %d\n", len(config.Orderers))
	fmt.Printf("  Channels: %d\n", len(config.Channels))

	if len(config.Channels) == 0 {
		fmt.Fprintf(os.Stderr, "❌ Erro: Nenhum canal encontrado no arquivo de configuração\n")
		os.Exit(1)
	}

	if len(config.Orderers) == 0 {
		fmt.Fprintf(os.Stderr, "❌ Erro: Nenhum orderer encontrado no arquivo de configuração\n")
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
		fmt.Println("⚠️  Continuando mesmo com erros nos certificados TLS...")
		for i := range tlsCerts {
			if errors[i] != nil {
				tlsCerts[i] = "        # Certificado não encontrado"
			}
		}
	}

	adminOrdererOrgs := buildAdminOrgsYAML(config.Channels, "orderer")
	adminPeerOrgs := buildAdminOrgsYAML(config.Channels, "peer")
	peerOrgs := buildPeerOrganizationsYAML(config.Peers, config.CAs)
	identities := buildIdentitiesYAML(config.Channels)
	ordererOrgs := buildOrdererOrganizationsYAML(config.Orderers, config.CAs)
	orderers := buildOrderersYAML(config.Orderers, tlsCerts)

	fmt.Printf("🔍 Debug - Conteúdo gerado:\n")
	fmt.Printf("  adminOrdererOrgs vazio: %t\n", strings.TrimSpace(adminOrdererOrgs) == "")
	fmt.Printf("  adminPeerOrgs vazio: %t\n", strings.TrimSpace(adminPeerOrgs) == "")
	fmt.Printf("  peerOrgs vazio: %t\n", strings.TrimSpace(peerOrgs) == "")
	fmt.Printf("  identities vazio: %t\n", strings.TrimSpace(identities) == "")
	fmt.Printf("  ordererOrgs vazio: %t\n", strings.TrimSpace(ordererOrgs) == "")
	fmt.Printf("  orderers vazio: %t\n", strings.TrimSpace(orderers) == "")

	if strings.TrimSpace(adminOrdererOrgs) == "" {
		fmt.Fprintf(os.Stderr, "❌ Erro: adminOrdererOrganizations está vazio. Verificar se existe canal com MspID 'OrdererMSP'\n")
	}
	if strings.TrimSpace(adminPeerOrgs) == "" {
		fmt.Fprintf(os.Stderr, "❌ Erro: adminPeerOrganizations está vazio. Verificar se existem canais com MspID diferente de 'OrdererMSP'\n")
	}
	if strings.TrimSpace(identities) == "" {
		fmt.Fprintf(os.Stderr, "❌ Erro: identities está vazio. Verificar dados dos canais\n")
	}
	if strings.TrimSpace(ordererOrgs) == "" {
		fmt.Fprintf(os.Stderr, "❌ Erro: ordererOrganizations está vazio. Verificar se orderers têm Mspid correspondente aos CAs\n")
	}
	if strings.TrimSpace(orderers) == "" {
		fmt.Fprintf(os.Stderr, "❌ Erro: orderers está vazio. Verificar dados dos orderers\n")
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
		adminOrdererOrgs,
		adminPeerOrgs,
		peerOrgs,
		identities,
		ordererOrgs,
		orderers)

	fmt.Println("\n🔍 YAML gerado:")
	fmt.Println("================")
	fmt.Println(yaml)
	fmt.Println("================\n")

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
