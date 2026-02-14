package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-02: DualSPHysics Solver Tool (GPU SPH Simulation)
// Unit tests only — Docker integration tests are in solver_docker_test.go

func TestSolverTool_Info(t *testing.T) {
	tool := NewSolverTool()
	info := tool.Info()

	assert.Equal(t, "solver", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "case_name")
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "out_dir")
	assert.Contains(t, info.Parameters, "gpu")
	assert.Contains(t, info.Required, "case_name")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "out_dir")
}

func TestSolverTool_Run(t *testing.T) {
	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewSolverTool()
		call := ToolCall{
			Name:  "solver",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing required parameters", func(t *testing.T) {
		tool := NewSolverTool()
		params := SolverParams{
			CaseName: "",
			DataDir:  "",
			OutDir:   "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "solver",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
