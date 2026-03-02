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

**M-A3 Parameter Fidelity**: 70% (32B) ≡ 70% (8B)
**Determinism**: stddev=0.0 across all 3 trials (temperature=0)

| 난이도 | 32B M-A3 | 8B M-A3 | 비고 |
|--------|:---:|:---:|------|
| Easy (S01-S03) | 72% | 71% | S03 fluid material 무시 |
| Medium (S04-S07) | 69% | 69% | S07 원통형 미지원 |
| Hard (S08-S10) | 69% | 69% | S08/S10 STL 미지원 |

**체계적 오류 패턴 (P1-P6)**:
- P1: Pitch/Roll 혼동 (mvrectsinu → mvrotsinu) — 5/10 시나리오
- P2: TimeMax 과소추정 — 4/10 시나리오
- P3: 원통형 형상 미지원 — 2/10 시나리오
- P4: 비수 유체(오일) 무시 — 2/10 시나리오
- P5: SPHERIC 도메인 지식 부족 — 1/10 시나리오
- P6: 8B tool calling 불안정 — S10에서 관찰

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

**적절성**: 완전히 뒷받침. 70% M-A3는 "작동하지만 개선 여지 있음"을 보여줌.

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
Section 4.1 (EXP-A): "파이프라인이 70% M-A3로 작동"
    → 모델 크기 무관 (32B ≡ 8B)
    → 도구 설계가 병목 (P1-P6 패턴)

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
2. [x] EXP-A stddev=0.0 확인 (32B + 8B)
3. [x] EXP-B B4 완료 (6/6 runs)
4. [x] EXP-B B2 완료 (6/6 runs, 재실행 포함)
5. [x] EXP-B B1 완료 (6/6 runs)
6. [x] EXP-B 결과 채점 및 2×2 factorial 분석
7. [x] EXP-C v2 데이터 확인 (이관 가능)
8. [x] 최종 결과 문서 (이 파일)
9. [ ] Git commit

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
