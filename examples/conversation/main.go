package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jeanhaley/go-openai-client"
	"github.com/jeanhaley/go-openai-client/chat"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not set, using mock backend for demonstration")
		runWithMockBackend()
		return
	}

	// Create OpenAI client
	client := openai.NewClient(openai.Config{
		APIKey: apiKey,
		Model:  "gpt-3.5-turbo",
	})

	runConversation(client)
}

func runWithMockBackend() {
	// Create mock backend for demonstration
	mockBackend := openai.NewMockBackend()
	runConversation(mockBackend)
}

func runConversation(backend openai.Backend) {
	// Create chat controller
	controller := chat.NewController(backend, &chat.ControllerConfig{
		DefaultModel: "gpt-3.5-turbo",
		MaxTokens:    150,
		Temperature:  0.7,
	})

	ctx := context.Background()

	fmt.Printf("Using backend: %s\n", backend.Name())
	fmt.Println("Starting conversation...")

	// Create a conversation with system prompt
	conversation := controller.CreateConversation("You are a helpful assistant that loves to tell jokes.")

	// Send first message
	resp1, err := controller.SendMessage(ctx, chat.ChatRequest{
		ConversationID: conversation.ID,
		Message:        "Hello! Can you tell me a joke?",
	})
	if err != nil {
		log.Fatalf("Error sending first message: %v", err)
	}

	fmt.Printf("\nUser: Hello! Can you tell me a joke?\n")
	fmt.Printf("Assistant: %s\n", resp1.Message.Content)

	// Send follow-up message
	resp2, err := controller.SendMessage(ctx, chat.ChatRequest{
		ConversationID: conversation.ID,
		Message:        "That was funny! Can you tell me another one?",
	})
	if err != nil {
		log.Fatalf("Error sending second message: %v", err)
	}

	fmt.Printf("\nUser: That was funny! Can you tell me another one?\n")
	fmt.Printf("Assistant: %s\n", resp2.Message.Content)

	// Get conversation summary
	summary, err := controller.GetConversationSummary(conversation.ID)
	if err != nil {
		log.Fatalf("Error getting summary: %v", err)
	}

	fmt.Printf("\n--- Conversation Summary ---\n")
	fmt.Printf("Messages: %d\n", summary.MessageCount)
	fmt.Printf("User messages: %d\n", summary.UserMessages)
	fmt.Printf("Assistant messages: %d\n", summary.AssistantMessages)
	fmt.Printf("Estimated tokens: %d\n", summary.EstimatedTokens)
}