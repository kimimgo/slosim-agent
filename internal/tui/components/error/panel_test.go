package errorpanel

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func makeErrors() []ErrorEntry {
	return []ErrorEntry{
		{
			JobID:    "job_001",
			Message:  "시뮬레이션 발산 감지",
			Severity: "error",
			FixActions: []string{
				"TimeStep 감소 (CFL 조정)",
				"점성 계수 증가",
			},
			Timestamp: "2026-02-14T10:00:00",
		},
		{
			JobID:      "job_001",
			Message:    "매우 긴 시뮬레이션 (10000+ 타임스텝)",
			Severity:   "warning",
			FixActions: nil,
			Timestamp:  "2026-02-14T10:01:00",
		},
		{
			JobID:    "job_002",
			Message:  "GPU 메모리 부족",
			Severity: "error",
			FixActions: []string{
				"파티클 수 감소 또는 dp 증가",
			},
			Timestamp: "2026-02-14T10:05:00",
		},
	}
}

func TestNewPanel(t *testing.T) {
	p := NewPanel()

	assert.Empty(t, p.Errors)
	assert.Equal(t, 0, p.Cursor)
	assert.Equal(t, 3, p.MaxRetries)
	assert.False(t, p.HasErrors())
}

func TestPanelErrorDetection(t *testing.T) {
	p := NewPanel()

	msg := ErrorDetectedMsg{
		Entry: ErrorEntry{
			JobID:    "job_001",
			Message:  "발산 감지",
			Severity: "error",
		},
		IsDivergent: true,
	}

	p, _ = p.Update(msg)
	assert.True(t, p.HasErrors())
	assert.Len(t, p.Errors, 1)
	assert.True(t, p.IsDivergent)
}

func TestPanelNavigation(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()
	p.Width = 80
	p.Height = 24

	// Down
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 1, p.Cursor)

	// Down again
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 2, p.Cursor)

	// Up
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 1, p.Cursor)

	// Boundary: can't go past end
	p.Cursor = 2
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 2, p.Cursor)
}

func TestPanelExpandCollapse(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()

	assert.False(t, p.Expanded)

	// Enter to expand
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, p.Expanded)

	// Enter again to collapse
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, p.Expanded)
}

func TestPanelRetry(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()

	// First retry
	p, cmd := p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	assert.Equal(t, 1, p.RetryCount)
	assert.NotNil(t, cmd)

	// Execute the command to verify it produces RetryRequestMsg
	msg := cmd()
	retryMsg, ok := msg.(RetryRequestMsg)
	assert.True(t, ok)
	assert.Equal(t, "job_001", retryMsg.JobID)

	// Retry max times
	p.RetryCount = 3
	p, cmd = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	assert.Equal(t, 3, p.RetryCount) // Should not increment
	assert.Nil(t, cmd)               // No command emitted
}

func TestPanelDismiss(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()
	assert.Len(t, p.Errors, 3)

	// Dismiss first error
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	assert.Len(t, p.Errors, 2)
	assert.Equal(t, 0, p.Cursor) // Stays at 0

	// Dismiss all
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	p, _ = p.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	assert.Len(t, p.Errors, 0)
	assert.False(t, p.HasErrors())
}

func TestPanelClear(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()
	p.IsDivergent = true
	p.RetryCount = 2
	p.Cursor = 1

	p, _ = p.Update(ErrorsClearedMsg{})
	assert.Empty(t, p.Errors)
	assert.Equal(t, 0, p.Cursor)
	assert.False(t, p.IsDivergent)
	assert.Equal(t, 0, p.RetryCount)
}

func TestPanelView(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()
	p.Width = 80
	p.Height = 24
	p.IsDivergent = true

	output := p.View()
	assert.Contains(t, output, "에러")
	assert.Contains(t, output, "경고")
	assert.Contains(t, output, "DIVERGENT")
}

func TestPanelViewNoErrors(t *testing.T) {
	p := NewPanel()
	output := p.View()
	assert.Contains(t, output, "에러 없음")
}

func TestPanelExpandedView(t *testing.T) {
	p := NewPanel()
	p.Errors = makeErrors()
	p.Width = 80
	p.Height = 24
	p.Expanded = true

	output := p.View()
	assert.Contains(t, output, "수정 제안")
	assert.Contains(t, output, "TimeStep")
}
