# E2E 테스트 심층 분석 보고서

> 테스트 일시: 2026-02-14
> 테스트 환경: RTX 4090, Docker dsph-agent:latest, pvpython 6.0.1, ffmpeg 4.4.2

---

## 1. 실행 이력 요약

### 에이전트(cc)가 실제 수행한 Action 시퀀스

| # | Action | 결과 | 소요시간 | 비고 |
|---|--------|------|----------|------|
| 1 | `go test ./internal/llm/tools/` | FAIL | 0.3s | 바이너리명 불일치 발견 |
| 2 | Docker 바이너리명 확인 (`docker compose run --rm dsph ls /opt/dsph/bin/`) | 소문자 symlink 확인 | 2s | |
| 3 | gencase.go 수정: `GenCase` → `gencase`, `-save:` → positional arg | 코드 수정 | - | |
| 4 | solver.go 수정: `DualSPHysics5.4_linux64` → `dualsphysics` | 코드 수정 | - | |
| 5 | partvtk.go 수정: `PartVTK` → `partvtk` | 코드 수정 | - | |
| 6 | measuretool.go 수정: `MeasureTool` → `measuretool` | 코드 수정 | - | |
| 7 | job_manager_test.go 수정 | 코드 수정 | - | |
| 8 | GenCase 실행 (0.1s sim) | FAIL - domain too small | 3s | |
| 9 | XML 수정: SimulationDomain X ±10% 마진 추가 | 코드 수정 | - | |
| 10 | GenCase 재실행 | SUCCESS - 168,681 particles | 3s | |
| 11 | DualSPHysics GPU solver 실행 (0.1s sim) | SUCCESS - 6 PART files | 5s | |
| 12 | PartVTK 실행 | SUCCESS - 6 VTK files | 2s | |
| 13 | MeasureTool 실행 | FAIL - "no valid points" | 2s | |
| 14 | probe_points.txt 수정: `POINTS` 헤더 추가 | 파일 수정 | - | |
| 15 | MeasureTool 재실행 | SUCCESS - CSV 생성 | 2s | |
| 16 | XML 수정: TimeMax 0.1→2.0, TimeOut 0.05→0.02 | 파일 수정 | - | |
| 17 | GenCase → Solver 재실행 (2.0s sim) | SUCCESS - 101 PART, 38s | 41s | |
| 18 | PartVTK 실행 (101 files) | SUCCESS - 101 VTK | 10s | |
| 19 | MeasureTool 실행 | SUCCESS - CSV 생성 | 3s | |
| 20 | pvpython_simple.py 작성 + 실행 | SUCCESS - 252KB PNG | 5s | |
| 21 | pvpython_animation.py 작성 (SaveAnimation) | FAIL - 0 frames | 10s | Mesa silent fail |
| 22 | pvpython_animation.py 수정 (manual loop) | FAIL - crash | 3s | Box+GradientBG |
| 23 | pvpython_anim_minimal.py (파티클만) | SUCCESS - 101 frames | 60s | |
| 24 | ffmpeg MP4 생성 (minimal) | SUCCESS - 6.2MB | 3s | |
| 25 | pvpython_test_box.py (Box+Gradient) | FAIL - crash | 3s | |
| 26 | pvpython_test_box.py (Box only) | SUCCESS | 5s | root cause 확인 |
| 27 | pvpython_animation.py 최종 (no gradient) | SUCCESS - 101 frames | 90s | 1920x1080 |
| 28 | ffmpeg MP4 + drawtext overlay | SUCCESS - 31MB MP4 | 5s | |
| 29 | snapshot_peak.png 생성 | SUCCESS - 1.1MB | 2s | |

**총 Action: 29개, 실패 후 재시도: 6회, 총 소요: ~15분**

---

## 2. 시나리오 ↔ 실제 테스트 대조

### 커버된 기능 (시나리오 Act 기준)

| Act | 시나리오 기능 | 실제 테스트 | 커버 수준 | 판정 |
|-----|-------------|-----------|----------|------|
| Act 3 | GenCase 파티클 전처리 | `gencase` Docker 실행 → 168,681 particles | **완전** | PASS (버그 수정 후) |
| Act 3 | GPU Solver 실행 | `dualsphysics` 2.0s sim, 101 PART | **완전** | PASS (버그 수정 후) |
| Act 5 | PartVTK 추출 | `partvtk` → 101 VTK files | **완전** | PASS (버그 수정 후) |
| Act 5 | MeasureTool 측정 | `measuretool` → CSV (Elevation 등) | **완전** | PASS (버그 수정 후) |
| Act 6 | ParaView 스냅샷 | pvpython → snapshot_peak.png (1.1MB) | **완전** | PASS |
| Act 6 | 애니메이션 생성 | pvpython → 101 frames → MP4 (31MB) | **완전** | PASS (3차 시도 후) |

### 커버되지 않은 기능

| Act | 시나리오 기능 | 미커버 이유 | 심각도 |
|-----|-------------|-----------|--------|
| Act 1 | 자연어 → AI 조건 추론 | TUI + LLM 통합 테스트 미수행 | **HIGH** |
| Act 1 | SloshingCoderPrompt | 프롬프트 ↔ Qwen3 응답 품질 미검증 | **HIGH** |
| Act 2 | 원통형 형상 변경 | geometry tool 미테스트 | MEDIUM |
| Act 2 | CSV 파도 데이터 임포트 | seismic_input tool 미테스트 | MEDIUM |
| Act 3 | Job Manager 백그라운드 실행 | 직접 docker compose 호출로 우회 | MEDIUM |
| Act 3 | 실시간 모니터링 Dashboard | TUI 컴포넌트 미테스트 | **HIGH** |
| Act 4 | 에러 감지 + 자동 복구 | error_recovery tool 미테스트 | **HIGH** |
| Act 5 | AI 물리 분석 (analysis) | analysis tool 미테스트 | MEDIUM |
| Act 5 | Markdown 리포트 | report tool 미테스트 | MEDIUM |
| Act 7 | 파라메트릭 스터디 | 테스트 제외 (유저 요청) | N/A |
| Act 8 | 결과 저장/검색 | result_store 미테스트 | LOW |
| Act 9 | STL 임포트 | stl_import 미테스트 | LOW |

---

## 3. 발견된 버그 및 수정

### Bug #1: Docker 바이너리명 불일치 (Critical)

| 항목 | 내용 |
|------|------|
| **증상** | 모든 Docker 도구 실행 시 `exit status 127` (command not found) |
| **원인** | Go 코드에서 CamelCase (`GenCase`, `PartVTK`, `MeasureTool`, `DualSPHysics5.4_linux64`) 사용, Docker 이미지 내부는 lowercase symlink (`gencase`, `partvtk`, `measuretool`, `dualsphysics`) |
| **영향** | 4개 핵심 도구 전부 작동 불가 → 시뮬레이션 파이프라인 완전 차단 |
| **수정** | 4개 소스파일 + 1개 테스트 파일 바이너리명 변경 |
| **근본원인** | Dockerfile에서 소문자 symlink를 만든 시점과 Go 코드 작성 시점 사이의 연동 부재 |

### Bug #2: GenCase 파라미터 문법 오류 (Critical)

| 항목 | 내용 |
|------|------|
| **증상** | GenCase 실행 시 `-save:path` 옵션 에러 |
| **원인** | `-save:`는 저장 형식 옵션 (vtkcells 등)이지 출력 경로 지정이 아님. GenCase는 positional args: `gencase <input> <output>` |
| **영향** | GenCase 완전 실패 |
| **수정** | `fmt.Sprintf("-save:%s", params.SavePath)` → `params.SavePath` (positional arg) |
| **근본원인** | DualSPHysics CLI 문서 미확인 / 다른 도구(PartVTK)의 `-save` 문법과 혼동 |

### Bug #3: SimulationDomain 범위 부족 (Major)

| 항목 | 내용 |
|------|------|
| **증상** | Solver 실행 시 "boundary particle out of domain" 에러 |
| **원인** | 탱크 모션 진폭 ±0.05m인데 SimulationDomain X 마진이 없었음 |
| **영향** | Solver 즉시 중단 |
| **수정** | `posmin x="default - 10%"`, `posmax x="default + 10%"` 추가 |
| **근본원인** | XML 케이스 작성 시 모션에 의한 도메인 확장 미고려 |

### Bug #4: MeasureTool 포인트 파일 형식 (Major)

| 항목 | 내용 |
|------|------|
| **증상** | "no valid points in file" 에러 |
| **원인** | DualSPHysics MeasureTool은 포인트 파일 첫 줄에 `POINTS` 키워드 필수 |
| **영향** | 측정 데이터 추출 불가 |
| **수정** | probe_points.txt 첫 줄에 `POINTS` 추가 |
| **근본원인** | 포인트 파일 포맷 스펙 미확인 |

### Bug #5: ParaView Mesa 호환성 (Major)

| 항목 | 내용 |
|------|------|
| **증상** | (a) SaveAnimation() silent fail, (b) Box()+GradientBackground crash |
| **원인** | Mesa 소프트웨어 렌더러에서 특정 VTK 파이프라인 조합 미지원 |
| **영향** | 애니메이션 생성 실패 |
| **수정** | (a) 수동 frame-by-frame 루프, (b) gradient background 제거 |
| **근본원인** | headless 서버에서 ParaView Mesa 백엔드의 제약사항 |

---

## 4. 객관적 평가

### 좋은 점 (Strengths)

**S1. 핵심 파이프라인 E2E 검증 완료**
- GenCase → Solver → PartVTK → MeasureTool → pvpython → ffmpeg 전체 체인이 실제 데이터로 검증됨
- 168,681 파티클, 2초 시뮬레이션, RTX 4090에서 38초 → 성능 벤치마크 확보

**S2. 실제 버그 4개 발견 및 즉시 수정**
- Bug #1~#4 모두 "코드는 있지만 실제로 한번도 돌려보지 않았다"는 사실을 드러냄
- 특히 Bug #1(바이너리명)은 모든 도구에 영향 → E2E 없이는 출시 불가능했을 치명적 결함

**S3. 시각화 품질 확인**
- 압력 컬러맵이 물리적으로 올바름 (정수압 분포: 바닥 고압, 상부 저압)
- 슬로싱 자유표면 형상이 명확
- 1920x1080 고해상도, 압력 컬러바 + 탱크 와이어프레임 + 텍스트 오버레이

**S4. 환경 의존성 문서화**
- Mesa 백엔드 제약사항 3건 발견 및 우회법 확립
- Docker 바이너리 네이밍 규칙 명확화
- DualSPHysics CLI 파라미터 정확한 문법 확인

### 나쁜 점 (Weaknesses)

**W1. 시나리오 18개 기능 중 6개만 실제 테스트 (33% 커버리지)**
- 커버: gencase, solver, partvtk, measuretool, pvpython, animation
- 미커버: **12개** (SloshingCoderPrompt, xml_generator, geometry, seismic_input, job_manager 백그라운드, monitor, error_recovery, analysis, report, parametric_study, result_store, stl_import)
- **특히 Act 1 (자연어→AI추론)과 Act 4 (에러복구)가 미검증 → 사용자 경험의 핵심이 테스트되지 않음**

**W2. TUI 컴포넌트 전혀 미테스트**
- 시나리오의 5개 TUI 컴포넌트 (dashboard, error panel, results viewer, parametric view, results browser) 중 0개 테스트
- slosim-agent의 차별화 요소인 BubbleTea TUI가 검증되지 않음

**W3. LLM 통합 미검증**
- Qwen3:32b ↔ SloshingCoderPrompt 조합의 응답 품질, 할루시네이션 여부, 도구 호출 정확도 등 미확인
- 이것이 실제 사용자 경험의 80%를 차지하는 부분

**W4. Docker 호출만 검증, slosim-agent 자체는 미실행**
- 모든 테스트가 `docker compose run` 직접 호출로 수행
- slosim-agent 바이너리를 통한 도구 실행 (LLM → Tool interface → Docker)은 미검증
- 즉, "도구가 Docker를 올바로 호출하는가"는 확인했지만 "에이전트가 도구를 올바로 호출하는가"는 미확인

**W5. Go 단위 테스트 여전히 실패**
- GenCase, PartVTK, MeasureTool 테스트: Docker 미설치 환경(CI) 대응 미비
- 테스트가 실제 Docker를 호출하므로 Docker 없는 환경에서는 항상 FAIL
- Mock/Stub 기반 단위 테스트 부재

**W6. 에러 재현 → 수정 사이클이 반복적**
- ParaView 이슈에서 6번의 시행착오 (SaveAnimation → minimal → Box test → gradient test → 최종)
- 체계적 이분법이 아닌 "하나씩 빼보기" 방식 → 효율성 낮음

---

## 5. 정량 지표

| 지표 | 값 | 평가 |
|------|-----|------|
| 시나리오 기능 커버리지 | 6/18 (33%) | 낮음 |
| TUI 컴포넌트 커버리지 | 0/5 (0%) | 미수행 |
| 발견 버그 수 | 5개 (Critical 2, Major 3) | 높은 발견율 |
| 버그 수정율 | 5/5 (100%) | 완전 |
| 실패 후 재시도 횟수 | 6회 | 보통 |
| 파이프라인 E2E 통과 | YES | 핵심 성과 |
| 시각화 물리적 정확성 | 정수압 분포 일치 | 양호 |
| 시뮬레이션 성능 | 168K particles, 38s/2s-sim | RTX 4090 기준 양호 |
| Go 단위 테스트 통과율 | Ls/JobManager PASS, GenCase/PartVTK/MeasureTool FAIL | Docker 의존성 문제 |

---

## 6. 인간 vs 에이전트 역할 회고

### 인간이 실제로 판단한 항목
- "장난해?" → 역할 분배 재정의 (에이전트가 모든 실행 담당)
- "paraview로 에니메이션 만들어봐" → 시각화 방향 지시
- "7번은 테스트에서 제외" → 테스트 범위 결정
- 시각화 결과물 품질 확인 (현재 진행 중)

### 에이전트가 자율적으로 판단한 항목
- 버그 root cause 분석 및 수정 방향
- Docker 바이너리명 확인 방법
- XML 파라미터 조정 (SimulationDomain 마진)
- ParaView Mesa 호환성 이슈 이분법 디버깅
- ffmpeg drawtext로 텍스트 오버레이 대체

### 역할 분배 효과성
- **효과적**: 에이전트가 반복적 디버깅/수정/재실행을 자율 수행 → 인간 시간 절약
- **비효과적**: 에이전트가 테스트 범위를 자율 확장하지 않음 → 인간이 "다음 뭐 할까" 지시 필요
- **개선점**: 테스트 체크리스트를 먼저 만들고, 에이전트가 자율적으로 순서대로 진행하는 방식이 더 효율적

---

## 7. 우선순위별 후속 작업

| 우선순위 | 작업 | 이유 |
|----------|------|------|
| P0 | slosim-agent 바이너리 통한 E2E (TUI + LLM + Tool) | 실제 사용자 경로 미검증 |
| P0 | Go 단위 테스트 Docker mock 추가 | CI에서 항상 실패 |
| P1 | Act 1 검증: Qwen3 + SloshingCoderPrompt 응답 품질 | UX 핵심 |
| P1 | Act 4 검증: error_recovery 도구 | 안정성 핵심 |
| P2 | TUI 컴포넌트 테스트 | 차별화 요소 |
| P2 | Act 2 검증: geometry, seismic_input | 기능 완성도 |
| P3 | analysis, report 도구 테스트 | 부가 기능 |
| P3 | result_store, stl_import 테스트 | 부가 기능 |
