package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"hlf/internal/fabric"
)

func CreateWallet(configFile string) error {
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

	pwd, _ := os.Getwd()
	
	foundFiles := make(map[string]string)

	for _, ch := range partialConfig.Channel {
		
		for _, fname := range []string{ch.FileOutput, ch.FileOutputTls} {
			if fname == "" {
				continue
			}
			
			fullPath := filepath.Join(pwd, fname)
			if _, err := os.Stat(fullPath); err == nil {
				
				fileName := filepath.Base(fname)
				foundFiles[fileName] = fullPath
			}
		}
	}

	if len(foundFiles) == 0 {
		return fmt.Errorf("nenhum arquivo encontrado para criar o Secret")
	}

	return createSecret(foundFiles)
}

func createSecret(files map[string]string) error {
	fmt.Printf("Criando Secret 'wallet' no namespace 'default'...\n")

	args := []string{
		"create", "secret", "generic", "wallet",
		"--namespace=default",
	}

	for fileName, fullPath := range files {
		args = append(args, fmt.Sprintf("--from-file=%s=%s", fileName, fullPath))
	}

	fmt.Printf("📋 Comando: kubectl %s\n", strings.Join(args, " "))

	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao criar Secret: %v", err)
	}

	fmt.Printf("Secret 'wallet' criado com sucesso no namespace 'default'.\n")
	fmt.Printf("Arquivos incluídos:\n")
	for fileName, fullPath := range files {
		fmt.Printf("   - %s (%s)\n", fileName, fullPath)
	}
	fmt.Println()
	
	return nil
}
