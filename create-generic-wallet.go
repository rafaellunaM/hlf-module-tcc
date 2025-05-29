package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SecretChannel struct {
	FileOutput    string `json:"fileOutput"`
	FileOutputTls string `json:"fileOutputTls"`
	Namespace     string `json:"namespace"`
}

type FullResources struct {
	Channels []SecretChannel `json:"Channels"`
}

func main() {
	raw, err := os.ReadFile("output.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro lendo output.json: %v\n", err)
		os.Exit(1)
	}

	var config FullResources
	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Erro no unmarshal: %v\n", err)
		os.Exit(1)
	}

	pwd, _ := os.Getwd()
	
	foundFiles := make(map[string]string)

	for _, ch := range config.Channels {
		
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
		fmt.Println("❌ Nenhum arquivo encontrado para criar o Secret.")
		return
	}

	createSecret(foundFiles)
}

func createSecret(files map[string]string) {
	fmt.Printf("🚀 Criando Secret 'wallet' no namespace 'default'...\n")

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
		fmt.Printf("❌ Erro ao criar Secret: %v\n\n", err)
		return
	}

	fmt.Printf("✅ Secret 'wallet' criado com sucesso no namespace 'default'.\n")
	fmt.Printf("📁 Arquivos incluídos:\n")
	for fileName, fullPath := range files {
		fmt.Printf("   - %s (%s)\n", fileName, fullPath)
	}
	fmt.Println()
}