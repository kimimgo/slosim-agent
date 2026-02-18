#!/usr/bin/env python3
"""
analyze_exp5.py — Analyze EXP-5 Industrial PoC results

Reads MeasureTool elevation data and generates:
  - Fig 5: Baffle comparison — SWL with/without baffle
  - Fig 6: Seismic scenario — SWL time series
  - Table 6: EXP-5 summary metrics
"""

import numpy as np
from pathlib import Path
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import json

OUT_DIR = Path("research/experiments/exp5_industrial")
OUT_DIR.mkdir(parents=True, exist_ok=True)
(OUT_DIR / "figures").mkdir(exist_ok=True)
(OUT_DIR / "comparison").mkdir(exist_ok=True)


def load_elevation(sim_dir, probe_idx=0):
    """Load MeasureTool elevation CSV → (time, elevation_mm).

    probe_idx: 0=Left, 1=Right (for baffle/seismic 2-probe setup)
    Tries elevation_fixed first (adjusted probe positions), then elevation.
    """
    filepath = Path(sim_dir) / "elevation_fixed_Elevation.csv"
    if not filepath.exists():
        filepath = Path(sim_dir) / "elevation_Elevation.csv"
    if not filepath.exists():
        print(f"  WARNING: {filepath} not found")
        return None, None

    lines = filepath.read_text().strip().split("\n")

    data_start = 0
    for i, line in enumerate(lines):
        if line.strip() and line.strip()[0].isdigit():
            data_start = i
            break

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


# ── Figure 5: Baffle Comparison ──────────────────────────────────

def fig_baffle_comparison():
    """Compare SWL with and without baffle."""
    fig, axes = plt.subplots(2, 1, figsize=(12, 8), height_ratios=[3, 1])

    # Left wall SWL
    ax1 = axes[0]
    for label, sim_dir, color in [
        ("No Baffle", "simulations/exp5_baffle_no", "red"),
        ("With Baffle", "simulations/exp5_baffle_yes", "blue"),
    ]:
        t, swl = load_elevation(sim_dir, probe_idx=0)  # Left
        if t is not None:
            ax1.plot(t, swl, color=color, linewidth=1.0, alpha=0.8, label=label)

    ax1.axhline(y=300, color="gray", linestyle="--", alpha=0.5,
                label="Still level (300mm)")
    ax1.set_ylabel("SWL at Left Wall [mm]")
    ax1.set_title("Baffle Effect — Left Wall Free Surface at Resonance "
                   "(f_1=0.758Hz, A=10mm)")
    ax1.legend(loc="upper right")
    ax1.grid(True, alpha=0.3)

    # Right wall SWL
    ax2 = axes[1]
    for label, sim_dir, color in [
        ("No Baffle", "simulations/exp5_baffle_no", "red"),
        ("With Baffle", "simulations/exp5_baffle_yes", "blue"),
    ]:
        t, swl = load_elevation(sim_dir, probe_idx=1)  # Right
        if t is not None:
            ax2.plot(t, swl, color=color, linewidth=1.0, alpha=0.8, label=label)

    ax2.set_xlabel("Time [s]")
    ax2.set_ylabel("SWL Right [mm]")
    ax2.legend(fontsize=8, loc="upper right")
    ax2.grid(True, alpha=0.3)

    plt.tight_layout()
    path = OUT_DIR / "figures" / "fig5_baffle_comparison.png"
    plt.savefig(path, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {path}")


# ── Figure 6: Seismic Scenario ───────────────────────────────────

def fig_seismic():
    """Seismic scenario SWL time series."""
    fig, ax = plt.subplots(figsize=(12, 5))

    for probe_idx, color, label in [
        (0, "blue", "Left Wall (x=0.1m)"),
        (1, "red", "Right Wall (x=9.9m)"),
    ]:
        t, swl = load_elevation("simulations/exp5_seismic", probe_idx=probe_idx)
        if t is not None:
            ax.plot(t, swl / 1000, color=color, linewidth=0.8, alpha=0.8,
                    label=label)

    ax.axhline(y=4.8, color="gray", linestyle="--", alpha=0.5,
               label="Still level (4.8m)")
    ax.set_xlabel("Time [s]")
    ax.set_ylabel("Free Surface Height [m]")
    ax.set_title("Seismic Sloshing — 10m x 5m x 8m Tank, 60% Fill, "
                 "f=0.3Hz, A=50mm")
    ax.legend(loc="upper right")
    ax.grid(True, alpha=0.3)

    plt.tight_layout()
    path = OUT_DIR / "figures" / "fig6_seismic.png"
    plt.savefig(path, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {path}")


# ── Table 6: Summary ─────────────────────────────────────────────

def generate_table():
    """Generate EXP-5 summary."""
    results = []

    for name, sim_dir, desc, still_mm in [
        ("Baffle_No", "simulations/exp5_baffle_no",
         "1m x 0.5m x 0.6m, 50% fill, f1=0.758Hz, no baffle", 300),
        ("Baffle_Yes", "simulations/exp5_baffle_yes",
         "1m x 0.5m x 0.6m, 50% fill, f1=0.758Hz, with baffle", 300),
        ("Seismic", "simulations/exp5_seismic",
         "10m x 5m x 8m, 60% fill, f=0.3Hz", 4800),
    ]:
        t_l, swl_l = load_elevation(sim_dir, probe_idx=0)
        if t_l is not None:
            mask = t_l > 3.0
            if np.sum(mask) > 10:
                swl_ss = swl_l[mask]
                amp = (np.max(swl_ss) - np.min(swl_ss)) / 2
                max_swl = np.max(swl_ss)
                min_swl = np.min(swl_ss)
            else:
                amp = max_swl = min_swl = 0.0
        else:
            amp = max_swl = min_swl = 0.0

        results.append({
            "case": name,
            "description": desc,
            "still_level_mm": still_mm,
            "amp_left_mm": f"{amp:.1f}",
            "max_swl_mm": f"{max_swl:.1f}",
            "min_swl_mm": f"{min_swl:.1f}",
        })

    # Compute baffle reduction
    if len(results) >= 2:
        amp_no = float(results[0]["amp_left_mm"])
        amp_yes = float(results[1]["amp_left_mm"])
        if amp_no > 0:
            reduction = (amp_no - amp_yes) / amp_no * 100
            results[1]["reduction_pct"] = f"{reduction:.1f}"

    json_path = OUT_DIR / "comparison" / "table6_summary.json"
    with open(json_path, "w") as f:
        json.dump(results, f, indent=2)

    md_path = OUT_DIR / "comparison" / "table6_summary.md"
    with open(md_path, "w") as f:
        f.write("# Table 6: EXP-5 Industrial PoC Results\n\n")
        f.write("## Baffle Comparison (Resonance Excitation)\n\n")
        f.write("| Case | Still Level [mm] | Amplitude [mm] | Max SWL [mm] | "
                "Min SWL [mm] | Reduction |\n")
        f.write("|------|-----------------|---------------|-------------|"
                "-------------|----------|\n")
        for r in results[:2]:
            red = r.get("reduction_pct", "--")
            f.write(f"| {r['case']} | {r['still_level_mm']} | "
                    f"{r['amp_left_mm']} | {r['max_swl_mm']} | "
                    f"{r['min_swl_mm']} | {red}% |\n")

        f.write("\n## Seismic Scenario\n\n")
        f.write("| Amplitude [mm] | Max SWL [mm] | Min SWL [mm] |\n")
        f.write("|---------------|-------------|-------------|\n")
        if len(results) > 2:
            r = results[2]
            f.write(f"| {r['amp_left_mm']} | {r['max_swl_mm']} | "
                    f"{r['min_swl_mm']} |\n")

    print(f"Saved: {json_path}")
    print(f"Saved: {md_path}")


# ── Main ──────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("Analyzing EXP-5 industrial PoC results...\n")
    fig_baffle_comparison()
    fig_seismic()
    generate_table()
    print("\nDone!")
