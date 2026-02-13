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

// STL-01: STL File Import Tool

func TestSTLImportTool_Info(t *testing.T) {
	tool := NewSTLImportTool()
	info := tool.Info()

	assert.Equal(t, STLImportToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "stl_file")
	assert.Contains(t, info.Parameters, "out_path")
	assert.Contains(t, info.Parameters, "dp")
	assert.Contains(t, info.Parameters, "scale")
	assert.Contains(t, info.Required, "stl_file")
	assert.Contains(t, info.Required, "out_path")
	assert.Contains(t, info.Required, "dp")
}

func TestSTLImportTool_Run(t *testing.T) {
	t.Run("STL-01: validates file existence", func(t *testing.T) {
		// 존재하지 않는 파일 → 에러 반환
		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: "/nonexistent/file.stl",
			OutPath: "/tmp/output",
			DP:      0.02,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "STL 파일을 찾을 수 없습니다")
	})

	t.Run("STL-01: processes valid STL file", func(t *testing.T) {
		// Create a temporary STL file
		tempDir, err := os.MkdirTemp("", "stl_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		stlPath := filepath.Join(tempDir, "test.stl")
		// Minimal STL content (ASCII format)
		stlContent := `solid test
facet normal 0 0 1
  outer loop
    vertex 0 0 0
    vertex 1 0 0
    vertex 0 1 0
  endloop
endfacet
endsolid test
`
		err = os.WriteFile(stlPath, []byte(stlContent), 0644)
		require.NoError(t, err)

		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: stlPath,
			OutPath: filepath.Join(tempDir, "output"),
			DP:      0.02,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "STL 파일 임포트 완료")
		assert.Contains(t, response.Content, "삼각형 수")
	})

	t.Run("STL-01: applies scale factor", func(t *testing.T) {
		// Create a temporary STL file
		tempDir, err := os.MkdirTemp("", "stl_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		stlPath := filepath.Join(tempDir, "test.stl")
		stlContent := `solid test
facet normal 0 0 1
  outer loop
    vertex 0 0 0
    vertex 1 0 0
    vertex 0 1 0
  endloop
endfacet
endsolid test
`
		err = os.WriteFile(stlPath, []byte(stlContent), 0644)
		require.NoError(t, err)

		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: stlPath,
			OutPath: filepath.Join(tempDir, "scaled"),
			DP:      0.02,
			Scale:   2.0,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// Check that XML file was created
		xmlPath := filepath.Join(tempDir, "scaled.xml")
		_, err = os.Stat(xmlPath)
		assert.NoError(t, err)
	})

	t.Run("STL-01: handles invalid JSON", func(t *testing.T) {
		tool := NewSTLImportTool()
		call := ToolCall{
			Name:  STLImportToolName,
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "파라미터 파싱 오류")
	})

	t.Run("STL-01: handles missing parameters", func(t *testing.T) {
		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: "", // Missing required parameter
			OutPath: "",
			DP:      0,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "필수 파라미터가 누락")
	})
}

func TestValidateSTLMesh(t *testing.T) {
	t.Run("validates ASCII STL format", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "stl_validate_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		stlPath := filepath.Join(tempDir, "test.stl")
		stlContent := `solid cube
facet normal 0 0 1
  outer loop
    vertex 0 0 0
    vertex 1 0 0
    vertex 0 1 0
  endloop
endfacet
facet normal 0 0 1
  outer loop
    vertex 1 0 0
    vertex 1 1 0
    vertex 0 1 0
  endloop
endfacet
endsolid cube
`
		err = os.WriteFile(stlPath, []byte(stlContent), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.NotNil(t, validation)
		assert.Equal(t, 2, validation.TriangleCount)
		assert.True(t, validation.IsWatertight)
	})

	t.Run("handles non-existent file", func(t *testing.T) {
		_, err := validateSTLMesh("/nonexistent/file.stl")
		assert.Error(t, err)
	})

	t.Run("handles empty file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "stl_empty_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		stlPath := filepath.Join(tempDir, "empty.stl")
		err = os.WriteFile(stlPath, []byte(""), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.NotNil(t, validation)
		assert.Equal(t, 0, validation.TriangleCount)
	})
}

func TestGenerateSTLXML(t *testing.T) {
	t.Run("generates valid XML for STL import", func(t *testing.T) {
		params := STLImportParams{
			STLFile:     "tank.stl",
			OutPath:     "/tmp/tank",
			DP:          0.02,
			FluidHeight: 1.5,
			Scale:       1.0,
		}

		xml := generateSTLXML(params)

		// Check essential XML components
		assert.Contains(t, xml, "<?xml version=\"1.0\"")
		assert.Contains(t, xml, "<case>")
		assert.Contains(t, xml, "<casedef>")
		assert.Contains(t, xml, "<geometry>")
		assert.Contains(t, xml, "dp=\"0.02\"")
		assert.Contains(t, xml, "<drawfilestl")
		assert.Contains(t, xml, "file=\"tank.stl\"")
		assert.Contains(t, xml, "scale=\"1\"")
		assert.Contains(t, xml, "<setmkbound mk=\"0\"")
		assert.Contains(t, xml, "<setmkfluid mk=\"0\"")
	})

	t.Run("uses default fluid height when not specified", func(t *testing.T) {
		params := STLImportParams{
			STLFile: "test.stl",
			OutPath: "/tmp/test",
			DP:      0.01,
			// FluidHeight not specified
		}

		xml := generateSTLXML(params)

		// Should use default fluid height (1.0)
		assert.Contains(t, xml, "dp=\"0.01\"")
		assert.Contains(t, xml, "<drawbox>")
	})
}
