package main

import (
	
	"encoding/json"
	"fmt"
	"regexp"
	"os"
	"hlf/internal/fabric"
	// "os/exec"
	// "strings"
	// "encoding/base64"
	// "io/ioutil"
	// "log"
)



func readJson(filename string) ([]fabric.Channel, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
	return nil, fmt.Errorf("erro ao ler Channel JSON: %w", err)
	}
	var partialConfig struct {
	Channel []fabric.Channel `json:"Channel"`
	}

	if err := json.Unmarshal(file, &partialConfig); err != nil {
	return nil, fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
	}
return partialConfig.Channel, nil
}

func main () {

	channels, _:= readJson("output.json")
	r, _ := regexp.Compile(`org(\d+)-[a-z]+`)
	
	for _, channel := range channels {
		if r.MatchString(channel.Name) {
			fmt.Printf("Encontrei: \n")
			fmt.Printf("  MspID: %s\n", channel.MspID)
			fmt.Printf("  FileOutput: %s\n", channel.FileOutput)
			fmt.Printf("  Namespace: %s\n", channel.Namespace)
		} else {
			fmt.Printf("  FileOutputTls: %s\n", channel.FileOutputTls)
		}
		// fmt.Printf("  Name: %s\n", channel.Name)
		// fmt.Printf("Canal %d: %+v\n", i, channel)
    // 
    // fmt.Printf("  UserAdmin: %s\n", channel.UserAdmin)
    // fmt.Printf("  FileOutput: %s\n", channel.FileOutput)
    // fmt.Printf("  FileOutputTls: %s\n", channel.FileOutputTls)
    // fmt.Printf("  Namespace: %s\n", channel.Namespace)
	}

	// cmd := exec.Command("kubectl", "hlf", "channelcrd", "main", "create",
  // "--name=demo",
  // "--channel-name=demo",
  // "--secret-name=wallet",
  // "--admin-orderer-orgs=OrdererMSP",
  // "--orderer-orgs=OrdererMSP",
  // "--identities="+"OrdererMSP;orderermsp.yaml",
  // "--identities="+"OrdererMSP-sign;orderermspsign.yaml",

	// "--admin-peer-orgs=" + Org1MSP,
  // "--peer-orgs=" +  Org1MSP,
  // "--identities=" +  "Org1MSP;org1msp.yaml",
  // "--secret-ns=" +  default,
  // "--consenters=" +  "orderer0-ord.localho.st:443",
  // "--consenter-certificates=" + "orderermsp.yaml" , // pegar do secret esse valor
	// )
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
}

// func findOrgs () {

// 	channels, _:= readJson("output.json")
// 	// for i, channel := range channels {
// 	// 		fmt.Printf("Canal %d: %+v\n", i, channel)
// 	// }

// 	for _, channels := range partialConfig.Channel {

// 	}
// }