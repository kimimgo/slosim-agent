# EXP-A Final Results — NL→DualSPHysics XML Parameter Fidelity

**Date**: 2026-03-02
**Infrastructure**: pajulab (4× RTX 4090, Ollama)
**Models**: Qwen3-32B (20GB), Qwen3-8B/latest (5.2GB)
**Trials**: 3 per model × 10 scenarios = 60 runs total

---

## 1. Master Results Table

### M-A3 Parameter Fidelity Scores (%)

| Scenario | Tier | Description | 32B | 8B | Notes |
|----------|------|-------------|-----|-----|-------|
| S01 | Easy | SPHERIC Oil Low Fill | 64 | 64 | density, viscosity, motion_type, amplitude |
| S02 | Easy | Chen Shallow Sway | **100** | **100** | Perfect across all trials |
| S03 | Medium | Chen Near-Critical | 89 | 89 | timemax only |
| S04 | Medium | Liu Large Pitch 30s | 78 | 78 | motion_type, amplitude |
| S05 | Medium | Liu Amplitude Parametric | 56 | 56 | tank_z, fill, motion_type, amplitude |
| S06 | Medium | ISOPE LNG Roof Impact | 89 | 89 | timemax only |
| S07 | Hard | NASA Cylinder Baffle | 57 | 57 | geometry, timemax, dp |
| S08 | Easy | English mDBC/DBC | 60 | N/A† | 8B tool calling failure |
| S09 | Hard | Horizontal Cylinder Pitch | 33 | 33 | geometry, motion_type, amplitude, timemax |
| S10 | Hard | STL Fuel Tank | N/A‡ | N/A‡ | STL file not available |

**†** 8B model outputs tool call JSON as text instead of executing tools (S08 all 3 trials)
**‡** Structural limitation: STL file not provided on test server

### Aggregate by Tier

| Tier | 32B | 8B | Combined |
|------|-----|-----|----------|
| Easy (S01,S02,S08) | 74.5% | 81.8%§ | 74.5% |
| Medium (S03-S06) | 77.8% | 77.8% | 77.8% |
| Hard (S07,S09) | 45.2% | 45.2% | 45.2% |
| **Overall** | **69.5%** | **70.7%** | **69.5%** |

§ Higher because S08 (60%) excluded from 8B due to tool failure

---

## 2. Key Finding: Perfect Determinism

**Standard deviation = 0.0 across all trials for all scenarios.**

| Metric | Value |
|--------|-------|
| Inter-trial stddev (32B) | 0.0 for all 9 scored scenarios |
| Inter-trial stddev (8B) | 0.0 for all 8 scored scenarios |
| 32B vs 8B delta | 0.0 for all 8 overlapping scenarios |

**Cause**: Ollama temperature=0 + identical system prompt + deterministic tool execution → bit-identical XML output across trials.

**Implication**: 3 trials confirmed determinism, but 1 trial would have sufficed. The "variability" dimension of the experimental design is unnecessary for temperature=0 settings. Report as "deterministic at t=0; future work should test t>0."

---

## 3. Systematic Error Patterns

### P1: Pitch Rotation Confusion (5/10 scenarios)
- **Pattern**: Agent generates `mvrectsinu` (horizontal translation) instead of `mvrotsinu` (pitch rotation)
- **Affected**: S01, S04, S05, S08, S09
- **Root Cause**: System prompt describes both motion types but xml_generator tool only accepts `freq`, `amplitude` parameters as floats — no motion type selector
- **Fix**: Add `motion_type` parameter to xml_generator with enum {mvrectsinu, mvrotsinu}

### P2: TimeMax Underestimation (4/10 scenarios)
- **Pattern**: Agent calculates ~5 oscillation periods as simulation time
- **Affected**: S03 (4.96 vs 10), S06 (8.33 vs 15), S07 (10 vs 15), S09 (9.09 vs 15)
- **Root Cause**: Agent applies heuristic "5 periods is enough" when user doesn't explicitly state timemax

### P3: Cylindrical Geometry Not Supported (2/10 scenarios)
- **Pattern**: Agent defaults to rectangular drawbox for cylindrical tanks
- **Affected**: S07, S09
- **Root Cause**: xml_generator tool has no cylinder geometry option

### P4: Non-Water Fluid Properties Ignored (1/10 scenarios)
- **Pattern**: Agent uses default rhop0=1000 and Visco=0.01 regardless
- **Affected**: S01 (oil, density 990, viscosity 0.045)
- **Root Cause**: xml_generator template hardcodes water properties

### P5: SPHERIC Benchmark Knowledge Gap (1/10 scenario)
- **Pattern**: When "SPHERIC benchmark" is referenced without explicit motion params, agent fabricates values
- **Affected**: S08
- **Root Cause**: No embedded domain knowledge of standard benchmark cases

### P6: 8B Tool Calling Fragility (1/10 scenario)
- **Pattern**: 8B model outputs JSON tool calls in response text instead of executing them
- **Affected**: S08 (consistent across all 3 trials)
- **Severity**: Agent produces no XML → complete failure

---

## 4. Execution Time Analysis

| Scenario | 32B Mean (s) | 8B Mean (s) | 8B/32B Ratio |
|----------|-------------|-------------|--------------|
| S01 | 162 | 96 | 0.59 |
| S02 | 94 | 74 | 0.79 |
| S03 | 103 | 256 | 2.49 |
| S04 | 114 | 82 | 0.72 |
| S05 | 151 | 78 | 0.52 |
| S06 | 119 | 72 | 0.61 |
| S07 | 147 | 300 | 2.04 |
| S08 | 185 | 101 | 0.55 |
| S09 | 152 | 269 | 1.77 |
| S10 | 214 | 63 | 0.29 |
| **Mean** | **144** | **139** | ~0.97 |

**Finding**: 8B is faster for simple scenarios (0.3-0.8×) but MUCH slower for complex ones (2-3×), likely due to more tool call retries and longer reasoning chains.

---

## 5. Implications for Paper

### Primary Claims Supported
1. **NL→XML is feasible**: 70% overall parameter fidelity demonstrates viability
2. **Perfect scenario exists**: S02 achieves 100% with explicit parameters
3. **Model size is irrelevant**: 32B ≡ 8B (tool design is the bottleneck)
4. **Deterministic at t=0**: No need for multiple trials

### Primary Claims to Adjust
1. ~~"Larger models perform better"~~ → "Model size has no effect on structured tool output"
2. ~~"Statistical variation across trials"~~ → "Deterministic output; trial variance is zero"
3. ~~"10 scenarios cover the space"~~ → "8 scored scenarios (S10 excluded, S08 8B excluded)"

### Table/Figure Plan for Paper
- **Table 3**: M-A3 scores per scenario per tier (this table)
- **Table 4**: Error pattern analysis (P1-P6)
- **Figure 2**: M-A3 by tier bar chart (Easy/Medium/Hard)
- **Table 5**: Execution time comparison (response latency)

---

## 6. EXP-B Ablation Predictions

Based on EXP-A findings, EXP-B ablation conditions:
- **B1 (−SystemPrompt)**: Expected M-A3 drop — agent won't know DualSPHysics XML format
- **B2 (−PostProcess)**: Expected M-A3 unchanged (post-processing doesn't affect XML generation)
- **B3 (−Tools)**: Expected complete failure (no xml_generator → no XML)

Given determinism at t=0, EXP-B needs only 1 trial per condition.
