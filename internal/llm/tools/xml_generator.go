package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type XMLGeneratorParams struct {
	TankLength  float64 `json:"tank_length"`
	TankWidth   float64 `json:"tank_width"`
	TankHeight  float64 `json:"tank_height"`
	FluidHeight float64 `json:"fluid_height"`
	Freq        float64 `json:"freq"`
	Amplitude   float64 `json:"amplitude"`
	DP          float64 `json:"dp"`
	TimeMax     float64 `json:"time_max"`
	OutPath     string  `json:"out_path"`
	Dimension   string  `json:"dimension,omitempty"` // "2D" or "3D" (default: "3D")
}

type xmlGeneratorTool struct{}

const (
	XMLGeneratorToolName    = "xml_generator"
	xmlGeneratorDescription = `DualSPHysics XML 케이스 자동 생성 도구 — 해석 조건을 XML 케이스 파일로 변환합니다.

입력 파라미터 (탱크 치수, 유체 높이, 주파수 등)를 받아
DualSPHysics 호환 XML 케이스 파일과 probe_points.txt를 생성합니다.

출력:
- {out_path}.xml — DualSPHysics 호환 XML 케이스 파일
- {out_path}_probe_points.txt — 자동 생성된 측정 포인트 파일`
)

func NewXMLGeneratorTool() BaseTool {
	return &xmlGeneratorTool{}
}

func (x *xmlGeneratorTool) Info() ToolInfo {
	return ToolInfo{
		Name:        XMLGeneratorToolName,
		Description: xmlGeneratorDescription,
		Parameters: map[string]any{
			"tank_length": map[string]any{
				"type":        "number",
				"description": "탱크 길이 (m)",
			},
			"tank_width": map[string]any{
				"type":        "number",
				"description": "탱크 너비 (m)",
			},
			"tank_height": map[string]any{
				"type":        "number",
				"description": "탱크 높이 (m)",
			},
			"fluid_height": map[string]any{
				"type":        "number",
				"description": "유체 높이 (m)",
			},
			"freq": map[string]any{
				"type":        "number",
				"description": "가진 주파수 (Hz)",
			},
			"amplitude": map[string]any{
				"type":        "number",
				"description": "가진 진폭 (m)",
			},
			"dp": map[string]any{
				"type":        "number",
				"description": "파티클 간격 (m)",
			},
			"time_max": map[string]any{
				"type":        "number",
				"description": "시뮬레이션 시간 (s)",
			},
			"out_path": map[string]any{
				"type":        "string",
				"description": "출력 경로 (.xml 확장자 자동 추가)",
			},
		},
		Required: []string{
			"tank_length", "tank_width", "tank_height",
			"fluid_height", "freq", "amplitude", "dp", "time_max", "out_path",
		},
	}
}

func (x *xmlGeneratorTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params XMLGeneratorParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	// Validate dimensions
	if params.TankLength <= 0 || params.TankWidth <= 0 || params.TankHeight <= 0 {
		return NewTextErrorResponse("탱크 치수는 0보다 커야 합니다"), nil
	}
	if params.FluidHeight <= 0 {
		return NewTextErrorResponse("유체 높이는 0보다 커야 합니다"), nil
	}
	if params.FluidHeight > params.TankHeight {
		return NewTextErrorResponse(fmt.Sprintf("유체 높이(%.2fm)가 탱크 높이(%.2fm)보다 클 수 없습니다", params.FluidHeight, params.TankHeight)), nil
	}
	if params.Freq <= 0 || params.Amplitude <= 0 || params.DP <= 0 || params.TimeMax <= 0 {
		return NewTextErrorResponse("주파수, 진폭, 파티클 간격, 시뮬레이션 시간은 모두 0보다 커야 합니다"), nil
	}
	if params.OutPath == "" {
		return NewTextErrorResponse("출력 경로(out_path)를 지정해주세요"), nil
	}

	// Ensure output directory exists
	outDir := filepath.Dir(params.OutPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("출력 디렉토리 생성 실패: %s", err)), nil
	}

	// Default dimension
	if params.Dimension == "" {
		params.Dimension = "3D"
	}

	// Generate XML (DIM-01: 2D support)
	xmlContent := generateSloshingXML(params)
	xmlPath := params.OutPath + ".xml"
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("XML 파일 생성 실패: %s", err)), nil
	}

	// Generate probe_points.txt
	probeContent := generateProbePoints(params)
	probePath := params.OutPath + "_probe_points.txt"
	if err := os.WriteFile(probePath, []byte(probeContent), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("측정 포인트 파일 생성 실패: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf(
		"XML 케이스 파일이 생성되었습니다.\n- 케이스: %s\n- 측정 포인트: %s",
		xmlPath, probePath)), nil
}

func generateSloshingXML(p XMLGeneratorParams) string {
	L := p.TankLength
	W := p.TankWidth
	H := p.TankHeight
	fH := p.FluidHeight
	margin := p.DP * 5
	timeOut := p.TimeMax / 50

	// DIM-01: For 2D, set width to minimal value
	if p.Dimension == "2D" {
		W = p.DP * 3 // Minimal width for 2D simulation
	}

	// SWL gauge positions
	leftX := fmt.Sprintf("%.4f", 0.05*L)
	centerX := fmt.Sprintf("%.4f", 0.50*L)
	rightX := fmt.Sprintf("%.4f", 0.95*L)
	centerY := fmt.Sprintf("%.4f", 0.50*W)

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" ?>
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="1000" />
            <rhopgradient value="2" />
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
            <definition dp="%g" units_comment="metres (m)">
                <pointmin x="%.4f" y="%.4f" z="%.4f" />
                <pointmax x="%.4f" y="%.4f" z="%.4f" />
            </definition>
            <commands>
                <mainlist>
                    <setshapemode>dp | bound</setshapemode>
                    <setdrawmode mode="full" />
                    <setmkfluid mk="0" />
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>
                    <setmkbound mk="0" />
                    <drawbox>
                        <boxfill>bottom | left | right | front | back</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>
                    <shapeout file="Tank" />
                </mainlist>
            </commands>
        </geometry>
        <motion>
            <objreal ref="0">
                <begin mov="1" start="0" />
                <mvrectsinu id="1" duration="%g" anglesunits="degrees">
                    <freq x="%g" y="0" z="0" units_comment="Hz" />
                    <ampl x="%g" y="0" z="0" units_comment="metres (m)" />
                    <phase x="0" y="0" z="0" units_comment="degrees" />
                </mvrectsinu>
            </objreal>
        </motion>
    </casedef>
    <execution>
        <special>
            <gauges>
                <swl name="SWL_Left">
                    <point0 x="%s" y="%s" z="0.0" />
                    <point2 x="%s" y="%s" z="%g" />
                    <pointdp value="0.5" />
                    <masslimit value="0.2" />
                    <savevtk value="false" />
                </swl>
                <swl name="SWL_Center">
                    <point0 x="%s" y="%s" z="0.0" />
                    <point2 x="%s" y="%s" z="%g" />
                    <pointdp value="0.5" />
                    <masslimit value="0.2" />
                    <savevtk value="false" />
                </swl>
                <swl name="SWL_Right">
                    <point0 x="%s" y="%s" z="0.0" />
                    <point2 x="%s" y="%s" z="%g" />
                    <pointdp value="0.5" />
                    <masslimit value="0.2" />
                    <savevtk value="false" />
                </swl>
            </gauges>
        </special>
        <parameters>
            <parameter key="SavePosDouble" value="0" />
            <parameter key="StepAlgorithm" value="2" comment="2:Symplectic" />
            <parameter key="Kernel" value="2" comment="2:Wendland" />
            <parameter key="ViscoTreatment" value="1" comment="1:Artificial" />
            <parameter key="Visco" value="0.01" />
            <parameter key="ViscoBoundFactor" value="1" />
            <parameter key="DensityDT" value="2" comment="2:Fourtakas" />
            <parameter key="DensityDTvalue" value="0.1" />
            <parameter key="Shifting" value="0" />
            <parameter key="RigidAlgorithm" value="1" />
            <parameter key="CoefDtMin" value="0.05" />
            <parameter key="DtIni" value="0.0001" />
            <parameter key="DtMin" value="0.0001" />
            <parameter key="TimeMax" value="%g" units_comment="seconds" />
            <parameter key="TimeOut" value="%g" units_comment="seconds" />
            <parameter key="PartsOutMax" value="1" />
            <parameter key="RhopOutMin" value="700" />
            <parameter key="RhopOutMax" value="1300" />
            <simulationdomain>
                <posmin x="default - 10%%" y="default" z="default" />
                <posmax x="default + 10%%" y="default" z="default + 50%%" />
            </simulationdomain>
        </parameters>
    </execution>
</case>
`,
		// definition dp, pointmin xyz, pointmax xyz
		p.DP, -margin, -margin, -margin, L+margin, W+margin, H+margin,
		// fluid drawbox size
		L, W, fH,
		// tank drawbox size
		L, W, H,
		// motion: duration, freq, amplitude
		p.TimeMax, p.Freq, p.Amplitude,
		// SWL_Left: point0 xy, point2 xy z
		leftX, centerY, leftX, centerY, H,
		// SWL_Center: point0 xy, point2 xy z
		centerX, centerY, centerX, centerY, H,
		// SWL_Right: point0 xy, point2 xy z
		rightX, centerY, rightX, centerY, H,
		// parameters: TimeMax, TimeOut
		p.TimeMax, timeOut,
	)
}

func generateProbePoints(p XMLGeneratorParams) string {
	// Generate probe points at 3 x-positions × 6 z-heights
	// x-positions: left wall (5%), center (50%), right wall (95%)
	// z-heights: evenly spaced from dp to tank_height
	xPositions := []float64{0.05 * p.TankLength, 0.50 * p.TankLength, 0.95 * p.TankLength}
	y := 0.5 * p.TankWidth
	numHeights := 6
	result := ""

	for _, x := range xPositions {
		for i := 0; i < numHeights; i++ {
			z := p.DP + float64(i)*(p.TankHeight-p.DP)/float64(numHeights-1)
			result += fmt.Sprintf("%.4f %.4f %.4f\n", x, y, z)
		}
	}

	return result
}
