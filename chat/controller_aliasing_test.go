package chat

import (
	"context"
	"sync"
	"testing"

	openai "github.com/jeanhaley32/go-openai-client"
)

// aliasBackend records the MaxTokens/Temperature it is handed and can mutate the
// pointee, so a test can detect whether SendMessage handed out a pointer that
// aliases the controller's internal default state.
type aliasBackend struct {
	mu      sync.Mutex
	sawMax  []int
	sawTemp []float64
	mutate  bool // when true, write through the received pointers
	mutMax  int
	mutTemp float64
}

func (b *aliasBackend) Name() string { return "alias-mock" }

func (b *aliasBackend) ChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	b.mu.Lock()
	if req.MaxTokens != nil {
		b.sawMax = append(b.sawMax, *req.MaxTokens)
	}
	if req.Temperature != nil {
		b.sawTemp = append(b.sawTemp, *req.Temperature)
	}
	mutate := b.mutate
	b.mu.Unlock()

	// A misbehaving/adversarial backend that writes to the optional-parameter
	// pointers must not be able to corrupt the controller's shared defaults.
	if mutate {
		if req.MaxTokens != nil {
			*req.MaxTokens = b.mutMax
		}
		if req.Temperature != nil {
			*req.Temperature = b.mutTemp
		}
	}

	return &openai.ChatCompletionResponse{
		Choices: []openai.Choice{{
			Index:   0,
			Message: openai.Message{Role: "assistant", Content: "ok"},
		}},
	}, nil
}

func (b *aliasBackend) SendMessage(ctx context.Context, req openai.Request) (*openai.Response, error) {
	return &openai.Response{}, nil
}
func (b *aliasBackend) IsAvailable(ctx context.Context) bool   { return true }
func (b *aliasBackend) Configure(map[string]interface{}) error { return nil }

// TestSendMessageDefaultsDoNotAliasControllerState verifies that when a request
// omits MaxTokens/Temperature, the controller falls back to its defaults by
// value — not by handing out &c.maxTokens / &c.temperature. If a backend writes
// through the pointer it receives, the controller's own defaults must be
// unaffected for subsequent requests.
func TestSendMessageDefaultsDoNotAliasControllerState(t *testing.T) {
	backend := &aliasBackend{mutate: true, mutMax: 999999, mutTemp: -3.14}
	c := NewController(backend, &ControllerConfig{
		DefaultModel: "mock-model",
		MaxTokens:    100,
		Temperature:  0.7,
	})

	// First default request: the backend corrupts the pointer it is handed.
	if _, err := c.SendMessage(context.Background(), ChatRequest{Message: "one"}); err != nil {
		t.Fatalf("first SendMessage failed: %v", err)
	}
	// Second default request: it must still see the original defaults, proving
	// the first request did not alias (and thus corrupt) controller state.
	if _, err := c.SendMessage(context.Background(), ChatRequest{Message: "two"}); err != nil {
		t.Fatalf("second SendMessage failed: %v", err)
	}

	backend.mu.Lock()
	defer backend.mu.Unlock()
	for i, v := range backend.sawMax {
		if v != 100 {
			t.Errorf("request %d saw MaxTokens=%d, want 100 (controller default was corrupted via aliased pointer)", i, v)
		}
	}
	for i, v := range backend.sawTemp {
		if v != 0.7 {
			t.Errorf("request %d saw Temperature=%v, want 0.7 (controller default was corrupted via aliased pointer)", i, v)
		}
	}
}

// TestSendMessageConcurrentDefaultsRaceFree runs many concurrent default
// requests; under `go test -race` a shared pointer into controller state that
// any goroutine writes would be reported as a data race.
func TestSendMessageConcurrentDefaultsRaceFree(t *testing.T) {
	backend := &aliasBackend{mutate: true, mutMax: 1, mutTemp: 0.1}
	c := NewController(backend, &ControllerConfig{
		DefaultModel: "mock-model",
		MaxTokens:    256,
		Temperature:  0.5,
	})

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = c.SendMessage(context.Background(), ChatRequest{Message: "concurrent"})
		}()
	}
	wg.Wait()
}
