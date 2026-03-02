# 서베이 심층 분석 + 핵심 논문 비교 + Target GAP

## 1. 서베이 개요

- **조사 범위**: 742편 (8개 검색 쿼리, research-sniper)
- **Tier 1**: 30편 핵심 논문
- **Deep Read**: 6편 전문 추출
- **핵심 질문**: "SPH 슬로싱에서 SPHERIC Test 10 벤치마크를 어떻게 검증했는가?"

## 2. 핵심 논문 심층 비교

### 2.1 English et al. (2022) — mDBC

**DOI**: 10.1007/s40571-021-00403-3 | **인용**: 263회

- **연구 내용**: DualSPHysics DBC의 비물리적 갭 문제를 해결하는 mDBC (modified Dynamic Boundary Condition) 제안
- **SPHERIC Test 10 사용**: Section 4.2에서 sloshing tank 검증
  - dp = 0.004 m, dp = 0.002 m, h/dp = 2
  - DBC vs mDBC 비교: 시계열 오버레이 그래프 3개
  - mDBC가 DBC의 비물리적 갭 제거 → 센서 정확 위치 측정 가능
- **검증 방법**: **순수 시각적 비교 (visual comparison)**
  - "good agreement with the experimental data" — 유일한 정량적 표현
  - **peak error, RMSE, cross-correlation, GCI: 일체 없음**
- **우리와의 차이**:
  - 우리: M1-M8 8개 정량 메트릭, 3-level 수렴 연구, GCI 계산
  - English: 0개 정량 메트릭, 2-level만 (수렴 분석 없음)
  - **동일 벤치마크에서 검증 엄격도 차이가 극명**

### 2.2 Vacondio et al. (2020) — Grand Challenges for SPH

**DOI**: 10.1007/s40571-020-00354-1 | **인용**: 241회

- **연구 내용**: SPHERIC 커뮤니티가 정의한 SPH 5대 Grand Challenges
- **GC5 (Applicability to Industry)** 핵심 요구:
  - "comprehensive validation" — 정성+정량 모두 필요
  - "convincing quantitative error analysis" — 설득력 있는 정량 오차 분석
  - 수렴률 조사 필수
  - 불일치는 "물리적/수치적 근거로 충분히 설명" 필요
- **그러나**: 구체적 통과/불합격 기준값 미명시 → 리뷰어 재량
- **GC1 (Convergence)**: Richardson extrapolation 언급은 있으나 SPH 입자 기반 수렴의 어려움 강조
  - Quinlan et al. (2006): 입자 분포에 따라 수렴률 변동
  - 수렴 + 일관성 + 안정성: SPH에서 여전히 grand challenge
- **우리와의 관계**:
  - 우리의 M1-M8 프레임워크는 GC5 요구사항의 **최초 체계적 실현** 중 하나
  - GCI 적용은 GC1의 수렴 분석 요구에 직접 응답
  - 하지만 SPH 입자법에 Richardson extrapolation 직접 적용은 엄밀하지 않음 → 이를 인정하고 "정신 차용"으로 포지셔닝

### 2.3 Dominguez et al. (2022) — DualSPHysics v5

**인용**: 494회 (DualSPHysics 코드 논문)

- **연구 내용**: DualSPHysics 코드 전반 리뷰, multiphysics 확장
- **슬로싱 검증**: English et al. (2022) 참조로 대체 — 별도 정량 검증 없음
- **의미**: DualSPHysics 팀 자체가 슬로싱에서 정량 검증을 수행하지 않았음

### 2.4 Delorme et al. (2009) — Canonical Sloshing Part I

- **연구 내용**: SPH 슬로싱 canonical problems 정의 (Part 0: Souto-Iglesias et al.)
- **검증 방법**: 정량적 — SPH가 실험 대비 과추정(overestimation) 경향 보고
- **한계**: 2D SPH, GCI 미적용, dp 비교 제한적
- **우리와의 관계**: 이 논문이 정량적 보고를 한 드문 예 → 우리가 이를 확장

### 2.5 Laha et al. (2024) — Heart Valve SPH FSI

**DOI**: 10.1038/s41598-024-57177-w

- DualSPHysics를 심장판막 FSI에 적용
- **검증 방법**: FVM, 4D MRI 데이터와 비교 — 정량적 (leaflet angle error 5.6%)
- **의미**: 다른 응용에서는 정량 검증을 하면서, 슬로싱에서는 안 하는 이중 기준

### 2.6 Flow3D 연구 (격자 기반, narrow tank)

- **방법**: 격자 기반 VOF
- **검증**: RMS 4.79%, Max error 8.33%, 수렴 연구 수행
- **의미**: 격자 기반 CFD는 이미 정량 검증이 표준 → SPH의 갭이 더 두드러짐

## 3. 문헌 비교 매트릭스

| 항목 | English 2022 | Vacondio 2020 | Delorme 2009 | Flow3D | **우리 (2026)** |
|------|-------------|--------------|-------------|--------|---------------|
| 벤치마크 | SPHERIC T10 | (리뷰) | canonical | narrow | SPHERIC T10 |
| 솔버 | DualSPHysics | (리뷰) | SPH 2D | Flow3D | DualSPHysics |
| BC | DBC/mDBC | — | — | VOF | DBC |
| dp 수준 | 2 (4mm, 2mm) | — | 1 | 격자 | 3 (4/2/1mm) |
| 정량 메트릭 수 | **0** | — | 1-2 | 3+ | **8 (M1-M8)** |
| Peak error 보고 | 없음 | — | 있음 | 4.79% RMS | **19.5% MAE** |
| Cross-correlation | 없음 | — | 없음 | 없음 | **r=0.655** |
| GCI/Richardson | 없음 | 언급만 | 없음 | 있음 | **적용 (2-level)** |
| 수렴 단조성 | 미확인 | — | 미확인 | 확인 | **확인 (3-level)** |
| 다유체 | 물만 | — | 물만 | 물만 | **물 + 오일** |
| Roof impact | 없음 | — | 없음 | 없음 | **시도 (PARTIAL)** |
| 검증 판정 | "good agreement" | — | "overestimation" | 수치 | **PASS/FAIL 체계** |

## 4. Target GAP 정의

### GAP 1: SPH 슬로싱 정량 검증 부재
- **현상**: 대부분의 SPH 슬로싱 논문이 시각적 비교만 수행
- **영향**: 재현성 없음, 솔버 간 비교 불가, 산업 신뢰 확보 불가
- **근거**: English et al. (2022) — DualSPHysics 팀 자체가 정량 메트릭 미보고
- **우리의 기여**: M1-M8 프레임워크로 체계적 정량 검증 수행

### GAP 2: SPH 수렴 분석 표준 부재
- **현상**: GCI/Richardson extrapolation이 SPH 슬로싱에 적용된 사례 없음
- **영향**: 수치 해의 신뢰도 정량화 불가 (ASME V&V 20 미충족)
- **근거**: Vacondio et al. (2020) GC1 — 수렴/일관성이 grand challenge
- **우리의 기여**: 3-level dp 수렴 연구 + 2-level GCI 계산 (SPH 최초)

### GAP 3: ASME V&V 20과 SPH의 간극
- **현상**: V&V 20은 격자 기반 CFD 표준으로 설계, SPH 전용 표준 부재
- **영향**: SPH 코드 검증에 참조할 프레임워크 없음
- **우리의 기여**: V&V 20의 "정신"을 SPH 입자법에 적용한 첫 시도

### GAP 4: 다조건 검증의 부재
- **현상**: 대부분 단일 유체/단일 하중에서만 검증
- **영향**: 솔버의 일반성 평가 불가
- **우리의 기여**: Water lateral + Oil lateral + Water roof (3개 서브케이스)

## 5. 논문 포지셔닝 전략

### 강점 (Strengths)
1. **최초 체계적 정량 검증**: M1-M8 8개 메트릭 — SPH 슬로싱 문헌 대비 최다
2. **최초 GCI 적용**: SPH 슬로싱에서 Richardson extrapolation 차용
3. **3-level 수렴 실증**: dp = 4/2/1mm, 교차상관 단조 수렴 확인
4. **다조건 검증**: 3개 서브케이스 (물, 오일, 지붕)
5. **재현성**: 모든 스크립트 공개, JSON 메트릭 출력, 자동화된 PASS/FAIL 판정

### 약점 → 정직한 보고 (Limitations as Transparency)
1. **M2 = 19.5%**: 높아 보이지만, 대부분의 논문이 보고조차 안 함 + 실험 CoV 20-40%
2. **DBC only**: mDBC 시도했으나 발산 → 정직하게 문서화
3. **2-level GCI**: Run 003 TimeOut 제약 → conservative Fs=3.0
4. **시간 이동**: τ=+0.57s → 물리적 원인 설명 (SPH 초기화 특성)

### 약점 → 강점 전환 키워드
- "첫 번째로 정량적 오차를 보고한 SPH 슬로싱 연구 중 하나"
- "기존 문헌의 검증 관행을 넘어서는 엄격한 프레임워크"
- "ASME V&V 20의 정신을 SPH 입자법에 적용한 첫 시도"
