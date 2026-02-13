package results

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func makeResults() []StoredResult {
	return []StoredResult{
		{
			ID: "sim_001", Name: "Sloshing Benchmark",
			Timestamp: time.Date(2026, 2, 10, 14, 0, 0, 0, time.UTC),
			Status: "completed", Tags: []string{"sloshing", "benchmark"},
			Duration: 120.5, Description: "Basic sloshing test",
		},
		{
			ID: "sim_002", Name: "Fill Ratio 0.6",
			Timestamp: time.Date(2026, 2, 11, 9, 30, 0, 0, time.UTC),
			Status: "completed", Tags: []string{"parametric"},
			Duration: 200.1, Description: "Parametric fill ratio case",
		},
		{
			ID: "sim_003", Name: "High Frequency Run",
			Timestamp: time.Date(2026, 2, 12, 16, 0, 0, 0, time.UTC),
			Status: "failed", Tags: []string{"sloshing", "high-freq"},
			Duration: 45.0, Description: "Near resonance test",
		},
		{
			ID: "sim_004", Name: "Guard Panel Test",
			Timestamp: time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC),
			Status: "running", Tags: []string{"guard"},
			Duration: 0, Description: "Testing guard panel effect",
		},
	}
}

func TestNewBrowser(t *testing.T) {
	b := NewBrowser(makeResults())

	assert.Len(t, b.AllResults, 4)
	assert.Len(t, b.Filtered, 4) // No filter applied
	assert.Equal(t, 0, b.Cursor)
}

func TestBrowserNavigation(t *testing.T) {
	b := NewBrowser(makeResults())
	b.Width = 80
	b.Height = 24

	// Down
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 1, b.Cursor)

	// Up
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 0, b.Cursor)

	// Boundary
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 0, b.Cursor)
}

func TestBrowserStatusFilter(t *testing.T) {
	b := NewBrowser(makeResults())

	// Press 't' to cycle to "completed"
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	assert.Equal(t, "completed", b.Filter.Status)
	assert.Len(t, b.Filtered, 2) // sim_001 and sim_002

	// Press 't' to cycle to "failed"
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	assert.Equal(t, "failed", b.Filter.Status)
	assert.Len(t, b.Filtered, 1) // sim_003

	// Press 't' to cycle to "running"
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	assert.Equal(t, "running", b.Filter.Status)
	assert.Len(t, b.Filtered, 1) // sim_004

	// Press 't' to cycle back to all
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")})
	assert.Equal(t, "", b.Filter.Status)
	assert.Len(t, b.Filtered, 4)
}

func TestBrowserSearchMode(t *testing.T) {
	b := NewBrowser(makeResults())
	b.Width = 80
	b.Height = 24

	// Enter search mode
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	assert.True(t, b.SearchMode)

	// Type "bench"
	for _, ch := range "bench" {
		b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, "bench", b.SearchBuf)
	assert.Len(t, b.Filtered, 1) // "Sloshing Benchmark"

	// Backspace
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	assert.Equal(t, "benc", b.SearchBuf)

	// Exit search
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.False(t, b.SearchMode)
}

func TestBrowserSearchByTag(t *testing.T) {
	b := NewBrowser(makeResults())

	// Search by tag
	b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")})
	for _, ch := range "guard" {
		b, _ = b.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Len(t, b.Filtered, 1) // "Guard Panel Test" matched by tag
	assert.Equal(t, "sim_004", b.Filtered[0].ID)
}

func TestBrowserGetSelectedResult(t *testing.T) {
	b := NewBrowser(makeResults())

	r, ok := b.GetSelectedResult()
	assert.True(t, ok)
	assert.Equal(t, "sim_001", r.ID)

	// Move down
	b.Cursor = 2
	r, ok = b.GetSelectedResult()
	assert.True(t, ok)
	assert.Equal(t, "sim_003", r.ID)
}

func TestBrowserEmptyResults(t *testing.T) {
	b := NewBrowser(nil)
	assert.Len(t, b.Filtered, 0)

	_, ok := b.GetSelectedResult()
	assert.False(t, ok)

	output := b.View()
	assert.Contains(t, output, "검색 결과가 없습니다")
}

func TestBrowserRender(t *testing.T) {
	b := NewBrowser(makeResults())
	b.Width = 80
	b.Height = 24

	output := b.View()
	assert.Contains(t, output, "Result Browser")
	assert.Contains(t, output, "4 / 4 결과")
}

func TestResultBrowserUpdateMsg(t *testing.T) {
	b := NewBrowser(nil)

	msg := ResultBrowserUpdateMsg{
		Results: []StoredResult{
			{ID: "new_001", Name: "New Result", Status: "completed"},
		},
	}

	b, _ = b.Update(msg)
	assert.Len(t, b.AllResults, 1)
	assert.Len(t, b.Filtered, 1)
}
