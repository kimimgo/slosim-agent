# Figure Index — SlosimAgent Research

총 98 파일, 10 카테고리

## expa/ (18) — EXP-A: M-A3 Parameter Fidelity
| File | Description |
|------|-------------|
| fig03_expa_heatmap | 논문 Fig.3: 10시나리오 × 2모델 히트맵 |
| fig04_v3v4_impact | 논문 Fig.4: v3→v4 tool fix 개선 |
| fig_expa_ma3 | EXP-A M-A3 bar chart |
| fig_v4_improvement | v4 improvement bar chart |
| 01_expa_heatmap | 히트맵 (viz 버전) |
| 02_expa_radar | 레이더 차트 |
| 03_expa_tier_grouped | 난이도별 그룹 |
| 09_parameter_breakdown | 파라미터별 분해 |
| 10_ceiling_analysis | 천장 분석 |

## expb/ (8) — EXP-B: 2×2 Factorial Ablation
| File | Description |
|------|-------------|
| fig05_expb_ablation | 논문 Fig.5: ablation 결과 |
| fig_expb_factorial | factorial 효과 bar |
| 04_expb_hierarchical | 계층적 의존 관계 |
| 05_expb_waterfall | 워터폴 차트 |

## expc_spheric/ (23) — EXP-C: SPHERIC Test 10
| File | Description |
|------|-------------|
| fig06_spheric_pressure | 논문 Fig.6: SPHERIC 압력 비교 |
| spheric_water_lateral | Water lateral 시계열 (실험 vs SPH) |
| spheric_oil_lateral | Oil lateral 시계열 (26.4% error) |
| spheric_peak_zoom | 피크 확대 |
| fig07_convergence_metrics | 논문 Fig.7: 수렴 메트릭 |
| 06_expc_viscosity_peaks | 점성 모델별 피크 비교 |
| 11_expc_pressure_detail | 압력 상세 |
| hires_comparison | 3-way: Artificial vs Laminar+SPS (0.1s/10ms) |
| viscosity_comparison | 점성 모델 비교 |
| visco_sweep_comparison | 점성 sweep |
| pressure_comparison_run001_002_exp | run001/002 vs 실험 |
| fig_oil_lateral | Oil lateral (초기) |
| fig_timeseries | 시계열 비교 |
| fig_convergence | 수렴 곡선 |
| convergence_study | 수렴 연구 |
| fig_water_roof | Water roof impact |

## expc_rafiee/ (4) — EXP-C: Rafiee 2011 Benchmark
| File | Description |
|------|-------------|
| rafiee2011_pressure_4way | dp004/dp002 × DBC 4-way 비교 |
| rafiee2011_pressure | Rafiee 시계열 |

## expc_bridge/ (4) — EXP-C: Bridge Multi-Physics
| File | Description |
|------|-------------|
| fig08_bridge_physics | 논문 Fig.8: 5기준 PASS |
| fig_bridge_pressure | Bridge 압력 시계열 |

## expd_baffle/ (6) — EXP-D: Baffle Optimization
| File | Description |
|------|-------------|
| fig09_expd_baffle | 논문 Fig.9: baffle vs baseline |
| expd_swl_comparison | SWL 28.5% 저감 비교 |
| 07_expd_baffle_summary | baffle 요약 |

## crossmodel/ (4) — Cross-Model Generalization
| File | Description |
|------|-------------|
| crossmodel_comparison | 5모델 × 10시나리오 bar chart |
| crossmodel_size_vs_score | 모델크기 vs 점수 scatter (8B>70B) |

## paraview_oil/ (9) — ParaView: Oil Peak Snapshots
| File | Description |
|------|-------------|
| expc_oil_Peak1_t1.54s | Peak 1 (t=1.54s) |
| expc_oil_Peak2_t3.07s | Peak 2 (t=3.07s) |
| expc_oil_Peak3_t4.61s | Peak 3 (t=4.61s) |
| expc_oil_Peak4_t6.14s | Peak 4 (t=6.14s) |
| expc_peaks_panel | 4-peak 패널 |

## paraview_baffle/ (16) — ParaView: Baffle vs Baseline
| File | Description |
|------|-------------|
| expd_baseline_t{3,4,5,6}.0s | Baseline 스냅샷 (4 timesteps) |
| expd_baffle_t{3,4,5,6}.0s | Baffle 스냅샷 (4 timesteps) |
| expd_comparison_panel | 비교 패널 |

## summary/ (6) — Summary & Dashboard
| File | Description |
|------|-------------|
| 08_gap_dashboard | GAP 증명 대시보드 |
| 12_summary_infographic | 전체 요약 인포그래픽 |
| experiment_results_summary | 실험 결과 종합 (LaTeX PDF) |
| slosim_agent_full_results | 전체 결과 (LaTeX PDF) |

---

## 논문별 사용 Figure

### Paper 1 (CS — Verification)
- `expa/fig03_expa_heatmap` — Table 3 시각화
- `expa/fig04_v3v4_impact` — Tool fix 영향
- `expb/fig05_expb_ablation` — Factorial ablation
- `crossmodel/crossmodel_comparison` — 5-model bar
- `crossmodel/crossmodel_size_vs_score` — 8B>70B scatter

### Paper 2 (PoF — Validation)
- `expc_spheric/spheric_water_lateral` — SPHERIC Water
- `expc_spheric/spheric_oil_lateral` — SPHERIC Oil 26.4%
- `expc_spheric/hires_comparison` — Artificial vs Laminar+SPS
- `expc_rafiee/rafiee2011_pressure_4way` — Rafiee 2011
- `expc_bridge/fig08_bridge_physics` — Bridge 5/5 PASS
- `expd_baffle/expd_swl_comparison` — Baffle 28.5%
- `paraview_oil/expc_peaks_panel` — Oil peak snapshots
- `paraview_baffle/expd_comparison_panel` — Baffle comparison
