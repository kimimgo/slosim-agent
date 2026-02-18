package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractInputSchema_ValidSchema(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "The name parameter",
			},
			"count": map[string]any{
				"type":        "integer",
				"description": "The count parameter",
			},
		},
		"required": []any{"name"},
	}

	props, required := extractInputSchema(schema)

	assert.Len(t, props, 2)
	assert.Contains(t, props, "name")
	assert.Contains(t, props, "count")
	assert.Equal(t, []string{"name"}, required)
}

func TestExtractInputSchema_EmptySchema(t *testing.T) {
	props, required := extractInputSchema(nil)

	assert.Empty(t, props)
	assert.Nil(t, required)
}

func TestExtractInputSchema_NoProperties(t *testing.T) {
	schema := map[string]any{
		"type": "object",
	}

	props, required := extractInputSchema(schema)

	assert.NotNil(t, props)
	assert.Empty(t, props)
	assert.Nil(t, required)
}

func TestExtractInputSchema_NoRequired(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	props, required := extractInputSchema(schema)

	assert.Len(t, props, 1)
	assert.Nil(t, required)
}

func TestExtractInputSchema_MultipleRequired(t *testing.T) {
	schema := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []any{"a", "b", "c"},
	}

	_, required := extractInputSchema(schema)

	require.Len(t, required, 3)
	assert.Equal(t, []string{"a", "b", "c"}, required)
}

func TestExtractInputSchema_NonMapSchema(t *testing.T) {
	// go-sdk returns any — could be a string or other type in edge cases
	props, required := extractInputSchema("not a map")

	assert.NotNil(t, props)
	assert.Empty(t, props)
	assert.Nil(t, required)
}

func TestSessionPool_InitialState(t *testing.T) {
	assert.NotNil(t, pool)
	assert.NotNil(t, pool.client)
	assert.NotNil(t, pool.sessions)
}
