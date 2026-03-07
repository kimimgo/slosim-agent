#!/usr/bin/env python3
"""Baffle placement evolution visualization — tank cross-section with baffle positions."""
import sys
import json
import re
from pathlib import Path

try:
    import matplotlib
    matplotlib.use('Agg')
    import matplotlib.pyplot as plt
    import matplotlib.patches as patches
    import numpy as np
except ImportError:
    print("matplotlib/numpy required: pip install matplotlib numpy")
    sys.exit(1)


def parse_baffle_from_xml(xml_path: str) -> list[dict]:
    """Extract baffle positions from DualSPHysics XML (drawbox after setmkbound mk>=10)."""
    baffles = []
    content = Path(xml_path).read_text()

    # Find all setmkbound + drawbox pairs for baffles (mk >= 10)
    pattern = r'<setmkbound\s+mk="(\d+)".*?<drawbox>\s*<boxfill>solid</boxfill>\s*<point\s+x="([^"]+)"\s+y="([^"]+)"\s+z="([^"]+)".*?<size\s+x="([^"]+)"\s+y="([^"]+)"\s+z="([^"]+)"'
    for match in re.finditer(pattern, content, re.DOTALL):
        mk = int(match.group(1))
        if mk >= 10:  # Baffle mk numbers start at 10
            baffles.append({
                'mk': mk,
                'x': float(match.group(2)),
                'y': float(match.group(3)),
                'z': float(match.group(4)),
                'sx': float(match.group(5)),
                'sy': float(match.group(6)),
                'sz': float(match.group(7)),
            })
    return baffles


def main():
    if len(sys.argv) < 2:
        results_dir = Path(__file__).parent.parent / "results"
    else:
        results_dir = Path(sys.argv[1])

    # Tank dimensions (Frosina2018 fuel_tank approximate)
    tank_x, tank_z = 0.500, 0.250  # Length x Height (XZ cross-section)
    fluid_height = 0.125  # 50% fill

    iterations = ['baseline', 'iter_1', 'iter_2', 'iter_3']
    n_plots = sum(1 for it in iterations if (results_dir / it).exists())
    if n_plots == 0:
        print("No iteration directories found")
        return

    fig, axes = plt.subplots(1, n_plots, figsize=(5 * n_plots, 4), squeeze=False)
    axes = axes[0]

    plot_idx = 0
    for iter_name in iterations:
        iter_dir = results_dir / iter_name
        if not iter_dir.exists():
            continue

        ax = axes[plot_idx]

        # Draw tank outline
        tank_rect = patches.Rectangle((0, 0), tank_x, tank_z, linewidth=2,
                                       edgecolor='black', facecolor='none')
        ax.add_patch(tank_rect)

        # Draw fluid
        fluid_rect = patches.Rectangle((0, 0), tank_x, fluid_height,
                                        alpha=0.3, facecolor='#2196F3', edgecolor='none')
        ax.add_patch(fluid_rect)

        # Draw baffles from XML files
        xml_files = list(iter_dir.glob("*.xml"))
        for xf in xml_files:
            baffles = parse_baffle_from_xml(str(xf))
            for b in baffles:
                baffle_rect = patches.Rectangle(
                    (b['x'], b['z']), b['sx'], b['sz'],
                    linewidth=1, edgecolor='red', facecolor='red', alpha=0.7
                )
                ax.add_patch(baffle_rect)

        ax.set_xlim(-0.02, tank_x + 0.02)
        ax.set_ylim(-0.02, tank_z + 0.02)
        ax.set_aspect('equal')
        ax.set_xlabel('X (m)')
        ax.set_ylabel('Z (m)')

        title = iter_name.replace('_', ' ').title()
        if iter_name == 'baseline':
            title = 'Baseline (No Baffle)'
        ax.set_title(title, fontsize=11)
        ax.grid(True, alpha=0.2)

        plot_idx += 1

    fig.suptitle('EXP-D: Baffle Placement Evolution', fontsize=14, fontweight='bold')
    fig.tight_layout()

    output_path = results_dir.parent / "analysis" / "baffle_evolution.png"
    output_path.parent.mkdir(parents=True, exist_ok=True)
    fig.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"Saved: {output_path}")
    plt.close()


if __name__ == "__main__":
    main()
