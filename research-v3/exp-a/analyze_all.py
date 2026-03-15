#!/usr/bin/env python3
"""EXP-A 논문용 종합 분석 — score_expb.py와 동일한 M-A3 기준 사용

M-A3 채점 파라미터 (8개 고정):
  1. tank_x    2. tank_y    3. tank_z (또는 geometry_type)
  4. fill_height
  5. motion_type    6. frequency    7. amplitude    8. timemax
"""
import json
import sys
from pathlib import Path
from collections import defaultdict

# score_expb.py의 동일한 채점 로직 재사용
sys.path.insert(0, str(Path(__file__).parent.parent / "exp-b"))
from score_expb import parse_xml_params, score_scenario, check_xml_validity

SCRIPT_DIR = Path(__file__).parent
GT_PATH = SCRIPT_DIR / "ground_truth.json"
RESULTS_DIR = SCRIPT_DIR / "results"

with open(GT_PATH) as f:
    GT = json.load(f)

SCENARIOS = [f"S{i:02d}" for i in range(1, 11)]
MODELS = ["qwen3_32b", "qwen3_latest"]
TRIALS = ["trial1", "trial2", "trial3"]
# 논문 기준 난이도 그룹 (시나리오 순서 기반)
TIER_GROUPS = {
    "Easy": ["S01", "S02", "S03"],
    "Medium": ["S04", "S05", "S06", "S07"],
    "Hard": ["S08", "S09", "S10"],
}


def find_xml(result_dir: Path):
    """EXP-A 결과 디렉토리에서 XML 찾기 — score_expb B0 로직과 동일."""
    candidates = [
        result_dir / "simulations" / "sloshing_case.xml",
        result_dir / "sloshing_case.xml",
    ]
    # 대안 파일명
    alt_names = ["parametric_case.xml", "ISOPE_LNG_Benchmark.xml",
                 "spheric_benchmark.xml", "case.xml"]
    for alt in alt_names:
        candidates.append(result_dir / "simulations" / alt)
        candidates.append(result_dir / alt)

    for c in candidates:
        if c.exists():
            return c
    return None


def main():
    scores = {}   # (scenario, model, trial) → M-A3 %
    details = {}  # (scenario, model, trial) → check list

    for s in SCENARIOS:
        for m in MODELS:
            for t in TRIALS:
                key = (s, m, t)
                result_dir = RESULTS_DIR / f"{s}_{m}_{t}"
                if not result_dir.exists():
                    scores[key] = -1
                    continue

                xml_path = find_xml(result_dir)
                if xml_path is None:
                    scores[key] = 0
                    details[key] = []
                    continue

                validity = check_xml_validity(str(xml_path))
                if not validity["parseable"]:
                    scores[key] = 0
                    details[key] = []
                    continue

                params, err = parse_xml_params(str(xml_path))
                if params is None:
                    scores[key] = 0
                    details[key] = []
                    continue

                checks, passed, total = score_scenario(s, params)
                pct = (passed / total * 100) if total > 0 else 0
                scores[key] = pct
                details[key] = checks

    # ═══════════════════════════════════════════════════════
    # 1. 시나리오별 M-A3 테이블
    # ═══════════════════════════════════════════════════════
    print("=" * 80)
    print("  EXP-A: M-A3 Parameter Fidelity (score_expb 기준, 8-param)")
    print("=" * 80)

    for model in MODELS:
        print(f"\n  Model: {model}")
        print(f"  {'Scenario':<6} {'T1':>6} {'T2':>6} {'T3':>6} {'Mean':>8} {'σ':>6}")
        print(f"  {'─' * 42}")

        all_means = []
        for s in SCENARIOS:
            trial_scores = []
            row = f"  {s:<6}"
            for t in TRIALS:
                sc = scores.get((s, model, t), -1)
                if sc < 0:
                    row += f" {'N/A':>6}"
                else:
                    row += f" {sc:>5.0f}%"
                    trial_scores.append(sc)

            if trial_scores:
                mean = sum(trial_scores) / len(trial_scores)
                stddev = (sum((x - mean) ** 2 for x in trial_scores) / len(trial_scores)) ** 0.5
                row += f" {mean:>7.1f}% {stddev:>5.1f}"
                all_means.append(mean)
            else:
                row += f" {'—':>8} {'—':>6}"
            print(row)

        if all_means:
            print(f"\n  {'Overall':>6} {'':>6} {'':>6} {'':>6} {sum(all_means)/len(all_means):>7.1f}%")

    # ═══════════════════════════════════════════════════════
    # 2. 난이도 그룹별 요약 (논문 Table 용)
    # ═══════════════════════════════════════════════════════
    print(f"\n\n{'=' * 80}")
    print("  Paper Table: Tier × Model (Mean M-A3 %)")
    print("=" * 80)

    print(f"\n  {'Tier':<10} {'32B':>8} {'8B':>8} {'Δ':>8}")
    print(f"  {'─' * 36}")

    overall_32b, overall_8b = [], []
    for tier, scenarios in TIER_GROUPS.items():
        vals_32b, vals_8b = [], []
        for s in scenarios:
            for m, lst in [("qwen3_32b", vals_32b), ("qwen3_latest", vals_8b)]:
                ts = [scores[(s, m, t)] for t in TRIALS if scores.get((s, m, t), -1) >= 0]
                if ts:
                    lst.append(sum(ts) / len(ts))

        m32 = sum(vals_32b) / len(vals_32b) if vals_32b else None
        m8 = sum(vals_8b) / len(vals_8b) if vals_8b else None
        diff = (m32 - m8) if (m32 is not None and m8 is not None) else None

        r32 = f"{m32:>7.1f}%" if m32 is not None else f"{'N/A':>8}"
        r8 = f"{m8:>7.1f}%" if m8 is not None else f"{'N/A':>8}"
        rd = f"{diff:>+7.1f}%" if diff is not None else f"{'—':>8}"
        print(f"  {tier:<10} {r32} {r8} {rd}")

        if m32 is not None:
            overall_32b.extend(vals_32b)
        if m8 is not None:
            overall_8b.extend(vals_8b)

    if overall_32b and overall_8b:
        m32a = sum(overall_32b) / len(overall_32b)
        m8a = sum(overall_8b) / len(overall_8b)
        print(f"  {'─' * 36}")
        print(f"  {'Overall':<10} {m32a:>7.1f}% {m8a:>7.1f}% {m32a - m8a:>+7.1f}%")

    # ═══════════════════════════════════════════════════════
    # 3. Determinism
    # ═══════════════════════════════════════════════════════
    print(f"\n\n{'=' * 80}")
    print("  Determinism: M-A3 StdDev Across 3 Trials")
    print("=" * 80)

    for model in MODELS:
        print(f"\n  {model}:")
        det_count = 0
        for s in SCENARIOS:
            ts = [scores.get((s, model, t), -1) for t in TRIALS]
            valid = [x for x in ts if x >= 0]
            if len(valid) >= 2:
                mean = sum(valid) / len(valid)
                std = (sum((x - mean) ** 2 for x in valid) / len(valid)) ** 0.5
                if std == 0:
                    det_count += 1
                mark = "✓" if std == 0 else f"σ={std:.1f}"
                print(f"    {s}: {valid} → {mark}")
        print(f"  Deterministic: {det_count}/10")

    # ═══════════════════════════════════════════════════════
    # 4. Parameter-Level Error Patterns
    # ═══════════════════════════════════════════════════════
    print(f"\n\n{'=' * 80}")
    print("  Error Patterns (Trial 1, Both Models)")
    print("=" * 80)

    patterns = defaultdict(lambda: {"count": 0, "scenarios": set()})

    for s in SCENARIOS:
        for m in MODELS:
            checks = details.get((s, m, "trial1"), [])
            for c in checks:
                if not c["pass"]:
                    p = patterns[c["name"]]
                    p["count"] += 1
                    p["scenarios"].add(s)

    print(f"\n  {'Parameter':<16} {'Fails':>6} {'Scenarios'}")
    print(f"  {'─' * 55}")
    for name, info in sorted(patterns.items(), key=lambda x: -x[1]["count"]):
        ss = ", ".join(sorted(info["scenarios"]))
        print(f"  {name:<16} {info['count']:>6} {ss}")

    # ═══════════════════════════════════════════════════════
    # 5. JSON export
    # ═══════════════════════════════════════════════════════
    output = {
        "experiment": "EXP-A",
        "scoring": "score_expb 8-parameter M-A3",
        "total_runs": sum(1 for v in scores.values() if v >= 0),
        "expected_runs": 60,
        "scores": {},
        "per_model_overall": {},
    }
    for s in SCENARIOS:
        for m in MODELS:
            ts = [scores[(s, m, t)] for t in TRIALS if scores.get((s, m, t), -1) >= 0]
            for t in TRIALS:
                k = f"{s}_{m}_{t}"
                sc = scores.get((s, m, t), -1)
                if sc >= 0:
                    output["scores"][k] = sc

    for m in MODELS:
        vals = [scores[(s, m, "trial1")] for s in SCENARIOS
                if scores.get((s, m, "trial1"), -1) >= 0]
        if vals:
            output["per_model_overall"][m] = sum(vals) / len(vals)

    out_path = RESULTS_DIR / "exp_a_full_analysis.json"
    with open(out_path, "w") as f:
        json.dump(output, f, indent=2)
    print(f"\n\nSaved: {out_path}")


if __name__ == "__main__":
    main()
