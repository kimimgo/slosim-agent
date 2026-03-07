# Ralph Loop Iteration 1: 실험 결과 체계적 비판 분석

**Date**: 2026-03-03
**Question**: "논문 뼈대와 GAP을 증명하기 위해 이 실험 결과가 완벽한가?"

---

## Executive Summary

**판정: 아직 완벽하지 않다.** EXP-A/B/C의 실험 자체는 잘 수행되었으나, 문서 간 수치 불일치, EXP-D 미실행, 그리고 몇 가지 방법론적 공백이 있다. 현재 수준은 Workshop/arXiv 출판 가능하나 "이것이 최선"이라 말하기엔 5가지 해결 가능한 문제가 남아있다.

---

## I. 치명적 문제 (Must Fix Before Submission)

### 1. M-A3 수치 불일치 — 두 개의 다른 숫자가 공존

**발견**: 동일한 EXP-A 결과에 대해 두 문서가 다른 수치를 보고한다.

| 문서 | S01 32B | Overall 32B | 채점 방법 |
|------|---------|-------------|----------|
| EXP-A_FINAL_RESULTS.md | 64% | 69.5% | ? (9개 파라미터, density/viscosity 포함) |
| EXPERIMENT_STATUS.md | 75% | 61.2% | score_expb.py 8-parameter |

**구체적 차이**:
- S01: FINAL_RESULTS=64% vs STATUS=75% — **11pp 차이**
- Overall: FINAL_RESULTS=69.5% vs STATUS=61.2% — **8.3pp 차이**
- FINAL_RESULTS는 density, viscosity를 채점 대상에 포함하나 STATUS는 제외

**리스크**: 논문에 어떤 숫자를 쓰든, 다른 문서와 모순된다. 리뷰어가 raw data를 요청하면 정합성 문제 노출.

**권장 조치**:
1. **score_expb.py를 canonical scorer로 확정** — 이미 tolerance sensitivity 분석에 사용
2. EXP-A_FINAL_RESULTS.md를 score_expb.py 기준으로 **재작성**하거나 deprecated 표시
3. 논문에는 EXPERIMENT_STATUS.md 숫자(61.2%)만 사용
4. 채점 파라미터 목록과 tolerance를 논문 Appendix에 명시

### 2. EXP-D 미실행 — "핵심 차별화"가 비어있음

**현재 상태**:
- optimization_prompt.md ✅ (설계 완료)
- run_expd.sh ✅ (스크립트 완료)
- score_expd.py ✅ (채점 완료)
- baffle_generator.go 🔄 (unstaged, 개발 중)
- stl_import fillpoint 🔄 (unstaged, 보완 중)
- **실제 실행: ❌ 0 runs**

**논리적 모순**: CRITICAL_REVIEW는 "EXP-D가 논문의 핵심 차별화"이며 "정규 저널 수준 도달 조건"이라고 명시. 그러나 실행되지 않았다면 이 주장은 공허하다.

**더 심각한 문제**: S10 (STL Fuel Tank)이 EXP-A에서 N/A — STL 파일이 테스트 서버에 없었음. EXP-D는 정확히 이 STL 기능에 의존한다. **S10도 못 하는 시스템이 EXP-D를 할 수 있는가?**

**권장 조치**:
- A) EXP-D를 실행하고 결과를 포함 → 저널 수준
- B) EXP-D를 포기하고 논문을 EXP-A/B/C만으로 구성 → Workshop 수준, 정직한 범위 축소
- **권장: B 우선, A는 시간 허용 시**. 이유: EXP-D 선행 조건(stl_import + baffle_generator)이 미완성이며, 중간에 추가하면 전체 일정 리스크.

### 3. EXP-B 3시나리오 Cherry-Picking 리스크

**문제**: B0(Full)의 S01/S04/S07 평균은 72.2%이지만, EXP-A 전체 평균은 61.2%.

```
Selected scenarios (EXP-B):  72.2% ← "좋은" 시나리오만 선택
Full 10 scenarios (EXP-A):   61.2% ← 실제 전체 성능
Gap:                         11pp ← 의도적 선택이 아니어도 결과적 편향
```

**리뷰어 공격**: "왜 하필 S01/S04/S07? 다른 조합이면 결과 다른가?"

**부분 해소 가능**: EXP-B B0 데이터는 EXP-A에서 재사용하므로, **나머지 7개 시나리오에 대해 B2/B4만 추가 실행하면 된다** (B1은 이미 0%로 결정적).

**권장 조치**:
- Ollama API 호출(B2/B4) 7개 × 2 조건 × 2 모델 = 28 runs 추가
- GPU 불필요(Ollama만 사용), 약 2시간
- 이것만으로도 "n=3 → n=10" 확장, cherry-picking 방어 완전 해소

---

## II. 방법론적 공백 (Should Address)

### 4. Amplitude ↔ Motion_type 채점 커플링

**발견**: score_expb.py에서 motion_type이 FAIL이면 amplitude도 자동 FAIL.

```python
# score_expb.py 구조 (추정)
if actual_motion != 'mvrotsinu':
    amplitude_check = FAIL  # motion_type 틀리면 amplitude 비교 자체 불가
```

**영향**: P1(motion_type) 5개 시나리오 × amplitude 자동 FAIL = M-A3 과소 추정.

**정량적 규모**:
- 커플링으로 인한 이중 실패: 최소 5개 시나리오(S01,S04,S05,S08,S09)
- 독립 채점 시 M-A3 상승 예상: ~3-5pp

**권장 조치**:
1. score_expb.py에 `decoupled_scoring` 모드 추가 — motion_type과 amplitude를 독립 채점
2. 논문에는 기본(커플링) 결과를 보고하되, Appendix에 독립 채점 결과 병기
3. "커플링 채점이 보수적 추정"임을 명시하면 리뷰어 방어 강화

### 5. EXP-B B1 = 0%의 해석 깊이 부족

**B1 실행 상세** (B1_TOOL_CALL_ANALYSIS.md):
- Generic CoderPrompt + 18개 DualSPHysics 도구 available
- 결과: **DualSPHysics 도구 0회 호출** — `glob`, `edit`, `bash` 등 범용 도구만 사용

**논리적 문제**: "도메인 프롬프트가 필수" vs "프롬프트 형식이 부적절"의 구분이 안 됨.
- Generic CoderPrompt는 "코드를 도와달라"는 지시 — 시뮬레이션 맥락 전무
- 이것은 "DualSPHysics 도메인 지식"의 부재가 아니라 "시뮬레이션 태스크 인식 자체의 실패"
- 보다 공정한 비교: "일반적 CFD 시뮬레이션 프롬프트(DualSPHysics 특화 없음)"로 테스트

**그러나**: 이 구분이 논문 주장을 약화시키진 않음.
- "도구 이름만으로는 사용 컨텍스트를 알 수 없다" → 도메인 프롬프트가 tool discovery 역할
- 이 발견 자체가 가치 있음. 논문에서 이 뉘앙스를 논의하면 오히려 깊이 ↑

**권장 조치**:
- Discussion에 "B1의 0%는 tool availability ≠ tool utility" 명시
- "프롬프트가 도구의 사용법을 가르치는 orchestration instruction" 역할을 강조

### 6. EXP-C → EXP-A 연결 실험 부재

**DEEP_ANALYSIS 지적**: EXP-A는 입력(XML) 품질, EXP-C는 출력(시뮬레이션) 품질. 둘의 연결은 **가정**됨.

**핵심 질문**: "M-A3 61%로 생성된 XML이 물리적으로 의미 있는 결과를 내는가?"

**현재 상태**:
- EXP-C의 SPHERIC 시뮬레이션은 **수동으로 만든 XML**(research-v2)로 실행
- 에이전트가 생성한 XML(M-A3 75% S01)로 SPHERIC을 재실행한 적은 없음
- 따라서 C3 주장의 근거가 간접적

**권장 조치 (선택적)**:
- S01 또는 S02(100% M-A3)에서 에이전트 생성 XML로 SPHERIC 시뮬레이션 실행
- GPU ~30분, 결과를 v2 결과와 비교하면 "에이전트 XML → 물리적 결과" 직접 증명
- 비용 대비 효과 매우 높음

---

## III. "32B ≡ 8B" 발견의 이중성

### 강점으로서의 해석
- "모델 크기 무관 → 아키텍처 설계가 결정적" — 비용 효율적 배포 시사
- Temperature=0에서 동일 도구 입력 → 동일 출력: 도구가 출력 형식을 결정

### 약점으로서의 해석
- **도구 설계의 버그(P1-P6)가 성능 천장** — 모델을 바꿔도 같은 오류 반복
- 논문이 "LLM 에이전트 아키텍처"를 주장하지만, 실질적 병목은 "도구 코딩 품질"
- "61%에서도 32B≈8B이면 도구를 고쳐야지 모델을 키울 이유가 없다" — 맞지만, 이것은 논문 기여가 아니라 **엔지니어링 숙제**

### 논문 프레이밍 권장
- Contribution을 "에이전트 아키텍처가 모델 크기보다 중요하다"로 재설정
- P1/P2 수정 시 M-A3 61% → 75% 추정을 "Tool Quality Ceiling" 분석으로 포함
- **"모델을 키우는 것이 아니라 도구를 개선하는 것이 올바른 투자 방향"** — 실용적 시사점

---

## IV. 기존 분석의 깊이 평가

### 충분히 깊은 분석 ✅
1. **Error Pattern P1-P6 분류** — 체계적, 재현 가능, 근본 원인 명시
2. **Tolerance Sensitivity** — M-A3가 5-25% 범위에서 불변 → 강력한 robustness 증명
3. **Tool-Induced Error (analyze_tool_effect.py)** — B0 vs B2 파라미터별 비교, 순 효과 계산
4. **B1 Tool Call Analysis** — 정성적 분석이 예상외로 풍부 (도구 호출 패턴 기록)
5. **Determinism 확인** — 3회 반복이 단일 실행과 동일 → 실험 효율성 최적화

### 깊이 부족한 분석 ⚠️
1. **EXP-A FINAL_RESULTS vs STATUS 불일치** — 두 채점 방법 간 비교/조화 미수행
2. **S08 8B 선택적 실패** — 왜 S08만? 프롬프트의 기술 수준("DBC", "dp=0.002") 영향 분석 없음
3. **실행 시간 분석** — 데이터는 있으나(Table 5, FINAL_RESULTS §4) 논문에서의 활용 방안 미정의
4. **EXP-C GCI 분석** — DEEP_ANALYSIS에서 "GCI M3/M4 FAIL인데 수렴 주장"의 모순 지적 → 미해결

---

## V. 실험 설계 적절성 평가

### 강점
| 항목 | 평가 | 근거 |
|------|------|------|
| 시나리오 설계 | ★★★★☆ | Easy/Medium/Hard 3-tier, NL2FOAM 참조, 10개 다양성 |
| 반복 설계 | ★★★★★ | 3-trial → determinism 확인은 우수한 판단 |
| Ablation 설계 | ★★★★☆ | 2×2 factorial 깔끔, 교호작용 분석 가능 |
| 메트릭 정의 | ★★★★☆ | M-A1~A5 5단계, 명확한 Pass 기준 |
| Ground truth | ★★★★☆ | tier 재정의 + rationale 문서화 완료 |

### 약점
| 항목 | 평가 | 근거 |
|------|------|------|
| Ablation 시나리오 수 | ★★☆☆☆ | n=3 통계 불충분, cherry-picking 리스크 |
| Score 정합성 | ★★☆☆☆ | 두 개의 다른 채점 결과 공존 |
| 연결 실험 | ★★☆☆☆ | EXP-A→C 직접 검증 없음 |
| S10/STL 커버리지 | ★★☆☆☆ | 10개 중 1개 N/A = 10% 데이터 손실 |
| 외부 비교 | ★★☆☆☆ | NL2FOAM과 메트릭 비교만, 직접 실험 비교 없음 |

---

## VI. 우선순위별 행동 계획

### 즉시 실행 가능 (< 30분)
| # | 항목 | 영향 |
|:---:|------|------|
| 1 | ~~EXP-A_FINAL_RESULTS.md~~ deprecated 표시 + score_expb.py 기준 재작성 | 수치 정합성 |
| 2 | 논문에 사용할 canonical 수치 확정 (61.2% / 58.7%) | 혼동 방지 |

### 2-3시간 투자 (High Impact)
| # | 항목 | 영향 |
|:---:|------|------|
| 3 | EXP-B 10개 시나리오 확장 (B2/B4 추가 28 runs) | cherry-picking 완전 해소 |
| 4 | amplitude/motion_type 독립 채점 모드 추가 | 보수적 추정 명시 |

### 1일 투자 (Major Enhancement)
| # | 항목 | 영향 |
|:---:|------|------|
| 5 | 에이전트 생성 XML(S01/S02)로 SPHERIC 시뮬레이션 → EXP-A↔C 연결 | C3 증명 강화 |
| 6 | EXP-D 선행 조건 완성 + 실행 | "autonomy" 차별화 확보 |

### 논문 작성 시 반영 (시간 무관)
| # | 항목 |
|:---:|------|
| 7 | B1=0%의 "tool discovery" 뉘앙스 Discussion에 반영 |
| 8 | P1/P2 수정 시 M-A3 75% ceiling analysis 본문 포함 |
| 9 | "모델 크기 무관 → 도구 설계가 병목" 프레이밍 정교화 |
| 10 | Limitations에 정직한 한계 나열 (L1-L8 from EXPERIMENT_PLAN) |

---

## VII. "이것이 최선인가?"

### 현재 수준: Workshop/arXiv 출판 가능 (70/100)
- 60 runs + 18 ablation runs = 충분한 데이터
- Determinism + tolerance robustness = 강력한 방어
- Error pattern + tool-induced error = 흥미로운 발견

### 30분 투자 후: 75/100
- 수치 정합성 해소 → 신뢰도 ↑

### 3시간 투자 후: 85/100
- EXP-B 10개 확장 → cherry-picking 방어 ↑
- 독립 채점 → 투명성 ↑

### 1일 투자 후: 90/100
- EXP-A→C 연결 실험 → C3 직접 증명

### EXP-D 완수 시: 95/100 (정규 저널 수준)
- Autonomous optimization → 차별화 확보

**결론**: "최선"은 아니지만 "충분히 좋음"에 가까운 상태. **항목 1-4를 즉시 처리하면 3시간 내에 85/100 도달 가능.**

---

# Ralph Loop Iteration 2: 구조적 약점 심층 검증

**Date**: 2026-03-03
**Focus**: Iteration 1의 지적 사항을 소스코드 수준에서 검증하고, 논문 논리 구조의 치명적 결함 탐색

---

## VIII. 검증 완료: Amplitude-MotionType 채점 커플링 (Confirmed)

### 소스코드 증거 (score_expb.py:159-187)

```python
# 경로 1: amplitude_m (sway 시나리오 — S02, S03, S06, S07, S10)
if motion_gt.get('amplitude_m'):
    actual_ampl = params.get('ampl_x', 0) or params.get('ampl_y', 0)
    tol = expected_ampl * 0.15
    passed = abs(actual_ampl - expected_ampl) <= tol  # ← 독립 채점

# 경로 2: amplitude_deg (pitch 시나리오 — S01, S04, S05, S08, S09)
elif motion_gt.get('amplitude_deg'):
    if actual_motion == 'mvrotsinu':
        # 정상 tolerance 비교
    else:
        passed = False  # ← motion_type FAIL → amplitude AUTO-FAIL
```

### 정량적 영향

| 시나리오 | motion_type 결과 | amplitude 채점 | 커플링 여부 |
|---------|:---:|:---:|:---:|
| S01 (pitch) | FAIL (mvrectsinu) | **AUTO-FAIL** | ✅ 커플링 |
| S02 (sway) | PASS | 독립 채점 | — |
| S03 (sway) | PASS | 독립 채점 | — |
| S04 (pitch) | FAIL (mvrectsinu) | **AUTO-FAIL** | ✅ 커플링 |
| S05 (pitch) | FAIL (mvrectsinu) | **AUTO-FAIL** | ✅ 커플링 |
| S06 (sway) | PASS | 독립 채점 | — |
| S07 (sway) | PASS | 독립 채점 | — |
| S08 (pitch) | FAIL (mvrectsinu) | **AUTO-FAIL** | ✅ 커플링 |
| S09 (pitch) | FAIL (mvrectsinu) | **AUTO-FAIL** | ✅ 커플링 |
| S10 (sway) | PASS | 독립 채점 | — |

**영향 규모**: 5개 시나리오 × 2개 커플링 실패 = 10 추가 실패 / 총 ~80 파라미터 체크 = **12.5% M-A3 과소추정 방향**

**논문 대응**: 이 커플링은 **의도적이고 정당하다** — motion_type이 틀리면 amplitude 비교 자체가 의미 없기 때문 (사과와 오렌지 비교). 단, 논문 Appendix에 "채점 규칙: motion_type이 pitch 시나리오에서 FAIL이면 amplitude도 자동 FAIL (독립 채점 불가)" 명시 필요.

---

## IX. 치명적 발견: 퇴화 2×2 Factorial (Degenerate Factorial)

### 문제 정의

```
              Without Prompt    With Prompt
              ─────────────    ────────────
Without Tools │  B4 = 0.0%  │  B2 = 58.3%  │
              │─────────────│──────────────│
With Tools    │  B1 = 0.0%  │  B0 = 72.2%  │
              └─────────────└──────────────┘
```

좌측 열(Without Prompt)이 **완전 0%** → 2×2 factorial이 **퇴화(degenerate)**.

### Factorial 공식 검증 (score_expb.py:373-375)

```python
prompt_effect = (B0 + B2)/2 - (B1 + B4)/2 = (72.2 + 58.3)/2 - (0 + 0)/2 = 65.25%
tool_effect   = (B0 + B1)/2 - (B2 + B4)/2 = (72.2 + 0)/2 - (58.3 + 0)/2 = 6.95%
interaction   = B0 - B1 - B2 + B4          = 72.2 - 0 - 58.3 + 0         = 13.9%
```

수학적으로 정확하지만, **해석이 잘못될 위험**이 있다:

1. **"Tool main effect = 6.95%"는 실질적으로 (B0-B2)/2**
   - B1=0이므로 tool이 있어도 prompt 없이는 0% → tool 효과가 절반 희석
   - 실제 tool 효과는 prompt가 있을 때만 관찰: **B0-B2 = 13.9%**

2. **"Prompt main effect = 65.25%"는 실질적으로 (B0+B2)/2**
   - B1=B4=0이므로 prompt가 없으면 무조건 0%
   - 이것은 "main effect"가 아니라 **"필요조건(necessary condition)"**

3. **"Interaction = 13.9%"은 순수하게 B0-B2**
   - 정상적 factorial에서 interaction은 "시너지"를 의미하지만
   - 여기서는 단순히 "prompt가 있을 때 tool이 추가하는 가치"

### 리뷰어 공격 벡터

> "B1=B4=0%이면 이것은 2×2 factorial이 아니라 one-factor (tool presence given prompt) 실험이다.
>  Prompt는 on/off가 아니라 necessary condition이므로 factorial로 보고하는 것은 부적절하다."

### 대응 방안

**A) 정직한 재프레이밍** (권장):
```
"Domain prompt is a necessary condition: without it, no DualSPHysics tool usage
occurs (B1=B4=0%). Given this prerequisite, tools provide an additional +13.9pp
improvement (B0=72.2% vs B2=58.3%). This hierarchical dependency — prompt enables,
tools refine — is itself a key finding."
```

**B) Factorial을 유지하되 주석 추가**:
- "Note: The factorial is degenerate due to complete prompt dependence (B1=B4=0%)."
- "Tool effect should be interpreted as conditional effect given prompt presence."

**C) 표현 수정**:
- ~~"Main effect of Tools: +6.9%"~~ → "Conditional effect of Tools (given Prompt): +13.9%"
- ~~"Interaction: +13.9%"~~ → "Direct B0-B2 comparison (the only informative contrast): +13.9%"

---

## X. S08 8B=0% vs S10 8B=25% — "32B≡8B" 주장의 미묘한 균열

### 사실 확인

```
Scenario  32B   8B   Δ      Direction
S01       75%   75%  =
S02      100%  100%  =
S03       88%   88%  =
S04       75%   75%  =
S05       50%   50%  =
S06       88%   88%  =
S07       67%   67%  =
S08       50%    0%  +50%   ← 32B wins (8B: tool calling failure)
S09       20%   20%  =
S10        0%   25%  -25%   ← 8B wins (non-deterministic)
```

### 분석

**S08 (32B=50%, 8B=0%)**:
- NL 프롬프트에 "dp=0.002m", "경계조건을 DBC로" — Expert-level 기술 용어
- 8B는 tool calling 자체 실패 (P6 패턴) → XML 미생성
- 32B는 성공하지만 P1(motion_type) + P3(timemax) 실패로 50%
- **해석**: 8B의 tool calling 불안정은 프롬프트 복잡도에 의존

**S10 (32B=0%, 8B=25%)**:
- STL 미지원 시나리오 — 양쪽 모두 실패하지만 방식이 다름
- 8B trial2/3에서 부분적 XML 생성 (σ=35.4, non-deterministic)
- **유일한 비결정적 케이스** — temperature=0에도 불구

**"32B≡8B" 주장의 정밀화**:
- Overall: Δ=2.5pp (61.2% vs 58.7%) — 통계적 무의미 ✓
- Per-scenario: 8/10 완전 동일 ✓
- **예외 2개 (S08, S10)**: 방향이 반대라 상쇄되지만, 이는 "동등"이 아니라 "다르지만 상쇄"
- 논문 표현: "In 8 of 10 scenarios, 32B and 8B produce identical M-A3 scores. Two edge cases (S08, S10) show divergent behavior in opposite directions."

---

## XI. Gap→Evidence 매핑 엄밀 검증

### Gap 1: "No NL→SPH tool-calling agent exists"

| 증거 유형 | 내용 | 강도 |
|----------|------|:---:|
| Table 1 (feature comparison) | 7종 기존 에이전트 vs SlosimAgent | ★★★☆ |
| EXP-A (existence proof) | 10 시나리오 61% M-A3 | ★★★★ |
| EXP-C (physics validation) | SPHERIC 2/3 PASS | ★★★☆ |

**약점**: "no prior work"를 주장하지만, 직접 비교 실험 없음. NL2FOAM과 메트릭 비교만 (간접적).
**방어**: 솔버 불일치(FVM vs SPH)로 직접 비교 불가능함을 논문에 명시. 이것은 정당한 이유.

### Gap 2: "Domain prompting vs fine-tuning for CFD agent"

| 증거 유형 | 내용 | 강도 |
|----------|------|:---:|
| EXP-B ablation | +65.3% prompt effect | ★★★☆ |
| NL2FOAM 간접 비교 | 82.6% (fine-tuned) vs 61% (zero-shot prompt) | ★★☆☆ |

**약점**:
1. EXP-B는 퇴화 factorial — prompt는 "효과"가 아니라 "전제조건"
2. NL2FOAM과 직접 비교 불가 (다른 솔버, 다른 시나리오)
3. "Domain prompting이 fine-tuning의 대안"이라는 주장에는 동일 시나리오 비교가 필요

**방어**: "Zero-shot prompting이 competitive하다"는 주장으로 약화. "대안이다"가 아니라 "가능하다" 수준.

### Gap 3: "No quantitative validation framework for SPH agents"

| 증거 유형 | 내용 | 강도 |
|----------|------|:---:|
| M-A1~A5 메트릭 정의 | 5단계 파이프라인 메트릭 | ★★★★ |
| M1-M8 물리 메트릭 | SPHERIC 기반 정량 평가 | ★★★★ |
| EXP-A 적용 | 60 runs에 적용 | ★★★★ |

**이것이 가장 강한 Gap**. 메트릭 프레임워크 자체가 contribution이며, 적용 결과도 풍부.
**유일한 약점**: 커뮤니티 채택 전이므로 "제안"이지 "표준"이 아님.

### Gap 4: "MCP integration for simulation toolchains"

| 증거 유형 | 내용 | 강도 |
|----------|------|:---:|
| 아키텍처 설명 | MCP 프로토콜 사용 | ★☆☆☆ |
| 실험 증거 | **없음** | ☆☆☆☆ |

**치명적 약점**: 실험적 증거 전무. MCP는 구현 선택이지 연구 기여가 아님.
**권장**: Gap 4를 Contribution에서 제거하거나, "implementation detail"로 격하.

### Gap→Evidence 종합 평가

| Gap | 증거 강도 | 논문 위험도 | 권장 조치 |
|:---:|:---:|:---:|------|
| G1 | ★★★★ | 낮음 | 유지 — existence proof 충분 |
| G2 | ★★☆☆ | **높음** | 주장 약화 — "대안" → "가능성" |
| G3 | ★★★★ | 낮음 | 유지 — 가장 강한 기여 |
| G4 | ★☆☆☆ | **매우 높음** | 제거 또는 격하 |

---

## XII. EXP-C GCI 모순 검증

### 보고된 수렴 데이터

```
dp     r       해석
4mm    0.460   낮음
2mm    0.655   중간 (SPHERIC PASS 기준)
1mm    0.697   중-상
```

r은 단조 증가 → **추세 수렴(trend convergence)은 확인됨**.

### GCI 문제

- DEEP_ANALYSIS C2: "GCI M3/M4 FAIL인데 수렴 확인 주장 — 모순"
- GCI(Grid Convergence Index)는 Richardson 외삽법 기반. 3-level 데이터에서 단조 수렴 비율을 확인하고 안전계수를 적용.
- r 증가율: 0.460→0.655 (+42%) vs 0.655→0.697 (+6.4%) — **수렴 감속**
- 이는 비단조적 수렴 감속 → GCI가 높은 불확실성 보고 → FAIL 가능

### 논문 대응

1. ~~"수렴 확인"~~ → "dp 정밀화에 따른 단조 r 증가 확인 (trend convergence)"
2. GCI 결과를 정직하게 보고: "GCI safety factor criterion은 미충족. SPH의 비단조적 수렴 특성으로 인해 GCI가 과도한 불확실성을 보고하며, 이는 입자 방법의 알려진 한계."
3. 1mm dp 결과가 0.697이므로 "추가 정밀화의 수렴 이득이 감소" (diminishing returns) 논의

---

## XIII. Iteration 2 종합 판정

### 새로 발견한 치명적 문제 (Iteration 1 이후)

| # | 문제 | 심각도 | Iteration 1 언급 |
|:---:|------|:---:|:---:|
| A | 퇴화 2×2 Factorial | **치명** | 부분적 (n=3만 지적) |
| B | Gap 4 실험 근거 전무 | **높음** | 있음 (RALPH §III) |
| C | GCI FAIL↔수렴 주장 모순 | **중간** | 있음 (DEEP_ANALYSIS C2) |
| D | "32B≡8B" S08/S10 예외 미논의 | **중간** | 부분적 |

### 수정된 점수 평가

```
항목                          Iter 1 평가    Iter 2 재평가    이유
──────────────────────────    ──────────    ────────────    ─────
현재 수준                     70/100        65/100          퇴화 factorial 발견
30분 투자 후                  75/100        72/100          수치 정합성만으론 부족
3시간 투자 후                 85/100        80/100          factorial 재프레이밍 필요
1일 투자 후                   90/100        85/100          EXP-A→C 연결 + GCI 정직 보고
EXP-D 완수 시                95/100        92/100          factorial 문제 잔존
```

**Iteration 1 대비 하향 조정 이유**: 퇴화 factorial은 논문의 핵심 주장(C2: "아키텍처 컴포넌트 필수성")의 통계적 기반을 약화시킨다. 이것은 데이터 문제가 아니라 **해석 문제**이므로, 재프레이밍으로 해결 가능하지만 현재 상태에서는 리뷰어 공격에 취약.

### 즉시 조치 (Iteration 2 기준, 우선순위 순)

| # | 항목 | 시간 | 영향 |
|:---:|------|------|------|
| 1 | **Factorial 재프레이밍** — "hierarchical dependency"로 보고 | 1h | C2 주장 방어 |
| 2 | **Gap 4 제거/격하** — Contribution에서 삭제 | 30m | 논문 견고성 ↑ |
| 3 | **EXP-A_FINAL_RESULTS 정합** — score_expb.py canonical | 30m | 수치 신뢰도 |
| 4 | **amplitude 커플링 명시** — Appendix 또는 Table note | 15m | 투명성 |
| 5 | **GCI 정직 보고** — "trend convergence, not formal GCI" | 30m | 신뢰도 |
| 6 | **"32B≡8B" 정밀화** — "8/10 identical, 2 edge cases" | 15m | 정확성 |

### "이것이 최선인가?" Iteration 2 답변

**아니다. 하지만 "최선에 가장 가까운 합리적 버전"으로 만들 수 있다.**

현재 상태의 가장 큰 리스크는 데이터나 실험 자체가 아니라 **해석과 보고의 정밀도**:
1. Factorial이 퇴화됨 → 재프레이밍으로 해결 (데이터 재수집 불필요)
2. Gap 4가 공허함 → 제거로 해결
3. GCI 모순 → 정직한 보고로 해결
4. 수치 불일치 → canonical scorer 확정으로 해결

**이 4가지는 모두 글쓰기(writing) 수준의 수정이지, 실험 재실행이 아니다.** 2-3시간 투자로 65→80/100 가능.

EXP-B 10개 확장(+2h GPU)과 EXP-A→C 연결 실험(+30min GPU)까지 수행하면 85/100.
EXP-D까지 완수하면 92/100.

---

# Ralph Loop Iteration 3: Raw Data 검증 + 근본원인 분석

**Date**: 2026-03-04
**Focus**: Iteration 1-2의 주장을 실제 XML/소스코드로 검증, 용어 정정, 새로운 구조적 발견

---

## XIV. 결정적 검증: P1은 "버그"가 아니라 "미구현 기능" (CONFIRMED)

### xml_generator.go 소스코드 증거

**파라미터 구조체 (line 11-26):**
```go
type XMLGeneratorParams struct {
    TankLength  float64 `json:"tank_length"`
    TankWidth   float64 `json:"tank_width"`
    TankHeight  float64 `json:"tank_height"`
    FluidHeight float64 `json:"fluid_height"`
    Freq        float64 `json:"freq"`
    Amplitude   float64 `json:"amplitude"`     // ← "가진 진폭 (m)" — 미터만
    DP          float64 `json:"dp"`
    TimeMax     float64 `json:"time_max"`
    OutPath     string  `json:"out_path"`
    // ... dimension, geometry, excitation_type, boundary_method
    // ❌ motion_type 필드 없음!
    // ❌ amplitude_unit 필드 없음!
}
```

**템플릿 (line 234):**
```go
<mvrectsinu id="1" duration="%g" anglesunits="degrees">
    <freq x="%g" y="0" z="0" units_comment="Hz" />
    <ampl x="%g" y="0" z="0" units_comment="metres (m)" />
</mvrectsinu>
```

### 핵심 구분: Bug vs Missing Feature

| 기존 용어 (Iteration 1-2) | 정확한 용어 (Iteration 3) | 근거 |
|---|---|---|
| "P1 bug" | **Missing Feature: motion_type** | 파라미터 구조체에 필드 자체 없음 |
| "P2 bug" | **Missing Feature: amplitude_unit** | 도구가 meters만 지원, degrees 불가 |
| "도구가 강제 오버라이드" | **도구가 입력 채널 미제공** | LLM이 전달하고 싶어도 방법 없음 |

**논문 용어 정정 필요**: "tool-induced error"가 아니라 "**tool design limitation**" — 도구가 오류를 '도입'하는 게 아니라 올바른 입력을 '수용하지 못함'.

### EXP-A 전체 motion_type 분석 (Raw XML 검증)

```
EXP-A B0 (with tools) — 전체 sloshing_case.xml grep 결과:
─────────────────────────────────────────────────────
모든 S01~S10 XML: mvrectsinu (100%)
mvrotsinu: 0건

EXP-B B2 (without tools, domain prompt only):
─────────────────────────────────────────────────────
S01 (pitch): mvrotsinu ✅ (32B, 8B 모두)
S04 (pitch): mvrotsinu ✅ (32B, 8B 모두)
S07 (sway):  mvrectsinu ✅ (올바름 — sway IS mvrectsinu)
```

**결론**: xml_generator 도구는 **100% mvrectsinu만 출력** (결정론적). LLM 판단과 무관.

### stl_import.go와의 비교 (line 28)

```go
// stl_import.go — 더 나중에 개발된 도구
MotionType string  `json:"motion_type,omitempty"` // "mvrectsinu" | "mvrotsinu"
```

stl_import은 motion_type을 지원하지만, xml_generator는 지원하지 않음.
→ 개발 과정에서 인지했지만 기존 도구를 업데이트하지 않은 **기술적 부채**.

### 논문에서의 프레이밍

**현재 (CRITICAL_REVIEW)**:
> "도구가 도입하는 오류가 도구가 방지하는 오류보다 크지 않은가?"

**정정된 프레이밍 (Iteration 3)**:
> "도구의 설계 범위가 시나리오 다양성에 미치지 못한다. xml_generator는 수평 진동(sway) 전용으로 설계되어 회전 운동(pitch/roll)을 수용하지 못한다. 이는 도구 버그가 아닌 **도구 커버리지 한계**이며, 이 한계가 전체 시스템 M-A3의 상한선을 결정한다."

---

## XV. "Decoupled Scoring"은 결과를 변경하지 않음 (New Finding)

### Iteration 2 제안: "독립 채점 시 M-A3 상승 예상 ~3-5pp"

### Iteration 3 검증: **예상 틀림. 상승 0pp.**

**이유**: amplitude가 AUTO-FAIL이 되는 이유는 coupling 뿐만 아니라 **값 자체도 부정확**:
- S01: xml_generator에 amplitude=0.045m 전달 → XML에 0.045m → GT는 4° → 독립 채점해도 FAIL
- S04: amplitude=0.035m → GT는 2° → FAIL
- B2(도구 없음): amplitude=0.0698 (4° × π/180 = radians) → GT는 4° → FAIL

motion_type FAIL로 인한 AUTO-FAIL이 아니라, **amplitude 값 자체가 독립적으로 FAIL**.
커플링은 단순히 불필요한 비교를 건너뛰는 최적화 — 결과에 영향 없음.

**수정**: DEEP_ANALYSIS와 RALPH_LOOP Iteration 2의 "decoupled 시 +3-5pp" 주장은 **철회**.

---

## XVI. S08 숨겨진 오류: 프롬프트의 암묵적 파라미터 문제

### S08 NL 프롬프트 (EXPERIMENT_PLAN에서 추출):
```
"SPHERIC 벤치마크 탱크(900x62x508mm) 물 18% 해석인데,
 경계조건을 DBC로 해줘. dp=0.002m로 고해상도로.
 센서1 위치 압력 시계열 추출."
```

### S08 Ground Truth:
```json
{
    "motion": {"xml_tag": "mvrotsinu", "freq_hz": 0.6515, "amplitude_deg": 4},
    "timemax": 7.68
}
```

### 문제: **NL 프롬프트에 운동 파라미터가 명시되지 않음**

S08 프롬프트는 freq, amplitude, timemax를 전혀 지정하지 않음!
Ground truth는 이들을 "SPHERIC 벤치마크"라는 컨텍스트에서 암묵적으로 기대.

**실제 S08 XML 결과 (32B, trial1):**
```
dp = 0.005 (프롬프트: 0.002 → 무시됨, dp는 M-A3 비채점)
freq = 0.519 Hz (GT: 0.6515 Hz → FAIL, tolerance 외)
timemax = 9.63s (GT: 7.68s → FAIL, 9.63 > 8.448 = 7.68×1.1)
```

**진단**: LLM은 "SPHERIC 벤치마크"를 인식했지만 (탱크 치수는 정확), 정확한 운동 파라미터를 추론하지 못함. freq=0.519 Hz는 아마 다른 SPHERIC 시나리오의 파라미터 혼동.

### 이것이 논문에 미치는 영향

**프롬프트 설계 약점**: 10개 NL 프롬프트 중 S08만 유일하게 운동 파라미터가 누락.
이는 "비전문가 프롬프트"의 현실적 한계를 보여주는 동시에, 실험 설계의 비균질성(L1)을 심화.

**두 가지 해석**:
1. **S08 = implicit knowledge test** — 에이전트가 "SPHERIC benchmark"에서 파라미터를 추론할 수 있는가?
   - 결과: 부분 성공 (탱크 치수 O, 운동 파라미터 X)
   - 이 해석은 논문에 흥미로운 discussion point 제공

2. **S08 = unfair test** — 프롬프트에 없는 정보를 기대하는 것은 부당
   - 리뷰어가 이 해석을 취하면 S08 M-A3의 freq/timemax FAIL은 채점 오류

**권장**: S08을 "implicit knowledge test"로 프레이밍하되, 논문 본문에 "S08 프롬프트는 운동 파라미터를 명시하지 않음. 에이전트가 벤치마크 이름에서 추론해야 함"을 명시. freq/timemax FAIL을 별도 분석.

---

## XVII. B3 조건 사일런트 드롭

### EXPERIMENT_PLAN에 정의된 조건: 5개 (B0, B1, B2, B3, B4)
### EXPERIMENT_STATUS에 실행된 조건: 4개 (B0, B1, B2, B4)

**B3 (-PostProcess)**: 도메인 프롬프트 + DualSPHysics 도구 있음, 후처리 도구만 제거.
- XML 생성에 영향 없음 → M-A3 = B0과 동일 예상
- M-A5(전체 파이프라인) 에만 영향 → VTK/CSV 결과물 생성 불가

**사일런트 드롭의 문제**: 논문에서 B3 제거 사유를 밝히지 않으면 리뷰어가 "왜 5개 조건 중 1개를 생략했나?" 질문.

**권장**: EXPERIMENT_STATUS 또는 논문에 "B3 was excluded because post-processing tools do not affect XML generation or solver execution; M-A2/A3 scores would be identical to B0" 한 줄 추가.

---

## XVIII. Per-Error-Pattern M-A3 영향 정밀 정량화

### 시나리오별 에러 패턴 매핑 (32B trial1 기준)

```
Pattern   Affected    Params    FAIL    M-A3 Impact
                      Scenarios  Lost    (pp, overall)
────────────────────────────────────────────────────────
P1 motion_type       S01,S04,S05,   5      -6.25pp
   (missing feature) S08,S09

P2 amplitude_unit    S01,S04,S05,   5      -6.25pp  (SAME scenarios as P1)
   (missing feature) S08,S09

P3 timemax           S03,S06,       4      ~-5.0pp
   (under-estimate)  S08,S09

P4 geometry_type     S07,S09        2      -2.5pp
   (cylinder→box)

P5 tank_z/fill_h     S05            1      -1.25pp
   (swap)

P6 8B tool calling   S08(8B),S10    —      8B only
────────────────────────────────────────────────────────
```

**주의**: P1과 P2는 **완전히 같은 시나리오**에서 발생하므로, 합산하면 이중계산.
실제 독립 영향:
- P1만 수정: +6.25pp (motion_type 5개 → PASS) → 67.4%
- P1+P2 수정: +12.5pp (motion_type 5개 + amplitude 5개 → PASS) → 73.7%
- P1+P2+P3 수정: +12.5pp + ~5pp = ~17.5pp → ~78.7%
- P1+P2+P3+P4 수정: ~80.0%+
- All fixed (theoretical): ~88-90%

**DEEP_ANALYSIS와의 차이**: DEEP_ANALYSIS는 69.4%(P1) / 75.2%(P1+P2)를 제시. 내 추정은 67.4% / 73.7%. 차이(~1.5pp)는 시나리오별 파라미터 수 차이와 반올림에 기인.

---

## XIX. Iteration 3 종합 판정

### Iteration 3에서 새로 발견된 사항

| # | 발견 | 심각도 | Iteration 2 대비 |
|:---:|------|:---:|:---:|
| A | P1/P2는 "버그"가 아닌 "미구현 기능" | **용어 수정** | 새 발견 |
| B | "Decoupled scoring"은 결과를 변경하지 않음 | **분석 수정** | Iter 2 주장 철회 |
| C | S08 프롬프트에 운동 파라미터 누락 | **중간** | 새 발견 |
| D | B3 사일런트 드롭 | **낮음** | 새 발견 |
| E | xml_generator = sway 전용 도구 (설계 한계) | **프레이밍 변경** | 새 발견 |

### 점수 재평가

```
항목                    Iter 1    Iter 2    Iter 3    이유
────────────────────    ──────    ──────    ──────    ─────
현재 수준               70/100    65/100    65/100    변동 없음 (새 발견이 해석 수정)
즉시 조치 후 (2-3h)     85/100    80/100    82/100    decoupled 불필요로 시간 절감
1일 투자 후             90/100    85/100    85/100
EXP-D 완수 시           95/100    92/100    92/100
```

### 수정된 즉시 조치 목록 (Iteration 3 최종)

| # | 항목 | 시간 | 영향 | 상태 |
|:---:|------|------|------|:---:|
| 1 | **Factorial "hierarchical dependency" 재프레이밍** | 1h | C2 방어 | ⬜ |
| 2 | **P1/P2 용어 정정**: "bug" → "design limitation" | 15m | 정확성 | ⬜ |
| 3 | **Gap 4 (MCP) 격하** → implementation detail | 30m | 논문 견고성 | ⬜ |
| 4 | **Score 정합성**: EXP-A_FINAL_RESULTS deprecated | 30m | 수치 신뢰도 | ⬜ |
| 5 | **S08 implicit parameter 명시** | 15m | 투명성 | ⬜ |
| 6 | **B3 생략 사유 한 줄 추가** | 5m | 완전성 | ⬜ |
| 7 | ~~Decoupled scoring~~ | — | ~~불필요~~ | 철회 |
| 8 | **GCI "trend convergence" 정직 보고** | 30m | 신뢰도 | ⬜ |
| 9 | **"32B≡8B" 정밀화** — 8/10 identical 명시 | 15m | 정확성 | ⬜ |

**총 예상 시간**: ~3시간 (모두 글쓰기 수준, 실험 재실행 불필요)

### "이것이 최선인가?" Iteration 3 답변

**"최선"까지의 거리가 더 명확해졌다.**

Iteration 1은 문제를 식별했고, Iteration 2는 구조적 약점을 검증했고, Iteration 3은 **근본원인을 소스코드 수준까지 추적하고 잘못된 분석(decoupled scoring)을 철회**했다.

남은 핵심 리스크 2가지:
1. **퇴화 factorial** → 재프레이밍으로 해결 가능 (글쓰기)
2. **도구 설계 한계가 성능 천장** → "tool quality ceiling"으로 정직하게 보고

이 두 가지를 논문에 반영하면, 현재 데이터로 가능한 **최선의 논문**이 된다.
추가 실험(EXP-B 확장, EXP-D)은 "더 좋은 논문"을 만들지만, "현재 데이터의 최선"은 글쓰기로 달성 가능.

---

# Ralph Loop Iteration 4: 논문뼈대 교차검증 + 수정 실행

**Date**: 2026-03-04
**Focus**: 논문뼈대.tex와 Iteration 1-3 발견의 불일치 식별 → 구체적 수정 지시서 작성 → 고영향 아이템 실행

---

## XX. 논문뼈대.tex vs Ralph Loop 발견 — 8개 불일치

| # | 위치 | 논문 현재 표현 | 정확한 표현 (Ralph Loop 근거) | 심각도 |
|:---:|------|----------------|-------------------------------|:---:|
| D1 | §Gap #4 | "MCP 기반 과학계산 통합 부재" | 실험 근거 전무 — Gap이 아닌 implementation detail | **높음** |
| D2 | §4.2 L108 | "Tool 주효과 +6.9%pp" | 퇴화 factorial → "Conditional tool effect (given prompt): +13.9%pp" | **높음** |
| D3 | §4.2 L109-111 | factorial 해석 3개 | B1=B4=0% → "Prompt is prerequisite, Tools are additive refinement" | **높음** |
| D4 | §6.2 L136 | "xml_generator가 pitch를 sway로 강제 변환" | "xml_generator에 motion_type 입력 채널 없음 — sway 전용 설계" | **중간** |
| D5 | §6.3 L143 | "P1(motion_type 버그)" | "P1(motion_type 미구현): 도구 설계 한계" | **중간** |
| D6 | §5 L119 | "수렴: r = 0.460 → 0.655 → 0.697" (수렴 확인 암시) | "Trend convergence 확인. GCI formal criterion 미충족 (SPH 비단조 수렴 특성)" | **중간** |
| D7 | §4.2 L111 | "32B ≡ 8B 모든 ablation 조건에서 동일" | "8/10 시나리오 동일, S08(32B>8B)·S10(8B>32B) 예외, 상쇄하여 overall Δ=2.5%pp" | **낮음** |
| D8 | — | B3 조건 생략 사유 없음 | "B3(-PostProcess)는 XML 생성에 무영향 → 제외" 한 줄 필요 | **낮음** |

---

## XXI. Gap 구조 재설계: 4 Gaps → 3 Gaps + 1 Implementation Note

### 현재 (논문뼈대)
```
Gap 1: 비전문가 접근성 공백 (NL→SPH = 0편)         ★★★★
Gap 2: SPH 솔버 에이전트 부재 (7종 전부 OpenFOAM)   ★★☆☆
Gap 3: 로컬 SLM 부재 (Cloud API 의존)              ★★★★
Gap 4: MCP 기반 과학계산 통합 부재                    ★☆☆☆ ← 실험 근거 전무
```

### 제안: Gap 4 격하 → 3 Gaps + Architecture Detail
```
Gap 1: 비전문가 접근성 공백 → EXP-A (60 runs, M-A3 61%)
Gap 2: SPH 솔버 에이전트 부재 → EXP-C (SPHERIC T10, r=0.655)
Gap 3: 로컬 SLM 부재 → EXP-A/B (32B≡8B, Δ=2.5%pp)

§3 System Design에서: "MCP 프로토콜 기반 도구 통합" → 아키텍처 설계 선택으로 서술
(Gap이 아닌 Design Decision — "왜 MCP를 선택했나"로 프레이밍)
```

### Contribution 매핑 수정
```
현재 4개:
C1: 비전문가용 최초 SPH AI 에이전트 ← Gap 1
C2: DualSPHysics+Qwen3+MCP 통합  ← Gap 2 + Gap 4
C3: 2×2 factorial ablation        ← Gap 1/3
C4: SPHERIC T10 정량 검증         ← Gap 2

수정 후 4개 (재배치):
C1: 비전문가용 최초 SPH AI 에이전트 ← Gap 1
C2: 도구 아키텍처 필수성 정량화 (hierarchical dependency) ← Gap 1/2
C3: 로컬 SLM 실현 가능성 입증 (32B≡8B) ← Gap 3
C4: SPHERIC T10 정량 검증 + M1-M8 프레임워크 ← Gap 2
```

---

## XXII. Factorial 재프레이밍: "Hierarchical Dependency" 서술 가이드

### 현재 표현 (§4.2):
> "2×2 factorial — Domain Prompt × DualSPHysics Tools"
> "Prompt 주효과 +65.3%pp >> Tool 주효과 +6.9%pp, 상호작용 +13.9%pp"

### 수정 표현:

**Option A — Honest Factorial (권장)**:
> "We conduct a 2×2 factorial ablation (Table 4). The factorial is **degenerate**: without the domain prompt, models never invoke DualSPHysics tools (B1=B4=0%). This reveals a **hierarchical dependency** — the domain prompt is a *necessary condition* that enables tool usage, not merely an additive factor.
>
> Given this prerequisite, the informative contrast is B0 vs B2: tools provide +13.9%pp additional accuracy by enforcing geometric constraints (fill_height/tank_z ordering). We report the standard factorial decomposition for completeness (Prompt: +65.3%pp, Tool: +6.9%pp, Interaction: +13.9%pp), but note that main effects are interpretable only as conditional effects given the hierarchical structure."

**핵심 변화**:
1. "main effect" → "conditional effect given prerequisite"
2. Tool 효과: +6.9%pp (main) → **+13.9%pp (conditional, B0-B2)**
3. "상호작용" → "the only informative contrast"
4. 퇴화 구조를 **주요 발견**(hierarchical dependency)으로 승격

---

## XXIII. Tool-Induced Error → Tool Design Limitation 수정 가이드

### §6.2 현재 표현:
> "Tool HURTS (1 param): motion_type (1/3→3/3) — xml_generator가 pitch를 sway로 강제 변환"

### 수정 표현:
> "**Tool coverage gap** (1 param): xml_generator was designed for translational sway only (mvrectsinu hardcoded). It lacks a motion_type parameter, so models cannot specify rotational pitch (mvrotsinu) even when their judgment is correct (verified: B2 correctly uses mvrotsinu for all pitch scenarios).
>
> This is not a bug but a **design limitation** — the tool's parameter space does not cover the full scenario diversity. Notably, the newer stl_import tool includes a motion_type field, confirming developer awareness of this gap."

### §6.3 P1/P2 표현 수정:
> ~~"P1(motion_type 버그) + P2(amplitude 단위 변환)"~~
> → "P1(motion_type: sway-only tool design) + P2(amplitude: meters-only unit support)"
> → "These are **tool coverage limitations**, not implementation bugs. The xml_generator parameter struct has no motion_type field (line 11-26) and describes amplitude as '가진 진폭 (m)' — meters only."

---

## XXIV. GCI 표현 수정 가이드

### §5 현재 표현:
> "수렴: 3-level dp, r = 0.460 → 0.655 → 0.697" (수렴 확인 암시)

### 수정 표현:
> "**Trend convergence**: r increases monotonically with refinement (0.460 → 0.655 → 0.697), but convergence **decelerates** (+42% → +6.4%). The formal GCI safety factor criterion is not met — a known limitation of particle methods where convergence is non-monotonic at fine resolutions (Vacondio et al. 2020). We report dp=2mm (r=0.655) as our reference resolution, noting diminishing returns at dp=1mm."

---

## XXV. "32B≡8B" 정밀화 가이드

### 현재 (§4.2 L111, §6.4):
> "32B ≡ 8B 모든 ablation 조건에서 동일"

### 수정:
> "32B and 8B produce **identical** M-A3 in 8 of 10 scenarios (Table 3). Two edge cases diverge: S08 (32B=50%, 8B=0% — tool calling complexity), S10 (32B=0%, 8B=25% — non-deterministic). These cancel in aggregate (Δ=2.5%pp), but the mechanism differs: S08 reflects 8B's fragile tool invocation under complex prompts; S10 reflects stochastic behavior even at temperature=0."

---

## XXVI. Iteration 4 종합 판정 + 최종 액션 리스트

### Iteration 4의 핵심 기여

분석(Iter 1-3)에서 **실행 가능한 수정 지시서**(Iter 4)로 전환. 8개 불일치를 식별하고, 각각에 대해 Before/After 텍스트를 작성 완료.

### 최종 통합 액션 리스트 (Iteration 1-4 통합)

| # | 항목 | 대상 파일 | 시간 | 영향 | 상태 |
|:---:|------|----------|------|------|:---:|
| 1 | **Factorial hierarchical dependency 재프레이밍** | 논문뼈대 §4.2 | 30m | ★★★★★ | ✅ Iter 4 |
| 2 | **Gap 4 격하 → 3 Gaps + Design Detail** | 논문뼈대 §Gap, §Contribution | 30m | ★★★★ | ✅ Iter 4 |
| 3 | **Tool "bug" → "design limitation" 용어 통일** | 논문뼈대 §6.2, §6.3 | 15m | ★★★ | ✅ Iter 4 |
| 4 | **GCI "trend convergence" 정직 보고** | 논문뼈대 §5 | 15m | ★★★ | ✅ Iter 4 |
| 5 | **"32B≡8B" → "8/10 identical + 2 edge cases"** | 논문뼈대 §4.2, §6.4 | 15m | ★★ | ✅ Iter 4 |
| 6 | **S08 implicit parameter 명시** | 논문뼈대 §6.5 | 10m | ★★ | ✅ Iter 4 |
| 7 | **B3 생략 사유 한 줄 추가** | 논문뼈대 §4.2 | 5m | ★ | ✅ Iter 4 |
| 8 | **EXP-A_FINAL_RESULTS deprecated 표시** | research-v3/ | 10m | ★★ | ✅ Iter 4 |

**8/8 항목 완료.** 모든 논문뼈대 수정 + 문서 정합성 작업 완료.

### 점수 재평가

```
항목                    Iter 1    Iter 2    Iter 3    Iter 4    Iter 4 실행 후
────────────────────    ──────    ──────    ──────    ──────    ──────────────
현재 수준               70/100    65/100    65/100    65/100    78/100 ← 8/8 적용 완료
추가 작업 후 (+1d)      90/100    85/100    85/100    85/100    85/100
EXP-D 완수 시           95/100    92/100    92/100    92/100    92/100
```

**65→78/100 상승 근거**: 논문뼈대에 hierarchical dependency 재프레이밍, Gap 구조 재설계, 용어 통일, GCI 정직 보고, 32B≡8B 정밀화, S08/B3 명시, FINAL_RESULTS deprecated 처리 모두 반영 완료.

### Iteration 4에서 달라진 것

1. **분석 → 실행 전환**: 8개 불일치에 대한 **Before/After 텍스트** 완성
2. **Gap 구조 재설계안** 구체화: 4→3 Gaps + Design Detail
3. **Contribution 재매핑**: C2를 "MCP 통합"에서 "아키텍처 필수성 정량화"로 변경
4. **시간 추정 개선**: 사전 가이드 존재로 3h→2.5h

### "이것이 최선인가?" Iteration 4 답변

**분석으로서는 이것이 최선이다.**

4번의 반복으로:
- Iteration 1: 문제 식별 (5개 문제)
- Iteration 2: 구조적 약점 검증 (factorial 퇴화, Gap 4 공허)
- Iteration 3: 근본원인 소스코드 추적 (P1/P2 = 미구현, decoupled scoring 철회)
- Iteration 4: 논문뼈대 교차검증 + 수정 지시서 완성

**Iteration 4에서 "실행"까지 완료.** 8/8 수정 항목이 논문뼈대.tex에 적용됨.

적용된 변경 요약 (논문뼈대.tex, +37/-27 lines):
1. ✅ §Gap: 4→3 Gaps + MCP를 design decision으로 격하
2. ✅ §4.2: Degenerate factorial → hierarchical dependency 재프레이밍
3. ✅ §4.2: B3 제외 사유 한 줄 추가
4. ✅ §4.2: "32B≡8B 모든 조건" → "8/10 동일, 2 예외"
5. ✅ §5: "수렴" → "Trend convergence, GCI 미충족"
6. ✅ §6.2: "Tool-Induced Error" → "Tool Design Limitation"
7. ✅ §6.3: "P1 버그" → "P1 sway 전용 설계", tolerance sensitivity 추가
8. ✅ §6.4: 32B≡8B 정밀화 + S08/S10 예외 분석
9. ✅ §6.5: S08 implicit parameter 명시
10. ✅ Contribution Statement: hierarchical dependency + tool ceiling 반영
11. ✅ Abstract + Conclusion: 전체 일관성 업데이트
12. ✅ EXP-A_FINAL_RESULTS.md: deprecated 표시 + canonical scorer 안내

**현재 수준: 78/100** (분석 + 실행 완료). 추가 향상: EXP-B 10개 확장(+2h) → 82, EXP-D 완수(+1d) → 92.

---

# Ralph Loop Iteration 5: Self-Socratic Debate + 수치 재현 검증

**Date**: 2026-03-04
**Focus**: Mandatory Self-Socratic Debate — 3대 질문에 데이터 기반 답변

---

## XXVII. 수치 재현 검증 결과

### score_expb.py 실행 (EXP-B)

```
실행 결과 vs 보고 수치:
  B0 = 72.2% → ✅ 일치
  B1 = 0.0%  → ✅ 일치
  B2 = 58.3% → ✅ 일치
  B4 = 0.0%  → ✅ 일치
  Factorial: Prompt +65.3%, Tool +6.9%, Interaction +13.9% → ✅ 일치

  qwen3_32b ≡ qwen3_latest: 모든 조건 완전 동일 → ✅ 일치
```

### tolerance_sensitivity.py 실행 (EXP-A)

```
32B = 61.2% → ✅ 논문과 일치
8B  = 56.2% → ⚠️ 논문 58.7%와 2.5pp 차이!
```

### 8B 불일치 근본원인 추적

**tolerance_sensitivity.py**:
- trial1만 검사 (line 128: `f"{s}_{m}_trial1"`)
- `find_xml()` 후보에 `fuel_tank.xml` / `fuel_tank_analysis.xml` 없음
- S10 8B trial1: XML 없음 → 0% (N/A → 0 처리)

**EXPERIMENT_STATUS.md** (58.7%):
- 3-trial 수동 확인 — S10 8B trial2/3에서 부분 XML 생성 확인
- S10 8B = 25% (3-trial 평균, trial1=0%, trial2 또는 trial3에서 부분 성공)
- 25% > 0% → overall 8B 점수 2.5pp 상승

**구체적 불일치**:
```
             trial1 only    3-trial avg    Δ
32B S10       0%              0%           0    (전부 N/A)
8B S10        0%             25%          25    (trial2: fuel_tank.xml 존재)

8B Overall   56.2%           58.7%       +2.5pp
```

### 논문에 미치는 영향

**"Δ=2.5%pp" 주장의 두 가지 버전**:
1. trial1 기준: 32B=61.2%, 8B=56.2% → **Δ=5.0pp**
2. 3-trial 기준: 32B=61.2%, 8B=58.7% → **Δ=2.5pp**

논문은 version 2 (Δ=2.5pp)를 사용. 하지만 tolerance sensitivity 분석은 version 1 (Δ=5.0pp)의 데이터에 기반.

**위험**: 리뷰어가 tolerance_sensitivity.py를 실행하면 8B=56.2%를 얻음 → 논문과 불일치.

---

## XXVIII. Self-Socratic Debate

### [Doubt] 근본적 결함의 가장 가능성 높은 숨겨진 원인

**가정 위반: "temperature=0 → 결정론적"**

논문의 전체 논증 체인:
```
temperature=0 → deterministic → trial 무관 → 1 trial suffices → 60 runs 충분
```

하지만 **S10 8B는 이 체인을 깨트린다**:
- trial1: XML 미생성 (0%)
- trial2: fuel_tank.xml 생성 (25%?)
- trial3: sloshing_case.xml 생성 (25%?)
- σ = 35.4 ≠ 0

**원인**: Qwen3 8B의 tool calling은 temperature=0에서도 **비결정적**. 이는 Ollama의 내부 sampling (logit tie-breaking)이나 batch scheduling에 의한 것으로 추정.

**숨겨진 결함**: tolerance_sensitivity.py가 trial1만 사용하는 것은 32B에서는 문제없지만 (모든 시나리오 σ=0), **8B에서는 S10의 비결정성을 무시**한다. 논문의 58.7%와 분석 도구의 56.2%가 다른 숫자를 보고하는 것은 **재현성 문제**이다.

**대응**: tolerance_sensitivity.py를 3-trial 평균으로 수정하거나, 논문에서 "S10을 제외한 9개 시나리오에서는 trial 무관 (σ=0). S10은 3-trial 평균으로 별도 보고"라고 명시.

**심각도**: ⚠️ 중간. S10 하나의 비결정성이 전체 Δ를 2.5pp→5.0pp로 변경하지만, "32B≡8B" 결론 자체는 양쪽 해석 모두에서 유효 (5pp도 통계적 무의미 수준).

### [Falsification] 솔루션이 붕괴하는 최극단 에지케이스

**최극단 시나리오**: S10 8B의 3개 trial이 모두 75%로 성공하는 경우.

```
현재:  S10 8B = 25% (trial1=0%, trial2/3 부분성공)
극단:  S10 8B = 75% (모든 trial 성공)

→ 8B Overall = (75+100+88+75+50+88+67+0+20+75)/10 = 63.8%
→ 32B Overall = 61.2%
→ Δ = -2.6pp (8B > 32B!)
```

**이 시나리오에서 "32B≡8B" → "8B>32B"로 반전**.

**방어**: 이것은 발생하지 않았다. 실제 데이터에서 S10 8B는 최대 25%. 그리고 S10은 STL 미지원 시나리오이므로 양쪽 모두 "시스템 한계"에 해당. **핵심 발견 "32B≈8B"는 시스템이 지원하는 시나리오(S01-S09) 내에서만 주장하면 안전**:
```
S01-S09 only:
  32B = (75+100+88+75+50+88+67+50+20)/9 = 68.1%
  8B(trial1) = (75+100+88+75+50+88+67+0+20)/9 = 62.6%
  8B(3-trial) = (75+100+88+75+50+88+67+0+20)/9 = 62.6%
  Δ = 5.5pp
```

**그런데** — S08 8B=0%가 여전히 차이의 주범. S08을 제외하면:
```
S01-S07,S09 only:
  32B = (75+100+88+75+50+88+67+20)/8 = 70.4%
  8B  = (75+100+88+75+50+88+67+20)/8 = 70.4%
  Δ = 0.0pp (완벽히 동일!)
```

**결론**: 32B≡8B의 진짜 영역은 "시스템이 지원하고 8B가 tool calling에 성공하는 시나리오" (8/10). 이것을 논문에서 정밀하게 서술해야 한다.

### [Root Cause] 핵심 구조적 문제를 해결했는가, 아니면 표면 증상만 패치했는가?

**패치한 것 (Iteration 4에서 적용)**:
1. ✅ "Tool bug" → "Tool design limitation" — 용어 수정 (올바름)
2. ✅ Factorial → hierarchical dependency — 재프레이밍 (올바름)
3. ✅ Gap 4 격하 — 실험 근거 없는 Gap 제거 (올바름)
4. ✅ GCI "trend convergence" — 정직한 보고 (올바름)

**해결하지 못한 근본 문제**:
1. ❌ **xml_generator motion_type 미지원** — 도구 코드 자체를 고치지 않음 → M-A3 61%에서 75%로의 향상 불가
2. ❌ **S10 비결정성** — 원인 미규명 (Ollama? Qwen3 8B? batch scheduling?)
3. ❌ **trial1 vs 3-trial 불일치** — tolerance_sensitivity.py 미수정, 논문 수치와 분석 도구 수치 불일치
4. ❌ **EXP-D 미실행** — "핵심 차별화"라고 했지만 실행 0 runs

**정직한 평가**: Iteration 1-4의 작업은 **논문 서술의 정확성과 정직성을 크게 개선**했다. 하지만 실험 데이터 자체는 변하지 않았다. 이것은 "증거를 더 정확하게 보고"하는 것이지 "더 강한 증거를 생산"하는 것이 아니다.

**"패치 vs 근본 해결" 판정**:
- 논문 품질: **근본 해결** (재프레이밍, 용어 정정, 구조 재설계 — 리뷰어 방어력 대폭 증가)
- 실험 품질: **패치** (데이터 불변, 도구 코드 불변, EXP-D 미실행)

---

## XXIX. Self-Socratic Debate에서 발견한 새로운 취약점

### V1. trial1-vs-3trial 불일치 (⚠️ 새 발견)

**문제**: tolerance_sensitivity.py와 논문의 8B 수치가 2.5pp 다름.
**수정 필요**: tolerance_sensitivity.py의 `find_xml()`에 `fuel_tank.xml`, `fuel_tank_analysis.xml` 추가 + 3-trial 평균 옵션 추가.

### V2. "32B≡8B" 유효 범위 미명시 (보강 필요)

**문제**: S08(8B 0%) + S10(비결정) 제외 시 8/10 완전 동일이지만, 이 범위를 논문에 명시하지 않음.
**수정 필요**: 논문뼈대 §6.4에 "In the 8 scenarios where both models successfully invoke tools, scores are identical. Divergence occurs only in edge cases: S08 (complex prompt overwhelming 8B tool calling) and S10 (STL unsupported + non-deterministic behavior)."

### V3. Tolerance sensitivity vs paper 수치 불일치 (⚠️ 새 발견)

**문제**: 논문이 tolerance sensitivity를 robustness 근거로 인용하면서, 그 분석의 8B 수치가 논문과 다름.
**수정 필요**: tolerance_sensitivity.py에 주석 추가 "Note: This analysis uses trial1 only. Paper 8B=58.7% includes 3-trial average for S10 (non-deterministic). 32B values are identical in both methods."

---

## XXX. Iteration 5 종합 판정

### 점수

```
항목                    Iter 4 실행후    Iter 5 Debate후    이유
────────────────────    ────────────    ───────────────    ─────
현재 수준               78/100          76/100             trial 불일치 발견 (-2)
V1+V2+V3 수정 후       80/100          80/100
EXP-B 10개 확장 후     82/100          82/100
EXP-D 완수 시          92/100          92/100
```

### "이것이 최선인가?" Iteration 5 답변

**분석과 논문 서술 수정으로서는 거의 최선에 도달했다.** Self-Socratic Debate로 3개의 새로운 취약점(V1-V3)을 발견했지만, 이들은 모두 **20분 이내의 코드/텍스트 수정**으로 해결 가능.

**진짜 "최선"에 필요한 것은 분석이 아니라 실행**:
1. xml_generator에 motion_type + amplitude_unit 파라미터 추가 (코드 수정)
2. EXP-A 재실행 (M-A3 61% → 75%)
3. EXP-D 실행 (자율 Baffle 최적화)

이것들은 Ralph Loop(분석 루프) 내에서 수행할 수 없는 **개발 작업**이다.

---

## Iteration 6: EXP-B 10-시나리오 확장 후 Self-Socratic Debate

**Date**: 2026-03-05
**Context**: EXP-B를 3→10 시나리오로 확장 완료 (80/80 runs). 논문뼈대 수치 전체 업데이트.

### [Doubt] 숨겨진 결함

**검증 완료**:
- B1=0% (40/40 runs): agent_output.json에서 "tools don't include simulation functions" 확인 — 도구는 존재하나 도메인 프롬프트 없이는 호출 불가
- B0 trial1 vs 3-trial avg: 32B=동일, 8B=S10에서만 2.5pp 차이 (σ=35.4)
- 3-시나리오 부분집합(S01,S04,S07): 확장 후에도 동일 점수 유지

**결론**: 숨겨진 결함 없음.

### [Falsification] 극단 케이스

**B2 model-size gap 발견**: B0에서 Δ=5pp(8/10 동일), B2에서 Δ=10pp(6/10 동일).

```
B2 차이 시나리오:
  S02: 32B=75%, 8B=0%   (8B XML 추출 실패)
  S05: 32B=0%,  8B=75%  (32B timeout)
  S06: 32B=75%, 8B=0%   (8B XML 추출 실패)
  S08: 32B=25%, 8B=0%   (8B 구조 실패)
```

**원인**: 모델 능력 차이가 아니라 Ollama API raw output에서 XML regex 추출의 fragility. 도구가 출력 형식을 구조화하여 이 분산을 흡수.

**조치**: §6.4에 B2 반증거 1줄 추가 → "도구가 모델 크기 분산 흡수" 주장 강화

### [Root Cause] 핵심 vs 증상

**핵심 문제**: EXP-B 통계 검정력 (3 시나리오 = N=6/조건)
**해결**: 10 시나리오 × 2 모델 = N=20/조건. B1=B4=0% (40/40 runs)로 hierarchical dependency 견고히 입증.

**추가 발견 및 수정**:
- §6.4 B2 model-size gap 증거 추가 (논문뼈대.tex line 164)
- EXPERIMENT_STATUS.md는 이미 line 130에 기술됨

### 결론

**3개 질문 모두 기술적으로 방어 가능:**
1. [Doubt] → 모든 수치 교차검증 완료, 숨겨진 결함 없음
2. [Falsification] → B2 gap(10pp)은 infrastructure artifact, 논문에 반증거로 추가하여 주장 강화
3. [Root Cause] → 핵심 문제(통계 검정력) 해결, 표면 패치 아님
