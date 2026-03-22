#!/usr/bin/env python3
"""SPHERIC T10 Oil ρ=920 — Full pressure-time timeseries for α=1.0~3.0 (0.1 step)"""
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
                    press.append(float(parts[1]) * 100.0)
                except: continue
    return np.array(times), np.array(press)

script_dir = os.path.dirname(os.path.abspath(__file__))
sim_base = os.path.join(script_dir, '../../simulations')
EXP_FILE = os.path.join(script_dir, '../../paper-pof/data/spheric/lateral_oil_1x.txt')
OUT_DIR = os.path.join(script_dir, '../figures/organized/expc_spheric')
os.makedirs(OUT_DIR, exist_ok=True)

t_exp, p_exp = load_exp(EXP_FILE)

# All alpha values
alphas = [round(i*0.1, 1) for i in range(5, 31)]  # 0.5 ~ 3.0
sph = {}
for a in alphas:
    tag = f"{a:.1f}".replace('.', '_')
    path = f"{sim_base}/oil_rho920_a{tag}/out/measure_lateral_Press.csv"
    if os.path.exists(path):
        sph[a] = load_csv(path)
    else:
        print(f"  WARN: missing α={a} at {path}")

cmap = plt.cm.coolwarm
denom = max(len(alphas) - 1, 1)
colors = {a: cmap(i / denom) for i, a in enumerate(alphas)}

# === Figure 1: Full timeseries — ALL 21 curves ===
fig, ax = plt.subplots(figsize=(16, 6))
ax.plot(t_exp, p_exp/1000, 'k-', lw=2.5, alpha=0.95, label='Experiment', zorder=20)

for a in alphas:
    if a not in sph: continue
    t, p = sph[a]
    ax.plot(t, p/1000, '-', color=colors[a], lw=0.6, alpha=0.6)

# Highlight best few
for a_hi in [0.5, 1.0, 1.4, 2.0, 3.0]:
    if a_hi in sph:
        t, p = sph[a_hi]
        ax.plot(t, p/1000, '-', color=colors[a_hi], lw=1.5, alpha=0.9, label=f'α={a_hi}')

ax.set_xlabel('Time (s)', fontsize=13)
ax.set_ylabel('Pressure (kPa)', fontsize=13)
ax.set_title('SPHERIC T10 Oil (ρ=920) — Pressure vs Time, α=0.5~3.0', fontsize=14, fontweight='bold')
ax.set_xlim(0, 7)
ax.legend(fontsize=10, loc='upper right')
ax.grid(alpha=0.2)

# Colorbar for alpha
sm = plt.cm.ScalarMappable(cmap=cmap, norm=plt.Normalize(0.5, 3.0))
sm.set_array([])
cbar = plt.colorbar(sm, ax=ax, pad=0.01, aspect=30)
cbar.set_label('Artificial Viscosity α', fontsize=11)

plt.tight_layout()
plt.savefig(f'{OUT_DIR}/oil_rho920_timeseries_all.pdf', dpi=300, bbox_inches='tight')
plt.savefig(f'{OUT_DIR}/oil_rho920_timeseries_all.png', dpi=150, bbox_inches='tight')
print("Saved: oil_rho920_timeseries_all")

# === Figure 2: Time-shifted timeseries (top 5) ===
# Exp peak detection: distance in samples, ~80% of sloshing period T=1.535s
dt_exp = np.mean(np.diff(t_exp)) if len(t_exp) > 1 else 0.001
dist_exp = max(int(0.8 * 1.535 / dt_exp), 10)
exp_idx, _ = find_peaks(p_exp, prominence=100, distance=dist_exp)
exp_peaks = [(t_exp[i], p_exp[i]) for i in exp_idx[:4]]

def compute_shift(t_s, p_s, exp_peaks):
    sph_idx, _ = find_peaks(p_s, prominence=50, distance=50)
    sp = sorted([(t_s[i], p_s[i]) for i in sph_idx if p_s[i] > 50], key=lambda x: x[0])[:4]
    shifts = [exp_peaks[i][0] - sp[i][0] for i in range(min(len(exp_peaks), len(sp)))]
    return np.mean(shifts) if shifts else 0

best_alphas = [0.5, 0.8, 1.0, 1.4, 2.0]
fig2, (ax2a, ax2b) = plt.subplots(2, 1, figsize=(16, 10))

# Top: raw
ax2a.plot(t_exp, p_exp/1000, 'k-', lw=2.5, label='Experiment', zorder=20)
for a in best_alphas:
    if a not in sph: continue
    t, p = sph[a]
    ax2a.plot(t, p/1000, '-', color=colors[a], lw=1.2, alpha=0.8, label=f'α={a}')
ax2a.set_ylabel('Pressure (kPa)', fontsize=12)
ax2a.set_title('Raw (no time-shift)', fontsize=12, fontweight='bold')
ax2a.set_xlim(0, 7)
ax2a.legend(fontsize=9, ncol=3)
ax2a.grid(alpha=0.2)

# Bottom: shifted
ax2b.plot(t_exp, p_exp/1000, 'k-', lw=2.5, label='Experiment', zorder=20)
for a in best_alphas:
    if a not in sph: continue
    t, p = sph[a]
    dt = compute_shift(t, p, exp_peaks)
    ax2b.plot(t + dt, p/1000, '-', color=colors[a], lw=1.2, alpha=0.8, label=f'α={a} (Δt={dt:+.2f}s)')
ax2b.set_xlabel('Time (s)', fontsize=12)
ax2b.set_ylabel('Pressure (kPa)', fontsize=12)
ax2b.set_title('Time-shifted (aligned to experiment peaks)', fontsize=12, fontweight='bold')
ax2b.set_xlim(0, 7)
ax2b.legend(fontsize=9, ncol=3)
ax2b.grid(alpha=0.2)

plt.suptitle('SPHERIC T10 Oil (ρ=920) — Raw vs Time-Shifted', fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig(f'{OUT_DIR}/oil_rho920_timeseries_shifted.pdf', dpi=300, bbox_inches='tight')
plt.savefig(f'{OUT_DIR}/oil_rho920_timeseries_shifted.png', dpi=150, bbox_inches='tight')
print("Saved: oil_rho920_timeseries_shifted")

# === Figure 3: Mean error vs α bar chart ===
results = {}
for a in alphas:
    if a not in sph: continue
    t, p = sph[a]
    sph_idx, _ = find_peaks(p, prominence=50, distance=50)
    sp = sorted([(t[i], p[i]/1000) for i in sph_idx if p[i] > 50], key=lambda x: x[0])[:4]
    errs = []
    for i in range(min(4, len(sp))):
        ep = exp_peaks[i][1] / 1000
        err = abs(sp[i][1] - ep) / ep * 100 if ep > 0 else 0
        errs.append(err)
    results[a] = {'peaks': sp, 'errs': errs, 'mean': np.mean(errs) if errs else 100}

fig3, ax3 = plt.subplots(figsize=(14, 5))
a_vals = sorted(results.keys())
means = [results[a]['mean'] for a in a_vals]
bar_colors = [colors[a] for a in a_vals]
bars = ax3.bar([f'{a:.1f}' for a in a_vals], means, color=bar_colors, edgecolor='black', linewidth=0.3)

best_a = a_vals[np.argmin(means)]
for i, (a, m) in enumerate(zip(a_vals, means)):
    marker = ' ★' if a == best_a else ''
    ax3.text(i, m + 1, f'{m:.0f}%{marker}', ha='center', fontsize=7, fontweight='bold')

ax3.set_xlabel('Artificial Viscosity α', fontsize=12)
ax3.set_ylabel('Mean Peak Error (%)', fontsize=12)
ax3.set_title(f'SPHERIC T10 Oil (ρ=920) — Mean Error vs α (best: α={best_a}, {results[best_a]["mean"]:.1f}%)', fontsize=13, fontweight='bold')
ax3.grid(axis='y', alpha=0.2)
plt.tight_layout()
plt.savefig(f'{OUT_DIR}/oil_rho920_error_bar.pdf', dpi=300, bbox_inches='tight')
plt.savefig(f'{OUT_DIR}/oil_rho920_error_bar.png', dpi=150, bbox_inches='tight')
print("Saved: oil_rho920_error_bar")

# === Print table ===
print(f"\n{'='*80}")
print(f"SPHERIC T10 Oil (ρ=920) — α=1.0~3.0 Parametric Study")
print(f"{'='*80}")
print(f"{'α':>5} | {'P1':>8} | {'P2':>8} | {'P3':>8} | {'P4':>8} | {'Mean':>6}")
def fmt_peak(peaks, i):
    return f"{peaks[i][1]/1000:>5.3f}kPa" if i < len(peaks) else f"{'N/A':>8}"
print(f"{'Exp':>5} | {fmt_peak(exp_peaks,0)} | {fmt_peak(exp_peaks,1)} | {fmt_peak(exp_peaks,2)} | {fmt_peak(exp_peaks,3)} |")
print('-'*80)
for a in a_vals:
    r = results[a]
    parts = []
    for i in range(4):
        if i < len(r['errs']):
            parts.append(f"{r['errs'][i]:>5.0f}%")
        else:
            parts.append(f"{'N/A':>6}")
    mark = ' ★' if a == best_a else ''
    print(f" {a:>4.1f} | {'  |  '.join(parts)}  | {r['mean']:>4.1f}%{mark}")
print(f"{'='*80}")
print(f"Best: α={best_a} ({results[best_a]['mean']:.1f}%)")
