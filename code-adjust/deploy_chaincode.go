package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCmd(name string, args ...string) error {
	fmt.Printf(" Executando: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCmdOutput(name string, args ...string) (string, error) {
	fmt.Printf(" Executando: %s %s\n", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro: %v, output: %s", err, string(out))
	}
	return string(out), nil
}

func main() {
	const (
		namespace       = "default"
		orgMSP          = "Org1MSP"
		ordererMSP      = "OrdererMSP"
		caName          = "org1-ca"
		caUser          = "admin"
		caUserSecret    = "adminpw"
		enrollID        = "enroll"
		enrollSecret    = "enrollpw"
		chaincodeName   = "asset"
		chaincodeLabel  = "asset"
		peerName        = "org1-peer0.default"
		chaincodeLang   = "golang"
		externalImage   = "kfsoftware/chaincode-external:latest"
	)

	if err := runCmd("kubectl", "hlf", "inspect", "--output", "org1.yaml", "-o", orgMSP, "-o", ordererMSP); err != nil {
		fmt.Println(" Erro no inspect")
		os.Exit(1)
	}

	if err := runCmd("kubectl", "hlf", "ca", "enroll",
		"--name="+caName,
		"--user="+caUser,
		"--secret="+caUserSecret,
		"--mspid="+orgMSP,
		"--ca-name=ca",
		"--output", "peer-org1.yaml",
	); err != nil {
		fmt.Println(" Erro ao fazer enroll")
		os.Exit(1)
	}

	if err := runCmd("kubectl", "hlf", "utils", "adduser",
		"--userPath=peer-org1.yaml",
		"--config=org1.yaml",
		"--username="+caUser,
		"--mspid="+orgMSP,
	); err != nil {
		fmt.Println(" Erro ao adicionar usuário na connection string")
		os.Exit(1)
	}

	os.Remove("code.tar.gz")
	os.Remove("chaincode.tgz")

	metadataJSON := fmt.Sprintf(`{
  "type": "ccaas",
  "label": "%s"
}
`, chaincodeLabel)
	if err := os.WriteFile("metadata.json", []byte(metadataJSON), 0644); err != nil {
		fmt.Printf(" Erro ao criar metadata.json: %v\n", err)
		os.Exit(1)
	}

	connectionJSON := fmt.Sprintf(`{
  "address": "%s:7052",
  "dial_timeout": "10s",
  "tls_required": false
}
`, chaincodeName)
	if err := os.WriteFile("connection.json", []byte(connectionJSON), 0644); err != nil {
		fmt.Printf(" Erro ao criar connection.json: %v\n", err)
		os.Exit(1)
	}

	if err := runCmd("tar", "cfz", "code.tar.gz", "connection.json"); err != nil {
		fmt.Println(" Erro ao criar code.tar.gz")
		os.Exit(1)
	}

	if err := runCmd("tar", "cfz", "chaincode.tgz", "metadata.json", "code.tar.gz"); err != nil {
		fmt.Println(" Erro ao criar chaincode.tgz")
		os.Exit(1)
	}
	out, err := runCmdOutput("kubectl", "hlf", "chaincode", "calculatepackageid",
		"--path=chaincode.tgz",
		"--language=node",
		"--label="+chaincodeLabel,
	)
	if err != nil {
		fmt.Printf(" Erro ao calcular package ID: %v\n", err)
		os.Exit(1)
	}
	packageID := strings.TrimSpace(out)
	fmt.Printf("📦 PACKAGE_ID=%s\n", packageID)

	if err := runCmd("kubectl", "hlf", "chaincode", "install",
		"--path=./chaincode.tgz",
		"--config=org1.yaml",
		"--language="+chaincodeLang,
		"--label="+chaincodeLabel,
		"--user="+caUser,
		"--peer="+peerName,
	); err != nil {
		fmt.Println(" Erro ao instalar chaincode")
		os.Exit(1)
	}

	if err := runCmd("kubectl", "hlf", "externalchaincode", "sync",
		"--image="+externalImage,
		"--name="+chaincodeName,
		"--namespace="+namespace,
		"--package-id="+packageID,
		"--tls-required=false",
		"--replicas=1",
	); err != nil {
		fmt.Println(" Erro ao sincronizar chaincode externo")
		os.Exit(1)
	}

	if err := runCmd("kubectl", "hlf", "chaincode", "queryinstalled",
		"--config=org1.yaml",
		"--user="+caUser,
		"--peer="+peerName,
	); err != nil {
		fmt.Println(" Erro ao consultar chaincodes instalados")
		os.Exit(1)
	}

	fmt.Println("Processo completo executado com sucesso!")
}
