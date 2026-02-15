package widgets

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// Panel renders a themed bordered container with title and focus state.
type Panel struct {
	Title   string
	Content string
	Width   int
	Height  int
	Focused bool
}

// View renders the panel using theme tokens.
func (p Panel) View() string {
	t := theme.CurrentTheme().Tokens()

	borderColor := t.PanelBorder
	if p.Focused {
		borderColor = t.PanelBorderFocus
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PanelTitle)

	title := titleStyle.Render(p.Title)

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(t.PanelPadding).
		Background(t.PanelBg)

	if p.Width > 0 {
		style = style.Width(p.Width)
	}
	if p.Height > 0 {
		style = style.Height(p.Height)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, p.Content)
	return style.Render(content)
}
