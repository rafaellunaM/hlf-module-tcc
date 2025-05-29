package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
)

type PeerConfig struct {
    EnrollIDpeer    string `json:"enrollIDpeer"`
    EnrollIPWpeer	string `json:"enrollIPWpeer"`
    StateDB			string `json:"stateDB"`
    Capacity    string `json:"capacity"`
    MSPID		    string `json:"mspid"`
    Name        string `json:"name"`
    CAName				string `json:"CAName"`
    Hosts       string `json:"hosts"`
    IstioPort   string `json:"istioPort"`    
}

type FullResources struct {
    Peers []PeerConfig `json:"Peers"`
}

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

    data, err := os.ReadFile("output.json")
    if err != nil {
        log.Fatalf("❌ não consegui ler output.json: %v", err)
    }

    var cfg FullResources
    if err := json.Unmarshal(data, &cfg); err != nil {
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

    for _, p := range cfg.Peers {
        fmt.Printf("🚀 Fazendo deploy do peer `%s`…\n", p.Name)
        args := []string{
            "hlf", "peer", "create",
            "--statedb=" + p.StateDB, 
            "--enroll-id=" + p.EnrollIDpeer,
            "--enroll-pw=" + p.EnrollIPWpeer,
            "--mspid=" + p.MSPID,
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
