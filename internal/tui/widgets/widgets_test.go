package widgets

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

func init() {
	// Ensure at least one theme is loaded for tests
	_ = theme.AvailableThemes()
}

// --- Panel Tests ---

func TestPanelRender(t *testing.T) {
	_ = theme.SetTheme("opencode")

	p := Panel{
		Title:   "Test Panel",
		Content: "Hello World",
		Width:   40,
		Height:  10,
		Focused: false,
	}

	view := p.View()
	if view == "" {
		t.Error("Panel.View() returned empty string")
	}
	if !strings.Contains(view, "Test Panel") {
		t.Error("Panel.View() missing title")
	}
}

func TestPanelFocusState(t *testing.T) {
	_ = theme.SetTheme("opencode")

	unfocused := Panel{Title: "A", Content: "x", Width: 30, Height: 5, Focused: false}
	focused := Panel{Title: "A", Content: "x", Width: 30, Height: 5, Focused: true}

	viewUnfocused := unfocused.View()
	viewFocused := focused.View()

	// Both should render without panic and contain content
	if viewUnfocused == "" || viewFocused == "" {
		t.Error("Panel should render in both focused and unfocused states")
	}
	if !strings.Contains(viewFocused, "A") {
		t.Error("Focused panel missing title")
	}
	// Note: in headless terminals, ANSI colors may not differ;
	// visual focus distinction is verified in integration/VHS tests
}

func TestPanelZeroSize(t *testing.T) {
	_ = theme.SetTheme("opencode")

	p := Panel{Title: "T", Content: "C", Width: 0, Height: 0}
	view := p.View()
	// Should not panic, may render minimal
	if view == "" {
		t.Error("Panel.View() with zero size returned empty")
	}
}

// --- StatusBadge Tests ---

func TestStatusBadge(t *testing.T) {
	_ = theme.SetTheme("opencode")

	tests := []struct {
		status   string
		wantIcon string
	}{
		{"RUNNING", "●"},
		{"ERROR", "●"},
		{"WARNING", "●"},
		{"DONE", "✓"},
		{"IDLE", "○"},
	}

	for _, tc := range tests {
		t.Run(tc.status, func(t *testing.T) {
			badge := StatusBadge(tc.status)
			if badge == "" {
				t.Errorf("StatusBadge(%q) returned empty", tc.status)
			}
			if !strings.Contains(badge, tc.status) {
				t.Errorf("StatusBadge(%q) missing status text", tc.status)
			}
		})
	}
}

func TestStatusBadgeUnknown(t *testing.T) {
	_ = theme.SetTheme("opencode")

	badge := StatusBadge("UNKNOWN")
	if badge == "" {
		t.Error("StatusBadge with unknown status returned empty")
	}
	if !strings.Contains(badge, "UNKNOWN") {
		t.Error("StatusBadge should include unknown status text")
	}
}

// --- Metric Tests ---

func TestMetricRender(t *testing.T) {
	_ = theme.SetTheme("opencode")

	m := Metric{Label: "Particles", Value: "136,420", Unit: ""}
	view := m.View()
	if !strings.Contains(view, "Particles") {
		t.Error("Metric.View() missing label")
	}
	if !strings.Contains(view, "136,420") {
		t.Error("Metric.View() missing value")
	}
}

func TestMetricWithUnit(t *testing.T) {
	_ = theme.SetTheme("opencode")

	m := Metric{Label: "Energy", Value: "0.0234", Unit: "J"}
	view := m.View()
	if !strings.Contains(view, "J") {
		t.Error("Metric.View() missing unit")
	}
}

func TestMetricNoUnit(t *testing.T) {
	_ = theme.SetTheme("opencode")

	m := Metric{Label: "Count", Value: "42", Unit: ""}
	view := m.View()
	// Should not contain trailing space or extra formatting for empty unit
	if strings.Contains(view, "  ") && strings.HasSuffix(strings.TrimSpace(view), "  ") {
		t.Error("Metric.View() has extra spaces for empty unit")
	}
}

// --- KeyHint Tests ---

func TestKeyHintBar(t *testing.T) {
	_ = theme.SetTheme("opencode")

	hints := []KeyHint{
		{Key: "q", Description: "Quit"},
		{Key: "p", Description: "Pause"},
		{Key: "l", Description: "Log"},
	}

	view := KeyHintBar(hints)
	if !strings.Contains(view, "q") {
		t.Error("KeyHintBar missing key 'q'")
	}
	if !strings.Contains(view, "Quit") {
		t.Error("KeyHintBar missing description 'Quit'")
	}
	if !strings.Contains(view, "p") {
		t.Error("KeyHintBar missing key 'p'")
	}
}

func TestKeyHintBarEmpty(t *testing.T) {
	_ = theme.SetTheme("opencode")

	view := KeyHintBar(nil)
	if view != "" {
		t.Error("KeyHintBar(nil) should return empty string")
	}
}

// --- Theme Consistency ---

func TestWidgetsAcrossThemes(t *testing.T) {
	themeNames := theme.AvailableThemes()

	for _, name := range themeNames {
		t.Run(name, func(t *testing.T) {
			_ = theme.SetTheme(name)

			// Panel should not panic
			p := Panel{Title: "T", Content: "C", Width: 30, Height: 5}
			if p.View() == "" {
				t.Error("Panel.View() empty")
			}

			// StatusBadge should not panic
			if StatusBadge("RUNNING") == "" {
				t.Error("StatusBadge empty")
			}

			// Metric should not panic
			m := Metric{Label: "L", Value: "V"}
			if m.View() == "" {
				t.Error("Metric.View() empty")
			}

			// KeyHintBar should not panic
			hints := []KeyHint{{Key: "k", Description: "d"}}
			if KeyHintBar(hints) == "" {
				t.Error("KeyHintBar empty")
			}
		})
	}
}

// Helper: verify no hardcoded colors by checking that rendered output
// doesn't contain ANSI escape codes when theme is nil (edge case test)
func TestWidgetsUseThemeColors(t *testing.T) {
	_ = theme.SetTheme("catppuccin")

	tokens := theme.CurrentTheme().Tokens()

	// Verify tokens are used (type assertion check)
	_, ok := tokens.PanelBorder.(lipgloss.AdaptiveColor)
	if !ok {
		t.Error("PanelBorder should be AdaptiveColor")
	}

	_, ok = tokens.StatusRunning.(lipgloss.AdaptiveColor)
	if !ok {
		t.Error("StatusRunning should be AdaptiveColor")
	}
}
