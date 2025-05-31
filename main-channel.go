package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	
	"hlf/internal/fabric"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("❌ Uso: go run main.go <config.json> [channel-name]")
	}

	configFile := os.Args[1]
	channelName := "demo"
	
	if len(os.Args) > 2 {
		channelName = os.Args[2]
	}

	// Lê o arquivo de configuração
	fmt.Printf("📖 Lendo configuração de %s...\n", configFile)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("❌ Erro ao ler arquivo: %v", err)
	}

	var config fabric.Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("❌ Erro ao fazer parse do JSON: %v", err)
	}

	fmt.Printf("✅ Configuração carregada com sucesso!\n")
	fmt.Printf("📋 Canal: %s\n", channelName)

	// Extrai certificados dos orderers
	fmt.Println("🔍 Extraindo certificados dos orderers...")
	certificates := extractOrdererCertificates(&config)
	fmt.Printf("✅ %d certificado(s) extraído(s)\n", len(certificates))

	// Monta o comando kubectl hlf channelcrd
	args := buildChannelCommand(&config, channelName, certificates)
	
	// Exibe o comando que será executado
	fmt.Println("\n🔧 Comando a ser executado:")
	fmt.Printf("kubectl %s\n", strings.Join(args, " \\\n  "))
	
	// Pergunta se deve executar
	fmt.Print("\n❓ Executar este comando? (s/n): ")
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) == "s" || strings.ToLower(response) == "sim" {
		fmt.Println("\n🚀 Executando comando...")
		
		cmd := exec.Command("kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				fmt.Printf("⚠️  Comando retornou código de saída %d\n", exitCode)
				if exitCode == 74 {
					fmt.Printf("⚠️  Canal %s já existe, continuando...\n", channelName)
					os.Exit(0)
				}
			}
			log.Fatalf("❌ Erro ao executar comando: %v", err)
		}
		
		fmt.Println("✅ Canal criado com sucesso!")
	} else {
		fmt.Println("ℹ️  Operação cancelada pelo usuário")
	}
}

func extractOrdererCertificates(config *fabric.Config) []string {
	var certificates []string
	
	for _, channel := range config.Channels {
		if !strings.Contains(strings.ToLower(channel.MspID), "orderer") {
			continue
		}
		
		// Busca o certificado no secret do Kubernetes
		cmd := exec.Command("kubectl", "get", "secret", "wallet", 
			"-n", channel.Namespace,
			"-o", fmt.Sprintf("jsonpath={.data.%s}", channel.FileOutputTls))
		
		base64Output, err := cmd.Output()
		if err != nil {
			log.Printf("⚠️  Não foi possível obter secret field %s: %v", channel.FileOutputTls, err)
			continue
		}
		
		// Decodifica o base64
		decodedData, err := base64.StdEncoding.DecodeString(string(base64Output))
		if err != nil {
			log.Printf("⚠️  Erro ao decodificar base64: %v", err)
			continue
		}
		
		// Extrai o conteúdo PEM
		cert := extractPEMContent(string(decodedData))
		if cert != "" {
			certificates = append(certificates, cert)
			fmt.Printf("  ✅ Certificado extraído de %s\n", channel.FileOutputTls)
		}
	}
	
	return certificates
}

func extractPEMContent(data string) string {
	lines := strings.Split(data, "\n")
	var pemLines []string
	foundPem := false
	
	for _, line := range lines {
		if strings.Contains(line, "pem: |") {
			foundPem = true
			// Pega tudo depois de "pem: |"
			parts := strings.Split(line, "pem: |")
			if len(parts) > 1 {
				trimmed := strings.TrimSpace(parts[1])
				if trimmed != "" {
					pemLines = append(pemLines, trimmed)
				}
			}
			continue
		}
		
		if foundPem {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				pemLines = append(pemLines, trimmed)
			}
		}
	}
	
	return strings.Join(pemLines, "\n")
}

func buildChannelCommand(config *fabric.Config, channelName string, certificates []string) []string {
	args := []string{"hlf", "channelcrd", "main", "create"}
	
	// Configurações básicas
	args = append(args, "--name", channelName)
	args = append(args, "--channel-name", channelName)
	
	// Organizações orderer (fixas)
	args = append(args, "--admin-orderer-orgs", "OrdererMSP")
	args = append(args, "--orderer-orgs", "OrdererMSP")
	
	// Identidades orderer (fixas)
	args = append(args, "--identities", "OrdererMSP;orderermsp.yaml")
	args = append(args, "--identities", "OrdererMSP-sign;orderermspsign.yaml")
	
	// Coleta MSP IDs dos peers
	peerMSPs := make(map[string]bool)
	for _, peer := range config.Peers {
		peerMSPs[peer.Mspid] = true
	}
	
	// Adiciona organizações peer
	var peerOrgs []string
	for mspID := range peerMSPs {
		peerOrgs = append(peerOrgs, mspID)
	}
	
	if len(peerOrgs) > 0 {
		args = append(args, "--admin-peer-orgs", strings.Join(peerOrgs, ","))
		args = append(args, "--peer-orgs", strings.Join(peerOrgs, ","))
	}
	
	// Adiciona identidades dos peers
	for _, channel := range config.Channels {
		if strings.Contains(strings.ToLower(channel.MspID), "orderer") {
			continue
		}
		
		// Identidade regular
		args = append(args, "--identities", fmt.Sprintf("%s;%s", channel.MspID, channel.FileOutput))
		
		// Identidade TLS
		if channel.FileOutputTls != "" {
			args = append(args, "--identities", fmt.Sprintf("%s-tls;%s", channel.MspID, channel.FileOutputTls))
		}
	}
	
	// Configuração do secret
	secretNamespace := "default"
	if len(config.Channels) > 0 {
		secretNamespace = config.Channels[0].Namespace
	}
	
	args = append(args, "--secret-name", "wallet")
	args = append(args, "--secret-ns", secretNamespace)
	
	// Adiciona consenters
	for _, orderer := range config.Orderers {
		consenter := fmt.Sprintf("%s:%s", orderer.Hosts, orderer.IstioPort)
		args = append(args, "--consenters", consenter)
	}
	
	// Adiciona certificados
	for _, cert := range certificates {
		args = append(args, "--consenter-certificates", cert)
	}
	
	return args
}