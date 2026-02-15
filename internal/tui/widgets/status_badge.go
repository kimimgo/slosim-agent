package widgets

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// StatusBadge renders a colored status indicator: "● STATUS" or "✓ DONE".
func StatusBadge(status string) string {
	t := theme.CurrentTheme().Tokens()

	var color lipgloss.TerminalColor
	var icon string

	switch status {
	case "RUNNING":
		color, icon = t.StatusRunning, "●"
	case "ERROR":
		color, icon = t.StatusError, "●"
	case "WARNING":
		color, icon = t.StatusWarning, "●"
	case "DONE":
		color, icon = t.StatusRunning, "✓"
	case "IDLE":
		color, icon = t.StatusIdle, "○"
	default:
		color, icon = t.StatusIdle, "?"
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Render(icon + " " + status)
}
