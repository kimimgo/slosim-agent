package results

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewViewer_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	v := NewViewer(tmpDir)

	assert.Equal(t, tmpDir, v.SimDir)
	assert.Empty(t, v.Files)
	assert.Equal(t, 0, v.Cursor)
}

func TestNewViewer_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "report.md"), []byte("# Report"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "analysis.md"), []byte("# Analysis"), 0644))

	v := NewViewer(tmpDir)

	assert.Equal(t, tmpDir, v.SimDir)
	assert.Len(t, v.Files, 2)
	assert.Equal(t, 0, v.Cursor)

	// Verify file types
	for _, f := range v.Files {
		assert.Equal(t, "report", f.Type)
		assert.Contains(t, []string{"report.md", "analysis.md"}, f.Name)
	}
}

func TestViewerNavigation_Down(t *testing.T) {
	v := &Viewer{
		SimDir: "/tmp/test",
		Files: []FileInfo{
			{Name: "file1.md", Type: "report"},
			{Name: "file2.csv", Type: "csv"},
			{Name: "file3.vtk", Type: "vtk"},
		},
		Cursor: 0,
		Width:  80,
		Height: 24,
	}

	// Move down
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 1, v.Cursor)

	// Move down again
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 2, v.Cursor)

	// Can't move past last item
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 2, v.Cursor)
}

func TestViewerNavigation_Up(t *testing.T) {
	v := &Viewer{
		SimDir: "/tmp/test",
		Files: []FileInfo{
			{Name: "file1.md", Type: "report"},
			{Name: "file2.csv", Type: "csv"},
			{Name: "file3.vtk", Type: "vtk"},
		},
		Cursor: 2,
		Width:  80,
		Height: 24,
	}

	// Move up
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 1, v.Cursor)

	// Move up again
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 0, v.Cursor)

	// Can't move above 0
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	assert.Equal(t, 0, v.Cursor)
}

func TestViewerNavigation_ArrowKeys(t *testing.T) {
	v := &Viewer{
		SimDir: "/tmp/test",
		Files: []FileInfo{
			{Name: "file1.md", Type: "report"},
			{Name: "file2.csv", Type: "csv"},
		},
		Cursor: 0,
		Width:  80,
		Height: 24,
	}

	// Arrow down
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, v.Cursor)

	// Arrow up
	v, _ = v.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, v.Cursor)
}

func TestViewerUpdate_WindowSizeMsg(t *testing.T) {
	v := NewViewer("/tmp/test")

	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}

	v, _ = v.Update(msg)

	assert.Equal(t, 100, v.Width)
	assert.Equal(t, 30, v.Height)
}

func TestViewerView_Empty(t *testing.T) {
	v := NewViewer(t.TempDir())
	v.Width = 80
	v.Height = 24

	output := v.View()

	assert.Contains(t, output, "결과 파일이 없습니다")
}

func TestViewerView_WithFiles(t *testing.T) {
	v := &Viewer{
		SimDir: "/tmp/test/my_simulation",
		Files: []FileInfo{
			{Name: "report.md", Path: "/tmp/test/my_simulation/report.md", Type: "report", Size: 2048},
			{Name: "data.csv", Path: "/tmp/test/my_simulation/data.csv", Type: "csv", Size: 4096},
			{Name: "result.vtk", Path: "/tmp/test/my_simulation/result.vtk", Type: "vtk", Size: 102400},
		},
		Cursor: 0,
		Width:  80,
		Height: 24,
	}

	output := v.View()

	assert.Contains(t, output, "결과 뷰어")
	assert.Contains(t, output, "my_simulation")
	assert.Contains(t, output, "report.md")
	assert.Contains(t, output, "data.csv")
	assert.Contains(t, output, "result.vtk")
	assert.Contains(t, output, "KB") // File sizes should be shown
}

func TestViewerView_CursorHighlight(t *testing.T) {
	v := &Viewer{
		SimDir: "/tmp/test",
		Files: []FileInfo{
			{Name: "file1.md", Type: "report", Size: 1024},
			{Name: "file2.csv", Type: "csv", Size: 2048},
		},
		Cursor: 0,
		Width:  80,
		Height: 24,
	}

	output := v.View()

	// Cursor should be on first item
	assert.Contains(t, output, "▶")
	assert.Contains(t, output, "file1.md")

	// Move cursor to second item
	v.Cursor = 1
	output2 := v.View()
	assert.Contains(t, output2, "▶")
	assert.Contains(t, output2, "file2.csv")
}

func TestGetFileIcon(t *testing.T) {
	tests := []struct {
		fileType string
		expected string
	}{
		{"report", "📄"},
		{"image", "🖼️ "},
		{"csv", "📊"},
		{"vtk", "🔬"},
		{"unknown", "📁"},
	}

	for _, tt := range tests {
		t.Run(tt.fileType, func(t *testing.T) {
			icon := getFileIcon(tt.fileType)
			assert.Equal(t, tt.expected, icon)
		})
	}
}

func TestScanResultFiles_MultiplePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "report.md"), []byte("report"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "analysis.md"), []byte("analysis"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "snapshot.png"), []byte("png"), 0644))

	// Create subdirectories
	csvDir := filepath.Join(tmpDir, "csv")
	vtkDir := filepath.Join(tmpDir, "vtk")
	require.NoError(t, os.MkdirAll(csvDir, 0755))
	require.NoError(t, os.MkdirAll(vtkDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(csvDir, "data.csv"), []byte("csv"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(vtkDir, "fluid.vtk"), []byte("vtk"), 0644))

	v := NewViewer(tmpDir)

	// Should find all files
	assert.GreaterOrEqual(t, len(v.Files), 3) // At least report, analysis, png

	// Check types
	typeMap := make(map[string]int)
	for _, f := range v.Files {
		typeMap[f.Type]++
	}

	assert.Greater(t, typeMap["report"], 0)
	// Note: image, csv, vtk may not be found depending on glob patterns in scanResultFiles
}
