# Session Handoff — 2026-02-18

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
- [x] **1a 어블레이션 코드**: `internal/llm/prompt/sloshing_coder.go`
  - `SLOSIM_ABLATION` 환경변수로 4 모드 전환: `full` / `no-domain` / `no-rules` / `generic`
  - 11개 섹션 조합 방식, 테스트 11개 전부 PASS
  - 사용법: `SLOSIM_ABLATION=no-domain ./slosim-agent -p "..."`
- [x] **1b 비교 스크립트**: `research/scripts/compare_spheric.py`
  - SPHERIC 실험 vs 시뮬레이션 비교 (r, NRMSE, Peak Error, Phase Error)
  - CLI: `python research/scripts/compare_spheric.py --sim-dir ... --exp-file ... --out-dir ...`
- [x] **1c 시나리오 JSON**: `research/experiments/exp2_nl2xml/scenarios.json`
  - 20개 시나리오 (5 Level × 4), ground truth 포함

---

## 다음 작업: Phase 2 — GPU 시뮬레이션 실행

### 사전 조건
1. **Ollama 중지**: `ollama stop` (GPU VRAM 해제, Qwen3 32B ~20GB)
2. **Docker 확인**: `docker compose ps` (dsph-agent:latest 이미지 준비)
3. **NVIDIA 확인**: `nvidia-smi` (RTX 4090 VRAM 여유 확인)

### 2a. EXP-1: SPHERIC 3 cases

```bash
# Low Resolution (dp=0.006, ~136K particles)
docker compose run --rm dsph gencase /cases/SPHERIC_Test10_Low_Def /data/spheric_low
docker compose run --rm dsph dualsphysics /data/spheric_low/SPHERIC_Test10_Low_Def -gpu
docker compose run --rm dsph partvtk -dirin /data/spheric_low -savevtk /data/spheric_low/vtk/PartFluid -onlytype:-all,+fluid
docker compose run --rm dsph measuretool -dirin /data/spheric_low -pointsdef /data/spheric_probes.txt -savecsv /data/spheric_low/measure/pressure

# High Resolution (dp=0.004, ~344K particles)
docker compose run --rm dsph gencase /cases/SPHERIC_Test10_High_Def /data/spheric_high
docker compose run --rm dsph dualsphysics /data/spheric_high/SPHERIC_Test10_High_Def -gpu
# ... 동일 후처리

# Oil (Low, sunflower oil)
docker compose run --rm dsph gencase /cases/SPHERIC_Test10_Oil_Low_Def /data/spheric_oil_low
docker compose run --rm dsph dualsphysics /data/spheric_oil_low/SPHERIC_Test10_Oil_Low_Def -gpu
# ... 동일 후처리
```

### 2b. EXP-5: Industrial PoC (2 cases)
- Baffle: 탱크+배플 XML 작성 필요
- Seismic: seismic_input tool 테스트 필요

### 2c. EXP-3: Chen2018 Parametric (6 cases)
- 6개 fill level XML 생성 (120, 150, 195, 260, 325, 390 mm)
- Chen2018 케이스 XML이 이미 일부 존재: `cases/Chen2018_*.xml`

### Phase 2 후 → Phase 3 (Agent 테스트)
- 3a. EXP-2: 20 scenarios × FULL prompt
- 3b. EXP-4: 10 scenarios × 4 ablation conditions (40 runs)
- 3c. EXP-3: Agent-driven parametric (시간 측정)

### Phase 4: 분석 및 시각화
```bash
# SPHERIC 비교 실행
python research/scripts/compare_spheric.py \
  --sim-dir simulations/spheric_low/measure \
  --exp-file datasets/spheric/case_1/lateral_water_1x.txt \
  --peak-file datasets/spheric/case_1/Water_4first_peak_lateral_impact_tto_0_85_H93_B1X.txt \
  --out-dir research/experiments/exp1_spheric \
  --label "Low Resolution (dp=0.006)"
```

---

## 핵심 파일 위치

| 파일 | 용도 |
|------|------|
| `research/output/experiment_design.md` | 실험 설계서 (5개 실험 상세) |
| `research/output/validation_data_index.md` | 검증 데이터 인덱스 |
| `research/output/paper_skeleton.md` | 논문 뼈대 |
| `research/output/competitor_analysis_verified.md` | 경쟁자 분석 |
| `research/output/gap_refinement_cycle3.md` | Gap 분석 |
| `internal/llm/prompt/sloshing_coder.go` | 어블레이션 코드 |
| `research/scripts/compare_spheric.py` | SPHERIC 비교 스크립트 |
| `research/experiments/exp2_nl2xml/scenarios.json` | 20 시나리오 ground truth |
| `datasets/spheric/case_1/` | SPHERIC 실험 데이터 (9파일, 35MB) |

## 커밋 히스토리 (이 브랜치)

```
a1d0dab Phase 1 인프라 준비 (어블레이션 + 비교 스크립트 + 시나리오)
092ddec SPHERIC 유체 수 오류 수정 (3→2종)
8a1a8fa 검증 데이터셋 인덱스
ff8ba5c 실험 설계서
29e8e8d paper skeleton 슬로싱 중심 리프레이밍
a067f86 Gap Cycle 3 + pain points
1a5c2d8 3차 deep analysis
7cde089 2차 사이클 완료
cef9a59 경쟁자 분석
10ea903 RECON→ZERO→SEND 완료
```
