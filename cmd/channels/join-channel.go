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
	for i := 1; i < len(lines); i++ {
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
		if !strings.HasPrefix(joinChannel.MspID, "Org") {
			continue
		}

		fmt.Printf("Gerando YAML para %s...\n", joinChannel.MspID)

		orderersSection := ""
		for i, ordererNode := range joinChannel.OrdererNodesList {
			certPath := fmt.Sprintf("/tmp/%s-cert.pem", ordererNode)
			cert, err := readCertificate(certPath)
			if err != nil {
				return fmt.Errorf("erro ao ler certificado %s: %v", certPath, err)
			}

			orderersSection += fmt.Sprintf(`    - certificate: |
%s
      url: %s`, cert, joinChannel.OrderNodeHost[i])

			if i < len(joinChannel.OrdererNodesList)-1 {
				orderersSection += "\n"
			}
		}

		peersToJoinSection := ""
		for i, peer := range joinChannel.PeersToJoin {
			peersToJoinSection += fmt.Sprintf(`    - name: %s
      namespace: %s`, peer, joinChannel.Namespace)

			if i < len(joinChannel.PeersToJoin)-1 {
				peersToJoinSection += "\n"
			}
		}

		yaml := fmt.Sprintf(`apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricFollowerChannel
metadata:
  name: %s
spec:
  anchorPeers:
    - host: %s
      port: 443
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
			joinChannel.FabricChannelFollower[0],
			joinChannel.AnchorPeers[0],          
			joinChannel.FileOutputTls[0],        
			joinChannel.Namespace,                
			joinChannel.MspID,                   
			orderersSection,                     
			peersToJoinSection,                  

		yamlFilename := fmt.Sprintf("fabric-follower-channel-%s.yaml", joinChannel.FabricChannelFollower[0])
		if err := saveYAMLToFile(yaml, yamlFilename); err != nil {
			return fmt.Errorf("erro ao salvar YAML %s: %v", yamlFilename, err)
		}

		if err := applyYAML(yaml); err != nil {
			return fmt.Errorf("erro ao aplicar YAML para %s: %v", joinChannel.MspID, err)
		}

		fmt.Printf("FabricFollowerChannel %s criado com sucesso!\n", joinChannel.FabricChannelFollower[0])
		fmt.Printf("YAML salvo em: %s\n\n", yamlFilename)
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
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yamlContent)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Printf("Executando: kubectl apply -f -\n")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar kubectl apply: %v", err)
	}
	
	return nil
}
