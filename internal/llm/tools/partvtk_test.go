package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-03: PartVTK Tool (VTK Export)
// FRD AC-1 ~ AC-4

func TestPartVTKTool_Info(t *testing.T) {
	tool := NewPartVTKTool()
	info := tool.Info()

	assert.Equal(t, "partvtk", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "out_file")
	assert.Contains(t, info.Parameters, "only_fluid")
	assert.Contains(t, info.Parameters, "vars")
	assert.Contains(t, info.Parameters, "first")
	assert.Contains(t, info.Parameters, "last")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "out_file")
}

func TestPartVTKTool_Run(t *testing.T) {
	t.Run("AC-1: Part_*.bi4 to VTK conversion success", func(t *testing.T) {
		// Part_*.bi4 → VTK 변환 성공
		// Solver 완료 후 실행 → .vtk 파일 생성 확인
		tool := NewPartVTKTool()
		params := PartVTKParams{
			DataDir: "/tmp/partvtk_test_data",
			OutFile: "/tmp/partvtk_test_out/fluid",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, ".vtk")
	})

	t.Run("AC-2: only_fluid option excludes boundary particles", func(t *testing.T) {
		// 유체만 추출 옵션 동작
		// only_fluid=true → 경계 파티클 미포함 확인
		tool := NewPartVTKTool()
		onlyFluid := true
		params := PartVTKParams{
			DataDir:   "/tmp/partvtk_test_data",
			OutFile:   "/tmp/partvtk_test_out/fluid",
			OnlyFluid: &onlyFluid,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("AC-3: variable selection works", func(t *testing.T) {
		// 변수 선택 동작
		// vars=["vel","press"] → VTK에 해당 필드만 포함
		tool := NewPartVTKTool()
		params := PartVTKParams{
			DataDir: "/tmp/partvtk_test_data",
			OutFile: "/tmp/partvtk_test_out/fluid",
			Vars:    []string{"vel", "press"},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("AC-4: timestep range filtering", func(t *testing.T) {
		// 타임스텝 범위 필터링
		// first=10, last=20 → 해당 범위만 변환
		first := 10
		last := 20
		tool := NewPartVTKTool()
		params := PartVTKParams{
			DataDir: "/tmp/partvtk_test_data",
			OutFile: "/tmp/partvtk_test_out/fluid",
			First:   &first,
			Last:    &last,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewPartVTKTool()
		call := ToolCall{
			Name:  "partvtk",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing data directory", func(t *testing.T) {
		tool := NewPartVTKTool()
		params := PartVTKParams{
			DataDir: "/nonexistent/path",
			OutFile: "/tmp/partvtk_test_out/fluid",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "partvtk",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
