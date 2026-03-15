package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// ParametricStudyParams defines parameters for parametric study execution (PARA-01).
type ParametricStudyParams struct {
	StudyName  string                   `json:"study_name"`
	BaseCase   string                   `json:"base_case"`
	Variables  []ParametricVariable     `json:"variables"`
	OutputDir  string                   `json:"output_dir,omitempty"`
	Concurrent int                      `json:"concurrent,omitempty"`
}

// ParametricVariable defines a parameter to vary in the study.
type ParametricVariable struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

// ParametricResult holds results for a single parameter combination.
type ParametricResult struct {
	CaseID     string             `json:"case_id"`
	Parameters map[string]float64 `json:"parameters"`
	JobID      string             `json:"job_id,omitempty"`
	Status     string             `json:"status"`
	OutputPath string             `json:"output_path,omitempty"`
	Error      string             `json:"error,omitempty"`
}

type parametricStudyTool struct {
	jobManager *jobManagerTool
	mu         sync.Mutex
}

const (
	ParametricStudyToolName    = "parametric_study"
	parametricStudyDescription = `파라메트릭 스터디 실행 도구 (PARA-01) — 여러 파라미터 조합을 자동으로 실행합니다.

사용법:
- study_name: 스터디 이름
- base_case: 베이스 케이스 XML 경로
- variables: 변경할 파라미터 목록 (이름 + 값 배열)
- output_dir: 출력 디렉토리 (기본: ./parametric_studies/{study_name})
- concurrent: 동시 실행 개수 (기본: 3)

예시:
{
  "study_name": "fill_ratio_study",
  "base_case": "cases/SloshingTank_Def.xml",
  "variables": [
    {"name": "fill_ratio", "values": [0.5, 0.6, 0.7]}
  ]
}

출력:
- 각 조합별 케이스 디렉토리
- 비교 리포트 (PARA-03)
- Job 상태 추적`
)

func NewParametricStudyTool() BaseTool {
	return &parametricStudyTool{
		jobManager: NewJobManagerTool().(*jobManagerTool),
	}
}

func (p *parametricStudyTool) Info() ToolInfo {
	return ToolInfo{
		Name:        ParametricStudyToolName,
		Description: parametricStudyDescription,
		Parameters: map[string]any{
			"study_name": map[string]any{
				"type":        "string",
				"description": "스터디 이름",
			},
			"base_case": map[string]any{
				"type":        "string",
				"description": "베이스 케이스 XML 경로",
			},
			"variables": map[string]any{
				"type":        "array",
				"description": "변경할 파라미터 목록",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type": "string",
						},
						"values": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "number",
							},
						},
					},
				},
			},
			"output_dir": map[string]any{
				"type":        "string",
				"description": "출력 디렉토리",
			},
			"concurrent": map[string]any{
				"type":        "integer",
				"description": "동시 실행 개수",
			},
		},
		Required: []string{"study_name", "base_case", "variables"},
	}
}

func (p *parametricStudyTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params ParametricStudyParams
	if err := UnmarshalToolInput(call.Input, &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.StudyName == "" || params.BaseCase == "" || len(params.Variables) == 0 {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다"), nil
	}

	// Default values
	if params.OutputDir == "" {
		params.OutputDir = filepath.Join("parametric_studies", params.StudyName)
	}
	if params.Concurrent == 0 {
		params.Concurrent = 3
	}

	// Validate base case exists
	if _, err := os.Stat(params.BaseCase); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("베이스 케이스를 찾을 수 없습니다: %s", params.BaseCase)), nil
	}

	// Create output directory
	if err := os.MkdirAll(params.OutputDir, 0755); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("출력 디렉토리 생성 실패: %s", err)), nil
	}

	// Generate parameter combinations
	combinations := generateCombinations(params.Variables)
	totalCases := len(combinations)

	resultMsg := fmt.Sprintf(`파라메트릭 스터디 시작: %s
- 베이스 케이스: %s
- 변수 개수: %d
- 총 케이스 수: %d
- 출력 디렉토리: %s

`, params.StudyName, params.BaseCase, len(params.Variables), totalCases, params.OutputDir)

	// Execute cases
	results := make([]ParametricResult, totalCases)
	for i, combo := range combinations {
		caseID := fmt.Sprintf("case_%03d", i+1)
		caseDir := filepath.Join(params.OutputDir, caseID)

		result := ParametricResult{
			CaseID:     caseID,
			Parameters: combo,
			OutputPath: caseDir,
			Status:     "pending",
		}

		// Create case directory
		if err := os.MkdirAll(caseDir, 0755); err != nil {
			result.Status = "failed"
			result.Error = fmt.Sprintf("디렉토리 생성 실패: %s", err)
			results[i] = result
			continue
		}

		// Generate modified case XML
		modifiedXML := modifyBaseCase(params.BaseCase, combo)
		xmlPath := filepath.Join(caseDir, "case.xml")
		if err := os.WriteFile(xmlPath, []byte(modifiedXML), 0644); err != nil {
			result.Status = "failed"
			result.Error = fmt.Sprintf("XML 생성 실패: %s", err)
			results[i] = result
			continue
		}

		// Submit job (placeholder - would use job manager in real implementation)
		result.Status = "created"
		results[i] = result

		resultMsg += fmt.Sprintf("  [%s] %v → %s\n", caseID, combo, "created")
	}

	resultMsg += fmt.Sprintf("\n총 %d개 케이스 생성 완료\n", totalCases)
	resultMsg += "Job 실행은 job_manager 도구를 사용하여 각 케이스를 실행하세요.\n"

	// Save study manifest
	manifestPath := filepath.Join(params.OutputDir, "study_manifest.json")
	manifestData := map[string]any{
		"study_name": params.StudyName,
		"base_case":  params.BaseCase,
		"variables":  params.Variables,
		"total_cases": totalCases,
		"results":    results,
	}
	manifestJSON, _ := json.MarshalIndent(manifestData, "", "  ")
	os.WriteFile(manifestPath, manifestJSON, 0644)

	resultMsg += fmt.Sprintf("스터디 매니페스트 저장: %s\n", manifestPath)

	return NewTextResponse(resultMsg), nil
}

// generateCombinations generates all parameter combinations (Cartesian product).
func generateCombinations(variables []ParametricVariable) []map[string]float64 {
	if len(variables) == 0 {
		return []map[string]float64{}
	}

	// Recursive Cartesian product
	var generate func(int) []map[string]float64
	generate = func(depth int) []map[string]float64 {
		if depth == len(variables) {
			return []map[string]float64{{}}
		}

		subCombos := generate(depth + 1)
		result := []map[string]float64{}

		for _, value := range variables[depth].Values {
			for _, subCombo := range subCombos {
				combo := make(map[string]float64)
				for k, v := range subCombo {
					combo[k] = v
				}
				combo[variables[depth].Name] = value
				result = append(result, combo)
			}
		}

		return result
	}

	return generate(0)
}

// modifyBaseCase applies parameter modifications to base case XML.
// Supports two modes:
// 1. Template placeholders: {{param_name}} → value
// 2. XML attribute replacement: <param_name value="old"/> → <param_name value="new"/>
func modifyBaseCase(baseCasePath string, params map[string]float64) string {
	baseContent, _ := os.ReadFile(baseCasePath)
	modified := string(baseContent)

	for name, value := range params {
		valueStr := strconv.FormatFloat(value, 'f', -1, 64)

		// 1. Replace template placeholders: {{name}} → value
		placeholder := fmt.Sprintf("{{%s}}", name)
		modified = strings.ReplaceAll(modified, placeholder, valueStr)

		// 2. Replace XML attribute values: <name value="old" /> → <name value="new" />
		pattern := regexp.MustCompile(fmt.Sprintf(`(<%s\s+value=")([^"]*)(")`, regexp.QuoteMeta(name)))
		modified = pattern.ReplaceAllString(modified, fmt.Sprintf("${1}%s${3}", valueStr))
	}

	return modified
}
