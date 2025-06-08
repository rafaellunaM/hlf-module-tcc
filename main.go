package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	Function    func(string) error
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

func selectConfigFile() string {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("SELEÇÃO DE ARQUIVO DE CONFIGURAÇÃO")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("1. hlf-config.json (configuração padrão)")
	fmt.Println("2. templates/4-orderers.json (4 orderers)")
	fmt.Println("3. templates/4-peers.json (4 peers)")
	fmt.Println(strings.Repeat("=", 50))
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Escolha o arquivo de configuração (1-3): ")
		
		if !scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			fmt.Println("Usando: hlf-config.json")
			return "hlf-config.json"
		case "2":
			fmt.Println("Usando: templates/4-orderers.json")
			return "templates/4-orderers.json"
		case "3":
			fmt.Println("Usando: templates/4-peers.json")
			return "templates/4-peers.json"
		default:
			fmt.Println("Opção inválida. Escolha 1, 2 ou 3.")
		}
	}
	return "hlf-config.json"
}

func main() {
	displayWelcomeBanner()
	
	configFile := selectConfigFile()
	
	steps := []Step{
		{1,  "Create   CAs", "Criar Certificate Authorities", func(config string) error { return ca.CreateCAs(config) }},
		{2,  "Register Orderers", "Registrar orderers", func(config string) error { return ca.RegisterOrderers(config) }},
		{3,  "Register Peers", "Registrar peers", func(config string) error { return ca.RegisterPeers(config) }},
		{4,  "Deploy   Peers", "Fazer deploy dos peers", func(config string) error { return node.DeployPeers(config) }},
		{5,  "Deploy   Orderers", "Fazer deploy dos orderers", func(config string) error { return node.DeployOrderers(config) }},
		{6,  "Register Channels", "Registrar channels", func(config string) error { return ca.RegisterChannels(config) }},
		{7,  "Enroll   Channels", "Fazer enroll dos channels", func(config string) error { return ca.EnrollChannels(config) }},
		{8,  "Create   Wallet", "Criar wallet", func(config string) error { return ca.CreateWallet(config) }},
		{9,  "Execute  PEM Script", "Extrair certificado PEM", func(config string) error { return scripts.ExecutePemScript(config) }},
		{10, "Create   Main Channel", "Criar canal principal", func(config string) error { return channels.CreateMainChannel(config) }},
		{11, "Create   Follower Channel", "Criar canal para os peers", func(config string) error { return channels.JoinChannel(config) }},
		{12, "Delete   Components", "Delete todos os Components HLF e secret", func(config string) error { return administration.DeleteAllResources() }},
		{13, "Mostrar  Components", "Mostra todos Components para o HLF", func(config string) error { return administration.ShowResources() }},
		{14, "Change   Config", "Alterar arquivo de configuração", func(config string) error { return changeConfig() }},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		displayMenu(steps, configFile)
		
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
			executeAllSteps(steps, configFile)
			continue
		}

		if strings.Contains(input, "-") {
			if err := executeStepRange(input, steps, configFile); err != nil {
				fmt.Printf("Erro ao executar range: %v\n", err)
			}
			continue
		}

		if strings.Contains(input, ",") {
			if err := executeMultipleSteps(input, steps, configFile); err != nil {
				fmt.Printf("Erro ao executar múltiplos passos: %v\n", err)
			}
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Entrada inválida. Digite um número, 'all', ou 'q' para sair.")
			continue
		}
		
		if choice == 14 {
			configFile = selectConfigFile()
			continue
		}
		
		if err := executeSingleStep(choice, steps, configFile); err != nil {
			fmt.Printf("Erro: %v\n", err)
		}
	}
}

func changeConfig() error {
	fmt.Println("Para alterar a configuração, escolha a opção 14 no menu.")
	return nil
}

func displayMenu(steps []Step, currentConfig string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("HYPERLEDGER FABRIC DEPLOYMENT CLI")
	fmt.Printf("Configuração atual: %s\n", currentConfig)
	fmt.Println(strings.Repeat("=", 60))
	
	for _, step := range steps {
		fmt.Printf("%2d. %-25s - %s\n", step.ID, step.Name, step.Description)
	}
	
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("Opções especiais:")
	fmt.Println("   all      - Executar todos os passos em sequência")
	fmt.Println("   1-5      - Executar passos de 1 a 5")
	fmt.Println("   1,3,5    - Executar passos 1, 3 e 5")
	fmt.Println("   q        - Sair")
	fmt.Println(strings.Repeat("=", 60))
}

func executeSingleStep(choice int, steps []Step, configFile string) error {
	if choice < 1 || choice > len(steps) {
		return fmt.Errorf("opção inválida. Escolha entre 1 e %d", len(steps))
	}
	
	step := steps[choice-1]
	fmt.Printf("\nExecutando: %s...\n", step.Name)
	fmt.Printf("Usando configuração: %s\n", configFile)
	
	if err := step.Function(configFile); err != nil {
		return fmt.Errorf("erro ao executar '%s': %v", step.Name, err)
	}
	
	fmt.Printf("%s executado com sucesso!\n", step.Name)
	return nil
}

func executeAllSteps(steps []Step, configFile string) {
	fmt.Println("\nExecutando todos os passos em sequência...")
	fmt.Printf("Usando configuração: %s\n", configFile)
	fmt.Println("Nota: Os passos de delete (12) e change config (14) serão pulados na execução completa.")
	
	for _, step := range steps {
		if step.ID == 12 || step.ID == 14 {
			fmt.Printf("\nPulando passo %d: %s (não executado no modo 'all')\n", step.ID, step.Name)
			continue
		}
		
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(configFile); err != nil {
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

		if step.ID != 11 {
			fmt.Println("Aguardando... ")
			time.Sleep(20 * time.Second)
		}
	}
	
	fmt.Println("\nTodos os passos de deployment foram executados!")
	fmt.Println("Para deletar recursos, execute o passo 12 individualmente.")
}

func executeStepRange(input string, steps []Step, configFile string) error {
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
	fmt.Printf("Usando configuração: %s\n", configFile)
	
	for i := start - 1; i < end; i++ {
		step := steps[i]
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(configFile); err != nil {
			return fmt.Errorf("erro no passo %d (%s): %v", step.ID, step.Name, err)
		}
		
		fmt.Printf("Passo %d concluído: %s\n", step.ID, step.Name)

		if i < end-1 {
			fmt.Println("Aguardando... ")
			time.Sleep(20 * time.Second)
		}
	}
	
	return nil
}

func executeMultipleSteps(input string, steps []Step, configFile string) error {
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
	fmt.Printf("Usando configuração: %s\n", configFile)
	
	for i, stepNum := range selectedSteps {
		step := steps[stepNum-1]
		fmt.Printf("\nExecutando passo %d: %s...\n", step.ID, step.Name)
		
		if err := step.Function(configFile); err != nil {
			return fmt.Errorf("erro no passo %d (%s): %v", step.ID, step.Name, err)
		}
		
		fmt.Printf("Passo %d concluído: %s\n", step.ID, step.Name)

		if i < len(selectedSteps)-1 {
			fmt.Println("Aguardando... ")
			time.Sleep(20 * time.Second)
		}
	}
	
	return nil
}
