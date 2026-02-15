package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestSemanticTokensStruct(t *testing.T) {
	// Verify SemanticTokens struct has all expected fields
	tokens := SemanticTokens{
		PanelBg:          lipgloss.AdaptiveColor{Dark: "#000", Light: "#fff"},
		PanelBorder:      lipgloss.AdaptiveColor{Dark: "#111", Light: "#eee"},
		PanelBorderFocus: lipgloss.AdaptiveColor{Dark: "#222", Light: "#ddd"},
		PanelTitle:       lipgloss.AdaptiveColor{Dark: "#333", Light: "#ccc"},

		StatusRunning: lipgloss.AdaptiveColor{Dark: "#0f0", Light: "#0a0"},
		StatusError:   lipgloss.AdaptiveColor{Dark: "#f00", Light: "#a00"},
		StatusWarning: lipgloss.AdaptiveColor{Dark: "#ff0", Light: "#aa0"},
		StatusIdle:    lipgloss.AdaptiveColor{Dark: "#888", Light: "#666"},

		ProgressFill:  lipgloss.AdaptiveColor{Dark: "#0f0", Light: "#0a0"},
		ProgressEmpty: lipgloss.AdaptiveColor{Dark: "#333", Light: "#ccc"},
		ProgressLabel: lipgloss.AdaptiveColor{Dark: "#fff", Light: "#000"},

		ListCursor:    lipgloss.AdaptiveColor{Dark: "#0ff", Light: "#099"},
		ListSelected:  lipgloss.AdaptiveColor{Dark: "#00f", Light: "#009"},
		ListItemNormal: lipgloss.AdaptiveColor{Dark: "#aaa", Light: "#555"},

		DataLabel: lipgloss.AdaptiveColor{Dark: "#888", Light: "#666"},
		DataValue: lipgloss.AdaptiveColor{Dark: "#fff", Light: "#000"},
		DataUnit:  lipgloss.AdaptiveColor{Dark: "#666", Light: "#888"},

		PanelPadding: 1,
		PanelMargin:  0,
		SectionGap:   1,
	}

	// Verify spacing defaults
	if tokens.PanelPadding != 1 {
		t.Errorf("PanelPadding = %d, want 1", tokens.PanelPadding)
	}
	if tokens.PanelMargin != 0 {
		t.Errorf("PanelMargin = %d, want 0", tokens.PanelMargin)
	}
	if tokens.SectionGap != 1 {
		t.Errorf("SectionGap = %d, want 1", tokens.SectionGap)
	}
}

func TestBaseThemeTokensDefaultMapping(t *testing.T) {
	// BaseTheme.Tokens() should map existing colors to semantic tokens
	bt := &BaseTheme{
		PrimaryColor:            lipgloss.AdaptiveColor{Dark: "#0000ff", Light: "#0000aa"},
		SecondaryColor:          lipgloss.AdaptiveColor{Dark: "#ff00ff", Light: "#aa00aa"},
		AccentColor:             lipgloss.AdaptiveColor{Dark: "#ff8800", Light: "#aa5500"},
		ErrorColor:              lipgloss.AdaptiveColor{Dark: "#ff0000", Light: "#aa0000"},
		WarningColor:            lipgloss.AdaptiveColor{Dark: "#ffff00", Light: "#aaaa00"},
		SuccessColor:            lipgloss.AdaptiveColor{Dark: "#00ff00", Light: "#00aa00"},
		TextColor:               lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		TextMutedColor:          lipgloss.AdaptiveColor{Dark: "#888888", Light: "#666666"},
		BackgroundColor:         lipgloss.AdaptiveColor{Dark: "#212121", Light: "#f8f8f8"},
		BackgroundSecondaryColor: lipgloss.AdaptiveColor{Dark: "#2c2c2c", Light: "#e0e0e0"},
		BorderNormalColor:       lipgloss.AdaptiveColor{Dark: "#444444", Light: "#cccccc"},
		BorderFocusedColor:      lipgloss.AdaptiveColor{Dark: "#0000ff", Light: "#0000aa"},
		BorderDimColor:          lipgloss.AdaptiveColor{Dark: "#333333", Light: "#dddddd"},
	}

	tokens := bt.Tokens()

	// Panel colors should map from background/border
	assertColorEquals(t, "PanelBg", tokens.PanelBg, bt.BackgroundSecondaryColor)
	assertColorEquals(t, "PanelBorder", tokens.PanelBorder, bt.BorderNormalColor)
	assertColorEquals(t, "PanelBorderFocus", tokens.PanelBorderFocus, bt.BorderFocusedColor)
	assertColorEquals(t, "PanelTitle", tokens.PanelTitle, bt.PrimaryColor)

	// Status colors should map from status
	assertColorEquals(t, "StatusRunning", tokens.StatusRunning, bt.SuccessColor)
	assertColorEquals(t, "StatusError", tokens.StatusError, bt.ErrorColor)
	assertColorEquals(t, "StatusWarning", tokens.StatusWarning, bt.WarningColor)
	assertColorEquals(t, "StatusIdle", tokens.StatusIdle, bt.TextMutedColor)

	// Progress colors
	assertColorEquals(t, "ProgressFill", tokens.ProgressFill, bt.PrimaryColor)
	assertColorEquals(t, "ProgressEmpty", tokens.ProgressEmpty, bt.BorderDimColor)
	assertColorEquals(t, "ProgressLabel", tokens.ProgressLabel, bt.TextColor)

	// List colors
	assertColorEquals(t, "ListCursor", tokens.ListCursor, bt.PrimaryColor)
	assertColorEquals(t, "ListSelected", tokens.ListSelected, bt.AccentColor)
	assertColorEquals(t, "ListItemNormal", tokens.ListItemNormal, bt.TextColor)

	// Data colors
	assertColorEquals(t, "DataLabel", tokens.DataLabel, bt.TextMutedColor)
	assertColorEquals(t, "DataValue", tokens.DataValue, bt.TextColor)
	assertColorEquals(t, "DataUnit", tokens.DataUnit, bt.TextMutedColor)

	// Spacing defaults
	if tokens.PanelPadding != 1 {
		t.Errorf("PanelPadding = %d, want 1", tokens.PanelPadding)
	}
	if tokens.PanelMargin != 0 {
		t.Errorf("PanelMargin = %d, want 0", tokens.PanelMargin)
	}
	if tokens.SectionGap != 1 {
		t.Errorf("SectionGap = %d, want 1", tokens.SectionGap)
	}
}

func TestAllThemesHaveTokens(t *testing.T) {
	themeNames := []string{
		"catppuccin", "dracula", "flexoki", "gruvbox", "monokai",
		"onedark", "opencode", "tokyonight", "tron",
	}

	for _, name := range themeNames {
		t.Run(name, func(t *testing.T) {
			theme := GetTheme(name)
			if theme == nil {
				t.Fatalf("Theme '%s' not registered", name)
			}

			tokens := theme.Tokens()

			// Verify all color tokens are non-nil
			assertTokenNotNil(t, "PanelBg", tokens.PanelBg)
			assertTokenNotNil(t, "PanelBorder", tokens.PanelBorder)
			assertTokenNotNil(t, "PanelBorderFocus", tokens.PanelBorderFocus)
			assertTokenNotNil(t, "PanelTitle", tokens.PanelTitle)

			assertTokenNotNil(t, "StatusRunning", tokens.StatusRunning)
			assertTokenNotNil(t, "StatusError", tokens.StatusError)
			assertTokenNotNil(t, "StatusWarning", tokens.StatusWarning)
			assertTokenNotNil(t, "StatusIdle", tokens.StatusIdle)

			assertTokenNotNil(t, "ProgressFill", tokens.ProgressFill)
			assertTokenNotNil(t, "ProgressEmpty", tokens.ProgressEmpty)
			assertTokenNotNil(t, "ProgressLabel", tokens.ProgressLabel)

			assertTokenNotNil(t, "ListCursor", tokens.ListCursor)
			assertTokenNotNil(t, "ListSelected", tokens.ListSelected)
			assertTokenNotNil(t, "ListItemNormal", tokens.ListItemNormal)

			assertTokenNotNil(t, "DataLabel", tokens.DataLabel)
			assertTokenNotNil(t, "DataValue", tokens.DataValue)
			assertTokenNotNil(t, "DataUnit", tokens.DataUnit)

			// Verify spacing
			if tokens.PanelPadding < 0 {
				t.Errorf("PanelPadding = %d, want >= 0", tokens.PanelPadding)
			}
		})
	}
}

// assertColorEquals checks that two TerminalColors are the same AdaptiveColor
func assertColorEquals(t *testing.T, name string, got, want lipgloss.TerminalColor) {
	t.Helper()
	gotAC, gotOK := got.(lipgloss.AdaptiveColor)
	wantAC, wantOK := want.(lipgloss.AdaptiveColor)
	if !gotOK || !wantOK {
		t.Errorf("%s: type mismatch got=%T want=%T", name, got, want)
		return
	}
	if gotAC.Dark != wantAC.Dark || gotAC.Light != wantAC.Light {
		t.Errorf("%s: got={Dark:%s,Light:%s} want={Dark:%s,Light:%s}",
			name, gotAC.Dark, gotAC.Light, wantAC.Dark, wantAC.Light)
	}
}

// assertTokenNotNil checks that a TerminalColor is not nil
func assertTokenNotNil(t *testing.T, name string, color lipgloss.TerminalColor) {
	t.Helper()
	if color == nil {
		t.Errorf("%s is nil", name)
	}
}
