---
name: go-tui-design
description: BubbleTea/Lipgloss TUI 디자인 품질 가이드
autoTrigger: true
triggerPattern: "internal/tui/"
---

# BubbleTea TUI Design Quality Rules

> `internal/tui/components/` 파일 편집 시 자동 발동

## 절대 금지

1. **하드코딩 색상 금지**: `lipgloss.Color("99")` ❌ → `theme.CurrentTheme().Primary()` ✅
2. **고정 크기 금지**: `.Width(80)` ❌ → `lipgloss.Width(content)` 또는 `msg.Width` ✅
3. **커스텀 스크롤 금지**: Bubbles `viewport.Model` 사용 ✅
4. **커스텀 테이블 금지**: Bubbles `table.Model` 사용 ✅
5. **커스텀 진행률 금지**: Bubbles `progress.Model` 사용 ✅

## 필수 패턴

### 색상
```go
// ALWAYS use theme system
t := theme.CurrentTheme()
style := lipgloss.NewStyle().Foreground(t.Primary())

// NEVER hardcode
style := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))  // ❌
```

### 레이아웃
```go
// ALWAYS measure dynamically
headerH := lipgloss.Height(header)
contentH := totalHeight - headerH - footerH

// ALWAYS handle WindowSizeMsg
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

### 보더
```go
// Focus state distinction
if focused {
    style = style.BorderForeground(t.BorderFocused())
} else {
    style = style.BorderForeground(t.BorderNormal())
}

// Border height accounting (-2 for top+bottom border)
contentH := panelH - 2
```

### 키바인딩
```go
// Use bubbles/key for remappable bindings
import "github.com/charmbracelet/bubbles/key"

key.NewBinding(
    key.WithKeys("up", "k"),
    key.WithHelp("↑/k", "move up"),
)
```

## 참조 앱 (컴포넌트별)

| 컴포넌트 | 참조 | 핵심 패턴 |
|----------|------|-----------|
| Dashboard | K9s | 실시간 메트릭 테이블 + 상태 색상 |
| File Browser | lazygit | 포커스 보더 + 2-패널 분할 |
| Progress | gum | 그라디언트 이징 |
| Forms | huh | 단계별 입력 + Catppuccin 테마 |
| Error | soft-serve | 인라인 에러 + 액션 버튼 |

## 기존 폴리시된 컴포넌트 참조

새 컴포넌트 작성 시 기존 코드를 반드시 참조:
- `chat/chat.go` — 테마 통합 패턴
- `dialog/commands.go` — 오버레이 + 키바인딩
- `core/status.go` — 실시간 업데이트 패턴
