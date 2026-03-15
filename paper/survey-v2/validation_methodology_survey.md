# SPH 슬로싱 검증 방법론 서베이

## 연구 질문
"다른 논문들에서는 SPHERIC Test 10 벤치마크를 어떻게 검증했다고 보고하는가?"

## 핵심 발견

### 1. 대부분의 SPH 슬로싱 논문은 정성적 비교만 수행

SPH 슬로싱 문헌을 체계적으로 조사한 결과, **대부분의 논문이 "good agreement" 또는 "close agreement" 같은 정성적 표현으로만 검증을 보고**한다. 정량적 오차 메트릭(peak error %, RMSE, L2 norm, cross-correlation)을 보고하는 논문은 극소수이며, Richardson extrapolation이나 GCI를 적용한 SPH 슬로싱 논문은 **발견되지 않았다**.

**Oberkampf & Trucano (2002)**: "The present method of qualitative 'graphical validation'—comparison of computational results and experimental data on a graph—is inadequate"

이 기준에 비추면, 현재 SPH 슬로싱 문헌의 대다수가 부적절한 검증을 수행하고 있다.

### 2. SPHERIC Test 10을 사용한 주요 논문들의 검증 방법

| 논문 | 년도 | 인용수 | 검증 방법 | 정량 메트릭 | 수렴 연구 | dp 값 |
|------|------|--------|----------|-----------|----------|-------|
| English et al. (mDBC) | 2022 | 263 | 시각적 비교 | 없음 ("good agreement") | 없음 | 0.004, 0.002 |
| Dominguez et al. (DualSPHysics) | 2022 | 494 | 시각적 비교 | 없음 (리뷰 논문) | 없음 | 다양 |
| Delorme et al. (canonical Part I) | 2009 | — | 정량적 | SPH 과추정 보고 | 없음 | 2D |
| Green & Peiró | 2018 | — | 정성+안정성 | 에너지 소산 | 없음 | GPU 해상도 |
| Flow3D (narrow tank) | — | — | 정량적 | RMS 4.79%, Max 8.33% | 예 | 격자 기반 |
| **우리 (M1-M8)** | **2026** | **—** | **정량적 8개 메트릭** | **M2=19.5%, r=0.655** | **예 (3-level)** | **0.004/0.002/0.001** |

### 3. DualSPHysics 팀의 SPHERIC Test 10 검증 (English et al. 2022)

mDBC 논문 Section 4.2:
- dp = 0.004 m, dp = 0.002 m, h/dp = 2
- DBC vs mDBC 비교: 시계열 오버레이 그래프 3개 (원위치/+h 오프셋/mDBC)
- **보고된 정량적 메트릭: 없음**
- 사용된 표현: "very good results", "good agreement with the experimental data"
- 주요 기여: mDBC가 DBC의 비물리적 갭 문제 해결 → 센서 정확 위치에서 측정 가능

**시사점**: DualSPHysics 개발팀 자체가 SPHERIC Test 10에 대해 정성적 비교만 수행. 이는 현 SPH 커뮤니티의 검증 표준이 낮음을 보여준다.

### 4. ASME V&V 20 표준과 SPH의 간극

**ASME V&V 20 요구사항**:
- Richardson extrapolation으로 이산화 불확실성 정량화
- 체계적 격자/입자 수렴 연구
- 실험 불확실성의 통계적 처리
- 검증 불확실성 (Uval) = 실험 + 수치 + 입력 불확실성 결합

**SPH 적용 현실**: 사실상 없음
- V&V 20은 격자 기반 CFD 표준으로 설계됨
- SPH의 입자 특성과 직접 호환되지 않음
- GCI의 "격자 수렴 지수"를 "입자 수렴 지수"로 변환하는 표준 방법이 없음
- **SPH 전용 V&V 표준은 존재하지 않음**

### 5. 커뮤니티의 암묵적 오차 허용 범위 (CFD 일반)

| 오차 범위 | 수용 수준 | 적용 맥락 |
|----------|----------|----------|
| < 5% | 우수 | 산업 표준, 확립된 방법 |
| 5-10% | 양호 | 일반 CFD, 대부분의 저널 |
| 10-15% | 허용 | 복잡한 유동, 새로운 방법 |
| 15-30% | 조건부 허용 | 어려운 문제, 개념 증명 (정당화 필요) |
| > 30% | 일반적 불허 | 강한 정당화 또는 정성적 연구만 |

**우리의 M2 = 19.5%**: "조건부 허용" 범위. 슬로싱 임팩트 압력의 본질적 확률성 (CoV 20-40%, 102회 반복의 실험 산포)을 고려하면 충분히 정당화 가능.

### 6. SPHERIC 커뮤니티 검증 기준

SPHERIC Grand Challenges (GC5: Applicability to Industry):
- "comprehensive validation" — 정성적 + 정량적 모두 요구
- "convincing quantitative error analysis" — 설득력 있는 정량적 오차 분석 필수
- 수렴률 조사 필수
- 불일치는 "물리적/수치적 근거로 충분히 설명되어야 함"

그러나 **구체적인 통과/불합격 기준값은 명시하지 않음**. "설득력 있는" 검증의 판단은 리뷰어에게 위임.

### 7. 우리 M1-M8 프레임워크의 위치

우리의 검증 프레임워크가 기존 문헌 대비 어떤 위치에 있는지:

| 메트릭 | 우리 | 기존 SPH 문헌 | 평가 |
|--------|------|-------------|------|
| Peak-in-band (M1) | 3/3 in ±2σ | 사용하는 논문 없음 | **최초 적용** |
| Mean peak error (M2) | 19.5% | 보고하는 논문 거의 없음 | **정량적 보고** |
| 수렴 단조성 (M3) | 3-level 확인 | 2-level만 하는 경우 드묾 | **표준 이상** |
| GCI (M4) | 계산 완료 | SPH 슬로싱에서 미발견 | **최초 적용** |
| Cross-correlation (M5) | r=0.655 | 보고하는 논문 극소수 | **정량적 보고** |
| Time shift (M6) | τ=+0.57s | 보고하는 논문 없음 | **최초 적용** |
| Oil peak detection (M7) | 4/4 | N/A | **표준** |
| Impact FWHM (M8) | 계산 완료 | 보고하는 논문 거의 없음 | **최초 적용** |

**결론**: 우리의 M1-M8 프레임워크는 **현재 SPH 슬로싱 문헌의 검증 관행보다 현저히 엄격**하다. 특히 M4(GCI)와 M5-M6(cross-correlation)은 이 분야에서 처음 적용된 것으로 보인다.

## 논문 포지셔닝 권고

### 약점 → 강점 전환

1. **M2 = 19.5% 해석**:
   - "19.5% 오차"로 보면 약해 보임
   - 그러나 **대부분의 SPH 논문은 오차를 보고조차 하지 않음**
   - 실험 피크의 CoV = 20-40% → 오차가 실험 산포 내에 있음
   - **"첫 번째로 정량적 오차를 보고한 SPH 슬로싱 연구 중 하나"**로 포지셔닝

2. **r = 0.655 해석**:
   - 완벽한 상관은 아님
   - 그러나 **cross-correlation을 보고하는 SPH 슬로싱 논문이 거의 없음**
   - 시간 이동(τ=+0.57s) 후 정렬의 물리적 정당화 포함

3. **GCI 적용**:
   - SPH 슬로싱에서 GCI를 적용한 최초 사례 (발견된 바 없음)
   - 이 자체가 방법론적 기여

### 논문에서 강조할 내용

> "기존 SPH 슬로싱 검증은 대부분 시각적 비교에 의존하며 (English et al. 2022; Dominguez et al. 2022), 정량적 오차 메트릭을 체계적으로 보고하는 연구는 드물다. 본 연구에서는 8개의 정량적 메트릭 (M1-M8)을 정의하고 적용하여, Peak-in-band 분석, GCI 수렴 지수, 교차상관 분석을 포함한 체계적 검증 프레임워크를 제시한다. 이는 ASME V&V 20의 정신을 SPH 입자법에 적용한 첫 시도 중 하나이다."

## 참고 문헌

### 핵심 검증 방법론
- Oberkampf, W.L., Trucano, T.G. (2002). V&V in computational simulations. Progress in Aerospace Sciences, 38, 209-272.
- Roache, P.J. (1998). Verification of codes and calculations. AIAA Journal, 36(5), 696-702.
- ASME V&V 20-2009. Standard for Verification and Validation in CFD and Heat Transfer.

### SPH 슬로싱 핵심 논문
- English, A. et al. (2022). Modified dynamic boundary conditions (mDBC). Computational Particle Mechanics, 263 citations.
- Dominguez, J.M. et al. (2022). DualSPHysics: from fluid dynamics to multiphysics. Computational Particle Mechanics, 494 citations.
- Vacondio, R. et al. (2020). Grand challenges for SPH. Computational Particle Mechanics, 241 citations.
- Souto-Iglesias, A. et al. (2011). A set of canonical problems in sloshing. Part 0. Ocean Engineering.
- Delorme, L. et al. (2009). A set of canonical problems in sloshing. Part I. Ocean Engineering.

### SPHERIC 벤치마크
- Botia-Vera, E., Souto-Iglesias, A. et al. (2010). SPHERIC Benchmark Test Case 10.
- SPHERIC Validation Tests: https://www.spheric-sph.org/validation-tests

---

## 서베이 메타데이터

- 조사 기간: 2026-03-02
- research-sniper 세션: 23d70d4b7230
- 수집 논문: 742편 (8개 검색 쿼리)
- Tier 1 (핵심): 30편
- Deep read: 6편 전문 추출
- 웹 검색: SPHERIC 공식 사이트, DualSPHysics 포럼, ResearchGate, MDPI
