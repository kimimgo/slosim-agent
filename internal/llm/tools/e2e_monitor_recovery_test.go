//go:build e2e

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// E2E Test: Monitor + Error Recovery Workflow (Task #8)
//
// Tests the monitor → error_recovery pipeline using synthetic simulation data.
// No Docker or GPU required — validates tool logic with realistic Run.csv data.
//
// Run with: go test -tags e2e -timeout 60s ./internal/llm/tools/ -run TestE2E_Monitor -v
// =============================================================================

// --- Test 1: Monitor reads healthy simulation progress ---

func TestE2E_MonitorHealthySimulation(t *testing.T) {
	simDir := setupSyntheticSimulation(t, syntheticHealthy)

	monitor := NewMonitorTool()
	params := MonitorParams{
		SimDir:  simDir,
		TimeMax: 2.0,
	}
	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := monitor.Run(context.Background(), ToolCall{
		Name:  MonitorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "정상적으로 진행 중")
	assert.Contains(t, resp.Content, "progress_pct")

	// Verify progress is approximately 50% (1.0 / 2.0)
	var status MonitorStatus
	extractJSON(t, resp.Content, &status)
	assert.InDelta(t, 50.0, status.ProgressPct, 1.0,
		"progress should be ~50%% at time 1.0 with TimeMax 2.0")
	assert.Equal(t, 10000, status.ParticleCount)
	assert.False(t, status.IsUnstable)
}

// --- Test 2: Monitor detects unstable simulation (mass particle loss) ---

func TestE2E_MonitorUnstableSimulation(t *testing.T) {
	simDir := setupSyntheticSimulation(t, syntheticUnstable)

	monitor := NewMonitorTool()
	params := MonitorParams{
		SimDir:  simDir,
		TimeMax: 2.0,
	}
	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := monitor.Run(context.Background(), ToolCall{
		Name:  MonitorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "불안정")
	assert.Contains(t, resp.Content, "이탈")
}

// --- Test 3: Error recovery detects divergence in Run.csv ---

func TestE2E_ErrorRecoveryDivergence(t *testing.T) {
	simDir := setupSyntheticSimulation(t, syntheticDivergent)

	recovery := NewErrorRecoveryTool()
	params := ErrorRecoveryParams{
		JobID:           "e2e_divergent",
		OutputDir:       simDir,
		CheckDivergence: true,
	}
	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := recovery.Run(context.Background(), ToolCall{
		Name:  ErrorRecoveryToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "발산 감지: true")
	assert.Contains(t, resp.Content, "TimeStep 감소")
}

// --- Test 4: Full workflow — monitor detects instability, error_recovery suggests fix ---

func TestE2E_MonitorToRecoveryWorkflow(t *testing.T) {
	// Simulate a failing simulation: particle loss + error log + divergent energy
	simDir := t.TempDir()

	// 1. Write Run.csv with divergent energy AND particle loss
	runCSV := generateDivergentWithParticleLossCSV(200)
	require.NoError(t, os.WriteFile(filepath.Join(simDir, "Run.csv"), []byte(runCSV), 0644))

	// 2. Write error log
	logContent := `[DSPH v5.4.023] DualSPHysics5.4_linux64
GPU: NVIDIA GeForce RTX 4090 (sm_89)
Loading case...
  Number of particles: 10000
Starting simulation...
  Step:     0  Time: 0.000000
  Step:   100  Time: 0.010000
*** Error: 500 particles out of domain bounds at step 150
*** Error: NaN detected in velocity field at step 151
Simulation terminated.
`
	require.NoError(t, os.WriteFile(filepath.Join(simDir, "DualSPHysics.log"), []byte(logContent), 0644))

	// 3. Write minimal XML for TimeMax auto-detection
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<case>
  <execution>
    <parameters>
      <parameter key="TimeMax" value="2.0" />
    </parameters>
  </execution>
</case>`
	require.NoError(t, os.WriteFile(filepath.Join(simDir, "case.xml"), []byte(xmlContent), 0644))

	// --- Step A: Monitor detects instability ---
	t.Run("step_A_monitor_detects_instability", func(t *testing.T) {
		monitor := NewMonitorTool()
		params := MonitorParams{SimDir: simDir}
		paramsJSON, _ := json.Marshal(params)

		resp, err := monitor.Run(context.Background(), ToolCall{
			Name:  MonitorToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)

		// Monitor should detect instability from particle loss
		assert.Contains(t, resp.Content, "불안정")

		// Monitor should auto-detect TimeMax from XML
		var status MonitorStatus
		extractJSON(t, resp.Content, &status)
		assert.Greater(t, status.ProgressPct, 0.0,
			"progress should be > 0 when TimeMax is auto-detected")
	})

	// --- Step B: Error recovery analyzes and suggests fixes ---
	t.Run("step_B_recovery_suggests_fixes", func(t *testing.T) {
		recovery := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			JobID:           "e2e_workflow",
			OutputDir:       simDir,
			CheckDivergence: true,
		}
		paramsJSON, _ := json.Marshal(params)

		resp, err := recovery.Run(context.Background(), ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)

		// Should detect divergence from Run.csv
		assert.Contains(t, resp.Content, "발산 감지: true")

		// Should detect log errors
		assert.Contains(t, resp.Content, "에러 목록")

		// Should suggest multiple fix strategies
		assert.Contains(t, resp.Content, "제안된 수정 방법")
		assert.Contains(t, resp.Content, "도메인 크기 증가", "particle-out should suggest domain expansion")
		assert.Contains(t, resp.Content, "NaN", "NaN error should be detected")

		// Should recommend retry
		assert.Contains(t, resp.Content, "재시도 권장")
	})

	// --- Step C: Error recovery with auto-fix enabled ---
	t.Run("step_C_autofix_applies_first_fix", func(t *testing.T) {
		recovery := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			JobID:     "e2e_autofix",
			OutputDir: simDir,
			AutoFix:   true,
		}
		paramsJSON, _ := json.Marshal(params)

		resp, err := recovery.Run(context.Background(), ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.Contains(t, resp.Content, "자동 수정 적용됨")
		assert.Contains(t, resp.Content, "재시도 권장")
	})
}

// --- Test 5: Monitor auto-detects TimeMax from XML ---

func TestE2E_MonitorTimeMaxAutoDetect(t *testing.T) {
	simDir := t.TempDir()

	// Write Run.csv
	runCSV := `#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;100.0;50.0
0.500;5000;10000;8000;2000;0;102.0;51.0
1.000;10000;10000;8000;2000;0;101.0;50.5
`
	require.NoError(t, os.WriteFile(filepath.Join(simDir, "Run.csv"), []byte(runCSV), 0644))

	// Write XML with TimeMax=5.0
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<case>
  <execution>
    <parameters>
      <parameter key="TimeMax" value="5.0" />
      <parameter key="TimeOut" value="0.01" />
    </parameters>
  </execution>
</case>`
	require.NoError(t, os.WriteFile(filepath.Join(simDir, "simulation.xml"), []byte(xmlContent), 0644))

	monitor := NewMonitorTool()
	params := MonitorParams{
		SimDir: simDir,
		// TimeMax not set → should auto-detect 5.0 from XML
	}
	paramsJSON, _ := json.Marshal(params)

	resp, err := monitor.Run(context.Background(), ToolCall{
		Name:  MonitorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)

	var status MonitorStatus
	extractJSON(t, resp.Content, &status)
	// Current time 1.0 / TimeMax 5.0 = 20%
	assert.InDelta(t, 20.0, status.ProgressPct, 1.0,
		"should auto-detect TimeMax=5.0 and calculate 20%% progress")
}

// =============================================================================
// Synthetic Simulation Data Generators
// =============================================================================

type syntheticSimType int

const (
	syntheticHealthy   syntheticSimType = iota
	syntheticUnstable                   // Mass particle loss
	syntheticDivergent                  // Exponential energy growth
)

func setupSyntheticSimulation(t *testing.T, simType syntheticSimType) string {
	t.Helper()
	simDir := t.TempDir()

	var csvContent string
	switch simType {
	case syntheticHealthy:
		csvContent = `#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;100.0;50.0
0.100;1000;10000;8000;2000;0;101.0;50.5
0.200;2000;10000;8000;2000;0;102.0;51.0
0.300;3000;10000;8000;2000;0;101.5;50.8
0.400;4000;10000;8000;2000;0;100.8;50.3
0.500;5000;10000;8000;2000;0;101.2;50.6
0.600;6000;10000;8000;2000;0;100.5;50.2
0.700;7000;10000;8000;2000;0;101.0;50.5
0.800;8000;10000;8000;2000;0;100.3;50.1
0.900;9000;10000;8000;2000;0;101.1;50.5
1.000;10000;10000;8000;2000;0;100.7;50.3
`
	case syntheticUnstable:
		// 2000+ particles out (> 10% of 10000 total) triggers instability
		csvContent = `#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;100.0;50.0
0.100;1000;10000;8000;2000;100;110.0;55.0
0.200;2000;10000;8000;2000;500;130.0;65.0
0.300;3000;10000;8000;2000;1500;180.0;90.0
0.400;4000;10000;8000;2000;2500;250.0;125.0
`
	case syntheticDivergent:
		csvContent = generateDivergentRunCSV(20)
	}

	require.NoError(t, os.WriteFile(filepath.Join(simDir, "Run.csv"), []byte(csvContent), 0644))
	return simDir
}

// generateDivergentRunCSV creates a Run.csv with exponentially growing energy.
// Growth factor 1.5x per step → triggers divergence detection after 5 consecutive steps.
func generateDivergentRunCSV(steps int) string {
	csv := "#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot\n"
	energy := 100.0
	for i := 0; i < steps; i++ {
		csv += fmt.Sprintf("%.3f;%d;10000;8000;2000;0;%.1f;50.0\n",
			float64(i)*0.01, i*100, energy)
		energy *= 1.5
	}
	return csv
}

// generateDivergentWithParticleLossCSV creates Run.csv with both energy divergence
// and significant particle loss (>10% of total), triggering both monitor instability
// and error_recovery divergence detection.
func generateDivergentWithParticleLossCSV(steps int) string {
	csv := "#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot\n"
	energy := 100.0
	for i := 0; i < steps; i++ {
		// PartOut grows linearly, exceeding 10% of 10000 around step 10
		partOut := i * 100
		if partOut > 5000 {
			partOut = 5000
		}
		csv += fmt.Sprintf("%.3f;%d;10000;8000;2000;%d;%.1f;50.0\n",
			float64(i)*0.01, i*100, partOut, energy)
		energy *= 1.5
	}
	return csv
}

// extractJSON parses the JSON portion from a monitor response string.
func extractJSON(t *testing.T, content string, out interface{}) {
	t.Helper()
	// Find JSON object in content (between { and })
	start := -1
	depth := 0
	end := -1
	for i, c := range content {
		if c == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if c == '}' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}
	if start >= 0 && end > start {
		err := json.Unmarshal([]byte(content[start:end]), out)
		require.NoError(t, err, "failed to parse JSON from response")
	} else {
		t.Fatalf("no JSON object found in response: %s", content[:min(200, len(content))])
	}
}
