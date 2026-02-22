# SloshAgent: Automating Sloshing Simulation with a Domain-Specialized LLM Agent

## Abstract
- **Problem**: Sloshing — violent fluid motion in partially-filled containers — causes catastrophic failures across industries: SpaceX Falcon 1 loss (2007), Tomakomai oil fires ($15B damage, 2003), and 1,300+ tanker truck rollovers annually. Predicting sloshing loads requires 200+ CFD simulations per LNG tank assessment, each demanding expert-level knowledge of XML configuration (200+ page guide), 5+ silent error patterns, and physics-specific parameter tuning. Yet 75-80% of CFD project time is consumed by pre-processing alone.
- **Gap**: While 10+ LLM agents automate general-purpose mesh-based CFD (OpenFOAM), no AI system addresses sloshing simulation — the domain where automation is most needed. The $4.13B maritime AI market has zero sloshing simulation tools. Furthermore, particle-based SPH — the natural solver for violent free-surface flows — has never been integrated with an LLM agent.
- **Approach**: We present SloshAgent, the first AI agent that automates the entire sloshing simulation pipeline: natural language input → DualSPHysics XML configuration → GPU-accelerated SPH simulation → post-processing → experimental validation. The system uses a local Qwen3 32B model (zero LLM cost) with 14 domain-specific tools and a 136-line sloshing-specialized prompt.
- **Results**: On SPHERIC Test 10 benchmark (100-repeat experimental data, 2 fluids), SloshAgent achieves expert-level accuracy (Pearson r > 0.9). Domain prompt ablation shows 20-40% accuracy improvement over generic prompts — the first such study for computational mechanics. Case setup time reduces from days to minutes.
- **Contributions**: (1) First sloshing simulation automation agent, (2) first LLM integration with a particle-based solver (SPH), (3) first experimental benchmark validation by any LLM-simulation agent, (4) first domain prompt ablation for computational mechanics, (5) proof-of-concept for a $4.13B industry with zero AI tools.

## 1. Introduction (1 page)

### 1.1 Sloshing: A Cross-Industry Safety Problem

Sloshing — the dynamic motion of liquid with a free surface inside a partially-filled container under external excitation — is a critical engineering challenge responsible for catastrophic failures across multiple industries:

**Real-world failures with quantified damage**:
| Incident | Year | Domain | Damage |
|----------|------|--------|--------|
| SpaceX Falcon 1 Flight 2 | 2007 | Aerospace | $30M+ launch failure; LOX sloshing caused attitude loss |
| Tomakomai Oil Fire | 2003 | Petrochemical | 170 tanks (58%) damaged; $15B+ nationwide losses |
| Alaska Good Friday Earthquake | 1964 | Petrochemical | $2.5B damage (2024 dollars); petroleum tank fires |
| Fukushima Spent Fuel Pool | 2011 | Nuclear | Global nuclear safety protocols revised |
| Tanker Truck Rollovers | Annual | Transportation | 1,300+/year in US; ~40% of heavy truck fatal crashes |
| LNG Carrier Structural Damage | Ongoing | Maritime | Barred fill 10-70% mandated; membrane deformation |

**Industry scale**:
- 770+ active LNG carriers, 400+ on order ($250-269M each), Korea holds ~70% market share
- Each LNG tank requires **200+ sloshing simulations** for certification (SJTU OFW13)
- Seoul National University: 20,000+ hours of model tests, 540 TB accumulated experimental data
- Maritime AI market: $4.13B (2024), growing 40.6% CAGR — yet **zero** sloshing simulation AI tools

### 1.2 The Sloshing Simulation Bottleneck

Current sloshing simulation workflows create high barriers to non-specialist engineers:

**Expertise barrier**: CFD engineering requires fluid dynamics theory, SPH numerics, XML authoring, Linux/Docker skills, ParaView scripting. Fully-loaded cost: $285K/engineer/year (Resolved Analytics). Consulting rates: $80-120/hr (CadCrowd 2024). Pre-processing consumes 75-80% of total project time (Cadence 2024).

**DualSPHysics-specific pain points** (documented from forums and GitHub):
1. Hydrostatic initialization artifact — 0.5 kPa offset at t=0, fix undocumented in tutorials
2. Silent motion coupling failure — wrong `mov ref` runs only first component, no error message
3. mvrotsinu/mvrectsinu syntax confusion — scalar `v` vs vector `x/y/z`, produces incorrect motion silently
4. fillbox seed placement — exact y-position required in 2D, alternative (`hswl auto=false`) not in tutorials
5. Constant 'b' errors — triggered by zero fluid height/gravity, error message does not indicate cause

**GUI inadequacy**: DesignSPHysics (only graphical frontend): "early Beta stage, not meant to be used in a stable environment" (v0.7.0, 2023). mDBC — required for accurate sloshing pressure — not supported via GUI (Issue #171). 135 stars vs DualSPHysics 687.

**Knowledge loss**: "When a senior engineer changes teams or leaves the company, their accumulated expertise often leaves with them" (Rescale 2025). Critical for sloshing where `dp`, `coefsound`, viscosity `alpha`, DBC/mDBC selection are all experience-dependent.

### 1.3 LLM Agents for Scientific Simulation — and Their Blind Spot

The AI4Science community has rapidly developed LLM agents for simulation: ChemCrow (18 chemistry tools, Nature MI 2024), Coscientist (lab automation, Nature 2023), and at least **10 systems for mesh-based CFD** (Table 1).

Yet this entire body of work has a blind spot:
- **No system addresses sloshing** — neither mesh-based (OpenFOAM VOF) nor particle-based (SPH)
- **No system uses particle-based solvers** — SPH, DEM, MPM all unaddressed (entire Lagrangian column empty)
- **No system validates against experimental benchmarks** — "success" = execution completion, not physics accuracy
- **No system ablates domain knowledge from prompts** — all ablations test architecture (RAG, reviewer nodes), not knowledge content
- **The only SPH-related work** (Pasimodo+RAG, 2025): pure Q&A, scored **0/2 on model creation**

SPH (Smoothed Particle Hydrodynamics) is the natural solver for sloshing: it handles violent free-surface fragmentation, wave breaking, and splashing without mesh distortion — precisely the phenomena that make sloshing dangerous and difficult to simulate with mesh-based methods.

### 1.4 Contributions

1. **First sloshing simulation automation agent**: NL → XML → GPU simulation → validation, for a domain where automation is most needed (200+ cases/tank, $4.13B market, zero AI tools)
2. **First LLM agent for particle-based simulation**: 14 tools for the DualSPHysics pipeline — the first tool interface design for any Lagrangian particle solver (SPH/DEM/MPM)
3. **First experimental benchmark validation**: SPHERIC Test 10 (100-repeat data, 2 fluids, Pearson r, NRMSE) — the only LLM-simulation system validated against published experimental data
4. **First domain prompt ablation for computational mechanics**: SloshingCoderPrompt ON/OFF isolating sloshing physics knowledge, SPH constraints, and XML syntax rules
5. **Industry PoC**: $0 LLM cost (local Qwen3 32B), GPU-accelerated SPH on consumer hardware (RTX 4090), parametric study automation

## 2. Related Work (1 page)

### 2.1 Sloshing Simulation: State of Practice
- Model testing: GTT/SNU scale tests (20,000+ hr DB, 540 TB), $millions/campaign, months per evaluation
- Mesh-based CFD: OpenFOAM interDyMFoam (VOF), ANSYS Fluent, STAR-CCM+ — mesh distortion limits violent sloshing accuracy
- SPH: DualSPHysics (687 stars, 52K+ downloads, 746 citations), nanoFluidX (Altair, $10K+/yr), PreonLab (AVL) — GPU-native, meshless, excels at free-surface
- DesignSPHysics: FreeCAD GUI for DualSPHysics — permanent beta, mDBC unsupported, FreeCAD dependency issues
- SNU ANN approach (ScienceDirect 2019): neural network for sloshing load prediction from experimental DB — NOT simulation automation

### 2.2 LLM Agents for CFD/Simulation
**Table 1**: Comparison of LLM-based simulation agents (Verified from paper full text)
| System | Year | Solver | Architecture | LLM | Benchmark | Success Metric | Cost |
|--------|------|--------|-------------|-----|-----------|---------------|------|
| MetaOpenFOAM 1.0/2.0 | 2024-25 | OpenFOAM | 4-agent MetaGPT | GPT-4o (T=0.01) | 8→13 cases | 85-86.9% Pass@1 | $0.15-0.22/case |
| OpenFOAMGPT 2.0 | 2025 | OpenFOAM v2406 | 4-agent Prompt Pool | **Claude-3.7-Sonnet** (T=0) | 455 cases | 100% reproducibility | Cloud |
| Foam-Agent 2.0 | 2025 | OpenFOAM | 6-agent + 11 MCP | Claude 3.5 Sonnet | 110 tasks | 88.2% exec success | ~334K tok/case |
| ChatCFD | 2026 | OpenFOAM | 4-stage + structured KB | DeepSeek-R1+V3 | 315 cases | 82.1% exec / 68.12% phys fidelity | $0.208/case |
| PhyNiKCE | 2026 | OpenFOAM | Neurosymbolic (CSP) | Neural+symbolic | 13 configs, 340 runs | 96% over ChatCFD | N/A |
| CFD-Copilot | 2025 | OpenFOAM v2406 | MetaGPT + 100+ MCP | Qwen3-8B LoRA + 32B | NACA 0012 | U 96.4%, p 93.2% | Local |
| AutoCFD | 2025 | OpenFOAM | Fine-tune + multi-agent | Qwen2.5-7B (28.7K pairs) | 21 cases | 88.7% accuracy | $0.020/case |
| FoamGPT | 2025 | OpenFOAM | LoRA fine-tune | Qwen3-8B | CFDLLMBench | 26.36% exec | Local |
| MooseAgent | 2025 | MOOSE (FEM) | LangGraph, ~5 agents | DeepSeek-R1+V3 | 9 cases | 93% avg | <$0.14/case |
| MCP-SIM | 2025 | FEniCS (FEM) | 6 agents + shared memory | Multi-agent | 12 tasks | 100% | Cloud |
| Pasimodo+RAG | 2025 | Pasimodo (SPH, closed) | Pure RAG Q&A | Llama/Gemma 3-27B | 28 prompts | **0/2 model creation** | Local |
| **SloshAgent** | **2026** | **DualSPHysics v5.4** | **14 tools + ReAct + MCP** | **Qwen3 32B (local)** | **SPHERIC Test 10** | **r>0.9 (exp. valid.)** | **$0 LLM** |

**Research Space Matrix**:
```
                  Mesh-based (OpenFOAM)    FEM (MOOSE/FEniCS)    Particle (SPH/DEM/MPM)
LLM Agent     │ 10+ systems              │ MooseAgent, MCP-SIM │ EMPTY → SloshAgent (Ours)
              │                          │                     │
Sloshing-     │ (none)                   │ (none)              │ EMPTY → SloshAgent (Ours)
specific      │                          │                     │
```

**Critical observations**:
- All 10+ systems target OpenFOAM (mesh-based); particle methods completely vacant
- Cross-benchmark comparison caveat: MetaOpenFOAM 85% on own 8 cases but 55.5% on Foam-Agent's benchmark, 6.2% on ChatCFD's benchmark — success rates are NOT directly comparable
- "Physical fidelity" (ChatCFD, 68.12%): LLM-as-judge, NOT experimental comparison
- **No system validates against published experimental data** — our strongest differentiator

### 2.3 LLM Agents for Scientific Discovery (Broader Context)
- Chemistry: ChemCrow (18 tools, Nature MI 2024), Coscientist (Nature 2023), MDCrow (40+ MD tools, ICLR 2025)
- Materials: MatSciAgent (2025), MechAgents (FEA, 2024)
- General: The AI Scientist (Sakana AI, 2024), ScienceAgentBench (102 tasks, best 32.4%)
- Surveys: arXiv:2503.08979, arXiv:2508.14111

## 3. System Design (1.5 pages)

### 3.1 Architecture Overview
- **Fig 1**: End-to-end pipeline: NL input → SloshingCoderPrompt → ReAct Agent Loop → 14 SPH tools + 12 pv-agent MCP tools → DualSPHysics Docker (CUDA 12.6, RTX 4090) → AI Analysis → Report
- **Design philosophy**: Sloshing domain drives architecture. The pipeline is not hard-coded; the system prompt guides tool selection based on sloshing physics (unlike MetaOpenFOAM's fixed role-based 4-agent pattern)
- **Single-agent + many tools** vs **multi-agent**: Simpler, more cost-effective, avoids inter-agent coordination overhead. Comparable to ChemCrow (18 tools, single ReAct) and MDCrow (40+ tools)

### 3.2 Tool Interface Design for Sloshing SPH
**14 DualSPHysics tools** addressing sloshing-specific challenges:

| Tool | Sloshing Pain Point Addressed |
|------|------|
| `xml_generator` | Eliminates hand-crafting of XML (200+ page guide) |
| `gencase` | Automates pre-processing (75-80% of CFD time) |
| `solver` | GPU-accelerated SPH execution (40-100x vs CPU) |
| `partvtk`, `measuretool` | Post-processing chain automation |
| `monitor` | Run.csv divergence detection (SPH-specific stability) |
| `error_recovery` | IsError pattern for LLM self-correction |
| `job_manager` | Async GPU jobs (SPH-specific: hours-long simulations) |
| `analysis` | AI-powered physics interpretation |
| `parametric_study` | Multi-case orchestration (200+ cases/tank) |
| `seismic_input` | Earthquake time-series parsing (Tomakomai scenario) |
| `stl_import` | CAD mesh → SPH particles (industrial geometry) |
| `result_store`, `report` | Result persistence and reporting |

**Key design patterns**:
- **IsError pattern**: SPH errors returned as `ToolResponse{IsError: true}` for LLM self-correction loop. Addresses DualSPHysics's 5 silent error types.
- **Run.csv monitoring**: Divergence detection (20% kinetic energy growth × 5 consecutive steps). SPH-specific — mesh-based solvers have different instability signatures.
- **Async GPU execution**: `context.Background()` for solver survival, max 3 concurrent. GPU-native SPH can run hours; agent must not block.
- **pv-agent MCP**: ParaView post-processing via Model Context Protocol — analogous to Foam-Agent 2.0's MCP architecture.

### 3.3 Domain-Specialized Sloshing Prompt
**SloshingCoderPrompt** (136 lines, 5 categories of sloshing-specific knowledge):

1. **Parameter inference** — `dp = min(L,W,H)/50`, `coefsound = 60`, viscosity `alpha` by fluid type
   → Addresses: beginner dp convergence failures, trial-and-error parameter tuning
2. **Tank presets** — LNG Mark III, automotive fuel tank, cylindrical storage, rectangular sloshing tank
   → Addresses: XML configuration from scratch (hours → seconds)
3. **Physics formulas** — `f1 = sqrt(g*pi/L * tanh(pi*h/L)) / (2*pi)`, SWL calculation, resonance detection
   → Addresses: natural frequency miscalculation, nonlinear spring behavior confusion
4. **Terminology mapping** — Korean/English CFD terminology disambiguation
   → Addresses: non-English-speaking engineers (Korea 70% of LNG shipbuilding)
5. **Docker path conventions** — `/cases/` for XML, `/data/` for output, lowercase binary names
   → Addresses: XML/binary path confusion (DualSPHysics Docker-specific)

**Contrast with existing approaches**:
- MetaOpenFOAM/ChatCFD: RAG retrieval of OpenFOAM documentation — general, not domain-specific
- OpenFOAMGPT 2.0: Prompt Pool — structured but not domain-knowledge-aware
- AutoCFD/FoamGPT: Fine-tuning on 28K+ labeled pairs — requires expensive dataset creation
- SloshAgent: Zero fine-tuning, prompt-only domain specialization — instant deployment

### 3.4 XML Generation Pipeline
- Template-based approach (336 lines Go code): structured tool input → DualSPHysics XML
- Auto-configuration: SWL gauge placement, mDBC/DBC switching, motion type selection
- Addresses the #1 DualSPHysics forum error: boundary particle exclusion from incorrect geometry setup

## 4. Experiments (2 pages)

### 4.1 Experimental Setup
- **Hardware**: NVIDIA RTX 4090 (24GB VRAM), CUDA 12.6
- **LLM**: Qwen3 32B via Ollama (local inference, zero cloud dependency, zero LLM cost)
- **Solver**: DualSPHysics v5.4 GPU in Docker (CUDA 12.6, lowercase symlinks)
- **Benchmarks**: SPHERIC Test 10 (FTP raw data, 100 repetitions, 2 fluids), 20 NL sloshing scenarios
- **Existing assets**: 20 XML cases, 17 simulation results (8 PASS, 2 PARTIAL), probe measurement files

### 4.2 EXP-1: SPHERIC Benchmark Reproduction (RQ2 — Validation Gap)
- **Goal**: Can agent-generated simulations match experimental sloshing data?
- NL description → XML generation → GPU simulation → probe measurement → comparison with SPHERIC Test 10
- **Metrics**: Pearson r, NRMSE, peak pressure error
- **Comparison**: agent results vs English2021 expert mDBC results vs experimental data (100-repeat statistics)
- **Fig/Table**: Pressure time series overlay (agent vs expert vs experimental ±2σ)
- **Significance**: First LLM-simulation system validated against published experimental benchmark data

### 4.3 EXP-2: NL→XML Generation Accuracy (RQ1 — Domain Gap)
- **Goal**: How accurately can the agent generate valid sloshing configurations from natural language?
- 20 scenarios × 5 complexity levels:
  - Level 1: "rectangular tank, water, 50% fill, horizontal sinusoidal motion"
  - Level 2: "LNG Mark III tank with 30% fill at resonance frequency"
  - Level 3: "Chen2018 case 3 with 40% fill level at f/f1 = 0.9"
  - Level 4: "cylindrical tank with anti-slosh baffle under seismic excitation"
  - Level 5: "STL-imported geometry with mDBC boundary and parametric fill sweep"
- **Metrics**: GenCase pass rate, parameter accuracy vs expert ground truth, physical validity
- **Comparison**: FoamGPT (26.36% on CFDLLMBench), AutoCFD (88.7%)

### 4.4 EXP-3: Parametric Study Automation (RQ3 — Industry Gap)
- **Goal**: Can the agent automate the 200+ case parametric studies required for LNG tank assessment?
- Reproduce Chen2018: 6 fill levels × multiple frequencies, automated from single NL prompt
- **Metrics**: Setup time (agent vs manual), result accuracy vs published data
  - Expert manual: ~2hr/case × 6 = 12hr
  - Agent: single NL prompt → 6 cases in minutes
- **Fig/Table**: Results overlay against Chen2018 published data
- **Significance**: Addresses the core LNG industry bottleneck (200+ cases/tank, 2-5 weeks/project)

### 4.5 EXP-4: Domain Prompt Ablation (RQ4 — Knowledge Gap)
- **Goal**: How much does sloshing-specific knowledge in the system prompt contribute to performance?
- Conditions: SloshingCoderPrompt ON vs OFF (generic coding prompt only)
- **Metrics**: XML generation accuracy delta, GenCase pass rate delta, physical validity delta
- **Ablation categories**: Remove sloshing physics formulas only / Remove XML syntax rules only / Remove parameter inference only / Remove all domain knowledge
- **Comparison with existing ablations**: MetaOpenFOAM (RAG removal), Foam-Agent (reviewer removal), MooseAgent (RAG removal) — all architectural, not knowledge ablation
- **Significance**: First domain prompt ablation for ANY computational mechanics domain

### 4.6 EXP-5: Industrial Scenario PoC (Contribution 5)
- Anti-slosh baffle design: automotive vertical baffle + NASA ring baffle reference
- Seismic excitation scenario: Tomakomai-type earthquake input via `seismic_input` tool
- **Fig/Table**: Force/pressure reduction with/without baffle, particle snapshots
- **Industry context**: Baffles reduce sloshing amplitude ~70%, wall pressure ~50%

## 5. Discussion (0.5 pages)

### 5.1 Limitations
- **Single GPU**: RTX 4090 particle count limited to ~500K (industrial cases may need millions)
- **Local LLM**: Qwen3 32B may underperform cloud models on complex reasoning (CFD-Copilot also uses Qwen3-32B)
- **Solver-specific**: DualSPHysics v5.4 only; not tested with SPHinXsys, GPUSPH, or mesh-based solvers
- **No user study**: non-expert accessibility claim is qualitative (target for future work)
- **Sloshing-specific**: intentionally not general-purpose CFD (a feature, not a limitation — domain depth over breadth)

### 5.2 From Sloshing to General Particle Simulation
- Tool interface patterns (IsError, async GPU, Run.csv monitoring) are transferable to other particle methods:
  - DEM (Discrete Element Method) — granular flow, powder processing
  - MPM (Material Point Method) — geomechanics, impact
  - SPHinXsys — multi-physics SPH
- The "particle-solver + LLM agent" paradigm applies wherever Lagrangian methods are used

### 5.3 Industry Impact Projection
- **LNG**: 200-case evaluation from months → days. Per-case setup: days → minutes. At 400+ carriers on order, total addressable market is significant.
- **Cost reduction**: BCG (2025) projects 10-20% R&D time-to-market reduction, up to 20% cost reduction from AI in R&D.
- **Open-source advantage**: DualSPHysics (open) + Qwen3 32B (open-weight) + SloshAgent = fully reproducible. Unlike cloud-dependent competitors (MetaOpenFOAM/GPT-4o, ChatCFD/DeepSeek API).

## 6. Conclusion (0.5 pages)
- Sloshing causes real-world disasters ($15B+ cumulative damage documented) and requires massive simulation effort (200+ cases/tank)
- No AI system previously automated sloshing simulation — or any particle-based simulation
- SloshAgent fills both gaps: domain-specialized sloshing automation + first particle-solver LLM agent
- SPHERIC benchmark validation: first LLM-simulation system validated against experimental data
- Domain prompt ablation: first knowledge ablation study for computational mechanics
- Future: multi-solver support (SPHinXsys), user study with naval architects, Qwen3 32B vs 8B comparison, LNG industry pilot

## References
[See references.bib — key papers organized by category]

**Sloshing Domain**:
- DualSPHysics (Dominguez et al., 2022, 485 citations)
- English et al. 2021 — mDBC sloshing validation
- SPHERIC benchmarks (FTP data)
- DesignSPHysics (Vieira et al., 2017)
- Pasimodo+RAG (arXiv:2502.03916, 2025) — 0/2 on model creation
- SNU sloshing DB (Kim et al., 2019, ScienceDirect)
- ISOPE 2012 sloshing benchmark (9 institutions)
- Chen et al. 2018 — OpenFOAM parametric sloshing
- Liu et al. 2024 — pitch sloshing 3×5×3

**LLM+CFD Agents** (primary contrast):
- MetaOpenFOAM 1.0/2.0 (Chen et al., 2024-25)
- Foam-Agent 2.0 (Yue et al., NeurIPS ML4PS 2025)
- ChatCFD (Fan et al., 2026)
- OpenFOAMGPT 2.0 (Pandey et al., 2025) — Claude-3.7-Sonnet
- FoamGPT (Yue et al., NeurIPS ML4PS 2025)
- CFD-Copilot (2025) — Qwen3-8B LoRA + Qwen3-32B
- AutoCFD (Dong et al., 2025)
- PhyNiKCE (PolyU, 2026) — neurosymbolic

**FEM/MD Agents**:
- MooseAgent (Zhang et al., 2025) — DeepSeek-R1+V3
- MCP-SIM (KAIST, npj AI 2025) — FEniCS
- MDCrow (2025) — 40+ tools

**AI4Science Foundations**:
- ChemCrow (Bran et al., Nature MI 2024)
- Coscientist (Boiko et al., Nature 2023)
- The AI Scientist (Sakana AI, 2024)
- ScienceAgentBench (ICLR 2025)
- CFDLLMBench (NeurIPS 2025)

**Industry & Cost Data**:
- Cadence 2024 — 75-80% pre-processing time
- CadCrowd 2024 — CFD consulting rates
- BCG 2025 — AI R&D impact (10-20% time, 20% cost reduction)
- Rescale 2025 — institutional knowledge loss
- Crespo et al. 2011 (PLOS ONE) — GPU 40-100x speedup

**Sloshing Accidents** (for Introduction):
- SpaceX Falcon 1 Flight 2 (2007)
- Tomakomai oil fire (2003) — Hatayama et al. 2004
- Alaska Good Friday (1964) — USGS
- Fukushima SFP (2011) — NRC/IAEA
- Tanker truck rollovers — FMCSA statistics

**LLM & Surveys**:
- arXiv:2503.08979 (Agentic AI for Scientific Discovery)
- arXiv:2508.14111 (From AI for Science to Agentic Science)
- arXiv:2504.02888 (Qwen vs GPT-4 for CFD)
