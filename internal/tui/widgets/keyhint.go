package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// KeyHint represents a single keybinding hint.
type KeyHint struct {
	Key         string
	Description string
}

// KeyHintBar renders a horizontal bar of keybinding hints: "[q] Quit  [p] Pause".
func KeyHintBar(hints []KeyHint) string {
	if len(hints) == 0 {
		return ""
	}

	t := theme.CurrentTheme().Tokens()

	keyStyle := lipgloss.NewStyle().
		Foreground(t.PanelTitle).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(t.DataLabel)

	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = keyStyle.Render("["+h.Key+"]") + " " + descStyle.Render(h.Description)
	}

	return strings.Join(parts, "  ")
}
