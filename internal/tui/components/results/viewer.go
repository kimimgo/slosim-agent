package results

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Viewer represents the Result Viewer TUI component (TUI-03).
// Displays simulation results (files, images, reports).
type Viewer struct {
	SimDir  string
	Files   []FileInfo
	Cursor  int
	Width   int
	Height  int
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

// View renders the result viewer.
func (v *Viewer) View() string {
	if len(v.Files) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("결과 파일이 없습니다.")
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Render(fmt.Sprintf("📁 결과 뷰어: %s", filepath.Base(v.SimDir)))

	var items []string
	for i, file := range v.Files {
		cursor := " "
		if i == v.Cursor {
			cursor = "▶"
		}
		icon := getFileIcon(file.Type)
		item := fmt.Sprintf("%s %s %s (%.1f KB)", cursor, icon, file.Name, float64(file.Size)/1024)
		if i == v.Cursor {
			item = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(item)
		}
		items = append(items, item)
	}

	fileList := strings.Join(items, "\n")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		fileList,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("↑/k: 위 | ↓/j: 아래 | Enter: 열기"),
	)
}

func scanResultFiles(simDir string) []FileInfo {
	var files []FileInfo

	// Scan for common result files
	patterns := map[string]string{
		"report.md":      "report",
		"analysis.md":    "report",
		"*.png":          "image",
		"*.jpg":          "image",
		"csv/*.csv":      "csv",
		"vtk/*.vtk":      "vtk",
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
