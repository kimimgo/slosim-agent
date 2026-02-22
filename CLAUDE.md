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

# Build with version info
go build -ldflags "-s -w -X github.com/opencode-ai/opencode/internal/version.Version=dev" ./main.go

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

# Docker: run simulation step by step
docker compose run --rm dsph GenCase /cases/SloshingTank_Def -save:/data/out
docker compose run --rm dsph DualSPHysics5.4_linux64 /data/out/SloshingTank_Def -gpu
docker compose run --rm dsph PartVTK -dirin /data/out -savevtk /data/out/PartFluid -onlytype:-all,+fluid
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
  ├── provider/       → LLM providers (Anthropic, OpenAI, Gemini, Copilot, Bedrock, Azure, VertexAI)
  ├── models/         → Provider-specific model definitions
  ├── prompt/         → System prompts (SloshingCoderPrompt for sloshing domain)
  └── tools/          → Built-in tools (13) + DualSPHysics tools (18)
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

### DualSPHysics Tools (18개)

| Tool | Feature ID | 역할 |
|------|-----------|------|
| `gencase` | — | GenCase 실행 (XML → 바이너리 파티클) |
| `solver` | — | DualSPHysics GPU 솔버 실행 |
| `partvtk` | — | VTK 시각화 데이터 추출 |
| `measuretool` | — | 측정점 데이터 추출 |
| `xml_generator` | — | 탱크 치수 → DualSPHysics XML 자동 생성 |
| `job_manager` | JOB-02 | 백그라운드 Job 관리 (장시간 시뮬레이션) |
| `monitor` | MON-01 | 실시간 시뮬레이션 모니터링 (Run.csv 파싱) |
| `analysis` | RPT-03 | AI 물리적 해석 (Qwen3 기반 도메인 분석) |
| `report` | — | Markdown 리포트 생성 |
| `pvpython` | — | ParaView CLI 렌더링 |
| `animation` | ANI-01 | 애니메이션 생성 (mp4/gif) |
| `parametric_study` | PARA-01 | 파라메트릭 스터디 오케스트레이션 |
| `result_store` | STORE-01 | 결과 저장/검색 |
| `error_recovery` | NFR-01 | 에러 감지 및 복구 전략 |
| `geometry` | GEO-01/02 | 탱크 형상 (rectangular, cylindrical, l_shaped) |
| `seismic_input` | EXC-01 | 지진파 시계열 파싱 |
| `stl_import` | STL-01 | STL 파일 입력 (CAD 메시 → DualSPHysics) |
| `isosurface` | — | IsoSurface 등위면 메시 생성 (미구현) |

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

### Case Files

`cases/` 디렉토리에 7개 사전 정의 슬로싱 시나리오:
- `SloshingTank_Def.xml` — 기본 벤치마크
- `Sloshing_{Normal,NearRes,Res}[_Guard]_Def.xml` — 주파수별/가드 유무 변형

### Docker Pipeline

```
Dockerfile:          CUDA 12.6 devel → MoorDynPlus 빌드 → DualSPHysics GPU 빌드 → runtime 이미지
Dockerfile.pvpython: ParaView pvpython 전용 이미지 (headless 렌더링)
docker-compose.yml:  NVIDIA GPU reservation, ./simulations:/data, ./cases:/cases 마운트
```

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
| `tui-test-runner` | TUI 컴포넌트 편집 후 자동 테스트 |
| `scenario-validator` | BDD 시나리오 vs 구현 검증 |
| `dsph-xml-validator` | DualSPHysics XML 케이스 검증 |
| `prompt-optimizer` | Qwen3 SLM 프롬프트 최적화 |

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

설정 파일: `.opencode/config.json`

핵심 설정 항목:
- `agents.coder.model` — 주 에이전트 모델
- `agents.task.model` — 백그라운드 작업 모델
- `providers.*` — LLM 프로바이더 API 키/설정
- `mcpServers.*` — MCP 서버 등록 (stdio/SSE)
- `lsp.*` — LSP 서버 설정 (gopls, pyright 등)
- `contextPaths` — 자동 로딩할 컨텍스트 파일 목록

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
main                    ← agent 개발, semver 태깅 (v0.4.0, v0.5.0...)
├── research/runs       ← long-lived, 시뮬레이션 실행 결과 (메타데이터만)
└── research/paper      ← long-lived, runs 결과 기반 논문 작성
```

### 데이터 흐름

```
main (v0.x.0 태그) ─ merge ──→ research/runs ─ merge ──→ research/paper
```

### 브랜치별 규칙

| 브랜치 | 역할 | 추적 대상 | 금지 대상 |
|--------|------|----------|----------|
| `main` | agent 코드 개발 | Go 소스, 테스트, CI, docs | 시뮬레이션 데이터, 논문 |
| `research/runs` | 시뮬레이션 실행 | 스크립트, 설정, CSV 요약, PNG | VTK, bi4, 대용량 바이너리 |
| `research/paper` | 논문 작성 | LaTeX, figures, bib, 분석 | 원본 시뮬레이션 바이너리 |
