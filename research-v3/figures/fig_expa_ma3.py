#!/usr/bin/env python3
"""EXP-A: M-A3 Parameter Fidelity bar chart (32B vs 8B, 10 scenarios)"""
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

# Data from analyze_all.py (score_expb 8-param M-A3, trial mean)
scenarios = ['S01', 'S02', 'S03', 'S04', 'S05', 'S06', 'S07', 'S08', 'S09', 'S10']
m32b = [75.0, 100.0, 87.5, 75.0, 50.0, 87.5, 66.7, 50.0, 20.0, 0.0]
m8b  = [75.0, 100.0, 87.5, 75.0, 50.0, 87.5, 66.7,  0.0, 20.0, 25.0]

# Tier coloring
tier_colors = {
    'Easy':   '#4CAF50',
    'Medium': '#FF9800',
    'Hard':   '#F44336',
}
tiers = ['Easy','Easy','Medium','Medium','Medium','Medium','Hard','Easy','Hard','Hard']
# For paper, group by scenario order: S01-S03=Easy, S04-S07=Medium, S08-S10=Hard
paper_tiers = ['Easy','Easy','Easy','Medium','Medium','Medium','Medium','Hard','Hard','Hard']

x = np.arange(len(scenarios))
width = 0.35

fig, ax = plt.subplots(figsize=(10, 4.5))

bars1 = ax.bar(x - width/2, m32b, width, label='Qwen3-32B', color='#1976D2', alpha=0.85)
bars2 = ax.bar(x + width/2, m8b, width, label='Qwen3-8B', color='#FF7043', alpha=0.85)

# Add tier background shading
ax.axvspan(-0.5, 2.5, alpha=0.06, color='green', label='_Easy')
ax.axvspan(2.5, 6.5, alpha=0.06, color='orange', label='_Medium')
ax.axvspan(6.5, 9.5, alpha=0.06, color='red', label='_Hard')

# Tier labels
ax.text(1.0, 103, 'Easy', ha='center', fontsize=9, color='#388E3C', fontweight='bold')
ax.text(4.5, 103, 'Medium', ha='center', fontsize=9, color='#E65100', fontweight='bold')
ax.text(8.0, 103, 'Hard', ha='center', fontsize=9, color='#C62828', fontweight='bold')

# Overall means
mean_32b = np.mean(m32b)
mean_8b = np.mean(m8b)
ax.axhline(y=mean_32b, color='#1976D2', linestyle='--', alpha=0.5, linewidth=1)
ax.axhline(y=mean_8b, color='#FF7043', linestyle='--', alpha=0.5, linewidth=1)
ax.text(9.7, mean_32b + 1.5, f'32B μ={mean_32b:.1f}%', fontsize=7, color='#1976D2')
ax.text(9.7, mean_8b - 4, f'8B μ={mean_8b:.1f}%', fontsize=7, color='#FF7043')

ax.set_xlabel('Scenario', fontsize=11)
ax.set_ylabel('M-A3 Parameter Fidelity (%)', fontsize=11)
ax.set_title('EXP-A: NL→XML Parameter Fidelity across 10 Sloshing Scenarios', fontsize=12, fontweight='bold')
ax.set_xticks(x)
ax.set_xticklabels(scenarios, fontsize=9)
ax.set_ylim(0, 115)
ax.legend(loc='upper right', fontsize=9)
ax.grid(axis='y', alpha=0.3)

# Value labels on bars
for bar in bars1:
    h = bar.get_height()
    if h > 0:
        ax.text(bar.get_x() + bar.get_width()/2., h + 1, f'{h:.0f}', ha='center', va='bottom', fontsize=7)

for bar in bars2:
    h = bar.get_height()
    if h > 0:
        ax.text(bar.get_x() + bar.get_width()/2., h + 1, f'{h:.0f}', ha='center', va='bottom', fontsize=7)

plt.tight_layout()
plt.savefig('fig_expa_ma3.png', dpi=300, bbox_inches='tight')
plt.savefig('fig_expa_ma3.pdf', bbox_inches='tight')
print("Saved: fig_expa_ma3.png, fig_expa_ma3.pdf")
