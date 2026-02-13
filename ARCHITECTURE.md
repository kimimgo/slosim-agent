# slosim-agent 개발 체계 v3.0

*2026-02-13 — CC Agent Teams 기반 전면 재설계*

---

## 0. 핵심 구분

| | 제품 (Product) | 개발팀 (Dev Team) |
|---|---|---|
| **정체** | slosim-agent (OpenCode fork + Qwen3 SLM) | Claude Code Agent Teams |
| **모델** | Qwen3 32B/Coder (오픈웨이트, 로컬 Ollama) | Claude Sonnet/Opus (Anthropic API) |
| **실행** | 호스트에서 직접 실행 (Go TUI) | CC TeamCreate → Task(team_name) → SendMessage |
| **환경변수** | — | `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` |

---

## 1. 제품 아키텍처: slosim-agent

OpenCode(`kimimgo/opencode-custom`) 포크 → 슬로싱 전문 에이전트로 커스텀

```
┌─────────────────────────────────────────────────┐
│           slosim-agent (Go + BubbleTea)         │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  TUI Layer (슬로싱 특화 UI)               │  │
│  │  ┌───────────┐ ┌────────────┐ ┌────────┐  │  │
│  │  │ Chat View │ │ Sim Dash   │ │ Result │  │  │
│  │  │           │ │ - Job List │ │ Viewer │  │  │
│  │  │           │ │ - Progress │ │ - VTK  │  │  │
│  │  │           │ │ - Logs     │ │ - Plot │  │  │
│  │  └───────────┘ └────────────┘ └────────┘  │  │
│  │  ┌────────────────┐ ┌──────────────────┐  │  │
│  │  │ Case Wizard    │ │ Parametric View  │  │  │
│  │  │ (Step-by-step) │ │ (Compare/Table)  │  │  │
│  │  └────────────────┘ └──────────────────┘  │  │
│  └──────────────────┬────────────────────────┘  │
│                     │ pubsub events              │
│  ┌──────────────────▼────────────────────────┐  │
│  │  Agent Core (Sloshing Coder)              │  │
│  │  - LLM: Qwen3 via Ollama (local)         │  │
│  │  - Sloshing domain system prompt          │  │
│  │  - Tool orchestration (ReAct loop)        │  │
│  │  - Job Manager (background goroutine)     │  │
│  └──────────────────┬────────────────────────┘  │
│                     │                            │
│  ┌──────────────────▼────────────────────────┐  │
│  │  Tools (DualSPHysics Pipeline)            │  │
│  │  ┌────────┐ ┌────────┐ ┌───────────────┐  │  │
│  │  │gencase │ │solver  │ │partvtk        │  │  │
│  │  │        │ │        │ │measuretool    │  │  │
│  │  └────────┘ └────────┘ │isosurface     │  │  │
│  │  ┌────────┐ ┌────────┐ └───────────────┘  │  │
│  │  │pvpython│ │report  │                    │  │
│  │  │ParaView│ │generat.│                    │  │
│  │  └────────┘ └────────┘                    │  │
│  └──────────────────┬────────────────────────┘  │
│                     │ docker exec / compose      │
│  ┌──────────────────▼────────────────────────┐  │
│  │  DualSPHysics v5.4 (CUDA 12.6, RTX 4090) │  │
│  │  ParaView pvpython (headless rendering)   │  │
│  └───────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

### 코드 레이아웃

```
internal/
├── tui/                          # TUI Layer (tui-dev 담당)
│   ├── components/
│   │   ├── chat/                 # 기존 Chat View
│   │   ├── dialog/               # 기존 Dialogs
│   │   ├── sim/                  # [신규] Sim Dashboard, Job List
│   │   ├── wizard/               # [신규] Case Wizard (step-by-step)
│   │   ├── result/               # [신규] Result Viewer
│   │   └── parametric/           # [신규] Parametric Comparison View
│   ├── theme/                    # 슬로싱 브랜드 테마
│   └── styles/                   # Lipgloss adaptive styles
├── llm/                          # Agent Core (agent-eng 담당)
│   ├── agent/                    # Agent loop + sloshing-coder type
│   ├── provider/                 # LLM providers (Ollama/Qwen3 중심)
│   ├── prompt/                   # System prompts (sloshing domain)
│   └── tools/                    # Built-in + DualSPHysics tools
│       ├── bash.go, edit.go, ... # 기존 OpenCode tools
│       ├── gencase.go            # [신규] GenCase XML → particles
│       ├── solver.go             # [신규] DualSPHysics GPU solver
│       ├── partvtk.go            # [신규] VTK export
│       ├── measuretool.go        # [신규] Measurement extraction
│       ├── pvpython.go           # [신규] ParaView rendering
│       ├── job_manager.go        # [신규] Background job orchestration
│       └── report.go             # [신규] Report generation
├── app/                          # App orchestrator
├── config/                       # Viper config (.opencode/config.json)
├── db/                           # SQLite3 + sqlc
├── session/                      # Session CRUD
├── message/                      # Message CRUD
├── pubsub/                       # Event broker
├── permission/                   # Tool permission checks
└── lsp/                          # LSP client (gopls)

tests/                            # E2E 테스트 (qa-tester 담당)
├── e2e/
│   ├── scenario_simple_sloshing_test.go
│   ├── scenario_parametric_test.go
│   └── scenario_stl_import_test.go
└── fixtures/
    ├── simple_tank.xml
    └── sample.stl

docs/scenarios/                   # BDD 시나리오 (qa-tester 담당)
├── simple_sloshing.feature
├── parametric_study.feature
└── error_handling.feature
```

---

## 2. CC Agent Team 구성

### Team API 사용법

```
# 1. Lead가 팀 생성
TeamCreate(team_name="slosim", description="slosim-agent 개발팀")

# 2. Teammate 스폰 (Task tool with team_name)
Task(
  subagent_type="general-purpose",
  team_name="slosim",
  name="tui-dev",
  mode="plan",           # plan approval 필수
  prompt="TUI 개발 담당. ARCHITECTURE.md의 tui-dev 역할 참조..."
)

# 3. 태스크 분배
TaskCreate(subject="Sim Dashboard 컴포넌트 구현", ...)
TaskUpdate(taskId="1", owner="tui-dev")

# 4. 커뮤니케이션
SendMessage(type="message", recipient="tui-dev", content="...", summary="...")
SendMessage(type="broadcast", content="develop 머지 완료", summary="...")
```

### 팀 구조

```
prj-slosim (Lead CC) ─── delegate mode, 코드 직접 수정 금지
├── tui-dev          ─── BubbleTea TUI 개발 (plan mode)
├── qa-tester        ─── UX 기획 + 제품 QA (plan mode)
└── agent-eng        ─── SLM 에이전트 엔지니어 (plan mode)
```

### Teammate 상세

---

### Teammate 1: `tui-dev` — TUI 개발

**담당 영역:** `internal/tui/`, `internal/app/`
**worktree 브랜치:** `feat/tui`
**mode:** `plan` (코드 변경 전 Lead 승인 필수)

**하는 일:**
- BubbleTea 컴포넌트 개발 (Sim Dashboard, Case Wizard, Result Viewer, Parametric View)
- Lipgloss 테마/스타일링 (슬로싱 브랜딩, AdaptiveColor)
- pubsub 이벤트 연결 (Agent Core ↔ TUI)
- 이미지/애니메이션 인라인 렌더링 (Sixel/Kitty 프로토콜)

**BubbleTea 컴포넌트 패턴:**
```go
// Model struct with Init/Update/View
type SimDashboard struct {
    jobs     []JobStatus
    spinner  spinner.Model
    viewport viewport.Model
    width    int
    height   int
}

func (m SimDashboard) Init() tea.Cmd {
    return tea.Batch(m.spinner.Tick, m.subscribeJobEvents())
}

func (m SimDashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
    case JobUpdateMsg:
        m.jobs = updateJobList(m.jobs, msg)
    }
    // 컴포넌트 위임
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
}

func (m SimDashboard) View() string {
    return lipgloss.JoinVertical(lipgloss.Left,
        m.renderHeader(),
        m.renderJobList(),
        m.spinner.View(),
    )
}
```

**CC 플러그인:** gopls-lsp, github, commit-commands
**MCP 서버:** context7 (BubbleTea, Lipgloss, Bubbles API 참조)
**Sub-agent:** `tui-test-runner` — 편집 후 자동 테스트 (haiku)
**Slash command:** `/bubbletea-component <name>` — 컴포넌트 스캐폴딩

---

### Teammate 2: `qa-tester` — UX 기획 & 제품 테스터

**담당 영역:** `tests/`, `docs/scenarios/`
**worktree 브랜치:** `feat/qa`
**mode:** `plan` (테스트/시나리오 작성 전 Lead 승인 필수)
**권한:** 소스 코드 Read-only (테스트·문서만 작성)

**하는 일:**
- BDD 시나리오 설계 (CFD 비전문가 관점)
- Go E2E 테스트 작성 + 실행
- UX 플로우 검증 (시나리오대로 구현됐는지)
- CFD 용어 → 쉬운 말 변환 검증
- 버그 발견 시 GitHub Issue 생성

**시나리오 형식:**
```gherkin
# docs/scenarios/simple_sloshing.feature
Feature: 단순 슬로싱 해석

  Scenario: 비전문가가 자연어로 해석 요청
    Given 사용자가 slosim-agent를 실행한다
    When "LNG 탱크 슬로싱 해석해줘"라고 입력한다
    Then AI가 탱크 치수를 제안한다
    And 제안에 CFD 전문 용어가 포함되지 않는다
    And 사용자에게 확인을 요청한다

  Scenario: 시뮬레이션 실패 시 알기 쉬운 에러 메시지
    Given 사용자가 잘못된 조건으로 해석을 요청한다
    When 시뮬레이션이 발산한다
    Then "시뮬레이션이 불안정해졌습니다" 메시지를 보여준다
    And 가능한 원인과 해결 방법을 쉬운 말로 안내한다
```

**CC 플러그인:** gopls-lsp, github, pr-review-toolkit
**MCP 서버:** context7 (Go testing, testify 문서 참조)
**Sub-agent:** `scenario-validator` — 시나리오 vs 구현 검증 (sonnet)
**Slash command:** `/write-scenario <feature>` — BDD 시나리오 자동 생성

---

### Teammate 3: `agent-eng` — SLM 에이전트 엔지니어

**담당 영역:** `internal/llm/` (agent, prompt, tools)
**worktree 브랜치:** `feat/agent-core`
**mode:** `plan` (코드 변경 전 Lead 승인 필수)

**하는 일:**
- DualSPHysics Tool 구현 (GenCase, Solver, PartVTK, MeasureTool, pvpython)
- Qwen3 전용 시스템 프롬프트 설계 (슬로싱 도메인 지식)
- Agent Loop 커스텀 (sloshing-coder agent type)
- Job Manager (백그라운드 실행 + goroutine 모니터링)
- 파라메트릭 스터디 오케스트레이션
- 리포트 생성기

**Tool 구현 패턴:**
```go
// internal/llm/tools/gencase.go
// BaseTool 인터페이스: Info() + Run(ctx, ToolCall) -> (ToolResponse, error)

type genCaseTool struct{}

func (t *genCaseTool) Info() ToolInfo {
    return ToolInfo{
        Name:        "gencase",
        Description: "Generate particle geometry from DualSPHysics XML case definition",
        Parameters: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "case_path": map[string]any{
                    "type":        "string",
                    "description": "Path to XML case file (WITHOUT .xml extension)",
                },
                "save_path": map[string]any{
                    "type":        "string",
                    "description": "Output directory for generated particles",
                },
            },
            "required": []string{"case_path", "save_path"},
        },
    }
}

func (t *genCaseTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
    // docker compose run --rm dsph GenCase <case_path> -save:<save_path>
    // GenCase 자동 .xml 추가 — 경로에 .xml 포함하지 않음
}
```

**CC 플러그인:** gopls-lsp, github, commit-commands
**MCP 서버:** context7 (DualSPHysics docs, Ollama Go SDK), docker-mcp
**Sub-agents:** `dsph-xml-validator` (haiku), `prompt-optimizer` (sonnet)
**Slash command:** `/implement-tool <name>` — DualSPHysics 도구 스캐폴딩

---

## 3. Git Worktree + TDD 협업 체계

### 브랜치 전략

```
main (protected — QA 통과 후 머지만)
├── develop          ← 통합 브랜치 (Lead가 머지)
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

# Worktree 생성 (각 Teammate용 독립 작업 디렉토리)
git worktree add ../slosim-wt-tui feat/tui
git worktree add ../slosim-wt-agent feat/agent-core
git worktree add ../slosim-wt-qa feat/qa
```

### TDD 흐름 (Red → Green → Refactor)

```
1. qa-tester   → 시나리오 작성 → 실패하는 E2E 테스트 작성 (RED)
2. agent-eng   → Tool 유닛 테스트 작성 → 구현 (RED → GREEN)
3. tui-dev     → UI 컴포넌트 teatest 작성 → 구현 (RED → GREEN)
4. Lead        → develop 머지 → `go test ./...` 통합 테스트
5. qa-tester   → E2E 테스트 전체 통과 확인 (GREEN)
6. Lead        → main 머지 (PR + review)
```

### 의존성 그래프 (빌드 순서)

```
Phase 1 — Foundation (모두 병렬)
┌──────────────────────────────────────────────────────────────┐
│ agent-eng                                                    │
│ ├── gencase tool + test                                      │
│ ├── solver tool + test                                       │
│ └── partvtk/measuretool/isosurface tools + tests             │
│                                                              │
│ qa-tester                                                    │
│ ├── docs/scenarios/*.feature (BDD 시나리오 4건)              │
│ └── tests/e2e/scenario_simple_sloshing_test.go (RED)         │
│                                                              │
│ tui-dev                                                      │
│ ├── 기존 OpenCode TUI 구조 분석                              │
│ └── internal/tui/theme/sloshing.go (슬로싱 브랜드 테마)      │
└──────────────────────────────────────────────────────────────┘
          │
          ▼
Phase 2 — Core Integration (agent-eng 선행 의존)
┌──────────────────────────────────────────────────────────────┐
│ agent-eng                                                    │
│ ├── job_manager.go (백그라운드 실행 + pubsub events)         │
│ ├── sloshing-coder agent type + prompt                       │
│ └── pvpython tool + test                                     │
│                                                              │
│ tui-dev (agent-eng의 pubsub events 필요)                     │
│ ├── internal/tui/components/sim/ (Sim Dashboard)             │
│ └── internal/tui/components/wizard/ (Case Wizard)            │
│                                                              │
│ qa-tester                                                    │
│ └── Tool 유닛 테스트 커버리지 리뷰 → Issue 생성              │
└──────────────────────────────────────────────────────────────┘
          │
          ▼
Phase 3 — Feature Complete (통합)
┌──────────────────────────────────────────────────────────────┐
│ Lead: develop 머지 (feat/agent-core + feat/tui)              │
│                                                              │
│ agent-eng                                                    │
│ ├── report.go (Markdown 리포트 생성)                         │
│ ├── 파라메트릭 스터디 오케스트레이션                         │
│ └── Qwen3 프롬프트 최적화 (prompt-optimizer sub-agent)       │
│                                                              │
│ tui-dev                                                      │
│ ├── internal/tui/components/result/ (Result Viewer)          │
│ └── internal/tui/components/parametric/ (Comparison View)    │
│                                                              │
│ qa-tester                                                    │
│ ├── E2E 시나리오 전체 검증                                   │
│ └── CFD 용어 검증 (scenario-validator sub-agent)             │
└──────────────────────────────────────────────────────────────┘
          │
          ▼
Phase 4 — Polish & Benchmark
┌──────────────────────────────────────────────────────────────┐
│ tui-dev: UX 개선 (qa-tester 피드백 반영)                     │
│ agent-eng: Qwen3 32b vs 8b 벤치마크                          │
│ qa-tester: 전체 회귀 테스트 + 논문 데모 검증                 │
│ Lead: main 머지 → v1.0 태그                                  │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Hooks (품질 게이트)

`.claude/settings.json`에 정의. CC의 hook 시스템으로 편집/커밋 시 자동 품질 검증.

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "command": "cd /home/imgyu/workspace/02_active/slosim-agent && go vet ./...",
        "description": "go vet after every code edit"
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Bash",
        "command": "echo \"$TOOL_INPUT\" | grep -qvE '(rm -rf|sudo|/mnt/storage)' || (echo 'BLOCKED: dangerous command' && exit 1)",
        "description": "Block dangerous bash commands"
      }
    ]
  }
}
```

---

## 5. Agent Team 실행 방법

### 1단계: 환경 설정

```bash
# 전역 CC settings에 Agent Teams 활성화 (~/.claude/settings.json)
# "env": { "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1" }

# Git worktree 생성
cd ~/workspace/02_active/slosim-agent
git branch develop && git branch feat/tui && git branch feat/agent-core && git branch feat/qa
git worktree add ../slosim-wt-tui feat/tui
git worktree add ../slosim-wt-agent feat/agent-core
git worktree add ../slosim-wt-qa feat/qa
```

### 2단계: 팀 시작 프롬프트

```
slosim-agent 프로젝트를 개발할 Agent Team을 구성해줘:

1. tui-dev: BubbleTea TUI 개발 담당.
   - worktree: ~/workspace/02_active/slosim-wt-tui (feat/tui)
   - 영역: internal/tui/, internal/app/
   - mode: plan (변경 전 승인 필수)

2. qa-tester: UX 시나리오 설계 + E2E 테스트 전담.
   - worktree: ~/workspace/02_active/slosim-wt-qa (feat/qa)
   - 영역: tests/, docs/scenarios/ (소스 코드 수정 금지)
   - mode: plan

3. agent-eng: SLM 에이전트 + DualSPHysics Tool 개발.
   - worktree: ~/workspace/02_active/slosim-wt-agent (feat/agent-core)
   - 영역: internal/llm/
   - mode: plan

나(Lead)는 delegate mode로 조율만 한다.
TDD: qa-tester가 먼저 실패 테스트 작성 → 개발자들이 구현.
각 teammate는 plan approval 필수.
```

### 3단계: Lead 워크플로우

```
# 태스크 생성 → 분배
TaskCreate(subject="GenCase tool 구현", description="...", activeForm="Implementing GenCase tool")
TaskUpdate(taskId="1", owner="agent-eng")

# 진행 확인
TaskList()

# Teammate에 메시지
SendMessage(type="message", recipient="tui-dev", content="Phase 2 시작. Sim Dashboard부터.", summary="Phase 2 kickoff")

# Plan 승인 (teammate가 ExitPlanMode 호출 시)
SendMessage(type="plan_approval_response", recipient="tui-dev", request_id="...", approve=true)

# 작업 완료 후 셧다운
SendMessage(type="shutdown_request", recipient="tui-dev", content="Phase 완료, 셧다운")
```

---

## 6. 플러그인/스킬 매트릭스

| | tui-dev | qa-tester | agent-eng |
|---|---|---|---|
| **gopls-lsp** | Y | Y | Y |
| **github** | Y | Y (Issue) | Y |
| **commit-commands** | Y | — | Y |
| **pr-review-toolkit** | — | Y | — |
| **Context7 MCP** | Y (BubbleTea) | Y (testing) | Y (DSPH, Ollama) |
| **docker-mcp** | — | — | Y |
| **Agent: tui-test-runner** | Y | — | — |
| **Agent: scenario-validator** | — | Y | — |
| **Agent: dsph-xml-validator** | — | — | Y |
| **Agent: prompt-optimizer** | — | — | Y |
| **Cmd: /bubbletea-component** | Y | — | — |
| **Cmd: /write-scenario** | — | Y | — |
| **Cmd: /implement-tool** | — | — | Y |
| **코드 수정 권한** | tui/ only | Read-only | llm/ only |

---

## 7. 파일 배치

```
slosim-agent/
├── .claude/
│   ├── settings.json               ← hooks, permissions
│   ├── settings.local.json         ← MCP 권한 (gitignore)
│   ├── agents/
│   │   ├── tui-test-runner.md      ← TUI 테스트 자동 실행
│   │   ├── scenario-validator.md   ← 시나리오 vs 구현 검증
│   │   ├── dsph-xml-validator.md   ← XML 케이스 검증
│   │   └── prompt-optimizer.md     ← Qwen3 프롬프트 최적화
│   └── commands/
│       ├── bubbletea-component.md  ← /bubbletea-component <name>
│       ├── write-scenario.md       ← /write-scenario <feature>
│       └── implement-tool.md       ← /implement-tool <name>
├── .mcp.json                       ← context7 MCP 서버
├── CLAUDE.md                       ← CC 프로젝트 컨텍스트
├── AGENTS.md                       ← 제품 에이전트 개요
├── PRD.md                          ← 제품 요구사항
├── ARCHITECTURE.md                 ← 이 파일
├── internal/                       ← Go 소스
├── tests/                          ← E2E 테스트
├── docs/scenarios/                 ← BDD 시나리오
├── cases/                          ← DualSPHysics XML 템플릿
├── Dockerfile                      ← GPU Docker 이미지
└── docker-compose.yml              ← NVIDIA runtime
```
