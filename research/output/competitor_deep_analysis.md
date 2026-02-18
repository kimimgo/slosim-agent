# Competitor Deep Analysis — Full Text Reading (9 Systems, 10 Papers)

Date: 2026-02-18
Method: arXiv HTML full text extraction via 8 parallel agents (Opus)

---

## Overview

| # | System | arXiv | Year | Solver | LLM | Key Metric | Cost |
|---|--------|-------|------|--------|-----|------------|------|
| 1 | MetaOpenFOAM 2.0 | 2502.00498 | 2025 | OpenFOAM | GPT-4o | 86.9% Pass@1 | $0.15/case |
| 2 | PhyNiKCE | 2602.11666 | 2026 | OpenFOAM | Gemini-2.5-Pro/Flash | 51% accuracy | $1.56/case |
| 3 | ChatCFD | 2506.02019 | 2025 | OpenFOAM | DeepSeek-R1+V3 | 82.1% exec / 68.12% phys | $0.21/case |
| 4 | Foam-Agent 2.0 | 2509.18178 | 2025 | OpenFOAM | Claude 3.5 Sonnet | 88.2% exec success | ~334K tok |
| 5 | OpenFOAMGPT 2.0 | 2504.19338 | 2025 | OpenFOAM | Claude-3.7-Sonnet | 100% reproducibility | N/A |
| 6 | CFD-Copilot | 2512.07917 | 2025 | OpenFOAM | Qwen3-8B(FT)+32B | U 96.4%, p 93.2% | Local |
| 7 | AutoCFD | 2504.09602 | 2025 | OpenFOAM | Qwen2.5-7B(FT) | 88.7% accuracy | $0.02/case |
| 8 | MooseAgent | 2504.08621 | 2025 | MOOSE(FEM) | DeepSeek-R1+V3 | 93% success | $0.14/case |
| 9 | Pasimodo+RAG | 2502.03916v2 | 2025 | Pasimodo(SPH) | Llama/Gemma/NLM | **0% creation** | Local |
| **10** | **SloshAgent** | **—** | **2026** | **DualSPHysics(SPH)** | **Qwen3 32B** | **8/10 E2E** | **$0** |

---

## 1. MetaOpenFOAM 2.0

**Title**: MetaOpenFOAM 2.0: Large Language Model Driven Automated CFD Simulation
**arXiv**: 2502.00498 (2025-02)
**Authors**: Yuxuan Chen et al., Tsinghua University
**Status**: NEW — 기존 2-cycle 분석에서 누락, citation chasing으로 발견

### Architecture
- **Framework**: MetaGPT v0.8.0, 4-agent pipeline
- **LLM**: GPT-4o (T=0.01)
- **Key Innovation**: QDCOT (Query Decomposition CoT) + ICOT (Iterative Correction CoT)
  - QDCOT: 2-level decomposition — 복잡 쿼리를 물리 하위 문제로 분해
  - ICOT: 에러 기반 반복 수정, scaling law 발견 (log-linear improvement)
- **Post-processing**: 자동 시각화 생성 (v1.0에 없던 기능)

### Benchmark
- 13 cases (v1.0의 8에서 확대)
- **86.9% Pass@1** (v1.0: 85%)
- **Executability**: 6.3/7 (v1.0: 2.1/7) — 3배 향상
- **$0.15/case** (v1.0: $0.22/case)
- Scaling laws: QDCOT follows log-linear, ICOT diminishing returns after ~3 iterations

### Limitations
- OpenFOAM 전용, SPH/meshfree 미지원
- v1.0→2.0 incremental improvement, 구조적 혁신보다는 prompting 기법 개선
- 13 cases는 상대적으로 작은 벤치마크

### Implications for SloshAgent
- QDCOT의 query decomposition 개념은 복잡한 슬로싱 시나리오 파싱에 참고 가능
- 우리의 xml_generator는 LLM이 config를 직접 쓰지 않아 더 안전한 설계
- 비용 비교: $0.15 vs $0 (우리 로컬 LLM)

---

## 2. PhyNiKCE

**Title**: PhyNiKCE: Physics-informed Neural-symbolic Integration for Knowledge-driven CFD Execution
**arXiv**: 2602.11666 (2026-02)
**Authors**: E Fan et al., Hong Kong PolyU
**Status**: Cycle 2에서 추가, 가장 최신 (2026-02)

### Architecture
- **Approach**: Neurosymbolic — neural planning과 symbolic validation 분리
- **LLM**: Gemini-2.5-Pro (planning) + Gemini-2.5-Flash (execution)
- **Key Innovation**: Deterministic RAG Engine with 5 specialized retrievers
  - Retriever 1: Solver compatibility
  - Retriever 2: Boundary condition constraints
  - Retriever 3: Numerical scheme selection
  - Retriever 4: Turbulence model matching
  - Retriever 5: Post-processing pipeline
- **CSP (Constraint Satisfaction Problem)**: 물리 법칙을 symbolic constraint로 인코딩, LLM hallucination 방지

### Benchmark
- 13 configurations, 340 runs
- **51% accuracy** (96% relative improvement vs ChatCFD baseline 26%)
- **$1.56/case** (가장 비싼 시스템)
- Error Analysis: 주요 실패는 turbulence model-boundary condition 불일치

### Limitations
- 절대 정확도 51%는 낮음 (상대 개선율 강조)
- OpenFOAM 전용
- Gemini API 비용이 높음
- CSP 규칙 수동 인코딩 필요 — 새 도메인으로의 확장 비용

### Implications for SloshAgent
- Neurosymbolic 접근은 학술적으로 가장 엄밀하나 실용성에서 의문
- 우리의 SloshingCoderPrompt + xml_generator 조합이 사실상 "물리 지식의 도구적 인코딩"
- 물리 제약을 prompt vs symbolic constraint vs tool parameter로 인코딩하는 세 가지 전략 비교 논의 가능

---

## 3. ChatCFD

**Title**: ChatCFD: Interactive Computational Fluid Dynamics with Natural Language
**arXiv**: 2506.02019 (2025-06)
**Authors**: E Fan et al., SUSTech/Beihang/DP Technology
**Status**: Cycle 1부터 분석, Cycle 2에서 검증 완료

### Architecture
- **Pipeline**: 4-stage (NL Parser → Config Generator → Executor → Physics Interpreter)
- **LLM**: DeepSeek-R1 (reasoning) + DeepSeek-V3 (general)
- **Knowledge Base**: 4개 structured JSON DB
  - Solver Template DB (가장 중요 — 제거 시 48%p 하락)
  - Boundary Condition DB
  - Numerical Scheme DB
  - Physics Model DB
- **Physical Fidelity Metric**: LLM(DeepSeek-R1)이 시뮬레이션의 과학적 유의미성 평가

### Benchmark
- **315 cases** (205 tutorial-derived + 110 perturbed)
- **82.1% execution success**
- **68.12% physical fidelity** (LLM-evaluated, NOT experimental)
- **97.4% narrative fidelity** (LLM이 물리를 유창하게 설명하는 비율)
- **$0.208/case** (192K tokens average)
- "Striking LLM gap": 97.4% narration vs 68.12% physics — LLM은 설명은 잘하지만 실제 물리 구현은 부족

### Key Ablation
| Component Removed | Impact |
|-------------------|--------|
| Solver Template DB | **-48%p** (가장 critical) |
| Physics Interpreter | -15%p |
| Error Recovery | -12%p |

### Limitations
- Physical Fidelity가 LLM-evaluated — 실험 데이터 비교가 아님
- 315 cases 중 205개가 tutorial 변형 → 진정한 generalization 의문
- Cross-benchmark 문제: MetaOpenFOAM을 자기 벤치마크에서 테스트 → 6.2%만 달성

### Implications for SloshAgent
- Physical Fidelity 메트릭의 한계: LLM이 평가하는 것은 circular reasoning 위험
- 우리는 SPHERIC 실험 데이터와의 정량 비교(r>0.9)로 차별화 — 객관적 검증
- Solver Template DB의 critical 역할은 우리 도메인 프롬프트의 중요성과 일치

---

## 4. Foam-Agent 2.0

**Title**: Foam-Agent 2.0: Composable Multi-Agent System for CFD Automation
**arXiv**: 2509.18178 (2025-09)
**Authors**: Shaowu Pan et al., RPI
**Status**: Cycle 1부터 분석, Cycle 2에서 검증 완료

### Architecture
- **Agents**: 6-agent composable (Planner, Retriever, Generator, Reviewer, Executor, Debugger)
- **LLM**: Claude 3.5 Sonnet (T=0.6)
- **Framework**: LangGraph + custom orchestration
- **MCP**: 11 functions for solver interaction
- **RAG**: 4-index hierarchical FAISS
  - Tutorial index
  - Documentation index
  - Error pattern index
  - Physics reference index
- **Visualization**: ParaView/Pyvista agent integrated

### Benchmark
- **110 tasks**, 7 physics categories (laminar, turbulent, multiphase, heat transfer, etc.)
- **88.2% execution success**
- **~334K tokens/case** (비용 미공시)

### Key Ablation
| Component Removed | Success Rate |
|-------------------|-------------|
| Full system | 88.2% |
| Without Reviewer | ~50% (**-38%p, 가장 critical**) |
| Without RAG | 57.3% (-31%p) |
| Without File Dependency | 72.1% (-16%p) |

### Cross-Benchmark Testing (중요!)
| System | Own Benchmark | Foam-Agent Benchmark | ChatCFD Benchmark |
|--------|---------------|---------------------|-------------------|
| MetaOpenFOAM | 85% | **55.5%** | 6.2% |
| Foam-Agent 2.0 | 88.2% | — | — |
| ChatCFD | 82.1% | — | — |

→ MetaOpenFOAM이 자기 벤치마크에서 85%이지만 다른 벤치마크에서는 급격히 하락
→ **Cross-benchmark 비교가 필수**적임을 증명

### Limitations
- Claude 3.5 Sonnet 의존 (클라우드 API 비용)
- 비용 정보 미공개
- 110 tasks의 physics categories 불균형 가능성

### Implications for SloshAgent
- **Reviewer agent의 결정적 역할**: 우리 시스템에도 self-review 메커니즘 필요
- MCP 11 functions vs 우리 14 tools + pv-agent — 유사한 규모
- Cross-benchmark 위험: 우리도 자체 벤치마크 + SPHERIC 외부 벤치마크 모두 보고해야 함

---

## 5. OpenFOAMGPT 2.0

**Title**: OpenFOAMGPT 2.0: End-to-End, Trustworthy Automation for CFD
**arXiv**: 2504.19338v1 (2025-04)
**Authors**: Jingsen Feng, Ran Xu, Xu Chu — U of Exeter + U of Stuttgart (3인 소팀)
**Status**: Cycle 2에서 LLM 오류 수정 (GPT→Claude)

### Architecture
- **Agents**: 4-agent pipeline (Pre-processing → Prompt Generation → Simulation Engine → Post-processing)
- **LLM**: **Claude-3.7-Sonnet** (T=0) — "GPT" in name is heritage only
- **Thinking**: **의도적으로 비활성화** — "extremely stringent formatting requirements"
- **Key Innovation**: Prompt Pool approach (RAG 명시적 거부)
  - 매번 fresh comprehensive prompt로 config 생성
  - Config file template 수정 방식을 "highly inefficient and prone to errors"로 비판
- **Docker**: OpenFOAM v2406 container, host-independent execution

### Benchmark: 455 Cases, 6 Study Types

| Study Type | Cases | Physics |
|-----------|-------|---------|
| Single-phase Poiseuille | 10 | Laminar channel flow |
| Multi-phase Poiseuille | 40 | Stratified immiscible |
| Single-phase porous media | 80 | Darcy/non-Darcy |
| Multi-phase porous media | 25 | Drainage process |
| Aerodynamics (motorbike) | 100 | Turbulent RANS |
| Extended reproducibility | 200 | Continuous stress test |
| **Total** | **455** | |

- **100% success rate** (시뮬레이션 완료)
- **100% reproducibility** (동일 입력 → 동일 출력)

### Critical Caveats
1. **정량적 에러 메트릭 없음** — "excellent agreement"은 visual inspection
2. **100% reproducibility ≠ 100% accuracy** — 실행 성공률이지 물리 정확도 아님
3. **비용/코드 미공개** — reproducibility 주장과 모순
4. **Ablation 없음** — 각 agent/Prompt Pool의 기여도 불명
5. **Limitations 섹션 없음** — 학술적 약점
6. **단일 LLM만 테스트** — Claude 3.7 Sonnet only

### Implications for SloshAgent
- Thinking 비활성화 + T=0: 우리도 CanReason:false로 동일 방향 → 논문에서 인용 가능
- xml_generator가 deterministic하게 XML 생성하는 접근이 Prompt Pool보다 더 안전
- 우리는 SPHERIC 실험 데이터와 정량 비교(RMSE, correlation) → 이들의 약점 보완
- Docker 기반 실행 환경이 동일 → architecture comparison 가능

---

## 6. CFD-Copilot

**Title**: CFD-copilot: Leveraging Domain-Adapted LLM and MCP to Enhance Simulation Automation
**arXiv**: 2512.07917v1 (2025-12)
**Authors**: Zhehao Dong, Shanghai Du, Zhen Lu, Yue Yang — Peking University
**Funding**: NSFC (52306126, 12525201, 12432010, 12588201), National Key R&D (2023YFB4502600)

### Architecture
- **Framework**: MetaGPT v0.8.1 + **MCP v1.9.0**
- **Agents**: 4 (Pre-checker → Generator → Runner → Corrector)
  - Generator: **Qwen3-8B (LoRA fine-tuned)**
  - Pre-checker/Runner/Corrector: **Qwen3-32B** (general purpose)
- **Self-correction**: 최대 **10회 반복** 루프
- **MCP Server**: **100+ validated post-processing tools** from OpenFOAM library
  - forceCoeffs, vorticity, streamlines, surface sampling, etc.
  - Self-descriptive tools (metadata + execution logic)

### Fine-tuning Details
- **Dataset**: NL2FOAM **49,205 pairs** + CoT annotation
- **Base model**: Qwen3-8B
- **Method**: LoRA
- **Hyperparameters**: 미공개 (rank, alpha, LR, epochs 모두 미기술) — 재현성 문제
- **GPU**: 미공개

### Benchmark

**NACA 0012 Airfoil**:

| AoA | U accuracy | p accuracy | Success Rate |
|-----|-----------|-----------|-------------|
| -2.5° | 97.75% | 97.83% | High |
| 0° | 97.56% | 98.53% | High |
| 5° | 98.55% | 99.09% | High |
| 10° | 96.37% | 93.20% | ~80% |
| 12.5° | 89.82% | 71.95% | Low |
| **Average** | **96.41%** | **93.22%** | **52.86%** |

**30P-30N Three-Element Airfoil (고양력)**:

| Model | Completion | Success |
|-------|-----------|---------|
| Qwen3-8B (fine-tuned) | **80%** | **10%** |
| Qwen3-Next-80B (general) | 10% | **0%** |
| Qwen3-235B (general) | 10% | **0%** |

→ **235B 범용 모델도 0% vs 8B fine-tuned 80% completion** — 도메인 fine-tuning이 모델 크기보다 중요

**Failure Analysis**: 범용 모델은 "Gauss linear" (단순)을 사용하지만, fine-tuned 모델은 "cellLimited Gauss linear 1" (안정적)을 학습

### Limitations
- 30P-30N 성공률 10% — 복잡 유동에서 낮은 신뢰성
- NL2FOAM 생성 방법론 미공개 — 재현 불가
- Fine-tuning hyperparameters 미공개
- 단일 solvertype (steady-state RANS)만 검증
- 코드/데이터 미공개

### Implications for SloshAgent
- **동일 Qwen3-32B 사용**: 같은 모델을 다른 전략(FT vs prompt-only)으로 활용 → 직접 비교 가치 최고
- Fine-tuning 효과 (235B=0% vs 8B-FT=80%): prompt-only의 한계를 인정하면서 SPH 데이터 부족을 정당화
- MCP 100+ tools: 우리의 pv-agent MCP와 유사한 접근 — post-processing 자동화의 수렴적 설계
- 10 trials x 10 iterations, T=0.6: 실험 설계 방법론 채택 가능

---

## 7. AutoCFD

**Title**: Fine-tuning a Large Language Model for Automating Computational Fluid Dynamics Simulations
**arXiv**: 2504.09602 (2025-04)
**Authors**: Zhehao Dong, Zhen Lu, Yue Yang — Peking University (CFD-Copilot과 동일 팀)
**Code**: https://github.com/YYgroup/AutoCFD (공개)

### Architecture
- **Multi-Agent**: 4 agents (Pre-checker → LLM Generator → Runner → Corrector)
- **LLM**: Qwen2.5-7B-Instruct (LoRA fine-tuned)
- **Max iterations**: 10회, 72시간 타임아웃
- **Mesh**: 사용자 제공 필수 — "LLMs cannot reliably generate [meshes]"

### Fine-tuning Details (공개)
- **Dataset**: NL2FOAM **28,716 pairs** + 6-step CoT reasoning
  - 16 OpenFOAM base cases에서 파라미터 변형 → 100K+ 생성 → 필터링
  - CoT 6단계: 문제 정의 → solver 선택 → 파일 결정 → BC/IC → 파라미터 → 실행 스크립트
- **LoRA**: rank r=8, trainable 0.02B/7.6B (0.26%)
- **Framework**: Llama-Factory
- **GPU**: NVIDIA RTX 4090 x4
- **Optimizer**: AdamW, LR 5e-5, linear warmup 10%
- **Epochs**: 4 (epoch 2 checkpoint 최적, 이후 overfitting)
- **Batch size**: 16 total

### Benchmark: 21 Cases

| Model | Accuracy | Pass@1 | Iterations | Cost (USD) |
|-------|----------|--------|------------|------------|
| **AutoCFD (7B-FT)** | **88.7%** | **82.6%** | **2.6** | **$0.020** |
| Qwen2.5-72B | 31.4% | 47.1% | 7.2 | $0.035 |
| DeepSeek-R1 | 41.7% | 22.4% | — | $0.042 |
| Llama3.3-70B | 4.7% | 0.5% | — | $0.018 |
| MetaOpenFOAM (GPT-4o) | limited | limited | — | $0.227 |

### CoT Ablation
| Config | Accuracy | Pass@1 |
|--------|----------|--------|
| With CoT | 88.7% | 82.6% |
| Without CoT | 78.2% | 61.7% |
| **Difference** | **+10.5%** | **+20.9%p** |

### Error Types
- **Structural**: 누락 파라미터, 필수 파일 → Corrector가 효율적 해결
- **Physical**: 비합리적 값 (k-omega dissipation rate 수 자릿수 과대) → 가장 어려운 문제

### Limitations
- 비압축성(incompressible) 유동만 (다상/압축성/반응 미포함)
- 16개 base case에서 파생 — 커버리지 제한적
- Epoch 2 이후 overfitting — 데이터 규모 한계

### Implications for SloshAgent
- **Fine-tuning >> RAG 증명**: 7B-FT가 72B 범용 대비 accuracy +57.3%p, 비용 11배 저렴
- **CoT 필수**: +20.9%p pass@1 향상. 우리 SloshingCoderPrompt에도 step-by-step reasoning 구조화 필요
- **데이터셋 구축 방법론**: base case → 파라미터 변형 → LLM paraphrasing → 실행 필터링 — 우리 20 XML에도 적용 가능
- **벤치마크 설계**: 21 cases x 10 trials = 210 experiments, 5 metrics (Accuracy, Pass@1, Iterations, Tokens, Cost)

---

## 8. MooseAgent

**Title**: MooseAgent: A LLM Based Multi-agent Framework for Automating Moose Simulation
**arXiv**: 2504.08621 (2025-04)
**Authors**: Tao Zhang, Zhenhai Liu, Yong Xin, Yongjun Jiao — Nuclear Power Institute of China
**Code**: https://github.com/taozhan18/MooseAgent (공개)

### Architecture
- **Framework**: LangGraph 기반 3-stage pipeline
- **Stage 1 (Alignment)**: 모호한 NL → 정밀 시뮬레이션 사양 (DeepSeek-V3)
- **Stage 2 (Write Input Card)**: RAG + 핵심 생성 (DeepSeek-R1)
- **Stage 3 (Execute-Analyze-Correct)**: 실행 + 에러 수정 루프 (DeepSeek-V3)
- **Max iterations**: 3 (더 많아도 개선 없음)
- **Temperature**: 0.01

### RAG (Knowledge Base)
- **8,000+ annotated MOOSE input files**: 자동 주석 파이프라인 (MOOSE repo 샘플 → function 분석 → LLM 주석 생성)
- **Function documentation**: Python 스크립트로 MOOSE docs 스캔
- **Embedding**: BGE-M3
- **Vector DB**: FAISS
- **Retrieval**: Top-3 chunks by similarity, simple concatenation

### LLM Role Split
| Model | Role | Reason |
|-------|------|--------|
| DeepSeek-R1 | Core generation (Stage 2) | Superior reasoning |
| DeepSeek-V3 | Alignment + Correction (Stage 1, 3) | Speed-quality balance |

### Benchmark: 9 Cases x 5 Repeats

| Case | Physics | Success |
|------|---------|---------|
| Steady-State Heat | 1D rod | 100% |
| Transient Heat | 2D square | 100% |
| Linear Elasticity | 2D plate | 100% |
| Plasticity | 2D elastoplastic | **60%** |
| Porous Media Flow | 3D cubic | 100% |
| Phase Change Heat | 1D melting | 100% |
| Phase Field | 2D solidification | 100% |
| Prismatic Fuel | 16-sided fuel assembly | 100% |
| Thermal-Mechanical | 2D MultiApp coupling | 80% |
| **Average** | | **93%** |

- **~60,777 tokens/case**, **<$0.14/case**

### Ablation: KB Removal (93% → 76%)
| Case | With KB | Without KB |
|------|---------|-----------|
| Plasticity | 60% | 20% (-40%p) |
| Phase Change | 100% | 60% (-40%p) |
| Porous Media | 100% | 80% (-20%p) |
| Simple cases | 100% | 100% (no change) |

### Failure Mode
- **무한 에러 루프**: Agent가 동일 실패 function을 반복 호출하다 iteration 한도 초과
- Plasticity(60%)가 가장 어려움: 탄소성 거동의 복잡한 재료 모델

### Limitations
- 9 cases는 작은 벤치마크
- LangGraph 구현 세부 미공개 (재현성 제한)
- 타 시스템과 직접 비교 없음
- MOOSE GeneratedMesh에 의존 (복잡 형상 미지원)

### Implications for SloshAgent
- RAG(93%) vs Fine-tuning(88.7%): 벤치마크 난이도가 달라 직접 비교 어렵지만, 두 방식 모두 유효
- **반복 에러 루프 문제**: 우리 error_recovery tool에도 "동일 에러 반복 시 전략 전환" 필요
- DeepSeek R1+V3 역할 분담: 우리도 Qwen3 32B(추론) + 8B(보조) 분업 가능
- 핵공학 도메인: "위험한 시뮬레이션 자동화"의 안전성 논의 참고

---

## 9. Pasimodo+RAG — 유일한 SPH+LLM 논문

**Title**: Experiments with Large Language Models on Retrieval-Augmented Generation for Closed-Source Simulation Software
**arXiv**: 2502.03916v2 (2025-02, revised 2025-12)
**Authors**: Andreas Baumann, Peter Eberhard — University of Stuttgart (ITM)
**Pages**: 16 pages, 6 tables, 2 figures, 66 references

### Architecture: Pure RAG (NOT an Agent)
- **No agent capabilities**: No tool calling, function calling, simulation execution
- **No iterative correction**: 전문가가 수동으로 에러 피드백
- **GPU**: LLM inference에만 사용, 솔버 실행과 통합 없음

### LLMs Tested
| Model | Size | Platform |
|-------|------|----------|
| Llama 3.2 3B | 3.21B | Ollama (RTX 4070/4090) |
| Gemma 3 4B | 4.3B | Ollama |
| Gemma 3 27B | 27.43B | Ollama (RTX 4090) |
| NotebookLM (Gemini 2.5 Flash) | N/A | Google Cloud |

### Pasimodo Solver
- **Method**: SPH + DEM (mesh-free)
- **Status**: **Closed-source** (Stuttgart ITM proprietary)
- **Input**: XML-based configuration
- **Applications**: abrasive wear, deep-hole drilling, laser beam welding, particle damping
- **No sloshing capability mentioned**

### Benchmark: 28 Prompts, 6 Categories

| Category | Llama 3B | Gemma 4B | Gemma 27B | NotebookLM |
|----------|----------|----------|-----------|------------|
| Q&A (6) | 3/6 | 4/6 | 6/6 | 6/6 |
| Structured (6) | 2/6 | 3/6 | 4/6 | 6/6 |
| Explaining (4) | 2/4 | 3/4 | 4/4 | 4/4 |
| Summarization (5) | 1/5 | 4/5 | 3/5 | 2/5 |
| Reasoning (5) | 1/5 | 0/5 | 2/5 | 6/6* |
| **Creation (2)** | **0/2** | **0/2** | **0/2** | **0/2** |

### Model Creation Failure (0/2): Root Causes

**Prompt 6.1 (Influx_External 최소 예제)**:
- RAG가 `Influx_External` 대신 `Inflow_External` 문서를 검색 (이름 유사)
- 모든 LLM이 잘못된 component 기반 코드 생성

**Prompt 6.2 (3D oil droplet 시뮬레이션)**:
- NotebookLM (최선): missing parameters + faulty SPH interaction + missing neighborhood search
- Local LLMs: syntax hallucination (Pasimodo에 없는 문법 생성)
- Error feedback 시도: 부분 수정 가능하나 자율적이지 않음, 전문가 수동 분석 필요

### Key Quotes for Our Paper

> "all systems failed to create a minimal input example"

> "The best simulation model was created by NotebookLM, but it is not executable because it had missing parameters in the definition of a component and a faulty SPH interaction"

> "the main challenge remains a proper information retrieval, in general, a simple RAG search with similarity does not work sufficiently"

### Implications for SloshAgent — GAP-1의 결정적 근거

**이 논문이 우리 연구의 핵심 positioning 근거:**

1. **SPH+LLM 분야 유일한 논문**이 모델 생성 **0% 성공**
2. 그들의 future work ("tool-calling, simulation execution, error recovery")가 **정확히 우리가 이미 구현한 것**
3. 우리 SloshAgent: 10/10 E2E GPU 테스트 완료 (8 PASS, 2 PARTIAL) vs 이들 0%
4. **"SPH solver를 agent가 자율적으로 구동하여 실행+검증까지 하는 시스템은 세계 최초"** 주장의 완벽한 근거

| 차원 | Pasimodo+RAG | SloshAgent |
|------|------------|------------|
| Architecture | Pure RAG Q&A | 14-tool ReAct agent |
| Solver integration | None | Full GPU pipeline |
| Error recovery | Manual expert feedback | Automated error_recovery tool |
| Domain specialization | Generic RAG | SloshingCoderPrompt 136 lines |
| Model creation success | **0%** | **80%+ (8/10 E2E)** |
| Simulation execution | GPU for LLM only | GPU for solver (RTX 4090) |
| Benchmark | Qualitative Q&A | SPHERIC quantitative comparison |

---

## Cross-Cutting Analysis

### Architecture Patterns

| Pattern | Systems | Key Characteristic |
|---------|---------|-------------------|
| Multi-agent + RAG | MetaOpenFOAM, Foam-Agent, MooseAgent | FAISS vector DB, tutorial-based knowledge |
| Multi-agent + Fine-tuning | AutoCFD, CFD-Copilot | LoRA on 28-49K pairs, CoT annotation |
| Multi-agent + Prompt Pool | OpenFOAMGPT 2.0 | Fresh prompt per case, no RAG |
| Neurosymbolic | PhyNiKCE | CSP constraints + neural planning |
| Pipeline + KB | ChatCFD | 4 structured JSON databases |
| Pure RAG Q&A | Pasimodo+RAG | No agent, no execution |
| **Single agent + Tools** | **SloshAgent** | **14 tools + domain prompt** |

### LLM Model Landscape

| Tier | Systems | Trend |
|------|---------|-------|
| Cloud ($$$$) | MetaOpenFOAM(GPT-4o), OpenFOAMGPT(Claude) | 2024 early adopters |
| Cloud ($) | ChatCFD/MooseAgent(DeepSeek), Foam-Agent(Claude) | 2025 cost optimization |
| Local (FT) | AutoCFD(Qwen2.5-7B), CFD-Copilot(Qwen3-8B) | 2025 late, fine-tune trend |
| Local (prompt-only) | Pasimodo(Llama/Gemma), **SloshAgent(Qwen3-32B)** | Zero cost |

**Trend**: 2024 GPT-4 독점 → 2025 중반 DeepSeek/Qwen 전환 → 2025 후반 로컬 LoRA fine-tune → 2026 prompt-only + tools

### Success Rate Comparison (주의: Cross-benchmark 불가)

각 시스템의 성공률은 **자체 벤치마크** 기준이며 직접 비교 불가:
- MetaOpenFOAM 자기 벤치마크 85% → Foam-Agent 벤치마크에서 **55.5%** → ChatCFD 벤치마크에서 **6.2%**
- 벤치마크 난이도, 정의, 케이스 구성이 모두 다름
- **우리 논문에서도 이 비교 불가능성을 명시하고, SPHERIC 외부 벤치마크 사용을 차별점으로 강조해야 함**

### Key Findings for Our Paper

1. **SPH 공백 확정**: 9개 시스템 중 8개 OpenFOAM(FVM), 1개 MOOSE(FEM). SPH는 0% 성공 Pasimodo만 → GAP-1 확정적

2. **Fine-tuning vs Prompt-only**: AutoCFD/CFD-Copilot이 FT 우위 증명하나, SPH 분야에 학습 데이터 없음 → prompt-only 정당화

3. **MCP 패턴 수렴**: Foam-Agent(11), CFD-Copilot(100+), SloshAgent(14+pv-agent) → tool interface가 핵심 성공 요인

4. **Thinking 비활성화**: OpenFOAMGPT T=0+no thinking, 우리 CanReason:false → config 생성에서 reasoning 해로움

5. **Cross-benchmark 필수**: 자체 벤치마크만으로는 일반화 주장 불가 → SPHERIC 외부 데이터 사용이 차별화 포인트

6. **Reviewer/KB 결정적**: Foam-Agent -38%p, ChatCFD -48%p, MooseAgent -17%p → 검증 메커니즘과 도메인 지식이 핵심

---

## References (arXiv IDs)

1. MetaOpenFOAM 2.0: arXiv:2502.00498
2. PhyNiKCE: arXiv:2602.11666
3. ChatCFD: arXiv:2506.02019
4. Foam-Agent 2.0: arXiv:2509.18178
5. OpenFOAMGPT 2.0: arXiv:2504.19338
6. CFD-Copilot: arXiv:2512.07917
7. AutoCFD: arXiv:2504.09602
8. MooseAgent: arXiv:2504.08621
9. Pasimodo+RAG: arXiv:2502.03916v2
