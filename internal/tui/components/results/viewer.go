package results

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/tui/theme"
	"github.com/opencode-ai/opencode/internal/tui/widgets"
)

// Viewer represents the Result Viewer TUI component (TUI-03).
// Displays simulation results in a 2-panel layout: file list + preview.
type Viewer struct {
	SimDir string
	Files  []FileInfo
	Cursor int
	Width  int
	Height int
}

// FileInfo holds metadata about a result file.
type FileInfo struct {
	Name string
	Path string
	Type string // "report", "image", "csv", "vtk"
	Size int64
}

// NewViewer creates a new result viewer.
func NewViewer(simDir string) *Viewer {
	files := scanResultFiles(simDir)
	return &Viewer{
		SimDir: simDir,
		Files:  files,
		Cursor: 0,
	}
}

// Update handles viewer updates.
func (v *Viewer) Update(msg tea.Msg) (*Viewer, tea.Cmd) {
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
			if v.Cursor < len(v.Files)-1 {
				v.Cursor++
			}
		}
	}
	return v, nil
}

// View renders the result viewer using theme tokens and widgets.
func (v *Viewer) View() string {
	t := v.getTokens()

	if len(v.Files) == 0 {
		return lipgloss.NewStyle().
			Foreground(t.StatusIdle).
			Render("결과 파일이 없습니다.")
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PanelTitle).
		Render(fmt.Sprintf("결과 뷰어: %s", filepath.Base(v.SimDir)))

	fileList := v.renderFileList(t)
	preview := v.renderPreview(t)

	listW := v.listPanelWidth()
	previewW := v.previewPanelWidth()
	panelH := v.bodyHeight()

	listPanel := widgets.Panel{
		Title:   "Files",
		Content: fileList,
		Width:   listW,
		Height:  panelH,
		Focused: true,
	}

	previewPanel := widgets.Panel{
		Title:   "Preview",
		Content: preview,
		Width:   previewW,
		Height:  panelH,
		Focused: false,
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		listPanel.View(), " ", previewPanel.View(),
	)

	footer := widgets.KeyHintBar([]widgets.KeyHint{
		{Key: "j/k", Description: "이동"},
		{Key: "Enter", Description: "열기"},
	})

	summary := lipgloss.NewStyle().
		Foreground(t.DataLabel).
		Render(fmt.Sprintf("%d files, %.1f KB total", len(v.Files), v.totalSizeKB()))

	return lipgloss.JoinVertical(lipgloss.Left,
		title, "", body, "", summary, footer,
	)
}

// renderFileList renders the file list with cursor.
func (v *Viewer) renderFileList(t theme.SemanticTokens) string {
	cursorStyle := lipgloss.NewStyle().Foreground(t.ListCursor).Bold(true)
	normalStyle := lipgloss.NewStyle().Foreground(t.ListItemNormal)

	var items []string
	for i, file := range v.Files {
		icon := getFileIcon(file.Type)
		sizeStr := fmt.Sprintf("%.1f KB", float64(file.Size)/1024)

		if i == v.Cursor {
			line := fmt.Sprintf("▶ %s %s (%s)", icon, file.Name, sizeStr)
			items = append(items, cursorStyle.Render(line))
		} else {
			line := fmt.Sprintf("  %s %s (%s)", icon, file.Name, sizeStr)
			items = append(items, normalStyle.Render(line))
		}
	}

	return strings.Join(items, "\n")
}

// renderPreview renders metadata for the selected file.
func (v *Viewer) renderPreview(t theme.SemanticTokens) string {
	if v.Cursor < 0 || v.Cursor >= len(v.Files) {
		return lipgloss.NewStyle().
			Foreground(t.StatusIdle).
			Render("No file selected")
	}

	f := v.Files[v.Cursor]
	metrics := []widgets.Metric{
		{Label: "Name", Value: f.Name},
		{Label: "Type", Value: f.Type},
		{Label: "Size", Value: fmt.Sprintf("%.1f", float64(f.Size)/1024), Unit: "KB"},
	}

	lines := make([]string, len(metrics))
	for i, m := range metrics {
		lines[i] = m.View()
	}
	return strings.Join(lines, "\n")
}

// --- Layout calculations ---

func (v *Viewer) listPanelWidth() int {
	w := v.contentWidth()
	return w * 40 / 100 // 40% for file list
}

func (v *Viewer) previewPanelWidth() int {
	w := v.contentWidth()
	return w - v.listPanelWidth() - 1 // rest for preview, -1 gap
}

func (v *Viewer) contentWidth() int {
	if v.Width <= 0 {
		return 60
	}
	return v.Width - 2
}

func (v *Viewer) bodyHeight() int {
	if v.Height <= 0 {
		return 8
	}
	return v.Height - 5 // title + gaps + summary + footer
}

func (v *Viewer) totalSizeKB() float64 {
	var total int64
	for _, f := range v.Files {
		total += f.Size
	}
	return float64(total) / 1024
}

func (v *Viewer) getTokens() theme.SemanticTokens {
	t := theme.CurrentTheme()
	if t != nil {
		return t.Tokens()
	}
	return theme.SemanticTokens{
		PanelBg:          lipgloss.AdaptiveColor{Dark: "#222", Light: "#eee"},
		PanelBorder:      lipgloss.AdaptiveColor{Dark: "#444", Light: "#ccc"},
		PanelBorderFocus: lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		PanelTitle:       lipgloss.AdaptiveColor{Dark: "#88f", Light: "#44a"},
		StatusIdle:       lipgloss.AdaptiveColor{Dark: "#888", Light: "#666"},
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

func scanResultFiles(simDir string) []FileInfo {
	var files []FileInfo

	patterns := map[string]string{
		"report.md":   "report",
		"analysis.md": "report",
		"*.png":       "image",
		"*.jpg":       "image",
		"csv/*.csv":   "csv",
		"vtk/*.vtk":   "vtk",
	}

	for pattern, fileType := range patterns {
		matches, _ := filepath.Glob(filepath.Join(simDir, pattern))
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			files = append(files, FileInfo{
				Name: filepath.Base(match),
				Path: match,
				Type: fileType,
				Size: info.Size(),
			})
		}
	}

	return files
}

func getFileIcon(fileType string) string {
	switch fileType {
	case "report":
		return "📄"
	case "image":
		return "🖼️ "
	case "csv":
		return "📊"
	case "vtk":
		return "🔬"
	default:
		return "📁"
	}
}
