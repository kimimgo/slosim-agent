#!/usr/bin/env python3
"""EXP-D: Baseline vs Baffle SWL comparison analysis.
Reads DualSPHysics GaugesSWL CSV files, computes amplitude reduction."""
import sys
from pathlib import Path
import numpy as np

try:
    import matplotlib
    matplotlib.use('Agg')
    import matplotlib.pyplot as plt
except ImportError:
    print("matplotlib/numpy required: pip install matplotlib numpy")
    sys.exit(1)

RESULTS_DIR = Path(__file__).parent.parent / "results"

def load_swl(csv_path):
    """Load SWL gauge CSV. Returns (time, swlz) arrays."""
    times, swlz = [], []
    with open(csv_path) as f:
        for line in f:
            if line.startswith('time') or line.startswith('#'):
                continue
            parts = line.strip().split(';')
            if len(parts) < 4:
                continue
            try:
                t = float(parts[0])
                z = float(parts[3])  # swlz column
                times.append(t)
                swlz.append(z)
            except ValueError:
                continue
    return np.array(times), np.array(swlz)

def compute_amplitude(times, swlz, t_start=2.0):
    """Compute peak-to-peak amplitude after t_start (skip transient)."""
    mask = times >= t_start
    if mask.sum() < 10:
        return 0, 0, 0
    z = swlz[mask]
    z_mean = np.mean(z)
    z_max = np.max(z)
    z_min = np.min(z)
    amplitude = (z_max - z_min) / 2  # half peak-to-peak
    return amplitude, z_mean, z_max - z_min

def main():
    files = {
        'baseline_left': RESULTS_DIR / 'baseline_SWL_Left.csv',
        'baseline_right': RESULTS_DIR / 'baseline_SWL_Right.csv',
        'baffle_left': RESULTS_DIR / 'baffle_SWL_Left.csv',
        'baffle_right': RESULTS_DIR / 'baffle_SWL_Right.csv',
    }

    # Check files exist
    missing = [k for k, v in files.items() if not v.exists()]
    if missing:
        print(f"Missing files: {missing}")
        print(f"Expected in: {RESULTS_DIR}")
        sys.exit(1)

    # Load data
    data = {}
    for name, path in files.items():
        t, z = load_swl(path)
        data[name] = (t, z)
        print(f"  Loaded {name}: {len(t)} points, t=[{t[0]:.3f}, {t[-1]:.3f}]s")

    # Compute amplitudes (skip first 2s of transient)
    print("\n" + "=" * 60)
    print("  EXP-D: SWL Amplitude Analysis")
    print("=" * 60)

    for side in ['left', 'right']:
        bl_amp, bl_mean, bl_pp = compute_amplitude(*data[f'baseline_{side}'])
        bf_amp, bf_mean, bf_pp = compute_amplitude(*data[f'baffle_{side}'])
        reduction = (1 - bf_amp / bl_amp) * 100 if bl_amp > 0 else 0

        print(f"\n  {side.upper()} wall (t>2s):")
        print(f"    Baseline:  mean={bl_mean:.4f}m, amplitude={bl_amp:.4f}m, p-p={bl_pp:.4f}m")
        print(f"    Baffle:    mean={bf_mean:.4f}m, amplitude={bf_amp:.4f}m, p-p={bf_pp:.4f}m")
        print(f"    Reduction: {reduction:.1f}%")

    # Combined (average of left and right)
    bl_amp_l, _, _ = compute_amplitude(*data['baseline_left'])
    bl_amp_r, _, _ = compute_amplitude(*data['baseline_right'])
    bf_amp_l, _, _ = compute_amplitude(*data['baffle_left'])
    bf_amp_r, _, _ = compute_amplitude(*data['baffle_right'])
    bl_avg = (bl_amp_l + bl_amp_r) / 2
    bf_avg = (bf_amp_l + bf_amp_r) / 2
    avg_reduction = (1 - bf_avg / bl_avg) * 100 if bl_avg > 0 else 0

    print(f"\n  COMBINED (avg L+R):")
    print(f"    Baseline amplitude: {bl_avg:.4f}m")
    print(f"    Baffle amplitude:   {bf_avg:.4f}m")
    print(f"    Average reduction:  {avg_reduction:.1f}%")
    print(f"    Expected range:     30-50%")
    print(f"    {'PASS' if 20 <= avg_reduction <= 70 else 'CHECK'}")

    # Plot
    fig, axes = plt.subplots(2, 1, figsize=(12, 8), sharex=True)

    for i, side in enumerate(['left', 'right']):
        ax = axes[i]
        t_bl, z_bl = data[f'baseline_{side}']
        t_bf, z_bf = data[f'baffle_{side}']

        # Subsample for plotting (every 10th point)
        step = max(1, len(t_bl) // 2000)
        ax.plot(t_bl[::step], z_bl[::step], 'b-', alpha=0.7, linewidth=0.8, label='Baseline')
        step = max(1, len(t_bf) // 2000)
        ax.plot(t_bf[::step], z_bf[::step], 'r-', alpha=0.7, linewidth=0.8, label='With Baffle')

        bl_amp, _, _ = compute_amplitude(t_bl, z_bl)
        bf_amp, _, _ = compute_amplitude(t_bf, z_bf)
        red = (1 - bf_amp / bl_amp) * 100 if bl_amp > 0 else 0

        ax.set_ylabel(f'SWL {side.capitalize()} (m)', fontsize=11)
        ax.set_title(f'{side.capitalize()} Wall — Amplitude Reduction: {red:.0f}%', fontsize=12)
        ax.legend(loc='upper right')
        ax.grid(True, alpha=0.3)
        ax.axvline(x=2.0, color='gray', linestyle='--', alpha=0.5, label='Transient cutoff')

    axes[1].set_xlabel('Time (s)', fontsize=11)
    fig.suptitle(f'EXP-D: Baffle Effectiveness — {avg_reduction:.0f}% Average SWL Reduction',
                 fontsize=14, fontweight='bold')
    fig.tight_layout()

    out_path = Path(__file__).parent / 'expd_swl_comparison.png'
    fig.savefig(out_path, dpi=150)
    print(f"\n  Figure saved: {out_path}")
    plt.close()

if __name__ == '__main__':
    main()
