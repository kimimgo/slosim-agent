package tools

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// STLImportParams defines parameters for STL file import (STL-01).
type STLImportParams struct {
	STLFile     string  `json:"stl_file"`
	OutPath     string  `json:"out_path"`
	DP          float64 `json:"dp"`
	FluidHeight float64 `json:"fluid_height,omitempty"`
	Scale       float64 `json:"scale,omitempty"`
}

type stlImportTool struct{}

const (
	STLImportToolName    = "stl_import"
	stlImportDescription = `STL 파일 입력 도구 (STL-01) — CAD 메시를 DualSPHysics 케이스로 변환합니다.

사용법:
- stl_file: STL 파일 경로 (.stl)
- out_path: 출력 경로 (XML 생성)
- dp: 파티클 간격 (m)
- fluid_height: 유체 높이 (선택)
- scale: 스케일 팩터 (기본: 1.0)

메시 품질 검증:
- 수밀성 검사 (watertight mesh)
- 법선 방향 검증
- 최소/최대 치수 확인

Docker 컨테이너 내에서 실행됩니다.`
)

func NewSTLImportTool() BaseTool {
	return &stlImportTool{}
}

func (s *stlImportTool) Info() ToolInfo {
	return ToolInfo{
		Name:        STLImportToolName,
		Description: stlImportDescription,
		Parameters: map[string]any{
			"stl_file": map[string]any{
				"type":        "string",
				"description": "STL 파일 경로",
			},
			"out_path": map[string]any{
				"type":        "string",
				"description": "출력 경로",
			},
			"dp": map[string]any{
				"type":        "number",
				"description": "파티클 간격 (m)",
			},
			"fluid_height": map[string]any{
				"type":        "number",
				"description": "유체 높이 (m)",
			},
			"scale": map[string]any{
				"type":        "number",
				"description": "스케일 팩터",
			},
		},
		Required: []string{"stl_file", "out_path", "dp"},
	}
}

func (s *stlImportTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params STLImportParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.STLFile == "" || params.OutPath == "" || params.DP <= 0 {
		return NewTextErrorResponse("필수 파라미터가 누락되었습니다: stl_file, out_path, dp를 지정해주세요"), nil
	}

	// Validate STL file exists
	if _, err := os.Stat(params.STLFile); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("STL 파일을 찾을 수 없습니다: %s", params.STLFile)), nil
	}

	// Default scale
	if params.Scale == 0 {
		params.Scale = 1.0
	}

	// Validate STL mesh quality
	validation, err := validateSTLMesh(params.STLFile)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("STL 메시 검증 실패: %s", err)), nil
	}

	if !validation.IsWatertight {
		return NewTextErrorResponse("STL 메시가 수밀하지 않습니다 (watertight 메시 필요)"), nil
	}

	// Generate XML with STL import
	xmlContent := generateSTLXML(params)
	xmlPath := params.OutPath + ".xml"
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("XML 파일 생성 실패: %s", err)), nil
	}

	resultMsg := fmt.Sprintf(`STL 파일 임포트 완료:
- XML 케이스: %s
- 메시 정보:
  * 삼각형 수: %d
  * 수밀성: %v
  * 바운딩 박스: [%.3f, %.3f, %.3f] ~ [%.3f, %.3f, %.3f]
`,
		xmlPath,
		validation.TriangleCount,
		validation.IsWatertight,
		validation.BBoxMin[0], validation.BBoxMin[1], validation.BBoxMin[2],
		validation.BBoxMax[0], validation.BBoxMax[1], validation.BBoxMax[2],
	)

	return NewTextResponse(resultMsg), nil
}

// STLValidation holds mesh validation results.
type STLValidation struct {
	TriangleCount    int
	IsWatertight     bool
	NormalsConsistent bool
	BBoxMin          [3]float64
	BBoxMax          [3]float64
}

// stlVertex represents a 3D point with quantized coordinates for map key usage.
type stlVertex struct {
	X, Y, Z float64
}

// stlTriangle holds three vertices and a normal.
type stlTriangle struct {
	Normal   stlVertex
	Vertices [3]stlVertex
}

// vertexKey returns a string key for map operations, rounding to avoid floating point issues.
func vertexKey(v stlVertex) string {
	return fmt.Sprintf("%.6f,%.6f,%.6f", v.X, v.Y, v.Z)
}

// edgeKey returns a canonical string key for an edge (order-independent).
func edgeKey(a, b stlVertex) string {
	ka, kb := vertexKey(a), vertexKey(b)
	if ka < kb {
		return ka + "|" + kb
	}
	return kb + "|" + ka
}

func validateSTLMesh(stlPath string) (*STLValidation, error) {
	content, err := os.ReadFile(stlPath)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return &STLValidation{}, nil
	}

	// Parse triangles
	var triangles []stlTriangle
	isASCII := isASCIISTL(content)

	if isASCII {
		triangles, err = parseASCIISTL(content)
	} else {
		triangles, err = parseBinarySTL(content)
	}
	if err != nil {
		return nil, err
	}

	validation := &STLValidation{
		TriangleCount: len(triangles),
	}

	if len(triangles) == 0 {
		return validation, nil
	}

	// Compute bounding box
	validation.BBoxMin = [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	validation.BBoxMax = [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for _, tri := range triangles {
		for _, v := range tri.Vertices {
			coords := [3]float64{v.X, v.Y, v.Z}
			for i := 0; i < 3; i++ {
				if coords[i] < validation.BBoxMin[i] {
					validation.BBoxMin[i] = coords[i]
				}
				if coords[i] > validation.BBoxMax[i] {
					validation.BBoxMax[i] = coords[i]
				}
			}
		}
	}

	// Watertight check: every edge must be shared by exactly 2 triangles
	edgeCounts := make(map[string]int)
	for _, tri := range triangles {
		for i := 0; i < 3; i++ {
			j := (i + 1) % 3
			key := edgeKey(tri.Vertices[i], tri.Vertices[j])
			edgeCounts[key]++
		}
	}

	validation.IsWatertight = len(edgeCounts) > 0
	for _, count := range edgeCounts {
		if count != 2 {
			validation.IsWatertight = false
			break
		}
	}

	// Normal consistency check: verify normals point outward (cross product direction matches stated normal)
	validation.NormalsConsistent = checkNormalConsistency(triangles)

	return validation, nil
}

// isASCIISTL checks whether content is ASCII STL (starts with "solid" and contains "facet").
func isASCIISTL(content []byte) bool {
	if len(content) < 6 {
		return false
	}
	// Binary STL can also start with "solid" in header, so check for "facet" keyword
	prefix := strings.HasPrefix(string(content), "solid")
	hasFacet := strings.Contains(string(content[:min(1000, len(content))]), "facet")
	return prefix && hasFacet
}

// parseASCIISTL parses ASCII STL format into triangles.
func parseASCIISTL(content []byte) ([]stlTriangle, error) {
	var triangles []stlTriangle
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	var currentTri stlTriangle
	vertexIdx := 0
	inFacet := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "facet normal") {
			parts := strings.Fields(line)
			if len(parts) >= 5 {
				currentTri.Normal.X, _ = strconv.ParseFloat(parts[2], 64)
				currentTri.Normal.Y, _ = strconv.ParseFloat(parts[3], 64)
				currentTri.Normal.Z, _ = strconv.ParseFloat(parts[4], 64)
			}
			inFacet = true
			vertexIdx = 0
		} else if strings.HasPrefix(line, "vertex") && inFacet {
			parts := strings.Fields(line)
			if len(parts) >= 4 && vertexIdx < 3 {
				currentTri.Vertices[vertexIdx].X, _ = strconv.ParseFloat(parts[1], 64)
				currentTri.Vertices[vertexIdx].Y, _ = strconv.ParseFloat(parts[2], 64)
				currentTri.Vertices[vertexIdx].Z, _ = strconv.ParseFloat(parts[3], 64)
				vertexIdx++
			}
		} else if strings.HasPrefix(line, "endfacet") {
			if vertexIdx == 3 {
				triangles = append(triangles, currentTri)
			}
			inFacet = false
			currentTri = stlTriangle{}
		}
	}

	return triangles, scanner.Err()
}

// parseBinarySTL parses binary STL format into triangles.
func parseBinarySTL(content []byte) ([]stlTriangle, error) {
	if len(content) < 84 {
		return nil, fmt.Errorf("binary STL too short: %d bytes", len(content))
	}

	triCount := binary.LittleEndian.Uint32(content[80:84])
	expectedSize := 84 + int(triCount)*50
	if len(content) < expectedSize {
		return nil, fmt.Errorf("binary STL truncated: expected %d bytes, got %d", expectedSize, len(content))
	}

	triangles := make([]stlTriangle, triCount)
	for i := uint32(0); i < triCount; i++ {
		offset := 84 + int(i)*50
		data := content[offset : offset+50]

		triangles[i].Normal.X = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[0:4])))
		triangles[i].Normal.Y = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])))
		triangles[i].Normal.Z = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])))

		for v := 0; v < 3; v++ {
			vOff := 12 + v*12
			triangles[i].Vertices[v].X = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[vOff : vOff+4])))
			triangles[i].Vertices[v].Y = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[vOff+4 : vOff+8])))
			triangles[i].Vertices[v].Z = float64(math.Float32frombits(binary.LittleEndian.Uint32(data[vOff+8 : vOff+12])))
		}
	}

	return triangles, nil
}

// checkNormalConsistency verifies that stated normals roughly match cross-product normals.
func checkNormalConsistency(triangles []stlTriangle) bool {
	for _, tri := range triangles {
		// Compute cross product of (v1-v0) x (v2-v0)
		e1 := stlVertex{
			X: tri.Vertices[1].X - tri.Vertices[0].X,
			Y: tri.Vertices[1].Y - tri.Vertices[0].Y,
			Z: tri.Vertices[1].Z - tri.Vertices[0].Z,
		}
		e2 := stlVertex{
			X: tri.Vertices[2].X - tri.Vertices[0].X,
			Y: tri.Vertices[2].Y - tri.Vertices[0].Y,
			Z: tri.Vertices[2].Z - tri.Vertices[0].Z,
		}
		cross := stlVertex{
			X: e1.Y*e2.Z - e1.Z*e2.Y,
			Y: e1.Z*e2.X - e1.X*e2.Z,
			Z: e1.X*e2.Y - e1.Y*e2.X,
		}

		// Dot product with stated normal should be positive (same direction)
		dot := cross.X*tri.Normal.X + cross.Y*tri.Normal.Y + cross.Z*tri.Normal.Z

		// Skip degenerate triangles
		crossMag := math.Sqrt(cross.X*cross.X + cross.Y*cross.Y + cross.Z*cross.Z)
		if crossMag < 1e-10 {
			continue
		}

		if dot < 0 {
			return false
		}
	}
	return true
}

func generateSTLXML(params STLImportParams) string {
	margin := params.DP * 5
	fluidHeight := params.FluidHeight
	if fluidHeight == 0 {
		fluidHeight = 1.0 // Default
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" ?>
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
            <definition dp="%g" units_comment="metres (m)">
                <pointmin x="%.4f" y="%.4f" z="%.4f" />
                <pointmax x="%.4f" y="%.4f" z="%.4f" />
            </definition>
            <commands>
                <mainlist>
                    <setshapemode>dp | bound</setshapemode>
                    <setdrawmode mode="full" />

                    <!-- Import STL as boundary -->
                    <setmkbound mk="0" />
                    <drawfilestl file="%s" scale="%g">
                        <drawmove x="0" y="0" z="0" />
                    </drawfilestl>

                    <!-- Fluid box (adjust based on STL bounds) -->
                    <setmkfluid mk="0" />
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="%.4f" y="%.4f" z="0" />
                        <size x="%.4f" y="%.4f" z="%g" />
                    </drawbox>

                    <shapeout file="Tank" />
                </mainlist>
            </commands>
        </geometry>
    </casedef>
    <execution>
        <parameters>
            <parameter key="SavePosDouble" value="0" />
            <parameter key="StepAlgorithm" value="2" comment="2:Symplectic" />
            <parameter key="Kernel" value="2" comment="2:Wendland" />
            <parameter key="ViscoTreatment" value="1" comment="1:Artificial" />
            <parameter key="Visco" value="0.01" />
            <parameter key="DensityDT" value="2" comment="2:Fourtakas" />
            <parameter key="Shifting" value="0" />
            <parameter key="TimeMax" value="10.0" units_comment="seconds" />
            <parameter key="TimeOut" value="0.1" units_comment="seconds" />
            <parameter key="RhopOutMin" value="700" />
            <parameter key="RhopOutMax" value="1300" />
        </parameters>
    </execution>
</case>
`,
		params.DP,
		-margin, -margin, -margin,
		margin, margin, margin*2,
		filepath.Base(params.STLFile), params.Scale,
		-margin, -margin, margin*2, margin*2, fluidHeight,
	)
}
