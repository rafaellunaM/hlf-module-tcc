package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func RegisterChannels(configFile string) error {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("erro ao ler o JSON: %v", err)
	}

	var partialConfig struct {
		Channel []fabric.Channel `json:"Channel"`
	}
	
	if err := json.Unmarshal(file, &partialConfig); err != nil {
		return fmt.Errorf("erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, channels := range partialConfig.Channel {
		fmt.Printf("Registrando a CA %s...\n", channels.Name)
		cmd := exec.Command("kubectl", "hlf", "ca", "register",
			"--name="+channels.Name,
			"--user="+ channels.UserAdmin,
			"--secret="+channels.Secretadmin,
			"--type="+channels.UserType,
			"--enroll-id="+channels.EnrollID,
			"--enroll-secret="+channels.EnrollPW,
			"--mspid="+channels.MspID,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				if exitCode == 74 {
					fmt.Printf("Identidade %s já registrada, continuando...\n", channels.UserAdmin)
					continue
				}
				fmt.Printf("Comando retornou código de saída %d\n", exitCode)
				continue
			}
			return fmt.Errorf("erro ao registrar a CA %s: %v", channels.Name, err)
		}
		fmt.Printf("Channel %s registrado com sucesso.\n\n", channels.Name)
	}
	
	return nil
}
