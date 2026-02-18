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
| **SPH particle-based** (DualSPHysics GPU) | MetaOpenFOAM/Foam-Agent/ChatCFD/OpenFOAMGPT/FoamGPT/AutoCFD/CFD-Copilot | All 7 mesh-based OpenFOAM |
| **Autonomous agent** (14 tools + ReAct) | Pasimodo+RAG | RAG = assistance, not autonomous |
| **Sloshing domain-specialized** (136-line prompt) | All above (general-purpose CFD) | No domain specialization for sloshing |
| **Benchmark-validated** (SPHERIC Test 10, r>0.9) | ChatCFD (68.12% phys plausibility) | First physical accuracy on SPH benchmark |
| **Local open-weight** (Qwen3 32B) | MetaOpenFOAM/ChatCFD (GPT-4/DeepSeek) | IP protection, zero cloud cost |
| **GPU-native execution** (RTX 4090) | MooseAgent (CPU FEM) | SPH is inherently GPU-parallel |

## Competitive Landscape (8 direct competitors)

| # | System | Solver | Best Metric | Weakness vs Ours |
|---|--------|--------|------------|-----------------|
| C1 | MetaOpenFOAM | OpenFOAM | 85% pass, $0.22 | Mesh-based, cloud LLM |
| C2 | OpenFOAMGPT 2.0 | OpenFOAM | 100% reproducibility | Academic cases only |
| C3 | Foam-Agent 2.0 | OpenFOAM | 88.2% success | MCP similar but mesh-only |
| C4 | ChatCFD | OpenFOAM | 82.1% exec, 68.1% phys | Cloud LLM, no SPH |
| C5 | FoamGPT | OpenFOAM | 26.36% CFDLLMBench | Low success rate |
| C6 | AutoCFD | OpenFOAM | 88.7% accuracy | 16 basic cases only |
| C7 | MooseAgent | MOOSE (FEM) | 93% success | FEM not SPH |
| C8 | Pasimodo+RAG | Pasimodo (SPH) | Qualitative only | RAG-only, closed-source |

---

## Approved Research Gaps (6)

### GAP-1: SPH + LLM Agent 조합 연구 전무 [CRITICAL]
- **Evidence**: 2,175편 survey "sph-specific" coverage gap; "sloshing+LLM" web search = 0건; 8개 LLM+CFD 시스템 전부 OpenFOAM
- **Implication**: 핵심 novelty claim의 근거. Research space matrix에서 SPH×LLM Agent 셀이 비어 있음.
- **Maps to**: EXP-1, EXP-2

### GAP-2: Particle-Based Solver용 Tool Interface 설계 패턴 부재 [HIGH]
- **Evidence**: MetaOpenFOAM=4 agents, Foam-Agent=MCP, MooseAgent=LangGraph; SPH 비동기 GPU + IsError + Run.csv 모니터링 패턴 전무
- **Implication**: 14개 DualSPHysics 도구 인터페이스가 논문 기여. ChemCrow(18 tools), MDCrow(40 tools)와 동일 레벨.
- **Maps to**: EXP-2, EXP-5

### GAP-3: 슬로싱 도메인 특화 프롬프트 체계적 평가 부재 [HIGH]
- **Evidence**: ChatCFD가 최초 physical plausibility 지표(68.12%) 도입했으나 도메인 프롬프트 ablation 없음; ChemCrow은 화학 특화
- **Implication**: Ablation study로 정량적 기여 가능 (20-40% 향상)
- **Maps to**: EXP-4

### GAP-4: NL→벤치마크 검증 E2E 파이프라인 부재 [MEDIUM-HIGH]
- **Evidence**: FoamGPT 최고 26.36% 실행률; OpenFOAMGPT 100% 재현성이나 학술 케이스만; 벤치마크 대비 정량 검증 없음
- **Implication**: SPHERIC Test 10 raw data (100회 반복)로 정량 검증 가능
- **Maps to**: EXP-1, EXP-3

### GAP-5: 슬로싱 산업 PoC 및 현장 기대효과 검증 부재 [MEDIUM-HIGH]
- **Evidence**: 4대 산업 정량 데이터 확보 — LNG(200+ cases/tank), 자동차($2K-$20K/project), 원자력(ASCE 4-98), 우주(Falcon 1 실패). AI PoC 부재.
- **Practitioner Pain Points**: 초보자 2-4주 vs 전문가 2-3일, CFD 컨설팅 $80-120/hr, DualSPHysics 포럼 상위 에러(boundary excluded, STL autofill, mvrotsinu 속성 혼동)
- **Implication**: 파라메트릭 자동화 + baffle 설계로 실용성 입증
- **Maps to**: EXP-3, EXP-5

### GAP-6: 로컬/오픈웨이트 LLM 시뮬레이션 적용 [MEDIUM]
- **Evidence**: arXiv:2504.02888 Qwen vs GPT-4 비교 — 오픈웨이트 LLM이 CFD에 경쟁력 있음 확인. MetaOpenFOAM/ChatCFD는 클라우드 의존.
- **Implication**: Qwen3 32B 로컬 배포 전략의 근거. 32B vs 8B 비교는 향후 연구.
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
