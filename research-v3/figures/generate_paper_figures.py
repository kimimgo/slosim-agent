#!/usr/bin/env python3
"""Unified paper figure generation — 7 publication-quality figures.

논문 논리 순서:
  Fig 3: EXP-A Parameter Fidelity Heatmap (Gap1 증명)
  Fig 4: v3→v4 Tool Design Impact (도구 설계 병목)
  Fig 5: EXP-B Hierarchical Ablation (프롬프트=필수, 도구=증폭)
  Fig 6: SPHERIC Pressure Validation (Gap2 물리적 타당성)
  Fig 7: Convergence Metrics (수치 신뢰성)
  Fig 8: Bridge Physics Checks (Agent→Solver 검증)
  Fig 9: EXP-D Baffle Optimization (자율 설계 능력)

Usage:
  python3 generate_paper_figures.py              # all figures
  python3 generate_paper_figures.py --only 3 5   # specific figures

Output: PDF (vector) + PNG (300 dpi)
"""
import json
import argparse
from pathlib import Path

import numpy as np

import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import matplotlib.gridspec as gridspec
from matplotlib.colors import LinearSegmentedColormap

# ═══════════════════════════════════════════════════════════════
# Style Configuration
# ═══════════════════════════════════════════════════════════════
plt.rcParams.update({
    'font.family': 'serif',
    'font.serif': ['DejaVu Serif', 'Computer Modern Roman', 'Times New Roman'],
    'font.size': 9,
    'axes.labelsize': 10,
    'axes.titlesize': 11,
    'axes.titleweight': 'bold',
    'xtick.labelsize': 8,
    'ytick.labelsize': 8,
    'legend.fontsize': 8,
    'legend.framealpha': 0.9,
    'figure.dpi': 150,
    'savefig.dpi': 300,
    'savefig.bbox': 'tight',
    'savefig.pad_inches': 0.08,
    'axes.grid': False,
    'axes.spines.top': False,
    'axes.spines.right': False,
    'axes.linewidth': 0.8,
    'grid.linewidth': 0.4,
    'grid.alpha': 0.3,
    'lines.linewidth': 1.2,
})

# Color palette
C_32B = '#2166AC'
C_8B = '#E08214'
C_EASY = '#66BD63'
C_MED = '#FEE08B'
C_HARD = '#F46D43'
C_V3 = '#BDBDBD'
C_PROMPT = '#4393C3'
C_TOOL = '#D6604D'
C_EXP = '#333333'
C_SIM = ['#2166AC', '#E08214', '#1B7837']
C_BASELINE = '#2166AC'
C_BAFFLE = '#D6604D'

# Journal sizing (inches) — Springer CPM
COL_SINGLE = 3.54
COL_DOUBLE = 7.48

# Paths
SCRIPT_DIR = Path(__file__).parent
ROOT = SCRIPT_DIR.parent
REPO = ROOT.parent
OUT_DIR = SCRIPT_DIR / "paper"

V4_JSON = ROOT / "exp-a" / "v4_analysis_results.json"
BRIDGE_CSV = ROOT / "exp-c" / "agent-bridge" / "gauges_Press.csv"
EXPD_DIR = ROOT / "exp-d" / "results"
SIMDATA = Path("/mnt/simdata/dualsphysics/exp1")
EXPDATA = REPO / "datasets" / "spheric" / "case_1"

SCENARIOS = [f"S{i:02d}" for i in range(1, 11)]
DESCRIPTIONS = {
    "S01": "SPHERIC Oil", "S02": "Chen Shallow",
    "S03": "Chen Near-Crit.", "S04": "Liu Pitch 30s",
    "S05": "Liu Parametric", "S06": "ISOPE LNG",
    "S07": "NASA Cylinder", "S08": "English mDBC",
    "S09": "Horiz. Cyl.", "S10": "STL Fuel Tank",
}
TIER_COLORS = {"Easy": C_EASY, "Medium": C_MED, "Hard": C_HARD}


def save_fig(fig, name):
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    for ext in ('pdf', 'png'):
        fig.savefig(OUT_DIR / f"{name}.{ext}")
    plt.close(fig)
    print(f"  OK  {name}.pdf + .png")


def load_v4():
    with open(V4_JSON) as f:
        return json.load(f)


# ═══════════════════════════════════════════════════════════════
# Fig 3: EXP-A M-A3 Heatmap
# ═══════════════════════════════════════════════════════════════
def fig03_expa_heatmap():
    print("Fig 3: EXP-A Parameter Fidelity Heatmap")
    d = load_v4()
    s32 = [d["models"]["32B"]["v4_scores"][s] for s in SCENARIOS]
    s8 = [d["models"]["8B"]["v4_scores"][s] for s in SCENARIOS]
    matrix = np.array([s32, s8]).T

    fig, ax = plt.subplots(figsize=(COL_SINGLE + 0.8, 4.2))
    cmap = LinearSegmentedColormap.from_list(
        'ma3', ['#D73027', '#FC8D59', '#FEE08B', '#D9EF8B', '#1A9850'], N=256)
    im = ax.imshow(matrix, cmap=cmap, aspect='auto', vmin=0, vmax=100)

    # Tier bands + separators
    for y0, y1, c in [(0, 3, C_EASY), (3, 7, C_MED), (7, 10, C_HARD)]:
        ax.axhspan(y0 - 0.5, y1 - 0.5, facecolor=c, alpha=0.08, zorder=0)
    ax.axhline(2.5, color='#666', lw=0.6, alpha=0.5)
    ax.axhline(6.5, color='#666', lw=0.6, alpha=0.5)

    # Cell annotations
    for i in range(10):
        for j in range(2):
            v = matrix[i, j]
            c = 'white' if v > 90 or v < 40 else 'black'
            ax.text(j, i, f'{v:.0f}', ha='center', va='center',
                    fontsize=8.5, color=c, fontweight='bold' if v == 100 else 'normal')

    # Axes
    ax.set_xticks([0, 1])
    ax.set_xticklabels(['Qwen3-32B', 'Qwen3-8B'], fontsize=9, fontweight='bold')
    ax.xaxis.tick_top()
    ylabels = [f'{s}  {DESCRIPTIONS[s]}' for s in SCENARIOS]
    ax.set_yticks(range(10))
    ax.set_yticklabels(ylabels, fontsize=7.5)

    # Colorbar
    cbar = fig.colorbar(im, ax=ax, shrink=0.5, aspect=20, pad=0.12)
    cbar.set_label('M-A3 (%)', fontsize=8)
    cbar.ax.tick_params(labelsize=7)

    # Tier labels — use figure transform to place after colorbar
    cbar_right = cbar.ax.get_position().x1
    for y_data, label in [(1, 'Easy'), (5, 'Medium'), (8.5, 'Hard')]:
        # Convert data y to axes fraction, then to figure coords
        y_frac = ax.transData.transform((0, y_data))[1]
        y_fig = fig.transFigure.inverted().transform((0, y_frac))[1]
        fig.text(cbar_right + 0.03, y_fig, label, ha='left', va='center',
                 fontsize=7, fontweight='bold', color=TIER_COLORS[label])
    # Mean summary below heatmap
    m32 = d["models"]["32B"]["v4_mean"]
    m8 = d["models"]["8B"]["v4_mean"]
    ax.text(0, 10.6, f'{m32:.1f}%', ha='center', fontsize=10,
            fontweight='bold', color=C_32B)
    ax.text(1, 10.6, f'{m8:.1f}%', ha='center', fontsize=10,
            fontweight='bold', color=C_8B)
    ax.text(0.5, 11.4, f'Delta = {abs(m32 - m8):.1f} pp  (9/10 identical)',
            ha='center', fontsize=7.5, color='#555')

    ax.set_ylim(9.5, -0.5)
    fig.tight_layout(rect=[0, 0.06, 1, 1])
    save_fig(fig, 'fig03_expa_heatmap')


# ═══════════════════════════════════════════════════════════════
# Fig 4: v3 → v4 Tool Impact
# ═══════════════════════════════════════════════════════════════
def fig04_v3v4_impact():
    print("Fig 4: v3->v4 Tool Design Impact")
    d = load_v4()
    pitch = {'S01', 'S04', 'S05', 'S08', 'S09'}

    fig, axes = plt.subplots(1, 2, figsize=(COL_DOUBLE, 2.8), sharey=True)
    for ax_i, (model, title, c_v4) in enumerate([
        ('32B', 'Qwen3-32B', C_32B), ('8B', 'Qwen3-8B', C_8B)
    ]):
        ax = axes[ax_i]
        v3 = [d["models"][model]["v3_scores"][s] for s in SCENARIOS]
        v4 = [d["models"][model]["v4_scores"][s] for s in SCENARIOS]
        x = np.arange(10)
        w = 0.35

        ax.bar(x - w/2, v3, w, color=C_V3, label='v3 (baseline)',
               edgecolor='white', linewidth=0.3)
        ax.bar(x + w/2, v4, w, color=c_v4, label='v4 (+P1/P2 fix)',
               edgecolor='white', linewidth=0.3)

        for i, s in enumerate(SCENARIOS):
            delta = v4[i] - v3[i]
            if abs(delta) > 0.5:
                ax.text(i + w/2, v4[i] + 2, f'+{delta:.0f}',
                        ha='center', va='bottom', fontsize=6,
                        fontweight='bold', color='#1A9850')
            if s in pitch:
                ax.text(i, -8, '†', ha='center', fontsize=7, color=C_HARD)

        mean_v4 = np.mean(v4)
        ax.axhline(mean_v4, color=c_v4, ls='--', lw=0.7, alpha=0.5)
        ax.text(9.6, mean_v4 + 2, f'{mean_v4:.1f}%', fontsize=6.5,
                color=c_v4, fontweight='bold')

        ax.set_xticks(x)
        ax.set_xticklabels([s.replace('S0', 'S') for s in SCENARIOS], fontsize=7)
        ax.set_ylabel('M-A3 (%)' if ax_i == 0 else '')
        ax.set_ylim(-12, 115)
        ax.set_title(title, fontsize=10)
        ax.legend(fontsize=7, loc='upper center', ncol=2)
        for x0, x1, c in [(-0.5, 2.5, C_EASY), (2.5, 6.5, C_MED), (6.5, 9.5, C_HARD)]:
            ax.axvspan(x0, x1, facecolor=c, alpha=0.06)

    fig.text(0.5, -0.01, '† Pitch-motion scenarios affected by P1+P2 tool fix',
             ha='center', fontsize=7, color='#666', style='italic')
    fig.tight_layout()
    save_fig(fig, 'fig04_v3v4_impact')


# ═══════════════════════════════════════════════════════════════
# Fig 5: EXP-B Hierarchical Ablation
# ═══════════════════════════════════════════════════════════════
def fig05_expb_ablation():
    print("Fig 5: EXP-B Hierarchical Ablation")
    conds = [
        ('B4  Bare LLM',     0.0,  '#BDBDBD'),
        ('B1  -Prompt',      0.0,  '#9E9E9E'),
        ('B2  -Tools',      46.1,  C_PROMPT),
        ('B0  Full System', 67.0,  '#1A9850'),
    ]

    fig, ax = plt.subplots(figsize=(COL_SINGLE, 2.4))
    y = np.arange(4)
    vals = [c[1] for c in conds]
    bars = ax.barh(y, vals, height=0.6, color=[c[2] for c in conds],
                   edgecolor='white', linewidth=0.5, zorder=3)

    # Value labels inside/outside bars
    for i, v in enumerate(vals):
        if v > 10:
            ax.text(v - 1.5, i, f'{v:.0f}%', ha='right', va='center',
                    fontsize=9, fontweight='bold', color='white')
        else:
            ax.text(v + 1, i, f'{v:.0f}%', ha='left', va='center',
                    fontsize=9, fontweight='bold', color='#555')

    # Effect arrows (above the bars, no overlap)
    # Prompt effect: B4(0) → B2(46.1)
    ax.annotate('', xy=(46.1, 2.55), xytext=(0, 2.55),
                arrowprops=dict(arrowstyle='->', color=C_PROMPT, lw=1.5))
    ax.text(23, 2.85, '+46.1 pp (Prompt)', ha='center', fontsize=6.5,
            color=C_PROMPT, fontweight='bold')

    # Tool effect: B2(46.1) → B0(67)
    ax.annotate('', xy=(67, 3.55), xytext=(46.1, 3.55),
                arrowprops=dict(arrowstyle='->', color=C_TOOL, lw=1.5))
    ax.text(56.5, 3.8, '+20.9 pp (Tools)', ha='center', fontsize=6.5,
            color=C_TOOL, fontweight='bold')

    # Key insight
    ax.text(0.98, 0.02, 'Without domain prompt\nTools produce 0%',
            transform=ax.transAxes, ha='right', va='bottom',
            fontsize=6.5, color='#D32F2F', style='italic',
            bbox=dict(boxstyle='round,pad=0.3', fc='#FFEBEE', ec='#D32F2F', alpha=0.8))

    ax.set_yticks(y)
    ax.set_yticklabels([c[0] for c in conds], fontsize=7.5)
    ax.set_xlabel('M-A3 Score (%)', fontsize=9)
    ax.set_xlim(-2, 82)
    ax.set_ylim(-0.8, 4.3)
    ax.grid(axis='x', alpha=0.2, zorder=0)
    ax.spines['left'].set_visible(False)
    ax.tick_params(axis='y', length=0)

    fig.tight_layout()
    save_fig(fig, 'fig05_expb_ablation')


# ═══════════════════════════════════════════════════════════════
# SPHERIC data loading utilities
# ═══════════════════════════════════════════════════════════════
def load_exp_timeseries(fp):
    data = np.loadtxt(str(fp), skiprows=1)
    return data[:, 0], data[:, 1] * 100

def load_exp_peaks(fp):
    return np.loadtxt(str(fp), skiprows=2) * 100

def load_dsph_csv(fp, probe=0):
    t, p = [], []
    with open(fp) as f:
        for line in f:
            line = line.strip()
            if not line or line[0] in '#P ':
                continue
            parts = line.split(';')
            try:
                t.append(float(parts[1]))
                p.append(float(parts[2 + probe]))
            except (ValueError, IndexError):
                continue
    return np.array(t), np.array(p)

def smooth_p(p, w=5):
    if len(p) < w:
        return p
    return np.convolve(p, np.ones(w)/w, mode='same')

def find_peaks_filtered(t, p, min_h, min_gap, max_h=None):
    from scipy.signal import find_peaks
    if len(t) < 2:
        return [], []
    ps = smooth_p(p, 5)
    dt = np.median(np.diff(t))
    d = max(1, int(min_gap / dt))
    pks, _ = find_peaks(ps, height=min_h, distance=d, prominence=min_h * 0.3)
    valid = [i for i in pks if t[i] > 0.5]
    if max_h is not None:
        valid = [i for i in valid if ps[i] <= max_h]
    return [t[i] for i in valid], [ps[i] for i in valid]


# ═══════════════════════════════════════════════════════════════
# Fig 6: SPHERIC Pressure — Water + Oil Lateral (2×2 grid)
# ═══════════════════════════════════════════════════════════════
def _plot_spheric_row(axes, te, pe, pmean, p2sig, ept, sim, fluid, ylim_ts):
    """Plot one fluid row: (left) time-series, (right) peak bars."""
    ax_ts, ax_pk = axes

    # Time-series
    step = max(1, len(te) // 2000)
    ax_ts.plot(te[::step], pe[::step], '-', color=C_EXP, lw=0.4, alpha=0.5,
               label='Experiment', zorder=1)
    for r in sim.values():
        ts = r['t'] + r['shift']
        ax_ts.plot(ts, smooth_p(r['p'], 5), '-', color=r['col'], lw=r['lw'],
                   alpha=0.85, label=f"SPH {r['lab']}", zorder=3)
    for i in range(min(4, len(ept))):
        ax_ts.axvline(ept[i], color='#999', lw=0.4, ls=':', alpha=0.4)
        ax_ts.text(ept[i], ylim_ts * 0.9, f'P{i+1}', ha='center', fontsize=6, color='#666')
    ax_ts.set_ylabel('Pressure (Pa)')
    ax_ts.legend(fontsize=6, loc='upper right', ncol=2)
    ax_ts.set_xlim(0, 8)
    ax_ts.set_ylim(-100, ylim_ts)
    ax_ts.grid(True, alpha=0.15)

    # Peak bars
    np_ = min(4, len(pmean))
    x = np.arange(np_)
    w = 0.25
    ax_pk.bar(x - w, pmean[:np_], w, color=C_EXP, alpha=0.7, label='Exp.', zorder=3)
    ax_pk.errorbar(x - w, pmean[:np_], yerr=p2sig[:np_], fmt='none',
                   color='black', capsize=2, capthick=0.6, lw=0.6, zorder=4)
    for k, r in enumerate(sim.values()):
        pk = list(r['pk_p'][:np_]) + [0] * (np_ - len(r['pk_p'][:np_]))
        ax_pk.bar(x + k * w, pk[:np_], w, color=r['col'], alpha=0.8,
                  label=f"SPH {r['lab']}", zorder=3)
    ax_pk.set_xticks(x)
    ax_pk.set_xticklabels([f'P{i+1}' for i in range(np_)], fontsize=7)
    ax_pk.set_ylabel('Peak (Pa)')
    ax_pk.legend(fontsize=6, ncol=2, loc='upper right')
    ax_pk.grid(axis='y', alpha=0.15)


def _load_fluid_data(exp_ts_path, exp_pk_path, runs, min_peak_h=800):
    """Load experiment + simulation data for one fluid. Returns (te, pe, pmean, p2sig, ept, sim)."""
    te, pe = load_exp_timeseries(exp_ts_path)
    pstats = load_exp_peaks(exp_pk_path)
    pmean = pstats.mean(axis=0)
    p2sig = 2 * pstats.std(axis=0)
    ept, epp = find_peaks_filtered(te, pe, min_peak_h, 1.0)

    sim = {}
    for name, csv, pi, lab, col, lw in runs:
        if not csv.exists():
            continue
        t, p = load_dsph_csv(csv, probe=pi)
        spt, spp = find_peaks_filtered(t, p, min_peak_h * 0.8, 1.0, max_h=min_peak_h * 12)
        shift = 0.0
        n = min(3, len(ept), len(spt))
        if n > 0:
            shift = np.mean([ept[i] - spt[i] for i in range(n)])
        sim[name] = dict(t=t, p=p, lab=lab, col=col, lw=lw,
                         shift=shift, pk_t=spt, pk_p=spp)
    return te, pe, pmean, p2sig, ept, sim


SIMDATA_C = Path("/mnt/simdata/dualsphysics/exp-c")

def fig06_spheric_pressure():
    print("Fig 6: SPHERIC Pressure Validation (Water + Oil)")

    # Water
    water_runs = [
        ('run_001', SIMDATA/'run_001'/'PressLateral_dp004_Press.csv', 0, 'dp=4mm', C_SIM[0], 0.8),
        ('run_002', SIMDATA/'run_002'/'PressLateral_dp002_Press.csv', 3, 'dp=2mm', C_SIM[1], 1.1),
    ]
    w_ts = EXPDATA / "lateral_water_1x.txt"
    w_pk = EXPDATA / "Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt"

    # Oil: Artificial (run_005) vs Laminar+SPS (run_009b)
    OIL_009B = ROOT / "exp-c" / "analysis" / "run_009b" / "pressure_009b.csv"
    oil_runs = [
        ('run_005', SIMDATA/'run_005'/'PressConsistent_Press.csv', 0,
         'Artificial (α=0.1)', C_SIM[0], 0.8),
        ('run_009b', OIL_009B, 0,
         'Laminar+SPS (ν)', C_SIM[2], 1.1),
    ]
    o_ts = EXPDATA / "lateral_oil_1x.txt"
    o_pk = EXPDATA / "Oil_4first_peak_lateral_impact_tto_0_8_H93_B1X.txt"

    # Check required files
    missing = []
    for f in [w_ts, w_pk, o_ts, o_pk]:
        if not f.exists():
            missing.append(f.name)
    for _, csv, *_ in water_runs + oil_runs:
        if not csv.exists():
            missing.append(csv.name)
    if missing:
        print(f"  SKIP: {missing}")
        return

    # Load data
    w_te, w_pe, w_pmean, w_p2sig, w_ept, w_sim = _load_fluid_data(
        w_ts, w_pk, water_runs, min_peak_h=1000)
    o_te, o_pe, o_pmean, o_p2sig, o_ept, o_sim = _load_fluid_data(
        o_ts, o_pk, oil_runs, min_peak_h=300)

    # 2×2 figure
    fig, axes = plt.subplots(2, 2, figsize=(COL_DOUBLE, 5.0),
                              gridspec_kw={'width_ratios': [2.5, 1], 'hspace': 0.35, 'wspace': 0.30})

    _plot_spheric_row(axes[0], w_te, w_pe, w_pmean, w_p2sig, w_ept, w_sim, 'Water', 5000)
    axes[0, 0].set_title('(a) Water Lateral — Time Series', fontsize=9, loc='left')
    axes[0, 1].set_title('(b) Water — Peak Comparison', fontsize=9, loc='left')

    _plot_spheric_row(axes[1], o_te, o_pe, o_pmean, o_p2sig, o_ept, o_sim, 'Oil', 1500)
    axes[1, 0].set_title('(c) Oil Lateral — Time Series', fontsize=9, loc='left')
    axes[1, 1].set_title('(d) Oil — Peak Comparison', fontsize=9, loc='left')
    axes[1, 0].set_xlabel('Time (s)')

    fig.subplots_adjust(left=0.08, right=0.97, top=0.95, bottom=0.08,
                        hspace=0.35, wspace=0.30)
    save_fig(fig, 'fig06_spheric_pressure')


# ═══════════════════════════════════════════════════════════════
# Fig 7: Convergence Metrics (documented values)
# ═══════════════════════════════════════════════════════════════
def fig07_convergence():
    """Convergence metrics across 3 dp levels (documented values)."""
    print("Fig 7: Convergence Metrics")

    # Documented convergence results from SPHERIC T10 Water Lateral
    # Source: research-v3 experiment logs, spheric_pressure_comparison.py output
    dp_labels = ['dp=4 mm', 'dp=2 mm', 'dp=1 mm']
    dp_vals = [4, 2, 1]

    # Cross-correlation with experiment (documented)
    xcorr = [0.460, 0.655, 0.697]
    # Mean peak error % vs experiment (documented)
    peak_err = [28.5, 19.5, 15.2]

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(COL_DOUBLE, 2.5))

    x = np.arange(3)
    colors = C_SIM

    # (a) Cross-correlation
    bars1 = ax1.bar(x, xcorr, width=0.5, color=colors, edgecolor='white',
                    linewidth=0.5, zorder=3)
    for i, v in enumerate(xcorr):
        ax1.text(i, v + 0.015, f'{v:.3f}', ha='center', fontsize=8, fontweight='bold')

    ax1.set_xticks(x)
    ax1.set_xticklabels(dp_labels, fontsize=8)
    ax1.set_ylabel('Cross-correlation r')
    ax1.set_title('(a) Time-Series Correlation', fontsize=10, loc='left')
    ax1.set_ylim(0, 0.85)
    ax1.axhline(1.0, color='#999', ls=':', lw=0.5, alpha=0.3)
    ax1.grid(axis='y', alpha=0.15)

    # Convergence arrow
    ax1.annotate('', xy=(2, 0.697), xytext=(0, 0.460),
                 arrowprops=dict(arrowstyle='->', color='#1A9850',
                                lw=1.5, connectionstyle='arc3,rad=0.15'))
    ax1.text(1, 0.38, '+51% improvement', ha='center', fontsize=7,
             color='#1A9850', fontweight='bold')

    # (b) Peak error
    bars2 = ax2.bar(x, peak_err, width=0.5, color=colors, edgecolor='white',
                    linewidth=0.5, zorder=3)
    for i, v in enumerate(peak_err):
        ax2.text(i, v + 0.5, f'{v:.1f}%', ha='center', fontsize=8, fontweight='bold')

    # Experimental scatter band
    ax2.axhspan(0, 50, facecolor='#eee', alpha=0.4, zorder=0)
    ax2.text(2.3, 42, 'Exp. CoV\n(20-50%)', fontsize=6.5, color='#888',
             ha='center', va='top')

    ax2.set_xticks(x)
    ax2.set_xticklabels(dp_labels, fontsize=8)
    ax2.set_ylabel('Mean Peak Error (%)')
    ax2.set_title('(b) Peak Pressure Error', fontsize=10, loc='left')
    ax2.set_ylim(0, 38)
    ax2.grid(axis='y', alpha=0.15)

    fig.tight_layout()
    save_fig(fig, 'fig07_convergence_metrics')


# ═══════════════════════════════════════════════════════════════
# Fig 8: Bridge Physics Checks
# ═══════════════════════════════════════════════════════════════
def fig08_bridge_physics():
    print("Fig 8: Bridge Physics Checks")
    if not BRIDGE_CSV.exists():
        print(f"  SKIP: {BRIDGE_CSV}")
        return

    times, probes = [], {i: [] for i in range(5)}
    with open(BRIDGE_CSV) as f:
        lines = f.readlines()
        pos_x = [float(x) for x in lines[0].strip().split(';')[2:]]
        for line in lines[4:]:
            parts = line.strip().split(';')
            if len(parts) < 7:
                continue
            try:
                times.append(float(parts[1]))
                for i in range(5):
                    probes[i].append(float(parts[2 + i]))
            except ValueError:
                continue

    t = np.array(times)
    p_left = np.array(probes[0])
    p_right = np.array(probes[4])
    p_center = np.array(probes[2])

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(COL_DOUBLE, 3.0),
                                    gridspec_kw={'width_ratios': [2.2, 1]})

    # (a) Anti-phase time-series
    ax1.plot(t, p_left, '-', color=C_32B, lw=0.8, alpha=0.85,
             label=f'Left wall (x={pos_x[0]:.2f} m)')
    ax1.plot(t, p_right, '-', color=C_BAFFLE, lw=0.8, alpha=0.85,
             label=f'Right wall (x={pos_x[4]:.2f} m)')
    hydro = np.mean(p_center[len(p_center)//2:])
    ax1.axhline(hydro, color='#999', ls=':', lw=0.6, alpha=0.5, label='Hydrostatic ref.')

    ax1.set_xlabel('Time (s)')
    ax1.set_ylabel('Pressure (Pa)')
    ax1.set_title('(a) Anti-Phase Wall Pressure', fontsize=10, loc='left')
    ax1.legend(fontsize=6.5, loc='upper right')
    ax1.grid(True, alpha=0.15)

    # (b) Physics scorecard — clean table-like layout
    ax2.axis('off')
    ax2.set_xlim(0, 1)
    ax2.set_ylim(0, 1)

    checks = [
        ('M1 Hydrostatic',    '3.0% err'),
        ('M2 Anti-phase',     'r = -0.72'),
        ('M3 Nat. frequency', '4.4% err'),
        ('M4 Peak pressure',  'within 2s'),
        ('M5 Convergence',    'monotonic'),
    ]

    ax2.text(0.5, 0.97, '(b) Physics Checks', ha='center', va='top',
             fontsize=10, fontweight='bold')

    for i, (name, val) in enumerate(checks):
        y = 0.82 - i * 0.15
        # Green background rectangle
        ax2.add_patch(plt.Rectangle((0.02, y - 0.05), 0.96, 0.12,
                      facecolor='#E8F5E9', edgecolor='#1A9850',
                      linewidth=0.6, alpha=0.7, transform=ax2.transAxes))
        ax2.text(0.05, y, 'PASS', transform=ax2.transAxes, fontsize=6.5,
                 fontweight='bold', color='#1A9850', va='center',
                 family='sans-serif')
        ax2.text(0.22, y, name, transform=ax2.transAxes, fontsize=7,
                 va='center', color='#333')
        ax2.text(0.95, y, val, transform=ax2.transAxes, fontsize=7,
                 ha='right', va='center', color='#555', fontweight='bold')

    # Verdict
    ax2.text(0.5, 0.05, '5 / 5 PASS', ha='center', va='center',
             transform=ax2.transAxes, fontsize=12, fontweight='bold',
             color='#1A9850',
             bbox=dict(boxstyle='round,pad=0.3', fc='#E8F5E9', ec='#1A9850'))

    fig.tight_layout()
    save_fig(fig, 'fig08_bridge_physics')


# ═══════════════════════════════════════════════════════════════
# Fig 9: EXP-D Baffle
# ═══════════════════════════════════════════════════════════════
def load_swl_probes(csv_path):
    Z = np.arange(0.0, 0.40, 0.02)
    TH = 10.0
    t_out, z_out = [], []
    with open(csv_path) as f:
        for line in f:
            line = line.strip()
            if not line or line[0] in '#P ':
                continue
            parts = line.split(';')
            if len(parts) < 22:
                continue
            try:
                t = float(parts[1])
                pr = [float(parts[2 + i]) for i in range(20)]
                mz = 0.0
                for i in range(19, -1, -1):
                    if pr[i] > TH:
                        mz = Z[i]
                        if i < 19 and pr[i+1] <= TH:
                            mz = Z[i] + (Z[i+1]-Z[i]) * (pr[i]/(pr[i]+TH)) * 0.5
                        break
                t_out.append(t)
                z_out.append(mz)
            except (ValueError, IndexError):
                continue
    return np.array(t_out), np.array(z_out)


def fig09_expd_baffle():
    print("Fig 9: EXP-D Baffle Optimization")
    files = {
        'bl_left':  EXPD_DIR / 'baseline_swl_left_Press.csv',
        'bl_right': EXPD_DIR / 'baseline_swl_right_Press.csv',
        'bf_left':  EXPD_DIR / 'baffle_swl_left_Press.csv',
        'bf_right': EXPD_DIR / 'baffle_swl_right_Press.csv',
    }
    missing = [k for k, v in files.items() if not v.exists()]
    if missing:
        print(f"  SKIP: {missing}")
        return

    data = {k: load_swl_probes(v) for k, v in files.items()}

    fig, axes = plt.subplots(2, 1, figsize=(COL_DOUBLE, 3.5), sharex=True)
    for i, (side, label) in enumerate([('left', 'Left Wall'), ('right', 'Right Wall')]):
        ax = axes[i]
        tb, zb = data[f'bl_{side}']
        tf, zf = data[f'bf_{side}']

        ax.plot(tb, zb, '-', color=C_BASELINE, lw=0.8, alpha=0.85, label='Baseline')
        ax.plot(tf, zf, '-', color=C_BAFFLE, lw=0.8, alpha=0.85, label='With baffle')
        ax.axhline(0.2, color='#999', ls=':', lw=0.5, alpha=0.4)

        mask_b = tb >= 2.0
        mask_f = tf >= 2.0
        if mask_b.any() and mask_f.any():
            ab = (zb[mask_b].max() - zb[mask_b].min()) / 2
            af = (zf[mask_f].max() - zf[mask_f].min()) / 2
            red = (1 - af/ab) * 100 if ab > 0 else 0
            ax.text(0.97, 0.90, f'-{red:.0f}%', transform=ax.transAxes,
                    ha='right', va='top', fontsize=12, fontweight='bold',
                    color=C_BAFFLE,
                    bbox=dict(boxstyle='round,pad=0.2', fc='white', ec=C_BAFFLE, alpha=0.9))

        ax.set_ylabel('SWL (m)')
        ax.set_title(f'({chr(97+i)}) {label}', fontsize=10, loc='left')
        if i == 0:
            ax.legend(fontsize=7, loc='upper left')
        ax.grid(True, alpha=0.15)
        ax.set_ylim(0, 0.40)

    axes[1].set_xlabel('Time (s)')
    fig.tight_layout()
    save_fig(fig, 'fig09_expd_baffle')


# ═══════════════════════════════════════════════════════════════
FIGURES = {
    3: ('EXP-A Heatmap', fig03_expa_heatmap),
    4: ('v3->v4 Impact', fig04_v3v4_impact),
    5: ('EXP-B Ablation', fig05_expb_ablation),
    6: ('SPHERIC Pressure', fig06_spheric_pressure),
    7: ('Convergence', fig07_convergence),
    8: ('Bridge Physics', fig08_bridge_physics),
    9: ('EXP-D Baffle', fig09_expd_baffle),
}

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--only', nargs='+', type=int)
    parser.add_argument('--outdir', type=str, default=None)
    args = parser.parse_args()

    global OUT_DIR
    if args.outdir:
        OUT_DIR = Path(args.outdir)
    OUT_DIR.mkdir(parents=True, exist_ok=True)

    targets = args.only or sorted(FIGURES.keys())
    print(f"Output: {OUT_DIR}\nGenerating: {targets}\n")

    for n in targets:
        if n not in FIGURES:
            continue
        name, func = FIGURES[n]
        try:
            func()
        except Exception as e:
            print(f"  FAIL Fig {n} ({name}): {e}")
            import traceback; traceback.print_exc()

    print(f"\nDone. {OUT_DIR}/")

if __name__ == '__main__':
    main()
