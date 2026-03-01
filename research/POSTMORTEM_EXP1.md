# EXP-1 POSTMORTEM — SPHERIC Test 10 검증 실패

**Date**: 2026-03-01
**Branch**: research-v1
**Status**: FAILED — 5/5 검증 메트릭 미달

---

## 실험 요약

SPHERIC Benchmark Test 10 (Souto-Iglesias & Botia-Vera, 2011)을 DualSPHysics v5.4 GPU로 재현.
3개 서브케이스 (Water Lateral, Oil Lateral, Water Roof) 검증 시도.

## 실패 모드 분석

### F1: 진폭 오차 +63.5%

| 항목 | 값 |
|------|-----|
| 원인 | dp=0.004 (23 입자층) — 수심 93mm 대비 너무 조잡 |
| 영향 | 피크 압력 과대 추정, 자유 표면 해상도 부족 |
| 교훈 | SPH 슬로싱은 최소 40+ 입자층 필요 (dp ≤ 0.002) |

### F2: 해상도 비수렴 (35-97% 오차)

| 항목 | 값 |
|------|-----|
| 원인 | 2개 해상도만 (dp=0.004, 0.002), 둘 다 조잡 |
| 영향 | 수렴 판단 불가, GCI 계산 불가 |
| 교훈 | 최소 3개 해상도 + Richardson extrapolation 필수 |

### F3: Oil 0/4 피크 검출

| 항목 | 값 |
|------|-----|
| 원인 | DBC 과감쇠 — Oil의 높은 점성(0.045 Pa·s)과 DBC 인공 점성이 중첩 |
| 영향 | Oil 파형 완전 억제, 피크 검출 0개 |
| 교훈 | 점성 유체는 mDBC noslip 필수 |

### F4: 시간 해상도 부족 (200Hz vs 20kHz)

| 항목 | 값 |
|------|-----|
| 원인 | computedt=0.005 (200Hz 게이지 출력) |
| 영향 | 임팩트 스파이크 반치폭 캡처 불가, 시계열 비교 의미 없음 |
| 교훈 | 게이지 출력 10kHz 이상 필요 (computedt=0.0001) |

### F5: 시계열 상관 없음 (r = -0.087)

| 항목 | 값 |
|------|-----|
| 원인 | F1-F4의 복합 효과 |
| 영향 | 시뮬레이션이 실험과 물리적으로 다른 현상 재현 |
| 교훈 | 개별 문제 수정 후 cross-correlation 재측정 필요 |

## 타이밍 불일치 (+0.5s)

교수님 조언: SPH 초기 입자 정착 시간으로 인한 위상 지연은 물리적으로 합리적.
→ v2에서도 time-shift 정렬 후 비교 방식 유지.

## v2 수정 방향

| 실패 모드 | v2 수정 |
|-----------|---------|
| F1 진폭 | dp=0.004/0.002/0.001 수렴 연구 |
| F2 비수렴 | 3개 해상도 + Richardson extrapolation + GCI |
| F3 Oil | mDBC noslip (`BoundaryMethod=2`, `-mdbc_noslip`) |
| F4 시간해상도 | computedt=0.0001 (10kHz) |
| F5 상관 | 위 수정 후 time-shift aligned cross-correlation |

## 결론

research-v1은 탐색적 실험으로 실패 모드를 식별하는 데 가치가 있었다.
v2에서는 이 5개 실패 모드를 체계적으로 수정하여 재검증한다.
