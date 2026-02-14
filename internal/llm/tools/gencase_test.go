package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-01: GenCase Tool (XML → Particle Pre-processing)
// Unit tests only — Docker integration tests are in gencase_docker_test.go

func TestGenCaseTool_Info(t *testing.T) {
	tool := NewGenCaseTool()
	info := tool.Info()

	assert.Equal(t, "gencase", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "case_path")
	assert.Contains(t, info.Parameters, "save_path")
	assert.Contains(t, info.Parameters, "dp")
	assert.Contains(t, info.Required, "case_path")
	assert.Contains(t, info.Required, "save_path")
}

func TestGenCaseTool_Run(t *testing.T) {
	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewGenCaseTool()
		call := ToolCall{
			Name:  "gencase",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing required parameters", func(t *testing.T) {
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "",
			SavePath: "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
