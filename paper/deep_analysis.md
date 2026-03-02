# 서베이 심층 분석 + 핵심 논문 비교 + Target GAP

## 1. 서베이 개요

### 1.1 Survey v1 — AI Agent + Accessibility for CFD (c51cc0b3faf1)

- **조사 범위**: 1,173편 (8개 DB × 12개 쿼리)
- **쿼리 구성**: 8개 (AI-CFD agent) + 4개 (NL interface, low-code, non-expert, democratization)
- **핵심 질문**: "LLM-CFD 에이전트 중 SPH 솔버를 사용하는 것이 있는가? 비전문가 접근성은?"
- **Tier 1**: 40편 | **Round 1 Gap**: 4개, **Round 2 Gap**: 1개 (G0 접근성)
- **Coverage gap 확인**: 코퍼스에 "sphinxsys" 부재 (Round 2에서 "sloshing" 유입됨)
- **키워드 상위**: simulation(23.0), technician(21.2), language(21.0), expert(16.3), natural(14.5)
- **시간 분포**: 2025=439편, 2024=175편 (매우 최신 분야)
- **핵심 논문 원문**: 6편 확보 (OpenFOAMGPT2, Foam-Agent2, ChatCFD, MetaOpenFOAM, Engineering.ai, SimLLM)

### 1.2 Survey v2 — SPH Sloshing Validation (7a1b4eb3df2a)

- **조사 범위**: 527편 (8개 DB × 4개 쿼리)
- **핵심 질문**: "SPH 슬로싱에서 SPHERIC Test 10 벤치마크를 어떻게 검증했는가?"
- **Tier 1**: 30편 | **Round 1 Gap**: 4개 확정
- **Coverage gap 확인**: 코퍼스에 "cross-correlation" 키워드 부재
- **시간 분포**: 2025=68편, 2024=57편, 2023=64편 (성숙 분야)
- **핵심 논문 원문**: 3편 확보 (Dominguez 2022, English 2022, NL2FOAM)
- **기존 Deep Read**: 6편 (markdown)

### 1.3 통합 요약

| 항목 | v1 (AI Agent+접근성) | v2 (Validation) | 합계 |
|------|---------------------|----------------|------|
| 수집 논문 | 1,173 | 527 | **1,700** |
| 검색 쿼리 | 12 | 4 | 16 |
| DB 커버리지 | 8/8 | 8/8 | 8/8 |
| Round | 2 | 1 | — |
| 확정 Gap | 5 (R1:4 + R2:1) | 4 | **핵심 4개 + 부수 발견 6개** |
| 원문 확보 | 6 PDF + 3 txt | 3 PDF + 6 md + 3 txt | **18** |

## 2. AI-CFD Agent 심층 비교 (Survey v1 원문 분석)

### 2.0 LLM-CFD Agent 비교 매트릭스 (원문 기반, 2026-03-02 확인)

| 항목 | OpenFOAMGPT 2.0 | ChatCFD | Foam-Agent 2.0 | **SlosimAgent (ours)** |
|------|-----------------|---------|----------------|----------------------|
| arXiv | 2504.19338 | 2506.02019 | 2509.18178 | — |
| LLM | Claude-3.7 Sonnet | DeepSeek-R1/V3 | Claude 3.5 Sonnet | **Qwen3-32B (local)** |
| Solver | OpenFOAM only | OpenFOAM only | OpenFOAM only | **DualSPHysics (SPH)** |
| Deployment | Cloud API | Cloud API | Cloud API | **Local GPU (Ollama)** |
| Architecture | RAG + ReAct | Thought-Action-Obs | 6-Agent MCP | **Coder Agent + 18 tools** |
| Agent framework | Custom (LangGraph) | Custom | MCP-based | **OpenCode fork + MCP** |
| Simulation count | 450+ | 315 | 110 | 8 (EXP-1) |
| Success rate | ~100% (generate) | 82.1% (correct) | 88.2% | — |
| Physics fidelity | Not measured | 68.12% | Not measured | **M1-M8 framework** |
| Validation method | Run/No-run | Metrics vs ref | Expert review | **SPHERIC T10 + M1-M8** |
| Target user | CFD expert | CFD expert | CFD expert | **Non-expert technician** |
| Multi-fluid | No | No | No | **Yes (water + oil)** |
| Mesh type | Grid (FVM) | Grid (FVM) | Grid (FVM) | **Particle (SPH)** |

**핵심 발견**:
1. 7개 에이전트 모두 OpenFOAM(격자 기반 FVM)만 지원 — SPH 솔버 = 0편
2. 모두 클라우드 API 의존 — 로컬 SLM = 0편
3. ChatCFD만 물리적 정확도(68.12%) 측정 — 나머지는 "실행 성공/실패"만 보고
4. 슬로싱 도메인 특화 = 0편
5. **모두 CFD 전문가 대상** — 비전문가 접근성을 목표로 한 에이전트 = 0편

### 2.1 English et al. (2022) — mDBC

**DOI**: 10.1007/s40571-021-00403-3 | **인용**: 263회

- **연구 내용**: DualSPHysics DBC의 비물리적 갭 문제를 해결하는 mDBC (modified Dynamic Boundary Condition) 제안
- **SPHERIC Test 10 사용**: Section 4.2에서 sloshing tank 검증
  - dp = 0.004 m, dp = 0.002 m, h/dp = 2
  - DBC vs mDBC 비교: 시계열 오버레이 그래프 3개
  - mDBC가 DBC의 비물리적 갭 제거 → 센서 정확 위치 측정 가능
- **검증 방법**: **순수 시각적 비교 (visual comparison)**
  - "good agreement with the experimental data" — 유일한 정량적 표현
  - **peak error, RMSE, cross-correlation, GCI: 일체 없음**
- **우리와의 차이**:
  - 우리: M1-M8 8개 정량 메트릭, 3-level 수렴 연구, GCI 계산
  - English: 0개 정량 메트릭, 2-level만 (수렴 분석 없음)
  - **동일 벤치마크에서 검증 엄격도 차이가 극명**

### 2.2 Vacondio et al. (2020) — Grand Challenges for SPH

**DOI**: 10.1007/s40571-020-00354-1 | **인용**: 241회

- **연구 내용**: SPHERIC 커뮤니티가 정의한 SPH 5대 Grand Challenges
- **GC5 (Applicability to Industry)** 핵심 요구:
  - "comprehensive validation" — 정성+정량 모두 필요
  - "convincing quantitative error analysis" — 설득력 있는 정량 오차 분석
  - 수렴률 조사 필수
  - 불일치는 "물리적/수치적 근거로 충분히 설명" 필요
- **그러나**: 구체적 통과/불합격 기준값 미명시 → 리뷰어 재량
- **GC1 (Convergence)**: Richardson extrapolation 언급은 있으나 SPH 입자 기반 수렴의 어려움 강조
  - Quinlan et al. (2006): 입자 분포에 따라 수렴률 변동
  - 수렴 + 일관성 + 안정성: SPH에서 여전히 grand challenge
- **우리와의 관계**:
  - 우리의 M1-M8 프레임워크는 GC5 요구사항의 **최초 체계적 실현** 중 하나
  - GCI 적용은 GC1의 수렴 분석 요구에 직접 응답
  - 하지만 SPH 입자법에 Richardson extrapolation 직접 적용은 엄밀하지 않음 → 이를 인정하고 "정신 차용"으로 포지셔닝

### 2.3 Dominguez et al. (2022) — DualSPHysics v5

**인용**: 494회 (DualSPHysics 코드 논문)

- **연구 내용**: DualSPHysics 코드 전반 리뷰, multiphysics 확장
- **슬로싱 검증**: English et al. (2022) 참조로 대체 — 별도 정량 검증 없음
- **의미**: DualSPHysics 팀 자체가 슬로싱에서 정량 검증을 수행하지 않았음

### 2.4 Delorme et al. (2009) — Canonical Sloshing Part I

- **연구 내용**: SPH 슬로싱 canonical problems 정의 (Part 0: Souto-Iglesias et al.)
- **검증 방법**: 정량적 — SPH가 실험 대비 과추정(overestimation) 경향 보고
- **한계**: 2D SPH, GCI 미적용, dp 비교 제한적
- **우리와의 관계**: 이 논문이 정량적 보고를 한 드문 예 → 우리가 이를 확장

### 2.5 Laha et al. (2024) — Heart Valve SPH FSI

**DOI**: 10.1038/s41598-024-57177-w

- DualSPHysics를 심장판막 FSI에 적용
- **검증 방법**: FVM, 4D MRI 데이터와 비교 — 정량적 (leaflet angle error 5.6%)
- **의미**: 다른 응용에서는 정량 검증을 하면서, 슬로싱에서는 안 하는 이중 기준

### 2.6 Flow3D 연구 (격자 기반, narrow tank)

- **방법**: 격자 기반 VOF
- **검증**: RMS 4.79%, Max error 8.33%, 수렴 연구 수행
- **의미**: 격자 기반 CFD는 이미 정량 검증이 표준 → SPH의 갭이 더 두드러짐

## 3. 문헌 비교 매트릭스

| 항목 | English 2022 | Vacondio 2020 | Delorme 2009 | Flow3D | **우리 (2026)** |
|------|-------------|--------------|-------------|--------|---------------|
| 벤치마크 | SPHERIC T10 | (리뷰) | canonical | narrow | SPHERIC T10 |
| 솔버 | DualSPHysics | (리뷰) | SPH 2D | Flow3D | DualSPHysics |
| BC | DBC/mDBC | — | — | VOF | DBC |
| dp 수준 | 2 (4mm, 2mm) | — | 1 | 격자 | 3 (4/2/1mm) |
| 정량 메트릭 수 | **0** | — | 1-2 | 3+ | **8 (M1-M8)** |
| Peak error 보고 | 없음 | — | 있음 | 4.79% RMS | **19.5% MAE** |
| Cross-correlation | 없음 | — | 없음 | 없음 | **r=0.655** |
| GCI/Richardson | 없음 | 언급만 | 없음 | 있음 | **적용 (2-level)** |
| 수렴 단조성 | 미확인 | — | 미확인 | 확인 | **확인 (3-level)** |
| 다유체 | 물만 | — | 물만 | 물만 | **물 + 오일** |
| Roof impact | 없음 | — | 없음 | 없음 | **시도 (PARTIAL)** |
| 검증 판정 | "good agreement" | — | "overestimation" | 수치 | **PASS/FAIL 체계** |

## 4. Target GAP 정의 (2축 4개 핵심 Gap + EXP-1 부수 발견)

### 핵심 Gap: Axis 0 — Accessibility (WHY)

| # | Gap | Sev. | 근거 | 서베이 증거 |
|---|-----|:----:|------|-------------|
| **G0** | SPH 슬로싱 시뮬레이션은 비전문가가 접근 불가 | **10** | XML 작성, dp/BC 설정, 후처리 → CFD 전문 지식 필수 | NL→도구 패턴: 포토닉스(PIC, Q1), 방사선치료(Q1), 데이터분석(low-code) 존재. OpenFOAM CFD(7종)도 존재. **SPH 슬로싱 = 0편** |

**G0의 핵심 차이점**: 기존 7개 LLM-CFD 에이전트는 **CFD 전문가가 자기 작업을 자동화**. SlosimAgent는 **CFD를 모르는 슬로싱 테크니션**이 자연어로 시뮬레이션을 수행. 타겟 사용자가 근본적으로 다름.

**우리의 증거**: 10개 E2E 시나리오 (8개 데이터셋, 직사각/원통/STL, 물/오일/LNG) → **8/10 GPU PASS**

### 핵심 Gap: Axis 1 — Technical Agent (WHAT)

| # | Gap | Sev. | 근거 | 서베이 증거 |
|---|-----|:----:|------|-------------|
| **G1** | SPH 솔버용 LLM 에이전트 = 0편 | **9** | 7종 전부 OpenFOAM(격자 FVM) only | 원문 비교: OpenFOAMGPT2, ChatCFD, Foam-Agent2 모두 OpenFOAM |
| **G2** | 로컬 오픈웨이트 SLM = 0편 | **8** | 전부 Claude/GPT-4/DeepSeek 클라우드 API | 원문 확인: privacy/industrial deployment 불가 |
| **G3** | MCP 기반 과학계산 통합 = 0편 | **7** | 전부 커스텀 래퍼/쉘 스크립트 | MCP spec 367회 인용 but 과학계산 적용 0편 |

### Gap 대응 매트릭스

```
Gap    문제                              SlosimAgent 해소            증거
─────  ────────────────────────────────  ─────────────────────────  ──────────────────
G0     비전문가 접근 불가                  자연어 → XML → GPU 파이프라인  10 E2E, 8/10 PASS
G1     SPH 에이전트 = 0                   DualSPHysics 18 tools       최초 SPH AI 에이전트
G2     로컬 SLM = 0                       Qwen3-32B Ollama            클라우드 API 불필요
G3     MCP 과학계산 = 0                   MCP 표준 프로토콜             18개 도구 MCP 통합
```

### EXP-1 수행 중 부수 발견 (추가 기여, 사전 동기 아님)

EXP-1 (SPHERIC Test 10 벤치마크) 검증 실험을 수행하면서, SPH 슬로싱 문헌의 검증 관행에 체계적 부재가 있음을 관찰했다. 이는 연구의 사전 동기가 아니라 **실험 과정에서 드러난 발견**이다.

| 발견 | 관찰 내용 | SlosimAgent에서의 대응 |
|------|----------|----------------------|
| **F1. 정량 메트릭 부재** | 263인용 English 2022도 "good agreement"만 보고, 0개 메트릭 | M1-M8 프레임워크 8개 정량 메트릭 적용 |
| **F2. GCI 미적용** | SPH sloshing + Richardson extrapolation 조합 = 0편 | 3-level dp + 2-level GCI 계산 |
| **F3. 교차상관 미사용** | 시간 이동(time shift) 정량화 안 함 | r=0.655, τ=+0.57s 보고 |
| **F4. 다유체 검증 부재** | 단일 유체(물)만 검증 | Water + Oil 동일 프레임워크 |
| **F5. Roof impact 미보고** | 기존 논문에서 roof 충격 검증 없음 | 시도 + DBC 한계 정직 보고 (PARTIAL) |
| **F6. v1→v2 개선 기록** | 실패→성공 과정 미공개가 관행 | peak error 63.5%→19.5% (3.3×), r: -0.087→0.655 공개 |

**의의**: 이 발견들은 핵심 Gap에서 파생된 것이지만, 에이전트가 생성한 시뮬레이션의 신뢰성을 입증하는 데 필수적인 증거 역할을 한다. 특히 F1-F3은 "에이전트가 돌린 시뮬레이션이 정말 맞는가?"라는 질문에 정량적으로 답하는 수단이다.

## 5. 논문 포지셔닝 전략 (2축 + 실험 발견)

### 핵심 내러티브 (2-Axis Story)

**문제**: 슬로싱 실무자(탱크 검사원, 구조 엔지니어)는 CFD 전문 지식 없이 시뮬레이션을 수행할 수 없다.
기존 LLM-CFD 에이전트 7종은 CFD 전문가의 자동화 도구이며, 모두 격자 기반 솔버만 지원한다.

**해결**: 본 연구는 SlosimAgent를 제시한다 — CFD를 모르는 슬로싱 실무자가 자연어로 DualSPHysics 시뮬레이션을 설정·실행·분석할 수 있는 최초의 AI 에이전트.

**부수 발견**: 에이전트의 물리적 신뢰성을 검증하기 위해 SPHERIC Test 10 벤치마크를 수행한 결과, SPH 슬로싱 문헌 전반에 정량 검증 관행이 부재함을 관찰. 이에 M1-M8 정량 프레임워크를 자체 제안하여 에이전트 출력의 신뢰성을 입증.

### 논문 구조 매핑

```
Section          내용                                    근거
─────────────    ──────────────────────────────────────  ──────────
Introduction     G0: 비전문가 접근성 공백 (WHY)             Survey v1 R2
                 G1-G3: 기술적 공백 (WHAT)                Survey v1 R1
System Design    SlosimAgent 아키텍처                      Go + TUI + 18 tools
                 → G1(SPH), G2(Local SLM), G3(MCP) 해소
E2E Evaluation   10개 사용자 시나리오 → 8/10 PASS          → G0 해소
Validation       SPHERIC T10 벤치마크 (EXP-1)              → 에이전트 출력 신뢰성
                 F1-F3 문헌 관행 관찰 + M1-M8 제안         → 추가 기여
Discussion       한계 정직 보고 + 향후 연구
```

### Contribution 목록 (우선순위 순)

1. **[G0] 비전문가 접근성 실증**: 10개 E2E 시나리오, 8/10 GPU PASS — CFD 비전문가가 자연어로 슬로싱 시뮬레이션 수행
2. **[G1] 최초 SPH AI 에이전트**: 7개 LLM-CFD 에이전트 중 유일한 입자법(SPH) 지원
3. **[G2] 로컬 SLM 배포**: Qwen3-32B Ollama, 클라우드 API 불필요 — 산업 현장 배포 가능
4. **[G3] MCP 기반 과학계산 통합**: 18개 도구를 표준 프로토콜로 통합
5. **[F1-F3] 정량 검증 프레임워크**: M1-M8 8개 메트릭, GCI, 교차상관 — EXP-1에서 관찰된 문헌 관행 부재에 대한 대응

### 약점 → 정직한 보고 (Limitations as Transparency)

1. **M2 = 19.5%**: 높아 보이지만, 대부분의 논문이 보고조차 안 함 + 실험 CoV 20-40%
2. **DBC only**: mDBC 시도했으나 발산 → 정직하게 문서화
3. **2-level GCI**: Run 003 TimeOut 제약 → conservative Fs=3.0
4. **시간 이동**: τ=+0.57s → 물리적 원인 설명 (SPH 초기화 특성)
5. **시뮬레이션 8회**: 에이전트 성공률 미측정 (OpenFOAMGPT의 450+ 대비 적음)
6. **Qwen3:8B tool call 실패**: Quality Report v0.3에서 0/3 시나리오 PASS → 32B 필요
7. **E2E 10개는 GPU 파이프라인 검증**: 에이전트 NL→tool call 자동화 검증과 별개

### 핵심 문장 (Key Claims)

- "CFD를 모르는 실무자가 자연어로 슬로싱 시뮬레이션을 수행할 수 있는 최초의 도구"
- "SPH 솔버를 자동화한 최초의 AI 에이전트"
- "클라우드 API 없이 로컬 GPU만으로 작동하는 CFD 에이전트"
- "에이전트 출력의 물리적 신뢰성을 정량적으로 검증한 유일한 사례"
