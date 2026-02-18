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
| **Agent** (autonomous pipeline) | Pasimodo+RAG | RAG = assistance tool, not autonomous |
| **Structured tool calling** (14 tools + IsError) | AutoCFD/OpenFOAMGPT/FoamGPT | All mesh-based (OpenFOAM), not SPH |
| **Benchmark-validated** (SPHERIC Test 10) | FoamGPT (26.36% exec rate) | No benchmark accuracy evaluation |
| **Local open-weight** (Qwen3 32B) | ChemCrow/Coscientist (GPT-4) | IP/data protection for industry |
| **Domain-specialized prompt** (136 lines, 5 categories) | Generic LLM agents | No ablation for comp. mechanics |

---

## Approved Research Gaps (6)

### GAP-1: SPH + LLM Agent 조합 연구 전무 [CRITICAL]
- **Evidence**: 2,175편 survey "sph-specific" coverage gap; "sloshing+LLM" web search = 0건
- **Implication**: 핵심 novelty claim의 근거. Mesh-based LLM+CFD만 존재.
- **Maps to**: EXP-1, EXP-2

### GAP-2: Particle-Based Solver용 Tool Interface 설계 패턴 부재 [HIGH]
- **Evidence**: MechAgents=FEA, ChemCrow=화학; SPH 비동기 GPU + IsError 패턴 전무
- **Implication**: 14개 DualSPHysics 도구 인터페이스가 논문 기여
- **Maps to**: EXP-2, EXP-5

### GAP-3: 슬로싱 도메인 특화 프롬프트 체계적 평가 부재 [HIGH]
- **Evidence**: ChemCrow은 화학 특화; 전산역학/SPH 프롬프트 평가 없음
- **Implication**: Ablation study로 정량적 기여 가능
- **Maps to**: EXP-4

### GAP-4: NL→벤치마크 검증 E2E 파이프라인 부재 [MEDIUM-HIGH]
- **Evidence**: FoamGPT 최고 26.36% 실행률, 벤치마크 대비 정량 검증 없음
- **Implication**: SPHERIC Test 10 raw data로 정량 검증 가능
- **Maps to**: EXP-1, EXP-3

### GAP-5: 슬로싱 산업 PoC 및 현장 기대효과 검증 부재 [MEDIUM-HIGH]
- **Evidence**: 4대 산업(LNG/자동차/원자력/우주)에서 슬로싱 중요, AI PoC 부재
- **Implication**: 파라메트릭 자동화 + baffle 설계로 실용성 입증
- **Maps to**: EXP-3, EXP-5

### GAP-6: 로컬/오픈웨이트 LLM 시뮬레이션 적용 [MEDIUM]
- **Evidence**: Pasimodo+RAG가 local LLM 관심 언급; 산업 IP 보호 이슈
- **Implication**: Qwen3 32B/8B 비교는 향후 연구로 위임
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
- **Dataset**: 20개 NL 시나리오 (5 복잡도 × 4 도메인)
  - 간단: "0.9m × 0.062m × 0.508m tank at 0.6Hz"
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

### EXP-3: 파라메트릭 스터디 자동화
- **Purpose**: 산업 가치 입증 — 자동 vs 수동 비교 (RQ3 일부)
- **Dataset**: Chen2018 6개 fill level (20%, 30%, 40%, 50%, 60%, 70%)
  - 탱크: 600×300×650mm, 주파수: 1st natural freq
- **Baseline**: 수동 XML 수정 × 6회 (시간 측정)
- **Procedure**:
  1. NL 입력: "Run parametric study: fill levels 20-70% in 10% steps"
  2. Agent generates 6 XMLs → 6 simulations → results comparison
  3. 수동 대비 시간 측정 (셋업 시간, 전체 시간)
- **Metrics**: Setup time reduction ratio, automation completeness
- **Table**: Table 3 — Agent vs Manual: setup time for N parametric cases
- **Figure**: Fig 5 — Free surface height vs fill level (6 cases overlay)
- **Expected outcome**: 셋업 시간 10x+ 단축, 전체 자동화율 80%+

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

---

## Paper Structure (AI4Science Workshop, 4-8 pages)

1. **Introduction** (1p) — Sloshing importance, SPH+LLM blue ocean, contributions
2. **Related Work** (1p) — LLM agents for science, CFD automation, SPH sloshing
3. **System Design** (1.5p) — Architecture, tool interface, domain prompt, agent loop
4. **Experiments** (2p) — EXP-1~5 results
5. **Discussion** (0.5p) — Limitations, industry implications
6. **Conclusion** (0.5p)

## Key Figures Plan

| Fig # | Content | Purpose |
|-------|---------|---------|
| Fig 1 | Sloshing domain motivation (4 industries) | Introduction |
| Fig 2 | System architecture (NL→XML→GPU→Analysis) | System Design |
| Fig 3 | SPHERIC Test 10 pressure comparison | EXP-1 |
| Fig 4 | Tool calling trace (ReAct loop visualization) | System Design |
| Fig 5 | Parametric study results (6 fill levels) | EXP-3 |
| Fig 6 | Ablation bar chart | EXP-4 |
| Fig 7 | Baffle effect particle snapshots | EXP-5 |

## Key Tables Plan

| Table # | Content | Purpose |
|---------|---------|---------|
| Table 1 | Competitor comparison (us vs AutoCFD vs FoamGPT vs Pasimodo+RAG) | Related Work |
| Table 2 | NL→XML success rate by complexity | EXP-2 |
| Table 3 | Agent vs Manual setup time | EXP-3 |
| Table 4 | Ablation: domain prompt ON vs OFF | EXP-4 |
| Table 5 | Baffle effect comparison | EXP-5 |
