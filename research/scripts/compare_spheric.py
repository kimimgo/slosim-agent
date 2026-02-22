#!/usr/bin/env python3
"""
compare_spheric.py — SPHERIC Test 10 실험 데이터 vs 시뮬레이션 결과 비교

EXP-1: SPHERIC Benchmark Reproduction (RQ2)
- 실험 데이터: datasets/spheric/case_1/ (100-repeat, 20kHz, Water + Oil)
- 시뮬레이션: simulations/spheric_{low,high,oil_low}/measure/
- 출력: research/experiments/exp1_spheric/

Usage:
    python research/scripts/compare_spheric.py --sim-dir simulations/spheric_low/measure \
        --exp-file datasets/spheric/case_1/lateral_water_1x.txt \
        --peak-file datasets/spheric/case_1/Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt \
        --out-dir research/experiments/exp1_spheric \
        --label "Low Resolution (dp=0.006)"
"""

import argparse
import os
import sys
from pathlib import Path

import numpy as np
import pandas as pd
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from scipy import stats, interpolate


# ── Data Loading ───────────────────────────────────────────────────────

def load_spheric_timeseries(filepath: str) -> pd.DataFrame:
    """Load SPHERIC Test 10 experimental time series.

    Columns: Time[s], Pressure[mbar], Position[deg], Velocity[deg/s],
             Acceleration[deg/s2], Position_original[deg]
    """
    df = pd.read_csv(filepath, sep="\t", header=0)
    df.columns = ["time", "pressure", "position", "velocity", "acceleration", "position_original"]
    return df


def load_spheric_peaks(filepath: str) -> np.ndarray:
    """Load SPHERIC 100-repeat peak pressure data.

    Returns: (N_repeats, 4) array of peak pressures in mbar.
    """
    lines = Path(filepath).read_text().strip().split("\n")
    # Skip header lines (first 2 lines: title + column names)
    data_lines = lines[2:]
    peaks = []
    for line in data_lines:
        vals = line.strip().split("\t")
        if len(vals) >= 4:
            try:
                peaks.append([float(v) for v in vals[:4]])
            except ValueError:
                continue
    return np.array(peaks)


def load_simulation_pressure(filepath: str) -> pd.DataFrame:
    """Load MeasureTool pressure output (CSV).

    MeasureTool format:
      Row 0-2: Position metadata (PosX/PosY/PosZ, starting with space)
      Row 3: Header (Part;Time [s];Press_0 [Pa];...)
      Row 4+: Data (semicolon-separated)
    """
    lines = Path(filepath).read_text().strip().split("\n")

    # Find header row: first row containing "Time" or "Part"
    header_idx = 0
    for i, line in enumerate(lines):
        if "Time" in line or "Part" in line:
            header_idx = i
            break

    # Try semicolon first (DualSPHysics standard)
    try:
        df = pd.read_csv(filepath, sep=";", skiprows=header_idx, header=0)
        # Drop non-numeric columns or "Part" column if present
        if "Part" in df.columns:
            df = df.drop(columns=["Part"])
        # Clean column names
        df.columns = [c.strip() for c in df.columns]
        # Ensure all columns are numeric
        df = df.apply(pd.to_numeric, errors="coerce")
        df = df.dropna(how="all")
        if len(df.columns) >= 2:
            return df
    except Exception:
        pass

    # Fallback to comma
    try:
        df = pd.read_csv(filepath, skiprows=header_idx, header=0)
        if "Part" in df.columns:
            df = df.drop(columns=["Part"])
        df = df.apply(pd.to_numeric, errors="coerce")
        df = df.dropna(how="all")
        if len(df.columns) >= 2:
            return df
    except Exception:
        pass

    raise ValueError(f"Cannot parse simulation file: {filepath}")


# ── Metrics ────────────────────────────────────────────────────────────

def compute_metrics(exp_time: np.ndarray, exp_pressure: np.ndarray,
                    sim_time: np.ndarray, sim_pressure: np.ndarray,
                    trim_start: float = 0.5) -> dict:
    """Compute comparison metrics between experimental and simulation pressure.

    Args:
        exp_time: experimental time array [s]
        exp_pressure: experimental pressure [mbar]
        sim_time: simulation time array [s]
        sim_pressure: simulation pressure [mbar or Pa, auto-detected]
        trim_start: trim initial transient [s]

    Returns:
        dict with r, nrmse, peak_error, phase_error
    """
    # Auto-detect unit: if sim pressure > 100, likely Pa → convert to mbar
    if np.nanmax(np.abs(sim_pressure)) > 500:
        sim_pressure = sim_pressure / 100.0  # Pa → mbar

    # Trim initial transient
    exp_mask = exp_time >= trim_start
    sim_mask = sim_time >= trim_start

    exp_t = exp_time[exp_mask]
    exp_p = exp_pressure[exp_mask]
    sim_t = sim_time[sim_mask]
    sim_p = sim_pressure[sim_mask]

    # Interpolate simulation onto experimental time grid
    if len(sim_t) < 2 or len(exp_t) < 2:
        return {"r": np.nan, "nrmse": np.nan, "peak_error": np.nan, "phase_error": np.nan}

    # Common time range
    t_start = max(exp_t.min(), sim_t.min())
    t_end = min(exp_t.max(), sim_t.max())

    common_mask = (exp_t >= t_start) & (exp_t <= t_end)
    common_t = exp_t[common_mask]
    common_exp = exp_p[common_mask]

    # Interpolate sim to common time
    f_sim = interpolate.interp1d(sim_t, sim_p, kind="linear", fill_value="extrapolate")
    common_sim = f_sim(common_t)

    # Pearson correlation
    r, p_value = stats.pearsonr(common_exp, common_sim)

    # NRMSE
    rmse = np.sqrt(np.mean((common_exp - common_sim) ** 2))
    nrmse = rmse / (np.max(common_exp) - np.min(common_exp)) if (np.max(common_exp) - np.min(common_exp)) > 0 else np.nan

    # Peak pressure error
    exp_peak = np.max(common_exp)
    sim_peak = np.max(common_sim)
    peak_error = abs(sim_peak - exp_peak) / abs(exp_peak) if abs(exp_peak) > 0 else np.nan

    # Phase error (time of first major peak)
    exp_peak_idx = np.argmax(common_exp)
    sim_peak_idx = np.argmax(common_sim)
    phase_error = abs(common_t[sim_peak_idx] - common_t[exp_peak_idx])

    return {
        "r": r,
        "r_pvalue": p_value,
        "nrmse": nrmse,
        "rmse_mbar": rmse,
        "peak_error": peak_error,
        "phase_error_s": phase_error,
        "exp_peak_mbar": exp_peak,
        "sim_peak_mbar": sim_peak,
        "n_points": len(common_t),
        "time_range": f"{t_start:.2f}-{t_end:.2f}s",
    }


def compute_peak_statistics(peaks: np.ndarray) -> dict:
    """Compute statistics from 100-repeat peak data.

    Args:
        peaks: (N, 4) array of peak pressures [mbar]

    Returns:
        dict with mean, std, 2sigma band for each peak
    """
    result = {}
    for i in range(peaks.shape[1]):
        col = peaks[:, i]
        result[f"peak{i+1}_mean"] = np.mean(col)
        result[f"peak{i+1}_std"] = np.std(col)
        result[f"peak{i+1}_2sigma_low"] = np.mean(col) - 2 * np.std(col)
        result[f"peak{i+1}_2sigma_high"] = np.mean(col) + 2 * np.std(col)
    return result


# ── Plotting ───────────────────────────────────────────────────────────

def plot_pressure_comparison(exp_df: pd.DataFrame, sim_df: pd.DataFrame,
                             metrics: dict, label: str, out_path: str,
                             time_range: tuple = None):
    """Generate pressure time series comparison figure."""
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 8), height_ratios=[3, 1])

    exp_t = exp_df["time"].values
    exp_p = exp_df["pressure"].values

    sim_t = sim_df.iloc[:, 0].values
    sim_p = sim_df.iloc[:, 1].values

    # Auto-detect Pa → mbar
    if np.nanmax(np.abs(sim_p)) > 500:
        sim_p = sim_p / 100.0

    # Time range
    if time_range:
        exp_mask = (exp_t >= time_range[0]) & (exp_t <= time_range[1])
        sim_mask = (sim_t >= time_range[0]) & (sim_t <= time_range[1])
        exp_t, exp_p = exp_t[exp_mask], exp_p[exp_mask]
        sim_t, sim_p = sim_t[sim_mask], sim_p[sim_mask]

    # Main comparison plot
    ax1.plot(exp_t, exp_p, "b-", alpha=0.6, linewidth=0.5, label="Experiment (SPHERIC)")
    ax1.plot(sim_t, sim_p, "r-", alpha=0.8, linewidth=1.0, label=f"SloshAgent ({label})")
    ax1.set_ylabel("Pressure [mbar]")
    ax1.set_title(f"SPHERIC Test 10 — {label}\n"
                  f"r={metrics['r']:.4f}, NRMSE={metrics['nrmse']:.4f}, "
                  f"Peak Error={metrics['peak_error']:.1%}")
    ax1.legend(loc="upper right")
    ax1.grid(True, alpha=0.3)

    # Residual plot
    if len(sim_t) > 1 and len(exp_t) > 1:
        t_start = max(exp_t.min(), sim_t.min())
        t_end = min(exp_t.max(), sim_t.max())
        common_mask = (exp_t >= t_start) & (exp_t <= t_end)
        common_t = exp_t[common_mask]
        common_exp = exp_p[common_mask]
        f_sim = interpolate.interp1d(sim_t, sim_p, kind="linear", fill_value="extrapolate")
        common_sim = f_sim(common_t)
        residual = common_sim - common_exp

        ax2.plot(common_t, residual, "g-", alpha=0.6, linewidth=0.5)
        ax2.axhline(y=0, color="k", linewidth=0.5)
        ax2.set_xlabel("Time [s]")
        ax2.set_ylabel("Residual [mbar]")
        ax2.grid(True, alpha=0.3)

    plt.tight_layout()
    plt.savefig(out_path, dpi=200, bbox_inches="tight")
    plt.close()
    print(f"  Saved: {out_path}")


def plot_peak_distribution(peaks: np.ndarray, sim_peaks: list,
                           label: str, out_path: str):
    """Plot peak pressure distribution with 100-repeat statistics."""
    fig, axes = plt.subplots(1, 4, figsize=(16, 4))

    peak_labels = ["1st Peak", "2nd Peak", "3rd Peak", "4th Peak"]

    for i, (ax, plabel) in enumerate(zip(axes, peak_labels)):
        col = peaks[:, i]
        ax.hist(col, bins=20, alpha=0.6, color="steelblue", edgecolor="white",
                label="Experiment (N=100)")
        ax.axvline(np.mean(col), color="blue", linewidth=2, label=f"Mean={np.mean(col):.1f}")
        ax.axvline(np.mean(col) - 2*np.std(col), color="blue", linewidth=1, linestyle="--")
        ax.axvline(np.mean(col) + 2*np.std(col), color="blue", linewidth=1, linestyle="--",
                   label=f"±2σ=[{np.mean(col)-2*np.std(col):.1f}, {np.mean(col)+2*np.std(col):.1f}]")

        if i < len(sim_peaks) and sim_peaks[i] is not None:
            ax.axvline(sim_peaks[i], color="red", linewidth=2,
                       label=f"SloshAgent={sim_peaks[i]:.1f}")

        ax.set_title(plabel)
        ax.set_xlabel("Pressure [mbar]")
        ax.legend(fontsize=7)

    fig.suptitle(f"Peak Pressure Distribution — {label}", fontsize=14)
    plt.tight_layout()
    plt.savefig(out_path, dpi=200, bbox_inches="tight")
    plt.close()
    print(f"  Saved: {out_path}")


# ── Main ───────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Compare SPHERIC Test 10 experiment vs simulation")
    parser.add_argument("--sim-dir", required=True, help="Simulation MeasureTool output directory")
    parser.add_argument("--exp-file", required=True, help="SPHERIC experimental time series file")
    parser.add_argument("--peak-file", help="SPHERIC 100-repeat peak file")
    parser.add_argument("--out-dir", default="research/experiments/exp1_spheric", help="Output directory")
    parser.add_argument("--label", default="Simulation", help="Label for this run")
    parser.add_argument("--trim-start", type=float, default=0.5, help="Trim initial transient [s]")
    parser.add_argument("--time-range", type=float, nargs=2, help="Plot time range [start end]")
    parser.add_argument("--probe-idx", type=int, nargs="+", help="Probe column indices (1-based) to compare. Default: all pressure probes")
    args = parser.parse_args()

    os.makedirs(args.out_dir, exist_ok=True)
    os.makedirs(os.path.join(args.out_dir, "figures"), exist_ok=True)
    os.makedirs(os.path.join(args.out_dir, "comparison"), exist_ok=True)

    # Load experimental data
    print(f"Loading experimental data: {args.exp_file}")
    exp_df = load_spheric_timeseries(args.exp_file)
    print(f"  {len(exp_df)} rows, time range: {exp_df['time'].min():.4f} - {exp_df['time'].max():.4f}s")

    # Find simulation pressure file (only *Press*.csv, not Rhop or Vel)
    sim_dir = Path(args.sim_dir)
    sim_files = sorted(sim_dir.glob("*Press*.csv"))

    if not sim_files:
        print(f"ERROR: No simulation pressure CSV found in {args.sim_dir}")
        sys.exit(1)

    # Deduplicate (in case glob finds the same file twice)
    sim_files = list(dict.fromkeys(sim_files))
    print(f"Found {len(sim_files)} simulation pressure file(s)")

    all_metrics = []

    for sim_file in sim_files:
        print(f"\nComparing: {sim_file.name}")
        sim_df = load_simulation_pressure(str(sim_file))
        print(f"  {len(sim_df)} rows, {len(sim_df.columns)} columns")

        # Determine which probes to compare
        if args.probe_idx:
            probe_indices = [i for i in args.probe_idx if i < len(sim_df.columns)]
        else:
            probe_indices = list(range(1, len(sim_df.columns)))

        for col_idx in probe_indices:
            col_name = sim_df.columns[col_idx]
            print(f"  Probe {col_idx} ({col_name}):")

            sim_t = sim_df.iloc[:, 0].values
            sim_p = sim_df.iloc[:, col_idx].values

            metrics = compute_metrics(
                exp_df["time"].values, exp_df["pressure"].values,
                sim_t, sim_p,
                trim_start=args.trim_start
            )
            metrics["probe"] = col_name
            metrics["sim_file"] = sim_file.name
            metrics["label"] = args.label
            all_metrics.append(metrics)

            print(f"    r = {metrics['r']:.4f}")
            print(f"    NRMSE = {metrics['nrmse']:.4f}")
            print(f"    Peak Error = {metrics['peak_error']:.1%}")
            print(f"    Phase Error = {metrics['phase_error_s']:.4f}s")

            # Generate figure
            safe_name = col_name.replace("/", "_").replace(" ", "_")
            fig_path = os.path.join(args.out_dir, "figures",
                                    f"pressure_{sim_file.stem}_{safe_name}.png")

            # Create single-column sim DataFrame for plotting
            plot_sim = pd.DataFrame({sim_df.columns[0]: sim_t, col_name: sim_p})
            plot_pressure_comparison(
                exp_df, plot_sim, metrics, f"{args.label} - {col_name}",
                fig_path, time_range=tuple(args.time_range) if args.time_range else None
            )

    # Peak distribution analysis
    if args.peak_file and os.path.exists(args.peak_file):
        print(f"\nLoading peak data: {args.peak_file}")
        peaks = load_spheric_peaks(args.peak_file)
        print(f"  {peaks.shape[0]} repeats × {peaks.shape[1]} peaks")

        peak_stats = compute_peak_statistics(peaks)
        for k, v in peak_stats.items():
            print(f"  {k}: {v:.2f}")

        # Plot peak distribution (sim peaks from first probe if available)
        sim_peaks = [None] * 4  # Placeholder until sim peaks extracted
        fig_path = os.path.join(args.out_dir, "figures", "peak_distribution.png")
        plot_peak_distribution(peaks, sim_peaks, args.label, fig_path)

    # Save metrics CSV
    metrics_df = pd.DataFrame(all_metrics)
    metrics_path = os.path.join(args.out_dir, "comparison", "metrics.csv")
    metrics_df.to_csv(metrics_path, index=False)
    print(f"\nMetrics saved: {metrics_path}")

    # Summary report
    report_path = os.path.join(args.out_dir, "comparison", "summary.md")
    with open(report_path, "w") as f:
        f.write(f"# SPHERIC Test 10 Comparison — {args.label}\n\n")
        f.write(f"Date: {pd.Timestamp.now().strftime('%Y-%m-%d %H:%M')}\n\n")
        f.write("## Metrics\n\n")
        f.write("| Probe | r | NRMSE | Peak Error | Phase Error [s] |\n")
        f.write("|-------|---|-------|------------|----------------|\n")
        for m in all_metrics:
            f.write(f"| {m['probe']} | {m['r']:.4f} | {m['nrmse']:.4f} | "
                    f"{m['peak_error']:.1%} | {m['phase_error_s']:.4f} |\n")

        # Pass/fail assessment
        f.write("\n## Assessment\n\n")
        for m in all_metrics:
            passed = []
            if m["r"] > 0.9:
                passed.append("r > 0.9 PASS")
            else:
                passed.append(f"r = {m['r']:.4f} FAIL (target > 0.9)")

            if m["nrmse"] < 0.15:
                passed.append("NRMSE < 0.15 PASS")
            else:
                passed.append(f"NRMSE = {m['nrmse']:.4f} FAIL (target < 0.15)")

            if m["peak_error"] < 0.20:
                passed.append("Peak Error < 20% PASS")
            else:
                passed.append(f"Peak Error = {m['peak_error']:.1%} FAIL (target < 20%)")

            f.write(f"- **{m['probe']}**: {'; '.join(passed)}\n")

    print(f"Report saved: {report_path}")

    # Print final summary
    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)
    for m in all_metrics:
        status = "PASS" if (m["r"] > 0.9 and m["nrmse"] < 0.15) else "FAIL"
        print(f"  [{status}] {m['probe']}: r={m['r']:.4f}, NRMSE={m['nrmse']:.4f}")


if __name__ == "__main__":
    main()
