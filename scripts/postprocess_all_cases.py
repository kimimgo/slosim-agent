#!/usr/bin/env python3
"""Post-process all sloshing simulation cases using IsoSurface + pvpython.

Pipeline per case:
  1. DualSPHysics IsoSurface (Docker) -> VTK triangle mesh (bi4 -> surface)
  2. pvpython (Mesa) -> snapshot PNGs (Vel + Press, isometric 3D view)
  3. pvpython (Mesa) -> animation frames (time-stepping through IsoSurface VTKs)
  4. ffmpeg -> MP4 with text overlay
  5. Field statistics extraction (JSON)
"""
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path

# Case Definitions
CASES = {
    "e2e_test": {
        "label": "Basic Sloshing Benchmark",
        "motion": "Sway 1.0Hz, A=0.05m",
        "tank": "1.0x0.5x0.6m rectangular",
        "dp": 0.01, "particles": 145530,
        "data_dir": "data",
        "text_overlay": "Sloshing Benchmark | 1.0Hz, A=0.05m",
    },
    "chen_nearcrit": {
        "label": "Chen 2018 Near-Resonance",
        "motion": "Sway 1.008Hz, A=0.007m",
        "tank": "0.6x0.3x0.65m rectangular, fill=50%",
        "dp": 0.005, "particles": 259777,
        "data_dir": "data",
        "text_overlay": "Chen 2018 Near-Resonance | f=1.008Hz",
    },
    "chen_shallow": {
        "label": "Chen 2018 Shallow Fill",
        "motion": "Sway 0.756Hz, A=0.007m",
        "tank": "0.6x0.3x0.65m rectangular, fill=20%",
        "dp": 0.005, "particles": 119357,
        "data_dir": "data",
        "text_overlay": "Chen 2018 Shallow Fill | 0.756Hz, h/L=0.2",
    },
    "english_dbc": {
        "label": "English 2021 DBC (SPHERIC)",
        "motion": "Roll 0.613Hz, theta=4deg",
        "tank": "0.9x0.062x0.508m quasi-2D",
        "dp": 0.002, "particles": 633059,
        "data_dir": "data",
        "text_overlay": "English 2021 DBC | Roll 0.613Hz, theta=4deg",
    },
    "isope_lng": {
        "label": "ISOPE 2012 LNG Mark III",
        "motion": "Sway+Roll 0.6Hz, A=0.02m + theta=3deg",
        "tank": "0.946x0.118x0.67m Mark III",
        "dp": 0.005, "particles": 301283,
        "data_dir": "data",
        "text_overlay": "ISOPE LNG Mark III | Sway+Roll 0.6Hz",
    },
    "liu_amp1": {
        "label": "Liu 2024 Pitch Low Amp",
        "motion": "Pitch 0.87Hz, theta=1deg",
        "tank": "1.0x0.5x1.0m rectangular",
        "dp": 0.01, "particles": 668856,
        "data_dir": "data",
        "text_overlay": "Liu 2024 Pitch | 0.87Hz, theta=1deg (Low Amp)",
    },
    "liu_large": {
        "label": "Liu 2024 Pitch Large Amp",
        "motion": "Pitch 0.66Hz, theta=2deg",
        "tank": "1.0x0.5x1.0m rectangular",
        "dp": 0.015, "particles": 192200,
        "data_dir": "data",
        "text_overlay": "Liu 2024 Pitch | 0.66Hz, theta=2deg (Large Amp)",
    },
    "nasa_cylinder": {
        "label": "NASA Cylindrical Anti-slosh Baffle",
        "motion": "Sway 0.5Hz, A=0.01m",
        "tank": "R=1.42m cylindrical (D=2.84m, H=3.0m)",
        "dp": 0.03, "particles": 143460,
        "data_dir": "data",
        "text_overlay": "NASA Cylindrical | Sway 0.5Hz, A=0.01m",
    },
    "spheric_low": {
        "label": "SPHERIC Test 10 Low Res",
        "motion": "Roll 0.613Hz, theta=4deg",
        "tank": "0.9x0.062x0.508m quasi-2D",
        "dp": 0.008, "particles": 72128,
        "data_dir": "data",
        "text_overlay": "SPHERIC Test 10 | Roll 0.613Hz, theta=4deg",
    },
    "zhao_horizcyl": {
        "label": "Zhao 2024 Horizontal Cylinder LNG",
        "motion": "Roll 0.5496Hz, theta=3deg",
        "tank": "L=3.0m, D=1.0m horizontal cylinder",
        "dp": 0.01, "particles": 529434,
        "data_dir": "data",
        "text_overlay": "Zhao 2024 Horiz Cyl | Roll 0.55Hz, theta=3deg",
    },
    "spheric_high": {
        "label": "SPHERIC Test 10 High Res",
        "motion": "Roll 0.613Hz, theta=4deg",
        "tank": "0.9x0.062x0.508m quasi-2D",
        "dp": 0.004, "particles": 344000,
        "data_dir": "data",
        "text_overlay": "SPHERIC Test 10 High | Roll 0.613Hz, theta=4deg",
    },
    "spheric_oil_low": {
        "label": "SPHERIC Test 10 Oil Low Res",
        "motion": "Roll 0.613Hz, theta=4deg",
        "tank": "0.9x0.062x0.508m quasi-2D",
        "dp": 0.008, "particles": 72128,
        "data_dir": "data",
        "text_overlay": "SPHERIC Oil Low | Roll 0.613Hz, theta=4deg",
    },
    "liu_amp2": {
        "label": "Liu 2024 Pitch Mid Amp",
        "motion": "Pitch 0.87Hz, theta=2deg",
        "tank": "1.0x0.5x1.0m rectangular",
        "dp": 0.01, "particles": 668856,
        "data_dir": "data",
        "text_overlay": "Liu 2024 Pitch | 0.87Hz, theta=2deg (Mid Amp)",
    },
    "liu_amp3": {
        "label": "Liu 2024 Pitch High Amp",
        "motion": "Pitch 0.87Hz, theta=3deg",
        "tank": "1.0x0.5x1.0m rectangular",
        "dp": 0.01, "particles": 668856,
        "data_dir": "data",
        "text_overlay": "Liu 2024 Pitch | 0.87Hz, theta=3deg (High Amp)",
    },
    "frosina_fueltank": {
        "label": "Frosina 2018 Fuel Tank",
        "motion": "Sway 0.5Hz, A=0.05m",
        "tank": "STL automotive fuel tank",
        "dp": 0.01, "particles": 145000,
        "data_dir": "data",
        "text_overlay": "Frosina 2018 Fuel Tank | Sway 0.5Hz",
    },
}

SIM_BASE = Path(__file__).parent.parent / "simulations"
PVPYTHON = "/opt/paraview/bin/pvpython"
DOCKER_IMAGE = "dsph-agent:latest"


def run_isosurface(case_name, info):
    """Run DualSPHysics IsoSurface to generate VTK surface meshes from bi4."""
    case_dir = SIM_BASE / case_name
    iso_dir = case_dir / "iso"
    data_dir = case_dir / info["data_dir"]

    if not data_dir.exists():
        print(f"    No data directory: {data_dir}")
        return 0

    if iso_dir.exists():
        shutil.rmtree(str(iso_dir), ignore_errors=True)
    iso_dir.mkdir(parents=True, exist_ok=True)

    cmd = [
        "docker", "run", "--rm",
        "-v", f"{case_dir}:/data",
        DOCKER_IMAGE,
        "isosurface",
        "-dirin", f"/data/{info['data_dir']}",
        "-saveiso", "/data/iso/Surface",
        "-vars:+vel,+press,+rhop",
        "-onlytype:+fluid",
    ]

    subprocess.run(cmd, capture_output=True, text=True, timeout=600)
    iso_files = sorted(iso_dir.glob("Surface_*.vtk"))
    return len(iso_files)


def pvpython_exec(script, timeout=300):
    """Execute a pvpython script with Mesa offscreen rendering."""
    env = os.environ.copy()
    env["DISPLAY"] = ""
    env["VTK_DEFAULT_RENDER_WINDOW_OFFSCREEN"] = "1"

    result = subprocess.run(
        [PVPYTHON, "--force-offscreen-rendering", "--mesa", "-c", script],
        stdout=subprocess.PIPE, stderr=subprocess.PIPE,
        timeout=timeout, env=env,
    )
    return result.stdout.decode("utf-8", errors="replace")


def render_snapshot(iso_file, field, colormap, output_path):
    """Render a single isosurface snapshot."""
    labels = {"Press": "Pressure [Pa]", "Vel": "Velocity [m/s]", "Rhop": "Density [kg/m3]"}
    field_label = labels.get(field, field)

    script = f"""
import json, os
from pathlib import Path
from paraview.simple import *
paraview.simple._DisableFirstRenderCameraReset()

r = LegacyVTKReader(FileNames=['{iso_file}'])
r.UpdatePipeline()

rv = GetActiveViewOrCreate('RenderView')
rv.ViewSize = [1920, 1080]
rv.UseColorPaletteForBackground = 0
rv.Background = [0.12, 0.12, 0.18]
rv.OrientationAxesVisibility = 0

d = Show(r, rv)
d.Representation = 'Surface'
d.Specular = 0.25
d.SpecularPower = 80

ColorBy(d, ('POINTS', '{field}'))
d.RescaleTransferFunctionToDataRange(True)
lut = GetColorTransferFunction('{field}')
lut.ApplyPreset('{colormap}', True)

d.SetScalarBarVisibility(rv, True)
sb = GetScalarBar(lut, rv)
sb.TitleColor = [1, 1, 1]
sb.LabelColor = [1, 1, 1]
sb.Title = '{field_label}'
sb.Orientation = 'Vertical'
sb.WindowLocation = 'Any Location'
sb.Position = [0.88, 0.12]
sb.ScalarBarLength = 0.65
sb.TitleFontSize = 22
sb.LabelFontSize = 16

rv.CameraPosition = [1, -1, 0.8]
rv.CameraFocalPoint = [0, 0, 0]
rv.CameraViewUp = [0, 0, 1]
rv.ResetCamera()
cam = GetActiveCamera()
cam.Dolly(1.5)
Render()
SaveScreenshot('{output_path}', rv, ImageResolution=[1920, 1080])
"""
    pvpython_exec(script)
    return Path(output_path).exists()


def render_animation_frames(iso_files, field, colormap, frames_dir, fps=15):
    """Render animation frames — loads each IsoSurface VTK individually per frame."""
    labels = {"Press": "Pressure [Pa]", "Vel": "Velocity [m/s]"}
    field_label = labels.get(field, field)
    file_list = json.dumps(iso_files)

    script = f"""
import json, os, sys
from pathlib import Path
from paraview.simple import *
paraview.simple._DisableFirstRenderCameraReset()

frames_dir = '{frames_dir}'
os.makedirs(frames_dir, exist_ok=True)
iso_files = {file_list}

# First pass: scan sample files to find global field range for consistent coloring
global_min = float('inf')
global_max = float('-inf')
_sample_indices = [0, len(iso_files)//4, len(iso_files)//2, 3*len(iso_files)//4, len(iso_files)-1]
for _si in _sample_indices:
    tmp = LegacyVTKReader(FileNames=[iso_files[_si]])
    tmp.UpdatePipeline()
    pa = tmp.GetDataInformation().GetPointDataInformation()
    for ai in range(pa.GetNumberOfArrays()):
        arr = pa.GetArrayInformation(ai)
        if arr.GetName() == '{field}':
            nc = arr.GetNumberOfComponents()
            for c in range(nc):
                lo, hi = arr.GetComponentRange(c)
                if abs(lo) > global_max: global_max = abs(lo)
                if abs(hi) > global_max: global_max = abs(hi)
    Delete(tmp)
if global_min == float('inf') or global_max == float('-inf'):
    global_min, global_max = 0, 1
else:
    global_min = 0  # magnitude always starts at 0
print(f'GLOBAL_RANGE:{{global_min}}:{{global_max}}')

# Setup view once
rv = GetActiveViewOrCreate('RenderView')
rv.ViewSize = [1920, 1080]
rv.UseColorPaletteForBackground = 0
rv.Background = [0.12, 0.12, 0.18]
rv.OrientationAxesVisibility = 0

# Camera setup from first file
r = LegacyVTKReader(FileNames=[iso_files[len(iso_files)//2]])
r.UpdatePipeline()
d = Show(r, rv)
d.Representation = 'Surface'
d.Specular = 0.25
d.SpecularPower = 80
ColorBy(d, ('POINTS', '{field}'))
lut = GetColorTransferFunction('{field}')
lut.ApplyPreset('{colormap}', True)
lut.RescaleTransferFunction(global_min, global_max)
d.SetScalarBarVisibility(rv, True)
sb = GetScalarBar(lut, rv)
sb.TitleColor = [1, 1, 1]
sb.LabelColor = [1, 1, 1]
sb.Title = '{field_label}'
sb.Orientation = 'Vertical'
sb.WindowLocation = 'Any Location'
sb.Position = [0.88, 0.12]
sb.ScalarBarLength = 0.65
sb.TitleFontSize = 22
sb.LabelFontSize = 16

rv.CameraPosition = [1, -1, 0.8]
rv.CameraFocalPoint = [0, 0, 0]
rv.CameraViewUp = [0, 0, 1]
rv.ResetCamera()
cam = GetActiveCamera()
cam.Dolly(1.5)
Render()

# Save camera state
cam_pos = list(rv.CameraPosition)
cam_fp = list(rv.CameraFocalPoint)
cam_up = list(rv.CameraViewUp)
cam_ps = rv.CameraParallelScale if hasattr(rv, 'CameraParallelScale') else None

Delete(d)
Delete(r)

# Render each frame by loading individual files
for fi, fpath_vtk in enumerate(iso_files):
    r = LegacyVTKReader(FileNames=[fpath_vtk])
    r.UpdatePipeline()
    d = Show(r, rv)
    d.Representation = 'Surface'
    d.Specular = 0.25
    d.SpecularPower = 80
    ColorBy(d, ('POINTS', '{field}'))
    lut = GetColorTransferFunction('{field}')
    lut.ApplyPreset('{colormap}', True)
    lut.RescaleTransferFunction(global_min, global_max)
    d.SetScalarBarVisibility(rv, True)

    rv.CameraPosition = cam_pos
    rv.CameraFocalPoint = cam_fp
    rv.CameraViewUp = cam_up
    if cam_ps is not None:
        rv.CameraParallelScale = cam_ps

    Render()
    fpath = os.path.join(frames_dir, f'frame_{{fi:06d}}.png')
    SaveScreenshot(fpath, rv, ImageResolution=[1920, 1080])

    Hide(r, rv)
    Delete(d)
    Delete(r)

    if fi % 10 == 0:
        print(f'FRAME:{{fi}}/{{len(iso_files)}}')

print(f'DONE:frames')
"""
    timeout = max(300, len(iso_files) * 8 + 120)
    pvpython_exec(script, timeout=timeout)
    return len(list(Path(frames_dir).glob("frame_*.png")))


def compile_ffmpeg(frames_dir, output_mp4, fps, text):
    """Compile PNG frames to MP4 with text overlay."""
    if not os.path.isdir(frames_dir) or not os.listdir(frames_dir):
        return False

    text_escaped = text.replace(":", "\\:").replace("'", "\\'")

    cmd = [
        "ffmpeg", "-y",
        "-framerate", str(fps),
        "-i", os.path.join(frames_dir, "frame_%06d.png"),
        "-vf", (
            f"drawtext=text='{text_escaped}'"
            ":fontfile=/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"
            ":fontsize=28:fontcolor=white:borderw=2:bordercolor=black"
            ":x=30:y=30"
        ),
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "-preset", "fast",
        "-crf", "23",
        output_mp4,
    ]
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=120)
    return result.returncode == 0


def extract_stats(vtk_file):
    """Extract field statistics from a VTK file."""
    script = f"""
import json
from paraview.simple import *
r = LegacyVTKReader(FileNames=['{vtk_file}'])
r.UpdatePipeline()
info = r.GetDataInformation()
pa = info.GetPointDataInformation()
stats = {{'points': info.GetNumberOfPoints(), 'cells': info.GetNumberOfCells()}}
for i in range(pa.GetNumberOfArrays()):
    arr = pa.GetArrayInformation(i)
    name = arr.GetName()
    nc = arr.GetNumberOfComponents()
    ranges = []
    for c in range(nc):
        lo, hi = arr.GetComponentRange(c)
        ranges.append({{'min': lo, 'max': hi}})
    stats[name] = {{'components': nc, 'ranges': ranges}}
print('STATS_JSON:' + json.dumps(stats))
"""
    stdout = pvpython_exec(script, timeout=60)
    for line in stdout.split("\n"):
        if line.startswith("STATS_JSON:"):
            return json.loads(line[len("STATS_JSON:"):])
    return {}


def process_case(case_name, info):
    """Full pipeline: isosurface -> render -> animate -> stats."""
    case_dir = SIM_BASE / case_name
    viz_dir = case_dir / "viz"
    iso_dir = case_dir / "iso"
    viz_dir.mkdir(parents=True, exist_ok=True)

    results = {"case": case_name, "label": info["label"], "outputs": []}

    # 1. IsoSurface generation
    iso_files = sorted(iso_dir.glob("Surface_*.vtk")) if iso_dir.exists() else []
    if len(iso_files) < 2:
        print(f"  [{case_name}] Generating isosurface meshes...")
        n_iso = run_isosurface(case_name, info)
        iso_files = sorted(iso_dir.glob("Surface_*.vtk"))
        print(f"    Generated {n_iso} surface meshes")
    else:
        print(f"  [{case_name}] IsoSurface cached ({len(iso_files)} files)")

    if not iso_files:
        return {"case": case_name, "error": "no isosurface files"}

    iso_paths = [str(f) for f in iso_files]
    snap_idx = min(int(len(iso_files) * 0.6), len(iso_files) - 1)
    snap_file = str(iso_files[snap_idx])

    # 2. Snapshots
    print(f"  [{case_name}] Rendering snapshots...")
    vel_out = str(viz_dir / "snapshot_vel.png")
    if render_snapshot(snap_file, "Vel", "Viridis", vel_out):
        sz = Path(vel_out).stat().st_size // 1024
        results["outputs"].append({"type": "snapshot", "field": "Vel", "path": vel_out})
        print(f"    Vel snapshot ({sz}KB)")
    else:
        print(f"    Vel snapshot FAILED")

    press_out = str(viz_dir / "snapshot_press.png")
    if render_snapshot(snap_file, "Press", "Cool to Warm", press_out):
        sz = Path(press_out).stat().st_size // 1024
        results["outputs"].append({"type": "snapshot", "field": "Press", "path": press_out})
        print(f"    Press snapshot ({sz}KB)")
    else:
        print(f"    Press snapshot FAILED")

    # 3. Animation
    print(f"  [{case_name}] Generating animation ({len(iso_files)} timesteps)...")
    frames_dir = str(viz_dir / "frames")
    n_frames = render_animation_frames(iso_paths, "Vel", "Viridis", frames_dir)
    print(f"    Frames: {n_frames}")

    if n_frames > 0:
        mp4_path = str(viz_dir / f"{case_name}_sloshing.mp4")
        dur = n_frames / 15
        print(f"    Compiling MP4 (15fps, {dur:.1f}s)...")
        ok = compile_ffmpeg(frames_dir, mp4_path, 15, info["text_overlay"])
        if ok and Path(mp4_path).exists():
            size_mb = Path(mp4_path).stat().st_size / (1024 * 1024)
            results["outputs"].append({
                "type": "animation", "path": mp4_path,
                "frames": n_frames, "size_mb": round(size_mb, 1),
            })
            print(f"    Animation: {size_mb:.1f}MB, {n_frames} frames")
        else:
            print(f"    ffmpeg FAILED")

    shutil.rmtree(frames_dir, ignore_errors=True)

    # 4. Statistics
    print(f"  [{case_name}] Statistics...")
    stats = extract_stats(snap_file)
    enriched = {
        "case": case_name, "label": info["label"],
        "motion": info["motion"], "tank": info["tank"],
        "dp": info["dp"], "particles": info["particles"],
        "isosurface_files": len(iso_files),
        "statistics": stats,
    }
    stats_path = viz_dir / "field_stats.json"
    stats_path.write_text(json.dumps(enriched, indent=2, ensure_ascii=False))
    results["outputs"].append({"type": "stats", "path": str(stats_path)})
    print(f"    Stats saved")

    return results


def main():
    print("=" * 70)
    print("slosim-agent Post-Processing (IsoSurface + pvpython)")
    print(f"Cases: {len(CASES)}")
    print("=" * 70)

    skip_cases = set()
    for cn in CASES:
        viz_dir = SIM_BASE / cn / "viz"
        if viz_dir.exists():
            mp4s = list(viz_dir.glob("*.mp4"))
            pngs = list(viz_dir.glob("*.png"))
            if len(mp4s) >= 1 and len(pngs) >= 2:
                skip_cases.add(cn)
                print(f"  SKIP: {cn}")

    all_results = []
    for i, (cn, info) in enumerate(CASES.items(), 1):
        if cn in skip_cases:
            all_results.append({"case": cn, "skipped": True})
            continue

        print(f"\n{'=' * 50}")
        print(f"[{i}/{len(CASES)}] {cn}: {info['label']}")
        print(f"  {info['tank']} | {info['particles']:,} particles")
        print(f"{'=' * 50}")

        try:
            result = process_case(cn, info)
            all_results.append(result)
            n = len(result.get("outputs", []))
            print(f"\n  [{cn}] DONE ({n} outputs)")
        except Exception as e:
            import traceback
            print(f"\n  [{cn}] ERROR: {e}")
            traceback.print_exc()
            all_results.append({"case": cn, "error": str(e)})

    summary_path = SIM_BASE / "postprocess_summary.json"
    summary_path.write_text(json.dumps(all_results, indent=2, ensure_ascii=False))

    print(f"\n{'=' * 70}")
    print("COMPLETE")
    print(f"{'=' * 70}")
    total = 0
    for r in all_results:
        n = len(r.get("outputs", []))
        total += n
        s = "SKIP" if r.get("skipped") else (f"ERR: {r['error']}" if "error" in r else f"{n} out")
        print(f"  {r['case']}: {s}")
    print(f"\nTotal: {total} outputs / {len(CASES)} cases")


if __name__ == "__main__":
    main()
