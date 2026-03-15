package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type MeasureToolParams struct {
	DataDir    string   `json:"data_dir"`
	PointsFile string   `json:"points_file"`
	OutCSV     string   `json:"out_csv"`
	Vars       []string `json:"vars,omitempty"`
	Elevation  *bool    `json:"elevation,omitempty"`
}

type measureToolTool struct{}

const (
	MeasureToolToolName    = "measuretool"
	measureToolDescription = `DualSPHysics MeasureTool — 특정 위치에서 수위, 압력, 속도 등의 시계열 데이터를 추출합니다.

사용법:
- data_dir: 시뮬레이션 출력 디렉토리
- points_file: 측정 포인트 파일 (x y z 좌표)
- out_csv: CSV 출력 파일 경로
- vars: 측정 변수 (기본: ["vel","rhop","press"])
- elevation: 수위 높이 계산 모드

Docker 컨테이너 내에서 실행됩니다.`
)

func NewMeasureToolTool() BaseTool {
	return &measureToolTool{}
}

func (m *measureToolTool) Info() ToolInfo {
	return ToolInfo{
		Name:        MeasureToolToolName,
		Description: measureToolDescription,
		Parameters: map[string]any{
			"data_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 출력 디렉토리",
			},
			"points_file": map[string]any{
				"type":        "string",
				"description": "측정 포인트 파일 경로 (예: /data/{case_name}/measure/probe_points.txt)",
			},
			"out_csv": map[string]any{
				"type":        "string",
				"description": "CSV 출력 파일 경로 (예: /data/{case_name}/measure/pressure)",
			},
			"vars": map[string]any{
				"type":        "array",
				"description": "측정 변수 (기본: [\"vel\",\"rhop\",\"press\"])",
				"items": map[string]any{
					"type": "string",
				},
			},
			"elevation": map[string]any{
				"type":        "boolean",
				"description": "수위 높이 계산 모드",
			},
		},
		Required: []string{"data_dir", "points_file", "out_csv"},
	}
}

func (m *measureToolTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params MeasureToolParams
	if err := UnmarshalToolInput(call.Input, &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.DataDir == "" || params.PointsFile == "" || params.OutCSV == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: data_dir, points_file, out_csv를 지정해주세요"), nil
	}

	// Validate absolute paths: check first directory component exists
	if isPathClearlyInvalid(params.PointsFile) {
		return NewTextErrorResponse(fmt.Sprintf("측정 포인트 파일 경로가 올바르지 않습니다: %s", params.PointsFile)), nil
	}
	if isPathClearlyInvalid(params.DataDir) {
		return NewTextErrorResponse(fmt.Sprintf("데이터 디렉토리 경로가 올바르지 않습니다: %s", params.DataDir)), nil
	}

	// Build Docker command
	args := []string{"compose", "run", "--rm", "dsph", "measuretool",
		fmt.Sprintf("-dirdata %s", params.DataDir),
		"-filexml AUTO",
		fmt.Sprintf("-points %s", params.PointsFile),
		fmt.Sprintf("-savecsv %s", params.OutCSV),
	}

	// Variables
	vars := params.Vars
	if len(vars) == 0 {
		vars = []string{"vel", "rhop", "press"}
	}
	args = append(args, fmt.Sprintf("-vars:%s", strings.Join(vars, ",")))

	// Elevation mode
	if params.Elevation != nil && *params.Elevation {
		args = append(args, "-elevation")
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = getWorkingDirectory()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("측정 데이터 추출 실패: %s\n출력: %s", err, string(output))), nil
	}

	resultDesc := "측정 데이터(.csv) 추출 완료."
	if params.Elevation != nil && *params.Elevation {
		resultDesc += " 수위 높이(Elevation) 데이터가 포함되었습니다."
	}

	return NewTextResponse(fmt.Sprintf("%s\n출력 경로: %s\n%s", resultDesc, params.OutCSV, string(output))), nil
}
