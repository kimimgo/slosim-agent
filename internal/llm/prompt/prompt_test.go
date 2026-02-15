package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContextFromPaths(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	_, err := config.Load(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	cfg := config.Get()
	cfg.WorkingDir = tmpDir
	cfg.ContextPaths = []string{
		"file.txt",
		"directory/",
	}
	testFiles := []string{
		"file.txt",
		"directory/file_a.txt",
		"directory/file_b.txt",
		"directory/file_c.txt",
	}

	createTestFiles(t, tmpDir, testFiles)

	context := getContextFromPaths()
	// File order from getContextFromPaths is non-deterministic (map/directory walk),
	// so check each file's content individually instead of exact string match.
	assert.Contains(t, context, fmt.Sprintf("# From:%s/file.txt\nfile.txt: test content", tmpDir))
	assert.Contains(t, context, fmt.Sprintf("# From:%s/directory/file_a.txt\ndirectory/file_a.txt: test content", tmpDir))
	assert.Contains(t, context, fmt.Sprintf("# From:%s/directory/file_b.txt\ndirectory/file_b.txt: test content", tmpDir))
	assert.Contains(t, context, fmt.Sprintf("# From:%s/directory/file_c.txt\ndirectory/file_c.txt: test content", tmpDir))
}

// AGENT-E2E: Sloshing system prompt structural quality tests

func TestSloshingPrompt_Structure(t *testing.T) {
	prompt := sloshingSystemPrompt

	t.Run("AGENT-E2E: prompt contains tool call ordering", func(t *testing.T) {
		assert.Contains(t, prompt, "xml_generator")
		assert.Contains(t, prompt, "gencase")
		assert.Contains(t, prompt, "solver")
		assert.Contains(t, prompt, "partvtk")
		assert.Contains(t, prompt, "measuretool")
		assert.Contains(t, prompt, "report")
	})

	t.Run("AGENT-E2E: prompt contains parameter auto-fill rules", func(t *testing.T) {
		assert.Contains(t, prompt, "dp")
		assert.Contains(t, prompt, "time_max")
		assert.Contains(t, prompt, "유체 높이") // fluid_height in Korean
		assert.Contains(t, prompt, "out_path")
	})

	t.Run("AGENT-E2E: prompt forbids error_recovery on new simulation", func(t *testing.T) {
		assert.Contains(t, prompt, "error_recovery")
		assert.Contains(t, prompt, "에러가 발생한 경우에만")
	})

	t.Run("AGENT-E2E: prompt has Korean user-facing language", func(t *testing.T) {
		assert.Contains(t, prompt, "한국어")
		assert.Contains(t, prompt, "사용자")
	})

	t.Run("AGENT-E2E: prompt token budget within Qwen3 limit", func(t *testing.T) {
		// Rough token estimate: 1 token ≈ 3-4 chars for Korean/mixed text
		// Qwen3 32B optimal: ≤ 8K tokens system prompt
		estimatedTokens := len(prompt) / 3
		assert.Less(t, estimatedTokens, 8000,
			"System prompt should be under ~8K tokens for Qwen3 32B (got ~%d)", estimatedTokens)
	})

	t.Run("AGENT-E2E: prompt mentions standard tank sizes", func(t *testing.T) {
		assert.Contains(t, prompt, "LNG")
		assert.Contains(t, prompt, "소형 탱크")
		assert.Contains(t, prompt, "실험 탱크")
	})

	t.Run("AGENT-E2E: prompt mentions domain formulas", func(t *testing.T) {
		assert.Contains(t, prompt, "공진 주파수")
		assert.Contains(t, prompt, "tanh")
	})

	t.Run("AGENT-E2E: prompt mentions v0.3 features", func(t *testing.T) {
		// v0.3: mDBC boundary method, cylindrical/L-shaped geometry, seismic input
		assert.Contains(t, prompt, "mDBC",
			"Prompt should mention mDBC boundary method (v0.3 feature)")
		assert.Contains(t, prompt, "원통형",
			"Prompt should mention cylindrical tank support")
		assert.Contains(t, prompt, "L",
			"Prompt should mention L-shaped tank")
		assert.Contains(t, prompt, "seismic_input",
			"Prompt should mention seismic_input tool")
		assert.Contains(t, prompt, "boundary_method",
			"Prompt should mention boundary_method parameter")
	})
}

func createTestFiles(t *testing.T, tmpDir string, testFiles []string) {
	t.Helper()
	for _, path := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if path[len(path)-1] == '/' {
			err := os.MkdirAll(fullPath, 0755)
			require.NoError(t, err)
		} else {
			dir := filepath.Dir(fullPath)
			err := os.MkdirAll(dir, 0755)
			require.NoError(t, err)
			err = os.WriteFile(fullPath, []byte(path+": test content"), 0644)
			require.NoError(t, err)
		}
	}
}
