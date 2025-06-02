package scripts

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecutePemScript() error {
	fmt.Printf("Extraindo certificado PEM do orderermsp.yaml...\n")

	cmd := exec.Command("sh", "-c", 
		"cat orderermsp.yaml | grep -A 100 'pem: |' | sed 's/.*pem: |//' | sed '/^[[:space:]]*$/d' | sed 's/^[[:space:]]*//' > /tmp/orderer-cert.pem")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao extrair certificado PEM: %v", err)
	}

	fmt.Printf("Certificado PEM extraído para /tmp/orderer-cert.pem\n")
	return nil
}
