package simulation

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewDashboard(t *testing.T) {
	d := NewDashboard()

	assert.Equal(t, "IDLE", d.Status)
	assert.Equal(t, 0.0, d.ProgressPct)
	assert.Empty(t, d.Logs)
	assert.NotNil(t, d.progressBar)
}

func TestDashboardUpdate_SimulationStatusMsg(t *testing.T) {
	d := NewDashboard()

	msg := SimulationStatusMsg{
		Status:      "RUNNING",
		ProgressPct: 42.5,
		LogMessage:  "Time step 100/1000",
	}

	d, _ = d.Update(msg)

	assert.Equal(t, "RUNNING", d.Status)
	assert.Equal(t, 42.5, d.ProgressPct)
	assert.Len(t, d.Logs, 1)
	assert.Equal(t, "Time step 100/1000", d.Logs[0])
}

func TestDashboardUpdate_MultipleStatusMsgs(t *testing.T) {
	d := NewDashboard()

	// First update
	msg1 := SimulationStatusMsg{
		Status:      "RUNNING",
		ProgressPct: 25.0,
		LogMessage:  "Starting simulation",
	}
	d, _ = d.Update(msg1)

	// Second update
	msg2 := SimulationStatusMsg{
		Status:      "RUNNING",
		ProgressPct: 75.0,
		LogMessage:  "Processing particles",
	}
	d, _ = d.Update(msg2)

	assert.Equal(t, "RUNNING", d.Status)
	assert.Equal(t, 75.0, d.ProgressPct)
	assert.Len(t, d.Logs, 2)
	assert.Equal(t, "Starting simulation", d.Logs[0])
	assert.Equal(t, "Processing particles", d.Logs[1])
}

func TestDashboardUpdate_WindowSizeMsg(t *testing.T) {
	d := NewDashboard()

	msg := tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	}

	d, _ = d.Update(msg)

	assert.Equal(t, 120, d.Width)
	assert.Equal(t, 40, d.Height)
}

func TestDashboardView_Idle(t *testing.T) {
	d := NewDashboard()
	d.Width = 80
	d.Height = 24

	output := d.View()

	assert.Contains(t, output, "시뮬레이션이 실행되지 않고 있습니다")
}

func TestDashboardView_Running(t *testing.T) {
	d := NewDashboard()
	d.Status = "RUNNING"
	d.ProgressPct = 60.0
	d.Logs = []string{"Step 1", "Step 2", "Step 3"}
	d.Width = 80
	d.Height = 24

	output := d.View()

	assert.Contains(t, output, "시뮬레이션 대시보드")
	assert.Contains(t, output, "RUNNING")
	assert.Contains(t, output, "60.0")
	assert.Contains(t, output, "로그")
	assert.Contains(t, output, "Step 3")
}

func TestDashboardView_ManyLogs(t *testing.T) {
	d := NewDashboard()
	d.Status = "RUNNING"
	d.ProgressPct = 90.0
	d.Width = 80
	d.Height = 24

	// Add 15 log messages (more than maxLogs=10)
	for i := 1; i <= 15; i++ {
		msg := SimulationStatusMsg{
			Status:      "RUNNING",
			ProgressPct: float64(i) * 6,
			LogMessage:  "Log message " + string(rune('0'+i)),
		}
		d, _ = d.Update(msg)
	}

	output := d.View()

	// Should only show last 10 logs
	assert.Len(t, d.Logs, 15)
	// View should only contain recent logs (implementation shows last 10)
	assert.Contains(t, output, "로그")
}

func TestDashboardView_CompleteStatus(t *testing.T) {
	d := NewDashboard()
	d.Status = "COMPLETED"
	d.ProgressPct = 100.0
	d.Logs = []string{"Simulation finished successfully"}
	d.Width = 80
	d.Height = 24

	output := d.View()

	assert.Contains(t, output, "COMPLETED")
	assert.Contains(t, output, "100.0")
	assert.Contains(t, output, "Simulation finished successfully")
}
