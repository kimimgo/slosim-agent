# slosim-agent v0.3 TUI Design Plan

> 2026-02-15 작성 | 딥리서치 기반 (BubbleTea/Lipgloss 고급 기법 + AI TUI 개발 전략)

---

## 1. 현재 상태 진단

### 1.1 컴포넌트 성숙도

| 컴포넌트 | 상태 | LOC | 문제점 |
|----------|------|-----|--------|
| Chat (editor/list/sidebar) | POLISHED | ~800 | 없음 |
| Dialog System (8종) | POLISHED | ~500 | 없음 |
| Status Bar | POLISHED | ~200 | 없음 |
| Theme System (10개) | EXCELLENT | ~2000 | 48-method interface, AdaptiveColor |
| Markdown Rendering | VERY GOOD | ~170 | Glamour + ansi config |
| Parametric View | FUNCTIONAL | ~240 | H/V scroll, sorting |
| Simulation Dashboard | **STUB** | ~96 | 하드코딩 색상 `"99"`, `"240"` |
| Results Viewer | **BASIC** | ~147 | 하드코딩 `"212"`, 미리보기 없음 |
| Error Panel | **STUB** | ~150 | 에러 라우팅 미연결 |
| Logs Viewer | **PARTIAL** | ~100 | 필터/검색 없음 |

### 1.2 핵심 문제: AI가 TUI를 못 만드는 이유

1. **참조 없는 생성**: AI가 "대시보드 만들어줘" → 무작위 색상 + 기본 레이아웃
2. **하드코딩 습관**: `lipgloss.Color("99")` 대신 `theme.CurrentTheme().Primary()` 사용해야 함
3. **Bubbles 활용 부족**: viewport, table, list, progress 등 기성 컴포넌트 대신 처음부터 작성
4. **측정 없는 레이아웃**: `lipgloss.Width()/Height()` 대신 magic number 사용
5. **디자인 시스템 부재**: 색상/간격/보더 규칙이 코드에 분산

---

## 2. AI 디자인 한계 극복 전략

### 2.1 Design Token 시스템 (최우선)

**목표**: AI가 임의 색상을 쓰지 못하도록 타입 안전 Design Token 강제.

```go
// internal/tui/theme/tokens.go (신규)
package theme

import "github.com/charmbracelet/lipgloss"

// SemanticTokens — 모든 컴포넌트가 이 토큰만 사용
type SemanticTokens struct {
    // Panel
    PanelBg          lipgloss.TerminalColor
    PanelBorder      lipgloss.TerminalColor
    PanelBorderFocus lipgloss.TerminalColor
    PanelTitle       lipgloss.TerminalColor

    // Status
    StatusRunning    lipgloss.TerminalColor  // Green
    StatusError      lipgloss.TerminalColor  // Red
    StatusWarning    lipgloss.TerminalColor  // Yellow
    StatusIdle       lipgloss.TerminalColor  // Muted

    // Progress
    ProgressFill     lipgloss.TerminalColor
    ProgressEmpty    lipgloss.TerminalColor
    ProgressLabel    lipgloss.TerminalColor

    // List
    ListCursor       lipgloss.TerminalColor
    ListSelected     lipgloss.TerminalColor
    ListItemNormal   lipgloss.TerminalColor

    // Data
    DataLabel        lipgloss.TerminalColor
    DataValue        lipgloss.TerminalColor
    DataUnit         lipgloss.TerminalColor

    // Spacing
    PanelPadding     int  // 1
    PanelMargin      int  // 0
    SectionGap       int  // 1
}

// 각 Theme이 SemanticTokens를 반환
func (t *CatppuccinTheme) Tokens() SemanticTokens { ... }
```

**AI 프롬프트 제약**: "Never use `lipgloss.Color()` directly. Always use `theme.CurrentTheme().Tokens().XXX`."

### 2.2 Reference-Driven Development (참조 기반 개발)

**원칙**: 모든 TUI 작업에 구체적 참조 앱을 지정.

| 컴포넌트 | 참조 앱 | 참조 요소 |
|----------|---------|-----------|
| Sim Dashboard | **K9s** pod monitor | 실시간 메트릭 테이블 + 상태 색상 |
| Results Browser | **lazygit** file panel | 포커스 보더 + 파일 트리 + 미리보기 분할 |
| Progress Bar | **gum** progress | 그라디언트 + 부드러운 이징 |
| Error Panel | **soft-serve** error view | 보더 색상 + 아이콘 + 액션 버튼 |
| Case Wizard | **huh** form library | 단계별 입력 + 검증 + 테마 |
| Log Viewer | **lazygit** commit log | 필터링 + 색상 코딩 + 스크롤 |

**AI 프롬프트 예시**:
```
Build the Simulation Dashboard. Reference:
- K9s pod monitor layout: real-time table with status colors
- lazygit focus states: green bold active borders, gray inactive
- Use bubbles/table for the job list
- Use bubbles/progress for the progress bar
- All colors from theme.CurrentTheme().Tokens()
```

### 2.3 Component-First Architecture (기성 컴포넌트 우선)

**규칙**: Bubbles 라이브러리에 있는 컴포넌트는 반드시 사용.

| 용도 | Bubbles 컴포넌트 | 커스텀 금지 |
|------|-----------------|------------|
| 스크롤 영역 | `viewport.Model` | 수동 offset 계산 |
| 테이블 | `table.Model` | 수동 열 정렬 |
| 목록 | `list.Model` | 수동 커서 관리 |
| 진행률 | `progress.Model` | 수동 바 그리기 |
| 입력 | `textinput.Model` | 수동 커서 관리 |
| 스피너 | `spinner.Model` | `tea.Tick` 직접 사용 |
| 키바인딩 | `key.Binding` | 하드코딩 키 문자열 |
| 도움말 | `help.Model` | 수동 help 텍스트 |

**신규 도입 고려**:
- `charmbracelet/huh` — Case Wizard용 폼 라이브러리
- `charmbracelet/harmonica` — 스프링 기반 애니메이션

### 2.4 Screenshot-Driven Validation (스크린샷 검증)

**워크플로우**:
1. `charmbracelet/vhs`로 TUI 스크린샷 자동 촬영
2. 레퍼런스 이미지와 diff 비교
3. 시각적 회귀 방지

```tape
# .vhs/sim_dashboard.tape
Output sim_dashboard.gif
Set Shell "bash"
Set FontSize 14
Set Width 120
Set Height 40
Type "./slosim-agent"
Sleep 2s
Type "1m x 0.5m x 0.6m 탱크 슬로싱 해석"
Enter
Sleep 5s
Screenshot sim_dashboard.png
```

### 2.5 Layout Constraint System (레이아웃 제약)

**규칙**: 모든 크기는 `tea.WindowSizeMsg` 기반으로 계산.

```go
// GOOD — 동적 측정
func (d *Dashboard) View() string {
    header := d.renderHeader()
    headerH := lipgloss.Height(header)

    tableH := d.height - headerH - footerH - 2  // 2 = border
    table := d.renderTable(tableH)

    return lipgloss.JoinVertical(lipgloss.Left, header, table, footer)
}

// BAD — 하드코딩
func (d *Dashboard) View() string {
    table := d.renderTable(20)  // ❌ Magic number
}
```

---

## 3. 컴포넌트별 리디자인 계획

### 3.1 Simulation Dashboard (TUI-05)

**현재**: 하드코딩 색상 스텁, 실제 Job 연결 없음
**목표**: K9s 스타일 실시간 모니터링 대시보드

```
┌─ Simulation Dashboard ─────────────────────────────────────────┐
│                                                                 │
│  Job: sloshing_001        Status: ● RUNNING     GPU: RTX 4090  │
│  TimeStep: 0.0012s        Time: 1.45/2.00s                     │
│                                                                 │
│  Progress ████████████████████░░░░░░░░░  72.5%  ETA: 38s       │
│                                                                 │
│  ┌─ Metrics ──────────────┐  ┌─ Recent Log ─────────────────┐  │
│  │ Particles: 136,420     │  │ [14:32:01] Part 0072 saved   │  │
│  │ Kinetic E: 0.0234 J    │  │ [14:32:00] Step dt=0.00012   │  │
│  │ dt:        0.00012s    │  │ [14:31:59] Interaction...     │  │
│  │ Memory:    2.1 GB      │  │ [14:31:58] Part 0071 saved   │  │
│  └────────────────────────┘  └──────────────────────────────┘  │
│                                                                 │
│  [q] Quit  [p] Pause  [l] Full Log  [r] Refresh                │
└─────────────────────────────────────────────────────────────────┘
```

**구현 요소**:
- `bubbles/progress` — 테마 그라디언트 적용
- `bubbles/table` — 메트릭 키-값 표시
- `bubbles/viewport` — 로그 스크롤
- `harmonica.Spring` — 진행률 바 부드러운 업데이트
- `theme.Tokens().StatusRunning/Error/Warning` — 상태 색상

### 3.2 Results Browser (TUI-03 고도화)

**현재**: 파일 목록만 표시, 하드코딩 커서 색상
**목표**: lazygit 스타일 2-패널 (목록 + 미리보기)

```
┌─ Results ────────────────┐┌─ Preview ───────────────────────┐
│ ● PartFluid_0000.vtk     ││                                 │
│   PartFluid_0001.vtk     ││  [VTK file - 136420 particles]  │
│   PartFluid_0002.vtk     ││                                 │
│   ...                    ││  Fields: Vel, Press, Rhop       │
│   PartFluid_0100.vtk     ││  Bounds: [0,1] x [0,0.5]       │
│ ● Run.csv                ││                                 │
│   report.md              ││  Size: 12.3 MB                  │
│   animation.mp4          ││  Created: 2026-02-15 14:32      │
│                          ││                                 │
│  101 files, 1.2 GB total ││  [Enter] Open  [v] View  [d] D │
└──────────────────────────┘└─────────────────────────────────┘
```

**구현 요소**:
- `bubbles/list` — 파일 목록 (필터링, 페이징)
- `layout.SplitPaneLayout` — 좌/우 분할 (40:60)
- 포커스 상태: `Tokens().PanelBorderFocus` (활성), `Tokens().PanelBorder` (비활성)
- 파일 타입별 아이콘 + 메타데이터 미리보기

### 3.3 Case Wizard (TUI-04)

**현재**: 미구현
**목표**: `huh` 폼 라이브러리 기반 단계별 입력

```
┌─ New Simulation ── Step 2/4: Tank Geometry ─────────────────────┐
│                                                                  │
│  Tank Type:  ○ Rectangular  ● Cylindrical  ○ L-shaped  ○ STL   │
│                                                                  │
│  Radius:     ┃ 0.5                                    (meters)  │
│  Height:     ┃ 0.8                                    (meters)  │
│  Fill Ratio: ┃ 50                                     (%)       │
│                                                                  │
│  ┌─ Preview ─────────────────────────┐                          │
│  │      ╭──────╮                     │                          │
│  │      │~~~~~~│ ← 50% fill         │                          │
│  │      │      │                     │                          │
│  │      ╰──────╯                     │                          │
│  │   R=0.5m, H=0.8m                 │                          │
│  └───────────────────────────────────┘                          │
│                                                                  │
│  [←] Back  [→] Next  [Esc] Cancel                               │
└──────────────────────────────────────────────────────────────────┘
```

**구현 요소**:
- `charmbracelet/huh` — 폼 필드 (Input, Select, Confirm)
- Catppuccin 테마 적용 (`huh.ThemeCatppuccin()`)
- 4단계: Tank Geometry → Fluid Properties → Excitation → Simulation Parameters
- ASCII 탱크 미리보기 (실시간 업데이트)

### 3.4 Error Recovery Panel (TUI-06 개선)

**현재**: 스텁
**목표**: 인라인 에러 표시 + 원클릭 복구

```
┌─ ⚠ Error Detected ─────────────────────────────────────────────┐
│                                                                  │
│  Type: Solver Divergence                                        │
│  Time: t=0.832s (Part 0041)                                    │
│                                                                  │
│  Message:                                                        │
│  Kinetic energy exceeds threshold (EKin=1.2e4 > limit=1e3)     │
│  Particle acceleration unstable at boundary                     │
│                                                                  │
│  Suggested Fix:                                                  │
│  Reduce dp from 0.005 to 0.003 and re-run from t=0             │
│                                                                  │
│  [r] Retry with fix  [i] Ignore  [d] Dismiss  [?] Details      │
└──────────────────────────────────────────────────────────────────┘
```

### 3.5 Log Viewer 고도화

**현재**: 기본 테이블 + 상세
**목표**: lazygit commit log 스타일 필터 + 색상 코딩

**구현 요소**:
- `bubbles/viewport` — 스크롤
- 레벨별 색상: `Tokens().StatusError` (ERROR), `Tokens().StatusWarning` (WARN), `Tokens().DataValue` (INFO)
- `/` 키로 인라인 검색
- `bubbles/key` — 표준 vim 키바인딩

---

## 4. 디자인 시스템 구축 계획

### 4.1 파일 구조

```
internal/tui/
├── theme/
│   ├── theme.go          — Theme interface (기존 48 methods)
│   ├── tokens.go         — 🆕 SemanticTokens struct
│   ├── manager.go        — Global registry (기존)
│   ├── catppuccin.go     — 기존 + Tokens() 구현
│   └── ... (10개 테마)
├── widgets/              — 🆕 재사용 가능 위젯
│   ├── panel.go          — 테마 적용 패널 (보더 + 타이틀 + 포커스)
│   ├── metric.go         — 키-값 메트릭 표시
│   ├── status_badge.go   — 상태 배지 (● RUNNING, ● ERROR)
│   └── keyhint.go        — 하단 키 힌트 바
├── styles/
│   ├── styles.go         — 기존 헬퍼 (token 기반으로 리팩토링)
│   └── markdown.go       — 기존 (변경 없음)
├── components/
│   ├── simulation/
│   │   └── dashboard.go  — 🔄 리디자인 (K9s 참조)
│   ├── results/
│   │   ├── browser.go    — 🔄 리디자인 (lazygit 참조)
│   │   └── viewer.go     — 🔄 메타데이터 미리보기
│   ├── wizard/           — 🆕 Case Wizard (huh 기반)
│   │   └── wizard.go
│   └── error/
│       └── panel.go      — 🔄 리디자인
└── layout/
    └── (기존 유지)
```

### 4.2 공통 위젯 라이브러리 (`widgets/`)

**Panel** — 모든 박스형 UI의 기본:
```go
type Panel struct {
    Title    string
    Focused  bool
    Content  string
    Width    int
    Height   int
}

func (p Panel) View() string {
    t := theme.CurrentTheme().Tokens()
    borderColor := t.PanelBorder
    if p.Focused {
        borderColor = t.PanelBorderFocus
    }

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(borderColor).
        Padding(t.PanelPadding).
        Width(p.Width).
        Height(p.Height).
        Render(p.Content)
}
```

**StatusBadge** — 상태 표시:
```go
func StatusBadge(status string) string {
    t := theme.CurrentTheme().Tokens()
    var color lipgloss.TerminalColor
    var icon string
    switch status {
    case "RUNNING": color, icon = t.StatusRunning, "●"
    case "ERROR":   color, icon = t.StatusError, "●"
    case "DONE":    color, icon = t.StatusRunning, "✓"
    case "IDLE":    color, icon = t.StatusIdle, "○"
    }
    return lipgloss.NewStyle().Foreground(color).Render(icon + " " + status)
}
```

### 4.3 신규 의존성

```go
// go.mod 추가
require (
    github.com/charmbracelet/huh v0.7.0       // Case Wizard 폼
    github.com/charmbracelet/harmonica v0.2.0  // 스프링 애니메이션
)
```

---

## 5. 개발 프로세스 (AI 품질 보장)

### 5.1 TUI 코딩 규칙 (.claude/agents/ 또는 CLAUDE.md 추가)

```markdown
## TUI Component Rules (AI MUST follow)

1. **NEVER** use `lipgloss.Color("...")` directly — use `theme.CurrentTheme().Tokens().XXX`
2. **ALWAYS** use Bubbles components (viewport, table, list, progress, spinner)
3. **ALWAYS** handle `tea.WindowSizeMsg` for responsive layout
4. **ALWAYS** measure with `lipgloss.Width()/Height()` — never hardcode dimensions
5. **ALWAYS** use `widgets.Panel()` for boxed containers
6. **ALWAYS** use `widgets.StatusBadge()` for status indicators
7. **REFERENCE** existing polished components (chat/chat.go, dialog/commands.go) for patterns
8. **TEST** with `Ctrl+T` theme switch — all colors must change correctly
```

### 5.2 개발 순서

```
Phase 1: 인프라 (Week 1)
├── tokens.go — SemanticTokens 구조체 + 10개 테마별 구현
├── widgets/ — Panel, StatusBadge, Metric, KeyHint
├── go.mod — huh, harmonica 추가
└── CLAUDE.md — TUI 코딩 규칙 추가

Phase 2: 컴포넌트 리디자인 (Week 2-3)
├── simulation/dashboard.go — K9s 참조, Bubbles table+progress+viewport
├── results/browser.go — lazygit 참조, 2-패널 분할
├── error/panel.go — 인라인 에러 + 복구 액션
└── logs/table.go — 필터 + 색상 코딩

Phase 3: 신규 컴포넌트 (Week 3-4)
├── wizard/wizard.go — huh 기반 4단계 Case Wizard
└── VHS 스크린샷 테스트 스크립트

Phase 4: 통합 QA (Week 4)
├── 전체 테마 10개 순회 테스트
├── 터미널 크기 변경 테스트 (80x24 ~ 200x60)
└── VHS GIF 데모 생성
```

### 5.3 품질 게이트

| 체크 | 기준 |
|------|------|
| 하드코딩 색상 | `grep -r 'lipgloss.Color("' components/` → 0건 |
| 테마 호환 | 10개 테마 전부 `Ctrl+T`로 전환 시 깨지지 않음 |
| 반응형 | 80x24 최소 크기에서 모든 컴포넌트 렌더 가능 |
| Bubbles 사용 | 커스텀 스크롤/테이블/진행률 0건 |
| 키바인딩 | `bubbles/key.Binding` 사용, 도움말 자동 생성 |

---

## 6. 레퍼런스 자료

### 6.1 디자인 참조 앱

| 앱 | 참조 요소 | GitHub |
|----|-----------|--------|
| **K9s** | 실시간 모니터링 대시보드 | github.com/derailed/k9s |
| **lazygit** | 포커스 상태, 패널 분할, 커밋 로그 | github.com/jesseduffield/lazygit |
| **gum** | 진행률 바, 보더, 간격 패턴 | github.com/charmbracelet/gum |
| **soft-serve** | 멀티 패널, 권한별 UI, 에러 표시 | github.com/charmbracelet/soft-serve |

### 6.2 Charm 생태계 도구

| 도구 | 용도 |
|------|------|
| **huh** | 대화형 폼 (Case Wizard) |
| **harmonica** | 스프링 기반 애니메이션 (진행률 바) |
| **glamour** | Markdown 렌더링 (이미 사용 중) |
| **vhs** | TUI 스크린샷/GIF 자동 생성 (테스트용) |
| **freeze** | 코드 → 이미지 변환 (문서용) |
| **log** | 구조화 로깅 (로그 뷰어 통합) |

### 6.3 디자인 토큰/팔레트

| 자료 | URL |
|------|-----|
| **Catppuccin Palette** | catppuccin.com/palette |
| **lipgloss-theme 참조 구현** | github.com/purpleclay/lipgloss-theme |
| **termenv 감지** | github.com/muesli/termenv |

### 6.4 BubbleTea 학습 자료

| 자료 | URL |
|------|-----|
| BubbleTea 51 예제 | github.com/charmbracelet/bubbletea/tree/main/examples |
| Bubbles 컴포넌트 | github.com/charmbracelet/bubbles |
| Tips for Building BubbleTea | leg100.github.io/en/posts/building-bubbletea-programs/ |
| State Machine 패턴 | zackproser.com/blog/bubbletea-state-machine |

---

## 7. BubbleTea v2 마이그레이션 고려

> 현재 v1.3.5 → v2 마이그레이션은 v0.3 범위 밖. v1.0에서 검토.

### 주요 변경점

| v1 | v2 |
|----|-----|
| `View() string` | `View() tea.View` (구조체) |
| 문자열 기반 렌더링 | 선언적 View API |
| 기본 렌더러 | Cursed Renderer (10x 빠름) |
| 기본 키보드 | Enhanced Keyboard (`shift+enter` 감지) |
| 별도 progress | 내장 `tea.ProgressBar` |
| 별도 clipboard | OSC52 클립보드 (SSH에서도 동작) |
| Import: `github.com/charmbracelet/bubbletea` | `charm.land/bubbletea/v2` |

**결론**: v0.3에서는 v1.3.5 유지하되, Design Token + Widget 추상화를 통해 v2 마이그레이션 용이하게 설계.
