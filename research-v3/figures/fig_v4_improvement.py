#!/usr/bin/env python3
"""Figure: v3 vs v4 M-A3 improvement per scenario, grouped by model."""
import json
from pathlib import Path
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

DATA_PATH = Path(__file__).parent.parent / "exp-a" / "v4_analysis_results.json"
with open(DATA_PATH) as f:
    data = json.load(f)

scenarios = [f"S{i:02d}" for i in range(1, 11)]
tiers = {
    "S01": "E", "S02": "E", "S03": "E",
    "S04": "M", "S05": "M", "S06": "M", "S07": "M",
    "S08": "H", "S09": "H", "S10": "H",
}
pitch_scenarios = {"S01", "S04", "S05", "S08", "S09"}

fig, axes = plt.subplots(2, 1, figsize=(12, 8), sharex=True)

for idx, (model_key, label) in enumerate([("32B", "Qwen3-32B"), ("8B", "Qwen3-8B")]):
    ax = axes[idx]
    v3 = [data['models'][model_key]['v3_scores'][s] for s in scenarios]
    v4 = [data['models'][model_key]['v4_scores'][s] for s in scenarios]

    x = np.arange(len(scenarios))
    w = 0.35

    bars_v3 = ax.bar(x - w/2, v3, w, label='v3 (baseline)', color='#90CAF9', edgecolor='#1565C0', linewidth=0.5)
    bars_v4 = ax.bar(x + w/2, v4, w, label='v4 (P1+P2 fix)', color='#EF5350', edgecolor='#B71C1C', linewidth=0.5)

    # Mark pitch scenarios with dagger
    for i, s in enumerate(scenarios):
        if s in pitch_scenarios:
            delta = v4[i] - v3[i]
            if delta > 0:
                ax.annotate(f'+{delta:.0f}', (x[i] + w/2, v4[i] + 1),
                           ha='center', va='bottom', fontsize=8, fontweight='bold', color='#B71C1C')

    ax.set_ylabel('M-A3 (%)', fontsize=11)
    ax.set_title(f'{label}  (v3: {data["models"][model_key]["v3_mean"]:.1f}% → v4: {data["models"][model_key]["v4_mean"]:.1f}%)',
                fontsize=12, fontweight='bold')
    ax.set_ylim(0, 115)
    ax.axhline(y=100, color='gray', linestyle='--', alpha=0.3)
    ax.legend(loc='upper left', fontsize=9)
    ax.grid(axis='y', alpha=0.2)

    # Tier separators
    ax.axvline(x=2.5, color='gray', linestyle=':', alpha=0.4)
    ax.axvline(x=6.5, color='gray', linestyle=':', alpha=0.4)
    ax.text(1, 108, 'Easy', ha='center', fontsize=9, color='gray')
    ax.text(4.5, 108, 'Medium', ha='center', fontsize=9, color='gray')
    ax.text(8, 108, 'Hard', ha='center', fontsize=9, color='gray')

xlabels = [f'{s}\n({tiers[s]})' for s in scenarios]
axes[1].set_xticks(x)
axes[1].set_xticklabels(xlabels, fontsize=9)

fig.suptitle('EXP-A: Tool Design Fix Impact (v3 → v4)', fontsize=14, fontweight='bold', y=0.98)
fig.tight_layout(rect=[0, 0, 1, 0.96])

out_pdf = Path(__file__).parent / 'fig_v4_improvement.pdf'
out_png = Path(__file__).parent / 'fig_v4_improvement.png'
fig.savefig(out_pdf, bbox_inches='tight')
fig.savefig(out_png, dpi=150, bbox_inches='tight')
print(f"Saved: {out_pdf}")
print(f"Saved: {out_png}")
plt.close()
