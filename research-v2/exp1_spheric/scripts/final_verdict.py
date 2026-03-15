#!/usr/bin/env python3
"""EXP-1 Final Verdict — Reads all metrics and produces overall PASS/FAIL.

Outputs a concise summary suitable for Telegram notification.

Usage:
    python final_verdict.py [--telegram]
"""

import json
import sys
from pathlib import Path

ANALYSIS_DIR = Path(__file__).resolve().parent.parent / "analysis"
FIG_DIR = Path(__file__).resolve().parent.parent / "figures"


def load_json(path):
    if not path.exists():
        return None
    with open(path) as f:
        return json.load(f)


def main():
    send_telegram = '--telegram' in sys.argv

    water = load_json(ANALYSIS_DIR / "metrics.json")
    oil_roof = load_json(ANALYSIS_DIR / "oil_roof_metrics.json")

    verdicts = {}
    summary_lines = []

    summary_lines.append("=== EXP-1 SPHERIC Test 10 — Final Verdict ===\n")

    # --- Water Lateral ---
    if water and "runs" in water:
        summary_lines.append("[ Water Lateral ]")
        best = water["runs"].get("run_003") or water["runs"].get("run_002")
        run_label = "Run 003 (dp=1mm)" if "run_003" in water["runs"] else "Run 002 (dp=2mm)"

        if best:
            m1 = best.get("M1_pass", False)
            m2 = best.get("M2_pass", False)
            m5 = best.get("M5_pass", False)
            m6 = best.get("M6_pass", False)

            summary_lines.append(f"  Best: {run_label}")
            summary_lines.append(f"  M1 (peak-in-band): {best.get('in_band', '?')} {'PASS' if m1 else 'FAIL'}")
            summary_lines.append(f"  M2 (error <30%): {best.get('mean_abs_error_pct', '?')}% {'PASS' if m2 else 'FAIL'}")
            summary_lines.append(f"  M5 (corr >0.5): r={best.get('cross_corr_r', '?')} {'PASS' if m5 else 'FAIL'}")
            summary_lines.append(f"  M6 (time shift): tau={best.get('time_shift_s', '?')}s {'PASS' if m6 else 'FAIL'}")

            water_pass = m1 and m2 and m5 and m6
            verdicts["water_lateral"] = water_pass
            summary_lines.append(f"  -> Water Lateral: {'PASS' if water_pass else 'FAIL'}")
        else:
            summary_lines.append("  No simulation data available")
            verdicts["water_lateral"] = False
    else:
        summary_lines.append("[ Water Lateral ] — metrics.json not found")
        verdicts["water_lateral"] = False

    summary_lines.append("")

    # --- Oil Lateral ---
    if oil_roof and "oil_lateral" in oil_roof:
        oil = oil_roof["oil_lateral"]
        summary_lines.append("[ Oil Lateral ]")
        summary_lines.append(f"  Peaks detected: {oil.get('peaks_detected', 0)}")
        summary_lines.append(f"  M7 (>=3/4 detected): {'PASS' if oil.get('M7_pass') else 'FAIL'}")
        summary_lines.append(f"  M1 (>=2/4 in-band): {oil.get('M1_in_band', '?')} {'PASS' if oil.get('M1_pass') else 'FAIL'}")

        oil_pass = oil.get("M7_pass", False) and oil.get("M1_pass", False)
        verdicts["oil_lateral"] = oil_pass
        summary_lines.append(f"  -> Oil Lateral: {'PASS' if oil_pass else 'FAIL'}")
    else:
        summary_lines.append("[ Oil Lateral ] — data not available")
        verdicts["oil_lateral"] = False

    summary_lines.append("")

    # --- Water Roof ---
    if oil_roof and "water_roof" in oil_roof:
        roof = oil_roof["water_roof"]
        summary_lines.append("[ Water Roof ]")
        summary_lines.append(f"  M1 (>=3/4 in-band): {roof.get('M1_in_band', '?')} {'PASS' if roof.get('M1_pass') else 'FAIL'}")
        summary_lines.append(f"  M2 (<40%): {roof.get('M2_mae_pct', '?')}% {'PASS' if roof.get('M2_pass') else 'FAIL'}")

        roof_pass = roof.get("M1_pass", False) and roof.get("M2_pass", False)
        verdicts["water_roof"] = roof_pass
        summary_lines.append(f"  -> Water Roof: {'PASS' if roof_pass else 'FAIL'}")
    else:
        summary_lines.append("[ Water Roof ] — data not available")
        verdicts["water_roof"] = False

    summary_lines.append("")

    # --- Overall ---
    # Validation plan: Water PASS + Oil PASS → overall PASS
    # Water Roof is optional bonus (DBC limitation documented)
    required_pass = (verdicts.get("water_lateral", False) and
                     verdicts.get("oil_lateral", False))
    all_pass = all(verdicts.values())
    summary_lines.append("=" * 40)
    if all_pass:
        summary_lines.append("EXP-1 OVERALL: PASS (all sub-cases pass)")
    elif required_pass:
        roof_status = "PARTIAL" if not verdicts.get("water_roof", False) else "PASS"
        summary_lines.append(f"EXP-1 OVERALL: PASS (Water+Oil PASS, Roof {roof_status} — optional)")
    elif verdicts.get("water_lateral", False):
        summary_lines.append("EXP-1 OVERALL: PARTIAL (Water PASS, Oil pending/fail)")
    else:
        summary_lines.append("EXP-1 OVERALL: FAIL")

    summary_text = "\n".join(summary_lines)
    print(summary_text)

    # --- Figures ---
    figs = list(FIG_DIR.glob("*.png"))
    if figs:
        print(f"\nFigures ({len(figs)}):")
        for f in sorted(figs):
            print(f"  {f.name}")

    # Return exit code based on verdict (0 if required sub-cases pass)
    return 0 if required_pass else 1


if __name__ == "__main__":
    sys.exit(main())
