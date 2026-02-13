# slosim-agent 개발 체계 v2.0

*2026-02-13 — 전면 재설계*

---

## 0. 핵심 구분

| | 제품 (Product) | 개발팀 (Dev Team) |
|---|---|---|
| **정체** | slosim-agent (OpenCode fork + Qwen3 SLM) | Claude Code (Anthropic) 인스턴스들 |
| **모델** | Qwen3 32B/Coder (오픈웨이트, 로컬) | Claude Sonnet/Opus (Anthropic API) |
| **실행** | 호스트에서 직접 실행 (Go TUI) | 호스트 tmux에서 CC Agent Team |

---

## 1. 제품 아키텍처: slosim-agent

OpenCode(kimimgo/opencode-custom) 포크 → 슬로싱 전문 에이전트로 커스텀

```
┌─────────────────────────────────────────────┐
│         slosim-agent (Go + BubbleTea)       │
│                                             │
│  ┌─────────────────────────────────────┐    │
│  │  TUI Layer (슬로싱 특화 UI)         │    │
│  │  - Chat + Sim Dashboard             │    │
│  │  - Case Wizard / Result Viewer      │    │
│  │  - Parametric Comparison View       │    │
│  └──────────────┬──────────────────────┘    │
│                 │ pubsub                     │
│  ┌──────────────▼──────────────────────┐    │
│  │  Agent Core (Sloshing Coder)        │    │
│  │  - LLM: Qwen3 via Ollama           │    │
│  │  - Sloshing domain prompt           │    │
│  │  - Tool orchestration (ReAct)       │    │
│  └──────────────┬──────────────────────┘    │
│                 │                            │
│  ┌──────────────▼──────────────────────┐    │
│  │  Tools (DualSPHysics Pipeline)      │    │
│  │  gencase / solver / partvtk /       │    │
│  │  measuretool / pvpython /           │    │
│  │  job_manager / report_generator     │    │
│  └──────────────┬──────────────────────┘    │
│                 │ Docker exec                │
│  ┌──────────────▼──────────────────────┐    │
│  │  DualSPHysics v5.4 (CUDA 12.6)     │    │
│  │  ParaView (pvpython)                │    │
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

---

## 2. CC 개발팀 구성

CC의 **Agent Teams** 기능 활용 (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1)

### 팀 구조

```
prj-slosim-main (Lead CC)
├── Teammate 1: tui-dev      ← TUI 개발 담당
├── Teammate 2: qa-tester    ← UX 기획 + 제품 테스터
└── Teammate 3: agent-eng    ← SLM 에이전트 엔지니어
```

**Lead (오케스트레이터):**
- Delegate mode (Shift+Tab) — 코드 직접 수정 금지, 조율만
- 태스크 분배, 의존성 관리, PR 머지, 품질 게이트
- Split-pane mode (tmux 기반)

### 역할별 상세

---

### 🖥️ Teammate 1: `tui-dev` — TUI 개발

**담당 영역:** `internal/tui/`, `internal/app/`
**worktree 브랜치:** `feat/tui`

**하는 일:**
- BubbleTea 컴포넌트 개발 (Sim Dashboard, Case Wizard, Result Viewer)
- Lipgloss 테마/스타일링 (슬로싱 브랜딩)
- pubsub 이벤트 연결 (Agent Core ↔ TUI)
- 이미지/애니메이션 인라인 렌더링

**CC 플러그인:**

| 플러그인 | 출처 | 용도 |
|---------|------|------|
| `gopls-lsp` | Official marketplace | Go 코드 인텔리전스 (타입 에러 즉시 감지) |
| `github` | Official marketplace | PR 생성, 리뷰 자동화 |
| `commit-commands` | Official marketplace | Git 커밋 워크플로우 |

**MCP 서버:**

| MCP | 용도 |
|-----|------|
| `context7` | BubbleTea, Lipgloss, Bubbles 최신 API 참조 |

**커스텀 Sub-agent:**
```yaml
# .claude/agents/tui-test-runner.md
---
name: tui-test-runner
description: Run TUI component tests after every edit
tools: Bash, Read, Glob
model: haiku
---
After code changes in internal/tui/, run:
go test ./internal/tui/... -v
Report failures concisely with file:line and fix suggestions.
```

**커스텀 Skill:**
```
# .claude/commands/bubbletea-component.md
Create a new BubbleTea component following the project pattern:
- Model struct with Init/Update/View
- Lipgloss adaptive styling
- pubsub event subscription
- teatest unit test
Component name: $ARGUMENTS
```

---

### 🧪 Teammate 2: `qa-tester` — UX 기획 & 제품 테스터

**담당 영역:** `tests/`, `docs/scenarios/`
**worktree 브랜치:** `feat/qa`
**권한:** Read-only (코드 수정 금지, 테스트/문서만 작성)

**하는 일:**
- BDD 시나리오 설계 (비전문가 관점)
- E2E 테스트 작성 + 실행
- UX 플로우 검증 (시나리오대로 구현됐는지)
- CFD 용어 → 쉬운 말 변환 검증
- 버그 발견 시 GitHub Issue 생성

**CC 플러그인:**

| 플러그인 | 출처 | 용도 |
|---------|------|------|
| `gopls-lsp` | Official | Go 테스트 코드 인텔리전스 |
| `github` | Official | Issue 생성, PR 리뷰 |
| `pr-review-toolkit` | Official | PR 품질 리뷰 자동화 |

**MCP 서버:**

| MCP | 용도 |
|-----|------|
| `context7` | Go testing, testify 문서 참조 |

**커스텀 Sub-agent:**
```yaml
# .claude/agents/scenario-validator.md
---
name: scenario-validator
description: Validate UX scenarios against implementation
tools: Read, Bash, Glob, Grep
model: sonnet
---
You are a product tester who has NO CFD knowledge.
Read docs/scenarios/*.feature files.
For each scenario, verify the implementation handles all steps.
Flag any error message or UI text that contains CFD jargon
without explanation. Report as GitHub issues.
```

**커스텀 Skill:**
```
# .claude/commands/write-scenario.md
Write a BDD scenario for slosim-agent from a NON-EXPERT user perspective.
Rules:
- User has zero CFD knowledge
- All interactions are natural language Korean
- Include error cases (invalid input, simulation failure)
- Verify error messages are understandable
Feature: $ARGUMENTS
```

---

### 🤖 Teammate 3: `agent-eng` — SLM 에이전트 엔지니어

**담당 영역:** `internal/llm/`, `internal/tools/` (신규 DualSPHysics tools)
**worktree 브랜치:** `feat/agent-core`

**하는 일:**
- DualSPHysics Tool 구현 (GenCase, Solver, PartVTK, MeasureTool, pvpython)
- Qwen3 전용 시스템 프롬프트 설계 (슬로싱 도메인 지식)
- Agent Loop 커스텀 (sloshing-coder agent type)
- Job Manager (백그라운드 실행 + 모니터링)
- 파라메트릭 스터디 오케스트레이션
- 리포트 생성기

**CC 플러그인:**

| 플러그인 | 출처 | 용도 |
|---------|------|------|
| `gopls-lsp` | Official | Go 코드 인텔리전스 |
| `github` | Official | PR, 코드 리뷰 |
| `commit-commands` | Official | 커밋 워크플로우 |

**MCP 서버:**

| MCP | 용도 |
|-----|------|
| `context7` | DualSPHysics docs, Ollama Go SDK, Go concurrency 패턴 |
| `docker-mcp` | Docker 컨테이너 제어 (DualSPHysics 빌드/실행/로그) |

**커스텀 Sub-agent:**
```yaml
# .claude/agents/dsph-xml-validator.md
---
name: dsph-xml-validator
description: Validate DualSPHysics XML case files
tools: Read, Bash
model: haiku
---
Validate XML files for DualSPHysics GenCase:
1. All values must be in attributes (never text content)
2. File path must NOT end in .xml (GenCase auto-appends)
3. Required sections: casedef, execution, geometry
4. Gravity, dp, kernel, viscotreatment must be present
Run: xmllint --noout <file> for well-formedness check.
```

```yaml
# .claude/agents/prompt-optimizer.md
---
name: prompt-optimizer
description: Optimize system prompts for Qwen3 SLM
tools: Read, Write, Bash
model: sonnet
---
You are an SLM prompt engineering specialist for Qwen3 (32B/8B).
Constraints:
- Total system prompt must fit in 8K tokens
- Use structured output (JSON mode) for tool calls
- Few-shot examples must be token-efficient
- Korean+English domain terms mixed
- Chain-of-Thought for physics reasoning
Test prompts against Ollama locally before finalizing.
```

**커스텀 Skill:**
```
# .claude/commands/implement-tool.md
Implement a new DualSPHysics tool following the OpenCode Tool interface:
1. Create internal/llm/tools/<name>.go
2. Implement BaseTool.Info() with name, description, input schema
3. Implement Execute(ctx, input) -> ToolResponse
4. All execution via Docker (docker compose run --rm dsph ...)
5. Create internal/llm/tools/<name>_test.go (TDD: test first)
6. Register in tools/tools.go
Tool name: $ARGUMENTS
```

---

## 3. Git Worktree + TDD 협업 체계

### 브랜치 전략

```
main (protected — QA 통과 후 머지만)
├── develop          ← 통합 브랜치
├── feat/tui         ← tui-dev worktree
├── feat/agent-core  ← agent-eng worktree
└── feat/qa          ← qa-tester worktree
```

### 셋업 명령

```bash
cd ~/workspace/02_active/slosim-agent

# 브랜치 생성
git branch develop
git branch feat/tui
git branch feat/agent-core
git branch feat/qa

# Worktree 생성 (각 Teammate용)
git worktree add ../slosim-wt-tui feat/tui
git worktree add ../slosim-wt-agent feat/agent-core
git worktree add ../slosim-wt-qa feat/qa
```

### TDD 흐름

```
1. qa-tester: 시나리오 작성 → 실패하는 E2E 테스트 작성 (RED)
2. agent-eng: Tool 유닛 테스트 작성 → 구현 (RED → GREEN)
3. tui-dev: UI 컴포넌트 테스트 작성 → 구현 (RED → GREEN)
4. Lead: develop 머지 → 통합 테스트
5. qa-tester: E2E 테스트 통과 확인 (GREEN)
6. Lead: main 머지
```

### 의존성 그래프 (빌드 순서)

```
Phase 1 (병렬):
  agent-eng: Tools (gencase, solver, partvtk, measuretool, pvpython)
  qa-tester: 시나리오 문서 + 실패 테스트
  tui-dev: 기존 OpenCode TUI 분석 + 테마 프로토타입

Phase 2 (agent-eng 선행):
  agent-eng: Job Manager + Result Store + Agent Loop
  tui-dev: Sim Dashboard + Case Wizard (Agent Core 이벤트 구독)
  qa-tester: Tool 유닛 테스트 리뷰

Phase 3 (통합):
  Lead: develop 머지
  tui-dev: Result Viewer + Parametric View
  agent-eng: Report Generator + 프롬프트 최적화
  qa-tester: E2E 시나리오 검증

Phase 4 (폴리시):
  tui-dev: UX 개선 (qa-tester 피드백 반영)
  agent-eng: Qwen3 32b vs 8b 벤치마크
  qa-tester: 전체 회귀 테스트 + 논문 데모 검증
```

### Hooks (품질 게이트)

```json
// .claude/settings.json
{
  "hooks": {
    "TeammateIdle": {
      "command": "go test ./... 2>&1 | tail -5",
      "description": "Run tests when teammate finishes",
      "exitCode2Action": "block"
    },
    "PostEdit": {
      "command": "go vet ./...",
      "description": "Vet after every edit"
    }
  }
}
```

---

## 4. Agent Team 실행 방법

### 환경 설정

```bash
# CC settings (~/.claude/settings.json)
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  },
  "teammateMode": "tmux"
}
```

### 팀 시작 프롬프트

```
slosim-agent 프로젝트를 개발할 Agent Team을 구성해줘:

1. tui-dev: BubbleTea TUI 개발 담당.
   - worktree: ~/workspace/02_active/slosim-wt-tui (feat/tui)
   - 영역: internal/tui/, internal/app/
   - 플러그인: gopls-lsp, github, commit-commands
   - MCP: context7

2. qa-tester: UX 시나리오 설계 + E2E 테스트 전담.
   - worktree: ~/workspace/02_active/slosim-wt-qa (feat/qa)
   - 영역: tests/, docs/scenarios/ (코드 수정 금지)
   - 플러그인: gopls-lsp, github, pr-review-toolkit
   - MCP: context7

3. agent-eng: SLM 에이전트 + DualSPHysics Tool 개발.
   - worktree: ~/workspace/02_active/slosim-wt-agent (feat/agent-core)
   - 영역: internal/llm/, internal/tools/
   - 플러그인: gopls-lsp, github, commit-commands
   - MCP: context7, docker-mcp

나(Lead)는 delegate mode로 조율만 한다.
TDD: qa-tester가 먼저 실패 테스트 작성 → 개발자들이 구현.
각 teammate는 plan approval 필수.
```

---

## 5. 플러그인/스킬 요약 매트릭스

| | tui-dev | qa-tester | agent-eng |
|---|---|---|---|
| **gopls-lsp** | ✅ | ✅ | ✅ |
| **github** | ✅ | ✅ (Issue 생성) | ✅ |
| **commit-commands** | ✅ | — | ✅ |
| **pr-review-toolkit** | — | ✅ | — |
| **Context7 MCP** | ✅ (BubbleTea) | ✅ (testing) | ✅ (DSPH, Ollama) |
| **docker-mcp** | — | — | ✅ |
| **Sub-agent: test-runner** | ✅ | — | — |
| **Sub-agent: scenario-validator** | — | ✅ | — |
| **Sub-agent: xml-validator** | — | — | ✅ |
| **Sub-agent: prompt-optimizer** | — | — | ✅ |
| **Skill: /bubbletea-component** | ✅ | — | — |
| **Skill: /write-scenario** | — | ✅ | — |
| **Skill: /implement-tool** | — | — | ✅ |
| **코드 수정 권한** | tui/ only | 없음 (Read-only) | llm/, tools/ only |

---

## 6. 파일 배치 계획

```
slosim-agent/
├── .claude/
│   ├── settings.json          ← hooks, team config
│   ├── agents/
│   │   ├── tui-test-runner.md
│   │   ├── scenario-validator.md
│   │   ├── dsph-xml-validator.md
│   │   └── prompt-optimizer.md
│   └── commands/
│       ├── bubbletea-component.md
│       ├── write-scenario.md
│       └── implement-tool.md
├── .mcp.json                  ← context7, docker-mcp
├── CLAUDE.md                  ← 전체 프로젝트 컨텍스트
├── PRD.md
├── ARCHITECTURE.md
├── internal/
│   ├── tui/                   ← tui-dev 영역
│   ├── llm/
│   │   ├── agent/             ← agent-eng 영역
│   │   ├── prompt/            ← agent-eng 영역
│   │   └── tools/             ← agent-eng 영역 (DSPH tools)
│   └── ...
├── tests/                     ← qa-tester 영역
├── docs/scenarios/            ← qa-tester 영역
├── cases/                     ← XML 케이스 템플릿
├── Dockerfile
└── docker-compose.yml
```
