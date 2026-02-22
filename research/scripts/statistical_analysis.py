#!/usr/bin/env python3
"""
Statistical Analysis for SloshAgent Paper
EXP-2 (NL→XML Generation) + EXP-4 (Ablation Study)
Bootstrap BCa CI + Non-parametric Tests

Authors: SloshAgent Research Team
Seed: 42 (fixed for reproducibility)
"""

import json
import os
import numpy as np
from scipy import stats

np.random.seed(42)

# ─────────────────────────────────────────
# Paths
# ─────────────────────────────────────────
BASE = "/home/imgyu/workspace/02_active/slosim-agent/research/experiments"
OUT_DIR = BASE

# ─────────────────────────────────────────
# Load Data
# ─────────────────────────────────────────
with open(f"{BASE}/exp2_nl2xml/results_32b.json") as f:
    exp2_32b = json.load(f)
with open(f"{BASE}/exp2_nl2xml/results_8b.json") as f:
    exp2_8b = json.load(f)
with open(f"{BASE}/exp4_ablation/results_32b.json") as f:
    exp4_32b = json.load(f)
with open(f"{BASE}/exp4_ablation/results_8b.json") as f:
    exp4_8b = json.load(f)

print(f"EXP-2: 32B={len(exp2_32b)}, 8B={len(exp2_8b)} scenarios")
print(f"EXP-4: 32B={len(exp4_32b)}, 8B={len(exp4_8b)} entries (10 scenarios × 4 conditions)")

# ─────────────────────────────────────────
# Statistical Helper Functions
# ─────────────────────────────────────────

def bca_ci(data, n_boot=10000, ci=0.95, stat_fn=np.mean):
    """Bias-Corrected and Accelerated (BCa) Bootstrap Confidence Interval."""
    data = np.array(data, dtype=float)
    n = len(data)
    if n < 2:
        m = float(np.mean(data)) if n == 1 else float('nan')
        return m, m, m

    theta_hat = stat_fn(data)

    # Bootstrap distribution
    boot_rng = np.random.default_rng(42)
    boot_stats = np.array([
        stat_fn(boot_rng.choice(data, size=n, replace=True))
        for _ in range(n_boot)
    ])

    # Bias correction z0
    prop_below = np.mean(boot_stats < theta_hat)
    prop_below = np.clip(prop_below, 1e-10, 1 - 1e-10)
    z0 = stats.norm.ppf(prop_below)

    # Acceleration parameter (jackknife)
    jack_stats = np.array([stat_fn(np.delete(data, i)) for i in range(n)])
    jack_mean = np.mean(jack_stats)
    diff = jack_mean - jack_stats
    numer = np.sum(diff ** 3)
    denom = 6 * (np.sum(diff ** 2) ** 1.5)
    acc = numer / denom if abs(denom) > 1e-12 else 0.0

    # Adjusted percentiles
    alpha = (1 - ci) / 2
    z_alpha = stats.norm.ppf(alpha)
    z_1alpha = stats.norm.ppf(1 - alpha)

    def adj_perc(z_a):
        num = z0 + z_a
        denom_adj = 1 - acc * num
        if abs(denom_adj) < 1e-12:
            return np.nan
        return stats.norm.cdf(z0 + num / denom_adj)

    a1 = adj_perc(z_alpha)
    a2 = adj_perc(z_1alpha)

    a1 = np.clip(a1, 0, 1)
    a2 = np.clip(a2, 0, 1)

    lower = float(np.percentile(boot_stats, 100 * a1))
    upper = float(np.percentile(boot_stats, 100 * a2))

    return float(theta_hat), lower, upper


def permutation_test(a, b, n_perm=10000):
    """Two-sided permutation test for difference in means."""
    a, b = np.array(a, dtype=float), np.array(b, dtype=float)
    obs_diff = np.mean(a) - np.mean(b)
    combined = np.concatenate([a, b])
    na = len(a)

    rng = np.random.default_rng(42)
    perm_diffs = np.array([
        np.mean(perm := rng.permutation(combined)[:na]) - np.mean(perm := rng.permutation(combined)[:na])
        for _ in range(n_perm)
    ])
    # Recalculate properly
    perm_diffs = []
    rng = np.random.default_rng(42)
    for _ in range(n_perm):
        perm = rng.permutation(combined)
        perm_diffs.append(np.mean(perm[:na]) - np.mean(perm[na:]))
    perm_diffs = np.array(perm_diffs)

    p_value = float(np.mean(np.abs(perm_diffs) >= np.abs(obs_diff)))
    return float(obs_diff), p_value


def cohens_d(a, b):
    """Cohen's d effect size (pooled SD)."""
    a, b = np.array(a, dtype=float), np.array(b, dtype=float)
    if len(a) < 2 or len(b) < 2:
        return float('nan')
    pooled_std = np.sqrt((np.var(a, ddof=1) + np.var(b, ddof=1)) / 2)
    return float((np.mean(a) - np.mean(b)) / pooled_std) if pooled_std > 1e-12 else 0.0


def bonferroni(p_values, n_comparisons):
    """Bonferroni correction."""
    return [min(float(p) * n_comparisons, 1.0) for p in p_values]


# ─────────────────────────────────────────
# A. EXP-2: NL→XML Generation
# ─────────────────────────────────────────
print("\n" + "=" * 65)
print("A. EXP-2: NL→XML Generation Statistics")
print("=" * 65)

def organize_by_level(data):
    by_level = {}
    for entry in data:
        lvl = entry["level"]
        by_level.setdefault(lvl, []).append(entry["param_accuracy"])
    return by_level


exp2_32b_by_level = organize_by_level(exp2_32b)
exp2_8b_by_level  = organize_by_level(exp2_8b)

results = {}
results["exp2"] = {
    "level_ci_32b": {},
    "level_ci_8b": {},
    "overall_ci": {}
}

# A1. Level-wise Bootstrap 95% CI (BCa)
print("\nA1. Level-wise Accuracy Bootstrap 95% CI (BCa, 10000 resamples)")
print(f"{'Level':<7} {'32B mean':>10} {'32B 95% CI':>22} {'8B mean':>10} {'8B 95% CI':>22}")
print("-" * 75)

for lvl in sorted(exp2_32b_by_level.keys()):
    acc_32b = exp2_32b_by_level[lvl]
    acc_8b  = exp2_8b_by_level.get(lvl, [])

    m32, lo32, hi32 = bca_ci(acc_32b)
    m8,  lo8,  hi8  = bca_ci(acc_8b)

    results["exp2"]["level_ci_32b"][lvl] = {"mean": m32, "ci_lo": lo32, "ci_hi": hi32, "n": len(acc_32b)}
    results["exp2"]["level_ci_8b"][lvl]  = {"mean": m8,  "ci_lo": lo8,  "ci_hi": hi8,  "n": len(acc_8b)}

    print(f"L{lvl:<6} {m32:>10.4f} [{lo32:.4f}, {hi32:.4f}]  {m8:>10.4f} [{lo8:.4f}, {hi8:.4f}]")

# A2. Overall Bootstrap CI
all_acc_32b = [e["param_accuracy"] for e in exp2_32b]
all_acc_8b  = [e["param_accuracy"] for e in exp2_8b]

m32_all, lo32_all, hi32_all = bca_ci(all_acc_32b)
m8_all,  lo8_all,  hi8_all  = bca_ci(all_acc_8b)

results["exp2"]["overall_ci"]["32b"] = {"mean": m32_all, "ci_lo": lo32_all, "ci_hi": hi32_all, "n": 20}
results["exp2"]["overall_ci"]["8b"]  = {"mean": m8_all,  "ci_lo": lo8_all,  "ci_hi": hi8_all,  "n": 20}

print(f"\nA2. Overall Accuracy Bootstrap 95% CI (BCa, n=20 each)")
print(f"  32B: {m32_all:.4f}  95% CI [{lo32_all:.4f}, {hi32_all:.4f}]")
print(f"  8B:  {m8_all:.4f}  95% CI [{lo8_all:.4f}, {hi8_all:.4f}]")

# A3. 32B vs 8B Paired Wilcoxon Signed-Rank Test
paired_32b = {e["scenario_id"]: e["param_accuracy"] for e in exp2_32b}
paired_8b  = {e["scenario_id"]: e["param_accuracy"] for e in exp2_8b}
common_ids = sorted(set(paired_32b) & set(paired_8b))

x32 = np.array([paired_32b[sid] for sid in common_ids])
x8  = np.array([paired_8b[sid]  for sid in common_ids])

# Wilcoxon requires non-zero differences; handle ties
diff_arr = x32 - x8
if np.all(diff_arr == 0):
    w_stat, w_p = 0.0, 1.0
else:
    w_stat, w_p = stats.wilcoxon(x32, x8, zero_method='wilcox', alternative='two-sided')

results["exp2"]["wilcoxon_32b_vs_8b"] = {
    "stat": float(w_stat),
    "p_value": float(w_p),
    "n_pairs": len(common_ids),
    "mean_diff_32b_minus_8b": float(np.mean(x32 - x8))
}

print(f"\nA3. 32B vs 8B Paired Wilcoxon Signed-Rank (n={len(common_ids)})")
print(f"  W = {w_stat:.4f},  p = {w_p:.4f}")
print(f"  Mean diff (32B - 8B) = {np.mean(x32 - x8):.4f}")

# A4. Tool Call Rate by Level — Fisher's Exact (Bonferroni n=5)
print(f"\nA4. Tool Call Rate by Level: Fisher's Exact (32B vs 8B, Bonferroni n=5)")
print(f"{'Level':<7} {'32B call/N':>12} {'8B call/N':>12} {'p_raw':>8} {'p_bonf':>8}")
print("-" * 52)

results["exp2"]["fisher_tool_call_by_level"] = {}
raw_ps_fisher = []
fisher_rows   = []

for lvl in sorted(exp2_32b_by_level.keys()):
    e32 = [e for e in exp2_32b if e["level"] == lvl]
    e8  = [e for e in exp2_8b  if e["level"] == lvl]
    c32, nc32 = sum(e["called_tool"] for e in e32), len(e32) - sum(e["called_tool"] for e in e32)
    c8,  nc8  = sum(e["called_tool"] for e in e8),  len(e8)  - sum(e["called_tool"] for e in e8)

    _, p_f = stats.fisher_exact([[c32, nc32], [c8, nc8]])
    raw_ps_fisher.append(float(p_f))
    fisher_rows.append((lvl, c32, len(e32), c8, len(e8), float(p_f)))

bonf_f = bonferroni(raw_ps_fisher, 5)
for i, (lvl, c32, n32, c8, n8, p_raw) in enumerate(fisher_rows):
    p_adj = bonf_f[i]
    results["exp2"]["fisher_tool_call_by_level"][lvl] = {
        "32b_called": c32, "32b_n": n32,
        "8b_called":  c8,  "8b_n":  n8,
        "p_raw": p_raw,    "p_bonferroni": p_adj
    }
    print(f"L{lvl:<6} {c32}/{n32:>2}          {c8}/{n8:>2}       {p_raw:>8.4f} {p_adj:>8.4f}")

# ─────────────────────────────────────────
# B. EXP-4: Ablation Study
# ─────────────────────────────────────────
print("\n" + "=" * 65)
print("B. EXP-4: Ablation Study Statistics")
print("=" * 65)

CONDITIONS = ["full", "no-domain", "no-rules", "generic"]
COND_LABEL = {"full": "FULL", "no-domain": "NO-DOMAIN", "no-rules": "NO-RULES", "generic": "GENERIC"}

def organize_by_ablation(data):
    by_abl = {}
    for entry in data:
        abl = entry["ablation"]
        by_abl.setdefault(abl, []).append(entry["param_accuracy"])
    return by_abl

exp4_32b_by_abl = organize_by_ablation(exp4_32b)
exp4_8b_by_abl  = organize_by_ablation(exp4_8b)

results["exp4"] = {
    "condition_ci_32b": {},
    "condition_ci_8b": {}
}

# B1. Condition-wise Bootstrap CI
print(f"\nB1. Condition-wise Accuracy Bootstrap 95% CI (BCa, n=10 each)")
print(f"{'Condition':<12} {'32B mean':>10} {'32B 95% CI':>22} {'8B mean':>10} {'8B 95% CI':>22}")
print("-" * 80)

for cond in CONDITIONS:
    acc32 = exp4_32b_by_abl.get(cond, [])
    acc8  = exp4_8b_by_abl.get(cond, [])
    m32, lo32, hi32 = bca_ci(acc32)
    m8,  lo8,  hi8  = bca_ci(acc8)
    results["exp4"]["condition_ci_32b"][cond] = {"mean": m32, "ci_lo": lo32, "ci_hi": hi32, "n": len(acc32)}
    results["exp4"]["condition_ci_8b"][cond]  = {"mean": m8,  "ci_lo": lo8,  "ci_hi": hi8,  "n": len(acc8)}
    label = COND_LABEL[cond]
    print(f"{label:<12} {m32:>10.4f} [{lo32:.4f}, {hi32:.4f}]  {m8:>10.4f} [{lo8:.4f}, {hi8:.4f}]")

# B2-4. Permutation Tests: FULL vs {GENERIC, NO-RULES, NO-DOMAIN} (Bonferroni n=3)
COMPARISONS = [
    ("full", "generic",   "FULL vs GENERIC"),
    ("full", "no-rules",  "FULL vs NO-RULES"),
    ("full", "no-domain", "FULL vs NO-DOMAIN"),
]

print(f"\nB2-4. Permutation Tests (10000 iterations, Bonferroni n=3)")
print(f"{'Comparison':<22} | {'32B diff':>9} {'p_raw':>7} {'p_bonf':>8} | {'8B diff':>9} {'p_raw':>7} {'p_bonf':>8}")
print("-" * 80)

results["exp4"]["permutation_tests_32b"] = {}
results["exp4"]["permutation_tests_8b"]  = {}
raw_p32_perm, raw_p8_perm = [], []
perm_rows = []

for cond_a, cond_b, label in COMPARISONS:
    a32, b32 = exp4_32b_by_abl.get(cond_a, []), exp4_32b_by_abl.get(cond_b, [])
    a8,  b8  = exp4_8b_by_abl.get(cond_a, []),  exp4_8b_by_abl.get(cond_b, [])
    diff32, p32 = permutation_test(a32, b32)
    diff8,  p8  = permutation_test(a8,  b8)
    raw_p32_perm.append(p32)
    raw_p8_perm.append(p8)
    perm_rows.append((label, cond_b, diff32, p32, diff8, p8))

bonf_p32_perm = bonferroni(raw_p32_perm, 3)
bonf_p8_perm  = bonferroni(raw_p8_perm,  3)

for i, (label, cond_b, diff32, p32, diff8, p8) in enumerate(perm_rows):
    p32a = bonf_p32_perm[i]
    p8a  = bonf_p8_perm[i]
    results["exp4"]["permutation_tests_32b"][f"full_vs_{cond_b}"] = {
        "obs_diff": diff32, "p_raw": p32, "p_bonferroni": p32a
    }
    results["exp4"]["permutation_tests_8b"][f"full_vs_{cond_b}"] = {
        "obs_diff": diff8, "p_raw": p8, "p_bonferroni": p8a
    }
    print(f"{label:<22} | {diff32:>9.4f} {p32:>7.4f} {p32a:>8.4f} | {diff8:>9.4f} {p8:>7.4f} {p8a:>8.4f}")

# B5. Fisher's Exact: Tool Call Rate FULL vs GENERIC
print(f"\nB5. Fisher's Exact: Tool Call Rate FULL vs GENERIC")
results["exp4"]["fisher_tool_call"] = {}

for model_name, exp4_data in [("32B", exp4_32b), ("8B", exp4_8b)]:
    def get_rate(data, abl):
        ents = [e for e in data if e["ablation"] == abl]
        called = sum(e["called_tool"] for e in ents)
        return called, len(ents) - called, len(ents)

    c_full, nc_full, n_full = get_rate(exp4_data, "full")
    c_gen,  nc_gen,  n_gen  = get_rate(exp4_data, "generic")

    _, p_f = stats.fisher_exact([[c_full, nc_full], [c_gen, nc_gen]])

    results["exp4"]["fisher_tool_call"][model_name.lower()] = {
        "full_called": c_full, "full_n": n_full,
        "generic_called": c_gen, "generic_n": n_gen,
        "p_value": float(p_f)
    }
    print(f"  {model_name}: FULL={c_full}/{n_full}, GENERIC={c_gen}/{n_gen},  p={p_f:.4f}")

# B6. Cohen's d: FULL vs GENERIC (both models)
print(f"\nB6. Cohen's d Effect Size (FULL vs GENERIC)")
acc_full_32b = exp4_32b_by_abl.get("full", [])
acc_gen_32b  = exp4_32b_by_abl.get("generic", [])
acc_full_8b  = exp4_8b_by_abl.get("full", [])
acc_gen_8b   = exp4_8b_by_abl.get("generic", [])

d_32b = cohens_d(acc_full_32b, acc_gen_32b)
d_8b  = cohens_d(acc_full_8b,  acc_gen_8b)

results["exp4"]["cohens_d"] = {
    "32b_full_vs_generic": d_32b,
    "8b_full_vs_generic":  d_8b
}
print(f"  32B: d = {d_32b:.4f}  ({abs(d_32b):.2f} — {'large' if abs(d_32b)>=0.8 else 'medium' if abs(d_32b)>=0.5 else 'small'})")
print(f"  8B:  d = {d_8b:.4f}  ({abs(d_8b):.2f} — {'large' if abs(d_8b)>=0.8 else 'medium' if abs(d_8b)>=0.5 else 'small'})")

# B7. 8B FULL vs NO-RULES permutation (inversion significance check)
print(f"\nB7. 8B FULL vs NO-RULES (checking inversion pattern)")
acc_norules_8b = exp4_8b_by_abl.get("no-rules", [])
diff_inv, p_inv = permutation_test(acc_full_8b, acc_norules_8b)
results["exp4"]["8b_inversion_test"] = {
    "mean_full":    float(np.mean(acc_full_8b)),
    "mean_no_rules": float(np.mean(acc_norules_8b)),
    "obs_diff":     diff_inv,
    "p_raw":        p_inv
}
direction = "FULL > NO-RULES" if diff_inv > 0 else "NO-RULES > FULL (inversion!)"
print(f"  8B FULL mean={np.mean(acc_full_8b):.4f}, NO-RULES mean={np.mean(acc_norules_8b):.4f}")
print(f"  diff = {diff_inv:.4f}, p = {p_inv:.4f}  → {direction}")

# ─────────────────────────────────────────
# C. Cross-Experiment Analysis
# ─────────────────────────────────────────
print("\n" + "=" * 65)
print("C. Cross-Experiment Analysis")
print("=" * 65)

results["cross"] = {}

# C1. L1+L2 vs L3+L5 accuracy (Mann-Whitney U)
print(f"\nC1. Mann-Whitney U: L1+L2 vs L3+L5 accuracy")
print(f"{'Model':<7} {'L1+L2 mean':>12} {'L3+L5 mean':>12} {'U':>10} {'p':>8}")
print("-" * 55)

for model_name, exp2_data in [("32B", exp2_32b), ("8B", exp2_8b)]:
    low_acc  = [e["param_accuracy"] for e in exp2_data if e["level"] in [1, 2]]
    high_acc = [e["param_accuracy"] for e in exp2_data if e["level"] in [3, 5]]
    u_stat, u_p = stats.mannwhitneyu(low_acc, high_acc, alternative='two-sided')
    results["cross"][f"mannwhitney_L12vsL35_{model_name.lower()}"] = {
        "mean_L12": float(np.mean(low_acc)),
        "mean_L35": float(np.mean(high_acc)),
        "U": float(u_stat),
        "p_value": float(u_p),
        "n_L12": len(low_acc),
        "n_L35": len(high_acc)
    }
    print(f"{model_name:<7} {np.mean(low_acc):>12.4f} {np.mean(high_acc):>12.4f} {u_stat:>10.1f} {u_p:>8.4f}")

# C2. Latency: 32B vs 8B Wilcoxon Signed-Rank
print(f"\nC2. Latency Comparison: 32B vs 8B (Wilcoxon Signed-Rank, n=20 pairs)")
lat32 = {e["scenario_id"]: e["latency_s"] for e in exp2_32b}
lat8  = {e["scenario_id"]: e["latency_s"] for e in exp2_8b}
common_lat = sorted(set(lat32) & set(lat8))
arr32 = np.array([lat32[sid] for sid in common_lat])
arr8  = np.array([lat8[sid]  for sid in common_lat])

w_lat, p_lat = stats.wilcoxon(arr32, arr8, alternative='two-sided')
speedup = float(np.median(arr32) / np.median(arr8))

results["cross"]["latency"] = {
    "median_32b_s":   float(np.median(arr32)),
    "median_8b_s":    float(np.median(arr8)),
    "mean_32b_s":     float(np.mean(arr32)),
    "mean_8b_s":      float(np.mean(arr8)),
    "speedup_factor": speedup,
    "wilcoxon_W":     float(w_lat),
    "p_value":        float(p_lat),
    "n_pairs":        len(common_lat)
}

print(f"  32B median latency: {np.median(arr32):.1f}s  (mean: {np.mean(arr32):.1f}s)")
print(f"  8B  median latency: {np.median(arr8):.1f}s  (mean: {np.mean(arr8):.1f}s)")
print(f"  Speedup 8B/32B:     {speedup:.1f}×")
print(f"  W = {w_lat:.4f},  p = {p_lat:.4f}")

# ─────────────────────────────────────────
# Descriptive Summary (for markdown)
# ─────────────────────────────────────────
def sig_stars(p):
    if p < 0.001: return "***"
    if p < 0.01:  return "**"
    if p < 0.05:  return "*"
    return "ns"

# ─────────────────────────────────────────
# Save JSON
# ─────────────────────────────────────────
def to_python(obj):
    if isinstance(obj, (np.floating, np.float64, np.float32)):
        return float(obj)
    if isinstance(obj, (np.integer, np.int64, np.int32)):
        return int(obj)
    if isinstance(obj, dict):
        return {k: to_python(v) for k, v in obj.items()}
    if isinstance(obj, (list, tuple)):
        return [to_python(v) for v in obj]
    return obj

json_path = os.path.join(OUT_DIR, "statistical_analysis.json")
with open(json_path, "w", encoding="utf-8") as f:
    json.dump(to_python(results), f, indent=2, ensure_ascii=False)
print(f"\n✓ JSON saved: {json_path}")

# ─────────────────────────────────────────
# Generate Markdown Report
# ─────────────────────────────────────────
md_lines = []
md_lines.append("# Statistical Analysis Report — SloshAgent Paper")
md_lines.append("")
md_lines.append("> **Seed**: 42 | **Bootstrap**: BCa, 10,000 resamples | "
                "**Correction**: Bonferroni | **Non-parametric**: Wilcoxon/Mann-Whitney/Permutation")
md_lines.append("")

# ── EXP-2 section ──
md_lines.append("## A. EXP-2: NL→XML Generation (20 Scenarios, 5 Levels)")
md_lines.append("")
md_lines.append("### A1. Level-wise Accuracy — BCa 95% CI")
md_lines.append("")
md_lines.append("| Level | 32B Mean | 32B 95% CI | 8B Mean | 8B 95% CI |")
md_lines.append("|-------|----------|------------|---------|-----------|")
for lvl in sorted(results["exp2"]["level_ci_32b"].keys()):
    r32 = results["exp2"]["level_ci_32b"][lvl]
    r8  = results["exp2"]["level_ci_8b"][lvl]
    md_lines.append(
        f"| L{lvl} (n={r32['n']}) | {r32['mean']:.4f} | [{r32['ci_lo']:.4f}, {r32['ci_hi']:.4f}] "
        f"| {r8['mean']:.4f} | [{r8['ci_lo']:.4f}, {r8['ci_hi']:.4f}] |"
    )
md_lines.append("")

md_lines.append("### A2. Overall Accuracy — BCa 95% CI")
md_lines.append("")
r32o = results["exp2"]["overall_ci"]["32b"]
r8o  = results["exp2"]["overall_ci"]["8b"]
md_lines.append(f"- **32B**: {r32o['mean']:.4f}  95% CI [{r32o['ci_lo']:.4f}, {r32o['ci_hi']:.4f}]  (n=20)")
md_lines.append(f"- **8B**: {r8o['mean']:.4f}  95% CI [{r8o['ci_lo']:.4f}, {r8o['ci_hi']:.4f}]  (n=20)")
md_lines.append("")

md_lines.append("### A3. 32B vs 8B — Paired Wilcoxon Signed-Rank Test")
md_lines.append("")
wr = results["exp2"]["wilcoxon_32b_vs_8b"]
md_lines.append(f"- W = {wr['stat']:.4f}, p = {wr['p_value']:.4f} {sig_stars(wr['p_value'])}")
md_lines.append(f"- Mean difference (32B − 8B) = {wr['mean_diff_32b_minus_8b']:.4f}")
md_lines.append(f"- n pairs = {wr['n_pairs']}")
md_lines.append("")

md_lines.append("### A4. Tool Call Rate by Level — Fisher's Exact (Bonferroni n=5)")
md_lines.append("")
md_lines.append("| Level | 32B called/N | 8B called/N | p_raw | p_bonf | sig |")
md_lines.append("|-------|-------------|------------|-------|--------|-----|")
for lvl, r in sorted(results["exp2"]["fisher_tool_call_by_level"].items()):
    md_lines.append(
        f"| L{lvl} | {r['32b_called']}/{r['32b_n']} | {r['8b_called']}/{r['8b_n']} "
        f"| {r['p_raw']:.4f} | {r['p_bonferroni']:.4f} | {sig_stars(r['p_bonferroni'])} |"
    )
md_lines.append("")

# ── EXP-4 section ──
md_lines.append("## B. EXP-4: Ablation Study (10 Scenarios × 4 Conditions)")
md_lines.append("")
md_lines.append("### B1. Condition-wise Accuracy — BCa 95% CI")
md_lines.append("")
md_lines.append("| Condition | 32B Mean | 32B 95% CI | 8B Mean | 8B 95% CI |")
md_lines.append("|-----------|----------|------------|---------|-----------|")
for cond in CONDITIONS:
    r32 = results["exp4"]["condition_ci_32b"][cond]
    r8  = results["exp4"]["condition_ci_8b"][cond]
    label = COND_LABEL[cond]
    md_lines.append(
        f"| {label} (n={r32['n']}) | {r32['mean']:.4f} | [{r32['ci_lo']:.4f}, {r32['ci_hi']:.4f}] "
        f"| {r8['mean']:.4f} | [{r8['ci_lo']:.4f}, {r8['ci_hi']:.4f}] |"
    )
md_lines.append("")

md_lines.append("### B2–4. Permutation Tests: FULL vs Degraded (Bonferroni n=3)")
md_lines.append("")
md_lines.append("| Comparison | 32B diff | 32B p_raw | 32B p_bonf | 32B sig | 8B diff | 8B p_raw | 8B p_bonf | 8B sig |")
md_lines.append("|------------|----------|-----------|------------|---------|---------|----------|-----------|--------|")
for cond_b_key in ["generic", "no-rules", "no-domain"]:
    label_map = {"generic": "FULL vs GENERIC", "no-rules": "FULL vs NO-RULES", "no-domain": "FULL vs NO-DOMAIN"}
    r32 = results["exp4"]["permutation_tests_32b"][f"full_vs_{cond_b_key}"]
    r8  = results["exp4"]["permutation_tests_8b"][f"full_vs_{cond_b_key}"]
    md_lines.append(
        f"| {label_map[cond_b_key]} | {r32['obs_diff']:.4f} | {r32['p_raw']:.4f} | {r32['p_bonferroni']:.4f} | {sig_stars(r32['p_bonferroni'])} "
        f"| {r8['obs_diff']:.4f} | {r8['p_raw']:.4f} | {r8['p_bonferroni']:.4f} | {sig_stars(r8['p_bonferroni'])} |"
    )
md_lines.append("")

md_lines.append("### B5. Tool Call Rate: FULL vs GENERIC — Fisher's Exact")
md_lines.append("")
md_lines.append("| Model | FULL called/N | GENERIC called/N | p-value | sig |")
md_lines.append("|-------|--------------|-----------------|---------|-----|")
for mdl in ["32b", "8b"]:
    r = results["exp4"]["fisher_tool_call"][mdl]
    md_lines.append(
        f"| {mdl.upper()} | {r['full_called']}/{r['full_n']} | {r['generic_called']}/{r['generic_n']} "
        f"| {r['p_value']:.4f} | {sig_stars(r['p_value'])} |"
    )
md_lines.append("")

md_lines.append("### B6. Cohen's d Effect Size: FULL vs GENERIC")
md_lines.append("")
cd = results["exp4"]["cohens_d"]
def d_interp(d_val):
    a = abs(d_val)
    if a >= 0.8: return "large"
    if a >= 0.5: return "medium"
    if a >= 0.2: return "small"
    return "negligible"
md_lines.append(f"- **32B**: d = {cd['32b_full_vs_generic']:.4f} ({d_interp(cd['32b_full_vs_generic'])})")
md_lines.append(f"- **8B**: d = {cd['8b_full_vs_generic']:.4f} ({d_interp(cd['8b_full_vs_generic'])})")
md_lines.append("")

md_lines.append("### B7. 8B: FULL vs NO-RULES — Inversion Test")
md_lines.append("")
inv = results["exp4"]["8b_inversion_test"]
direction = "FULL > NO-RULES" if inv["obs_diff"] > 0 else "**NO-RULES > FULL (inversion)**"
md_lines.append(f"- FULL mean = {inv['mean_full']:.4f}, NO-RULES mean = {inv['mean_no_rules']:.4f}")
md_lines.append(f"- diff = {inv['obs_diff']:.4f}, p = {inv['p_raw']:.4f} {sig_stars(inv['p_raw'])}")
md_lines.append(f"- Direction: {direction}")
md_lines.append("")

# ── Cross-experiment section ──
md_lines.append("## C. Cross-Experiment Analysis")
md_lines.append("")

md_lines.append("### C1. Easy (L1+L2) vs Hard (L3+L5) — Mann-Whitney U")
md_lines.append("")
md_lines.append("| Model | L1+L2 mean (n=8) | L3+L5 mean (n=8) | U | p-value | sig |")
md_lines.append("|-------|-----------------|-----------------|---|---------|-----|")
for mdl in ["32b", "8b"]:
    r = results["cross"][f"mannwhitney_L12vsL35_{mdl}"]
    md_lines.append(
        f"| {mdl.upper()} | {r['mean_L12']:.4f} | {r['mean_L35']:.4f} "
        f"| {r['U']:.1f} | {r['p_value']:.4f} | {sig_stars(r['p_value'])} |"
    )
md_lines.append("")

md_lines.append("### C2. Latency: 32B vs 8B — Wilcoxon Signed-Rank")
md_lines.append("")
lat = results["cross"]["latency"]
md_lines.append(f"| Metric | 32B | 8B |")
md_lines.append(f"|--------|-----|-----|")
md_lines.append(f"| Median latency | {lat['median_32b_s']:.1f}s | {lat['median_8b_s']:.1f}s |")
md_lines.append(f"| Mean latency | {lat['mean_32b_s']:.1f}s | {lat['mean_8b_s']:.1f}s |")
md_lines.append(f"| Speedup (32B/8B) | {lat['speedup_factor']:.1f}× | — |")
md_lines.append(f"| W | {lat['wilcoxon_W']:.1f} | — |")
md_lines.append(f"| p-value | {lat['p_value']:.4f} {sig_stars(lat['p_value'])} | — |")
md_lines.append("")

md_lines.append("---")
md_lines.append("")
md_lines.append("*Significance levels: \\*\\*\\* p<0.001, \\*\\* p<0.01, \\* p<0.05, ns = not significant*")
md_lines.append("")
md_lines.append("> Generated by `research/scripts/statistical_analysis.py`")

md_path = os.path.join(OUT_DIR, "statistical_analysis.md")
with open(md_path, "w", encoding="utf-8") as f:
    f.write("\n".join(md_lines))
print(f"✓ Markdown saved: {md_path}")
print("\n✅ All statistical analyses complete.")
