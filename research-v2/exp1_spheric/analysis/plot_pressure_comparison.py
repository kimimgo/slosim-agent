#!/usr/bin/env python3
"""EXP-1 압력 시계열 비교: Run 001 vs Run 002 vs Experimental"""

import numpy as np
import matplotlib.pyplot as plt
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]  # slosim-agent/
SIM_DIR = ROOT / "simulations" / "exp1"
EXP_DIR = ROOT / "datasets" / "spheric" / "case_1"
FIG_DIR = Path(__file__).resolve().parent.parent / "figures"
FIG_DIR.mkdir(exist_ok=True)

# --- Load simulation data (semicolon-delimited, 4 header lines) ---
def load_sim(run_id):
    path = SIM_DIR / f"run_{run_id:03d}" / "PressureLateral_Press.csv"
    data = np.genfromtxt(path, delimiter=';', skip_header=4)
    t = data[:, 1]          # Time [s]
    p1 = data[:, 2] / 100   # Press_0 [Pa] → [mbar]
    return t, p1

t001, p001 = load_sim(1)
t002, p002 = load_sim(2)

# --- Load experimental data (tab-delimited, 1 header line) ---
exp_path = EXP_DIR / "lateral_water_1x.txt"
exp_data = np.genfromtxt(exp_path, delimiter='\t', skip_header=1)
t_exp = exp_data[:, 0]
p_exp = exp_data[:, 1]  # already mbar

# --- Load experimental peak statistics (103 repetitions, 2 header lines) ---
peak_path = EXP_DIR / "Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt"
peak_data = np.genfromtxt(peak_path, delimiter='\t', skip_header=2)
peak_means = np.nanmean(peak_data, axis=0)  # 4 peaks
peak_stds = np.nanstd(peak_data, axis=0)
print(f"Exp peak means: {peak_means}")
print(f"Exp peak ±2σ:   {2*peak_stds}")

# --- Peak detection ---
from scipy.signal import find_peaks

def detect_peaks(t, p, min_height=3.0, min_dist_s=1.0):
    dt = np.median(np.diff(t))
    min_dist = int(min_dist_s / dt) if dt > 0 else 1000
    peaks, props = find_peaks(p, height=min_height, distance=min_dist)
    return peaks

# 시뮬레이션은 rest에서 시작하므로 초기 작은 피크 제외 (>15 mbar)
peaks_001 = detect_peaks(t001, p001, min_height=15.0, min_dist_s=1.2)
peaks_002 = detect_peaks(t002, p002, min_height=15.0, min_dist_s=1.2)
peaks_exp = detect_peaks(t_exp, p_exp, min_height=15.0, min_dist_s=0.8)

# --- Plot ---
fig, axes = plt.subplots(2, 1, figsize=(14, 10), gridspec_kw={'height_ratios': [3, 1]})

# Top panel: time series
ax = axes[0]
ax.plot(t_exp, p_exp, color='black', alpha=0.5, linewidth=0.5, label='Experimental (20kHz)')
ax.plot(t001, p001, color='tab:blue', linewidth=1.0, label=f'Run 001 (dp=0.004, DBC)')
ax.plot(t002, p002, color='tab:red', linewidth=1.0, label=f'Run 002 (dp=0.002, DBC)')

# Peak markers
for i, idx in enumerate(peaks_001[:4]):
    marker_label = f'Run001 P{i+1}={p001[idx]:.1f}mbar' if i == 0 else None
    ax.plot(t001[idx], p001[idx], 'v', color='tab:blue', markersize=8)
    ax.annotate(f'{p001[idx]:.1f}', (t001[idx], p001[idx]),
                textcoords="offset points", xytext=(0, 10), fontsize=8, color='tab:blue', ha='center')

for i, idx in enumerate(peaks_002[:4]):
    ax.plot(t002[idx], p002[idx], 'v', color='tab:red', markersize=8)
    ax.annotate(f'{p002[idx]:.1f}', (t002[idx], p002[idx]),
                textcoords="offset points", xytext=(0, -15), fontsize=8, color='tab:red', ha='center')

# Experimental peak mean ± 2σ bands (horizontal lines at approximate peak times)
ax.set_xlabel('Time [s]', fontsize=12)
ax.set_ylabel('Pressure [mbar]', fontsize=12)
ax.set_title('SPHERIC Test 10 — Lateral Water P1: Sim vs Experiment', fontsize=14, fontweight='bold')
ax.legend(loc='upper left', fontsize=10)
ax.set_xlim(0, 7)
ax.set_ylim(-10, max(p002.max(), p001.max(), p_exp.max()) * 1.1)
ax.grid(True, alpha=0.3)

# Annotate the Run 002 spike
spike_idx = peaks_002[1] if len(peaks_002) > 1 else None
if spike_idx is not None and p002[spike_idx] > 100:
    ax.annotate(f'DBC spike\n{p002[spike_idx]:.0f} mbar',
                xy=(t002[spike_idx], p002[spike_idx]),
                xytext=(t002[spike_idx] + 0.5, p002[spike_idx] - 30),
                arrowprops=dict(arrowstyle='->', color='tab:red'),
                fontsize=9, color='tab:red', fontweight='bold')

# Bottom panel: peak comparison bar chart
ax2 = axes[1]
n_peaks = 4
x = np.arange(n_peaks)
width = 0.2

# Experimental means ± 2σ
bars_exp = ax2.bar(x - width, peak_means[:n_peaks], width,
                    yerr=2*peak_stds[:n_peaks], label='Exp (mean±2σ)',
                    color='gray', alpha=0.7, capsize=4)

# Run 001 peaks
if len(peaks_001) >= n_peaks:
    sim001_vals = [p001[peaks_001[i]] for i in range(n_peaks)]
else:
    sim001_vals = [p001[peaks_001[i]] if i < len(peaks_001) else 0 for i in range(n_peaks)]
bars_001 = ax2.bar(x, sim001_vals, width, label='Run 001 (dp=0.004)',
                    color='tab:blue', alpha=0.7)

# Run 002 peaks
if len(peaks_002) >= n_peaks:
    sim002_vals = [p002[peaks_002[i]] for i in range(n_peaks)]
else:
    sim002_vals = [p002[peaks_002[i]] if i < len(peaks_002) else 0 for i in range(n_peaks)]
bars_002 = ax2.bar(x + width, sim002_vals, width, label='Run 002 (dp=0.002)',
                    color='tab:red', alpha=0.7)

ax2.set_xlabel('Peak Number', fontsize=12)
ax2.set_ylabel('Pressure [mbar]', fontsize=12)
ax2.set_title('First 4 Peak Comparison', fontsize=12, fontweight='bold')
ax2.set_xticks(x)
ax2.set_xticklabels([f'Peak {i+1}' for i in range(n_peaks)])
ax2.legend(fontsize=9)
ax2.grid(True, alpha=0.3, axis='y')

# Add error percentages on top of sim bars
for i in range(n_peaks):
    if peak_means[i] > 0:
        err001 = (sim001_vals[i] - peak_means[i]) / peak_means[i] * 100
        err002 = (sim002_vals[i] - peak_means[i]) / peak_means[i] * 100
        ax2.text(x[i], sim001_vals[i] + 2, f'{err001:+.0f}%', ha='center', fontsize=7, color='tab:blue')
        ax2.text(x[i] + width, sim002_vals[i] + 2, f'{err002:+.0f}%', ha='center', fontsize=7, color='tab:red')

plt.tight_layout()
out_path = FIG_DIR / "pressure_comparison_run001_002_exp.png"
plt.savefig(out_path, dpi=150, bbox_inches='tight')
print(f"Saved: {out_path}")
print(f"\n=== Peak Summary ===")
print(f"{'Peak':<8} {'Exp Mean':>10} {'Exp ±2σ':>10} {'Run001':>10} {'Err001':>10} {'Run002':>10} {'Err002':>10}")
for i in range(n_peaks):
    err001 = (sim001_vals[i] - peak_means[i]) / peak_means[i] * 100 if peak_means[i] > 0 else float('inf')
    err002 = (sim002_vals[i] - peak_means[i]) / peak_means[i] * 100 if peak_means[i] > 0 else float('inf')
    print(f"Peak {i+1:<3} {peak_means[i]:>10.1f} {2*peak_stds[i]:>10.1f} {sim001_vals[i]:>10.1f} {err001:>+9.1f}% {sim002_vals[i]:>10.1f} {err002:>+9.1f}%")
