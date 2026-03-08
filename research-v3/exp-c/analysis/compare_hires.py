#!/usr/bin/env python3
"""Compare run_005 (Artificial, 0.1s) vs run_009 (Laminar+SPS, 0.1s) vs
run_009b (Laminar+SPS, 10ms) pressure at z=0.093m.
Shows high-resolution peak shape from run_009b."""
import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from pathlib import Path

ANALYSIS_DIR = Path(__file__).parent

# Experimental peaks (SPHERIC Test 10, Oil Lateral, z=0.093m)
EXP_PEAKS = [
    {'time': 1.54, 'pressure': 2200, 'label': 'Peak 1'},
    {'time': 3.07, 'pressure': 3800, 'label': 'Peak 2'},
    {'time': 4.61, 'pressure': 5500, 'label': 'Peak 3'},
    {'time': 6.14, 'pressure': 6200, 'label': 'Peak 4'},
]


def load_press_csv(path, probe_idx=0):
    """Load DualSPHysics MeasureTool pressure CSV (semicolon-delimited)."""
    times, pressures = [], []
    with open(path) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith('#') or line.startswith('Part'):
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


def find_peak(times, pressures, t_center, window=0.5):
    """Find peak pressure within a time window."""
    mask = (times >= t_center - window) & (times <= t_center + window)
    if mask.sum() == 0:
        return 0, t_center
    idx = np.argmax(pressures[mask])
    t_peak = times[mask][idx]
    p_peak = pressures[mask][idx]
    return p_peak, t_peak


# Load data
runs = {}
SIMDATA = Path('/mnt/simdata/dualsphysics/exp1')

# run_005 (Artificial, 0.1s resolution)
f005 = SIMDATA / 'run_005' / 'PressConsistent_Press.csv'
if f005.exists():
    t005, p005 = load_press_csv(str(f005))
    runs['run_005'] = (t005, p005, 'Artificial (α=0.1), 0.1s', '#2196F3', '-')

# run_009 (Laminar+SPS, 0.1s resolution)
f009 = SIMDATA / 'run_009' / 'PressConsistent_Press.csv'
if f009.exists():
    t009, p009 = load_press_csv(str(f009))
    runs['run_009'] = (t009, p009, 'Laminar+SPS, 0.1s', '#FF9800', 'o-')

# run_009b (Laminar+SPS, 10ms resolution) — from MeasureTool after solver
f009b = ANALYSIS_DIR / 'run_009b' / 'pressure_009b.csv'
if f009b.exists():
    t009b, p009b = load_press_csv(str(f009b))
    runs['run_009b'] = (t009b, p009b, 'Laminar+SPS, 10ms (hi-res)', '#F44336', '-')

if not runs:
    print("No data files found!")
    print(f"Checked: {f005}, {f009}, {f009b}")
    exit(1)

print(f"Loaded {len(runs)} runs: {list(runs.keys())}")
for name, (t, p, *_) in runs.items():
    print(f"  {name}: {len(t)} points, t=[{t[0]:.3f}, {t[-1]:.3f}]")

# ================================================================
# Figure: Time series comparison (all runs + experimental peaks)
# ================================================================
fig, axes = plt.subplots(2, 1, figsize=(14, 10), height_ratios=[3, 1])

ax = axes[0]
for name, (t, p, label, color, style) in runs.items():
    if 'o' in style:
        ax.plot(t, p, style, color=color, label=label, markersize=3, linewidth=0.8, alpha=0.8)
    else:
        ax.plot(t, p, style, color=color, label=label, linewidth=0.8, alpha=0.9)

# Experimental peaks as markers
for peak in EXP_PEAKS:
    ax.plot(peak['time'], peak['pressure'], 'g^', markersize=12, zorder=10)
    ax.annotate(f"{peak['pressure']} Pa", (peak['time'], peak['pressure']),
                textcoords='offset points', xytext=(10, 5), fontsize=8, color='green')

ax.set_ylabel('Pressure (Pa)', fontsize=12)
ax.set_title('Oil Lateral — Pressure at z=0.093m: Resolution Comparison', fontsize=13)
ax.legend(loc='upper left', fontsize=10)
ax.grid(True, alpha=0.3)
ax.set_xlim(0, 7)
ax.set_ylim(-200, 7000)

# Peak error bar chart
ax2 = axes[1]
peak_data = {}
for name, (t, p, label, color, style) in runs.items():
    errors = []
    for peak in EXP_PEAKS:
        p_sim, _ = find_peak(t, p, peak['time'])
        err = abs(p_sim - peak['pressure']) / peak['pressure'] * 100
        errors.append(err)
    peak_data[name] = (errors, color, label)
    print(f"  {name} peak errors: {[f'{e:.1f}%' for e in errors]}, mean={np.mean(errors):.1f}%")

x = np.arange(len(EXP_PEAKS))
w = 0.8 / len(peak_data)
for i, (name, (errors, color, label)) in enumerate(peak_data.items()):
    ax2.bar(x + i*w - 0.4 + w/2, errors, w, label=label, color=color, alpha=0.7)

ax2.set_xticks(x)
ax2.set_xticklabels([p['label'] for p in EXP_PEAKS])
ax2.set_ylabel('Peak Error (%)', fontsize=11)
ax2.set_title('Peak-by-Peak Error Comparison', fontsize=12)
ax2.legend(loc='upper right', fontsize=9)
ax2.grid(axis='y', alpha=0.3)

fig.tight_layout()
out_path = ANALYSIS_DIR / 'hires_comparison.png'
fig.savefig(out_path, dpi=150)
print(f"\nSaved: {out_path}")
plt.close()
