# Research-v3: 실험 상태 보고서

**Date**: 2026-03-07
**Status**: ALL EXPERIMENTS COMPLETE (EXP-A/B/C + Bridge)

---

## 1. 전체 실험 완료 상태

| 실험 | 목적 | 상태 | 실행 수 | 비고 |
|------|------|:---:|------:|------|
| **EXP-A** | NL->XML 파이프라인 검증 | DONE | 60/60 | 3 trials x 2 models x 10 scenarios |
| **EXP-B** | 아키텍처 Ablation | DONE | 80/80 | 4 conditions x 10 scenarios x 2 models |
| **EXP-C** | SPHERIC 물리 검증 | DONE | 7 runs | research-v3/exp-c/ (이관 완료) |
| **Bridge** | Agent XML -> Physics | DONE | 1 run | S02 agent-gen XML -> Solver -> 5/5 PASS |

---

## 2. EXP-A 결과 요약

**채점 기준**: score_expb.py 8-parameter M-A3 (tank x3, fill_height, motion_type, frequency, amplitude, timemax)

### 2.1 전체 M-A3

**Overall M-A3**: 69.5% (32B), 64.5% (8B) -- Delta=5.0%pp

| 난이도 | 시나리오 | 32B M-A3 | 8B M-A3 | 비고 |
|--------|----------|:---:|:---:|------|
| Easy (S01-S03) | S01,S02,S03 | 87.7% | 87.7% | 완전 동일 |
| Medium (S04-S07) | S04,S05,S06,S07 | 70.0% | 70.0% | 완전 동일 |
| Hard (S08-S10) | S08,S09,S10 | 51.0% | 34.3% | S08만 차이 |

### 2.2 시나리오별 상세

```
Scenario  32B   8B   Delta   Key Issue
--------------------------------------------
S01       75%   75%  =       P1: mvrectsinu (pitch->sway)
S02      100%  100%  =       Perfect
S03       88%   88%  =       timemax underestimate
S04       75%   75%  =       P1: mvrectsinu
S05       50%   50%  =       tank_z/fill_height confusion
S06       88%   88%  =       timemax underestimate
S07       67%   67%  =       cylinder->box conversion
S08       50%    0%  +50%    8B: XML not generated (tool calling failure)
S09       20%   20%  =       cylinder + P1 + timemax
S10       83%   83%  =       STL import OK, timemax=10 vs 15
```

**9/10 scenarios identical** between 32B and 8B. Only S08 differs.

### 2.3 Determinism (temperature=0)

**M-A3 score reproducibility**: 32B 10/10, 8B 10/10 deterministic (sigma=0.0)
- S10 8B: trial1=trial2=trial3 (verified md5 identical)
**Full XML reproducibility**: 32B 9/10 identical, 8B 9/10 identical
- Differences in non-scored areas (comments, GenCase post-processing, simulationdomain)

### 2.4 체계적 오류 패턴

| 패턴 | 파라미터 | 영향 시나리오 | 빈도(20 runs) |
|------|----------|--------------|:---:|
| P1 | motion_type (pitch->mvrectsinu) | S01,S04,S05,S08,S09 | 9 |
| P2 | amplitude (deg->rad conversion) | S01,S04,S05,S08,S09 | 9 |
| P3 | timemax underestimate | S03,S06,S07,S08,S09,S10 | 11 |
| P4 | geometry_type (cylinder unsupported) | S07,S09 | 4 |
| P5 | tank_z/fill_height confusion | S05 | 2 |
| P6 | 8B tool calling instability | S08 | -- |

---

## 3. EXP-B 결과 -- 2x2 Factorial Ablation

### 3.1 실험 설계

| Condition | Domain Prompt | Tools | 실행 방법 |
|-----------|:---:|:---:|------|
| B0 Full | Y | Y | slosim (EXP-A 결과 활용) |
| B1 -Prompt | N | Y | slosim-b1 (SLOSIM_GENERIC_PROMPT=1) |
| B2 -Tool | Y | N | Ollama API + domain_prompt.txt |
| B4 Bare | N | N | Ollama API + generic prompt |

### 3.2 M-A3 결과 (10 scenarios x 2 models)

```
  qwen3_32b:
  Cond   S01(E) S02(E) S03(E) S04(M) S05(M) S06(M) S07(M) S08(H) S09(H) S10(H)    Mean
  ----
  B0        75%   100%    88%    75%    50%    88%    67%    50%    20%    83%   69.5%
  B1         0%     0%     0%     0%     0%     0%     0%     0%     0%     0%    0.0%
  B2        62%    75%    62%    62%     0%    75%    50%    25%    40%    67%   51.9%
  B4         0%     0%     0%     0%     0%     0%     0%     0%     0%     0%    0.0%

  qwen3_latest (8B):
  Cond   S01(E) S02(E) S03(E) S04(M) S05(M) S06(M) S07(M) S08(H) S09(H) S10(H)    Mean
  ----
  B0        75%   100%    88%    75%    50%    88%    67%     0%    20%    83%   64.5%
  B1         0%     0%     0%     0%     0%     0%     0%     0%     0%     0%    0.0%
  B2        62%     0%    62%    62%    75%     0%    50%     0%    40%    50%   40.2%
  B4         0%     0%     0%     0%     0%     0%     0%     0%     0%     0%    0.0%
```

**Combined**: B0=67.0%, B1=0.0%, B2=46.1%, B4=0.0%

### 3.3 2x2 Factorial Analysis

```
  Combined (20 runs per condition):
  Main effect of Domain Prompt: +56.5%pp
  Main effect of Tools:         +10.5%pp
  Interaction (Prompt x Tools): +21.0%pp

  Per-model:
  32B -- Prompt: +60.7%, Tool: +8.8%, Interaction: +17.6%
  8B  -- Prompt: +52.4%, Tool: +12.1%, Interaction: +24.2%
```

### 3.4 핵심 발견

1. **Domain Prompt가 지배적** (+56.5%pp): 프롬프트 없이는 아무것도 작동하지 않음
   - B1=B4=0%: 10개 전 시나리오 x 2 모델 = 40 runs 전부 0%
   - B1: 모델이 DualSPHysics 도구를 전혀 호출하지 않음
   - Hierarchical dependency: prompt는 tool 사용의 전제 조건

2. **Tools는 Prompt와 결합할 때만 유의미** (+10.5%pp main, +21.0%pp interaction)
   - B0(67.0%) vs B2(46.1%): 도구가 +20.9pp 순기여
   - 도구의 역할: 구조화된 입력으로 geometry/fill height 정확도 보장
   - S10에서 tool 기여 극명: B0=83% vs B2=67%(32B)/50%(8B) (STL import)

3. **모델 크기 차이 미미**: 9/10 시나리오 동일
   - B0: 32B=69.5% vs 8B=64.5% (Delta=5.0pp)
   - 유일한 차이 S08: 32B=50% vs 8B=0% (tool calling 실패)
   - 아키텍처 효과(+56.5pp) >> 모델 크기 효과(5.0pp)

4. **B2의 한계**: 도메인 프롬프트만으로 XML 생성 가능하나:
   - tank_z와 fill_height를 뒤바꿈 (도구가 이를 방지)
   - 원통형 형상은 여전히 박스로 생성
   - 8B에서 B2 timeout/파싱 실패 더 빈번 (S02,S06,S08)

---

## 4. Bridge Experiment (EXP-A -> EXP-C)

**Purpose**: Agent-generated XML이 물리적으로 유효한 시뮬레이션을 생성하는지 직접 검증

**Method**: S02 agent-generated XML (M-A3=100%) -> GenCase -> DualSPHysics GPU -> MeasureTool

**Results**: 5/5 physics checks PASSED

| Check | Metric | Result | Status |
|-------|--------|--------|:---:|
| Hydrostatic init | P_expected=323.7Pa | P_measured=333.5Pa (err 3.0%) | PASS |
| Left-Right anti-phase | r(left,right) | r=-0.720 | PASS |
| Oscillation frequency | f_expected=0.756Hz | f_measured=0.789Hz (err 4.4%) | PASS |
| Wall > Center amplitude | wall_range vs center | 465 vs 213 Pa | PASS |
| Simulation stability | particle retention | 105,465/105,465 (100%) | PASS |

**Significance**: Bridges the gap between EXP-A (XML generation quality) and EXP-C (physics validity).
Agent M-A3=100% XML -> correct hydrostatic equilibrium -> correct sloshing dynamics.

---

## 5. 논문 Claim -> 실험 매핑

### C1: "NL->SPH 파이프라인이 작동한다" -> EXP-A + Bridge

**Evidence**: 69.5% M-A3 (32B), 64.5% (8B). Bridge 5/5 physics validation.
**Strength**: Strong. Easy/Medium 70-88% success. Bridge proves physical validity.

### C2: "아키텍처 컴포넌트가 각각 필수적이다" -> EXP-B

**Evidence**: 2x2 factorial. Prompt +56.5%pp, Tool +10.5%pp, Interaction +21.0%pp.
B1=B4=0% (40/40 runs). Hierarchical dependency structure.
**Strength**: Very Strong. Cleanest experiment design.

### C3: "물리적으로 신뢰할 수 있다" -> EXP-C + Bridge

| Sub-case | Status | Metric |
|----------|:---:|--------|
| Water Lateral | PASS | r=0.655, M2=19.5% |
| Oil Lateral | PASS | r=0.570, M7=4/4 |
| Water Roof | PARTIAL | DBC boundary limitation |
| Convergence | PASS | dp 4mm->2mm->1mm trend |
| Bridge (Agent XML) | PASS | 5/5 physics checks |

**Strength**: Strong. SPHERIC benchmark + agent-generated XML validation.

### C4: "로컬 SLM으로 충분하다" -> EXP-A + EXP-B

**Evidence**: 9/10 identical scores. Delta=5.0pp. Model-invariant condition ranks.
**Strength**: Strong.

---

## 6. 논문 논리 흐름

```
Section 4.1 (EXP-A): "파이프라인이 67% M-A3로 작동"
    -> 모델 크기 무관 (32B 69.5% ~ 8B 64.5%, Delta=5.0pp)
    -> 도구 설계가 병목 (P1-P6 패턴, motion_type/amplitude가 최다 실패)

Section 4.2 (EXP-B): "도메인 프롬프트가 핵심, 도구가 보조"
    -> Domain Prompt main effect: +56.5%pp
    -> Tools main effect: +10.5%pp (with +21.0%pp interaction)
    -> Prompt 없이는 도구 호출 불가 (B1=B4=0%, 40/40 runs)
    -> Hierarchical dependency structure

Section 4.3 (EXP-C + Bridge): "물리적으로 정확"
    -> SPHERIC PASS (2/3 sub-case)
    -> Agent-generated XML -> 5/5 physics validation
    -> 수렴 확인 (dp refinement)
```

---

## 7. 완료 체크리스트

1. [x] EXP-A 3 trials 완료 (60/60 runs, S10 stl_import fix 반영)
2. [x] EXP-A determinism 확인 (32B 10/10, 8B 10/10)
3. [x] EXP-B B0/B1/B2/B4 완료 (80/80 runs)
4. [x] EXP-B 결과 채점 및 2x2 factorial 분석 (10-시나리오)
5. [x] EXP-C 데이터 이관 완료 (research-v3/exp-c/)
6. [x] Bridge experiment 완료 (S02 -> GenCase -> Solver -> 5/5 PASS)
7. [x] 최종 결과 문서 (이 파일)
8. [x] Git commit
9. [x] 채점 기준 통일 (score_expb 8-param -> EXP-A 전체 적용)
10. [x] EXP-A 60개 결과 로컬 동기화 완료

---

## 8. 파일 구조

```
research-v3/
|-- EXPERIMENT_STATUS.md     # 이 파일
|-- GAP_PROOF_ANALYSIS.md    # Gap-Experiment 충분성 분석
|-- exp-a/
|   |-- prompts/             # 10개 NL 시나리오 (S01-S10)
|   |-- configs/             # Ollama 모델 설정
|   |-- ground_truth.json    # 평가 기준
|   |-- run_scenario.sh      # 단일 실행 스크립트
|   |-- analyze_all.py       # EXP-A 결과 분석
|   `-- results/             # 60개 실행 결과 (pajulab)
|-- exp-b/
|   |-- domain_prompt.txt    # B2용 도메인 프롬프트
|   |-- ollama_generate.py   # B2/B4 공용 Ollama API 호출
|   |-- run_b1_batch.sh      # B1 배치
|   |-- run_b2_notool.sh     # B2 배치
|   |-- run_b4_bare.sh       # B4 배치
|   |-- score_expb.py        # EXP-B M-A3 채점 + factorial
|   `-- results/             # 80개 실행 결과
|-- exp-c/                   # SPHERIC 물리 검증
|   |-- analysis/            # 메트릭 (M2=19.5%, r=0.655)
|   |-- figures/             # 검증 그래프
|   |-- agent-bridge/        # Bridge experiment output
|   |   |-- AgentS02_Def.xml # Agent-generated input
|   |   |-- gauges_*.csv     # MeasureTool output (Press, Vel, Rhop)
|   |   `-- AgentS02/        # Solver output (51 Part files)
|   |-- analyze_bridge.py    # Bridge physics validation script
|   `-- run_agent_bridge.sh  # Bridge experiment runner
`-- figures/                 # 논문용 그래프
    |-- fig_expb_factorial.py
    `-- generate_tables.py
```
