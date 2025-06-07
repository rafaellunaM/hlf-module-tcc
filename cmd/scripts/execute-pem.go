package scripts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"hlf/internal/fabric"
)

func ExecutePemScript(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	var partialConfig struct {
		Orderer []fabric.Orderer `json:"Orderer"`
	}

	if err := json.Unmarshal(data, &partialConfig); err != nil {
		return fmt.Errorf("erro ao parsear JSON: %v", err)
	}

	for _, orderer := range partialConfig.Orderer {
		fmt.Printf("Extraindo certificado PEM para %s...\n", orderer.Name)

		cmd := exec.Command("kubectl", "get", "fabricorderernodes", orderer.Name, "-o=jsonpath={.status.tlsCert}")

		outputFileName := fmt.Sprintf("/tmp/%s-cert.pem", orderer.Name)
		
		outFile, err := os.Create(outputFileName)
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de saída para %s: %v", orderer.Name, err)
		}

		cmd.Stdout = outFile
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			outFile.Close()
			return fmt.Errorf("erro ao extrair certificado PEM para %s: %v", orderer.Name, err)
		}

		outFile.Close()
		fmt.Printf("Certificado PEM extraído para %s\n", outputFileName)
	}

	fmt.Printf("Processados %d orderers com sucesso!\n", len(partialConfig.Orderer))
	return nil
}
