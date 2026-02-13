package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorToolInfo(t *testing.T) {
	tool := NewMonitorTool()
	info := tool.Info()

	assert.Equal(t, MonitorToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "sim_dir")
	assert.Contains(t, info.Required, "sim_dir")
}

func TestMonitorToolRun_NoRunCSV(t *testing.T) {
	tool := NewMonitorTool()
	ctx := context.Background()

	tmpDir := t.TempDir()

	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-1",
		Name:  MonitorToolName,
		Input: `{"sim_dir": "` + tmpDir + `"}`,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content, "Run.csv")
}

func TestMonitorToolRun_ValidRunCSV(t *testing.T) {
	tool := NewMonitorTool()
	ctx := context.Background()

	tmpDir := t.TempDir()
	runCSV := filepath.Join(tmpDir, "Run.csv")

	// Create minimal Run.csv
	content := `# Run.csv
Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot
0.000;0;10000;8000;2000;0;0.0;0.0
0.100;100;10000;8000;2000;5;1.234;5.678
`
	err := os.WriteFile(runCSV, []byte(content), 0644)
	assert.NoError(t, err)

	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-2",
		Name:  MonitorToolName,
		Input: `{"sim_dir": "` + tmpDir + `"}`,
	})

	assert.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "모니터링 결과")
}
