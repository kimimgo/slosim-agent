#!/usr/bin/env python3
"""SWL time-series comparison plot: baseline vs optimized iterations."""
import sys
import csv
import glob
from pathlib import Path

try:
    import matplotlib
    matplotlib.use('Agg')
    import matplotlib.pyplot as plt
    import numpy as np
except ImportError:
    print("matplotlib/numpy required: pip install matplotlib numpy")
    sys.exit(1)


def load_swl_csv(csv_path: str) -> tuple:
    """Load time-series SWL data from MeasureTool CSV output.
    Returns (time_array, swl_array)."""
    times, values = [], []
    with open(csv_path) as f:
        reader = csv.reader(f, delimiter=';')
        for row in reader:
            if not row or row[0].startswith('#'):
                continue
            try:
                t = float(row[0])
                # Take max absolute value across all measurement points
                vals = [abs(float(v)) for v in row[1:] if v.strip()]
                if vals:
                    times.append(t)
                    values.append(max(vals))
            except (ValueError, IndexError):
                continue
    return np.array(times), np.array(values)


def find_swl_file(iteration_dir: Path) -> str | None:
    """Find the SWL/height CSV file in an iteration directory."""
    patterns = ["*[Hh]eight*.csv", "*swl*.csv", "*[Pp]oint*.csv"]
    for pattern in patterns:
        files = list(iteration_dir.glob(pattern))
        files = [f for f in files if f.name != "Run.csv"]
        if files:
            return str(files[0])
    return None


def main():
    if len(sys.argv) < 2:
        results_dir = Path(__file__).parent.parent / "results"
    else:
        results_dir = Path(sys.argv[1])

    fig, ax = plt.subplots(figsize=(10, 6))

    colors = {'baseline': '#2196F3', 'iter_1': '#FF9800', 'iter_2': '#4CAF50', 'iter_3': '#F44336'}
    labels = {'baseline': 'Baseline (no baffle)', 'iter_1': 'Iteration 1', 'iter_2': 'Iteration 2', 'iter_3': 'Iteration 3'}

    for iter_name in ['baseline', 'iter_1', 'iter_2', 'iter_3']:
        iter_dir = results_dir / iter_name
        if not iter_dir.exists():
            continue

        swl_file = find_swl_file(iter_dir)
        if swl_file is None:
            continue

        times, swl = load_swl_csv(swl_file)
        if len(times) == 0:
            continue

        ax.plot(times, swl, label=labels.get(iter_name, iter_name),
                color=colors.get(iter_name, 'gray'), linewidth=1.5)

    ax.set_xlabel('Time (s)', fontsize=12)
    ax.set_ylabel('Surface Wave Level (m)', fontsize=12)
    ax.set_title('EXP-D: SWL Comparison — Baseline vs Optimized', fontsize=14)
    ax.legend(loc='upper right', fontsize=10)
    ax.grid(True, alpha=0.3)
    ax.set_xlim(left=0)

    output_path = results_dir.parent / "analysis" / "swl_comparison.png"
    output_path.parent.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"Saved: {output_path}")
    plt.close()


if __name__ == "__main__":
    main()
