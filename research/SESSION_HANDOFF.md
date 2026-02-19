# Session Handoff — 2026-02-19 (Updated: Phase 5 Complete)

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

### Phase 3a: EXP-2 NL→XML 생성 (완료, Qwen3 8B)
- [x] 20 시나리오 × FULL 프롬프트, Ollama tool calling API
- [x] L1 Basic: **96%** 정확도 (4/4 tool called)
- [x] L2 Domain: **42%** (3/4), L3 Paper: **15%** (3/4)
- [x] L4 Complex: **50%** (3/4), L5 Edge: **25%** (2/4)
- [x] **전체: 46% 정확도, 15/20 tool called, 14/20 physical valid**
- [x] Table 3, results.json 생성

### Phase 3b: EXP-4 도메인 프롬프트 어블레이션 (완료, Qwen3 8B)
- [x] 10 시나리오 × 4 ablation conditions = 40 runs
- [x] FULL: 46%, NO-DOMAIN: 44%, **NO-RULES: 55%**, GENERIC: 39%
- [x] Table 5, Fig 4 (ablation bar chart), results.json 생성
- [x] **주요 발견**: NO-RULES(규칙 제거, 도메인 지식 유지)가 최고 성능!

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

### Qwen3 Thinking Mode 이슈
- Qwen3 32B/8B 모델 모두 thinking 모드 기본 활성화
- 1000+ thinking 토큰 생성 후 tool call → 32B에서 호출당 10분+
- `/no_think` 지시가 Ollama 환경에서 비효과적
- 8B 모델은 GPU에 완전 로딩되어 124 tok/s → 호출당 5-15초
- **32B 모델**: 14GB/20GB GPU → 나머지 CPU offload → 호출당 10분+ (백그라운드 실행 중)

### EXP-2/4 주요 분석
- **L1 (명시적 파라미터)**: 8B도 96% 정확도 — 탱크 치수, 주파수 등 직접 지정 시 우수
- **L3 (논문 재현)**: 15% — "Chen2018", "SPHERIC Test 10" 언급만으로 파라미터 추론 불가
- **Tool calling 실패 패턴**: 긴 system prompt + 모호한 입력 → 텍스트 응답 (tool call 안 함)
- **어블레이션 역설**: NO-RULES(55%) > FULL(46%) — 과도한 규칙이 오히려 tool calling 방해

---

### Phase 4: 통합 분석 (완료, Agent Team)
- [x] EXP-1 SPHERIC 정량 분석: **"Pearson r > 0.9" 주장 기각** → peak-within-±2σ (7/7=100%)
  - 시계열 r은 impact pressure에 부적절 (CV=20-36%, 실험 반복 간에도 불가)
  - SPHERIC/ISOPE 표준 peak-in-band metric으로 전환
- [x] GAP-RQ-EXP 통합 결과표: 5 GAP × 5 EXP 매핑 (3 COVERED, 2 PARTIAL)
- [x] Fig 1 아키텍처 다이어그램 (matplotlib)
- [x] 8B vs 32B 비교 분석: 32B를 primary model로 결정

### Phase 5: 논문 작성 (완료, Agent Team 5명 병렬)
- [x] Abstract (305 words)
- [x] Introduction (1,360 words): 슬로싱 중심 내러티브, 실사고 6건 테이블
- [x] Related Work (1,355 words): 11개 경쟁 시스템 Table 1
- [x] System Design (1,238 words): 14 tools, SloshingCoderPrompt, IsError 패턴
- [x] Experiments (3,007 words): EXP-1~5 전체 결과 + 정직한 한계 기술
- [x] Discussion (1,000 words): 8개 한계점 (DBC, 47% accuracy, 8B inversion)
- [x] Conclusion (455 words): 5 contributions + future work
- [x] LaTeX shell + BibTeX 35+ entries
- **총 8,720 words**

---

## 다음 작업 (Phase 6: 리비전 및 제출)

### 선택적 보강
- [ ] LaTeX 본문 완전 변환 (현재 shell + markdown 참조 형태)
- [ ] Figure 해상도/스타일 통일 (publication-ready)
- [ ] Chen2018 Figure digitization (EXP-3 정량 비교)
- [ ] mDBC 지원 추가 후 Oil 재검증
- [ ] 타겟 학회/저널 선정 및 포맷 맞춤

### 제출 준비
- [ ] Camera-ready LaTeX 완성
- [ ] Supplementary materials (코드, 데이터)
- [ ] Cover letter

---

## 핵심 파일 위치

| 파일 | 용도 |
|------|------|
| **논문 초안** | |
| `research/paper_draft/sloshagent_paper.md` | 통합 논문 (8,720 words) |
| `research/paper_draft/sloshagent_paper.tex` | LaTeX shell |
| `research/paper_draft/00_abstract.md` | Abstract |
| `research/paper_draft/01_introduction.md` | Introduction |
| `research/paper_draft/02_related_work.md` | Related Work |
| `research/paper_draft/03_system_design.md` | System Design |
| `research/paper_draft/04_experiments.md` | Experiments |
| `research/paper_draft/05_discussion.md` | Discussion |
| `research/paper_draft/06_conclusion.md` | Conclusion |
| `research/references.bib` | BibTeX 참고문헌 (35+ entries) |
| **분석** | |
| `research/experiments/unified_results.md` | GAP-RQ-EXP 통합 결과 |
| `research/experiments/exp1_spheric/analysis_summary.md` | EXP-1 정량 분석 |
| `research/experiments/model_comparison.md` | 8B vs 32B 비교 |
| **실험 설계** | |
| `research/output/experiment_design.md` | 실험 설계서 |
| `research/output/paper_skeleton.md` | 논문 뼈대 (Phase 5 이전) |
| `research/output/gap_refinement_cycle3.md` | GAP analysis (5 gaps) |
| **실험 결과** | |
| `research/experiments/exp1_spheric/` | EXP-1 SPHERIC 결과 |
| `research/experiments/exp2_nl2xml/` | EXP-2 NL→XML 결과 |
| `research/experiments/exp3_parametric/` | EXP-3 파라메트릭 결과 |
| `research/experiments/exp4_ablation/` | EXP-4 어블레이션 결과 |
| `research/experiments/exp5_industrial/` | EXP-5 산업 PoC 결과 |

## 커밋 히스토리

```
691b631 Phase 4-5 완료 — 논문 초안 전체 (8,720 words)
2f7d0e7 32B 결과 + 8B vs 32B 비교 분석
c7612d3 Phase 3 EXP-2/EXP-4 (Qwen3 8B)
3c498a4 Phase 2 완료 세션 핸드오프
e576e26 EXP-3 Chen2018 + EXP-5 Industrial PoC
3231de6 EXP-1 SPHERIC Test 10 벤치마크
fc0ed7e 세션 핸드오프 문서
a1d0dab Phase 1 인프라 준비
092ddec SPHERIC 유체 수 오류 수정
8a1a8fa 검증 데이터셋 인덱스
ff8ba5c 실험 설계서
```

## 브랜치 구조

```
research/paper  ← 메인 연구 브랜치 (Phase 1-5 전체)
paper/draft     ← 워크트리 작업 브랜치 (논문 초안, research/paper에 머지됨)
```
