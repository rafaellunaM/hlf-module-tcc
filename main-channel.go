package channels

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"log"
	"hlf/internal/fabric"
)

func main() {

	file, err := os.ReadFile("hlf-config.json")
	if err != nil {
		log.Fatalf("Erro ao ler o JSON: %v", err)
	}

	var partialConfig struct {
		Channel []fabric.Channel `json:"Channel"`
	}

	if err := json.Unmarshal(file, &partialConfig); err != nil {
		log.Fatalf("Erro ao fazer parse do JSON: %v", err)
	}

	peerRegex, _ := regexp.Compile(`org(\d+)-[a-z]+`)
	ordRegex, _ := regexp.Compile(`ord-[a-z]+`)

	var peerOrgs []string
	var peerIdentities []string
	var ordererConsenters []string
	var ordererCerts []string
	var namespace = "default"

	for _, channel := range partialConfig.Channel {
		if peerRegex.MatchString(channel.Name) {
			peerOrgs = append(peerOrgs, channel.MspID)
			peerIdentities = append(peerIdentities, channel.MspID+";"+channel.FileOutput)
		}
		if ordRegex.MatchString(channel.Name) {
			ordererConsenters = append(ordererConsenters, channel.OrdererNodeEndpoint+":443")
			ordererCerts = append(ordererCerts, channel.FileOutput)
		}
	}

	fmt.Printf("Orderer Consenters: %v\n", ordererConsenters)
	fmt.Printf("Peer Orgs: %v\n", peerOrgs)

	args := []string{
		"hlf", "channelcrd", "main", "create",
		"--name", "demo",
		"--channel-name", "demo", 
		"--secret-name", "wallet",
		"--admin-orderer-orgs", "OrdererMSP",
		"--orderer-orgs", "OrdererMSP",
		"--consenter-certificates", "/tmp/orderer-cert.pem",
		"--identities", "OrdererMSP;orderermsp.yaml",
		"--identities", "OrdererMSP-sign;orderermspsign.yaml",
		"--admin-peer-orgs", strings.Join(peerOrgs, ","),
		"--peer-orgs", strings.Join(peerOrgs, ","),
		"--secret-ns", namespace,
		"--consenters", strings.Join(ordererConsenters, ","),
	}

	for _, identity := range peerIdentities {
		args = append(args, "--identities", identity)
	}

	for _, channel := range partialConfig.Channel {
		if peerRegex.MatchString(channel.Name) {
			args = append(args, "--identities", channel.MspID+"-tls;"+channel.FileOutputTls)
		}
	}

	fmt.Printf("Comando: %s\n", strings.Join(args, " "))

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Erro ao criar o channel: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("channel demo criado com sucesso\n")
}
