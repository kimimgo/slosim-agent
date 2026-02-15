package widgets

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// Metric renders a key-value pair: "Label: Value Unit".
type Metric struct {
	Label string
	Value string
	Unit  string
}

// View renders the metric using theme tokens.
func (m Metric) View() string {
	t := theme.CurrentTheme().Tokens()

	labelStyle := lipgloss.NewStyle().Foreground(t.DataLabel)
	valueStyle := lipgloss.NewStyle().Foreground(t.DataValue).Bold(true)
	unitStyle := lipgloss.NewStyle().Foreground(t.DataUnit)

	result := labelStyle.Render(m.Label+":") + " " + valueStyle.Render(m.Value)
	if m.Unit != "" {
		result += " " + unitStyle.Render(m.Unit)
	}
	return result
}
