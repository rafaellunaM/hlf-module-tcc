package chaincode

import (
	"fmt"
	"os/exec" 
)

func DeployChaincode (pathFile string) error {
	scriptPath, err := selectChaincode(pathFile)
	if err != nil {
			return fmt.Errorf("erro ao selecionar chaincode: %v", err)
	}

	fmt.Printf("Aplicando chaincode: %s\n", scriptPath)
	cmd := exec.Command("bash", scriptPath)
	
	if err := cmd.Run(); err != nil {
			return fmt.Errorf("erro ao aplicar chaincode: %v", err)
	}
	
	fmt.Println("Chaincode deployado com sucesso!")
	return nil
}

func selectChaincode(configFile string) (string, error) {
	switch configFile {
	case "hlf-config.json":
		fmt.Println("Usando: chain-code")
		return "cmd/chaincode/chain-code.sh", nil
	case "templates/4-orderers.json":
		fmt.Println("Usando: chain-code")
		return "cmd/chaincode/chain-code.sh", nil
	case "templates/4-peers.json":
		fmt.Println("Usando: chain-code-4peers")
		return "cmd/chaincode/chain-code-4peers.sh", nil
	default:
		return "", fmt.Errorf("arquivo de configuração não reconhecido: %s", configFile)
	}
}
