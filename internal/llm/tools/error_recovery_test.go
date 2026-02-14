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

// NFR-01: Error Recovery

func TestErrorRecoveryTool_Info(t *testing.T) {
	tool := NewErrorRecoveryTool()
	info := tool.Info()

	assert.Equal(t, ErrorRecoveryToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "job_id")
	assert.Contains(t, info.Parameters, "output_dir")
	assert.Contains(t, info.Parameters, "max_retries")
	assert.Contains(t, info.Parameters, "auto_fix")
	assert.Contains(t, info.Required, "output_dir")
}

func TestErrorRecoveryTool_Run(t *testing.T) {
	t.Run("NFR-01: analyzes output directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "error_recovery_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create mock Run.csv (DualSPHysics format: semicolon separator, # header)
		runCSVPath := filepath.Join(tempDir, "Run.csv")
		runContent := `#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;100.0;50.0
0.100;100;10000;8000;2000;0;101.0;51.0
0.200;200;10000;8000;2000;0;102.0;52.0
`
		err = os.WriteFile(runCSVPath, []byte(runContent), 0644)
		require.NoError(t, err)

		tool := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			JobID:           "test_job",
			OutputDir:       tempDir,
			CheckDivergence: true,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "에러 복구 분석 결과")
	})

	t.Run("NFR-01: detects errors in log file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "error_log_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		// Create mock log with errors
		logPath := filepath.Join(tempDir, "DualSPHysics.log")
		logContent := `[INFO] Starting simulation
[ERROR] Particle out of bounds
[WARN] High CFL number detected
[ERROR] NaN value detected in pressure field
`
		err = os.WriteFile(logPath, []byte(logContent), 0644)
		require.NoError(t, err)

		tool := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			OutputDir: tempDir,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "에러 목록")
		assert.Contains(t, response.Content, "제안된 수정 방법")
	})

	t.Run("NFR-01: suggests fixes for common errors", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "error_fix_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "DualSPHysics.log")
		logContent := `[ERROR] GPU memory allocation failed
[ERROR] Particle out of bounds at timestep 1234
`
		err = os.WriteFile(logPath, []byte(logContent), 0644)
		require.NoError(t, err)

		tool := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			OutputDir: tempDir,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "GPU 메모리 부족")
		assert.Contains(t, response.Content, "도메인 크기 증가")
	})

	t.Run("NFR-01: handles missing output directory", func(t *testing.T) {
		tool := NewErrorRecoveryTool()
		params := ErrorRecoveryParams{
			OutputDir: "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ErrorRecoveryToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "output_dir")
	})
}

func TestCheckDivergence(t *testing.T) {
	t.Run("analyzes Run.csv for divergence", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "divergence_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		csvPath := filepath.Join(tempDir, "Run.csv")
		csvContent := `#Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;100.0;50.0
0.100;100;10000;8000;2000;0;105.0;52.0
0.200;200;10000;8000;2000;0;150.0;55.0
`
		err = os.WriteFile(csvPath, []byte(csvContent), 0644)
		require.NoError(t, err)

		divergent, warnings := checkDivergence(csvPath)
		assert.False(t, divergent) // Simplified check
		assert.NotNil(t, warnings)
	})

	t.Run("handles non-existent CSV", func(t *testing.T) {
		divergent, warnings := checkDivergence("/nonexistent/Run.csv")
		assert.False(t, divergent)
		assert.Len(t, warnings, 0)
	})
}

func TestParseLogErrors(t *testing.T) {
	t.Run("extracts errors from log file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "log_parse_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "test.log")
		logContent := `[INFO] Simulation started
[ERROR] Critical failure detected
[WARN] Low particle count
[ERROR] NaN in velocity field
Normal log line
`
		err = os.WriteFile(logPath, []byte(logContent), 0644)
		require.NoError(t, err)

		errors := parseLogErrors(logPath)
		assert.Len(t, errors, 2)
		assert.Contains(t, errors[0], "ERROR")
		assert.Contains(t, errors[1], "NaN")
	})

	t.Run("handles empty log file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "empty_log_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		logPath := filepath.Join(tempDir, "empty.log")
		err = os.WriteFile(logPath, []byte(""), 0644)
		require.NoError(t, err)

		errors := parseLogErrors(logPath)
		assert.Len(t, errors, 0)
	})
}
