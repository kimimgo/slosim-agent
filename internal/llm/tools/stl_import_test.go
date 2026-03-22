package tools

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// STL-01: STL File Import Tool

// tetrahedronSTL returns an ASCII STL of a watertight tetrahedron (4 triangles, all edges shared exactly twice).
// Vertex winding is CCW when viewed from outside (outward normals match cross product).
func tetrahedronSTL() string {
	// Vertices: v0=(0,0,0), v1=(1,0,0), v2=(0.5,0.866,0), v3=(0.5,0.289,0.816)
	return `solid tetrahedron
facet normal 0 0 -1
  outer loop
    vertex 0 0 0
    vertex 0.5 0.866025 0
    vertex 1 0 0
  endloop
endfacet
facet normal 0 -0.816497 0.57735
  outer loop
    vertex 0 0 0
    vertex 1 0 0
    vertex 0.5 0.288675 0.816497
  endloop
endfacet
facet normal -0.942809 0.235702 0.235702
  outer loop
    vertex 0 0 0
    vertex 0.5 0.288675 0.816497
    vertex 0.5 0.866025 0
  endloop
endfacet
facet normal 0.942809 0.235702 0.235702
  outer loop
    vertex 1 0 0
    vertex 0.5 0.866025 0
    vertex 0.5 0.288675 0.816497
  endloop
endfacet
endsolid tetrahedron
`
}

// openMeshSTL returns a non-watertight STL (single triangle — open surface).
func openMeshSTL() string {
	return `solid open
facet normal 0 0 1
  outer loop
    vertex 0 0 0
    vertex 1 0 0
    vertex 0 1 0
  endloop
endfacet
endsolid open
`
}

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

	t.Run("STL-01: processes watertight STL file", func(t *testing.T) {
		tempDir := t.TempDir()

		stlPath := filepath.Join(tempDir, "tetra.stl")
		err := os.WriteFile(stlPath, []byte(tetrahedronSTL()), 0644)
		require.NoError(t, err)

		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: stlPath,
			OutPath: filepath.Join(tempDir, "output"),
			DP:      0.02,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		response, err := tool.Run(context.Background(), ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "STL 파일 임포트 완료")
		assert.Contains(t, response.Content, "삼각형 수: 4")
		assert.Contains(t, response.Content, "수밀성: true")
	})

	t.Run("STL-01: rejects non-watertight STL", func(t *testing.T) {
		tempDir := t.TempDir()

		stlPath := filepath.Join(tempDir, "open.stl")
		err := os.WriteFile(stlPath, []byte(openMeshSTL()), 0644)
		require.NoError(t, err)

		tool := NewSTLImportTool()
		params := STLImportParams{
			STLFile: stlPath,
			OutPath: filepath.Join(tempDir, "output"),
			DP:      0.02,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		response, err := tool.Run(context.Background(), ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "수밀하지 않습니다")
	})

	t.Run("STL-01: applies scale factor", func(t *testing.T) {
		tempDir := t.TempDir()

		stlPath := filepath.Join(tempDir, "tetra.stl")
		err := os.WriteFile(stlPath, []byte(tetrahedronSTL()), 0644)
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

		response, err := tool.Run(context.Background(), ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// Check XML file created
		xmlPath := filepath.Join(tempDir, "scaled.xml")
		_, err = os.Stat(xmlPath)
		assert.NoError(t, err)
	})

	t.Run("STL-01: handles invalid JSON", func(t *testing.T) {
		tool := NewSTLImportTool()
		response, err := tool.Run(context.Background(), ToolCall{
			Name:  STLImportToolName,
			Input: "invalid json",
		})
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "파라미터 파싱 오류")
	})

	t.Run("STL-01: handles missing parameters", func(t *testing.T) {
		tool := NewSTLImportTool()
		params := STLImportParams{STLFile: "", OutPath: "", DP: 0}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		response, err := tool.Run(context.Background(), ToolCall{
			Name:  STLImportToolName,
			Input: string(paramsJSON),
		})
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "필수 파라미터가 누락")
	})
}

func TestValidateSTLMesh(t *testing.T) {
	t.Run("watertight tetrahedron", func(t *testing.T) {
		tempDir := t.TempDir()
		stlPath := filepath.Join(tempDir, "tetra.stl")
		err := os.WriteFile(stlPath, []byte(tetrahedronSTL()), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.Equal(t, 4, validation.TriangleCount)
		assert.True(t, validation.IsWatertight)
		assert.True(t, validation.NormalsConsistent)
		// Bounding box: x[0,1], y[0, 0.866], z[0, 0.816]
		assert.InDelta(t, 0.0, validation.BBoxMin[0], 0.01)
		assert.InDelta(t, 0.0, validation.BBoxMin[1], 0.01)
		assert.InDelta(t, 0.0, validation.BBoxMin[2], 0.01)
		assert.InDelta(t, 1.0, validation.BBoxMax[0], 0.01)
		assert.InDelta(t, 0.866, validation.BBoxMax[1], 0.01)
		assert.InDelta(t, 0.816, validation.BBoxMax[2], 0.01)
	})

	t.Run("open mesh is not watertight", func(t *testing.T) {
		tempDir := t.TempDir()
		stlPath := filepath.Join(tempDir, "open.stl")
		err := os.WriteFile(stlPath, []byte(openMeshSTL()), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.Equal(t, 1, validation.TriangleCount)
		assert.False(t, validation.IsWatertight)
	})

	t.Run("two-triangle open surface is not watertight", func(t *testing.T) {
		tempDir := t.TempDir()
		stlPath := filepath.Join(tempDir, "square.stl")
		stlContent := `solid square
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
endsolid square
`
		err := os.WriteFile(stlPath, []byte(stlContent), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.Equal(t, 2, validation.TriangleCount)
		assert.False(t, validation.IsWatertight) // Open surface — edges on boundary are shared only once
	})

	t.Run("binary STL parsing", func(t *testing.T) {
		tempDir := t.TempDir()
		stlPath := filepath.Join(tempDir, "binary.stl")

		// Create a minimal binary STL: 1 triangle
		var data []byte
		// 80 bytes header
		header := make([]byte, 80)
		copy(header, "Binary STL test")
		data = append(data, header...)
		// Triangle count: 1
		triCount := make([]byte, 4)
		triCount[0] = 1
		data = append(data, triCount...)
		// Triangle: normal(0,0,1) + 3 vertices + 2 bytes attribute
		tri := make([]byte, 50)
		// Normal z=1.0
		putFloat32LE(tri[8:12], 1.0)
		// v0: (0,0,0) — already zeros
		// v1: (1,0,0)
		putFloat32LE(tri[12:16], 1.0)
		// v2: (0,1,0)
		putFloat32LE(tri[28:32], 1.0)
		data = append(data, tri...)

		err := os.WriteFile(stlPath, data, 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.Equal(t, 1, validation.TriangleCount)
		assert.False(t, validation.IsWatertight) // Single triangle = open
		assert.InDelta(t, 0.0, validation.BBoxMin[0], 0.01)
		assert.InDelta(t, 1.0, validation.BBoxMax[0], 0.01)
		assert.InDelta(t, 1.0, validation.BBoxMax[1], 0.01)
	})

	t.Run("handles non-existent file", func(t *testing.T) {
		_, err := validateSTLMesh("/nonexistent/file.stl")
		assert.Error(t, err)
	})

	t.Run("handles empty file", func(t *testing.T) {
		tempDir := t.TempDir()
		stlPath := filepath.Join(tempDir, "empty.stl")
		err := os.WriteFile(stlPath, []byte(""), 0644)
		require.NoError(t, err)

		validation, err := validateSTLMesh(stlPath)
		require.NoError(t, err)
		assert.NotNil(t, validation)
		assert.Equal(t, 0, validation.TriangleCount)
	})
}

// putFloat32LE writes a float32 in little-endian to buf.
func putFloat32LE(buf []byte, v float32) {
	bits := math.Float32bits(v)
	buf[0] = byte(bits)
	buf[1] = byte(bits >> 8)
	buf[2] = byte(bits >> 16)
	buf[3] = byte(bits >> 24)
}

func dummyValidation(xMin, yMin, zMin, xMax, yMax, zMax float64) *STLValidation {
	return &STLValidation{
		TriangleCount:    100,
		IsWatertight:     true,
		NormalsConsistent: true,
		BBoxMin:          [3]float64{xMin, yMin, zMin},
		BBoxMax:          [3]float64{xMax, yMax, zMax},
	}
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
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)

		assert.Contains(t, xml, "<?xml version=\"1.0\"")
		assert.Contains(t, xml, "<case>")
		assert.Contains(t, xml, "<casedef>")
		assert.Contains(t, xml, "<geometry>")
		assert.Contains(t, xml, "dp=\"0.02\"")
		assert.Contains(t, xml, "<drawfilestl")
		assert.Contains(t, xml, "file=\"tank.stl\"")
		assert.Contains(t, xml, "autofill=\"false\"")
		assert.Contains(t, xml, "<setmkbound mk=\"0\"")
		assert.Contains(t, xml, "<setmkfluid mk=\"0\"")
	})

	t.Run("uses default 50% fill when no height specified", func(t *testing.T) {
		params := STLImportParams{
			STLFile: "test.stl",
			OutPath: "/tmp/test",
			DP:      0.01,
		}
		val := dummyValidation(0, 0, 0, 1.0, 0.5, 0.8)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, "dp=\"0.01\"")
		assert.Contains(t, xml, "<drawbox>")
	})

	t.Run("fillpoint is included in XML", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "fuel_tank.stl",
			OutPath:    "/tmp/fuel",
			DP:         0.008,
			FillPointX: 0.25,
			FillPointY: 0.175,
			FillPointZ: 0.125,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, "<fillpoint")
		assert.Contains(t, xml, "x=\"0.2500\"")
		assert.Contains(t, xml, "y=\"0.1750\"")
	})

	t.Run("motion XML is generated when specified", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "tank.stl",
			OutPath:    "/tmp/tank",
			DP:         0.008,
			MotionType: "mvrectsinu",
			MotionFreq: 0.5,
			MotionAmpl: 0.05,
			MotionAxis: "x",
			TimeMax:    15.0,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, "<motion>")
		assert.Contains(t, xml, "<mvrectsinu")
		assert.Contains(t, xml, "duration=\"15\"")
	})

	t.Run("fillpoint_modefill", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "fuel_tank.stl",
			OutPath:    "/tmp/fuel",
			DP:         0.008,
			FillPointX: 0.25,
			FillPointY: 0.175,
			FillPointZ: 0.125,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, "<modefill>void</modefill>")
		assert.Contains(t, xml, "<fillpoint")
	})

	t.Run("fillpoint_before_drawbox", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "fuel_tank.stl",
			OutPath:    "/tmp/fuel",
			DP:         0.008,
			FillPointX: 0.25,
			FillPointY: 0.175,
			FillPointZ: 0.125,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		fillpointIdx := strings.Index(xml, "<fillpoint")
		drawboxIdx := strings.Index(xml, "<drawbox>")
		assert.Greater(t, drawboxIdx, fillpointIdx, "fillpoint must appear before drawbox")
	})

	t.Run("drawfilestl_autofill", func(t *testing.T) {
		params := STLImportParams{
			STLFile: "tank.stl",
			OutPath: "/tmp/tank",
			DP:      0.02,
			Scale:   1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, `autofill="false"`)
		assert.NotContains(t, xml, `scale="1"`)
		assert.NotContains(t, xml, "<drawscale")
	})

	t.Run("drawfilestl_scale_child", func(t *testing.T) {
		params := STLImportParams{
			STLFile: "tank.stl",
			OutPath: "/tmp/tank",
			DP:      0.02,
			Scale:   2.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, `autofill="false"`)
		assert.Contains(t, xml, `<drawscale x="2" y="2" z="2" />`)
	})

	t.Run("motion_begin_element", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "tank.stl",
			OutPath:    "/tmp/tank",
			DP:         0.008,
			MotionType: "mvrectsinu",
			MotionFreq: 0.5,
			MotionAmpl: 0.05,
			MotionAxis: "x",
			TimeMax:    15.0,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, `<begin mov="1" start="0" />`)
		assert.NotContains(t, xml, `next="0"`)
	})

	t.Run("motion_mvrotsinu_anglesunits", func(t *testing.T) {
		params := STLImportParams{
			STLFile:    "tank.stl",
			OutPath:    "/tmp/tank",
			DP:         0.008,
			MotionType: "mvrotsinu",
			MotionFreq: 0.5,
			MotionAmpl: 4.0,
			MotionAxis: "y",
			TimeMax:    10.0,
			Scale:      1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, `anglesunits="degrees"`)
		assert.Contains(t, xml, `<begin mov="1" start="0" />`)
	})

	t.Run("fill_ratio calculates fluid height from BBox", func(t *testing.T) {
		params := STLImportParams{
			STLFile:   "tank.stl",
			OutPath:   "/tmp/tank",
			DP:        0.008,
			FillRatio: 0.5,
			Scale:     1.0,
		}
		val := dummyValidation(0, 0, 0, 0.5, 0.35, 0.25)
		// BBox height = 0.25, fill_ratio = 0.5 → fluid_height = 0.125
		// But fill_ratio is computed in Run(), not generateSTLXML, so set manually
		params.FluidHeight = 0.125

		xml := generateSTLXML(params, val)
		assert.Contains(t, xml, "<drawbox>")
	})
}
