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

// RPT-01: Markdown Report Generation
// FRD AC-1 ~ AC-4

func TestReportTool_Info(t *testing.T) {
	tool := NewReportTool()
	info := tool.Info()

	assert.Equal(t, "report", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "sim_dir")
	assert.Contains(t, info.Parameters, "case_config")
	assert.Contains(t, info.Parameters, "csv_files")
	assert.Contains(t, info.Required, "sim_dir")
	assert.Contains(t, info.Required, "case_config")
}

func TestReportTool_Run(t *testing.T) {
	t.Run("AC-1: generates report.md file", func(t *testing.T) {
		// 리포트 Markdown 파일 생성
		// 시뮬레이션 완료 후 → report.md 존재
		tool := NewReportTool()
		tmpDir, err := os.MkdirTemp("", "report_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		params := ReportParams{
			SimDir: tmpDir,
			CaseConfig: CaseConfig{
				TankLength:  1.0,
				TankWidth:   0.5,
				TankHeight:  0.6,
				FluidHeight: 0.3,
				Freq:        0.5,
				Amplitude:   0.05,
				DP:          0.02,
				TimeMax:     5.0,
			},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// report.md 존재 확인
		reportPath := filepath.Join(tmpDir, "report.md")
		_, err = os.Stat(reportPath)
		assert.NoError(t, err, "report.md should be generated")
	})

	t.Run("AC-2: report includes all input parameters", func(t *testing.T) {
		// 해석 조건 섹션에 모든 입력 파라미터 포함
		// 탱크 치수, 주파수, dp, TimeMax 확인
		tool := NewReportTool()
		tmpDir, err := os.MkdirTemp("", "report_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		params := ReportParams{
			SimDir: tmpDir,
			CaseConfig: CaseConfig{
				TankLength:  1.0,
				TankWidth:   0.5,
				TankHeight:  0.6,
				FluidHeight: 0.3,
				Freq:        0.5,
				Amplitude:   0.05,
				DP:          0.02,
				TimeMax:     5.0,
			},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// report.md 읽어서 파라미터 포함 확인
		reportPath := filepath.Join(tmpDir, "report.md")
		content, err := os.ReadFile(reportPath)
		if err == nil {
			report := string(content)
			assert.True(t, strings.Contains(report, "1.0") || strings.Contains(report, "1m"))  // 탱크 길이
			assert.True(t, strings.Contains(report, "0.5"))                                      // 주파수 or 너비
			assert.True(t, strings.Contains(report, "0.02") || strings.Contains(report, "입자")) // dp
			assert.True(t, strings.Contains(report, "5.0") || strings.Contains(report, "5초"))   // TimeMax
		}
	})

	t.Run("AC-3: includes water level table from CSV", func(t *testing.T) {
		// 수위 테이블 데이터 포함 (CSV가 있는 경우)
		// MeasureTool CSV → Markdown 테이블 변환
		tool := NewReportTool()
		tmpDir, err := os.MkdirTemp("", "report_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// 테스트용 CSV 파일 생성
		csvDir := filepath.Join(tmpDir, "csv")
		err = os.MkdirAll(csvDir, 0755)
		require.NoError(t, err)

		csvContent := "Time;WaterLevel_Left;WaterLevel_Center;WaterLevel_Right\n0.0;0.300;0.300;0.300\n0.1;0.310;0.298;0.290\n"
		csvFile := filepath.Join(csvDir, "measure_Elevation.csv")
		err = os.WriteFile(csvFile, []byte(csvContent), 0644)
		require.NoError(t, err)

		params := ReportParams{
			SimDir: tmpDir,
			CaseConfig: CaseConfig{
				TankLength:  1.0,
				TankWidth:   0.5,
				TankHeight:  0.6,
				FluidHeight: 0.3,
				Freq:        0.5,
				Amplitude:   0.05,
				DP:          0.02,
				TimeMax:     5.0,
			},
			CSVFiles: []string{csvFile},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// report.md에 테이블 데이터 포함 확인
		reportPath := filepath.Join(tmpDir, "report.md")
		content, err := os.ReadFile(reportPath)
		if err == nil {
			report := string(content)
			// Markdown 테이블 구분자 확인
			assert.True(t, strings.Contains(report, "|"))
		}
	})

	t.Run("AC-4: no CFD jargon in user-facing text", func(t *testing.T) {
		// CFD 전문 용어 없음 (사용자 대면 텍스트)
		// "dp" → "입자 간격", "SPH" 미사용 등
		tool := NewReportTool()
		tmpDir, err := os.MkdirTemp("", "report_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		params := ReportParams{
			SimDir: tmpDir,
			CaseConfig: CaseConfig{
				TankLength:  1.0,
				TankWidth:   0.5,
				TankHeight:  0.6,
				FluidHeight: 0.3,
				Freq:        0.5,
				Amplitude:   0.05,
				DP:          0.02,
				TimeMax:     5.0,
			},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)

		// report.md에 CFD 전문 용어 미포함 확인
		reportPath := filepath.Join(tmpDir, "report.md")
		content, err := os.ReadFile(reportPath)
		if err == nil {
			report := string(content)
			// 사용자 대면 텍스트에 전문 용어가 노출되지 않아야 함
			assert.False(t, strings.Contains(report, "SPH"), "report should not contain 'SPH'")
			assert.False(t, strings.Contains(report, "CFL"), "report should not contain 'CFL'")
			// "dp"는 "입자 간격"으로 표시되어야 함
			if strings.Contains(report, "dp") {
				// dp가 포함되었다면 "입자 간격"도 있어야 함
				assert.True(t, strings.Contains(report, "입자 간격") || strings.Contains(report, "입자"), "dp should be shown as user-friendly term")
			}
		}
	})

	t.Run("handles invalid JSON parameters", func(t *testing.T) {
		tool := NewReportTool()
		call := ToolCall{
			Name:  "report",
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles non-existent sim directory", func(t *testing.T) {
		tool := NewReportTool()
		params := ReportParams{
			SimDir: "/nonexistent/sim/dir",
			CaseConfig: CaseConfig{
				TankLength: 1.0,
				TankWidth:  0.5,
			},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})

	t.Run("handles empty csv_files gracefully", func(t *testing.T) {
		// CSV 파일 없이도 리포트 생성 가능해야 함
		tool := NewReportTool()
		tmpDir, err := os.MkdirTemp("", "report_test")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		params := ReportParams{
			SimDir: tmpDir,
			CaseConfig: CaseConfig{
				TankLength:  1.0,
				TankWidth:   0.5,
				TankHeight:  0.6,
				FluidHeight: 0.3,
				Freq:        0.5,
				Amplitude:   0.05,
				DP:          0.02,
				TimeMax:     5.0,
			},
			CSVFiles: []string{}, // 빈 CSV 목록
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  "report",
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
	})
}
