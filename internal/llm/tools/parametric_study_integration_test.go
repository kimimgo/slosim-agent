//go:build integration

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// PARA-01 Integration Tests: Parametric Study Tool
// These tests verify multi-case generation, manifest creation, and XML
// modification at the integration level (beyond unit-level existing tests).
// ============================================================================

// --- Multi-variable Cartesian Product ---

func TestParametricStudy_Integration_MultiVariable(t *testing.T) {
	t.Run("PARA-01: 3-variable Cartesian product generates correct case count", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "multi_var_study",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.6, 0.7}},
				{Name: "frequency", Values: []float64{0.5, 1.0}},
				{Name: "amplitude", Values: []float64{0.02, 0.05}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		response := runParametricStudy(t, tool, params)

		// 3 * 2 * 2 = 12 cases
		assert.Contains(t, response.Content, "총 케이스 수: 12",
			"3-variable Cartesian product should produce 3*2*2=12 combinations")
		assert.False(t, response.IsError)
	})

	t.Run("PARA-01: 2-variable study creates all expected case directories", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		outputDir := filepath.Join(tempDir, "output")
		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "dir_check_study",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.7}},
				{Name: "frequency", Values: []float64{0.5, 1.0}},
			},
			OutputDir: outputDir,
		}

		runParametricStudy(t, tool, params)

		// Verify 4 case directories exist (2*2)
		for i := 1; i <= 4; i++ {
			caseDir := filepath.Join(outputDir, fmt.Sprintf("case_%03d", i))
			_, err := os.Stat(caseDir)
			assert.NoError(t, err, "case_%03d directory should exist", i)

			// Each directory should contain case.xml
			xmlPath := filepath.Join(caseDir, "case.xml")
			_, err = os.Stat(xmlPath)
			assert.NoError(t, err, "case_%03d/case.xml should exist", i)
		}

		// case_005 should NOT exist
		_, err := os.Stat(filepath.Join(outputDir, "case_005"))
		assert.True(t, os.IsNotExist(err), "case_005 should not exist for 4-case study")
	})
}

// --- Study Manifest ---

func TestParametricStudy_Integration_Manifest(t *testing.T) {
	t.Run("PARA-01: study manifest contains all required fields", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		outputDir := filepath.Join(tempDir, "output")
		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "manifest_test",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.6}},
			},
			OutputDir: outputDir,
		}

		runParametricStudy(t, tool, params)

		// Read and parse manifest
		manifestPath := filepath.Join(outputDir, "study_manifest.json")
		data, err := os.ReadFile(manifestPath)
		require.NoError(t, err, "manifest file should be readable")

		var manifest map[string]any
		err = json.Unmarshal(data, &manifest)
		require.NoError(t, err, "manifest should be valid JSON")

		assert.Equal(t, "manifest_test", manifest["study_name"])
		assert.Equal(t, baseCasePath, manifest["base_case"])
		assert.Equal(t, float64(2), manifest["total_cases"])

		// Verify results array
		results, ok := manifest["results"].([]any)
		require.True(t, ok, "results should be an array")
		assert.Len(t, results, 2)

		// Verify each result has correct structure
		for i, r := range results {
			result, ok := r.(map[string]any)
			require.True(t, ok, "each result should be an object")
			assert.Equal(t, fmt.Sprintf("case_%03d", i+1), result["case_id"])
			assert.Equal(t, "created", result["status"])
			assert.NotEmpty(t, result["output_path"])
		}
	})

	t.Run("PARA-01: manifest variables field preserves original variable definitions", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		outputDir := filepath.Join(tempDir, "output")
		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "var_preserve_test",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.6, 0.7}},
				{Name: "frequency", Values: []float64{0.8, 1.2}},
			},
			OutputDir: outputDir,
		}

		runParametricStudy(t, tool, params)

		data, err := os.ReadFile(filepath.Join(outputDir, "study_manifest.json"))
		require.NoError(t, err)

		var manifest map[string]any
		require.NoError(t, json.Unmarshal(data, &manifest))

		variables, ok := manifest["variables"].([]any)
		require.True(t, ok)
		assert.Len(t, variables, 2, "should preserve both variable definitions")
	})
}

// --- XML Modification (RED TESTS - Known bugs) ---

func TestParametricStudy_Integration_XMLModification(t *testing.T) {
	t.Run("RED/PARA-01: modified XML should contain actual parameter values in XML attributes", func(t *testing.T) {
		// BUG: modifyBaseCase only appends comments, doesn't actually modify XML attributes.
		// This test documents the expected behavior.
		tempDir := t.TempDir()

		baseCasePath := filepath.Join(tempDir, "base.xml")
		baseContent := `<?xml version="1.0"?>
<case>
  <parameters>
    <fill_ratio value="0.5" />
  </parameters>
</case>`
		require.NoError(t, os.WriteFile(baseCasePath, []byte(baseContent), 0644))

		params := map[string]float64{"fill_ratio": 0.7}
		modified := modifyBaseCase(baseCasePath, params)

		// Expected behavior: the value attribute should be changed from 0.5 to 0.7
		// Current behavior: only appends a comment, doesn't modify the attribute
		assert.Contains(t, modified, `value="0.7"`,
			"RED: modifyBaseCase should actually replace XML attribute values, not just append comments")
	})

	t.Run("RED/PARA-01: template placeholders should be substituted in XML", func(t *testing.T) {
		// BUG: placeholder variable is computed but never used for actual substitution
		tempDir := t.TempDir()

		baseCasePath := filepath.Join(tempDir, "template.xml")
		templateContent := `<?xml version="1.0"?>
<case>
  <parameters>
    <fill_ratio value="{{fill_ratio}}" />
    <frequency value="{{frequency}}" />
  </parameters>
</case>`
		require.NoError(t, os.WriteFile(baseCasePath, []byte(templateContent), 0644))

		params := map[string]float64{
			"fill_ratio": 0.65,
			"frequency":  1.2,
		}
		modified := modifyBaseCase(baseCasePath, params)

		// Template placeholders should be replaced
		assert.NotContains(t, modified, "{{fill_ratio}}",
			"RED: {{fill_ratio}} placeholder should be replaced with actual value")
		assert.NotContains(t, modified, "{{frequency}}",
			"RED: {{frequency}} placeholder should be replaced with actual value")
		assert.Contains(t, modified, "0.65",
			"RED: modified XML should contain the fill_ratio value 0.65")
		assert.Contains(t, modified, "1.2",
			"RED: modified XML should contain the frequency value 1.2")
	})

	t.Run("PARA-01: each case directory XML has unique parameter values", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		outputDir := filepath.Join(tempDir, "output")
		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "unique_xml_test",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.3, 0.5, 0.7}},
			},
			OutputDir: outputDir,
		}

		runParametricStudy(t, tool, params)

		// Read all 3 case XMLs and verify they differ
		xmlContents := make([]string, 3)
		for i := 0; i < 3; i++ {
			xmlPath := filepath.Join(outputDir, fmt.Sprintf("case_%03d", i+1), "case.xml")
			data, err := os.ReadFile(xmlPath)
			require.NoError(t, err)
			xmlContents[i] = string(data)
		}

		// Each case XML should be unique (at minimum, the appended comment differs)
		assert.NotEqual(t, xmlContents[0], xmlContents[1],
			"case_001 and case_002 XML should differ")
		assert.NotEqual(t, xmlContents[1], xmlContents[2],
			"case_002 and case_003 XML should differ")
	})
}

// --- Default Values & Edge Cases ---

func TestParametricStudy_Integration_Defaults(t *testing.T) {
	t.Run("PARA-01: default output_dir uses study name", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		// Change working directory temporarily
		origWd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(tempDir))
		defer os.Chdir(origWd)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "default_dir_study",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5}},
			},
			// OutputDir intentionally omitted
		}

		response := runParametricStudy(t, tool, params)
		assert.Contains(t, response.Content, "parametric_studies/default_dir_study",
			"default output_dir should be parametric_studies/{study_name}")
	})

	t.Run("PARA-01: default concurrent value is 3", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "concurrent_test",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
			// Concurrent intentionally omitted
		}

		response := runParametricStudy(t, tool, params)
		assert.False(t, response.IsError)
		// No direct assertion on concurrent value since it doesn't appear in output,
		// but the tool should function without error
	})

	t.Run("PARA-01: single value per variable creates 1 case", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "single_case",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		response := runParametricStudy(t, tool, params)
		assert.Contains(t, response.Content, "총 케이스 수: 1")
	})

	t.Run("PARA-01: rejects empty study_name", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		response := runParametricStudy(t, tool, params)
		assert.True(t, response.IsError, "empty study_name should produce error")
	})

	t.Run("PARA-01: rejects empty variables list", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "empty_vars",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		response := runParametricStudy(t, tool, params)
		assert.True(t, response.IsError, "empty variables should produce error")
	})
}

// --- Large Cartesian Product ---

func TestParametricStudy_Integration_LargeStudy(t *testing.T) {
	t.Run("PARA-01: handles large Cartesian product (4 vars, 2 values each = 16 cases)", func(t *testing.T) {
		tempDir := t.TempDir()
		baseCasePath := createTestBaseCase(t, tempDir)

		tool := NewParametricStudyTool()
		params := ParametricStudyParams{
			StudyName: "large_study",
			BaseCase:  baseCasePath,
			Variables: []ParametricVariable{
				{Name: "fill_ratio", Values: []float64{0.5, 0.7}},
				{Name: "frequency", Values: []float64{0.5, 1.0}},
				{Name: "amplitude", Values: []float64{0.02, 0.05}},
				{Name: "viscosity", Values: []float64{1e-6, 1e-5}},
			},
			OutputDir: filepath.Join(tempDir, "output"),
		}

		response := runParametricStudy(t, tool, params)
		assert.Contains(t, response.Content, "총 케이스 수: 16")
		assert.Contains(t, response.Content, "총 16개 케이스 생성 완료")
	})
}

// --- Combination Generator (unit-level supplement) ---

func TestGenerateCombinations_Integration(t *testing.T) {
	t.Run("PARA-01: 3-variable combinations are exhaustive", func(t *testing.T) {
		variables := []ParametricVariable{
			{Name: "a", Values: []float64{1, 2}},
			{Name: "b", Values: []float64{10, 20}},
			{Name: "c", Values: []float64{100, 200}},
		}

		combos := generateCombinations(variables)
		assert.Len(t, combos, 8, "2*2*2 should produce 8 combinations")

		// Verify all unique
		seen := make(map[string]bool)
		for _, combo := range combos {
			key := fmt.Sprintf("%.0f_%.0f_%.0f", combo["a"], combo["b"], combo["c"])
			assert.False(t, seen[key], "duplicate combination found: %s", key)
			seen[key] = true
		}

		// Verify specific combinations
		expected := []string{
			"1_10_100", "1_10_200", "1_20_100", "1_20_200",
			"2_10_100", "2_10_200", "2_20_100", "2_20_200",
		}
		for _, exp := range expected {
			assert.True(t, seen[exp], "expected combination %s not found", exp)
		}
	})

	t.Run("PARA-01: single-variable preserves order", func(t *testing.T) {
		variables := []ParametricVariable{
			{Name: "x", Values: []float64{3.0, 1.0, 2.0}},
		}

		combos := generateCombinations(variables)
		require.Len(t, combos, 3)

		assert.Equal(t, 3.0, combos[0]["x"])
		assert.Equal(t, 1.0, combos[1]["x"])
		assert.Equal(t, 2.0, combos[2]["x"])
	})
}

// ============================================================================
// Test Helpers
// ============================================================================

func createTestBaseCase(t *testing.T, dir string) string {
	t.Helper()
	baseCasePath := filepath.Join(dir, "base_case.xml")
	baseContent := `<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="1000" />
            <cflnumber value="0.2" />
        </constantsdef>
        <geometry>
            <definition dp="0.01" />
        </geometry>
    </casedef>
    <parameters>
        <fill_ratio value="{{fill_ratio}}" />
        <frequency value="{{frequency}}" />
    </parameters>
</case>`
	require.NoError(t, os.WriteFile(baseCasePath, []byte(baseContent), 0644))
	return baseCasePath
}

func runParametricStudy(t *testing.T, tool BaseTool, params ParametricStudyParams) ToolResponse {
	t.Helper()
	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  ParametricStudyToolName,
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err, "tool.Run should not return Go error")
	return response
}
