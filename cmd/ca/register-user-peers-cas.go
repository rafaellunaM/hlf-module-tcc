package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func RegisterPeers(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("não consegui ler %s: %v", configFile, err)
	}

	var partialConfig struct {
		Peers []fabric.Peer `json:"Peers"`
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

	for _, peer := range partialConfig.Peers {
		fmt.Printf("Registrando peer `%s` na CA `%s`…\n", peer.User, peer.CAName)
		args := []string{
			"hlf", "ca", "register",
			"--name=" + peer.CAName,
			"--user=" + peer.User,
			"--secret=" + peer.Secret,
			"--type=" + peer.UserType,
			"--enroll-id=" + peer.EnrollId,
			"--enroll-secret=" + peer.EnrollPw,
			"--mspid=" + peer.Mspid,
		}
		
		if err := run(args...); err != nil {
			return err
		}
		fmt.Printf(" Peer `%s` registrado.\n\n", peer.User)
	}
	
	return nil
}
