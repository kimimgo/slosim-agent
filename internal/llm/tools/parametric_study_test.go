package tools

import (
	"fmt"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PARA-01: Parametric Study Execution

func TestParametricStudyTool_Info(t *testing.T) {
	tool := NewParametricStudyTool()
	info := tool.Info()

	assert.Equal(t, ParametricStudyToolName, info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Parameters, "study_name")
	assert.Contains(t, info.Parameters, "base_case")
	assert.Contains(t, info.Parameters, "variables")
	assert.Contains(t, info.Required, "study_name")
	assert.Contains(t, info.Required, "base_case")
	assert.Contains(t, info.Required, "variables")
}

func TestParametricStudyTool_Run(t *testing.T) {
	t.Run("PARA-01: generates parameter combinations", func(t *testing.T) {
		// Create temporary base case
		tempDir, err := os.MkdirTemp("", "param_study_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		baseCasePath := filepath.Join(tempDir, "base.xml")
		baseContent := `<?xml version="1.0"?>
<case>
  <parameters>
    <fill_ratio value="0.5" />
  </parameters>
</case>
`
		err = os.WriteFile(baseCasePath, []byte(baseContent), 0644)
		require.NoError(t, err)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "test_study",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.6, 0.7}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ParametricStudyToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content, "파라메트릭 스터디 시작")
		assert.Contains(t, response.Content, "총 케이스 수: 3")
	})

	t.Run("PARA-01: validates base case existence", func(t *testing.T) {
		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "test",
			BaseCase:  "/nonexistent/case.xml",
			Variables: []ParametricVariable{
				{Name: "param", Values: []float64{1.0}},
			},
		}

		paramsJSON, err := json.Marshal(params)
		require.NoError(t, err)

		call := ToolCall{
			Name:  ParametricStudyToolName,
			Input: string(paramsJSON),
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content, "베이스 케이스를 찾을 수 없습니다")
	})

	t.Run("PARA-01: handles invalid JSON", func(t *testing.T) {
		tool := NewParametricStudyTool()
		call := ToolCall{
			Name:  ParametricStudyToolName,
			Input: "invalid json",
		}

		response, err := tool.Run(context.Background(), call)
		require.NoError(t, err)
		assert.True(t, response.IsError)
	})
}

func TestGenerateCombinations(t *testing.T) {
	t.Run("generates single variable combinations", func(t *testing.T) {
		variables := []ParametricVariable{
			{Name: "param1", Values: []float64{1.0, 2.0, 3.0}},
		}

		combos := generateCombinations(variables)
		assert.Len(t, combos, 3)

		assert.Equal(t, 1.0, combos[0]["param1"])
		assert.Equal(t, 2.0, combos[1]["param1"])
		assert.Equal(t, 3.0, combos[2]["param1"])
	})

	t.Run("generates multi-variable combinations (Cartesian product)", func(t *testing.T) {
		variables := []ParametricVariable{
			{Name: "fill_ratio", Values: []float64{0.5, 0.7}},
			{Name: "frequency", Values: []float64{0.5, 1.0}},
		}

		combos := generateCombinations(variables)
		assert.Len(t, combos, 4) // 2 * 2

		// Check all combinations exist
		foundCombos := make(map[string]bool)
		for _, combo := range combos {
			key := fmt.Sprintf("%.1f_%.1f", combo["fill_ratio"], combo["frequency"])
			foundCombos[key] = true
		}

		assert.True(t, foundCombos["0.5_0.5"])
		assert.True(t, foundCombos["0.5_1.0"])
		assert.True(t, foundCombos["0.7_0.5"])
		assert.True(t, foundCombos["0.7_1.0"])
	})

	t.Run("handles empty variables", func(t *testing.T) {
		variables := []ParametricVariable{}
		combos := generateCombinations(variables)
		assert.Len(t, combos, 0)
	})
}

func TestModifyBaseCase(t *testing.T) {
	t.Run("modifies base case with parameters", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "modify_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		baseCasePath := filepath.Join(tempDir, "base.xml")
		baseContent := `<case><param value="{{fill_ratio}}" /></case>`
		err = os.WriteFile(baseCasePath, []byte(baseContent), 0644)
		require.NoError(t, err)

		params := map[string]float64{
			"fill_ratio": 0.65,
		}

		modified := modifyBaseCase(baseCasePath, params)
		assert.NotEmpty(t, modified)
		assert.Contains(t, modified, "fill_ratio = 0.650")
	})
}
