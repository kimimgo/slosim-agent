# Research-v3: 실험 계획 (Detailed)

## 논문 주장 → 실험 매핑

```
Claim                          Experiment    Status
──────────────────────────     ──────────    ──────
C1. NL→SPH 파이프라인 작동     EXP-A         NEW
C2. 아키텍처 컴포넌트 필수성    EXP-B         NEW
C3. 물리적 신뢰성              EXP-C         v2 완료 (이관)
```

---

## EXP-A: E2E Pipeline Test — NL→Tool Call→GPU

### 목적

비전문가의 자연어 입력이 SlosimAgent를 통해 DualSPHysics GPU 시뮬레이션으로
자율 변환되는지 검증한다. NL2FOAM (Xu et al. 2025)의 21-case 벤치마크 방법론을
tool-calling 패러다임에 적용한다.

### NL2FOAM과의 방법론적 차이

| 항목 | NL2FOAM | SlosimAgent (본 연구) |
|------|---------|----------------------|
| 패러다임 | Fine-tuning (text→config) | Tool-calling (NL→tool→XML→solver→post) |
| 솔버 | OpenFOAM (격자 FVM) | DualSPHysics (입자 SPH) |
| 단계 수 | 1단계 (config 생성) | 5단계 (tool 선택→XML→GenCase→Solver→Post) |
| 모델 | Qwen2.5-7B fine-tuned | Qwen3-32B zero-shot (도메인 프롬프트) |
| 데이터셋 | 28,716 합성 쌍 | 10 실험 벤치마크 시나리오 |
| 난이도 분류 | Easy/Medium/Hard (21-case) | Easy/Medium/Hard (10-case, 아래 정의) |

### 시나리오 분류 (Easy / Medium / Hard)

NL2FOAM의 난이도 분류를 참고하되, SPH 도메인 특성을 반영:
- **Easy**: 직사각 탱크, 단일 운동, 표준 유체
- **Medium**: 대형 탱크, 공진/비선형, 파라메트릭, 복합 운동
- **Hard**: 비표준 형상(원통형, STL), 고급 경계조건

| Tier | # | 시나리오 | 핵심 난이도 요소 | GPU 결과 (수동) |
|------|---|---------|-----------------|----------------|
| **Easy** | S1 | SPHERIC Oil Low Fill | 유체 물성 변경만 | PASS (131s) |
| | S2 | Chen Shallow Sway | 수평 병진 운동 | PASS (174s) |
| | S8 | English mDBC vs DBC | BC 파라미터 변경 | PARTIAL (DBC 한계) |
| **Medium** | S3 | Chen Near-Critical | 비선형 전이 구간 | PASS (430s) |
| | S4 | Liu Large Pitch 30s | 대형 탱크 + 공진 | PASS (738s) |
| | S5 | Liu Amplitude Parametric | 3-case 파라메트릭 | PASS/PASS/PARTIAL |
| | S6 | ISOPE LNG Roof Impact | 복합 운동 + 고충전 | PASS (sway only) |
| **Hard** | S7 | NASA Cylinder Baffle | 원통형 geometry | PASS (216s) |
| | S9 | Zhao Horizontal Cyl | 수평 원통 STL | PASS (~90% fill 타협) |
| | S10 | Frosina Fuel Tank | 복잡 3D STL + 급제동 | PASS (STL+fillpoint) |

### 실행 방법

비대화식 모드로 에이전트에 NL 프롬프트를 입력하고 전체 파이프라인 자율 실행:

```bash
# 단일 시나리오 실행
./slosim-agent -p "${NL_PROMPT}" -q -f json 2>&1 | tee exp-a/S${N}_log.json

# 타임아웃 (1시간)
timeout 3600 ./slosim-agent -p "${NL_PROMPT}" -q -f json
```

각 시나리오는 **3회 반복** 실행하여 재현성 확인 (LLM non-determinism 감안).

### NL 프롬프트 (10개)

각 프롬프트는 비전문가가 실제로 입력할 수준으로 작성. 전문 용어 최소화.

```
S1:  "SPHERIC 벤치마크 탱크(900x62x508mm)에 해바라기 오일을 18% 채워서
      슬로싱 해석해줘. 오일 밀도 990, 점도 0.045 Pa·s.
      주기 1.535초로 바닥 중심 기준 피치 회전 4도."

S2:  "600x300x650mm 직사각형 탱크에 물 83mm 높이로 채워줘.
      수평 방향(x축)으로 진폭 7mm, 주파수 0.756Hz로 좌우 흔들어.
      10초 동안 해석."

S3:  "같은 탱크(600x300x650mm)인데 이번엔 물 185mm로 채워줘.
      수평 진동 7mm, 자연주파수(1.008Hz)로 가진.
      이 수위가 critical depth 근처라서 비선형 효과가 클 거야."

S4:  "1미터 정사각 탱크 (1000x500x1000mm)에 물 200mm 채워서
      피치 공진 해석해줘. 공진주파수가 0.66Hz이고 진폭 2도야.
      회전축은 탱크 바닥 중심. 30초 해석."

S5:  "1미터 탱크에 물 70% 채우고 공진주파수(0.87Hz) 피치 회전인데,
      진폭을 1도, 2도, 3도로 바꿔가며 파라메트릭 해석 해줘.
      각각 20초씩."

S6:  "ISOPE 벤치마크 LNG 탱크 (946x118x670mm) 슬로싱 해석.
      물 90% 채우고 수평 진동 f=0.6Hz 진폭 20mm.
      천장 압력 충격을 보고 싶어."

S7:  "직경 2.84m, 높이 3m 원통형 탱크에 물 50% 채워줘.
      수평 방향 0.5Hz 사인파 10mm 진폭으로 흔들어."

S8:  "SPHERIC 벤치마크 탱크(900x62x508mm) 물 18% 해석인데,
      경계조건을 DBC로 해줘. dp=0.002m로 고해상도로.
      센서1 위치 압력 시계열 추출."

S9:  "수평 원통형 탱크 (직경 1m, 길이 3m)에 물 25% 채워서
      피치 3도, 주파수 0.55Hz로 슬로싱 해석해줘."

S10: "자동차 연료탱크 STL 파일(cases/fuel_tank.stl)로 슬로싱 해석해줘.
      물 절반 채우고, 수평 방향으로 진폭 50mm, 0.5Hz로 흔들어."
```

### 메트릭 정의 (5-Stage Pipeline)

NL2FOAM의 first-attempt success rate를 5단계 파이프라인으로 확장:

| 메트릭 | 정의 | 측정 방법 | Pass 기준 |
|--------|------|----------|----------|
| **M-A1: Tool Selection** | 올바른 도구를 올바른 순서로 호출 | 로그에서 tool call 시퀀스 추출 | 필수 도구 100% 호출 |
| **M-A2: XML Validity** | 생성된 XML이 GenCase 통과 | GenCase exit code = 0 | 0 errors |
| **M-A3: Parameter Fidelity** | 사용자 지정 파라미터 정확 반영 | XML 파싱 → NL 프롬프트 대조 | 핵심 파라미터 100% 일치 |
| **M-A4: Solver Completion** | GPU 솔버 정상 완주 | DualSPHysics exit code + PART 수 | 예상 PART ≥95% |
| **M-A5: Pipeline Success** | 전체 파이프라인 완주 (= first-attempt success) | M-A1~A4 모두 PASS | 모두 PASS |

#### M-A3 Parameter Fidelity 상세 체크리스트

각 시나리오별로 NL 프롬프트에서 지정한 파라미터가 XML에 정확히 반영되었는지 확인:

| 파라미터 유형 | 체크 항목 | 허용 오차 |
|-------------|----------|----------|
| 탱크 치수 | pointmin/pointmax에서 계산 | ±1% |
| 유체 높이 | fluid drawbox z 범위 | ±5% |
| 유체 물성 | rhop0, visco | 정확 일치 |
| 여진 주파수 | mvrotsinu/mvrectsinu freq | ±1% |
| 여진 진폭 | ampl 값 | ±1% |
| 시뮬레이션 시간 | TimeMax | 정확 일치 |

#### 예상 Tool Call 시퀀스 (M-A1 기준)

```
Easy (S1,S2,S8):
  xml_generator → gencase → solver → partvtk/measuretool

Medium (S3,S4,S6):
  xml_generator → gencase → solver → partvtk/measuretool

Medium-Parametric (S5):
  parametric_study → [xml_generator → gencase → solver] × 3 → result_store

Hard-Cylinder (S7):
  geometry(cylindrical) → xml_generator → gencase → solver → partvtk

Hard-STL (S9,S10):
  stl_import → xml_generator → gencase → solver → partvtk
```

### 모델 비교

| 모델 | 파라미터 | 용도 | 비고 |
|------|---------|------|------|
| **Qwen3-32B** | 32B | Primary | zero-shot + 도메인 프롬프트 |
| **Qwen3-8B** | 8B | Comparison | Quality Report v0.3에서 0/3 PASS |

NL2FOAM은 fine-tuned 7B vs zero-shot 72B/DeepSeek-R1을 비교.
본 연구는 **zero-shot 도메인 프롬프트**의 모델 크기별 성능 차를 비교.
8B는 이미 3개 시나리오에서 0/3 PASS (무한 루프, 도구 오선택, 파싱 실패) 기록.
→ 10개 전체에서 확인하여 "32B 필수" 주장 뒷받침.

### 산출물

| 산출물 | 형태 | 내용 |
|-------|------|------|
| **Table 3** | 10×5 매트릭스 | 시나리오별 M-A1~A5 PASS/FAIL (32B) |
| **Table 3b** | 10×5 매트릭스 | 동일, Qwen3-8B |
| **Fig 2** | 흐름도 | S1 파이프라인 상세 예시 (NL→tool calls→XML→GPU→결과) |
| **Fig 3** | Grouped bar chart | 난이도별 Pipeline Success Rate (32B vs 8B) |

### GPU 자원 추정

```
10 시나리오 × 3회 반복 × 2 모델 = 60 실행
평균 GPU 시간 ~10분/실행 → ~10시간 총 GPU 시간
RTX 4090 1대, Ollama VRAM 교대 관리 필수
```

---

## EXP-B: Ablation Study — 아키텍처 컴포넌트 필수성

### 목적

SlosimAgent의 도메인 특화 아키텍처에서 각 컴포넌트(도메인 프롬프트, XML 도구,
후처리 도구)를 제거했을 때 파이프라인 성공률이 어떻게 변하는지 검증.
"없으면 안 된다"를 증명하여 Contribution 2 (아키텍처 유효성)를 지지.

### 실험 조건 (5가지)

| ID | 조건 | 도메인 프롬프트 | DualSPHysics 도구 18개 | 후처리 도구 | XML 템플릿 |
|----|------|:---:|:---:|:---:|:---:|
| B0 | **Full (SlosimAgent)** | ✓ | ✓ | ✓ | ✓ |
| B1 | **−DomainPrompt** | ✗ (generic) | ✓ | ✓ | ✓ |
| B2 | **−DSPHTool** | ✓ | ✗ (raw LLM XML) | ✓ | ✗ |
| B3 | **−PostProcess** | ✓ | ✓ | ✗ | ✓ |
| B4 | **Bare LLM** | ✗ | ✗ | ✗ | ✗ |

#### 조건별 구현 방법

- **B0 Full**: 현재 SlosimAgent 그대로 실행
- **B1 −DomainPrompt**: `SloshingCoderPrompt`를 generic `CoderPrompt`로 교체.
  도구 18개는 그대로 제공하되, SPH 도메인 지식 프롬프트만 제거.
- **B2 −DSPHTool**: xml_generator/gencase/solver 등 DualSPHysics 전용 도구 제거.
  LLM이 직접 XML 텍스트를 생성하도록 유도 (file_write 도구만 허용).
  GenCase는 생성된 XML에 수동 적용.
- **B3 −PostProcess**: partvtk/measuretool/analysis/report 제거.
  XML→GenCase→Solver까지만 실행. 후처리 단계 성공률은 측정 불가 → N/A.
- **B4 Bare LLM**: 모든 도구 제거. Qwen3-32B에 NL 프롬프트만 입력.
  "DualSPHysics XML을 생성해줘"라고 추가 지시.
  출력된 텍스트를 XML 파일로 저장 → GenCase 수동 적용.

### 시나리오 (3개 — 난이도별 대표)

| Tier | 시나리오 | 선정 이유 |
|------|---------|----------|
| Easy | S1: SPHERIC Oil | 기본 직사각형, 단일 운동 |
| Medium | S4: Liu Large Pitch | 대형 탱크, 공진, 장시간 |
| Hard | S10: Frosina Fuel Tank | STL 복잡 형상 |

### 메트릭 (EXP-A 메트릭 재사용 + 추가)

| 메트릭 | B0 | B1 | B2 | B3 | B4 |
|--------|:--:|:--:|:--:|:--:|:--:|
| M-A1 Tool Selection | ✓ | ✓ | N/A | ✓ | N/A |
| M-A2 XML Validity | ✓ | ✓ | ✓ | ✓ | ✓ |
| M-A3 Parameter Fidelity | ✓ | ✓ | ✓ | ✓ | ✓ |
| M-A4 Solver Completion | ✓ | ✓ | ✓ | ✓ | ✓ |
| **M-B1 Error Count** | ✓ | ✓ | ✓ | ✓ | ✓ |

**M-B1 Error Count**: LLM hallucination, XML 문법 오류, 잘못된 단위, 누락 태그 등의 총 개수.

### 예상 결과 (가설)

```
           S1(Easy)   S4(Med)   S10(Hard)
B0 Full:    PASS       PASS       PASS
B1 -Prompt: PASS       PARTIAL    FAIL     ← 도메인 지식 없이 공진/STL 실패
B2 -Tool:   FAIL       FAIL       FAIL     ← raw LLM XML은 GenCase 통과 불가
B3 -Post:   PASS*      PASS*      PASS*    ← 솔버까지는 성공, 결과 해석 불가 (*후처리 N/A)
B4 Bare:    FAIL       FAIL       FAIL     ← 도구 없이 XML 생성 자체 불가
```

### 산출물

| 산출물 | 형태 | 내용 |
|-------|------|------|
| **Table 4** | 5×3 매트릭스 | 조건별 × 시나리오별 PASS/PARTIAL/FAIL |
| **Fig 4** | Stacked bar chart | 조건별 M-A2~A4 성공률 |

### GPU 자원 추정

```
5 조건 × 3 시나리오 × 1회 = 15 실행 (B2, B4는 GenCase 실패 시 Solver 불필요)
실질 GPU 실행: ~8회 × 평균 10분 = ~1.5시간
```

---

## EXP-C: SPHERIC Test 10 물리 검증 (v2에서 이관)

### 목적

에이전트가 생성한 시뮬레이션의 물리적 정확도를 SPHERIC Test 10 벤치마크로 정량 검증.
Souto-Iglesias et al. (2011)의 102회 반복 실험 데이터와 비교.

### v2에서 이관할 데이터

| 데이터 | 경로 | 내용 |
|--------|------|------|
| run_001 | `research-v2/exp1_spheric/run_001/` | Water Lat dp=4mm |
| run_002 | `research-v2/exp1_spheric/run_002/` | Water Lat dp=2mm |
| run_003 | `research-v2/exp1_spheric/run_003/` | Water Lat dp=1mm |
| run_005 | `research-v2/exp1_spheric/run_005/` | Oil Lat dp=2mm |
| run_006~008 | `research-v2/exp1_spheric/run_006~008/` | Water Roof |
| 분석 스크립트 | `research-v2/exp1_spheric/scripts/` | convergence, paper_figures |
| 메트릭 JSON | `research-v2/exp1_spheric/analysis/` | metrics.json |
| Figure | `research-v2/exp1_spheric/figures/` | fig_timeseries, convergence 등 |

### M1-M8 메트릭 프레임워크

| 메트릭 | 정의 | Water Lat | Oil Lat | Water Roof |
|--------|------|:---------:|:-------:|:----------:|
| M1 Crash-free | 솔버 완주 | ✓ | ✓ | ✓ |
| M2 Peak Error | 첫 피크 오차 (%) | 19.5% | 18.7% | — |
| M3 Phase Error | 위상 오차 (ms) | <50ms | <60ms | — |
| M4 Cross-corr | 피어슨 상관계수 r | 0.655 | 0.61 | — |
| M5 RMS Error | NRMSE | 0.35 | 0.38 | — |
| M6 Trend Match | 물리적 경향 일치 | ✓ | ✓ | PARTIAL |
| M7 Conservation | 에너지/질량 보존 | ✓ | ✓ | ✓ |
| M8 Convergence | dp 수렴 | r: 0.460→0.655→0.697 | — | — |

### 추가 작업: 없음

v2 결과를 정리하여 논문에 포함. 추가 시뮬레이션 불필요.

### 산출물

| 산출물 | 형태 | 내용 |
|-------|------|------|
| **Table 2** | 정의표 | M1-M8 메트릭 정의 + Pass 기준 |
| **Table 5** | 결과표 | Water/Oil/Roof × M1-M8 결과 |
| **Table 6** | 수렴표 | dp별 r, GCI 값 |
| **Fig 5** | 시계열 | Water Lateral 3-resolution vs 실험 |
| **Fig 6** | 수렴 곡선 | r vs dp (3-level) |
| **Fig 7** | 시계열 | Oil Lateral vs 실험 |

---

## 실험 순서 및 일정

```
Phase 1: EXP-A (NL→Tool Call)        ← 즉시 시작
  ├─ Step 1: Qwen3-32B × 10 시나리오 × 3회 (GPU ~5h)
  ├─ Step 2: Qwen3-8B × 10 시나리오 × 1회 (GPU ~3h, 대부분 조기종료 예상)
  └─ Step 3: 결과 정리 → Table 3, Fig 2-3

Phase 2: EXP-B (Ablation)            ← EXP-A B0 결과 재사용
  ├─ Step 1: B1~B4 조건 구현 (코드 분기)
  ├─ Step 2: 3 시나리오 × 4 조건 실행 (GPU ~1.5h)
  └─ Step 3: 결과 정리 → Table 4, Fig 4

Phase 3: EXP-C (이관)                ← v2 데이터 복사 + 정리
  └─ Step 1: v2 결과 → Table 2,5,6 + Fig 5-7

Phase 4: 논문 작성                    ← 모든 산출물 완성 후
```

## Figure 목록 (최종)

| # | 내용 | 실험 | 상태 |
|---|------|------|------|
| Fig 1 | 시스템 아키텍처 다이어그램 | — | NEW |
| Fig 2 | NL→tool call 파이프라인 예시 (S1 상세) | EXP-A | NEW |
| Fig 3 | E2E Pipeline Success: Grouped bar (Easy/Med/Hard × 32B/8B) | EXP-A | NEW |
| Fig 4 | Ablation: Stacked bar (B0~B4 × S1/S4/S10) | EXP-B | NEW |
| Fig 5 | Water Lateral 시계열 3-resolution vs 실험 | EXP-C | v2 완료 |
| Fig 6 | 수렴 분석 (r vs dp) | EXP-C | v2 완료 |
| Fig 7 | Oil Lateral vs 실험 | EXP-C | v2 완료 |

## Table 목록 (최종)

| # | 내용 | 실험 |
|---|------|------|
| Table 1 | 경쟁 에이전트 비교 (7종 vs Ours) | 서베이 |
| Table 2 | M1-M8 메트릭 정의 + Pass 기준 | EXP-C |
| Table 3 | E2E 10시나리오 × 5단계 성공률 (32B/8B) | EXP-A |
| Table 4 | Ablation 5조건 × 3시나리오 매트릭스 | EXP-B |
| Table 5 | SPHERIC T10 결과 (Water/Oil/Roof × M1-M8) | EXP-C |
| Table 6 | 수렴 분석 (dp, r, GCI) | EXP-C |

## 결론 (논문에서 주장할 것)

1. **SlosimAgent는 비전문가의 자연어를 SPH 시뮬레이션으로 변환할 수 있다**
   - 증거: EXP-A — 10 시나리오 중 N/10 Pipeline Success (32B)
   - 난이도별 성공률: Easy > Medium > Hard (NL2FOAM 패턴과 유사 예상)
   - 한계: 8B 모델은 0~1/10 성공 → 32B 이상 필수

2. **도메인 특화 설계의 각 컴포넌트가 필수적이다**
   - 증거: EXP-B — Full(B0)만 모든 시나리오 성공
   - B2(−DSPHTool): LLM이 직접 생성한 XML은 GenCase 통과 불가
   - B4(Bare LLM): DualSPHysics XML 문법 자체를 모름
   - 핵심 발견: 도구 없는 LLM은 attribute-only SPH XML을 생성할 수 없다

3. **에이전트가 생성한 시뮬레이션은 물리적으로 신뢰할 수 있다**
   - 증거: EXP-C — SPHERIC T10 Water PASS (M2=19.5%, r=0.655), Oil PASS
   - 수렴: dp 감소 → r 단조 증가 (0.460→0.655→0.697)
   - 부수 발견: SPH 슬로싱 문헌의 정량 검증 관행 부재 → M1-M8 제안

---

## Appendix: NL2FOAM 참조 정보

- 논문: Xu et al. (2025), "NL2FOAM: Natural Language to OpenFOAM"
- 데이터: HuggingFace `Xu-Zimu/NL2FOAM` (28,716 NL→config 쌍, 16 base cases)
- 벤치마크: 21-case (7 Easy + 7 Medium + 7 Hard)
- 주요 메트릭: First-attempt Success Rate (82.6%), Solution Accuracy (88.7%)
- 모델: fine-tuned Qwen2.5-7B vs Qwen2.5-72B/DeepSeek-R1-671B/Llama-70B
- 핵심 차이: NL2FOAM은 fine-tuning, 본 연구는 zero-shot tool-calling

---

## 알려진 제한점 및 대응

> 이 섹션은 실험 설계의 논리적 약점을 사전에 식별하고, 논문 Discussion에서 정직하게 보고하기 위한 것이다.

### L1. NL 프롬프트 기술 수준의 비균질성

**문제**: 논문은 "비전문가 접근성"을 핵심 주장으로 내세우지만, 10개 프롬프트의 기술 수준이 균일하지 않다.

| 수준 | 프롬프트 | 근거 |
|------|---------|------|
| **Novice** (비전문가) | S07, S09 | 기본 치수+운동만 기술. 전문 용어 없음 |
| **Intermediate** (준전문가) | S02, S04, S10 | "공진주파수", "STL 파일" 등 일부 전문 용어 사용 |
| **Expert** (전문가) | S01, S03, S05, S06, S08 | "SPHERIC 벤치마크", "critical depth", "DBC", "dp", "ISOPE", "파라메트릭 해석" |

특히 문제가 되는 프롬프트:
- **S08**: "경계조건을 DBC로", "dp=0.002m" — SPH 전문가만 아는 용어
- **S03**: "critical depth 근처라서 비선형 효과가 클 거야" — 유체역학 이론 지식 전제
- **S06**: "ISOPE 벤치마크 LNG 탱크" — 학회 벤치마크 인용은 연구자 수준

**대응**:
1. 논문 Table 3에 각 프롬프트의 **기술 수준 라벨** (Novice/Intermediate/Expert) 열 추가
2. Discussion에서 명시: "본 연구의 프롬프트는 비전문가~전문가 스펙트럼을 포함하며, 순수 비전문가(Novice) 프롬프트는 S07, S09 2개에 해당한다. 타겟 사용자인 탱크 검사원/구조 엔지니어는 Intermediate 수준으로 가정한다."
3. 향후 연구로 "프롬프트 기술 수준별 성공률 분석" 제안 (Novice-only 벤치마크)

### L2. 기존 에이전트 직접 실험 비교 부재

**문제**: 7종 LLM-CFD 에이전트(OpenFOAMGPT, ChatCFD, Foam-Agent 등)와 동일 태스크에서 직접 실험 비교가 없다. Table 1은 기능 비교(feature comparison)에 불과하다.

**직접 비교가 불가능한 구조적 이유**:
1. **솔버 불일치**: 7종 전부 OpenFOAM(격자 FVM) 전용. DualSPHysics(SPH)와 동일 태스크를 공평하게 비교할 기준이 없음
2. **접근 불가**: OpenFOAMGPT, CFD-GPT 등 대부분 오픈소스가 아니거나 재현 불가
3. **도메인 불일치**: 기존 에이전트는 범용 CFD를 목표. 슬로싱 특화 태스크 자체가 테스트 범위 밖

**NL2FOAM과의 간접 비교는 유효**:
- NL2FOAM은 유일하게 벤치마크+메트릭을 공개한 연구
- 동일 메트릭(first-attempt success rate)으로 **방법론적 비교** 가능:
  - NL2FOAM: fine-tuned 7B → 82.6% (21 cases)
  - SlosimAgent: zero-shot 32B → N/10 (10 cases)
- 솔버가 다르므로 "동급 성능" 주장은 삼가되, "tool-calling도 competitive"임을 보임

**대응**:
1. Discussion에 "Direct Comparison Limitations" 서브섹션 추가
2. Table 1에 비교 불가 사유 각주 추가: "Direct experimental comparison is infeasible due to solver incompatibility (FVM vs SPH) and limited source availability"
3. NL2FOAM과의 메트릭 정렬(M-A5 ≈ first-attempt success rate)을 명시하여 간접 비교 근거 강화

### L3. EXP-B B3(−PostProcess) 해석의 뉘앙스

**문제**: B3에서 후처리 도구(partvtk, measuretool, analysis, report) 제거 시 "솔버까지 성공"으로 예상. 이를 단순히 "PASS"로 기록하면 ablation의 의미가 희석된다.

**핵심 구분**:

| 관점 | B3 결과 | 해석 |
|------|---------|------|
| **계산적 필수성** (Computational Necessity) | PASS* — 솔버 완주 가능 | 후처리는 시뮬레이션 자체에 불필요 |
| **실용적 필수성** (Practical Necessity) | FAIL — 비전문가가 결과 해석 불가 | 후처리 없이는 에이전트의 존재 의의 상실 |

**대응**:
1. B3 결과를 "PASS*" 표기하고, 각주: "Solver completes but results are inaccessible without post-processing tools"
2. M-A5 (Pipeline Success) 정의를 확장: **"사용자가 결과를 해석 가능한 형태로 받는 것"** 까지 포함
3. Table 4에 **M-B2 (User Interpretability)** 메트릭 추가 고려:
   - "사용자에게 유의미한 출력(VTK, 시계열 plot, 요약 리포트) 제공 여부"
   - B3는 M-B2 = FAIL → "후처리는 계산적으로 선택이나, 사용자 경험 관점에서 필수"

### L4. S03 프롬프트의 문맥 의존성

**문제**: S03 "같은 탱크(600x300x650mm)인데 이번엔..." — "같은 탱크"는 S02를 전제로 하는 대화 문맥이다. 그러나 EXP-A는 비대화식 모드(단일 프롬프트)로 실행하므로, 에이전트는 S02 문맥이 없다.

**영향**: 괄호 안 치수(600x300x650mm)가 있어 기능적으로는 문제 없을 수 있으나, "같은 탱크"라는 표현이 LLM의 혼란을 유발할 가능성 존재.

**대응**:
1. S03 프롬프트를 자립적(self-contained)으로 수정 권장:
   - 변경 전: "같은 탱크(600x300x650mm)인데 이번엔..."
   - 변경 후: "600x300x650mm 직사각형 탱크에 물 185mm로 채워줘..."
2. 수정하지 않는 경우, 논문에서 "프롬프트는 독립적으로 실행되며, 대화 문맥 없는 단일 입력"임을 명시

### L5. EXP-A 반복 횟수(3회)의 통계적 한계

**문제**: 3회 반복은 통계적 유의성을 주장하기에 불충분하다. 표준편차, 신뢰구간 등의 통계 지표를 제시하기 어렵다.

**비교**: NL2FOAM은 single-attempt(1회)으로 first-attempt success rate를 보고. 본 연구의 3회는 NL2FOAM보다 robust하지만 여전히 적다.

**대응**:
1. 3회 반복의 목적을 "통계적 유의성"이 아닌 **"LLM non-determinism에 의한 변동 관찰"** 로 명시
2. 보고 형식: "3/3 PASS", "2/3 PASS", "0/3 PASS" — 비율이 아닌 raw count로 투명 보고
3. Discussion에서 "소규모 반복의 한계"를 인정하고, 향후 연구로 "대규모 반복(N≥30) + bootstrap 분석" 제안

### L6. EXP-B 반복 없음 (1회)

**문제**: EXP-B는 조건당 1회만 실행. LLM의 비결정성을 감안하면 1회 결과는 재현성이 보장되지 않는다.

**정당화**: Ablation의 목적이 "정확한 성공률 측정"이 아니라 "컴포넌트 제거 시 질적 실패 모드 관찰"이므로, 1회로도 "LLM이 유효한 XML을 생성하지 못한다"는 질적 결론은 도출 가능.

**대응**:
1. EXP-B를 **정성적 ablation** (qualitative ablation)으로 위치 설정
2. 핵심 조건(B2, B4)은 실패가 결정적(XML 문법 자체가 불가)이므로 1회로 충분
3. B1(−DomainPrompt)은 확률적 결과이므로, 가능하면 3회 반복 추가 권장

### L7. 논문뼈대와 EXPERIMENT_PLAN 간 Figure/Table 번호 불일치

**문제**: EXPERIMENT_PLAN은 Table 2-6, Fig 1-7을 정의하지만, 논문뼈대.tex는 Table 1-4, Fig 1-4만 있다. 특히 **EXP-B ablation 결과(Table 4, Fig 4)가 논문뼈대에 반영되지 않았다**.

**대응**:
1. 논문뼈대.tex Section 4에 ablation 서브섹션 추가 (또는 별도 Section으로 분리)
2. Figure/Table 번호를 EXPERIMENT_PLAN 기준으로 통일
3. 논문 작성 단계(Phase 4)에서 번호 최종 확정

### L8. "8/10 PASS" 사전 결론 리스크

**문제**: 논문뼈대와 Contribution Statement에 이미 "8/10 GPU PASS"가 기재되어 있으나, EXP-A는 아직 미실행이다. 실제 결과가 다를 경우 전체 서사가 붕괴한다.

**대응**:
1. 논문뼈대의 "8/10"은 **v2 수동 테스트 기반 예측**임을 명시
2. EXP-A 실행 후 실제 결과로 교체 — 결과가 6/10이든 10/10이든 정직 보고
3. 서사 유연성 확보: "N/10 scenarios achieved full pipeline success" 형식으로 N을 변수로 유지
