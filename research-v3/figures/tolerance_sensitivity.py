#!/usr/bin/env python3
"""Tolerance Sensitivity Analysis — M-A3를 여러 tolerance 수준에서 재계산

목적: "tolerance 선택이 자의적이다"라는 리뷰어 공격 방어
방법: score_expb.py의 tolerance를 5%, 10%, 15%, 20%로 변경하며 M-A3 재산정
"""
import json
import sys
from pathlib import Path
from copy import deepcopy

sys.path.insert(0, str(Path(__file__).parent.parent / "exp-b"))
from score_expb import parse_xml_params, check_xml_validity

SCRIPT_DIR = Path(__file__).parent.parent
GT_PATH = SCRIPT_DIR / "exp-a" / "ground_truth.json"
RESULTS_DIR = SCRIPT_DIR / "exp-a" / "results"

with open(GT_PATH) as f:
    GT = json.load(f)

SCENARIOS = [f"S{i:02d}" for i in range(1, 11)]
MODELS = ["qwen3_32b", "qwen3_latest"]
TOLERANCES = [0.05, 0.10, 0.15, 0.20, 0.25]  # 5% to 25%


def find_xml(result_dir: Path):
    candidates = [
        result_dir / "simulations" / "sloshing_case.xml",
        result_dir / "sloshing_case.xml",
    ]
    alt_names = ["parametric_case.xml", "ISOPE_LNG_Benchmark.xml",
                 "spheric_benchmark.xml", "case.xml"]
    for alt in alt_names:
        candidates.append(result_dir / "simulations" / alt)
        candidates.append(result_dir / alt)
    for c in candidates:
        if c.exists():
            return c
    return None


def score_with_tolerance(scenario_id, params, tol_pct):
    """score_expb.score_scenario와 동일하되, tolerance를 통일 파라미터로 받음"""
    gt = GT.get(scenario_id, {})
    if not gt:
        return [], 0, 0

    checks = []

    # Tank dimensions
    if 'x' in gt.get('tank', {}):
        tank_gt = gt['tank']
        tank_act = params.get('tank', {})
        for dim in ['x', 'y', 'z']:
            if dim in tank_gt:
                expected = tank_gt[dim]
                actual = tank_act.get(dim, 0)
                tol = expected * tol_pct
                passed = abs(actual - expected) <= tol
                checks.append({'name': f'tank_{dim}', 'pass': passed})
    elif gt.get('tank', {}).get('shape') in ('cylinder', 'horizontal_cylinder'):
        has_cyl = params.get('has_cylinder', False)
        checks.append({'name': 'geometry_type', 'pass': has_cyl})

    # Fill height
    fill_h = gt.get('fluid', {}).get('fill_height')
    if fill_h:
        fluid_box = params.get('fluid_box', {})
        actual_fill = fluid_box.get('z', 0) if fluid_box else 0
        tol = fill_h * tol_pct
        passed = abs(actual_fill - fill_h) <= tol
        checks.append({'name': 'fill_height', 'pass': passed})

    # Motion type (always exact)
    motion_gt = gt.get('motion', {})
    if motion_gt.get('xml_tag'):
        actual_motion = params.get('motion_type', 'none')
        passed = actual_motion == motion_gt['xml_tag']
        checks.append({'name': 'motion_type', 'pass': passed})

    # Frequency
    if motion_gt.get('freq_hz'):
        actual_freq = params.get('freq_x', 0) or params.get('freq_y', 0)
        expected_freq = motion_gt['freq_hz']
        tol = expected_freq * tol_pct
        passed = abs(actual_freq - expected_freq) <= tol
        checks.append({'name': 'frequency', 'pass': passed})

    # Amplitude
    if motion_gt.get('amplitude_m'):
        actual_ampl = params.get('ampl_x', 0) or params.get('ampl_y', 0)
        expected_ampl = motion_gt['amplitude_m']
        tol = expected_ampl * tol_pct
        passed = abs(actual_ampl - expected_ampl) <= tol
        checks.append({'name': 'amplitude', 'pass': passed})
    elif motion_gt.get('amplitude_deg'):
        actual_motion = params.get('motion_type', 'none')
        if actual_motion == 'mvrotsinu':
            actual_ampl = params.get('ampl_x', 0) or params.get('ampl_y', 0)
            exp_deg = motion_gt['amplitude_deg']
            if isinstance(exp_deg, list):
                exp_deg = exp_deg[0]
            tol = exp_deg * tol_pct
            passed = abs(actual_ampl - exp_deg) <= tol
        else:
            passed = False
        checks.append({'name': 'amplitude', 'pass': passed})

    # TimeMax
    if gt.get('timemax'):
        actual_tm = params.get('timemax', 0)
        expected_tm = gt['timemax']
        tol = expected_tm * tol_pct
        passed = abs(actual_tm - expected_tm) <= tol
        checks.append({'name': 'timemax', 'pass': passed})

    passed_count = sum(1 for c in checks if c['pass'])
    total = len(checks)
    return checks, passed_count, total


def main():
    # Parse all XMLs once
    all_params = {}
    for s in SCENARIOS:
        for m in MODELS:
            result_dir = RESULTS_DIR / f"{s}_{m}_trial1"
            if not result_dir.exists():
                continue
            xml_path = find_xml(result_dir)
            if xml_path is None:
                continue
            validity = check_xml_validity(str(xml_path))
            if not validity["parseable"]:
                continue
            params, err = parse_xml_params(str(xml_path))
            if params:
                all_params[(s, m)] = params

    # Calculate M-A3 for each tolerance
    print("=" * 80)
    print("  Tolerance Sensitivity Analysis: M-A3 vs. Tolerance Threshold")
    print("=" * 80)

    results = {}  # (tol, model) → overall M-A3

    for tol in TOLERANCES:
        for model in MODELS:
            scenario_scores = []
            for s in SCENARIOS:
                params = all_params.get((s, model))
                if params is None:
                    scenario_scores.append(0)
                    continue
                checks, passed, total = score_with_tolerance(s, params, tol)
                pct = (passed / total * 100) if total > 0 else 0
                scenario_scores.append(pct)
            overall = sum(scenario_scores) / len(scenario_scores) if scenario_scores else 0
            results[(tol, model)] = overall

    # Print table
    print(f"\n  {'Tolerance':>10}", end="")
    for m in MODELS:
        label = "32B" if "32b" in m else "8B"
        print(f"  {label:>8}", end="")
    print(f"  {'Δ':>6}")
    print(f"  {'─' * 38}")

    for tol in TOLERANCES:
        v32 = results.get((tol, "qwen3_32b"), 0)
        v8 = results.get((tol, "qwen3_latest"), 0)
        delta = v32 - v8
        marker = " ◀ paper" if tol == 0.10 else ""
        # Note: actual paper uses mixed 10%/15%, this uses uniform
        print(f"  {tol*100:>8.0f}%  {v32:>7.1f}%  {v8:>7.1f}%  {delta:>+5.1f}%{marker}")

    # Per-scenario breakdown at key tolerances
    for tol in [0.05, 0.10, 0.15]:
        print(f"\n\n  === Per-Scenario M-A3 at {tol*100:.0f}% Tolerance ===")
        print(f"  {'Scenario':>8}", end="")
        for m in MODELS:
            label = "32B" if "32b" in m else "8B"
            print(f"  {label:>6}", end="")
        print()
        print(f"  {'─' * 26}")

        for s in SCENARIOS:
            print(f"  {s:>8}", end="")
            for m in MODELS:
                params = all_params.get((s, m))
                if params is None:
                    print(f"  {'N/A':>6}", end="")
                    continue
                checks, passed, total = score_with_tolerance(s, params, tol)
                pct = (passed / total * 100) if total > 0 else 0
                print(f"  {pct:>5.0f}%", end="")
            print()

    # Paper-relevant comparison: uniform 10% vs mixed (current paper)
    print(f"\n\n{'=' * 80}")
    print("  Comparison: Uniform vs Mixed Tolerance (Paper)")
    print("=" * 80)
    print("""
  Current paper uses mixed tolerance:
    tank: 10%, fill_height: 15%, freq: 10%, amplitude: 15%, timemax: 10%
    → 32B=61.2%, 8B=58.7%

  This analysis uses uniform tolerance per level.
  Key insight: M-A3 is most sensitive between 5% and 10% threshold.
  At ≥15%, tolerance has diminishing returns — primary failures are
  motion_type (exact match) and amplitude (unit conversion), which
  are tolerance-INDEPENDENT errors.
""")

    # Export for plotting
    export = {
        "tolerances": [t * 100 for t in TOLERANCES],
        "qwen3_32b": [results[(t, "qwen3_32b")] for t in TOLERANCES],
        "qwen3_latest": [results[(t, "qwen3_latest")] for t in TOLERANCES],
    }
    out_path = Path(__file__).parent / "tolerance_sensitivity.json"
    with open(out_path, "w") as f:
        json.dump(export, f, indent=2)
    print(f"\nExported: {out_path}")


if __name__ == "__main__":
    main()
