package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CaseConfig holds the simulation parameters for report generation.
type CaseConfig struct {
	TankLength  float64 `json:"tank_length"`
	TankWidth   float64 `json:"tank_width"`
	TankHeight  float64 `json:"tank_height"`
	FluidHeight float64 `json:"fluid_height"`
	Freq        float64 `json:"freq"`
	Amplitude   float64 `json:"amplitude"`
	DP          float64 `json:"dp"`
	TimeMax     float64 `json:"time_max"`
}

type ReportParams struct {
	SimDir     string     `json:"sim_dir"`
	CaseConfig CaseConfig `json:"case_config"`
	CSVFiles   []string   `json:"csv_files,omitempty"`
	ImageFiles []string   `json:"image_files,omitempty"` // RPT-02: Include images
}

type reportTool struct{}

const (
	ReportToolName    = "report"
	reportDescription = `슬로싱 해석 리포트 생성 도구 — 시뮬레이션 결과를 Markdown 리포트로 정리합니다.

입력:
- sim_dir: 시뮬레이션 결과 디렉토리
- case_config: 해석 조건 (탱크 치수, 주파수, 입자 간격 등)
- csv_files: MeasureTool CSV 파일 경로 (선택)

출력:
- {sim_dir}/report.md — Markdown 형식 해석 리포트`
)

func NewReportTool() BaseTool {
	return &reportTool{}
}

func (r *reportTool) Info() ToolInfo {
	return ToolInfo{
		Name:        ReportToolName,
		Description: reportDescription,
		Parameters: map[string]any{
			"sim_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 결과 디렉토리 (예: simulations/{case_name})",
			},
			"case_config": map[string]any{
				"type":        "object",
				"description": "해석 조건 (탱크 치수, 주파수, 입자 간격 등)",
			},
			"csv_files": map[string]any{
				"type":        "array",
				"description": "MeasureTool CSV 파일 경로",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		Required: []string{"sim_dir", "case_config"},
	}
}

func (r *reportTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params ReportParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.SimDir == "" {
		return NewTextErrorResponse("시뮬레이션 디렉토리(sim_dir)를 지정해주세요"), nil
	}

	// Validate sim_dir exists
	if _, err := os.Stat(params.SimDir); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("시뮬레이션 디렉토리를 찾을 수 없습니다: %s", params.SimDir)), nil
	}

	// Generate report content
	report := generateReport(params)

	// Write report.md
	reportPath := filepath.Join(params.SimDir, "report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("리포트 파일 생성 실패: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("해석 리포트가 생성되었습니다: %s", reportPath)), nil
}

func generateReport(params ReportParams) string {
	cc := params.CaseConfig
	fillRatio := 0.0
	if cc.TankHeight > 0 {
		fillRatio = cc.FluidHeight / cc.TankHeight * 100
	}

	var sb strings.Builder

	// Header
	sb.WriteString("# 슬로싱 해석 리포트\n\n")

	// Section 1: Analysis conditions
	sb.WriteString("## 1. 해석 조건\n\n")
	sb.WriteString(fmt.Sprintf("- 탱크: %.2fm × %.2fm × %.2fm (직사각형)\n", cc.TankLength, cc.TankWidth, cc.TankHeight))
	sb.WriteString(fmt.Sprintf("- 유체: 물 (1000 kg/m³), 높이 %.2fm (%.0f%% 충진)\n", cc.FluidHeight, fillRatio))
	sb.WriteString(fmt.Sprintf("- 가진: 정현파 %.2fHz, 진폭 %.3fm\n", cc.Freq, cc.Amplitude))
	sb.WriteString(fmt.Sprintf("- 입자 간격: %.4fm\n", cc.DP))
	sb.WriteString(fmt.Sprintf("- 시뮬레이션 시간: %.1f초\n", cc.TimeMax))
	sb.WriteString("\n")

	// Section 2: Execution info
	sb.WriteString("## 2. 실행 정보\n\n")
	sb.WriteString("- GPU: NVIDIA RTX 4090\n")
	sb.WriteString("- 솔버: 입자 기반 유체 시뮬레이션 v5.4\n")
	sb.WriteString("\n")

	// Section 3: Results summary
	sb.WriteString("## 3. 결과 요약\n\n")

	// Count VTK files if vtk/ directory exists
	vtkDir := filepath.Join(params.SimDir, "vtk")
	vtkCount := countFiles(vtkDir, ".vtk")
	if vtkCount > 0 {
		sb.WriteString(fmt.Sprintf("- VTK 파일: %d개 생성 (%s)\n", vtkCount, vtkDir))
	} else {
		sb.WriteString("- VTK 파일: 미생성\n")
	}

	if len(params.CSVFiles) > 0 {
		sb.WriteString(fmt.Sprintf("- 측정 데이터: %d개 파일\n", len(params.CSVFiles)))
		for _, f := range params.CSVFiles {
			sb.WriteString(fmt.Sprintf("  - %s\n", filepath.Base(f)))
		}
	} else {
		sb.WriteString("- 측정 데이터: 없음\n")
	}

	// RPT-02: Include images
	if len(params.ImageFiles) > 0 {
		sb.WriteString(fmt.Sprintf("- 렌더링 이미지: %d개 파일\n", len(params.ImageFiles)))
	}
	sb.WriteString("\n")

	// RPT-02: Embed images in report
	if len(params.ImageFiles) > 0 {
		sb.WriteString("### 시각화 결과\n\n")
		for _, imgPath := range params.ImageFiles {
			relPath, _ := filepath.Rel(params.SimDir, imgPath)
			sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", filepath.Base(imgPath), relPath))
		}
	}

	// Section 4: Water level time series (if CSV available)
	if len(params.CSVFiles) > 0 {
		sb.WriteString("## 4. 수위 시계열 데이터\n\n")
		for _, csvFile := range params.CSVFiles {
			table := csvToMarkdownTable(csvFile)
			if table != "" {
				sb.WriteString(fmt.Sprintf("### %s\n\n", filepath.Base(csvFile)))
				sb.WriteString(table)
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

func countFiles(dir string, ext string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ext) {
			count++
		}
	}
	return count
}

func csvToMarkdownTable(csvPath string) string {
	file, err := os.Open(csvPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	lineNum := 0
	maxRows := 20 // Limit output rows

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Detect delimiter (semicolon or comma)
		delimiter := ";"
		if !strings.Contains(line, ";") {
			delimiter = ","
		}

		fields := strings.Split(line, delimiter)
		row := "| " + strings.Join(fields, " | ") + " |"
		sb.WriteString(row + "\n")

		// Add separator after header
		if lineNum == 0 {
			sep := "|"
			for range fields {
				sep += "---------|"
			}
			sb.WriteString(sep + "\n")
		}

		lineNum++
		if lineNum > maxRows {
			sb.WriteString("| ... | (이하 생략) |\n")
			break
		}
	}

	return sb.String()
}
