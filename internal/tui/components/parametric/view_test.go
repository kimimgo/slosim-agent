package parametric

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func makeCases() []CaseResult {
	return []CaseResult{
		{CaseID: "case_001", Parameters: map[string]float64{"fill_ratio": 0.5, "frequency": 1.0}, Status: "completed"},
		{CaseID: "case_002", Parameters: map[string]float64{"fill_ratio": 0.6, "frequency": 1.0}, Status: "running"},
		{CaseID: "case_003", Parameters: map[string]float64{"fill_ratio": 0.7, "frequency": 1.5}, Status: "failed"},
		{CaseID: "case_004", Parameters: map[string]float64{"fill_ratio": 0.5, "frequency": 1.5}, Status: "pending"},
	}
}

func TestNewView(t *testing.T) {
	cases := makeCases()
	v := NewView("test_study", cases)

	assert.Equal(t, "test_study", v.StudyName)
	assert.Len(t, v.Cases, 4)
	assert.Equal(t, []string{"fill_ratio", "frequency"}, v.ParamKeys)
	assert.Equal(t, 0, v.Cursor)
}

func TestViewNavigation(t *testing.T) {
	v := NewView("nav_test", makeCases())
	v.Width = 80
	v.Height = 24

	// Move down
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 1, v.Cursor)

	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 2, v.Cursor)

	// Move up
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 1, v.Cursor)

	// Boundary check: can't go below 0
	v.Cursor = 0
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 0, v.Cursor)

	// Boundary check: can't go past last
	v.Cursor = 3
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 3, v.Cursor)
}

func TestViewSelection(t *testing.T) {
	v := NewView("sel_test", makeCases())

	// Select first case
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	assert.True(t, v.Selected[0])

	// Move to second and select
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	assert.True(t, v.Selected[1])

	// Both selected
	selected := v.GetSelectedCases()
	assert.Len(t, selected, 2)

	// Deselect first
	v.Cursor = 0
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
	assert.False(t, v.Selected[0])
	selected = v.GetSelectedCases()
	assert.Len(t, selected, 1)
}

func TestViewColumnScroll(t *testing.T) {
	v := NewView("scroll_test", makeCases())
	v.Width = 40 // Narrow width forces column scrolling

	// Scroll right
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	assert.Equal(t, 1, v.ColOffset)

	// Scroll left
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	assert.Equal(t, 0, v.ColOffset)

	// Can't scroll left past 0
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	assert.Equal(t, 0, v.ColOffset)
}

func TestExtractParamKeys(t *testing.T) {
	cases := []CaseResult{
		{Parameters: map[string]float64{"z_var": 1.0, "a_var": 2.0}},
		{Parameters: map[string]float64{"a_var": 3.0, "m_var": 4.0}},
	}

	keys := extractParamKeys(cases)
	assert.Equal(t, []string{"a_var", "m_var", "z_var"}, keys)
}

func TestSetCases(t *testing.T) {
	v := NewView("test", makeCases())
	v.Cursor = 3

	newCases := []CaseResult{
		{CaseID: "new_001", Parameters: map[string]float64{"dp": 0.01}, Status: "pending"},
	}
	v.SetCases(newCases)

	assert.Len(t, v.Cases, 1)
	assert.Equal(t, 0, v.Cursor) // Reset since old cursor was out of range
	assert.Equal(t, []string{"dp"}, v.ParamKeys)
}

func TestViewRender(t *testing.T) {
	v := NewView("render_test", makeCases())
	v.Width = 80
	v.Height = 24

	output := v.View()
	assert.Contains(t, output, "Parametric Study: render_test")
	assert.Contains(t, output, "4 케이스")
	assert.Contains(t, output, "Case")
}

func TestEmptyView(t *testing.T) {
	v := NewView("empty", nil)
	output := v.View()
	assert.Contains(t, output, "결과가 없습니다")
}

func TestParametricStudyUpdateMsg(t *testing.T) {
	v := NewView("old", nil)

	msg := ParametricStudyUpdateMsg{
		StudyName: "new_study",
		Cases: []CaseResult{
			{CaseID: "c1", Parameters: map[string]float64{"x": 1}, Status: "completed"},
		},
	}

	v, _ = v.Update(msg)
	assert.Equal(t, "new_study", v.StudyName)
	assert.Len(t, v.Cases, 1)
}
