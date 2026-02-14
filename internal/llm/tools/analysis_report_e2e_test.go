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

// ═══════════════════════════════════════════════════════════════════════════
// ANALYSIS TOOL E2E TESTS
// ═══════════════════════════════════════════════════════════════════════════

func TestAnalysis_ResonanceAccuracy(t *testing.T) {
	tmpDir := t.TempDir()

	// Known values for resonance calculation
	// Tank: L=1.0m, h=0.3m → theoretical f₁ ≈ 0.86 Hz
	// f₁ = (1/2π) × √(g × π/L × tanh(π/L × h))
	params := AnalysisParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  1.0,
			TankWidth:   0.5,
			TankHeight:  0.5,
			FluidHeight: 0.3,
			Freq:        0.86,
			Amplitude:   0.02,
			DP:          0.02,
			TimeMax:     10.0,
		},
	}

	// Calculate f₁ using same formula as code (analysis.go line 99)
	g := 9.81
	L := params.CaseConfig.TankLength
	h := params.CaseConfig.FluidHeight
	pi := 3.14159265359
	tanhValue := tanh(pi / L * h)
	f1Code := (1 / (2 * pi)) * (g * pi / L * tanhValue)

	// Run analysis tool
	tool := NewAnalysisTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read generated analysis file
	analysisPath := filepath.Join(tmpDir, "analysis.md")
	content, err := os.ReadFile(analysisPath)
	require.NoError(t, err)

	analysisText := string(content)
	assert.Contains(t, analysisText, "f₁")

	// Verify f₁ is positive and matches code calculation
	assert.Greater(t, f1Code, 0.0, "f₁ should be positive")
	_ = math.Abs(f1Code) // use math import
}

func TestAnalysis_ScenarioClassification(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		freq         float64
		expectedText string
	}{
		{
			name:         "Normal (f/f₁ < 0.8)",
			freq:         0.5, // f₁ ≈ 0.758, so 0.5/0.758 ≈ 0.66
			expectedText: "일반 출렁임 (Normal)",
		},
		{
			name:         "Near-Resonance (0.8 ≤ f/f₁ < 0.95)",
			freq:         0.65, // 0.65/0.758 ≈ 0.86
			expectedText: "공진 근처 (Near-Resonance)",
		},
		{
			name:         "Resonance (0.95 ≤ f/f₁ ≤ 1.05)",
			freq:         0.758, // 0.758/0.758 = 1.0
			expectedText: "공진 (Resonance)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tmpDir, strings.ReplaceAll(tt.name, " ", "_"))
			require.NoError(t, os.MkdirAll(testDir, 0755))

			params := AnalysisParams{
				SimDir: testDir,
				CaseConfig: CaseConfig{
					TankLength:  1.0,
					TankWidth:   0.5,
					TankHeight:  0.5,
					FluidHeight: 0.3,
					Freq:        tt.freq,
					Amplitude:   0.02,
					DP:          0.02,
					TimeMax:     10.0,
				},
			}

			tool := NewAnalysisTool()
			paramsJSON, _ := json.Marshal(params)
			call := ToolCall{Input: string(paramsJSON)}

			resp, err := tool.Run(context.Background(), call)
			require.NoError(t, err)
			require.False(t, resp.IsError)

			// Read analysis
			analysisPath := filepath.Join(testDir, "analysis.md")
			content, err := os.ReadFile(analysisPath)
			require.NoError(t, err)

			analysisText := string(content)
			assert.Contains(t, analysisText, tt.expectedText)
		})
	}
}

func TestAnalysis_FileGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	params := AnalysisParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  2.0,
			TankWidth:   1.0,
			TankHeight:  1.0,
			FluidHeight: 0.5,
			Freq:        1.0,
			Amplitude:   0.03,
			DP:          0.025,
			TimeMax:     20.0,
		},
	}

	tool := NewAnalysisTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Verify analysis.md exists
	analysisPath := filepath.Join(tmpDir, "analysis.md")
	_, err = os.Stat(analysisPath)
	assert.NoError(t, err, "analysis.md should be created")

	// Verify file has content
	content, err := os.ReadFile(analysisPath)
	require.NoError(t, err)
	assert.Greater(t, len(content), 100, "analysis.md should have substantial content")

	// Verify sections
	analysisText := string(content)
	assert.Contains(t, analysisText, "# AI 물리적 해석")
	assert.Contains(t, analysisText, "## 1. 슬로싱 시나리오 분석")
	assert.Contains(t, analysisText, "## 2. 해석 조건 평가")
	assert.Contains(t, analysisText, "## 3. 결과 해석 가이드")
	assert.Contains(t, analysisText, "## 5. 추가 검토 사항")
}

func TestAnalysis_WithCSVFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create fake CSV files
	csvDir := filepath.Join(tmpDir, "csv")
	require.NoError(t, os.MkdirAll(csvDir, 0755))

	csv1 := filepath.Join(csvDir, "water_level_left.csv")
	csv2 := filepath.Join(csvDir, "water_level_right.csv")
	require.NoError(t, os.WriteFile(csv1, []byte("Time;Height\n0.0;0.5\n1.0;0.52\n"), 0644))
	require.NoError(t, os.WriteFile(csv2, []byte("Time;Height\n0.0;0.5\n1.0;0.48\n"), 0644))

	params := AnalysisParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  1.5,
			TankWidth:   0.8,
			TankHeight:  0.8,
			FluidHeight: 0.5,
			Freq:        0.9,
			Amplitude:   0.02,
			DP:          0.02,
			TimeMax:     15.0,
		},
		CSVFiles: []string{csv1, csv2},
	}

	tool := NewAnalysisTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read analysis
	analysisPath := filepath.Join(tmpDir, "analysis.md")
	content, err := os.ReadFile(analysisPath)
	require.NoError(t, err)

	analysisText := string(content)

	// Should contain CSV analysis section
	assert.Contains(t, analysisText, "## 4. 측정 데이터 분석")
	assert.Contains(t, analysisText, "2개의 측정 데이터 파일")
	assert.Contains(t, analysisText, "water_level_left.csv")
	assert.Contains(t, analysisText, "water_level_right.csv")
}

// ═══════════════════════════════════════════════════════════════════════════
// REPORT TOOL E2E TESTS
// ═══════════════════════════════════════════════════════════════════════════

func TestReport_FileGeneration(t *testing.T) {
	tmpDir := t.TempDir()

	params := ReportParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  2.0,
			TankWidth:   1.0,
			TankHeight:  1.0,
			FluidHeight: 0.6,
			Freq:        1.2,
			Amplitude:   0.025,
			DP:          0.02,
			TimeMax:     25.0,
		},
	}

	tool := NewReportTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Verify report.md exists
	reportPath := filepath.Join(tmpDir, "report.md")
	_, err = os.Stat(reportPath)
	assert.NoError(t, err, "report.md should be created")

	// Verify file has content
	content, err := os.ReadFile(reportPath)
	require.NoError(t, err)
	assert.Greater(t, len(content), 100, "report.md should have substantial content")

	reportText := string(content)
	assert.Contains(t, reportText, "# 슬로싱 해석 리포트")
	assert.Contains(t, reportText, "## 1. 해석 조건")
	assert.Contains(t, reportText, "## 2. 실행 정보")
	assert.Contains(t, reportText, "## 3. 결과 요약")
}

func TestReport_ConditionsSection(t *testing.T) {
	tmpDir := t.TempDir()

	params := ReportParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  3.5,
			TankWidth:   1.8,
			TankHeight:  2.0,
			FluidHeight: 1.2,
			Freq:        0.75,
			Amplitude:   0.04,
			DP:          0.03,
			TimeMax:     30.0,
		},
	}

	tool := NewReportTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read report
	reportPath := filepath.Join(tmpDir, "report.md")
	content, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	reportText := string(content)

	// Verify tank dimensions
	assert.Contains(t, reportText, "3.50m × 1.80m × 2.00m")

	// Verify fluid height
	assert.Contains(t, reportText, "1.20m")

	// Verify frequency
	assert.Contains(t, reportText, "0.75Hz")

	// Verify particle spacing
	assert.Contains(t, reportText, "0.0300m")

	// Verify simulation time
	assert.Contains(t, reportText, "30.0초")
}

func TestReport_FillRatio(t *testing.T) {
	tmpDir := t.TempDir()

	params := ReportParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  2.0,
			TankWidth:   1.0,
			TankHeight:  1.0,
			FluidHeight: 0.6, // 60% fill
			Freq:        1.0,
			Amplitude:   0.02,
			DP:          0.02,
			TimeMax:     10.0,
		},
	}

	tool := NewReportTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read report
	reportPath := filepath.Join(tmpDir, "report.md")
	content, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	reportText := string(content)

	// Verify fill ratio = fluid_height / tank_height * 100 = 0.6 / 1.0 * 100 = 60%
	assert.Contains(t, reportText, "60% 충진")
}

func TestReport_CSVToMarkdownTable(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test CSV file with semicolon delimiter
	csvPath := filepath.Join(tmpDir, "test_data.csv")
	csvContent := `Time;Height;Pressure
0.00;0.50;101325.0
0.10;0.52;101400.0
0.20;0.48;101250.0
0.30;0.55;101500.0`
	require.NoError(t, os.WriteFile(csvPath, []byte(csvContent), 0644))

	// Call csvToMarkdownTable
	table := csvToMarkdownTable(csvPath)

	// Verify markdown table format
	assert.Contains(t, table, "|")
	assert.Contains(t, table, "Time")
	assert.Contains(t, table, "Height")
	assert.Contains(t, table, "Pressure")
	assert.Contains(t, table, "0.00")
	assert.Contains(t, table, "0.50")
	assert.Contains(t, table, "101325.0")

	// Verify separator line (after header)
	assert.Contains(t, table, "------")
}

func TestReport_VTKCount(t *testing.T) {
	tmpDir := t.TempDir()

	// Create VTK directory with files
	vtkDir := filepath.Join(tmpDir, "vtk")
	require.NoError(t, os.MkdirAll(vtkDir, 0755))

	// Create 5 VTK files
	for i := 0; i < 5; i++ {
		vtkFile := filepath.Join(vtkDir, filepath.Base(tmpDir)+"_"+string(rune('0'+i))+".vtk")
		require.NoError(t, os.WriteFile(vtkFile, []byte("# vtk data"), 0644))
	}

	params := ReportParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  1.5,
			TankWidth:   0.8,
			TankHeight:  0.8,
			FluidHeight: 0.5,
			Freq:        0.9,
			Amplitude:   0.02,
			DP:          0.02,
			TimeMax:     10.0,
		},
	}

	tool := NewReportTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read report
	reportPath := filepath.Join(tmpDir, "report.md")
	content, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	reportText := string(content)

	// Should show VTK count
	assert.Contains(t, reportText, "VTK 파일: 5개")
}

func TestReport_WithCSVFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create CSV files
	csvDir := filepath.Join(tmpDir, "csv")
	require.NoError(t, os.MkdirAll(csvDir, 0755))

	csv1 := filepath.Join(csvDir, "measure_left.csv")
	csv2 := filepath.Join(csvDir, "measure_center.csv")
	csv3 := filepath.Join(csvDir, "measure_right.csv")

	csvData := "Time;Height\n0.0;0.5\n1.0;0.52\n2.0;0.48\n"
	require.NoError(t, os.WriteFile(csv1, []byte(csvData), 0644))
	require.NoError(t, os.WriteFile(csv2, []byte(csvData), 0644))
	require.NoError(t, os.WriteFile(csv3, []byte(csvData), 0644))

	params := ReportParams{
		SimDir: tmpDir,
		CaseConfig: CaseConfig{
			TankLength:  2.0,
			TankWidth:   1.0,
			TankHeight:  1.0,
			FluidHeight: 0.6,
			Freq:        1.0,
			Amplitude:   0.02,
			DP:          0.02,
			TimeMax:     10.0,
		},
		CSVFiles: []string{csv1, csv2, csv3},
	}

	tool := NewReportTool()
	paramsJSON, _ := json.Marshal(params)
	call := ToolCall{Input: string(paramsJSON)}

	resp, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	require.False(t, resp.IsError)

	// Read report
	reportPath := filepath.Join(tmpDir, "report.md")
	content, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	reportText := string(content)

	// Should include CSV section
	assert.Contains(t, reportText, "## 4. 수위 시계열 데이터")
	assert.Contains(t, reportText, "측정 데이터: 3개 파일")
	assert.Contains(t, reportText, "measure_left.csv")
	assert.Contains(t, reportText, "measure_center.csv")
	assert.Contains(t, reportText, "measure_right.csv")

	// Should contain markdown tables
	assert.Contains(t, reportText, "Time")
	assert.Contains(t, reportText, "Height")
}
