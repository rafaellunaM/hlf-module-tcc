package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/exec"
)

type OrdererConfig struct {
    EnrollIDorderer    string `json:"enrollIDorderer"`
    EnrollPWorderer    string `json:"enrollPWorderer"`
    MSPID       string `json:"mspid"`
    Capacity    string `json:"capacity"`
    Name        string `json:"name"`
    Hosts       string `json:"hosts"`
    IstioPort   string `json:"istioPort"`
    AdminHosts  string `json:"admin-hosts"`
    CAName      string `json:"CAName"`
}

type FullResources struct {
    Orderers []OrdererConfig `json:"Orderes"`
}

func main() {
    ordererImage := os.Getenv("ORDERER_IMAGE")
    if ordererImage == "" {
        log.Fatal("❌ variável de ambiente ORDERER_IMAGE não definida")
    }

    ordererVersion := os.Getenv("ORDERER_VERSION")
    if ordererVersion == "" {
        log.Fatal("❌ variável de ambiente ORDERER_VERSION não definida")
    }

    storageClass := os.Getenv("SC_NAME")
    if storageClass == "" {
        log.Fatal("❌ variável de ambiente SC_NAME não definida")
    }

    data, err := os.ReadFile("output.json")
    if err != nil {
        log.Fatalf("❌ não consegui ler output.json: %v", err)
    }

    var cfg FullResources
    if err := json.Unmarshal(data, &cfg); err != nil {
        log.Fatalf("❌ erro ao parsear JSON: %v", err)
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

    for _, o := range cfg.Orderers {
        fmt.Printf("🚀 Criando o orderer `%s`…\n", o.Name)
        args := []string{
            "hlf", "ordnode", "create",
            "--image=" + ordererImage,
            "--version=" + ordererVersion,
            "--storage-class=" + storageClass,
            "--enroll-id=" + o.EnrollIDorderer,
            "--enroll-pw=" + o.EnrollPWorderer,
            "--mspid=" + o.MSPID,
            "--capacity=" + o.Capacity,
            "--name=" + o.Name,
            "--ca-name=" + o.CAName + ".default",
            "--hosts=" + o.Hosts,
            "--admin-hosts=" + o.AdminHosts,
            "--istio-port=" + o.IstioPort,
        }
        run(args...)
        fmt.Printf("✅ Orderer `%s` criado com sucesso.\n\n", o.Name)
    }

    fmt.Println("⏳ Aguardando todos os orderer nodes ficarem em estado Running...")
    waitCmd := exec.Command("kubectl", "wait",
        "--timeout=180s",
        "--for=condition=Running",
        "fabricorderernodes.hlf.kungfusoftware.es",
        "--all",
    )
    waitCmd.Stdout = os.Stdout
    waitCmd.Stderr = os.Stderr
    if err := waitCmd.Run(); err != nil {
        log.Fatalf("❌ Erro ao aguardar orderers: %v", err)
    }
    fmt.Println("✅ Todos os orderers estão em execução.")
}
