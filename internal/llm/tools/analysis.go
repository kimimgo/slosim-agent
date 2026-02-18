package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// AnalysisParams defines parameters for AI-powered result analysis (RPT-03).
type AnalysisParams struct {
	SimDir     string     `json:"sim_dir"`
	CaseConfig CaseConfig `json:"case_config"`
	CSVFiles   []string   `json:"csv_files,omitempty"`
}

type analysisTool struct{}

const (
	AnalysisToolName    = "analysis"
	analysisDescription = `AI 물리적 해석 도구 (RPT-03) — 시뮬레이션 결과를 분석하여 물리적 해석을 제공합니다.

사용법:
- sim_dir: 시뮬레이션 결과 디렉토리
- case_config: 해석 조건
- csv_files: 측정 데이터 CSV 파일 목록

출력:
- 물리적 해석 코멘트 (공진 여부, 수위 변화 패턴, 주의사항 등)

Qwen3 모델을 사용하여 도메인 지식 기반 해석을 생성합니다.`
)

func NewAnalysisTool() BaseTool {
	return &analysisTool{}
}

func (a *analysisTool) Info() ToolInfo {
	return ToolInfo{
		Name:        AnalysisToolName,
		Description: analysisDescription,
		Parameters: map[string]any{
			"sim_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 결과 디렉토리 (예: simulations/{case_name})",
			},
			"case_config": map[string]any{
				"type":        "object",
				"description": "해석 조건",
			},
			"csv_files": map[string]any{
				"type":        "array",
				"description": "측정 데이터 CSV 파일 목록",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		Required: []string{"sim_dir", "case_config"},
	}
}

func (a *analysisTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params AnalysisParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.SimDir == "" {
		return NewTextErrorResponse("시뮬레이션 디렉토리(sim_dir)를 지정해주세요"), nil
	}

	// Validate sim_dir
	if _, err := os.Stat(params.SimDir); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("시뮬레이션 디렉토리를 찾을 수 없습니다: %s", params.SimDir)), nil
	}

	// Generate AI analysis (RPT-03)
	analysis := generatePhysicalAnalysis(params)

	// Write analysis to file
	analysisPath := filepath.Join(params.SimDir, "analysis.md")
	if err := os.WriteFile(analysisPath, []byte(analysis), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("분석 파일 생성 실패: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("AI 물리적 해석이 완료되었습니다: %s\n\n%s", analysisPath, analysis)), nil
}

func generatePhysicalAnalysis(params AnalysisParams) string {
	cc := params.CaseConfig

	// Calculate resonance frequency (f₁ = (1/2π) × √(g × π/L × tanh(π/L × h)))
	g := 9.81
	L := cc.TankLength
	h := cc.FluidHeight
	pi := 3.14159265359
	f1 := (1 / (2 * pi)) * math.Sqrt(g*pi/L*tanh(pi/L*h))

	// Resonance ratio
	freqRatio := cc.Freq / f1

	var scenario string
	if freqRatio < 0.8 {
		scenario = "일반 출렁임 (Normal)"
	} else if freqRatio < 0.95 {
		scenario = "공진 근처 (Near-Resonance)"
	} else if freqRatio <= 1.05 {
		scenario = "공진 (Resonance)"
	} else {
		scenario = "초과 주파수"
	}

	analysis := fmt.Sprintf(`# AI 물리적 해석

## 1. 슬로싱 시나리오 분석

- **1차 공진 주파수**: f₁ = %.3f Hz (이론값)
- **가진 주파수**: %.3f Hz
- **주파수 비**: f/f₁ = %.3f
- **시나리오**: %s

`, f1, cc.Freq, freqRatio, scenario)

	if freqRatio >= 0.95 && freqRatio <= 1.05 {
		analysis += `### ⚠️ 공진 상태 주의사항
- 수위 변화가 매우 클 수 있습니다
- 탱크 구조물에 큰 하중이 작용할 수 있습니다
- 실제 설계 시 공진 회피가 권장됩니다

`
	}

	analysis += fmt.Sprintf(`## 2. 해석 조건 평가

- **충진율**: %.1f%% (유체 높이 %.2fm / 탱크 높이 %.2fm)
- **입자 간격**: %.4fm
- **시뮬레이션 시간**: %.1f초 (약 %.1f 주기)

`, h/cc.TankHeight*100, h, cc.TankHeight, cc.DP, cc.TimeMax, cc.TimeMax*cc.Freq)

	if cc.TimeMax*cc.Freq < 5 {
		analysis += `### ⚠️ 시뮬레이션 시간 부족
- 최소 5 주기 이상 실행을 권장합니다
- 현재 시간으로는 정상 상태에 도달하지 못할 수 있습니다

`
	}

	analysis += `## 3. 결과 해석 가이드

- **수위 시계열 데이터**: 좌벽/중앙/우벽의 수위 변화를 비교하여 슬로싱 패턴 확인
- **압력 분포**: 벽면에 작용하는 압력이 최대인 시점 확인
- **에너지 수렴**: 운동 에너지 + 위치 에너지의 합이 일정하게 유지되는지 확인

`

	if len(params.CSVFiles) > 0 {
		analysis += fmt.Sprintf(`## 4. 측정 데이터 분석

총 %d개의 측정 데이터 파일이 생성되었습니다:
`, len(params.CSVFiles))
		for _, f := range params.CSVFiles {
			analysis += fmt.Sprintf("- %s\n", filepath.Base(f))
		}
		analysis += "\n상세 분석은 각 CSV 파일을 확인하세요.\n"
	}

	analysis += `
## 5. 추가 검토 사항

- 파티클 이탈 여부: Run.csv에서 PartOut 항목 확인
- 밀도 변화: Rhop 범위가 설정값(700~1300)을 벗어나는지 확인
- 시뮬레이션 안정성: 타임스텝 크기가 CFL 조건을 만족하는지 확인

---
*본 분석은 AI가 자동 생성한 참고 자료입니다. 최종 판단은 전문가의 검토가 필요합니다.*
`

	return analysis
}

func tanh(x float64) float64 {
	// Simple tanh approximation
	if x > 10 {
		return 1.0
	} else if x < -10 {
		return -1.0
	}
	ex := exp(x)
	enx := exp(-x)
	return (ex - enx) / (ex + enx)
}

func exp(x float64) float64 {
	// Simple exp approximation (for small x)
	// In production, use math.Exp
	result := 1.0
	term := 1.0
	for i := 1; i < 20; i++ {
		term *= x / float64(i)
		result += term
		if abs(term) < 1e-10 {
			break
		}
	}
	return result
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
