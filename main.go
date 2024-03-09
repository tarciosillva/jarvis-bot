package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// State representa os estados da máquina de estado
type State struct {
	SessionId string
}

type Node struct {
	BotMessage       string            `json:"botMessage"`
	NextNode         string            `json:"nextNode"`
	VariableResponse map[string]string `json:"variableResponse"`
}

type ConversationContext struct {
	PreviousUserInput string
	CurrentState      string
	Variables         map[string]string
}

func main() {
	// Carrega o arquivo que contém o fluxo do bot
	flowJSON, err := os.ReadFile("conversation_flow.json")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo de fluxo da conversa:", err)
		return
	}

	startConversation(flowJSON)
}

// Função principal para iniciar a conversa
func startConversation(flowJSON []byte) {
	context := &ConversationContext{
		CurrentState: "Inicio",
	}

	// Carrega o fluxo da conversa a partir do JSON
	flow, err := loadConversationFlow(flowJSON)
	if err != nil {
		fmt.Println("Erro ao carregar o fluxo da conversa:", err)
		return
	}

	processNode(flow[context.CurrentState], context, flow)
}

// Função para carregar o fluxo da conversa a partir do JSON
func loadConversationFlow(flowJSON []byte) (map[string]Node, error) {
	var flow map[string]Node
	if err := json.Unmarshal(flowJSON, &flow); err != nil {
		return nil, fmt.Errorf("Erro ao fazer unmarshal do JSON: %v", err)
	}
	return flow, nil
}

// Função para processar um nó da árvore de decisão
func processNode(node Node, context *ConversationContext, flow map[string]Node) {
	// Exibe a mensagem do bot
	fmt.Println(node.BotMessage)

	userInput := getUserInput()

	context.PreviousUserInput = userInput

	// Verifica se o próximo nó existe no fluxo
	if nextNode, ok := flow[node.NextNode]; ok {
		if len(node.VariableResponse) > 0 {
			// Se houver variáveis de resposta, atualize a mensagem do próximo nó
			newBotMessage := replaceVariables(nextNode.BotMessage, node.VariableResponse, userInput)
			nextNode.BotMessage = newBotMessage
		}
		processNode(nextNode, context, flow)
	} else {
		fmt.Println("Interação encerrada")
	}
}

// Função para substituir variáveis dinâmicas em uma string
func replaceVariables(input string, variables map[string]string, userInput string) string {
	for _, value := range variables {
		// Substitui todas as ocorrências de {key} pelo valor correspondente
		input = strings.ReplaceAll(input, "{"+value+"}", userInput)
	}
	return input
}

// Função para obter a entrada do usuário
func getUserInput() string {
	fmt.Print("Usuário: ")
	var userInput string
	fmt.Scanln(&userInput)
	return strings.ToLower(userInput)
}
