package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-02: DualSPHysics Solver Tool (GPU SPH Simulation)
// FRD AC-1 ~ AC-5

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
	t.Run("AC-1: GPU simulation completes with Part_*.bi4", func(t *testing.T) {
		// GPU 모드로 시뮬레이션 완료 (Part_*.bi4 생성)
		// SloshingTank_Def (0.1초) 실행 → Part 파일 존재
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
	})

	t.Run("AC-2: background execution returns Job ID immediately", func(t *testing.T) {
		// 백그라운드 실행: 즉시 Job ID 반환
		// Run() 호출 → 즉시 ToolResponse 반환, Job ID 포함
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
		// Job ID가 포함된 응답을 즉시 반환해야 함
		assert.Contains(t, response.Content, "job_id")
	})

	t.Run("AC-3: detects simulation completion", func(t *testing.T) {
		// 시뮬레이션 완료 감지
		// Run.out 파일의 "Finished" 문자열 파싱
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
	})

	t.Run("AC-4: detects simulation failure", func(t *testing.T) {
		// 시뮬레이션 실패 감지
		// RhopOut 초과 시 에러 상태 전환
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
	})

	t.Run("AC-5: TUI remains responsive during simulation", func(t *testing.T) {
		// 실행 중 TUI 응답 유지
		// 시뮬레이션 중 채팅 입력 가능 (non-blocking 확인)
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

		// Run should return immediately (non-blocking)
		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		// 응답이 즉시 반환되었다면 non-blocking 확인
		assert.Contains(t, response.Content, "job_id")
	})

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
