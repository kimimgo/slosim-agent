# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Agent Identity

- **이 에이전트 이름**: **cc-slosim-1**
- **역할**: slosim-agent 개발 에이전트
- **응답 시 이름 표시**: 매 응답 첫 줄에 `[cc-slosim-1]` 태그 포함

## Project Overview

slosim-agent는 자연어로 슬로싱(sloshing) 시뮬레이션을 설정·실행·분석하는 AI 에이전트다. OpenCode(Go + BubbleTea TUI) 포크 위에 DualSPHysics v5.4 GPU solver를 통합했다.

**핵심 흐름**: 자연어 입력 → DualSPHysics XML 케이스 생성 → GPU 시뮬레이션 실행 → 후처리 → AI 분석

**타겟 LLM**: Qwen3 32B/8B (오픈웨이트, 로컬 Ollama) — 슬로싱 도메인 특화 시스템 프롬프트 사용

## Build & Development Commands

```bash
# Build
go build -o slosim-agent ./main.go

# Build for deployment (no CGO, stripped)
CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/opencode-ai/opencode/internal/version.Version=dev" -o slosim ./main.go

# Run tests (all)
go test ./...

# Run a single test
go test ./internal/llm/tools/ -run TestGenCase

# Run with race detector + coverage
go test ./... -race -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Lint / vet
go vet ./...

# Generate DB code from SQL (requires sqlc)
sqlc generate

# Generate config JSON schema
go run ./cmd/schema/main.go > opencode-schema.json

# Docker: build DualSPHysics GPU image
docker compose build

# Docker: run simulation step by step (binaries at /opt/dsph/bin/)
docker compose run --rm --entrypoint '' -w /cases dsph /opt/dsph/bin/GenCase_linux64 /cases/SloshingTank_Def /data/out
docker compose run --rm --entrypoint '' dsph /opt/dsph/bin/DualSPHysics5.4_linux64 /data/out/SloshingTank_Def /data/out/SloshingTank_Def -gpu -svres
docker compose run --rm --entrypoint '' dsph /opt/dsph/bin/PartVTK_linux64 -dirin /data/out/SloshingTank_Def -savevtk /data/out/PartFluid -onlytype:-all,+fluid

# Non-interactive mode (prompt-driven)
./slosim -c . -p "0.6m x 0.3m 탱크에 물 50% 채우고 0.5Hz로 흔들어" -q -f json
```

## Architecture

### Module & Go Version

- Module: `github.com/opencode-ai/opencode`
- Go: 1.24.0
- CGO: disabled for release builds (`CGO_ENABLED=0`)

### Core Layers

```
cmd/root.go           → Cobra CLI entry (interactive TUI / non-interactive -p mode)
internal/app/         → App orchestrator (wires all services together)
internal/config/      → Viper-based config loader (.opencode/config.json)
internal/llm/
  ├── agent/          → Agent loop: coder, task, title, summarizer + MCP integration
  ├── provider/       → LLM providers (OpenAI, Gemini, GROQ, OpenRouter, XAI, Local/Ollama)
  ├── models/         → Provider-specific model definitions (openai, gemini, groq, openrouter, xai, local)
  ├── prompt/         → System prompts (SloshingCoderPrompt for sloshing domain)
  └── tools/          → Built-in tools (11: bash, view, glob, grep, ls, write, edit, fetch, sourcegraph, patch, diagnostics) + DualSPHysics tools (15)
internal/tui/         → BubbleTea TUI
  └── components/     → chat, core, dialog, error, logs, simulation, results, parametric, util
internal/db/          → SQLite3 + sqlc generated code + goose migrations
internal/session/     → Session CRUD
internal/message/     → Message CRUD
internal/pubsub/      → Event broadcast system
internal/permission/  → Tool execution permission checks
internal/lsp/         → LSP client (gopls, pyright integration)
internal/history/     → File version history
internal/logging/     → Structured logging
internal/diff/        → Diff utilities
internal/fileutil/    → File utilities
internal/format/      → Output formatting
internal/completions/ → CLI autocompletion
```

### Key Patterns

- **Service interface + pub/sub**: 모든 서비스는 `Service` interface를 구현하고, `pubsub.Subscriber[T]`로 상태 변경을 브로드캐스트
- **Tool interface**: `BaseTool.Info()` + `Run(ctx, ToolCall) → (ToolResponse, error)` — 모든 built-in 도구와 MCP 도구가 동일 인터페이스
- **Provider abstraction**: `Provider.SendMessages()` / `StreamResponse()` — 모든 LLM 프로바이더가 동일 인터페이스
- **Agent composition**: coder는 전체 도구 접근, task는 제한된 도구셋, title/summarizer는 도구 없음
- **Feature ID 체계**: FRD 문서와 연동된 기능 ID가 코드 주석에 표기 (예: `ANI-01`, `PARA-01`, `GEO-01`)

### DualSPHysics Tools (15개, CoderAgentTools에 등록된 도구)

| Tool | Feature ID | 역할 |
|------|-----------|------|
| `xml_generator` | — | 탱크 치수 → DualSPHysics XML 자동 생성 |
| `stl_import` | STL-01 | STL 파일 입력 (CAD 메시 → DualSPHysics) |
| `gencase` | — | GenCase 실행 (XML → 바이너리 파티클) |
| `solver` | — | DualSPHysics GPU 솔버 실행 |
| `partvtk` | — | VTK 시각화 데이터 추출 |
| `measuretool` | — | 측정점 데이터 추출 |
| `baffle_generator` | — | 배플(내부 벽) STL 자동 생성 |
| `report` | — | Markdown 리포트 생성 |
| `analysis` | RPT-03 | AI 물리적 해석 (Qwen3 기반 도메인 분석) |
| `monitor` | MON-01 | 실시간 시뮬레이션 모니터링 (Run.csv 파싱) |
| `error_recovery` | NFR-01 | 에러 감지 및 복구 전략 |
| `job_manager` | JOB-02 | 백그라운드 Job 관리 (장시간 시뮬레이션) |
| `seismic_input` | EXC-01 | 지진파 시계열 파싱 |
| `parametric_study` | PARA-01 | 파라메트릭 스터디 오케스트레이션 |
| `result_store` | STORE-01 | 결과 저장/검색 |

**헬퍼 모듈** (독립 도구가 아닌 내부 유틸리티):
- `geometry.go` — 탱크 형상 생성 함수 (rectangular, cylindrical, l_shaped) → xml_generator가 사용
- `file.go` — 파일 읽기/쓰기 시간 추적 유틸리티

### Database

SQLite3 with goose migrations (`internal/db/migrations/`). 3개 테이블:
- `sessions` — 대화 세션
- `messages` — 메시지 (JSON parts 배열)
- `files` — 파일 버전 히스토리

SQL 쿼리는 `internal/db/sql/`에 작성 → `sqlc generate`로 Go 코드 자동 생성.

## DualSPHysics Integration

### Critical Rules

1. **GenCase는 `.xml` 자동 추가** — 경로에 `.xml` 확장자 포함하지 않는다
2. **XML은 attribute-only 문법** — `<gravity x="0" y="0" z="-9.81" />` (자식 노드 아님)
3. **GPU 타겟**: NVIDIA RTX 4090 (sm_89, Ada Lovelace), CUDA 12.6
4. **STL Import (Frosina 패턴)** — STL 형상을 경계 파티클로 가져올 때의 필수 순서:
   - `drawfilestl autofill="false"` → `fillpoint(modefill void)` → `drawbox(fluid)` 순서 엄수
   - `autofill="false"` 속성 필수 (scale 속성이 아님)
   - scale≠1.0일 때만 `<drawscale>` 자식 요소 추가
   - motion에 `<begin mov="1" start="0" />` 요소 필수
   - `<simulationdomain>` 확장 (posmin -20%/-10%/-10%, posmax +20%/+10%/+50%)으로 BoundaryOut 방지

### Case Files

`cases/` 디렉토리에 21개+ 사전 정의 슬로싱 시나리오:
- `SloshingTank_Def.xml` — 기본 벤치마크
- `Sloshing_{Normal,NearRes,Res}[_Guard]_Def.xml` — 주파수별/가드 유무 변형 (6개)
- `SPHERIC_Test10_{High,Low,Oil_Low}_Def.xml` — SPHERIC 벤치마크 (3개)
- `Chen2018_*.xml`, `English2021_*.xml`, `Liu2024_*.xml` — 논문 검증 케이스 (6개)
- `Frosina2018_FuelTank_Def.xml` — STL import 레퍼런스 (Frosina 패턴)
- `NASA2023_Cylinder_Def.xml`, `Zhao2024_HorizCyl_Def.xml` — 원통형 탱크
- `ISOPE2012_LNG_Def.xml` — LNG 대형 탱크
- `cases/fuel_tank.stl` — STL 형상 파일 (S10 시나리오용)

### Docker Pipeline

```
Dockerfile:          CUDA 12.6 devel → MoorDynPlus 빌드 → DualSPHysics GPU 빌드 → runtime 이미지
Dockerfile.pvpython: ParaView pvpython 전용 이미지 (headless 렌더링)
docker-compose.yml:  NVIDIA GPU reservation, ./simulations:/data, ./cases:/cases 마운트
```

**컨테이너 바이너리 경로**: `/opt/dsph/bin/` (GenCase_linux64, DualSPHysics5.4_linux64, PartVTK_linux64, MeasureTool_linux64, IsoSurface_linux64)
**주의**: docker-compose.yml의 기본 entrypoint는 `sleep infinity`이므로 `--entrypoint ''`로 오버라이드 필요

### Solver Toolchain

| Binary | 역할 |
|--------|------|
| GenCase | XML → 바이너리 파티클 전처리 |
| DualSPHysics5.4_linux64 | GPU SPH 솔버 |
| PartVTK | VTK 시각화 데이터 추출 |
| MeasureTool | 측정점 데이터 추출 |
| IsoSurface | 등위면 메시 생성 |

## TUI Components

```
internal/tui/components/
├── chat/        → 메인 대화 인터페이스
├── core/        → 공통 UI 요소
├── dialog/      → Dialog 시스템 (custom commands 포함)
├── error/       → 에러 표시 컴포넌트
├── logs/        → 로그 뷰어
├── simulation/  → Sim Dashboard (Job 목록, 진행률, 로그)
├── results/     → 결과 뷰어 (VTK, Plot)
├── parametric/  → 파라메트릭 비교 뷰
└── util/        → 유틸리티 컴포넌트
```

BubbleTea 패턴: Model struct + Init/Update/View, Lipgloss 스타일링, pubsub 이벤트 연동.

## Custom Agents & Commands

### Sub-agents (`.claude/agents/`)

| Agent | 역할 |
|-------|------|
| `tui-architect` | TUI 디자인 전문가 (SemanticTokens, 레이아웃 패턴) |
| `dsph-engineer` | DualSPHysics 도구 구현/디버깅 (XML, Docker GPU, 솔버) |
| `provider-specialist` | LLM Provider 코드 (Ollama/gemini-cli, 스트리밍) |
| `qa-tester` | E2E/단위 테스트 (testify, mock, coverage) |
| `scenario-validator` | UX 품질 스캐너 (전문용어, 에러메시지, 비전문가 검증) |
| `prompt-optimizer` | Qwen3 SLM 프롬프트 최적화 |
| `e2e-nl-validator` | NL→tool call 파이프라인 검증 (Ollama API 직접 호출) |
| `pipeline-tester` | Docker GPU E2E 파이프라인 테스트 |
| `dsph-xml-validator` | DualSPHysics XML 케이스 검증 |
| `tui-test-runner` | TUI 컴포넌트 편집 후 자동 테스트 |

### Slash commands (`.claude/commands/`)

| Command | 역할 |
|---------|------|
| `/bubbletea-component <name>` | BubbleTea 컴포넌트 스캐폴딩 |
| `/write-scenario <feature>` | BDD 시나리오 자동 생성 |
| `/implement-tool <name>` | DualSPHysics 도구 구현 템플릿 |

## CI/CD & Release

- `.github/workflows/build.yml` — main 브랜치 푸시 시 빌드
- `.github/workflows/release.yml` — 태그 푸시 시 goreleaser 릴리스
- `.goreleaser.yml` — 멀티플랫폼 빌드 (linux/darwin × amd64/arm64), AUR, Homebrew tap, deb/rpm

## Configuration

설정 파일: `.opencode/config.json` (또는 `~/.opencode.json` 글로벌)

핵심 설정 항목:
- `agents.coder.model` — 주 에이전트 모델 (예: `local.qwen3:32b`)
- `agents.task.model` — 백그라운드 작업 모델
- `providers.local` — Ollama/LM Studio 연결 (`LOCAL_ENDPOINT` 환경변수)
- `providers.*` — LLM 프로바이더 API 키/설정
- `mcpServers.*` — MCP 서버 등록 (stdio/SSE)
- `lsp.*` — LSP 서버 설정 (gopls, pyright 등)
- `contextPaths` — 자동 로딩할 컨텍스트 파일 목록

### SloshingCoderPrompt (시스템 프롬프트)

`internal/llm/prompt/sloshing_coder.go` — Qwen3용 도메인 특화 프롬프트:
- Tool 호출 순서 강제 (A: 직사각형, B: STL, C: Baffle)
- 파라미터 자동 결정 규칙 (dp, time_max, fill_height, amplitude)
- 경로 규칙 (simulations/{case_name}/vtk/, measure/, viz/)
- pv-agent MCP 후처리 도구 안내

## Testing

테스트 프레임워크: stdlib `testing` + `stretchr/testify`

도구별 테스트 파일이 `internal/llm/tools/*_test.go`에 위치. Integration 테스트는 `*_integration_test.go` 접미사 사용 (error_recovery, parametric_study, result_store).

BDD 시나리오: `docs/scenarios/*.feature` (Gherkin 형식, CFD 비전문가 관점)

## Project Documents

| 문서 | 역할 |
|------|------|
| `ARCHITECTURE.md` | CC Agent Teams 기반 개발 체계 (팀 구성, 워크플로우) |
| `AGENTS.md` | 제품 에이전트 시스템 개요 |
| `PRD.md` | Product Requirements Document |
| `docs/FRD_v0.1.md` | Functional Requirements (Feature ID 정의) |
| `docs/ROADMAP.md` | 개발 로드맵 |
| `docs/USER_MANUAL.md` | 비전문가 대상 사용자 매뉴얼 |

## Git Branch Strategy

```
main                        ← slosim-agent 구현체 (Go 소스, TUI, 도구)
├── research-v1             ← EXP 수행 이력 v1 (CLOSED)
├── research-v2             ← EXP 수행 이력 v2 (CLOSED)
├── research-v3             ← EXP 수행 이력 v3 (CLOSED)
├── research-v4             ← 현재 활성 (cross-model 실험 + 코드 수정)
├── paper                   ← (LEGACY) 구 통합 논문 브랜치
├── paper-cs                ← 1차 CS 저널 (EAAI/EwC) — 에이전트 아키텍처 + ablation
└── paper-pof               ← 2차 PoF — 물리 검증 + SPH 수렴성 + 배플 최적화
```

### 브랜치별 규칙

| 브랜치 | 역할 | 추적 대상 | 금지 대상 |
|--------|------|----------|----------|
| `main` | agent 코드 개발 | Go 소스, 테스트, CI, docs | 시뮬레이션 데이터, 논문, 실험 이력 |
| `research-vN` | EXP 수행 이력 (try & error) | 케이스 XML, RUNBOOK, 실행 로그, 스크립트 | 대용량 바이너리 (VTK, bi4) |
| `paper-cs` | 1차 CS 논문 | LaTeX 원고, Table/Figure, ablation 분석, cross-model 데이터 | 대용량 바이너리 |
| `paper-pof` | 2차 PoF 논문 | LaTeX 원고, 물리 검증 심화, 수렴 분석, 배플 파라메트릭 | 대용량 바이너리 |

### 데이터 흐름

```
research-v4 (실험 완료) ──┬→ paper-cs  (에이전트 아키텍처, ablation, cross-model)
                         └→ paper-pof (물리 검증, 수렴성, 배플 최적화 확장)
main (구현 완료)        ──→ paper-cs  (에이전트 아키텍처 설명 참조)
```

### 시뮬레이션 데이터

대용량 PART/VTK 파일은 git에 포함하지 않음.
- 로컬: `/mnt/simdata/dualsphysics/` (심볼릭 링크 `./simulations` → `/mnt/simdata/dualsphysics/`)
- 논문에 필요한 CSV/메트릭만 `paper/analysis/`에 포함

## Research Infrastructure (research-v3)

### 실험 체계

| 실험 | 목적 | 상태 |
|------|------|------|
| EXP-A | M-A3 parameter fidelity baseline (10 시나리오 × 2 모델 × 3 trials) | 완료 |
| EXP-B | 2×2 factorial ablation (B0=full, B1=-prompt, B2=-tool, B4=bare) | 완료 |
| EXP-C | 이관 완료 | 완료 |
| EXP-D | 미실행 | — |

### 채점 도구

- `research-v3/exp-b/score_expb.py` — M-A3 채점 (XML 파싱 → ground truth 비교)
  - `--xml <file> --scenario S10` 단일 XML 채점
  - 전체 실행: `python3 score_expb.py` (results/ 디렉토리 자동 스캔)
- `research-v3/exp-a/ground_truth.json` — 10개 시나리오 정답 (S01-S10)
- `research-v3/exp-a/run_scenario.sh` — EXP-A 시나리오 실행 스크립트
  - Usage: `./run_scenario.sh <scenario> <model> <trial> [gpu_id]`

### pajulab 원격 인프라

- **SSH alias**: `pajulab` (icap@192.168.0.3)
- **GPU**: 4× RTX 4090 (GPU 0,1=Ollama, GPU 2=DualSPHysics Docker, GPU 3=여유)
- **Ollama**: `http://localhost:11434` (qwen3:32b, qwen3:latest)
- **Docker**: `dsph-agent:latest` 이미지 (DualSPHysics GPU 빌드)
- **바이너리**: `~/slosim-agent/slosim` (로컬 빌드 후 scp 배포)
- **환경변수**: `.env` → `LOCAL_ENDPOINT=http://localhost:11434`
