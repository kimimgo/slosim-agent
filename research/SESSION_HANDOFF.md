# Session Handoff — 2026-02-19

## 브랜치: `research/paper`

## 완료된 작업

### 연구 기획 (전체 완료)
- [x] RECON 서베이 (2,175편) → ZERO gap analysis → SEND 논문 계획
- [x] 경쟁자 12개 시스템 분석 (3차 deep analysis, PDF 전문 검증)
- [x] Gap Cycle 3 — 슬로싱 중심 리프레이밍 (5 gaps)
- [x] Paper skeleton 슬로싱 중심 업데이트 (93% 재작성)
- [x] 실험 설계서 (EXP-1~5, 5개 실험)
- [x] 검증 데이터셋 인덱스 (Tier A/B/C/D)
- [x] SPHERIC Test 10 유체 수 오류 수정 (3종→2종, 6곳)

### Phase 1 인프라 준비 (완료, 커밋 `a1d0dab`)
- [x] 어블레이션 코드 (4 모드: full/no-domain/no-rules/generic)
- [x] 비교 스크립트 (`compare_spheric.py`)
- [x] 시나리오 JSON (20개, ground truth 포함)

### Phase 2a: EXP-1 SPHERIC (완료, 커밋 `3231de6`)
- [x] SPHERIC Test 10 Low/High/Oil 3개 시뮬레이션 GPU 실행
- [x] 실험 데이터 vs 시뮬레이션 비교 (r=0.84~0.93, NRMSE=18~33%)
- [x] Fig 1, 2 + Table 2 생성

### Phase 2b: EXP-3 Chen2018 Parametric (완료, 커밋 `e576e26`)
- [x] 6개 fill level (120-390mm) XML 생성 및 GPU 시뮬레이션 (총 67분)
- [x] MeasureTool -elevation으로 free surface 추출
- [x] Fig 3 (6-panel SWL 시계열), Fig 3b (진폭 vs fill level)
- [x] Table 4: 진폭 35-67mm, A/h 0.13~0.56 (shallow fill에서 비선형 강화)

### Phase 2c: EXP-5 Industrial PoC (완료, 커밋 `e576e26`)
- [x] Baffle 비교: No Baffle=158.9mm → With Baffle=12.8mm (**91.9% 감소**)
- [x] Seismic: 10m 탱크 극한 슬로싱 (2362mm 진폭)
- [x] Fig 5 (배플 효과), Fig 6 (지진 시나리오)
- [x] Table 6: 산업 PoC 요약

---

## 기술 노트

### SWL Gauge vs MeasureTool -elevation
- DualSPHysics 내장 SWL gauge가 `swlz=0` 반환하는 버그 발견
- **해결책**: MeasureTool `-elevation` + grid `pointsdef` 방식
- 개별 점(`-points`)이 아닌 grid 정의(`ptls[x=...,y=...,z=...]`)가 필수
- Probe가 sway 진폭 범위 내 벽에 너무 가까우면 건조됨 → 벽에서 3dp 이상 이격

### 시뮬레이션 출력 경로
- DualSPHysics는 모든 출력을 `simulations/data/`에 덮어씀
- 배치 실행 시 각 케이스 후 `simulations/$case/data/`로 복사 필요

---

## 다음 작업: Phase 3 — Agent 테스트

### 3a. EXP-2: NL→XML 생성 테스트 (20 scenarios)
```bash
# 각 시나리오를 FULL 프롬프트로 실행
SLOSIM_ABLATION=full ./slosim-agent -p "시나리오 텍스트..."
# XML validity, physical correctness, simulation success 측정
```

### 3b. EXP-4: 어블레이션 테스트 (10 scenarios × 4 conditions = 40 runs)
```bash
for mode in full no-domain no-rules generic; do
    SLOSIM_ABLATION=$mode ./slosim-agent -p "..."
done
```

### 3c. EXP-3: Agent-driven parametric (시간 측정)
- Agent가 자동으로 6개 fill level 파라메트릭 스터디 수행 시간 측정

### Phase 4: 분석 및 시각화
- 전체 실험 결과 통합 분석
- 논문용 최종 Figure/Table 생성

### Phase 5: 논문 작성
- paper_skeleton.md 기반 본문 작성

---

## 핵심 파일 위치

| 파일 | 용도 |
|------|------|
| `research/output/experiment_design.md` | 실험 설계서 |
| `research/output/paper_skeleton.md` | 논문 뼈대 |
| `research/scripts/compare_spheric.py` | SPHERIC 비교 |
| `research/scripts/analyze_chen2018.py` | Chen2018 분석 |
| `research/scripts/analyze_exp5.py` | EXP-5 분석 |
| `research/experiments/exp1_spheric/` | EXP-1 결과 |
| `research/experiments/exp3_parametric/` | EXP-3 결과 |
| `research/experiments/exp5_industrial/` | EXP-5 결과 |
| `research/experiments/exp2_nl2xml/scenarios.json` | 20 시나리오 |

## 커밋 히스토리

```
e576e26 EXP-3 Chen2018 파라메트릭 + EXP-5 Industrial PoC 완료
3231de6 EXP-1 SPHERIC Test 10 벤치마크 검증 완료
fc0ed7e 세션 핸드오프 문서 작성
a1d0dab Phase 1 인프라 준비
092ddec SPHERIC 유체 수 오류 수정
8a1a8fa 검증 데이터셋 인덱스
ff8ba5c 실험 설계서
```
