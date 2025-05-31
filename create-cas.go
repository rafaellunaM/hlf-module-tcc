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
	caImage := os.Getenv("CA_IMAGE")
	caVersion := os.Getenv("CA_VERSION")
	storageClass := os.Getenv("SC_NAME")

	if caImage == "" || caVersion == "" || storageClass == "" {
		fmt.Println("Erro: CA_IMAGE, CA_VERSION e SC_NAME devem estar definidas nas variáveis de ambiente.")
		os.Exit(1)
	}

	file, err := os.ReadFile("hlf-config.json")
	if err != nil {
		log.Fatalf("❌ Erro ao ler o JSON: %v", err)
	}

	var partialConfig struct {
		CA []fabric.CA `json:"CA"`
	}
	
	if err := json.Unmarshal(file, &partialConfig); err != nil {
		log.Fatalf("❌ Erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, ca := range partialConfig.CA {
		fmt.Printf("🔧 Criando a CA %s...\n", ca.Name)
		cmd := exec.Command("kubectl", "hlf", "ca", "create",
			"--image="+caImage,
			"--version="+caVersion,
			"--storage-class="+storageClass,
			"--capacity="+ca.Capacity,
			"--name="+ca.Name,
			"--enroll-id="+ca.EnrollId, 
			"--enroll-pw="+ca.EnrollPw,
			"--hosts="+ca.Hosts,
			"--istio-port="+ca.IstioPort,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Erro ao criar a CA %s: %v\n", ca.Name, err)
			os.Exit(1)
		}
		fmt.Printf("✅ CA %s criada com sucesso.\n", ca.Name)
	}
}
