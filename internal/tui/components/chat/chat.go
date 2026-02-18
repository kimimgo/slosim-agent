package chat

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/opencode-ai/opencode/internal/session"
	"github.com/opencode-ai/opencode/internal/tui/styles"
	"github.com/opencode-ai/opencode/internal/tui/theme"
	"github.com/opencode-ai/opencode/internal/version"
)

type SendMsg struct {
	Text        string
	Attachments []message.Attachment
}

type SessionSelectedMsg = session.Session

type SessionClearedMsg struct{}

type EditorFocusMsg bool

func header(width int) string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		logo(width),
		repo(width),
		"",
		cwd(width),
	)
}

func lspsConfigured(width int) string {
	cfg := config.Get()
	title := "LSP Configuration"
	title = ansi.Truncate(title, width, "…")

	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	lsps := baseStyle.
		Width(width).
		Foreground(t.Primary()).
		Bold(true).
		Render(title)

	// Get LSP names and sort them for consistent ordering
	var lspNames []string
	for name := range cfg.LSP {
		lspNames = append(lspNames, name)
	}
	sort.Strings(lspNames)

	var lspViews []string
	for _, name := range lspNames {
		lsp := cfg.LSP[name]
		lspName := baseStyle.
			Foreground(t.Text()).
			Render(fmt.Sprintf("• %s", name))

		cmd := lsp.Command
		cmd = ansi.Truncate(cmd, width-lipgloss.Width(lspName)-3, "…")

		lspPath := baseStyle.
			Foreground(t.TextMuted()).
			Render(fmt.Sprintf(" (%s)", cmd))

		lspViews = append(lspViews,
			baseStyle.
				Width(width).
				Render(
					lipgloss.JoinHorizontal(
						lipgloss.Left,
						lspName,
						lspPath,
					),
				),
		)
	}

	return baseStyle.
		Width(width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lsps,
				lipgloss.JoinVertical(
					lipgloss.Left,
					lspViews...,
				),
			),
		)
}

// slosimASCII is the ASCII art banner for slosim-agent.
// Displayed on the initial screen when terminal width >= 44 chars.
const slosimASCII = `     _           _
 ___| | ___  ___(_)_ __ ___
/ __| |/ _ \/ __| | '_ ` + "`" + ` _ \
\__ \ | (_) \__ \ | | | | | |
|___/_|\___/|___/_|_| |_| |_|`

// slosimWave is a decorative wave line rendered below the ASCII art.
const slosimWave = `~^~._.~^~._.~^~._.~^~._.~^~`

func logo(width int) string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	versionText := baseStyle.
		Foreground(t.TextMuted()).
		Render(version.Version)

	// Narrow terminal: single-line fallback
	if width < 44 {
		line := fmt.Sprintf("%s slosim %s", styles.OpenCodeIcon, version.Version)
		return baseStyle.Bold(true).Width(width).Render(line)
	}

	// Wide terminal: ASCII art banner
	artStyle := baseStyle.Bold(true).Foreground(t.Primary())
	waveStyle := baseStyle.Foreground(t.Accent())
	subtitleStyle := baseStyle.Foreground(t.TextMuted())

	var lines []string
	for _, line := range strings.Split(slosimASCII, "\n") {
		lines = append(lines, artStyle.Render(line))
	}
	lines = append(lines, waveStyle.Render(slosimWave))
	lines = append(lines, lipgloss.JoinHorizontal(
		lipgloss.Left,
		subtitleStyle.Render("  Sloshing Simulator "),
		versionText,
	))

	return baseStyle.Width(width).Render(
		lipgloss.JoinVertical(lipgloss.Left, lines...),
	)
}

func repo(width int) string {
	repo := "DualSPHysics GPU Sloshing Analysis Agent"
	t := theme.CurrentTheme()

	return styles.BaseStyle().
		Foreground(t.TextMuted()).
		Width(width).
		Render(repo)
}

func cwd(width int) string {
	cwd := fmt.Sprintf("cwd: %s", config.WorkingDirectory())
	t := theme.CurrentTheme()

	return styles.BaseStyle().
		Foreground(t.TextMuted()).
		Width(width).
		Render(cwd)
}

