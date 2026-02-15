package simulation

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
	"github.com/opencode-ai/opencode/internal/tui/widgets"
)

// Dashboard represents the Sim Dashboard TUI component (TUI-02).
// Displays detailed simulation status, progress, and logs.
// Redesigned with K9s pod monitor layout: status header, progress bar,
// metrics + log split, and keybinding hints.
type Dashboard struct {
	JobID       string
	Status      string
	ProgressPct float64
	Logs        []string
	Width       int
	Height      int
	progressBar progress.Model
	logViewport viewport.Model
}

// NewDashboard creates a new simulation dashboard.
func NewDashboard() *Dashboard {
	pb := progress.New(progress.WithDefaultGradient())
	vp := viewport.New(0, 0)
	return &Dashboard{
		Status:      "IDLE",
		ProgressPct: 0,
		Logs:        []string{},
		progressBar: pb,
		logViewport: vp,
	}
}

// Update handles dashboard updates.
func (d *Dashboard) Update(msg tea.Msg) (*Dashboard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.Width = msg.Width
		d.Height = msg.Height
		d.progressBar.Width = d.contentWidth()
		d.logViewport.Width = d.logPanelWidth()
		d.logViewport.Height = d.logPanelHeight()
	case SimulationStatusMsg:
		d.Status = msg.Status
		d.ProgressPct = msg.ProgressPct
		d.Logs = append(d.Logs, msg.LogMessage)
		d.syncLogViewport()
	}
	return d, nil
}

// View renders the dashboard using theme tokens and widgets.
func (d *Dashboard) View() string {
	t := d.getThemeOrFallback()

	if d.Status == "IDLE" {
		idleStyle := lipgloss.NewStyle().Foreground(t.StatusIdle)
		return idleStyle.Render("시뮬레이션이 실행되지 않고 있습니다.")
	}

	header := d.renderHeader(t)
	progressView := d.renderProgress(t)
	body := d.renderBody(t)
	footer := d.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", progressView, "", body, "", footer,
	)
}

// SimulationStatusMsg is a message for updating dashboard.
type SimulationStatusMsg struct {
	Status      string
	ProgressPct float64
	LogMessage  string
}

// --- Private rendering methods ---

func (d *Dashboard) renderHeader(t theme.SemanticTokens) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PanelTitle)

	title := titleStyle.Render("시뮬레이션 대시보드")
	badge := widgets.StatusBadge(d.Status)
	status := fmt.Sprintf("진행률: %.1f%%", d.ProgressPct)
	statusView := lipgloss.NewStyle().Foreground(t.DataValue).Render(status)

	jobInfo := ""
	if d.JobID != "" {
		jobLabel := lipgloss.NewStyle().Foreground(t.DataLabel).Render("Job:")
		jobValue := lipgloss.NewStyle().Foreground(t.DataValue).Bold(true).Render(d.JobID)
		jobInfo = jobLabel + " " + jobValue + "  "
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		title, "  ", jobInfo, badge, "  ", statusView,
	)
}

func (d *Dashboard) renderProgress(t theme.SemanticTokens) string {
	w := d.contentWidth()
	if w > 0 {
		d.progressBar.Width = w
	}
	return d.progressBar.ViewAs(d.ProgressPct / 100)
}

func (d *Dashboard) renderBody(t theme.SemanticTokens) string {
	metrics := d.renderMetrics()
	logs := d.renderLogs(t)

	metricsW := d.metricsPanelWidth()
	logsW := d.logPanelWidth()
	panelH := d.bodyPanelHeight()

	metricsPanel := widgets.Panel{
		Title:   "Metrics",
		Content: metrics,
		Width:   metricsW,
		Height:  panelH,
		Focused: false,
	}

	logsPanel := widgets.Panel{
		Title:   "로그",
		Content: logs,
		Width:   logsW,
		Height:  panelH,
		Focused: true,
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		metricsPanel.View(), " ", logsPanel.View(),
	)
}

func (d *Dashboard) renderMetrics() string {
	m := []widgets.Metric{
		{Label: "Status", Value: d.Status},
		{Label: "Progress", Value: fmt.Sprintf("%.1f%%", d.ProgressPct)},
		{Label: "Logs", Value: fmt.Sprintf("%d", len(d.Logs)), Unit: "entries"},
	}

	lines := make([]string, len(m))
	for i, metric := range m {
		lines[i] = metric.View()
	}
	return strings.Join(lines, "\n")
}

func (d *Dashboard) renderLogs(t theme.SemanticTokens) string {
	maxLogs := 10
	start := len(d.Logs) - maxLogs
	if start < 0 {
		start = 0
	}

	logStyle := lipgloss.NewStyle().Foreground(t.DataValue)
	lines := make([]string, 0, maxLogs)
	for _, log := range d.Logs[start:] {
		lines = append(lines, logStyle.Render(log))
	}
	return strings.Join(lines, "\n")
}

func (d *Dashboard) renderFooter() string {
	hints := []widgets.KeyHint{
		{Key: "q", Description: "Quit"},
		{Key: "p", Description: "Pause"},
		{Key: "l", Description: "Full Log"},
		{Key: "r", Description: "Refresh"},
	}
	return widgets.KeyHintBar(hints)
}

// --- Layout calculations ---

func (d *Dashboard) contentWidth() int {
	if d.Width <= 0 {
		return 60
	}
	return d.Width - 4 // border + padding
}

func (d *Dashboard) metricsPanelWidth() int {
	w := d.contentWidth()
	return w * 35 / 100 // 35% for metrics
}

func (d *Dashboard) logPanelWidth() int {
	w := d.contentWidth()
	return w - d.metricsPanelWidth() - 1 // rest for logs, -1 for gap
}

func (d *Dashboard) logPanelHeight() int {
	h := d.bodyPanelHeight()
	if h < 3 {
		return 3
	}
	return h - 2 // subtract border
}

func (d *Dashboard) bodyPanelHeight() int {
	if d.Height <= 0 {
		return 8
	}
	// header(1) + gap(1) + progress(1) + gap(1) + body + gap(1) + footer(1)
	return d.Height - 6
}

func (d *Dashboard) syncLogViewport() {
	content := strings.Join(d.Logs, "\n")
	d.logViewport.SetContent(content)
	d.logViewport.GotoBottom()
}

func (d *Dashboard) getThemeOrFallback() theme.SemanticTokens {
	t := theme.CurrentTheme()
	if t != nil {
		return t.Tokens()
	}
	// Fallback: neutral tokens if no theme loaded
	return theme.SemanticTokens{
		PanelBg:          lipgloss.AdaptiveColor{Dark: "#222", Light: "#eee"},
		PanelBorder:      lipgloss.AdaptiveColor{Dark: "#444", Light: "#ccc"},
		PanelBorderFocus: lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		PanelTitle:       lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		StatusRunning:    lipgloss.AdaptiveColor{Dark: "#0f0", Light: "#0a0"},
		StatusError:      lipgloss.AdaptiveColor{Dark: "#f00", Light: "#a00"},
		StatusWarning:    lipgloss.AdaptiveColor{Dark: "#ff0", Light: "#aa0"},
		StatusIdle:       lipgloss.AdaptiveColor{Dark: "#888", Light: "#666"},
		ProgressFill:     lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		ProgressEmpty:    lipgloss.AdaptiveColor{Dark: "#333", Light: "#ccc"},
		ProgressLabel:    lipgloss.AdaptiveColor{Dark: "#fff", Light: "#000"},
		ListCursor:       lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		ListSelected:     lipgloss.AdaptiveColor{Dark: "#f80", Light: "#a50"},
		ListItemNormal:   lipgloss.AdaptiveColor{Dark: "#ccc", Light: "#333"},
		DataLabel:        lipgloss.AdaptiveColor{Dark: "#888", Light: "#666"},
		DataValue:        lipgloss.AdaptiveColor{Dark: "#fff", Light: "#000"},
		DataUnit:         lipgloss.AdaptiveColor{Dark: "#666", Light: "#888"},
		PanelPadding:     1,
		PanelMargin:      0,
		SectionGap:       1,
	}
}
