package node

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "hlf/internal/fabric"
)

func DeployOrderers(configFile string) error {
    ordererImage := os.Getenv("ORDERER_IMAGE")
    if ordererImage == "" {
        return fmt.Errorf("variável de ambiente ORDERER_IMAGE não definida")
    }

    ordererVersion := os.Getenv("ORDERER_VERSION")
    if ordererVersion == "" {
        return fmt.Errorf("variável de ambiente ORDERER_VERSION não definida")
    }

    storageClass := os.Getenv("SC_NAME")
    if storageClass == "" {
        return fmt.Errorf("variável de ambiente SC_NAME não definida")
    }

    data, err := os.ReadFile(configFile)
    if err != nil {
        return fmt.Errorf("não consegui ler %s: %v", configFile, err)
    }

    var partialConfig struct {
        Orderer []fabric.Orderer `json:"Orderer"`
    }
    
    if err := json.Unmarshal(data, &partialConfig); err != nil {
        return fmt.Errorf("erro ao parsear JSON: %v", err)
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

    for _, o := range partialConfig.Orderer {
        fmt.Printf(" Criando o orderer `%s`…\n", o.Name)
        args := []string{
            "hlf", "ordnode", "create",
            "--image=" + ordererImage,
            "--version=" + ordererVersion,
            "--storage-class=" + storageClass,
            "--enroll-id=" + o.EnrollIDorderer,
            "--enroll-pw=" + o.EnrollPWorderer,
            "--mspid=" + o.Mspid,
            "--capacity=" + o.Capacity,
            "--name=" + o.Name,
            "--ca-name=" + o.CAName + ".default",
            "--hosts=" + o.Hosts,
            "--admin-hosts=" + o.AdminHosts,
            "--istio-port=" + o.IstioPort,
        }
        
        if err := run(args...); err != nil {
            return err
        }
        fmt.Printf(" Orderer `%s` criado com sucesso.\n\n", o.Name)
    }

    fmt.Println("Aguardando todos os orderer nodes ficarem em estado Running...")
    waitCmd := exec.Command("kubectl", "wait",
        "--timeout=180s",
        "--for=condition=Running",
        "fabricorderernodes.hlf.kungfusoftware.es",
        "--all",
    )
    waitCmd.Stdout = os.Stdout
    waitCmd.Stderr = os.Stderr
    if err := waitCmd.Run(); err != nil {
        return fmt.Errorf("erro ao aguardar orderers: %v", err)
    }
    fmt.Println(" Todos os orderers estão em execução.")
    
    return nil
}
