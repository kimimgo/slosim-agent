package tools

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GEO-01, GEO-02, EXC-01: Geometry + Seismic E2E Tests

func TestCylindricalGeometry_FullXML(t *testing.T) {
	// Test complete XML case assembly with cylindrical geometry
	radius := 0.5
	height := 1.0
	fluidHeight := 0.4
	dp := 0.01

	// Generate cylindrical geometry XML fragment
	geomXML := CylindricalGeometry(radius, height, fluidHeight, dp)

	// Wrap in complete XML case structure
	fullXML := `<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="1000" />
            <hswl value="0" auto="true" />
            <gamma value="7" />
            <speedsystem value="0" auto="true" />
            <coefsound value="20" />
            <speedsound value="0" auto="true" />
            <coefh value="1.2" />
            <cflnumber value="0.2" />
        </constantsdef>
        <mkconfig boundcount="240" fluidcount="9" />
        <geometry>
` + geomXML + `
        </geometry>
    </casedef>
    <execution>
        <parameters>
            <parameter key="SavePosDouble" value="0" />
            <parameter key="StepAlgorithm" value="2" comment="2:Symplectic" />
            <parameter key="Kernel" value="2" comment="2:Wendland" />
            <parameter key="ViscoTreatment" value="1" comment="1:Artificial" />
            <parameter key="Visco" value="0.01" />
            <parameter key="TimeMax" value="10.0" units_comment="seconds" />
            <parameter key="TimeOut" value="0.1" units_comment="seconds" />
        </parameters>
    </execution>
</case>
`

	// Validate XML structure by parsing
	var xmlStruct struct {
		XMLName xml.Name `xml:"case"`
	}
	err := xml.Unmarshal([]byte(fullXML), &xmlStruct)
	require.NoError(t, err, "Generated XML should be valid")

	// Verify dp attribute is present
	assert.Contains(t, geomXML, `dp="0.01"`)

	// Verify cylinder dimensions in drawcylinder elements
	assert.Contains(t, geomXML, `radius="0.49"`) // radius - dp for fluid
	assert.Contains(t, geomXML, `radius="0.5"`)  // full radius for boundary

	// Verify fluid height appears in drawcylinder
	assert.Contains(t, geomXML, `z="0.4"`) // fluidHeight

	// Verify boundary height
	assert.Contains(t, geomXML, `z="1"`) // height

	// Verify setmkfluid and setmkbound markers
	assert.Contains(t, geomXML, `<setmkfluid mk="0" />`)
	assert.Contains(t, geomXML, `<setmkbound mk="0" />`)

	// Verify shapeout present
	assert.Contains(t, geomXML, `<shapeout file="Tank" />`)
}

func TestCylindricalGeometry_FluidHeight(t *testing.T) {
	// Verify fluid height appears correctly in drawcylinder element
	radius := 0.3
	height := 0.8
	fluidHeight := 0.25
	dp := 0.005

	geomXML := CylindricalGeometry(radius, height, fluidHeight, dp)

	// Check fluid height in the fluid cylinder definition
	assert.Contains(t, geomXML, `z="0.25"`)

	// Check that fluid cylinder uses reduced radius (radius - dp)
	_ = radius - dp // expectedFluidRadius = 0.295
	assert.Contains(t, geomXML, `radius="0.295"`)

	// Verify fluid height is less than total height
	assert.True(t, fluidHeight < height, "Fluid height should be less than tank height")
}

func TestLShapedGeometry_FullXML(t *testing.T) {
	// Test complete XML case assembly with L-shaped geometry
	L1 := 1.0
	W1 := 0.5
	L2 := 0.6
	W2 := 0.4
	H := 0.8
	fluidHeight := 0.3
	dp := 0.01

	// Generate L-shaped geometry XML fragment
	geomXML := LShapedGeometry(L1, W1, L2, W2, H, fluidHeight, dp)

	// Wrap in complete XML case structure
	fullXML := `<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="1000" />
        </constantsdef>
        <geometry>
` + geomXML + `
        </geometry>
    </casedef>
</case>
`

	// Validate XML structure
	var xmlStruct struct {
		XMLName xml.Name `xml:"case"`
	}
	err := xml.Unmarshal([]byte(fullXML), &xmlStruct)
	require.NoError(t, err, "Generated XML should be valid")

	// Verify two drawbox elements in fluid section
	fluidBoxCount := strings.Count(geomXML, "<setmkfluid mk=\"0\" />")
	assert.Equal(t, 1, fluidBoxCount, "Should have one setmkfluid marker")

	// Count drawbox after setmkfluid
	fluidSection := geomXML[strings.Index(geomXML, "<setmkfluid mk=\"0\" />"):]
	boundSection := geomXML[strings.Index(geomXML, "<setmkbound mk=\"0\" />"):]

	fluidDrawBoxCount := strings.Count(fluidSection[:strings.Index(fluidSection, "<setmkbound")], "<drawbox>")
	assert.Equal(t, 2, fluidDrawBoxCount, "Should have two drawbox elements for L-shaped fluid")

	// Verify boundary section also has two drawbox elements
	boundDrawBoxCount := strings.Count(boundSection, "<drawbox>")
	assert.Equal(t, 2, boundDrawBoxCount, "Should have two drawbox elements for L-shaped boundary")

	// Verify dp attribute
	assert.Contains(t, geomXML, `dp="0.01"`)

	// Verify dimensions in the XML
	assert.Contains(t, geomXML, `x="1"`)   // L1
	assert.Contains(t, geomXML, `y="0.5"`) // W1
	assert.Contains(t, geomXML, `z="0.3"`) // fluidHeight

	assert.Contains(t, geomXML, `x="0.6"`) // L2
	assert.Contains(t, geomXML, `y="0.4"`) // W2
}

func TestSeismicToMotion_E2E(t *testing.T) {
	// E2E test: CSV file → seismic_input tool → output .dat file → MotionSeismicWave XML
	tempDir := t.TempDir()

	// Step 1: Create CSV seismic data file
	csvPath := filepath.Join(tempDir, "seismic.csv")
	csvContent := `# Seismic wave data
0.0 0.0
0.1 0.5
0.2 1.0
0.3 0.5
0.4 0.0
`
	err := os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Step 2: Parse with seismic_input tool
	tool := NewSeismicInputTool()
	params := SeismicInputParams{
		FilePath:     csvPath,
		ValidateOnly: false, // Generate output file
		OutputFormat: "dsph",
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  SeismicInputToolName,
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError, "Seismic input tool should succeed")

	// Step 3: Verify output file exists
	outputPath := filepath.Join(tempDir, "seismic_dsph.dat")
	_, err = os.Stat(outputPath)
	require.NoError(t, err, "Output .dat file should exist")

	// Verify file has correct format (2 columns)
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	dataLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			fields := strings.Fields(trimmed)
			assert.Len(t, fields, 2, "Each data line should have 2 columns (time, accel)")
			dataLines++
		}
	}
	assert.Equal(t, 5, dataLines, "Should have 5 data lines")

	// Step 4: Generate MotionSeismicWave XML using the output file
	duration := 0.4
	motionXML := MotionSeismicWave("seismic_dsph.dat", duration)

	// Verify XML structure
	assert.Contains(t, motionXML, `<objreal ref="0">`)
	assert.Contains(t, motionXML, `<mvrectfile id="1"`)
	assert.Contains(t, motionXML, `duration="0.4"`)
	assert.Contains(t, motionXML, `<file name="seismic_dsph.dat"`)
	assert.Contains(t, motionXML, `fields="2"`)
	assert.Contains(t, motionXML, `fieldtime="0"`)
	assert.Contains(t, motionXML, `fieldx="1"`)
}

func TestSeismicToMotion_ScaleIntegration(t *testing.T) {
	// Verify scaled values in output file match input * scale_factor
	tempDir := t.TempDir()

	// Create input file with known values
	inputPath := filepath.Join(tempDir, "input.txt")
	inputContent := `0.0 1.0
0.1 2.0
0.2 3.0
`
	err := os.WriteFile(inputPath, []byte(inputContent), 0644)
	require.NoError(t, err)

	// Apply scale factor
	scaleFactor := 2.5
	tool := NewSeismicInputTool()
	params := SeismicInputParams{
		FilePath:     inputPath,
		ScaleFactor:  scaleFactor,
		ValidateOnly: false,
		OutputFormat: "dsph",
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	call := ToolCall{
		Name:  SeismicInputToolName,
		Input: string(paramsJSON),
	}

	response, err := tool.Run(context.Background(), call)
	require.NoError(t, err)
	assert.False(t, response.IsError)

	// Read output file and verify scaled values
	outputPath := filepath.Join(tempDir, "input_dsph.dat")
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	expectedScaled := []float64{2.5, 5.0, 7.5} // 1.0*2.5, 2.0*2.5, 3.0*2.5
	scaledIdx := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		fields := strings.Fields(trimmed)
		require.Len(t, fields, 2)

		var accel float64
		_, err := fmt.Sscanf(fields[1], "%f", &accel)
		require.NoError(t, err)

		assert.InDelta(t, expectedScaled[scaledIdx], accel, 0.001, "Scaled acceleration should match expected value")
		scaledIdx++
	}

	assert.Equal(t, 3, scaledIdx, "Should have processed 3 data lines")
}
