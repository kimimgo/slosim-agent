package parametric

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// CaseResult represents a single parametric case for display.
type CaseResult struct {
	CaseID     string
	Parameters map[string]float64
	Status     string // "pending", "running", "completed", "failed"
	OutputPath string
}

// View is the Parametric View TUI component (TUI-05).
// Displays a comparison table with parameter columns for side-by-side results.
type View struct {
	StudyName string
	Cases     []CaseResult
	ParamKeys []string // ordered parameter names
	Cursor    int
	ColOffset int // horizontal scroll offset
	Width     int
	Height    int
	Selected  map[int]bool // multi-select for comparison
}

// NewView creates a new parametric comparison view.
func NewView(studyName string, cases []CaseResult) *View {
	keys := extractParamKeys(cases)
	return &View{
		StudyName: studyName,
		Cases:     cases,
		ParamKeys: keys,
		Cursor:    0,
		Selected:  make(map[int]bool),
	}
}

// SetCases updates the case list and recalculates parameter keys.
func (v *View) SetCases(cases []CaseResult) {
	v.Cases = cases
	v.ParamKeys = extractParamKeys(cases)
	if v.Cursor >= len(cases) {
		v.Cursor = max(0, len(cases)-1)
	}
}

// Update handles view updates.
func (v *View) Update(msg tea.Msg) (*View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Width = msg.Width
		v.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.Cursor > 0 {
				v.Cursor--
			}
		case "down", "j":
			if v.Cursor < len(v.Cases)-1 {
				v.Cursor++
			}
		case "left", "h":
			if v.ColOffset > 0 {
				v.ColOffset--
			}
		case "right", "l":
			maxOffset := len(v.ParamKeys) - 1
			if maxOffset < 0 {
				maxOffset = 0
			}
			if v.ColOffset < maxOffset {
				v.ColOffset++
			}
		case " ":
			// Toggle selection for comparison
			if v.Selected[v.Cursor] {
				delete(v.Selected, v.Cursor)
			} else {
				v.Selected[v.Cursor] = true
			}
		}
	case ParametricStudyUpdateMsg:
		v.SetCases(msg.Cases)
		v.StudyName = msg.StudyName
	}
	return v, nil
}

// GetSelectedCases returns the currently selected cases for comparison.
func (v *View) GetSelectedCases() []CaseResult {
	var selected []CaseResult
	for idx := range v.Selected {
		if idx < len(v.Cases) {
			selected = append(selected, v.Cases[idx])
		}
	}
	return selected
}

// View renders the parametric comparison table.
func (v *View) View() string {
	t := theme.CurrentTheme()

	if len(v.Cases) == 0 {
		return lipgloss.NewStyle().
			Foreground(t.TextMuted()).
			Render("파라메트릭 스터디 결과가 없습니다.")
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary()).
		Render(fmt.Sprintf("Parametric Study: %s", v.StudyName))

	summary := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render(fmt.Sprintf("총 %d 케이스 | 선택됨: %d | 변수: %d",
			len(v.Cases), len(v.Selected), len(v.ParamKeys)))

	table := v.renderTable(t)

	help := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render("j/k: 이동 | Space: 선택 | h/l: 열 스크롤 | Enter: 상세보기")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		summary,
		"",
		table,
		"",
		help,
	)
}

func (v *View) renderTable(t theme.Theme) string {
	colWidth := 14
	idWidth := 10

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.TextEmphasized()).
		Width(colWidth).
		Align(lipgloss.Center)

	idHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.TextEmphasized()).
		Width(idWidth).
		Align(lipgloss.Left)

	statusHeaderStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.TextEmphasized()).
		Width(colWidth).
		Align(lipgloss.Center)

	// Build header row
	headerCells := []string{idHeaderStyle.Render("Case")}

	// Show visible parameter columns
	maxCols := v.maxVisibleCols(idWidth, colWidth)
	endCol := min(v.ColOffset+maxCols, len(v.ParamKeys))

	for i := v.ColOffset; i < endCol; i++ {
		headerCells = append(headerCells, headerStyle.Render(v.ParamKeys[i]))
	}
	headerCells = append(headerCells, statusHeaderStyle.Render("Status"))

	header := strings.Join(headerCells, " ")

	// Separator
	sep := lipgloss.NewStyle().
		Foreground(t.BorderDim()).
		Render(strings.Repeat("─", lipgloss.Width(header)))

	// Build data rows
	var rows []string
	rows = append(rows, header)
	rows = append(rows, sep)

	maxRows := v.Height - 6
	if maxRows < 5 {
		maxRows = 15
	}
	startRow, endRow := v.visibleRange(maxRows)

	for i := startRow; i < endRow; i++ {
		c := v.Cases[i]
		isCursor := i == v.Cursor
		isSelected := v.Selected[i]

		row := v.renderRow(t, c, i, isCursor, isSelected, idWidth, colWidth, endCol)
		rows = append(rows, row)
	}

	// Scroll indicator
	if len(v.ParamKeys) > maxCols {
		indicator := lipgloss.NewStyle().
			Foreground(t.TextMuted()).
			Render(fmt.Sprintf("  [열 %d-%d / %d]", v.ColOffset+1, endCol, len(v.ParamKeys)))
		rows = append(rows, indicator)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (v *View) renderRow(t theme.Theme, c CaseResult, idx int, isCursor, isSelected bool, idWidth, colWidth, endCol int) string {
	// Cursor indicator
	prefix := "  "
	if isCursor && isSelected {
		prefix = ">*"
	} else if isCursor {
		prefix = "> "
	} else if isSelected {
		prefix = " *"
	}

	cellStyle := lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Right)
	idStyle := lipgloss.NewStyle().Width(idWidth).Align(lipgloss.Left)

	if isCursor {
		cellStyle = cellStyle.Foreground(t.Accent())
		idStyle = idStyle.Foreground(t.Accent()).Bold(true)
	} else if isSelected {
		cellStyle = cellStyle.Foreground(t.Info())
		idStyle = idStyle.Foreground(t.Info())
	}

	cells := []string{idStyle.Render(prefix + c.CaseID)}

	for i := v.ColOffset; i < endCol; i++ {
		key := v.ParamKeys[i]
		val, ok := c.Parameters[key]
		if ok {
			cells = append(cells, cellStyle.Render(fmt.Sprintf("%.4g", val)))
		} else {
			cells = append(cells, cellStyle.Render("-"))
		}
	}

	// Status column with color
	statusStyle := lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Center)
	switch c.Status {
	case "completed":
		statusStyle = statusStyle.Foreground(t.Success())
	case "running":
		statusStyle = statusStyle.Foreground(t.Info())
	case "failed":
		statusStyle = statusStyle.Foreground(t.Error())
	default:
		statusStyle = statusStyle.Foreground(t.TextMuted())
	}
	cells = append(cells, statusStyle.Render(c.Status))

	return strings.Join(cells, " ")
}

func (v *View) maxVisibleCols(idWidth, colWidth int) int {
	available := v.Width - idWidth - colWidth - 10 // reserve for id + status + margins
	if available <= 0 {
		available = 60
	}
	n := available / (colWidth + 1)
	if n < 1 {
		n = 1
	}
	return n
}

func (v *View) visibleRange(maxRows int) (int, int) {
	if len(v.Cases) <= maxRows {
		return 0, len(v.Cases)
	}

	half := maxRows / 2
	start := v.Cursor - half
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > len(v.Cases) {
		end = len(v.Cases)
		start = end - maxRows
	}
	return start, end
}

// extractParamKeys collects all unique parameter names from cases, sorted.
func extractParamKeys(cases []CaseResult) []string {
	keySet := make(map[string]bool)
	for _, c := range cases {
		for k := range c.Parameters {
			keySet[k] = true
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ParametricStudyUpdateMsg notifies the view of new study data.
type ParametricStudyUpdateMsg struct {
	StudyName string
	Cases     []CaseResult
}
