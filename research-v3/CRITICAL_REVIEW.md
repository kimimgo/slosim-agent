# 실험 설계 비판적 검토 (Self-Review)

**Date**: 2026-03-03
**Purpose**: 논문 제출 전, 리뷰어 관점에서 실험 설계/결과/논리의 약점 식별

---

## 1. 심각도 높음 (Must Address)

### 1.1 Tier 사후 재분류 문제

**문제**: ground_truth.json과 논문의 tier 분류가 6/10 시나리오에서 불일치.

```
Scenario  GT(사전)    Paper(사후)   방향
S03       Medium      Easy         ↓ (88% 달성 → "쉬움"으로 재분류)
S07       Hard        Medium       ↓ (67% 달성 → "중간"으로 재분류)
S08       Easy        Hard         ↑ (50%/0% → "어려움"으로 재분류)
```

**리뷰어 공격 벡터**: "난이도 분류가 결과 기반 사후 조정이 아닌가? 이는 순환 논증이다."

**대응 방안**:
- A) ground_truth.json의 원래 tier를 사용하고, 결과 기반 "emergent difficulty"를 별도 논의
- B) tier 기준을 명시적으로 재정의 (예: geometry complexity + parameter count)하고, 사전 판단임을 문서화
- **권장: B)** — S08이 "Easy geometry but Hard dp+BC requirements"라는 논거로 재분류 정당화. ground_truth.json도 업데이트

### 1.2 도구가 오류를 도입하는 역설 (Tool-Induced Error)

**발견**: B0(Full, with tools)는 pitch를 `mvrectsinu`(오답)로 생성하지만, B2(no tools)는 `mvrotsinu`(정답)로 생성.

```
           motion_type    Correct?
B0 (tools)  mvrectsinu     ✗ (sway로 잘못 생성)
B2 (no tool) mvrotsinu     ✓ (rotation 정확)
```

**함의**: xml_generator 도구가 motion_type을 강제 오버라이드하여 모델의 올바른 판단을 무시함.

**리뷰어 공격 벡터**: "도구가 도입하는 오류가 도구가 방지하는 오류보다 크지 않은가?"

**정량 분석**:
- Tool이 방지하는 오류: geometry 정확도 (tank_z/fill_height 뒤바뀜) → ~+13.9pp
- Tool이 도입하는 오류: motion_type 강제 변환 → ~-12.5pp (5개 시나리오 × 2개 파라미터)
- **순 효과 ≈ +1.4pp** (거의 상쇄)

**대응**: Discussion에서 솔직히 논의. "도구는 geometry 정확성 보장과 motion_type 오류 도입의 트레이드오프. P1 버그 수정 시 M-A3 61% → 75%로 도약 가능."

### 1.3 도구 버그가 성능 천장을 결정

**발견**: P1(motion_type) + P2(amplitude) 버그만 수정하면 M-A3 **61.2% → 75.2% (+14pp)**.

```
현재:    61.2% — 도구 버그 포함
수정 후: 75.2% — 도구 버그 2개 수정
상한선: ~88%  — P3(timemax) + P4(geometry) 추가 수정
```

**리뷰어 공격 벡터**: "61%는 시스템 성능이 아니라 버그 있는 도구의 성능이다. 버그 수정 후 재평가하라."

**대응**:
- A) 버그를 수정하고 재실행 (가장 강력하지만 시간 소요)
- B) Discussion에서 "도구 품질이 성능 상한선 → 모델 스케일링보다 도구 개선이 우선"이라는 결론으로 연결
- **권장: B)** — 버그 포함 결과가 "모델 크기 무관성" 논지를 더 강화함. "61%에서도 32B≈8B이면, 75%에서도 32B≈8B일 것" 추론.

---

## 2. 심각도 중간 (Should Address)

### 2.1 EXP-B 통계적 검정력 부족

**문제**: 3개 시나리오 × 1 trial = factorial cell당 n=3.

- **반복(replication) 없음**: B1/B2/B4는 single trial. 비결정성 존재 시 신뢰도 저하
- **시나리오 선택 편향**: S01/S04/S07만 사용. 다른 조합(S02/S06/S09)이면 결과 다를 수 있음

**리뷰어 공격 벡터**: "n=3으로는 통계적 유의성 검정 불가. p-value를 보고하라."

**대응**:
- 2×2 factorial with 3 scenarios는 NL2FOAM 방법론과 유사 (설계 수준에서 정당화)
- "모든 조건에서 32B ≡ 8B 동일 결과"가 효과 크기의 강력한 증거
- 실험 비용(각 run ≈ 5분 LLM 호출)을 고려하면 추가 trials 가능하나, 결정론적 결과로 인해 불필요

### 2.2 사용자 연구 부재

**문제**: "비전문가용" 에이전트를 주장하면서 실제 비전문가로 테스트하지 않음.

**리뷰어 공격 벡터**: "비전문가가 실제로 NL 프롬프트를 올바르게 작성할 수 있는지 검증되지 않았다."

**대응**: Limitations에 명시. "본 연구는 시스템 역량(에이전트가 올바른 입력 시 얼마나 정확한가)을 평가함. 사용자 연구(비전문가가 올바른 입력을 생성할 수 있는가)는 future work."

### 2.3 인간 베이스라인 부재

**문제**: 61% M-A3가 실용적으로 어떤 의미인지 해석 불가.

- CFD 전문가가 같은 NL 설명을 읽고 XML 작성하면 몇 %?
- 초보자가 DualSPHysics 매뉴얼만 보고 작성하면 몇 %?
- 비전문가에게는 0% (접근 자체 불가)

**대응**: "비전문가 baseline = 0% (도구 없이는 DualSPHysics XML 작성 불가). 따라서 0% → 61%로의 도약 자체가 contribution."

### 2.4 EXP-A/EXP-C 연결 부재

**문제**: EXP-A는 XML 생성 품질(input), EXP-C는 시뮬레이션 결과 품질(output). 둘의 연결이 가정됨.

**리뷰어 공격 벡터**: "61% M-A3로 생성된 XML이 물리적으로 의미 있는 결과를 내는가?"

**대응**: EXP-C가 바로 이 질문에 답함. "M-A3 100% 시나리오(S01과 동일 탱크)에서 SPHERIC PASS." 하지만 61% 시나리오에서의 물리적 결과는 미검증.

---

## 3. 심각도 낮음 (Nice to Address)

### 3.1 Gap 4 (MCP) 약함

**문제**: MCP는 구현 프로토콜이지 연구 기여가 아님. 도구 통합이 핵심이면 REST API든 MCP든 동등.

**대응**: Gap 4를 "표준 프로토콜 기반 확장성"으로 프레이밍. 다른 솔버(OpenFOAM, COMSOL)로의 확장 가능성 강조.

### 3.2 단일 모델 패밀리

**문제**: Qwen3만 테스트. LLaMA 3.1, Gemma 2, Phi-3 등 미검증.

**대응**: "32B ≡ 8B 발견은 Qwen3 내에서 크기 무관성을 보임. 모델 패밀리 비교는 future work."

### 3.3 M-A3 Tolerance 관대함

**문제**: 10-15% tolerance는 공학적으로 관대. 탱크 크기 10% 오차는 실무에서 용납 불가.

**대응**: "tolerance는 NL→XML 번역의 '올바른 방향' 판단용. 정밀 튜닝은 도구의 역할(xml_generator가 정확한 값 입력 보장)."

---

## 4. 실험 설계 개선 가능성

### 4.1 즉시 실행 가능 (Low Effort, High Impact)

| 개선 | 예상 효과 | 필요 시간 |
|------|----------|----------|
| ~~ground_truth.json tier 재정의 + 기준 문서화~~ | ~~리뷰어 신뢰도 ↑~~ | ✅ 완료 |
| EXP-B를 10개 전체 시나리오로 확장 (B0 재사용) | 통계 검정력 ↑ | 2시간 |
| ~~Discussion에 tool-induced error 분석 추가~~ | ~~논문 깊이 ↑~~ | ✅ 완료 (analyze_tool_effect.py) |
| ~~P1+P2 수정 시 M-A3 75% 추정을 본문에 포함~~ | ~~ceiling analysis~~ | ✅ 완료 (§6.3) |

### 4.2 시간 투자 필요 (Future Work)

| 개선 | 예상 효과 | 필요 시간 |
|------|----------|----------|
| xml_generator P1+P2 버그 수정 + 재실행 | M-A3 61% → ~75% | 1-2일 |
| 다른 모델 패밀리 1종 테스트 | "model-agnostic" 주장 강화 | 1일 |
| 비전문가 3인 사용자 연구 | "비전문가용" 주장 강화 | 1주 |

---

## 5. 결론: "이것이 최선인가?"

### 강점 (What's already good)
1. **60 runs** — 충분한 데이터 볼륨
2. **2×2 factorial** — 깔끔한 실험 설계, 효과 크기 명확
3. **Temperature=0 determinism** — 재현성 확보
4. **SPHERIC 벤치마크** — 국제 표준 검증
5. **32B ≡ 8B** — 강력하고 실용적인 발견
6. **Error pattern 분석** — P1-P6 체계적 분류

### 약점 (What needs work)
1. **Tier 재분류** — 정당화 필요 (30분으로 해결 가능)
2. **Tool-induced error** — 핵심 발견이나 미충분 분석
3. **EXP-B 통계** — 3 시나리오는 약함 (확장 가능)
4. **사용자 연구 없음** — Limitations에 명시

### 판정: "출판 가능하나, 2-3가지 보강이 필요"

현재 결과로 **Workshop paper 또는 arXiv 수준은 충분**.
Top venue (CHI, EMNLP, OE) 제출에는:
1. Tier 기준 명시적 정의 + ground_truth 업데이트
2. Tool-induced error에 대한 심층 분석 (Section 6에 추가)
3. EXP-B 10 시나리오 확장 (가능하면)

이 3가지만 보강하면 "61%가 도구 한계이지 모델 한계가 아님"이라는 **더 강력한 스토리**가 완성됨.
