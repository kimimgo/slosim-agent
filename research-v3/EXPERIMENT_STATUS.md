# Research-v3: 실험 상태 보고서

**Date**: 2026-03-03
**Status**: ✅ 전체 실험 완료

---

## 1. 전체 실험 완료 상태

| 실험 | 목적 | 상태 | 실행 수 | 비고 |
|------|------|:---:|------:|------|
| **EXP-A** | NL→XML 파이프라인 검증 | ✅ 완료 | 60/60 | 3 trials × 2 models × 10 scenarios |
| **EXP-B** | 아키텍처 Ablation | ✅ 완료 | 18/18 | 4 conditions × 3 scenarios × 2 models |
| **EXP-C** | SPHERIC 물리 검증 | ✅ 데이터 존재 | 7 runs | research-v2에서 이관 가능 |

---

## 2. EXP-A 결과 요약

**채점 기준**: score_expb.py 8-parameter M-A3 (tank×3, fill_height, motion_type, frequency, amplitude, timemax)

### 2.1 전체 M-A3

**Overall M-A3**: 61.2% (32B) ≈ 58.7% (8B) — Δ=2.5%pp

| 난이도 | 시나리오 | 32B M-A3 | 8B M-A3 | 비고 |
|--------|----------|:---:|:---:|------|
| Easy (S01-S03) | S01,S02,S03 | 87.5% | 87.5% | 완전 동일 |
| Medium (S04-S07) | S04,S05,S06,S07 | 69.8% | 69.8% | 완전 동일 |
| Hard (S08-S10) | S08,S09,S10 | 23.3% | 15.0% | S10 비결정성 |

### 2.2 시나리오별 상세

```
Scenario  32B   8B   Δ      Key Issue
────────────────────────────────────────────────
S01       75%   75%  =      P1: mvrectsinu (pitch→sway 혼동)
S02      100%  100%  =      완벽
S03       88%   88%  =      timemax 과소추정
S04       75%   75%  =      P1: mvrectsinu
S05       50%   50%  =      tank_z/fill_height 혼동
S06       88%   88%  =      timemax 과소추정
S07       67%   67%  =      원통형 → 박스 변환
S08       50%    0%  +50%   8B: XML 미생성 (tool calling 실패)
S09       20%   20%  =      원통형 + P1 + timemax
S10        0%   25%  -25%   STL 미지원, 8B T2/T3만 부분 성공
```

### 2.3 Determinism (temperature=0)

**M-A3 점수 재현성**: 32B 10/10, 8B 9/10 deterministic (σ=0.0)
- 유일한 비결정적 케이스: S10_8B (trial1=0%, trial2/3=75%, σ=35.4)
**Full XML 재현성**: 32B 9/10 identical, 8B 4/10 identical
- 8B 차이는 비채점 영역 (comments, GenCase 후처리, simulationdomain 설정)
- M-A3 채점 대상 파라미터는 S10_8B 제외 전부 동일

### 2.4 체계적 오류 패턴

| 패턴 | 파라미터 | 영향 시나리오 | 빈도(20 runs) |
|------|----------|--------------|:---:|
| P1 | motion_type (pitch→mvrectsinu) | S01,S04,S05,S08,S09 | 9 |
| P2 | amplitude (deg→rad 변환) | S01,S04,S05,S08,S09 | 9 |
| P3 | timemax 과소추정 | S03,S06,S07,S08,S09 | 9 |
| P4 | geometry_type (원통형 미지원) | S07,S09 | 4 |
| P5 | tank_z/fill_height 혼동 | S05 | 2 |
| P6 | 8B tool calling 불안정 | S08,S10 | — |

---

## 3. EXP-B 결과 — 2×2 Factorial Ablation

### 3.1 실험 설계

| Condition | Domain Prompt | Tools | 실행 방법 |
|-----------|:---:|:---:|------|
| B0 Full | ✓ | ✓ | slosim (EXP-A 결과 활용) |
| B1 −Prompt | ✗ | ✓ | slosim-b1 (SLOSIM_GENERIC_PROMPT=1) |
| B2 −Tool | ✓ | ✗ | Ollama API + domain_prompt.txt |
| B4 Bare | ✗ | ✗ | Ollama API + generic prompt |

### 3.2 M-A3 결과

```
  Condition     S01(Easy)   S04(Med)  S07(Hard)     Mean
  ────────────────────────────────────────────────────
  B0 Full            75%        75%        67%     72.2%
  B1 −Prompt          0%         0%         0%      0.0%
  B2 −Tool           62%        62%        50%     58.3%
  B4 Bare             0%         0%         0%      0.0%
```

**32B ≡ 8B**: 모든 조건에서 동일 결과.

### 3.3 2×2 Factorial Analysis

```
  Main effect of Domain Prompt: +65.3%
  Main effect of Tools:          +6.9%
  Interaction (Prompt × Tools): +13.9%
```

### 3.4 핵심 발견

1. **Domain Prompt가 지배적** (+65.3%): 프롬프트 없이는 아무것도 작동하지 않음
   - B1: 모델이 DualSPHysics 도구(xml_generator 등)를 전혀 호출하지 않음
   - 대신 generic tools(glob, grep, edit)로 기존 XML 편집을 시도 → 실패

2. **Tools는 Prompt와 결합할 때만 유의미** (+6.9%→+13.9% with interaction)
   - B0(72.2%) vs B2(58.3%): 도구가 +13.9pp 기여
   - B2의 재미있는 패턴: motion_type은 B0보다 정확 (mvrotsinu vs mvrectsinu)
   - 도구의 역할: 구조화된 입력으로 geometry/fill height 정확도 보장

3. **모델 크기 무의미**: 32B ≡ 8B — 아키텍처 설계가 성능을 결정

4. **B2의 한계**: 도메인 프롬프트만으로 DualSPHysics XML 구조를 생성하나:
   - tank_z와 fill_height를 뒤바꿈 (도구가 이를 방지)
   - amplitude를 radians로 변환 (도구는 degrees 유지)
   - 원통형 형상은 여전히 박스로 생성

---

## 4. 논문 Claim → 실험 매핑

### C1: "NL→SPH 파이프라인이 작동한다" → EXP-A ✅

**적절성**: 완전히 뒷받침. 61% M-A3는 "작동하지만 개선 여지 있음"을 보여줌.

**잠재적 약점**:
- Qwen3 계열만 테스트 (LLaMA, Gemma 비교 없음) → Future work
- temperature=0만 — 확률적 출력 미탐색 → 정당화: 재현성 우선
- 10개 시나리오 — NL2FOAM의 21개보다 적음 → SPH 특화로 정당화

### C2: "아키텍처 컴포넌트가 각각 필수적이다" → EXP-B ✅

**적절성**: 2×2 factorial이 깔끔한 ablation 제공.

**핵심 논증**:
```
Section 4.2 (EXP-B): "도메인 프롬프트가 핵심 드라이버"
    → Prompt main effect (+65.3%) >> Tools main effect (+6.9%)
    → Prompt 없이는 도구 호출 자체가 불가 (B1=0%)
    → Prompt+Tools 시너지 (+13.9%)
```

### C3: "물리적으로 신뢰할 수 있다" → EXP-C ✅

**적절성**: SPHERIC Test 10은 국제 벤치마크. 2/3 sub-case PASS.

| Sub-case | 상태 | 메트릭 |
|----------|:---:|--------|
| Water Lateral | ✅ PASS | r=0.655, M2=19.5% |
| Oil Lateral | ✅ PASS | r=0.570, M7=4/4 |
| Water Roof | ⚠️ PARTIAL | DBC 한계 (mDBC diverged) |
| Convergence | ✅ | dp 4mm→2mm→1mm, GCI 분석 |

---

## 5. 논문 논리 흐름

```
Section 4.1 (EXP-A): "파이프라인이 61% M-A3로 작동"
    → 모델 크기 무관 (32B 61.2% ≈ 8B 58.7%, Δ=2.5%pp)
    → 도구 설계가 병목 (P1-P6 패턴, motion_type/amplitude가 최다 실패)

Section 4.2 (EXP-B): "도메인 프롬프트가 핵심, 도구가 보조"
    → Domain Prompt main effect: +65.3%
    → Tools main effect: +6.9% (with +13.9% interaction)
    → Prompt 없이는 도구 호출 불가 (B1=0%)
    → 모델 크기 ALL 조건에서 무관

Section 4.3 (EXP-C): "물리적으로 정확"
    → SPHERIC PASS (2/3 sub-case)
    → 수렴 확인 (dp refinement)
```

**논리적 일관성**: ✓ C1→C2→C3가 자연스럽게 연결됨.
"작동한다(C1) → 핵심은 프롬프트+도구 시너지(C2) → 결과가 물리적으로 타당(C3)"

---

## 6. 완료 체크리스트

1. [x] EXP-A 3 trials 완료 (60/60 runs)
2. [x] EXP-A determinism 확인 (32B 10/10, 8B 9/10)
3. [x] EXP-B B4 완료 (6/6 runs)
4. [x] EXP-B B2 완료 (6/6 runs, 재실행 포함)
5. [x] EXP-B B1 완료 (6/6 runs)
6. [x] EXP-B 결과 채점 및 2×2 factorial 분석
7. [x] EXP-C v2 데이터 확인 (이관 가능)
8. [x] 최종 결과 문서 (이 파일)
9. [x] Git commit (64078e9, a93a3fd, e423813, 12103a5)
10. [x] 채점 기준 통일 (score_expb 8-param → EXP-A 전체 적용)
11. [x] EXP-A 60개 결과 로컬 동기화 완료
12. [x] analyze_all.py — EXP-A 논문용 종합 분석 스크립트

---

## 7. 파일 구조

```
research-v3/
├── EXPERIMENT_PLAN.md       # 실험 설계서
├── EXPERIMENT_STATUS.md     # 이 파일
├── exp-a/
│   ├── prompts/             # 10개 NL 시나리오 (S01-S10)
│   ├── configs/             # Ollama 모델 설정
│   ├── ground_truth.json    # 평가 기준
│   ├── run_scenario.sh      # 단일 실행 스크립트
│   ├── analyze_all.py       # EXP-A 결과 분석
│   └── results/             # 60개 실행 결과 (pajulab)
└── exp-b/
    ├── domain_prompt.txt    # B2용 도메인 프롬프트
    ├── ollama_generate.py   # B2/B4 공용 Ollama API 호출
    ├── run_b1_noprompt.sh   # B1 단일 실행
    ├── run_b1_batch.sh      # B1 배치
    ├── run_b2_notool.sh     # B2 배치
    ├── run_b4_bare.sh       # B4 배치
    ├── run_all_expb.sh      # 전체 EXP-B 오케스트레이터
    ├── score_expb.py        # EXP-B M-A3 채점 + factorial
    └── results/             # 18개 실행 결과
        ├── B1_S{01,04,07}_qwen3_{32b,latest}/
        ├── B2_S{01,04,07}_qwen3_{32b,latest}/
        └── B4_S{01,04,07}_qwen3_{32b,latest}/
```
