//go:build integration

package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-01: GenCase Docker Integration Tests
// These tests require Docker with dsph-agent image.
// Run with: go test -tags integration ./internal/llm/tools/ -run TestGenCaseDocker

func TestGenCaseDocker_ValidXML(t *testing.T) {
	tool := NewGenCaseTool()
	params := GenCaseParams{
		CasePath: "cases/SloshingTank_Def",
		SavePath: "/tmp/gencase_test_out",
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "gencase",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, ".bi4")
}

func TestGenCaseDocker_InvalidXML(t *testing.T) {
	tool := NewGenCaseTool()
	params := GenCaseParams{
		CasePath: "/nonexistent/invalid_case",
		SavePath: "/tmp/gencase_test_out",
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
	assert.NotEmpty(t, response.Content)
}

func TestGenCaseDocker_AutoStripXMLExtension(t *testing.T) {
	tool := NewGenCaseTool()
	params := GenCaseParams{
		CasePath: "cases/SloshingTank_Def.xml",
		SavePath: "/tmp/gencase_test_out",
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "gencase",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestGenCaseDocker_DockerExecution(t *testing.T) {
	tool := NewGenCaseTool()
	params := GenCaseParams{
		CasePath: "cases/SloshingTank_Def",
		SavePath: "/tmp/gencase_test_out",
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "gencase",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestGenCaseDocker_DPOverride(t *testing.T) {
	dpValue := 0.02
	tool := NewGenCaseTool()
	params := GenCaseParams{
		CasePath: "cases/SloshingTank_Def",
		SavePath: "/tmp/gencase_test_out",
		DP:       &dpValue,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "gencase",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}
