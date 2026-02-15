# slosim-agent 전체 버전 로드맵

*Last updated: 2026-02-15*

---

## 버전 요약

| 버전 | 테마 | 핵심 목표 | 상태 |
|------|------|----------|------|
| **v0.1** | MVP: E2E 파이프라인 | 자연어 → 시뮬레이션 → 텍스트 결과 | **COMPLETE** |
| **v0.2** | 시각화 + AI 분석 + 도구 검증 | 결과를 "보고 이해하는" 기능 + 18도구 E2E | **COMPLETE** (19/19 이슈) |
| **v0.3** | TUI 고도화 + 프로덕션 안정성 | 대화형 UX, mDBC, 에이전트 E2E | PLANNED |
| **v1.0** | 파라메트릭 스터디 + 배포 | 비교 분석, 패키징, 문서 | PLANNED |
| **v1.x+** | 확장 플랫폼 | 새 유체/솔버/클라우드/API | 지속적 |

---

## 버전 간 의존성 흐름

```
v0.1 MVP ✅               v0.2 시각화+검증 ✅       v0.3 TUI+안정성           v1.0 완성
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
자연어 입력 ✅      ───→  이미지 결과 확인 ✅  ───→  대화형 Case Wizard ───→  케이스 간 비교
텍스트 리포트 ✅     ───→  Mesa 렌더링 ✅      ───→  Sim Dashboard ✅  ───→  비교 리포트
VTK/CSV 생성 ✅     ───→  pvpython+ffmpeg ✅  ───→  에이전트 통합 E2E ───→  Result Store
단일 Job ✅         ───→  실시간 모니터링 ✅  ───→  mDBC 지원        ───→  파라메트릭 스터디
직사각형+정현파 ✅   ───→  원통형+STL+Mesa ✅  ───→  L형+지진파 E2E   ───→  플러그인 템플릿
                         GPU 10시나리오 ✅
```

---

## v0.1 — MVP: E2E 최소 파이프라인 ✅ COMPLETE

**테마:** "자연어 한 줄 → 시뮬레이션 → 결과 확인"이 동작하는 최소 제품
**상세 명세:** `docs/FRD_v0.1.md`

### 포함 기능 (전부 완료)

| FRD ID | 기능명 | 상태 |
|--------|--------|------|
| AGENT-01 | Sloshing Coder 시스템 프롬프트 (Qwen3 32B) | ✅ |
| AGENT-02 | XML 케이스 자동 생성 (직사각형) | ✅ |
| TOOL-01 | GenCase Tool (Docker) | ✅ |
| TOOL-02 | Solver Tool (GPU) | ✅ |
| TOOL-03 | PartVTK Tool | ✅ |
| TOOL-04 | MeasureTool | ✅ |
| JOB-01 | Job Manager (단일) | ✅ |
| RPT-01 | Markdown 리포트 | ✅ |
| TUI-01 | 상태 바 시뮬레이션 표시 | ✅ |
| PIPE-01 | E2E 파이프라인 오케스트레이션 | ✅ |

---

## v0.2 — 시각화 + AI 분석 + 도구 E2E 검증 ✅ COMPLETE

**테마:** 시뮬레이션 결과를 "보고 이해하는" 기능 + 18개 도구 전체 검증
**완료일:** 2026-02-15
**GitHub Issues:** #1~#19 전부 종료 (alpha 13개 + beta 6개)

### 포함 기능 (전부 완료)

| ID | 기능명 | 상태 | 비고 |
|----|--------|------|------|
| VIS-01 | pvpython 후처리 (Mesa 호환) | ✅ | offscreen + Mesa backend |
| VIS-02 | 정지 이미지 스냅샷 | ✅ | frame-by-frame loop (SaveAnimation 불가) |
| VIS-03 | 시각화 항목 (압력 컬러맵) | ✅ | Press 컬러바, Points representation |
| MON-01 | 실시간 모니터링 (Run.csv) | ✅ | semicolon 파서, TimeMax 동적 읽기 |
| RPT-02 | 이미지 포함 리포트 | ✅ | |
| RPT-03 | AI 물리적 해석 | ✅ | Qwen3 32B 기반 |
| ANI-01 | 애니메이션 (MP4) | ✅ | pvpython + ffmpeg drawtext |
| STL-01 | STL 파일 입력 + watertight 검증 | ✅ | ASCII/Binary 파서 + edge-sharing |
| GEO-01 | 원통형 탱크 | ✅ | drawcylinder primitives |
| NFR-01 | 에러 감지 및 복구 | ✅ | |
| TUI-02 | Sim Dashboard (기본) | ✅ | BubbleTea 컴포넌트 |
| TUI-03 | Result Viewer (기본) | ✅ | |

### E2E GPU 검증 (10/10 시나리오)

| # | 시나리오 | 파티클 | GPU 시간 | 결과 |
|---|---------|--------|---------|------|
| 1 | SPHERIC Oil Low Fill | 136K | 131s | PASS |
| 2 | Chen Shallow Sway | 173K | 174s | PASS |
| 3 | Chen Near-Critical | 313K | 430s | PASS |
| 4 | Liu Large Pitch | 247K | 738s | PASS |
| 5 | Liu Amplitude Parametric (3개) | 225-723K | 300-2194s | PASS (5c partial) |
| 6 | ISOPE LNG Mark III | 347K | 1016s | PASS |
| 7 | NASA Cylindrical | 323K | 216s | PASS |
| 8 | English DBC | 891K | ~3000s | PARTIAL (mDBC 필요) |
| 9 | Zhao Horizontal Cyl | 863K | 833s | PASS (primitives) |
| 10 | Frosina Fuel Tank (STL) | ~200K | ~300s | PASS |

상세 결과: `docs/scenarios/e2e_dataset_scenarios.md` (`test/e2e-gpu-scenarios` 브랜치)

### 주요 기술 결정 (확정)

| 사항 | 결정 |
|------|------|
| ParaView 렌더링 | Mesa offscreen (Docker 없이 호스트 pvpython) |
| SaveAnimation() | 사용 불가 → frame-by-frame Render()+SaveScreenshot() |
| Text/AnnotateTime | Mesa crash → ffmpeg drawtext 대체 |
| STL 경계 패턴 | `autofill=false` + `fillpoint modefill=void` + `drawbox` fluid |
| Qwen3 모델 | 32B 검증 완료, 8B 미검증 |

---

## v0.3 — TUI 고도화 + 에이전트 통합 E2E + 프로덕션 안정성

**테마:** 사용자가 실제로 쓸 수 있는 수준으로 UX 완성 + 에이전트 호출 E2E 검증
**예상 기간:** 4주
**의존 근거:** v0.2에서 18개 도구와 GPU 파이프라인은 검증됨. 그러나 "에이전트가 도구를 올바르게 호출하는" 통합 E2E와 TUI가 미완.

### 포함 기능

| ID | 기능명 | 설명 | 전제 |
|----|--------|------|------|
| TUI-04 | Case Wizard | 대화형 단계별 조건 설정 UI | TUI-02 |
| TUI-05 | Sim Dashboard 고도화 | 실시간 에너지/수렴 그래프, GPU 메모리 | MON-01 ✅ |
| AGENT-E2E | 에이전트 통합 E2E | `./slosim-agent -p "..."` → 도구 자동 호출 → 결과 | AGENT-01 ✅ |
| MDBC-01 | mDBC 경계조건 지원 | `BoundaryMethod=2` 옵션, dp<0.003 안정 | Scenario 8 결과 |
| GEO-02 | L형 탱크 E2E 검증 | 코드 구현됨, GPU E2E 미검증 | GEO-01 ✅ |
| EXC-E2E | 지진파 가진 E2E 검증 | seismic_input 코드 구현됨, GPU E2E 미검증 | EXC-01 코드 ✅ |
| QWEN-8B | Qwen3 8B 모델 검증 | 경량 모델 프롬프트 품질 테스트 | AGENT-01 ✅ |
| JOB-02 | 동시 다중 Job 관리 | Job 큐, 동시 실행 제한 | JOB-01 ✅ |

### v0.3에서 이미 v0.2에 앞당겨 완료된 항목

| 원래 v0.3 계획 | 상태 | 비고 |
|---------------|------|------|
| GEO-01 원통형 탱크 | ✅ v0.2에서 완료 | NASA 시나리오 + Zhao HorizCyl |
| STL-01 STL 입력 | ✅ v0.2에서 완료 | watertight 검증 구현 + Frosina E2E |
| ANI-01 애니메이션 | ✅ v0.2에서 완료 | Mesa 호환 pvpython + ffmpeg |
| EXC-01 지진파 가진 | ✅ 코드 완료 | GPU E2E 미검증 → v0.3 |

### 구현 순서

```
Week 1: AGENT-E2E (에이전트 통합 E2E — Qwen3 32B + 도구 호출 검증)
Week 2: QWEN-8B (8B 검증) + TUI-04 (Case Wizard)
Week 3: MDBC-01 + GEO-02/EXC-E2E (남은 도메인 GPU 검증)
Week 4: TUI-05 + JOB-02 + 통합 QA
```

---

## v1.0 — 파라메트릭 스터디 + 배포 패키징

**테마:** 전문가 수준 활용, 비교 분석, 프로덕션 배포
**예상 기간:** 4주
**의존 근거:** 파라메트릭 스터디는 멀티 Job(v0.3) + 시각화(v0.2 ✅)이 전제.

### 포함 기능

| ID | 기능명 | 설명 | 전제 |
|----|--------|------|------|
| PARA-01 | 파라메트릭 스터디 실행 | "충진율 50/60/70% 비교" → 자동 3건 | JOB-02 |
| PARA-02 | 이전 결과 비교 | 시뮬레이션 이력에서 선택 비교 | STORE-01 |
| PARA-03 | 비교 리포트 | 비교 테이블 + 오버레이 시각화 | PARA-01+02 |
| STORE-01 | Result Store (DB) | SQLite 메타 + 파일시스템 원본 | — |
| TUI-06 | Parametric View | 비교 테이블 + 결과 오버레이 뷰 | PARA-01 |
| NFR-02 | 템플릿 플러그인 | 새 케이스 파일 드롭 추가 | — |
| NFR-03 | 데이터 보존 | 시뮬레이션 결과 + 리포트 이력 | STORE-01 |
| PKG-01 | goreleaser 배포 | deb/rpm + Docker Hub 이미지 | CI ✅ |
| DOC-02 | 튜토리얼 + 온보딩 | 비전문가 step-by-step 가이드 | USER_MANUAL ✅ |

### 구현 순서

```
Week 1: STORE-01 (Result Store DB 스키마 + 마이그레이션)
Week 2: PARA-01 (파라메트릭 실행) + PARA-02 (결과 비교)
Week 3: PARA-03 (비교 리포트) + TUI-06 (Parametric View)
Week 4: PKG-01 + NFR-02/03 + DOC-02 + 전체 QA + 릴리스
```

---

## v1.x+ — 확장 플랫폼

**테마:** 새로운 유체/솔버/배포 환경으로의 확장
**예상 기간:** 지속적 개발 (기능별 2~4주)
**의존 근거:** v1.0의 안정된 기반 + 플러그인 시스템(NFR-02) 위에 확장.

### 후보 기능 (우선순위 미정)

| ID | 기능명 | 설명 | 전제 조건 |
|----|--------|------|----------|
| FLUID-01 | 추가 유체 물성 | LNG, 오일, 화학물질 (밀도/점성 DB) | DOM-02 확장 |
| CLOUD-01 | 클라우드 GPU 실행 | AWS/GCP GPU 인스턴스 원격 실행 | JOB-02 + SSH/API 연동 |
| MULTI-01 | 멀티 유저 지원 | 사용자별 세션/결과 격리 | STORE-01 확장 |
| API-01 | REST API 서버 모드 | TUI 없이 API로 시뮬레이션 제출/조회 | Agent Core 분리 |
| I18N-01 | 국제화 (영어/일본어) | 다국어 프롬프트 + UI 번역 | AGENT-01 확장 |
| SOLVER-01 | OpenFOAM VOF 솔버 통합 | SPH vs VOF 비교 분석 | Tool 인터페이스 확장 |
| COLLAB-01 | 팀 협업 기능 | 시뮬레이션 설정/결과 공유 | MULTI-01 + API-01 |

---

## PRD 기능 추적 매트릭스

### 입력 (IN)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| IN-01 | 자연어 입력 → AI 조건 추론 | v0.1 | ✅ |
| IN-02 | STL 파일 직접 입력 | v0.2 | ✅ (앞당겨 완료) |
| IN-03 | 단순 입력 → AI 표준 치수 제안 | v0.1 | ✅ |

### 시뮬레이션 엔진 (SIM)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| SIM-01 | DualSPHysics v5.4 GPU 통합 | v0.1 | ✅ |
| SIM-02 | GenCase XML 자동 생성 | v0.1 | ✅ |
| SIM-03 | AI 파라미터 자율 판단 | v0.1 | ✅ |
| SIM-04 | Job 백그라운드 실행 | v0.1 (단일) → v0.3 (다중) | ✅ 단일 / 다중 미완 |
| SIM-05 | 실시간 모니터링 | v0.2 | ✅ |

### 후처리 (POST)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| POST-01 | pvpython 후처리 | v0.2 | ✅ |
| POST-02 | 정지 이미지 스냅샷 | v0.2 | ✅ |
| POST-03 | 애니메이션 (mp4) | v0.2 | ✅ (앞당겨 완료) |
| POST-04 | 시각화 항목 | v0.2 | ✅ |

### 리포트 (RPT)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| RPT-01 | Markdown 리포트 자동 생성 | v0.1 | ✅ |
| RPT-02 | 해석 조건 요약 | v0.1 | ✅ |
| RPT-03 | 결과 이미지 포함 | v0.2 | ✅ |
| RPT-04 | AI 물리적 해석 코멘트 | v0.2 | ✅ |

### 파라메트릭 스터디 (PARA)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| PARA-01 | 파라메트릭 스터디 실행 | v1.0 | 코드 완료, E2E 미검증 |
| PARA-02 | 이전 결과와 비교 | v1.0 | 코드 완료, E2E 미검증 |
| PARA-03 | 비교 리포트 생성 | v1.0 | 미구현 |

### 슬로싱 도메인 (DOM)

| PRD ID | 기능 | 배치 버전 | 상태 |
|--------|------|----------|------|
| DOM-01 | 탱크 형상 다양화 | v0.1 직사각형 → v0.2 원통형 | ✅ L형은 v0.3 |
| DOM-02 | 유체: 물 | v0.1 | ✅ |
| DOM-03 | 가진 조건 다양화 | v0.1 정현파 → v0.3 지진파 | 정현파 ✅, 지진파 코드만 |
| DOM-04 | 2D / 3D | v0.1 (3D) | ✅ 3D, 2D 미검증 |

### PRD §10 미결 사항 (결정 완료)

| 미결 사항 | 결정 |
|----------|------|
| ParaView Docker vs 호스트 | **호스트 Mesa offscreen** (Docker 아님) |
| STL 메시 품질 검증 | **수밀성(edge-sharing) + normal consistency** (v0.2 구현 완료) |
| Qwen3 32b vs 8b | **32B 검증 완료**, 8B는 v0.3에서 검증 예정 |
| 결과 저장소 구조 | **하이브리드** (SQLite 메타 + 파일시스템 원본) |

---

## 전체 ID 체계

| 접두사 | 범주 | 버전 범위 |
|--------|------|----------|
| AGENT-* | AI 에이전트/프롬프트 | v0.1 |
| TOOL-* | DualSPHysics CLI 래핑 도구 | v0.1 |
| JOB-* | 백그라운드 Job 관리 | v0.1 ~ v0.3 |
| RPT-* | 리포트 생성 | v0.1 ~ v0.2 |
| TUI-* | BubbleTea TUI 컴포넌트 | v0.1 ~ v1.0 |
| PIPE-* | 파이프라인 오케스트레이션 | v0.1 |
| VIS-* | pvpython 시각화 | v0.2 |
| MON-* | 실시간 모니터링 | v0.2 |
| DIM-* | 차원 확장 (2D/3D) | v0.2 |
| GEO-* | 탱크 지오메트리 확장 | v0.3 |
| EXC-* | 가진 조건 확장 | v0.3 |
| STL-* | 외부 형상 입력 | v0.3 |
| ANI-* | 애니메이션 생성 | v0.3 |
| PARA-* | 파라메트릭 스터디 | v1.0 |
| STORE-* | 데이터 저장소 | v1.0 |
| NFR-* | 비기능 요구사항 | v1.0 |
| DOC-* | 문서/온보딩 | v1.0 |
| FLUID/CLOUD/... | 향후 확장 | v1.x+ |
