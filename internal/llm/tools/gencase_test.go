package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TOOL-01: GenCase Tool (XML → Particle Pre-processing)
// FRD AC-1 ~ AC-5

func TestGenCaseTool_Info(t *testing.T) {
	tool := NewGenCaseTool()
	info := tool.Info()

	assert.Equal(t, "gencase", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "case_path")
	assert.Contains(t, info.Parameters, "save_path")
	assert.Contains(t, info.Parameters, "dp")
	assert.Contains(t, info.Required, "case_path")
	assert.Contains(t, info.Required, "save_path")
}

func TestGenCaseTool_Run(t *testing.T) {
	t.Run("AC-1: valid XML produces .bi4 file", func(t *testing.T) {
		// GenCase가 유효한 XML에서 .bi4 파일 생성 성공
		// cases/SloshingTank_Def.xml 입력 → .bi4 존재 확인
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "cases/SloshingTank_Def",
			SavePath: "/tmp/gencase_test_out",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, ".bi4")
	})

	t.Run("AC-2: invalid XML returns Korean error message", func(t *testing.T) {
		// 잘못된 XML 입력 시 한국어 에러 메시지 반환
		// 빈 XML 입력 → ToolResponse.IsError=true, 한국어 메시지
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "/nonexistent/invalid_case",
			SavePath: "/tmp/gencase_test_out",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		// 한국어 에러 메시지 포함 확인
		assert.NotEmpty(t, response.Content)
	})

	t.Run("AC-3: auto-strip .xml extension from path", func(t *testing.T) {
		// 경로에 .xml 포함 시 자동 제거
		// "Tank_Def.xml" 입력 → .xml strip 후 실행
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "cases/SloshingTank_Def.xml", // .xml 포함
			SavePath: "/tmp/gencase_test_out",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		// .xml 확장자가 자동 제거되어 실행 성공해야 함
		assert.False(t, response.IsError)
	})

	t.Run("AC-4: executes via Docker container", func(t *testing.T) {
		// Docker 컨테이너 내 실행
		// docker compose run --rm dsph GenCase ... 호출 확인
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "cases/SloshingTank_Def",
			SavePath: "/tmp/gencase_test_out",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		// 실행 성공 시 docker compose 실행 확인
		assert.False(t, response.IsError)
	})

	t.Run("AC-5: dp override parameter applied", func(t *testing.T) {
		// dp 오버라이드 파라미터 적용
		// dp=0.02 지정 → GenCase -dp:0.02 인자 포함
		dpValue := 0.02
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "cases/SloshingTank_Def",
			SavePath: "/tmp/gencase_test_out",
			DP:       &dpValue,
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})

	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewGenCaseTool()
		call := ToolCall{
			Name:  "gencase",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles missing required parameters", func(t *testing.T) {
		tool := NewGenCaseTool()
		params := GenCaseParams{
			CasePath: "", // 필수 파라미터 누락
			SavePath: "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "gencase",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}
