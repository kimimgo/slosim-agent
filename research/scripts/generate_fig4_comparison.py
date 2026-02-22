#!/usr/bin/env python3
"""Generate Fig 4: 8B vs 32B ablation comparison bar chart."""

import json
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import numpy as np
from pathlib import Path

EXP4_DIR = Path("research/experiments/exp4_ablation")

def load_results(path):
    with open(path) as f:
        return json.load(f)

def compute_metrics(results):
    ablations = ["full", "no-domain", "no-rules", "generic"]
    metrics = {}
    for ab in ablations:
        abr = [r for r in results if r["ablation"] == ab]
        n = len(abr) or 1
        metrics[ab] = {
            "tool_pct": sum(1 for r in abr if r["called_tool"]) / n * 100,
            "acc_pct": sum(r["param_accuracy"] for r in abr) / n * 100,
            "valid_pct": sum(1 for r in abr if r["physical_valid"]) / n * 100,
        }
    return metrics

def main():
    r8b = load_results(EXP4_DIR / "results_8b.json")
    r32b = load_results(EXP4_DIR / "results_32b.json")

    m8b = compute_metrics(r8b)
    m32b = compute_metrics(r32b)

    ablations = ["full", "no-domain", "no-rules", "generic"]
    labels = ["FULL", "NO-DOMAIN", "NO-RULES", "GENERIC"]

    fig, axes = plt.subplots(1, 3, figsize=(16, 5))

    metrics_names = [
        ("acc_pct", "Parameter Accuracy (%)"),
        ("tool_pct", "Tool Call Rate (%)"),
        ("valid_pct", "Physical Validity (%)"),
    ]

    x = np.arange(len(labels))
    width = 0.35

    for ax, (key, title) in zip(axes, metrics_names):
        vals_8b = [m8b[ab][key] for ab in ablations]
        vals_32b = [m32b[ab][key] for ab in ablations]

        bars1 = ax.bar(x - width / 2, vals_8b, width, label="Qwen3-8B",
                       color="steelblue", edgecolor="navy", alpha=0.8)
        bars2 = ax.bar(x + width / 2, vals_32b, width, label="Qwen3-32B",
                       color="coral", edgecolor="darkred", alpha=0.8)

        ax.set_ylabel(title.split("(")[0].strip())
        ax.set_title(title)
        ax.set_xticks(x)
        ax.set_xticklabels(labels, rotation=15, ha="right")
        ax.legend()
        ax.set_ylim(0, 110)
        ax.grid(True, alpha=0.3, axis="y")

        for bars in [bars1, bars2]:
            for bar in bars:
                h = bar.get_height()
                ax.annotate(f"{h:.0f}", xy=(bar.get_x() + bar.get_width() / 2, h),
                            xytext=(0, 3), textcoords="offset points",
                            ha="center", va="bottom", fontsize=8)

    fig.suptitle("Fig 4: Domain Prompt Ablation — Qwen3 8B vs 32B", fontsize=14, y=1.02)
    plt.tight_layout()

    out = EXP4_DIR / "figures" / "fig4_comparison.png"
    plt.savefig(out, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {out}")


if __name__ == "__main__":
    main()
