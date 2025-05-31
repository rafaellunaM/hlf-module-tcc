package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
    "hlf/internal/fabric"
)

func main() {
    peerImage := os.Getenv("PEER_IMAGE")
    if peerImage == "" {
        log.Fatal("❌ variável de ambiente PEER_IMAGE não definida")
    }
    peerVersion := os.Getenv("PEER_VERSION")
    if peerVersion == "" {

        log.Fatal("❌ variável de ambiente PEER_VERSION não definida")
    }

    sc_name := os.Getenv("SC_NAME")
    if peerImage == "" {
        log.Fatal("❌ variável de ambiente SC_NAME não definida")
    }

    data, err := os.ReadFile("hlf-config.json")
    if err != nil {
        log.Fatalf("❌ não consegui ler hlf-config.json: %v", err)
    }

    var partialConfig struct {
		Peers []fabric.Peer `json:"Peers"`
	}
	
    if err := json.Unmarshal(data, &partialConfig); err != nil {
        log.Fatalf("❌ não consegui parsear JSON: %v", err)
    }

    run := func(args ...string) {
        fmt.Printf("🔧 Executando: kubectl %v\n", args)
        cmd := exec.Command("kubectl", args...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            log.Fatalf("❌ erro ao executar %v: %v", args, err)
        }
    }

    for _, p := range partialConfig.Peers {
        fmt.Printf("🚀 Fazendo deploy do peer `%s`…\n", p.Name)
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
            "--storage-class=" + sc_name,
            "--capacity=" + p.Capacity,
            "--image=" + peerImage,
            "--version=" + peerVersion,
        }
        run(args...)
        fmt.Printf("✅ Peer `%s` deployado.\n\n", p.Name)
    }
}
