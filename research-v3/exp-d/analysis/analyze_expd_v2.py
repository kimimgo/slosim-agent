#!/usr/bin/env python3
"""EXP-D v2: Baseline vs Baffle SWL from MeasureTool pressure probes.
Detects surface level as highest z where pressure > threshold."""
import sys
from pathlib import Path
import numpy as np

try:
    import matplotlib
    matplotlib.use('Agg')
    import matplotlib.pyplot as plt
except ImportError:
    print("matplotlib/numpy required")
    sys.exit(1)

RESULTS_DIR = Path(__file__).parent.parent / "results"
Z_PROBES = np.arange(0.0, 0.40, 0.02)  # 20 probes: 0.00, 0.02, ..., 0.38
P_THRESHOLD = 10.0  # Pa — minimum pressure to consider as "fluid present"


def load_press_probes(csv_path):
    """Load MeasureTool multi-probe pressure CSV. Returns (times, surface_z)."""
    times, surface_z = [], []
    with open(csv_path) as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith('#') or line.startswith('Part') or line.startswith(' '):
                if line.startswith('Part'):
                    continue  # skip header
                continue
            parts = line.split(';')
            if len(parts) < 22:
                continue
            try:
                t = float(parts[1])
                pressures = [float(parts[2 + i]) for i in range(20)]

                # Find highest z where pressure > threshold
                max_z = 0.0
                for i in range(19, -1, -1):  # top to bottom
                    if pressures[i] > P_THRESHOLD:
                        max_z = Z_PROBES[i]
                        # Interpolate: if next probe up has p < threshold
                        if i < 19 and pressures[i + 1] <= P_THRESHOLD:
                            # Linear interpolation between probe i and i+1
                            p_i = pressures[i]
                            dz = Z_PROBES[i + 1] - Z_PROBES[i]
                            max_z = Z_PROBES[i] + dz * (p_i / (p_i + P_THRESHOLD)) * 0.5
                        break

                times.append(t)
                surface_z.append(max_z)
            except (ValueError, IndexError):
                continue
    return np.array(times), np.array(surface_z)


def compute_amplitude(times, swlz, t_start=2.0):
    """Compute peak-to-peak amplitude after t_start (skip transient)."""
    mask = times >= t_start
    if mask.sum() < 5:
        return 0, 0, 0
    z = swlz[mask]
    z_mean = np.mean(z)
    z_max = np.max(z)
    z_min = np.min(z)
    amplitude = (z_max - z_min) / 2
    return amplitude, z_mean, z_max - z_min


def main():
    files = {
        'baseline_left': RESULTS_DIR / 'baseline_swl_left_Press.csv',
        'baseline_right': RESULTS_DIR / 'baseline_swl_right_Press.csv',
        'baffle_left': RESULTS_DIR / 'baffle_swl_left_Press.csv',
        'baffle_right': RESULTS_DIR / 'baffle_swl_right_Press.csv',
    }

    missing = [k for k, v in files.items() if not v.exists()]
    if missing:
        print(f"Missing: {missing}")
        sys.exit(1)

    data = {}
    for name, path in files.items():
        t, z = load_press_probes(path)
        data[name] = (t, z)
        print(f"  {name}: {len(t)} points, t=[{t[0]:.2f}, {t[-1]:.2f}]s, "
              f"z=[{z.min():.3f}, {z.max():.3f}]m")

    print("\n" + "=" * 60)
    print("  EXP-D: SWL Amplitude Analysis (MeasureTool probes)")
    print("=" * 60)

    results = {}
    for side in ['left', 'right']:
        bl_amp, bl_mean, bl_pp = compute_amplitude(*data[f'baseline_{side}'])
        bf_amp, bf_mean, bf_pp = compute_amplitude(*data[f'baffle_{side}'])
        reduction = (1 - bf_amp / bl_amp) * 100 if bl_amp > 0 else 0
        results[side] = {'bl_amp': bl_amp, 'bf_amp': bf_amp, 'reduction': reduction}

        print(f"\n  {side.upper()} wall (t>2s):")
        print(f"    Baseline:  mean={bl_mean:.4f}m, half-amp={bl_amp:.4f}m, p-p={bl_pp:.4f}m")
        print(f"    Baffle:    mean={bf_mean:.4f}m, half-amp={bf_amp:.4f}m, p-p={bf_pp:.4f}m")
        print(f"    Reduction: {reduction:.1f}%")

    # Combined average
    bl_avg = (results['left']['bl_amp'] + results['right']['bl_amp']) / 2
    bf_avg = (results['left']['bf_amp'] + results['right']['bf_amp']) / 2
    avg_reduction = (1 - bf_avg / bl_avg) * 100 if bl_avg > 0 else 0

    print(f"\n  COMBINED (avg L+R):")
    print(f"    Baseline amplitude: {bl_avg:.4f}m")
    print(f"    Baffle amplitude:   {bf_avg:.4f}m")
    print(f"    Average reduction:  {avg_reduction:.1f}%")
    print(f"    Expected range:     30-50% (literature)")
    verdict = 'PASS' if 15 <= avg_reduction <= 70 else 'CHECK'
    print(f"    Verdict:            {verdict}")

    # Plot
    fig, axes = plt.subplots(2, 1, figsize=(12, 8), sharex=True)

    for i, side in enumerate(['left', 'right']):
        ax = axes[i]
        t_bl, z_bl = data[f'baseline_{side}']
        t_bf, z_bf = data[f'baffle_{side}']

        ax.plot(t_bl, z_bl, 'b-', alpha=0.8, linewidth=1.0, label='Baseline (no baffle)')
        ax.plot(t_bf, z_bf, 'r-', alpha=0.8, linewidth=1.0, label='With center baffle')

        # Equilibrium line
        ax.axhline(y=0.2, color='gray', linestyle='--', alpha=0.4, label='Equilibrium (0.2m)')
        ax.axvline(x=2.0, color='green', linestyle=':', alpha=0.3)

        red = results[side]['reduction']
        ax.set_ylabel(f'Surface Level (m)', fontsize=11)
        ax.set_title(f'{side.capitalize()} Wall — Amplitude Reduction: {red:.0f}%', fontsize=12)
        ax.legend(loc='upper right', fontsize=9)
        ax.grid(True, alpha=0.3)
        ax.set_ylim(0, 0.40)

    axes[1].set_xlabel('Time (s)', fontsize=11)
    fig.suptitle(f'EXP-D: Baffle Effectiveness — {avg_reduction:.0f}% Average SWL Reduction',
                 fontsize=14, fontweight='bold')
    fig.tight_layout()

    out_png = Path(__file__).parent / 'expd_swl_comparison.png'
    out_pdf = Path(__file__).parent / 'expd_swl_comparison.pdf'
    fig.savefig(out_png, dpi=150)
    fig.savefig(out_pdf, bbox_inches='tight')
    print(f"\n  Figures saved: {out_png}")
    print(f"                 {out_pdf}")
    plt.close()


if __name__ == '__main__':
    main()
