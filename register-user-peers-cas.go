package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type PeerConfig struct {
	CAName   string `json:"CAName"`
	EnrollID string `json:"enrollID"`
	EnrollPW string `json:"enrollPW"`
	User     string `json:"user"`
	Secret   string `json:"secret"`
	UserType string `json:"userType"`
	MSPID    string `json:"mspid"`
}

type FullResources struct {
	Peers []PeerConfig `json:"Peers"`
}

func main() {
	data, err := os.ReadFile("output.json")
	if err != nil {
		log.Fatalf("❌ não consegui ler output.json: %v", err)
	}

	var cfg FullResources
	if err := json.Unmarshal(data, &cfg); err != nil {
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

	for _, p := range cfg.Peers {
		fmt.Printf("🔐 Registrando peer `%s` na CA `%s`…\n", p.User, p.CAName)
		args := []string{
			"hlf", "ca", "register",
			"--name=" + p.CAName,
			"--user=" + p.User,
			"--secret=" + p.Secret,
			"--type=" + p.UserType,
			"--enroll-id=" + p.EnrollID,
			"--enroll-secret=" + p.EnrollPW,
			"--mspid=" + p.MSPID,
		}
		run(args...)
		fmt.Printf("✅ Peer `%s` registrado.\n\n", p.User)
	}
}
