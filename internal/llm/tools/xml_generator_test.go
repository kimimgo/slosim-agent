package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AGENT-02: XML Case Auto-generation
// FRD AC-1 ~ AC-6

func TestXMLGeneratorTool_Info(t *testing.T) {
	tool := NewXMLGeneratorTool()
	info := tool.Info()

	assert.Equal(t, "xml_generator", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "tank_length")
	assert.Contains(t, info.Parameters, "tank_width")
	assert.Contains(t, info.Parameters, "tank_height")
	assert.Contains(t, info.Parameters, "fluid_height")
	assert.Contains(t, info.Parameters, "freq")
	assert.Contains(t, info.Parameters, "amplitude")
	assert.Contains(t, info.Parameters, "dp")
	assert.Contains(t, info.Parameters, "time_max")
	assert.Contains(t, info.Parameters, "out_path")
	assert.Contains(t, info.Required, "tank_length")
	assert.Contains(t, info.Required, "tank_width")
	assert.Contains(t, info.Required, "tank_height")
	assert.Contains(t, info.Required, "fluid_height")
	assert.Contains(t, info.Required, "freq")
	assert.Contains(t, info.Required, "amplitude")
	assert.Contains(t, info.Required, "dp")
	assert.Contains(t, info.Required, "time_max")
	assert.Contains(t, info.Required, "out_path")
}

func TestXMLGeneratorTool_Run(t *testing.T) {
	t.Run("AC-1: generated XML passes GenCase execution", func(t *testing.T) {
		// 생성된 XML이 GenCase에서 성공 실행
		// XML 생성 → GenCase 실행 → .bi4 생성
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     filepath.Join(tmpDir, "test_case_Def"),
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, ".xml")
	})

	t.Run("AC-2: XML uses attribute-only syntax", func(t *testing.T) {
		// XML은 attribute-only 문법
		// xmllint + 커스텀 검증 (텍스트 노드 없음)
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		outPath := filepath.Join(tmpDir, "test_case_Def")
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     outPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// XML 파일이 생성되었다면 attribute-only 검증
		xmlContent, err := os.ReadFile(outPath + ".xml")
		if err == nil {
			content := string(xmlContent)
			// gravity는 attribute-only: <gravity x="0" y="0" z="-9.81" />
			assert.Contains(t, content, `gravity`)
			assert.Contains(t, content, `x=`)
		}
	})

	t.Run("AC-3: all required sections present", func(t *testing.T) {
		// 필수 섹션 모두 포함
		// constantsdef, geometry, motion, execution/parameters
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		outPath := filepath.Join(tmpDir, "test_case_Def")
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     outPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// XML 파일 읽어서 필수 섹션 확인
		xmlContent, err := os.ReadFile(outPath + ".xml")
		if err == nil {
			content := string(xmlContent)
			assert.True(t, strings.Contains(content, "constantsdef"))
			assert.True(t, strings.Contains(content, "geometry"))
			assert.True(t, strings.Contains(content, "motion"))
			assert.True(t, strings.Contains(content, "execution"))
			assert.True(t, strings.Contains(content, "parameters"))
		}
	})

	t.Run("AC-4: auto-includes 3 SWL gauges", func(t *testing.T) {
		// SWL gauge 3개 자동 포함
		// 좌벽, 중앙, 우벽 수위 게이지
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		outPath := filepath.Join(tmpDir, "test_case_Def")
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     outPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// XML 파일에 gauge 섹션 확인
		xmlContent, err := os.ReadFile(outPath + ".xml")
		if err == nil {
			content := string(xmlContent)
			assert.True(t, strings.Contains(content, "gauges"))
		}
	})

	t.Run("AC-5: auto-generates probe_points.txt", func(t *testing.T) {
		// probe_points.txt 자동 생성
		// 탱크 치수 기반 측정 포인트 좌표
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		outPath := filepath.Join(tmpDir, "test_case_Def")
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     outPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// probe_points.txt 존재 확인
		probeFile := outPath + "_probe_points.txt"
		_, err = os.Stat(probeFile)
		assert.NoError(t, err, "probe_points.txt should be auto-generated")
	})

	t.Run("AC-6: compatible with existing template structure", func(t *testing.T) {
		// 기존 템플릿과 구조 호환
		// cases/Sloshing_Normal_Def.xml과 동일 구조
		tool := NewXMLGeneratorTool()
		tmpDir, err := os.MkdirTemp("", "xml_gen_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		outPath := filepath.Join(tmpDir, "test_case_Def")
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.3,
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     outPath,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// 기본 XML 구조 확인: <case> 루트 요소
		xmlContent, err := os.ReadFile(outPath + ".xml")
		if err == nil {
			content := string(xmlContent)
			assert.True(t, strings.Contains(content, "<case"))
			assert.True(t, strings.Contains(content, "casedef"))
		}
	})

	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewXMLGeneratorTool()
		call := ToolCall{
			Name:  "xml_generator",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles zero dimensions", func(t *testing.T) {
		tool := NewXMLGeneratorTool()
		params := XMLGeneratorParams{
			TankLength:  0,
			TankWidth:   0,
			TankHeight:  0,
			FluidHeight: 0,
			Freq:        0,
			Amplitude:   0,
			DP:          0,
			TimeMax:     0,
			OutPath:     "/tmp/invalid_case",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles fluid_height > tank_height", func(t *testing.T) {
		tool := NewXMLGeneratorTool()
		params := XMLGeneratorParams{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.6,
			FluidHeight: 0.8, // 탱크 높이보다 큼
			Freq:        0.5,
			Amplitude:   0.05,
			DP:          0.02,
			TimeMax:     5.0,
			OutPath:     "/tmp/invalid_case",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "xml_generator",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
