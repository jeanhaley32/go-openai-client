# go-openai-client

A Go library for OpenAI Chat Completions API with conversation management and testing utilities.

## Features

- ✅ **OpenAI Chat Completions API**: Complete implementation following OpenAI's standard
- ✅ **Conversation Management**: High-level chat controller with conversation history
- ✅ **Mock Backend**: Built-in testing backend for development and testing
- ✅ **Thread-safe**: All operations are safe for concurrent use
- ✅ **Context Support**: Proper context cancellation throughout
- ✅ **Flexible Configuration**: Easy backend switching and configuration

## Installation

```bash
go get github.com/jeanhaley/go-openai-client
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jeanhaley/go-openai-client"
)

func main() {
    // Create OpenAI client
    client := openai.NewClient(openai.Config{
        APIKey: "your-api-key",
        Model:  "gpt-3.5-turbo",
    })

    // Send a chat completion request
    req := openai.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []openai.Message{
            {Role: "user", Content: "Hello!"},
        },
    }

    resp, err := client.ChatCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

### Conversation Management

```go
package main

import (
    "context"
    "fmt"

    "github.com/jeanhaley/go-openai-client"
    "github.com/jeanhaley/go-openai-client/chat"
)

func main() {
    // Create backend and controller
    client := openai.NewClient(openai.Config{
        APIKey: "your-api-key",
        Model:  "gpt-3.5-turbo",
    })

    controller := chat.NewController(client, &chat.ControllerConfig{
        DefaultModel: "gpt-3.5-turbo",
        MaxTokens:    150,
        Temperature:  0.7,
    })

    // Create conversation with system prompt
    conv := controller.CreateConversation("You are a helpful assistant.")

    // Send messages
    resp, err := controller.SendMessage(context.Background(), chat.ChatRequest{
        ConversationID: conv.ID,
        Message:        "Hello! How can you help me?",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Assistant: %s\n", resp.Message.Content)
}
```

### Testing with Mock Backend

```go
package main

import (
    "context"
    "fmt"

    "github.com/jeanhaley/go-openai-client"
)

func main() {
    // Use mock backend for testing
    mockBackend := openai.NewMockBackend()

    req := openai.ChatCompletionRequest{
        Model: "mock-model",
        Messages: []openai.Message{
            {Role: "user", Content: "Test message"},
        },
    }

    resp, err := mockBackend.ChatCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Mock response: %s\n", resp.Choices[0].Message.Content)
}
```

## API Reference

### Core Types

#### `Backend` Interface
All backends implement this interface:
- `ChatCompletion(ctx, req) (*ChatCompletionResponse, error)` - Send chat completion request
- `SendMessage(ctx, req) (*Response, error)` - Legacy method for backward compatibility
- `IsAvailable(ctx) bool` - Health check
- `Configure(config) error` - Set configuration
- `Name() string` - Get backend name

#### `ChatCompletionRequest`
Standard OpenAI Chat Completions request format:
```go
type ChatCompletionRequest struct {
    Model       string     `json:"model"`        // Required
    Messages    []Message  `json:"messages"`     // Required
    MaxTokens   *int       `json:"max_tokens,omitempty"`
    Temperature *float64   `json:"temperature,omitempty"`
    TopP        *float64   `json:"top_p,omitempty"`
    Stream      bool       `json:"stream,omitempty"`
}
```

#### `Message`
Individual message in a conversation:
```go
type Message struct {
    Role    string `json:"role"`    // "system", "user", or "assistant"
    Content string `json:"content"` // Message text
}
```

### Backends

#### OpenAI Client
```go
client := openai.NewClient(openai.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://api.openai.com/v1", // Optional
    Model:   "gpt-4",                     // Optional
    Timeout: 30 * time.Second,           // Optional
})
```

#### Mock Backend
```go
mockBackend := openai.NewMockBackend()
```

### Chat Controller

#### Configuration
```go
config := &chat.ControllerConfig{
    DefaultModel: "gpt-3.5-turbo",
    MaxTokens:    500,
    Temperature:  0.7,
}
controller := chat.NewController(backend, config)
```

#### Conversation Management
- `CreateConversation(systemPrompt) *Conversation` - Create new conversation
- `GetConversation(id) (*Conversation, error)` - Get existing conversation
- `ListConversations() []*Conversation` - List all conversations
- `DeleteConversation(id) error` - Delete conversation
- `SendMessage(ctx, request) (*ChatResponse, error)` - Send message

## Examples

See the [examples](./examples/) directory for complete working examples:

- [**basic**](./examples/basic/) - Simple OpenAI API usage
- [**conversation**](./examples/conversation/) - Conversation management
- [**mock-testing**](./examples/mock-testing/) - Testing with mock backend

## Environment Variables

For examples using real OpenAI API:
```bash
export OPENAI_API_KEY="your-api-key-here"
```

## Testing

The library includes a mock backend for testing:

```go
// In your tests
mockBackend := openai.NewMockBackend()
// Use mockBackend in place of real OpenAI client
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.