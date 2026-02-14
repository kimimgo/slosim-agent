//go:build integration

package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-04: MeasureTool Docker Integration Tests
// Run with: go test -tags integration ./internal/llm/tools/ -run TestMeasureToolDocker

func TestMeasureToolDocker_ProbePointsMeasurement(t *testing.T) {
	tool := NewMeasureToolTool()
	params := MeasureToolParams{
		DataDir:    "/tmp/measuretool_test_data",
		PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
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
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, ".csv")
}

func TestMeasureToolDocker_ElevationMode(t *testing.T) {
	elevation := true
	tool := NewMeasureToolTool()
	params := MeasureToolParams{
		DataDir:    "/tmp/measuretool_test_data",
		PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
		OutCSV:     "/tmp/measuretool_test_out/measure",
		Elevation:  &elevation,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "measuretool",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, "Elevation")
}

func TestMeasureToolDocker_CSVFormat(t *testing.T) {
	tool := NewMeasureToolTool()
	params := MeasureToolParams{
		DataDir:    "/tmp/measuretool_test_data",
		PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
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
	assert.False(t, response.IsError)
}

func TestMeasureToolDocker_CustomVariables(t *testing.T) {
	tool := NewMeasureToolTool()
	params := MeasureToolParams{
		DataDir:    "/tmp/measuretool_test_data",
		PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
		OutCSV:     "/tmp/measuretool_test_out/measure",
		Vars:       []string{"vel", "press"},
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "measuretool",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}
