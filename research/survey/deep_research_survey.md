# Deep Literature Survey: Domain-Specific AI Agent Architecture for Automated Sloshing Simulation

**Generated:** 2026-02-18
**Source:** MCP Survey (2,144편) + Deep Web Research (65+ sources)

---

## Very High Relevance Papers (8편)

### 1. OpenFOAMGPT 2.0 (April 2025)
- **arXiv:** 2504.19338
- End-to-end multi-agent framework (Pre-processing, Prompt Gen, Simulation, Post-processing)
- 100% success rate on 450+ simulations

### 2. Foam-Agent 2.0 (September 2025)
- **arXiv:** 2509.18178
- Claude 3.5 Sonnet 기반 88.2% success rate
- HPC deployment capability

### 3. ChatCFD (2025-2026)
- **arXiv:** 2506.02019
- DeepSeek-R1/V3 기반 82.1% execution success, **68.12% physics fidelity** (최초 물리 메트릭)
- 315 benchmark cases

### 4. NL2FOAM (2025)
- **ScienceDirect:** S2095034925000261
- Qwen2.5-7B fine-tuned on 28,716 CFD cases with CoT
- 88.7% accuracy, 82.6% pass@1
- CoT ablation: +10.5% accuracy, +20.9% pass@1

### 5. SimLLM (January 2025)
- **arXiv:** 2601.06543
- Multi-stage SFT+DPO for SimPy simulation
- Qwen2.5-Coder-7B: 80.4% → 86.0% executability

### 6. MCP Specification (Anthropic, November 2024)
- **GitHub:** modelcontextprotocol
- LLM-tool integration open standard (JSON-RPC 2.0)
- OpenAI, Google DeepMind 등 채택

### 7. SPH Method for Long-Time LNG Sloshing (Pilloton et al., 2022)
- **Journal:** European Journal of Mechanics - B/Fluids
- 3-hour real-time simulation with severe sea-state forcing

### 8. ReAct Framework (Yao et al., 2022)
- Interleaved reasoning + action framework
- 에이전트 아키텍처 기반

---

## High Relevance Papers (23편)

### LLM + CFD
| Paper | Year | Key |
|-------|------|-----|
| OpenFOAMGPT | 2025 | RAG + GPT-4o, iterative correction |
| MetaOpenFOAM | 2024 | MetaGPT assembly line paradigm |
| Engineering.ai | 2025 | 100% success, Gmsh mesh gen |
| Coding Agents for CFD | 2026 | Lightweight, no multi-agent overhead |
| ChatCFD Knowledge Base | 2025 | JSON DB from OpenFOAM manuals |

### Tool-Augmented Agents
| Paper | Year | Key |
|-------|------|-----|
| AstaBench (Allen AI) | 2025 | 57 agents, 22 architectures |
| Multi-Dimensional Agent Eval | 2025 | ReAct-GPT4 best in SW dev (73.3%) |
| ChemCrow | 2023 | 13 domain-specific tools |

### Sloshing / DualSPHysics
| Paper | Year | Key |
|-------|------|-----|
| SPHERIC Xi'an | 2022 | 3D multi-depth, stochastic pressure |
| Liquid Sloshing Elastic Tank | 2024 | FSI + multi-phase SPH |
| SPH Validation (MDPI) | 2019 | Prismatic tank, <1% error |
| State-of-the-art DualSPHysics | 2021 | Multi-phase, GPU/CPU |
| DualSPHysics GPU | Wiki | 100x speedup, CUDA 11.7 |
| SPHERIC 2024 Berlin | 2024 | AI+CFD convergence keynote |

### Local SLM
| Paper | Year | Key |
|-------|------|-----|
| SLM Landscape (HuggingFace) | 2024 | Gartner: 3x SLM > LLM by 2027 |
| SLM Fine-Tuning Benchmark | 2024 | Qwen3-4B = GPT-OSS-120B |
| Qwen2.5-Coder Report | 2024 | 73.7 Aider (≈GPT-4o) |

### Democratization
| Source | Key |
|--------|-----|
| NAFEMS | "Not simply deployment" — expert knowledge capture |
| EASA SPDM | <5% global adoption |

### Multi-Agent Frameworks
| Framework | Key |
|-----------|-----|
| AutoGen | 40% faster agent comms |
| CrewAI | 6x faster execution |
| LangGraph | Graph-based, complex workflows |

---

## Research Gap Analysis

| Gap | Current State | Our Contribution |
|-----|---------------|-----------------|
| DualSPHysics + LLM | **0 papers** | First DualSPHysics agent |
| Local SLM for CFD | Central deployment only | Local SLM (field engineer laptop) |
| Sloshing automation | General CFD only | Sloshing-specific (LNG/fuel tank) |
| MCP for scientific | IDE/DevOps only | MCP for CFD workflow |
| Physics metrics | ChatCFD 68.12% | Sloshing-specific validation |

---

## Key arXiv IDs for Citation
```
2504.19338  OpenFOAMGPT 2.0
2509.18178  Foam-Agent 2.0
2506.02019  ChatCFD
2501.06327  OpenFOAMGPT
2407.21320  MetaOpenFOAM
2601.06543  SimLLM
2602.11689  Coding Agents for CFD
2511.00122  Engineering.ai
2409.16226  Sloshing Elastic Tank
2104.00537  State-of-the-art DualSPHysics
2409.12186  Qwen2.5-Coder
2602.12049  HPC Code Gen with RL
```
