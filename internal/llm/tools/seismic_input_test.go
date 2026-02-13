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

// EXC-01: Seismic Wave Input Tool

func TestSeismicInputTool_Info(t *testing.T) {
	tool := NewSeismicInputTool()
	info := tool.Info()

	assert.Equal(t, SeismicInputToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "file_path")
	assert.Contains(t, info.Parameters, "time_column")
	assert.Contains(t, info.Parameters, "accel_column")
	assert.Contains(t, info.Parameters, "scale_factor")
	assert.Contains(t, info.Parameters, "validate_only")
	assert.Contains(t, info.Required, "file_path")
}

func TestSeismicInputTool_Run(t *testing.T) {
	t.Run("EXC-01: parses space-delimited seismic file", func(t *testing.T) {
		// Create temporary seismic data file
		tempDir, err := os.MkdirTemp("", "seismic_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "seismic.txt")
		dataContent := `# Seismic wave data
0.00 0.000
0.01 0.500
0.02 1.000
0.03 0.500
0.04 0.000
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath:     dataPath,
			ValidateOnly: true,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "지진파 데이터 검증 완료")
		assert.Contains(t, response.Content, "샘플 수: 5")
		assert.Contains(t, response.Content, "지속 시간")
	})

	t.Run("EXC-01: parses comma-delimited file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "seismic.csv")
		dataContent := `time,acceleration
0.0,0.0
0.1,0.5
0.2,1.0
0.3,0.5
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath:     dataPath,
			SkipRows:     1, // Skip header
			ValidateOnly: true,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "샘플 수: 4")
	})

	t.Run("EXC-01: applies scale factor", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "seismic.txt")
		dataContent := `0.0 1.0
0.1 2.0
0.2 3.0
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath:     dataPath,
			ScaleFactor:  0.5,
			ValidateOnly: true,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "스케일 계수: 0.50")
		// Scaled acceleration range: 1.0*0.5=0.5, 3.0*0.5=1.5
		assert.Contains(t, response.Content, "가속도 범위")
	})

	t.Run("EXC-01: converts to DualSPHysics format", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "seismic.txt")
		dataContent := `0.0 0.0
0.1 1.0
0.2 2.0
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath:     dataPath,
			ValidateOnly: false, // Perform conversion
			OutputFormat: "dsph",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "변환된 파일 경로")

		// Check output file exists
		outputPath := filepath.Join(tempDir, "seismic_dsph.dat")
		_, err = os.Stat(outputPath)
		assert.NoError(t, err)
	})

	t.Run("EXC-01: handles non-existent file", func(t *testing.T) {
		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath: "/nonexistent/file.txt",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "파일을 찾을 수 없습니다")
	})

	t.Run("EXC-01: handles invalid JSON", func(t *testing.T) {
		tool := NewSeismicInputTool()
		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "파라미터 파싱 오류")
	})

	t.Run("EXC-01: handles empty file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "empty.txt")
		err = os.WriteFile(dataPath, []byte(""), 0644)
		require.NoError(t, err)

		tool := NewSeismicInputTool()
		params := SeismicInputParams{
			FilePath: dataPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  SeismicInputToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "데이터가 비어 있습니다")
	})
}

func TestParseSeismicFile(t *testing.T) {
	t.Run("parses tab-delimited data", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_parse_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "data.txt")
		dataContent := "0.0\t1.5\n0.1\t2.5\n0.2\t3.5\n"
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		params := SeismicInputParams{
			TimeColumn:  0,
			AccelColumn: 1,
		}

		timeData, accelData, err := parseSeismicFile(dataPath, params)
		require.NoError(t, err)
		assert.Len(t, timeData, 3)
		assert.Len(t, accelData, 3)
		assert.Equal(t, 0.0, timeData[0])
		assert.Equal(t, 1.5, accelData[0])
		assert.Equal(t, 0.2, timeData[2])
		assert.Equal(t, 3.5, accelData[2])
	})

	t.Run("skips comment lines", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_parse_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "data.txt")
		dataContent := `# Comment line
0.0 1.0
# Another comment
0.1 2.0
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		params := SeismicInputParams{}

		timeData, accelData, err := parseSeismicFile(dataPath, params)
		require.NoError(t, err)
		assert.Len(t, timeData, 2)
		assert.Len(t, accelData, 2)
	})

	t.Run("handles header skip", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "seismic_parse_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		dataPath := filepath.Join(tempDir, "data.txt")
		dataContent := `Time Acceleration
seconds m/s^2
0.0 1.0
0.1 2.0
`
		err = os.WriteFile(dataPath, []byte(dataContent), 0644)
		require.NoError(t, err)

		params := SeismicInputParams{
			SkipRows: 2,
		}

		timeData, accelData, err := parseSeismicFile(dataPath, params)
		require.NoError(t, err)
		assert.Len(t, timeData, 2)
		assert.Len(t, accelData, 2)
	})
}

func TestFindMinMax(t *testing.T) {
	t.Run("finds correct min and max", func(t *testing.T) {
		data := []float64{1.0, 5.0, -3.0, 7.0, 2.0}
		min, max := findMinMax(data)
		assert.Equal(t, -3.0, min)
		assert.Equal(t, 7.0, max)
	})

	t.Run("handles empty array", func(t *testing.T) {
		data := []float64{}
		min, max := findMinMax(data)
		assert.Equal(t, 0.0, min)
		assert.Equal(t, 0.0, max)
	})

	t.Run("handles single element", func(t *testing.T) {
		data := []float64{42.0}
		min, max := findMinMax(data)
		assert.Equal(t, 42.0, min)
		assert.Equal(t, 42.0, max)
	})
}
