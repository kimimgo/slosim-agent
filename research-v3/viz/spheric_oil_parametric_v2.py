#!/usr/bin/env python3
"""SPHERIC T10 Oil — Parametric study v2: independent peak detection (true time-shifting)"""
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

def find_sph_peaks(t, p, n=4, prom=50, dist=50, min_kpa=0.05):
    idx, _ = find_peaks(p, prominence=prom, distance=dist)
    peaks = sorted([(t[i], p[i]/1000) for i in idx], key=lambda x: x[0])
    return [(tp, pp) for tp, pp in peaks if pp > min_kpa][:n]

script_dir = os.path.dirname(os.path.abspath(__file__))
sim_base = os.path.join(script_dir, "../../simulations")
EXP_FILE = os.path.join(script_dir, "../../paper-pof/data/spheric/lateral_oil_1x.txt")
OUT_DIR = os.path.join(script_dir, "../figures/organized/expc_spheric")
os.makedirs(OUT_DIR, exist_ok=True)

# Load experiment
t_exp, p_exp = load_exp(EXP_FILE)
exp_idx, _ = find_peaks(p_exp, prominence=100, distance=10000)
exp_peaks = [(t_exp[i], p_exp[i]/1000) for i in exp_idx[:4]]

# All cases
case_defs = [
    ("0.1",  f"{sim_base}/spheric_oil_low_art/out/measure_lateral_Press.csv"),
    ("1.5",  f"{sim_base}/spheric_oil_art_a1_5/out/measure_lateral_Press.csv"),
    ("2.0",  f"{sim_base}/spheric_oil_art_a2_0/out/measure_lateral_Press.csv"),
    ("2.5",  f"{sim_base}/spheric_oil_art_a2_5/out/measure_lateral_Press.csv"),
    ("3.0",  f"{sim_base}/spheric_oil_art_a3_0/out/measure_lateral_Press.csv"),
    ("4.0",  f"{sim_base}/spheric_oil_art_a4_0/out/measure_lateral_Press.csv"),
    ("5.0",  f"{sim_base}/spheric_oil_art_a5_0/out/measure_lateral_Press.csv"),
    ("10.0", f"{sim_base}/spheric_oil_art_a10_0/out/measure_lateral_Press.csv"),
    ("1.0",  f"{sim_base}/spheric_oil_art_a1_0/out/measure_lateral_Press.csv"),
    ("1.1",  f"{sim_base}/spheric_oil_art_a1_1/out/measure_lateral_Press.csv"),
    ("1.2",  f"{sim_base}/spheric_oil_art_a1_2/out/measure_lateral_Press.csv"),
    ("1.3",  f"{sim_base}/spheric_oil_art_a1_3/out/measure_lateral_Press.csv"),
    ("1.4",  f"{sim_base}/spheric_oil_art_a1_4/out/measure_lateral_Press.csv"),
]

sph = {}
sph_peaks_all = {}
for alpha, path in case_defs:
    if os.path.exists(path):
        t, p = load_csv(path)
        sph[alpha] = (t, p)
        sph_peaks_all[alpha] = find_sph_peaks(t, p)

alpha_list = ["0.1", "1.0", "1.1", "1.2", "1.3", "1.4", "1.5", "2.0", "2.5", "3.0", "4.0", "5.0", "10.0"]
cmap = plt.cm.viridis
colors = {a: cmap(i / (len(alpha_list)-1)) for i, a in enumerate(alpha_list)}

# === Figure 1: Timeseries with peak markers ===
fig, ax = plt.subplots(figsize=(14, 5.5))
ax.plot(t_exp, p_exp/1000, 'k-', lw=2, alpha=0.9, label='Experiment', zorder=10)
for et, ep in exp_peaks:
    ax.plot(et, ep, 'k^', markersize=8, zorder=11)

for alpha in alpha_list:
    if alpha not in sph: continue
    t, p = sph[alpha]
    ax.plot(t, p/1000, '-', color=colors[alpha], lw=1.0, alpha=0.75, label=f'α={alpha}')
    for st, sp in sph_peaks_all[alpha]:
        ax.plot(st, sp, 'v', color=colors[alpha], markersize=5, zorder=9)

ax.set_xlabel('Time (s)', fontsize=12)
ax.set_ylabel('Pressure (kPa)', fontsize=12)
ax.set_title('SPHERIC T10 Oil — Artificial Viscosity α Sweep (dp=0.004, true time-shift)', fontsize=13, fontweight='bold')
ax.set_xlim(0, 7)
ax.legend(fontsize=9, loc='upper right', ncol=2)
ax.grid(alpha=0.2)
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_timeseries.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_timeseries.png"), dpi=150, bbox_inches='tight')
print("Saved: param_v2_timeseries")

# === Figure 2: Error bar chart + per-peak lines ===
results = {}  # alpha → [(peak_idx, sph_kpa, err)]
for alpha in alpha_list:
    if alpha not in sph_peaks_all: continue
    sp = sph_peaks_all[alpha]
    errs = []
    for i in range(4):
        if i < len(sp):
            err = abs(sp[i][1] - exp_peaks[i][1]) / exp_peaks[i][1] * 100
            errs.append((i+1, sp[i][1], err, sp[i][0]))
        else:
            errs.append((i+1, 0, 100, 0))
    results[alpha] = errs

fig2, (ax2a, ax2b) = plt.subplots(1, 2, figsize=(13, 5.5))

# Mean error bar
mean_errs = []
for alpha in alpha_list:
    if alpha in results:
        mean_errs.append(np.mean([e for _, _, e, _ in results[alpha]]))
    else:
        mean_errs.append(0)

bar_colors = [colors[a] for a in alpha_list]
bars = ax2a.bar(alpha_list, mean_errs, color=bar_colors, edgecolor='black', linewidth=0.5)
ax2a.set_xlabel('Artificial Viscosity α', fontsize=12)
ax2a.set_ylabel('Mean Peak Error (%)', fontsize=12)
ax2a.set_title('Mean Error (independent peak matching)', fontsize=12, fontweight='bold')
ax2a.grid(axis='y', alpha=0.2)
best_idx = np.argmin(mean_errs)
for i, (a, e) in enumerate(zip(alpha_list, mean_errs)):
    marker = ' ★' if i == best_idx else ''
    ax2a.text(i, e + 1.5, f'{e:.0f}%{marker}', ha='center', fontsize=10, fontweight='bold')

# Per-peak lines
for pidx in range(4):
    vals = []
    for alpha in alpha_list:
        if alpha in results:
            vals.append(results[alpha][pidx][2])
    ax2b.plot(alpha_list[:len(vals)], vals, 'o-', label=f'Peak {pidx+1}', lw=1.5, markersize=7)

ax2b.set_xlabel('Artificial Viscosity α', fontsize=12)
ax2b.set_ylabel('Peak Error (%)', fontsize=12)
ax2b.set_title('Per-Peak Error (independent matching)', fontsize=12, fontweight='bold')
ax2b.legend(fontsize=9)
ax2b.grid(alpha=0.2)
ax2b.axhline(y=20, color='gray', ls='--', alpha=0.3, lw=0.8)
ax2b.text(4.1, 21, '20%', fontsize=8, color='gray')

plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_error.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_error.png"), dpi=150, bbox_inches='tight')
print("Saved: param_v2_error")

# === Figure 3: Peak-by-peak zoom (show time-shift) ===
fig3, axes = plt.subplots(2, 2, figsize=(13, 9))
for pidx in range(4):
    ax = axes[pidx//2][pidx%2]
    et, ep = exp_peaks[pidx]

    # Wide window centered on average of exp and sph peak times
    all_times = [et]
    for alpha in alpha_list:
        if alpha in results and pidx < len(sph_peaks_all.get(alpha, [])):
            all_times.append(sph_peaks_all[alpha][pidx][0])
    center = np.mean(all_times)
    w = 1.0
    t0, t1 = center - w, center + w

    m = (t_exp >= t0) & (t_exp <= t1)
    ax.plot(t_exp[m], p_exp[m]/1000, 'k-', lw=2, label='Experiment')

    for alpha in alpha_list:
        if alpha not in sph: continue
        t, p = sph[alpha]
        ms = (t >= t0) & (t <= t1)
        if np.any(ms):
            ax.plot(t[ms], p[ms]/1000, '-', color=colors[alpha], lw=1.0, alpha=0.8, label=f'α={alpha}')

    # Annotate errors
    y_pos = 0.95
    for alpha in alpha_list:
        if alpha in results and pidx < len(results[alpha]):
            _, val, err, st = results[alpha][pidx]
            dt = st - et
            ax.text(0.02, y_pos, f'α={alpha}: {err:.0f}% (Δt={dt:+.2f}s)',
                    transform=ax.transAxes, fontsize=7, va='top', color=colors[alpha], fontweight='bold')
            y_pos -= 0.065

    ax.set_title(f'Peak {pidx+1} (exp: {ep:.3f} kPa @ {et:.2f}s)', fontsize=10, fontweight='bold')
    ax.set_xlabel('Time (s)', fontsize=9)
    ax.set_ylabel('Pressure (kPa)', fontsize=9)
    ax.grid(alpha=0.2)
    if pidx == 0:
        ax.legend(fontsize=7, ncol=2, loc='upper right')

plt.suptitle('SPHERIC T10 Oil — Peak Zoom with Time-Shift (Δt shown)', fontsize=13, fontweight='bold')
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_peaks.pdf"), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, "spheric_oil_param_v2_peaks.png"), dpi=150, bbox_inches='tight')
print("Saved: param_v2_peaks")

# === Print final summary ===
print("\n" + "=" * 100)
print("SPHERIC T10 Oil — Parametric Study v2 (Independent Peak Detection, True Time-Shift)")
print("=" * 100)
print(f"{'α':>5} | {'Peak1':>14} {'Δt':>6} | {'Peak2':>14} {'Δt':>6} | {'Peak3':>14} {'Δt':>6} | {'Peak4':>14} {'Δt':>6} | {'Mean':>5}")
print("-" * 100)
for alpha in alpha_list:
    if alpha not in results: continue
    parts = []
    for pidx, (_, val, err, st) in enumerate(results[alpha]):
        dt = st - exp_peaks[pidx][0]
        parts.append(f"{val:.3f} {err:>4.0f}% {dt:>+5.2f}s")
    mean_e = np.mean([e for _, _, e, _ in results[alpha]])
    best = " ★" if alpha == alpha_list[best_idx] else ""
    print(f"  {alpha:>3} | {' | '.join(parts)} | {mean_e:>4.1f}%{best}")
print("=" * 100)
