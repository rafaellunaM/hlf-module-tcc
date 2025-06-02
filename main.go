package main

import (
	"log"
	"hlf/cmd/ca"  
	"hlf/cmd/node"
	"hlf/cmd/channels"
	"hlf/cmd/scripts"
)

func main() {

	if err := ca.CreateCAs("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao criar CAs: %v", err)
	}
	
	if err := ca.RegisterOrderers("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao registrar orderers: %v", err)
	}

	if err := ca.RegisterPeers("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao registrar peers: %v", err)
	}

	if err := node.DeployPeers("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao fazer deploy dos peers: %v", err)
	}
	
	if err := node.DeployOrderers("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao fazer deploy dos orderers: %v", err)
	}

	if err := ca.RegisterChannels("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao registrar channels: %v", err)
	}

	if err := ca.EnrollChannels("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao fazer enroll dos channels: %v", err)
	}

	if err := ca.CreateWallet("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao criar wallet: %v", err)
	}

	if err := scripts.ExecutePemScript(); err != nil {
			log.Fatalf("Erro ao extrair certificado PEM: %v", err)
	}

	if err := channels.CreateMainChannel("hlf-config.json"); err != nil {
		log.Fatalf("Erro ao criar wallet: %v", err)
	}

	log.Println("Processo completo finalizado!")
}
