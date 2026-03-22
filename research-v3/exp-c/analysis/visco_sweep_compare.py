#!/usr/bin/env python3
"""Oil Visco Sweep Quick Comparison
Compare pressure peaks at left wall (x=0.005, z=0.093) across alpha values.
Also include existing Artificial α=0.1 and Laminar+SPS ν=4.55e-5 data.
"""
import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

SWEEP_DIR = '/home/imgyu/workspace/02_active/sim/slosim-agent/simulations/exp-c/visco_sweep'
SIMDATA = '/home/imgyu/workspace/02_active/sim/slosim-agent/simulations/exp-c'

def load_csv(path, probe_idx=0):
    """Load DualSPHysics MeasureTool CSV, return (time, pressure) at probe_idx"""
    data = []
    with open(path) as f:
        for line in f:
            if line.startswith((' ', ';', 'Part')):
                if line.strip().startswith('Part'):
                    continue  # skip header
                continue
            parts = line.strip().split(';')
            if len(parts) >= probe_idx + 3:
                try:
                    t = float(parts[1])
                    p = float(parts[probe_idx + 2])
                    data.append((t, p))
                except (ValueError, IndexError):
                    continue
    if not data:
        return np.array([]), np.array([])
    data = np.array(data)
    return data[:, 0], data[:, 1]

def find_peaks_simple(t, p, min_height=100, min_spacing=0.3):
    """Find peaks in pressure signal"""
    peaks = []
    for i in range(1, len(p)-1):
        if p[i] > p[i-1] and p[i] > p[i+1] and p[i] > min_height:
            if not peaks or (t[i] - t[peaks[-1]]) > min_spacing:
                peaks.append(i)
    return peaks

# Experimental reference peaks (SPHERIC Test 10, Oil, from Souto-Iglesias)
# These are approximate from the paper
EXP_PEAKS = {
    'times': [1.0, 2.3],  # approximate peak times
    'pressures': [700, 600],  # approximate peak pressures (Pa)
}

# Load sweep data
alphas = [0.5, 1.0, 2.0, 5.0]
colors = ['#E91E63', '#FF9800', '#4CAF50', '#2196F3']

fig, axes = plt.subplots(2, 2, figsize=(14, 10))
fig.suptitle('Oil Visco Sweep: Artificial α Comparison (dp=0.002, 3s)',
             fontsize=14, fontweight='bold')

# Summary stats
print("=" * 60)
print("Oil Visco Sweep Results (Probe P1: x=0.005, z=0.093)")
print("=" * 60)

peak_data = {}

for idx, (alpha, color) in enumerate(zip(alphas, colors)):
    ax = axes[idx // 2, idx % 2]
    csv_path = f"{SWEEP_DIR}/alpha_{alpha}_Press.csv"
    t, p = load_csv(csv_path, probe_idx=0)

    if len(t) == 0:
        print(f"alpha={alpha}: NO DATA")
        continue

    # Find peaks
    peaks = find_peaks_simple(t, p, min_height=50)
    peak_times = t[peaks] if peaks else []
    peak_pressures = p[peaks] if peaks else []

    # Plot
    ax.plot(t, p, color=color, linewidth=1.5, label=f'α={alpha}')
    if len(peaks) > 0:
        ax.plot(peak_times, peak_pressures, 'v', color='red', markersize=8)
        for pt, pp in zip(peak_times, peak_pressures):
            ax.annotate(f'{pp:.0f} Pa', (pt, pp), textcoords="offset points",
                       xytext=(0, 10), ha='center', fontsize=8, color='red')

    ax.set_xlabel('Time (s)')
    ax.set_ylabel('Pressure (Pa)')
    ax.set_title(f'α = {alpha}')
    ax.grid(True, alpha=0.3)
    ax.set_xlim(0, 3.1)

    # Stats
    max_p = np.max(p) if len(p) > 0 else 0
    mean_p = np.mean(p[p > 0]) if np.any(p > 0) else 0
    n_peaks = len(peaks)

    peak_data[alpha] = {
        'max_p': max_p,
        'mean_p': mean_p,
        'n_peaks': n_peaks,
        'peak_pressures': list(peak_pressures),
        'peak_times': list(peak_times),
    }

    print(f"\nα = {alpha}:")
    print(f"  Max pressure: {max_p:.1f} Pa")
    print(f"  Mean pressure (>0): {mean_p:.1f} Pa")
    print(f"  Number of peaks: {n_peaks}")
    for i, (pt, pp) in enumerate(zip(peak_times, peak_pressures)):
        print(f"  Peak {i+1}: t={pt:.2f}s, P={pp:.0f} Pa")

plt.tight_layout()
outpath = '/home/imgyu/workspace/02_active/sim/slosim-agent/research-v3/exp-c/analysis/visco_sweep_comparison.png'
plt.savefig(outpath, dpi=150, bbox_inches='tight')
print(f"\nPlot saved: {outpath}")

# Summary table
print("\n" + "=" * 60)
print("SUMMARY: α → Peak Pressure Comparison")
print("=" * 60)
print(f"{'Alpha':<8} {'Max P (Pa)':<12} {'Peak 1 (Pa)':<14} {'Peak 2 (Pa)':<14} {'Peaks':<6}")
print("-" * 54)
for alpha in alphas:
    d = peak_data.get(alpha, {})
    p1 = f"{d['peak_pressures'][0]:.0f}" if d.get('peak_pressures') and len(d['peak_pressures']) > 0 else "N/A"
    p2 = f"{d['peak_pressures'][1]:.0f}" if d.get('peak_pressures') and len(d['peak_pressures']) > 1 else "N/A"
    print(f"{alpha:<8} {d.get('max_p', 0):<12.1f} {p1:<14} {p2:<14} {d.get('n_peaks', 0):<6}")
print(f"\nExp. ref: Peak 1 ≈ 700 Pa, Peak 2 ≈ 600 Pa (approximate)")
