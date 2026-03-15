#!/usr/bin/env python3
"""EXP-C: Oil lateral sloshing pressure animation.
Renders particle snapshots colored by pressure, stitches into MP4.
"""
import subprocess
import sys
from pathlib import Path
import numpy as np

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from matplotlib.colors import Normalize
import matplotlib.cm as cm

try:
    import pyvista as pv
    pv.OFF_SCREEN = True
except ImportError:
    print("ERROR: pyvista required")
    sys.exit(1)

VTK_DIR = Path("/tmp/vtk_render/expc_anim")
FRAME_DIR = Path("/tmp/vtk_render/frames_expc")
FRAME_DIR.mkdir(exist_ok=True)
OUT_DIR = Path(__file__).parent / "videos"
OUT_DIR.mkdir(exist_ok=True)

DT = 0.01  # 10ms per PART
STEP = 5   # every 5th frame
N_SAMPLE = 6000
PMAX = 5000  # Pa


def load_vtk(filepath):
    mesh = pv.read(str(filepath))
    pos = np.array(mesh.points)
    vel = None
    press = None
    if 'Vel' in mesh.point_data:
        vel = np.array(mesh.point_data['Vel'])
    if 'Press' in mesh.point_data:
        press = np.array(mesh.point_data['Press'])
    return pos, vel, press


def render_frame(vtk_path, frame_idx, t_sec):
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 5),
                                     gridspec_kw={'width_ratios': [1, 1]})

    pos, vel, press = load_vtk(vtk_path)

    # Subsample
    if len(pos) > N_SAMPLE:
        idx = np.random.choice(len(pos), N_SAMPLE, replace=False)
        pos = pos[idx]
        vel = vel[idx] if vel is not None else None
        press = press[idx] if press is not None else None

    # Left panel: velocity
    norm_v = Normalize(0, 1.5)
    if vel is not None:
        vmag = np.linalg.norm(vel, axis=1)
        ax1.scatter(pos[:, 0], pos[:, 2], c=vmag, cmap='coolwarm',
                   norm=norm_v, s=0.8, alpha=0.6, rasterized=True)
    else:
        ax1.scatter(pos[:, 0], pos[:, 2], c='#2196F3', s=0.8, alpha=0.5, rasterized=True)

    ax1.set_title('Velocity', fontsize=12, fontweight='bold')
    ax1.set_xlabel('X (m)', fontsize=9)
    ax1.set_ylabel('Z (m)', fontsize=9)

    # Right panel: pressure
    norm_p = Normalize(0, PMAX)
    if press is not None:
        ax2.scatter(pos[:, 0], pos[:, 2], c=press, cmap='RdYlBu_r',
                   norm=norm_p, s=0.8, alpha=0.6, rasterized=True)
    else:
        ax2.scatter(pos[:, 0], pos[:, 2], c='#2196F3', s=0.8, alpha=0.5, rasterized=True)

    ax2.set_title('Pressure', fontsize=12, fontweight='bold')
    ax2.set_xlabel('X (m)', fontsize=9)

    for ax in [ax1, ax2]:
        ax.set_aspect('equal')
        ax.grid(True, alpha=0.15)
        ax.tick_params(labelsize=7)

    # Colorbars
    sm_v = cm.ScalarMappable(cmap='coolwarm', norm=norm_v)
    sm_v.set_array([])
    cb1 = fig.colorbar(sm_v, ax=ax1, shrink=0.7, pad=0.02)
    cb1.set_label('|V| (m/s)', fontsize=8)

    sm_p = cm.ScalarMappable(cmap='RdYlBu_r', norm=norm_p)
    sm_p.set_array([])
    cb2 = fig.colorbar(sm_p, ax=ax2, shrink=0.7, pad=0.02)
    cb2.set_label('P (Pa)', fontsize=8)

    # Peak indicators
    peak_times = [1.54, 3.07, 4.61, 6.14]
    peak_labels = ["P1", "P2", "P3", "P4"]
    for pt, pl in zip(peak_times, peak_labels):
        if abs(t_sec - pt) < 0.03:
            fig.text(0.5, 0.02, f'>>> IMPACT PEAK {pl} <<<',
                    ha='center', fontsize=14, fontweight='bold', color='red')

    fig.suptitle(f'EXP-C: SPHERIC Oil Lateral (Laminar+SPS, 10ms)  |  t = {t_sec:.2f} s',
                 fontsize=13, fontweight='bold')
    fig.tight_layout(rect=[0, 0.04, 1, 0.95])

    outpath = FRAME_DIR / f"frame_{frame_idx:04d}.png"
    fig.savefig(outpath, dpi=150, bbox_inches='tight',
                facecolor='white', edgecolor='none')
    plt.close(fig)
    return outpath


def main():
    print("=" * 50)
    print("  EXP-C Animation Renderer")
    print(f"  Output: {OUT_DIR}")
    print("=" * 50)

    np.random.seed(42)

    files = sorted(VTK_DIR.glob("Fluid_????.vtk"))
    print(f"  {len(files)} frames to render")

    for i, f in enumerate(files):
        part_num = int(f.stem.split('_')[1])
        t_sec = part_num * DT
        render_frame(f, i, t_sec)
        if (i + 1) % 20 == 0 or i == 0:
            print(f"  Rendered {i+1}/{len(files)} (t={t_sec:.2f}s)")

    print(f"\n  All {len(files)} frames rendered")

    # MP4
    mp4_path = OUT_DIR / "expc_oil_lateral.mp4"
    cmd = [
        "ffmpeg", "-y", "-framerate", "12",
        "-i", str(FRAME_DIR / "frame_%04d.png"),
        "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
        "-c:v", "libx264", "-pix_fmt", "yuv420p",
        "-crf", "20", "-preset", "medium",
        str(mp4_path)
    ]
    print("  Encoding MP4...")
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode == 0:
        size_mb = mp4_path.stat().st_size / (1024 * 1024)
        duration = len(files) / 12
        print(f"  OK: {mp4_path.name} ({size_mb:.1f} MB, {duration:.1f}s)")
    else:
        print(f"  ERROR: {result.stderr[-300:]}")

    # GIF (lower quality for embedding)
    gif_path = OUT_DIR / "expc_oil_lateral.gif"
    cmd_gif = [
        "ffmpeg", "-y", "-framerate", "12",
        "-i", str(FRAME_DIR / "frame_%04d.png"),
        "-vf", "fps=8,scale=800:-1:flags=lanczos",
        str(gif_path)
    ]
    print("  Encoding GIF...")
    result = subprocess.run(cmd_gif, capture_output=True, text=True)
    if result.returncode == 0:
        size_mb = gif_path.stat().st_size / (1024 * 1024)
        print(f"  OK: {gif_path.name} ({size_mb:.1f} MB)")

    print("\nDone!")


if __name__ == "__main__":
    main()
