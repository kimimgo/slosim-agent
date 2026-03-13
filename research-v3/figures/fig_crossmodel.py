#!/usr/bin/env python3
"""Cross-model generalization figure: grouped bar chart (5 models × 10 scenarios)"""
import matplotlib.pyplot as plt
import matplotlib
import numpy as np

matplotlib.use('Agg')

# Data from M-A3 scoring (v4 for Qwen3, crossmodel for LLaMA)
SCENARIOS = [f"S{i:02d}" for i in range(1, 11)]

scores = {
    "Qwen3:32B":  [100, 100, 88, 100, 88, 88, 67, 50, 60, 83],
    "Qwen3:14B":  [  0, 100,100, 100,  0, 88, 67, 50,  0,  0],
    "Qwen3:8B":   [100, 100, 88, 100, 50, 88, 67, 50, 60, 83],
    "LLaMA 3.3:70B": [75, 38, 38, 88, 88, 25, 67, 38, 60,  0],
    "LLaMA 3.1:8B":  [ 0,  0,  0,  0,  0,  0,  0,  0,  0,  0],
}

means = {k: np.mean(v) for k, v in scores.items()}

# Colors: Qwen3=blues, LLaMA=oranges
colors = {
    "Qwen3:32B":     "#1f77b4",
    "Qwen3:14B":     "#6baed6",
    "Qwen3:8B":      "#9ecae1",
    "LLaMA 3.3:70B": "#e6550d",
    "LLaMA 3.1:8B":  "#fdae6b",
}

model_names = list(scores.keys())
n_models = len(model_names)
n_scenarios = len(SCENARIOS)
x = np.arange(n_scenarios)
width = 0.15

fig, ax = plt.subplots(figsize=(14, 5))

for i, model in enumerate(model_names):
    offset = (i - n_models / 2 + 0.5) * width
    bars = ax.bar(x + offset, scores[model], width,
                  label=f"{model} ({means[model]:.1f}%)",
                  color=colors[model], edgecolor='white', linewidth=0.5)

ax.set_xlabel("Scenario", fontsize=12)
ax.set_ylabel("M-A3 Score (%)", fontsize=12)
ax.set_title("Cross-Model Generalization: M-A3 Parameter Fidelity", fontsize=14, fontweight='bold')
ax.set_xticks(x)
ax.set_xticklabels(SCENARIOS)
ax.set_ylim(0, 110)
ax.axhline(y=50, color='gray', linestyle='--', alpha=0.3, linewidth=0.8)

# Tier separators
ax.axvline(x=2.5, color='gray', linestyle=':', alpha=0.4)
ax.axvline(x=6.5, color='gray', linestyle=':', alpha=0.4)
ax.text(1.0, 105, "Easy", ha='center', fontsize=9, color='gray')
ax.text(4.5, 105, "Medium", ha='center', fontsize=9, color='gray')
ax.text(8.0, 105, "Hard", ha='center', fontsize=9, color='gray')

ax.legend(loc='upper right', fontsize=9, ncol=2, framealpha=0.9)
ax.grid(axis='y', alpha=0.2)

plt.tight_layout()
out_dir = "research-v3/figures/paper"
import os; os.makedirs(out_dir, exist_ok=True)
plt.savefig(f"{out_dir}/crossmodel_comparison.pdf", dpi=300, bbox_inches='tight')
plt.savefig(f"{out_dir}/crossmodel_comparison.png", dpi=150, bbox_inches='tight')
print(f"Saved: {out_dir}/crossmodel_comparison.{{pdf,png}}")

# --- Figure 2: Mean score vs model size scatter ---
fig2, ax2 = plt.subplots(figsize=(7, 5))

sizes = [32, 14, 8, 70, 8]
mean_scores = [means[m] for m in model_names]
families = ["Qwen3", "Qwen3", "Qwen3", "LLaMA", "LLaMA"]
family_colors = {"Qwen3": "#1f77b4", "LLaMA": "#e6550d"}

for m, s, sc, fam in zip(model_names, sizes, mean_scores, families):
    ax2.scatter(s, sc, s=200, c=family_colors[fam], edgecolors='black', linewidth=1.5, zorder=5)
    offset_x = 2 if s < 50 else -5
    offset_y = 3
    if m == "LLaMA 3.1:8B":
        offset_y = -8
    ax2.annotate(m, (s, sc), fontsize=9, ha='left',
                 xytext=(offset_x, offset_y), textcoords='offset points')

ax2.set_xlabel("Model Parameters (B)", fontsize=12)
ax2.set_ylabel("Mean M-A3 Score (%)", fontsize=12)
ax2.set_title("Prompt Optimization vs. Model Size", fontsize=14, fontweight='bold')
ax2.set_xlim(0, 80)
ax2.set_ylim(-5, 100)
ax2.axhline(y=50, color='gray', linestyle='--', alpha=0.3)
ax2.grid(alpha=0.2)

# Add annotation arrow: Qwen3:8B > LLaMA:70B
ax2.annotate("",
             xy=(8, 78.5), xytext=(70, 51.4),
             arrowprops=dict(arrowstyle="->", color="red", lw=2, linestyle="--"))
ax2.text(35, 70, "8B > 70B\n(prompt-optimized)", fontsize=10, color="red",
         ha='center', style='italic')

plt.tight_layout()
plt.savefig(f"{out_dir}/crossmodel_size_vs_score.pdf", dpi=300, bbox_inches='tight')
plt.savefig(f"{out_dir}/crossmodel_size_vs_score.png", dpi=150, bbox_inches='tight')
print(f"Saved: {out_dir}/crossmodel_size_vs_score.{{pdf,png}}")
