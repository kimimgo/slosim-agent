package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalysisToolInfo(t *testing.T) {
	tool := NewAnalysisTool()
	info := tool.Info()

	assert.Equal(t, AnalysisToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "sim_dir")
	assert.Contains(t, info.Parameters, "case_config")
	assert.Contains(t, info.Required, "sim_dir")
	assert.Contains(t, info.Required, "case_config")
}

func TestAnalysisToolRun_MissingParams(t *testing.T) {
	tool := NewAnalysisTool()
	ctx := context.Background()

	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-1",
		Name:  AnalysisToolName,
		Input: `{}`,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content, "시뮬레이션 디렉토리")
}

func TestAnalysisToolRun_ValidParams(t *testing.T) {
	tool := NewAnalysisTool()
	ctx := context.Background()

	tmpDir := t.TempDir()

	input := `{
		"sim_dir": "` + tmpDir + `",
		"case_config": {
			"tank_length": 1.0,
			"tank_width": 0.5,
			"tank_height": 0.6,
			"fluid_height": 0.3,
			"freq": 0.5,
			"amplitude": 0.05,
			"dp": 0.01,
			"time_max": 10.0
		}
	}`

	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-2",
		Name:  AnalysisToolName,
		Input: input,
	})

	assert.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "AI 물리적 해석")
	assert.Contains(t, resp.Content, "공진")
}
