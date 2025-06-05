package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func runCommand(name string, args ...string) error {
	fmt.Printf("⏳ Executando: %s %v\n", name, args)
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandOutput(name string, args ...string) (string, error) {
	fmt.Printf("⏳ Executando: %s %v\n", name, args)
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func getNextSequence(config, user, peer, channel string) (int, error) {
	out, err := runCommandOutput("kubectl", "hlf", "chaincode", "querycommitted",
		"--config="+config,
		"--user="+user,
		"--peer="+peer,
		"--channel="+channel,
	)
	if err != nil {
		fmt.Printf("❌ Erro ao consultar definição comprometida: %v\nOutput: %s\n", err, out)
		return 1, nil // Assume 1 se nada estiver comprometido ainda
	}

	re := regexp.MustCompile(`Sequence: (\d+)`)
	matches := re.FindStringSubmatch(out)
	if len(matches) == 2 {
		seq, err := strconv.Atoi(matches[1])
		if err != nil {
			return 1, err
		}
		return seq + 1, nil
	}

	return 1, nil // Nenhum encontrado
}

func main() {
	packageID := os.Getenv("PACKAGE_ID")
	if packageID == "" {
		fmt.Println("ERRO: Variável PACKAGE_ID não definida.")
		os.Exit(1)
	}

	version := "1.0"
	channel := "demo"
	chaincodeName := "asset"
	user := "admin"
	peer := "org1-peer0.default"
	config := "org1.yaml"
	mspID := "Org1MSP"
	policy := "OR('Org1MSP.member')"

	// 🔁 Descobrir próxima sequência
	sequenceInt, err := getNextSequence(config, user, peer, channel)
	if err != nil {
		fmt.Printf("❌ Erro ao obter próxima sequência: %v\n", err)
		os.Exit(1)
	}
	sequence := fmt.Sprintf("%d", sequenceInt)

	// ✅ 1. Approve chaincode
	err = runCommand("kubectl", "hlf", "chaincode", "approveformyorg",
		"--config="+config,
		"--user="+user,
		"--peer="+peer,
		"--package-id="+packageID,
		"--version="+version,
		"--sequence="+sequence,
		"--name="+chaincodeName,
		"--policy="+policy,
		"--channel="+channel,
	)
	if err != nil {
		fmt.Printf("❌ Erro ao aprovar chaincode: %v\n", err)
		os.Exit(1)
	}

	// ✅ 2. Commit chaincode
	err = runCommand("kubectl", "hlf", "chaincode", "commit",
		"--config="+config,
		"--user="+user,
		"--mspid="+mspID,
		"--version="+version,
		"--sequence="+sequence,
		"--name="+chaincodeName,
		"--policy="+policy,
		"--channel="+channel,
	)
	if err != nil {
		fmt.Printf("❌ Erro ao commitar chaincode: %v\n", err)
		os.Exit(1)
	}

	// ✅ 3. Invoke initLedger
	err = runCommand("kubectl", "hlf", "chaincode", "invoke",
		"--config="+config,
		"--user="+user,
		"--peer="+peer,
		"--chaincode="+chaincodeName,
		"--channel="+channel,
		"--fcn=initLedger",
		"-a", "[]",
	)
	if err != nil {
		fmt.Printf("❌ Erro ao invocar chaincode: %v\n", err)
		os.Exit(1)
	}

	// ✅ 4. Query GetAllAssets
	err = runCommand("kubectl", "hlf", "chaincode", "query",
		"--config="+config,
		"--user="+user,
		"--peer="+peer,
		"--chaincode="+chaincodeName,
		"--channel="+channel,
		"--fcn=GetAllAssets",
		"-a", "[]",
	)
	if err != nil {
		fmt.Printf("❌ Erro ao consultar chaincode: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Chaincode aprovado, commitado, invocado e consultado com sucesso.")
}
