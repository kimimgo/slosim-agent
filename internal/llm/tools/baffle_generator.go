package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// BaffleGeneratorParams defines parameters for the baffle generator tool.
type BaffleGeneratorParams struct {
	XMLFile    string     `json:"xml_file"`
	Baffles    []Baffle   `json:"baffles"`
	TankBounds [6]float64 `json:"tank_bounds"` // [xmin,ymin,zmin,xmax,ymax,zmax]
	DP         float64    `json:"dp"`
}

// Baffle defines a single baffle (internal wall) geometry.
type Baffle struct {
	BaffleType string  `json:"baffle_type"` // "vertical" | "horizontal"
	PositionX  float64 `json:"position_x"`
	PositionY  float64 `json:"position_y"`
	Height     float64 `json:"height"`
	Thickness  float64 `json:"thickness"`
	MK         int     `json:"mk"`
}

type baffleGeneratorTool struct{}

const (
	BaffleGeneratorToolName    = "baffle_generator"
	baffleGeneratorDescription = `Baffle 생성 도구 — 슬로싱 저감을 위한 내부 격벽(baffle)을 DualSPHysics XML에 추가합니다.

사용법:
- baffles: baffle 배열 (위치, 높이, 타입 지정)
- tank_bounds: 탱크 경계 [xmin,ymin,zmin,xmax,ymax,zmax]
- xml_file: 기존 XML 파일 경로 (선택 — 없으면 snippet만 반환)
- dp: 파티클 간격 (baffle 두께 기본값 = dp*2)

baffle 최적화 워크플로우:
1. baseline 시뮬레이션 실행 후 SWL 측정
2. baffle_generator로 격벽 추가
3. 재시뮬레이션 후 SWL 비교
4. 위치/높이 조정하여 반복`
)

func NewBaffleGeneratorTool() BaseTool {
	return &baffleGeneratorTool{}
}

func (b *baffleGeneratorTool) Info() ToolInfo {
	return ToolInfo{
		Name:        BaffleGeneratorToolName,
		Description: baffleGeneratorDescription,
		Parameters: map[string]any{
			"xml_file": map[string]any{
				"type":        "string",
				"description": "기존 XML 파일 경로 (선택 — 없으면 snippet만 반환)",
			},
			"baffles": map[string]any{
				"type":        "array",
				"description": "baffle 배열",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"baffle_type": map[string]any{
							"type":        "string",
							"description": "baffle 타입: vertical 또는 horizontal",
						},
						"position_x": map[string]any{
							"type":        "number",
							"description": "baffle 중심 X 위치 (m)",
						},
						"position_y": map[string]any{
							"type":        "number",
							"description": "baffle 중심 Y 위치 (m, 0=전체 너비)",
						},
						"height": map[string]any{
							"type":        "number",
							"description": "baffle 높이 (m)",
						},
						"thickness": map[string]any{
							"type":        "number",
							"description": "baffle 두께 (m, 기본: dp*2)",
						},
						"mk": map[string]any{
							"type":        "integer",
							"description": "boundary mk 번호 (기본: 10+인덱스)",
						},
					},
				},
			},
			"tank_bounds": map[string]any{
				"type":        "array",
				"description": "탱크 경계 [xmin,ymin,zmin,xmax,ymax,zmax] (m)",
				"items": map[string]any{
					"type": "number",
				},
			},
			"dp": map[string]any{
				"type":        "number",
				"description": "파티클 간격 (m)",
			},
		},
		Required: []string{"baffles", "tank_bounds", "dp"},
	}
}

func (b *baffleGeneratorTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params BaffleGeneratorParams
	if err := UnmarshalToolInput(call.Input, &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if len(params.Baffles) == 0 {
		return NewTextErrorResponse("baffles 배열이 비어 있습니다"), nil
	}
	if params.DP <= 0 {
		return NewTextErrorResponse("dp는 0보다 커야 합니다"), nil
	}

	tb := params.TankBounds
	if tb[3] <= tb[0] || tb[4] <= tb[1] || tb[5] <= tb[2] {
		return NewTextErrorResponse("tank_bounds가 올바르지 않습니다: max > min 이어야 합니다"), nil
	}

	var snippets []string
	for i, baffle := range params.Baffles {
		if err := validateBaffle(baffle, tb, i); err != nil {
			return NewTextErrorResponse(err.Error()), nil
		}

		// Apply defaults
		if baffle.Thickness <= 0 {
			baffle.Thickness = params.DP * 2
		}
		if baffle.MK == 0 {
			baffle.MK = 10 + i
		}

		snippet := generateBaffleSnippet(baffle, tb, params.DP)
		snippets = append(snippets, snippet)
	}

	combined := strings.Join(snippets, "\n")

	// If xml_file is provided, insert into existing XML
	if params.XMLFile != "" {
		result, err := insertBaffleIntoXML(params.XMLFile, combined)
		if err != nil {
			return NewTextErrorResponse(err.Error()), nil
		}
		return NewTextResponse(result), nil
	}

	return NewTextResponse(fmt.Sprintf("Baffle XML snippet (%d개):\n\n%s", len(params.Baffles), combined)), nil
}

func validateBaffle(baffle Baffle, tb [6]float64, index int) error {
	if baffle.BaffleType != "vertical" && baffle.BaffleType != "horizontal" {
		return fmt.Errorf("baffle[%d]: baffle_type은 'vertical' 또는 'horizontal'이어야 합니다 (현재: '%s')", index, baffle.BaffleType)
	}
	if baffle.Height <= 0 {
		return fmt.Errorf("baffle[%d]: height는 0보다 커야 합니다", index)
	}

	// Check baffle is inside tank bounds
	if baffle.BaffleType == "vertical" {
		if baffle.PositionX <= tb[0] || baffle.PositionX >= tb[3] {
			return fmt.Errorf("baffle[%d]: position_x(%.4f)가 탱크 경계(%.4f ~ %.4f) 밖에 있습니다", index, baffle.PositionX, tb[0], tb[3])
		}
		maxHeight := tb[5] - tb[2]
		if baffle.Height > maxHeight {
			return fmt.Errorf("baffle[%d]: height(%.4f)가 탱크 높이(%.4f)보다 큽니다", index, baffle.Height, maxHeight)
		}
	} else { // horizontal
		if baffle.Height <= tb[2] || baffle.Height >= tb[5] {
			return fmt.Errorf("baffle[%d]: height(%.4f, 설치 높이)가 탱크 z 경계(%.4f ~ %.4f) 밖에 있습니다", index, baffle.Height, tb[2], tb[5])
		}
	}

	return nil
}

func generateBaffleSnippet(baffle Baffle, tb [6]float64, dp float64) string {
	var pointX, pointY, pointZ float64
	var sizeX, sizeY, sizeZ float64

	switch baffle.BaffleType {
	case "vertical":
		pointX = baffle.PositionX - baffle.Thickness/2
		pointY = tb[1] + dp
		pointZ = tb[2] + dp
		sizeX = baffle.Thickness
		sizeY = (tb[4] - tb[1]) - 2*dp
		sizeZ = baffle.Height

	case "horizontal":
		pointX = tb[0] + dp
		pointY = tb[1] + dp
		pointZ = baffle.Height - baffle.Thickness/2
		sizeX = (tb[3] - tb[0]) - 2*dp
		sizeY = (tb[4] - tb[1]) - 2*dp
		sizeZ = baffle.Thickness
	}

	return fmt.Sprintf(`                    <!-- %s baffle at %s -->
                    <setmkbound mk="%d" />
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="%.4f" y="%.4f" z="%.4f" />
                        <size x="%.4f" y="%.4f" z="%.4f" />
                    </drawbox>`,
		baffle.BaffleType, bafflePositionLabel(baffle),
		baffle.MK,
		pointX, pointY, pointZ,
		sizeX, sizeY, sizeZ,
	)
}

func bafflePositionLabel(baffle Baffle) string {
	if baffle.BaffleType == "vertical" {
		return fmt.Sprintf("x=%.4f", baffle.PositionX)
	}
	return fmt.Sprintf("z=%.4f", baffle.Height)
}

func insertBaffleIntoXML(xmlFile string, snippet string) (string, error) {
	content, err := os.ReadFile(xmlFile)
	if err != nil {
		return "", fmt.Errorf("XML 파일 읽기 실패: %s", err)
	}

	xmlStr := string(content)

	// Find <shapeout> tag and insert before it
	shapeoutIdx := strings.Index(xmlStr, "<shapeout")
	if shapeoutIdx == -1 {
		return "", fmt.Errorf("XML에서 <shapeout> 태그를 찾을 수 없습니다")
	}

	// Insert snippet before <shapeout>
	modified := xmlStr[:shapeoutIdx] + snippet + "\n" + xmlStr[shapeoutIdx:]

	// Write back
	if err := os.WriteFile(xmlFile, []byte(modified), 0644); err != nil {
		return "", fmt.Errorf("XML 파일 쓰기 실패: %s", err)
	}

	return fmt.Sprintf("XML 파일에 baffle이 삽입되었습니다: %s\n\n삽입된 snippet:\n%s", xmlFile, snippet), nil
}
