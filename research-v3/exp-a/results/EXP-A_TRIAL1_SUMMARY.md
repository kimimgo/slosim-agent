# EXP-A Trial 1 Complete Summary

**Date**: 2026-03-02
**Infrastructure**: pajulab (4× RTX 4090), Ollama
**Models**: Qwen3-32B (20GB), Qwen3-8B/latest (5.2GB)

## M-A3 Parameter Fidelity Results

| Scenario | Tier | 32B Score | 8B Score | Delta | Key Failures |
|----------|------|-----------|----------|-------|--------------|
| S01 | Easy | 64% (7/11) | 64% (7/11) | 0% | density, viscosity, motion_type, amplitude |
| S02 | Easy | **100%** (9/9) | **100%** (9/9) | 0% | None |
| S03 | Medium | 89% (8/9) | 89% (8/9) | 0% | timemax |
| S04 | Medium | 78% (7/9) | 78% (7/9) | 0% | motion_type, amplitude |
| S05 | Medium | 56% (5/9) | 56% (5/9) | 0% | tank_z, fill_height, motion_type, amplitude |
| S06 | Medium | 89% (8/9) | 89% (8/9) | 0% | timemax |
| S07 | Hard | 57% (4/7) | 57% (4/7) | 0% | geometry, timemax, dp |
| S08 | Easy | 60% (6/10) | N/A (no XML) | — | motion_type, freq, ampl, timemax (32B); tool fail (8B) |
| S09 | Hard | 33% (2/6) | 33% (2/6) | 0% | geometry, motion_type, amplitude, timemax |
| S10 | Hard | N/A | N/A | — | STL file not provided |

### Aggregate Scores

| | 32B | 8B |
|---------|------|------|
| **Overall Mean** | **69.5%** (9 scenarios) | **70.7%** (8 scenarios) |
| Easy | 74.5% (3) | 81.8% (2) |
| Medium | 77.8% (4) | 77.8% (4) |
| Hard | 45.2% (2) | 45.2% (2) |

**8B aggregate is higher only because S08 (60%) is excluded due to tool calling failure.**
**When comparing only overlapping 7 scenarios: 32B = 8B = identical scores.**

## Execution Time Comparison

| Scenario | 32B (s) | 8B (s) | Ratio |
|----------|---------|--------|-------|
| S01 | 157 | 121 | 0.77× |
| S02 | 99 | 82 | 0.83× |
| S03 | 109 | 373 | 3.42× |
| S04 | 120 | 69 | 0.58× |
| S05 | 150 | 66 | 0.44× |
| S06 | 135 | 62 | 0.46× |
| S07 | 151 | 164 | 1.09× |
| S08 | 180 | 105 | 0.58× |
| S09 | 153 | 433 | 2.83× |
| S10 | 213 | 44 | 0.21× |
| **Mean** | **147** | **152** | ~1.0× |

**Note**: 8B is faster on simple scenarios but MUCH slower on Hard scenarios (S03: 3.4×, S09: 2.8×), likely due to more tool call retries.

## Key Findings

### F1: Model Size Has No Impact on Parameter Fidelity
32B and 8B produce **identical M-A3 scores** across all comparable scenarios. The bottleneck is the tool (xml_generator) template design, not the LLM's reasoning capability.

### F2: Systematic Error Patterns (Model-Independent)
1. **P1 Pitch Rotation Confusion** — Agent always generates `mvrectsinu` (translation) instead of `mvrotsinu` (rotation). Affects S01, S04, S05, S08, S09 (5/10 scenarios).
2. **P2 TimeMax Underestimation** — Agent calculates ~5 periods instead of using the specified simulation time. Affects S03, S06, S07, S09.
3. **P3 Cylindrical Geometry Not Supported** — xml_generator only produces rectangular drawbox. Affects S07, S09.
4. **P4 Non-Water Fluid Properties Ignored** — Agent uses default rhop0=1000 and Visco=0.01 regardless of specified fluid. Affects S01.
5. **P5 Degree-to-Meter Amplitude Confusion** — When rotation is needed, agent converts degrees to a small linear displacement. Co-occurs with P1.

### F3: Tier-Dependent Difficulty Curve
- Easy: 74.5% — basic parameter extraction works well for horizontal sway scenarios
- Medium: 77.8% — comparable, but timemax and rotation are failure modes
- Hard: 45.2% — cylinder geometry and complex scenarios collapse below 50%

### F4: 8B Tool Calling Fragility
S08 8B failed to generate XML entirely. The 8B model sometimes outputs tool call JSON in response text instead of executing actual tool calls. This is a significant finding for production deployment.

## Implications for Paper

1. **Primary claim shift**: "NL→XML parameter extraction achieves ~70% fidelity, with model size being irrelevant; tool design is the bottleneck"
2. **Ablation study (EXP-B) predictions**: Removing system prompt should primarily affect parameter interpretation, not tool mechanics
3. **S10 should be excluded or redesigned** (structural limitation, not model failure)
4. **S08 needs separate analysis** for 8B (tool calling failure, not parameter fidelity)

## Remaining Work

- [ ] Trial 2-3 for 32B and 8B (statistical variation)
- [ ] EXP-B ablation (remove system prompt, remove post-processing)
- [ ] S10 redesign or exclusion decision
- [ ] M-A1 (tool sequence), M-A2 (XML validity), M-A4/A5 (solver metrics) analysis
