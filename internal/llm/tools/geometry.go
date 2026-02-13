package tools

import (
	"fmt"
	"math"
)

// GeometryType represents different tank geometries (GEO-01, GEO-02).
type GeometryType string

const (
	GeometryRectangular GeometryType = "rectangular"
	GeometryCylindrical GeometryType = "cylindrical"
	GeometryLShaped     GeometryType = "l_shaped"
)

// CylindricalGeometry generates XML geometry for cylindrical tank (GEO-01).
func CylindricalGeometry(radius, height, fluidHeight, dp float64) string {
	margin := dp * 5
	segments := 36 // TODO: use for curved wall generation
	_ = segments // Angular segments for cylinder approximation

	return fmt.Sprintf(`
            <definition dp="%g" units_comment="metres (m)">
                <pointmin x="%.4f" y="%.4f" z="%.4f" />
                <pointmax x="%.4f" y="%.4f" z="%.4f" />
            </definition>
            <commands>
                <mainlist>
                    <setshapemode>dp | bound</setshapemode>
                    <setdrawmode mode="full" />

                    <!-- Fluid: Cylindrical volume -->
                    <setmkfluid mk="0" />
                    <drawcylinder radius="%g" mask="1 | 2 | 3">
                        <point x="0" y="0" z="0" />
                        <point x="0" y="0" z="%g" />
                    </drawcylinder>

                    <!-- Boundary: Cylindrical shell + bottom -->
                    <setmkbound mk="0" />
                    <drawcylinder radius="%g" mask="2">
                        <point x="0" y="0" z="0" />
                        <point x="0" y="0" z="%g" />
                    </drawcylinder>
                    <drawbox>
                        <boxfill>bottom</boxfill>
                        <point x="%.4f" y="%.4f" z="0" />
                        <size x="%.4f" y="%.4f" z="%g" />
                    </drawbox>

                    <shapeout file="Tank" />
                </mainlist>
            </commands>
`,
		dp,
		-radius-margin, -radius-margin, -margin,
		radius+margin, radius+margin, height+margin,
		radius-dp, fluidHeight,
		radius, height,
		-radius, -radius, 2*radius, 2*radius, dp,
	)
}

// LShapedGeometry generates XML geometry for L-shaped tank (GEO-02).
func LShapedGeometry(L1, W1, L2, W2, H, fluidHeight, dp float64) string {
	margin := dp * 5

	return fmt.Sprintf(`
            <definition dp="%g" units_comment="metres (m)">
                <pointmin x="%.4f" y="%.4f" z="%.4f" />
                <pointmax x="%.4f" y="%.4f" z="%.4f" />
            </definition>
            <commands>
                <mainlist>
                    <setshapemode>dp | bound</setshapemode>
                    <setdrawmode mode="full" />

                    <!-- Fluid: L-shaped (two boxes) -->
                    <setmkfluid mk="0" />
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="%g" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>

                    <!-- Boundary: L-shaped walls -->
                    <setmkbound mk="0" />
                    <drawbox>
                        <boxfill>bottom | left | right | front | back</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>
                    <drawbox>
                        <boxfill>bottom | right | front | back</boxfill>
                        <point x="%g" y="0" z="0" />
                        <size x="%g" y="%g" z="%g" />
                    </drawbox>

                    <shapeout file="Tank" />
                </mainlist>
            </commands>
`,
		dp,
		-margin, -margin, -margin,
		math.Max(L1, L2)+margin, W1+W2+margin, H+margin,
		L1, W1, fluidHeight,
		L1, L2, W2, fluidHeight,
		L1, W1, H,
		L1, L2, W2, H,
	)
}

// MotionSeismicWave generates motion XML for seismic wave input (EXC-01).
func MotionSeismicWave(waveFile string, duration float64) string {
	return fmt.Sprintf(`
            <objreal ref="0">
                <begin mov="1" start="0" />
                <mvrectfile id="1" duration="%g" anglesunits="degrees">
                    <file name="%s" fields="2" fieldtime="0" fieldx="1" />
                </mvrectfile>
            </objreal>
`, duration, waveFile)
}

// MotionCustom generates motion XML for custom excitation (EXC-02).
func MotionCustom(timeValues, xValues []float64, duration float64) string {
	// Generate inline motion data
	dataStr := ""
	for i := 0; i < len(timeValues) && i < len(xValues); i++ {
		dataStr += fmt.Sprintf("                <data t=\"%g\" x=\"%g\" />\n", timeValues[i], xValues[i])
	}

	return fmt.Sprintf(`
            <objreal ref="0">
                <begin mov="1" start="0" />
                <mvrectfile id="1" duration="%g" anglesunits="degrees">
                    <values>
%s                    </values>
                </mvrectfile>
            </objreal>
`, duration, dataStr)
}
