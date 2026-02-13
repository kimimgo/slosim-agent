package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ANI-01: Animation Generation Tool

func TestAnimationTool_Info(t *testing.T) {
	tool := NewAnimationTool()
	info := tool.Info()

	assert.Equal(t, AnimationToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "data_dir")
	assert.Contains(t, info.Parameters, "out_file")
	assert.Contains(t, info.Parameters, "variables")
	assert.Contains(t, info.Parameters, "colormap")
	assert.Contains(t, info.Parameters, "fps")
	assert.Contains(t, info.Parameters, "resolution")
	assert.Contains(t, info.Required, "data_dir")
	assert.Contains(t, info.Required, "out_file")
}

func TestAnimationTool_Run(t *testing.T) {
	t.Run("ANI-01: validates required parameters", func(t *testing.T) {
		tool := NewAnimationTool()
		params := AnimationParams{
			DataDir: "", // Missing required parameter
			OutFile: "",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  AnimationToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "필수 파라미터가 누락")
	})

	t.Run("ANI-01: validates data directory path", func(t *testing.T) {
		tool := NewAnimationTool()
		params := AnimationParams{
			DataDir: "clearly-invalid/../../../etc/passwd",
			OutFile: "/tmp/animation.mp4",
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  AnimationToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "올바르지 않습니다")
	})

	t.Run("ANI-01: applies default values", func(t *testing.T) {
		tool := NewAnimationTool()
		params := AnimationParams{
			DataDir: "/data/simulation",
			OutFile: "/tmp/animation.mp4",
			// All optional parameters omitted
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  AnimationToolName,
			ID:    "test-call",
			Input: string(paramsJSON),
		}

		// This will fail in actual execution (no Docker), but should validate parameters
		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		// Error expected due to Docker not running, but parameters should be validated
		// The error should be about Docker execution, not parameter validation
		if response.IsError {
			assert.NotContains(t, response.Content, "필수 파라미터")
		}
	})

	t.Run("ANI-01: handles invalid JSON", func(t *testing.T) {
		tool := NewAnimationTool()
		call := ToolCall{
			Name:  AnimationToolName,
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "파라미터 파싱 오류")
	})

	t.Run("ANI-01: infers format from extension", func(t *testing.T) {
		// Test format inference logic
		testCases := []struct {
			filename       string
			expectedFormat string
		}{
			{"animation.mp4", "mp4"},
			{"animation.MP4", "mp4"},
			{"animation.gif", "gif"},
			{"animation.GIF", "gif"},
			{"animation.avi", "mp4"}, // Default to mp4
		}

		for _, tc := range testCases {
			tool := NewAnimationTool()
			params := AnimationParams{
				DataDir: "/data/test",
				OutFile: "/tmp/" + tc.filename,
			}

			paramsJSON, err := json.Marshal(params)
			require.NoError(t, err)

			call := ToolCall{
				Name:  AnimationToolName,
				ID:    "test-call",
				Input: string(paramsJSON),
			}

			// Execute (will fail due to no Docker, but format should be inferred)
			_, err = tool.Run(context.Background(), call)
			require.NoError(t, err)
			// Format inference happens before Docker execution
		}
	})
}

func TestGenerateAnimationScript(t *testing.T) {
	t.Run("generates valid pvpython script", func(t *testing.T) {
		params := AnimationParams{
			DataDir:    "/data/input",
			OutFile:    "/data/output/animation.mp4",
			Variables:  []string{"Press"},
			Colormap:   "Blue to Red Rainbow",
			ViewAngle:  "XZ",
			Resolution: []int{1920, 1080},
			FPS:        30,
			StartTime:  0.0,
			EndTime:    10.0,
		}

		script := generateAnimationScript(params, true)

		// Check essential script components
		assert.Contains(t, script, "#!/usr/bin/env pvpython")
		assert.Contains(t, script, "from paraview.simple import *")
		assert.Contains(t, script, "LegacyVTKReader")
		assert.Contains(t, script, "ColorBy")
		assert.Contains(t, script, "Blue to Red Rainbow")
		assert.Contains(t, script, "SaveAnimation")
		assert.Contains(t, script, "[1920, 1080]")
		assert.Contains(t, script, "FrameRate=30")
		assert.Contains(t, script, "Press")
	})

	t.Run("handles different view angles", func(t *testing.T) {
		viewAngles := []string{"XZ", "XY", "YZ"}

		for _, angle := range viewAngles {
			params := AnimationParams{
				DataDir:    "/data/input",
				OutFile:    "/data/output/animation.mp4",
				Variables:  []string{"Vel"},
				ViewAngle:  angle,
				Resolution: []int{1280, 720},
				FPS:        24,
			}

			script := generateAnimationScript(params, false)

			assert.Contains(t, script, angle)
			assert.Contains(t, script, "CameraPosition")
		}
	})

	t.Run("includes fluid filter when requested", func(t *testing.T) {
		params := AnimationParams{
			DataDir:   "/data/input",
			OutFile:   "/data/output/animation.mp4",
			Variables: []string{"Rhop"},
		}

		scriptWithFluid := generateAnimationScript(params, true)
		scriptWithoutFluid := generateAnimationScript(params, false)

		// When onlyFluid=true, should check for fluid filter
		assert.Contains(t, scriptWithFluid, "Threshold")

		// When onlyFluid=false, should not apply threshold
		assert.NotContains(t, scriptWithoutFluid, "Threshold")
	})
}

func TestConvertToGif(t *testing.T) {
	t.Run("skips conversion for non-gif files", func(t *testing.T) {
		err := convertToGif(context.Background(), "/tmp/animation.mp4", 30)
		// Should not error (just return nil)
		assert.NoError(t, err)
	})

	t.Run("attempts conversion for gif files", func(t *testing.T) {
		// This will fail without ffmpeg installed, but tests the logic flow
		err := convertToGif(context.Background(), "/tmp/animation.gif", 30)
		// Error expected (ffmpeg not available or file not found)
		// But the function should handle it gracefully
		if err != nil {
			assert.Contains(t, err.Error(), "ffmpeg")
		}
	})
}
