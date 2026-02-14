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

func TestBrowserApplyFilter_StatusFilter(t *testing.T) {
	b := NewBrowser(makeResults())

	// Set status filter manually and apply
	b.Filter.Status = "completed"
	b.applyFilter()

	assert.Len(t, b.Filtered, 2) // sim_001 and sim_002
	assert.Equal(t, "sim_001", b.Filtered[0].ID)
	assert.Equal(t, "sim_002", b.Filtered[1].ID)

	// Change to failed
	b.Filter.Status = "failed"
	b.applyFilter()

	assert.Len(t, b.Filtered, 1) // sim_003 only
	assert.Equal(t, "sim_003", b.Filtered[0].ID)
}

func TestBrowserApplyFilter_QueryFilter(t *testing.T) {
	b := NewBrowser(makeResults())

	// Set query filter and apply
	b.Filter.Query = "benchmark"
	b.applyFilter()

	assert.Len(t, b.Filtered, 1) // "Sloshing Benchmark"
	assert.Equal(t, "sim_001", b.Filtered[0].ID)

	// Case insensitive search
	b.Filter.Query = "FREQ"
	b.applyFilter()

	assert.Len(t, b.Filtered, 1) // "High Frequency Run"
	assert.Equal(t, "sim_003", b.Filtered[0].ID)
}

func TestBrowserApplyFilter_TagFilter(t *testing.T) {
	b := NewBrowser(makeResults())

	// Filter by tag
	b.Filter.Tags = []string{"sloshing"}
	b.applyFilter()

	assert.Len(t, b.Filtered, 2) // sim_001 and sim_003 have "sloshing" tag
	assert.Equal(t, "sim_001", b.Filtered[0].ID)
	assert.Equal(t, "sim_003", b.Filtered[1].ID)

	// Filter by another tag
	b.Filter.Tags = []string{"parametric"}
	b.applyFilter()

	assert.Len(t, b.Filtered, 1) // sim_002
	assert.Equal(t, "sim_002", b.Filtered[0].ID)
}

func TestBrowserApplyFilter_DateFilter(t *testing.T) {
	b := NewBrowser(makeResults())

	// Filter by FromDate
	b.Filter.FromDate = time.Date(2026, 2, 11, 0, 0, 0, 0, time.UTC)
	b.applyFilter()

	assert.Len(t, b.Filtered, 3) // sim_002, sim_003, sim_004 (after Feb 11)

	// Filter by ToDate
	b.Filter.FromDate = time.Time{} // Reset FromDate
	b.Filter.ToDate = time.Date(2026, 2, 11, 23, 59, 59, 0, time.UTC)
	b.applyFilter()

	assert.Len(t, b.Filtered, 2) // sim_001, sim_002 (before Feb 11 end)
}

func TestBrowserApplyFilter_CombinedFilters(t *testing.T) {
	b := NewBrowser(makeResults())

	// Combine status and query
	b.Filter.Status = "completed"
	b.Filter.Query = "sloshing"
	b.applyFilter()

	assert.Len(t, b.Filtered, 1) // Only sim_001 (completed + contains "sloshing")
	assert.Equal(t, "sim_001", b.Filtered[0].ID)
}

func TestBrowserApplyFilter_CursorReset(t *testing.T) {
	b := NewBrowser(makeResults())
	b.Cursor = 3 // Last item

	// Apply filter that reduces results
	b.Filter.Status = "completed"
	b.applyFilter()

	// Cursor should be adjusted to valid range
	assert.LessOrEqual(t, b.Cursor, len(b.Filtered)-1)
	assert.GreaterOrEqual(t, b.Cursor, 0)
}
