package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Channels struct {
	Name       		string `json:"name"`
	UserAdmin 		string `json:"userAdmin"`
	Secretadmin 	string `json:"secretadmin"`
	Type 			 		string `json:"userType"`
	EnrollID   		string `json:"enrollId"`
	EnrollPW   		string `json:"enrollPw"`
	MPSID      		string `json:"hosts"`
}

type FullResources struct {
	Channels []Channels `json:"Channels"`
}

func main() {

	file, err := os.ReadFile("output.json")
	if err != nil {
			log.Fatalf("❌ Erro ao ler o JSON: %v", err)
	}

	var config FullResources
	if err := json.Unmarshal(file, &config); err != nil {
			log.Fatalf("❌ Erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, channels := range config.Channels {
		fmt.Printf("🔧 Registrando a CA %s...\n", channels.Name)
		cmd := exec.Command("kubectl", "hlf", "ca", "register",
				"--name="+channels.Name,
				"--user="+ channels.UserAdmin,
				"--secret="+channels.Secretadmin,
				"--type="+channels.Type,
				"--enroll-id="+channels.EnrollID,
				"--enroll-secret="+channels.EnrollPW,
				"--mspid="+channels.MPSID,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode := exitErr.ExitCode()
						fmt.Printf("⚠️ Comando retornou código de saída %d\n", exitCode)
						continue
					if exitCode == 74 {
						fmt.Printf("⚠️ Identidade %s já registrada, continuando...\n", channels.UserAdmin)
						continue
					}
				}
			fmt.Printf("❌ Erro ao registrar a CA %s: %v\n", channels.Name, err)
			os.Exit(1)
		}
	}
}
