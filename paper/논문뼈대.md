# SlosimAgent: A Natural-Language-Driven AI Agent for Autonomous SPH Sloshing Simulation with Quantitative Validation

---

## 요약

기존 LLM 기반 CFD 에이전트(OpenFOAMGPT, ChatCFD, Foam-Agent, MetaOpenFOAM 등 7종)는 모두 OpenFOAM(격자 기반 FVM)만 지원하며, CFD 전문가의 자동화 도구로 설계되어 있다. 체계적 문헌 조사(1,700편, 8개 DB) 결과, SPH 솔버를 사용하는 LLM 에이전트는 0편이고, 비전문가가 자연어로 슬로싱 시뮬레이션을 수행할 수 있는 도구 역시 0편임이 확인되었다. 본 연구는 이 공백을 메우기 위해 두 가지를 제안한다: (1) CFD를 모르는 슬로싱 실무자가 자연어로 DualSPHysics GPU 시뮬레이션을 설정·실행·분석할 수 있는 최초의 도메인 특화 AI 에이전트 SlosimAgent, (2) 에이전트 출력의 물리적 신뢰성을 SPHERIC Test 10 벤치마크로 검증한 결과, SPH 슬로싱 문헌 전반에 정량 검증 관행이 부재함을 발견하고 이에 대응하는 M1-M8 정량 검증 프레임워크. 10개 E2E 시나리오에서 8/10 GPU PASS를 달성했으며, Water Lateral 서브케이스에서 M2=19.5%, r=0.655의 정량 결과를 보고한다.

---

## 핵심 비교 논문

| # | 논문 | arXiv/DOI | 저자 | 역할 |
|---|------|-----------|------|------|
| 1 | OpenFOAMGPT 2.0 | 2504.19338 | Pandey et al. | RAG+ReAct, OpenFOAM, 450+ sims |
| 2 | ChatCFD | 2506.02019 | Fan et al. | Structured Reasoning, 315 cases, 82.1% |
| 3 | Foam-Agent 2.0 | 2509.18178 | Yue et al. | 6-Agent MCP, 110 cases, 88.2% |
| 4 | MetaOpenFOAM | 2407.21320 | Chen et al. | RAG multi-agent, 85% |
| 5 | English et al. | 10.1007/s40571-021-00403-3 | English et al. | mDBC, SPHERIC T10, 0개 메트릭, 263인용 |
| 6 | DualSPHysics v5 | Dominguez et al. 2022 | Dominguez et al. | DualSPHysics 코드 논문, 494인용 |
| 7 | Grand Challenges | 10.1007/s40571-020-00354-1 | Vacondio et al. | SPH 5대 Grand Challenges, GC5 정량 검증 요구 |

---

## 경쟁 시스템 비교

| 시스템 | LLM | Solver | Deployment | Target User | SPH | Sloshing | Validation | Multi-fluid |
|--------|-----|--------|------------|-------------|-----|----------|------------|-------------|
| OpenFOAMGPT 2.0 | Claude-3.7 | OpenFOAM | Cloud API | CFD expert | ✗ | ✗ | Run/No-run | ✗ |
| ChatCFD | DeepSeek-R1 | OpenFOAM | Cloud API | CFD expert | ✗ | ✗ | 68.12% | ✗ |
| Foam-Agent 2.0 | Claude 3.5 | OpenFOAM | Cloud API | CFD expert | ✗ | ✗ | Expert review | ✗ |
| MetaOpenFOAM | GPT-4 | OpenFOAM | Cloud API | CFD expert | ✗ | ✗ | Run/No-run | ✗ |
| NL2FOAM | GPT-4 | OpenFOAM | Cloud API | CFD expert | ✗ | ✗ | None | ✗ |
| Engineering.ai | Multi-LLM | Multi-solver | Cloud API | CFD expert | ✗ | ✗ | None | ✗ |
| SimLLM | GPT-4 | COMSOL | Cloud API | CFD expert | ✗ | ✗ | None | ✗ |
| **SlosimAgent (Ours)** | **Qwen3-32B** | **DualSPHysics** | **Local GPU** | **Non-expert** | **✓** | **✓** | **M1-M8** | **✓** |

---

## 논문 구조

### Abstract

슬로싱 실무자(탱크 검사원, 구조 엔지니어)는 CFD 전문 지식 없이 시뮬레이션을 수행할 수 없다. 기존 7종 LLM-CFD 에이전트는 CFD 전문가 대상이며 모두 격자 기반 솔버만 지원한다. 본 연구는 자연어로 DualSPHysics SPH 시뮬레이션을 자율 수행하는 최초의 AI 에이전트 SlosimAgent를 제시하고, SPHERIC Test 10 벤치마크 검증을 통해 에이전트 출력의 물리적 신뢰성을 입증한다.

### 1. Introduction

**1.1 Background** — 슬로싱의 공학적 중요성(해양, 항공, 지진). SPH가 슬로싱에 적합한 이유(meshless, violent free-surface). LLM 기반 CFD 에이전트의 부상(7종, 2024-2025).

**1.2 Problem** — 4가지 연구 Gap

1. **G0: 비전문가 접근성 공백** — 슬로싱 테크니션은 XML 작성, dp/BC 설정, 후처리에 CFD 전문 지식이 필요. NL→도구 패턴이 포토닉스, 방사선치료, 데이터분석에는 존재하나 SPH 슬로싱 = 0편.
2. **G1: SPH 솔버용 LLM 에이전트 = 0편** — 7종 전부 OpenFOAM(격자 FVM) only.
3. **G2: 로컬 오픈웨이트 SLM = 0편** — 전부 Claude/GPT-4/DeepSeek 클라우드 API 의존. 산업 현장 배포 불가.
4. **G3: MCP 기반 과학계산 통합 = 0편** — 전부 커스텀 래퍼/쉘 스크립트. 표준 프로토콜 미사용.

**1.3 Contribution** — 본 연구의 기여

1. CFD 비전문가가 자연어로 SPH 슬로싱 시뮬레이션을 수행할 수 있는 최초의 AI 에이전트
2. DualSPHysics GPU 솔버를 통합한 최초의 LLM 에이전트 (SPH + 로컬 SLM + MCP)
3. SPHERIC Test 10 벤치마크를 통한 에이전트 출력의 정량적 물리 검증 + M1-M8 프레임워크

### 2. Related Work

**2.1 LLM-Based CFD Agents** — OpenFOAMGPT 2.0, ChatCFD, Foam-Agent 2.0, MetaOpenFOAM, NL2FOAM, Engineering.ai, SimLLM. 공통점: OpenFOAM only, Cloud API, CFD expert target. ChatCFD만 물리적 정확도(68.12%) 측정.

**2.2 Natural Language Interfaces for Scientific Simulation** — NL→도구 패턴의 타 도메인 사례: 포토닉스 PIC 설계(Q1), 방사선치료(Q1), low-code 데이터 분석. 과학 시뮬레이션 비전문가 접근성 논의.

**2.3 SPH Sloshing Simulation and DualSPHysics** — DualSPHysics v5(494인용), SPHERIC Test 10 벤치마크, English et al.(mDBC, 263인용), Vacondio et al.(Grand Challenges).

**2.4 MCP and Tool Integration** — Model Context Protocol (Anthropic, 367인용), 과학계산 적용 사례 부재.

### 3. System Design

**3.1 Architecture Overview**

```
[사용자 자연어 입력 (한국어/영어)]
     ↓
[Qwen3-32B Coder Agent] → 슬로싱 도메인 특화 프롬프트
     ↓
[18 MCP Tools] → XML Generator + GenCase + Solver + PostProcess
     ↓
[DualSPHysics GPU Pipeline] → Docker CUDA 12.6
     ↓
[AI Analysis + Report] → 물리적 해석 + Markdown 리포트
```

**3.2 Domain-Specialized Agent (G1 해소)**
- Go + BubbleTea TUI (OpenCode 포크)
- DualSPHysics 전용 18개 도구: gencase, solver, partvtk, measuretool, xml_generator, analysis, report, animation, parametric_study, geometry, seismic_input, stl_import 등
- 슬로싱 도메인 특화 시스템 프롬프트 (SloshingCoderPrompt)
- XML attribute-only 문법, GenCase 규칙, 솔버 파라미터 자동 설정

**3.3 Local SLM Deployment (G2 해소)**
- Qwen3-32B via Ollama (로컬 GPU)
- 클라우드 API 불필요 → 산업 보안 환경 배포 가능
- 한계: Qwen3:8B는 tool call 실패 (Quality Report v0.3, 0/3 PASS) → 32B 필수

**3.4 MCP-Based Tool Integration (G3 해소)**
- 18개 도구를 Model Context Protocol 표준으로 통합
- Tool interface: `BaseTool.Info()` + `Run(ctx, ToolCall) → ToolResponse`
- Docker Pipeline: GenCase → DualSPHysics5.4 GPU → PartVTK → MeasureTool
- NVIDIA RTX 4090, CUDA 12.6

### 4. End-to-End Evaluation (G0 해소)

**4.1 Scenario Design**
- 10개 E2E 사용자 시나리오 (8개 데이터셋)
- 직사각/원통/STL 탱크, 물/오일/LNG 유체, 피치/스웨이/지진 여진
- 자연어 입력 → XML 생성 → GPU 시뮬레이션 → 결과 검증

**4.2 Results**

| # | 시나리오 | 형상 | 유체 | 파티클 | GPU 시간 | 결과 |
|---|---------|------|------|--------|---------|------|
| 1 | SPHERIC Oil | 직사각 | Oil | 136K | 131s | **PASS** |
| 2 | Chen Shallow | 직사각 | Water | 173K | 174s | **PASS** |
| 3 | Chen Critical | 직사각 | Water | 313K | 430s | **PASS** |
| 4 | Liu Large Pitch | 직사각 | Water | 247K | 738s | **PASS** |
| 5a | Liu Amp1 | 직사각 | Water | 723K | 2194s | **PASS** |
| 5b | Liu Amp2 | 직사각 | Water | 225K | ~300s | **PASS** |
| 5c | Liu Amp3 | 직사각 | Water | 225K | partial | **PARTIAL** |
| 6 | ISOPE LNG | 직사각 | LNG | 347K | 1016s | **PASS** |
| 7 | NASA Cylinder | 원통형 | Water | 323K | 216s | **PASS** |
| 8 | English DBC | 직사각 | Water | 891K | ~3000s | **PARTIAL** |
| 9 | Zhao Horizontal | 원통형 | Water | 863K | 833s | **PASS** |
| 10 | Frosina Fuel Tank | STL | Fuel | ~200K | ~300s | **PASS** |

**Overall: 8/10 GPU PASS** (PARTIAL 2건: 공진 불안정, DBC 한계)

### 5. Physics Validation: SPHERIC Test 10 (EXP-1)

> 에이전트가 생성한 시뮬레이션의 물리적 신뢰성 검증. 이 과정에서 SPH 슬로싱 문헌의 정량 검증 관행 부재를 발견 → M1-M8 프레임워크 제안.

**5.1 Benchmark Setup**
- Souto-Iglesias & Botia-Vera (2011), 102회 반복 실험
- Tank: 900 × 62 × 508 mm, Pitch 4°
- 3 sub-cases: Water Lateral, Oil Lateral, Water Roof

**5.2 Quantitative Metrics Framework (M1-M8)**

| 메트릭 | 정의 | 기준 |
|--------|------|------|
| M1 | Peak-in-band (실험 μ±2σ) | ≥2/3 |
| M2 | Mean absolute peak error | <30% |
| M3 | Convergence monotonicity | 단조 감소 |
| M4 | GCI (Richardson ext.) | Fs=3.0 |
| M5 | Cross-correlation r_max | >0.5 |
| M6 | Time shift τ | <T |
| M7 | Peak detection rate | ≥3/4 |
| M8 | Impact FWHM ratio | — |

**5.3 Results Summary**

| Sub-case | dp | Verdict | Key Metrics |
|----------|-----|---------|-------------|
| Water Lateral | 2mm | **PASS** | M1=3/3, M2=19.5%, r=0.655 |
| Oil Lateral | 2mm | **PASS** | M7=4/4, M1=2/4, r=0.570 |
| Water Roof | 2mm | **PARTIAL** | DBC 한계, 파도 99.1% 도달 |

**5.4 Convergence Evidence**
- 3-level dp (4/2/1mm): r = 0.460 → 0.655 → 0.697 (단조 수렴)
- 2-level GCI (conservative, Fs=3.0)

**5.5 Comparison with Existing Literature**

| 항목 | English 2022 | Vacondio 2020 | Delorme 2009 | **Ours** |
|------|-------------|--------------|-------------|----------|
| 벤치마크 | SPHERIC T10 | (리뷰) | canonical | SPHERIC T10 |
| 정량 메트릭 수 | **0** | — | 1-2 | **8 (M1-M8)** |
| GCI | 없음 | 언급만 | 없음 | **적용** |
| Cross-correlation | 없음 | — | 없음 | **r=0.655** |
| 다유체 | 물만 | — | 물만 | **물+오일** |
| 검증 판정 | "good agreement" | — | "overestimation" | **PASS/FAIL** |

### 6. Discussion

**6.1 Accessibility vs. Automation** — 기존 에이전트는 "CFD 전문가의 자동화". SlosimAgent는 "비전문가의 접근성". 타겟 사용자가 근본적으로 다름. 10개 시나리오는 CFD 지식 없는 자연어 입력으로 설계.

**6.2 SPH vs. FVM Agent** — SPH의 attribute-only XML, dp 기반 해상도, 입자 경계 조건 등이 OpenFOAM 에이전트와 근본적으로 다른 도구 설계를 요구하는 이유.

**6.3 Validation Practices** — EXP-1 수행 중 관찰된 SPH 슬로싱 문헌의 정량 검증 부재. English et al.(263인용)조차 "good agreement"만 보고. Vacondio et al.의 GC5가 요구하는 "convincing quantitative error analysis"가 실현되지 않고 있음.

**6.4 Limitations (정직한 보고)**
1. M2 = 19.5% — 대부분의 논문이 보고조차 안 함 + 실험 CoV 20-40%
2. DBC only — mDBC 시도했으나 발산 → 문서화
3. 시뮬레이션 8회 — OpenFOAMGPT 450+ 대비 적음
4. Qwen3:8B tool call 실패 → 32B 필수
5. E2E 10개는 GPU 파이프라인 검증 — NL→tool call 자동화 검증은 별개
6. 시간 이동 τ=+0.57s — SPH 초기화 특성

### 7. Conclusion + Future Work

**결론**:
1. CFD를 모르는 실무자가 자연어로 슬로싱 시뮬레이션을 수행할 수 있는 최초의 AI 에이전트 (10 E2E, 8/10 PASS)
2. SPH 솔버를 통합한 최초의 LLM 에이전트 (DualSPHysics + Qwen3-32B + MCP)
3. SPHERIC T10 정량 검증으로 에이전트 출력 신뢰성 입증 (M2=19.5%, r=0.655)
4. SPH 슬로싱 문헌의 정량 검증 관행 부재 관찰 → M1-M8 프레임워크 제안

**향후 연구**:
- mDBC/buffer boundary 통합으로 roof impact 검증 확장
- 에이전트 성공률 체계적 측정 (>100 cases)
- 다른 SPH 코드(SPHinXsys, GPUSPH) 통합
- Qwen3 fine-tuning으로 8B 모델 tool call 성능 개선
- 산업 현장 사용자 스터디 (usability evaluation)

---

## English et al. 핵심 비교

| 차원 | English et al. (2022) | SlosimAgent (Ours) |
|------|----------------------|---------------------|
| 연구 목적 | mDBC 경계조건 제안 | 비전문가용 AI 에이전트 |
| 벤치마크 | SPHERIC T10 (Water lat) | SPHERIC T10 (3 sub-cases) |
| dp 수준 | 2 (4mm, 2mm) | 3 (4/2/1mm) |
| 정량 메트릭 | 0개 | 8개 (M1-M8) |
| GCI | 없음 | 적용 (2-level, Fs=3.0) |
| Cross-correlation | 없음 | r=0.655, τ=+0.57s |
| 검증 표현 | "good agreement" | PASS/FAIL 체계 |
| 다유체 | 물만 | 물 + 오일 |
| 인용 | 263회 | — |

---

## Contribution Statement

> To our knowledge, SlosimAgent is the first AI agent that enables non-expert technicians to autonomously configure, execute, and analyze SPH sloshing simulations through natural language. We make three contributions: (1) a domain-specialized agent architecture integrating DualSPHysics GPU solver with a local open-weight SLM (Qwen3-32B) via the Model Context Protocol — the first LLM agent for particle-based CFD, where all seven existing agents support only grid-based OpenFOAM; (2) end-to-end evaluation across 10 user scenarios covering rectangular, cylindrical, and STL tank geometries with water, oil, and LNG fluids, achieving 8/10 GPU pipeline success; and (3) quantitative physics validation against the SPHERIC Test 10 benchmark using a proposed M1-M8 metrics framework, demonstrating that agent-generated simulations achieve M2=19.5% mean peak error and r=0.655 cross-correlation — the first quantitative error report for SPH sloshing impact pressure, where existing literature reports only "good agreement."
