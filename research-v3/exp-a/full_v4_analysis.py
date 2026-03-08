#!/usr/bin/env python3
"""Full v3→v4 analysis: all 10 scenarios, both models, comprehensive comparison."""
import sys
sys.path.insert(0, str(__import__('pathlib').Path(__file__).parent.parent / 'exp-b'))
from score_expb import parse_xml_params, score_scenario, GROUND_TRUTH
from pathlib import Path
import json

V4_DIR = Path(__file__).parent / "results_v4"
V3_DIR = Path(__file__).parent / "results"

ALL = ['S01','S02','S03','S04','S05','S06','S07','S08','S09','S10']
PITCH = ['S01','S04','S05','S08','S09']
MODELS = ['qwen3_32b', 'qwen3_latest']

def find_best_xml(result_dir):
    if not result_dir.exists():
        return None
    candidates = []
    for xml in sorted(result_dir.rglob('*.xml')):
        rel = str(xml.relative_to(result_dir))
        if 'simulations/simulations' in rel:
            continue
        candidates.append(xml)
    if not candidates:
        return None
    root_xmls = [c for c in candidates if c.parent == result_dir]
    nested_xmls = [c for c in candidates if c.parent != result_dir]
    priority = ['fuel_tank.xml', 'pitch_case.xml', 'sloshing_case_1deg.xml',
                'spheric_benchmark.xml', 'cylindrical_tank.xml',
                'sloshing_case.xml', 'case.xml', 'base.xml']
    for pref in priority:
        for c in root_xmls:
            if c.name == pref:
                return c
    for pref in priority:
        for c in nested_xmls:
            if c.name == pref:
                return c
    return candidates[0]

def score_xml(xml_path, scenario):
    if xml_path is None:
        return 0, []
    params, _ = parse_xml_params(str(xml_path))
    if not params:
        return 0, []
    checks, passed, total = score_scenario(scenario, params)
    return (passed / total * 100) if total > 0 else 0, checks

# ======================================================================
# Collect all scores
# ======================================================================
results = {}  # {model: {scenario: {v3: pct, v4: pct}}}

for model in MODELS:
    results[model] = {}
    for s in ALL:
        # v3
        v3_dir = V3_DIR / f"{s}_{model}_trial1"
        v3_xml = find_best_xml(v3_dir)
        v3_pct, _ = score_xml(v3_xml, s)

        # v4: only pitch scenarios have v4 results
        if s in PITCH:
            v4_dir = V4_DIR / f"{s}_{model}"
            v4_xml = find_best_xml(v4_dir)
            v4_pct, checks = score_xml(v4_xml, s)
        else:
            v4_pct = v3_pct  # unchanged

        results[model][s] = {'v3': v3_pct, 'v4': v4_pct}

# ======================================================================
# Print comprehensive table
# ======================================================================
print("=" * 90)
print("  COMPREHENSIVE M-A3 ANALYSIS: v3 (baseline) vs v4 (P1+P2 fix)")
print("=" * 90)

tier_map = {s: GROUND_TRUTH.get(s, {}).get('tier', '?') for s in ALL}
tier_abbr = {'Easy': 'E', 'Medium': 'M', 'Hard': 'H'}

for model in MODELS:
    label = '32B' if '32b' in model else '8B'
    print(f"\n  Model: {label}")
    print(f"  {'':>5}", end='')
    for s in ALL:
        t = tier_abbr.get(tier_map[s], '?')
        print(f" {s}({t})", end='')
    print(f"  {'Mean':>6}")
    print(f"  {'─'*85}")

    for ver in ['v3', 'v4']:
        row = f"  {ver:>5}"
        vals = []
        for s in ALL:
            v = results[model][s][ver]
            row += f"  {v:>4.0f}%"
            vals.append(v)
        mean = sum(vals) / len(vals)
        row += f"  {mean:>5.1f}%"
        print(row)

    # Delta
    row = f"  {'Δ':>5}"
    vals = []
    for s in ALL:
        d = results[model][s]['v4'] - results[model][s]['v3']
        if d != 0:
            row += f"  {d:>+4.0f}p"
        else:
            row += f"  {'  · ':>5}"
        vals.append(d)
    mean_d = sum(vals) / len(vals)
    row += f"  {mean_d:>+5.1f}p"
    print(row)

# ======================================================================
# GAP Evidence Summary
# ======================================================================
print(f"\n\n{'=' * 90}")
print("  GAP EVIDENCE SUMMARY")
print(f"{'=' * 90}")

# GAP 1: Non-expert accessibility
for model in MODELS:
    label = '32B' if '32b' in model else '8B'
    v4_mean = sum(results[model][s]['v4'] for s in ALL) / len(ALL)
    print(f"\n  GAP 1 — Non-expert accessibility ({label}):")
    print(f"    Baseline (no agent): 0% M-A3")
    print(f"    With agent v3:       {sum(results[model][s]['v3'] for s in ALL)/len(ALL):.1f}% M-A3")
    print(f"    With agent v4:       {v4_mean:.1f}% M-A3")
    print(f"    Improvement:         {v4_mean:.1f}pp over baseline")

# GAP 3: Local SLM sufficiency
print(f"\n  GAP 3 — Local SLM sufficiency:")
v4_32b = sum(results['qwen3_32b'][s]['v4'] for s in ALL) / len(ALL)
v4_8b = sum(results['qwen3_latest'][s]['v4'] for s in ALL) / len(ALL)
delta = v4_32b - v4_8b
print(f"    32B v4: {v4_32b:.1f}%")
print(f"     8B v4: {v4_8b:.1f}%")
print(f"    Delta:  {delta:.1f}pp")

# Scenario-by-scenario agreement
agree = sum(1 for s in ALL if results['qwen3_32b'][s]['v4'] == results['qwen3_latest'][s]['v4'])
print(f"    Identical scores: {agree}/{len(ALL)} scenarios")

# Per-scenario comparison
print(f"\n  v4 Per-scenario 32B vs 8B:")
for s in ALL:
    v32 = results['qwen3_32b'][s]['v4']
    v8 = results['qwen3_latest'][s]['v4']
    diff = v32 - v8
    marker = "=" if diff == 0 else f"32B+{diff:.0f}" if diff > 0 else f"8B+{abs(diff):.0f}"
    print(f"    {s}: 32B={v32:.0f}% 8B={v8:.0f}% [{marker}]")

# ======================================================================
# Save results JSON
# ======================================================================
output = {
    'models': {},
    'gap_evidence': {}
}
for model in MODELS:
    label = '32B' if '32b' in model else '8B'
    output['models'][label] = {
        'v3_scores': {s: results[model][s]['v3'] for s in ALL},
        'v4_scores': {s: results[model][s]['v4'] for s in ALL},
        'v3_mean': sum(results[model][s]['v3'] for s in ALL) / len(ALL),
        'v4_mean': sum(results[model][s]['v4'] for s in ALL) / len(ALL),
    }
output['gap_evidence'] = {
    'gap1_v4_32b': v4_32b,
    'gap1_v4_8b': v4_8b,
    'gap3_delta': delta,
    'gap3_identical': agree,
}

out_path = Path(__file__).parent / "v4_analysis_results.json"
with open(out_path, 'w') as f:
    json.dump(output, f, indent=2)
print(f"\n  Saved: {out_path}")
