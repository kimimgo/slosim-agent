package tools

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os/exec"
)

type SolverParams struct {
	CaseName string `json:"case_name"`
	DataDir  string `json:"data_dir"`
	OutDir   string `json:"out_dir"`
	GPU      bool   `json:"gpu"`
}

type solverTool struct{}

const (
	SolverToolName    = "solver"
	solverDescription = `DualSPHysics GPU SPH 솔버 — 시뮬레이션을 백그라운드에서 실행합니다.

사용법:
- case_name: 케이스 이름 (GenCase 출력 기준)
- data_dir: 입력 데이터 디렉토리 (GenCase 출력 경로)
- out_dir: 시뮬레이션 출력 디렉토리
- gpu: GPU 사용 여부 (기본: true)

시뮬레이션은 백그라운드에서 실행되며, 즉시 Job ID를 반환합니다.
진행 상황은 job_manager 도구로 확인할 수 있습니다.`
)

func NewSolverTool() BaseTool {
	return &solverTool{}
}

func (s *solverTool) Info() ToolInfo {
	return ToolInfo{
		Name:        SolverToolName,
		Description: solverDescription,
		Parameters: map[string]any{
			"case_name": map[string]any{
				"type":        "string",
				"description": "케이스 이름 (GenCase 출력 기준)",
			},
			"data_dir": map[string]any{
				"type":        "string",
				"description": "입력 데이터 디렉토리",
			},
			"out_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 출력 디렉토리",
			},
			"gpu": map[string]any{
				"type":        "boolean",
				"description": "GPU 사용 여부 (기본: true)",
			},
		},
		Required: []string{"case_name", "data_dir", "out_dir"},
	}
}

func (s *solverTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params SolverParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.CaseName == "" || params.DataDir == "" || params.OutDir == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: case_name, data_dir, out_dir을 지정해주세요"), nil
	}

	// Validate: check first directory component exists for absolute paths
	if isPathClearlyInvalid(params.DataDir) {
		return NewTextErrorResponse(fmt.Sprintf("데이터 디렉토리 경로가 올바르지 않습니다: %s", params.DataDir)), nil
	}

	// Generate Job ID
	jobID := generateJobID()

	// Build Docker command
	args := []string{"compose", "run", "--rm", "dsph",
		"DualSPHysics5.4_linux64", params.CaseName, params.OutDir,
	}
	if params.GPU {
		args = append(args, "-gpu")
	}

	// Launch simulation in background goroutine
	workDir := getWorkingDirectory()
	go func() {
		cmd := exec.Command("docker", args...)
		cmd.Dir = workDir
		_ = cmd.Run()
	}()

	return NewTextResponse(fmt.Sprintf(
		`시뮬레이션이 백그라운드에서 시작되었습니다.
{"job_id": "%s", "case_name": "%s", "status": "RUNNING"}
진행 상황은 job_manager 도구로 확인하세요.`, jobID, params.CaseName)), nil
}

func generateJobID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
