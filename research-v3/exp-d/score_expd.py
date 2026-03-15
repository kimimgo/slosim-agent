#!/usr/bin/env python3
"""EXP-D Scoring: M-D1 through M-D5 metrics."""
import sys
import os
import csv
import json
import glob
from pathlib import Path

def score_md1_stl_loading(results_dir: str) -> dict:
    """M-D1: STL loading success — GenCase produced output without error."""
    baseline = Path(results_dir) / "baseline"
    # Check for any generated case files (bi4 or XML output)
    xml_files = list(baseline.glob("*.xml"))
    log_files = list(baseline.glob("*.log"))

    success = len(xml_files) > 0
    error_msg = ""
    if log_files:
        for lf in log_files:
            content = lf.read_text()
            if "error" in content.lower() or "exception" in content.lower():
                success = False
                error_msg = f"Error found in {lf.name}"
                break

    return {"metric": "M-D1", "name": "STL Loading Success", "pass": success, "detail": error_msg or "OK"}


def score_md2_baseline_complete(results_dir: str) -> dict:
    """M-D2: Baseline simulation completed — Run.csv exists and has data."""
    baseline = Path(results_dir) / "baseline"
    run_csv = baseline / "Run.csv"

    if not run_csv.exists():
        # Try finding it in subdirectories
        candidates = list(baseline.rglob("Run.csv"))
        if candidates:
            run_csv = candidates[0]
        else:
            return {"metric": "M-D2", "name": "Baseline Complete", "pass": False, "detail": "Run.csv not found"}

    try:
        with open(run_csv) as f:
            lines = f.readlines()
        # Run.csv should have header + data lines
        data_lines = [l for l in lines if l.strip() and not l.startswith("#") and not l.startswith("Run")]
        return {
            "metric": "M-D2",
            "name": "Baseline Complete",
            "pass": len(data_lines) > 1,
            "detail": f"{len(data_lines)} data lines in Run.csv"
        }
    except Exception as e:
        return {"metric": "M-D2", "name": "Baseline Complete", "pass": False, "detail": str(e)}


def score_md3_swl_reduction(results_dir: str) -> dict:
    """M-D3: SWL reduction rate = (baseline - best_optimized) / baseline."""
    results = Path(results_dir)

    def get_max_swl(iteration_dir: Path) -> float:
        """Extract max SWL from MeasureTool CSV output."""
        csv_files = list(iteration_dir.glob("*[Hh]eight*.csv")) + list(iteration_dir.glob("*swl*.csv"))
        if not csv_files:
            csv_files = list(iteration_dir.glob("*.csv"))

        max_swl = 0.0
        for cf in csv_files:
            if cf.name == "Run.csv":
                continue
            try:
                with open(cf) as f:
                    reader = csv.reader(f, delimiter=';')
                    for row in reader:
                        try:
                            values = [abs(float(v)) for v in row if v.strip()]
                            if values:
                                max_swl = max(max_swl, max(values))
                        except ValueError:
                            continue
            except Exception:
                continue
        return max_swl

    baseline_swl = get_max_swl(results / "baseline")

    best_swl = float('inf')
    best_iter = ""
    for iter_name in ["iter_1", "iter_2", "iter_3"]:
        iter_dir = results / iter_name
        if iter_dir.exists():
            swl = get_max_swl(iter_dir)
            if 0 < swl < best_swl:
                best_swl = swl
                best_iter = iter_name

    if baseline_swl <= 0 or best_swl == float('inf'):
        return {"metric": "M-D3", "name": "SWL Reduction", "pass": False,
                "detail": f"Cannot compute: baseline={baseline_swl:.4f}, best={best_swl}", "value": 0.0}

    reduction = (baseline_swl - best_swl) / baseline_swl
    return {
        "metric": "M-D3",
        "name": "SWL Reduction",
        "pass": reduction > 0,
        "detail": f"baseline={baseline_swl:.4f}m, best={best_swl:.4f}m ({best_iter}), reduction={reduction:.1%}",
        "value": reduction
    }


def score_md4_iteration_count(results_dir: str) -> dict:
    """M-D4: Number of autonomous iterations the agent performed."""
    results = Path(results_dir)
    count = 0
    for iter_name in ["iter_1", "iter_2", "iter_3"]:
        iter_dir = results / iter_name
        if iter_dir.exists() and any(iter_dir.iterdir()):
            count += 1

    return {
        "metric": "M-D4",
        "name": "Autonomous Iterations",
        "pass": count >= 1,
        "detail": f"{count} iterations completed",
        "value": count
    }


def score_md5_wall_clock(results_dir: str) -> dict:
    """M-D5: Total wall-clock time from NL input to final result."""
    wc_file = Path(results_dir) / "wall_clock_seconds.txt"
    if not wc_file.exists():
        return {"metric": "M-D5", "name": "Wall Clock Time", "pass": False, "detail": "wall_clock_seconds.txt not found"}

    try:
        seconds = int(wc_file.read_text().strip())
        minutes = seconds / 60
        return {
            "metric": "M-D5",
            "name": "Wall Clock Time",
            "pass": True,
            "detail": f"{minutes:.1f} minutes ({seconds}s)",
            "value": seconds
        }
    except Exception as e:
        return {"metric": "M-D5", "name": "Wall Clock Time", "pass": False, "detail": str(e)}


def main():
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <results_dir>")
        sys.exit(1)

    results_dir = sys.argv[1]
    if not os.path.isdir(results_dir):
        print(f"Error: {results_dir} is not a directory")
        sys.exit(1)

    scores = [
        score_md1_stl_loading(results_dir),
        score_md2_baseline_complete(results_dir),
        score_md3_swl_reduction(results_dir),
        score_md4_iteration_count(results_dir),
        score_md5_wall_clock(results_dir),
    ]

    print("\n" + "=" * 60)
    print("EXP-D Scoring Results")
    print("=" * 60)

    passed = 0
    for s in scores:
        status = "PASS" if s["pass"] else "FAIL"
        icon = "✓" if s["pass"] else "✗"
        print(f"  {icon} {s['metric']}: {s['name']} — {status}")
        print(f"    {s['detail']}")
        if s["pass"]:
            passed += 1

    print(f"\nTotal: {passed}/{len(scores)} metrics passed")
    print("=" * 60)

    # Save JSON
    output_file = os.path.join(results_dir, "scores.json")
    with open(output_file, "w") as f:
        json.dump(scores, f, indent=2, ensure_ascii=False)
    print(f"\nScores saved to {output_file}")

    return 0 if passed == len(scores) else 1


if __name__ == "__main__":
    sys.exit(main())
