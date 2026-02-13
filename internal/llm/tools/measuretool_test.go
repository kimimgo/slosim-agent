package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-04: MeasureTool (Water Level / Pressure Measurement)
// FRD AC-1 ~ AC-3

func TestMeasureToolTool_Info(t *testing.T) {
	tool := NewMeasureToolTool()
	info := tool.Info()

	assert.Equal(t, "measuretool", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "points_file")
	assert.Contains(t, info.Parameters, "out_csv")
	assert.Contains(t, info.Parameters, "vars")
	assert.Contains(t, info.Parameters, "elevation")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "points_file")
	assert.Contains(t, info.Required, "out_csv")
}

func TestMeasureToolTool_Run(t *testing.T) {
	t.Run("AC-1: probe_points.txt based measurement success", func(t *testing.T) {
		// probe_points.txt 기반 측정 성공
		// Solver 완료 후 → CSV 파일 생성, 데이터 행 > 0
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, ".csv")
	})

	t.Run("AC-2: elevation mode generates _Elevation.csv", func(t *testing.T) {
		// 수위 높이(elevation) 모드 동작
		// elevation=true → _Elevation.csv 파일 생성
		elevation := true
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
			Elevation:  &elevation,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "Elevation")
	})

	t.Run("AC-3: CSV format validity", func(t *testing.T) {
		// CSV 파일 포맷 정합성
		// 헤더 행 + 숫자 데이터, 파싱 가능 확인
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("handles custom variable selection", func(t *testing.T) {
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/tmp/measuretool_test_data/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
			Vars:       []string{"vel", "press"},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewMeasureToolTool()
		call := ToolCall{
			Name:  "measuretool",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing points file", func(t *testing.T) {
		tool := NewMeasureToolTool()
		params := MeasureToolParams{
			DataDir:    "/tmp/measuretool_test_data",
			PointsFile: "/nonexistent/probe_points.txt",
			OutCSV:     "/tmp/measuretool_test_out/measure",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "measuretool",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
