#!/usr/bin/env python3
"""
analyze_chen2018.py — Analyze Chen2018 parametric study results

Reads MeasureTool elevation data from 6 fill level simulations and generates:
  - Fig 3: 6-panel SWL time series (one per fill level)
  - Fig 3b: SWL amplitude vs fill level
  - Table 4: Summary metrics (max wave height, dominant frequency)
"""

import numpy as np
import pandas as pd
from pathlib import Path
from scipy.fft import fft, fftfreq
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import json
import math

OUT_DIR = Path("research/experiments/exp3_parametric")
OUT_DIR.mkdir(parents=True, exist_ok=True)
(OUT_DIR / "figures").mkdir(exist_ok=True)
(OUT_DIR / "comparison").mkdir(exist_ok=True)

# ── Config ────────────────────────────────────────────────────────
FILLS = [120, 150, 195, 260, 325, 390]
TANK_H = 650  # mm
TANK_L = 600  # mm


def natural_freq(h_m, L=0.6):
    """First mode natural frequency [Hz]."""
    omega1 = math.sqrt(9.81 * (math.pi / L) * math.tanh(math.pi * h_m / L))
    return omega1 / (2 * math.pi)


# ── Data Loading ──────────────────────────────────────────────────

def load_elevation(sim_dir, probe_idx=0):
    """Load MeasureTool elevation CSV → (time, elevation_mm).

    probe_idx: 0=Left, 1=Center, 2=Right
    """
    filepath = Path(sim_dir) / "elevation_Elevation.csv"
    if not filepath.exists():
        print(f"  WARNING: {filepath} not found")
        return None, None

    # Format: header lines start with space, data starts at "Part;Time..."
    lines = filepath.read_text().strip().split("\n")

    # Find data start (line starting with digit)
    data_start = 0
    for i, line in enumerate(lines):
        if line.strip() and line.strip()[0].isdigit():
            data_start = i
            break

    # Parse header to get column count
    header_line = lines[data_start - 1]  # "Part;Time [s];Elevation_0 [m];..."

    # Parse data
    data = []
    for line in lines[data_start:]:
        parts = line.strip().split(";")
        if len(parts) >= 3 + probe_idx:
            try:
                t = float(parts[1])
                elev = float(parts[2 + probe_idx])
                data.append((t, elev))
            except (ValueError, IndexError):
                continue

    if not data:
        return None, None

    data = np.array(data)
    t = data[:, 0]
    elev_mm = data[:, 1] * 1000  # m → mm
    return t, elev_mm


def compute_wave_amplitude(t, swl, h_mm, t_start=3.0):
    """Compute wave amplitude (peak-to-trough/2) in steady state."""
    mask = t >= t_start
    if np.sum(mask) < 10:
        return 0.0

    swl_ss = swl[mask]
    swl_centered = swl_ss - np.mean(swl_ss)
    amplitude = (np.max(swl_centered) - np.min(swl_centered)) / 2
    return amplitude


def compute_dominant_freq(t, swl, t_start=3.0):
    """FFT to find dominant frequency."""
    mask = t >= t_start
    if np.sum(mask) < 20:
        return 0.0

    t_ss = t[mask]
    swl_ss = swl[mask] - np.mean(swl[mask])

    dt = np.mean(np.diff(t_ss))
    N = len(swl_ss)
    yf = np.abs(fft(swl_ss))[:N // 2]
    xf = fftfreq(N, dt)[:N // 2]

    valid = xf > 0.1
    if np.any(valid):
        peak_idx = np.argmax(yf[valid])
        return xf[valid][peak_idx]
    return 0.0


# ── Figure 3: SWL Time Series ────────────────────────────────────

def fig_swl_timeseries():
    """6-panel SWL time series for each fill level."""
    fig, axes = plt.subplots(3, 2, figsize=(14, 12))
    axes = axes.flatten()

    for idx, h_mm in enumerate(FILLS):
        ax = axes[idx]
        sim_dir = f"simulations/chen_h{h_mm}"

        t, swl = load_elevation(sim_dir, probe_idx=0)  # Left wall
        if t is None:
            ax.text(0.5, 0.5, "No data", transform=ax.transAxes, ha="center")
            ax.set_title(f"h={h_mm}mm ({h_mm / TANK_H * 100:.0f}%)")
            continue

        fill_pct = h_mm / TANK_H * 100
        f1 = natural_freq(h_mm / 1000)
        f_exc = 0.9 * f1

        ax.plot(t, swl, "b-", linewidth=0.8)
        ax.axhline(y=h_mm, color="gray", linestyle="--", alpha=0.5,
                    label=f"Still level ({h_mm}mm)")
        ax.set_ylabel("SWL [mm]")
        ax.set_title(f"h={h_mm}mm ({fill_pct:.0f}%), f={f_exc:.3f}Hz (0.9f₁)")
        ax.grid(True, alpha=0.3)
        ax.legend(fontsize=8, loc="upper right")

        if idx >= 4:
            ax.set_xlabel("Time [s]")

    fig.suptitle("Chen2018 Parametric — Left Wall SWL at f/f₁=0.9",
                 fontsize=14, fontweight="bold")
    plt.tight_layout()
    path = OUT_DIR / "figures" / "fig3_swl_timeseries.png"
    plt.savefig(path, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {path}")


# ── Figure 3b: Amplitude vs Fill Level ───────────────────────────

def fig_amplitude_vs_fill():
    """Wave amplitude as function of fill level."""
    amplitudes_left = []
    amplitudes_right = []
    freqs = []

    for h_mm in FILLS:
        sim_dir = f"simulations/chen_h{h_mm}"

        t_l, swl_l = load_elevation(sim_dir, probe_idx=0)  # Left
        t_r, swl_r = load_elevation(sim_dir, probe_idx=2)  # Right

        if t_l is not None:
            amplitudes_left.append(compute_wave_amplitude(t_l, swl_l, h_mm))
        else:
            amplitudes_left.append(0)

        if t_r is not None:
            amplitudes_right.append(compute_wave_amplitude(t_r, swl_r, h_mm))
        else:
            amplitudes_right.append(0)

        if t_l is not None:
            freqs.append(compute_dominant_freq(t_l, swl_l))
        else:
            freqs.append(0)

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(12, 5))

    fill_pcts = [h / TANK_H * 100 for h in FILLS]

    ax1.bar(np.arange(len(FILLS)) - 0.15, amplitudes_left, 0.3,
            label="Left wall", color="steelblue", edgecolor="navy")
    ax1.bar(np.arange(len(FILLS)) + 0.15, amplitudes_right, 0.3,
            label="Right wall", color="tomato", edgecolor="darkred")
    ax1.set_xticks(range(len(FILLS)))
    ax1.set_xticklabels([f"{h}mm\n({p:.0f}%)" for h, p in zip(FILLS, fill_pcts)])
    ax1.set_ylabel("Wave Amplitude [mm]")
    ax1.set_title("Sloshing Wave Amplitude vs Fill Level")
    ax1.legend()
    ax1.grid(True, alpha=0.3, axis="y")

    norm_amp = [a / h if h > 0 else 0 for a, h in zip(amplitudes_left, FILLS)]
    ax2.plot(fill_pcts, norm_amp, "bo-", markersize=8, linewidth=2)
    ax2.set_xlabel("Fill Level [%]")
    ax2.set_ylabel("Normalized Amplitude (A/h)")
    ax2.set_title("Normalized Wave Amplitude")
    ax2.grid(True, alpha=0.3)

    fig.suptitle("Chen2018 Parametric — Wave Amplitude Analysis at f/f₁=0.9",
                 fontsize=13, fontweight="bold")
    plt.tight_layout()
    path = OUT_DIR / "figures" / "fig3b_amplitude_vs_fill.png"
    plt.savefig(path, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {path}")

    return amplitudes_left, amplitudes_right, freqs


# ── Table 4: Summary ─────────────────────────────────────────────

def generate_table(amplitudes_left, amplitudes_right, freqs):
    """Generate Table 4 summary."""
    rows = []
    for idx, h_mm in enumerate(FILLS):
        f1 = natural_freq(h_mm / 1000)
        f_exc = 0.9 * f1

        t_l, swl_l = load_elevation(f"simulations/chen_h{h_mm}", probe_idx=0)
        max_swl = np.max(swl_l) if swl_l is not None else 0
        min_swl = np.min(swl_l) if swl_l is not None else 0

        row = {
            "fill_mm": h_mm,
            "fill_pct": f"{h_mm / TANK_H * 100:.1f}",
            "f1_hz": f"{f1:.4f}",
            "f_exc_hz": f"{f_exc:.4f}",
            "amp_left_mm": f"{amplitudes_left[idx]:.1f}",
            "amp_right_mm": f"{amplitudes_right[idx]:.1f}",
            "max_swl_mm": f"{max_swl:.1f}",
            "min_swl_mm": f"{min_swl:.1f}",
            "dom_freq_hz": f"{freqs[idx]:.3f}",
        }
        rows.append(row)

    json_path = OUT_DIR / "comparison" / "table4_summary.json"
    with open(json_path, "w") as f:
        json.dump(rows, f, indent=2)

    md_path = OUT_DIR / "comparison" / "table4_summary.md"
    with open(md_path, "w") as f:
        f.write("# Table 4: Chen2018 Parametric Study — SWL Analysis\n\n")
        f.write("| Fill [mm] | Fill [%] | f₁ [Hz] | f_exc [Hz] | "
                "Amp Left [mm] | Amp Right [mm] | Max SWL [mm] | Min SWL [mm] | "
                "Dom Freq [Hz] |\n")
        f.write("|-----------|---------|---------|-----------|"
                "--------------|---------------|-------------|-------------|"
                "---------------|\n")
        for r in rows:
            f.write(f"| {r['fill_mm']} | {r['fill_pct']} | {r['f1_hz']} | "
                    f"{r['f_exc_hz']} | {r['amp_left_mm']} | {r['amp_right_mm']} | "
                    f"{r['max_swl_mm']} | {r['min_swl_mm']} | {r['dom_freq_hz']} |\n")

        f.write("\n## Parameters\n\n")
        f.write("- Tank: 600 x 300 x 650 mm\n")
        f.write("- Excitation: Horizontal sway, f/f1 = 0.9, A = 7mm\n")
        f.write("- Solver: DualSPHysics v5.4, DBC, dp=0.005m\n")
        f.write("- Duration: 10s per case\n")

    print(f"Saved: {json_path}")
    print(f"Saved: {md_path}")


# ── Main ──────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("Analyzing Chen2018 parametric results...\n")
    fig_swl_timeseries()
    amps_l, amps_r, freqs = fig_amplitude_vs_fill()
    generate_table(amps_l, amps_r, freqs)
    print("\nDone!")
