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
		Peers []fabric.Peer `json:"Peers"`
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

	for _, peer := range partialConfig.Peers {
		fmt.Printf("🔐 Registrando peer `%s` na CA `%s`…\n", peer.User, peer.CAName)
		args := []string{
			"hlf", "ca", "register",
			"--name=" + peer.CAName,
			"--user=" + peer.User,
			"--secret=" + peer.Secret,
			"--type=" + peer.UserType,
			"--enroll-id=" + peer.EnrollId,
			"--enroll-secret=" + peer.EnrollPw,
			"--mspid=" + peer.Mspid,
		}
		run(args...)
		fmt.Printf("✅ Peer `%s` registrado.\n\n", peer.User)
	}
}
