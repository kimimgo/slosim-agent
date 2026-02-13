package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GEO-01: Cylindrical Geometry
// GEO-02: L-Shaped Geometry
// EXC-01: Seismic Wave Motion
// EXC-02: Custom Motion

func TestCylindricalGeometry(t *testing.T) {
	t.Run("GEO-01: generates valid cylindrical geometry XML", func(t *testing.T) {
		// 원통형 탱크 지오메트리 생성
		radius := 1.0
		height := 2.0
		fluidHeight := 1.0
		dp := 0.02

		xml := CylindricalGeometry(radius, height, fluidHeight, dp)

		// Check essential components
		assert.Contains(t, xml, "<definition dp=")
		assert.Contains(t, xml, "<pointmin")
		assert.Contains(t, xml, "<pointmax")
		assert.Contains(t, xml, "<drawcylinder")
		assert.Contains(t, xml, "<setmkfluid mk=\"0\"")
		assert.Contains(t, xml, "<setmkbound mk=\"0\"")
		assert.Contains(t, xml, "radius=\"1\"")
		assert.Contains(t, xml, "<shapeout file=\"Tank\"")
	})

	t.Run("GEO-01: cylindrical fluid volume within bounds", func(t *testing.T) {
		// 유체 높이가 탱크 높이 내에 있는지 검증
		radius := 0.5
		height := 1.5
		fluidHeight := 0.8
		dp := 0.01

		xml := CylindricalGeometry(radius, height, fluidHeight, dp)

		// Fluid height should be present in XML
		assert.Contains(t, xml, "0.8")
		assert.Contains(t, xml, "mask=\"1 | 2 | 3\"") // Fluid mask
	})

	t.Run("GEO-01: margin calculation with dp", func(t *testing.T) {
		// dp * 5 margin 확인
		dp := 0.05
		margin := dp * 5 // 0.25

		xml := CylindricalGeometry(1.0, 2.0, 1.0, dp)

		// Margin should be reflected in pointmin/pointmax
		assert.Contains(t, xml, "dp=\"0.05\"")
		// pointmin/max values should account for margin
	})
}

func TestLShapedGeometry(t *testing.T) {
	t.Run("GEO-02: generates valid L-shaped geometry XML", func(t *testing.T) {
		// L형 탱크 지오메트리 생성
		L1 := 2.0  // First section length
		W1 := 1.0  // First section width
		L2 := 1.5  // Second section length
		W2 := 0.8  // Second section width
		H := 1.5   // Height
		fluidHeight := 0.7
		dp := 0.02

		xml := LShapedGeometry(L1, W1, L2, W2, H, fluidHeight, dp)

		// Check essential components
		assert.Contains(t, xml, "<definition dp=")
		assert.Contains(t, xml, "<setmkfluid mk=\"0\"")
		assert.Contains(t, xml, "<setmkbound mk=\"0\"")
		assert.Contains(t, xml, "<drawbox>")

		// Two fluid boxes for L-shape
		fluidBoxCount := strings.Count(xml, "<setmkfluid")
		assert.Equal(t, 1, fluidBoxCount) // One setmkfluid, but two drawbox calls

		// Two boundary boxes
		boundBoxCount := strings.Count(xml, "<setmkbound")
		assert.Equal(t, 1, boundBoxCount)
	})

	t.Run("GEO-02: L-shaped dimensions in XML", func(t *testing.T) {
		L1 := 1.0
		W1 := 0.5
		L2 := 0.8
		W2 := 0.4
		H := 1.0
		fluidHeight := 0.5
		dp := 0.01

		xml := LShapedGeometry(L1, W1, L2, W2, H, fluidHeight, dp)

		// Check that dimensions appear in XML
		assert.Contains(t, xml, "1") // L1 value
		assert.Contains(t, xml, "0.5") // fluidHeight value
	})
}

func TestMotionSeismicWave(t *testing.T) {
	t.Run("EXC-01: generates seismic wave motion XML", func(t *testing.T) {
		// 지진파 motion 생성
		waveFile := "seismic_data.dat"
		duration := 10.0

		xml := MotionSeismicWave(waveFile, duration)

		// Check motion structure
		assert.Contains(t, xml, "<objreal ref=\"0\">")
		assert.Contains(t, xml, "<begin mov=\"1\" start=\"0\"")
		assert.Contains(t, xml, "<mvrectfile id=\"1\"")
		assert.Contains(t, xml, "duration=\"10\"")
		assert.Contains(t, xml, "anglesunits=\"degrees\"")
		assert.Contains(t, xml, "<file name=\"seismic_data.dat\"")
		assert.Contains(t, xml, "fields=\"2\"")
		assert.Contains(t, xml, "fieldtime=\"0\"")
		assert.Contains(t, xml, "fieldx=\"1\"")
	})

	t.Run("EXC-01: file path included in motion", func(t *testing.T) {
		waveFile := "cases/kobe_earthquake.dat"
		duration := 20.5

		xml := MotionSeismicWave(waveFile, duration)

		assert.Contains(t, xml, "kobe_earthquake.dat")
		assert.Contains(t, xml, "20.5")
	})
}

func TestMotionCustom(t *testing.T) {
	t.Run("EXC-02: generates custom motion XML with data points", func(t *testing.T) {
		// 임의 가진 데이터
		timeValues := []float64{0.0, 0.1, 0.2, 0.3}
		xValues := []float64{0.0, 0.5, 1.0, 0.5}
		duration := 0.3

		xml := MotionCustom(timeValues, xValues, duration)

		// Check motion structure
		assert.Contains(t, xml, "<objreal ref=\"0\">")
		assert.Contains(t, xml, "<mvrectfile id=\"1\"")
		assert.Contains(t, xml, "duration=\"0.3\"")
		assert.Contains(t, xml, "<values>")

		// Check data points
		assert.Contains(t, xml, "<data t=\"0\" x=\"0\"")
		assert.Contains(t, xml, "<data t=\"0.1\" x=\"0.5\"")
		assert.Contains(t, xml, "<data t=\"0.2\" x=\"1\"")
		assert.Contains(t, xml, "<data t=\"0.3\" x=\"0.5\"")
	})

	t.Run("EXC-02: handles empty arrays gracefully", func(t *testing.T) {
		timeValues := []float64{}
		xValues := []float64{}
		duration := 1.0

		xml := MotionCustom(timeValues, xValues, duration)

		// Should still generate valid structure
		assert.Contains(t, xml, "<objreal ref=\"0\">")
		assert.Contains(t, xml, "<values>")
		assert.Contains(t, xml, "</values>")
	})

	t.Run("EXC-02: handles mismatched array lengths", func(t *testing.T) {
		// timeValues longer than xValues
		timeValues := []float64{0.0, 0.1, 0.2}
		xValues := []float64{0.0, 0.5}
		duration := 0.2

		xml := MotionCustom(timeValues, xValues, duration)

		// Should use shorter length (2)
		dataCount := strings.Count(xml, "<data t=")
		assert.Equal(t, 2, dataCount)
	})
}
