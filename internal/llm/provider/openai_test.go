package provider

import (
	"testing"

	"github.com/openai/openai-go"
	"github.com/opencode-ai/opencode/internal/llm/models"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTemperature(t *testing.T) {
	t.Parallel()
	opts := openaiOptions{}
	WithTemperature(0)(&opts)
	require.NotNil(t, opts.temperature)
	assert.Equal(t, float64(0), *opts.temperature)
}

func TestWithMaxTools(t *testing.T) {
	t.Parallel()
	opts := openaiOptions{}
	WithMaxTools(5)(&opts)
	assert.Equal(t, 5, opts.maxTools)
}

func TestWithReasoningEffort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"low", "low", "low"},
		{"medium", "medium", "medium"},
		{"high", "high", "high"},
		{"invalid falls back to medium", "invalid", "medium"},
		{"empty falls back to medium", "", "medium"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opts := openaiOptions{}
			WithReasoningEffort(tt.input)(&opts)
			assert.Equal(t, tt.expected, opts.reasoningEffort)
		})
	}
}

func TestWithOpenAIBaseURL(t *testing.T) {
	t.Parallel()
	opts := openaiOptions{}
	WithOpenAIBaseURL("http://localhost:11434/v1")(&opts)
	assert.Equal(t, "http://localhost:11434/v1", opts.baseURL)
}

func TestWithOpenAIExtraHeaders(t *testing.T) {
	t.Parallel()
	headers := map[string]string{"X-Custom": "value"}
	opts := openaiOptions{}
	WithOpenAIExtraHeaders(headers)(&opts)
	assert.Equal(t, "value", opts.extraHeaders["X-Custom"])
}

func TestPreparedParamsMaxToolsLimit(t *testing.T) {
	t.Parallel()
	client := &openaiClient{
		providerOptions: providerClientOptions{
			model:     models.Model{APIModel: "test-model"},
			maxTokens: 4096,
		},
		options: openaiOptions{
			maxTools: 3,
		},
	}

	tools := make([]openai.ChatCompletionToolParam, 5)
	for i := range tools {
		tools[i] = openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name: "tool_" + string(rune('a'+i)),
			},
		}
	}

	params := client.preparedParams(nil, tools)
	assert.Len(t, params.Tools, 3, "tools should be truncated to maxTools")
}

func TestPreparedParamsNoMaxToolsLimit(t *testing.T) {
	t.Parallel()
	client := &openaiClient{
		providerOptions: providerClientOptions{
			model:     models.Model{APIModel: "test-model"},
			maxTokens: 4096,
		},
		options: openaiOptions{
			maxTools: 0,
		},
	}

	tools := make([]openai.ChatCompletionToolParam, 10)
	params := client.preparedParams(nil, tools)
	assert.Len(t, params.Tools, 10, "tools should not be truncated when maxTools is 0")
}

func TestPreparedParamsTemperature(t *testing.T) {
	t.Parallel()
	temp := float64(0)
	client := &openaiClient{
		providerOptions: providerClientOptions{
			model:     models.Model{APIModel: "test-model"},
			maxTokens: 4096,
		},
		options: openaiOptions{
			temperature: &temp,
		},
	}

	params := client.preparedParams(nil, nil)
	assert.NotEmpty(t, params.Temperature)
}

func TestPreparedParamsNoTemperature(t *testing.T) {
	t.Parallel()
	client := &openaiClient{
		providerOptions: providerClientOptions{
			model:     models.Model{APIModel: "test-model"},
			maxTokens: 4096,
		},
		options: openaiOptions{},
	}

	params := client.preparedParams(nil, nil)
	assert.Empty(t, params.Temperature)
}

func TestPreparedParamsCanReason(t *testing.T) {
	t.Parallel()
	client := &openaiClient{
		providerOptions: providerClientOptions{
			model:     models.Model{APIModel: "test-model", CanReason: true},
			maxTokens: 4096,
		},
		options: openaiOptions{
			reasoningEffort: "high",
		},
	}

	params := client.preparedParams(nil, nil)
	assert.NotEmpty(t, params.MaxCompletionTokens)
	assert.Empty(t, params.MaxTokens)
}

func TestFinishReason(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected message.FinishReason
	}{
		{"stop", "stop", message.FinishReasonEndTurn},
		{"length", "length", message.FinishReasonMaxTokens},
		{"tool_calls", "tool_calls", message.FinishReasonToolUse},
		{"unknown", "other", message.FinishReasonUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := &openaiClient{}
			result := client.finishReason(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
