# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Agent Identity

- **이 에이전트 이름**: **cc-slosim-1**
- **역할**: slosim-agent 논문 연구 에이전트
- **응답 시 이름 표시**: 매 응답 첫 줄에 `[cc-slosim-1]` 태그 포함

## Branch: `research/paper`

이 브랜치는 slosim-agent 시스템을 기반으로 **CS 분야 학술 논문**을 작성하기 위한 연구 브랜치다.

### 논문 주제

**AI Agent-Driven Sloshing Simulation: Automating DualSPHysics CFD Workflows with Domain-Specialized LLM Agents**

### Research Questions

- **RQ1**: LLM 에이전트가 자연어에서 유효한 DualSPHysics XML 설정을 생성할 수 있는가?
- **RQ2**: 에이전트 생성 케이스의 시뮬레이션 정확도가 전문가 기준(SPHERIC, 논문 재현)과 비교하여 어떠한가?
- **RQ3**: 반복적 CFD 워크플로우에서 신뢰성 있는 에이전트-솔버 통합을 위한 도구 인터페이스 설계 패턴은?
- **RQ4**: 도메인 특화 프롬프팅이 범용 LLM 대비 에이전트 성능을 어떻게 향상시키는가?

### 타겟 출판

CS 분야 (AI for Science, HCI+AI, Software Engineering 저널/학회)

### Novelty Claim

자연어에서 검증된 결과까지 DualSPHysics 슬로싱 시뮬레이션 파이프라인 전체를 자율적으로 구동하는 최초의 도메인 특화 LLM 에이전트.

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

# Docker: run simulation step by step (binary names are lowercase symlinks)
docker compose run --rm dsph gencase /cases/SloshingTank_Def /data/out
docker compose run --rm dsph dualsphysics /data/out/SloshingTank_Def -gpu
docker compose run --rm dsph partvtk -dirin /data/out -savevtk /data/out/PartFluid -onlytype:-all,+fluid
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
  └── tools/          → Built-in tools (13) + DualSPHysics tools (14)
internal/tui/         → BubbleTea TUI
  └── components/     → chat, core, dialog, error, logs, simulation, results, parametric, wizard, util
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

### DualSPHysics Tools (14개)

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
| `parametric_study` | PARA-01 | 파라메트릭 스터디 오케스트레이션 |
| `result_store` | STORE-01 | 결과 저장/검색 |
| `error_recovery` | NFR-01 | 에러 감지 및 복구 전략 |
| `seismic_input` | EXC-01 | 지진파 시계열 파싱 |
| `stl_import` | STL-01 | STL 파일 입력 (CAD 메시 → DualSPHysics) |

참고: `geometry.go`는 독립 도구가 아닌 `xml_generator` 내부에서 사용하는 helper 모듈.

### pv-agent MCP Server (후처리·렌더링)

렌더링/애니메이션은 pv-agent MCP 서버로 위임 (기존 `pvpython.go`, `animation.go` 삭제).

설정: `.mcp.json` → `pv` 서버 (Python uv, `tools/pv-agent/`)

| MCP Tool | 역할 |
|----------|------|
| `render` | ParaView 스냅샷 렌더링 (PNG) |
| `animate` | 애니메이션 생성 (MP4) |
| `pv_isosurface` | IsoSurface 등위면 메시 생성 |
| `slice` / `clip` / `contour` | 후처리 필터 |
| `inspect_data` | VTK 데이터 검사 |
| `extract_stats` | 통계 추출 |
| `plot_over_line` / `streamlines` | 라인 플롯 / 유선 |

### Database

SQLite3 with goose migrations (`internal/db/migrations/`). 3개 테이블:
- `sessions` — 대화 세션
- `messages` — 메시지 (JSON parts 배열)
- `files` — 파일 버전 히스토리

SQL 쿼리는 `internal/db/sql/`에 작성 → `sqlc generate`로 Go 코드 자동 생성.

## DualSPHysics Integration

### Critical Rules

1. **GenCase는 `.xml` 자동 추가** — 경로에 `.xml` 확장자 포함하지 않는다
2. **GenCase는 positional args** — `gencase <case_path> <save_path>` (`-save:` 문법 아님)
3. **XML은 attribute-only 문법** — `<gravity x="0" y="0" z="-9.81" />` (자식 노드 아님)
4. **Docker 바이너리명은 소문자** — `gencase`, `dualsphysics`, `partvtk`, `measuretool` (대문자 아님)
5. **GPU 타겟**: NVIDIA RTX 4090 (sm_89, Ada Lovelace), CUDA 12.6

### Case Files

`cases/` 디렉토리에 사전 정의 슬로싱 시나리오 (논문 벤치마크 기반):
- `SloshingTank_Def.xml` — 기본 벤치마크
- `Sloshing_{Normal,NearRes,Res}[_Guard]_Def.xml` — 주파수별/가드 유무 변형
- `SPHERIC_Test10_{Low,High}_Def.xml` — SPHERIC 벤치마크 (Low/High res)
- `Chen2018_*.xml`, `Liu2024_*.xml`, `English2021_*.xml` 등 — 논문 재현 케이스
- `*.stl` — STL 형상 파일 (fuel_tank.stl, horiz_cylinder.stl)

### Docker Pipeline

```
Dockerfile:          CUDA 12.6 devel → MoorDynPlus 빌드 → DualSPHysics GPU 빌드 → runtime 이미지
Dockerfile.pvpython: ParaView pvpython 전용 이미지 (headless 렌더링)
docker-compose.yml:  NVIDIA GPU reservation, ./simulations:/data, ./cases:/cases 마운트
```

### Solver Toolchain

Docker image (`dsph-agent:latest`)에서 모든 바이너리는 소문자 symlink:

| Docker Binary | 역할 |
|---------------|------|
| `gencase` | XML → 바이너리 파티클 전처리 |
| `dualsphysics` | GPU SPH 솔버 |
| `partvtk` | VTK 시각화 데이터 추출 |
| `measuretool` | 측정점 데이터 추출 |
| `isosurface` | 등위면 메시 생성 |

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
├── wizard/      → Case Wizard (huh 폼 기반 시뮬레이션 설정)
└── util/        → 유틸리티 컴포넌트
```

BubbleTea 패턴: Model struct + Init/Update/View, Lipgloss 스타일링, pubsub 이벤트 연동.

## Research Pipeline (research-sniper v2)

논문 연구는 5단계 전술 파이프라인으로 관리한다.

```
0 목표선정(SCOPE)    — 연구 프로파일 정의        ✅ 완료 (research/profile.json)
1 정찰(RECON)        — 논문 수집/분석             ⬜ 대기
2 영점조정(ZERO)     — Gap reasoning              ⬜ 대기
3 조준(AIM)          — 연구 결론 도출             ⬜ 대기
4 발사(SEND)         — BibTeX + 연구 계획 출력    ⬜ 대기
```

### 연구 디렉토리 구조

```
research/
├── profile.json          — 연구 프로파일 (SCOPE 단계 산출물)
├── progress.json         — 파이프라인 진행 상황
├── survey_analysis.md    — RECON 단계 서베이 분석 (자동 생성)
├── research_plan.md      — SEND 단계 실험 TODO (자동 생성)
├── paper_skeleton.md     — 논문 뼈대 (자동 생성)
└── references.bib        — 수집된 논문 BibTeX (자동 생성)
```

### 기존 데이터 자산 (연구에 활용)

| 자산 | 위치 | 내용 |
|------|------|------|
| **논문 PDF 7편** | `datasets/papers/` | Chen2018, Liu2024, English2021, ISOPE2012, NASA2023, Frosina2018, Zhao2024 |
| **SPHERIC 실험 데이터** | `datasets/spheric/case_{1,2,3}/` | Test 10: 100회 반복, 2종 유체 (water + oil), 압력 시계열 |
| **XML 케이스 20개** | `cases/` | 벤치마크 재현 케이스 (SPHERIC, Chen, Liu, English 등) |
| **시뮬레이션 결과 17개** | `simulations/` | GPU 실행 완료 (8 PASS, 2 PARTIAL) |
| **탱크 파라미터** | `datasets/TANK_GEOMETRIES_FROM_PAPERS.md` | 8개 논문 상세 치수 + 10개 JSON 프리셋 |
| **STL/데이터셋 조사** | `docs/STL_DATASETS_AND_PARAMETRICS.md` | 535줄, Tier 1/2/3 분류 |

### 핵심 참고 문헌

- English et al. (2021) — mDBC sloshing validation with DualSPHysics
- Chen et al. (2018) — OpenFOAM sloshing parametric study
- Liu et al. (2024) — Large tank pitch sloshing
- Bran et al. (2024) — ChemCrow: LLM agent for chemistry
- Boiko et al. (2023) — Coscientist: AI-driven scientific discovery
- Schick et al. (2024) — Toolformer: LLMs learn to use tools

## Custom Agents & Commands

### Sub-agents (`.claude/agents/`)

| Agent | 역할 |
|-------|------|
| `tui-test-runner` | TUI 컴포넌트 편집 후 자동 테스트 |
| `scenario-validator` | BDD 시나리오 vs 구현 검증 |
| `dsph-xml-validator` | DualSPHysics XML 케이스 검증 |
| `prompt-optimizer` | Qwen3 SLM 프롬프트 최적화 |
| **`survey-orchestrator`** | **논문 서베이 전체 조율 (Opus, 3 approval gates)** |
| **`survey-paper-searcher`** | **학술 DB 논문 검색 워커 (Sonnet)** |
| **`survey-citation-chaser`** | **인용 체인 추적 워커 (Sonnet)** |
| **`survey-dataset-searcher`** | **데이터셋 검색 워커 (Sonnet)** |

### Slash commands (`.claude/commands/`)

| Command | 역할 |
|---------|------|
| `/bubbletea-component <name>` | BubbleTea 컴포넌트 스캐폴딩 |
| `/write-scenario <feature>` | BDD 시나리오 자동 생성 |
| `/implement-tool <name>` | DualSPHysics 도구 구현 템플릿 |
| **`/research [resume\|status]`** | **연구 파이프라인 대화형 실행** |

## CI/CD & Release

- `.github/workflows/build.yml` — main 브랜치 푸시 시 빌드
- `.github/workflows/release.yml` — 태그 푸시 시 goreleaser 릴리스
- `.goreleaser.yml` — 멀티플랫폼 빌드 (linux/darwin × amd64/arm64), AUR, Homebrew tap, deb/rpm

## Configuration

설정 파일: `.opencode/config.json` (앱 설정), `.mcp.json` (MCP 서버 설정)

핵심 설정 항목:
- `agents.coder.model` — 주 에이전트 모델
- `agents.task.model` — 백그라운드 작업 모델
- `providers.*` — LLM 프로바이더 API 키/설정
- `mcpServers.*` — MCP 서버 등록 (stdio/SSE)
- `lsp.*` — LSP 서버 설정 (gopls, pyright 등)
- `contextPaths` — 자동 로딩할 컨텍스트 파일 목록

MCP 서버 (`.mcp.json`):
- `context7` — 라이브러리 문서 검색 (npx @upstash/context7-mcp)
- `pv` — ParaView 후처리 서버 (uv, `tools/pv-agent/`)
- `research-sniper` — 논문 서베이 파이프라인 (survey, gap analysis, export)

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
main                  ← 제품 코드 (개발)
research/paper        ← 논문 연구 (현재 브랜치)
```

`research/paper`는 main의 코드 자산을 기반으로 연구 관련 파일만 추가한다. 코드 변경은 main에서, 연구는 이 브랜치에서 진행.
