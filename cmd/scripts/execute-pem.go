package scripts

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecutePemScript() error {
	fmt.Printf("Extraindo certificado PEM do orderermsp.yaml...\n")

	cmd := exec.Command("kubectl", "get", "fabricorderernodes", "ord-node1", "-o=jsonpath={.status.tlsCert}")
	
	outFile, err := os.Create("/tmp/orderer-cert.pem")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de saída: %v", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao extrair certificado PEM: %v", err)
	}

	fmt.Printf("Certificado PEM extraído para /tmp/orderer-cert.pem\n")
	return nil
}
