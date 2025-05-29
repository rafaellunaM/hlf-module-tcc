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

	file, err := os.ReadFile("output.json")
	if err != nil {
			log.Fatalf("❌ Erro ao ler o JSON: %v", err)
	}

	var partialConfig struct {
		Channel []fabric.Channel `json:"Channel"`
	}
	
	if err := json.Unmarshal(file, &partialConfig); err != nil {
			log.Fatalf("❌ Erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, channels := range partialConfig.Channel {
		fmt.Printf("🔧 Enroll admin user %s...\n", channels.CaNameTls)
		cmd := exec.Command("kubectl", "hlf", "ca", "enroll",
				"--name=" + channels.Name,
				"--namespace=" + channels.Namespace,
				"--user=" + channels.UserAdmin,
				"--secret=" + channels.Secretadmin,
				"--mspid=" + channels.MspID,
				"--ca-name=" + channels.CaNameTls,
				"--output="+ channels.FileOutput,		
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode := exitErr.ExitCode()
						fmt.Printf("⚠️ Comando retornou código de saída %d\n", exitCode)
						continue
					if exitCode == 74 {
						fmt.Printf("⚠️ Identidade %s já foi feito enroll, continuando...\n", channels.UserAdmin)
						continue
					}
				}
			fmt.Printf("❌ Erro ao fazer enroll do usuário %s: %v\n", channels.Name, err)
			os.Exit(1)
		}
	}

	for _, channels := range partialConfig.Channel {
		fmt.Printf("🔧 Enroll admin user %s...\n", channels.CaNameTls)
		cmd := exec.Command("kubectl", "hlf", "ca", "enroll",
				"--name=" + channels.Name,
				"--namespace=" + channels.Namespace,
				"--user=" + channels.UserAdmin,
				"--secret=" + channels.Secretadmin,
				"--mspid=" + channels.MspID,
				"--ca-name=" + channels.CaNameTls,
				"--output="+ channels.FileOutputTls,		
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
						exitCode := exitErr.ExitCode()
						fmt.Printf("⚠️ Comando retornou código de saída %d\n", exitCode)
						continue
					if exitCode == 74 {
						fmt.Printf("⚠️ Identidade %s já foi feito enroll, continuando...\n", channels.UserAdmin)
						continue
					}
				}
			fmt.Printf("❌ Erro ao fazer enroll do usuário %s: %v\n", channels.Name, err)
			os.Exit(1)
		}
	}
}
