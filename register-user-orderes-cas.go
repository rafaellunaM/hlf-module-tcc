package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func main() {
	data, err := os.ReadFile("hlf-config.json")
	if err != nil {
		log.Fatalf("❌ não consegui ler hlf-config.json: %v", err)
	}

	var partialConfig struct {
		Orderer []fabric.Orderer `json:"Orderer"`
	}

	if err := json.Unmarshal(data, &partialConfig); err != nil {
		log.Fatalf("❌ não consegui parsear JSON: %v", err)
	}

	run := func(args ...string) {
		fmt.Printf("🔧 Executando: kubectl %v\n", args) 
		cmd := exec.Command("kubectl", args...)   
		cmd.Stdout = os.Stdout                       
		cmd.Stderr = os.Stderr                      
		if err := cmd.Run(); err != nil {
			log.Fatalf("❌ erro ao executar %v: %v", args, err)
		}
	}

	for _, orderer := range partialConfig.Orderer {
		fmt.Printf("🔐 Registrando Orderer `%s` na CA `%s`…\n", orderer.User, orderer.CAName)
		args := []string{
			"hlf", "ca", "register",
			"--name=" + orderer.CAName,
			"--user=" + orderer.User,
			"--secret=" + orderer.Secret,
			"--type=" + orderer.UserType,
			"--enroll-id=" + orderer.EnrollID,
			"--enroll-secret=" + orderer.EnrollPW,
			"--mspid=" + orderer.Mspid,
			"--ca-url=" + orderer.CaURL,
		}
		run(args...)
		fmt.Printf("✅ Orderer `%s` registrado.\n\n", orderer.User)
	}
}
