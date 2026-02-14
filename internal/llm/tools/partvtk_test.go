package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-03: PartVTK Tool (VTK Export)
// Unit tests only — Docker integration tests are in partvtk_docker_test.go

func TestPartVTKTool_Info(t *testing.T) {
	tool := NewPartVTKTool()
	info := tool.Info()

	assert.Equal(t, "partvtk", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "out_file")
	assert.Contains(t, info.Parameters, "only_fluid")
	assert.Contains(t, info.Parameters, "vars")
	assert.Contains(t, info.Parameters, "first")
	assert.Contains(t, info.Parameters, "last")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "out_file")
}

func TestPartVTKTool_Run(t *testing.T) {
	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewPartVTKTool()
		call := ToolCall{
			Name:  "partvtk",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing data directory", func(t *testing.T) {
		tool := NewPartVTKTool()
		params := PartVTKParams{
			DataDir: "/nonexistent/path",
			OutFile: "/tmp/partvtk_test_out/fluid",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
