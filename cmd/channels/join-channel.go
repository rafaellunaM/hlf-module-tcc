package channels

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"hlf/internal/fabric"
)

func readCertificate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo %s: %v", filename, err)
	}

	certContent := strings.TrimSpace(string(content))

	lines := strings.Split(certContent, "\n")
	for i := 0; i < len(lines); i++ {
		lines[i] = "        " + lines[i]
	}
	
	return strings.Join(lines, "\n"), nil
}

func JoinChannel(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("não consegui ler %s: %v", configFile, err)
	}

	var partialConfig struct {
		JoinChannel []fabric.JoinChannel `json:"joinChannel"`
	}

	if err := json.Unmarshal(data, &partialConfig); err != nil {
		return fmt.Errorf("não consegui parsear JSON: %v", err)
	}

	for _, joinChannel := range partialConfig.JoinChannel {
		// Itera sobre cada MspID no array
		for mspIndex, mspID := range joinChannel.MspID {
			if !strings.HasPrefix(mspID, "Org") {
				continue
			}

			fmt.Printf("Gerando YAML para %s...\n", mspID)

			// Seleciona o secretKey correto baseado no mspIndex
			var secretKey string
			if mspIndex < len(joinChannel.FileOutputTls) {
				secretKey = joinChannel.FileOutputTls[mspIndex]
			} else {
				secretKey = joinChannel.FileOutputTls[0]
			}

			// Seleciona os anchor peers desta organização (mspIndex)
			var orgAnchorPeers []string
			if mspIndex < len(joinChannel.AnchorPeers) {
				orgAnchorPeers = joinChannel.AnchorPeers[mspIndex]
			} else if len(joinChannel.AnchorPeers) > 0 {
				orgAnchorPeers = joinChannel.AnchorPeers[0]
			}

			// Monta a seção dos anchor peers
			anchorPeersSection := ""
			for i, peer := range orgAnchorPeers {
				anchorPeersSection += fmt.Sprintf(`    - host: %s
      port: 443`, peer)
				if i < len(orgAnchorPeers)-1 {
					anchorPeersSection += "\n"
				}
			}

			// Seleciona os peers to join desta organização (mspIndex)
			var orgPeersToJoin []string
			if mspIndex < len(joinChannel.PeersToJoin) {
				orgPeersToJoin = joinChannel.PeersToJoin[mspIndex]
			} else if len(joinChannel.PeersToJoin) > 0 {
				orgPeersToJoin = joinChannel.PeersToJoin[0]
			}

			// Seleciona os orderer hosts desta organização (mspIndex)
			var orgOrdererHosts []string
			if mspIndex < len(joinChannel.OrderNodeHost) {
				orgOrdererHosts = joinChannel.OrderNodeHost[mspIndex]
			} else if len(joinChannel.OrderNodeHost) > 0 && len(joinChannel.OrderNodeHost[0]) > 0 {
				orgOrdererHosts = joinChannel.OrderNodeHost[0]
			}

			// Seleciona os orderer nodes desta organização (mspIndex)
			var orgOrdererNodes []string
			if mspIndex < len(joinChannel.OrdererNodesList) {
				orgOrdererNodes = joinChannel.OrdererNodesList[mspIndex]
			} else if len(joinChannel.OrdererNodesList) > 0 && len(joinChannel.OrdererNodesList[0]) > 0 {
				orgOrdererNodes = joinChannel.OrdererNodesList[0]
			}

			orderersSection := ""
			for i, ordererNode := range orgOrdererNodes {
				certPath := fmt.Sprintf("/tmp/%s-cert.pem", ordererNode)
				cert, err := readCertificate(certPath)
				if err != nil {
					return fmt.Errorf("erro ao ler certificado %s: %v", certPath, err)
				}

				var ordererURL string
				if i < len(orgOrdererHosts) {
					ordererURL = orgOrdererHosts[i]
				} else if len(orgOrdererHosts) > 0 {
					ordererURL = orgOrdererHosts[0]
				} else {
					ordererURL = "grpcs://ord-node1.default:7050" // fallback
				}

				orderersSection += fmt.Sprintf(`    - certificate: |
%s
      url: %s`, cert, ordererURL)

				if i < len(orgOrdererNodes)-1 {
					orderersSection += "\n"
				}
			}

			peersToJoinSection := ""
			for i, peer := range orgPeersToJoin {
				peersToJoinSection += fmt.Sprintf(`    - name: %s
      namespace: %s`, peer, joinChannel.Namespace)

				if i < len(orgPeersToJoin)-1 {
					peersToJoinSection += "\n"
				}
			}

			// Usa o mspIndex para selecionar o nome do canal
			var channelName string
			if mspIndex < len(joinChannel.FabricChannelFollower) {
				channelName = joinChannel.FabricChannelFollower[mspIndex]
			} else {
				channelName = fmt.Sprintf("%s-%s", joinChannel.FabricChannelFollower[0], strings.ToLower(mspID))
			}

			yaml := fmt.Sprintf(`apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricFollowerChannel
metadata:
  name: %s
spec:
  anchorPeers:
%s
  hlfIdentity:
    secretKey: %s
    secretName: wallet
    secretNamespace: %s
  mspId: %s
  name: demo
  externalPeersToJoin: []
  orderers:
%s
  peersToJoin:
%s`,
				channelName,
				anchorPeersSection,          
				secretKey,        
				joinChannel.Namespace,                
				mspID,                   
				orderersSection,                     
				peersToJoinSection,                  
			)

			yamlFilename := fmt.Sprintf("fabric-follower-channel-%s.yaml", channelName)
			if err := saveYAMLToFile(yaml, yamlFilename); err != nil {
				return fmt.Errorf("erro ao salvar YAML %s: %v", yamlFilename, err)
			}

			fmt.Printf("Aplicando FabricFollowerChannel %s no cluster...\n", channelName)
			if err := applyYAML(yaml); err != nil {
				return fmt.Errorf("erro ao aplicar YAML para %s: %v", mspID, err)
			}

			fmt.Printf("FabricFollowerChannel %s criado com sucesso!\n", channelName)
			fmt.Printf("YAML salvo em: %s\n\n", yamlFilename)
		}
	}

	return nil
}

func saveYAMLToFile(yamlContent, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo %s: %v", filename, err)
	}
	defer file.Close()

	_, err = file.WriteString(yamlContent)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo %s: %v", filename, err)
	}

	fmt.Printf("YAML salvo em: %s\n", filename)
	return nil
}

func applyYAML(yamlContent string) error {
	fmt.Printf("Executando: kubectl apply -f -\n")
	
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yamlContent)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar kubectl apply: %v", err)
	}
	
	return nil
}
