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
	FPS          int      `json:"fps,omitempty"`         // ANI-01: Frames per second for animation
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
	if params.FPS == 0 {
		params.FPS = 30 // ANI-01: Default 30 FPS for animations
	}

	// Generate pvpython script
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("render_%s.py", call.ID))
	script := generatePvpythonScript(params, onlyFluid)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("렌더링 스크립트 생성 실패: %s", err)), nil
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
	// Exit code 1 is normal for Mesa cleanup — only treat as error if no output was generated
	if err != nil && !isMesaCleanupExit(err, params.OutFile) {
		return NewTextErrorResponse(fmt.Sprintf("렌더링 실행 실패: %s\n출력: %s", err, string(output))), nil
	}

	resultMsg := fmt.Sprintf("렌더링 완료. 출력: %s\n", params.OutFile)
	if params.RenderMode == "animation" {
		resultMsg += "애니메이션 파일이 생성되었습니다.\n"
	} else {
		resultMsg += "스냅샷 이미지가 생성되었습니다.\n"
	}

	return NewTextResponse(resultMsg + string(output)), nil
}

// isMesaCleanupExit checks if the error is just a Mesa cleanup exit code 1 (not a real error).
func isMesaCleanupExit(err error, outFile string) bool {
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		// Check if output file(s) exist — if so, it's just cleanup noise
		if _, statErr := os.Stat(outFile); statErr == nil {
			return true
		}
		// For frame-based output, check if frames directory exists
		framesDir := filepath.Join(filepath.Dir(outFile), "frames")
		if entries, _ := os.ReadDir(framesDir); len(entries) > 0 {
			return true
		}
	}
	return false
}

func generatePvpythonScript(params PvpythonParams, onlyFluid bool) string {
	// Camera setup based on view_angle
	cameraBlock := pvpythonCameraBlock(params.ViewAngle)

	// VTK file pattern — use the actual data_dir path (not Docker mount)
	vtkPattern := filepath.Join(params.DataDir, "PartFluid_*.vtk")
	if !onlyFluid {
		vtkPattern = filepath.Join(params.DataDir, "Part_*.vtk")
	}

	outDir := filepath.Dir(params.OutFile)
	framesDir := filepath.Join(outDir, "frames")

	if params.RenderMode == "animation" {
		// Animation mode: frame-by-frame rendering + ffmpeg compilation
		return fmt.Sprintf(`#!/usr/bin/env pvpython
# Auto-generated Mesa-compatible ParaView animation script
import os, sys, glob, subprocess
from paraview.simple import *

# --- Configuration ---
VTK_PATTERN = %q
OUTPUT_MP4 = %q
OUTPUT_PNG_DIR = %q
RESOLUTION = [%d, %d]
FPS = %d

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
timesteps = list(scene.TimeKeeper.TimestepValues)
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

# --- Compile to MP4 with ffmpeg ---
print("Compiling animation with ffmpeg...")
vf = "pad=ceil(iw/2)*2:ceil(ih/2)*2"
ffmpeg_result = subprocess.run([
    "ffmpeg", "-y",
    "-framerate", str(FPS),
    "-i", os.path.join(OUTPUT_PNG_DIR, "frame_%%04d.png"),
    "-c:v", "libx264", "-preset", "medium", "-crf", "18",
    "-pix_fmt", "yuv420p", "-vf", vf,
    OUTPUT_MP4
], capture_output=True, text=True)

if ffmpeg_result.returncode == 0 and os.path.exists(OUTPUT_MP4):
    size_mb = os.path.getsize(OUTPUT_MP4) / (1024 * 1024)
    duration = frame_count / FPS
    print(f"Animation saved: {OUTPUT_MP4} ({size_mb:.1f} MB, {frame_count} frames, {duration:.1f}s)")
else:
    print(f"WARNING: ffmpeg failed: {ffmpeg_result.stderr[:500]}")

print("Done!")
`,
			vtkPattern, params.OutFile, framesDir,
			params.Resolution[0], params.Resolution[1], params.FPS,
			params.Variables[0], params.Variables[0], params.Colormap,
			params.Variables[0],
			cameraBlock,
		)
	}

	// Snapshot mode: single frame rendering
	return fmt.Sprintf(`#!/usr/bin/env pvpython
# Auto-generated Mesa-compatible ParaView snapshot script
import os, sys, glob
from paraview.simple import *

VTK_PATTERN = %q
OUTPUT_FILE = %q
RESOLUTION = [%d, %d]

vtk_files = sorted(glob.glob(VTK_PATTERN))
if not vtk_files:
    print(f"ERROR: No VTK files found at {VTK_PATTERN}")
    sys.exit(1)

print(f"Found {len(vtk_files)} VTK files")

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
colorBar.Visibility = 1
colorBar.WindowLocation = 'Upper Right Corner'

%s

# Render last (or mid) timestep for snapshot
scene = GetAnimationScene()
scene.UpdateAnimationUsingDataTimeSteps()
timesteps = list(scene.TimeKeeper.TimestepValues)
mid = len(timesteps) // 2
scene.AnimationTime = timesteps[mid]
Render()

os.makedirs(os.path.dirname(OUTPUT_FILE), exist_ok=True)
SaveScreenshot(OUTPUT_FILE, renderView, ImageResolution=RESOLUTION)
print(f"Snapshot saved: {OUTPUT_FILE}")
`,
		vtkPattern, params.OutFile,
		params.Resolution[0], params.Resolution[1],
		params.Variables[0], params.Variables[0], params.Colormap,
		params.Variables[0],
		cameraBlock,
	)
}

// pvpythonCameraBlock returns Python code to set up the camera for a given view angle.
func pvpythonCameraBlock(viewAngle string) string {
	switch viewAngle {
	case "XY":
		return `# Camera: top-down (XY plane)
renderView.CameraPosition = [0.5, 0.25, 2.0]
renderView.CameraFocalPoint = [0.5, 0.25, 0.25]
renderView.CameraViewUp = [0, 1, 0]
renderView.ResetCamera()`
	case "YZ":
		return `# Camera: side view (YZ plane)
renderView.CameraPosition = [2.0, 0.25, 0.3]
renderView.CameraFocalPoint = [0.5, 0.25, 0.3]
renderView.CameraViewUp = [0, 0, 1]
renderView.ResetCamera()`
	case "ISO":
		return `# Camera: isometric view
renderView.CameraPosition = [0.5, -1.2, 0.6]
renderView.CameraFocalPoint = [0.5, 0.25, 0.25]
renderView.CameraViewUp = [0, 0, 1]
renderView.CameraParallelScale = 0.6`
	default: // "XZ" — default side view
		return `# Camera: side view (XZ plane)
renderView.CameraPosition = [0.5, -1.5, 0.3]
renderView.CameraFocalPoint = [0.5, 0.25, 0.3]
renderView.CameraViewUp = [0, 0, 1]
renderView.ResetCamera()`
	}
}
