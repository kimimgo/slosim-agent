package errorpanel

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/styles"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// ErrorEntry represents a single error with fix suggestions.
type ErrorEntry struct {
	JobID      string
	Message    string
	Severity   string // "error", "warning"
	FixActions []string
	Timestamp  string
}

// Panel is the Error Recovery Panel TUI component.
// Displays real-time error notifications with fix suggestions and retry capability.
type Panel struct {
	Errors     []ErrorEntry
	Cursor     int
	Width      int
	Height     int
	Expanded   bool   // show fix details for selected error
	IsDivergent bool
	RetryCount int
	MaxRetries int
}

// NewPanel creates a new error recovery panel.
func NewPanel() *Panel {
	return &Panel{
		Errors:     []ErrorEntry{},
		MaxRetries: 3,
	}
}

// Update handles panel updates.
func (p *Panel) Update(msg tea.Msg) (*Panel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.Width = msg.Width
		p.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if p.Cursor > 0 {
				p.Cursor--
			}
		case "down", "j":
			if p.Cursor < len(p.Errors)-1 {
				p.Cursor++
			}
		case "enter", "tab":
			p.Expanded = !p.Expanded
		case "r", "R":
			// Request retry — emit retry command
			if len(p.Errors) > 0 && p.RetryCount < p.MaxRetries {
				entry := p.Errors[p.Cursor]
				p.RetryCount++
				return p, retryCmd(entry.JobID)
			}
		case "d":
			// Dismiss selected error
			if len(p.Errors) > 0 {
				p.Errors = append(p.Errors[:p.Cursor], p.Errors[p.Cursor+1:]...)
				if p.Cursor >= len(p.Errors) {
					p.Cursor = max(0, len(p.Errors)-1)
				}
			}
		}
	case ErrorDetectedMsg:
		p.Errors = append(p.Errors, msg.Entry)
		p.IsDivergent = p.IsDivergent || msg.IsDivergent
	case ErrorsClearedMsg:
		p.Errors = nil
		p.Cursor = 0
		p.IsDivergent = false
		p.RetryCount = 0
	}
	return p, nil
}

// HasErrors returns true if there are any active errors.
func (p *Panel) HasErrors() bool {
	return len(p.Errors) > 0
}

// View renders the error panel.
func (p *Panel) View() string {
	t := theme.CurrentTheme()

	if len(p.Errors) == 0 {
		return lipgloss.NewStyle().
			Foreground(t.Success()).
			Render(fmt.Sprintf("%s 에러 없음", styles.CheckIcon))
	}

	// Title with error count
	errorCount := 0
	warnCount := 0
	for _, e := range p.Errors {
		if e.Severity == "error" {
			errorCount++
		} else {
			warnCount++
		}
	}

	titleParts := []string{}
	if errorCount > 0 {
		titleParts = append(titleParts,
			lipgloss.NewStyle().Foreground(t.Error()).Bold(true).
				Render(fmt.Sprintf("%s %d 에러", styles.ErrorIcon, errorCount)))
	}
	if warnCount > 0 {
		titleParts = append(titleParts,
			lipgloss.NewStyle().Foreground(t.Warning()).Bold(true).
				Render(fmt.Sprintf("%s %d 경고", styles.WarningIcon, warnCount)))
	}

	title := strings.Join(titleParts, "  ")

	if p.IsDivergent {
		title += lipgloss.NewStyle().
			Foreground(t.Error()).Bold(true).
			Render("  [DIVERGENT]")
	}

	// Error list
	listView := p.renderList(t)

	// Expanded fix suggestions
	fixView := ""
	if p.Expanded && p.Cursor < len(p.Errors) {
		fixView = p.renderFixSuggestions(t, p.Errors[p.Cursor])
	}

	// Retry status
	retryInfo := ""
	if p.RetryCount > 0 {
		retryInfo = lipgloss.NewStyle().
			Foreground(t.Info()).
			Render(fmt.Sprintf("재시도: %d/%d", p.RetryCount, p.MaxRetries))
	}

	help := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render("j/k: 이동 | Enter: 수정 제안 | r: 재시도 | d: 무시")

	parts := []string{title, "", listView}
	if fixView != "" {
		parts = append(parts, "", fixView)
	}
	if retryInfo != "" {
		parts = append(parts, retryInfo)
	}
	parts = append(parts, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (p *Panel) renderList(t theme.Theme) string {
	var lines []string

	maxVisible := p.Height - 8
	if maxVisible < 5 {
		maxVisible = 10
	}
	start, end := visibleRange(p.Cursor, len(p.Errors), maxVisible)

	for i := start; i < end; i++ {
		e := p.Errors[i]
		isCursor := i == p.Cursor

		icon := styles.WarningIcon
		color := t.Warning()
		if e.Severity == "error" {
			icon = styles.ErrorIcon
			color = t.Error()
		}

		cursor := " "
		if isCursor {
			cursor = ">"
		}

		msgStyle := lipgloss.NewStyle().Foreground(color)
		if isCursor {
			msgStyle = msgStyle.Bold(true)
		}

		fixHint := ""
		if len(e.FixActions) > 0 {
			fixHint = lipgloss.NewStyle().
				Foreground(t.TextMuted()).
				Render(fmt.Sprintf(" (%d fix)", len(e.FixActions)))
		}

		line := fmt.Sprintf("%s %s %s%s",
			cursor,
			lipgloss.NewStyle().Foreground(color).Render(icon),
			msgStyle.Render(e.Message),
			fixHint,
		)
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (p *Panel) renderFixSuggestions(t theme.Theme, entry ErrorEntry) string {
	if len(entry.FixActions) == 0 {
		return lipgloss.NewStyle().
			Foreground(t.TextMuted()).
			Render("  수정 제안 없음")
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Info()).
		Render("  수정 제안:")

	var suggestions []string
	suggestions = append(suggestions, header)
	for i, action := range entry.FixActions {
		suggestion := lipgloss.NewStyle().
			Foreground(t.Text()).
			Render(fmt.Sprintf("  %d. %s", i+1, action))
		suggestions = append(suggestions, suggestion)
	}

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderDim()).
		Padding(0, 1)

	return border.Render(lipgloss.JoinVertical(lipgloss.Left, suggestions...))
}

func visibleRange(cursor, total, maxRows int) (int, int) {
	if total <= maxRows {
		return 0, total
	}

	half := maxRows / 2
	start := cursor - half
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > total {
		end = total
		start = end - maxRows
	}
	return start, end
}

func retryCmd(jobID string) tea.Cmd {
	return func() tea.Msg {
		return RetryRequestMsg{JobID: jobID}
	}
}

// ErrorDetectedMsg notifies the panel of a new error.
type ErrorDetectedMsg struct {
	Entry       ErrorEntry
	IsDivergent bool
}

// ErrorsClearedMsg clears all errors (e.g., after successful retry).
type ErrorsClearedMsg struct{}

// RetryRequestMsg is emitted when user requests a retry.
type RetryRequestMsg struct {
	JobID string
}
