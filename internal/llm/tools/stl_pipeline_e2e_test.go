package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// STL Pipeline E2E Tests — Quality Gate validation for S10 fuel tank scenario.
// Gate 1: XML generation correctness (unit-level, no Docker)
// Gate 2: Generated XML structural validation against Frosina2018 reference pattern

func TestSTLPipeline_Gate1_XMLGeneration(t *testing.T) {
	// Gate 1: Generate XML from fuel_tank.stl and verify all 4 bug fixes are applied
	stlPath := filepath.Join("..", "..", "..", "cases", "fuel_tank.stl")
	if _, err := os.Stat(stlPath); os.IsNotExist(err) {
		t.Skip("fuel_tank.stl not found in cases/ — skipping E2E test")
	}

	tempDir := t.TempDir()
	tool := NewSTLImportTool()
	params := STLImportParams{
		STLFile:       stlPath,
		OutPath:       filepath.Join(tempDir, "FuelTank_Def"),
		DP:            0.008,
		AutoFillPoint: true,
		FillRatio:     0.5,
		MotionType:    "mvrectsinu",
		MotionFreq:    0.5,
		MotionAmpl:    0.05,
		MotionAxis:    "x",
		TimeMax:       5.0,
		Scale:         1.0,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	response, err := tool.Run(context.Background(), ToolCall{
		Name:  STLImportToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, response.IsError, "STL import should succeed: %s", response.Content)

	// Read generated XML
	xmlPath := filepath.Join(tempDir, "FuelTank_Def.xml")
	xmlBytes, err := os.ReadFile(xmlPath)
	require.NoError(t, err, "Generated XML file should exist")
	xml := string(xmlBytes)

	// --- Quality Gate 1a: autofill="false" (no scale attribute) ---
	t.Run("Gate1a_autofill_no_scale_attr", func(t *testing.T) {
		assert.Contains(t, xml, `autofill="false"`, "drawfilestl must have autofill=false")
		assert.NotContains(t, xml, `scale="1"`, "drawfilestl must NOT have scale attribute for scale=1.0")
		assert.NotContains(t, xml, `<drawscale`, "scale=1.0 should not produce drawscale element")
	})

	// --- Quality Gate 1b: modefill void ---
	t.Run("Gate1b_modefill_void", func(t *testing.T) {
		assert.Contains(t, xml, "<modefill>void</modefill>", "fillpoint must contain modefill void")
		assert.Contains(t, xml, "<fillpoint", "fillpoint element must exist")
	})

	// --- Quality Gate 1c: fillpoint before drawbox ---
	t.Run("Gate1c_fillpoint_order", func(t *testing.T) {
		fpIdx := strings.Index(xml, "<fillpoint")
		dbIdx := strings.Index(xml, "<drawbox>")
		assert.Greater(t, fpIdx, 0, "fillpoint must exist")
		assert.Greater(t, dbIdx, 0, "drawbox must exist")
		assert.Less(t, fpIdx, dbIdx, "fillpoint must appear BEFORE drawbox")
	})

	// --- Quality Gate 1d: begin element + no next="0" ---
	t.Run("Gate1d_motion_begin", func(t *testing.T) {
		assert.Contains(t, xml, `<begin mov="1" start="0" />`, "motion must have begin element")
		assert.NotContains(t, xml, `next="0"`, "motion must NOT have next=0 attribute")
	})

	// --- Quality Gate 1e: motion parameters ---
	t.Run("Gate1e_motion_params", func(t *testing.T) {
		assert.Contains(t, xml, "<mvrectsinu", "motion type must be mvrectsinu")
		assert.Contains(t, xml, `duration="5"`, "duration must match TimeMax")
		assert.Contains(t, xml, `<freq x="0.5"`, "frequency must be 0.5 Hz")
		assert.Contains(t, xml, `<ampl x="0.05"`, "amplitude must be 0.05 m")
	})

	// --- Quality Gate 1f: auto fillpoint computed from BBox center ---
	t.Run("Gate1f_auto_fillpoint", func(t *testing.T) {
		assert.Contains(t, xml, "<fillpoint", "auto fillpoint should generate fillpoint element")
		assert.Contains(t, response.Content, "Fillpoint:", "response should report fillpoint coordinates")
	})

	// --- Quality Gate 1g: TimeMax in execution params ---
	t.Run("Gate1g_timemax", func(t *testing.T) {
		assert.Contains(t, xml, `key="TimeMax" value="5.0"`, "TimeMax execution parameter must be 5.0")
	})

	// --- Quality Gate 1h: fluid height from fill_ratio ---
	t.Run("Gate1h_fill_ratio", func(t *testing.T) {
		assert.Contains(t, xml, "<drawbox>", "drawbox for fluid must exist")
		assert.Contains(t, xml, "<boxfill>solid</boxfill>", "boxfill must be solid")
		// Fill ratio 0.5 should produce non-zero fluid height
		assert.Contains(t, response.Content, "유체 높이:", "response should report fluid height")
	})

	// Export generated XML for Gate 2 Docker test
	t.Logf("Generated XML path: %s", xmlPath)
	t.Logf("Response:\n%s", response.Content)
}

func TestSTLPipeline_Gate1_ScaleHandling(t *testing.T) {
	// Verify scale != 1.0 produces drawscale child element
	stlPath := filepath.Join("..", "..", "..", "cases", "fuel_tank.stl")
	if _, err := os.Stat(stlPath); os.IsNotExist(err) {
		t.Skip("fuel_tank.stl not found")
	}

	tempDir := t.TempDir()
	tool := NewSTLImportTool()
	params := STLImportParams{
		STLFile:       stlPath,
		OutPath:       filepath.Join(tempDir, "scaled_tank"),
		DP:            0.008,
		AutoFillPoint: true,
		Scale:         0.001, // mm → m conversion
	}

	paramsJSON, _ := json.Marshal(params)
	response, err := tool.Run(context.Background(), ToolCall{
		Name:  STLImportToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, response.IsError)

	xmlBytes, err := os.ReadFile(filepath.Join(tempDir, "scaled_tank.xml"))
	require.NoError(t, err)
	xml := string(xmlBytes)

	assert.Contains(t, xml, `<drawscale x="0.001" y="0.001" z="0.001" />`,
		"scale=0.001 must produce drawscale child element")
	assert.Contains(t, xml, `autofill="false"`)
	assert.NotContains(t, xml, `scale="0.001"`,
		"scale must NOT be a drawfilestl attribute")
}

func TestSTLPipeline_Gate1_RotationMotion(t *testing.T) {
	// Verify mvrotsinu gets anglesunits="degrees"
	stlPath := filepath.Join("..", "..", "..", "cases", "fuel_tank.stl")
	if _, err := os.Stat(stlPath); os.IsNotExist(err) {
		t.Skip("fuel_tank.stl not found")
	}

	tempDir := t.TempDir()
	tool := NewSTLImportTool()
	params := STLImportParams{
		STLFile:       stlPath,
		OutPath:       filepath.Join(tempDir, "rot_tank"),
		DP:            0.008,
		AutoFillPoint: true,
		MotionType:    "mvrotsinu",
		MotionFreq:    0.65,
		MotionAmpl:    4.0,
		MotionAxis:    "y",
		TimeMax:       10.0,
		Scale:         1.0,
	}

	paramsJSON, _ := json.Marshal(params)
	response, err := tool.Run(context.Background(), ToolCall{
		Name:  STLImportToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, response.IsError)

	xmlBytes, err := os.ReadFile(filepath.Join(tempDir, "rot_tank.xml"))
	require.NoError(t, err)
	xml := string(xmlBytes)

	assert.Contains(t, xml, `anglesunits="degrees"`, "mvrotsinu must have anglesunits=degrees")
	assert.Contains(t, xml, `<begin mov="1" start="0" />`)
	assert.Contains(t, xml, "<mvrotsinu")
}
