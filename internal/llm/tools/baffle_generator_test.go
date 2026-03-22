package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaffleGenerator_Info(t *testing.T) {
	tool := NewBaffleGeneratorTool()
	info := tool.Info()

	assert.Equal(t, "baffle_generator", info.Name)
	assert.NotEmpty(t, info.Description)
	assert.Contains(t, info.Description, "baffle 최적화")
	assert.Contains(t, info.Parameters, "baffles")
	assert.Contains(t, info.Parameters, "tank_bounds")
	assert.Contains(t, info.Parameters, "dp")
	assert.Contains(t, info.Parameters, "xml_file")
	assert.Contains(t, info.Required, "baffles")
	assert.Contains(t, info.Required, "tank_bounds")
	assert.Contains(t, info.Required, "dp")
}

func TestBaffleGenerator_SingleVerticalBaffle(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	params := BaffleGeneratorParams{
		Baffles: []Baffle{
			{
				BaffleType: "vertical",
				PositionX:  0.25,
				Height:     0.2,
				Thickness:  0.008,
				MK:         10,
			},
		},
		TankBounds: [6]float64{0, 0, 0, 0.5, 0.35, 0.4},
		DP:         0.004,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := tool.Run(context.Background(), ToolCall{
		Name:  BaffleGeneratorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)

	// Verify XML snippet content
	assert.Contains(t, resp.Content, "vertical baffle")
	assert.Contains(t, resp.Content, `mk="10"`)
	assert.Contains(t, resp.Content, "<drawbox>")
	assert.Contains(t, resp.Content, "<boxfill>solid</boxfill>")
	assert.Contains(t, resp.Content, "x=0.25")

	// Verify point.x = position_x - thickness/2 = 0.25 - 0.004 = 0.246
	assert.Contains(t, resp.Content, `x="0.2460"`)
	// Verify size.x = thickness = 0.008
	assert.Contains(t, resp.Content, `x="0.0080"`)
}

func TestBaffleGenerator_MultipleBaffles(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	params := BaffleGeneratorParams{
		Baffles: []Baffle{
			{
				BaffleType: "vertical",
				PositionX:  0.15,
				Height:     0.2,
			},
			{
				BaffleType: "vertical",
				PositionX:  0.35,
				Height:     0.15,
				MK:         20,
			},
		},
		TankBounds: [6]float64{0, 0, 0, 0.5, 0.35, 0.4},
		DP:         0.004,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := tool.Run(context.Background(), ToolCall{
		Name:  BaffleGeneratorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)

	// First baffle: default mk = 10+0 = 10
	assert.Contains(t, resp.Content, `mk="10"`)
	// Second baffle: explicit mk = 20
	assert.Contains(t, resp.Content, `mk="20"`)
	assert.Contains(t, resp.Content, "2개")
}

func TestBaffleGenerator_InsertIntoXML(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	// Create a temporary XML file with <shapeout> tag
	tmpDir, err := os.MkdirTemp("", "baffle_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	xmlPath := filepath.Join(tmpDir, "test_case_Def.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <geometry>
            <commands>
                <mainlist>
                    <setmkbound mk="0" />
                    <drawbox>
                        <boxfill>bottom | left | right | front | back</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="0.5" y="0.35" z="0.4" />
                    </drawbox>
                    <shapeout file="Tank" />
                </mainlist>
            </commands>
        </geometry>
    </casedef>
</case>`
	require.NoError(t, os.WriteFile(xmlPath, []byte(xmlContent), 0644))

	params := BaffleGeneratorParams{
		XMLFile: xmlPath,
		Baffles: []Baffle{
			{
				BaffleType: "vertical",
				PositionX:  0.25,
				Height:     0.2,
				Thickness:  0.008,
				MK:         10,
			},
		},
		TankBounds: [6]float64{0, 0, 0, 0.5, 0.35, 0.4},
		DP:         0.004,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := tool.Run(context.Background(), ToolCall{
		Name:  BaffleGeneratorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "삽입되었습니다")

	// Read modified XML and verify structure
	modified, err := os.ReadFile(xmlPath)
	require.NoError(t, err)
	modifiedStr := string(modified)

	// Baffle snippet should appear before <shapeout>
	assert.Contains(t, modifiedStr, `setmkbound mk="10"`)
	assert.Contains(t, modifiedStr, "<shapeout")

	// Verify order: baffle before shapeout
	baffleIdx := len(modifiedStr) - len(modifiedStr) // just for reference
	_ = baffleIdx
	assert.Less(t,
		indexOf(modifiedStr, `mk="10"`),
		indexOf(modifiedStr, "<shapeout"),
	)
}

func TestBaffleGenerator_ValidationErrors(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	tests := []struct {
		name   string
		params BaffleGeneratorParams
		errMsg string
	}{
		{
			name: "empty baffles",
			params: BaffleGeneratorParams{
				Baffles:    []Baffle{},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "비어 있습니다",
		},
		{
			name: "invalid baffle_type",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "diagonal", PositionX: 0.5, Height: 0.3},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "vertical",
		},
		{
			name: "baffle outside tank bounds (x)",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "vertical", PositionX: 1.5, Height: 0.3},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "밖에 있습니다",
		},
		{
			name: "baffle height exceeds tank",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "vertical", PositionX: 0.5, Height: 2.0},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "높이",
		},
		{
			name: "zero height",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "vertical", PositionX: 0.5, Height: 0},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "0보다 커야",
		},
		{
			name: "zero dp",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "vertical", PositionX: 0.5, Height: 0.3},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0,
			},
			errMsg: "dp는 0보다",
		},
		{
			name: "invalid tank_bounds (max <= min)",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "vertical", PositionX: 0.5, Height: 0.3},
				},
				TankBounds: [6]float64{1, 1, 1, 0, 0, 0},
				DP:         0.01,
			},
			errMsg: "tank_bounds",
		},
		{
			name: "horizontal baffle height outside z bounds",
			params: BaffleGeneratorParams{
				Baffles: []Baffle{
					{BaffleType: "horizontal", Height: 1.5},
				},
				TankBounds: [6]float64{0, 0, 0, 1, 1, 1},
				DP:         0.01,
			},
			errMsg: "밖에 있습니다",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paramsJSON, err := json.Marshal(tt.params)
			require.NoError(t, err)

			resp, err := tool.Run(context.Background(), ToolCall{
				Name:  BaffleGeneratorToolName,
				Input: string(paramsJSON),
			})
			require.NoError(t, err)
			assert.True(t, resp.IsError, "expected error response")
			assert.Contains(t, resp.Content, tt.errMsg)
		})
	}
}

func TestBaffleGenerator_DefaultThickness(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	dp := 0.004
	params := BaffleGeneratorParams{
		Baffles: []Baffle{
			{
				BaffleType: "vertical",
				PositionX:  0.25,
				Height:     0.2,
				Thickness:  0, // should default to dp*2 = 0.008
			},
		},
		TankBounds: [6]float64{0, 0, 0, 0.5, 0.35, 0.4},
		DP:         dp,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := tool.Run(context.Background(), ToolCall{
		Name:  BaffleGeneratorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)

	// Default thickness = dp*2 = 0.008
	// point.x = 0.25 - 0.008/2 = 0.246
	assert.Contains(t, resp.Content, `x="0.2460"`)
	// size.x = 0.008
	assert.Contains(t, resp.Content, `x="0.0080"`)
}

func TestBaffleGenerator_HorizontalBaffle(t *testing.T) {
	tool := NewBaffleGeneratorTool()

	params := BaffleGeneratorParams{
		Baffles: []Baffle{
			{
				BaffleType: "horizontal",
				Height:     0.2, // installation height (z position)
				Thickness:  0.008,
				MK:         15,
			},
		},
		TankBounds: [6]float64{0, 0, 0, 0.5, 0.35, 0.4},
		DP:         0.004,
	}

	paramsJSON, err := json.Marshal(params)
	require.NoError(t, err)

	resp, err := tool.Run(context.Background(), ToolCall{
		Name:  BaffleGeneratorToolName,
		Input: string(paramsJSON),
	})
	require.NoError(t, err)
	assert.False(t, resp.IsError)
	assert.Contains(t, resp.Content, "horizontal baffle")
	assert.Contains(t, resp.Content, `mk="15"`)
	// point.z = height - thickness/2 = 0.2 - 0.004 = 0.196
	assert.Contains(t, resp.Content, `z="0.1960"`)
}

// indexOf returns the index of the first occurrence of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
