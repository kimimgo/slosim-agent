package results

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// StoredResult represents a simulation result entry for the browser.
type StoredResult struct {
	ID          string
	Name        string
	Timestamp   time.Time
	Status      string // "completed", "failed", "running"
	Tags        []string
	Duration    float64
	Description string
}

// Render implements SimpleListItem-like rendering for the browser.
func (r StoredResult) Render(selected bool, width int) string {
	t := theme.CurrentTheme()

	statusIcon := "  "
	statusColor := t.TextMuted()
	switch r.Status {
	case "completed":
		statusIcon = "OK"
		statusColor = t.Success()
	case "failed":
		statusIcon = "!!"
		statusColor = t.Error()
	case "running":
		statusIcon = ".."
		statusColor = t.Info()
	}

	cursor := " "
	if selected {
		cursor = ">"
	}

	nameStyle := lipgloss.NewStyle()
	if selected {
		nameStyle = nameStyle.Foreground(t.Accent()).Bold(true)
	}

	statusStyle := lipgloss.NewStyle().Foreground(statusColor)
	tagStyle := lipgloss.NewStyle().Foreground(t.TextMuted())
	dateStyle := lipgloss.NewStyle().Foreground(t.TextMuted())

	tagStr := ""
	if len(r.Tags) > 0 {
		tagStr = tagStyle.Render(" [" + strings.Join(r.Tags, ",") + "]")
	}

	dateStr := dateStyle.Render(r.Timestamp.Format("01/02 15:04"))

	line := fmt.Sprintf("%s %s %s%s  %s",
		cursor,
		statusStyle.Render(statusIcon),
		nameStyle.Render(r.Name),
		tagStr,
		dateStr,
	)

	return line
}

// Filter holds the current filter settings for the result browser.
type Filter struct {
	Tags     []string
	Status   string
	FromDate time.Time
	ToDate   time.Time
	Query    string // free-text search
}

// Browser is the Result Browser TUI component.
// Displays a searchable, filterable list of saved simulation results.
type Browser struct {
	AllResults []StoredResult
	Filtered   []StoredResult
	Filter     Filter
	Cursor     int
	Width      int
	Height     int
	SearchMode bool
	SearchBuf  string
}

// NewBrowser creates a new result browser.
func NewBrowser(results []StoredResult) *Browser {
	b := &Browser{
		AllResults: results,
		Cursor:     0,
	}
	b.applyFilter()
	return b
}

// SetResults replaces all results and re-applies the current filter.
func (b *Browser) SetResults(results []StoredResult) {
	b.AllResults = results
	b.applyFilter()
}

// Update handles browser updates.
func (b *Browser) Update(msg tea.Msg) (*Browser, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.Width = msg.Width
		b.Height = msg.Height
	case tea.KeyMsg:
		if b.SearchMode {
			return b.updateSearchMode(msg)
		}
		return b.updateNormalMode(msg)
	case ResultBrowserUpdateMsg:
		b.SetResults(msg.Results)
	}
	return b, nil
}

func (b *Browser) updateNormalMode(msg tea.KeyMsg) (*Browser, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if b.Cursor > 0 {
			b.Cursor--
		}
	case "down", "j":
		if b.Cursor < len(b.Filtered)-1 {
			b.Cursor++
		}
	case "/":
		b.SearchMode = true
		b.SearchBuf = ""
	case "t":
		// Cycle status filter: all -> completed -> failed -> running -> all
		switch b.Filter.Status {
		case "":
			b.Filter.Status = "completed"
		case "completed":
			b.Filter.Status = "failed"
		case "failed":
			b.Filter.Status = "running"
		default:
			b.Filter.Status = ""
		}
		b.applyFilter()
	}
	return b, nil
}

func (b *Browser) updateSearchMode(msg tea.KeyMsg) (*Browser, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc":
		b.SearchMode = false
	case "backspace":
		if len(b.SearchBuf) > 0 {
			b.SearchBuf = b.SearchBuf[:len(b.SearchBuf)-1]
		}
		b.Filter.Query = b.SearchBuf
		b.applyFilter()
	default:
		if len(msg.String()) == 1 {
			b.SearchBuf += msg.String()
			b.Filter.Query = b.SearchBuf
			b.applyFilter()
		}
	}
	return b, nil
}

func (b *Browser) applyFilter() {
	b.Filtered = nil
	for _, r := range b.AllResults {
		if b.Filter.Status != "" && r.Status != b.Filter.Status {
			continue
		}
		if b.Filter.Query != "" {
			q := strings.ToLower(b.Filter.Query)
			matched := strings.Contains(strings.ToLower(r.Name), q) ||
				strings.Contains(strings.ToLower(r.Description), q) ||
				strings.Contains(strings.ToLower(r.ID), q)
			if !matched {
				// Check tags
				tagMatch := false
				for _, tag := range r.Tags {
					if strings.Contains(strings.ToLower(tag), q) {
						tagMatch = true
						break
					}
				}
				if !tagMatch {
					continue
				}
			}
		}
		if !b.Filter.FromDate.IsZero() && r.Timestamp.Before(b.Filter.FromDate) {
			continue
		}
		if !b.Filter.ToDate.IsZero() && r.Timestamp.After(b.Filter.ToDate) {
			continue
		}
		if len(b.Filter.Tags) > 0 {
			hasTag := false
			for _, ft := range b.Filter.Tags {
				for _, rt := range r.Tags {
					if rt == ft {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		b.Filtered = append(b.Filtered, r)
	}

	if b.Cursor >= len(b.Filtered) {
		b.Cursor = max(0, len(b.Filtered)-1)
	}
}

// GetSelectedResult returns the currently highlighted result.
func (b *Browser) GetSelectedResult() (StoredResult, bool) {
	if b.Cursor >= 0 && b.Cursor < len(b.Filtered) {
		return b.Filtered[b.Cursor], true
	}
	return StoredResult{}, false
}

// View renders the result browser.
func (b *Browser) View() string {
	t := theme.CurrentTheme()

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary()).
		Render("Result Browser")

	// Filter bar
	filterParts := []string{}
	if b.Filter.Status != "" {
		filterParts = append(filterParts, fmt.Sprintf("status:%s", b.Filter.Status))
	}
	if b.Filter.Query != "" {
		filterParts = append(filterParts, fmt.Sprintf("query:\"%s\"", b.Filter.Query))
	}
	if len(b.Filter.Tags) > 0 {
		filterParts = append(filterParts, fmt.Sprintf("tags:%s", strings.Join(b.Filter.Tags, ",")))
	}

	filterBar := ""
	if len(filterParts) > 0 {
		filterBar = lipgloss.NewStyle().
			Foreground(t.Info()).
			Render("Filter: " + strings.Join(filterParts, " | "))
	}

	summary := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render(fmt.Sprintf("%d / %d 결과", len(b.Filtered), len(b.AllResults)))

	// Search mode indicator
	searchLine := ""
	if b.SearchMode {
		searchLine = lipgloss.NewStyle().
			Foreground(t.Warning()).
			Render(fmt.Sprintf("/ %s_", b.SearchBuf))
	}

	// Result list
	listView := b.renderList(t)

	help := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render("j/k: 이동 | /: 검색 | t: 상태 필터 | Enter: 열기")

	parts := []string{title}
	if filterBar != "" {
		parts = append(parts, filterBar)
	}
	parts = append(parts, summary)
	if searchLine != "" {
		parts = append(parts, searchLine)
	}
	parts = append(parts, "", listView, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (b *Browser) renderList(t theme.Theme) string {
	if len(b.Filtered) == 0 {
		return lipgloss.NewStyle().
			Foreground(t.TextMuted()).
			Render("  검색 결과가 없습니다.")
	}

	maxVisible := b.Height - 8
	if maxVisible < 5 {
		maxVisible = 15
	}

	start, end := b.visibleRange(maxVisible)

	var lines []string
	for i := start; i < end; i++ {
		r := b.Filtered[i]
		line := r.Render(i == b.Cursor, b.Width)
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (b *Browser) visibleRange(maxRows int) (int, int) {
	total := len(b.Filtered)
	if total <= maxRows {
		return 0, total
	}

	half := maxRows / 2
	start := b.Cursor - half
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

// ResultBrowserUpdateMsg provides new result data to the browser.
type ResultBrowserUpdateMsg struct {
	Results []StoredResult
}
