package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/opencode-ai/opencode/internal/config"
)

type GenCaseParams struct {
	CasePath string   `json:"case_path"`
	SavePath string   `json:"save_path"`
	DP       *float64 `json:"dp,omitempty"`
}

type genCaseTool struct{}

const (
	GenCaseToolName    = "gencase"
	genCaseDescription = `DualSPHysics GenCase 도구 — XML 케이스 파일에서 파티클 지오메트리를 생성합니다.

사용법:
- case_path: XML 케이스 파일 경로 (.xml 확장자는 자동 제거됨)
- save_path: 출력 디렉토리 경로
- dp: 파티클 간격 오버라이드 (선택사항, 단위: m)

출력물:
- {save_path}/{case_name}.bi4 — 바이너리 파티클 데이터
- {save_path}/{case_name}.xml — 처리된 XML 복사본

Docker 컨테이너 내에서 실행됩니다.`
)

func NewGenCaseTool() BaseTool {
	return &genCaseTool{}
}

func (g *genCaseTool) Info() ToolInfo {
	return ToolInfo{
		Name:        GenCaseToolName,
		Description: genCaseDescription,
		Parameters: map[string]any{
			"case_path": map[string]any{
				"type":        "string",
				"description": "XML 케이스 파일 경로 (.xml 확장자 제외)",
			},
			"save_path": map[string]any{
				"type":        "string",
				"description": "출력 디렉토리 경로",
			},
			"dp": map[string]any{
				"type":        "number",
				"description": "파티클 간격 오버라이드 (m)",
			},
		},
		Required: []string{"case_path", "save_path"},
	}
}

func (g *genCaseTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params GenCaseParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.CasePath == "" || params.SavePath == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: case_path와 save_path를 지정해주세요"), nil
	}

	// Strip .xml extension (GenCase adds it automatically)
	casePath := strings.TrimSuffix(params.CasePath, ".xml")

	// Validate: for absolute paths, check that the XML file exists on host
	if filepath.IsAbs(casePath) {
		if _, err := os.Stat(casePath + ".xml"); os.IsNotExist(err) {
			return NewTextErrorResponse(fmt.Sprintf("케이스 파일을 찾을 수 없습니다: %s.xml", casePath)), nil
		}
	}

	// Build Docker command
	args := []string{"compose", "run", "--rm", "dsph", "GenCase", casePath, fmt.Sprintf("-save:%s", params.SavePath)}
	if params.DP != nil {
		args = append(args, fmt.Sprintf("-dp:%g", *params.DP))
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = getWorkingDirectory()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("해석 준비(GenCase) 실행 실패: %s\n출력: %s", err, string(output))), nil
	}

	return NewTextResponse(fmt.Sprintf("해석 준비 완료. 파티클 데이터(.bi4)가 %s에 생성되었습니다.\n%s", params.SavePath, string(output))), nil
}

// getWorkingDirectory returns the project working directory.
// Falls back to os.Getwd() if config is not loaded (avoids panic in tests).
func getWorkingDirectory() string {
	var dir string
	func() {
		defer func() { recover() }()
		dir = config.WorkingDirectory()
	}()
	if dir != "" {
		return dir
	}
	wd, _ := os.Getwd()
	return wd
}

// isPathClearlyInvalid checks if an absolute path is clearly invalid
// by verifying the first directory component exists.
// e.g., "/nonexistent/foo" → true (first component doesn't exist)
//
//	"/tmp/foo/bar" → false (first component /tmp exists)
func isPathClearlyInvalid(p string) bool {
	if !filepath.IsAbs(p) {
		return false
	}
	parts := strings.SplitN(filepath.Clean(p), string(filepath.Separator), 3)
	if len(parts) >= 2 && parts[1] != "" {
		firstDir := string(filepath.Separator) + parts[1]
		if _, err := os.Stat(firstDir); os.IsNotExist(err) {
			return true
		}
	}
	return false
}
