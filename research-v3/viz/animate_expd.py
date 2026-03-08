#!/usr/bin/env python3
"""EXP-D: Side-by-side Baseline vs Baffle animation.
Renders each frame pair with pyvista+matplotlib, stitches with ffmpeg.
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
    print("ERROR: pyvista required. pip install pyvista")
    sys.exit(1)

VTK_BL = Path("/tmp/vtk_render/expd_baseline")
VTK_BF = Path("/tmp/vtk_render/expd_baffle")
FRAME_DIR = Path("/tmp/vtk_render/frames_expd")
FRAME_DIR.mkdir(exist_ok=True)
OUT_DIR = Path(__file__).parent / "videos"
OUT_DIR.mkdir(exist_ok=True)

DT = 0.1  # seconds per PART file (timemax=8, 80 parts)
N_SAMPLE = 5000
VMAX = 1.5


def load_vtk(filepath):
    mesh = pv.read(str(filepath))
    pos = np.array(mesh.points)
    vel = None
    if 'Vel' in mesh.point_data:
        vel = np.array(mesh.point_data['Vel'])
    return pos, vel


def render_frame(bl_path, bf_path, frame_idx, t_sec):
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 5))

    norm = Normalize(0, VMAX)

    for ax, path, label in [(ax1, bl_path, "Baseline"), (ax2, bf_path, "Baffle")]:
        if path is not None and path.exists():
            pos, vel = load_vtk(path)
            # Subsample
            if len(pos) > N_SAMPLE:
                idx = np.random.choice(len(pos), N_SAMPLE, replace=False)
                pos = pos[idx]
                vel = vel[idx] if vel is not None else None

            if vel is not None:
                vmag = np.linalg.norm(vel, axis=1)
                ax.scatter(pos[:, 0], pos[:, 2], c=vmag, cmap='coolwarm',
                          norm=norm, s=1.0, alpha=0.6, rasterized=True)
            else:
                ax.scatter(pos[:, 0], pos[:, 2], c='#2196F3', s=1.0, alpha=0.5, rasterized=True)
        else:
            ax.text(0.5, 0.5, "N/A", ha='center', va='center', transform=ax.transAxes, fontsize=20)

        ax.set_xlim(-0.02, 0.62)
        ax.set_ylim(-0.02, 0.42)
        ax.set_aspect('equal')
        ax.set_title(label, fontsize=13, fontweight='bold')
        ax.set_xlabel('X (m)', fontsize=10)
        ax.set_ylabel('Z (m)', fontsize=10)
        ax.grid(True, alpha=0.15)
        # Tank outline
        ax.plot([0, 0.6, 0.6, 0, 0], [0, 0, 0.4, 0.4, 0],
                'k-', linewidth=1.2, alpha=0.4)

    # Baffle indicator on right panel
    ax2.plot([0.3, 0.3], [0, 0.25], 'r-', linewidth=2.5, alpha=0.7)

    # Colorbar
    sm = cm.ScalarMappable(cmap='coolwarm', norm=norm)
    sm.set_array([])
    cbar = fig.colorbar(sm, ax=[ax1, ax2], shrink=0.8, pad=0.02, aspect=30)
    cbar.set_label('|V| (m/s)', fontsize=10)

    fig.suptitle(f'EXP-D: Baseline vs Center Baffle  |  t = {t_sec:.1f} s',
                 fontsize=14, fontweight='bold')
    fig.tight_layout(rect=[0, 0, 0.92, 0.95])

    outpath = FRAME_DIR / f"frame_{frame_idx:04d}.png"
    fig.savefig(outpath, dpi=150, bbox_inches='tight',
                facecolor='white', edgecolor='none')
    plt.close(fig)
    return outpath


def main():
    print("=" * 50)
    print("  EXP-D Animation Renderer")
    print(f"  Output: {OUT_DIR}")
    print("=" * 50)

    np.random.seed(42)

    # Find common frame range
    bl_files = sorted(VTK_BL.glob("Fluid_????.vtk"))
    bf_files = sorted(VTK_BF.glob("Fluid_????.vtk"))

    bl_nums = {int(f.stem.split('_')[1]) for f in bl_files}
    bf_nums = {int(f.stem.split('_')[1]) for f in bf_files}
    common = sorted(bl_nums & bf_nums)

    print(f"  Baseline frames: {len(bl_files)}")
    print(f"  Baffle frames:   {len(bf_files)}")
    print(f"  Common frames:   {len(common)}")
    print(f"  Range: PART_{common[0]:04d} - PART_{common[-1]:04d}")

    # Render each frame
    for i, part_num in enumerate(common):
        t_sec = part_num * DT
        bl_path = VTK_BL / f"Fluid_{part_num:04d}.vtk"
        bf_path = VTK_BF / f"Fluid_{part_num:04d}.vtk"
        render_frame(bl_path, bf_path, i, t_sec)
        if (i + 1) % 10 == 0 or i == 0:
            print(f"  Rendered {i+1}/{len(common)} (t={t_sec:.1f}s)")

    print(f"\n  All {len(common)} frames rendered")

    # Stitch with ffmpeg
    mp4_path = OUT_DIR / "expd_baseline_vs_baffle.mp4"
    cmd = [
        "ffmpeg", "-y",
        "-framerate", "10",
        "-i", str(FRAME_DIR / "frame_%04d.png"),
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "-crf", "20",
        "-preset", "medium",
        str(mp4_path)
    ]
    print(f"\n  Encoding MP4...")
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode == 0:
        size_mb = mp4_path.stat().st_size / (1024 * 1024)
        print(f"  OK: {mp4_path} ({size_mb:.1f} MB)")
    else:
        print(f"  ERROR: ffmpeg failed\n{result.stderr[-500:]}")

    # Also create GIF (smaller, for embedding)
    gif_path = OUT_DIR / "expd_baseline_vs_baffle.gif"
    cmd_gif = [
        "ffmpeg", "-y",
        "-framerate", "10",
        "-i", str(FRAME_DIR / "frame_%04d.png"),
        "-vf", "fps=8,scale=800:-1:flags=lanczos",
        str(gif_path)
    ]
    print(f"  Encoding GIF...")
    result = subprocess.run(cmd_gif, capture_output=True, text=True)
    if result.returncode == 0:
        size_mb = gif_path.stat().st_size / (1024 * 1024)
        print(f"  OK: {gif_path} ({size_mb:.1f} MB)")

    print("\nDone!")


if __name__ == "__main__":
    main()
