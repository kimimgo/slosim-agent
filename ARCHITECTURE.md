# slosim-agent 아키텍처 & 개발 체계

*v1.0 — 2026-02-13*

---

## 1. 팀 구성 (3-Agent 체계)

### Agent 1: `cc-slosim-tui` — TUI 개발 담당

**역할:** BubbleTea TUI 커스터마이징, 슬로싱 도메인 특화 UI 컴포넌트 개발
**worktree 브랜치:** `feat/tui`
**작업 디렉토리:** `internal/tui/`, `internal/app/`

**담당 범위:**
- Simulation Dashboard 패널 (Job 상태, 진행률, 결과 미리보기)
- Case Setup Wizard (단계별 시뮬레이션 설정 UI)
- Result Viewer (이미지/애니메이션 인라인 표시)
- Parametric Study 비교 뷰
- 테마/스타일 (슬로싱 도메인 브랜딩)

**MCP/플러그인 설정:**
```json
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  }
}
```
- Context7: BubbleTea, Lipgloss, Glamour 최신 API 문서 참조
- LSP: gopls (Go 자동완성)

**CLAUDE.md 지침:**
- `internal/tui/` 하위만 수정
- 기존 OpenCode TUI 패턴 (Model-Update-View) 준수
- `pubsub` 이벤트로 Agent Core와 통신
- 새 컴포넌트는 반드시 `_test.go` 동반

---

### Agent 2: `cc-slosim-qa` — UX 기획 & QA 담당

**역할:** 사용자 시나리오 설계, E2E 테스트, 제품 품질 검증
**worktree 브랜치:** `feat/qa`
**작업 디렉토리:** `tests/`, `docs/scenarios/`, `internal/tui/` (읽기 전용)

**담당 범위:**
- 사용자 시나리오 문서 작성 (Gherkin/BDD 형식)
- E2E 테스트 스크립트 (비전문가 관점)
- UX 플로우 검증 (시나리오대로 동작하는지)
- 접근성/에러 핸들링 검증
- 리포트 품질 검증

**시나리오 예시:**
```gherkin
Feature: 비전문가 슬로싱 해석
  Scenario: 간단한 자연어로 시뮬레이션 실행
    Given 사용자가 에이전트를 시작한다
    When "2m x 1m 탱크에 물 70% 채우고 0.5Hz 가진" 입력
    Then AI가 해석 조건을 요약하여 확인 요청
    And 사용자가 확인하면 시뮬레이션 Job 제출
    And Dashboard에 진행률 표시
    And 완료 후 리포트 + 시각화 제시

  Scenario: STL 파일로 커스텀 지오메트리
    Given 사용자가 STL 파일 경로를 입력한다
    When AI가 메시 품질을 검증한다
    Then 적절한 dp와 파라미터를 제안
    And 사용자 확인 후 시뮬레이션 진행

  Scenario: 파라메트릭 스터디
    Given 이전 시뮬레이션 결과가 존재한다
    When "충진율 50%, 60%, 70%로 비교" 입력
    Then 3개 Job을 병렬 제출
    And 비교 리포트 (테이블 + 시각화) 생성
```

**MCP/플러그인 설정:**
```json
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  }
}
```
- Context7: Go testing, testify 문서 참조

**CLAUDE.md 지침:**
- 코드 수정 금지 (읽기 + 테스트만)
- `tests/` 디렉토리에 테스트 작성
- 시나리오 위반 발견 시 GitHub Issue 생성
- 매 빌드 후 시나리오 체크리스트 실행

---

### Agent 3: `cc-slosim-agent` — SLM 에이전트 엔지니어

**역할:** Qwen3 전용 프롬프트 엔지니어링, DualSPHysics 도구/스킬 개발, Agent Loop 커스텀
**worktree 브랜치:** `feat/agent-core`
**작업 디렉토리:** `internal/llm/`, `internal/tools/` (신규), `skills/`

**담당 범위:**
- Qwen3 전용 시스템 프롬프트 설계 (슬로싱 도메인 지식 주입)
- DualSPHysics Tool 구현 (GenCase, Solver, PartVTK, MeasureTool)
- pvpython 후처리 Tool 구현
- Agent Loop 커스터마이징 (sloshing-coder agent type)
- 파라메트릭 스터디 오케스트레이션
- Job Manager (백그라운드 실행 + 모니터링)

**Tool 구현 목록:**
```
internal/llm/tools/
├── gencase.go          ← XML 생성 + GenCase 실행
├── dualsphysics.go     ← GPU 솔버 실행 (Docker)
├── partvtk.go          ← 파티클 → VTK 변환
├── measuretool.go      ← 측정 데이터 추출
├── pvpython.go         ← ParaView CLI 후처리
├── job_manager.go      ← 백그라운드 Job 관리
├── result_store.go     ← 결과 이력 관리
└── report_generator.go ← Markdown 리포트 생성
```

**프롬프트 구조:**
```
internal/llm/prompt/
├── sloshing_system.go  ← 슬로싱 도메인 시스템 프롬프트
├── sloshing_tools.go   ← 도구 사용 지침
└── sloshing_report.go  ← 리포트 생성 프롬프트
```

**MCP/플러그인 설정:**
```json
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  }
}
```
- Context7: DualSPHysics docs, Ollama API, Go concurrency 참조

**CLAUDE.md 지침:**
- Tool interface 준수 (`BaseTool.Info()` + `Execute()`)
- 모든 Tool은 Docker 컨테이너 내 실행 (호스트 직접 실행 금지)
- Qwen3 context window 64K 제한 고려
- 프롬프트는 한국어/영어 혼용 가능 (도메인 용어는 영어)

---

## 2. Git Worktree 협업 체계

### 브랜치 전략

```
main (protected)
├── feat/tui          ← cc-slosim-tui worktree
├── feat/agent-core   ← cc-slosim-agent worktree  
├── feat/qa           ← cc-slosim-qa worktree
└── develop           ← 통합 브랜치 (PR merge target)
```

### Worktree 셋업

```bash
# 기본 디렉토리
cd ~/workspace/02_active/slosim-agent

# 브랜치 생성
git branch feat/tui
git branch feat/agent-core
git branch feat/qa
git branch develop

# Worktree 생성
git worktree add ../slosim-tui feat/tui
git worktree add ../slosim-agent-core feat/agent-core
git worktree add ../slosim-qa feat/qa

# tmux 세션 매핑
# prj-slosim-tui      → ~/workspace/02_active/slosim-tui
# prj-slosim-agent    → ~/workspace/02_active/slosim-agent-core  
# prj-slosim-qa       → ~/workspace/02_active/slosim-qa
```

### 머지 규칙

1. 각 에이전트는 자기 브랜치에만 커밋
2. `develop`로 PR → 마스터 승인 필요
3. 충돌 발생 시 MCC(나)가 중재
4. `main` 머지는 QA 시나리오 전체 통과 후에만

---

## 3. 의존성 그래프

```
                    ┌─────────────┐
                    │   main.go   │
                    │  (cmd/root) │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  internal/  │
                    │    app/     │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
      ┌───────▼──────┐ ┌──▼─────┐ ┌───▼────────┐
      │ internal/tui │ │internal│ │ internal/  │
      │              │ │  /db   │ │  session/  │
      │ [TUI Agent]  │ │        │ │  message/  │
      └───────┬──────┘ └────────┘ └────────────┘
              │
              │ pubsub events
              │
      ┌───────▼──────────────────────────────┐
      │         internal/llm/                 │
      │                                       │
      │  ┌─────────┐  ┌──────────────────┐   │
      │  │ agent/  │  │    provider/     │   │
      │  │         │  │  (Ollama/Qwen3)  │   │
      │  │ sloshing│  └──────────────────┘   │
      │  │ _coder  │                          │
      │  └────┬────┘  ┌──────────────────┐   │
      │       │       │    prompt/       │   │
      │       │       │  sloshing_*      │   │
      │       │       └──────────────────┘   │
      │       │                               │
      │  ┌────▼────────────────────────────┐  │
      │  │        tools/ [Agent Core]      │  │
      │  │                                  │  │
      │  │ gencase ──► dualsphysics ──►    │  │
      │  │              partvtk ──►        │  │
      │  │              measuretool ──►    │  │
      │  │              pvpython           │  │
      │  │                 │               │  │
      │  │           ┌─────▼─────┐         │  │
      │  │           │job_manager│         │  │
      │  │           └─────┬─────┘         │  │
      │  │           ┌─────▼──────┐        │  │
      │  │           │result_store│        │  │
      │  │           └─────┬──────┘        │  │
      │  │           ┌─────▼───────────┐   │  │
      │  │           │report_generator │   │  │
      │  │           └─────────────────┘   │  │
      │  └─────────────────────────────────┘  │
      └───────────────────────────────────────┘
              │
              │ Docker exec
              ▼
      ┌──────────────────┐
      │  DualSPHysics    │
      │  Docker Container│
      │  (CUDA 12.6)     │
      └──────────────────┘
```

### 모듈 간 의존 방향

| From | To | Type |
|------|----|------|
| tui → llm/agent | pubsub 이벤트 구독 | 약결합 |
| llm/agent → tools | Tool interface 호출 | 강결합 |
| tools → Docker | exec/API 호출 | 프로세스 |
| tools → result_store | 파일시스템 I/O | 데이터 |
| qa/tests → tui, tools | import (테스트용) | 테스트 |

### 빌드 순서 (의존성 기반)

```
Phase 1: tools/ (GenCase, Solver, PartVTK, MeasureTool, pvpython)
    ↓
Phase 2: job_manager, result_store
    ↓  
Phase 3: prompt/sloshing_*, agent/sloshing_coder
    ↓
Phase 4: tui/ (Dashboard, Wizard, ResultViewer)
    ↓
Phase 5: Integration + QA scenarios
```

---

## 4. TDD 기반 로드맵

### Phase 1: Foundation (Week 1-2) — Agent Core

**cc-slosim-agent 담당:**

| # | Task | Test First | Implementation |
|---|------|-----------|----------------|
| 1.1 | GenCase Tool | `gencase_test.go` — XML 생성 검증, GenCase 실행 mock | `gencase.go` |
| 1.2 | DualSPHysics Tool | `dualsphysics_test.go` — Docker exec mock, 출력 파싱 | `dualsphysics.go` |
| 1.3 | PartVTK Tool | `partvtk_test.go` — VTK 변환 검증 | `partvtk.go` |
| 1.4 | MeasureTool | `measuretool_test.go` — 측정 데이터 파싱 | `measuretool.go` |
| 1.5 | pvpython Tool | `pvpython_test.go` — ParaView 스크립트 생성 검증 | `pvpython.go` |
| 1.6 | Sloshing 시스템 프롬프트 | 프롬프트 렌더링 테스트 | `sloshing_system.go` |

**테스트 기준:** 각 Tool이 독립적으로 Docker 컨테이너 호출 가능

---

### Phase 2: Orchestration (Week 2-3) — Agent Core + TUI

**cc-slosim-agent:**

| # | Task | Test First | Implementation |
|---|------|-----------|----------------|
| 2.1 | Job Manager | `job_manager_test.go` — submit/monitor/cancel | `job_manager.go` |
| 2.2 | Result Store | `result_store_test.go` — CRUD, 비교 쿼리 | `result_store.go` |
| 2.3 | Sloshing Agent Loop | `sloshing_coder_test.go` — tool 선택 로직 | `agent/sloshing_coder.go` |
| 2.4 | Parametric Study | `parametric_test.go` — 다중 Job 오케스트레이션 | `parametric.go` |

**cc-slosim-tui:**

| # | Task | Test First | Implementation |
|---|------|-----------|----------------|
| 2.5 | Sim Dashboard 컴포넌트 | `dashboard_test.go` — 상태 렌더링 | `tui/components/dashboard.go` |
| 2.6 | Case Setup Wizard | `wizard_test.go` — 단계별 입력 | `tui/components/wizard.go` |

**cc-slosim-qa:**

| # | Task |
|---|------|
| 2.7 | 시나리오 문서 v1.0 작성 |
| 2.8 | E2E 테스트 프레임워크 셋업 |

---

### Phase 3: Integration (Week 3-4) — 전체

| # | Task | 담당 |
|---|------|------|
| 3.1 | Agent Loop + Tools 통합 | cc-slosim-agent |
| 3.2 | TUI + Agent 이벤트 연결 | cc-slosim-tui |
| 3.3 | Result Viewer + pvpython 통합 | cc-slosim-tui |
| 3.4 | Report Generator | cc-slosim-agent |
| 3.5 | E2E 시나리오 1: 기본 슬로싱 | cc-slosim-qa |
| 3.6 | E2E 시나리오 2: STL 입력 | cc-slosim-qa |
| 3.7 | E2E 시나리오 3: 파라메트릭 | cc-slosim-qa |

---

### Phase 4: Polish (Week 4-5) — UX + 논문

| # | Task | 담당 |
|---|------|------|
| 4.1 | TUI 테마/브랜딩 | cc-slosim-tui |
| 4.2 | 에러 핸들링 UX | cc-slosim-tui + qa |
| 4.3 | 프롬프트 최적화 (Qwen3 32b vs 8b) | cc-slosim-agent |
| 4.4 | 비교 리포트 고도화 | cc-slosim-agent |
| 4.5 | 전체 시나리오 회귀 테스트 | cc-slosim-qa |
| 4.6 | 논문용 데모 케이스 준비 | 전체 |

---

## 5. 각 에이전트별 OpenCode 설정

### cc-slosim-tui (.opencode.json)

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "ollama/qwen3-coder:64k",
  "provider": {
    "ollama": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Ollama",
      "options": { "baseURL": "http://localhost:11434/v1" },
      "models": {
        "qwen3-coder:64k": {
          "name": "Qwen3 Coder 64k",
          "limit": { "context": 65536, "output": 16384 }
        }
      }
    }
  },
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  },
  "lsp": {
    "gopls": { "command": "gopls" }
  },
  "permission": {
    "skill": { "of-*": "allow" },
    "tool": { "of_*": "allow" }
  }
}
```

### cc-slosim-agent (.opencode.json)

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "ollama/qwen3:32b-64k",
  "provider": {
    "ollama": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Ollama",
      "options": { "baseURL": "http://localhost:11434/v1" },
      "models": {
        "qwen3:32b-64k": {
          "name": "Qwen3 32B 64k",
          "limit": { "context": 65536, "output": 16384 }
        },
        "qwen3-coder:64k": {
          "name": "Qwen3 Coder 64k",
          "limit": { "context": 65536, "output": 16384 }
        }
      }
    }
  },
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  },
  "lsp": {
    "gopls": { "command": "gopls" }
  },
  "permission": {
    "skill": { "of-*": "allow" },
    "tool": { "of_*": "allow" }
  }
}
```

### cc-slosim-qa (.opencode.json)

```json
{
  "$schema": "https://opencode.ai/config.json",
  "model": "ollama/qwen3-coder:64k",
  "provider": {
    "ollama": {
      "npm": "@ai-sdk/openai-compatible",
      "name": "Ollama",
      "options": { "baseURL": "http://localhost:11434/v1" },
      "models": {
        "qwen3-coder:64k": {
          "name": "Qwen3 Coder 64k",
          "limit": { "context": 65536, "output": 16384 }
        }
      }
    }
  },
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp@latest"]
    }
  },
  "lsp": {
    "gopls": { "command": "gopls" }
  },
  "permission": {
    "skill": { "of-*": "allow" },
    "tool": { "of_*": "allow", "bash": "deny-write" }
  }
}
```

---

## 6. tmux 세션 매핑

| tmux 세션 | Worktree | 브랜치 | 에이전트 |
|-----------|----------|--------|----------|
| `prj-slosim-tui` | `~/workspace/02_active/slosim-tui` | `feat/tui` | cc-slosim-tui |
| `prj-slosim-agent` | `~/workspace/02_active/slosim-agent-core` | `feat/agent-core` | cc-slosim-agent |
| `prj-slosim-qa` | `~/workspace/02_active/slosim-qa` | `feat/qa` | cc-slosim-qa |
| `prj-slosim-main` | `~/workspace/02_active/slosim-agent` | `main` | MCC (통합) |

---

## 7. MCC 오케스트레이션 역할

나(MCC)의 역할:
1. **위임:** 각 에이전트에 Phase별 태스크 할당 (`mcc-delegate.sh`)
2. **모니터링:** 진행 상황 폴링 + 완료 보고
3. **충돌 중재:** worktree 간 merge conflict 해결
4. **통합:** develop 브랜치 머지 + 빌드 검증
5. **보고:** 마스터에게 진행 상황 보고

---

## 8. 미결 / 다음 단계

- [ ] Worktree 셋업 실행
- [ ] 각 에이전트 CLAUDE.md 작성
- [ ] Phase 1 태스크 위임 시작
- [ ] ParaView Docker 이미지 빌드 또는 호스트 설치 결정
- [ ] Context7에 DualSPHysics 문서 등록 여부 확인
