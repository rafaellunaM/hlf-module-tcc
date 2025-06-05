package administration

import (
	"fmt"
//	"os"
	"os/exec"
)

func ShowAllResources() error {
	fmt.Printf("Mostrando recursos\n")

	output, err = exec.Command("kubectl", "get", "pods")
	
	if err != nil {
			return fmt.Errorf("Erro ao Mostrar recurso %s: %v\nSaída: %s", resource, err, string(output))
	}

	fmt.Printf("Recursos: \n %s", output)
	return nil
}
