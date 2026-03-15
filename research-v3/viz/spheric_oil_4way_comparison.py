#!/usr/bin/env python3
"""SPHERIC T10 Oil Lateral — 4-way comparison: Experiment vs Art-0.1 vs Art-30 vs Lam+SPS-30
Time-shifting allowed: peak magnitude comparison with ±0.8T window matching."""
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
SIM_V30 = os.path.join(script_dir, "../../simulations/spheric_oil_low_visco30/out/measure_lateral_Press.csv")
SIM_ART01 = os.path.join(script_dir, "../../simulations/spheric_oil_low_art/out/measure_lateral_Press.csv")
SIM_ART30 = os.path.join(script_dir, "../../simulations/spheric_oil_low_art30/out/measure_lateral_Press.csv")
EXP_FILE = os.path.join(script_dir, "../../paper-pof/data/spheric/lateral_oil_1x.txt")
OUT_DIR = os.path.join(script_dir, "../figures/organized/expc_spheric")
os.makedirs(OUT_DIR, exist_ok=True)

# Load all
t_exp, p_exp = load_exp(EXP_FILE)
t_art01, p_art01 = load_csv(SIM_ART01)
t_art30, p_art30 = load_csv(SIM_ART30)
t_lam30, p_lam30 = load_csv(SIM_V30)

colors = {
    'exp': '#000000',
    'art01': '#e6550d',
    'art30': '#2ca02c',
    'lam30': '#1f77b4',
}
labels = {
    'exp': 'Experiment (SPHERIC T10)',
    'art01': 'Artificial α=0.1',
    'art30': 'Artificial α=30',
    'lam30': 'Laminar+SPS (α_SPS=30)',
}

# === Figure 1: Full timeseries ===
fig, ax = plt.subplots(figsize=(14, 5))
ax.plot(t_exp, p_exp/1000, '-', color=colors['exp'], lw=1.5, alpha=0.8, label=labels['exp'])
ax.plot(t_art01, p_art01/1000, '-', color=colors['art01'], lw=1.0, alpha=0.7, label=labels['art01'])
ax.plot(t_art30, p_art30/1000, '-', color=colors['art30'], lw=1.0, alpha=0.7, label=labels['art30'])
ax.plot(t_lam30, p_lam30/1000, '-', color=colors['lam30'], lw=1.0, alpha=0.6, label=labels['lam30'])

ax.set_xlabel('Time (s)', fontsize=12)
ax.set_ylabel('Pressure (kPa)', fontsize=12)
ax.set_title('SPHERIC Test 10 — Oil Lateral Impact (dp=0.004, time-shift allowed)', fontsize=14, fontweight='bold')
ax.set_xlim(0, 7)
ax.set_ylim(-0.2, 7)
ax.legend(fontsize=9, loc='upper right')
ax.grid(alpha=0.2)
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_4way.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_4way.png"), dpi=150, bbox_inches='tight')
print("Saved: spheric_oil_4way.{pdf,png}")

# === Figure 2: Peak zoom 2×2 ===
# Find experiment peaks
peaks_idx, _ = find_peaks(p_exp, prominence=100, distance=10000)
exp_peaks = [(t_exp[i], p_exp[i]) for i in peaks_idx[:4]]
print(f"\nExperiment peaks: {[(f'{t:.2f}s', f'{p/1000:.3f}kPa') for t,p in exp_peaks]}")

fig2, axes = plt.subplots(2, 2, figsize=(13, 9))
window = 0.8  # ±0.8s around experiment peak

results = {}
for pidx, (exp_t, exp_p) in enumerate(exp_peaks):
    ax = axes[pidx // 2][pidx % 2]
    t0, t1 = exp_t - window, exp_t + window

    # Experiment
    m = (t_exp >= t0) & (t_exp <= t1)
    ax.plot(t_exp[m], p_exp[m]/1000, '-', color=colors['exp'], lw=1.5, label='Experiment')

    row = {'exp': exp_p/1000}
    for key, (t_sim, p_sim) in [('art01', (t_art01, p_art01)),
                                  ('art30', (t_art30, p_art30)),
                                  ('lam30', (t_lam30, p_lam30))]:
        ms = (t_sim >= t0) & (t_sim <= t1)
        if np.any(ms):
            ax.plot(t_sim[ms], p_sim[ms]/1000, '-', color=colors[key], lw=1.0, alpha=0.8, label=labels[key])
            peak_val = np.max(p_sim[ms])/1000
            peak_time = t_sim[ms][np.argmax(p_sim[ms])]
            err = abs(peak_val - exp_p/1000) / (exp_p/1000) * 100 if exp_p > 0 else 0
            row[key] = (peak_val, err, peak_time)

    results[f'Peak {pidx+1}'] = row
    ax.set_title(f'Peak {pidx+1} (exp@{exp_t:.2f}s, {exp_p/1000:.3f} kPa)', fontsize=10, fontweight='bold')
    ax.set_xlabel('Time (s)', fontsize=9)
    ax.set_ylabel('Pressure (kPa)', fontsize=9)
    ax.grid(alpha=0.2)
    if pidx == 0:
        ax.legend(fontsize=7, loc='upper right')

    # Annotate errors
    y_pos = 0.95
    for key, col in [('art01', colors['art01']), ('art30', colors['art30']), ('lam30', colors['lam30'])]:
        if key in row and isinstance(row[key], tuple):
            _, err, _ = row[key]
            ax.text(0.02, y_pos, f'{labels[key]}: {err:.1f}%', transform=ax.transAxes,
                    fontsize=7, va='top', color=col, fontweight='bold')
            y_pos -= 0.08

plt.suptitle('SPHERIC T10 Oil — Peak Comparison (time-shift allowed, dp=0.004)', fontsize=13, fontweight='bold')
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_4way_peaks.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_4way_peaks.png"), dpi=150, bbox_inches='tight')
print("Saved: spheric_oil_4way_peaks.{pdf,png}")

# === Print summary table ===
print("\n" + "=" * 80)
print("SPHERIC T10 Oil — Time-Shifted Peak Error Summary (dp=0.004)")
print("=" * 80)
print(f"{'Peak':>8} | {'Experiment':>15} | {'Art α=0.1':>15} | {'Art α=30':>15} | {'Lam+SPS':>15}")
print("-" * 80)

all_errs = {'art01': [], 'art30': [], 'lam30': []}
for name, row in results.items():
    exp_str = f"{row['exp']:.3f} kPa"
    parts = [f"{name:>8}", f"{exp_str:>15}"]
    for key in ['art01', 'art30', 'lam30']:
        if key in row and isinstance(row[key], tuple):
            val, err, _ = row[key]
            parts.append(f"{val:.3f}kPa {err:>5.1f}%")
            all_errs[key].append(err)
        else:
            parts.append(f"{'N/A':>15}")
    print(" | ".join(parts))

print("-" * 80)
avg_parts = [f"{'Average':>8}", f"{'':>15}"]
for key in ['art01', 'art30', 'lam30']:
    if all_errs[key]:
        avg_parts.append(f"{'':>9}{np.mean(all_errs[key]):>5.1f}%")
    else:
        avg_parts.append(f"{'N/A':>15}")
print(" | ".join(avg_parts))
print("=" * 80)
