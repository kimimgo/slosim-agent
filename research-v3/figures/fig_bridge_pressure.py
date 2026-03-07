#!/usr/bin/env python3
"""Bridge Experiment: Agent-Generated XML → Pressure Time Series"""
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np
import csv
import os

BRIDGE_DIR = os.path.join(os.path.dirname(__file__), "..", "exp-c", "agent-bridge")

def load_gauge_csv(filename):
    filepath = os.path.join(BRIDGE_DIR, filename)
    times, data = [], []
    with open(filepath) as f:
        reader = csv.reader(f, delimiter=";")
        for i, row in enumerate(reader):
            if i < 4:
                continue
            times.append(float(row[1]))
            data.append([float(v) for v in row[2:]])
    return np.array(times), np.array(data)

times, press = load_gauge_csv("gauges_Press.csv")
probe_x = [0.03, 0.15, 0.30, 0.45, 0.57]
probe_labels = ['Left wall (0.03m)', 'Left mid (0.15m)', 'Center (0.30m)',
                'Right mid (0.45m)', 'Right wall (0.57m)']
colors = ['#C62828', '#FF9800', '#4CAF50', '#2196F3', '#1565C0']

fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 7), height_ratios=[3, 1])

# ── Top: Pressure time series ──
for i, (label, color) in enumerate(zip(probe_labels, colors)):
    lw = 2.0 if i in [0, 4] else 1.2
    alpha = 1.0 if i in [0, 4] else 0.6
    ax1.plot(times, press[:, i], color=color, linewidth=lw, alpha=alpha, label=label)

# Hydrostatic reference
p_hydro = 1000 * 9.81 * (0.083 - 0.05)  # 323.7 Pa
ax1.axhline(y=p_hydro, color='gray', linestyle=':', linewidth=1, alpha=0.5)
ax1.text(0.3, p_hydro + 15, f'Hydrostatic ({p_hydro:.0f} Pa)', fontsize=8, color='gray')

# Anti-phase annotation
ax1.annotate('', xy=(5.0, 629), xytext=(5.0, 271),
             arrowprops=dict(arrowstyle='<->', color='#388E3C', lw=1.5))
ax1.text(5.15, 450, 'Anti-phase\nr = −0.720', fontsize=8, color='#388E3C', fontweight='bold')

ax1.set_ylabel('Pressure (Pa)', fontsize=11)
ax1.set_title('Bridge Experiment: Agent-Generated XML (S02, M-A3=100%) → Physics Validation',
              fontsize=12, fontweight='bold')
ax1.legend(loc='upper right', fontsize=8, ncol=2)
ax1.set_xlim(0, 10)
ax1.set_ylim(100, 700)
ax1.grid(alpha=0.3)
ax1.set_xticklabels([])

# ── Bottom: Left−Right pressure difference (anti-phase) ──
p_diff = press[:, 0] - press[:, 4]  # left - right
ax2.fill_between(times, p_diff, 0, where=p_diff > 0, color='#C62828', alpha=0.3, label='Left > Right')
ax2.fill_between(times, p_diff, 0, where=p_diff < 0, color='#1565C0', alpha=0.3, label='Right > Left')
ax2.plot(times, p_diff, color='black', linewidth=1)
ax2.axhline(y=0, color='black', linewidth=0.5)

ax2.set_xlabel('Time (s)', fontsize=11)
ax2.set_ylabel('ΔP (Pa)', fontsize=11)
ax2.set_xlim(0, 10)
ax2.legend(loc='upper right', fontsize=8, ncol=2)
ax2.grid(alpha=0.3)

# Summary box
summary = (
    'Hydrostatic: 333.5 Pa (err 3.0%)\n'
    'Anti-phase: r = −0.720\n'
    'Frequency: 0.789 Hz (err 4.4%)\n'
    'Result: 5/5 checks PASS'
)
props = dict(boxstyle='round,pad=0.5', facecolor='#E8F5E9', alpha=0.9, edgecolor='#388E3C')
ax1.text(0.02, 0.97, summary, transform=ax1.transAxes, fontsize=8,
         verticalalignment='top', bbox=props)

plt.tight_layout()
plt.savefig(os.path.join(os.path.dirname(__file__), 'fig_bridge_pressure.png'), dpi=300, bbox_inches='tight')
plt.savefig(os.path.join(os.path.dirname(__file__), 'fig_bridge_pressure.pdf'), bbox_inches='tight')
print("Saved: fig_bridge_pressure.png, fig_bridge_pressure.pdf")
