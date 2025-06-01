package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func EnrollChannels(configFile string) error {
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

	for _, channel := range partialConfig.Channel {
		fmt.Printf("Fazendo enroll TLS para %s -> %s...\n", channel.Name, channel.FileOutputTls)
		
		cmd := exec.Command("kubectl", "hlf", "ca", "enroll",
			"--name="+channel.Name,
			"--namespace="+channel.Namespace,
			"--user="+channel.UserAdmin,
			"--secret="+channel.Secretadmin,
			"--mspid="+channel.MspID,
			"--ca-name="+channel.CaNameTls,
			"--output="+channel.FileOutputTls,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				if exitCode == 74 {
					fmt.Printf("Identidade TLS %s já foi feito enroll, continuando...\n", channel.UserAdmin)
					continue
				}
				fmt.Printf("Comando TLS retornou código de saída %d\n", exitCode)
			}
			fmt.Printf("Erro ao fazer enroll TLS do usuário %s: %v\n", channel.Name, err)
			continue
		}
		fmt.Printf(" Enroll TLS concluído para %s -> %s\n", channel.Name, channel.FileOutputTls)
	}

	for _, channel := range partialConfig.Channel {
		fmt.Printf("Fazendo enroll CA (signing) para %s -> %s...\n", channel.Name, channel.FileOutput)
		
		cmd := exec.Command("kubectl", "hlf", "ca", "enroll",
			"--name="+channel.Name,
			"--namespace="+channel.Namespace,
			"--user="+channel.UserAdmin,
			"--secret="+channel.Secretadmin,
			"--mspid="+channel.MspID,
			"--ca-name="+channel.CaName,
			"--output="+channel.FileOutput,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode := exitErr.ExitCode()
				if exitCode == 74 {
					fmt.Printf("Identidade CA %s já foi feito enroll, continuando...\n", channel.UserAdmin)
					continue
				}
				fmt.Printf("Comando CA retornou código de saída %d\n", exitCode)
			}
			fmt.Printf("Erro ao fazer enroll CA do usuário %s: %v\n", channel.Name, err)
			continue
		}
		fmt.Printf(" Enroll CA concluído para %s -> %s\n", channel.Name, channel.FileOutput)
	}

	fmt.Println("Processo de enroll concluído!")
	return nil
}
