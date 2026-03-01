# Research-v2 — Session Handoff

## Status: ACTIVE
## Current Phase: EXP-1 Phase 1 (Convergence Study Design)
## Last Updated: 2026-03-01

### Quick Context

research-v1에서 SPHERIC Test 10 검증 실패 (진폭 +63.5%, Oil 0/4, r=-0.087).
v2는 5개 실패 모드를 수정한 6개 케이스로 수렴 연구를 수행한다.

### Experiment Registry

| EXP | Name | Status | Depends | Key Result |
|-----|------|--------|---------|-----------|
| 1 | SPHERIC Test 10 재검증 | DESIGNED | — | — |
| 2 | NL → XML Generation | PLANNED | EXP-1 | — |
| 3 | Parametric Study | PLANNED | EXP-1,2 | — |
| 4 | Ablation Study | PLANNED | EXP-1,2 | — |
| 5 | Industrial Application | PLANNED | EXP-1,2,3 | — |

### EXP-1 Run Log

| Run | Case | dp | BC | Status | M1 | M2 | M5 | Notes |
|-----|------|----|----|--------|----|----|----|----|
| 001 | Water Lat | 0.004 | DBC | PENDING | — | — | — | coarse baseline |
| 002 | Water Lat | 0.002 | DBC | PENDING | — | — | — | medium |
| 003 | Water Lat | 0.001 | DBC | PENDING | — | — | — | fine |
| 004 | Water Lat | 0.002 | mDBC | PENDING | — | — | — | BC comparison |
| 005 | Oil Lat | 0.002 | mDBC | PENDING | — | — | — | viscous fluid |
| 006 | Water Roof | 0.002 | mDBC | PENDING | — | — | — | roof impact |

### Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-03-01 | v1 CLOSED, v2 시작 | 5개 실패 모드 체계적 수정 필요 |
| 2026-03-01 | computedt=0.0001 | 게이지 10kHz, 실험 20kHz의 절반 — 임팩트 캡처 충분 |
| 2026-03-01 | mDBC noslip for Oil/Roof | DBC 과감쇠 해결, 점성 유체 필수 |
| 2026-03-01 | 3-level convergence | ASME V&V 20 GCI 요구사항 충족 |

### Known Issues

- [ ] dp=0.001 (Run 003) 예상 2시간 — GPU 메모리 확인 필요 (RTX 4090 24GB)
- [ ] mDBC 플래그 `-mdbc_noslip`이 DualSPHysics v5.4에서 지원 확인 필요
- [ ] 실험 데이터 case_2 (Water Roof) 존재 여부 확인 필요

### File Map

```
research-v2/
├── README.md                    ← 이 파일 (SSOT)
├── experiment_registry.json     ← 전체 실험 상태 추적
├── POSTMORTEM_V1.md             ← v1 실패 분석
├── exp1_spheric/
│   ├── README.md                ← EXP-1 설계서
│   ├── cases/                   ← 6개 XML 케이스
│   ├── runs/                    ← 각 run의 결과 메타데이터
│   └── analysis/                ← 분석 스크립트
├── exp2_nl2xml/README.md
├── exp3_parametric/README.md
├── exp4_ablation/README.md
├── exp5_industrial/README.md
└── scripts/common/              ← 공유 유틸리티
```
