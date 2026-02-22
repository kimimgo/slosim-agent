package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedModelsContainsGemini(t *testing.T) {
	t.Parallel()
	model, ok := SupportedModels[Gemini25Flash]
	assert.True(t, ok, "Gemini25Flash should be in SupportedModels")
	assert.Equal(t, ProviderGemini, model.Provider)
	assert.Equal(t, "gemini-2.5-flash", model.APIModel)
}

func TestSupportedModelsContainsOpenAI(t *testing.T) {
	t.Parallel()
	model, ok := SupportedModels[GPT41]
	assert.True(t, ok, "GPT41 should be in SupportedModels")
	assert.Equal(t, ProviderOpenAI, model.Provider)
}

func TestSupportedModelsContainsOpenRouter(t *testing.T) {
	t.Parallel()
	model, ok := SupportedModels[OpenRouterClaude37Sonnet]
	assert.True(t, ok, "OpenRouterClaude37Sonnet should be in SupportedModels")
	assert.Equal(t, ProviderOpenRouter, model.Provider)
	assert.True(t, model.CanReason)
}

func TestSupportedModelsDoesNotContainRemovedProviders(t *testing.T) {
	t.Parallel()
	for id, model := range SupportedModels {
		assert.NotEqual(t, ModelProvider("anthropic"), model.Provider,
			"anthropic provider should not be in SupportedModels: %s", id)
		assert.NotEqual(t, ModelProvider("copilot"), model.Provider,
			"copilot provider should not be in SupportedModels: %s", id)
		assert.NotEqual(t, ModelProvider("bedrock"), model.Provider,
			"bedrock provider should not be in SupportedModels: %s", id)
		assert.NotEqual(t, ModelProvider("azure"), model.Provider,
			"azure provider should not be in SupportedModels: %s", id)
		assert.NotEqual(t, ModelProvider("vertexai"), model.Provider,
			"vertexai provider should not be in SupportedModels: %s", id)
	}
}

func TestProviderPopularityContainsActiveProviders(t *testing.T) {
	t.Parallel()
	expectedProviders := []ModelProvider{
		ProviderOpenAI,
		ProviderGemini,
		ProviderGROQ,
		ProviderOpenRouter,
		ProviderXAI,
		ProviderLocal,
	}
	for _, p := range expectedProviders {
		_, ok := ProviderPopularity[p]
		assert.True(t, ok, "provider %s should be in ProviderPopularity", p)
	}
}

func TestProviderPopularityDoesNotContainRemovedProviders(t *testing.T) {
	t.Parallel()
	removedProviders := []ModelProvider{
		"anthropic", "copilot", "bedrock", "azure", "vertexai",
	}
	for _, p := range removedProviders {
		_, ok := ProviderPopularity[p]
		assert.False(t, ok, "removed provider %s should not be in ProviderPopularity", p)
	}
}

func TestGeminiModelIDsAreStable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id       ModelID
		expected string
	}{
		{Gemini25Flash, "gemini-2.5-flash"},
		{Gemini25, "gemini-2.5-pro"},
		{Gemini20Flash, "gemini-2.0-flash"},
		{Gemini20FlashLite, "gemini-2.0-flash-lite"},
	}
	for _, tt := range tests {
		model, ok := GeminiModels[tt.id]
		assert.True(t, ok, "model %s should exist", tt.id)
		assert.Equal(t, tt.expected, model.APIModel, "APIModel for %s should be stable (no preview suffix)", tt.id)
	}
}

func TestOpenRouterGeminiModelIDsAreStable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id       ModelID
		expected string
	}{
		{OpenRouterGemini25Flash, "google/gemini-2.5-flash"},
		{OpenRouterGemini25, "google/gemini-2.5-pro"},
	}
	for _, tt := range tests {
		model, ok := OpenRouterModels[tt.id]
		assert.True(t, ok, "model %s should exist", tt.id)
		assert.Equal(t, tt.expected, model.APIModel, "APIModel for %s should be stable", tt.id)
	}
}
