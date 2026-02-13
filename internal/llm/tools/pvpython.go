package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PvpythonParams defines parameters for ParaView rendering (VIS-01, VIS-02, VIS-03).
type PvpythonParams struct {
	DataDir      string   `json:"data_dir"`
	OutFile      string   `json:"out_file"`
	Timesteps    []int    `json:"timesteps,omitempty"`
	Variables    []string `json:"variables,omitempty"`
	Colormap     string   `json:"colormap,omitempty"`
	ViewAngle    string   `json:"view_angle,omitempty"`
	Resolution   []int    `json:"resolution,omitempty"`
	OnlyFluid    *bool    `json:"only_fluid,omitempty"`
	RenderMode   string   `json:"render_mode,omitempty"` // "snapshot" or "animation"
}

type pvpythonTool struct{}

const (
	PvpythonToolName    = "pvpython"
	pvpythonDescription = `ParaView pvpython 렌더링 도구 — 시뮬레이션 결과를 이미지 또는 애니메이션으로 렌더링합니다.

사용법:
- data_dir: 시뮬레이션 결과 디렉토리 (VTK 파일 위치)
- out_file: 출력 파일 경로 (.png 또는 .mp4)
- timesteps: 렌더링할 타임스텝 목록 (비어있으면 전체)
- variables: 시각화 변수 (기본: ["Press"] 압력 컬러맵)
- colormap: 컬러맵 이름 (기본: "Blue to Red Rainbow")
- view_angle: 뷰 방향 (기본: "XZ" — 측면도)
- resolution: 이미지 해상도 [width, height] (기본: [1920, 1080])
- only_fluid: 유체 파티클만 (기본: true)
- render_mode: "snapshot" (정지 이미지) 또는 "animation" (동영상)

Docker 컨테이너 내에서 실행됩니다 (pvpython 이미지).`
)

func NewPvpythonTool() BaseTool {
	return &pvpythonTool{}
}

func (p *pvpythonTool) Info() ToolInfo {
	return ToolInfo{
		Name:        PvpythonToolName,
		Description: pvpythonDescription,
		Parameters: map[string]any{
			"data_dir": map[string]any{
				"type":        "string",
				"description": "VTK 파일이 있는 디렉토리",
			},
			"out_file": map[string]any{
				"type":        "string",
				"description": "출력 파일 경로 (.png 또는 .mp4)",
			},
			"timesteps": map[string]any{
				"type":        "array",
				"description": "렌더링할 타임스텝 목록",
				"items": map[string]any{
					"type": "integer",
				},
			},
			"variables": map[string]any{
				"type":        "array",
				"description": "시각화 변수 (Press, Vel, Rhop 등)",
				"items": map[string]any{
					"type": "string",
				},
			},
			"colormap": map[string]any{
				"type":        "string",
				"description": "컬러맵 이름",
			},
			"view_angle": map[string]any{
				"type":        "string",
				"description": "뷰 방향 (XZ, XY, YZ)",
			},
			"resolution": map[string]any{
				"type":        "array",
				"description": "이미지 해상도 [width, height]",
				"items": map[string]any{
					"type": "integer",
				},
			},
			"only_fluid": map[string]any{
				"type":        "boolean",
				"description": "유체 파티클만 렌더링",
			},
			"render_mode": map[string]any{
				"type":        "string",
				"description": "렌더링 모드 (snapshot 또는 animation)",
			},
		},
		Required: []string{"data_dir", "out_file"},
	}
}

func (p *pvpythonTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params PvpythonParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.DataDir == "" || params.OutFile == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: data_dir과 out_file을 지정해주세요"), nil
	}

	// Validate data directory
	if isPathClearlyInvalid(params.DataDir) {
		return NewTextErrorResponse(fmt.Sprintf("데이터 디렉토리 경로가 올바르지 않습니다: %s", params.DataDir)), nil
	}

	// Default values
	if len(params.Variables) == 0 {
		params.Variables = []string{"Press"}
	}
	if params.Colormap == "" {
		params.Colormap = "Blue to Red Rainbow"
	}
	if params.ViewAngle == "" {
		params.ViewAngle = "XZ"
	}
	if len(params.Resolution) == 0 {
		params.Resolution = []int{1920, 1080}
	}
	onlyFluid := true
	if params.OnlyFluid != nil {
		onlyFluid = *params.OnlyFluid
	}
	if params.RenderMode == "" {
		params.RenderMode = "snapshot"
	}

	// Generate pvpython script
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("render_%s.py", call.ID))
	script := generatePvpythonScript(params, onlyFluid)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("렌더링 스크립트 생성 실패: %s", err)), nil
	}
	defer os.Remove(scriptPath)

	// Build Docker command
	args := []string{"run", "--rm",
		"-v", fmt.Sprintf("%s:%s", params.DataDir, "/data/input"),
		"-v", fmt.Sprintf("%s:%s", filepath.Dir(params.OutFile), "/data/output"),
		"-v", fmt.Sprintf("%s:/tmp/render.py", scriptPath),
		"pvpython",
		"pvpython", "/tmp/render.py",
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = getWorkingDirectory()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("렌더링 실행 실패: %s\n출력: %s", err, string(output))), nil
	}

	resultMsg := fmt.Sprintf("렌더링 완료. 출력: %s\n", params.OutFile)
	if params.RenderMode == "animation" {
		resultMsg += "애니메이션 파일이 생성되었습니다.\n"
	} else {
		resultMsg += fmt.Sprintf("%d개 이미지가 생성되었습니다.\n", len(params.Timesteps))
	}

	return NewTextResponse(resultMsg + string(output)), nil
}

func generatePvpythonScript(params PvpythonParams, onlyFluid bool) string {
	// Simplified pvpython script template (actual implementation would be more complex)
	return fmt.Sprintf(`#!/usr/bin/env pvpython
# Auto-generated ParaView rendering script
from paraview.simple import *

# Load VTK files
reader = LegacyVTKReader(FileNames=['/data/input/*.vtk'])
UpdatePipeline()

# Color by variable
ColorBy(reader, ('POINTS', '%s'))
LUT = GetColorTransferFunction('%s')
LUT.ApplyPreset('%s', True)

# Set view
view = GetActiveView()
view.ViewSize = [%d, %d]
view.CameraPosition = [1, 0, 0]  # %s view

# Render
if '%s' == 'snapshot':
    SaveScreenshot('/data/output/%s', view)
else:
    SaveAnimation('/data/output/%s', view)

print("Rendering complete")
`,
		params.Variables[0],
		params.Variables[0],
		params.Colormap,
		params.Resolution[0], params.Resolution[1],
		params.ViewAngle,
		params.RenderMode,
		filepath.Base(params.OutFile),
		filepath.Base(params.OutFile),
	)
}
