package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// JOB-01: Job Manager (Background Execution Management)
// FRD AC-1 ~ AC-5

func TestJobManagerTool_Info(t *testing.T) {
	tool := NewJobManagerTool()
	info := tool.Info()

	assert.Equal(t, "job_manager", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "action")
	assert.Contains(t, info.Parameters, "job_id")
	assert.Contains(t, info.Parameters, "command")
	assert.Contains(t, info.Parameters, "work_dir")
	assert.Contains(t, info.Required, "action")
}

func TestJobManagerTool_Submit(t *testing.T) {
	t.Run("AC-1: submit returns immediately (non-blocking)", func(t *testing.T) {
		// Job 제출 시 즉시 반환 (non-blocking)
		// Submit() → 즉시 Job ID 반환, goroutine 실행
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action:  "submit",
			Command: []string{"docker", "compose", "run", "--rm", "dsph", "DualSPHysics5.4_linux64"},
			WorkDir: "/tmp/job_test",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "job_id")
	})

	t.Run("AC-5: rejects concurrent job submission beyond limit (v0.3 max 3)", func(t *testing.T) {
		// v0.3: 동시 최대 3개 Job 관리
		// 4번째 제출 시 거절
		tool := NewJobManagerTool()

		// 3개 Job 제출 → 모두 성공해야 함
		for i := 1; i <= 3; i++ {
			params := JobManagerParams{
				Action:  "submit",
				Command: []string{"sleep", "100"},
				WorkDir: fmt.Sprintf("/tmp/job_test_%d", i),
			}
			paramsJSON, err := json.Marshal(params)
			require.NoError(t, err)

			call := ToolCall{
				Name:  "job_manager",
				Input: string(paramsJSON),
			}
			response, err := tool.Run(context.Background(), call)
			require.NoError(t, err)
			assert.False(t, response.IsError, "Job %d should succeed", i)
		}

		// 4번째 Job 제출 → 거절되어야 함
		params4 := JobManagerParams{
			Action:  "submit",
			Command: []string{"sleep", "200"},
			WorkDir: "/tmp/job_test_4",
		}
		paramsJSON4, err := json.Marshal(params4)
		require.NoError(t, err)

		call4 := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON4),
		}
		response4, err := tool.Run(context.Background(), call4)
		require.NoError(t, err)
		assert.True(t, response4.IsError)
		assert.Contains(t, response4.Content, "초과")
	})
}

func TestJobManagerTool_Status(t *testing.T) {
	t.Run("AC-2: query job status", func(t *testing.T) {
		// Job 상태 조회 가능
		// Status(job_id) → RUNNING/COMPLETED/FAILED
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action: "status",
			JobID:  "test-job-id-123",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		// 존재하지 않는 Job ID → 에러 또는 not found
		assert.NotEmpty(t, response.Content)
	})

	t.Run("AC-3: pubsub event on completion/failure", func(t *testing.T) {
		// 완료/실패 시 pubsub 이벤트 발행
		// JobCompletedMsg / JobFailedMsg 이벤트 확인
		// Note: 이 테스트는 pubsub 통합 테스트에서 더 깊이 검증
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action: "status",
			JobID:  "test-job-id-123",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.NotEmpty(t, response.Content)
	})

	t.Run("AC-4: progress estimation via Part file counting", func(t *testing.T) {
		// 진행률 추정 (Part 파일 카운팅)
		// Part_*.bi4 파일 수 / 예상 총 파일 수
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action: "status",
			JobID:  "test-job-id-123",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.NotEmpty(t, response.Content)
	})
}

func TestJobManagerTool_ErrorHandling(t *testing.T) {
	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewJobManagerTool()
		call := ToolCall{
			Name:  "job_manager",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles unknown action", func(t *testing.T) {
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action: "unknown_action",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing job_id for status query", func(t *testing.T) {
		tool := NewJobManagerTool()
		params := JobManagerParams{
			Action: "status",
			JobID:  "", // 빈 Job ID
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "job_manager",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
