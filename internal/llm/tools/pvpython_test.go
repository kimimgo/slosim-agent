package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPvpythonToolInfo(t *testing.T) {
	tool := NewPvpythonTool()
	info := tool.Info()

	assert.Equal(t, PvpythonToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "out_file")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "out_file")
}

func TestPvpythonToolRun_MissingParams(t *testing.T) {
	tool := NewPvpythonTool()
	ctx := context.Background()

	// Missing data_dir
	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-1",
		Name:  PvpythonToolName,
		Input: `{"out_file": "/tmp/out.png"}`,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content, "필수 파라미터가 누락")
}

func TestPvpythonToolRun_InvalidPath(t *testing.T) {
	tool := NewPvpythonTool()
	ctx := context.Background()

	resp, err := tool.Run(ctx, ToolCall{
		ID:    "test-2",
		Name:  PvpythonToolName,
		Input: `{"data_dir": "/nonexistent/path", "out_file": "/tmp/out.png"}`,
	})

	assert.NoError(t, err)
	assert.True(t, resp.IsError)
	assert.Contains(t, resp.Content, "경로가 올바르지 않습니다")
}
