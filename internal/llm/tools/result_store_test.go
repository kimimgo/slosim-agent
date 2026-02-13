package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// STORE-01: Result Store (DB)

func TestResultStoreTool_Info(t *testing.T) {
	tool := NewResultStoreTool()
	info := tool.Info()

	assert.Equal(t, ResultStoreToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "action")
	assert.Contains(t, info.Parameters, "result_id")
	assert.Contains(t, info.Parameters, "result_data")
	assert.Contains(t, info.Parameters, "query")
	assert.Contains(t, info.Required, "action")
}

func TestResultStoreTool_Run(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "result_store_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override store directory for testing
	tool := &resultStoreTool{
		storeDir: filepath.Join(tempDir, "store"),
	}

	t.Run("STORE-01: saves simulation result", func(t *testing.T) {
		result := &SimulationResult{
			ID:         "test_sim_001",
			Name:       "Test Simulation",
			CaseFile:   "case.xml",
			OutputDir:  "/tmp/sim_out",
			Status:     "completed",
			Duration:   123.45,
			Tags:       []string{"test", "benchmark"},
		}

		params := ResultStoreParams{
			Action:     "save",
			ResultData: result,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "결과 저장 완료")
		assert.Contains(t, response.Content, "test_sim_001")

		// Verify file was created
		metadataPath := filepath.Join(tool.storeDir, "test_sim_001.json")
		_, err = os.Stat(metadataPath)
		assert.NoError(t, err)
	})

	t.Run("STORE-01: retrieves saved result", func(t *testing.T) {
		params := ResultStoreParams{
			Action:   "get",
			ResultID: "test_sim_001",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "test_sim_001")
		assert.Contains(t, response.Content, "Test Simulation")
	})

	t.Run("STORE-01: lists all results", func(t *testing.T) {
		params := ResultStoreParams{
			Action: "list",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "개의 결과가 저장")
	})

	t.Run("STORE-01: searches results by tags", func(t *testing.T) {
		query := &ResultQuery{
			Tags:  []string{"test"},
			Limit: 10,
		}

		params := ResultStoreParams{
			Action: "search",
			Query:  query,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "검색 결과")
	})

	t.Run("STORE-01: deletes result", func(t *testing.T) {
		params := ResultStoreParams{
			Action:   "delete",
			ResultID: "test_sim_001",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "삭제 완료")

		// Verify file was deleted
		metadataPath := filepath.Join(tool.storeDir, "test_sim_001.json")
		_, err = os.Stat(metadataPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("STORE-01: handles invalid action", func(t *testing.T) {
		params := ResultStoreParams{
			Action: "invalid_action",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "알 수 없는 액션")
	})
}

func TestContainsAnyTag(t *testing.T) {
	t.Run("finds matching tag", func(t *testing.T) {
		resultTags := []string{"test", "benchmark", "sloshing"}
		queryTags := []string{"benchmark"}

		assert.True(t, containsAnyTag(resultTags, queryTags))
	})

	t.Run("no matching tags", func(t *testing.T) {
		resultTags := []string{"test", "benchmark"}
		queryTags := []string{"production", "validated"}

		assert.False(t, containsAnyTag(resultTags, queryTags))
	})

	t.Run("empty tags", func(t *testing.T) {
		resultTags := []string{}
		queryTags := []string{"test"}

		assert.False(t, containsAnyTag(resultTags, queryTags))
	})
}
