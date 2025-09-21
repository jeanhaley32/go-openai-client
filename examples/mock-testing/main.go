package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jeanhaley32/go-openai-client"
	"github.com/jeanhaley32/go-openai-client/chat"
)

func main() {
	fmt.Println("=== Mock Backend Testing Example ===")

	// Create mock backend
	mockBackend := openai.NewMockBackend()

	// Test basic functionality
	testBasicChat(mockBackend)

	// Test chat controller
	testChatController(mockBackend)

	// Test configuration
	testConfiguration(mockBackend)
}

func testBasicChat(backend openai.Backend) {
	fmt.Printf("\n--- Testing Basic Chat with %s ---\n", backend.Name())

	ctx := context.Background()

	// Test ChatCompletion
	req := openai.ChatCompletionRequest{
		Model: "mock-model-v1",
		Messages: []openai.Message{
			{Role: "user", Content: "Hello, testing mock backend!"},
		},
	}

	resp, err := backend.ChatCompletion(ctx, req)
	if err != nil {
		log.Fatalf("ChatCompletion failed: %v", err)
	}

	fmt.Printf("Request: %s\n", req.Messages[0].Content)
	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)

	// Test legacy SendMessage
	legacyReq := openai.Request{
		Model: "mock-model-v1",
		Messages: []openai.Message{
			{Role: "user", Content: "Testing legacy format"},
		},
	}

	legacyResp, err := backend.SendMessage(ctx, legacyReq)
	if err != nil {
		log.Fatalf("SendMessage failed: %v", err)
	}

	fmt.Printf("\nLegacy Request: %s\n", legacyReq.Messages[0].Content)
	fmt.Printf("Legacy Response: %s\n", legacyResp.Content)
	fmt.Printf("Tokens used: %d\n", legacyResp.TokensUsed)
}

func testChatController(backend openai.Backend) {
	fmt.Printf("\n--- Testing Chat Controller with %s ---\n", backend.Name())

	// Create controller
	controller := chat.NewController(backend, &chat.ControllerConfig{
		DefaultModel: "mock-model-v1",
		MaxTokens:    100,
		Temperature:  0.5,
	})

	ctx := context.Background()

	// Test conversation creation and messaging
	conv := controller.CreateConversation("You are a testing assistant.")

	resp, err := controller.SendMessage(ctx, chat.ChatRequest{
		ConversationID: conv.ID,
		Message:        "Test message for conversation",
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Conversation ID: %s\n", conv.ID)
	fmt.Printf("User: Test message for conversation\n")
	fmt.Printf("Assistant: %s\n", resp.Message.Content)

	// Test conversation management
	conversations := controller.ListConversations()
	fmt.Printf("\nTotal conversations: %d\n", len(conversations))

	// Test stats
	stats := controller.GetStats()
	fmt.Printf("Controller stats - Messages: %d, Conversations: %d, Backend: %s\n",
		stats.TotalMessages, stats.TotalConversations, stats.BackendName)
}

func testConfiguration(backend openai.Backend) {
	fmt.Printf("\n--- Testing Configuration with %s ---\n", backend.Name())

	// Test availability
	ctx := context.Background()
	available := backend.IsAvailable(ctx)
	fmt.Printf("Backend available: %t\n", available)

	// Test configuration
	config := map[string]interface{}{
		"name":    "CustomMock",
		"setting": "test_value",
	}

	err := backend.Configure(config)
	if err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	fmt.Printf("Backend name after config: %s\n", backend.Name())
	fmt.Println("Configuration test completed successfully!")
}