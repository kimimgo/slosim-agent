# Research-v3: 실험 계획

## 논문 주장 → 실험 매핑

```
Claim                          Experiment    Status
──────────────────────────     ──────────    ──────
C1. NL→SPH 파이프라인 작동     EXP-A         NEW
C2. 아키텍처 컴포넌트 필수성    EXP-B         NEW
C3. 물리적 신뢰성              EXP-C         v2 완료 (가져옴)
```

---

## EXP-A: NL→Tool Call 파이프라인 테스트

**목적**: 비전문가의 자연어 입력이 실제로 올바른 시뮬레이션으로 변환되는지 검증

**방법**:
1. 10개 E2E 시나리오의 자연어 프롬프트를 Qwen3-32B에 입력
2. 에이전트가 자율적으로 tool call → XML 생성 → 솔버 실행 → 후처리까지 완주하는지 측정
3. 각 단계별 성공/실패 기록

**메트릭**:
- Tool Selection Accuracy: 올바른 도구를 올바른 순서로 호출했는가
- XML Validity: 생성된 XML이 GenCase를 통과하는가
- Pipeline Completion: GenCase → Solver → PostProcess 완주율
- Parameter Fidelity: 사용자가 지정한 파라미터(탱크 크기, 유체, 주파수)가 정확히 반영되는가

**조건**:
- Qwen3-32B (primary)
- Qwen3-8B (비교, Quality Report v0.3 확장)

**시나리오**: E2E 10개 중 대표 5개 선정
- S1: SPHERIC Water Lateral (기본)
- S2: Chen Shallow Sway (수평 가진)
- S4: Liu Large Pitch (장시간)
- S7: NASA Cylinder (원통형)
- S10: Frosina Fuel Tank (STL)

**선행 조건**: 없음 (즉시 수행 가능, Ollama + Docker 준비됨)

**산출물**:
- Table: 시나리오별 단계 성공률
- Figure: 파이프라인 흐름 예시 (1개 시나리오 상세)

---

## EXP-B: Ablation Study

**목적**: 도메인 특화 아키텍처의 각 컴포넌트가 필수적인지 검증

**방법**: 동일 시나리오 3개에 대해 5가지 조건으로 실행

**조건**:
| 조건 | 도메인 프롬프트 | DualSPHysics 도구 | 후처리 도구 | XML 템플릿 |
|------|:---:|:---:|:---:|:---:|
| Full (SlosimAgent) | ✓ | ✓ | ✓ | ✓ |
| -DomainPrompt | ✗ | ✓ | ✓ | ✓ |
| -XMLTool | ✓ | ✗ (raw LLM) | ✓ | ✗ |
| -PostProcess | ✓ | ✓ | ✗ | ✓ |
| Bare LLM | ✗ | ✗ | ✗ | ✗ |

**시나리오**: 3개
- S1: SPHERIC Water Lateral (기본, 직사각, 피치)
- S7: NASA Cylinder (원통형)
- S10: Frosina Fuel Tank (STL)

**메트릭**:
- XML Validity (GenCase 통과 여부)
- Solver Completion (GPU 솔버 완주 여부)
- Physical Plausibility (결과가 물리적으로 타당한가 — 압력 범위, 파형 등)
- Error Count (LLM hallucination, XML 문법 오류 등)

**선행 조건**: EXP-A 완료 후 (Full 조건 = EXP-A 결과 재사용)

**산출물**:
- Table: 조건별 × 시나리오별 성공/실패 매트릭스
- Figure: Ablation 결과 bar chart

---

## EXP-C: SPHERIC Test 10 물리 검증 (v2에서 이관)

**목적**: 에이전트가 생성한 시뮬레이션의 물리적 정확도 정량 검증

**방법**: research-v2의 run_001~008 결과 + M1-M8 프레임워크

**v2에서 가져올 데이터**:
- run_001 (dp=4mm), run_002 (dp=2mm), run_003 (dp=1mm): Water Lateral
- run_005 (dp=2mm): Oil Lateral
- run_006, 007, 008: Water Roof
- 분석 스크립트: convergence_analysis.py, oil_roof_analysis.py, paper_figures.py
- 메트릭 JSON: metrics.json, oil_roof_metrics.json
- Figure: fig_timeseries.png, fig_convergence.png, fig_oil_lateral.png

**추가 작업 없음** — v2 결과를 정리하여 논문에 포함

**산출물**:
- Table: M1-M8 결과 요약
- Figure: 시계열 비교 (3개), 수렴 분석 (1개)

---

## 실험 순서

```
EXP-A (NL→tool call)    ← 즉시 시작 가능
  ↓
EXP-B (Ablation)        ← EXP-A Full 결과 재사용
  ↓
EXP-C (물리 검증)        ← v2에서 이관, 추가 실험 없음
  ↓
논문 작성                 ← 모든 Figure/Table 완성 후
```

## Figure 목록 (최종)

| # | 내용 | 실험 | 상태 |
|---|------|------|------|
| Fig 1 | 시스템 아키텍처 다이어그램 | — | NEW |
| Fig 2 | NL→tool call 파이프라인 예시 | EXP-A | NEW |
| Fig 3 | E2E 시나리오 결과 bar chart | EXP-A | NEW |
| Fig 4 | Ablation 결과 | EXP-B | NEW |
| Fig 5 | Water Lateral 시계열 3-resolution | EXP-C | v2 완료 |
| Fig 6 | 수렴 분석 (r vs dp) | EXP-C | v2 완료 |
| Fig 7 | Oil Lateral vs 실험 | EXP-C | v2 완료 |

## Table 목록 (최종)

| # | 내용 | 실험 |
|---|------|------|
| Table 1 | 경쟁 에이전트 비교 (7종 vs Ours) | 서베이 |
| Table 2 | M1-M8 메트릭 정의 | EXP-C |
| Table 3 | NL→tool call 단계별 성공률 | EXP-A |
| Table 4 | Ablation 매트릭스 | EXP-B |
| Table 5 | SPHERIC T10 결과 | EXP-C |
| Table 6 | 수렴 분석 (r, GCI) | EXP-C |

## 결론 (논문에서 주장할 것)

1. **SlosimAgent는 비전문가의 자연어를 SPH 시뮬레이션으로 변환할 수 있다**
   - 증거: EXP-A — N/M 시나리오 파이프라인 완주
   - 한계: 8B 모델 실패, 32B 필수

2. **도메인 특화 설계의 각 컴포넌트가 필수적이다**
   - 증거: EXP-B — ablation에서 어느 하나 빠지면 성공률 급감
   - 핵심: Bare LLM은 DualSPHysics XML 생성 불가

3. **에이전트가 생성한 시뮬레이션은 물리적으로 신뢰할 수 있다**
   - 증거: EXP-C — SPHERIC T10 Water PASS (M2=19.5%, r=0.655), Oil PASS
   - 부수 발견: SPH 슬로싱 문헌의 정량 검증 관행 부재 → M1-M8 제안
