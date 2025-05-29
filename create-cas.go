package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type CAConfig struct {
	Capacity   string `json:"capacity"`
	Name       string `json:"name"`
	EnrollID   string `json:"enrollId"`
	EnrollPW   string `json:"enrollPw"`
	Hosts      string `json:"hosts"`
	IstioPort  string `json:"istioPort"`
}

type FullResources struct {
	CAs []CAConfig `json:"CAs"`
}

func main() {
	caImage := os.Getenv("CA_IMAGE")
	caVersion := os.Getenv("CA_VERSION")
	storageClass := os.Getenv("SC_NAME")

	if caImage == "" || caVersion == "" || storageClass == "" {
			fmt.Println("Erro: CA_IMAGE, CA_VERSION e SC_NAME devem estar definidas nas variáveis de ambiente.")
			os.Exit(1)
	}

	file, err := os.ReadFile("output.json")
	if err != nil {
			log.Fatalf("❌ Erro ao ler o JSON: %v", err)
	}

	var config FullResources
	if err := json.Unmarshal(file, &config); err != nil {
			log.Fatalf("❌ Erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, ca := range config.CAs {
			fmt.Printf("🔧 Criando a CA %s...\n", ca.Name)
			cmd := exec.Command("kubectl", "hlf", "ca", "create",
					"--image="+caImage,
					"--version="+caVersion,
					"--storage-class="+storageClass,
					"--capacity="+ca.Capacity,
					"--name="+ca.Name,
					"--enroll-id="+ca.EnrollID,
					"--enroll-pw="+ca.EnrollPW,
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
