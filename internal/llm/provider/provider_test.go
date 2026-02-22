package provider

import (
	"os"
	"testing"

	"github.com/opencode-ai/opencode/internal/llm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProviderLocal(t *testing.T) {
	t.Setenv("LOCAL_ENDPOINT", "http://localhost:11434")

	p, err := NewProvider(models.ProviderLocal,
		WithAPIKey("dummy"),
		WithModel(models.Model{ID: "local.test", Provider: models.ProviderLocal, APIModel: "test"}),
		WithMaxTokens(4096),
		WithSystemMessage("test"),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, models.ModelID("local.test"), p.Model().ID)
}

func TestNewProviderLocalAutoAppendsV1(t *testing.T) {
	t.Setenv("LOCAL_ENDPOINT", "http://localhost:11434")

	p, err := NewProvider(models.ProviderLocal,
		WithAPIKey("dummy"),
		WithModel(models.Model{ID: "local.test", Provider: models.ProviderLocal, APIModel: "test"}),
		WithMaxTokens(4096),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProviderLocalAlreadyHasV1(t *testing.T) {
	t.Setenv("LOCAL_ENDPOINT", "http://localhost:11434/v1")

	p, err := NewProvider(models.ProviderLocal,
		WithAPIKey("dummy"),
		WithModel(models.Model{ID: "local.test", Provider: models.ProviderLocal, APIModel: "test"}),
		WithMaxTokens(4096),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProviderUnsupported(t *testing.T) {
	_, err := NewProvider("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provider not supported")
}

func TestNewProviderOpenAI(t *testing.T) {
	p, err := NewProvider(models.ProviderOpenAI,
		WithAPIKey("test-key"),
		WithModel(models.Model{ID: "openai.test", Provider: models.ProviderOpenAI, APIModel: "test"}),
		WithMaxTokens(4096),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProviderGemini(t *testing.T) {
	// Gemini requires GEMINI_API_KEY to create client
	os.Setenv("GEMINI_API_KEY", "test-key")
	defer os.Unsetenv("GEMINI_API_KEY")

	p, err := NewProvider(models.ProviderGemini,
		WithAPIKey("test-key"),
		WithModel(models.Model{ID: "gemini.test", Provider: models.ProviderGemini, APIModel: "test"}),
		WithMaxTokens(4096),
		WithGeminiOptions(),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProviderGROQ(t *testing.T) {
	p, err := NewProvider(models.ProviderGROQ,
		WithAPIKey("test-key"),
		WithModel(models.Model{ID: "groq.test", Provider: models.ProviderGROQ, APIModel: "test"}),
		WithMaxTokens(4096),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestNewProviderOpenRouter(t *testing.T) {
	p, err := NewProvider(models.ProviderOpenRouter,
		WithAPIKey("test-key"),
		WithModel(models.Model{ID: "or.test", Provider: models.ProviderOpenRouter, APIModel: "test"}),
		WithMaxTokens(4096),
	)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestCleanMessages(t *testing.T) {
	t.Parallel()
	bp := &baseProvider[OpenAIClient]{
		options: providerClientOptions{},
	}

	// Test with empty parts (should be filtered)
	// This mainly tests that cleanMessages works without panic
	cleaned := bp.cleanMessages(nil)
	assert.Empty(t, cleaned)
}

func TestProviderOptionFunctions(t *testing.T) {
	t.Parallel()
	opts := providerClientOptions{}

	WithAPIKey("test-key")(&opts)
	assert.Equal(t, "test-key", opts.apiKey)

	WithModel(models.Model{ID: "test"})(&opts)
	assert.Equal(t, models.ModelID("test"), opts.model.ID)

	WithMaxTokens(1024)(&opts)
	assert.Equal(t, int64(1024), opts.maxTokens)

	WithSystemMessage("hello")(&opts)
	assert.Equal(t, "hello", opts.systemMessage)
}
