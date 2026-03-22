# Research Gap Analysis Report

**Topic:** Domain-Specific AI Agent Architecture for Automated Sloshing Simulation
**Date:** 2026-02-18
**Papers Analyzed:** 150 (survey) + 31 (deep research)

---

## Gap 1: No LLM Agent for Particle-Based SPH Solvers
- **Severity:** 9/10 (Critical)
- **Evidence:** All 7 existing LLM-CFD works (OpenFOAMGPT, OpenFOAMGPT 2.0, Foam-Agent 2.0, ChatCFD, NL2FOAM, MetaOpenFOAM, Engineering.ai) exclusively target OpenFOAM or other FVM solvers.
- **Zero** publications address DualSPHysics or any SPH solver.
- **Our contribution:** SlosimAgent — first AI agent for DualSPHysics with 18 specialized tools.

## Gap 2: No Local SLM Deployment for CFD Automation
- **Severity:** 8/10 (High)
- **Evidence:** All existing approaches rely on cloud-based proprietary models:
  - OpenFOAMGPT 2.0: GPT-4o (OpenAI API)
  - Foam-Agent 2.0: Claude 3.5 Sonnet (Anthropic API)
  - ChatCFD: DeepSeek-R1/V3 (DeepSeek API)
  - NL2FOAM: Qwen2.5-7B (cloud fine-tuning, not local deployment)
- **Our contribution:** Qwen3-32B via Ollama, fully local inference on RTX 4090.

## Gap 3: No Sloshing-Specific AI Agent
- **Severity:** 9/10 (Critical)
- **Evidence:** No AI agent targets sloshing simulation specifically. Existing agents handle general CFD (turbulence, heat transfer, external aerodynamics).
- **Unique requirements:** Resonance frequency (f₁), fill-level dependent dynamics, excitation characterization, standard tank configurations (LNG/ship/lab).
- **Our contribution:** Sloshing domain prompts with auto-parameter rules, resonance calculation, 4 standard tank configs.

## Gap 4: No MCP-Based Tool Integration for Scientific Computing
- **Severity:** 7/10 (Medium-High)
- **Evidence:** MCP (Anthropic, Nov 2024) adopted by OpenAI, Google DeepMind, but exclusively for IDE/DevOps tooling. No scientific computing application found.
- **Our contribution:** 18 MCP tools for DualSPHysics workflow (XML generation → solver → post-processing → visualization → reporting).

## Gap 5: No SPH Physics Fidelity Metrics
- **Severity:** 7/10 (Medium-High)
- **Evidence:** ChatCFD introduced physics fidelity metric (68.12%) but only for FVM/OpenFOAM. No equivalent metrics for particle-based SPH simulations.
- **SPH-specific challenges:** Free surface tracking accuracy, particle distribution uniformity, conservation properties, pressure oscillation filtering.
- **Our contribution:** Proposed metrics — free surface NRMSE, wall pressure fidelity, particle conservation, energy conservation.

---

## Gap Interconnection Matrix

| | Gap 1 (SPH) | Gap 2 (Local) | Gap 3 (Sloshing) | Gap 4 (MCP) | Gap 5 (Metrics) |
|---|:---:|:---:|:---:|:---:|:---:|
| Gap 1 | — | ○ | ● | ● | ● |
| Gap 2 | ○ | — | ○ | ○ | ○ |
| Gap 3 | ● | ○ | — | ● | ● |
| Gap 4 | ● | ○ | ● | — | ○ |
| Gap 5 | ● | ○ | ● | ○ | — |

● = Strongly connected, ○ = Weakly connected

---

## Key Insight

Gaps 1, 3, 4 form a tightly coupled cluster: building the first SPH agent (Gap 1) requires sloshing domain knowledge (Gap 3) and tool integration (Gap 4). Gap 2 (local deployment) is independently addressable but amplifies practical impact. Gap 5 (metrics) enables rigorous evaluation of all other contributions.
