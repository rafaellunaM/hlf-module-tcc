package ca

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	runWithErrorCheck := func(ordererUser string, args ...string) error {
		fmt.Printf("Executando: kubectl %v\n", args) 
		cmd := exec.Command("kubectl", args...)
		
		// Captura stderr para verificar a mensagem de erro
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			errorMsg := string(output)
			// Verifica se o erro é sobre identidade já registrada
			if strings.Contains(errorMsg, "is already registered") {
				fmt.Printf("Orderer `%s` já está registrado, continuando...\n", ordererUser)
				return nil
			}
			// Se for outro tipo de erro, retorna o erro
			fmt.Fprintf(os.Stderr, "%s", output)
			return fmt.Errorf("erro ao executar %v: %v", args, err)
		}
		
		// Se não houve erro, mostra a saída normalmente
		fmt.Print(string(output))
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
		
		if err := runWithErrorCheck(orderer.User, args...); err != nil {
			return err
		}
		fmt.Printf("Orderer `%s` processado com sucesso.\n\n", orderer.User)
	}
	
	return nil
}
