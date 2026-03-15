#!/usr/bin/env python3
"""Render SPH particle snapshots using pyvista + matplotlib.
Reads binary VTK files with pyvista, renders 2D projections with matplotlib.
"""
import sys
from pathlib import Path
import numpy as np

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from matplotlib.colors import Normalize

try:
    import pyvista as pv
    pv.OFF_SCREEN = True
except ImportError:
    print("ERROR: pyvista required. pip install pyvista")
    sys.exit(1)

VTK_BASE = Path("/tmp/vtk_render")
OUT_DIR = Path(__file__).parent / "paraview"
OUT_DIR.mkdir(exist_ok=True)


def load_vtk(filepath):
    """Load VTK file and return positions, velocity, pressure."""
    mesh = pv.read(str(filepath))
    pos = np.array(mesh.points)

    vel = None
    if 'Vel' in mesh.point_data:
        vel = np.array(mesh.point_data['Vel'])
    elif 'Velocity' in mesh.point_data:
        vel = np.array(mesh.point_data['Velocity'])

    press = None
    if 'Press' in mesh.point_data:
        press = np.array(mesh.point_data['Press'])
    elif 'Pressure' in mesh.point_data:
        press = np.array(mesh.point_data['Pressure'])

    return pos, vel, press


def render_side(pos, color_data, output_name, title="", cbar_label="",
                vmin=0, vmax=1, cmap='coolwarm', n_sample=6000):
    """Render X-Z side view."""
    if len(pos) == 0:
        return

    # Subsample
    if len(pos) > n_sample:
        idx = np.random.choice(len(pos), n_sample, replace=False)
        pos = pos[idx]
        if color_data is not None:
            color_data = color_data[idx]

    fig, ax = plt.subplots(figsize=(10, 6))

    if color_data is not None:
        sc = ax.scatter(pos[:, 0], pos[:, 2], c=color_data, cmap=cmap,
                       vmin=vmin, vmax=vmax, s=1.5, alpha=0.65, rasterized=True)
        cbar = fig.colorbar(sc, ax=ax, shrink=0.8, pad=0.02)
        cbar.set_label(cbar_label, fontsize=10)
    else:
        ax.scatter(pos[:, 0], pos[:, 2], c='#2196F3', s=1.5, alpha=0.5, rasterized=True)

    ax.set_xlabel('X (m)', fontsize=11)
    ax.set_ylabel('Z (m)', fontsize=11)
    ax.set_title(title, fontsize=13, fontweight='bold')
    ax.set_aspect('equal')
    ax.grid(True, alpha=0.2)

    fig.tight_layout()
    fig.savefig(OUT_DIR / f"{output_name}.png", dpi=200, bbox_inches='tight')
    fig.savefig(OUT_DIR / f"{output_name}.pdf", bbox_inches='tight')
    plt.close(fig)
    print(f"  OK: {output_name}")


def render_expd_panel():
    """EXP-D: 2×4 comparison panel (Baseline vs Baffle × 4 timesteps)."""
    fig, axes = plt.subplots(2, 4, figsize=(16, 6), sharex=True, sharey=True)

    for col, (part, t_label) in enumerate([("0030", "t=3s"), ("0040", "t=4s"),
                                            ("0050", "t=5s"), ("0060", "t=6s")]):
        for row, (case, label) in enumerate([("expd_baseline", "Baseline"),
                                              ("expd_baffle", "Baffle")]):
            ax = axes[row, col]
            path = VTK_BASE / case / f"Fluid_{part}.vtk"
            if not path.exists():
                ax.text(0.5, 0.5, "N/A", ha='center', transform=ax.transAxes)
                continue

            pos, vel, _ = load_vtk(path)
            n = min(4000, len(pos))
            idx = np.random.choice(len(pos), n, replace=False)

            vmag = np.linalg.norm(vel[idx], axis=1) if vel is not None else np.zeros(n)
            ax.scatter(pos[idx, 0], pos[idx, 2], c=vmag, cmap='coolwarm',
                      vmin=0, vmax=1.5, s=0.6, alpha=0.55, rasterized=True)

            ax.set_aspect('equal')
            if row == 0:
                ax.set_title(t_label, fontsize=10, fontweight='bold')
            if col == 0:
                ax.set_ylabel(f'{label}\nZ (m)', fontsize=9)
            if row == 1:
                ax.set_xlabel('X (m)', fontsize=8)
            ax.tick_params(labelsize=7)
            ax.grid(True, alpha=0.15)

    # Add colorbar
    norm = Normalize(0, 1.5)
    sm = plt.cm.ScalarMappable(cmap='coolwarm', norm=norm)
    sm.set_array([])
    cbar = fig.colorbar(sm, ax=axes, shrink=0.6, pad=0.02, label='|V| (m/s)')

    fig.suptitle('EXP-D: Baseline vs Center Baffle — Velocity Field', fontsize=14, fontweight='bold')
    fig.tight_layout(rect=[0, 0, 0.92, 0.95])
    fig.savefig(OUT_DIR / "expd_comparison_panel.png", dpi=200, bbox_inches='tight')
    fig.savefig(OUT_DIR / "expd_comparison_panel.pdf", bbox_inches='tight')
    plt.close(fig)
    print("  OK: expd_comparison_panel")


def render_expc_panel():
    """EXP-C: 1×4 oil lateral pressure at peak moments."""
    fig, axes = plt.subplots(1, 4, figsize=(16, 4.5))

    for col, (part, t_sec, label, p_exp) in enumerate([
        ("0154", "1.54", "Peak 1", 2200), ("0307", "3.07", "Peak 2", 3800),
        ("0461", "4.61", "Peak 3", 5500), ("0614", "6.14", "Peak 4", 6200)
    ]):
        ax = axes[col]
        path = VTK_BASE / "expc_009b" / f"Fluid_{part}.vtk"
        if not path.exists():
            ax.text(0.5, 0.5, "N/A", ha='center', transform=ax.transAxes)
            continue

        pos, _, press = load_vtk(path)
        n = min(6000, len(pos))
        idx = np.random.choice(len(pos), n, replace=False)

        p_data = press[idx] if press is not None else np.zeros(n)
        sc = ax.scatter(pos[idx, 0], pos[idx, 2], c=p_data, cmap='RdYlBu_r',
                       vmin=0, vmax=6000, s=0.5, alpha=0.55, rasterized=True)
        ax.set_aspect('equal')
        ax.set_title(f'{label}\nt={t_sec}s (exp: {p_exp} Pa)', fontsize=9, fontweight='bold')
        ax.set_xlabel('X (m)', fontsize=8)
        if col == 0:
            ax.set_ylabel('Z (m)', fontsize=8)
        ax.tick_params(labelsize=7)

    cbar = fig.colorbar(sc, ax=axes, shrink=0.8, pad=0.02, label='Pressure (Pa)')
    fig.suptitle('EXP-C: Oil Lateral — Pressure at Impact Peaks (Laminar+SPS, 10ms)',
                 fontsize=13, fontweight='bold')
    fig.tight_layout()
    fig.savefig(OUT_DIR / "expc_peaks_panel.png", dpi=200, bbox_inches='tight')
    fig.savefig(OUT_DIR / "expc_peaks_panel.pdf", bbox_inches='tight')
    plt.close(fig)
    print("  OK: expc_peaks_panel")


def main():
    print("=" * 50)
    print("  SPH Particle Rendering (pyvista + matplotlib)")
    print(f"  Output: {OUT_DIR}")
    print("=" * 50)
    np.random.seed(42)

    # EXP-D panel
    print("\n--- EXP-D ---")
    render_expd_panel()

    # EXP-D individual high-res
    for part, t_sec in [("0040", "4.0"), ("0050", "5.0"), ("0060", "6.0")]:
        for case, label in [("expd_baseline", "Baseline"), ("expd_baffle", "Baffle")]:
            path = VTK_BASE / case / f"Fluid_{part}.vtk"
            if path.exists():
                pos, vel, _ = load_vtk(path)
                vmag = np.linalg.norm(vel, axis=1) if vel is not None else None
                render_side(pos, vmag, f"expd_{label.lower()}_t{t_sec}s",
                           f"EXP-D {label} t={t_sec}s", "|V| (m/s)", 0, 1.5)

    # EXP-C panel
    print("\n--- EXP-C ---")
    render_expc_panel()

    # EXP-C individual
    for part, t_sec, label in [("0307", "3.07", "Peak2"), ("0461", "4.61", "Peak3"),
                                ("0614", "6.14", "Peak4")]:
        path = VTK_BASE / "expc_009b" / f"Fluid_{part}.vtk"
        if path.exists():
            pos, _, press = load_vtk(path)
            render_side(pos, press, f"expc_oil_{label}_t{t_sec}s",
                       f"Oil Lateral {label} t={t_sec}s", "Pressure (Pa)",
                       0, 6000, 'RdYlBu_r', n_sample=10000)

    total = len(list(OUT_DIR.glob("*.png")))
    print(f"\nDone! {total} images rendered")


if __name__ == "__main__":
    main()
