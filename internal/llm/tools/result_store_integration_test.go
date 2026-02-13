package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// STORE-01 Integration Tests: Result Store Tool
// Tests save→search→compare workflow, JSON persistence, and concurrent access.
// ============================================================================

// --- Full Workflow: Save → Get → Search → Compare ---

func TestResultStore_Integration_FullWorkflow(t *testing.T) {
	store := createTestStore(t)

	// Step 1: Save multiple results
	results := []*SimulationResult{
		{
			ID: "sim_fill50", Name: "Fill Ratio 50%", CaseFile: "case_001.xml",
			OutputDir: "/data/sim_fill50", Status: "completed", Duration: 120.5,
			Tags: []string{"sloshing", "benchmark", "fill50"},
			Parameters: map[string]any{"fill_ratio": 0.5, "frequency": 1.0},
		},
		{
			ID: "sim_fill60", Name: "Fill Ratio 60%", CaseFile: "case_002.xml",
			OutputDir: "/data/sim_fill60", Status: "completed", Duration: 135.2,
			Tags: []string{"sloshing", "benchmark", "fill60"},
			Parameters: map[string]any{"fill_ratio": 0.6, "frequency": 1.0},
		},
		{
			ID: "sim_fill70", Name: "Fill Ratio 70%", CaseFile: "case_003.xml",
			OutputDir: "/data/sim_fill70", Status: "failed", Duration: 45.0,
			Tags: []string{"sloshing", "fill70"},
			Parameters: map[string]any{"fill_ratio": 0.7, "frequency": 1.0},
		},
	}

	for _, result := range results {
		response := runResultStoreAction(t, store, "save", result.ID, result, nil)
		assert.False(t, response.IsError, "save %s should succeed", result.ID)
		assert.Contains(t, response.Content, "결과 저장 완료")
	}

	// Step 2: Get each result back and verify data integrity
	for _, expected := range results {
		response := runResultStoreAction(t, store, "get", expected.ID, nil, nil)
		assert.False(t, response.IsError)

		var retrieved SimulationResult
		err := json.Unmarshal([]byte(response.Content), &retrieved)
		require.NoError(t, err, "get response should be valid JSON for %s", expected.ID)

		assert.Equal(t, expected.ID, retrieved.ID)
		assert.Equal(t, expected.Name, retrieved.Name)
		assert.Equal(t, expected.Status, retrieved.Status)
		assert.Equal(t, expected.Duration, retrieved.Duration)
		assert.Equal(t, expected.CaseFile, retrieved.CaseFile)
	}

	// Step 3: List all results
	listResponse := runResultStoreAction(t, store, "list", "", nil, nil)
	assert.False(t, listResponse.IsError)
	assert.Contains(t, listResponse.Content, "3개의 결과가 저장")

	// Step 4: Search by tag
	searchResponse := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{
		Tags:  []string{"benchmark"},
		Limit: 10,
	})
	assert.False(t, searchResponse.IsError)
	assert.Contains(t, searchResponse.Content, "sim_fill50")
	assert.Contains(t, searchResponse.Content, "sim_fill60")
	// sim_fill70 doesn't have "benchmark" tag
	assert.NotContains(t, searchResponse.Content, "sim_fill70",
		"sim_fill70 lacks 'benchmark' tag and should not appear in search results")

	// Step 5: Search by status
	failedSearch := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{
		Status: "failed",
	})
	assert.False(t, failedSearch.IsError)
	assert.Contains(t, failedSearch.Content, "sim_fill70")
	assert.NotContains(t, failedSearch.Content, "sim_fill50")

	// Step 6: Compare results
	compareResponse := runResultStoreAction(t, store, "compare", "sim_fill50,sim_fill60,sim_fill70", nil, nil)
	assert.False(t, compareResponse.IsError)
	assert.Contains(t, compareResponse.Content, "결과 비교")
}

// --- JSON Persistence ---

func TestResultStore_Integration_Persistence(t *testing.T) {
	t.Run("STORE-01: saved data persists across tool instances", func(t *testing.T) {
		tempDir := t.TempDir()
		storeDir := filepath.Join(tempDir, "store")

		// First instance: save
		store1 := &resultStoreTool{storeDir: storeDir}
		result := &SimulationResult{
			ID: "persist_test", Name: "Persistence Test", CaseFile: "test.xml",
			OutputDir: "/data/persist", Status: "completed",
			Tags: []string{"persist"},
		}
		response := runResultStoreAction(t, store1, "save", "", result, nil)
		assert.False(t, response.IsError)

		// Second instance: read back
		store2 := &resultStoreTool{storeDir: storeDir}
		getResponse := runResultStoreAction(t, store2, "get", "persist_test", nil, nil)
		assert.False(t, getResponse.IsError)
		assert.Contains(t, getResponse.Content, "Persistence Test")
	})

	t.Run("STORE-01: saved JSON file has correct structure", func(t *testing.T) {
		tempDir := t.TempDir()
		storeDir := filepath.Join(tempDir, "store")

		store := &resultStoreTool{storeDir: storeDir}
		result := &SimulationResult{
			ID: "json_check", Name: "JSON Structure Test", CaseFile: "test.xml",
			OutputDir: "/data/json", Status: "completed", Duration: 99.9,
			Tags:       []string{"test", "structure"},
			Parameters: map[string]any{"fill_ratio": 0.5},
		}
		runResultStoreAction(t, store, "save", "", result, nil)

		// Read raw JSON file
		data, err := os.ReadFile(filepath.Join(storeDir, "json_check.json"))
		require.NoError(t, err)

		var parsed map[string]any
		require.NoError(t, json.Unmarshal(data, &parsed))

		assert.Equal(t, "json_check", parsed["id"])
		assert.Equal(t, "JSON Structure Test", parsed["name"])
		assert.Equal(t, "completed", parsed["status"])
		assert.Equal(t, 99.9, parsed["duration"])
		assert.NotEmpty(t, parsed["timestamp"], "timestamp should be auto-generated")

		tags, ok := parsed["tags"].([]any)
		require.True(t, ok)
		assert.Len(t, tags, 2)
	})

	t.Run("STORE-01: auto-generates ID when not provided", func(t *testing.T) {
		store := createTestStore(t)
		result := &SimulationResult{
			Name: "Auto ID Test", CaseFile: "test.xml",
			OutputDir: "/data/auto", Status: "completed",
		}

		response := runResultStoreAction(t, store, "save", "", result, nil)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "sim_", "auto-generated ID should start with sim_")
	})

	t.Run("STORE-01: auto-generates timestamp when not provided", func(t *testing.T) {
		store := createTestStore(t)
		result := &SimulationResult{
			ID: "ts_test", Name: "Timestamp Test", CaseFile: "test.xml",
			OutputDir: "/data/ts", Status: "completed",
		}

		runResultStoreAction(t, store, "save", "", result, nil)

		// Read back and check timestamp
		getResp := runResultStoreAction(t, store, "get", "ts_test", nil, nil)
		var retrieved SimulationResult
		require.NoError(t, json.Unmarshal([]byte(getResp.Content), &retrieved))
		assert.NotEmpty(t, retrieved.Timestamp, "timestamp should be auto-populated")
	})
}

// --- Concurrent Access ---

func TestResultStore_Integration_ConcurrentAccess(t *testing.T) {
	t.Run("STORE-01: concurrent saves don't corrupt data", func(t *testing.T) {
		store := createTestStore(t)
		numWorkers := 10

		var wg sync.WaitGroup
		errors := make([]error, numWorkers)

		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				result := &SimulationResult{
					ID:        fmt.Sprintf("concurrent_%03d", idx),
					Name:      fmt.Sprintf("Concurrent Test %d", idx),
					CaseFile:  fmt.Sprintf("case_%03d.xml", idx),
					OutputDir: fmt.Sprintf("/data/concurrent_%03d", idx),
					Status:    "completed",
					Tags:      []string{"concurrent"},
				}

				params := ResultStoreParams{
					Action:     "save",
					ResultData: result,
				}

				paramsJSON, err := json.Marshal(params)
				if err != nil {
					errors[idx] = err
					return
				}

				call := ToolCall{
					Name:  ResultStoreToolName,
					Input: string(paramsJSON),
				}

				response, err := store.Run(context.Background(), call)
				if err != nil {
					errors[idx] = err
					return
				}
				if response.IsError {
					errors[idx] = fmt.Errorf("save error: %s", response.Content)
				}
			}(i)
		}

		wg.Wait()

		// Check no errors occurred
		for i, err := range errors {
			assert.NoError(t, err, "worker %d should not error", i)
		}

		// Verify all results are readable
		listResp := runResultStoreAction(t, store, "list", "", nil, nil)
		assert.Contains(t, listResp.Content, fmt.Sprintf("%d개의 결과가 저장", numWorkers))
	})

	t.Run("STORE-01: concurrent reads during writes are safe", func(t *testing.T) {
		store := createTestStore(t)

		// Pre-populate some data
		for i := 0; i < 5; i++ {
			result := &SimulationResult{
				ID: fmt.Sprintf("preload_%03d", i), Name: fmt.Sprintf("Preload %d", i),
				CaseFile: "test.xml", OutputDir: "/data", Status: "completed",
				Tags: []string{"preload"},
			}
			runResultStoreAction(t, store, "save", "", result, nil)
		}

		// Concurrent reads and writes
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			// Writer
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				result := &SimulationResult{
					ID: fmt.Sprintf("write_%03d", idx), Name: fmt.Sprintf("Write %d", idx),
					CaseFile: "test.xml", OutputDir: "/data", Status: "completed",
				}
				runResultStoreAction(t, store, "save", "", result, nil)
			}(i)

			// Reader
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				runResultStoreAction(t, store, "get", fmt.Sprintf("preload_%03d", idx), nil, nil)
			}(i)
		}

		wg.Wait()
		// If we get here without panics or data corruption, the test passes
	})
}

// --- Search Edge Cases ---

func TestResultStore_Integration_SearchEdgeCases(t *testing.T) {
	store := createTestStore(t)

	// Setup test data
	testResults := []*SimulationResult{
		{ID: "s1", Name: "Alpha", Status: "completed", Tags: []string{"a", "b"}},
		{ID: "s2", Name: "Beta", Status: "completed", Tags: []string{"b", "c"}},
		{ID: "s3", Name: "Gamma", Status: "failed", Tags: []string{"c", "d"}},
		{ID: "s4", Name: "Delta", Status: "running", Tags: []string{"d", "e"}},
		{ID: "s5", Name: "Alpha", Status: "completed", Tags: []string{"a", "e"}},
	}
	for _, r := range testResults {
		r.CaseFile = "test.xml"
		r.OutputDir = "/data"
		runResultStoreAction(t, store, "save", "", r, nil)
	}

	t.Run("STORE-01: search by name returns exact matches", func(t *testing.T) {
		resp := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{Name: "Alpha"})
		assert.Contains(t, resp.Content, "s1")
		assert.Contains(t, resp.Content, "s5")
		assert.NotContains(t, resp.Content, "s2")
	})

	t.Run("STORE-01: search limit is respected", func(t *testing.T) {
		resp := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{Limit: 2})
		// Count occurrences of "id" in the response to estimate result count
		count := strings.Count(resp.Content, `"id"`)
		assert.LessOrEqual(t, count, 2, "limit=2 should return at most 2 results")
	})

	t.Run("STORE-01: search with no matches returns empty", func(t *testing.T) {
		resp := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{
			Tags: []string{"nonexistent_tag"},
		})
		assert.Contains(t, resp.Content, "검색 결과 0개")
	})

	t.Run("STORE-01: search with multiple tag filters (OR logic)", func(t *testing.T) {
		resp := runResultStoreAction(t, store, "search", "", nil, &ResultQuery{
			Tags: []string{"a", "d"},
		})
		// "a" matches s1,s5; "d" matches s3,s4
		assert.Contains(t, resp.Content, "s1")
		assert.Contains(t, resp.Content, "s3")
	})
}

// --- Delete ---

func TestResultStore_Integration_DeleteWorkflow(t *testing.T) {
	t.Run("STORE-01: delete then get returns error", func(t *testing.T) {
		store := createTestStore(t)

		result := &SimulationResult{
			ID: "delete_me", Name: "To Delete", CaseFile: "test.xml",
			OutputDir: "/data", Status: "completed",
		}
		runResultStoreAction(t, store, "save", "", result, nil)

		// Delete
		delResp := runResultStoreAction(t, store, "delete", "delete_me", nil, nil)
		assert.False(t, delResp.IsError)
		assert.Contains(t, delResp.Content, "삭제 완료")

		// Get should fail
		getResp := runResultStoreAction(t, store, "get", "delete_me", nil, nil)
		assert.True(t, getResp.IsError, "get after delete should return error")
	})

	t.Run("STORE-01: delete non-existent result returns error", func(t *testing.T) {
		store := createTestStore(t)
		resp := runResultStoreAction(t, store, "delete", "nonexistent_id", nil, nil)
		assert.True(t, resp.IsError)
	})
}

// --- Compare (RED TEST - Known bug) ---

func TestResultStore_Integration_Compare(t *testing.T) {
	t.Run("RED/STORE-01: compare should correctly split comma-separated IDs", func(t *testing.T) {
		// BUG: handleCompare iterates over individual runes instead of splitting by comma.
		// "sim1,sim2" produces individual characters, not ["sim1", "sim2"].
		store := createTestStore(t)

		// Save two results
		for _, id := range []string{"sim_a", "sim_b"} {
			result := &SimulationResult{
				ID: id, Name: id, CaseFile: "test.xml",
				OutputDir: "/data/" + id, Status: "completed",
			}
			runResultStoreAction(t, store, "save", "", result, nil)
		}

		resp := runResultStoreAction(t, store, "compare", "sim_a,sim_b", nil, nil)
		assert.False(t, resp.IsError)

		// Expected: 2 IDs being compared
		// BUG: current implementation splits by rune, resulting in many single-char "IDs"
		assert.Contains(t, resp.Content, "비교 대상 ID 수: 2",
			"RED: compare should split 'sim_a,sim_b' into 2 IDs, not individual characters")
	})

	t.Run("RED/STORE-01: compare should load actual result data for comparison", func(t *testing.T) {
		store := createTestStore(t)

		result1 := &SimulationResult{
			ID: "cmp_1", Name: "Compare Test 1", CaseFile: "case1.xml",
			OutputDir: "/data/cmp1", Status: "completed", Duration: 100.0,
			Parameters: map[string]any{"fill_ratio": 0.5},
		}
		result2 := &SimulationResult{
			ID: "cmp_2", Name: "Compare Test 2", CaseFile: "case2.xml",
			OutputDir: "/data/cmp2", Status: "completed", Duration: 200.0,
			Parameters: map[string]any{"fill_ratio": 0.7},
		}
		runResultStoreAction(t, store, "save", "", result1, nil)
		runResultStoreAction(t, store, "save", "", result2, nil)

		resp := runResultStoreAction(t, store, "compare", "cmp_1,cmp_2", nil, nil)

		// Compare should include actual parameter/duration differences
		assert.Contains(t, resp.Content, "Compare Test 1",
			"RED: compare should include result names in comparison output")
		assert.Contains(t, resp.Content, "Compare Test 2",
			"RED: compare should include result names in comparison output")
	})
}

// --- Error Handling ---

func TestResultStore_Integration_Errors(t *testing.T) {
	t.Run("STORE-01: save without result_data returns error", func(t *testing.T) {
		store := createTestStore(t)
		resp := runResultStoreAction(t, store, "save", "", nil, nil)
		assert.True(t, resp.IsError)
		assert.Contains(t, resp.Content, "result_data")
	})

	t.Run("STORE-01: get without result_id returns error", func(t *testing.T) {
		store := createTestStore(t)
		resp := runResultStoreAction(t, store, "get", "", nil, nil)
		assert.True(t, resp.IsError)
		assert.Contains(t, resp.Content, "result_id")
	})

	t.Run("STORE-01: search without query returns error", func(t *testing.T) {
		store := createTestStore(t)
		resp := runResultStoreAction(t, store, "search", "", nil, nil)
		assert.True(t, resp.IsError)
		assert.Contains(t, resp.Content, "query")
	})

	t.Run("STORE-01: invalid JSON input returns error", func(t *testing.T) {
		store := createTestStore(t)
		call := ToolCall{
			Name:  ResultStoreToolName,
			Input: "not valid json",
		}
		resp, err := store.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, resp.IsError)
	})
}

// ============================================================================
// Test Helpers
// ============================================================================

func createTestStore(t *testing.T) *resultStoreTool {
	t.Helper()
	tempDir := t.TempDir()
	return &resultStoreTool{
		storeDir: filepath.Join(tempDir, "store"),
	}
}

func runResultStoreAction(t *testing.T, store *resultStoreTool, action, resultID string, resultData *SimulationResult, query *ResultQuery) ToolResponse {
	t.Helper()

	params := ResultStoreParams{
		Action:     action,
		ResultID:   resultID,
		ResultData: resultData,
		Query:      query,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  ResultStoreToolName,
		Input: string(paramsJSON),
	}

	response, err := store.Run(context.Background(), call)
	require.NoError(t, err, "tool.Run should not return Go error")
	return response
}
