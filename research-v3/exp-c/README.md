# EXP-1: SPHERIC Benchmark Test 10 재검증

**Status**: DESIGNED
**Phase**: Convergence Study + BC Comparison

## 목표

DualSPHysics v5.4 GPU 솔버로 SPHERIC Test 10 (Souto-Iglesias & Botia-Vera, 2011) 재현.
research-v1의 5개 실패 모드를 체계적으로 수정하여 정량적 검증 달성.

## 실험 참조

- **DOI**: 10.1016/j.oceaneng.2015.05.013
- **탱크**: 900 x 62 x 508 mm (L x B x H), quasi-2D
- **충전율**: Low 18.3% (h=93mm), High 70% (h=355.6mm)
- **유체**: Water (rho=998), Sunflower Oil (rho=990, mu=0.045)
- **운동**: Pitch rotation, A=4 deg

## v1 실패 원인 → v2 수정

| 실패 모드 | v1 값 | v2 수정 |
|-----------|-------|---------|
| 진폭 +63.5% | dp=0.004 (23층) | dp=0.004/0.002/0.001 수렴 연구 |
| 비수렴 35-97% | 2개 해상도 | 3개 해상도 + GCI |
| Oil 0/4 검출 | DBC | mDBC noslip |
| 200Hz 게이지 | computedt=0.005 | computedt=0.0001 (10kHz) |
| r=-0.087 | 복합 | 위 수정 후 재측정 |

## 수렴 연구 매트릭스

| Run | 케이스 | dp | BC | 입자 수 | 수심 입자층 | 예상시간 |
|-----|--------|----|----|---------|-----------|---------|
| 001 | Water Lat | 0.004 | DBC | ~136K | 23 | ~3분 |
| 002 | Water Lat | 0.002 | DBC | ~545K | 46 | ~15분 |
| 003 | Water Lat | 0.001 | DBC | ~2.2M | 93 | ~2시간 |
| 004 | Water Lat | 0.002 | mDBC | ~545K | 46 | ~15분 |
| 005 | Oil Lat | 0.002 | mDBC | ~545K | 46 | ~15분 |
| 006 | Water Roof | 0.002 | mDBC | ~545K | 175 | ~20분 |

## 검증 기준 (Pass/Fail)

| ID | 메트릭 | 통과 기준 |
|----|--------|----------|
| M1 | Peak-in-band | sim_peak ∈ [mu-2sigma, mu+2sigma], >= 3/4 peaks |
| M2 | Mean absolute peak error | < 30% |
| M3 | 수렴 단조성 | \|p_fine - p_med\| < \|p_med - p_coarse\| |
| M4 | GCI (ASME V&V 20) | GCI_fine < 10% |
| M5 | Cross-correlation (정렬 후) | r_max > 0.5 |
| M6 | Time shift | < 1.0 period |
| M7 | Oil 피크 검출 | >= 3/4 |
| M8 | Impact spike FWHM | 실험 대비 2x 이내 |

## 서브케이스별 판정

- **Water Lateral**: M1>=3/4 AND M2<30% AND M3 단조 AND M4<10% → PASS
- **Oil Lateral**: M7>=3/4 AND M1>=2/4 → PASS
- **Water Roof**: M1>=3/4 AND M2<40% → PASS
- **EXP-1 전체**: Water PASS + Oil PASS (또는 PARTIAL+개선 입증) → PASS

## 파일 구조

```
exp1_spheric/
├── README.md              ← 이 파일
├── cases/
│   ├── run_001_water_lat_dp004_dbc.xml
│   ├── run_002_water_lat_dp002_dbc.xml
│   ├── run_003_water_lat_dp001_dbc.xml
│   ├── run_004_water_lat_dp002_mdbc.xml
│   ├── run_005_oil_lat_dp002_mdbc.xml
│   └── run_006_water_roof_dp002_mdbc.xml
├── runs/
│   └── run_NNN/config.json + status.json + metrics.json
├── analysis/
│   └── (수렴, 피크 비교, 시계열 비교 스크립트)
└── validation_report.md   ← 최종 PASS/FAIL
```

## 실험 데이터 위치

- 프로브 정의: `simulations/spheric_probes.txt` (16 points)
- 실험 데이터: `datasets/spheric/case_1/` (Water Lat, Low fill, 4 repetitions)
- 실험 데이터: `datasets/spheric/case_3/` (Oil Lat, Low fill, 4 repetitions)
- 실험 데이터: `datasets/spheric/case_2/` (Water Roof, High fill, 4 repetitions)
