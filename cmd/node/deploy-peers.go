package node

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "hlf/internal/fabric"
)

func DeployPeers(configFile string) error {
	peerImage := os.Getenv("PEER_IMAGE")
	if peerImage == "" {
			return fmt.Errorf("variável de ambiente PEER_IMAGE não definida")
	}
	
	peerVersion := os.Getenv("PEER_VERSION")
	if peerVersion == "" {
			return fmt.Errorf("variável de ambiente PEER_VERSION não definida")
	}

	scName := os.Getenv("SC_NAME")
	if scName == "" {
			return fmt.Errorf("variável de ambiente SC_NAME não definida")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
			return fmt.Errorf("não consegui ler %s: %v", configFile, err)
	}

	var partialConfig struct {
			Peers []fabric.Peer `json:"Peers"`
	}
	
	if err := json.Unmarshal(data, &partialConfig); err != nil {
			return fmt.Errorf("não consegui parsear JSON: %v", err)
	}

	run := func(args ...string) error {
			fmt.Printf("Executando: kubectl %v\n", args)
			cmd := exec.Command("kubectl", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
					return fmt.Errorf("erro ao executar %v: %v", args, err)
			}
			return nil
	}

	for _, p := range partialConfig.Peers {
			fmt.Printf(" Fazendo deploy do peer `%s`…\n", p.Name)
			args := []string{
					"hlf", "peer", "create",
					"--statedb=" + p.StateDB, 
					"--enroll-id=" + p.EnrollIDpeer,
					"--enroll-pw=" + p.EnrollIPWpeer,
					"--mspid=" + p.Mspid,
					"--name=" + p.Name,
					"--ca-name=" + p.CAName + ".default",
					"--hosts=" + p.Hosts,
					"--istio-port=" + p.IstioPort,
					"--storage-class=" + scName,
					"--capacity=" + p.Capacity,
					"--image=" + peerImage,
					"--version=" + peerVersion,
			}
			
			if err := run(args...); err != nil {
					return err
			}
			fmt.Printf(" Peer `%s` deployado.\n\n", p.Name)
	}
	
	return nil
}
