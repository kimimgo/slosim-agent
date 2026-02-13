---
description: Create a new BubbleTea TUI component following project patterns
argument-hint: "<component-name>"
---

Create a new BubbleTea component named `$ARGUMENTS` for the slosim-agent TUI.

## Component Structure

Generate two files:
1. `internal/tui/components/<name>/<name>.go` — Component implementation
2. `internal/tui/components/<name>/<name>_test.go` — teatest unit test

## Implementation Pattern

Follow the existing OpenCode TUI conventions in `internal/tui/components/`:

```go
package <name>

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/opencode-ai/opencode/internal/pubsub"
)

// Messages
type <Name>Msg struct { /* ... */ }

// Interface
type <Name>Cmp interface {
    tea.Model
}

// Model
type <name>Cmp struct {
    width  int
    height int
    // component-specific fields
}

func New() <Name>Cmp {
    return &<name>Cmp{}
}

func (m *<name>Cmp) Init() tea.Cmd { return nil }

func (m *<name>Cmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m *<name>Cmp) View() string {
    // Use Lipgloss adaptive styling
    style := lipgloss.NewStyle().
        Width(m.width).
        Foreground(lipgloss.AdaptiveColor{Light: "#333", Dark: "#EEE"})
    return style.Render("...")
}
```

## Test Pattern

```go
package <name>_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) { /* verify initial state */ }
func TestUpdate_WindowSize(t *testing.T) { /* verify responsive layout */ }
func TestView_Render(t *testing.T) { /* verify output contains expected content */ }
```

## Rules

- Subscribe to pubsub events for Agent Core communication
- Use Lipgloss AdaptiveColor for theme compatibility
- Handle tea.WindowSizeMsg for responsive layout
- Component name must match directory name
