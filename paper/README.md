# Paper — SlosimAgent 논문 자료

## 브랜치 전략

```
main          ← slosim-agent 구현체 (Go 소스, TUI, 도구)
research-v1   ← EXP 수행 이력 v1 (SPHERIC T10 첫 시도, FAILED)
research-v2   ← EXP 수행 이력 v2 (재검증, run_001-008)
paper         ← 논문 작성 (현재 브랜치) — 항상 최신 뼈대 유지
```

## 디렉토리 구조

```
paper/
├── README.md                  ← 이 파일
├── outline.md                 ← 논문 뼈대 (Section별 가이드)
├── deep_analysis.md           ← Gap 4개 + 핵심 논문 비교 매트릭스
├── validation_report.md       ← EXP-1 검증 결과 (PASS 판정)
├── POSTMORTEM_V1.md           ← v1→v2 개선 근거
│
├── survey-v1/                 ← research-sniper v1 (AI Agent 관점)
│   ├── gap_report.md          ← 5개 Gap (LLM+SPH=0편, 로컬SLM, 슬로싱에이전트, MCP, 메트릭)
│   ├── deep_research_survey.md ← 핵심 논문 31편 (2,144편 수집)
│   ├── survey_result.json     ← 전체 서베이 데이터
│   └── profile.json           ← 서베이 프로파일
│
├── survey-v2/                 ← research-sniper v2 (검증 프레임워크 관점)
│   ├── validation_methodology_survey.md ← 742편 서베이, 정량 메트릭 부재 확인
│   ├── survey_session.json    ← 세션 데이터
│   └── fulltexts/             ← Deep read 6편 전문
│
├── analysis/                  ← EXP-1 분석 스크립트 + 메트릭
│   ├── convergence_analysis.py
│   ├── oil_roof_analysis.py
│   ├── paper_figures.py
│   ├── plot_pressure_comparison.py
│   ├── metrics.json           ← run_001/002 정량 메트릭
│   └── oil_roof_metrics.json
│
├── figures/                   ← 논문 그래프 (PNG)
│   ├── fig_timeseries.png     ← Water lateral 3-resolution
│   ├── fig_convergence.png    ← Peak pressure convergence
│   ├── fig_oil_lateral.png    ← Oil lateral vs experiment
│   ├── fig_water_roof.png     ← Water near-roof
│   └── convergence_study.png  ← 5-panel convergence
│
├── scenarios/                 ← E2E 유저 시나리오 (에이전트 범용성 증명)
│   ├── e2e_dataset_scenarios.md ← 10개 시나리오 정의 + GPU 결과
│   ├── e2e_scenario_kim.md
│   ├── e2e_test_analysis.md
│   ├── e2e_gpu_report.md
│   └── E2E_QUALITY_REPORT_v0.3.md
│
├── data/                      ← 실험 데이터 + 케이스 설정
│   ├── spheric/               ← SPHERIC Test 10 실험 데이터 (102회 반복)
│   └── cases/                 ← DualSPHysics XML 케이스 파일
│
└── drafts/                    ← 논문 초고 (버전별)
    ├── v1_paper.typ           ← Typst 초고 (v1)
    ├── v1_paper.pdf
    ├── v1_references.bib
    ├── v1_sloshagent_paper.pdf ← LaTeX 초고 (v1)
    └── v1_sloshagent_paper.bbl

## 논문 트랙 2개

### Track A: AI Agent for SPH Sloshing (survey-v1 기반)
- Gap: LLM+DualSPHysics 0편, 슬로싱 전문 에이전트 0편
- 증거: E2E 10개 시나리오 8/10 PASS
- 비교: OpenFOAMGPT, Foam-Agent, ChatCFD, NL2FOAM

### Track B: Quantitative Validation Framework (survey-v2 기반)
- Gap: SPH 슬로싱 정량 검증 부재, GCI 미적용
- 증거: M1-M8 프레임워크, run_001/002/003 수렴 분석
- 비교: English et al. (0개 메트릭) vs 우리 (8개 메트릭)
```

## 시뮬레이션 데이터 위치

PART/VTK 등 대용량 데이터는 이 브랜치에 포함하지 않음.
- 로컬: `/mnt/simdata/dualsphysics/exp1/`
- 분석에 필요한 CSV/메트릭만 `analysis/`에 포함
