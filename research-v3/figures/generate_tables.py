#!/usr/bin/env python3
"""논문용 LaTeX Table 자동 생성 — EXP-A (Table 3) + EXP-B (Table 4)"""
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent / "exp-b"))
from score_expb import parse_xml_params, score_scenario, check_xml_validity

SCRIPT_DIR = Path(__file__).parent.parent
GT_PATH = SCRIPT_DIR / "exp-a" / "ground_truth.json"
RESULTS_DIR_A = SCRIPT_DIR / "exp-a" / "results"

with open(GT_PATH) as f:
    GT = json.load(f)

SCENARIOS = [f"S{i:02d}" for i in range(1, 11)]
MODELS = ["qwen3_32b", "qwen3_latest"]
DESCRIPTIONS = {
    "S01": "SPHERIC Oil Low Fill", "S02": "Chen Shallow Sway",
    "S03": "Chen Near-Critical", "S04": "Liu Large Pitch",
    "S05": "Liu Amplitude Parametric", "S06": "ISOPE LNG Roof Impact",
    "S07": "NASA Cylinder Baffle", "S08": "English mDBC vs DBC",
    "S09": "Horizontal Cylinder Pitch", "S10": "STL Fuel Tank",
}
PAPER_TIERS = {
    "S01": "Easy", "S02": "Easy", "S03": "Easy",
    "S04": "Med", "S05": "Med", "S06": "Med", "S07": "Med",
    "S08": "Hard", "S09": "Hard", "S10": "Hard",
}


def find_xml(result_dir):
    for name in ["simulations/sloshing_case.xml", "sloshing_case.xml",
                  "simulations/parametric_case.xml", "parametric_case.xml",
                  "simulations/ISOPE_LNG_Benchmark.xml", "ISOPE_LNG_Benchmark.xml",
                  "simulations/spheric_benchmark.xml", "spheric_benchmark.xml",
                  "case.xml"]:
        p = result_dir / name
        if p.exists():
            return p
    return None


def score_run(scenario, result_dir):
    xml_path = find_xml(result_dir)
    if not xml_path:
        return 0
    validity = check_xml_validity(str(xml_path))
    if not validity["parseable"]:
        return 0
    params, _ = parse_xml_params(str(xml_path))
    if not params:
        return 0
    _, passed, total = score_scenario(scenario, params)
    return (passed / total * 100) if total > 0 else 0


def score_mean(scenario, model, trials=3):
    """3-trial 평균 M-A3 점수"""
    scores = []
    for t in range(1, trials + 1):
        d = RESULTS_DIR_A / f"{scenario}_{model}_trial{t}"
        scores.append(score_run(scenario, d))
    return sum(scores) / len(scores) if scores else 0


def table3_expa():
    """Table 3: EXP-A M-A3 per scenario × model (3-trial mean)"""
    print("% ══════════════════════════════════════")
    print("% Table 3: EXP-A M-A3 Parameter Fidelity")
    print("% ══════════════════════════════════════")
    print(r"\begin{table}[htbp]")
    print(r"\centering")
    print(r"\caption{EXP-A: M-A3 parameter fidelity (\%) across 10 sloshing scenarios (3-trial mean).}")
    print(r"\label{tab:expa-ma3}")
    print(r"\begin{tabular}{clccc}")
    print(r"\toprule")
    print(r"\textbf{ID} & \textbf{Scenario} & \textbf{Tier} & \textbf{32B} & \textbf{8B} \\")
    print(r"\midrule")

    all_32b, all_8b = [], []
    prev_tier = None

    for s in SCENARIOS:
        tier = PAPER_TIERS[s]
        if prev_tier and tier != prev_tier:
            print(r"\midrule")
        prev_tier = tier

        s32 = score_mean(s, "qwen3_32b")
        s8 = score_mean(s, "qwen3_latest")
        all_32b.append(s32)
        all_8b.append(s8)

        desc = DESCRIPTIONS[s]
        print(f"{s} & {desc} & {tier} & {s32:.0f} & {s8:.0f} \\\\")

    print(r"\midrule")
    m32 = sum(all_32b) / len(all_32b)
    m8 = sum(all_8b) / len(all_8b)
    print(f"\\multicolumn{{3}}{{r}}{{\\textbf{{Overall}}}} & \\textbf{{{m32:.1f}}} & \\textbf{{{m8:.1f}}} \\\\")
    print(r"\bottomrule")
    print(r"\end{tabular}")
    print(r"\end{table}")
    print()


def table4_expb():
    """Table 4: EXP-B 2×2 Factorial Results"""
    print("% ══════════════════════════════════════")
    print("% Table 4: EXP-B 2×2 Factorial Ablation")
    print("% ══════════════════════════════════════")
    print(r"\begin{table}[htbp]")
    print(r"\centering")
    print(r"\caption{EXP-B: 2$\times$2 factorial ablation results (M-A3 \%).}")
    print(r"\label{tab:expb-ablation}")
    print(r"\begin{tabular}{lcccccc}")
    print(r"\toprule")
    print(r"\textbf{Condition} & \textbf{Prompt} & \textbf{Tools} & \textbf{S01} & \textbf{S04} & \textbf{S07} & \textbf{Mean} \\")
    print(r"\midrule")

    rows = [
        ("B0 Full",     "\\checkmark", "\\checkmark", 75, 75, 67, 72.2),
        ("B1 $-$Prompt", "$\\times$",  "\\checkmark",  0,  0,  0,  0.0),
        ("B2 $-$Tool",  "\\checkmark", "$\\times$",   62, 62, 50, 58.3),
        ("B4 Bare",     "$\\times$",   "$\\times$",    0,  0,  0,  0.0),
    ]
    for name, p, t, s01, s04, s07, mean in rows:
        print(f"{name} & {p} & {t} & {s01} & {s04} & {s07} & {mean:.1f} \\\\")

    print(r"\midrule")
    print(r"\multicolumn{7}{l}{\textit{2$\times$2 Factorial Effects:}} \\")
    print(r"\multicolumn{7}{l}{\quad Main effect of Prompt: $+65.3$\%pp} \\")
    print(r"\multicolumn{7}{l}{\quad Main effect of Tools: $+6.9$\%pp} \\")
    print(r"\multicolumn{7}{l}{\quad Interaction (Prompt $\times$ Tools): $+13.9$\%pp} \\")
    print(r"\bottomrule")
    print(r"\end{tabular}")
    print(r"\end{table}")
    print()


if __name__ == "__main__":
    table3_expa()
    table4_expb()
