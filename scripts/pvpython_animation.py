"""
slosim-agent E2E 시각화: 슬로싱 시뮬레이션 통합 애니메이션
- 3D 파티클 뷰 (압력 컬러맵 + 컬러바)
- 탱크 와이어프레임 오버레이
- 타임스탬프 오버레이 (ffmpeg drawtext)

Usage: /opt/paraview/bin/pvpython --force-offscreen-rendering --mesa scripts/pvpython_animation.py
Output: simulations/e2e_test/sloshing_animation.mp4
"""

import os
import sys
import glob
import subprocess

from paraview.simple import *

# --- Configuration ---
DATA_DIR = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                        "simulations", "e2e_test")
VTK_PATTERN = os.path.join(DATA_DIR, "PartFluid_*.vtk")
OUTPUT_DIR = DATA_DIR
OUTPUT_MP4 = os.path.join(OUTPUT_DIR, "sloshing_animation.mp4")
OUTPUT_PNG_DIR = os.path.join(OUTPUT_DIR, "frames")

RESOLUTION = [1920, 1080]
FPS = 25

# Tank dimensions from XML
TANK_L, TANK_W, TANK_H = 1.0, 0.5, 0.6

# --- Collect VTK files ---
vtk_files = sorted(glob.glob(VTK_PATTERN))
if not vtk_files:
    print(f"ERROR: No VTK files found at {VTK_PATTERN}")
    sys.exit(1)

print(f"Found {len(vtk_files)} VTK files")
os.makedirs(OUTPUT_PNG_DIR, exist_ok=True)

# --- Setup ParaView pipeline ---
paraview.simple._DisableFirstRenderCameraReset()

# Load VTK series
reader = LegacyVTKReader(FileNames=vtk_files)

# Create render view
renderView = CreateRenderView()
renderView.ViewSize = RESOLUTION
renderView.Background = [0.1, 0.1, 0.15]
# NOTE: UseGradientBackground crashes with Box() on Mesa backend

# Display particles
display = Show(reader, renderView)
display.Representation = 'Points'
display.PointSize = 4.0

# Color by Pressure
ColorBy(display, ('POINTS', 'Press'))
pressLUT = GetColorTransferFunction('Press')
pressLUT.ApplyPreset('Blue to Red Rainbow', True)
pressLUT.RescaleTransferFunction(0, 5000)

# Color bar
pressBar = GetScalarBar(pressLUT, renderView)
pressBar.Title = 'Pressure [Pa]'
pressBar.ComponentTitle = ''
pressBar.TitleFontSize = 18
pressBar.LabelFontSize = 14
pressBar.Visibility = 1
pressBar.WindowLocation = 'Upper Right Corner'

# Camera: isometric view
renderView.CameraPosition = [0.5, -1.2, 0.6]
renderView.CameraFocalPoint = [0.5, 0.25, 0.25]
renderView.CameraViewUp = [0, 0, 1]
renderView.CameraParallelScale = 0.6

# Tank wireframe box
tank = Box()
tank.XLength = TANK_L
tank.YLength = TANK_W
tank.ZLength = TANK_H
tank.Center = [TANK_L/2, TANK_W/2, TANK_H/2]

tankDisplay = Show(tank, renderView)
tankDisplay.Representation = 'Wireframe'
tankDisplay.AmbientColor = [0.7, 0.7, 0.8]
tankDisplay.DiffuseColor = [0.7, 0.7, 0.8]
tankDisplay.LineWidth = 2.0
tankDisplay.Opacity = 0.5

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
    if (i + 1) % 10 == 0 or i == 0:
        sz = os.path.getsize(frame_path) if os.path.exists(frame_path) else 0
        print(f"  Frame {i+1}/{len(timesteps)} (t={t:.3f}s) {sz//1024}KB")

frame_count = len(glob.glob(os.path.join(OUTPUT_PNG_DIR, "frame_*.png")))
print(f"Rendered {frame_count} frames")

if frame_count == 0:
    print("ERROR: No frames rendered")
    sys.exit(1)

# --- Compile to MP4 with ffmpeg (add text overlay here) ---
# Use ffmpeg drawtext for timestamp + title (more reliable than ParaView Text source with Mesa)
dt = 0.02  # TimeOut from XML
print("Compiling animation with ffmpeg (+ text overlay)...")
drawtext_title = (
    "drawtext=text='Sloshing Tank  1.0m x 0.5m x 0.6m  |  dp=0.01m  |  168K particles'"
    ":fontsize=28:fontcolor=white:x=30:y=30:fontfile=/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
)
drawtext_time = (
    "drawtext=text='t = %{eif\\:n*" + f"{dt:.4f}" + "\\:d\\:3} s'"
    ":fontsize=32:fontcolor=yellow:x=30:y=(h-60):fontfile=/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
)
vf = f"{drawtext_title},{drawtext_time},pad=ceil(iw/2)*2:ceil(ih/2)*2"

ffmpeg_result = subprocess.run(
    [
        "ffmpeg", "-y",
        "-framerate", str(FPS),
        "-i", os.path.join(OUTPUT_PNG_DIR, "frame_%04d.png"),
        "-c:v", "libx264",
        "-preset", "medium",
        "-crf", "18",
        "-pix_fmt", "yuv420p",
        "-vf", vf,
        OUTPUT_MP4
    ],
    capture_output=True,
    text=True
)

if ffmpeg_result.returncode == 0 and os.path.exists(OUTPUT_MP4):
    size_mb = os.path.getsize(OUTPUT_MP4) / (1024 * 1024)
    duration = frame_count / FPS
    print(f"\nAnimation saved: {OUTPUT_MP4}")
    print(f"  {size_mb:.1f} MB | {frame_count} frames | {FPS} FPS | {duration:.1f}s duration")
else:
    print(f"\nWARNING: ffmpeg failed: {ffmpeg_result.stderr[:500]}")
    # Fallback: no text overlay
    subprocess.run([
        "ffmpeg", "-y",
        "-framerate", str(FPS),
        "-i", os.path.join(OUTPUT_PNG_DIR, "frame_%04d.png"),
        "-c:v", "libx264", "-preset", "medium", "-crf", "18",
        "-pix_fmt", "yuv420p",
        "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
        OUTPUT_MP4
    ], capture_output=True)
    if os.path.exists(OUTPUT_MP4):
        print(f"Fallback (no overlay) saved: {OUTPUT_MP4}")

# --- Peak snapshot ---
mid_step = len(timesteps) // 2
scene.AnimationTime = timesteps[mid_step]
Render()
snapshot_path = os.path.join(OUTPUT_DIR, "snapshot_peak.png")
SaveScreenshot(snapshot_path, renderView, ImageResolution=RESOLUTION)
print(f"Peak snapshot: {snapshot_path}")

print("\nDone!")
