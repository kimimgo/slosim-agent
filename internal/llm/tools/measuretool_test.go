package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-04: MeasureTool (Water Level / Pressure Measurement)
// Unit tests only — Docker integration tests are in measuretool_docker_test.go

func TestMeasureToolTool_Info(t *testing.T) {
	tool := NewMeasureToolTool()
	info := tool.Info()

	assert.Equal(t, "measuretool", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "points_file")
	assert.Contains(t, info.Parameters, "out_csv")
	assert.Contains(t, info.Parameters, "vars")
	assert.Contains(t, info.Parameters, "elevation")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "points_file")
	assert.Contains(t, info.Required, "out_csv")
}

func TestMeasureToolTool_Run(t *testing.T) {
	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewMeasureToolTool()
		call := ToolCall{
			Name:  "measuretool",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing points file", func(t *testing.T) {
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/nonexistent/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
