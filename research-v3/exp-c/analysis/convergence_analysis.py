#!/usr/bin/env python3
"""EXP-1 수렴 분석: Richardson Extrapolation + GCI
   Run 001 (dp=0.004), Run 002 (dp=0.002), Run 003 (dp=0.001)
   ASME V&V 20-2009 Grid Convergence Index

   Run 003 note: TimeOut=0.01 (100Hz) limits peak capture. Peak-based metrics
   use 2-level GCI (Run 001+002) with Fs=3.0. Cross-correlation uses all 3 runs.

Usage:
    python convergence_analysis.py [--no-run003]
"""

import numpy as np
import matplotlib.pyplot as plt
from pathlib import Path
from scipy.signal import find_peaks, medfilt
from scipy.ndimage import uniform_filter1d
import sys
import json

ROOT = Path(__file__).resolve().parents[3]  # slosim-agent/
SIM_DIR = ROOT / "simulations" / "exp1"
EXP_DIR = ROOT / "datasets" / "spheric" / "case_1"
FIG_DIR = Path(__file__).resolve().parent.parent / "figures"
FIG_DIR.mkdir(exist_ok=True)

# --- Configuration ---
DPS = {1: 0.004, 2: 0.002, 3: 0.001}
REFINEMENT_RATIO = 2.0
N_PEAKS = 3  # 4th peak at t>7s, beyond simulation TimeMax=7.0

# SPH pressure filtering: median filter removes single-point spikes,
# then moving average smooths remaining noise.
MEDFILT_KERNEL = 5     # median filter window (removes isolated spikes)
SMOOTH_WINDOW = 11     # moving average window for smoothing (~10ms at 1kHz)


# --- Data loading ---
def load_gauge_pressure(run_id, apply_filter=True):
    """Load pressure from MeasureTool CSV (semicolon-delimited, 4 header lines).
    Prefers PressConsistent_Press.csv (fixed probe location x=0.005) over legacy data.
    """
    path = SIM_DIR / f"run_{run_id:03d}" / "PressConsistent_Press.csv"
    if not path.exists():
        path = SIM_DIR / f"run_{run_id:03d}" / "PressureLateral_Press.csv"
    if not path.exists():
        return None, None
    data = np.genfromtxt(path, delimiter=';', skip_header=4)
    t = data[:, 1]          # Time [s]
    p = data[:, 2] / 100    # Press_0 [Pa] → [mbar]

    if apply_filter:
        # 1) Median filter: removes single-point SPH noise spikes
        p = medfilt(p, kernel_size=MEDFILT_KERNEL)
        # 2) Moving average: smooths residual noise
        p = uniform_filter1d(p, size=SMOOTH_WINDOW)

    return t, p


def load_experimental():
    """Load experimental time series and peak statistics."""
    exp_path = EXP_DIR / "lateral_water_1x.txt"
    exp_data = np.genfromtxt(exp_path, delimiter='\t', skip_header=1)
    t_exp = exp_data[:, 0]
    p_exp = exp_data[:, 1]  # mbar

    peak_path = EXP_DIR / "Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt"
    peak_data = np.genfromtxt(peak_path, delimiter='\t', skip_header=2)
    peak_means = np.nanmean(peak_data, axis=0)
    peak_stds = np.nanstd(peak_data, axis=0)
    return t_exp, p_exp, peak_means, peak_stds


# --- Peak detection ---
def detect_peaks(t, p, min_height=10.0, min_dist_s=1.0):
    """Detect pressure peaks with minimum height and distance."""
    dt = np.median(np.diff(t))
    min_dist = int(min_dist_s / dt) if dt > 0 else 1000
    peaks, props = find_peaks(p, height=min_height, distance=min_dist)
    return peaks


def get_peak_vals(t, p, peaks, n=3):
    """Extract first n peak values."""
    return [(t[peaks[i]], p[peaks[i]]) if i < len(peaks) else (np.nan, np.nan)
            for i in range(n)]


# --- GCI calculation (ASME V&V 20-2009) ---
def compute_gci(f1, f2, f3=None, r=2.0, Fs=None):
    """Compute Grid Convergence Index.

    3-grid: f1=coarse, f2=medium, f3=fine. Fs=1.25 (default).
    2-grid: f1=coarse, f2=fine, f3=None. Fs=3.0 (conservative), assumed p=2.
    """
    if f3 is not None:
        # 3-level GCI
        if Fs is None:
            Fs = 1.25
        eps32 = f3 - f2
        eps21 = f2 - f1
        if abs(eps32) < 1e-12 or abs(eps21) < 1e-12:
            return np.nan, f3, np.nan, np.nan
        p = np.log(abs(eps21 / eps32)) / np.log(r)
        f_exact = f3 + eps32 / (r**p - 1)
        e_fine = abs(eps32 / f3) if abs(f3) > 1e-12 else np.nan
        GCI_fine = Fs * e_fine / (r**p - 1) * 100
        return p, f_exact, GCI_fine, None
    else:
        # 2-level GCI: assume p=2 (Symplectic is 2nd order), Fs=3.0
        if Fs is None:
            Fs = 3.0
        p = 2.0  # assumed order
        eps21 = f2 - f1
        e_fine = abs(eps21 / f2) if abs(f2) > 1e-12 else np.nan
        f_exact = f2 + eps21 / (r**p - 1)
        GCI_fine = Fs * e_fine / (r**p - 1) * 100
        return p, f_exact, GCI_fine, None


# --- Cross-correlation with time shift ---
def cross_correlation(t_sim, p_sim, t_exp, p_exp, max_shift_s=2.0):
    """Compute cross-correlation with optimal time shift."""
    dt_sim = np.median(np.diff(t_sim))
    dt_exp = np.median(np.diff(t_exp))

    dt = max(dt_sim, dt_exp)
    t_common = np.arange(0, min(t_sim[-1], t_exp[-1]), dt)
    p_sim_interp = np.interp(t_common, t_sim, p_sim)
    p_exp_interp = np.interp(t_common, t_exp, p_exp)

    p_sim_n = (p_sim_interp - np.mean(p_sim_interp)) / (np.std(p_sim_interp) + 1e-12)
    p_exp_n = (p_exp_interp - np.mean(p_exp_interp)) / (np.std(p_exp_interp) + 1e-12)

    n = len(t_common)
    corr = np.correlate(p_sim_n, p_exp_n, mode='full') / n
    lags = np.arange(-n + 1, n) * dt

    mask = np.abs(lags) <= max_shift_s
    corr_masked = corr[mask]
    lags_masked = lags[mask]

    idx_max = np.argmax(corr_masked)
    return corr_masked[idx_max], lags_masked[idx_max]


# === MAIN ===
def main():
    no_run003 = '--no-run003' in sys.argv

    # Load data (with filtering)
    t001, p001 = load_gauge_pressure(1)
    t002, p002 = load_gauge_pressure(2)
    t003, p003 = (None, None) if no_run003 else load_gauge_pressure(3)

    # Also load unfiltered for plot overlay
    _, p001_raw = load_gauge_pressure(1, apply_filter=False)
    _, p002_raw = load_gauge_pressure(2, apply_filter=False)

    t_exp, p_exp, peak_means, peak_stds = load_experimental()

    # Detect peaks (on filtered data)
    peaks_001 = detect_peaks(t001, p001, min_height=10.0, min_dist_s=1.2)
    peaks_002 = detect_peaks(t002, p002, min_height=10.0, min_dist_s=1.2)
    peaks_exp = detect_peaks(t_exp, p_exp, min_height=10.0, min_dist_s=0.8)
    peaks_003 = detect_peaks(t003, p003, min_height=10.0, min_dist_s=1.2) if t003 is not None else np.array([])

    # Peak values (time, pressure)
    pv_001 = get_peak_vals(t001, p001, peaks_001, N_PEAKS)
    pv_002 = get_peak_vals(t002, p002, peaks_002, N_PEAKS)
    pv_003 = get_peak_vals(t003, p003, peaks_003, N_PEAKS) if t003 is not None else [(np.nan, np.nan)] * N_PEAKS
    pv_exp = get_peak_vals(t_exp, p_exp, peaks_exp, N_PEAKS)

    vals_001 = [v[1] for v in pv_001]
    vals_002 = [v[1] for v in pv_002]
    vals_003 = [v[1] for v in pv_003]

    # Cross-correlation
    r001, tau001 = cross_correlation(t001, p001, t_exp, p_exp)
    r002, tau002 = cross_correlation(t002, p002, t_exp, p_exp)
    r003, tau003 = cross_correlation(t003, p003, t_exp, p_exp) if t003 is not None else (np.nan, np.nan)

    # === Print results ===
    print("=" * 80)
    print("EXP-1 SPHERIC Test 10 — Convergence Analysis (Filtered)")
    print(f"  Median filter: kernel={MEDFILT_KERNEL}, Moving avg: window={SMOOTH_WINDOW}")
    print(f"  Comparing first {N_PEAKS} peaks (4th peak at t>7s, beyond TimeMax)")
    print("=" * 80)

    print(f"\n{'Peak':<8} {'Exp Mean':>10} {'±2σ':>8} {'R001(4mm)':>12} {'R002(2mm)':>12} {'R003(1mm)':>12}")
    print("-" * 62)
    for i in range(N_PEAKS):
        v1 = f"{vals_001[i]:>12.1f}" if not np.isnan(vals_001[i]) else f"{'N/A':>12}"
        v2 = f"{vals_002[i]:>12.1f}" if not np.isnan(vals_002[i]) else f"{'N/A':>12}"
        v3 = f"{vals_003[i]:>12.1f}" if not np.isnan(vals_003[i]) else f"{'N/A':>12}"
        print(f"Peak {i+1:<3} {peak_means[i]:>10.1f} {2*peak_stds[i]:>8.1f} {v1} {v2} {v3}")

    # Peak timing
    print(f"\n{'Peak':<8} {'Exp t[s]':>10} {'R001 t[s]':>12} {'R002 t[s]':>12} {'R003 t[s]':>12}")
    print("-" * 58)
    for i in range(N_PEAKS):
        te = f"{pv_exp[i][0]:>10.3f}" if not np.isnan(pv_exp[i][0]) else f"{'N/A':>10}"
        t1 = f"{pv_001[i][0]:>12.3f}" if not np.isnan(pv_001[i][0]) else f"{'N/A':>12}"
        t2 = f"{pv_002[i][0]:>12.3f}" if not np.isnan(pv_002[i][0]) else f"{'N/A':>12}"
        t3 = f"{pv_003[i][0]:>12.3f}" if not np.isnan(pv_003[i][0]) else f"{'N/A':>12}"
        print(f"Peak {i+1:<3} {te} {t1} {t2} {t3}")

    # Error table
    print(f"\n{'Metric':<20} {'R001(4mm)':>12} {'R002(2mm)':>12} {'R003(1mm)':>12}")
    print("-" * 58)
    for i in range(N_PEAKS):
        def err(v):
            return f"{(v - peak_means[i]) / peak_means[i] * 100:+.1f}%" if not np.isnan(v) else "N/A"
        print(f"Peak {i+1} error     {err(vals_001[i]):>12} {err(vals_002[i]):>12} {err(vals_003[i]):>12}")

    # === M5: Cross-Correlation ===
    print(f"\n--- M5: Cross-Correlation (threshold: r>0.5) ---")
    print(f"{'Run':<20} {'r_max':>10} {'tau [s]':>10} {'Result':>10}")
    for name, r, tau in [("Run 001 (dp=0.004)", r001, tau001),
                          ("Run 002 (dp=0.002)", r002, tau002),
                          ("Run 003 (dp=0.001)", r003, tau003)]:
        if np.isnan(r):
            continue
        print(f"{name:<20} {r:>10.3f} {tau:>10.3f} {'PASS' if r > 0.5 else 'FAIL':>10}")

    # === M6: Time Shift ===
    period = 1.0 / 0.613  # 1.631s
    print(f"\n--- M6: Time Shift (threshold: |tau| < {period:.2f}s = 1 period) ---")
    for name, tau in [("Run 001", tau001), ("Run 002", tau002), ("Run 003", tau003)]:
        if np.isnan(tau):
            continue
        print(f"  {name}: tau={tau:+.3f}s {'PASS' if abs(tau) < period else 'FAIL'}")

    # === GCI ===
    gci_pass = None
    mono_pass = None

    # Check if Run 003 has valid peaks
    run003_peaks_valid = (t003 is not None and not np.isnan(vals_003[0])
                          and vals_003[0] > 10.0)  # must capture actual impact peaks

    if run003_peaks_valid:
        # 3-level GCI
        print(f"\n--- M3/M4: Grid Convergence Index (3-level) ---")
        print(f"{'Peak':<8} {'p (order)':>10} {'f_exact':>10} {'GCI_fine%':>10} {'Monotone':>10}")
        gci_pass = True
        mono_pass = True
        for i in range(N_PEAKS):
            f1, f2, f3 = vals_001[i], vals_002[i], vals_003[i]
            if np.isnan(f1) or np.isnan(f2) or np.isnan(f3):
                print(f"Peak {i+1:<3} {'N/A':>10} {'N/A':>10} {'N/A':>10} {'N/A':>10}")
                continue
            p, f_ex, gci_f, _ = compute_gci(f1, f2, f3, r=REFINEMENT_RATIO)
            mono = abs(f3 - f2) < abs(f2 - f1)
            if not np.isnan(gci_f) and gci_f >= 10:
                gci_pass = False
            if not mono:
                mono_pass = False
            print(f"Peak {i+1:<3} {p:>10.2f} {f_ex:>10.1f} {gci_f:>10.1f} {'YES' if mono else 'NO':>10}")
        print(f"\n  M3 Monotone convergence: {'PASS' if mono_pass else 'FAIL'}")
        print(f"  M4 GCI_fine < 10%:       {'PASS' if gci_pass else 'FAIL'}")
    else:
        # 2-level GCI (Run 001 + Run 002, assumed p=2, Fs=3.0)
        print(f"\n--- M3/M4: Grid Convergence Index (2-level, Fs=3.0, assumed p=2) ---")
        if t003 is not None:
            print(f"  Note: Run 003 TimeOut=0.01s (100Hz) → peaks not resolved. Using 2-level GCI.")
        print(f"{'Peak':<8} {'f_coarse':>10} {'f_fine':>10} {'f_exact':>10} {'GCI_fine%':>10} {'Monotone':>10}")
        gci_pass = True
        mono_pass = True
        for i in range(N_PEAKS):
            f1, f2 = vals_001[i], vals_002[i]
            if np.isnan(f1) or np.isnan(f2):
                print(f"Peak {i+1:<3} {'N/A':>10} {'N/A':>10} {'N/A':>10} {'N/A':>10} {'N/A':>10}")
                continue
            p, f_ex, gci_f, _ = compute_gci(f1, f2, r=REFINEMENT_RATIO)
            # Monotone: error decreases from coarse to fine (vs experiment)
            err_coarse = abs(f1 - peak_means[i])
            err_fine = abs(f2 - peak_means[i])
            mono = err_fine < err_coarse
            if not np.isnan(gci_f) and gci_f >= 10:
                gci_pass = False
            if not mono:
                mono_pass = False
            print(f"Peak {i+1:<3} {f1:>10.1f} {f2:>10.1f} {f_ex:>10.1f} {gci_f:>10.1f} {'YES' if mono else 'NO':>10}")
        # Cross-correlation convergence as supplementary evidence
        if t003 is not None:
            print(f"\n  Supplementary: Cross-correlation convergence")
            print(f"  Run 001 r={r001:.3f} → Run 002 r={r002:.3f} → Run 003 r={r003:.3f}")
            if r003 > r002 > r001:
                print(f"  Monotone improvement in waveform correlation: YES")
        print(f"\n  M3 Monotone convergence: {'PASS' if mono_pass else 'FAIL'}")
        print(f"  M4 GCI_fine < 10%:       {'PASS' if gci_pass else 'FAIL'} (conservative Fs=3.0)")

    # === M1: Peak-in-band ===
    print(f"\n--- M1: Peak-in-band (threshold: >={N_PEAKS-1}/{N_PEAKS} within exp mean±2σ) ---")
    m1_results = {}
    for run_name, vals in [("Run 001", vals_001), ("Run 002", vals_002), ("Run 003", vals_003)]:
        if np.isnan(vals[0]):
            continue
        in_band = 0
        for i in range(N_PEAKS):
            if np.isnan(vals[i]):
                continue
            lo = peak_means[i] - 2 * peak_stds[i]
            hi = peak_means[i] + 2 * peak_stds[i]
            if lo <= vals[i] <= hi:
                in_band += 1
        passes = in_band >= N_PEAKS - 1
        m1_results[run_name] = passes
        print(f"  {run_name}: {in_band}/{N_PEAKS} {'PASS' if passes else 'FAIL'}")

    # === M2: Mean absolute peak error ===
    print(f"\n--- M2: Mean Absolute Peak Error (threshold: <30%) ---")
    m2_results = {}
    for run_name, vals in [("Run 001", vals_001), ("Run 002", vals_002), ("Run 003", vals_003)]:
        if np.isnan(vals[0]):
            continue
        errors = [abs(vals[i] - peak_means[i]) / peak_means[i] * 100
                  for i in range(N_PEAKS) if not np.isnan(vals[i])]
        mae = np.mean(errors) if errors else np.nan
        passes = mae < 30
        m2_results[run_name] = (mae, passes)
        print(f"  {run_name}: {mae:.1f}% {'PASS' if passes else 'FAIL'}")

    # === OVERALL VERDICT ===
    print(f"\n{'='*80}")
    print("OVERALL VERDICT — Water Lateral")
    print(f"{'='*80}")
    m1_str = ', '.join(f'{k}: {"PASS" if v else "FAIL"}' for k, v in m1_results.items())
    print(f"  M1 (Peak-in-band):   {m1_str}")
    for k, (mae, v) in m2_results.items():
        result = "PASS" if v else "FAIL"
        print(f"  M2 (Peak error):     {k}: {mae:.1f}% {result}")
    m5_parts = []
    for name, r_val in [("Run 001", r001), ("Run 002", r002), ("Run 003", r003)]:
        if not np.isnan(r_val):
            m5_parts.append(f'{name}: r={r_val:.3f} {"PASS" if r_val > 0.5 else "FAIL"}')
    print(f"  M5 (Cross-corr):     {', '.join(m5_parts)}")
    best_tau = tau003 if not np.isnan(tau003) else tau002
    print(f"  M6 (Time shift):     tau ≈ {best_tau:+.3f}s (< 1 period) {'PASS' if abs(best_tau) < period else 'FAIL'}")
    if gci_pass is not None:
        print(f"  M3 (Monotone):       {'PASS' if mono_pass else 'FAIL'}")
        print(f"  M4 (GCI_fine<10%):   {'PASS' if gci_pass else 'FAIL'}")

    # === PLOT: Publication quality ===
    fig = plt.figure(figsize=(16, 14))
    gs = fig.add_gridspec(3, 2, height_ratios=[3, 1.5, 1.5], hspace=0.35, wspace=0.3)

    # Panel A: Full time series comparison
    ax1 = fig.add_subplot(gs[0, :])
    ax1.plot(t_exp, p_exp, color='black', alpha=0.3, linewidth=0.3, label='Experimental (20kHz)')
    ax1.plot(t001, p001, color='#1f77b4', linewidth=1.0, alpha=0.8,
             label=f'Run 001 (dp=4mm, DBC)')
    ax1.plot(t002, p002, color='#d62728', linewidth=1.0, alpha=0.8,
             label=f'Run 002 (dp=2mm, DBC)')
    if t003 is not None:
        ax1.plot(t003, p003, color='#2ca02c', linewidth=1.2,
                 label=f'Run 003 (dp=1mm, DBC)')

    # Experimental peak markers
    for i in range(min(N_PEAKS, len(peaks_exp))):
        idx = peaks_exp[i]
        ax1.axvline(t_exp[idx], color='gray', alpha=0.3, linestyle='--', linewidth=0.5)

    # Simulation peak markers (best available run)
    best_name = "R003" if t003 is not None else "R002"
    best_pv = pv_003 if t003 is not None else pv_002
    best_color = '#2ca02c' if t003 is not None else '#d62728'
    for i in range(N_PEAKS):
        tp, pp = best_pv[i]
        if not np.isnan(tp):
            ax1.plot(tp, pp, 'v', color=best_color, markersize=10, zorder=5)
            ax1.annotate(f'{pp:.1f}', (tp, pp), textcoords="offset points",
                        xytext=(0, 12), fontsize=9, color=best_color,
                        ha='center', fontweight='bold')

    ax1.set_xlabel('Time [s]', fontsize=12)
    ax1.set_ylabel('Pressure [mbar]', fontsize=12)
    ax1.set_title('SPHERIC Test 10 — Lateral Water Impact Pressure (Filtered)',
                   fontsize=14, fontweight='bold')
    ax1.legend(loc='upper left', fontsize=10, framealpha=0.9)
    ax1.set_xlim(0, 7)
    ymax = max(p001.max(), p002.max(), p_exp.max())
    if t003 is not None:
        ymax = max(ymax, p003.max())
    ax1.set_ylim(-15, min(ymax * 1.3, 120))
    ax1.grid(True, alpha=0.3)
    ax1.text(0.02, 0.95, '(a)', transform=ax1.transAxes, fontsize=14, fontweight='bold', va='top')

    # Panel B: Peak bar chart
    ax2 = fig.add_subplot(gs[1, 0])
    x = np.arange(N_PEAKS)
    width = 0.18

    ax2.bar(x - 1.5*width, peak_means[:N_PEAKS], width,
            yerr=2*peak_stds[:N_PEAKS], label='Exp (mean±2σ)',
            color='gray', alpha=0.7, capsize=4)
    ax2.bar(x - 0.5*width, vals_001, width, label='R001 (4mm)',
            color='#1f77b4', alpha=0.7)
    ax2.bar(x + 0.5*width, vals_002, width, label='R002 (2mm)',
            color='#d62728', alpha=0.7)
    if t003 is not None:
        ax2.bar(x + 1.5*width, vals_003, width, label='R003 (1mm)',
                color='#2ca02c', alpha=0.7)

    ax2.set_xlabel('Impact Peak', fontsize=11)
    ax2.set_ylabel('Pressure [mbar]', fontsize=11)
    ax2.set_title('Peak Pressure Comparison', fontsize=12, fontweight='bold')
    ax2.set_xticks(x)
    ax2.set_xticklabels([f'Peak {i+1}' for i in range(N_PEAKS)])
    ax2.legend(fontsize=8, loc='upper right')
    ax2.grid(True, alpha=0.3, axis='y')
    ax2.text(0.02, 0.95, '(b)', transform=ax2.transAxes, fontsize=14, fontweight='bold', va='top')

    # Panel C: Peak error vs resolution
    ax3 = fig.add_subplot(gs[1, 1])
    dps = [4.0, 2.0]
    colors_line = ['#1f77b4', '#d62728']
    if t003 is not None:
        dps.append(1.0)
        colors_line.append('#2ca02c')

    for i in range(N_PEAKS):
        errs = []
        for vals in [vals_001, vals_002] + ([vals_003] if t003 is not None else []):
            if not np.isnan(vals[i]):
                errs.append(abs(vals[i] - peak_means[i]) / peak_means[i] * 100)
            else:
                errs.append(np.nan)
        ax3.plot(dps[:len(errs)], errs, 'o-', label=f'Peak {i+1}', markersize=6)

    ax3.axhline(30, color='red', linestyle='--', alpha=0.5, label='M2 threshold (30%)')
    ax3.set_xlabel('dp [mm]', fontsize=11)
    ax3.set_ylabel('Absolute Error [%]', fontsize=11)
    ax3.set_title('Convergence: Error vs Resolution', fontsize=12, fontweight='bold')
    ax3.legend(fontsize=8)
    ax3.set_xlim(max(dps) * 1.2, min(dps) * 0.8)  # reversed x-axis (finer → right)
    ax3.grid(True, alpha=0.3)
    ax3.text(0.02, 0.95, '(c)', transform=ax3.transAxes, fontsize=14, fontweight='bold', va='top')

    # Panel D: Zoomed first impact (t ≈ 2.5-3.5s)
    ax4 = fig.add_subplot(gs[2, 0])
    t_start, t_end = 2.0, 4.0
    mask_exp = (t_exp >= t_start) & (t_exp <= t_end)
    mask_001 = (t001 >= t_start) & (t001 <= t_end)
    mask_002 = (t002 >= t_start) & (t002 <= t_end)

    ax4.plot(t_exp[mask_exp], p_exp[mask_exp], color='black', alpha=0.4, linewidth=0.5,
             label='Exp')
    ax4.plot(t001[mask_001], p001[mask_001], color='#1f77b4', linewidth=1.0,
             alpha=0.8, label='R001 (4mm)')
    ax4.plot(t002[mask_002], p002[mask_002], color='#d62728', linewidth=1.0,
             alpha=0.8, label='R002 (2mm)')
    if t003 is not None:
        mask_003 = (t003 >= t_start) & (t003 <= t_end)
        ax4.plot(t003[mask_003], p003[mask_003], color='#2ca02c', linewidth=1.2,
                 label='R003 (1mm)')

    ax4.set_xlabel('Time [s]', fontsize=11)
    ax4.set_ylabel('Pressure [mbar]', fontsize=11)
    ax4.set_title('First Impact (Zoomed)', fontsize=12, fontweight='bold')
    ax4.legend(fontsize=8, loc='upper right')
    ax4.set_xlim(t_start, t_end)
    ax4.grid(True, alpha=0.3)
    ax4.text(0.02, 0.95, '(d)', transform=ax4.transAxes, fontsize=14, fontweight='bold', va='top')

    # Panel E: Cross-correlation summary
    ax5 = fig.add_subplot(gs[2, 1])
    runs = ['R001\n(4mm)', 'R002\n(2mm)']
    r_vals = [r001, r002]
    tau_vals = [tau001, tau002]
    bar_colors = ['#1f77b4', '#d62728']
    if t003 is not None:
        runs.append('R003\n(1mm)')
        r_vals.append(r003)
        tau_vals.append(tau003)
        bar_colors.append('#2ca02c')

    bars = ax5.bar(runs, r_vals, color=bar_colors, alpha=0.7, edgecolor='black', linewidth=0.5)
    ax5.axhline(0.5, color='red', linestyle='--', alpha=0.5, label='M5 threshold (r=0.5)')
    for bar, tau in zip(bars, tau_vals):
        ax5.text(bar.get_x() + bar.get_width()/2, bar.get_height() + 0.02,
                f'τ={tau:+.2f}s', ha='center', fontsize=9, fontweight='bold')
    ax5.set_ylabel('Cross-correlation r_max', fontsize=11)
    ax5.set_title('Signal Correlation & Time Shift', fontsize=12, fontweight='bold')
    ax5.set_ylim(0, 1.0)
    ax5.legend(fontsize=8)
    ax5.grid(True, alpha=0.3, axis='y')
    ax5.text(0.02, 0.95, '(e)', transform=ax5.transAxes, fontsize=14, fontweight='bold', va='top')

    plt.savefig(FIG_DIR / "convergence_study.png", dpi=150, bbox_inches='tight')
    print(f"\nSaved: {FIG_DIR / 'convergence_study.png'}")

    # === Save metrics JSON ===
    metrics = {
        "experiment": "EXP-1 SPHERIC Test 10 — Water Lateral",
        "filter": {"median_kernel": MEDFILT_KERNEL, "smooth_window": SMOOTH_WINDOW},
        "n_peaks_compared": N_PEAKS,
        "runs": {}
    }
    for run_name, vals, r_val, tau_val in [
        ("run_001", vals_001, r001, tau001),
        ("run_002", vals_002, r002, tau002),
        ("run_003", vals_003, r003, tau003)
    ]:
        if np.isnan(vals[0]):
            continue
        errors = [abs(vals[i] - peak_means[i]) / peak_means[i] * 100
                  for i in range(N_PEAKS) if not np.isnan(vals[i])]
        in_band = sum(1 for i in range(N_PEAKS)
                      if not np.isnan(vals[i])
                      and peak_means[i] - 2*peak_stds[i] <= vals[i] <= peak_means[i] + 2*peak_stds[i])
        metrics["runs"][run_name] = {
            "dp_mm": {"run_001": 4.0, "run_002": 2.0, "run_003": 1.0}[run_name],
            "peaks_mbar": [round(v, 1) if not np.isnan(v) else None for v in vals],
            "peak_errors_pct": [round(e, 1) for e in errors],
            "mean_abs_error_pct": round(np.mean(errors), 1) if errors else None,
            "in_band": f"{in_band}/{N_PEAKS}",
            "cross_corr_r": round(float(r_val), 3),
            "time_shift_s": round(float(tau_val), 3),
            "M1_pass": bool(in_band >= N_PEAKS - 1),
            "M2_pass": bool(np.mean(errors) < 30) if errors else False,
            "M5_pass": bool(float(r_val) > 0.5),
            "M6_pass": bool(abs(float(tau_val)) < period)
        }
    metrics["exp_peak_means_mbar"] = [round(float(v), 1) for v in peak_means[:N_PEAKS]]
    metrics["exp_peak_2sigma_mbar"] = [round(float(2*v), 1) for v in peak_stds[:N_PEAKS]]

    # GCI results
    if gci_pass is not None:
        metrics["gci"] = {
            "n_levels": 3 if run003_peaks_valid else 2,
            "safety_factor": 1.25 if run003_peaks_valid else 3.0,
            "M3_monotone_pass": bool(mono_pass),
            "M4_gci_pass": bool(gci_pass),
        }
        if not run003_peaks_valid and t003 is not None:
            metrics["gci"]["note"] = "Run 003 TimeOut=0.01 limits peak capture; 2-level GCI used"

    metrics_path = FIG_DIR.parent / "analysis" / "metrics.json"
    with open(metrics_path, 'w') as f:
        json.dump(metrics, f, indent=2)
    print(f"Saved: {metrics_path}")


if __name__ == "__main__":
    main()
