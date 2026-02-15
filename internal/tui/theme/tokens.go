package theme

import "github.com/charmbracelet/lipgloss"

// SemanticTokens defines role-based design tokens for TUI components.
// All components MUST use these tokens instead of lipgloss.Color() directly.
type SemanticTokens struct {
	// Panel
	PanelBg          lipgloss.TerminalColor
	PanelBorder      lipgloss.TerminalColor
	PanelBorderFocus lipgloss.TerminalColor
	PanelTitle       lipgloss.TerminalColor

	// Status
	StatusRunning lipgloss.TerminalColor
	StatusError   lipgloss.TerminalColor
	StatusWarning lipgloss.TerminalColor
	StatusIdle    lipgloss.TerminalColor

	// Progress
	ProgressFill  lipgloss.TerminalColor
	ProgressEmpty lipgloss.TerminalColor
	ProgressLabel lipgloss.TerminalColor

	// List
	ListCursor     lipgloss.TerminalColor
	ListSelected   lipgloss.TerminalColor
	ListItemNormal lipgloss.TerminalColor

	// Data display
	DataLabel lipgloss.TerminalColor
	DataValue lipgloss.TerminalColor
	DataUnit  lipgloss.TerminalColor

	// Spacing
	PanelPadding int
	PanelMargin  int
	SectionGap   int
}

// Tokens returns semantic design tokens derived from the theme's base colors.
func (t *BaseTheme) Tokens() SemanticTokens {
	return SemanticTokens{
		PanelBg:          t.BackgroundSecondaryColor,
		PanelBorder:      t.BorderNormalColor,
		PanelBorderFocus: t.BorderFocusedColor,
		PanelTitle:       t.PrimaryColor,

		StatusRunning: t.SuccessColor,
		StatusError:   t.ErrorColor,
		StatusWarning: t.WarningColor,
		StatusIdle:    t.TextMutedColor,

		ProgressFill:  t.PrimaryColor,
		ProgressEmpty: t.BorderDimColor,
		ProgressLabel: t.TextColor,

		ListCursor:     t.PrimaryColor,
		ListSelected:   t.AccentColor,
		ListItemNormal: t.TextColor,

		DataLabel: t.TextMutedColor,
		DataValue: t.TextColor,
		DataUnit:  t.TextMutedColor,

		PanelPadding: 1,
		PanelMargin:  0,
		SectionGap:   1,
	}
}
