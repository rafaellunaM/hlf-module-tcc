package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func RegisterOrderers(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("não consegui ler %s: %v", configFile, err)
	}

	var partialConfig struct {
		Orderer []fabric.Orderer `json:"Orderer"`
	}

	if err := json.Unmarshal(data, &partialConfig); err != nil {
		return fmt.Errorf("não consegui parsear JSON: %v", err)
	}

	run := func(args ...string) error {
		fmt.Printf("Executando: kubectl %v\n", args) 
		cmd := exec.Command("kubectl", args...)   
		cmd.Stdout = os.Stdout                       
		cmd.Stderr = os.Stderr                      
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao executar %v: %v", args, err)
		}
		return nil
	}

	for _, orderer := range partialConfig.Orderer {
		fmt.Printf("Registrando Orderer `%s` na CA `%s`…\n", orderer.User, orderer.CAName)
		args := []string{
			"hlf", "ca", "register",
			"--name=" + orderer.CAName,
			"--user=" + orderer.User,
			"--secret=" + orderer.Secret,
			"--type=" + orderer.UserType,
			"--enroll-id=" + orderer.EnrollID,
			"--enroll-secret=" + orderer.EnrollPW,
			"--mspid=" + orderer.Mspid,
			"--ca-url=" + orderer.CaURL,
		}
		
		if err := run(args...); err != nil {
			return err
		}
		fmt.Printf("Orderer `%s` registrado.\n\n", orderer.User)
	}
	
	return nil
}
