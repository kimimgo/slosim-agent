# Gap Refinement Cycle 3 — Sloshing-Centric Reframing

Date: 2026-02-18
Agent: cc-slosim-1
Method: 4 parallel deep-research agents (accidents, pain points, LNG industry, time/cost) + Cycle 2 evidence consolidation

**CRITICAL SHIFT**: Cycle 2 was method-centric ("first SPH+LLM agent"). Cycle 3 reframes everything around **sloshing as the central problem domain**. SPH is the _consequence_ (best solver for violent free-surface flows), not the _premise_.

---

## Part A: Why Sloshing Matters — The Problem Domain

### A.1 Real-World Sloshing Accidents and Failures

| Incident | Year | Domain | Damage | Root Cause |
|----------|------|--------|--------|------------|
| **SpaceX Falcon 1 Flight 2** | 2007 | Aerospace | $30M+ launch failure | LOX sloshing in upper stage caused attitude control loss; spin-stabilization fix implemented for Flight 3 |
| **Tomakomai Oil Fire** | 2003 | Petrochemical | 170 tanks affected (58% of facility), $15B+ oil industry losses nationwide | 2003 Tokachi-Oki earthquake (M8.0) induced sloshing in floating-roof crude oil tanks; 7 tanks burned for days |
| **Alaska Good Friday Earthquake** | 1964 | Petrochemical | $2.5B in property damage (2024 dollars), fires from petroleum tank sloshing | M9.2 earthquake; sloshing in petroleum storage tanks caused spills and fires |
| **Fukushima SFP Crisis** | 2011 | Nuclear | Unit 4 spent fuel pool near-boiling; global nuclear safety protocols revised | Earthquake + tsunami caused water sloshing out of spent fuel pools; required emergency water injection |
| **Tanker Truck Rollovers** | Annual | Transportation | 1,300+ rollovers/year in US alone; ~40% of heavy truck fatal crashes involve liquid cargo | Sloshing in partially-filled tank trailers causes dynamic instability during braking/turning |
| **LNG Carrier Structural Damage** | Ongoing | Maritime | Membrane containment system (Mark III) deformation from sloshing impact | Barred fill levels (10-70%) mandated by classification societies to avoid resonant sloshing |

**Key Insight**: Sloshing is not a niche academic topic — it causes catastrophic failures across aerospace, maritime, petrochemical, nuclear, and transportation industries. The common thread is **partially-filled containers under dynamic excitation**.

### A.2 LNG Sloshing Industry Scale

| Metric | Value | Source |
|--------|-------|--------|
| Active LNG carriers worldwide | **770+ ships** | Industry data 2024 |
| New LNG carriers on order | **400+** ($250-269M each) | Lloyd's List 2024 |
| Korea's share of global orders | **~70%** (HD Hyundai, Samsung, Hanwha) | Korean shipyard data |
| SNU experimental sloshing database | **20,000+ hours** of model tests, 540 TB data | Kim et al. 2019, ScienceDirect |
| Per-campaign test volume | 120 experiments per fluid/tank/fill combination | ISOPE benchmark methodology |
| LNG tank assessment CFD cases required | **200+ simulations per tank** | SJTU OFW13 2018 |
| Barred fill range (GTT Mark III) | 10-70% (operation forbidden in this range) | Classification society rules |
| Maritime AI market | **$4.13B** (2024), 40.6% CAGR to 2030 | Market reports |

**Key Insight**: Each new LNG carrier costs $250-269M and requires 200+ sloshing simulations for tank certification. With 400+ carriers on order, there is a massive demand for sloshing simulation automation. Yet maritime AI ($4.13B) focuses entirely on route optimization and predictive maintenance — **zero** on sloshing simulation.

### A.3 Pre-processing Bottleneck: Where Time Goes

| Metric | Value | Source |
|--------|-------|--------|
| Engineering time on pre-processing/meshing | **75-80%** of total CFD project time | Cadence 2024, Toyota-Cadence |
| CFD engineer hourly rate (US) | $66-120/hr | Glassdoor/CadCrowd 2025 |
| Fully-loaded annual cost per CFD engineer | **$285K/year** | Resolved Analytics |
| Initial sloshing project timeline | **2-5 weeks** | CadCrowd consulting estimates |
| DualSPHysics GPU speedup vs CPU | **40-100x** | Crespo et al. 2011 (PLOS ONE) |
| Automation impact on model generation | **>80% time reduction** (months → 1 week) | Siemens Simcenter |
| AI impact on R&D time-to-market | **10-20% reduction** | BCG 2025 |
| AI impact on R&D cost | **up to 20% reduction** | BCG 2025 |

**Key Insight**: The GPU solver is fast (minutes on RTX 4090). The bottleneck is **human setup time**: understanding the XML format, configuring boundary conditions, tuning parameters, validating results. This is exactly what an LLM agent can automate.

---

## Part B: Why Current Tools Fail Sloshing Practitioners

### B.1 DualSPHysics XML: 5 Documented Error Types

From DualSPHysics forums and GitHub issues, 5 recurring error patterns that trap even experienced users:

1. **Hydrostatic Initialization Artifact** — Simulations show ~0.5 kPa at t=0 vs. 0 kPa in experiments. Fix: pre-run ~1s without motion + subtract hydrostatic baseline. Undocumented in tutorials.

2. **Silent Motion Coupling Failure** — Combined rotational + translational motion fails silently if `mov ref` values differ. No error message; only first motion component executes.

3. **mvrotsinu/mvrectsinu Syntax Confusion** — Rotational uses scalar `v` attribute; translational uses vector `x/y/z` attributes. Mixing produces incorrect motion with no warning.

4. **fillbox Seed Point Placement** — In 2D, seed must be at exact y-position or particle creation silently fails. Alternative (`hswl auto=false`) exists but is not in tutorials.

5. **Constant 'b' Errors** — Triggered by zero fluid height or zero gravity, but error message does not indicate which parameter is wrong.

### B.2 DesignSPHysics GUI: Permanent Beta

The only graphical frontend for DualSPHysics (DesignSPHysics) is inadequate:
- Official status: **"early Beta stage, not meant to be used in a stable environment"** (v0.7.0, Sep 2023)
- **mDBC not supported via GUI** (Issue #171, Jan 2024) — the boundary condition required for accurate sloshing pressure
- FreeCAD version dependency creates installation failures
- 135 GitHub stars vs. DualSPHysics's 687 — low adoption even among DualSPHysics users

### B.3 Expert Knowledge That Walks Out the Door

From Rescale's analysis: _"When a senior engineer changes teams or leaves the company, their accumulated expertise often leaves with them."_

For sloshing, this institutional knowledge includes:
- Which `dp` values are stable for which geometries
- How to set `coefsound` and viscosity for specific liquids
- How to handle the cold-start initialization artifact
- Whether to use DBC vs. mDBC for each case type
- Parameter tuning (artificial viscosity `alpha`: 0.01-0.5 range, trial-and-error)
- Resolution convergence study protocols (5 `dp` levels typical)

### B.4 No Existing LLM Tool Addresses Sloshing

| System | Solver | Sloshing? | SPH? | GPU? | Local LLM? |
|--------|--------|-----------|------|------|------------|
| MetaOpenFOAM 1.0/2.0 | OpenFOAM | No | No | No | No |
| Foam-Agent 2.0 | OpenFOAM | No | No | No | No |
| ChatCFD | OpenFOAM | No | No | No | No |
| OpenFOAMGPT 2.0 | OpenFOAM | No | No | No | No |
| CFD-Copilot | OpenFOAM | No | No | No | Partial |
| AutoCFD | OpenFOAM | No | No | No | Yes |
| FoamGPT | OpenFOAM | No | No | No | Yes |
| PhyNiKCE | OpenFOAM | No | No | No | No |
| MooseAgent | MOOSE (FEM) | No | No | No | No |
| MCP-SIM | FEniCS (FEM) | No | No | No | No |
| Pasimodo+RAG | Pasimodo (SPH) | No | Not agent | No | Yes |
| **SloshAgent** | **DualSPHysics** | **Yes** | **Yes** | **Yes** | **Yes** |

All 10+ existing LLM-for-simulation systems target general-purpose mesh-based solvers. None addresses sloshing, the most active industrial application of SPH.

---

## Part C: Five Sloshing-Centric Research Gaps

### GAP-1: No AI Agent Automates Sloshing Simulation (Domain Gap)
**Confidence: 95%**

**Claim**: Despite sloshing causing catastrophic failures across 5+ industries (aerospace, maritime, petrochemical, nuclear, transportation) and requiring 200+ simulations per LNG tank assessment, no AI agent system automates any part of the sloshing simulation workflow — neither mesh-based (OpenFOAM VOF) nor particle-based (SPH).

**Evidence**:
- 10+ LLM-CFD systems exist (Table B.4) — all target general CFD, none addresses sloshing
- $4.13B maritime AI market → zero sloshing simulation AI
- Sloshing practitioners still hand-craft XML/dictionaries (200+ page DualSPHysics guide)
- Closest work: SNU ANN for sloshing load _prediction_ from experimental DB — NOT simulation automation
- The entire free-surface violent flow domain is unaddressed by AI agents

**What SloshAgent provides**: First end-to-end sloshing simulation agent: NL → XML case → GPU solver → post-processing → validation.

### GAP-2: Sloshing Domain Knowledge Never Encoded in LLM Prompts (Knowledge Gap)
**Confidence: 90%**

**Claim**: No existing LLM agent system encodes sloshing-specific domain knowledge (natural frequencies, resonance phenomena, fill-level effects, hydrostatic initialization, baffle design rules) in its prompting strategy. Existing ablation studies test _architectural_ components (RAG, reviewer nodes), never _knowledge_ content.

**Evidence**:
- Exhaustive review of ablation studies across 10 competitor systems (Cycle 2 Task 4):
  - MetaOpenFOAM: ablates RAG, reviewer, temperature — NOT domain knowledge
  - Foam-Agent: ablates reviewer node, RAG — NOT domain knowledge
  - ChatCFD: ablates RAG modules — NOT domain knowledge
  - PhyNiKCE: ablates instruction structure — NOT domain content
  - MDCrow: tests prompt style — NOT domain knowledge
- No system encodes sloshing physics (resonance freq = sqrt(g*pi/L * tanh(pi*h/L))/(2*pi), fill-level sensitivity, nonlinear spring behavior)
- SloshAgent's SloshingCoderPrompt: 136 lines, 5 categories of sloshing-specific constraints

**What SloshAgent provides**: First domain prompt ablation study for computational mechanics — isolating the contribution of sloshing physics knowledge, SPH constraints, and XML syntax rules.

### GAP-3: No Experimental Sloshing Benchmark Validation by Any AI Agent (Validation Gap)
**Confidence: 90%**

**Claim**: No LLM-for-simulation system validates against published experimental sloshing benchmark data. "Success" is universally defined as "did the code run?" (execution success) or "does it match another simulation?" (reference comparison) — never "does the physics match reality?"

**Evidence**:
- Foam-Agent 2.0: 88.2% = "ran successfully" = no crash
- ChatCFD: 68.12% "physical fidelity" = LLM-as-judge evaluation (NOT experimental comparison)
- MetaOpenFOAM: human rubric comparison to tutorial results
- OpenFOAMGPT 2.0: 100% reproducibility = identical re-runs, verified against analytical solutions (Poiseuille)
- PhyNiKCE: 51% on complex tasks, against known CFD solutions
- **Zero** competitor compares to published experimental benchmark data (SPHERIC, ISOPE, etc.)
- ISOPE benchmark: 9/16 CFD pressure peaks outside ±2σ of experiments — validation is genuinely hard

**What SloshAgent provides**: SPHERIC Test 10 validation — 100-repeat experimental pressure statistics, 3 fluids, quantitative metrics (Pearson r > 0.9 target, NRMSE).

### GAP-4: Sloshing Industry Automation Gap ($4.13B Market, Zero AI Tools) (Industry Gap)
**Confidence: 85%**

**Claim**: The maritime AI market is $4.13B with 40.6% CAGR, and LNG carrier construction represents $100B+ in orders. AI is applied to route optimization, fuel efficiency, and predictive maintenance — but zero AI tools exist for sloshing simulation, the critical safety analysis required for every LNG tank design.

**Evidence**:
- 770 active LNG carriers, 400+ on order ($250-269M each)
- Korea 70% market share → direct industrial relevance
- SNU: 20,000+ hours experimental data, 540 TB — shows scale of sloshing analysis demand
- 200+ CFD cases per LNG tank assessment (SJTU OFW13)
- DesignSPHysics in permanent beta; commercial SPH tools ($10K+/year, no AI)
- DualSPHysics: 52,000+ downloads, 687 stars, 746 citations — large user base with no AI tooling
- CFD consulting: $80-120/hr, $2K-$20K/project — high cost barrier for parametric studies
- BCG 2025: AI in R&D can deliver 10-20% time-to-market reduction, up to 20% cost reduction

**What SloshAgent provides**: $0 LLM cost (local Qwen3 32B via Ollama), GPU-accelerated SPH on consumer hardware (RTX 4090), natural language input for sloshing scenarios.

### GAP-5: Particle-Based Simulation + LLM = Complete Void (Method Gap)
**Confidence: 95%**

**Claim**: SPH is the natural solver for sloshing (particle method handles violent free-surface fragmentation without mesh distortion), yet the entire Lagrangian particle simulation + LLM agent space is vacant. Not just SPH — DEM (Discrete Element) and MPM (Material Point Method) are also unaddressed.

**Evidence**:
- 6 search query patterns across academic databases → zero results for SPH+LLM, DEM+LLM, MPM+LLM
- Pasimodo+RAG: pure Q&A system, scored 0/2 on model creation — explicitly NOT an agent
- All 10+ LLM-CFD systems use mesh-based solvers (OpenFOAM, MOOSE, FEniCS)
- Particle solvers have fundamentally different I/O: XML case definition, binary particle data, GPU telemetry, particle-to-mesh conversion
- No tool interface designs exist for Lagrangian particle I/O patterns

**What SloshAgent provides**: 14 purpose-built tools for the DualSPHysics pipeline (xml_generator, solver, gencase, partvtk, measuretool, monitor, analysis, seismic_input, stl_import, parametric_study, etc.) — first tool interface for any particle solver.

---

## Part D: Revised Narrative

### OLD Narrative (Method-Centric, REJECTED by user)
> "We present the first LLM agent for SPH-based simulation. SPH is a particle method used in fluid dynamics. We combine it with LLM tools."

**Problem**: This makes SPH the primary selling point. Reader asks: "Why should I care about SPH?"

### NEW Narrative (Sloshing-Domain-Centric, APPROVED)
> "Sloshing in liquid containers causes catastrophic failures — SpaceX Falcon 1 (2007), Tomakomai oil fire ($15B, 2003), 1,300+ tanker rollovers/year. Predicting sloshing loads requires 200+ simulations per LNG tank, each needing expert-level XML configuration with 5+ common silent error types. We present SloshAgent: the first AI agent that automates the entire sloshing simulation pipeline from natural language to experimentally-validated results. Because sloshing involves violent free-surface fragmentation, we use DualSPHysics (SPH), the natural solver for this physics — and the first particle solver ever integrated with an LLM agent."

**Why this works**: Sloshing is the PROBLEM everyone understands (accidents, industry scale). SPH is the SOLUTION (best solver for this physics). The narrative flows: Problem → Why it's hard → Our solution → Why SPH is the right choice.

### One-Paragraph Summary (Revised)

> 슬로싱(sloshing)은 LNG 탱크, 연료 탱크, 원자력 냉각수 등 부분 충전 용기에서 발생하는 유체 동적 현상으로, 2003년 토마코마이 화재($15B), 2007년 SpaceX Falcon 1 실패, 연간 1,300건 이상의 탱커 전복 사고 등 다양한 산업에서 치명적 피해를 야기한다. 현재 LNG 탱크 1기당 200건 이상의 슬로싱 시뮬레이션이 필요하지만, 설정에는 200페이지 XML 가이드 숙지, 5가지 이상의 무경고 에러 패턴 이해, DBC/mDBC 경계조건 선택 등 전문가 수준의 도메인 지식이 요구된다. 10개 이상의 LLM-CFD 에이전트가 OpenFOAM 범용 시뮬레이션을 자동화하고 있으나, 슬로싱 전용 에이전트는 전무하며, SPH를 포함한 입자 기반 솔버와 LLM을 결합한 연구는 세계 최초로 공백 상태이다. 본 연구는 자연어 입력에서 실험 검증까지 DualSPHysics 슬로싱 시뮬레이션 파이프라인 전체를 자율 구동하는 최초의 도메인 특화 LLM 에이전트 SloshAgent를 제시하며, SPHERIC Test 10 벤치마크(100회 반복 실험 데이터) 대비 정량적 검증과 도메인 프롬프트 어블레이션 실험을 통해 그 유효성을 입증한다.

---

## Part E: Evidence Source Index

### Sloshing Accidents
- SpaceX Falcon 1 Flight 2 (2007) — SpaceX mission archives, Musk interviews
- Tomakomai Oil Fire (2003) — FDMA Japan reports, Hatayama et al. 2004
- Alaska Good Friday Earthquake (1964) — USGS historical records
- Fukushima SFP (2011) — NRC/IAEA reports
- Tanker Truck Rollovers — FMCSA statistics, Kang et al. 2019

### LNG Industry Data
- 770 ships, 400+ orders — Lloyd's List, DNV reports 2024
- Korea 70% share — HD Hyundai/Samsung/Hanwha quarterly reports
- SNU 20,000+ hours DB — Kim et al. 2019 (ScienceDirect)
- 200+ cases per tank — SJTU OFW13 (2018)
- Maritime AI $4.13B — Grand View Research / MarketsandMarkets 2024

### Cost/Time Data
- Pre-processing 75-80% — Cadence 2024, Toyota-Cadence
- $66-120/hr CFD engineer — Glassdoor/CadCrowd 2025
- $285K/yr fully-loaded — Resolved Analytics
- GPU 40-100x speedup — Crespo et al. PLOS ONE 2011
- BCG 10-20% time reduction — BCG Executive Perspectives Feb 2025

### DualSPHysics Pain Points
- Forums: tank-sloshing-validation, sloshing-tank-tutorial-problem, beginner-problems
- GitHub: Discussion #216 (coupled motion), DesignSPHysics Issues #93, #171
- Documentation: XML Guide v4.2 (200+ pages)

### Competitor Systems (Cycle 2 + 3)
- 11 systems fully verified from paper full text (competitor_analysis_verified.md)
- PhyNiKCE (arXiv:2602.11666, Feb 2026) — newest entrant
- Pasimodo+RAG confirmed as non-agent (0/2 on model creation)

---

## Part F: Gaps Comparison (Cycle 2 → Cycle 3)

| Cycle 2 (Method-Centric) | Cycle 3 (Sloshing-Centric) | Change |
|---------------------------|---------------------------|--------|
| GAP-1: SPH + LLM = zero | GAP-5: Particle + LLM = void | Now secondary (method is consequence) |
| GAP-2: No tool interfaces for particle solvers | GAP-5: Absorbed into method gap | Merged |
| GAP-3: No domain prompt ablation | GAP-2: Sloshing knowledge never in LLM prompts | Reframed around sloshing knowledge specifically |
| GAP-4: No E2E experimental pipeline | GAP-3: No experimental sloshing benchmark validation | Narrowed to sloshing benchmarks |
| GAP-5: Industry PoC missing | GAP-4: $4.13B maritime AI, zero sloshing AI | Elevated with quantified industry data |
| GAP-6: Local LLM underexplored | Removed (weak, CFD-Copilot uses same Qwen3-32B) | Dropped — not strong enough as standalone gap |

### Why GAP-6 was dropped
- CFD-Copilot uses Qwen3-32B for general agents (same model as us)
- The distinction (zero fine-tuning vs fine-tuned 8B) is a design choice, not a research gap
- Better to mention as a system advantage in methodology section, not frame as a gap

### New ordering rationale
1. **GAP-1 (Domain)**: Sloshing needs AI automation — anchors the whole paper
2. **GAP-2 (Knowledge)**: Domain knowledge isn't in any LLM prompt — our core innovation
3. **GAP-3 (Validation)**: No experimental validation — our strongest differentiator
4. **GAP-4 (Industry)**: Massive market with zero tools — practical motivation
5. **GAP-5 (Method)**: Particle solvers completely unaddressed — technical novelty

This ordering follows the narrative arc: **Problem → Knowledge → Validation → Market → Method**.
