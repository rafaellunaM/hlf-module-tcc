package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"hlf/internal/fabric"
	"log"
)

func CreateCAs(configFile string) error {
	caImage := os.Getenv("CA_IMAGE")
	caVersion := os.Getenv("CA_VERSION")
	storageClass := os.Getenv("SC_NAME")

	if caImage == "" || caVersion == "" || storageClass == "" {
		return fmt.Errorf("CA_IMAGE, CA_VERSION e SC_NAME devem estar definidas nas variáveis de ambiente")
	}

	file, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("erro ao ler o JSON: %v", err)
	}

	var partialConfig struct {
		CA []fabric.CA `json:"CA"`
	}
	
	if err := json.Unmarshal(file, &partialConfig); err != nil {
		return fmt.Errorf("erro ao fazer unmarshal do JSON: %v", err)
	}

	for _, ca := range partialConfig.CA {
		fmt.Printf("Criando a CA %s...\n", ca.Name)
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
			return fmt.Errorf("erro ao criar a CA %s: %v", ca.Name, err)
		}
		fmt.Printf(" CA %s criada com sucesso.\n", ca.Name)
	}
	fmt.Println("⏳ Aguardando todos os CAs nodes ficarem em estado Running...")
	waitCmd := exec.Command("kubectl", "wait",
			"--timeout=60s",
			"--for=condition=Running",
			"fabriccas.hlf.kungfusoftware.es",
			"--all",
	)
	waitCmd.Stdout = os.Stdout
	waitCmd.Stderr = os.Stderr
	if err := waitCmd.Run(); err != nil {
		log.Fatalf("Erro ao aguardar orderers: %v", err)
	}
	fmt.Println(" Todos os CAs estão em execução.")
	
	return nil
}
