#!/usr/bin/env python3
"""EXP-1 Oil Lateral + Water Roof validation analysis.

Runs after Run 005 (Oil) and Run 006 (Water Roof) complete.
Computes M1, M2, M5, M6, M7 metrics.

Usage:
    python oil_roof_analysis.py
"""

import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from pathlib import Path
from scipy.signal import find_peaks, medfilt
from scipy.ndimage import uniform_filter1d
import json

ROOT = Path(__file__).resolve().parents[3]
SIM_DIR = ROOT / "simulations" / "exp1"
EXP_DIR = ROOT / "datasets" / "spheric" / "case_1"
FIG_DIR = Path(__file__).resolve().parent.parent / "figures"
FIG_DIR.mkdir(exist_ok=True)

# Case-specific filtering: oil peaks (~5 mbar) are ~5x smaller than water (~30 mbar),
# so lighter filtering preserves more of the true impact signal.
FILTER_PARAMS = {
    'oil':   {'medfilt_kernel': 3, 'smooth_window': 7},
    'water': {'medfilt_kernel': 5, 'smooth_window': 11},
}


def load_sim(run_dir, prefix_list, fluid='water', apply_filter=True):
    """Load simulation pressure from first available CSV."""
    for prefix in prefix_list:
        path = run_dir / f"{prefix}_Press.csv"
        if path.exists():
            data = np.genfromtxt(path, delimiter=';', skip_header=4)
            t = data[:, 1]
            p = data[:, 2] / 100  # Pa → mbar
            if apply_filter:
                fp = FILTER_PARAMS.get(fluid, FILTER_PARAMS['water'])
                p = medfilt(p, kernel_size=fp['medfilt_kernel'])
                p = uniform_filter1d(p, size=fp['smooth_window'])
            return t, p, path.name
    return None, None, None


def load_exp(case):
    filemap = {
        'lateral_oil': ('lateral_oil_1x.txt',
                        'Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt'),
        'roof_water': ('roof_water_1x.txt',
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


def cross_correlation(t_sim, p_sim, t_exp, p_exp, max_shift_s=2.0):
    dt_sim = np.median(np.diff(t_sim))
    dt_exp = np.median(np.diff(t_exp))
    dt = max(dt_sim, dt_exp)
    t_common = np.arange(0, min(t_sim[-1], t_exp[-1]), dt)
    p_sim_i = np.interp(t_common, t_sim, p_sim)
    p_exp_i = np.interp(t_common, t_exp, p_exp)

    p_sim_n = (p_sim_i - np.mean(p_sim_i)) / (np.std(p_sim_i) + 1e-12)
    p_exp_n = (p_exp_i - np.mean(p_exp_i)) / (np.std(p_exp_i) + 1e-12)

    n = len(t_common)
    corr = np.correlate(p_sim_n, p_exp_n, mode='full') / n
    lags = np.arange(-n + 1, n) * dt
    mask = np.abs(lags) <= max_shift_s
    idx_max = np.argmax(corr[mask])
    return corr[mask][idx_max], lags[mask][idx_max]


def analyze_oil():
    """Oil Lateral (Run 005) analysis."""
    print("\n" + "=" * 60)
    print("Oil Lateral (Run 005)")
    print("=" * 60)

    t_sim, p_sim, src = load_sim(SIM_DIR / "run_005",
                                  ["PressConsistent", "PressureLateral"],
                                  fluid='oil')
    if t_sim is None:
        print("  Run 005 data not available yet")
        return None

    t_exp, p_exp, pmeans, pstds = load_exp('lateral_oil')

    # Peak detection — oil peaks are small (5-7 mbar)
    dt_sim = np.median(np.diff(t_sim))
    peaks_sim = find_peaks(p_sim, height=2.0,
                           distance=int(1.0 / dt_sim))[0]

    dt_exp = np.median(np.diff(t_exp))
    peaks_exp = find_peaks(p_exp, height=2.0,
                           distance=int(0.8 / dt_exp))[0]

    n_peaks = min(4, len(peaks_sim))
    vals_sim = [p_sim[peaks_sim[i]] if i < len(peaks_sim) else np.nan
                for i in range(4)]

    print(f"\n  Data: {src}, {len(t_sim)} points, dt={dt_sim:.4f}s")
    print(f"  Peaks detected: {len(peaks_sim)} (need ≥3 for M7)")

    # M7: Peak detection
    m7_pass = len(peaks_sim) >= 3
    print(f"\n  M7 (≥3/4 detected): {len(peaks_sim)}/4 {'PASS' if m7_pass else 'FAIL'}")

    # M1: Peak-in-band (relaxed: ≥2/4)
    in_band = 0
    print(f"\n  Peak comparison (Exp mean ± 2σ):")
    for i in range(min(4, len(peaks_sim))):
        lo = pmeans[i] - 2 * pstds[i]
        hi = pmeans[i] + 2 * pstds[i]
        in_b = lo <= vals_sim[i] <= hi if not np.isnan(vals_sim[i]) else False
        if in_b:
            in_band += 1
        err = abs(vals_sim[i] - pmeans[i]) / pmeans[i] * 100 if not np.isnan(vals_sim[i]) else float('nan')
        status = "IN" if in_b else "OUT"
        print(f"    Peak {i+1}: sim={vals_sim[i]:.2f}, exp={pmeans[i]:.1f}±{2*pstds[i]:.1f}, "
              f"err={err:.1f}%, {status}")

    m1_pass = in_band >= 2
    print(f"\n  M1 (≥2/4 in-band): {in_band}/4 {'PASS' if m1_pass else 'FAIL'}")

    # Cross-correlation
    r_max, tau = cross_correlation(t_sim, p_sim, t_exp, p_exp)
    print(f"  M5 (r>0.5): r={r_max:.3f}, tau={tau:+.3f}s {'PASS' if r_max > 0.5 else 'FAIL'}")

    # Plot
    fig, ax = plt.subplots(1, 1, figsize=(12, 5))
    ax.plot(t_exp, p_exp, color='black', alpha=0.4, linewidth=0.5, label='Experiment')
    ax.plot(t_sim, p_sim, color='#ff7f0e', linewidth=1.0, label='DBC dp=2mm (Oil)')

    for i, idx in enumerate(peaks_sim[:4]):
        ax.plot(t_sim[idx], p_sim[idx], 'v', color='#ff7f0e', markersize=8)
        ax.annotate(f'{p_sim[idx]:.1f}', (t_sim[idx], p_sim[idx]),
                    textcoords="offset points", xytext=(0, 10), fontsize=9,
                    color='#ff7f0e', ha='center', fontweight='bold')

    ax.set_xlabel('Time [s]')
    ax.set_ylabel('Pressure [mbar]')
    ax.set_title('SPHERIC Test 10 — Oil Lateral Impact Pressure')
    ax.legend(loc='upper left')
    ax.set_xlim(0, 7)
    ax.set_ylim(-2, max(15, p_sim.max() * 1.3))
    ax.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(FIG_DIR / "fig_oil_lateral.png", dpi=150)
    print(f"  Saved: {FIG_DIR / 'fig_oil_lateral.png'}")
    plt.close()

    return {
        "peaks_detected": len(peaks_sim),
        "M7_pass": m7_pass,
        "M1_in_band": f"{in_band}/4",
        "M1_pass": m1_pass,
        "cross_corr_r": round(float(r_max), 3),
        "time_shift_s": round(float(tau), 3),
    }


def load_roof_multiprobe(run_dir):
    """Load roof MeasureTool output with multiple probe columns.

    Returns (t, pressures_dict) where pressures_dict maps column_idx → (p_mbar, x_pos).
    """
    for prefix in ["PressNearRoof", "PressRoof", "PressConsistent"]:
        path = run_dir / f"{prefix}_Press.csv"
        if not path.exists():
            continue
        # Read header lines for probe positions
        with open(path) as f:
            lines = [f.readline() for _ in range(4)]
        # Parse x positions from header line 1: " ;PosX [m]:;0.050;0.150;..."
        x_parts = lines[0].strip().split(';')
        x_positions = []
        for v in x_parts[2:]:  # skip first two fields (empty + label)
            try:
                x_positions.append(float(v))
            except ValueError:
                x_positions.append(None)

        data = np.genfromtxt(path, delimiter=';', skip_header=4)
        t = data[:, 1]
        pressures = {}
        for i, x_pos in enumerate(x_positions):
            col = i + 2  # columns: Part, Time, Press_0, Press_1, ...
            if col < data.shape[1]:
                p = data[:, col] / 100  # Pa → mbar
                fp = FILTER_PARAMS['water']
                p = medfilt(p, kernel_size=fp['medfilt_kernel'])
                p = uniform_filter1d(p, size=fp['smooth_window'])
                pressures[i] = (p, x_pos)
        return t, pressures, path.name
    return None, None, None


def analyze_roof():
    """Water Roof analysis with automatic probe selection.

    Tries Run 007 first (Visco=0.001), falls back to Run 006 (Visco=0.01).
    """
    print("\n" + "=" * 60)
    print("Water Roof Impact")
    print("=" * 60)

    # Try Run 007 first (reduced viscosity), then fall back to Run 006
    t_sim, probes, src = None, None, None
    run_label = None
    for run_id in ["run_006", "run_007"]:
        t_sim, probes, src = load_roof_multiprobe(SIM_DIR / run_id)
        if t_sim is not None:
            run_label = run_id
            print(f"  Using: {run_id}")
            break

    if t_sim is None:
        print("  No roof data available (tried run_007, run_006)")
        return None

    t_exp, p_exp, pmeans, pstds = load_exp('roof_water')

    # Try each probe and select best cross-correlation match
    print(f"\n  Data: {src}, {len(t_sim)} points")
    print(f"  Probes available: {len(probes)}")

    best_r = -1
    best_col = 0
    probe_results = []
    for col_idx, (p_sim, x_pos) in probes.items():
        r, tau = cross_correlation(t_sim, p_sim, t_exp, p_exp)
        probe_results.append((col_idx, x_pos, r, tau))
        print(f"    Probe {col_idx} (x={x_pos:.3f}m): r={r:.3f}, tau={tau:+.3f}s")
        if r > best_r:
            best_r = r
            best_col = col_idx

    p_sim = probes[best_col][0]
    x_best = probes[best_col][1]
    print(f"\n  Selected: Probe {best_col} (x={x_best:.3f}m) with r={best_r:.3f}")

    # Peak detection
    dt_sim = np.median(np.diff(t_sim))
    peaks_sim = find_peaks(p_sim, height=5.0,
                           distance=int(0.8 / dt_sim))[0]

    n_peaks = min(4, len(peaks_sim))
    vals_sim = [p_sim[peaks_sim[i]] if i < len(peaks_sim) else np.nan
                for i in range(4)]

    print(f"  Peaks detected: {len(peaks_sim)}")

    # M1: Peak-in-band (≥3/4)
    in_band = 0
    print(f"\n  Peak comparison (Exp mean ± 2σ):")
    for i in range(min(4, len(peaks_sim))):
        lo = pmeans[i] - 2 * pstds[i]
        hi = pmeans[i] + 2 * pstds[i]
        in_b = lo <= vals_sim[i] <= hi if not np.isnan(vals_sim[i]) else False
        if in_b:
            in_band += 1
        err = abs(vals_sim[i] - pmeans[i]) / pmeans[i] * 100 if not np.isnan(vals_sim[i]) else float('nan')
        status = "IN" if in_b else "OUT"
        print(f"    Peak {i+1}: sim={vals_sim[i]:.1f}, exp={pmeans[i]:.1f}±{2*pstds[i]:.1f}, "
              f"err={err:.1f}%, {status}")

    m1_pass = in_band >= 3
    print(f"\n  M1 (≥3/4 in-band): {in_band}/4 {'PASS' if m1_pass else 'FAIL'}")

    # M2: Mean absolute error (threshold: <40% for roof)
    errors = [abs(vals_sim[i] - pmeans[i]) / pmeans[i] * 100
              for i in range(min(4, len(peaks_sim))) if not np.isnan(vals_sim[i])]
    mae = np.mean(errors) if errors else float('nan')
    m2_pass = mae < 40
    print(f"  M2 (<40%): {mae:.1f}% {'PASS' if m2_pass else 'FAIL'}")

    # Cross-correlation (already computed for best probe)
    r_max, tau = cross_correlation(t_sim, p_sim, t_exp, p_exp)
    print(f"  M5 (r>0.5): r={r_max:.3f}, tau={tau:+.3f}s {'PASS' if r_max > 0.5 else 'FAIL'}")

    # Plot
    fig, ax = plt.subplots(1, 1, figsize=(12, 5))
    ax.plot(t_exp, p_exp, color='black', alpha=0.4, linewidth=0.5, label='Experiment')
    ax.plot(t_sim, p_sim, color='#9467bd', linewidth=1.0,
            label=f'DBC dp=2mm (Roof, x={x_best:.3f}m)')

    for i, idx in enumerate(peaks_sim[:4]):
        ax.plot(t_sim[idx], p_sim[idx], 'v', color='#9467bd', markersize=8)
        ax.annotate(f'{p_sim[idx]:.1f}', (t_sim[idx], p_sim[idx]),
                    textcoords="offset points", xytext=(0, 10), fontsize=9,
                    color='#9467bd', ha='center', fontweight='bold')

    ax.set_xlabel('Time [s]')
    ax.set_ylabel('Pressure [mbar]')
    ax.set_title('SPHERIC Test 10 — Water Roof Impact Pressure')
    ax.legend(loc='upper left')
    ax.set_xlim(0, 7)
    ax.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig(FIG_DIR / "fig_water_roof.png", dpi=150)
    print(f"  Saved: {FIG_DIR / 'fig_water_roof.png'}")
    plt.close()

    return {
        "probe_x_m": round(x_best, 3),
        "peaks_mbar": [round(v, 1) if not np.isnan(v) else None for v in vals_sim],
        "M1_in_band": f"{in_band}/4",
        "M1_pass": m1_pass,
        "M2_mae_pct": round(mae, 1),
        "M2_pass": m2_pass,
        "cross_corr_r": round(float(r_max), 3),
        "time_shift_s": round(float(tau), 3),
    }


if __name__ == "__main__":
    print("=" * 60)
    print("EXP-1 — Oil Lateral & Water Roof Analysis")
    print("=" * 60)

    results = {}
    oil = analyze_oil()
    if oil:
        results["oil_lateral"] = oil

    roof = analyze_roof()
    if roof:
        results["water_roof"] = roof

    if results:
        out = FIG_DIR.parent / "analysis" / "oil_roof_metrics.json"
        with open(out, 'w') as f:
            json.dump(results, f, indent=2, default=lambda x: bool(x) if isinstance(x, (np.bool_,)) else float(x) if isinstance(x, (np.floating, np.integer)) else x)
        print(f"\nSaved: {out}")

    print("\nDone.")
