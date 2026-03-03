#!/usr/bin/env python3
"""Tool-Induced Error Analysis — B0(Full) vs B2(−Tool) 파라미터별 비교

핵심 질문: 도구가 각 파라미터의 정확도를 올리는가, 내리는가?
B0 = Prompt + Tools, B2 = Prompt only → 차이 = Tool의 순 효과
"""
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent / "exp-b"))
from score_expb import parse_xml_params, score_scenario, check_xml_validity

SCRIPT_DIR = Path(__file__).parent.parent
GT_PATH = SCRIPT_DIR / "exp-a" / "ground_truth.json"
RESULTS_A = SCRIPT_DIR / "exp-a" / "results"
RESULTS_B = SCRIPT_DIR / "exp-b" / "results"

with open(GT_PATH) as f:
    GT = json.load(f)

ABLATION_SCENARIOS = ["S01", "S04", "S07"]


def find_xml(result_dir):
    for name in ["simulations/sloshing_case.xml", "sloshing_case.xml",
                  "simulations/parametric_case.xml", "parametric_case.xml",
                  "simulations/ISOPE_LNG_Benchmark.xml", "ISOPE_LNG_Benchmark.xml",
                  "simulations/spheric_benchmark.xml", "spheric_benchmark.xml",
                  "generated.xml", "case.xml"]:
        p = result_dir / name
        if p.exists():
            return p
    return None


def get_checks(scenario, result_dir):
    xml_path = find_xml(result_dir)
    if not xml_path:
        return None, None
    validity = check_xml_validity(str(xml_path))
    if not validity["parseable"]:
        return None, None
    params, _ = parse_xml_params(str(xml_path))
    if not params:
        return None, None
    checks, _, _ = score_scenario(scenario, params)
    return checks, params


def main():
    print("=" * 72)
    print("  Tool-Induced Error Analysis: B0 (Full) vs B2 (−Tool)")
    print("=" * 72)
    print()

    # Per-parameter comparison
    param_stats = {}  # param_name → {b0_pass, b0_total, b2_pass, b2_total}

    for scenario in ABLATION_SCENARIOS:
        print(f"\n{'─' * 60}")
        print(f"  {scenario}: {GT[scenario]['description']}")
        print(f"{'─' * 60}")

        # B0 (from EXP-A trial1)
        b0_dir = RESULTS_A / f"{scenario}_qwen3_32b_trial1"
        b0_checks, b0_params = get_checks(scenario, b0_dir)

        # B2 (from EXP-B)
        b2_dir = RESULTS_B / f"B2_{scenario}_qwen3_32b"
        b2_checks, b2_params = get_checks(scenario, b2_dir)

        if not b0_checks or not b2_checks:
            print("  Missing data")
            continue

        # Compare parameter-by-parameter
        b0_dict = {c['name']: c for c in b0_checks}
        b2_dict = {c['name']: c for c in b2_checks}

        all_params = sorted(set(list(b0_dict.keys()) + list(b2_dict.keys())))

        print(f"  {'Parameter':<18} {'B0 (Tools)':>14} {'B2 (No Tool)':>14} {'Tool Effect':>14}")
        print(f"  {'─' * 60}")

        for p in all_params:
            b0 = b0_dict.get(p, {})
            b2 = b2_dict.get(p, {})

            b0_pass = b0.get('pass', False) if b0 else False
            b2_pass = b2.get('pass', False) if b2 else False

            b0_val = b0.get('actual', '—')
            b2_val = b2.get('actual', '—')
            exp_val = b0.get('expected', b2.get('expected', '—'))

            if b0_pass and not b2_pass:
                effect = "Tool HELPS"
            elif not b0_pass and b2_pass:
                effect = "Tool HURTS"
            elif b0_pass and b2_pass:
                effect = "Both OK"
            else:
                effect = "Both FAIL"

            b0_mark = "✓" if b0_pass else "✗"
            b2_mark = "✓" if b2_pass else "✗"

            print(f"  {p:<18} {b0_mark} ({b0_val!s:>8}) {b2_mark} ({b2_val!s:>8}) {effect}")

            if p not in param_stats:
                param_stats[p] = {'b0_pass': 0, 'b0_total': 0, 'b2_pass': 0, 'b2_total': 0}
            param_stats[p]['b0_total'] += 1
            param_stats[p]['b2_total'] += 1
            if b0_pass:
                param_stats[p]['b0_pass'] += 1
            if b2_pass:
                param_stats[p]['b2_pass'] += 1

    # Summary
    print(f"\n\n{'=' * 72}")
    print("  SUMMARY: Tool Effect per Parameter (across S01, S04, S07)")
    print(f"{'=' * 72}")
    print(f"\n  {'Parameter':<18} {'B0 Pass':>10} {'B2 Pass':>10} {'Δ':>8} {'Verdict':>14}")
    print(f"  {'─' * 60}")

    tool_helps = 0
    tool_hurts = 0
    neutral = 0

    for p in sorted(param_stats.keys()):
        s = param_stats[p]
        b0_rate = s['b0_pass'] / s['b0_total'] if s['b0_total'] > 0 else 0
        b2_rate = s['b2_pass'] / s['b2_total'] if s['b2_total'] > 0 else 0
        delta = b0_rate - b2_rate

        if delta > 0:
            verdict = "Tool HELPS"
            tool_helps += 1
        elif delta < 0:
            verdict = "Tool HURTS"
            tool_hurts += 1
        else:
            verdict = "Neutral"
            neutral += 1

        print(f"  {p:<18} {s['b0_pass']}/{s['b0_total']:>2}       {s['b2_pass']}/{s['b2_total']:>2}       {delta:>+.2f}  {verdict}")

    print(f"\n  Tool helps: {tool_helps} parameters")
    print(f"  Tool hurts: {tool_hurts} parameters")
    print(f"  Neutral:    {neutral} parameters")

    # Calculate overall net effect
    b0_total_pass = sum(s['b0_pass'] for s in param_stats.values())
    b0_total_checks = sum(s['b0_total'] for s in param_stats.values())
    b2_total_pass = sum(s['b2_pass'] for s in param_stats.values())
    b2_total_checks = sum(s['b2_total'] for s in param_stats.values())

    b0_pct = b0_total_pass / b0_total_checks * 100
    b2_pct = b2_total_pass / b2_total_checks * 100

    print(f"\n  Overall: B0={b0_pct:.1f}% vs B2={b2_pct:.1f}% → Tool net effect: {b0_pct - b2_pct:+.1f}%pp")
    print(f"  (B0 = {b0_total_pass}/{b0_total_checks} checks, B2 = {b2_total_pass}/{b2_total_checks} checks)")

    # Key insight
    print(f"\n{'=' * 72}")
    print("  KEY INSIGHT")
    print(f"{'=' * 72}")
    print("""
  Tool는 geometry (tank_x/y/z, fill_height)를 보호하지만,
  motion_type을 강제 오버라이드하여 정확도를 낮춘다.

  → Tool 순 효과 = geometry 보호 − motion_type 오류 도입
  → Prompt가 의미(무엇을), Tool이 구조(어떻게) 담당하나,
    Tool이 의미 영역(motion_type)까지 침범할 때 오류 발생.
""")


if __name__ == "__main__":
    main()
