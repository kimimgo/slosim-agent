#!/usr/bin/env python3
"""Compare run_005 (Artificial) vs run_009 (Laminar+SPS) pressure peaks.
Oil Lateral SPHERIC Test 10: probe at x=0.005m, z=0.093m (fill height).
"""
import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

# Experimental peak values from Souto-Iglesias (2011) — Oil Lateral
# Peaks at z=0.093m (fill height sensor), first 4 impact events
EXP_PEAKS = {
    'peak1': {'time': 1.54, 'pressure': 2200},  # Pa
    'peak2': {'time': 3.07, 'pressure': 3800},
    'peak3': {'time': 4.61, 'pressure': 5500},
    'peak4': {'time': 6.14, 'pressure': 6200},
}

def load_press_csv(path, probe_idx=0):
    """Load DualSPHysics MeasureTool pressure CSV."""
    times, pressures = [], []
    with open(path) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith(';') or line.startswith(' ') or line.startswith('Part'):
                continue
            parts = line.split(';')
            if len(parts) < 3 + probe_idx:
                continue
            try:
                t = float(parts[1])
                p = float(parts[2 + probe_idx])
                times.append(t)
                pressures.append(p)
            except (ValueError, IndexError):
                continue
    return np.array(times), np.array(pressures)

def find_peaks_in_windows(times, pressures, windows):
    """Find peak pressure in each time window."""
    peaks = {}
    for name, w in windows.items():
        mask = (times >= w[0]) & (times <= w[1])
        if mask.sum() == 0:
            peaks[name] = {'time': np.nan, 'pressure': 0}
            continue
        idx = np.argmax(pressures[mask])
        t_peak = times[mask][idx]
        p_peak = pressures[mask][idx]
        peaks[name] = {'time': t_peak, 'pressure': p_peak}
    return peaks

# Peak search windows (±0.3s around expected peak times)
WINDOWS = {
    'peak1': (1.2, 1.9),
    'peak2': (2.7, 3.5),
    'peak3': (4.2, 5.0),
    'peak4': (5.8, 6.6),
}

# Load data
base = '/mnt/simdata/dualsphysics/exp1'
t5, p5 = load_press_csv(f'{base}/run_005/PressConsistent_Press.csv', probe_idx=0)
t9, p9 = load_press_csv(f'{base}/run_009/PressConsistent_Press.csv', probe_idx=0)

print(f"run_005 (Artificial): {len(t5)} points, dt≈{np.median(np.diff(t5)):.4f}s")
print(f"run_009 (Laminar+SPS): {len(t9)} points, dt≈{np.median(np.diff(t9)):.4f}s")
print()

# Find peaks
peaks5 = find_peaks_in_windows(t5, p5, WINDOWS)
peaks9 = find_peaks_in_windows(t9, p9, WINDOWS)

print(f"{'Peak':<8} {'Exp (Pa)':<10} {'run005 (Pa)':<12} {'err005':<8} {'run009 (Pa)':<12} {'err009':<8} {'Improved?'}")
print("-" * 80)
for name in ['peak1', 'peak2', 'peak3', 'peak4']:
    exp_p = EXP_PEAKS[name]['pressure']
    p5_val = peaks5[name]['pressure']
    p9_val = peaks9[name]['pressure']
    err5 = (p5_val - exp_p) / exp_p * 100
    err9 = (p9_val - exp_p) / exp_p * 100
    improved = "YES" if abs(err9) < abs(err5) else "NO"
    print(f"{name:<8} {exp_p:<10.0f} {p5_val:<12.1f} {err5:<8.1f}% {p9_val:<12.1f} {err9:<8.1f}% {improved}")

# Summary
print()
avg_err5 = np.mean([abs((peaks5[n]['pressure'] - EXP_PEAKS[n]['pressure'])/EXP_PEAKS[n]['pressure']*100) for n in EXP_PEAKS])
avg_err9 = np.mean([abs((peaks9[n]['pressure'] - EXP_PEAKS[n]['pressure'])/EXP_PEAKS[n]['pressure']*100) for n in EXP_PEAKS])
print(f"Average absolute error: run_005={avg_err5:.1f}%, run_009={avg_err9:.1f}%")
print(f"Improvement: {avg_err5 - avg_err9:.1f}pp")

# Plot comparison
fig, axes = plt.subplots(2, 1, figsize=(12, 8), sharex=True)

# Top: time series
ax = axes[0]
ax.plot(t5, p5, 'b-', alpha=0.7, linewidth=0.5, label='run_005 (Artificial, α=0.1)')
ax.plot(t9, p9, 'r-o', markersize=3, alpha=0.8, linewidth=1.0, label='run_009 (Laminar+SPS)')
for name, ep in EXP_PEAKS.items():
    ax.axvline(ep['time'], color='green', linestyle='--', alpha=0.3)
    ax.plot(ep['time'], ep['pressure'], 'g^', markersize=10, zorder=5)
ax.set_ylabel('Pressure (Pa)')
ax.set_title('Oil Lateral — Pressure at z=0.093m: Artificial vs Laminar+SPS')
ax.legend()
ax.grid(True, alpha=0.3)

# Bottom: peak comparison bar chart
ax = axes[1]
peak_names = list(EXP_PEAKS.keys())
x = np.arange(len(peak_names))
w = 0.25
exp_vals = [EXP_PEAKS[n]['pressure'] for n in peak_names]
r5_vals = [peaks5[n]['pressure'] for n in peak_names]
r9_vals = [peaks9[n]['pressure'] for n in peak_names]

ax.bar(x - w, exp_vals, w, label='Experiment', color='green', alpha=0.7)
ax.bar(x, r5_vals, w, label='run_005 (Artificial)', color='blue', alpha=0.7)
ax.bar(x + w, r9_vals, w, label='run_009 (Laminar+SPS)', color='red', alpha=0.7)
ax.set_xticks(x)
ax.set_xticklabels([f'Peak {i+1}' for i in range(4)])
ax.set_ylabel('Peak Pressure (Pa)')
ax.legend()
ax.grid(True, alpha=0.3, axis='y')

plt.tight_layout()
out_path = '/home/imgyu/workspace/02_active/sim/slosim-agent/research-v3/exp-c/analysis/viscosity_comparison.png'
plt.savefig(out_path, dpi=150)
print(f"\nFigure saved: {out_path}")
