#!/usr/bin/env python3
"""SPHERIC T10 Oil — 2D Parametric Study (ρ × α) with Time-Shifting
Grid: ρ = {880, 900, 920, 940, 960, 980} × α = {0.5, 0.7, 0.8, 0.9, 1.0, 1.2, 1.5}
Probe: x=0.007 (dp/2 inward), y=0.031, z=0.093"""
import matplotlib.pyplot as plt
import matplotlib
import numpy as np
import os
from scipy.signal import find_peaks
from scipy.optimize import minimize_scalar

matplotlib.use('Agg')

script_dir = os.path.dirname(os.path.abspath(__file__))
sim_base = os.path.join(script_dir, '../../simulations')
EXP_FILE = os.path.join(script_dir, '../../paper-pof/data/spheric/lateral_oil_1x.txt')
OUT_DIR = os.path.join(script_dir, '../figures/organized/expc_spheric')
os.makedirs(OUT_DIR, exist_ok=True)

# ── Load functions ──
def load_csv(path):
    t, p = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i < 4: continue
            parts = line.strip().split(';')
            if len(parts) >= 3:
                t.append(float(parts[1]))
                p.append(float(parts[2]))
    return np.array(t), np.array(p)

def load_exp(path):
    t, p = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i == 0: continue
            parts = line.strip().split('\t')
            if len(parts) >= 2:
                try:
                    t.append(float(parts[0]))
                    p.append(float(parts[1]) * 100.0)
                except: continue
    return np.array(t), np.array(p)

# ── Time-shifting: minimize L2 error over time window ──
def compute_optimal_shift(t_sph, p_sph, t_exp, p_exp, shift_range=(-1.0, 1.0)):
    """Find optimal time shift that minimizes RMS error between SPH and experiment."""
    def rms_error(dt):
        t_shifted = t_sph + dt
        # Interpolate SPH onto experiment time grid (within overlap)
        mask = (t_exp >= max(t_shifted[0], t_exp[0])) & (t_exp <= min(t_shifted[-1], t_exp[-1]))
        if np.sum(mask) < 10:
            return 1e6
        p_interp = np.interp(t_exp[mask], t_shifted, p_sph)
        return np.sqrt(np.mean((p_interp - p_exp[mask])**2))

    result = minimize_scalar(rms_error, bounds=shift_range, method='bounded')
    return result.x, result.fun

def compute_peak_error_shifted(t_sph, p_sph, t_exp, p_exp, dt):
    """Compute peak errors after time-shifting."""
    exp_idx, _ = find_peaks(p_exp, prominence=100, distance=10000)
    exp_peaks = [(t_exp[i], p_exp[i]/1000) for i in exp_idx[:4]]

    t_shifted = t_sph + dt
    sph_idx, _ = find_peaks(p_sph, prominence=30, distance=30)
    sph_peaks = sorted([(t_shifted[i], p_sph[i]/1000) for i in sph_idx if p_sph[i] > 30],
                       key=lambda x: x[0])

    # Match peaks by proximity
    errs = []
    for et, ep in exp_peaks:
        # Find closest SPH peak within ±0.5s
        candidates = [(st, sp) for st, sp in sph_peaks if abs(st - et) < 0.5]
        if candidates:
            best = min(candidates, key=lambda x: abs(x[0] - et))
            err = abs(best[1] - ep) / ep * 100 if ep > 0 else 100
            errs.append(err)
        else:
            errs.append(100)  # missing peak

    return errs, np.mean(errs) if errs else 100

# ── Load experiment ──
t_exp, p_exp = load_exp(EXP_FILE)

# ── Auto-discover all oil cases ──
import re
data = {}  # (rho, alpha) → (t, p)
pattern = re.compile(r'oil_rho(\d+)_a(\d+_\d+)')
for entry in sorted(os.listdir(sim_base)):
    m = pattern.match(entry)
    if not m: continue
    rho = int(m.group(1))
    alpha = float(m.group(2).replace('_', '.'))
    path = os.path.join(sim_base, entry, 'out', 'measure_lateral_Press.csv')
    if os.path.exists(path):
        t, p = load_csv(path)
        if len(t) > 50:
            data[(rho, alpha)] = (t, p)

rhos = sorted(set(r for r, _ in data.keys()))
alphas = sorted(set(a for _, a in data.keys()))

print(f"Loaded: {len(data)}/{len(rhos)*len(alphas)} cases")
missing = [(r, a) for r in rhos for a in alphas if (r, a) not in data]
if missing:
    print(f"Missing: {missing}")

# ── Compute metrics for each case ──
results = {}  # (rho, alpha) → {dt, rms, peak_errs, mean_peak_err}
for (rho, a), (t_sph, p_sph) in data.items():
    dt, rms = compute_optimal_shift(t_sph, p_sph, t_exp, p_exp)
    peak_errs, mean_peak = compute_peak_error_shifted(t_sph, p_sph, t_exp, p_exp, dt)
    results[(rho, a)] = {
        'dt': dt, 'rms': rms, 'peak_errs': peak_errs, 'mean_peak': mean_peak
    }

# ── Find best ──
best_key = min(results.keys(), key=lambda k: results[k]['mean_peak'])
best = results[best_key]
print(f"\nBest: ρ={best_key[0]}, α={best_key[1]} → {best['mean_peak']:.1f}% (Δt={best['dt']:+.3f}s)")

# ═══════════════════════════════════════════════════
# Figure 1: 2D Heatmap — Mean Peak Error (ρ × α)
# ═══════════════════════════════════════════════════
fig1, (ax1a, ax1b) = plt.subplots(1, 2, figsize=(16, 6))

# Build matrix
err_matrix = np.full((len(rhos), len(alphas)), np.nan)
dt_matrix = np.full((len(rhos), len(alphas)), np.nan)
for i, rho in enumerate(rhos):
    for j, a in enumerate(alphas):
        if (rho, a) in results:
            err_matrix[i, j] = results[(rho, a)]['mean_peak']
            dt_matrix[i, j] = results[(rho, a)]['dt']

# Error heatmap
im1 = ax1a.imshow(err_matrix, cmap='RdYlGn_r', aspect='auto',
                   vmin=np.nanmin(err_matrix), vmax=np.nanpercentile(err_matrix, 90))
ax1a.set_xticks(range(len(alphas)))
ax1a.set_xticklabels([f'{a}' for a in alphas])
ax1a.set_yticks(range(len(rhos)))
ax1a.set_yticklabels([str(r) for r in rhos])
ax1a.set_xlabel('Artificial Viscosity α', fontsize=12)
ax1a.set_ylabel('Density ρ (kg/m³)', fontsize=12)
ax1a.set_title('Mean Peak Error (%) — Time-Shifted', fontsize=13, fontweight='bold')
for i in range(len(rhos)):
    for j in range(len(alphas)):
        if not np.isnan(err_matrix[i, j]):
            color = 'white' if err_matrix[i, j] > 30 else 'black'
            marker = ' ★' if (rhos[i], alphas[j]) == best_key else ''
            ax1a.text(j, i, f'{err_matrix[i, j]:.0f}%{marker}',
                     ha='center', va='center', fontsize=9, fontweight='bold', color=color)
plt.colorbar(im1, ax=ax1a, label='Mean Peak Error (%)')

# Time-shift heatmap
im2 = ax1b.imshow(dt_matrix, cmap='coolwarm', aspect='auto',
                   vmin=-0.5, vmax=0.5)
ax1b.set_xticks(range(len(alphas)))
ax1b.set_xticklabels([f'{a}' for a in alphas])
ax1b.set_yticks(range(len(rhos)))
ax1b.set_yticklabels([str(r) for r in rhos])
ax1b.set_xlabel('Artificial Viscosity α', fontsize=12)
ax1b.set_ylabel('Density ρ (kg/m³)', fontsize=12)
ax1b.set_title('Optimal Time Shift Δt (s)', fontsize=13, fontweight='bold')
for i in range(len(rhos)):
    for j in range(len(alphas)):
        if not np.isnan(dt_matrix[i, j]):
            ax1b.text(j, i, f'{dt_matrix[i, j]:+.2f}',
                     ha='center', va='center', fontsize=8, color='black')
plt.colorbar(im2, ax=ax1b, label='Δt (s)')

plt.suptitle(f'SPHERIC T10 Oil — 2D Parametric (ρ×α), Best: ρ={best_key[0]} α={best_key[1]} ({best["mean_peak"]:.1f}%)',
             fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_heatmap.pdf'), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_heatmap.png'), dpi=150, bbox_inches='tight')
print("Saved: oil_2d_parametric_heatmap")

# ═══════════════════════════════════════════════════
# Figure 2: Best case — time-shifted timeseries overlay
# ═══════════════════════════════════════════════════
fig2, (ax2a, ax2b) = plt.subplots(2, 1, figsize=(16, 10))

# Top 5 cases
sorted_results = sorted(results.items(), key=lambda x: x[1]['mean_peak'])
top5 = sorted_results[:5]
cmap5 = plt.cm.tab10

# Raw
ax2a.plot(t_exp, p_exp/1000, 'k-', lw=2.5, label='Experiment', zorder=20)
for idx, ((rho, a), r) in enumerate(top5):
    t, p = data[(rho, a)]
    ax2a.plot(t, p/1000, '-', color=cmap5(idx), lw=1.2, alpha=0.8,
              label=f'ρ={rho} α={a} ({r["mean_peak"]:.0f}%)')
ax2a.set_ylabel('Pressure (kPa)', fontsize=12)
ax2a.set_title('Raw (no time-shift)', fontsize=12, fontweight='bold')
ax2a.set_xlim(0, 7)
ax2a.legend(fontsize=9, ncol=2)
ax2a.grid(alpha=0.2)

# Time-shifted
ax2b.plot(t_exp, p_exp/1000, 'k-', lw=2.5, label='Experiment', zorder=20)
for idx, ((rho, a), r) in enumerate(top5):
    t, p = data[(rho, a)]
    dt = r['dt']
    ax2b.plot(t + dt, p/1000, '-', color=cmap5(idx), lw=1.2, alpha=0.8,
              label=f'ρ={rho} α={a} (Δt={dt:+.2f}s, {r["mean_peak"]:.0f}%)')
ax2b.set_xlabel('Time (s)', fontsize=12)
ax2b.set_ylabel('Pressure (kPa)', fontsize=12)
ax2b.set_title('Time-Shifted (optimal Δt per case)', fontsize=12, fontweight='bold')
ax2b.set_xlim(0, 7)
ax2b.legend(fontsize=9, ncol=2)
ax2b.grid(alpha=0.2)

plt.suptitle('SPHERIC T10 Oil — Top 5 Configurations (ρ×α)', fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_timeseries.pdf'), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_timeseries.png'), dpi=150, bbox_inches='tight')
print("Saved: oil_2d_parametric_timeseries")

# ═══════════════════════════════════════════════════
# Figure 3: Per-density α sweep (line plot)
# ═══════════════════════════════════════════════════
fig3, ax3 = plt.subplots(figsize=(12, 6))
rho_colors = plt.cm.viridis(np.linspace(0, 1, len(rhos)))
for i, rho in enumerate(rhos):
    a_vals, errs = [], []
    for a in alphas:
        if (rho, a) in results:
            a_vals.append(a)
            errs.append(results[(rho, a)]['mean_peak'])
    if a_vals:
        ax3.plot(a_vals, errs, 'o-', color=rho_colors[i], lw=2, markersize=8,
                 label=f'ρ={rho}')
        # Mark minimum
        best_idx = np.argmin(errs)
        ax3.plot(a_vals[best_idx], errs[best_idx], '*', color=rho_colors[i],
                 markersize=15, zorder=10)

ax3.set_xlabel('Artificial Viscosity α', fontsize=12)
ax3.set_ylabel('Mean Peak Error (%) — Time-Shifted', fontsize=12)
ax3.set_title('Per-Density α Sweep', fontsize=13, fontweight='bold')
ax3.legend(fontsize=10)
ax3.grid(alpha=0.2)
plt.tight_layout()
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_sweep.pdf'), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(OUT_DIR, 'oil_2d_parametric_sweep.png'), dpi=150, bbox_inches='tight')
print("Saved: oil_2d_parametric_sweep")

# ═══════════════════════════════════════════════════
# Print summary table
# ═══════════════════════════════════════════════════
print(f"\n{'='*100}")
print(f"SPHERIC T10 Oil — 2D Parametric Study (ρ × α) with Time-Shifting")
print(f"Probe: (0.007, 0.031, 0.093) — dp/2 inward from wall")
print(f"{'='*100}")
print(f"{'ρ':>5} {'α':>5} | {'Δt':>7} | {'P1':>6} {'P2':>6} {'P3':>6} {'P4':>6} | {'Mean':>6}")
print(f"{'-'*70}")

for (rho, a), r in sorted(results.items()):
    pe = r['peak_errs']
    pe_str = '  '.join([f'{e:>5.0f}%' if e < 100 else f'{"N/A":>6}' for e in pe])
    mark = ' ★' if (rho, a) == best_key else ''
    print(f" {rho:>4} {a:>5.1f} | {r['dt']:>+6.3f}s | {pe_str} | {r['mean_peak']:>5.1f}%{mark}")

print(f"{'='*100}")
print(f"Best: ρ={best_key[0]}, α={best_key[1]} → {best['mean_peak']:.1f}% (Δt={best['dt']:+.3f}s)")
