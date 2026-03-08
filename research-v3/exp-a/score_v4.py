#!/usr/bin/env python3
"""Score v4 (P1+P2 fix) results and compare with v3 baseline."""
import sys
sys.path.insert(0, str(__import__('pathlib').Path(__file__).parent.parent / 'exp-b'))
from score_expb import parse_xml_params, score_scenario, check_xml_validity, GROUND_TRUTH
from pathlib import Path
import glob

V4_DIR = Path(__file__).parent / "results_v4"
V3_DIR = Path(__file__).parent / "results"

SCENARIOS = ['S01', 'S04', 'S05', 'S08', 'S09']
MODELS = ['qwen3_32b', 'qwen3_latest']

def find_best_xml(result_dir):
    """Find the most relevant XML in a result directory.
    Prefers root-level XMLs over nested ones (agent may overwrite initial attempt).
    """
    candidates = []
    for xml in sorted(result_dir.rglob('*.xml')):
        rel = str(xml.relative_to(result_dir))
        if 'simulations/simulations' in rel:
            continue
        candidates.append(xml)

    if not candidates:
        return None

    # Separate root-level vs nested
    root_xmls = [c for c in candidates if c.parent == result_dir]
    nested_xmls = [c for c in candidates if c.parent != result_dir]

    # Prefer root-level XMLs first (agent's final output)
    priority = ['pitch_case.xml', 'sloshing_case_1deg.xml', 'spheric_benchmark.xml',
                'cylindrical_tank.xml', 'sloshing_case.xml', 'case.xml', 'base.xml']
    for pref in priority:
        for c in root_xmls:
            if c.name == pref:
                return c
    # Fallback to nested
    for pref in priority:
        for c in nested_xmls:
            if c.name == pref:
                return c
    return candidates[0]


print("=" * 70)
print("  EXP-A v4 (P1+P2 Fix) — Pitch Scenario Retest")
print("=" * 70)

v3_scores = {}
v4_scores = {}

for model in MODELS:
    model_label = '32B' if '32b' in model else '8B'
    print(f"\n{'─' * 60}")
    print(f"  Model: {model_label} ({model})")
    print(f"{'─' * 60}")

    for scenario in SCENARIOS:
        # v3 score
        v3_dir = V3_DIR / f"{scenario}_{model}_trial1"
        v3_xml = find_best_xml(v3_dir) if v3_dir.exists() else None
        if v3_xml:
            params_v3, _ = parse_xml_params(str(v3_xml))
            if params_v3:
                _, passed_v3, total_v3 = score_scenario(scenario, params_v3)
                pct_v3 = (passed_v3 / total_v3 * 100) if total_v3 > 0 else 0
            else:
                pct_v3 = 0
        else:
            pct_v3 = 0
        v3_scores[f"{scenario}_{model}"] = pct_v3

        # v4 score
        v4_dir = V4_DIR / f"{scenario}_{model}"
        v4_xml = find_best_xml(v4_dir) if v4_dir.exists() else None
        if v4_xml:
            params_v4, _ = parse_xml_params(str(v4_xml))
            if params_v4:
                checks, passed_v4, total_v4 = score_scenario(scenario, params_v4)
                pct_v4 = (passed_v4 / total_v4 * 100) if total_v4 > 0 else 0
            else:
                checks, pct_v4 = [], 0
        else:
            checks, pct_v4 = [], 0
        v4_scores[f"{scenario}_{model}"] = pct_v4

        delta = pct_v4 - pct_v3
        arrow = "↑" if delta > 0 else "↓" if delta < 0 else "="
        print(f"\n  {scenario}: v3={pct_v3:.0f}% → v4={pct_v4:.0f}% ({arrow}{abs(delta):.0f}pp)")
        if v4_xml:
            print(f"    XML: {v4_xml.name}")
        for c in checks:
            mark = 'PASS' if c['pass'] else 'FAIL'
            print(f"    [{mark}] {c['name']}: expected={c['expected']}, actual={c['actual']}")

# Summary
print(f"\n\n{'=' * 70}")
print("  SUMMARY — Pitch Scenarios Only")
print(f"{'=' * 70}")

header = f"  {'':>6}"
for s in SCENARIOS:
    header += f"  {s:>5}"
header += f"  {'Mean':>6}"
print(header)

for model in MODELS:
    model_label = '32B' if '32b' in model else '8B'
    for version, scores in [('v3', v3_scores), ('v4', v4_scores)]:
        row = f"  {model_label} {version}"
        vals = []
        for s in SCENARIOS:
            key = f"{s}_{model}"
            v = scores.get(key, 0)
            row += f"  {v:>4.0f}%"
            vals.append(v)
        mean = sum(vals) / len(vals)
        row += f"  {mean:>5.1f}%"
        print(row)
    # Delta
    row = f"  {model_label} Δ "
    vals = []
    for s in SCENARIOS:
        key = f"{s}_{model}"
        d = v4_scores.get(key, 0) - v3_scores.get(key, 0)
        row += f"  {d:>+4.0f}p"
        vals.append(d)
    mean_d = sum(vals) / len(vals)
    row += f"  {mean_d:>+5.1f}p"
    print(row)
    print()

# Combined improvement for all 10 scenarios
print(f"\n{'=' * 70}")
print("  PROJECTED FULL M-A3 (all 10 scenarios)")
print(f"{'=' * 70}")

# Non-pitch scenarios keep their v3 scores
NON_PITCH = ['S02', 'S03', 'S06', 'S07', 'S10']
ALL_SCENARIOS = ['S01', 'S02', 'S03', 'S04', 'S05', 'S06', 'S07', 'S08', 'S09', 'S10']

for model in MODELS:
    model_label = '32B' if '32b' in model else '8B'
    full_v3 = []
    full_v4 = []
    for s in ALL_SCENARIOS:
        key = f"{s}_{model}"
        if s in SCENARIOS:
            full_v3.append(v3_scores.get(key, 0))
            full_v4.append(v4_scores.get(key, 0))
        else:
            # Use v3 score for non-pitch scenarios (unchanged)
            v3_dir = V3_DIR / f"{s}_{model}_trial1"
            v3_xml = find_best_xml(v3_dir) if v3_dir.exists() else None
            if v3_xml:
                params, _ = parse_xml_params(str(v3_xml))
                if params:
                    _, p, t = score_scenario(s, params)
                    pct = (p / t * 100) if t > 0 else 0
                else:
                    pct = 0
            else:
                pct = 0
            full_v3.append(pct)
            full_v4.append(pct)  # same for non-pitch

    mean_v3 = sum(full_v3) / len(full_v3)
    mean_v4 = sum(full_v4) / len(full_v4)
    print(f"  {model_label}: v3={mean_v3:.1f}% → v4={mean_v4:.1f}% (Δ={mean_v4-mean_v3:+.1f}pp)")
