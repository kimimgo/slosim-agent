# Experiment Design — SloshAgent Paper

Date: 2026-02-18
Status: APPROVED (gap analysis cycle 3 기반)
Branch: research/paper

---

## Overview

5개 실험, 4개 Research Question에 대응:

| EXP | RQ | 대응 Gap | 목표 | 산출물 |
|-----|----|---------|----|--------|
| EXP-1 | RQ2 | GAP-3 (Validation) | SPHERIC Test 10 실험 데이터 재현 | Fig: pressure time series, Table: r/NRMSE |
| EXP-2 | RQ1 | GAP-1 (Domain) | NL→XML 생성 정확도 | Table: 20 scenarios × pass/accuracy |
| EXP-3 | RQ3 | GAP-4 (Industry) | 파라메트릭 자동화 | Table: time comparison, Fig: results overlay |
| EXP-4 | RQ4 | GAP-2 (Knowledge) | 도메인 프롬프트 어블레이션 | Table: ON/OFF accuracy delta |
| EXP-5 | — | GAP-4 (Industry) | 산업 시나리오 PoC | Fig: baffle/seismic snapshots |

---

## EXP-1: SPHERIC Benchmark Reproduction (RQ2)

### 목표
에이전트 생성 시뮬레이션이 SPHERIC Test 10 실험 데이터를 재현할 수 있는가?

### 실험 데이터 (이미 보유)
- `datasets/spheric/case_1/` — Water & Oil, lateral/roof impact pressure (100-repeat raw data)
- `datasets/spheric/case_2/`, `case_3/` — 추가 유체 조건

### 에이전트 생성 케이스 (이미 존재)
- `cases/SPHERIC_Test10_Low_Def.xml` — dp=0.006m, ~136K particles
- `cases/SPHERIC_Test10_High_Def.xml` — dp=0.004m, ~344K particles
- `cases/SPHERIC_Test10_Oil_Low_Def.xml` — Oil 유체 (Low res)

### 실행 계획
```
Step 1: GenCase → Solver (GPU) → PartVTK → MeasureTool
        ├─ SPHERIC_Test10_Low  (~136K, 예상 ~5-10분)
        ├─ SPHERIC_Test10_High (~344K, 예상 ~20-40분)
        └─ SPHERIC_Test10_Oil_Low (~136K, 예상 ~5-10분)

Step 2: MeasureTool로 lateral/roof 프로브 위치에서 압력 시계열 추출
        → simulations/spheric_{low,high,oil_low}/measure/pressure_Press.csv

Step 3: 실험 데이터와 비교
        → datasets/spheric/case_1/lateral_water_1x.txt vs agent lateral pressure
        → datasets/spheric/case_1/roof_water_1x.txt vs agent roof pressure

Step 4: English2021 expert mDBC 결과와도 비교 (논문 Fig 3 digitize)
```

### 평가 지표
| 지표 | 정의 | 목표 |
|------|------|------|
| Pearson r | 시계열 상관계수 | r > 0.9 |
| NRMSE | Normalized Root Mean Square Error | < 0.15 |
| Peak Pressure Error | (P_sim - P_exp) / P_exp | < 20% |
| Phase Error | 피크 시간 차이 | < 0.1s |

### 비교 대상
1. SPHERIC 실험 데이터 (100-repeat 평균 ± 2σ)
2. English2021 expert mDBC 결과 (있으면 digitize)
3. Agent Low vs High (해상도 수렴 확인)

### 예상 Figure/Table
- **Fig 2**: Pressure time series — agent(Low) vs agent(High) vs experimental mean ± 2σ
- **Table 2**: r, NRMSE, Peak Error for each probe × each resolution

### 필요 스크립트
- `research/scripts/compare_spheric.py` — 실험 데이터 파싱 + 비교 + 플롯 생성

---

## EXP-2: NL→XML Generation Accuracy (RQ1)

### 목표
자연어 슬로싱 시나리오 설명에서 유효한 DualSPHysics XML을 생성할 수 있는가?

### 20 시나리오 설계 (4 도메인 × 5 복잡도)

#### Level 1: 기본 (명시적 파라미터)
| ID | 도메인 | NL Input | Ground Truth |
|----|--------|----------|-------------|
| S01 | General | "가로 0.9m, 세로 0.062m, 높이 0.508m 직사각형 탱크, 물 93mm 채움, 0.613Hz로 좌우 흔들기" | SPHERIC Test 10 파라미터 |
| S02 | General | "600mm × 300mm × 650mm 탱크에 물 150mm, 주파수 1.2Hz, 진폭 30mm 수평 정현파" | Chen2018 shallow |
| S03 | General | "1m × 0.5m × 0.6m 소형 탱크, 50% 충전, 공진 주파수로 가진" | 기본 공진 케이스 |
| S04 | General | "실험 탱크 크기로 물 반 채우고 천천히 흔들어줘" | 기본 default 추론 |

#### Level 2: 도메인 용어 사용 (파라미터 추론 필요)
| ID | 도메인 | NL Input | Ground Truth |
|----|--------|----------|-------------|
| S05 | LNG | "LNG 화물창 30% 충전 상태에서 공진 주파수 가진" | LNG preset + f1 계산 |
| S06 | Automotive | "자동차 연료탱크 크기로 절반 채우고 급제동 시나리오" | 소형탱크 + braking motion |
| S07 | Nuclear | "원자력 사용후연료 저장조 크기, 수심 10m, 0.3Hz 지진 가진" | 대형 탱크 + 저주파 |
| S08 | Aerospace | "소형 원통형 추진제 탱크, LOX, 70% 충전, 1Hz 횡방향 가진" | 원통형 + 높은 충전 |

#### Level 3: 논문 케이스 재현 (정확한 파라미터 매칭)
| ID | 도메인 | NL Input | Ground Truth |
|----|--------|----------|-------------|
| S09 | Paper | "Chen2018 Case 3 재현: 40% fill level, f/f1=0.9" | Chen2018 정확한 값 |
| S10 | Paper | "Liu2024 pitch sloshing: 50% fill, 주파수 0.8Hz, 진폭 4도" | Liu2024 정확한 값 |
| S11 | Paper | "SPHERIC Test 10 water 조건 재현" | SPHERIC 정확한 값 |
| S12 | Paper | "English2021 mDBC 검증 케이스 재현" | English2021 파라미터 |

#### Level 4: 복합 조건 (다중 기능 조합)
| ID | 도메인 | NL Input | Ground Truth |
|----|--------|----------|-------------|
| S13 | LNG | "LNG 탱크에 수직 배플 2개 설치, 50% 충전, 공진 가진" | 배플 geometry + motion |
| S14 | Seismic | "2003년 토카치오키 지진파로 석유 탱크 슬로싱 시뮬레이션" | seismic_input + 대형 탱크 |
| S15 | STL | "fuel_tank.stl 형상 임포트하고 50% 물 채워서 0.5Hz 가진" | stl_import + fill |
| S16 | Multi | "Chen2018 6개 충전 레벨 파라메트릭 스터디 실행" | parametric_study orchestration |

#### Level 5: 엣지 케이스/에러 복구
| ID | 도메인 | NL Input | Ground Truth |
|----|--------|----------|-------------|
| S17 | Edge | "dp 0.001로 아주 정밀하게 시뮬레이션" | dp 하한 경고 → 조정 |
| S18 | Edge | "공기만 있는 빈 탱크 흔들기" | 유체 없음 → 에러 안내 |
| S19 | Edge | "탱크 100% 가득 채우고 흔들기" | 자유표면 없음 → 경고 |
| S20 | Edge | "10m × 10m × 10m 탱크, dp=0.002" | ~125M particles → 스케일 경고 |

### 평가 방법
```
각 시나리오에 대해:
1. Agent에 NL input 전달
2. xml_generator 호출 여부 및 파라미터 확인
3. GenCase 실행 (pass/fail)
4. 생성된 파라미터를 ground truth와 비교

평가 기준:
- GenCase Pass: XML이 GenCase를 통과하는가? (binary)
- Parameter Accuracy: 각 파라미터의 정확도 (%)
  - tank dimensions: exact match 필요
  - dp: ±20% 허용
  - frequency: ±5% 허용
  - fill level: ±10% 허용
  - motion type: exact match 필요
- Physical Validity: 생성된 설정이 물리적으로 합리적인가? (expert judgment)
```

### 실행 방법
```bash
# Non-interactive mode (-p flag)로 각 시나리오 실행
./slosim-agent -p "S01의 NL input" 2>&1 | tee research/logs/exp2_s01.log

# 또는 Ollama API 직접 호출로 XML 생성만 테스트
```

### 예상 Table
- **Table 3**: NL→XML Generation Results (20 scenarios)

| Level | n | GenCase Pass | Param Accuracy | Physical Valid |
|-------|---|-------------|---------------|---------------|
| L1 (기본) | 4 | ?/4 | ?% | ?/4 |
| L2 (도메인) | 4 | ?/4 | ?% | ?/4 |
| L3 (논문) | 4 | ?/4 | ?% | ?/4 |
| L4 (복합) | 4 | ?/4 | ?% | ?/4 |
| L5 (엣지) | 4 | ?/4 | ?% | ?/4 |
| **Total** | **20** | **?/20** | **?%** | **?/20** |

---

## EXP-3: Parametric Study Automation (RQ3)

### 목표
에이전트가 단일 NL 프롬프트로 다중 케이스 파라메트릭 스터디를 자동화할 수 있는가?

### 실험 설계
```
NL Input: "Chen2018 논문의 6개 충전 레벨(20%, 30%, 40%, 50%, 60%, 70%)에 대해
          f/f1=0.9 조건으로 파라메트릭 스터디를 실행하고 결과를 비교해줘"

Expected Agent Behavior:
1. parametric_study tool 호출 (또는 반복 xml_generator)
2. 6개 XML 생성 (fill level만 변경)
3. 6개 GenCase + Solver 실행
4. MeasureTool로 각 케이스 압력/수위 추출
5. 결과 비교 리포트 생성
```

### 평가 지표
| 지표 | Manual (Expert) | Agent | 목표 |
|------|----------------|-------|------|
| Setup Time | ~2hr/case × 6 = 12hr | Single prompt → ? min | 10x+ 감소 |
| XML Correctness | 6/6 | ?/6 | 6/6 |
| Solver Completion | 6/6 | ?/6 | 6/6 |
| Result Accuracy | Chen2018 published | vs Chen2018 | NRMSE < 0.2 |

### 비교 데이터
- Chen et al. 2018: OpenFOAM 결과 (Figure digitize 필요)
- 탱크: 600mm × 300mm × 650mm
- 주파수: f/f1 = 0.9 (near-critical)
- 6 fill levels: 120, 150, 195, 260, 325, 390 mm

### 예상 Figure
- **Fig 3**: 6-panel 결과 — agent vs Chen2018 (각 fill level별 수위 시계열)
- **Table 4**: Setup time comparison (manual 12hr vs agent ?min)

---

## EXP-4: Domain Prompt Ablation (RQ4)

### 목표
SloshingCoderPrompt의 도메인 지식이 시뮬레이션 성능에 얼마나 기여하는가?

### 어블레이션 조건

| Condition | Prompt Content | 코드 변경 |
|-----------|---------------|----------|
| **FULL** | SloshingCoderPrompt 전체 (135줄) | 현재 그대로 |
| **NO-DOMAIN** | 도메인 지식 섹션 제거 (공진 공식, 탱크 프리셋, 용어 매핑) | 50-66줄 제거 |
| **NO-RULES** | Tool 호출 순서/경로 규칙 제거 | 26-86줄 제거 |
| **GENERIC** | "You are a helpful AI assistant for DualSPHysics simulation." (1줄) | 전체 교체 |

### 테스트 시나리오 (EXP-2에서 선별)
EXP-2의 20개 중 대표 10개 사용 (Level별 2개씩):
- L1: S01, S03
- L2: S05, S07
- L3: S09, S11
- L4: S13, S15
- L5: S17, S19

### 평가 지표
각 condition × 10 scenarios = 40 실행

| Metric | FULL | NO-DOMAIN | NO-RULES | GENERIC |
|--------|------|-----------|----------|---------|
| GenCase Pass Rate | ?/10 | ?/10 | ?/10 | ?/10 |
| Parameter Accuracy | ?% | ?% | ?% | ?% |
| Tool Call Correctness | ?/10 | ?/10 | ?/10 | ?/10 |
| Physical Validity | ?/10 | ?/10 | ?/10 | ?/10 |

### 구현 계획
```go
// sloshing_coder.go에 ablation mode 추가
func SloshingCoderPrompt(provider models.ModelProvider, ablation string) string {
    switch ablation {
    case "full":
        return sloshingSystemPrompt  // 현재 그대로
    case "no-domain":
        return sloshingNoDomainPrompt  // 도메인 지식 제거
    case "no-rules":
        return sloshingNoRulesPrompt  // 규칙 제거
    case "generic":
        return genericPrompt  // 최소 프롬프트
    }
}
```

환경변수 또는 config로 ablation mode 전환:
```bash
SLOSIM_ABLATION=no-domain ./slosim-agent -p "scenario..."
```

### 예상 Figure
- **Fig 4**: Ablation bar chart — 4 conditions × GenCase pass rate + accuracy
- **Table 5**: Detailed ablation results (10 scenarios × 4 conditions)

---

## EXP-5: Industrial Scenario PoC

### 목표
실제 산업 시나리오에서 SloshAgent의 실용성을 시연

### 시나리오 A: Anti-Slosh Baffle Design
```
NL: "직사각형 탱크(1m × 0.5m × 0.6m)에 수직 배플 1개를 중앙에 설치하고,
     물 50% 채운 상태에서 공진 주파수로 가진했을 때 배플 유무에 따른
     벽면 압력 차이를 비교해줘"

Expected:
- 2 simulations: with/without baffle
- Pressure comparison at wall probe
- Force reduction quantification
```

### 시나리오 B: Seismic Excitation
```
NL: "10m × 5m × 8m 석유 저장탱크, 물 60% 채움, 0.3Hz 지진파 가진으로
     슬로싱 시뮬레이션 실행"

Expected:
- seismic_input tool 사용
- Large tank configuration
- Low-frequency excitation
```

### 평가 (정성적)
- 에이전트가 시나리오를 올바르게 해석하는가?
- 적절한 tool sequence를 선택하는가?
- 결과가 물리적으로 합리적인가?
- 비전문가가 이해할 수 있는 분석 리포트를 생성하는가?

### 예상 Figure
- **Fig 5**: Baffle comparison — particle snapshots (with/without) + pressure time series
- **Fig 6**: Seismic scenario — particle snapshots at key time steps

---

## 실행 순서 및 의존성

```
Phase 1: 인프라 준비 (1-2일)
├─ 1a. 어블레이션 코드 구현 (sloshing_coder.go 수정)
├─ 1b. 비교 스크립트 작성 (compare_spheric.py)
└─ 1c. NL 시나리오 20개 ground truth JSON 작성

Phase 2: GPU 시뮬레이션 실행 (2-3일)
├─ 2a. EXP-1: SPHERIC 3 cases 실행 (Low/High/Oil)
├─ 2b. EXP-5: Baffle + Seismic 2 cases 실행
└─ 2c. EXP-3: Chen2018 6-case parametric 실행

Phase 3: Agent 테스트 (2-3일)
├─ 3a. EXP-2: 20 scenarios × FULL prompt
├─ 3b. EXP-4: 10 scenarios × 4 ablation conditions (40 runs)
└─ 3c. EXP-3: Agent-driven parametric (시간 측정)

Phase 4: 분석 및 시각화 (1-2일)
├─ 4a. SPHERIC 비교 (r, NRMSE, peak error)
├─ 4b. EXP-2/4 결과 테이블 집계
├─ 4c. Figure 생성 (matplotlib/plotly)
└─ 4d. Chen2018 digitize + 비교

Phase 5: 논문 실험 섹션 작성 (1-2일)
```

### 총 예상 소요: 7-12일 (병렬 실행 시 5-7일)

---

## Hardware/Software Requirements

| 항목 | 사양 | 비고 |
|------|------|------|
| GPU | RTX 4090 (24GB VRAM) | Ollama 언로드 후 solver 실행 |
| CUDA | 12.6 | Docker image |
| Docker | dsph-agent:latest | DualSPHysics v5.4 GPU |
| LLM | Qwen3 32B via Ollama | ~20GB VRAM |
| Python | 3.10+ | 비교 스크립트용 |
| Packages | numpy, pandas, matplotlib, scipy | 데이터 분석 |

**GPU VRAM 충돌**: Ollama(Qwen3 32B, ~20GB)와 DualSPHysics GPU는 동시 실행 불가.
- EXP-1/3/5 (solver 실행): Ollama 중지 → solver 실행 → Ollama 재시작
- EXP-2/4 (agent 테스트): Ollama 활성 → XML 생성만 (solver 불필요 또는 순차)

---

## 데이터 저장 구조

```
research/
├─ experiments/
│  ├─ exp1_spheric/
│  │  ├─ results/          ← 시뮬레이션 측정 데이터
│  │  ├─ comparison/       ← 실험 vs 시뮬레이션 비교
│  │  └─ figures/          ← 생성된 Figure
│  ├─ exp2_nl2xml/
│  │  ├─ scenarios.json    ← 20 시나리오 정의 + ground truth
│  │  ├─ logs/             ← Agent 실행 로그
│  │  └─ results.csv       ← 평가 결과
│  ├─ exp3_parametric/
│  │  ├─ results/
│  │  └─ comparison/
│  ├─ exp4_ablation/
│  │  ├─ logs/             ← 40 실행 로그
│  │  └─ results.csv
│  └─ exp5_industrial/
│     ├─ baffle/
│     └─ seismic/
├─ scripts/
│  ├─ compare_spheric.py
│  ├─ evaluate_nl2xml.py
│  └─ plot_results.py
└─ output/
   └─ (기존 문서들)
```
