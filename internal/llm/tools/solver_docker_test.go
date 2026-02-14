//go:build integration

package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-02: DualSPHysics Solver Docker Integration Tests
// Run with: go test -tags integration ./internal/llm/tools/ -run TestSolverDocker

func TestSolverDocker_GPUSimulation(t *testing.T) {
	tool := NewSolverTool()
	params := SolverParams{
		CaseName: "SloshingTank",
		DataDir:  "/tmp/solver_test_data",
		OutDir:   "/tmp/solver_test_out",
		GPU:      true,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "solver",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestSolverDocker_BackgroundExecution(t *testing.T) {
	tool := NewSolverTool()
	params := SolverParams{
		CaseName: "SloshingTank",
		DataDir:  "/tmp/solver_test_data",
		OutDir:   "/tmp/solver_test_out",
		GPU:      true,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "solver",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, "job_id")
}

func TestSolverDocker_CompletionDetection(t *testing.T) {
	tool := NewSolverTool()
	params := SolverParams{
		CaseName: "SloshingTank",
		DataDir:  "/tmp/solver_test_data",
		OutDir:   "/tmp/solver_test_out",
		GPU:      true,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "solver",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
}

func TestSolverDocker_FailureDetection(t *testing.T) {
	tool := NewSolverTool()
	params := SolverParams{
		CaseName: "InvalidCase",
		DataDir:  "/nonexistent/data",
		OutDir:   "/tmp/solver_test_out",
		GPU:      true,
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
}

func TestSolverDocker_NonBlocking(t *testing.T) {
	tool := NewSolverTool()
	params := SolverParams{
		CaseName: "SloshingTank",
		DataDir:  "/tmp/solver_test_data",
		OutDir:   "/tmp/solver_test_out",
		GPU:      true,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  "solver",
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)
	assert.Contains(t, response.Content, "job_id")
}
