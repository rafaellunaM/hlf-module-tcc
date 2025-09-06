package administration

import (
	"fmt"
	"os/exec"
)

func ShowResources() error {
	fmt.Printf("Mostrando recursos do cluster\n")
	fmt.Printf(string(make([]rune, 50)) + "\n")

	resources := []struct {
		name    string
		command []string
	}{
		{"Pods", []string{"kubectl", "get", "pods", "-n", "default"}},
		{"Fabric CAs", []string{"kubectl", "get", "fabriccas", "--all-namespaces"}},
		{"Fabric Peers", []string{"kubectl", "get", "fabricpeers", "--all-namespaces"}},
		{"Fabric Orderers", []string{"kubectl", "get", "fabricorderernodes", "--all-namespaces"}},
		{"Fabric Main Channels", []string{"kubectl", "get", "fabricmainchannels", "--all-namespaces"}},
		{"Fabric Follower Channels", []string{"kubectl", "get", "fabricfollowerchannels", "--all-namespaces"}},
		{"Fabric Chaincode", []string{"kubectl", "get", "fabricchaincode", "--all-namespaces"}},
	}

	for _, resource := range resources {
		fmt.Printf("\n--- %s ---\n", resource.name)
		
		output, err := exec.Command(resource.command[0], resource.command[1:]...).CombinedOutput()
		
		if err != nil {
			fmt.Printf("Nenhum recurso %s encontrado ou erro: %v\n", resource.name, err)
			continue
		}

		if len(output) == 0 {
			fmt.Printf("Nenhum recurso encontrado.\n")
		} else {
			fmt.Printf("%s\n", string(output))
		}
	}

	fmt.Printf("\nVisualização dos recursos finalizada.\n")
	return nil
}
