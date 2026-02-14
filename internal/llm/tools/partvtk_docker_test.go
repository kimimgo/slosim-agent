//go:build integration

package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-03: PartVTK Docker Integration Tests
// Run with: go test -tags integration ./internal/llm/tools/ -run TestPartVTKDocker

func TestPartVTKDocker_ConvertBi4ToVTK(t *testing.T) {
	tool := NewPartVTKTool()
	params := PartVTKParams{
		DataDir: "/tmp/partvtk_test_data",
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
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, ".vtk")
}

func TestPartVTKDocker_OnlyFluid(t *testing.T) {
	tool := NewPartVTKTool()
	onlyFluid := true
	params := PartVTKParams{
		DataDir:   "/tmp/partvtk_test_data",
		OutFile:   "/tmp/partvtk_test_out/fluid",
		OnlyFluid: &onlyFluid,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "partvtk",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestPartVTKDocker_VariableSelection(t *testing.T) {
	tool := NewPartVTKTool()
	params := PartVTKParams{
		DataDir: "/tmp/partvtk_test_data",
		OutFile: "/tmp/partvtk_test_out/fluid",
		Vars:    []string{"vel", "press"},
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "partvtk",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestPartVTKDocker_TimestepRange(t *testing.T) {
	first := 10
	last := 20
	tool := NewPartVTKTool()
	params := PartVTKParams{
		DataDir: "/tmp/partvtk_test_data",
		OutFile: "/tmp/partvtk_test_out/fluid",
		First:   &first,
		Last:    &last,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "partvtk",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}
