package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jeanhaley/go-openai-client"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI client
	client := openai.NewClient(openai.Config{
		APIKey: apiKey,
		Model:  "gpt-3.5-turbo",
	})

	// Create a simple chat completion request
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "Hello! Can you tell me a brief joke?",
			},
		},
		MaxTokens:   func(i int) *int { return &i }(100),
		Temperature: func(f float64) *float64 { return &f }(0.7),
	}

	// Send the request
	ctx := context.Background()
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	// Print the response
	if len(resp.Choices) > 0 {
		fmt.Printf("AI Response: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
	} else {
		fmt.Println("No response received")
	}
}