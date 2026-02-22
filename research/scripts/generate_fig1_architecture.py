#!/usr/bin/env python3
"""
generate_fig1_architecture.py — Generate Fig 1: SloshAgent System Architecture Diagram

Publication-quality architecture diagram showing the full pipeline:
  Natural Language → LLM Agent → Tool Integration → Simulation → Validation

Output: research/experiments/figures/fig1_architecture.png (300 DPI)
"""

import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.patches import FancyBboxPatch, FancyArrowPatch
from pathlib import Path


# ── Color Palette ────────────────────────────────────────────────────

C = {
    "input":     "#3B7DD8",   # Blue
    "prompt":    "#5194D4",   # Light blue
    "llm":       "#5FAA46",   # Green
    "agent":     "#3D7A28",   # Dark green
    "tool_sph":  "#E87424",   # Orange
    "tool_pv":   "#D4A017",   # Gold/amber
    "exec_gpu":  "#C0392B",   # Red
    "exec_pv":   "#A0522D",   # Sienna
    "output":    "#7B2D8E",   # Purple
    "arrow":     "#404040",
    "border":    "#2C2C2C",
    "text":      "#1A1A1A",
    "sub":       "#555555",
    "white":     "#FFFFFF",
    "black":     "#000000",
    "bg":        "#FFFFFF",
    "loop_bg":   "#F0F4E8",   # Very faint green for agent loop backdrop
}

plt.rcParams.update({
    "font.family": "sans-serif",
    "font.sans-serif": ["DejaVu Sans", "Helvetica", "Arial"],
    "font.size": 9,
    "axes.linewidth": 0,
    "figure.facecolor": C["bg"],
})


# ── Drawing Primitives ───────────────────────────────────────────────

def box(ax, x, y, w, h, label, color, fc="#FFF", fs=10, fw="bold",
        sub=None, sub_fs=7, rad=0.018, alpha=0.95, bw=1.2):
    """Rounded box with optional subtitle."""
    p = FancyBboxPatch(
        (x, y), w, h, boxstyle=f"round,pad=0,rounding_size={rad}",
        facecolor=color, edgecolor=C["border"], linewidth=bw,
        alpha=alpha, zorder=2,
    )
    ax.add_patch(p)
    if sub:
        ax.text(x + w/2, y + h*0.63, label, ha="center", va="center",
                fontsize=fs, fontweight=fw, color=fc, zorder=3)
        ax.text(x + w/2, y + h*0.28, sub, ha="center", va="center",
                fontsize=sub_fs, fontstyle="italic", color=fc, alpha=0.88, zorder=3)
    else:
        ax.text(x + w/2, y + h/2, label, ha="center", va="center",
                fontsize=fs, fontweight=fw, color=fc, zorder=3)


def arrow(ax, x1, y1, x2, y2, lw=1.5, color=None, cs="arc3,rad=0", ms=14):
    """Curved arrow."""
    a = FancyArrowPatch(
        (x1, y1), (x2, y2), arrowstyle="-|>",
        color=color or C["arrow"], linewidth=lw,
        connectionstyle=cs, mutation_scale=ms, zorder=4,
    )
    ax.add_patch(a)


def tool_box(ax, x, y, w, h, title, tools, bg, fc, ts=9.5, tfs=7):
    """Tool group with monospace tool listing."""
    p = FancyBboxPatch(
        (x, y), w, h, boxstyle="round,pad=0,rounding_size=0.012",
        facecolor=bg, edgecolor=C["border"], linewidth=1.0,
        alpha=0.92, zorder=2,
    )
    ax.add_patch(p)
    ax.text(x + w/2, y + h - 0.018, title, ha="center", va="top",
            fontsize=ts, fontweight="bold", color=fc, zorder=3)
    ax.text(x + w/2, y + h*0.38, "\n".join(tools), ha="center", va="center",
            fontsize=tfs, fontfamily="monospace", color=fc, alpha=0.92,
            linespacing=1.3, zorder=3)


# ── Main ─────────────────────────────────────────────────────────────

def main():
    fig, ax = plt.subplots(figsize=(9, 11.5))
    ax.set_xlim(0, 1)
    ax.set_ylim(0, 1)
    ax.set_aspect("auto")
    ax.axis("off")

    # Layout
    CX = 0.50
    W = 0.40            # Standard box width
    H = 0.050           # Standard box height
    TW = 0.30           # Tool box width
    TH = 0.135          # Tool box height

    # Y positions (top → bottom)
    Y = {
        "input":  0.890,
        "prompt": 0.825,
        "llm":    0.760,
        "agent":  0.693,
        "tools":  0.525,
        "exec":   0.370,
        "out":    0.283,
        "valid":  0.175,
    }

    # ── Title ────────────────────────────────────────────────────
    ax.text(CX, 0.990, "SloshAgent System Architecture",
            ha="center", va="top", fontsize=16, fontweight="bold", color=C["text"])
    ax.text(CX, 0.968, "End-to-end pipeline: Natural Language \u2192 Validated CFD Results",
            ha="center", va="top", fontsize=9.5, color=C["sub"], fontstyle="italic")

    # ── 1  Natural Language Input ────────────────────────────────
    box(ax, CX-W/2, Y["input"], W, H,
        "Natural Language Input", C["input"],
        fs=12, sub='"600mm tank, 30% fill, sway at 0.9f\u2081"', sub_fs=8.5)
    arrow(ax, CX, Y["input"], CX, Y["prompt"]+H)

    # ── 2  SloshingCoderPrompt ───────────────────────────────────
    box(ax, CX-W/2, Y["prompt"], W, H,
        "SloshingCoderPrompt", C["prompt"],
        fs=11.5, sub="136 lines  \u2022  5 domain categories", sub_fs=8)
    arrow(ax, CX, Y["prompt"], CX, Y["llm"]+H)

    # ── 3  Qwen3 LLM ────────────────────────────────────────────
    box(ax, CX-W/2, Y["llm"], W, H,
        "Qwen3 32B / 8B", C["llm"],
        fs=12, sub="Local Ollama  \u2022  Zero API Cost", sub_fs=8)
    arrow(ax, CX, Y["llm"], CX, Y["agent"]+H)

    # ── 4  ReAct Agent Loop ──────────────────────────────────────
    AW = 0.48
    box(ax, CX-AW/2, Y["agent"], AW, H,
        "ReAct Agent Loop", C["agent"],
        fs=12, sub="Tool Use \u2192 Observe \u2192 Reason \u2192 Act", sub_fs=8.5)

    # Feedback loop (right side, dashed)
    rx = CX + AW/2 + 0.025
    ax.annotate("", xy=(rx, Y["agent"]+H/2),
                xytext=(rx, Y["tools"]+TH/2),
                arrowprops=dict(arrowstyle="-|>", color=C["agent"],
                                linewidth=1.3, linestyle="--"))
    ax.text(rx + 0.018, (Y["agent"] + Y["tools"]+TH)/2,
            "feedback", rotation=90, ha="left", va="center",
            fontsize=7.5, color=C["agent"], fontstyle="italic")

    # Arrows: agent → tools (diverge)
    sph_cx = CX - TW/2 - 0.025 + TW/2   # left tool center x
    pv_cx  = CX + 0.025 + TW/2           # right tool center x
    arrow(ax, CX-0.06, Y["agent"], sph_cx, Y["tools"]+TH, cs="arc3,rad=0.08")
    arrow(ax, CX+0.06, Y["agent"], pv_cx, Y["tools"]+TH, cs="arc3,rad=-0.08")

    # ── 5  Tool Groups ──────────────────────────────────────────
    # Left: DualSPHysics Tools
    sx = CX - TW - 0.025
    tool_box(ax, sx, Y["tools"], TW, TH,
             "14 DualSPHysics Tools",
             ["xml_generator  gencase",
              "solver         partvtk",
              "measuretool    monitor",
              "error_recovery job_manager",
              "analysis       parametric_study",
              "seismic_input  stl_import",
              "result_store   report"],
             C["tool_sph"], C["white"])

    # Right: pv-agent MCP Tools
    px = CX + 0.025
    tool_box(ax, px, Y["tools"], TW, TH,
             "12 pv-agent MCP Tools",
             ["render       animate",
              "slice        clip",
              "contour      pv_isosurface",
              "inspect_data extract_stats",
              "plot_over_line",
              "streamlines"],
             C["tool_pv"], C["black"], ts=9.5, tfs=7)

    # ── 6  Execution Layer ───────────────────────────────────────
    EW = 0.27
    EH = 0.048
    gx = sx + (TW - EW)/2
    pvx = px + (TW - EW)/2

    box(ax, gx, Y["exec"], EW, EH,
        "DualSPHysics Docker", C["exec_gpu"],
        fs=9.5, sub="CUDA 12.6  \u2022  RTX 4090", sub_fs=7.5)
    box(ax, pvx, Y["exec"], EW, EH,
        "ParaView pvpython", C["exec_pv"],
        fs=9.5, sub="Mesa Headless Rendering", sub_fs=7.5)

    arrow(ax, sx+TW/2, Y["tools"], gx+EW/2, Y["exec"]+EH)
    arrow(ax, px+TW/2, Y["tools"], pvx+EW/2, Y["exec"]+EH)

    # ── 7  Outputs ───────────────────────────────────────────────
    OW = 0.24
    OH = 0.042
    olx = gx + (EW - OW)/2
    orx = pvx + (EW - OW)/2

    box(ax, olx, Y["out"], OW, OH,
        "GPU SPH Simulation", C["exec_gpu"], alpha=0.78,
        fs=9, sub="Particles + Pressure", sub_fs=6.5)
    box(ax, orx, Y["out"], OW, OH,
        "Post-processing", C["exec_pv"], alpha=0.78,
        fs=9, sub="VTK + MP4 Animation", sub_fs=6.5)

    arrow(ax, gx+EW/2, Y["exec"], olx+OW/2, Y["out"]+OH, lw=1.2)
    arrow(ax, pvx+EW/2, Y["exec"], orx+OW/2, Y["out"]+OH, lw=1.2)

    # ── 8  Validation ────────────────────────────────────────────
    VW = 0.58
    VH = 0.055
    box(ax, CX-VW/2, Y["valid"], VW, VH,
        "Experimental Validation & Analysis", C["output"],
        fs=11.5, sub="SPHERIC Test 10  \u2022  Chen2018 Parametric  \u2022  AI Report", sub_fs=8)

    arrow(ax, olx+OW/2, Y["out"], CX-0.05, Y["valid"]+VH, cs="arc3,rad=0.10", lw=1.2)
    arrow(ax, orx+OW/2, Y["out"], CX+0.05, Y["valid"]+VH, cs="arc3,rad=-0.10", lw=1.2)

    # ── Layer circled numbers (left margin) ──────────────────────
    circled = ["\u2460", "\u2461", "\u2462", "\u2463", "\u2464", "\u2465", "\u2466"]
    ys = [Y["input"]+H/2, Y["prompt"]+H/2, Y["llm"]+H/2, Y["agent"]+H/2,
          Y["tools"]+TH/2, Y["exec"]+EH/2, Y["valid"]+VH/2]
    for num, yy in zip(circled, ys):
        ax.text(0.045, yy, num, ha="center", va="center",
                fontsize=12, color=C["sub"], fontweight="bold")

    # ── Legend ────────────────────────────────────────────────────
    items = [
        ("User Input", C["input"]),
        ("LLM / Agent", C["llm"]),
        ("Tool Interface", C["tool_sph"]),
        ("Execution", C["exec_gpu"]),
        ("Output", C["output"]),
    ]
    ly = 0.095
    lx0 = 0.115
    sp = 0.162
    ax.plot([0.08, 0.92], [ly+0.035, ly+0.035], color="#CCCCCC", lw=0.5, zorder=1)
    for i, (label, color) in enumerate(items):
        lx = lx0 + i * sp
        p = FancyBboxPatch(
            (lx, ly), 0.024, 0.018,
            boxstyle="round,pad=0,rounding_size=0.004",
            facecolor=color, edgecolor=C["border"], linewidth=0.8, zorder=2,
        )
        ax.add_patch(p)
        ax.text(lx + 0.032, ly + 0.009, label, ha="left", va="center",
                fontsize=8, color=C["text"], zorder=3)

    # ── Save ─────────────────────────────────────────────────────
    plt.tight_layout(pad=0.3)
    out = Path("research/experiments/figures/fig1_architecture.png")
    out.parent.mkdir(parents=True, exist_ok=True)
    plt.savefig(out, dpi=300, bbox_inches="tight",
                facecolor=C["bg"], edgecolor="none")
    plt.close()
    print(f"Saved: {out}  ({out.stat().st_size / 1024:.0f} KB)")


if __name__ == "__main__":
    main()
