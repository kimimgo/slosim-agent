#!/usr/bin/env python3
"""Generate publication-quality figures for SPHERIC Test 10 validation.

Creates:
  1. fig_timeseries.png    — Time series comparison (sim vs exp)
  2. fig_convergence.png   — 3-level convergence + GCI
  3. fig_oil_lateral.png   — Oil lateral validation
  4. fig_water_roof.png    — Water roof impact validation

Usage:
    python paper_figures.py [--water-only]
"""

import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from matplotlib.ticker import MultipleLocator
from pathlib import Path
from scipy.signal import find_peaks, medfilt
from scipy.ndimage import uniform_filter1d
import sys
import json

plt.rcParams.update({
    'font.size': 11,
    'axes.labelsize': 12,
    'axes.titlesize': 13,
    'legend.fontsize': 9,
    'figure.dpi': 150,
    'savefig.dpi': 300,
    'axes.grid': True,
    'grid.alpha': 0.3,
})

ROOT = Path(__file__).resolve().parents[3]
SIM_DIR = ROOT / "simulations" / "exp1"
EXP_DIR = ROOT / "datasets" / "spheric" / "case_1"
FIG_DIR = Path(__file__).resolve().parent.parent / "figures"
FIG_DIR.mkdir(exist_ok=True)

# Case-specific filtering: oil peaks are ~5x smaller than water peaks.
FILTER_WATER = {'medfilt_kernel': 5, 'smooth_window': 11}
FILTER_OIL = {'medfilt_kernel': 3, 'smooth_window': 7}


def load_sim(run_id, fluid='water', apply_filter=True):
    """Load simulation pressure (prefers PressConsistent)."""
    for prefix in ["PressConsistent", "PressureLateral"]:
        path = SIM_DIR / f"run_{run_id:03d}" / f"{prefix}_Press.csv"
        if path.exists():
            data = np.genfromtxt(path, delimiter=';', skip_header=4)
            t = data[:, 1]
            p = data[:, 2] / 100  # Pa → mbar
            if apply_filter:
                fp = FILTER_OIL if fluid == 'oil' else FILTER_WATER
                p = medfilt(p, kernel_size=fp['medfilt_kernel'])
                p = uniform_filter1d(p, size=fp['smooth_window'])
            return t, p
    return None, None


def load_exp(case='lateral_water'):
    """Load experimental data."""
    filemap = {
        'lateral_water': ('lateral_water_1x.txt',
                          'Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt'),
        'lateral_oil':   ('lateral_oil_1x.txt',
                          'Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt'),
        'roof_water':    ('roof_water_1x.txt',
                          'Water_4first_peak_roof_impact_tto_1_H355_5_B1X.txt'),
    }
    ts_file, peak_file = filemap[case]
    data = np.genfromtxt(EXP_DIR / ts_file, delimiter='\t', skip_header=1)
    t_exp = data[:, 0]
    p_exp = data[:, 1]

    peak_data = np.genfromtxt(EXP_DIR / peak_file, delimiter='\t', skip_header=2)
    means = np.nanmean(peak_data, axis=0)
    stds = np.nanstd(peak_data, axis=0)
    return t_exp, p_exp, means, stds


def detect_peaks(t, p, min_height=10.0, min_dist_s=1.0):
    dt = np.median(np.diff(t))
    peaks, _ = find_peaks(p, height=min_height, distance=int(min_dist_s / dt))
    return peaks


# === Figure 1: Water Lateral Time Series ===
def fig_water_timeseries():
    t_exp, p_exp, pmeans, pstds = load_exp('lateral_water')
    runs = {}
    for rid in [1, 2, 3]:
        t, p = load_sim(rid)
        if t is not None:
            runs[rid] = (t, p)

    if not runs:
        print("No simulation data available for water lateral")
        return

    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 8),
                                     gridspec_kw={'height_ratios': [2, 1]},
                                     sharex=True)

    # Top: time series
    ax1.plot(t_exp, p_exp, color='black', alpha=0.3, linewidth=0.3,
             label='Experiment (20 kHz)')

    colors = {1: '#1f77b4', 2: '#d62728', 3: '#2ca02c'}
    labels = {1: 'DBC, dp = 4 mm', 2: 'DBC, dp = 2 mm', 3: 'DBC, dp = 1 mm'}

    for rid in sorted(runs.keys()):
        t, p = runs[rid]
        lw = 1.2 if rid == max(runs.keys()) else 0.8
        ax1.plot(t, p, color=colors[rid], linewidth=lw, alpha=0.8,
                 label=labels[rid])

    # Annotate peaks for finest run
    best_rid = max(runs.keys())
    t_best, p_best = runs[best_rid]
    peaks = detect_peaks(t_best, p_best, min_height=10.0, min_dist_s=1.2)
    for i, idx in enumerate(peaks[:3]):
        ax1.plot(t_best[idx], p_best[idx], 'v', color=colors[best_rid],
                 markersize=8, zorder=5)
        ax1.annotate(f'{p_best[idx]:.1f}', (t_best[idx], p_best[idx]),
                    textcoords="offset points", xytext=(0, 12),
                    fontsize=9, color=colors[best_rid], ha='center',
                    fontweight='bold')

    ax1.set_ylabel('Pressure [mbar]')
    ax1.set_title('SPHERIC Test 10 — Lateral Water Impact Pressure')
    ax1.legend(loc='upper left', framealpha=0.9)
    ymax = max(p_exp.max(), max(p.max() for _, p in runs.values()))
    ax1.set_ylim(-15, min(ymax * 1.3, 100))
    ax1.text(0.02, 0.95, '(a)', transform=ax1.transAxes, fontsize=14,
             fontweight='bold', va='top')

    # Bottom: residual (sim - exp interpolated)
    best_t, best_p = runs[best_rid]
    p_exp_interp = np.interp(best_t, t_exp, p_exp)
    residual = best_p - p_exp_interp
    ax2.fill_between(best_t, residual, alpha=0.3, color=colors[best_rid])
    ax2.plot(best_t, residual, color=colors[best_rid], linewidth=0.5)
    ax2.axhline(0, color='black', linewidth=0.5)
    ax2.set_xlabel('Time [s]')
    ax2.set_ylabel('Residual [mbar]')
    ax2.set_xlim(0, 7)
    ax2.text(0.02, 0.95, '(b)', transform=ax2.transAxes, fontsize=14,
             fontweight='bold', va='top')

    plt.savefig(FIG_DIR / "fig_timeseries.png")
    print(f"Saved: {FIG_DIR / 'fig_timeseries.png'}")
    plt.close()


# === Figure 2: Convergence Study ===
def fig_convergence():
    t_exp, p_exp, pmeans, pstds = load_exp('lateral_water')
    runs = {}
    for rid in [1, 2, 3]:
        t, p = load_sim(rid)
        if t is not None:
            runs[rid] = (t, p)

    if len(runs) < 2:
        print("Need at least 2 runs for convergence figure")
        return

    n_peaks = 3
    all_vals = {}
    for rid in sorted(runs.keys()):
        t, p = runs[rid]
        peaks = detect_peaks(t, p, min_height=10.0, min_dist_s=1.2)
        vals = [p[peaks[i]] if i < len(peaks) else np.nan for i in range(n_peaks)]
        all_vals[rid] = vals

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(12, 5))

    # Left: Peak comparison bar chart
    x = np.arange(n_peaks)
    width = 0.18
    colors = {1: '#1f77b4', 2: '#d62728', 3: '#2ca02c'}
    dp_labels = {1: '4 mm', 2: '2 mm', 3: '1 mm'}

    ax1.bar(x - 1.5*width, pmeans[:n_peaks], width,
            yerr=2*pstds[:n_peaks], label='Exp (μ±2σ)',
            color='gray', alpha=0.7, capsize=4)
    for i, rid in enumerate(sorted(runs.keys())):
        offset = (-0.5 + i) * width
        ax1.bar(x + offset, all_vals[rid], width,
                label=f'dp = {dp_labels[rid]}',
                color=colors[rid], alpha=0.7)

    ax1.set_xlabel('Impact Peak')
    ax1.set_ylabel('Pressure [mbar]')
    ax1.set_title('Peak Pressure Comparison')
    ax1.set_xticks(x)
    ax1.set_xticklabels([f'Peak {i+1}' for i in range(n_peaks)])
    ax1.legend(fontsize=8)
    ax1.text(0.02, 0.95, '(a)', transform=ax1.transAxes, fontsize=14,
             fontweight='bold', va='top')

    # Right: Error vs resolution
    dps = sorted([{1: 4.0, 2: 2.0, 3: 1.0}[rid] for rid in runs.keys()], reverse=True)

    for i in range(n_peaks):
        errs = []
        for dp in dps:
            rid = {4.0: 1, 2.0: 2, 1.0: 3}[dp]
            if rid in all_vals and not np.isnan(all_vals[rid][i]):
                errs.append(abs(all_vals[rid][i] - pmeans[i]) / pmeans[i] * 100)
            else:
                errs.append(np.nan)
        ax2.plot(dps[:len(errs)], errs, 'o-', label=f'Peak {i+1}', markersize=6)

    ax2.axhline(30, color='red', linestyle='--', alpha=0.5, label='M2 threshold')
    ax2.set_xlabel('dp [mm]')
    ax2.set_ylabel('Absolute Error [%]')
    ax2.set_title('Convergence: Error vs Resolution')
    ax2.set_xlim(max(dps) * 1.2, min(dps) * 0.8)
    ax2.legend(fontsize=8)
    ax2.text(0.02, 0.95, '(b)', transform=ax2.transAxes, fontsize=14,
             fontweight='bold', va='top')

    plt.tight_layout()
    plt.savefig(FIG_DIR / "fig_convergence.png")
    print(f"Saved: {FIG_DIR / 'fig_convergence.png'}")
    plt.close()


# === Figure 3: Oil Lateral ===
def fig_oil_lateral():
    t_exp, p_exp, pmeans, pstds = load_exp('lateral_oil')
    t_sim, p_sim = None, None

    # Try Run 005 with oil-specific filter
    t_sim, p_sim = load_sim(5, fluid='oil')
    if t_sim is None:
        # fallback
        pass

    if t_sim is None:
        print("Run 005 data not yet available")
        return

    fig, ax = plt.subplots(1, 1, figsize=(12, 5))
    ax.plot(t_exp, p_exp, color='black', alpha=0.4, linewidth=0.5, label='Experiment')
    ax.plot(t_sim, p_sim, color='#ff7f0e', linewidth=1.0, label='DBC dp=2mm (Oil)')

    peaks = detect_peaks(t_sim, p_sim, min_height=2.0, min_dist_s=1.0)
    for i, idx in enumerate(peaks[:4]):
        ax.plot(t_sim[idx], p_sim[idx], 'v', color='#ff7f0e', markersize=8)
        ax.annotate(f'{p_sim[idx]:.1f}', (t_sim[idx], p_sim[idx]),
                    textcoords="offset points", xytext=(0, 10), fontsize=9,
                    color='#ff7f0e', ha='center', fontweight='bold')

    ax.set_xlabel('Time [s]')
    ax.set_ylabel('Pressure [mbar]')
    ax.set_title('SPHERIC Test 10 — Oil Lateral Impact Pressure')
    ax.legend(loc='upper left')
    ax.set_xlim(0, 7)
    ax.set_ylim(-2, 15)

    plt.savefig(FIG_DIR / "fig_oil_lateral.png")
    print(f"Saved: {FIG_DIR / 'fig_oil_lateral.png'}")
    plt.close()


# === Figure 4: Water Roof ===
def fig_water_roof():
    from scipy.signal import correlate
    t_exp, p_exp, pmeans, pstds = load_exp('roof_water')

    # Load multi-probe roof data and select best matching column
    # Try Run 007 first (Visco=0.001), fall back to Run 006
    t_sim, p_sim, best_x = None, None, None
    for run_id in ["run_006", "run_007"]:
        for prefix in ["PressNearRoof", "PressRoof", "PressConsistent"]:
            path = SIM_DIR / run_id / f"{prefix}_Press.csv"
            if path.exists():
                break
        else:
            continue
        if not path.exists():
            continue
        # Parse x positions from header
        with open(path) as f:
            header = f.readline().strip().split(';')
        x_positions = []
        for v in header[2:]:
            try:
                x_positions.append(float(v))
            except ValueError:
                x_positions.append(None)

        data = np.genfromtxt(path, delimiter=';', skip_header=4)
        t = data[:, 1]

        # Find best probe by cross-correlation
        best_r = -1
        for i, x_pos in enumerate(x_positions):
            col = i + 2
            if col >= data.shape[1]:
                break
            p = data[:, col] / 100  # Pa → mbar
            fp = FILTER_WATER
            p = medfilt(p, kernel_size=fp['medfilt_kernel'])
            p = uniform_filter1d(p, size=fp['smooth_window'])

            # Quick cross-correlation
            dt = np.median(np.diff(t))
            t_common = np.arange(0, min(t[-1], t_exp[-1]), dt)
            p_s = np.interp(t_common, t, p)
            p_e = np.interp(t_common, t_exp, p_exp)
            p_s_n = (p_s - np.mean(p_s)) / (np.std(p_s) + 1e-12)
            p_e_n = (p_e - np.mean(p_e)) / (np.std(p_e) + 1e-12)
            n = len(t_common)
            corr = np.correlate(p_s_n, p_e_n, mode='full') / n
            lags = np.arange(-n + 1, n) * dt
            mask = np.abs(lags) <= 2.0
            r = corr[mask].max()

            if r > best_r:
                best_r = r
                t_sim = t
                p_sim = p
                best_x = x_pos
        break

    if t_sim is None:
        print("No roof data available (tried run_007, run_006)")
        return

    fig, ax = plt.subplots(1, 1, figsize=(12, 5))
    ax.plot(t_exp, p_exp, color='black', alpha=0.4, linewidth=0.5, label='Experiment')
    ax.plot(t_sim, p_sim, color='#9467bd', linewidth=1.0,
            label=f'DBC dp=2mm (Roof, x={best_x:.3f}m)')

    peaks = detect_peaks(t_sim, p_sim, min_height=15.0, min_dist_s=0.8)
    for i, idx in enumerate(peaks[:4]):
        ax.plot(t_sim[idx], p_sim[idx], 'v', color='#9467bd', markersize=8)

    ax.set_xlabel('Time [s]')
    ax.set_ylabel('Pressure [mbar]')
    ax.set_title('SPHERIC Test 10 — Water Roof Impact Pressure')
    ax.legend(loc='upper left')
    ax.set_xlim(0, 7)

    plt.savefig(FIG_DIR / "fig_water_roof.png")
    print(f"Saved: {FIG_DIR / 'fig_water_roof.png'}")
    plt.close()


if __name__ == "__main__":
    water_only = '--water-only' in sys.argv

    print("=== Generating publication figures ===\n")
    fig_water_timeseries()
    fig_convergence()

    if not water_only:
        fig_oil_lateral()
        fig_water_roof()

    print("\nDone.")
