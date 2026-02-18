# Competitor Analysis — Verified from Paper Full Text (Cycle 2)

Date: 2026-02-18
Source: arXiv PDF full text extraction via PyMuPDF

---

## Table 1 (Revised): LLM-based Simulation Systems — Verified Metrics

| System | Year | Solver | Architecture | LLM | Benchmark | Success Metric | Cost | Ablation |
|--------|------|--------|-------------|-----|-----------|---------------|------|---------|
| MetaOpenFOAM | 2024 | OpenFOAM 10 | 4-agent MetaGPT v0.8.0 | GPT-4o (T=0.01) | 8 self-built cases × n=10 | 85% avg Pass@1 (score 4/4, human-verified) | $0.22/case (44K tok avg) | RAG removal + temp sensitivity |
| Foam-Agent 2.0 | 2025 | OpenFOAM | 6-agent composable + MCP | Claude 3.5 Sonnet (T=0.6) | 110 tasks, 7 physics categories | 88.2% execution success | ~334K tok/case (cost N/A) | Reviewer + RAG + File Dependency |
| ChatCFD | 2026 | OpenFOAM | 4-stage pipeline + structured KB | DeepSeek-R1 + V3 (dual) | 315 cases (205 tutorial + 110 perturbed) | 82.1% exec / 68.12% physical fidelity | $0.208/case (192K tok) | Solver Template DB removal → 48% |
| OpenFOAMGPT 2.0 | 2025 | OpenFOAM v2406 Docker | 4-agent pipeline + Prompt Pool | **Claude-3.7-Sonnet** (T=0) | 455 cases (6 study types) | 100% reproducibility | Cloud (Claude API) | None |
| FoamGPT | 2025 | OpenFOAM | LoRA fine-tune | Qwen3-8B (LoRA) | CFDLLMBench | 26.36% execution success | Local | N/A |
| CFD-Copilot | 2025 | OpenFOAM v2406 | 4-agent MetaGPT v0.8.1 + MCP (100+ tools) | Qwen3-8B (LoRA, 49K pairs) + Qwen3-32B (T=0.6) | NACA 0012 (AoA -2.5°~12.5°) | U 96.4%, p 93.2% accuracy | Local | N/A |
| AutoCFD | 2025 | OpenFOAM | Fine-tune + multi-agent | Qwen2.5-7B | NL2FOAM 28.7K pairs | 88.7% accuracy | Local | N/A |
| MooseAgent | 2025 | MOOSE (FEM) | 3-part LangGraph, ~5 agents | DeepSeek-R1 + V3 (dual, T=0.01) | 9 cases × n=5 | 93% avg success | <$0.14/case (61K tok) | RAG removal: 93%→76% |
| Pasimodo+RAG | 2025 | Pasimodo (SPH+DEM, closed) | Pure RAG Q&A (NOT agent) | Llama 3.2 3B / Gemma 3 27B / NotebookLM | 28 prompts, 6 categories | **0/2 on model creation** | Local/Free | N/A |
| **SloshAgent** | **2026** | **DualSPHysics v5.4 GPU** | **14 tools + ReAct single agent** | **Qwen3 32B (local Ollama)** | **SPHERIC Test 10 + 20 NL scenarios** | **Target: r>0.9, 8/10 E2E pass** | **$0 LLM cost** | **Domain prompt ON/OFF** |

---

## Key Corrections from Paper Full Text

### 1. OpenFOAMGPT 2.0 — LLM is Claude, NOT GPT
- Paper explicitly states: "All intelligent agents within the framework are powered by the Claude-3.7-Sonnet"
- "GPT" in the product name is heritage naming only
- Temperature=0 explicitly for determinism
- NO RAG — uses "Prompt Pool" approach instead

### 2. MooseAgent — LLM is DeepSeek, NOT GPT-4
- DeepSeek-R1 for core input file generation (reasoning-intensive)
- DeepSeek-V3 for remaining modules (alignment, error analysis)
- Cost <1 yuan ≈ <$0.14 at DeepSeek pricing
- RAG: 8,000+ annotated MOOSE inputs + function docs (FAISS + BGE-M3)

### 3. CFD-Copilot — Much More Substantial Than Initially Reported
- MCP v1.9.0 with **100+ validated post-processing tools**
- Fine-tuned Qwen3-8B via LoRA on **49,205 NL2FOAM pairs** (chain-of-thought annotated)
- General agents use Qwen3-32B (same model as our system!)
- MetaGPT v0.8.1 framework
- NACA 0012 results: U accuracy 96.4%, p accuracy 93.2% (avg across AoA)

### 4. ChatCFD — Physical Fidelity Definition Clarified
- 68.12% = LLM (DeepSeek-R1 "Physics Interpreter") evaluates whether runnable simulation is "scientifically meaningful"
- 97.4% = summary text fidelity (LLM can narrate physics fluently)
- The "striking LLM gap": 97.4% narration vs 68.12% actual physics enforcement
- Self-described as "the first rigorous metric capturing whether a runnable simulation is scientifically meaningful"
- **Weakness**: Metric is LLM-evaluated, not experimental comparison

### 5. Foam-Agent 2.0 — Detailed Ablation Available
- Reviewer node is "most significant factor": without it, success drops to ~50%
- Hierarchical multi-index RAG: 57.3% → 88.2% (with reviewer)
- ParaView/Pyvista visualization agent integrated
- **MetaOpenFOAM on their benchmark**: only 55.5% (not 85% as MetaOpenFOAM self-reports on their own benchmark)
- ChatCFD reports "Foam-Agent" at 42.3% — but this is v1.x, not 2.0

### 6. Pasimodo+RAG — CRITICAL: Weakest Competitor
- **NOT an agent**: Pure RAG Q&A, no tool use, no simulation execution
- **ALL systems (including NotebookLM) scored 0/2 on model creation**
- Best result (NotebookLM): produced non-executable model with missing parameters and faulty SPH interaction
- Local LLMs (3B-27B) scored 0-2/5 on compositional reasoning
- **No sloshing capability mentioned**
- **GPU used only for LLM inference**, not simulation
- Code not open-source
- This MASSIVELY strengthens GAP-1: even the only SPH+LLM paper cannot execute simulations

---

## Competitive Position Summary

### What SloshAgent Uniquely Offers (No Competitor Has):
1. **GPU SPH execution** — RTX 4090, CUDA 12.6, DualSPHysics v5.4
2. **Sloshing domain specialization** — SloshingCoderPrompt 136 lines, 5 categories
3. **Free-surface flow simulation** from natural language
4. **SPHERIC benchmark validation** — experimental data comparison (r>0.9 target)
5. **STL import + seismic input** tools
6. **Real-time GPU monitoring** — Run.csv divergence detection
7. **Zero LLM cost** — fully local Qwen3 32B via Ollama

### Architecture Comparison:
| Pattern | MetaOpenFOAM | Foam-Agent | ChatCFD | CFD-Copilot | SloshAgent |
|---------|-------------|-----------|--------|------------|-----------|
| Multi-agent | 4 agents | 6 agents | 4-stage | 4 agents | Single agent |
| Framework | MetaGPT | Custom+MCP | Custom | MetaGPT+MCP | ReAct loop |
| RAG | FAISS | Hierarchical | Structured KB (4 JSON) | Fine-tuned | Domain prompt |
| MCP | No | Yes (exposure) | "Compatible" | Yes (100+ tools) | Yes (pv-agent) |
| Post-process | None | ParaView/Pyvista | Physics Interpreter | MCP 100+ tools | pv-agent MCP |
| Error loop | Max 20 iter | Reviewer agent | Reflection | Max 10 iter | IsError pattern |
| GPU | No | No | No | No | **Yes (CUDA)** |

### Cross-Benchmark Comparison Caveat:
Success rates are NOT directly comparable across papers:
- MetaOpenFOAM: 85% on own 8 cases, but 55.5% on Foam-Agent's benchmark, 6.2% on ChatCFD's benchmark
- Foam-Agent: 88.2% on own 110 cases
- ChatCFD: 82.1% on own 315 cases
- Each paper uses different benchmarks, success definitions, and evaluation criteria

---

## LLM Model Landscape (Verified)

| System | Primary LLM | Cost Tier | Local? |
|--------|------------|----------|--------|
| MetaOpenFOAM | GPT-4o | Cloud ($$$) | No |
| Foam-Agent 2.0 | Claude 3.5 Sonnet | Cloud ($$) | No |
| ChatCFD | DeepSeek-R1 + V3 | Cloud ($) | No |
| OpenFOAMGPT 2.0 | Claude-3.7-Sonnet | Cloud ($$) | No |
| FoamGPT | Qwen3-8B (LoRA) | Local | Yes |
| CFD-Copilot | Qwen3-8B (LoRA) + Qwen3-32B | Local | Yes |
| AutoCFD | Qwen2.5-7B | Local | Yes |
| MooseAgent | DeepSeek-R1 + V3 | Cloud ($) | No |
| Pasimodo+RAG | Llama/Gemma (3B-27B) | Local | Yes |
| **SloshAgent** | **Qwen3 32B** | **Local ($0)** | **Yes** |

Trend: 2025 후반부터 로컬/오픈웨이트(Qwen3, Gemma) 전환 가속. 2024 초 GPT-4 독점 → 2025 후반 DeepSeek/Qwen 이중 모델 → 2026 로컬 LoRA fine-tune 시대.
