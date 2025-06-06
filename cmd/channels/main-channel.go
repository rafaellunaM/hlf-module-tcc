package channels

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"hlf/internal/fabric"
)

type PeerOrganization struct {
	MspID        string `yaml:"mspID"`
	CaName       string `yaml:"caName"`
	CaNamespace  string `yaml:"caNamespace"`
}

type OrdererOrganization struct {
	CaName                 string            `yaml:"caName"`
	CaNamespace            string            `yaml:"caNamespace"`
	ExternalOrderersToJoin []ExternalOrderer `yaml:"externalOrderersToJoin"`
	MspID                  string            `yaml:"mspID"`
	OrdererEndpoints       []string          `yaml:"ordererEndpoints"`
	OrderersToJoin         []string          `yaml:"orderersToJoin"`
}

type ExternalOrderer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type OrdererNode struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	TlsCert string `yaml:"tlsCert"`
}

type Identity struct {
	SecretKey       string `yaml:"secretKey"`
	SecretName      string `yaml:"secretName"`
	SecretNamespace string `yaml:"secretNamespace"`
}

type AdminOrganization struct {
	MspID string `yaml:"mspID"`
}

type FabricMainChannelData struct {
	ChannelName               string                    `yaml:"channelName"`
	AdminOrdererOrganizations []AdminOrganization       `yaml:"adminOrdererOrganizations"`
	AdminPeerOrganizations    []AdminOrganization       `yaml:"adminPeerOrganizations"`
	PeerOrganizations         []PeerOrganization        `yaml:"peerOrganizations"`
	OrdererOrganizations      []OrdererOrganization     `yaml:"ordererOrganizations"`
	Orderers                  []OrdererNode             `yaml:"orderers"`
	Identities                map[string]Identity       `yaml:"identities"`
}

const fabricMainChannelTemplate = `apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricMainChannel
metadata:
  name: {{.ChannelName}}
spec:
  name: {{.ChannelName}}
  adminOrdererOrganizations:{{range .AdminOrdererOrganizations}}
    - mspID: {{.MspID}}{{end}}
  adminPeerOrganizations:{{range .AdminPeerOrganizations}}
    - mspID: {{.MspID}}{{end}}
  channelConfig:
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
  peerOrganizations:{{range .PeerOrganizations}}
    - mspID: {{.MspID}}
      caName: "{{.CaName}}"
      caNamespace: "{{.CaNamespace}}"{{end}}
  identities:{{range $key, $value := .Identities}}
    {{$key}}:
      secretKey: {{$value.SecretKey}}
      secretName: {{$value.SecretName}}
      secretNamespace: {{$value.SecretNamespace}}{{end}}
  ordererOrganizations:{{range .OrdererOrganizations}}
    - caName: "{{.CaName}}"
      caNamespace: "{{.CaNamespace}}"
      externalOrderersToJoin:{{range .ExternalOrderersToJoin}}
        - host: {{.Host}}
          port: {{.Port}}{{end}}
      mspID: {{.MspID}}
      ordererEndpoints:{{range .OrdererEndpoints}}
        - {{.}}{{end}}
      orderersToJoin: []{{end}}
  orderers:{{range .Orderers}}
    - host: {{.Host}}
      port: {{.Port}}
      tlsCert: |-
        {{.TlsCert}}{{end}}
`

func CreateMainChannel(configFile string) error {
	fmt.Println("Iniciando criação do Main Channel...")
	
	file, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo de configuração: %v", err)
	}

	var config struct {
		CA       []fabric.CA      `json:"CA"`
		Peers    []fabric.Peer    `json:"Peers"`
		Orderers []fabric.Orderer `json:"Orderer"`
		Channels []fabric.Channel `json:"Channel"`
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	fmt.Println("Processando configurações do channel...")

	channelData, err := generateFabricMainChannelData(config)
	if err != nil {
		return fmt.Errorf("erro ao processar dados do channel: %v", err)
	}
	channelData.ChannelName = "demo"

	fmt.Printf("Gerando YAML para o channel '%s'...\n", channelData.ChannelName)

	outputFile := "demo-channel.yaml"
	err = generateYAMLFile(channelData, outputFile)
	if err != nil {
		return fmt.Errorf("erro ao gerar arquivo YAML: %v", err)
	}

	fmt.Printf("Channel YAML '%s' gerado com sucesso!\n", outputFile)

	fmt.Printf("Resumo: %d peer org(s), %d orderer org(s), %d identidades\n", 
		len(channelData.PeerOrganizations), 
		len(channelData.OrdererOrganizations),
		len(channelData.Identities))

	fmt.Println("Aplicando channel no cluster Kubernetes...")
	err = applyChannelYAML(outputFile)
	if err != nil {
		return fmt.Errorf("erro ao aplicar YAML no cluster: %v", err)
	}
	
	fmt.Println("Channel aplicado com sucesso no cluster!")
	return nil
}

func readTLSCertificate(filename string) (string, error) {
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

func generateTlsCertVar(host string) string {
	varName := strings.ToUpper(host)
	varName = strings.ReplaceAll(varName, ".", "_")
	varName = strings.ReplaceAll(varName, "-", "_")
	return varName + "_TLS_CERT"
}

func generateFabricMainChannelData(config struct {
	CA       []fabric.CA      `json:"CA"`
	Peers    []fabric.Peer    `json:"Peers"`
	Orderers []fabric.Orderer `json:"Orderer"`
	Channels []fabric.Channel `json:"Channel"`
}) (FabricMainChannelData, error) {
	data := FabricMainChannelData{
		Identities: make(map[string]Identity),
	}

	ordererMSPs := make(map[string]string)
	peerMSPs := make(map[string]string)

	for _, ca := range config.CA {
		if ca.UserType == "orderer" {
			ordererMSPs[ca.MspID] = ca.Name
		} else if ca.UserType == "peer" {
			peerMSPs[ca.MspID] = ca.Name
		}
	}

	for _, peer := range config.Peers {
		if peer.Mspid != "" && peer.CAName != "" {
			peerMSPs[peer.Mspid] = peer.CAName
		}
	}

	fmt.Printf("Debug - Peer MSPs detectados: %v\n", peerMSPs)
	fmt.Printf("Debug - Orderer MSPs detectados: %v\n", ordererMSPs)

	for mspID := range ordererMSPs {
		data.AdminOrdererOrganizations = append(data.AdminOrdererOrganizations, AdminOrganization{MspID: mspID})
	}

	for mspID := range peerMSPs {
		data.AdminPeerOrganizations = append(data.AdminPeerOrganizations, AdminOrganization{MspID: mspID})
	}

	for mspID, caName := range peerMSPs {
		data.PeerOrganizations = append(data.PeerOrganizations, PeerOrganization{
			MspID:       mspID,
			CaName:      caName,
			CaNamespace: "default",
		})
	}

	ordererOrgMap := make(map[string]*OrdererOrganization)
	
	for _, orderer := range config.Orderers {
		if existingOrg, exists := ordererOrgMap[orderer.Mspid]; exists {
			externalOrderer := ExternalOrderer{
				Host: fmt.Sprintf("%s.default", orderer.Name),
				Port: 7053,
			}

			found := false
			for _, existing := range existingOrg.ExternalOrderersToJoin {
				if existing.Host == externalOrderer.Host {
					found = true
					break
				}
			}
			if !found {
				existingOrg.ExternalOrderersToJoin = append(existingOrg.ExternalOrderersToJoin, externalOrderer)
			}
		} else {
			var ordererEndpoints []string
			var externalOrderers []ExternalOrderer
			
			for _, channel := range config.Channels {
				if channel.MspID == orderer.Mspid {
					ordererEndpoints = append(ordererEndpoints, channel.OrdererNodeEndpoint...)
				}
			}
			externalOrderers = append(externalOrderers, ExternalOrderer{
				Host: fmt.Sprintf("%s.default", orderer.Name),
				Port: 7053,
			})

			ordererOrgMap[orderer.Mspid] = &OrdererOrganization{
				CaName:                 orderer.CAName,
				CaNamespace:            "default",
				ExternalOrderersToJoin: externalOrderers,
				MspID:                  orderer.Mspid,
				OrdererEndpoints:       ordererEndpoints,
				OrderersToJoin:         []string{},
			}
		}
	}

	for _, org := range ordererOrgMap {
		data.OrdererOrganizations = append(data.OrdererOrganizations, *org)
	}

	for _, orderer := range config.Orderers {
		host := orderer.Hosts
		port := 443

		certFileName := fmt.Sprintf("/tmp/%s-cert.pem", orderer.Name)
		tlsCert, err := readTLSCertificate(certFileName)
		if err != nil {
			return FabricMainChannelData{}, fmt.Errorf("erro ao ler certificado TLS do orderer %s: %v", orderer.Name, err)
		}

		data.Orderers = append(data.Orderers, OrdererNode{
			Host:    host,
			Port:    port,
			TlsCert: tlsCert,
		})
	}

	for _, channel := range config.Channels {
		if channel.FileOutput != "" {
			data.Identities[channel.MspID] = Identity{
				SecretKey:       channel.FileOutput,
				SecretName:      "wallet",
				SecretNamespace: channel.Namespace,
			}
		}

		if channel.FileOutputTls != "" {
			data.Identities[channel.MspID+"-tls"] = Identity{
				SecretKey:       channel.FileOutputTls,
				SecretName:      "wallet",
				SecretNamespace: channel.Namespace,
			}
		}

		if channel.MspID == "OrdererMSP" {
			data.Identities[channel.MspID+"-sign"] = Identity{
				SecretKey:       "orderermspsign.yaml",
				SecretName:      "wallet",
				SecretNamespace: channel.Namespace,
			}
		}
	}

	return data, nil
}

func applyChannelYAML(filename string) error {
	cmd := exec.Command("kubectl", "apply", "-f", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao executar kubectl apply: %v", err)
	}

	return nil
}

func generateYAMLFile(data FabricMainChannelData, filename string) error {
	tmpl, err := template.New("fabricMainChannel").Parse(fabricMainChannelTemplate)
	if err != nil {
		return fmt.Errorf("erro ao criar template: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("erro ao executar template: %v", err)
	}

	return nil
}
