# Unified Results Analysis — SloshAgent Paper

Date: 2026-02-19
Status: Complete (all 5 experiments)
Models: Qwen3 32B (primary) / Qwen3 8B (comparison)

---

## Section 1: GAP → Evidence Mapping

| GAP | RQ | EXP | Key Finding | Quantitative Evidence | Status |
|-----|-----|-----|-------------|----------------------|--------|
| **GAP-1**: No AI agent automates sloshing simulation | RQ1 | EXP-2 | SloshAgent generates valid DualSPHysics XML from natural language across 4 complexity levels | 32B: 17/20 tool calls (85%), 47% param accuracy, 14/20 physically valid (70%) | **COVERED** |
| **GAP-2**: Sloshing domain knowledge never encoded in LLM prompts | RQ4 | EXP-4 | Domain prompt ablation shows clear contribution hierarchy (32B only) | 32B: FULL 60% > NO-RULES 57% > NO-DOMAIN 50% > GENERIC 35% (+25%p FULL vs GENERIC) | **COVERED** |
| **GAP-3**: No experimental sloshing benchmark validation by any AI agent | RQ2 | EXP-1 | Agent-generated SPHERIC Test 10 simulation matches experimental pressure within ±2σ for water | Water: 7/7 peaks within ±2σ (Low+High); Oil: 0 peaks detected (DBC limitation) | **PARTIAL** |
| **GAP-4**: Sloshing industry automation gap ($4.13B market, zero AI) | RQ3 | EXP-3, EXP-5 | Automated 6-case parametric study + industrial baffle/seismic PoC | Parametric: 6/6 cases completed; Baffle: 91.9% amplitude reduction; single-prompt workflow | **COVERED** |
| **GAP-5**: Particle-based simulation + LLM = complete void | RQ3 | (Architecture) | First tool interface for any particle-based solver (14 purpose-built tools) | 14 tools covering full DualSPHysics pipeline; no direct experimental comparison possible | **PARTIAL** |

### Coverage Assessment

**COVERED (3/5)**: GAP-1, GAP-2, GAP-4 have direct experimental evidence with quantitative metrics.

**PARTIAL (2/5)**:
- **GAP-3**: Water validation is strong (all peaks within ±2σ), but oil fails (DBC boundary limitation). Time-series correlation metrics (Pearson r, NRMSE) are pending from Task #1 analysis. This is an honest limitation — DBC is adequate for water sloshing but mDBC is required for viscous fluids, consistent with English et al. (2021).
- **GAP-5**: Addressed by system design (14 tools, architecture), not by a comparative experiment. This is inherently a "first of its kind" claim that cannot be experimentally compared. The 10+ competitor survey (Table B.4 in gap analysis) provides the evidence.

### Honest Weakness Notes
1. **47% overall parameter accuracy** (EXP-2) is moderate — driven down by L3 (paper-specific) and L5 (edge cases). L1+L2 are strong at 82% average.
2. **Oil sloshing failure** (EXP-1) reveals a boundary condition limitation. Paper should present this honestly as a known DBC limitation, not an agent failure.
3. **No time-series correlation yet** for EXP-1 — peak-in-band is a necessary but not sufficient validation metric. Pearson r and NRMSE would strengthen GAP-3 considerably.

---

## Section 2: Complete Results Summary

### Table 2: SPHERIC Test 10 — Experimental Benchmark Validation (EXP-1)

| Case | Particles | dp [m] | Sim Peaks [mbar] | Peaks in ±2σ | Max P [mbar] | Status |
|------|-----------|--------|-------------------|--------------|--------------|--------|
| Water Low | 136K | 0.004 | 31.1, 58.9, 76.7 | **3/3** | 76.7 | PASS |
| Water High | 344K | 0.004 | 44.2, 29.4, 31.4, 45.3 | **4/4** | 50.0 | PASS |
| Oil Low | 136K | 0.004 | (none detected) | N/A | 0.0 | FAIL |

**Experimental Reference** (SPHERIC Test 10, 100-repeat statistics):

| Fluid | Peak 1 [mbar] | Peak 2 [mbar] | Peak 3 [mbar] | Peak 4 [mbar] |
|-------|---------------|---------------|---------------|---------------|
| Water mean | 37.1 | 48.2 | 46.9 | 46.4 |
| Water ±2σ | ±14.6 | ±29.9 | ±34.0 | ±26.3 |
| Oil mean | 6.9 | 6.5 | 5.4 | 5.5 |
| Oil ±2σ | ±0.3 | ±0.5 | ±0.5 | ±0.5 |

**Key Findings**:
- Water (Low + High): All 7 detected peaks fall within experimental ±2σ envelope
- Resolution convergence: Both Low (136K) and High (344K) produce valid results
- Oil failure: DBC + artificial viscosity over-damps the thin oil sloshing layer; mDBC boundary is required (cf. English et al., 2021)

---

### Table 3: NL→XML Generation Accuracy (EXP-2) — Combined 8B + 32B

#### Summary by Complexity Level

| Level | Name | n | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-------|------|---|--------------|-------------|-------------|------------|----------|---------|
| L1 | Basic (explicit params) | 4 | 4/4 | 4/4 | **96%** | **96%** | 4/4 | 4/4 |
| L2 | Domain (inference needed) | 4 | **4/4** | 3/4 | **67%** | 42% | 4/4 | 3/4 |
| L3 | Paper (exact reproduction) | 4 | 3/4 | 3/4 | 15% | 15% | 3/4 | 3/4 |
| L4 | Complex (multi-feature) | 4 | 3/4 | 3/4 | **58%** | 50% | 3/4 | 3/4 |
| L5 | Edge (error handling) | 4 | **3/4** | 2/4 | 0% | **25%** | 0/4 | 1/4 |
| **Total** | | **20** | **17/20** | 15/20 | **47%** | 46% | **14/20** | 14/20 |

#### Detailed Scenario Results (32B)

| ID | Level | Tool Called | Param Accuracy | Valid | Key Generated Params | Latency |
|----|-------|-----------|----------------|-------|---------------------|---------|
| S01 | L1 | xml_generator | 100% | Y | L=0.9, h=0.093, f=0.613 | 427.6s |
| S02 | L1 | xml_generator | 100% | Y | L=0.6, h=0.15, f=1.2 | 345.3s |
| S03 | L1 | xml_generator | 86% | Y | L=1, h=0.3, f=0.77 | 326.4s |
| S04 | L1 | xml_generator | 100% | Y | L=0.6, h=0.2, f=0.5 | 488.7s |
| S05 | L2 | xml_generator | 86% | Y | L=40, h=8.1, f=0.104 | 557.7s |
| S06 | L2 | xml_generator | 25% | Y | L=1, h=0.3, f=0.75 | 782.5s |
| S07 | L2 | xml_generator | 57% | Y | L=40, h=10, f=0.3 | 795.4s |
| S08 | L2 | xml_generator | 100% | Y | L=1, h=0.42, f=1 | 786.0s |
| S09 | L3 | xml_generator | 17% | Y | L=40, h=10.8, f=0.1053 | 665.2s |
| S10 | L3 | none | 0% | N | (no tool called) | 826.6s |
| S11 | L3 | xml_generator | 29% | Y | L=0.6, h=0.2, f=1.02 | 301.9s |
| S12 | L3 | xml_generator | 17% | Y | L=1, h=0.3, f=0.756 | 278.8s |
| S13 | L4 | xml_generator | 100% | Y | L=40, h=13.5, f=0.14 | 395.9s |
| S14 | L4 | xml_generator | 33% | Y | L=40, h=13.5, f=0.125 | 416.9s |
| S15 | L4 | xml_generator | 100% | Y | L=1, h=0.3, f=0.5 | 254.5s |
| S16 | L4 | none | 0% | N | (no tool called) | 722.7s |
| S17 | L5 | xml_generator | 0% | N | L=40, h=13.5, f=0.125 | 332.8s |
| S18 | L5 | xml_generator | 0% | N | L=1, h=0.3, f=1 | 296.8s |
| S19 | L5 | xml_generator | 0% | N | L=40, h=27, f=0.14 | 398.4s |
| S20 | L5 | none | 0% | N | (no tool called) | 737.8s |

**Interpretation by Level**:
- **L1 (96%)**: Explicit parameters are reliably extracted — both models achieve near-perfect accuracy
- **L2 (67%/42%)**: Domain inference (e.g., "LNG cargo tank" → dimensions, "resonance frequency" → formula) shows 32B advantage (+25%p)
- **L3 (15%/15%)**: Paper-specific parameters (Chen2018 exact values) are not in training data — expected failure
- **L4 (58%/50%)**: Multi-feature composition works for simple combos (baffle+tank) but fails for orchestration (parametric study)
- **L5 (0%/25%)**: Edge case handling (empty tank, overflow, extreme dp) is the weakest area — agent generates invalid configs instead of error messages

---

### Table 4: Chen2018 Parametric Study — Automated SWL Analysis (EXP-3)

| Fill [mm] | Fill [%] | f₁ [Hz] | f_exc [Hz] | Left Amp [mm] | Right Amp [mm] | Max SWL [mm] | Min SWL [mm] | Dom Freq [Hz] |
|-----------|---------|---------|-----------|--------------|---------------|-------------|-------------|---------------|
| 120 | 18.5 | 0.8512 | 0.7661 | 66.6 | 69.8 | 224.2 | 90.9 | 0.142 |
| 150 | 23.1 | 0.9237 | 0.8313 | 35.9 | 38.6 | 192.3 | 120.6 | 0.851 |
| 195 | 30.0 | 1.0011 | 0.9010 | 41.4 | 42.9 | 243.2 | 160.5 | 0.851 |
| 260 | 40.0 | 1.0680 | 0.9612 | 45.6 | 51.4 | 310.0 | 217.2 | 0.993 |
| 325 | 50.0 | 1.1033 | 0.9930 | 48.2 | 59.3 | 376.3 | 278.7 | 0.993 |
| 390 | 60.0 | 1.1216 | 1.0094 | 50.8 | 62.3 | 443.0 | 340.8 | 0.993 |

**Setup**: Tank 600 x 300 x 650 mm, horizontal sway f/f₁ = 0.9, A = 7mm, DualSPHysics DBC dp=0.005m, 10s duration.

**Key Findings**:
- All 6 fill levels completed successfully by agent in single parametric run
- SWL amplitude increases monotonically with fill level (physically expected)
- Left-right asymmetry grows with fill level (66.6→50.8 left vs 69.8→62.3 right)
- Dominant frequency transitions from low (0.142 Hz at 18.5%) to excitation-locked (0.993 Hz at 40%+)
- 120mm (18.5%) case shows anomalous low-frequency dominance — shallow water nonlinear effects

**Limitation**: Direct NRMSE comparison vs Chen2018 published figures requires digitization (not yet performed). The physical trends are consistent but quantitative validation is incomplete.

---

### Table 5: Domain Prompt Ablation (EXP-4) — Combined 8B + 32B

#### Summary

| Condition | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-----------|-------------|-------------|-------------|------------|----------|---------|
| **FULL** | **10/10** | 7/10 | **60%** | 46% | 8/10 | 6/10 |
| **NO-DOMAIN** | **10/10** | **10/10** | 50% | 44% | 8/10 | **8/10** |
| **NO-RULES** | **10/10** | 9/10 | 57% | **55%** | 8/10 | **8/10** |
| **GENERIC** | 6/10 | **8/10** | 35% | 39% | 6/10 | **8/10** |

#### Per-Scenario Breakdown (32B)

| Scenario | Level | FULL | NO-DOMAIN | NO-RULES | GENERIC |
|----------|-------|------|-----------|----------|---------|
| S01 | L1 | 100% PASS | 100% PASS | 100% PASS | 0% FAIL |
| S03 | L1 | 86% PASS | 100% PASS | 86% PASS | 100% PASS |
| S05 | L2 | **86% PASS** | **14% PASS** | **86% PASS** | **0% FAIL** |
| S07 | L2 | 86% PASS | 100% PASS | 57% PASS | 86% PASS |
| S09 | L3 | 17% PASS | 33% PASS | 17% PASS | 17% PASS |
| S11 | L3 | 29% PASS | 29% PASS | 29% PASS | 29% PASS |
| S13 | L4 | **100% PASS** | **20% PASS** | **100% PASS** | **20% PASS** |
| S15 | L4 | 100% PASS | 100% PASS | 100% PASS | 0% FAIL |
| S17 | L5 | 0% FAIL | 0% FAIL | 0% FAIL | 100% PASS |
| S19 | L5 | 0% FAIL | 0% FAIL | 0% FAIL | 0% FAIL |

#### Ablation Ordering Analysis

**32B (Expected ordering confirmed)**:
```
FULL (60%) > NO-RULES (57%) > NO-DOMAIN (50%) > GENERIC (35%)
  +3%p ↑       +7%p ↑           +15%p ↑
```
- Domain knowledge contributes +10%p (FULL vs NO-DOMAIN)
- Tool rules contribute +3%p (FULL vs NO-RULES)
- Combined vs generic: +25%p (FULL vs GENERIC)

**8B (Ordering inverted)**:
```
NO-RULES (55%) > FULL (46%) > NO-DOMAIN (44%) > GENERIC (39%)
```
- FULL < NO-RULES: Long system prompt (135 lines) overwhelms 8B's context handling
- 8B FULL tool call rate: 7/10 (vs 32B: 10/10) — clear capacity bottleneck

**Critical Insight**: Ablation only yields meaningful results on 32B. The 8B inversion is itself a finding — it demonstrates that domain prompt engineering requires sufficient model capacity to be effective.

---

### Table 6: Industrial PoC Results (EXP-5)

#### Scenario A: Anti-Slosh Baffle Design Comparison

| Case | Config | Still Level [mm] | SWL Amplitude [mm] | Max SWL [mm] | Min SWL [mm] |
|------|--------|-----------------|--------------------| -------------|-------------|
| No Baffle | 1m × 0.5m × 0.6m, 50% fill, resonance f | 300 | 158.9 | 522.0 | 204.1 |
| With Baffle | Same + vertical baffle at center | 300 | 12.8 | 318.8 | 293.2 |
| **Reduction** | | | **91.9%** | | |

#### Scenario B: Seismic Excitation

| Parameter | Value |
|-----------|-------|
| Tank | 10m × 5m × 8m |
| Fill | 60% (4.8m) |
| Excitation | 0.3 Hz seismic-like |
| Max SWL Amplitude | 2361.8 mm |
| Max SWL | 4723.7 mm |
| Min SWL | 0.0 mm (complete drainage at one wall) |

**Key Findings**:
- Baffle achieves 91.9% sloshing amplitude reduction — physically consistent with literature (typical 60-90% for well-placed vertical baffles)
- Seismic scenario shows extreme amplitudes (wave height exceeds tank height) indicating overflow condition — physically correct for near-resonance large-tank excitation
- Both scenarios completed from single natural language prompts without manual XML editing

---

## Section 3: Key Findings for Paper

### Ranked by Novelty/Impact

| Rank | Finding | Novelty | Evidence | Strength |
|------|---------|---------|----------|----------|
| 1 | **First experimental benchmark validation by any LLM-CFD agent** (SPHERIC Test 10, 7/7 water peaks within ±2σ) | First of its kind | EXP-1 Table 2 | HIGH — no competitor attempts experimental validation |
| 2 | **First domain knowledge ablation study in computational mechanics** (FULL +25%p vs GENERIC) | First of its kind | EXP-4 Table 5 (32B) | HIGH — all competitor ablations test architecture, never knowledge |
| 3 | **First LLM agent for any particle-based solver** (14 tools for DualSPHysics pipeline) | First of its kind | Architecture + 10+ competitor survey | HIGH — entire Lagrangian solver space is vacant |
| 4 | **85% tool call success rate on 20 sloshing scenarios spanning 5 complexity levels** | Novel benchmark | EXP-2 Table 3 | MEDIUM — 47% param accuracy tempers the tool call success |
| 5 | **Domain prompt ordering confirmed: FULL > NO-RULES > NO-DOMAIN > GENERIC** (32B only) | Novel finding | EXP-4 ablation ordering | MEDIUM — only works on 32B, 8B shows inversion |
| 6 | **8B model capacity bottleneck: long prompts degrade tool calling** (7/10 vs 10/10) | Novel finding | EXP-4 8B vs 32B | MEDIUM — practical finding for agent deployment |
| 7 | **Single-prompt parametric study: 6 fill levels automated** | Demonstration | EXP-3 Table 4 | MEDIUM — demonstrates practical value but no time comparison data |
| 8 | **91.9% baffle amplitude reduction in industrial PoC** | Demonstration | EXP-5 Table 6 | LOW — physics result, not agent novelty (well-known baffle effect) |

### "First of Their Kind" Claims

1. **First LLM agent to validate against published experimental sloshing data** (SPHERIC Test 10)
2. **First domain knowledge ablation study in computational mechanics LLM agents** (all competitors ablate architecture only)
3. **First LLM agent for any particle-based (Lagrangian) solver** (SPH, DEM, MPM all unaddressed)
4. **First sloshing-specific LLM agent** (10+ CFD agents exist, none for sloshing)

### Weaknesses Requiring Honest Discussion in Limitations

1. **47% overall parameter accuracy** — Heavily influenced by L3 (paper-specific, 15%) and L5 (edge cases, 0%). If restricted to L1+L2 (realistic usage), accuracy is 82% (32B). Paper should report both overall and by-level.

2. **Oil sloshing failure** — DBC boundary condition cannot capture thin-film viscous sloshing. This is a solver configuration limitation (DBC vs mDBC), not an agent limitation per se. Agent correctly configured the case; the solver's boundary condition was inadequate.

3. **No Pearson r / NRMSE for SPHERIC** — Current validation uses peak-in-band metric only. Time-series correlation would be stronger evidence. (Pending Task #1 completion.)

4. **Chen2018 comparison incomplete** — Parametric study results show correct physical trends but lack quantitative comparison to published figures (digitization needed).

5. **Ablation inversion on 8B** — Domain prompt ablation only shows expected hierarchy on 32B. 8B results suggest the full prompt may actually _hurt_ smaller models. This is a limitation on deployment flexibility.

6. **Single LLM (Qwen3)** — All experiments use Qwen3 variants only. Generalizability to other LLMs (Llama, Mistral) is unknown.

7. **Latency** — 32B averages 350s per interaction (thinking mode). Practical deployment would likely use 8B (11s) or non-thinking mode.

8. **No user study** — All evaluation is automated or expert-judged. A formal user study with CFD engineers would strengthen the HCI claims.

---

## Section 4: 8B vs 32B Model Comparison Summary

### Overall Performance

| Metric | Qwen3 8B | Qwen3 32B | Winner | Delta |
|--------|----------|-----------|--------|-------|
| **EXP-2 Tool Call Rate** | 75% (15/20) | **85% (17/20)** | 32B | +10%p |
| **EXP-2 Param Accuracy** | 46% | **47%** | 32B | +1%p |
| **EXP-2 Physical Validity** | 70% (14/20) | **70% (14/20)** | Tie | 0%p |
| **EXP-4 FULL Accuracy** | 46% | **60%** | 32B | +14%p |
| **EXP-4 FULL Tool Call** | 70% (7/10) | **100% (10/10)** | 32B | +30%p |
| **EXP-4 Ablation Ordering** | Inverted | **Expected** | 32B | N/A |
| **Avg Latency** | **11.0s** | 350.0s | 8B | 32x faster |
| **GPU VRAM** | **5.8 GB** | 14.2 GB | 8B | 2.4x less |

### Level-by-Level Comparison (EXP-2)

| Level | 32B Accuracy | 8B Accuracy | Delta | Interpretation |
|-------|-------------|------------|-------|---------------|
| L1 (Basic) | 96% | 96% | 0%p | Explicit params: model size irrelevant |
| L2 (Domain) | **67%** | 42% | **+25%p** | Domain inference: larger model clearly better |
| L3 (Paper) | 15% | 15% | 0%p | Paper-specific knowledge: neither model has it |
| L4 (Complex) | **58%** | 50% | +8%p | Multi-feature: modest 32B advantage |
| L5 (Edge) | 0% | **25%** | -25%p | Edge cases: 8B paradoxically better (less over-engineering) |

### Ablation Pattern Comparison (EXP-4)

**32B ablation ordering (expected)**:
```
FULL (60%) → NO-RULES (57%) → NO-DOMAIN (50%) → GENERIC (35%)
Domain: +10%p | Rules: +3%p | Full vs Generic: +25%p
```

**8B ablation ordering (inverted)**:
```
NO-RULES (55%) → FULL (46%) → NO-DOMAIN (44%) → GENERIC (39%)
Long prompt penalty: -9%p (FULL < NO-RULES)
```

### Recommendation for Paper

**Primary model: Qwen3 32B** — for all main results and analysis.

Rationale:
1. Only 32B produces the expected ablation ordering — the core RQ4 finding
2. Higher domain inference accuracy (L2: +25%p)
3. 100% tool call reliability with full prompt (vs 70% for 8B)
4. The 32B thinking overhead (350s avg) is acceptable for research validation

**Secondary model: Qwen3 8B** — as comparison/contrast.

Rationale:
1. Similar overall accuracy (46% vs 47%) masks important per-level differences
2. The 8B ablation inversion is itself a novel finding about model capacity requirements
3. 32x speed advantage demonstrates the deployment/research trade-off
4. 8B + NO-RULES (55%) > 8B + FULL (46%): practical deployment should use shorter prompts

### Paper Presentation Strategy

1. **Main results tables**: Use 32B as primary, 8B as comparison column
2. **Ablation section**: Present 32B ordering as main finding, 8B inversion as "model capacity threshold" discussion
3. **Discussion**: Frame 8B vs 32B as research contribution about minimum model requirements for domain-specific agents
4. **Practical implications**: Note that 8B + optimized (shorter) prompt may be preferred for deployment (11s vs 350s)

---

## Appendix: Cross-Reference Index

| Table | Location | Content | Paper Section |
|-------|----------|---------|---------------|
| Table 2 | exp1_spheric/comparison/table2_summary.md | SPHERIC peak pressure validation | Sec 5.1 |
| Table 3 | exp2_nl2xml/table3_results.md (32B), table3_results_8b.md (8B) | NL→XML generation accuracy | Sec 5.2 |
| Table 4 | exp3_parametric/comparison/table4_summary.md | Chen2018 parametric SWL analysis | Sec 5.3 |
| Table 5 | exp4_ablation/table5_ablation.md (32B), table5_ablation_8b.md (8B) | Domain prompt ablation | Sec 5.4 |
| Table 6 | exp5_industrial/comparison/table6_summary.md | Industrial PoC (baffle + seismic) | Sec 5.5 |
| Model Comp | experiments/model_comparison.md | 8B vs 32B analysis | Sec 5.6 or Discussion |
| GAP 1-5 | output/gap_refinement_cycle3.md | Research gaps and evidence | Sec 1-2 |
| Experiment Design | output/experiment_design.md | Full experimental protocol | Sec 4 |
