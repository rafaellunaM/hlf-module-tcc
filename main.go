package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"hlf/cmd/ca"  
	"hlf/cmd/node"
	"hlf/cmd/channels"
	"hlf/cmd/scripts"
	"hlf/cmd/administration"
)

type Step struct {
	ID          int
	Name        string
	Description string
	Function    func() error
}

func displayWelcomeBanner() {
	banner := `
 ██   ██ ██      ███████      █████  ██    ██ ████████  ██████  
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ███████ ██      █████       ███████ ██    ██    ██    ██    ██ 
 ██   ██ ██      ██          ██   ██ ██    ██    ██    ██    ██ 
 ██   ██ ███████ ██          ██   ██  ██████     ██     ██████  

           Hyperledger Fabric Automated Deployment Tool
               Version 1.0 - Created by Rafael Luna
`
	fmt.Println(banner)
}

func main() {
	displayWelcomeBanner()
	
	configFile := "hlf-config.json"
	
	steps := []Step{
		{1, "Create CAs", "Criar Certificate Authorities", func() error { return ca.CreateCAs(configFile) }},
		{2, "Register Orderers", "Registrar orderers", func() error { return ca.RegisterOrderers(configFile) }},
		{3, "Register Peers", "Registrar peers", func() error { return ca.RegisterPeers(configFile) }},
		{4, "Deploy Peers", "Fazer deploy dos peers", func() error { return node.DeployPeers(configFile) }},
		{5, "Deploy Orderers", "Fazer deploy dos orderers", func() error { return node.DeployOrderers(configFile) }},
		{6, "Register Channels", "Registrar channels", func() error { return ca.RegisterChannels(configFile) }},
		{7, "Enroll Channels", "Fazer enroll dos channels", func() error { return ca.EnrollChannels(configFile) }},
		{8, "Create Wallet", "Criar wallet", func() error { return ca.CreateWallet(configFile) }},
		{9, "Execute PEM Script", "Extrair certificado PEM", func() error { return scripts.ExecutePemScript(configFile) }},
		{10, "Create Main Channel", "Criar canal principal", func() error { return channels.CreateMainChannel(configFile) }},
		{11, "Delete todos os recursos", "Delete todos os recursos HLF e secret", func() error { return administration.DeleteAllResources() }},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		displayMenu(steps)
		
		fmt.Print("\nEscolha uma opção (ou 'q' para sair): ")
		
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		
		if input == "q" || input == "Q" || input == "quit" || input == "exit" {
			fmt.Println("Saindo...")
			break
		}
		
		if input == "all" || input == "ALL" {
			executeAllSteps(steps)
			continue
		}

		if strings.Contains(input, "-") {
			if err := executeStepRange(input, steps); err != nil {
				fmt.Printf("Erro ao executar range: %v\n", err)
			}
			continue
		}

		if strings.Contains(input, ",") {
			if err := executeMultipleSteps(input, steps); err != nil {
				fmt.Printf("Erro ao executar múltiplos passos: %v\n", err)
			}
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Entrada inválida. Digite um número, 'all', ou 'q' para sair.")
			continue
		}
		
		if err := executeSingleStep(choice, steps); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	}
}

func displayMenu(steps []Step) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("HYPERLEDGER FABRIC DEPLOYMENT CLI")
	fmt.Println(strings.Repeat("=", 60))
	
	for _, step := range steps {
		fmt.Printf("%2d. %-20s - %s\n", step.ID, step.Name, step.Description)
	}
	
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Opções especiais:")
	fmt.Println("   all      - Executar todos os passos em sequência")
	fmt.Println("   1-5      - Executar passos de 1 a 5")
	fmt.Println("   1,3,5    - Executar passos 1, 3 e 5")
	fmt.Println("   q        - Sair")
	fmt.Println(strings.Repeat("=", 60))
}

func executeSingleStep(choice int, steps []Step) error {
	if choice < 1 || choice > len(steps) {
		return fmt.Errorf("opção inválida. Escolha entre 1 e %d", len(steps))
	}
	
	step := steps[choice-1]
	fmt.Printf("\nExecutando: %s...\n", step.Name)
	
	if err := step.Function(); err != nil {
		return fmt.Errorf("erro ao executar '%s': %v", step.Name, err)
	}
	
	fmt.Printf("%s executado com sucesso!\n", step.Name)
	return nil
}

func executeAllSteps(steps []Step) {
	fmt.Println("\nExecutando todos os passos em sequência...")
	
	for _, step := range steps {
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(); err != nil {
			fmt.Printf("Erro no passo %d (%s): %v\n", step.ID, step.Name, err)
			fmt.Print("\nDeseja continuar com os próximos passos? (y/N): ")
			
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				response := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if response != "y" && response != "yes" {
					fmt.Println("Execução interrompida pelo usuário.")
					return
				}
			}
			continue
		}
		
		fmt.Printf("Passo %d concluído: %s\n", step.ID, step.Name)
	}
	
	fmt.Println("\nTodos os passos foram executados!")
}

func executeStepRange(input string, steps []Step) error {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return fmt.Errorf("formato de range inválido. Use: 1-5")
	}
	
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("número inicial inválido: %s", parts[0])
	}
	
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("número final inválido: %s", parts[1])
	}
	
	if start < 1 || end > len(steps) || start > end {
		return fmt.Errorf("range inválido. Use valores entre 1 e %d", len(steps))
	}
	
	fmt.Printf("\nExecutando passos %d a %d...\n", start, end)
	
	for i := start - 1; i < end; i++ {
		step := steps[i]
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(); err != nil {
			return fmt.Errorf("erro no passo %d (%s): %v", step.ID, step.Name, err)
		}
		
		fmt.Printf("Passo %d concluído: %s\n", step.ID, step.Name)
	}
	
	return nil
}

func executeMultipleSteps(input string, steps []Step) error {
	stepNumbers := strings.Split(input, ",")
	var selectedSteps []int
	
	for _, numStr := range stepNumbers {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return fmt.Errorf("número inválido: %s", numStr)
		}
		
		if num < 1 || num > len(steps) {
			return fmt.Errorf("passo %d inválido. Use valores entre 1 e %d", num, len(steps))
		}
		
		selectedSteps = append(selectedSteps, num)
	}
	
	fmt.Printf("\nExecutando passos selecionados: %v...\n", selectedSteps)
	
	for _, stepNum := range selectedSteps {
		step := steps[stepNum-1]
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(); err != nil {
			return fmt.Errorf("erro no passo %d (%s): %v", step.ID, step.Name, err)
		}
		
		fmt.Printf("Passo %d concluído: %s\n", step.ID, step.Name)
	}
	
	return nil
}
