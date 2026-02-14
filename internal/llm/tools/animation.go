package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AnimationParams defines parameters for animation generation (ANI-01).
type AnimationParams struct {
	DataDir    string   `json:"data_dir"`
	OutFile    string   `json:"out_file"`
	Variables  []string `json:"variables,omitempty"`
	Colormap   string   `json:"colormap,omitempty"`
	ViewAngle  string   `json:"view_angle,omitempty"`
	Resolution []int    `json:"resolution,omitempty"`
	OnlyFluid  *bool    `json:"only_fluid,omitempty"`
	FPS        int      `json:"fps,omitempty"`
	StartTime  float64  `json:"start_time,omitempty"`
	EndTime    float64  `json:"end_time,omitempty"`
	Format     string   `json:"format,omitempty"` // "mp4" or "gif"
}

type animationTool struct{}

const (
	AnimationToolName    = "animation"
	animationDescription = `ParaView pvpython 애니메이션 생성 도구 (ANI-01) — 시뮬레이션 결과를 동영상으로 렌더링합니다.

사용법:
- data_dir: 시뮬레이션 결과 디렉토리 (VTK 파일 위치)
- out_file: 출력 파일 경로 (.mp4 또는 .gif)
- variables: 시각화 변수 (기본: ["Press"] 압력 컬러맵)
- colormap: 컬러맵 이름 (기본: "Blue to Red Rainbow")
- view_angle: 뷰 방향 (기본: "XZ" — 측면도)
- resolution: 이미지 해상도 [width, height] (기본: [1920, 1080])
- only_fluid: 유체 파티클만 (기본: true)
- fps: 프레임레이트 (기본: 30)
- start_time: 시작 시간 (초, 기본: 0)
- end_time: 종료 시간 (초, 기본: 마지막 타임스텝)
- format: 출력 포맷 ("mp4" 또는 "gif", 기본: 확장자에서 추론)

출력:
- 애니메이션 파일 (.mp4 또는 .gif)

Docker 컨테이너 내에서 pvpython + ffmpeg를 사용하여 생성합니다.`
)

func NewAnimationTool() BaseTool {
	return &animationTool{}
}

func (a *animationTool) Info() ToolInfo {
	return ToolInfo{
		Name:        AnimationToolName,
		Description: animationDescription,
		Parameters: map[string]any{
			"data_dir": map[string]any{
				"type":        "string",
				"description": "VTK 파일이 있는 디렉토리",
			},
			"out_file": map[string]any{
				"type":        "string",
				"description": "출력 파일 경로 (.mp4 또는 .gif)",
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
			"fps": map[string]any{
				"type":        "integer",
				"description": "프레임레이트 (FPS)",
			},
			"start_time": map[string]any{
				"type":        "number",
				"description": "시작 시간 (초)",
			},
			"end_time": map[string]any{
				"type":        "number",
				"description": "종료 시간 (초)",
			},
			"format": map[string]any{
				"type":        "string",
				"description": "출력 포맷 (mp4 또는 gif)",
			},
		},
		Required: []string{"data_dir", "out_file"},
	}
}

func (a *animationTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params AnimationParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.DataDir == "" || params.OutFile == "" {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: data_dir과 out_file을 지정해주세요"), nil
	}

	// Validate data directory
	if isPathClearlyInvalid(params.DataDir) || strings.Contains(params.DataDir, "..") {
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
	if params.FPS == 0 {
		params.FPS = 30
	}
	if params.Format == "" {
		ext := strings.ToLower(filepath.Ext(params.OutFile))
		if ext == ".gif" {
			params.Format = "gif"
		} else {
			params.Format = "mp4"
		}
	}

	// Generate pvpython animation script
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("animate_%s.py", call.ID))
	script := generateAnimationScript(params, onlyFluid)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("애니메이션 스크립트 생성 실패: %s", err)), nil
	}
	defer os.Remove(scriptPath)

	// Run pvpython locally with Mesa backend (NOT Docker)
	cmd := exec.CommandContext(ctx,
		"/opt/paraview/bin/pvpython", "--force-offscreen-rendering", "--mesa",
		scriptPath,
	)
	cmd.Dir = getWorkingDirectory()
	cmd.Env = append(os.Environ(),
		"DISPLAY=",
		"VTK_DEFAULT_RENDER_WINDOW_OFFSCREEN=1",
	)

	output, err := cmd.CombinedOutput()
	// Exit code 1 is normal for Mesa cleanup
	if err != nil && !isMesaCleanupExit(err, params.OutFile) {
		return NewTextErrorResponse(fmt.Sprintf("애니메이션 생성 실패: %s\n출력: %s", err, string(output))), nil
	}

	// Post-process with ffmpeg if needed (for gif conversion)
	if params.Format == "gif" {
		if err := convertToGif(ctx, params.OutFile, params.FPS); err != nil {
			return NewTextErrorResponse(fmt.Sprintf("GIF 변환 실패: %s", err)), nil
		}
	}

	resultMsg := fmt.Sprintf(`애니메이션 생성 완료:
- 출력 파일: %s
- 포맷: %s
- FPS: %d
- 해상도: %dx%d
`, params.OutFile, params.Format, params.FPS, params.Resolution[0], params.Resolution[1])

	return NewTextResponse(resultMsg + string(output)), nil
}

func generateAnimationScript(params AnimationParams, onlyFluid bool) string {
	// Apply defaults for safe direct calls
	if len(params.Variables) == 0 {
		params.Variables = []string{"Press"}
	}
	if params.Colormap == "" {
		params.Colormap = "Blue to Red Rainbow"
	}
	if params.ViewAngle == "" {
		params.ViewAngle = "XZ"
	}
	if len(params.Resolution) < 2 {
		params.Resolution = []int{1920, 1080}
	}
	if params.FPS == 0 {
		params.FPS = 30
	}

	// VTK file pattern — use the actual data_dir path
	vtkPattern := filepath.Join(params.DataDir, "PartFluid_*.vtk")
	if !onlyFluid {
		vtkPattern = filepath.Join(params.DataDir, "Part_*.vtk")
	}

	framesDir := filepath.Join(filepath.Dir(params.OutFile), "frames")
	cameraBlock := pvpythonCameraBlock(params.ViewAngle)

	// Mesa-compatible animation: frame-by-frame rendering + ffmpeg compilation
	// NO SaveAnimation (crashes Mesa), NO Text/AnnotateTimeFilter (crashes Mesa)
	return fmt.Sprintf(`#!/usr/bin/env pvpython
# Auto-generated Mesa-compatible animation script (ANI-01)
import os, sys, glob, subprocess
from paraview.simple import *

# --- Configuration ---
VTK_PATTERN = %q
OUTPUT_FILE = %q
OUTPUT_PNG_DIR = %q
RESOLUTION = [%d, %d]
FPS = %d
START_TIME = %g
END_TIME = %g
OUTPUT_FORMAT = %q

vtk_files = sorted(glob.glob(VTK_PATTERN))
if not vtk_files:
    print(f"ERROR: No VTK files found at {VTK_PATTERN}")
    sys.exit(1)

print(f"Found {len(vtk_files)} VTK files")
os.makedirs(OUTPUT_PNG_DIR, exist_ok=True)

# --- Setup ParaView pipeline ---
paraview.simple._DisableFirstRenderCameraReset()
reader = LegacyVTKReader(FileNames=vtk_files)

renderView = CreateRenderView()
renderView.ViewSize = RESOLUTION
renderView.Background = [0.1, 0.1, 0.15]

display = Show(reader, renderView)
display.Representation = 'Points'
display.PointSize = 4.0

ColorBy(display, ('POINTS', %q))
lut = GetColorTransferFunction(%q)
lut.ApplyPreset(%q, True)
lut.RescaleTransferFunction(0, 5000)

colorBar = GetScalarBar(lut, renderView)
colorBar.Title = '%s'
colorBar.ComponentTitle = ''
colorBar.TitleFontSize = 18
colorBar.LabelFontSize = 14
colorBar.Visibility = 1
colorBar.WindowLocation = 'Upper Right Corner'

%s

# --- Frame-by-frame rendering (Mesa compatible) ---
scene = GetAnimationScene()
scene.UpdateAnimationUsingDataTimeSteps()
all_timesteps = list(scene.TimeKeeper.TimestepValues)

# Filter timesteps by start/end time
if START_TIME > 0 or END_TIME > 0:
    end = END_TIME if END_TIME > 0 else all_timesteps[-1]
    timesteps = [t for t in all_timesteps if START_TIME <= t <= end]
else:
    timesteps = all_timesteps

print(f"Rendering {len(timesteps)} frames at {RESOLUTION[0]}x{RESOLUTION[1]}...")

for i, t in enumerate(timesteps):
    scene.AnimationTime = t
    Render()
    frame_path = os.path.join(OUTPUT_PNG_DIR, f"frame_{i:04d}.png")
    SaveScreenshot(frame_path, renderView, ImageResolution=RESOLUTION)
    if (i + 1) %% 10 == 0 or i == 0:
        print(f"  Frame {i+1}/{len(timesteps)} (t={t:.3f}s)")

frame_count = len(glob.glob(os.path.join(OUTPUT_PNG_DIR, "frame_*.png")))
print(f"Rendered {frame_count} frames")

if frame_count == 0:
    print("ERROR: No frames rendered")
    sys.exit(1)

# --- Compile with ffmpeg (text overlay via drawtext, not ParaView Text) ---
print("Compiling animation with ffmpeg...")
vf = "pad=ceil(iw/2)*2:ceil(ih/2)*2"

if OUTPUT_FORMAT == "gif":
    ffmpeg_cmd = [
        "ffmpeg", "-y",
        "-framerate", str(FPS),
        "-i", os.path.join(OUTPUT_PNG_DIR, "frame_%%04d.png"),
        "-vf", f"fps={FPS},scale={RESOLUTION[0]}:-1:flags=lanczos,{vf}",
        "-c:v", "gif",
        OUTPUT_FILE
    ]
else:
    ffmpeg_cmd = [
        "ffmpeg", "-y",
        "-framerate", str(FPS),
        "-i", os.path.join(OUTPUT_PNG_DIR, "frame_%%04d.png"),
        "-c:v", "libx264", "-preset", "medium", "-crf", "18",
        "-pix_fmt", "yuv420p", "-vf", vf,
        OUTPUT_FILE
    ]

ffmpeg_result = subprocess.run(ffmpeg_cmd, capture_output=True, text=True)
if ffmpeg_result.returncode == 0 and os.path.exists(OUTPUT_FILE):
    size_mb = os.path.getsize(OUTPUT_FILE) / (1024 * 1024)
    duration = frame_count / FPS
    print(f"Animation saved: {OUTPUT_FILE} ({size_mb:.1f} MB, {frame_count} frames, {duration:.1f}s)")
else:
    print(f"WARNING: ffmpeg failed: {ffmpeg_result.stderr[:500]}")

print("Done!")
`,
		vtkPattern, params.OutFile, framesDir,
		params.Resolution[0], params.Resolution[1], params.FPS,
		params.StartTime, params.EndTime, params.Format,
		params.Variables[0], params.Variables[0], params.Colormap,
		params.Variables[0],
		cameraBlock,
	)
}

func convertToGif(ctx context.Context, outFile string, fps int) error {
	// Check if file is already .gif (pvpython might output as .avi or image sequence)
	if strings.HasSuffix(outFile, ".gif") {
		// Assume pvpython created .mp4 or .avi, convert to gif
		tempVideo := strings.TrimSuffix(outFile, ".gif") + ".mp4"

		// Convert with ffmpeg
		cmd := exec.CommandContext(ctx, "ffmpeg",
			"-i", tempVideo,
			"-vf", fmt.Sprintf("fps=%d,scale=1920:-1:flags=lanczos", fps),
			"-c:v", "gif",
			"-y",
			outFile,
		)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("ffmpeg 변환 실패: %w", err)
		}

		// Clean up temp video
		os.Remove(tempVideo)
	}
	return nil
}
