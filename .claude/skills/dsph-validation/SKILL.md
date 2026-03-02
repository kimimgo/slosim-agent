# DualSPHysics Validation Pipeline (Orchestrator)

description: DualSPHysics 시뮬레이션 벤치마크 검증 오케스트레이터. 케이스 설계(dsph-case-design) → 실행(dsph-run) → 후처리(dsph-analysis) → 정량 검증(dsph-metrics) → 보고서 생성까지 전체 파이프라인 조율. "검증", "벤치마크", "validation", "SPHERIC" 시 자동 발동.

---

## 1. 전체 워크플로우

```
Phase 1: 케이스 설계     → dsph-case-design 스킬
Phase 2: 시뮬레이션 실행  → dsph-run 스킬
Phase 3: 후처리 + 분석   → dsph-analysis 스킬
Phase 4: 정량 검증       → dsph-metrics 스킬
Phase 5: 보고서 + 그래프  → 이 스킬 (아래 §3-4)
```

### 스킬 간 관계도

```
dsph-validation (이 스킬 — 오케스트레이터)
├── Phase 1 → dsph-case-design (dp, BC, 프로브, domain)
├── Phase 2 → dsph-run (GenCase → Solver → 모니터링)
├── Phase 3 → dsph-analysis (MeasureTool, 필터링, 피크, 상관)
├── Phase 4 → dsph-metrics (M1-M8, GCI, 판정, 문헌)
└── Phase 5 → 보고서/그래프 (이 스킬)
```

## 2. Phase별 상세 가이드

### Phase 1: 케이스 설계

**스킬**: `dsph-case-design`

1. 실험 데이터 확보 (시계열 + 피크 통계)
2. 목적에 맞는 dp 선택 (정성/정량/수렴)
3. BC 선택 (DBC/mDBC + 의사결정 매트릭스)
4. 수렴 연구 설계 (3-level: coarse/medium/fine)
5. 게이지 출력 해상도 설정 (computedt, TimeOut 분리)
6. 프로브 파일 생성 (1.5h rule)
7. Simulation domain 설정

**산출물**: XML 케이스 3개 + 프로브 파일

### Phase 2: 시뮬레이션 실행

**스킬**: `dsph-run`

1. GenCase 실행 → 입자 수, normals 확인
2. Solver 실행 → Run.csv 모니터링
3. 각 run 완료 확인

**확인사항**:
- GenCase: `Normals` 줄에서 zero normals 수 (mDBC 사용 시 0이어야 함)
- Solver: Run.csv에서 시간 진행률, Run.out에서 에러

### Phase 3: 후처리 + 분석

**스킬**: `dsph-analysis`

1. MeasureTool 실행 (일관된 프로브 파일)
2. CSV 로딩 → Pa를 mbar로 변환
3. 필터링 (케이스별 median + moving avg)
4. 과필터링 감지
5. Cross-Correlation → 시간 정렬 (M5, M6)
6. 피크 검출 → 실험 비교 (M1, M2)
7. (Roof) 다중 프로브 선택

**산출물**: 필터링된 시계열, 피크 값, 상관 계수

### Phase 4: 정량 검증

**스킬**: `dsph-metrics`

1. M1-M8 메트릭 계산
2. GCI 계산 (3-level 또는 2-level)
3. 서브케이스별 PASS/FAIL 판정
4. 전체 판정 (Water + Oil + Roof)
5. 문헌 대비 포지셔닝

**산출물**: metrics.json, oil_roof_metrics.json, 판정 결과

### Phase 5: 보고서 + 그래프

이 스킬이 직접 수행. 아래 §3-4 참조.

## 3. 논문 그래프 생성

### 3.1 표준 4개 그래프

| 그래프 | 파일 | 내용 |
|--------|------|------|
| 시계열 | `fig_timeseries.png` | 3-해상도 + residual (2패널) |
| 수렴 | `fig_convergence.png` | 피크 바차트 + 오차 vs 해상도 (2패널) |
| Oil | `fig_oil_lateral.png` | Oil 시계열 + 피크 어노테이션 |
| Roof | `fig_water_roof.png` | Roof 시계열 (한계 문서화) |

### 3.2 스타일 설정

```python
plt.rcParams.update({
    'font.size': 11, 'axes.labelsize': 12,
    'legend.fontsize': 9, 'savefig.dpi': 300,
    'axes.grid': True, 'grid.alpha': 0.3,
})
```

### 3.3 생성 방법

```bash
# 수렴 분석 (Water Lateral) → convergence_study.png + metrics.json
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py

# Oil + Roof 분석 → fig_oil_lateral.png, fig_water_roof.png + oil_roof_metrics.json
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py

# 논문용 그래프 (4개 표준 그래프)
python3 research-v2/exp1_spheric/analysis/paper_figures.py

# 최종 판정
python3 research-v2/exp1_spheric/scripts/final_verdict.py
```

## 4. validation_report.md 템플릿

```markdown
# Validation Report — EXP-1 SPHERIC Test 10

## 실험 개요
- 벤치마크: SPHERIC Test 10 (Botia-Vera et al. 2010)
- 탱크: 0.9m × 0.062m × (fill: 0.093m)
- 유체: Water (ρ=1000) / Oil (ρ=917)
- 운동: Lateral harmonic (f=0.613Hz, A=50mm)

## 시뮬레이션 설정
| Run | dp [mm] | BC | 입자 수 | GPU 시간 |
|-----|---------|-----|---------|---------|
| 001 | 4 | DBC | ~80K | ~5min |
| 002 | 2 | DBC | ~650K | ~30min |
| 003 | 1 | DBC | ~5M | ~4h |

## Water Lateral 결과
- M1 (Peak-in-band): X/3 PASS/FAIL
- M2 (Mean error): XX.X% PASS/FAIL
- M5 (Cross-corr): r=X.XXX PASS/FAIL
- M6 (Time shift): τ=±X.XXXs PASS/FAIL
- M3 (Monotone): PASS/FAIL
- M4 (GCI): X.X% PASS/FAIL

## Oil Lateral 결과
- M7 (Peak detection): X/4 PASS/FAIL
- M1 (Peak-in-band): X/4 PASS/FAIL

## Water Roof 결과
- M1 (Peak-in-band): X/4 PASS/FAIL
- M2 (Mean error): XX.X% PASS/FAIL

## 전체 판정
**OVERALL: PASS / PARTIAL / FAIL**

## 첨부 그래프
- fig_timeseries.png
- fig_convergence.png
- fig_oil_lateral.png
- fig_water_roof.png
```

## 5. 흔한 함정과 해결책 (전체 요약)

| 함정 | 증상 | 해결 | 상세 스킬 |
|------|------|------|----------|
| computedt 너무 큼 | 임팩트 피크 놓침 | ≤ 0.0001 (10kHz) | dsph-case-design §3.1 |
| 프로브 벽 너무 가까움 | 0 또는 노이즈만 | 1.5h rule 적용 | dsph-case-design §3.3 |
| 프로브 dp 스냅 | run 간 다른 위치 비교 | 동일 프로브 파일 + 헤더 확인 | dsph-analysis §1.4 |
| Oil 과필터링 | M1 FAIL (±0.016 mbar) | Med3+Avg7 (Water보다 약하게) | dsph-analysis §2.2 |
| mDBC zero normals | 솔버 즉시 crash | point dp/2 오프셋 | dsph-case-design §2.3 |
| mDBC 발산 | 입자 대량 손실 | 공명 운동에서는 DBC 사용 | dsph-case-design §2.2 |
| TimeOut 작음 + fine dp | 디스크 1TB+ | TimeOut↑, computedt 별도 | dsph-case-design §3.2 |
| numpy bool JSON | TypeError | default handler 추가 | dsph-analysis §6 |
| dp=0.001 피크 미캡처 | 100Hz로 제한 | 2-level GCI (Fs=3.0) | dsph-metrics §3.4 |

## 6. 전체 워크플로우 체크리스트

```
□  1. 실험 데이터 확보 (시계열 + 피크 통계)
□  2. XML 케이스 3개 (coarse/medium/fine)        → dsph-case-design
□  3. GenCase → 입자 수, normals 확인             → dsph-run
□  4. Solver 실행 → Run.csv 모니터링              → dsph-run + sim-monitor
□  5. MeasureTool → 일관된 프로브 파일             → dsph-analysis §1
□  6. CSV 로딩 → Pa를 mbar로 변환                 → dsph-analysis §1.3
□  7. 필터링 (case-specific median + moving avg)   → dsph-analysis §2
□  8. 피크 검출 + 실험 비교 (M1, M2)              → dsph-analysis §4
□  9. Cross-correlation (M5, M6)                   → dsph-analysis §3
□ 10. 수렴 분석 (M3, M4/GCI)                      → dsph-metrics §3
□ 11. 논문 그래프 생성                              → §3 (이 스킬)
□ 12. validation_report.md 작성                    → §4 (이 스킬)
```
