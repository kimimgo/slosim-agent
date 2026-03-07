# 실험 결과 심층 비판적 분석

**Date**: 2026-03-03
**Purpose**: CRITICAL_REVIEW.md를 넘어서는 구조적/방법론적 약점 식별
**Status**: Ralph Loop Iteration 1

---

## A. 실험 설계의 구조적 문제 (CRITICAL_REVIEW 미포함)

### A1. EXP-A/B 수치 불일치

- EXP-A Overall M-A3: **61.2%** (10 시나리오)
- EXP-B B0 Mean: **72.2%** (S01/S04/S07만)
- 11pp 차이 — Cherry-picking 리스크

**대응**: 논문에 명시 "B0 values from 3 representative scenarios, not full mean" + EXP-B 10개 확장 검토.

### A2. per_model_overall 계산 불일치

- exp_a_full_analysis.json: trial1만 사용 → 8B=56.2%
- EXPERIMENT_STATUS.md: 3-trial 평균 → 8B=58.7%
- 논문에서 58.7% 사용 → JSON과 불일치

**수정 필요**: JSON의 per_model_overall을 3-trial 평균으로 갱신하거나, 차이를 주석 처리.

### A3. M-A3 tolerance 비대칭

```
tank: 10%, fill_height: 15%, frequency: 10%, amplitude: 15%, timemax: 10%
motion_type: exact, geometry_type: exact
```

근거 미문서화. Tolerance sensitivity analysis (5/10/15/20%) 필요.

### A4. motion_type ↔ amplitude 채점 커플링

score_expb.py에서 motion_type이 FAIL이면 amplitude도 자동 FAIL.
→ P1과 P2가 독립적이지 않음
→ M-A3를 과소 추정하는 방향으로 작용
→ "P1 수정 시 M-A3 75%" 추정은 이 커플링 고려 필요

---

## B. EXP-B 설계의 근본적 한계

### B1. n=3 통계 검정력

- Prompt 효과 +65.3pp: n=3에서도 유의미 (효과 크기 극대)
- Tools 효과 +6.9pp: n=3에서는 noise와 구분 불가능
- **해결**: 10개 시나리오로 확장 (B2/B4 7개 추가, ~2시간)

### B2. B1=0%의 해석

"도구가 쓸모없어서"가 아니라 "도구를 인식하지 못해서".
Availability ≠ Utility 구분이 논문에서 불충분.

### B3. B2/B4 구현 방식 차이

B0/B1: slosim-agent (TUI), B2/B4: Ollama raw API.
confounding variable — 도구 부재가 아니라 시스템 차이일 수 있음.
→ domain_prompt.txt ≡ SloshingCoderPrompt 동일성 검증 필요.

---

## C. EXP-C 분석 깊이 부족

### C1. 3 peaks만 비교 — 근거 불충분
### C2. GCI M3/M4 FAIL인데 "수렴 확인" 주장 — 모순
### C3. EXP-A→EXP-C 연결 실험 부재 (에이전트 생성 XML로 SPHERIC)

---

## D. 논문 논리 약점

### D1. Gap 4 (MCP) — 실험 근거 없는 Gap
### D2. "1,700편 서베이" — 자동 서베이 과대 표현 리스크

---

## E. 우선순위 보강 목록

| # | 항목 | 영향 | 노력 | 상태 |
|:---:|------|------|------|:---:|
| 1 | Tolerance sensitivity analysis | M-A3 robustness | 1h | ⬜ |
| 2 | EXP-B 10개 시나리오 확장 | 통계 검정력 | 2-3h | ⬜ |
| 3 | EXP-A→EXP-C 연결 실험 | 핵심 gap | 2h | ⬜ |
| 4 | GCI FAIL 솔직 기술 | 신뢰도 | 30m | ⬜ |
| 5 | B0 72.2% vs 61.2% 설명 | 리뷰어 방어 | 30m | ⬜ |
| 6 | Amplitude 채점 커플링 문서화 | 투명성 | 30m | ⬜ |
| 7 | Gap 4 약화/제거 | 견고성 | 1h | ⬜ |

---

## F. Tolerance Sensitivity 분석 결과 (즉시 실행)

### 핵심 발견: M-A3는 tolerance에 거의 불변

```
Tolerance     32B      8B       Δ
   5%      61.2%   56.2%   +5.0%
  10%      61.2%   56.2%   +5.0%  ◀ paper
  15%      61.2%   56.2%   +5.0%
  20%      61.2%   56.2%   +5.0%
  25%      62.4%   56.2%   +6.2%
```

**해석**: 실패하는 파라미터들은 tolerance-independent한 범주적 오류
(motion_type: exact match, amplitude: 수배 차이 단위 변환).
통과하는 파라미터들은 정확한 값 생성 (tolerance 내 충분한 여유).

→ **M-A3 메트릭이 tolerance-robust** = 리뷰어 방어 매우 강력.
→ 논문에 Appendix 또는 Table note로 추가 권장.

### P1/P2 수정 시 M-A3 정밀 추정 (기존 CRITICAL_REVIEW 수정)

```
Scenario | Actual | Fix P1 | Fix P1+P2
     S01 |  75.0% |  87.5% |   100.0%
     S04 |  75.0% |  87.5% |   100.0%
     S05 |  50.0% |  75.0% |    75.0%
     S09 |  20.0% |  40.0% |    60.0%
 Overall |  61.2% |  69.4% |    75.2%
```

**기존 주장 수정**:
- CRITICAL_REVIEW: "P1 수정 시 61% → 75%" — ✗ 과대
- 정밀 계산: P1만 수정 시 → **69.4%** (+8.2pp)
- P1+P2 모두 수정 시 → **75.2%** (+14.0pp) — 이것은 정확

이유: motion_type 수정 시 amplitude 비교가 활성화되나,
amplitude 값 자체가 radian이라 대부분 여전히 FAIL (S01: 0.045 ≈ 2.58° vs GT 4°).

---

## G. 사용자 피드백 반영: 실험 계획 수정

### 제외 항목
- ~~다른 모델 패밀리 (LLaMA/Gemma)~~ → 제외
- ~~비전문가 사용자 연구~~ → 제외 (오프라인)

### 추가 실험: EXP-D (Autonomous Baffle Optimization)

**핵심 시나리오**: STL CAD → 자연어 슬로싱 지시 → 에이전트 자율 Baffle 최적화

**의의**: EXP-A/B가 "파라미터 정확도"를 측정했다면, EXP-D는 **"자율 의사결정 능력"**을 측정.
- 에이전트가 결과를 해석하고 스스로 개선 방향을 결정하는 **closed-loop autonomy**
- 기존 7종 에이전트와 차별화되는 핵심 기여 (automation → autonomy)

**현재 상태**:
- optimization_prompt.md 작성 완료
- run_expd.sh, score_expd.py 구현 완료
- fuel_tank.stl 존재 (cases/)
- results/ 디렉토리 빈 상태 (.gitkeep)
- baffle_generator.go 구현 진행 중 (unstaged)
- stl_import.go fillpoint 보완 진행 중 (unstaged)

**선행 조건 (blocked by)**:
1. stl_import fillpoint 기능 완성
2. baffle_generator 도구 완성 + 테스트
3. slosim-agent 빌드
4. GPU Docker 환경 확인

---

## H. 수정된 우선순위 보강 목록

| # | 항목 | 영향 | 노력 | 상태 |
|:---:|------|------|------|:---:|
| 1 | ~~Tolerance sensitivity analysis~~ | M-A3 robustness | 1h | ✅ 완료 |
| 2 | **EXP-D 실행** (STL+Baffle 자율 최적화) | 핵심 새 기여 | 1-2일 | 🔄 실행중 |
| 3 | ~~CRITICAL_REVIEW P1 추정치 수정~~ | 정확성 | 15m | ✅ 완료 |
| 4 | EXP-B 10개 시나리오 확장 | 통계 검정력 | 2-3h | ⬜ |
| 5 | GCI FAIL 솔직 기술 | 신뢰도 | 30m | ⬜ |
| 6 | B0 72.2% vs 61.2% 설명 | 리뷰어 방어 | 30m | ⬜ |
| 7 | Gap 4 약화 | 견고성 | 1h | ⬜ |

---

## I. 결론

현재 수준: **Workshop/arXiv 출판 가능**.
사용자의 새 방향: EXP-D (자율 Baffle 최적화)가 논문의 핵심 차별화 기여가 될 수 있음.
"이것이 최선인가?" → **아직 아니다.** EXP-D 완수 + tolerance robustness 포함 시 "정규 저널 수준" 도달 가능.
