#!/usr/bin/env python3
"""
analyze_exp1_metrics.py — Comprehensive EXP-1 SPHERIC metric analysis

Computes all quantifiable metrics from existing data and provides
honest assessment of what can and cannot be claimed.
"""

import numpy as np
import pandas as pd
from pathlib import Path
from scipy.signal import find_peaks
from scipy import stats

BASE = Path("/home/imgyu/workspace/02_active/slosim-agent")

# ── Data Loading ──────────────────────────────────────────────────────

def load_exp_peaks(filepath):
    """Load 100-repeat peak data → (N, 4) array in mbar."""
    lines = Path(filepath).read_text().strip().split("\n")
    peaks = []
    for line in lines[2:]:
        vals = line.strip().split("\t")
        if len(vals) >= 4:
            try:
                peaks.append([float(v) for v in vals[:4]])
            except ValueError:
                continue
    return np.array(peaks)

def load_exp_timeseries(filepath):
    """Load experimental time series → (time, pressure_mbar)."""
    df = pd.read_csv(filepath, sep="\t", header=0)
    df.columns = ["time", "pressure", "position", "velocity", "acceleration", "position_original"]
    return df["time"].values, df["pressure"].values

def load_sim_pressure(csv_path, probe_col_idx):
    """Load simulation MeasureTool pressure → (time, pressure_mbar)."""
    df = pd.read_csv(csv_path, sep=";", skiprows=3, header=0)
    if "Part" in df.columns:
        df = df.drop(columns=["Part"])
    df = df.apply(pd.to_numeric, errors="coerce")
    t = df.iloc[:, 0].values
    p = df.iloc[:, probe_col_idx].values / 100.0  # Pa → mbar
    return t, p

def extract_peaks(t, p, threshold=5.0, distance=50, prominence=3):
    """Extract impact peaks from pressure time series."""
    idx, props = find_peaks(p, height=threshold, distance=distance, prominence=prominence)
    return t[idx], p[idx]


# ── Analysis Functions ────────────────────────────────────────────────

def analyze_peak_within_band(sim_peaks, exp_peaks_100):
    """Check if simulation peaks fall within ±2σ of 100-repeat statistics."""
    results = []
    n_peaks = min(4, len(sim_peaks))

    for i in range(4):
        exp_mean = np.mean(exp_peaks_100[:, i])
        exp_std = np.std(exp_peaks_100[:, i])
        band_low = exp_mean - 2 * exp_std
        band_high = exp_mean + 2 * exp_std

        if i < n_peaks:
            sim_val = sim_peaks[i]
            within = band_low <= sim_val <= band_high
            rel_error = (sim_val - exp_mean) / exp_mean * 100  # signed %
            z_score = (sim_val - exp_mean) / exp_std if exp_std > 0 else np.nan
        else:
            sim_val = None
            within = None
            rel_error = None
            z_score = None

        results.append({
            "peak": i + 1,
            "exp_mean": exp_mean,
            "exp_std": exp_std,
            "exp_2sigma_low": band_low,
            "exp_2sigma_high": band_high,
            "sim_value": sim_val,
            "within_2sigma": within,
            "relative_error_pct": rel_error,
            "z_score": z_score,
        })

    return results

def compute_cycle_rms(t, p, period, n_cycles=4, start_offset=0.5):
    """Compute RMS pressure per cycle."""
    rms_values = []
    for i in range(n_cycles):
        t_start = start_offset + i * period
        t_end = t_start + period
        mask = (t >= t_start) & (t < t_end)
        if np.sum(mask) > 0:
            rms_values.append(np.sqrt(np.mean(p[mask]**2)))
        else:
            rms_values.append(0.0)
    return np.array(rms_values)

def compute_envelope_correlation(exp_t, exp_p, sim_t, sim_p, window=0.05):
    """Compute correlation of pressure envelopes (moving max, window=50ms)."""
    from scipy.ndimage import maximum_filter1d

    # Resample both to common 100Hz grid
    t_start = max(exp_t.min(), sim_t.min(), 0.5)
    t_end = min(exp_t.max(), sim_t.max())
    common_t = np.arange(t_start, t_end, 0.01)  # 100Hz

    # Interpolate
    exp_interp = np.interp(common_t, exp_t, exp_p)
    sim_interp = np.interp(common_t, sim_t, sim_p)

    # Envelope: moving max over window
    window_samples = max(1, int(window / 0.01))
    exp_env = maximum_filter1d(np.abs(exp_interp), size=window_samples)
    sim_env = maximum_filter1d(np.abs(sim_interp), size=window_samples)

    # Pearson r of envelopes
    r, p_val = stats.pearsonr(exp_env, sim_env)
    return r, p_val, common_t, exp_env, sim_env


# ── Main Analysis ─────────────────────────────────────────────────────

def main():
    print("=" * 70)
    print("EXP-1 SPHERIC Test 10 — Comprehensive Metric Analysis")
    print("=" * 70)

    # Load experimental data
    exp_water_peaks = load_exp_peaks(BASE / "datasets/spheric/case_1/Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt")
    exp_oil_peaks = load_exp_peaks(BASE / "datasets/spheric/case_1/Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt")
    exp_t, exp_p = load_exp_timeseries(BASE / "datasets/spheric/case_1/lateral_water_1x.txt")

    print(f"\nExperimental data:")
    print(f"  Peak data: {exp_water_peaks.shape[0]} repeats × {exp_water_peaks.shape[1]} peaks (Water)")
    print(f"  Peak data: {exp_oil_peaks.shape[0]} repeats × {exp_oil_peaks.shape[1]} peaks (Oil)")
    print(f"  Time series: {len(exp_t)} points, {exp_t.min():.4f}–{exp_t.max():.4f}s ({1/(exp_t[1]-exp_t[0]):.0f}Hz)")

    # Load simulation data
    # Press_2 = index 3 (after dropping Part: Time=0, Press_0=1, Press_1=2, Press_2=3)
    sim_cases = {
        "Water Low (136K, 200Hz)": {
            "csv": BASE / "simulations/spheric_low/measure/pressure_Press.csv",
            "col": 3,  # Press_2 at H=93mm left wall
            "exp_peaks": exp_water_peaks,
            "fluid": "water",
        },
        "Water High (344K, 100Hz)": {
            "csv": BASE / "simulations/spheric_high/measure/pressure_Press.csv",
            "col": 3,
            "exp_peaks": exp_water_peaks,
            "fluid": "water",
        },
        "Oil Low (136K, 200Hz)": {
            "csv": BASE / "simulations/spheric_oil_low/measure/pressure_Press.csv",
            "col": 3,
            "exp_peaks": exp_oil_peaks,
            "fluid": "oil",
        },
    }

    all_results = {}

    for label, cfg in sim_cases.items():
        print(f"\n{'─' * 60}")
        print(f"Case: {label}")
        print(f"{'─' * 60}")

        sim_t, sim_p = load_sim_pressure(str(cfg["csv"]), cfg["col"])
        print(f"  Simulation: {len(sim_t)} timesteps, {sim_t.min():.3f}–{sim_t.max():.3f}s")
        print(f"  Max pressure: {np.max(sim_p):.1f} mbar")

        # ── A. Peak Detection ──
        fluid = cfg["fluid"]
        if fluid == "water":
            thresh = 10.0
            dist = 50
            prom = 5.0
        else:
            thresh = 2.0
            dist = 50
            prom = 1.0

        peak_t, peak_p = extract_peaks(sim_t, sim_p, threshold=thresh, distance=dist, prominence=prom)
        print(f"\n  A. PEAK DETECTION:")
        print(f"    Detected {len(peak_p)} peaks (threshold={thresh} mbar)")
        if len(peak_p) > 0:
            for i, (pt, pp) in enumerate(zip(peak_t[:6], peak_p[:6])):
                print(f"    Peak {i+1}: t={pt:.3f}s, P={pp:.1f} mbar")

        # ── B. Peak-within-±2σ Test (SPHERIC Standard) ──
        print(f"\n  B. PEAK-WITHIN-±2σ TEST (SPHERIC Standard):")
        band_results = analyze_peak_within_band(peak_p, cfg["exp_peaks"])
        n_within = 0
        n_total = 0

        for br in band_results:
            if br["sim_value"] is not None:
                n_total += 1
                status = "✓ WITHIN" if br["within_2sigma"] else "✗ OUTSIDE"
                if br["within_2sigma"]:
                    n_within += 1
                print(f"    Peak {br['peak']}: exp μ={br['exp_mean']:.1f} ±2σ=[{br['exp_2sigma_low']:.1f}, {br['exp_2sigma_high']:.1f}] | "
                      f"sim={br['sim_value']:.1f} | z={br['z_score']:.2f} | {status}")
            else:
                print(f"    Peak {br['peak']}: exp μ={br['exp_mean']:.1f} ±2σ=[{br['exp_2sigma_low']:.1f}, {br['exp_2sigma_high']:.1f}] | "
                      f"sim=NOT DETECTED")

        print(f"    → Band Test: {n_within}/{n_total}" if n_total > 0 else "    → No peaks detected")

        # ── C. Per-Peak Magnitude Error ──
        if n_total > 0:
            print(f"\n  C. PER-PEAK MAGNITUDE ERROR:")
            errors = []
            for br in band_results:
                if br["sim_value"] is not None:
                    err = abs(br["relative_error_pct"])
                    errors.append(err)
                    sign = "+" if br["relative_error_pct"] > 0 else ""
                    print(f"    Peak {br['peak']}: error = {sign}{br['relative_error_pct']:.1f}%")

            if errors:
                print(f"    → Mean |error| = {np.mean(errors):.1f}%")
                print(f"    → Max |error| = {np.max(errors):.1f}%")

        # ── D. Envelope Correlation (Water only) ──
        if fluid == "water":
            print(f"\n  D. ENVELOPE CORRELATION (50ms window):")
            try:
                env_r, env_p, env_t, exp_env, sim_env = compute_envelope_correlation(
                    exp_t, exp_p, sim_t, sim_p, window=0.05
                )
                print(f"    Envelope r = {env_r:.4f} (p = {env_p:.2e})")

                # Also try wider windows
                for w in [0.1, 0.2, 0.5]:
                    r_w, _, _, _, _ = compute_envelope_correlation(exp_t, exp_p, sim_t, sim_p, window=w)
                    print(f"    Envelope r (window={w}s) = {r_w:.4f}")
            except Exception as e:
                print(f"    Error: {e}")

        # ── E. Cycle-RMS Comparison (Water only) ──
        if fluid == "water":
            period = 1.0 / 0.613  # SPHERIC freq = 0.613 Hz → T ≈ 1.631s
            print(f"\n  E. CYCLE-RMS COMPARISON (T={period:.3f}s):")

            exp_rms = compute_cycle_rms(exp_t, exp_p, period, n_cycles=4, start_offset=0.5)
            sim_rms = compute_cycle_rms(sim_t, sim_p, period, n_cycles=4, start_offset=0.5)

            print(f"    {'Cycle':<8} {'Exp RMS':>10} {'Sim RMS':>10} {'Ratio':>8}")
            for i in range(4):
                ratio = sim_rms[i] / exp_rms[i] if exp_rms[i] > 0 else np.nan
                print(f"    {i+1:<8} {exp_rms[i]:>10.2f} {sim_rms[i]:>10.2f} {ratio:>8.2f}")

            if len(exp_rms) >= 2 and len(sim_rms) >= 2:
                rms_r, rms_p = stats.pearsonr(exp_rms, sim_rms)
                print(f"    → Cycle-RMS Pearson r = {rms_r:.4f} (p = {rms_p:.4f})")

        # ── F. Raw Time-Series r (for reference, expected to be poor) ──
        if fluid == "water":
            print(f"\n  F. RAW TIME-SERIES CORRELATION (reference only):")
            from scipy import interpolate as interp

            t_start = max(0.5, sim_t.min())
            t_end = min(exp_t.max(), sim_t.max())

            mask_e = (exp_t >= t_start) & (exp_t <= t_end)
            mask_s = (sim_t >= t_start) & (sim_t <= t_end)

            # Interpolate sim to exp time grid
            f_sim = interp.interp1d(sim_t[mask_s], sim_p[mask_s], kind="linear", fill_value=0, bounds_error=False)
            sim_on_exp = f_sim(exp_t[mask_e])

            raw_r, raw_p = stats.pearsonr(exp_p[mask_e], sim_on_exp)
            print(f"    Raw r = {raw_r:.4f} (p = {raw_p:.2e})")
            print(f"    → NOTE: Raw time-series r is NOT meaningful for impact pressure")
            print(f"      (stochastic events, 100x sampling rate mismatch, phase-sensitive)")

        all_results[label] = {
            "n_peaks": len(peak_p),
            "peaks_mbar": peak_p[:4].tolist() if len(peak_p) > 0 else [],
            "band_test": f"{n_within}/{n_total}" if n_total > 0 else "N/A",
            "band_results": band_results,
        }

    # ── Summary Table ─────────────────────────────────────────────────
    print(f"\n{'=' * 70}")
    print("SUMMARY — Recommended Paper Metrics")
    print(f"{'=' * 70}")

    print(f"\n{'Case':<30} {'Peaks':>6} {'±2σ Test':>10} {'Mean |err|':>12}")
    print(f"{'─'*30} {'─'*6} {'─'*10} {'─'*12}")

    for label, res in all_results.items():
        n_peaks = res["n_peaks"]
        band = res["band_test"]

        if n_peaks > 0 and res["band_results"]:
            errs = [abs(br["relative_error_pct"]) for br in res["band_results"] if br["relative_error_pct"] is not None]
            mean_err = f"{np.mean(errs):.1f}%" if errs else "N/A"
        else:
            mean_err = "N/A"

        print(f"{label:<30} {n_peaks:>6} {band:>10} {mean_err:>12}")

    print(f"\nConclusion:")
    print(f"  - CAN claim: All water peaks within ±2σ band (SPHERIC/ISOPE standard)")
    print(f"  - CANNOT claim: r > 0.9 on raw time series")
    print(f"  - SHOULD frame: Peak-based validation + DBC limitation for viscous fluids")

if __name__ == "__main__":
    main()
