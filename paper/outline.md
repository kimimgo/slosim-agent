# 논문 뼈대 — Section별 가이드

## 메타 정보

- **제목**: A Quantitative Validation Framework for SPH Sloshing Impact Pressure: Systematic Metrics, Convergence Analysis, and SPHERIC Test 10 Benchmark
- **타겟 저널**: Computational Particle Mechanics (Springer) — English et al. (2022), Vacondio et al. (2020) 게재지
- **예상 분량**: 20-25 pages (figures 포함)
- **핵심 메시지**: SPH 슬로싱 검증은 정성적 비교에서 정량적 프레임워크로 전환이 필요하며, 본 연구가 이를 최초로 체계화

---

## Section 1: Introduction (~2.5 pages)

### 구조: Context → Gap → Contribution → Organization

**Context** (0.5p):
- 슬로싱의 공학적 중요성 (해양, 항공, 지진)
- SPH가 슬로싱에 적합한 이유 (meshless, Lagrangian, violent free-surface)
- DualSPHysics의 부상 (494 인용, open-source GPU)

**Gap 1** (0.5p): 정량 검증 부재
- 서베이 결과: 742편 조사, 대부분 "good agreement" 정성적 비교
- English et al. (2022): SPHERIC T10 사용, 정량 메트릭 0개
- Oberkampf & Trucano (2002) 인용: "graphical validation is inadequate"

**Gap 2** (0.5p): GCI/수렴 분석 미적용
- ASME V&V 20: 격자 수렴 연구 표준 → SPH 슬로싱에 적용 사례 없음
- Vacondio et al. (2020) GC1+GC5: 수렴 + 정량 검증 요구
- Richardson extrapolation → SPH 입자법 적응 필요

**Gap 3** (0.3p): 표준 프레임워크 부재
- 재현성, 코드 간 비교, 산업 채택 장벽

**Contribution** (0.4p):
1. M1-M8 프레임워크 제안
2. SPH 슬로싱 최초 GCI 적용
3. SPHERIC T10 3개 서브케이스 적용
4. 스크립트 공개 (재현성)

**Organization** (0.3p): 섹션 구조 안내

---

## Section 2: Background and Related Work (~2 pages)

### 2.1 SPHERIC Benchmark Test Case 10
- 탱크 형상: 900 × 62 × 508 mm
- 3개 서브케이스: Water lat (h=93mm, 0.613Hz), Oil lat (990 kg/m³, 0.651Hz), Water roof (h=355.6mm, 0.856Hz resonance)
- 실험 데이터: 102 반복, 피크 통계 (mean, std)
- 센서 위치: Sensor 1 (left wall, x=0.005m)
- **그림**: 탱크 스키매틱 + 센서 위치

### 2.2 SPH Validation Practices in Sloshing Literature
- **Table 1**: 문헌 비교 매트릭스 (5개 논문 vs 우리)
- English et al. 분석: mDBC 기여는 크지만 검증은 정성적
- Delorme et al.: 정량 보고 시도한 드문 예
- Flow3D: 격자 기반이므로 GCI 적용 → SPH 대비 기준점

### 2.3 ASME V&V 20 and SPH
- V&V 20 핵심: Richardson extrapolation, 체계적 격자 수렴, 실험 불확실성
- SPH 적용 난점: 입자 분포 비균일, 수렴 차수 불명확
- 우리의 적응: dp refinement → "particle convergence index", 실험 102회 통계

---

## Section 3: Validation Metrics Framework (M1-M8) (~2.5 pages)

### 핵심: 각 메트릭의 정의, 수식, Pass 기준, 적용 서브케이스

- **Table 2**: M1-M8 요약 테이블
- **M1**: Peak-in-band — 시뮬레이션 피크가 실험 μ±2σ 범위 내인지
- **M2**: Mean absolute peak error — 피크 진폭 오차 %
- **M3**: Convergence monotonicity — GCI 전제 조건
- **M4**: GCI — Roache (1998) 공식, Fs 선택 기준
- **M5-M6**: Cross-correlation + time shift — 정규화 교차상관, 물리적 범위 제한
- **M7**: Peak detection — scipy find_peaks 기반, 검출률
- **M8**: Impact FWHM — 임팩트 폭 비교

---

## Section 4: Numerical Setup (~1.5 pages)

### 4.1 DualSPHysics Configuration
- v5.4, GPU (RTX 4090), WCSPH, Wendland C2, h/dp=2
- δ-SPH (δ=0.1), α=0.01, Symplectic, CFL=0.2
- DBC (standard), TimeMax=7.0s

### 4.2 Case Configurations
- **Table 3**: Run matrix (5 runs)
- dp = 4, 2, 1 mm (Water lat convergence) + 2mm (Oil, Roof)

### 4.3 Pressure Measurement and Filtering
- MeasureTool, 10kHz, 1.5h rule probe placement
- SPH noise filtering: Median + Moving average (case-specific parameters)
- 과필터링 감지 기준

---

## Section 5: Results (~4 pages)

### 5.1 Water Lateral Impact (primary)
- **Figure 2**: 3-resolution time-series overlay + residual panel
- **Table 4**: Peak comparison (Run 002 vs experiment)
- M1=3/3, M2=19.5%, r=0.655, τ=+0.57s → PASS
- 피크 1 정확도: -0.1% (거의 정확), 피크 2-3 과소추정 경향 논의

### 5.2 Oil Lateral Impact
- **Figure 3**: Oil time-series vs experiment
- M7=4/4 (모든 피크 검출), M1=2/4 (피크 1-2 out-of-band)
- DBC 경계 소산 영향: 절대 피크 ~5mbar로 작아 상대 영향 큼
- r=0.570, τ=-1.047s → PASS

### 5.3 Water Roof Impact
- **Figure 4**: Near-roof probe 시계열
- 파도 높이 99.1% 도달 (503.6mm / 508mm)
- mDBC 시도 → 발산 (26% 입자 손실) 문서화
- PARTIAL 판정 — DBC 한계

---

## Section 6: Convergence Analysis (~2 pages)

### 6.1 Three-Level Resolution Study
- **Table 5**: Cross-correlation convergence (0.460 → 0.655 → 0.697)
- **Figure 5**: Peak pressure + error vs dp
- 피크 진폭: 비단조 (Peak 2) — SPH 해상도 역설 논의
- 교차상관: 단조 수렴 — 더 robust한 수렴 지표

### 6.2 GCI Calculation
- 2-level (Run 001→002), Fs=3.0
- 교차상관 기반 GCI 계산
- M3=FAIL for Peak 2 → 해당 피크 GCI 미적용

### 6.3 Convergence Discussion
- **SPH 해상도 역설**: 해상도↑ → 더 격렬한 임팩트 포착 → 피크 비단조
- **교차상관 우위**: 전체 시계열 형태 비교 → 피크 진폭보다 robust
- **권고**: SPH 슬로싱에서는 교차상관을 1차 수렴 지표로 사용

---

## Section 7: Discussion (~3 pages)

### 7.1 Positioning in SPH Sloshing Literature
- English et al.과 직접 비교: 동일 벤치마크, 검증 수준 차이
- 우리 19.5%의 맥락: (1) 대부분 보고 안 함 (2) 실험 CoV 20-40% (3) Flow3D 4.79% (격자 기반)
- **"첫 번째로 정량적 오차를 보고한 SPH 슬로싱 연구 중 하나"**

### 7.2 Relationship to ASME V&V 20
- 정신 차용 (spirit), 엄격 준수 아님
- dp refinement ≈ grid refinement — 합리적 유추
- 한계 인정: SPH 입자 분포 비균일, 수렴 차수 가정 필요
- **기여**: V&V 20 개념이 SPH에 적용 가능함을 실증

### 7.3 DBC Limitations and mDBC Prospects
- DBC ~5mm 소산층 → roof impact 불가, oil peaks 과소추정
- mDBC 시도 + 실패 → 정직한 보고 (GC2 boundary challenge)
- 향후: 안정적 mDBC 또는 대안 BC로 roof 검증 가능

### 7.4 Improvement from Previous Work (v1 → v2)
- **Table 6**: v1 vs v2 비교
- 핵심 개선: probe 배치 (1.5h rule), gauge 해상도 (50x), 필터링, 3-level 수렴
- **방법론적 개선** (솔버 변경 아님) → 다른 SPH 연구에도 적용 가능

### 7.5 Recommendations for SPH Sloshing Validation
1. M1, M2, M5 최소 보고
2. 교차상관을 1차 수렴 지표로 사용
3. 최소 2-level dp 수렴 연구
4. 실험 다중 반복의 불확실성 반영
5. 분석 스크립트 공개

---

## Section 8: Conclusions (~0.5 pages)

6개 핵심 결론:
1. M1-M8 프레임워크 제안 — 정성→정량 전환
2. SPHERIC T10 적용: Water PASS, Oil PASS, Roof PARTIAL
3. 3-level 수렴: 교차상관 단조 (0.460→0.655→0.697)
4. SPH 슬로싱 최초 GCI 적용
5. 스크립트 공개 (재현성)
6. 프레임워크는 솔버 비의존적

---

## 예상 그림 목록

| 번호 | 내용 | 소스 |
|------|------|------|
| Fig 1 | 탱크 스키매틱 + 센서 위치 | 새로 작성 |
| Fig 2 | Water lateral 3-resolution + residual | fig_timeseries.png |
| Fig 3 | Oil lateral time-series | fig_oil_lateral.png |
| Fig 4 | Water roof near-roof probes | fig_water_roof.png |
| Fig 5 | Peak convergence + error vs dp | fig_convergence.png |
| Fig 6 | 5-panel convergence analysis | convergence_study.png |

## 예상 테이블 목록

| 번호 | 내용 |
|------|------|
| Table 1 | 문헌 비교 매트릭스 |
| Table 2 | M1-M8 메트릭 정의 |
| Table 3 | Run matrix |
| Table 4 | Water lateral peak comparison |
| Table 5 | Cross-correlation convergence |
| Table 6 | v1 → v2 improvement |

---

## Overleaf 연동 계획

1. research-sniper의 `survey_export_plan` 또는 수동으로 Overleaf에 업로드
2. Overleaf 프로젝트: https://ko.overleaf.com/project/699b122eb24a9f1d4ecfaaeb
3. 파일 구조:
   ```
   main.tex          → 이 뼈대
   figures/           → exp1_spheric/figures/ 에서 복사
   references.bib     → 별도 생성 필요
   ```
4. 현재 Overleaf API 접근: research-sniper 업데이트 필요 (사용자 언급)
