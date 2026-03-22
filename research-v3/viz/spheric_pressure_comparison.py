#!/usr/bin/env python3
"""SPHERIC Test 10 — Simulation vs Experiment pressure time-series comparison.
Peak-based time alignment: align first major impact peak, then compare.

Water Lateral: run_001 (dp=0.004), run_002 (dp=0.002), run_003v2 (dp=0.001)
Oil Lateral:   run_005 (dp=0.002, mDBC), run_010 (probe fix), run_011 (dp=0.001)

Probe positions corrected for 1.5h rule:
  dp=0.004 → x=0.013 (1.56h)
  dp=0.002 → x=0.007 (1.68h)
  dp=0.001 → x=0.005 (1.51h)
"""
from pathlib import Path
import numpy as np
from scipy.signal import find_peaks

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import matplotlib.gridspec as gridspec

OUT_DIR = Path(__file__).parent / "plots"
OUT_DIR.mkdir(exist_ok=True)

SIMDATA = Path("/mnt/simdata/dualsphysics/exp1")
SIMDATA_C = Path("/mnt/simdata/dualsphysics/exp-c")
EXPDATA = Path(__file__).resolve().parents[2] / "datasets" / "spheric" / "case_1"
ANALYSIS_DIR = Path(__file__).resolve().parents[1] / "exp-c" / "analysis"


def load_exp_timeseries(filepath):
    """Load SPHERIC experimental time series (tab-sep, mbar → Pa)."""
    data = np.loadtxt(str(filepath), skiprows=1)
    return data[:, 0], data[:, 1] * 100  # t, P(Pa)


def load_exp_peaks(filepath):
    """Load SPHERIC peak statistics (103 trials, mbar → Pa)."""
    data = np.loadtxt(str(filepath), skiprows=2)
    return data * 100  # (103, 4) in Pa


def load_dsph_csv(filepath, probe=0):
    """Load DualSPHysics MeasureTool pressure CSV."""
    times, pressures = [], []
    with open(filepath) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith(('#', 'Part', ' ')):
                continue
            parts = line.split(';')
            try:
                times.append(float(parts[1]))
                pressures.append(float(parts[2 + probe]))
            except (ValueError, IndexError):
                continue
    return np.array(times), np.array(pressures)


def find_impact_peaks(t, p, min_height, min_spacing_sec, max_height=None,
                      smooth_win=5):
    """Detect impact peaks on smoothed signal to avoid SPH spikes."""
    if len(t) < 2:
        return [], []
    p_sm = smooth_pressure(p, smooth_win) if smooth_win > 1 else p
    dt = np.median(np.diff(t))
    dist = max(1, int(min_spacing_sec / dt))
    peaks, _ = find_peaks(p_sm, height=min_height, distance=dist,
                          prominence=min_height * 0.3)
    # Skip startup transients: require t > 0.5s
    valid = [pk for pk in peaks if t[pk] > 0.5]
    # Filter remaining SPH spikes
    if max_height is not None:
        valid = [pk for pk in valid if p_sm[pk] <= max_height]
    return [t[pk] for pk in valid], [p_sm[pk] for pk in valid]


def compute_peak_shift(exp_peak_times, sim_peak_times, n_match=3):
    """Compute mean time shift from matched peaks (exp - sim)."""
    n = min(n_match, len(exp_peak_times), len(sim_peak_times))
    if n == 0:
        return 0.0
    shifts = [exp_peak_times[i] - sim_peak_times[i] for i in range(n)]
    return np.mean(shifts)


def smooth_pressure(p, window=5):
    """Moving average smoothing to suppress SPH pressure spikes."""
    if len(p) < window:
        return p
    kernel = np.ones(window) / window
    return np.convolve(p, kernel, mode='same')


def find_peak_in_window(t, p, t_center, window=0.4, smooth_win=5):
    """Find max pressure within time window, with spike filtering."""
    mask = (t >= t_center - window) & (t <= t_center + window)
    if mask.sum() == 0:
        return 0.0, t_center
    p_smooth = smooth_pressure(p[mask], smooth_win) if smooth_win > 1 else p[mask]
    idx = np.argmax(p_smooth)
    return p_smooth[idx], t[mask][idx]


# ============================================================
# Water Lateral
# ============================================================

def plot_water_lateral():
    """Water Lateral: dp convergence study vs experiment."""
    print("\n=== Water Lateral ===")

    # Experimental data
    t_exp, p_exp = load_exp_timeseries(EXPDATA / "lateral_water_1x.txt")
    peak_stats = load_exp_peaks(EXPDATA / "Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt")
    pk_mean = peak_stats.mean(axis=0)
    pk_2sig = 2 * peak_stats.std(axis=0)

    # Experimental peak times
    exp_pk_t, exp_pk_p = find_impact_peaks(t_exp, p_exp, min_height=1000, min_spacing_sec=1.0)
    print(f"  Exp: {len(t_exp)} pts, {len(exp_pk_t)} peaks")
    for i in range(min(4, len(exp_pk_t))):
        print(f"    P{i+1}: t={exp_pk_t[i]:.3f}s, P={exp_pk_p[i]:.0f} Pa "
              f"(stats: {pk_mean[i]:.0f} +/- {pk_2sig[i]:.0f})")

    # Simulation runs — consistent probe at x=0.013m (safe for all dp)
    # probe index differs per CSV: run_001→0, run_002→3, run_003v2→TBD
    runs = [
        ('run_001', SIMDATA / 'run_001' / 'PressLateral_dp004_Press.csv',
         0, 'dp=4mm (DBC)', '#1976D2', 1.0),
        ('run_002', SIMDATA / 'run_002' / 'PressLateral_dp002_Press.csv',
         3, 'dp=2mm (DBC)', '#F57C00', 1.0),
        ('run_003v2', SIMDATA_C / 'run_003v2' / 'PressLateral_Press.csv',
         4, 'dp=1mm (DBC)', '#388E3C', 1.0),  # probe 4 = x=0.013
    ]

    sim_data = {}
    for name, csv, probe_idx, label, color, lw in runs:
        if not csv.exists():
            print(f"  {name}: SKIP (not found: {csv})")
            continue
        t, p = load_dsph_csv(csv, probe=probe_idx)
        # Water peaks ~3700 Pa; cap at 10000 Pa to filter SPH spikes
        sim_pk_t, sim_pk_p = find_impact_peaks(t, p, min_height=800,
                                                min_spacing_sec=1.0, max_height=10000)

        # Peak-based time shift
        shift = compute_peak_shift(exp_pk_t, sim_pk_t, n_match=3)
        sim_data[name] = {
            't': t, 'p': p, 'label': label, 'color': color, 'lw': lw,
            'shift': shift, 'pk_t': sim_pk_t, 'pk_p': sim_pk_p
        }
        print(f"  {name} ({label}): shift={shift:.3f}s, peaks={len(sim_pk_t)}")
        for i in range(min(4, len(sim_pk_t))):
            print(f"    P{i+1}: t={sim_pk_t[i]:.3f}→{sim_pk_t[i]+shift:.3f}s, P={sim_pk_p[i]:.0f} Pa")

    if not sim_data:
        print("  No simulation data available, skipping figure")
        return

    # ---- Figure: 3 rows (timeseries + peak zoom + peak bars) ----
    fig = plt.figure(figsize=(16, 14))
    gs = gridspec.GridSpec(3, 1, height_ratios=[2.5, 2, 1], hspace=0.22)

    # Row 1: Full time series
    ax1 = fig.add_subplot(gs[0])
    step = 100  # downsample 20kHz → 200Hz
    ax1.plot(t_exp[::step], p_exp[::step], '-', color='#333333', linewidth=0.5,
             alpha=0.6, label='Experiment (20kHz)', zorder=1)

    for i in range(min(4, len(exp_pk_t))):
        ax1.plot(exp_pk_t[i], exp_pk_p[i], 'k^', markersize=9, zorder=10)
        ax1.annotate(f'P{i+1}: {exp_pk_p[i]:.0f}',
                     (exp_pk_t[i], exp_pk_p[i]),
                     textcoords='offset points', xytext=(8, 5), fontsize=8)

    for name, r in sim_data.items():
        t_s = r['t'] + r['shift']
        ax1.plot(t_s, r['p'], '-', color=r['color'], linewidth=r['lw'], alpha=0.8,
                 label=f"{r['label']} (shift={r['shift']:+.2f}s)", zorder=3)

    ax1.set_ylabel('Pressure (Pa)', fontsize=12)
    ax1.set_title('SPHERIC Test 10 — Water Lateral — Pressure at z=0.093m',
                   fontsize=14, fontweight='bold')
    ax1.legend(loc='upper left', fontsize=9, framealpha=0.9)
    ax1.grid(True, alpha=0.2)
    ax1.set_xlim(0, 8)

    # Row 2: Peak zoom (4 subplots)
    gs_inner = gs[1].subgridspec(1, 4, wspace=0.25)
    zoom_half = 0.25  # +/- 0.25s
    n_peaks = min(4, len(exp_pk_t))

    for col in range(n_peaks):
        ax = fig.add_subplot(gs_inner[col])
        tc = exp_pk_t[col]
        tlo, thi = tc - zoom_half, tc + zoom_half

        # Exp full-res in zoom
        mask = (t_exp >= tlo) & (t_exp <= thi)
        ax.plot(t_exp[mask], p_exp[mask], '-', color='#333333', linewidth=0.8,
                alpha=0.7, label='Exp' if col == 0 else None)

        # 2-sigma band
        ax.axhspan(pk_mean[col] - pk_2sig[col], pk_mean[col] + pk_2sig[col],
                   alpha=0.12, color='gray', label='2sigma' if col == 0 else None)
        ax.axhline(pk_mean[col], color='gray', ls='--', lw=0.7, alpha=0.4)

        # Sim
        for name, r in sim_data.items():
            ts = r['t'] + r['shift']
            mask_s = (ts >= tlo) & (ts <= thi)
            if mask_s.sum() > 0:
                ax.plot(ts[mask_s], r['p'][mask_s], '-', color=r['color'],
                        linewidth=1.2, alpha=0.85,
                        label=r['label'] if col == 0 else None)

        ax.set_title(f'Peak {col+1}\nExp: {exp_pk_p[col]:.0f} Pa', fontsize=10, fontweight='bold')
        ax.set_xlim(tlo, thi)
        ax.grid(True, alpha=0.2)
        ax.tick_params(labelsize=8)
        if col == 0:
            ax.set_ylabel('P (Pa)', fontsize=10)
            ax.legend(fontsize=7, loc='upper left')
        ax.set_xlabel('Time (s)', fontsize=8)

    # Row 3: Peak bar comparison
    ax3 = fig.add_subplot(gs[2])
    x = np.arange(n_peaks)
    bw = 0.8 / (len(sim_data) + 1)

    # Exp bars
    ax3.bar(x - 0.4 + bw/2, pk_mean[:n_peaks], bw, yerr=pk_2sig[:n_peaks],
            color='#333333', alpha=0.6, label='Exp (103 trials)', capsize=3)

    for i, (name, r) in enumerate(sim_data.items()):
        ts = r['t'] + r['shift']
        vals = []
        for j in range(n_peaks):
            pv, _ = find_peak_in_window(ts, r['p'], exp_pk_t[j], 0.35)
            vals.append(pv)
            err = abs(pv - pk_mean[j]) / pk_mean[j] * 100 if pk_mean[j] > 0 else 0
            in_band = (pk_mean[j] - pk_2sig[j]) <= pv <= (pk_mean[j] + pk_2sig[j])
            marker = 'OK' if in_band else 'MISS'
            print(f"  {name} P{j+1}: sim={pv:.0f} exp={pk_mean[j]:.0f} err={err:.1f}% [{marker}]")

        ax3.bar(x - 0.4 + bw/2 + (i+1)*bw, vals, bw,
                color=r['color'], alpha=0.7, label=r['label'])

    ax3.set_xticks(x)
    ax3.set_xticklabels([f'Peak {i+1}' for i in range(n_peaks)])
    ax3.set_ylabel('Peak Pressure (Pa)', fontsize=11)
    ax3.legend(loc='upper left', fontsize=8, ncol=2)
    ax3.grid(axis='y', alpha=0.2)

    fig.savefig(OUT_DIR / "spheric_water_lateral.png", dpi=200, bbox_inches='tight')
    fig.savefig(OUT_DIR / "spheric_water_lateral.pdf", bbox_inches='tight')
    plt.close(fig)
    print("  Saved: spheric_water_lateral.png/pdf")


# ============================================================
# Oil Lateral
# ============================================================

def plot_oil_lateral():
    """Oil Lateral: Artificial vs Laminar+SPS vs experiment."""
    print("\n=== Oil Lateral ===")

    t_exp, p_exp = load_exp_timeseries(EXPDATA / "lateral_oil_1x.txt")
    peak_stats = load_exp_peaks(EXPDATA / "Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt")
    pk_mean = peak_stats.mean(axis=0)
    pk_2sig = 2 * peak_stats.std(axis=0)

    exp_pk_t, exp_pk_p = find_impact_peaks(t_exp, p_exp, min_height=300, min_spacing_sec=1.2)
    print(f"  Exp: {len(t_exp)} pts, {len(exp_pk_t)} peaks")
    for i in range(min(4, len(exp_pk_t))):
        print(f"    P{i+1}: t={exp_pk_t[i]:.3f}s, P={exp_pk_p[i]:.0f} Pa "
              f"(stats: {pk_mean[i]:.0f} +/- {pk_2sig[i]:.0f})")

    sim_configs = [
        ('run_005', SIMDATA / 'run_005' / 'PressConsistent_Press.csv',
         'dp=2mm 1kHz (mDBC)', '#1976D2', 1.0),
        ('run_010', SIMDATA_C / 'run_010' / 'PressLateral_Press.csv',
         'dp=2mm 100Hz probe@7mm', '#90CAF9', 0.7),
        ('run_011', SIMDATA_C / 'run_011' / 'PressLateral_Press.csv',
         'dp=1mm 100Hz (mDBC)', '#388E3C', 1.0),
    ]

    sim_data = {}
    for name, csv, label, color, lw in sim_configs:
        if not csv.exists():
            continue
        t, p = load_dsph_csv(csv)
        # For oil, use lower threshold for first peak detection
        sim_pk_t, sim_pk_p = find_impact_peaks(t, p, min_height=200, min_spacing_sec=1.0)

        # Align on first peak only (subsequent peaks may drift)
        shift = compute_peak_shift(exp_pk_t[:1], sim_pk_t[:1], n_match=1)
        sim_data[name] = {
            't': t, 'p': p, 'label': label, 'color': color, 'lw': lw,
            'shift': shift, 'pk_t': sim_pk_t, 'pk_p': sim_pk_p
        }
        print(f"  {name}: shift={shift:.3f}s, peaks={len(sim_pk_t)}")
        for i in range(min(4, len(sim_pk_t))):
            print(f"    P{i+1}: t={sim_pk_t[i]:.3f}→{sim_pk_t[i]+shift:.3f}s, P={sim_pk_p[i]:.0f} Pa")

    # ---- Figure ----
    fig = plt.figure(figsize=(16, 14))
    gs = gridspec.GridSpec(3, 1, height_ratios=[2.5, 2, 1], hspace=0.22)

    # Row 1: Full time series
    ax1 = fig.add_subplot(gs[0])
    step = 100
    ax1.plot(t_exp[::step], p_exp[::step], '-', color='#333333', linewidth=0.5,
             alpha=0.6, label='Experiment (20kHz)', zorder=1)

    for i in range(min(4, len(exp_pk_t))):
        ax1.plot(exp_pk_t[i], exp_pk_p[i], 'k^', markersize=9, zorder=10)
        ax1.annotate(f'P{i+1}: {exp_pk_p[i]:.0f}',
                     (exp_pk_t[i], exp_pk_p[i]),
                     textcoords='offset points', xytext=(8, 5), fontsize=8)

    for name, r in sim_data.items():
        t_s = r['t'] + r['shift']
        ax1.plot(t_s, r['p'], '-', color=r['color'], linewidth=r['lw'], alpha=0.8,
                 label=f"{r['label']} (shift={r['shift']:+.2f}s)", zorder=3)

    ax1.set_ylabel('Pressure (Pa)', fontsize=12)
    ax1.set_title('SPHERIC Test 10 — Oil Lateral — Pressure at z=0.093m',
                   fontsize=14, fontweight='bold')
    ax1.legend(loc='upper left', fontsize=9, framealpha=0.9)
    ax1.grid(True, alpha=0.2)
    ax1.set_xlim(0, 8)

    # Row 2: Peak zoom
    gs_inner = gs[1].subgridspec(1, 4, wspace=0.25)
    zoom_half = 0.3
    n_peaks = min(4, len(exp_pk_t))

    for col in range(n_peaks):
        ax = fig.add_subplot(gs_inner[col])
        tc = exp_pk_t[col]
        tlo, thi = tc - zoom_half, tc + zoom_half

        mask = (t_exp >= tlo) & (t_exp <= thi)
        ax.plot(t_exp[mask], p_exp[mask], '-', color='#333333', linewidth=0.8,
                alpha=0.7, label='Exp' if col == 0 else None)

        ax.axhspan(pk_mean[col] - pk_2sig[col], pk_mean[col] + pk_2sig[col],
                   alpha=0.12, color='gray', label='2sigma' if col == 0 else None)
        ax.axhline(pk_mean[col], color='gray', ls='--', lw=0.7, alpha=0.4)

        for name, r in sim_data.items():
            ts = r['t'] + r['shift']
            mask_s = (ts >= tlo) & (ts <= thi)
            if mask_s.sum() > 0:
                ax.plot(ts[mask_s], r['p'][mask_s], '-', color=r['color'],
                        linewidth=1.2, alpha=0.85,
                        label=r['label'] if col == 0 else None)

        ax.set_title(f'Peak {col+1}\nExp: {exp_pk_p[col]:.0f} Pa', fontsize=10, fontweight='bold')
        ax.set_xlim(tlo, thi)
        ax.grid(True, alpha=0.2)
        ax.tick_params(labelsize=8)
        if col == 0:
            ax.set_ylabel('P (Pa)', fontsize=10)
            ax.legend(fontsize=7, loc='upper left')
        ax.set_xlabel('Time (s)', fontsize=8)

    # Row 3: Peak bars
    ax3 = fig.add_subplot(gs[2])
    x = np.arange(n_peaks)
    bw = 0.8 / (len(sim_data) + 1)

    ax3.bar(x - 0.4 + bw/2, pk_mean[:n_peaks], bw, yerr=pk_2sig[:n_peaks],
            color='#333333', alpha=0.6, label='Exp (103 trials)', capsize=3)

    for i, (name, r) in enumerate(sim_data.items()):
        ts = r['t'] + r['shift']
        vals = []
        for j in range(n_peaks):
            pv, _ = find_peak_in_window(ts, r['p'], exp_pk_t[j], 0.4)
            vals.append(pv)
            err = abs(pv - pk_mean[j]) / pk_mean[j] * 100 if pk_mean[j] > 0 else 0
            in_band = (pk_mean[j] - pk_2sig[j]) <= pv <= (pk_mean[j] + pk_2sig[j])
            marker = 'OK' if in_band else 'MISS'
            print(f"  {name} P{j+1}: sim={pv:.0f} exp={pk_mean[j]:.0f} "
                  f"err={err:.1f}% [{marker}]")

        ax3.bar(x - 0.4 + bw/2 + (i+1)*bw, vals, bw,
                color=r['color'], alpha=0.7, label=r['label'])

    ax3.set_xticks(x)
    ax3.set_xticklabels([f'Peak {i+1}' for i in range(n_peaks)])
    ax3.set_ylabel('Peak Pressure (Pa)', fontsize=11)
    ax3.legend(loc='upper left', fontsize=8, ncol=2)
    ax3.grid(axis='y', alpha=0.2)

    fig.savefig(OUT_DIR / "spheric_oil_lateral.png", dpi=200, bbox_inches='tight')
    fig.savefig(OUT_DIR / "spheric_oil_lateral.pdf", bbox_inches='tight')
    plt.close(fig)
    print("  Saved: spheric_oil_lateral.png/pdf")


# ============================================================
# Combined summary table figure
# ============================================================

def plot_peak_summary():
    """Publication figure: peak comparison summary table."""
    print("\n=== Peak Summary ===")

    # Water
    w_stats = load_exp_peaks(EXPDATA / "Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt")
    w_mean = w_stats.mean(axis=0)
    w_2sig = 2 * w_stats.std(axis=0)

    t_exp_w, p_exp_w = load_exp_timeseries(EXPDATA / "lateral_water_1x.txt")
    exp_w_t, _ = find_impact_peaks(t_exp_w, p_exp_w, 1000, 1.0)

    # Oil
    o_stats = load_exp_peaks(EXPDATA / "Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt")
    o_mean = o_stats.mean(axis=0)
    o_2sig = 2 * o_stats.std(axis=0)

    t_exp_o, p_exp_o = load_exp_timeseries(EXPDATA / "lateral_oil_1x.txt")
    exp_o_t, _ = find_impact_peaks(t_exp_o, p_exp_o, 300, 1.2)

    # Collect all simulation peak values
    all_results = {}

    # Water runs — consistent probe at x=0.013m
    water_runs = [
        ('run_001', SIMDATA / 'run_001' / 'PressLateral_dp004_Press.csv', 0),
        ('run_002', SIMDATA / 'run_002' / 'PressLateral_dp002_Press.csv', 3),
        ('run_003v2', SIMDATA_C / 'run_003v2' / 'PressLateral_Press.csv', 4),
    ]
    for name, csv, probe_idx in water_runs:
        if not csv.exists():
            continue
        t, p = load_dsph_csv(csv, probe=probe_idx)
        spt, spp = find_impact_peaks(t, p, 800, 1.0, max_height=10000)
        shift = compute_peak_shift(exp_w_t, spt, 3)
        ts = t + shift

        vals, errs, bands = [], [], []
        for j in range(4):
            if j < len(exp_w_t):
                pv, _ = find_peak_in_window(ts, p, exp_w_t[j], 0.4)
            else:
                pv = 0
            vals.append(pv)
            err = abs(pv - w_mean[j]) / w_mean[j] * 100 if w_mean[j] > 0 else 0
            errs.append(err)
            bands.append((w_mean[j] - w_2sig[j]) <= pv <= (w_mean[j] + w_2sig[j]))

        all_results[name] = {
            'fluid': 'Water', 'peaks': vals, 'errors': errs, 'in_band': bands
        }

    # Oil runs
    for name, csv in [
        ('run_005', SIMDATA / 'run_005' / 'PressConsistent_Press.csv'),
        ('run_010', SIMDATA_C / 'run_010' / 'PressLateral_Press.csv'),
        ('run_011', SIMDATA_C / 'run_011' / 'PressLateral_Press.csv'),
    ]:
        if not csv.exists():
            continue
        t, p = load_dsph_csv(csv)
        spt, spp = find_impact_peaks(t, p, 200, 1.0)
        shift = compute_peak_shift(exp_o_t[:1], spt[:1], 1)
        ts = t + shift

        vals, errs, bands = [], [], []
        for j in range(4):
            if j < len(exp_o_t):
                pv, _ = find_peak_in_window(ts, p, exp_o_t[j], 0.4)
            else:
                pv = 0
            vals.append(pv)
            err = abs(pv - o_mean[j]) / o_mean[j] * 100 if o_mean[j] > 0 else 0
            errs.append(err)
            bands.append((o_mean[j] - o_2sig[j]) <= pv <= (o_mean[j] + o_2sig[j]))

        all_results[name] = {
            'fluid': 'Oil', 'peaks': vals, 'errors': errs, 'in_band': bands
        }

    # Print table
    print(f"\n  {'Run':12s} {'Fluid':6s} {'P1 (Pa)':>10s} {'P2 (Pa)':>10s} "
          f"{'P3 (Pa)':>10s} {'P4 (Pa)':>10s} {'Mean Err':>10s} {'In-band':>8s}")
    print(f"  {'-'*12} {'-'*6} {'-'*10} {'-'*10} {'-'*10} {'-'*10} {'-'*10} {'-'*8}")

    # Experiment rows
    print(f"  {'Exp(Water)':12s} {'Water':6s} {w_mean[0]:>10.0f} {w_mean[1]:>10.0f} "
          f"{w_mean[2]:>10.0f} {w_mean[3]:>10.0f} {'—':>10s} {'—':>8s}")
    print(f"  {'Exp(Oil)':12s} {'Oil':6s} {o_mean[0]:>10.0f} {o_mean[1]:>10.0f} "
          f"{o_mean[2]:>10.0f} {o_mean[3]:>10.0f} {'—':>10s} {'—':>8s}")

    for name, r in all_results.items():
        mean_err = np.mean(r['errors'])
        n_band = sum(r['in_band'])
        print(f"  {name:12s} {r['fluid']:6s} {r['peaks'][0]:>10.0f} {r['peaks'][1]:>10.0f} "
              f"{r['peaks'][2]:>10.0f} {r['peaks'][3]:>10.0f} {mean_err:>9.1f}% {n_band:>5d}/4")


def main():
    print("=" * 60)
    print("  SPHERIC Test 10 — Pressure Comparison Figures")
    print(f"  Output: {OUT_DIR}")
    print("=" * 60)

    plot_water_lateral()
    plot_oil_lateral()
    plot_peak_summary()

    print(f"\nDone!")


if __name__ == "__main__":
    main()
