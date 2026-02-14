package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PartVTKParams struct {
	DataDir   string   `json:"data_dir"`
	OutFile   string   `json:"out_file"`
	OnlyFluid *bool    `json:"only_fluid,omitempty"`
	Vars      []string `json:"vars,omitempty"`
	First     *int     `json:"first,omitempty"`
	Last      *int     `json:"last,omitempty"`
}

type partVTKTool struct{}

const (
	PartVTKToolName    = "partvtk"
	partVTKDescription = `DualSPHysics PartVTK 도구 — 시뮬레이션 결과를 VTK 형식으로 변환합니다.

사용법:
- data_dir: 시뮬레이션 출력 디렉토리 (Part_*.bi4 위치)
- out_file: VTK 출력 파일 경로 (접미사 자동 추가)
- only_fluid: 유체 파티클만 추출 (기본: true)
- vars: 출력 변수 목록 (기본: ["vel","rhop","press"])
- first: 시작 타임스텝
- last: 종료 타임스텝

Docker 컨테이너 내에서 실행됩니다.`
)

func NewPartVTKTool() BaseTool {
	return &partVTKTool{}
}

func (p *partVTKTool) Info() ToolInfo {
	return ToolInfo{
		Name:        PartVTKToolName,
		Description: partVTKDescription,
		Parameters: map[string]any{
			"data_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 출력 디렉토리 (Part_*.bi4 위치)",
			},
			"out_file": map[string]any{
				"type":        "string",
				"description": "VTK 출력 파일 경로",
			},
			"only_fluid": map[string]any{
				"type":        "boolean",
				"description": "유체 파티클만 추출 (기본: true)",
			},
			"vars": map[string]any{
				"type":        "array",
				"description": "출력 변수 (기본: [\"vel\",\"rhop\",\"press\"])",
				"items": map[string]any{
					"type": "string",
				},
			},
			"first": map[string]any{
				"type":        "integer",
				"description": "시작 타임스텝",
			},
			"last": map[string]any{
				"type":        "integer",
				"description": "종료 타임스텝",
			},
		},
		Required: []string{"data_dir", "out_file"},
	}
}

func (p *partVTKTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params PartVTKParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.DataDir == "" || params.OutFile == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: data_dir과 out_file을 지정해주세요"), nil
	}

	// Validate absolute paths: check first directory component exists
	if isPathClearlyInvalid(params.DataDir) {
		return NewTextErrorResponse(fmt.Sprintf("데이터 디렉토리 경로가 올바르지 않습니다: %s", params.DataDir)), nil
	}

	// Build Docker command
	args := []string{"compose", "run", "--rm", "dsph", "partvtk",
		fmt.Sprintf("-dirdata %s", params.DataDir),
		"-filexml AUTO",
		fmt.Sprintf("-savevtk %s", params.OutFile),
	}

	// Only fluid (default: true)
	onlyFluid := true
	if params.OnlyFluid != nil {
		onlyFluid = *params.OnlyFluid
	}
	if onlyFluid {
		args = append(args, "-onlytype:-bound")
	}

	// Variables
	vars := params.Vars
	if len(vars) == 0 {
		vars = []string{"vel", "rhop", "press"}
	}
	args = append(args, fmt.Sprintf("-vars:%s", strings.Join(vars, ",")))

	// Timestep range
	if params.First != nil {
		args = append(args, fmt.Sprintf("-first:%d", *params.First))
	}
	if params.Last != nil {
		args = append(args, fmt.Sprintf("-last:%d", *params.Last))
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = getWorkingDirectory()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("결과 변환(PartVTK) 실행 실패: %s\n출력: %s", err, string(output))), nil
	}

	return NewTextResponse(fmt.Sprintf("결과 변환 완료. VTK 파일(.vtk)이 %s에 생성되었습니다.\n%s", params.OutFile, string(output))), nil
}
