package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// v1.0 E2E Scenario Tests
//
// These tests validate complete user workflows that span multiple tools.
// Each scenario represents a real-world usage pattern of slosim-agent v1.0.
// ============================================================================

// Scenario 1: Parametric Study → Save Results → Compare
//
// User wants to study 3 fill ratios (0.5, 0.6, 0.7), save all results,
// then compare them to find the optimal configuration.
func TestE2E_ParametricStudy_SaveAndCompare(t *testing.T) {
	tempDir := t.TempDir()

	// === Phase 1: Setup base case ===
	baseCasePath := filepath.Join(tempDir, "sloshing_tank.xml")
	baseCaseContent := `<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="1000" />
            <cflnumber value="0.2" />
        </constantsdef>
        <geometry>
            <definition dp="0.01" />
        </geometry>
    </casedef>
    <parameters>
        <fill_ratio value="{{fill_ratio}}" />
    </parameters>
</case>`
	require.NoError(t, os.WriteFile(baseCasePath, []byte(baseCaseContent), 0644))

	// === Phase 2: Run parametric study with 3 fill ratios ===
	studyTool := tools.NewParametricStudyTool()
	studyOutputDir := filepath.Join(tempDir, "parametric_study")

	studyParams := tools.ParametricStudyParams{
		StudyName: "fill_ratio_comparison",
		BaseCase:  baseCasePath,
		Variables: []tools.ParametricVariable{
			{Name: "fill_ratio", Values: []float64{0.5, 0.6, 0.7}},
		},
		OutputDir: studyOutputDir,
	}

	studyParamsJSON, err := json.Marshal(studyParams)
	require.NoError(t, err)

	studyResponse, err := studyTool.Run(context.Background(), tools.ToolCall{
		Name:  tools.ParametricStudyToolName,
		Input: string(studyParamsJSON),
	})
	require.NoError(t, err)
	assert.False(t, studyResponse.IsError, "parametric study should succeed")
	assert.Contains(t, studyResponse.Content, "총 케이스 수: 3")
	assert.Contains(t, studyResponse.Content, "총 3개 케이스 생성 완료")

	// Verify manifest was created
	manifestPath := filepath.Join(studyOutputDir, "study_manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	require.NoError(t, err)
	var manifest map[string]any
	require.NoError(t, json.Unmarshal(manifestData, &manifest))
	assert.Equal(t, float64(3), manifest["total_cases"])

	// Verify 3 case directories with XML files
	for i := 1; i <= 3; i++ {
		caseXML := filepath.Join(studyOutputDir, fmt.Sprintf("case_%03d", i), "case.xml")
		_, err := os.Stat(caseXML)
		assert.NoError(t, err, "case_%03d/case.xml should exist", i)
	}

	// === Phase 3: Save each case result to result store ===
	storeDir := filepath.Join(tempDir, "result_store")
	store := tools.NewResultStoreToolWithDir(storeDir)

	fillRatios := []float64{0.5, 0.6, 0.7}
	resultIDs := make([]string, 3)

	for i, fr := range fillRatios {
		resultID := fmt.Sprintf("fill_ratio_%02d", int(fr*100))
		resultIDs[i] = resultID

		saveParams := tools.ResultStoreParams{
			Action: "save",
			ResultData: &tools.SimulationResult{
				ID:        resultID,
				Name:      fmt.Sprintf("Fill Ratio %.0f%%", fr*100),
				CaseFile:  filepath.Join(studyOutputDir, fmt.Sprintf("case_%03d", i+1), "case.xml"),
				OutputDir: filepath.Join(studyOutputDir, fmt.Sprintf("case_%03d", i+1)),
				Status:    "completed",
				Duration:  float64(100 + i*20),
				Tags:      []string{"sloshing", "parametric", fmt.Sprintf("fill%.0f", fr*100)},
				Parameters: map[string]any{
					"fill_ratio": fr,
					"frequency":  1.0,
				},
			},
		}

		saveJSON, err := json.Marshal(saveParams)
		require.NoError(t, err)

		saveResp, err := store.Run(context.Background(), tools.ToolCall{
			Name:  tools.ResultStoreToolName,
			Input: string(saveJSON),
		})
		require.NoError(t, err)
		assert.False(t, saveResp.IsError, "saving result %s should succeed", resultID)
	}

	// === Phase 4: Search results by tag ===
	searchParams := tools.ResultStoreParams{
		Action: "search",
		Query: &tools.ResultQuery{
			Tags:  []string{"parametric"},
			Limit: 10,
		},
	}
	searchJSON, err := json.Marshal(searchParams)
	require.NoError(t, err)

	searchResp, err := store.Run(context.Background(), tools.ToolCall{
		Name:  tools.ResultStoreToolName,
		Input: string(searchJSON),
	})
	require.NoError(t, err)
	assert.False(t, searchResp.IsError)
	assert.Contains(t, searchResp.Content, "검색 결과 3개")

	// === Phase 5: Compare all 3 results ===
	compareParams := tools.ResultStoreParams{
		Action:   "compare",
		ResultID: strings.Join(resultIDs, ","),
	}
	compareJSON, err := json.Marshal(compareParams)
	require.NoError(t, err)

	compareResp, err := store.Run(context.Background(), tools.ToolCall{
		Name:  tools.ResultStoreToolName,
		Input: string(compareJSON),
	})
	require.NoError(t, err)
	assert.False(t, compareResp.IsError)
	assert.Contains(t, compareResp.Content, "결과 비교")
}

// Scenario 2: Simulation Error → Auto-Detect → Suggest Fix → Retry
//
// User runs a simulation that fails with particle-out-of-bounds error.
// Error recovery detects the issue, suggests fixes, and auto-fix is applied.
func TestE2E_SimulationError_DetectFixRetry(t *testing.T) {
	tempDir := t.TempDir()

	// === Phase 1: Simulate a failed simulation output ===
	simOutputDir := filepath.Join(tempDir, "sim_output")
	require.NoError(t, os.MkdirAll(simOutputDir, 0755))

	// Create DualSPHysics log with particle-out-of-bounds error
	logContent := `[DSPH v5.4.023] DualSPHysics5.4_linux64
GPU: NVIDIA GeForce RTX 4090 (sm_89)
Loading case from /data/case_001...
  Number of particles: 500000
Starting simulation...
  Step:     0  Time: 0.000000  Particles: 500000
  Step:  1000  Time: 0.100000  Particles: 500000
  Step:  2000  Time: 0.200000  Particles: 499995
*** Error: 5 particles out of domain bounds at step 2001
  Particle IDs: [100234, 100235, 100236, 100237, 100238]
  Max velocity: 15.2 m/s (exceeds CFL limit)
Simulation terminated.
`
	require.NoError(t, os.WriteFile(filepath.Join(simOutputDir, "DualSPHysics.log"), []byte(logContent), 0644))

	// Create Run.csv with normal data (not divergent)
	csvContent := "Time,TotalMass,Energy\n"
	for i := 0; i < 2000; i++ {
		csvContent += fmt.Sprintf("%.4f,500000,%.1f\n", float64(i)*0.0001, 1000.0+float64(i)*0.01)
	}
	require.NoError(t, os.WriteFile(filepath.Join(simOutputDir, "Run.csv"), []byte(csvContent), 0644))

	// === Phase 2: Run error recovery with auto-detect ===
	recoveryTool := tools.NewErrorRecoveryTool()
	recoveryParams := tools.ErrorRecoveryParams{
		JobID:           "failed_sim_001",
		OutputDir:       simOutputDir,
		CheckDivergence: true,
		AutoFix:         false, // First, just analyze
	}

	recoveryJSON, err := json.Marshal(recoveryParams)
	require.NoError(t, err)

	detectResp, err := recoveryTool.Run(context.Background(), tools.ToolCall{
		Name:  tools.ErrorRecoveryToolName,
		Input: string(recoveryJSON),
	})
	require.NoError(t, err)
	assert.False(t, detectResp.IsError)

	// Should detect the particle-out-of-bounds error
	assert.Contains(t, detectResp.Content, "에러 목록", "should list detected errors")
	assert.Contains(t, detectResp.Content, "제안된 수정 방법", "should suggest fixes")
	assert.Contains(t, detectResp.Content, "도메인 크기 증가",
		"should specifically suggest domain expansion for particle-out-of-bounds")

	// === Phase 3: Run again with auto-fix enabled ===
	recoveryParams.AutoFix = true
	recoveryJSON, err = json.Marshal(recoveryParams)
	require.NoError(t, err)

	fixResp, err := recoveryTool.Run(context.Background(), tools.ToolCall{
		Name:  tools.ErrorRecoveryToolName,
		Input: string(recoveryJSON),
	})
	require.NoError(t, err)
	assert.False(t, fixResp.IsError)
	assert.Contains(t, fixResp.Content, "자동 수정 적용됨", "auto-fix should be applied")
	assert.Contains(t, fixResp.Content, "재시도 권장", "should recommend retry after fix")

	// === Phase 4: Save failed result to result store for tracking ===
	storeDir := filepath.Join(tempDir, "result_store")
	store := tools.NewResultStoreToolWithDir(storeDir)

	saveParams := tools.ResultStoreParams{
		Action: "save",
		ResultData: &tools.SimulationResult{
			ID:          "failed_sim_001",
			Name:        "Failed Sloshing Simulation",
			CaseFile:    "case_001.xml",
			OutputDir:   simOutputDir,
			Status:      "failed",
			Duration:    20.5,
			Tags:        []string{"sloshing", "failed", "particle-out"},
			Description: "Particle out of bounds at step 2001",
		},
	}

	saveJSON, err := json.Marshal(saveParams)
	require.NoError(t, err)

	saveResp, err := store.Run(context.Background(), tools.ToolCall{
		Name:  tools.ResultStoreToolName,
		Input: string(saveJSON),
	})
	require.NoError(t, err)
	assert.False(t, saveResp.IsError)

	// === Phase 5: Verify the failure is searchable ===
	searchParams := tools.ResultStoreParams{
		Action: "search",
		Query: &tools.ResultQuery{
			Status: "failed",
		},
	}
	searchJSON, err := json.Marshal(searchParams)
	require.NoError(t, err)

	searchResp, err := store.Run(context.Background(), tools.ToolCall{
		Name:  tools.ResultStoreToolName,
		Input: string(searchJSON),
	})
	require.NoError(t, err)
	assert.Contains(t, searchResp.Content, "failed_sim_001",
		"failed simulation should be searchable by status=failed")
}

// Scenario 3: Multi-variable parametric study with mixed success/failure
//
// User runs a 2-variable study (fill_ratio x frequency).
// Some cases succeed, some fail. Results are saved with appropriate status.
func TestE2E_ParametricStudy_MixedResults(t *testing.T) {
	tempDir := t.TempDir()

	// Setup base case
	baseCasePath := filepath.Join(tempDir, "base.xml")
	require.NoError(t, os.WriteFile(baseCasePath, []byte(`<?xml version="1.0"?>
<case><parameters><fill value="0.5" /></parameters></case>`), 0644))

	// Run parametric study
	studyTool := tools.NewParametricStudyTool()
	studyOutputDir := filepath.Join(tempDir, "study")

	studyParams := tools.ParametricStudyParams{
		StudyName: "mixed_study",
		BaseCase:  baseCasePath,
		Variables: []tools.ParametricVariable{
			{Name: "fill_ratio", Values: []float64{0.5, 0.7}},
			{Name: "frequency", Values: []float64{0.8, 1.2}},
		},
		OutputDir: studyOutputDir,
	}

	studyJSON, _ := json.Marshal(studyParams)
	studyResp, err := studyTool.Run(context.Background(), tools.ToolCall{
		Name:  tools.ParametricStudyToolName,
		Input: string(studyJSON),
	})
	require.NoError(t, err)
	assert.Contains(t, studyResp.Content, "총 케이스 수: 4")

	// Save results with mixed statuses
	storeDir := filepath.Join(tempDir, "store")
	store := tools.NewResultStoreToolWithDir(storeDir)

	statuses := []string{"completed", "completed", "failed", "completed"}
	for i := 1; i <= 4; i++ {
		result := &tools.SimulationResult{
			ID:        fmt.Sprintf("mixed_%03d", i),
			Name:      fmt.Sprintf("Mixed Case %d", i),
			CaseFile:  filepath.Join(studyOutputDir, fmt.Sprintf("case_%03d", i), "case.xml"),
			OutputDir: filepath.Join(studyOutputDir, fmt.Sprintf("case_%03d", i)),
			Status:    statuses[i-1],
			Tags:      []string{"mixed_study", statuses[i-1]},
		}

		saveParams := tools.ResultStoreParams{Action: "save", ResultData: result}
		saveJSON, _ := json.Marshal(saveParams)
		resp, err := store.Run(context.Background(), tools.ToolCall{
			Name: tools.ResultStoreToolName, Input: string(saveJSON),
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
	}

	// Search for completed results only
	searchParams := tools.ResultStoreParams{
		Action: "search",
		Query:  &tools.ResultQuery{Status: "completed"},
	}
	searchJSON, _ := json.Marshal(searchParams)
	searchResp, err := store.Run(context.Background(), tools.ToolCall{
		Name: tools.ResultStoreToolName, Input: string(searchJSON),
	})
	require.NoError(t, err)
	assert.Contains(t, searchResp.Content, "검색 결과 3개",
		"should find exactly 3 completed results")

	// Search for failed results
	failedParams := tools.ResultStoreParams{
		Action: "search",
		Query:  &tools.ResultQuery{Status: "failed"},
	}
	failedJSON, _ := json.Marshal(failedParams)
	failedResp, err := store.Run(context.Background(), tools.ToolCall{
		Name: tools.ResultStoreToolName, Input: string(failedJSON),
	})
	require.NoError(t, err)
	assert.Contains(t, failedResp.Content, "검색 결과 1개",
		"should find exactly 1 failed result")
}
