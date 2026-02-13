package simulation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashboard represents the Sim Dashboard TUI component (TUI-02).
// Displays detailed simulation status, progress, and logs.
type Dashboard struct {
	JobID       string
	Status      string
	ProgressPct float64
	Logs        []string
	Width       int
	Height      int
	progressBar progress.Model
}

// NewDashboard creates a new simulation dashboard.
func NewDashboard() *Dashboard {
	pb := progress.New(progress.WithDefaultGradient())
	return &Dashboard{
		Status:      "IDLE",
		ProgressPct: 0,
		Logs:        []string{},
		progressBar: pb,
	}
}

// Update handles dashboard updates.
func (d *Dashboard) Update(msg tea.Msg) (*Dashboard, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.Width = msg.Width
		d.Height = msg.Height
	case SimulationStatusMsg:
		d.Status = msg.Status
		d.ProgressPct = msg.ProgressPct
		d.Logs = append(d.Logs, msg.LogMessage)
	}
	return d, nil
}

// View renders the dashboard.
func (d *Dashboard) View() string {
	if d.Status == "IDLE" {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("시뮬레이션이 실행되지 않고 있습니다.")
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Render("🔄 시뮬레이션 대시보드")

	status := fmt.Sprintf("상태: %s | 진행률: %.1f%%", d.Status, d.ProgressPct)
	progressView := d.progressBar.ViewAs(d.ProgressPct / 100)

	logsTitle := lipgloss.NewStyle().
		Bold(true).
		Render("📋 로그:")

	logsView := ""
	maxLogs := 10
	start := len(d.Logs) - maxLogs
	if start < 0 {
		start = 0
	}
	for _, log := range d.Logs[start:] {
		logsView += fmt.Sprintf("  %s\n", log)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		status,
		progressView,
		"",
		logsTitle,
		logsView,
	)
}

// SimulationStatusMsg is a message for updating dashboard.
type SimulationStatusMsg struct {
	Status      string
	ProgressPct float64
	LogMessage  string
}
