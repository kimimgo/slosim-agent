# Research Plan: LLM Agent Framework for Autonomous SPH-based Sloshing Simulation

## Target Venue
NeurIPS/ICML AI4Science Workshop (4-8 pages, English)

## Deadline
2026-04-15 (ASAP)

## Novelty Claim
First domain-specialized LLM agent system that autonomously configures, executes, and validates SPH sloshing simulations through structured tool interfaces, achieving expert-level accuracy on SPHERIC benchmarks with a local open-weight model.

## Key Differentiators (vs Competitors)

| Our System | Competitor | Differentiator |
|-----------|-----------|---------------|
| **Lagrangian particle-based** (DualSPHysics GPU) | All 10+ systems (MetaOpenFOAM~PhyNiKCE) | All mesh-based OpenFOAM; entire particle column empty |
| **Autonomous agent** (14 tools + ReAct) | Pasimodo+RAG (0/2 model creation) | RAG = Q&A assistance, not autonomous execution |
| **Experimental validation** (SPHERIC Test 10, r>0.9) | All (execution success only) | First NL→experimentally-validated physics pipeline |
| **Knowledge ablation** (136-line domain prompt) | All (architectural ablation only) | First domain-knowledge ablation for simulation |
| **Zero fine-tuning** (Qwen3 32B prompt-only) | AutoCFD/FoamGPT/CFD-Copilot (28K+ pairs) | Prompt engineering vs expensive dataset creation |
| **GPU-native execution** (RTX 4090, CUDA 12.6) | All OpenFOAM systems (CPU) | SPH is inherently GPU-parallel, $0 LLM + consumer GPU |

## Competitive Landscape (11 direct competitors — Verified Cycle 2)

| # | System | Solver | LLM (Verified) | Best Metric | Weakness vs Ours |
|---|--------|--------|----------------|------------|-----------------|
| C1 | MetaOpenFOAM | OpenFOAM 10 | GPT-4o (T=0.01) | 85% avg Pass@1, $0.22 | Mesh-based, cloud LLM, 8 self-built cases |
| C2 | OpenFOAMGPT 2.0 | OpenFOAM v2406 | **Claude-3.7-Sonnet** | 100% repro (455 cases) | Cloud LLM, academic cases, no ablation |
| C3 | Foam-Agent 2.0 | OpenFOAM | Claude 3.5 Sonnet | 88.2% exec (110 tasks) | Execution only, no physical validation |
| C4 | ChatCFD | OpenFOAM | DeepSeek-R1+V3 | 82.1% exec / 68.12% phys (LLM-judge) | Physical fidelity=LLM-evaluated, not exp. |
| C5 | PhyNiKCE | OpenFOAM | Neurosymbolic | 96% over ChatCFD (340 runs) | Still OpenFOAM-only, no particle methods |
| C6 | FoamGPT | OpenFOAM | Qwen3-8B (LoRA) | 26.36% CFDLLMBench | Very low success rate |
| C7 | CFD-Copilot | OpenFOAM v2406 | Qwen3-8B LoRA + Qwen3-32B | U 96.4%, p 93.2% | Fine-tuning needed (49K pairs), mesh-only |
| C8 | AutoCFD | OpenFOAM | Qwen2.5-7B (fine-tuned) | 88.7% accuracy (21 cases) | Fine-tuning needed (28.7K pairs), mesh-only |
| C9 | MooseAgent | MOOSE (FEM) | **DeepSeek-R1+V3** | 93% success (9 cases) | FEM not SPH, cloud LLM |
| C10 | MCP-SIM | FEniCS (FEM) | N/A (6 agents) | 100% on 12 PDE tasks | FEM/PDE, no CFD/SPH |
| C11 | Pasimodo+RAG | Pasimodo (SPH) | Llama/Gemma (3-27B) | **0/2 model creation** | RAG-only, closed-source, cannot execute sims |

---

## Approved Research Gaps (6)

### GAP-1: Lagrangian Particle Simulation + LLM Agent = Zero Research [CRITICAL] — Confidence 95%
- **Evidence**: 6가지 검색 패턴(SPH+LLM, DualSPHysics+AI, particle+LLM, DEM+LLM, MPM+LLM)에서 zero results; 10+ OpenFOAM 시스템 전부 mesh-based; Pasimodo+RAG는 Pure RAG(0/2 model creation)
- **Refined Claim**: "To our knowledge, no prior work has combined LLM agents with particle-based simulation methods (SPH, DEM, MPM). All existing LLM-for-simulation agents target mesh-based solvers."
- **Implication**: 핵심 novelty — SPH뿐 아니라 전체 Lagrangian 입자법 영역이 비어 있음
- **Maps to**: EXP-1, EXP-2

### GAP-2: Lagrangian Particle Solver용 Tool Interface 설계 패턴 부재 [HIGH] — Confidence 95%
- **Evidence**: 11개 경쟁 시스템 tool counts 검증 — Foam-Agent 11 MCP, ChemCrow 18, MDCrow 40+, CFD-Copilot 100+; 전부 mesh-based I/O 패턴 (config dict → solver → postProcess). Particle solver의 XML→binary particle→GPU telemetry→particle-to-mesh 패턴 전무
- **Refined Claim**: "No prior work has designed tool abstractions for Lagrangian particle solvers, which require fundamentally different patterns: XML-based case definition, binary particle I/O, GPU-specific telemetry monitoring, and post-processing via particle-to-mesh conversion."
- **Implication**: 14개 DualSPHysics 도구가 ChemCrow(18)/MDCrow(40+) 수준의 기여
- **Maps to**: EXP-2, EXP-5

### GAP-3: 도메인 지식 Ablation 연구 전무 (전체 시뮬레이션 분야) [HIGH] — Confidence 90%
- **Evidence**: 9개 경쟁자 ablation 전수 조사 — MetaOpenFOAM(RAG/Reviewer/Temperature), Foam-Agent(Reviewer/RAG/File dep), ChatCFD(RAG/reflection), MooseAgent(RAG만); **전부 architectural ablation**. Domain knowledge(방정식, 제약조건, 용어) ablation = zero.
- **Refined Claim**: "While ablation studies for simulation agents exist (RAG, reviewer, temperature), no work has systematically ablated domain-specific knowledge from system prompts. Existing ablations test architectural components, not knowledge components."
- **Closest**: PhyNiKCE (structured vs zero-shot) — instruction FORMAT이지 knowledge CONTENT 아님
- **Implication**: Ablation study로 정량적 기여 가능 (20-40% 향상) — 최초의 knowledge ablation
- **Maps to**: EXP-4

### GAP-4: NL→실험 데이터 검증 E2E 파이프라인 부재 [MEDIUM-HIGH] — Confidence 90%
- **Evidence**: 모든 경쟁자가 "execution success"(코드 실행 여부)만 측정. ChatCFD의 68.12% physical fidelity도 LLM-as-judge(실험 데이터 아님). OpenFOAMGPT/AutoCFD만 reference simulation 비교 — **published experimental benchmark 대비 검증은 zero**
- **Refined Claim**: "No existing LLM-for-simulation system validates against published experimental benchmark data. 'Success' is universally measured as execution completion. SloshAgent is the first to close the loop from natural language to experimentally-validated physics (SPHERIC Test 10, 100-repeat pressure statistics)."
- **Implication**: SPHERIC Test 10 raw data로 Pearson r, NRMSE 정량 검증 — 가장 강력한 차별화
- **Maps to**: EXP-1, EXP-3

### GAP-5: 슬로싱 시뮬레이션 AI 자동화 부재 (산업 PoC) [MEDIUM-HIGH] — Confidence 85%
- **Evidence**: Maritime AI 시장 $4.3B('24), CAGR 40.6%; DualSPHysics 52K+ downloads, 687 stars, 746 citations; DesignSPHysics는 permanent beta (mDBC 미지원, FreeCAD 호환 문제, complex geometry 제한); nanoFluidX/PreonLab $10K+/year enterprise pricing; LNG 300+ 선박 발주 중, 한국 75% 점유
- **Refined Claim**: "Despite a $4.3B maritime AI market and 52,000+ DualSPHysics users, no AI tool automates sloshing simulation from natural language. DesignSPHysics is in permanent beta with critical limitations, commercial alternatives cost $10K+/year."
- **Practitioner Pain Points**: 초보자 2-4주 vs 전문가 2-3일, $80-120/hr consulting, SNU ANN은 experimental DB prediction만 (시뮬레이션 자동화 아님)
- **Implication**: 파라메트릭 자동화 + baffle 설계로 실용성 입증 — $0 LLM cost 대비 $10K+ 상용
- **Maps to**: EXP-3, EXP-5

### GAP-6: Fine-tuning 없이 오픈웨이트 LLM으로 시뮬레이션 오케스트레이션 [MEDIUM] — Confidence 80%
- **Evidence**: AutoCFD(Qwen2.5-7B, 28.7K pairs fine-tuning), FoamGPT(Qwen3-8B LoRA), CFD-Copilot(Qwen3-8B LoRA 49K pairs) — 전부 fine-tuning 필요. MDCrow(Llama3-405B) — 로컬 배포 비실용적. **Zero fine-tuning으로 general-purpose 오픈웨이트 모델이 시뮬레이션 오케스트레이션에 성공한 사례 없음**
- **Refined Claim**: "While fine-tuned open-weight models have been explored for CFD, no prior work demonstrates that a general-purpose open-weight model without domain-specific fine-tuning can orchestrate scientific simulations through prompt engineering alone."
- **Implication**: Prompt-only vs fine-tuning trade-off가 genuine contribution. Qwen3 32B vs 8B 비교는 향후 연구.
- **Note**: CFD-Copilot도 Qwen3-32B 사용하지만 general agent용이며, generation은 fine-tuned Qwen3-8B 사용
- **Maps to**: Future Work

---

## Experiment TODO List

### EXP-1: SPHERIC Test 10 벤치마크 재현
- **Purpose**: 에이전트 E2E 정확도 검증 (RQ2 직접 답변)
- **Dataset**: SPHERIC Test 10 FTP raw data (100회 반복, 3종 유체)
  - `datasets/spheric/case_1/` (압력 시계열 8파일)
- **Baseline**: English2021 mDBC 검증 결과 (DualSPHysics 직접 실행)
- **Procedure**:
  1. NL 입력: "Reproduce SPHERIC Test 10 with water at 0.613Hz roll"
  2. Agent generates XML → GenCase → DualSPHysics GPU → MeasureTool
  3. 압력 시계열 추출 (16 probe points)
  4. 실험 데이터 대비 RMSE/상관계수 계산
- **Metrics**: RMSE, Pearson correlation, peak pressure error %
- **Table**: Table 1 — Agent vs Expert vs Experimental (pressure at wall probes)
- **Figure**: Fig 3 — Pressure time series overlay (agent / expert / exp)
- **Expected outcome**: Agent-generated 시뮬레이션이 expert-level 정확도 달성 (r > 0.9)

### EXP-2: NL→XML 생성 정확도 평가
- **Purpose**: XML 생성 성공률 및 유효성 (RQ1 직접 답변)
- **Dataset**: 20개 NL 시나리오 (5 복잡도 x 4 도메인)
  - 간단: "0.9m x 0.062m x 0.508m tank at 0.6Hz"
  - 보통: "50% fill LNG tank with roll motion"
  - 복잡: "L-shaped tank with seismic input"
  - 고급: "Horizontal cylinder with vertical baffle"
- **Baseline**: Manual expert XML 작성 (ground truth)
- **Procedure**:
  1. 20개 NL 프롬프트를 에이전트에 입력
  2. 생성된 XML의 유효성 검증 (GenCase 성공 여부)
  3. 파라미터 정확도 비교 (dp, 치수, 운동 조건, 경계 조건)
- **Metrics**: XML validity rate, parameter accuracy, GenCase pass rate
- **Table**: Table 2 — NL→XML success rate by complexity level
- **Figure**: Fig 2 — System architecture diagram
- **Expected outcome**: 단순/보통 시나리오 90%+ 성공, 복잡 70%+, 고급 50%+
- **Contrast**: FoamGPT 26.36%, MetaOpenFOAM 85%, Foam-Agent 88.2% (all OpenFOAM)

### EXP-3: 파라메트릭 스터디 자동화
- **Purpose**: 산업 가치 입증 — 자동 vs 수동 비교 (RQ3 일부)
- **Dataset**: Chen2018 6개 fill level (20%, 30%, 40%, 50%, 60%, 70%)
  - 탱크: 600x300x650mm, 주파수: 1st natural freq
- **Baseline**: 수동 XML 수정 x 6회 (시간 측정)
  - Expert manual: ~2hr/case x 6 = 12hr
  - CFD consulting rate: $80-120/hr = $960-$1,440 for 6 cases
- **Procedure**:
  1. NL 입력: "Run parametric study: fill levels 20-70% in 10% steps"
  2. Agent generates 6 XMLs → 6 simulations → results comparison
  3. 수동 대비 시간 측정 (셋업 시간, 전체 시간)
- **Metrics**: Setup time reduction ratio, automation completeness
- **Table**: Table 3 — Agent vs Manual: setup time for N parametric cases
- **Figure**: Fig 5 — Free surface height vs fill level (6 cases overlay)
- **Expected outcome**: 셋업 시간 10x+ 단축, 전체 자동화율 80%+
- **Industry context**: LNG 탱크 평가 200+ 케이스 → 이 접근법으로 수 개월→수 일 압축 가능

### EXP-4: 도메인 프롬프트 Ablation
- **Purpose**: 도메인 특화 프롬프트의 효과 정량화 (RQ4 직접 답변)
- **Dataset**: EXP-2와 동일 20개 NL 시나리오
- **Baseline**: Generic CoderPrompt (SloshingCoderPrompt 제거)
- **Procedure**:
  1. SloshingCoderPrompt (136줄, 5카테고리) ON → 20개 시나리오
  2. SloshingCoderPrompt OFF (Generic) → 동일 20개 시나리오
  3. XML 생성 정확도, 물리 파라미터 추론 정확도 비교
- **Metrics**: XML validity delta, dp inference accuracy, parameter correctness
- **Table**: Table 4 — Ablation: Domain prompt ON vs OFF
- **Figure**: Fig 6 — Ablation bar chart (validity rate, parameter accuracy)
- **Expected outcome**: 도메인 프롬프트 ON 시 20-40% 정확도 향상
- **Novelty**: ChatCFD의 physical plausibility(68.12%)와 다른 접근 — 우리는 도메인 프롬프트 자체의 효과를 분리 측정

### EXP-5: Anti-Slosh Baffle 설계 시나리오 (산업 PoC)
- **Purpose**: 슬로싱 산업 응용 PoC (GAP-5 해결)
- **Dataset**: NASA 2023 ring baffle 파라미터 + 자동차 vertical baffle
- **Baseline**: 수동 XML 작성 (baffle geometry 포함)
- **Procedure**:
  1. NL: "Add vertical baffle at tank center, 80% height"
  2. Agent generates XML with baffle → simulation → force comparison
  3. With/without baffle 압력 비교
- **Metrics**: Baffle effect (force reduction %), XML generation accuracy
- **Table**: Table 5 — With vs Without baffle: peak force, sloshing amplitude
- **Figure**: Fig 7 — Particle snapshot comparison (with/without baffle)
- **Expected outcome**: Baffle가 피크 하중 30-50% 저감, 에이전트가 올바른 XML 생성
- **Industry context**: 자동차 연료탱크 baffle 최적화 → 슬로싱 진폭 70% 감소, 벽면 압력 50% 감소 (문헌 기준)

---

## Paper Structure (AI4Science Workshop, 4-8 pages)

1. **Introduction** (1p) — Sloshing importance (4 industries, quantified pain points), SPH+LLM blue ocean (8 competitors all mesh-based), contributions
2. **Related Work** (1p) — LLM agents for science (ChemCrow, MooseAgent), CFD automation (7 OpenFOAM systems), SPH sloshing (DualSPHysics, DesignSPHysics, Pasimodo+RAG)
3. **System Design** (1.5p) — Architecture (ReAct + 14 tools), tool interface patterns (IsError, async GPU, Run.csv), domain prompt (136 lines, 5 categories), XML generation
4. **Experiments** (2p) — EXP-1~5 results with competitor comparison
5. **Discussion** (0.5p) — Limitations, quantified industry implications, broader impact
6. **Conclusion** (0.5p) — Research space positioning, future work

## Key Figures Plan

| Fig # | Content | Purpose |
|-------|---------|---------|
| Fig 1 | Sloshing domain motivation (4 industries + quantified pain) | Introduction |
| Fig 2 | System architecture (NL→XML→GPU→Analysis) with tool taxonomy | System Design |
| Fig 3 | SPHERIC Test 10 pressure comparison (agent/expert/exp) | EXP-1 |
| Fig 4 | Tool calling trace (ReAct loop visualization) | System Design |
| Fig 5 | Parametric study results (6 fill levels overlay) | EXP-3 |
| Fig 6 | Ablation bar chart (domain prompt ON vs OFF) | EXP-4 |
| Fig 7 | Baffle effect particle snapshots (with/without) | EXP-5 |

## Key Tables Plan

| Table # | Content | Purpose |
|---------|---------|---------|
| Table 1 | 9-system competitor comparison (comprehensive) | Related Work |
| Table 2 | NL→XML success rate by complexity (5 levels) | EXP-2 |
| Table 3 | Agent vs Manual setup time (with cost analysis) | EXP-3 |
| Table 4 | Ablation: domain prompt ON vs OFF | EXP-4 |
| Table 5 | Baffle effect comparison (force, amplitude) | EXP-5 |

---

## Industrial Data Summary (from deep research agents)

### Practitioner Pain Points (Quantified)
| Pain Point | Evidence | Impact |
|-----------|---------|--------|
| Setup time | Expert 2-3 days, beginner 2-4 weeks | 10x+ difference |
| CFD consulting | $80-120/hr, $2K-$20K/project | High cost barrier |
| Expertise barrier | MSc+ degree, 5+ years experience | Limits adoption |
| Parameter tuning | dp, alpha, DensityDT, Shifting all empirical | Error-prone |
| Common errors | "boundary excluded", STL autofill, motion syntax | Forum recurring |

### Industry-Specific Data
| Industry | Key Metric | Source |
|---------|-----------|--------|
| LNG | 200+ cases/tank, model tests $millions, SNU 20K+ hr DB | Classification societies |
| Automotive | NVH complaints rising, baffle -70% amplitude, -50% pressure | IJRASET, AVL |
| Nuclear | 0.5% critical damping (10x sensitivity), SFP meltdown risk | ASCE 4-98, NRC RG 1.29 |
| Aerospace | Falcon 1 loss 2007, NASA SP-8009 standard | NASA NTRS |

### AI Automation Expected Benefits (BCG/Rescale)
- R&D cost: **20% reduction**
- Time-to-market: **10-20% reduction**
- Design exploration: **weeks → days** (parallel execution)
- Setup automation: **80% of engineer time** currently spent on pre-processing
