package minimax

import (
	"bufio"
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// clampMiniMaxTemperature applies MiniMax temperature clamping logic
// same as in chat_handler.go
func clampMiniMaxTemperature(modelValue string, temp float32) float32 {
	if strings.HasPrefix(strings.ToLower(modelValue), "minimax") {
		if temp > 1.0 {
			temp = 1.0
		}
		if temp < 0 {
			temp = 0
		}
	}
	return temp
}

// TestMiniMaxTemperatureClamping tests that MiniMax model temperature is clamped to [0, 1]
func TestMiniMaxTemperatureClamping(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		input    float32
		expected float32
	}{
		{"normal temperature", "MiniMax-M2.7", 0.9, 0.9},
		{"temperature too high", "MiniMax-M2.7", 1.5, 1.0},
		{"temperature at max", "MiniMax-M2.5-highspeed", 1.0, 1.0},
		{"temperature at zero", "MiniMax-M2.5", 0.0, 0.0},
		{"negative temperature", "MiniMax-M2.7-highspeed", -0.5, 0.0},
		{"extreme high temp", "MiniMax-M2.7", 2.0, 1.0},
		{"case insensitive", "minimax-m2.7", 2.0, 1.0},
		{"non-minimax not clamped", "gpt-4o", 1.5, 1.5},
		{"deepseek not clamped", "deepseek-chat", 1.5, 1.5},
		{"claude not clamped", "claude-3-5-sonnet-20240620", 1.8, 1.8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clampMiniMaxTemperature(tt.model, tt.input)
			if result != tt.expected {
				t.Errorf("clampMiniMaxTemperature(%q, %v) = %v, want %v",
					tt.model, tt.input, result, tt.expected)
			}
		})
	}
}

// TestMiniMaxModelPrefix tests MiniMax model identification
func TestMiniMaxModelPrefix(t *testing.T) {
	minimax := []string{
		"MiniMax-M2.7", "MiniMax-M2.7-highspeed",
		"MiniMax-M2.5", "MiniMax-M2.5-highspeed", "minimax-m2.7",
	}
	nonMinimax := []string{
		"gpt-4o", "gpt-4o-mini", "claude-3-5-sonnet-20240620",
		"deepseek-chat", "qwen-plus", "o1-mini",
	}

	for _, m := range minimax {
		if !strings.HasPrefix(strings.ToLower(m), "minimax") {
			t.Errorf("%q should be identified as MiniMax", m)
		}
	}
	for _, m := range nonMinimax {
		if strings.HasPrefix(strings.ToLower(m), "minimax") {
			t.Errorf("%q should NOT be identified as MiniMax", m)
		}
	}
}

// TestMiniMaxSSEParsing tests parsing MiniMax OpenAI-compatible streaming response
func TestMiniMaxSSEParsing(t *testing.T) {
	sseData := []string{
		`data: {"id":"chatcmpl-test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" from"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" MiniMax"},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":"stop"}]}`,
		`data: [DONE]`,
	}

	var contents []string
	for _, line := range sseData {
		if !strings.Contains(line, "data:") || len(line) < 30 {
			continue
		}
		if strings.Contains(line, "[DONE]") {
			break
		}

		var resp types.ApiResponse
		if err := json.Unmarshal([]byte(line[6:]), &resp); err != nil {
			t.Fatalf("parse error: %v", err)
		}

		if len(resp.Choices) == 0 {
			continue
		}
		if resp.Choices[0].Delta.Content != nil {
			c := fmt.Sprintf("%v", resp.Choices[0].Delta.Content)
			if c != "" {
				contents = append(contents, c)
			}
		}
	}

	if got := strings.Join(contents, ""); got != "Hello from MiniMax!" {
		t.Errorf("expected 'Hello from MiniMax!', got %q", got)
	}
}

// TestMiniMaxReasoningContent tests reasoning_content parsing for MiniMax M2.5
func TestMiniMaxReasoningContent(t *testing.T) {
	sseData := []string{
		`data: {"id":"chatcmpl-test","choices":[{"index":0,"delta":{"role":"assistant","reasoning_content":"Let me think..."},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","choices":[{"index":0,"delta":{"reasoning_content":" about this."},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","choices":[{"index":0,"delta":{"content":"The answer is 42."},"finish_reason":null}]}`,
		`data: {"id":"chatcmpl-test","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
	}

	var reasoning, content []string
	for _, line := range sseData {
		if !strings.Contains(line, "data:") || len(line) < 30 {
			continue
		}
		var resp types.ApiResponse
		if err := json.Unmarshal([]byte(line[6:]), &resp); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if len(resp.Choices) == 0 {
			continue
		}
		d := resp.Choices[0].Delta
		if d.ReasoningContent != "" {
			reasoning = append(reasoning, d.ReasoningContent)
		}
		if d.Content != nil {
			c := fmt.Sprintf("%v", d.Content)
			if c != "" {
				content = append(content, c)
			}
		}
	}

	if got := strings.Join(reasoning, ""); got != "Let me think... about this." {
		t.Errorf("reasoning: expected 'Let me think... about this.', got %q", got)
	}
	if got := strings.Join(content, ""); got != "The answer is 42." {
		t.Errorf("content: expected 'The answer is 42.', got %q", got)
	}
}

// TestMiniMaxSSEServerIntegration simulates a MiniMax-compatible SSE endpoint
func TestMiniMaxSSEServerIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Error("missing Bearer auth")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("wrong content-type")
		}

		var req types.ApiRequest
		json.NewDecoder(r.Body).Decode(&req)

		if !strings.HasPrefix(strings.ToLower(req.Model), "minimax") {
			t.Errorf("expected minimax model, got %s", req.Model)
		}
		if req.Temperature < 0 || req.Temperature > 1.0 {
			t.Errorf("temp out of range: %f", req.Temperature)
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		chunks := []string{
			`data: {"id":"mm","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
			`data: {"id":"mm","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
			`data: {"id":"mm","choices":[{"index":0,"delta":{"content":" World"},"finish_reason":null}]}`,
			`data: {"id":"mm","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
			`data: [DONE]`,
		}
		for _, c := range chunks {
			fmt.Fprintf(w, "%s\n\n", c)
		}
	}))
	defer server.Close()

	body := types.ApiRequest{
		Model:       "MiniMax-M2.7",
		Temperature: 0.9,
		MaxTokens:   4096,
		Stream:      true,
		Messages:    []any{map[string]string{"role": "user", "content": "Hi"}},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", server.URL+"/v1/chat/completions", strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var contents []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "data:") || len(line) < 30 {
			continue
		}
		if strings.Contains(line, "[DONE]") {
			break
		}
		var apiResp types.ApiResponse
		if err := json.Unmarshal([]byte(line[6:]), &apiResp); err != nil {
			continue
		}
		if len(apiResp.Choices) > 0 && apiResp.Choices[0].Delta.Content != nil {
			c := fmt.Sprintf("%v", apiResp.Choices[0].Delta.Content)
			if c != "" {
				contents = append(contents, c)
			}
		}
	}

	if got := strings.Join(contents, ""); got != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", got)
	}
}

// TestMiniMaxNonStreamingIntegration tests MiniMax non-streaming response
func TestMiniMaxNonStreamingIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "chatcmpl-minimax",
			"model": "MiniMax-M2.7",
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]string{"role": "assistant", "content": "Hello from MiniMax!"},
				"finish_reason": "stop",
			}},
			"usage": map[string]int{
				"prompt_tokens":     10,
				"completion_tokens": 5,
				"total_tokens":      15,
			},
		})
	}))
	defer server.Close()

	body := types.ApiRequest{
		Model:       "MiniMax-M2.7",
		Temperature: 0.9,
		MaxTokens:   4096,
		Stream:      false,
		Messages:    []any{map[string]string{"role": "user", "content": "Hello"}},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", server.URL, strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}
	json.NewDecoder(resp.Body).Decode(&respBody)

	if len(respBody.Choices) == 0 || respBody.Choices[0].Message.Content != "Hello from MiniMax!" {
		t.Errorf("unexpected response content")
	}
	if respBody.Usage.TotalTokens != 15 {
		t.Errorf("expected 15 tokens, got %d", respBody.Usage.TotalTokens)
	}
}

// TestMiniMaxAPIEndpoint tests URL construction for MiniMax API
func TestMiniMaxAPIEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		apiURL   string
		expected string
	}{
		{"base URL", "https://api.minimax.io", "https://api.minimax.io/v1/chat/completions"},
		{"full URL", "https://api.minimax.io/v1/chat/completions", "https://api.minimax.io/v1/chat/completions"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := strings.TrimPrefix(tt.apiURL, "https://")
			p = strings.TrimPrefix(p, "http://")
			var result string
			if !strings.Contains(p, "/") {
				result = fmt.Sprintf("%s/v1/chat/completions", tt.apiURL)
			} else {
				result = tt.apiURL
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
