#!/usr/bin/env python3
"""
EXP-A XML Parameter Fidelity Analyzer (M-A3)
Compares generated XML against ground truth parameters.
Usage: python3 analyze_xml.py <results_dir> [--ground-truth ground_truth.json]
"""
import json
import sys
import os
import re
import glob
from xml.etree import ElementTree as ET
from pathlib import Path


def find_xml_files(result_dir):
    """Find XML files in result directory."""
    patterns = [
        os.path.join(result_dir, "*.xml"),
        os.path.join(result_dir, "simulations", "**", "*.xml"),
    ]
    xmls = []
    for pat in patterns:
        xmls.extend(glob.glob(pat, recursive=True))
    return [x for x in xmls if "_Def" not in os.path.basename(x) or "sloshing" in os.path.basename(x).lower()]


def parse_xml_params(xml_path):
    """Extract key parameters from DualSPHysics XML."""
    try:
        tree = ET.parse(xml_path)
        root = tree.getroot()
    except ET.ParseError as e:
        return {"error": f"XML parse error: {e}", "valid": False}

    params = {"valid": True, "file": xml_path}

    # Tank dimensions (from boundary drawbox)
    for drawbox in root.iter("drawbox"):
        boxfill = drawbox.find("boxfill")
        if boxfill is not None and "bottom" in (boxfill.text or ""):
            size = drawbox.find("size")
            if size is not None:
                params["tank_x"] = float(size.get("x", 0))
                params["tank_y"] = float(size.get("y", 0))
                params["tank_z"] = float(size.get("z", 0))

    # Fluid dimensions (from fluid drawbox)
    for drawbox in root.iter("drawbox"):
        boxfill = drawbox.find("boxfill")
        if boxfill is not None and "solid" in (boxfill.text or ""):
            size = drawbox.find("size")
            if size is not None:
                params["fluid_x"] = float(size.get("x", 0))
                params["fluid_y"] = float(size.get("y", 0))
                params["fluid_z"] = float(size.get("z", 0))

    # dp
    definition = root.find(".//definition")
    if definition is not None:
        params["dp"] = float(definition.get("dp", 0))

    # rhop0
    for rhop in root.iter("rhop0"):
        params["rhop0"] = float(rhop.get("value", 0))

    # Motion
    for motion_tag in ["mvrectsinu", "mvrotsinu", "mvrectunif", "mvrotunif"]:
        for elem in root.iter(motion_tag):
            params["motion_type"] = motion_tag
            params["motion_duration"] = float(elem.get("duration", 0))

            freq = elem.find("freq")
            if freq is not None:
                params["freq_x"] = float(freq.get("x", 0))
                params["freq_y"] = float(freq.get("y", 0))
                params["freq_z"] = float(freq.get("z", 0))

            ampl = elem.find("ampl")
            if ampl is not None:
                params["ampl_x"] = float(ampl.get("x", 0))
                params["ampl_y"] = float(ampl.get("y", 0))
                params["ampl_z"] = float(ampl.get("z", 0))

    # Viscosity
    for param in root.iter("parameter"):
        key = param.get("key", "")
        val = param.get("value", "")
        if key == "Visco":
            params["visco"] = float(val)
        elif key == "ViscoTreatment":
            params["visco_treatment"] = int(val)
        elif key == "TimeMax":
            params["timemax"] = float(val)
        elif key == "TimeOut":
            params["timeout"] = float(val)

    # Gauges
    gauge_count = 0
    for gauge in root.iter("swl"):
        gauge_count += 1
    for gauge in root.iter("fixed"):
        gauge_count += 1
    params["gauge_count"] = gauge_count

    return params


def check_parameter(name, expected, actual, tolerance=0.05):
    """Check if a parameter matches within tolerance."""
    if expected is None:
        return {"name": name, "status": "SKIP", "expected": None, "actual": actual}
    if actual is None:
        return {"name": name, "status": "FAIL", "expected": expected, "actual": None, "reason": "missing"}

    if isinstance(expected, str):
        match = str(actual).lower() == expected.lower()
        return {"name": name, "status": "PASS" if match else "FAIL",
                "expected": expected, "actual": actual}

    try:
        exp_f = float(expected)
        act_f = float(actual)
        if exp_f == 0:
            match = abs(act_f) < 1e-6
        else:
            match = abs(act_f - exp_f) / abs(exp_f) <= tolerance
        return {"name": name, "status": "PASS" if match else "FAIL",
                "expected": exp_f, "actual": act_f,
                "error_pct": abs(act_f - exp_f) / max(abs(exp_f), 1e-10) * 100}
    except (ValueError, TypeError):
        return {"name": name, "status": "FAIL", "expected": expected, "actual": actual,
                "reason": "type_mismatch"}


def analyze_scenario(scenario_id, gt, xml_params):
    """Analyze one scenario against ground truth."""
    results = {"scenario": scenario_id, "tier": gt["tier"], "checks": []}

    if not xml_params.get("valid", False):
        results["ma2_xml_valid"] = False
        results["ma3_score"] = 0
        results["total_checks"] = 0
        results["passed_checks"] = 0
        return results

    results["ma2_xml_valid"] = True

    # Tank dimensions
    if "x" in gt.get("tank", {}):
        results["checks"].append(check_parameter("tank_x", gt["tank"]["x"], xml_params.get("tank_x")))
        results["checks"].append(check_parameter("tank_y", gt["tank"]["y"], xml_params.get("tank_y")))
        results["checks"].append(check_parameter("tank_z", gt["tank"]["z"], xml_params.get("tank_z")))

    # Fluid height
    if gt.get("fluid", {}).get("fill_height"):
        results["checks"].append(check_parameter("fill_height", gt["fluid"]["fill_height"],
                                                   xml_params.get("fluid_z"), tolerance=0.1))

    # Density (strict check: if non-default specified, actual must NOT be default 1000)
    if gt.get("fluid", {}).get("density") and gt["fluid"]["density"] != 1000:
        actual_rho = xml_params.get("rhop0")
        if actual_rho == 1000:
            results["checks"].append({"name": "rhop0", "status": "FAIL",
                                       "expected": gt["fluid"]["density"], "actual": actual_rho,
                                       "reason": "used_default_water_density"})
        else:
            results["checks"].append(check_parameter("rhop0", gt["fluid"]["density"], actual_rho, tolerance=0.02))

    # Motion type
    if gt.get("motion", {}).get("xml_tag"):
        results["checks"].append(check_parameter("motion_type", gt["motion"]["xml_tag"],
                                                   xml_params.get("motion_type")))

    # Frequency
    if gt.get("motion", {}).get("freq_hz"):
        # Determine which axis has the frequency
        actual_freq = max(xml_params.get("freq_x", 0), xml_params.get("freq_y", 0), xml_params.get("freq_z", 0))
        results["checks"].append(check_parameter("frequency", gt["motion"]["freq_hz"], actual_freq))

    # Amplitude (motion-type dependent)
    motion = gt.get("motion", {})
    if motion.get("amplitude_m"):
        actual_ampl = max(xml_params.get("ampl_x", 0), xml_params.get("ampl_y", 0), xml_params.get("ampl_z", 0))
        results["checks"].append(check_parameter("amplitude_m", motion["amplitude_m"], actual_ampl))
    elif motion.get("amplitude_deg"):
        amp_deg = motion["amplitude_deg"]
        if isinstance(amp_deg, list):
            amp_deg = amp_deg[0]  # Check first value for parametric
        actual_ampl = max(xml_params.get("ampl_x", 0), xml_params.get("ampl_y", 0), xml_params.get("ampl_z", 0))
        results["checks"].append(check_parameter("amplitude_deg", amp_deg, actual_ampl, tolerance=0.15))

    # TimeMax
    if gt.get("timemax"):
        results["checks"].append(check_parameter("timemax", gt["timemax"],
                                                   xml_params.get("timemax") or xml_params.get("motion_duration"),
                                                   tolerance=0.2))

    # Viscosity (special case for non-water fluids)
    if gt.get("viscosity", {}).get("value"):
        results["checks"].append(check_parameter("viscosity", gt["viscosity"]["value"],
                                                   xml_params.get("visco"), tolerance=0.2))

    # dp range
    if gt.get("dp_range") and xml_params.get("dp"):
        dp = xml_params["dp"]
        in_range = gt["dp_range"][0] <= dp <= gt["dp_range"][1]
        results["checks"].append({"name": "dp_range", "status": "PASS" if in_range else "WARN",
                                   "expected": gt["dp_range"], "actual": dp})

    # Calculate scores
    scoreable = [c for c in results["checks"] if c["status"] in ("PASS", "FAIL")]
    results["total_checks"] = len(scoreable)
    results["passed_checks"] = sum(1 for c in scoreable if c["status"] == "PASS")
    results["ma3_score"] = results["passed_checks"] / max(results["total_checks"], 1)

    return results


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 analyze_xml.py <results_dir> [--ground-truth gt.json]")
        sys.exit(1)

    results_dir = sys.argv[1]
    gt_file = "ground_truth.json"

    if "--ground-truth" in sys.argv:
        idx = sys.argv.index("--ground-truth")
        gt_file = sys.argv[idx + 1]

    # Load ground truth
    script_dir = os.path.dirname(os.path.abspath(__file__))
    gt_path = os.path.join(script_dir, gt_file) if not os.path.isabs(gt_file) else gt_file
    with open(gt_path) as f:
        ground_truth = json.load(f)

    # Process each result directory
    all_results = []

    if os.path.isdir(results_dir):
        # Single result directory (e.g., S01_qwen3_32b_trial1)
        scenario_match = re.match(r"(S\d+)_", os.path.basename(results_dir))
        if scenario_match:
            dirs = [results_dir]
        else:
            # Batch directory — find all subdirectories
            dirs = sorted(glob.glob(os.path.join(results_dir, "S*_*/")))
    else:
        print(f"ERROR: {results_dir} not found")
        sys.exit(1)

    for result_path in dirs:
        dirname = os.path.basename(result_path.rstrip("/"))
        scenario_match = re.match(r"(S\d+)_", dirname)
        if not scenario_match:
            continue

        scenario_id = scenario_match.group(1)
        if scenario_id not in ground_truth:
            print(f"WARNING: No ground truth for {scenario_id}")
            continue

        # Find XML files
        xml_files = find_xml_files(result_path)

        if not xml_files:
            result = {"scenario": scenario_id, "tier": ground_truth[scenario_id]["tier"],
                      "ma2_xml_valid": False, "ma3_score": 0, "error": "No XML found"}
            all_results.append(result)
            continue

        # Analyze first (main) XML
        xml_params = parse_xml_params(xml_files[0])
        result = analyze_scenario(scenario_id, ground_truth[scenario_id], xml_params)
        result["xml_file"] = xml_files[0]
        all_results.append(result)

    # Summary output
    print("\n" + "=" * 70)
    print("EXP-A: M-A3 Parameter Fidelity Analysis")
    print("=" * 70)

    for r in all_results:
        status = "PASS" if r.get("ma3_score", 0) >= 0.8 else ("PARTIAL" if r.get("ma3_score", 0) >= 0.5 else "FAIL")
        score_str = f"{r.get('passed_checks', 0)}/{r.get('total_checks', 0)}"
        print(f"\n{r['scenario']} [{r['tier']}] — M-A3: {status} ({score_str} = {r.get('ma3_score', 0):.0%})")

        if r.get("error"):
            print(f"  ERROR: {r['error']}")
            continue

        for check in r.get("checks", []):
            icon = "✓" if check["status"] == "PASS" else ("!" if check["status"] == "WARN" else "✗")
            print(f"  {icon} {check['name']}: expected={check.get('expected')}, actual={check.get('actual')}", end="")
            if "error_pct" in check:
                print(f" (err={check['error_pct']:.1f}%)", end="")
            print()

    # JSON summary
    summary = {
        "total_scenarios": len(all_results),
        "ma2_pass": sum(1 for r in all_results if r.get("ma2_xml_valid")),
        "ma3_pass": sum(1 for r in all_results if r.get("ma3_score", 0) >= 0.8),
        "ma3_partial": sum(1 for r in all_results if 0.5 <= r.get("ma3_score", 0) < 0.8),
        "ma3_fail": sum(1 for r in all_results if r.get("ma3_score", 0) < 0.5),
        "results": all_results,
    }

    summary_path = os.path.join(results_dir if os.path.isdir(results_dir) else ".", "ma3_summary.json")
    with open(summary_path, "w") as f:
        json.dump(summary, f, indent=2, ensure_ascii=False)
    print(f"\nSummary saved to {summary_path}")

    return summary


if __name__ == "__main__":
    main()
