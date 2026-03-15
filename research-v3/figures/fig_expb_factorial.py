#!/usr/bin/env python3
"""EXP-B: 2×2 Factorial Ablation — Interaction Plot + Bar Chart"""
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

# Data from score_expb.py (10 scenarios × 2 models combined)
conditions = ['B0\nFull', 'B1\n−Prompt', 'B2\n−Tool', 'B4\nBare']
ma3 = [67.0, 0.0, 46.1, 0.0]
has_prompt = [True, False, True, False]
has_tools = [True, True, False, False]

# Colors
colors = ['#1976D2', '#9E9E9E', '#FF9800', '#E0E0E0']

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(11, 4.5))

# ── Left: Bar Chart ──
bars = ax1.bar(conditions, ma3, color=colors, edgecolor='black', linewidth=0.5, width=0.6)

for bar, val in zip(bars, ma3):
    ax1.text(bar.get_x() + bar.get_width()/2., bar.get_height() + 1.5,
             f'{val:.1f}%', ha='center', va='bottom', fontsize=10, fontweight='bold')

# Prompt/Tool indicators
for i, (p, t) in enumerate(zip(has_prompt, has_tools)):
    y_base = -12
    ax1.text(i, y_base, f'P:{"✓" if p else "✗"}  T:{"✓" if t else "✗"}',
             ha='center', fontsize=8, color='gray')

ax1.set_ylabel('M-A3 Parameter Fidelity (%)', fontsize=11)
ax1.set_title('(a) Ablation Conditions', fontsize=12, fontweight='bold')
ax1.set_ylim(-15, 90)
ax1.grid(axis='y', alpha=0.3)
ax1.axhline(y=0, color='black', linewidth=0.5)

# ── Right: Interaction Plot ──
# Prompt present vs absent, lines for Tools present/absent
prompt_labels = ['Absent', 'Present']

# Tools present: B1(0%) → B0(67.0%)
tools_present = [0.0, 67.0]
# Tools absent: B4(0%) → B2(46.1%)
tools_absent = [0.0, 46.1]

ax2.plot(prompt_labels, tools_present, 'o-', color='#1976D2', linewidth=2.5,
         markersize=10, label='Tools ✓', zorder=5)
ax2.plot(prompt_labels, tools_absent, 's--', color='#FF9800', linewidth=2.5,
         markersize=10, label='Tools ✗', zorder=5)

# Annotations
ax2.annotate('B0: 67.0%', xy=(1, 67.0), xytext=(0.7, 73),
             fontsize=9, fontweight='bold', color='#1976D2')
ax2.annotate('B2: 46.1%', xy=(1, 46.1), xytext=(0.7, 39),
             fontsize=9, fontweight='bold', color='#FF9800')
ax2.annotate('B1=B4=0%\n(40/40 runs)', xy=(0, 0), xytext=(0.05, 8),
             fontsize=9, color='gray')

# Effect annotations
ax2.annotate('', xy=(1.08, 46.1), xycoords='data',
             xytext=(1.08, 67.0), textcoords='data',
             arrowprops=dict(arrowstyle='<->', color='#388E3C', lw=1.5))
ax2.text(1.15, 56, 'Tool\n+21.0%', fontsize=8, color='#388E3C', fontweight='bold')

ax2.annotate('', xy=(0.5, 2), xycoords='data',
             xytext=(0.5, 55), textcoords='data',
             arrowprops=dict(arrowstyle='<->', color='#C62828', lw=1.5))
ax2.text(0.3, 27, 'Prompt\n+56.5%', fontsize=8, color='#C62828', fontweight='bold',
         ha='center')

ax2.set_xlabel('Domain Prompt', fontsize=11)
ax2.set_ylabel('M-A3 Parameter Fidelity (%)', fontsize=11)
ax2.set_title('(b) Prompt × Tool Interaction', fontsize=12, fontweight='bold')
ax2.set_ylim(-5, 85)
ax2.legend(loc='upper left', fontsize=9)
ax2.grid(alpha=0.3)

plt.tight_layout()
plt.savefig('fig_expb_factorial.png', dpi=300, bbox_inches='tight')
plt.savefig('fig_expb_factorial.pdf', bbox_inches='tight')
print("Saved: fig_expb_factorial.png, fig_expb_factorial.pdf")
