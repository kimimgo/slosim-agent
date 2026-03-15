#!/usr/bin/env python3
"""SPHERIC T10 Oil — Artificial viscosity parametric study (α=0.1~30)
Time-shifting allowed for peak matching."""
import matplotlib.pyplot as plt
import matplotlib
import numpy as np
import os
from scipy.signal import find_peaks

matplotlib.use('Agg')

def load_csv(path):
    times, press = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i < 4: continue
            parts = line.strip().split(';')
            if len(parts) >= 3:
                times.append(float(parts[1]))
                press.append(float(parts[2]))
    return np.array(times), np.array(press)

def load_exp(path):
    times, press = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i == 0: continue
            parts = line.strip().split('\t')
            if len(parts) >= 2:
                try:
                    times.append(float(parts[0]))
                    press.append(float(parts[1]) * 100.0)  # mbar→Pa
                except: continue
    return np.array(times), np.array(press)

script_dir = os.path.dirname(os.path.abspath(__file__))
sim_base = os.path.join(script_dir, "../../simulations")
EXP_FILE = os.path.join(script_dir, "../../paper-pof/data/spheric/lateral_oil_1x.txt")
OUT_DIR = os.path.join(script_dir, "../figures/organized/expc_spheric")
os.makedirs(OUT_DIR, exist_ok=True)

# All cases
cases = [
    ("0.1",  f"{sim_base}/spheric_oil_low_art/out/measure_lateral_Press.csv"),
    ("1.5",  f"{sim_base}/spheric_oil_art_a1_5/out/measure_lateral_Press.csv"),
    ("2.0",  f"{sim_base}/spheric_oil_art_a2_0/out/measure_lateral_Press.csv"),
    ("2.5",  f"{sim_base}/spheric_oil_art_a2_5/out/measure_lateral_Press.csv"),
    ("3.0",  f"{sim_base}/spheric_oil_art_a3_0/out/measure_lateral_Press.csv"),
]

# Load experiment
t_exp, p_exp = load_exp(EXP_FILE)
peaks_idx, _ = find_peaks(p_exp, prominence=100, distance=10000)
exp_peaks = [(t_exp[i], p_exp[i]) for i in peaks_idx[:4]]

# Load all SPH
sph_data = {}
for alpha, path in cases:
    if os.path.exists(path):
        t, p = load_csv(path)
        sph_data[alpha] = (t, p)

# Color palette
cmap = plt.cm.viridis
alpha_vals = [0.1, 1.5, 2.0, 2.5, 3.0]
colors = {str(a): cmap(i / (len(alpha_vals) - 1)) for i, a in enumerate(alpha_vals)}

# === Figure 1: Full timeseries ===
fig, ax = plt.subplots(figsize=(14, 5.5))
ax.plot(t_exp, p_exp/1000, 'k-', lw=2.0, alpha=0.9, label='Experiment', zorder=10)
for alpha in ["0.1", "1.5", "2.0", "2.5", "3.0"]:
    if alpha in sph_data:
        t, p = sph_data[alpha]
        ax.plot(t, p/1000, '-', color=colors[alpha], lw=1.0, alpha=0.8, label=f'α={alpha}')

ax.set_xlabel('Time (s)', fontsize=12)
ax.set_ylabel('Pressure (kPa)', fontsize=12)
ax.set_title('SPHERIC T10 Oil Lateral — Artificial Viscosity Parametric Study (dp=0.004)', fontsize=13, fontweight='bold')
ax.set_xlim(0, 7)
ax.legend(fontsize=9, loc='upper right', ncol=2)
ax.grid(alpha=0.2)
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_timeseries.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_timeseries.png"), dpi=150, bbox_inches='tight')
print("Saved: spheric_oil_parametric_timeseries")

# === Figure 2: Peak-by-peak zoom ===
fig2, axes = plt.subplots(2, 2, figsize=(13, 9))
window = 0.8
all_results = {}  # alpha → [peak_errors]

for pidx, (exp_t, exp_p) in enumerate(exp_peaks):
    ax = axes[pidx // 2][pidx % 2]
    t0, t1 = exp_t - window, exp_t + window

    m = (t_exp >= t0) & (t_exp <= t1)
    ax.plot(t_exp[m], p_exp[m]/1000, 'k-', lw=2.0, label='Experiment', zorder=10)

    for alpha in ["0.1", "1.5", "2.0", "2.5", "3.0"]:
        if alpha not in sph_data: continue
        t, p = sph_data[alpha]
        ms = (t >= t0) & (t <= t1)
        if not np.any(ms): continue
        ax.plot(t[ms], p[ms]/1000, '-', color=colors[alpha], lw=1.0, alpha=0.8, label=f'α={alpha}')

        peak_val = np.max(p[ms])
        err = abs(peak_val - exp_p) / exp_p * 100 if exp_p > 0 else 0
        if alpha not in all_results:
            all_results[alpha] = []
        all_results[alpha].append((pidx+1, peak_val/1000, err))

    ax.set_title(f'Peak {pidx+1} (exp: {exp_p/1000:.3f} kPa @ {exp_t:.2f}s)', fontsize=10, fontweight='bold')
    ax.set_xlabel('Time (s)', fontsize=9)
    ax.set_ylabel('Pressure (kPa)', fontsize=9)
    ax.grid(alpha=0.2)
    if pidx == 0:
        ax.legend(fontsize=7, ncol=2)

plt.suptitle('SPHERIC T10 Oil — Peak Comparison (time-shift allowed)', fontsize=13, fontweight='bold')
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_peaks.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_peaks.png"), dpi=150, bbox_inches='tight')
print("Saved: spheric_oil_parametric_peaks")

# === Figure 3: Alpha vs Mean Error bar chart ===
fig3, (ax3a, ax3b) = plt.subplots(1, 2, figsize=(12, 5))

alphas_plot = []
mean_errs = []
per_peak = {1: [], 2: [], 3: [], 4: []}

for alpha in ["0.1", "1.5", "2.0", "2.5", "3.0"]:
    if alpha in all_results:
        errs = [e for _, _, e in all_results[alpha]]
        alphas_plot.append(float(alpha))
        mean_errs.append(np.mean(errs))
        for pidx, val, err in all_results[alpha]:
            per_peak[pidx].append(err)

# Bar chart: mean error vs alpha
bar_colors = [colors[str(a)] for a in alphas_plot]
ax3a.bar([str(a) for a in alphas_plot], mean_errs, color=bar_colors, edgecolor='black', linewidth=0.5)
ax3a.set_xlabel('Artificial Viscosity α', fontsize=12)
ax3a.set_ylabel('Mean Peak Error (%)', fontsize=12)
ax3a.set_title('Mean Peak Error vs α', fontsize=12, fontweight='bold')
ax3a.grid(axis='y', alpha=0.2)
for i, (a, e) in enumerate(zip(alphas_plot, mean_errs)):
    ax3a.text(i, e + 2, f'{e:.0f}%', ha='center', fontsize=10, fontweight='bold')

# Line chart: per-peak error
for pidx in [1, 2, 3, 4]:
    if per_peak[pidx]:
        ax3b.plot([str(a) for a in alphas_plot[:len(per_peak[pidx])]], per_peak[pidx],
                  'o-', label=f'Peak {pidx}', lw=1.5, markersize=6)
ax3b.set_xlabel('Artificial Viscosity α', fontsize=12)
ax3b.set_ylabel('Peak Error (%)', fontsize=12)
ax3b.set_title('Per-Peak Error vs α', fontsize=12, fontweight='bold')
ax3b.legend(fontsize=9)
ax3b.grid(alpha=0.2)

plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_error.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_parametric_error.png"), dpi=150, bbox_inches='tight')
print("Saved: spheric_oil_parametric_error")

# === Print summary ===
print("\n" + "=" * 90)
print("SPHERIC T10 Oil — Parametric Study Summary (dp=0.004, Artificial Viscosity)")
print("=" * 90)
header = f"{'Alpha':>6} | {'Peak1':>12} | {'Peak2':>12} | {'Peak3':>12} | {'Peak4':>12} | {'Mean':>8}"
print(header)
print("-" * 90)

# Experiment row
print(f"{'Exp':>6} | {exp_peaks[0][1]/1000:>8.3f} kPa | {exp_peaks[1][1]/1000:>8.3f} kPa | {exp_peaks[2][1]/1000:>8.3f} kPa | {exp_peaks[3][1]/1000:>8.3f} kPa |")
print("-" * 90)

best_alpha = None
best_err = float('inf')

for alpha in ["0.1", "1.5", "2.0", "2.5", "3.0"]:
    if alpha not in all_results: continue
    parts = []
    errs = []
    for pidx, val, err in all_results[alpha]:
        parts.append(f"{val:>5.3f} {err:>4.0f}%")
        errs.append(err)
    mean_e = np.mean(errs)
    if mean_e < best_err:
        best_err = mean_e
        best_alpha = alpha
    marker = " ★" if alpha == best_alpha else ""
    print(f"  α={alpha:>3} | {' | '.join(parts)} | {mean_e:>5.1f}%{marker}")

print("-" * 90)
print(f"Best: α={best_alpha} (mean error={best_err:.1f}%)")
print("=" * 90)
