package administration

import (
	"fmt"
	"os/exec"
)

func DeleteAllResources() error {
	fmt.Printf("Deletando recursos\n")

	resources := []string{
		"fabricorderernodes.hlf.kungfusoftware.es",
		"fabricpeers.hlf.kungfusoftware.es",
		"fabriccas.hlf.kungfusoftware.es",
		"fabricchaincode.hlf.kungfusoftware.es",
		"fabricmainchannels",
		"fabricfollowerchannels",
		"secret/wallet",
	}

	for _, resource := range resources {
		var output []byte
		var err error
		
		if resource != "secret/wallet" {
			output, err = exec.Command("kubectl", "delete", resource, "--all", "--all-namespaces").CombinedOutput()
		} else {
			output, err = exec.Command("kubectl", "delete", "secret", "wallet", "-n", "default").CombinedOutput()
		}
		
		if err != nil {
			return fmt.Errorf("Erro ao deletar recurso %s: %v\nSaída: %s", resource, err, string(output))
		}

		fmt.Printf("Recurso %s deletado com sucesso\n", resource)
	}

	fmt.Printf("Todos os recursos deletados com sucesso\n")
	return nil
}
