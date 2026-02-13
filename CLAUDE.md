# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Agent Identity

- **이 에이전트 이름**: **cc-slosim-1**
- **역할**: slosim-agent 개발 에이전트
- **응답 시 이름 표시**: 매 응답 첫 줄에 `[cc-slosim-1]` 태그 포함

## Project Overview

slosim-agent는 자연어로 슬로싱(sloshing) 시뮬레이션을 설정·실행·분석하는 AI 에이전트다. OpenCode(Go + BubbleTea TUI) 기반 위에 DualSPHysics v5.4 GPU solver를 통합했다.

**핵심 흐름**: 자연어 입력 → DualSPHysics XML 케이스 생성 → GPU 시뮬레이션 실행 → AI 분석

## Build & Development Commands

```bash
# Build
go build -o slosim-agent ./main.go

# Build with version info
go build -ldflags "-s -w -X github.com/opencode-ai/opencode/internal/version.Version=dev" ./main.go

# Run tests
go test ./...

# Run a single test
go test ./internal/llm/tools/ -run TestLs

# Generate DB code from SQL (requires sqlc)
sqlc generate

# Generate config JSON schema
go run ./cmd/schema/main.go > opencode-schema.json

# Docker: build DualSPHysics GPU image
docker compose build

# Docker: run simulation
docker compose run --rm dsph GenCase /cases/SloshingTank_Def -save:/data/out
```

## Architecture

### Module Path & Go Version

- Module: `github.com/opencode-ai/opencode`
- Go: 1.24.0
- CGO: disabled for release builds (`CGO_ENABLED=0`)

### Core Layers

```
cmd/root.go          → Cobra CLI entry (interactive TUI / non-interactive -p mode)
internal/app/        → App orchestrator (wires all services together)
internal/config/     → Viper-based config loader (.opencode/config.json)
internal/llm/
  ├── agent/         → Agent loop: coder, task, title, summarizer
  ├── provider/      → 10 LLM providers (Anthropic, OpenAI, Gemini, Groq, etc.)
  ├── models/        → Provider-specific model definitions
  ├── prompt/        → System prompts per agent type
  └── tools/         → Built-in tools (bash, edit, glob, grep, view, write, patch, fetch, ls)
internal/tui/        → BubbleTea TUI (chat, dialogs, themes)
internal/db/         → SQLite3 + sqlc generated code
internal/session/    → Session CRUD
internal/message/    → Message CRUD
internal/pubsub/     → Event broadcast system
internal/permission/ → Tool execution permission checks
internal/lsp/        → LSP client (gopls, pyright integration)
```

### Key Patterns

- **Service interface + pub/sub**: 모든 서비스는 `Service` interface를 구현하고, `pubsub.Subscriber[T]`로 상태 변경을 브로드캐스트
- **Tool interface**: `BaseTool.Info()` + `Execute(ctx, input) → ToolResponse` — 모든 built-in 도구와 MCP 도구가 동일 인터페이스
- **Provider abstraction**: `Provider.SendMessages()` / `StreamResponse()` — 10개 LLM 프로바이더가 동일 인터페이스
- **Agent composition**: coder는 전체 도구 접근, task는 제한된 도구셋, title/summarizer는 도구 없음

### Database

SQLite3 with goose migrations (`internal/db/migrations/`). 3개 테이블:
- `sessions` — 대화 세션
- `messages` — 메시지 (JSON parts 배열)
- `files` — 파일 버전 히스토리

SQL 쿼리는 `internal/db/sql/`에 작성 → `sqlc generate`로 Go 코드 자동 생성.

### MCP Integration

`internal/llm/agent/mcp-tools.go`에서 MCP 서버를 stdio/SSE 방식으로 연결하고, 도구를 동적 발견하여 에이전트에 등록.

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
Dockerfile: CUDA 12.6 devel → MoorDynPlus 빌드 → DualSPHysics GPU 빌드 → runtime 이미지
docker-compose.yml: NVIDIA GPU reservation, ./simulations:/data, ./cases:/cases 마운트
```

### Solver Toolchain

| Binary | 역할 |
|--------|------|
| GenCase | XML → 바이너리 파티클 전처리 |
| DualSPHysics5.4_linux64 | GPU SPH 솔버 |
| PartVTK | VTK 시각화 데이터 추출 |
| MeasureTool | 측정점 데이터 추출 |
| IsoSurface | 등위면 메시 생성 |

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

기존 테스트 파일:
- `internal/llm/prompt/prompt_test.go` — 프롬프트 생성
- `internal/llm/tools/ls_test.go` — 디렉토리 리스팅 (가장 상세)
- `internal/tui/components/dialog/custom_commands_test.go` — 커스텀 명령어
- `internal/tui/theme/theme_test.go` — 테마 시스템
