#!/usr/bin/env python3
"""SPHERIC T10 Oil Lateral — Visco30 vs Artificial vs Experiment comparison"""
import matplotlib.pyplot as plt
import matplotlib
import numpy as np
import csv
import sys
import os

matplotlib.use('Agg')

def load_measure_csv(path):
    """Load MeasureTool pressure CSV (skip 4 header lines)."""
    times, press = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i < 4:
                continue
            parts = line.strip().split(';')
            if len(parts) >= 3:
                times.append(float(parts[1]))
                press.append(float(parts[2]))
    return np.array(times), np.array(press)

def load_experiment(path):
    """Load SPHERIC experiment data (TSV: time pressure[mbar] ...)."""
    times, press = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i == 0:  # skip header
                continue
            line = line.strip()
            if not line or line.startswith('#') or line.startswith('//'):
                continue
            parts = line.split('\t')
            if len(parts) >= 2:
                try:
                    times.append(float(parts[0]))
                    # mbar → Pa (1 mbar = 100 Pa)
                    press.append(float(parts[1]) * 100.0)
                except ValueError:
                    continue
    return np.array(times), np.array(press)

# Paths
BASE = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SIM_V30 = os.path.join(BASE, "../../simulations/spheric_oil_low_visco30/out/measure_lateral_Press.csv")
SIM_ART = os.path.join(BASE, "../../simulations/spheric_oil_low_art/out/measure_lateral_Press.csv")
EXP_FILE = os.path.join(BASE, "../paper-pof/data/spheric/lateral_oil_1x.txt")
OUT_DIR = os.path.join(BASE, "figures/organized/expc_spheric")

# Fix paths to be relative to script location
script_dir = os.path.dirname(os.path.abspath(__file__))
SIM_V30 = os.path.join(script_dir, "../../simulations/spheric_oil_low_visco30/out/measure_lateral_Press.csv")
SIM_ART = os.path.join(script_dir, "../../simulations/spheric_oil_low_art/out/measure_lateral_Press.csv")
EXP_FILE = os.path.join(script_dir, "../../paper-pof/data/spheric/lateral_oil_1x.txt")
OUT_DIR = os.path.join(script_dir, "../figures/organized/expc_spheric")
os.makedirs(OUT_DIR, exist_ok=True)

# Load data
print("Loading visco30...")
t_v30, p_v30 = load_measure_csv(SIM_V30)

# Try loading Artificial if available
has_art = os.path.exists(SIM_ART)
if has_art:
    print("Loading artificial...")
    t_art, p_art = load_measure_csv(SIM_ART)
else:
    print("Artificial data not yet available, plotting without it")

print("Loading experiment...")
t_exp, p_exp = load_experiment(EXP_FILE)

# Normalize pressure to kPa
p_v30_kpa = p_v30 / 1000.0
if has_art:
    p_art_kpa = p_art / 1000.0
p_exp_kpa = p_exp / 1000.0

# === Figure 1: Full timeseries comparison ===
fig, ax = plt.subplots(figsize=(14, 5))

ax.plot(t_exp, p_exp_kpa, 'k-', linewidth=1.5, alpha=0.8, label='Experiment (SPHERIC T10)')
ax.plot(t_v30, p_v30_kpa, '-', color='#1f77b4', linewidth=1.2, alpha=0.8, label='SPH Laminar+SPS (α_SPS=30)')
if has_art:
    ax.plot(t_art, p_art_kpa, '-', color='#e6550d', linewidth=1.2, alpha=0.7, label='SPH Artificial (α=0.1)')

ax.set_xlabel('Time (s)', fontsize=12)
ax.set_ylabel('Pressure (kPa)', fontsize=12)
ax.set_title('SPHERIC Test 10 — Oil Lateral Impact Pressure', fontsize=14, fontweight='bold')
ax.set_xlim(0, 7)
ax.legend(fontsize=10, loc='upper right')
ax.grid(alpha=0.2)

plt.tight_layout()
fname = "spheric_oil_visco30_comparison"
plt.savefig(os.path.join(OUT_DIR, f"{fname}.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, f"{fname}.png"), dpi=150, bbox_inches='tight')
print(f"Saved: {OUT_DIR}/{fname}.{{pdf,png}}")

# === Figure 2: Peak zoom (first 4 peaks) ===
fig2, axes = plt.subplots(2, 2, figsize=(12, 8))
# Approximate peak times for Oil lateral (T=1.535s, period ≈ 1.535s)
peak_windows = [
    (1.0, 2.2, "Peak 1"),
    (2.5, 3.7, "Peak 2"),
    (4.0, 5.2, "Peak 3"),
    (5.5, 6.7, "Peak 4"),
]

for idx, (t_start, t_end, title) in enumerate(peak_windows):
    ax = axes[idx // 2][idx % 2]

    # Experiment
    mask_exp = (t_exp >= t_start) & (t_exp <= t_end)
    ax.plot(t_exp[mask_exp], p_exp_kpa[mask_exp], 'k-', linewidth=1.5, alpha=0.8, label='Experiment')

    # Visco30
    mask_v30 = (t_v30 >= t_start) & (t_v30 <= t_end)
    ax.plot(t_v30[mask_v30], p_v30_kpa[mask_v30], '-', color='#1f77b4', linewidth=1.2, label='Laminar+SPS')

    # Artificial
    if has_art:
        mask_art = (t_art >= t_start) & (t_art <= t_end)
        ax.plot(t_art[mask_art], p_art_kpa[mask_art], '-', color='#e6550d', linewidth=1.0, alpha=0.7, label='Artificial')

    ax.set_title(title, fontsize=11, fontweight='bold')
    ax.set_xlabel('Time (s)', fontsize=9)
    ax.set_ylabel('Pressure (kPa)', fontsize=9)
    ax.grid(alpha=0.2)
    if idx == 0:
        ax.legend(fontsize=8)

    # Compute peak error for visco30
    if np.any(mask_exp) and np.any(mask_v30):
        exp_peak = np.max(p_exp_kpa[mask_exp])
        # Find closest sim time to exp peak time
        exp_peak_time = t_exp[mask_exp][np.argmax(p_exp_kpa[mask_exp])]
        v30_at_window = p_v30_kpa[mask_v30]
        v30_peak = np.max(v30_at_window)
        err = abs(v30_peak - exp_peak) / exp_peak * 100 if exp_peak > 0 else 0
        ax.text(0.02, 0.95, f'Peak err: {err:.1f}%', transform=ax.transAxes,
                fontsize=8, va='top', color='#1f77b4')

plt.suptitle('SPHERIC T10 Oil — Peak-by-Peak Comparison (Visco α_SPS=30)', fontsize=13, fontweight='bold')
plt.tight_layout()
fname2 = "spheric_oil_visco30_peaks"
plt.savefig(os.path.join(OUT_DIR, f"{fname2}.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, f"{fname2}.png"), dpi=150, bbox_inches='tight')
print(f"Saved: {OUT_DIR}/{fname2}.{{pdf,png}}")

# Print peak errors
print("\n=== Peak Error Summary (Laminar+SPS α_SPS=30) ===")
for t_start, t_end, title in peak_windows:
    mask_exp = (t_exp >= t_start) & (t_exp <= t_end)
    mask_v30 = (t_v30 >= t_start) & (t_v30 <= t_end)
    if np.any(mask_exp) and np.any(mask_v30):
        exp_peak = np.max(p_exp_kpa[mask_exp])
        v30_peak = np.max(p_v30_kpa[mask_v30])
        err = abs(v30_peak - exp_peak) / exp_peak * 100
        print(f"  {title}: Exp={exp_peak:.2f} kPa, SPH={v30_peak:.2f} kPa, Error={err:.1f}%")
