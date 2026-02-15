package wizard

import (
	"fmt"
	"math"

	"github.com/opencode-ai/opencode/internal/llm/tools"
)

// Geometry types
type GeometryType = string

const (
	GeometryRectangular GeometryType = "rectangular"
	GeometryCylindrical GeometryType = "cylindrical"
	GeometryLShaped     GeometryType = "l_shaped"
	GeometrySTL         GeometryType = "stl"
)

// Fluid types
type FluidType = string

const (
	FluidWater FluidType = "water"
	FluidOil   FluidType = "oil"
	FluidLNG   FluidType = "lng"
)

// Excitation types
type ExcitationType = string

const (
	ExcitationSinusoidal ExcitationType = "sinusoidal"
	ExcitationSeismic    ExcitationType = "seismic"
)

// Boundary method types
type BoundaryMethodType = string

const (
	BoundaryDBC  BoundaryMethodType = "dbc"
	BoundaryMDBC BoundaryMethodType = "mdbc"
)

// Option represents a selectable option in the wizard.
type Option struct {
	Value string
	Label string
}

// WizardParams holds all parameters collected from the 4-step wizard.
type WizardParams struct {
	// Step 1: Tank Geometry
	Geometry   GeometryType
	TankLength float64
	TankWidth  float64
	TankHeight float64
	STLPath    string // Only for GeometrySTL

	// Step 2: Fluid Properties
	FluidType   FluidType
	FluidHeight float64
	Density     float64
	Viscosity   float64

	// Step 3: Excitation
	ExcitationType ExcitationType
	Freq           float64
	Amplitude      float64
	SeismicFile    string // Only for ExcitationSeismic

	// Step 4: Simulation Parameters
	DP             float64
	TimeMax        float64
	CFL            float64
	BoundaryMethod BoundaryMethodType
	OutPath        string

	// Internal: form state for huh string↔float conversion
	formState *formState
}

// DefaultWizardParams returns sensible defaults for a new simulation.
func DefaultWizardParams() WizardParams {
	return WizardParams{
		Geometry:       GeometryRectangular,
		TankLength:     1.0,
		TankWidth:      0.5,
		TankHeight:     0.6,
		FluidType:      FluidWater,
		FluidHeight:    0.3,
		Density:        1000.0,
		Viscosity:      0.01,
		ExcitationType: ExcitationSinusoidal,
		Freq:           0.5,
		Amplitude:      0.05,
		DP:             0.02,
		TimeMax:        10.0,
		CFL:            0.2,
		BoundaryMethod: BoundaryDBC,
		OutPath:        "simulations/sloshing_case",
	}
}

// AutoDP calculates dp from tank dimensions: min(L,W,H)/50, clamped to [0.005, 0.05].
func (p *WizardParams) AutoDP() float64 {
	minDim := math.Min(p.TankLength, math.Min(p.TankWidth, p.TankHeight))
	dp := minDim / 50.0
	dp = math.Max(dp, 0.005)
	dp = math.Min(dp, 0.05)
	return dp
}

// AutoTimeMax calculates time_max from frequency: 5/freq.
func (p *WizardParams) AutoTimeMax() float64 {
	if p.Freq <= 0 {
		return 10.0
	}
	return 5.0 / p.Freq
}

// Validate checks the wizard parameters for consistency.
func (p *WizardParams) Validate() error {
	if p.TankLength <= 0 {
		return fmt.Errorf("탱크 길이는 0보다 커야 합니다")
	}
	if p.TankWidth <= 0 {
		return fmt.Errorf("탱크 너비는 0보다 커야 합니다")
	}
	if p.TankHeight <= 0 {
		return fmt.Errorf("탱크 높이는 0보다 커야 합니다")
	}
	if p.FluidHeight <= 0 {
		return fmt.Errorf("유체 높이는 0보다 커야 합니다")
	}
	if p.FluidHeight > p.TankHeight {
		return fmt.Errorf("유체 높이(%.2fm)가 탱크 높이(%.2fm)보다 클 수 없습니다", p.FluidHeight, p.TankHeight)
	}
	if p.Freq <= 0 {
		return fmt.Errorf("가진 주파수는 0보다 커야 합니다")
	}
	if p.Amplitude <= 0 {
		return fmt.Errorf("가진 진폭은 0보다 커야 합니다")
	}
	return nil
}

// ToXMLGeneratorParams converts wizard parameters to XMLGeneratorParams for tool invocation.
func (p *WizardParams) ToXMLGeneratorParams() tools.XMLGeneratorParams {
	dp := p.DP
	if dp <= 0 {
		dp = p.AutoDP()
	}

	timeMax := p.TimeMax
	if timeMax <= 0 {
		timeMax = p.AutoTimeMax()
	}

	outPath := p.OutPath
	if outPath == "" {
		outPath = "simulations/sloshing_case"
	}

	return tools.XMLGeneratorParams{
		TankLength:     p.TankLength,
		TankWidth:      p.TankWidth,
		TankHeight:     p.TankHeight,
		FluidHeight:    p.FluidHeight,
		Freq:           p.Freq,
		Amplitude:      p.Amplitude,
		DP:             dp,
		TimeMax:        timeMax,
		OutPath:        outPath,
		Geometry:       string(p.Geometry),
		BoundaryMethod: string(p.BoundaryMethod),
	}
}

// FluidDensity returns the density for a given fluid type.
func FluidDensity(ft FluidType) float64 {
	switch ft {
	case FluidWater:
		return 1000.0
	case FluidOil:
		return 850.0
	case FluidLNG:
		return 450.0
	default:
		return 1000.0
	}
}

// GeometryOptions returns the available geometry options for the wizard.
func GeometryOptions() []Option {
	return []Option{
		{Value: GeometryRectangular, Label: "직사각형 탱크"},
		{Value: GeometryCylindrical, Label: "원통형 탱크"},
		{Value: GeometryLShaped, Label: "L형 탱크"},
		{Value: GeometrySTL, Label: "STL 파일 (CAD)"},
	}
}

// ExcitationOptions returns the available excitation options.
func ExcitationOptions() []Option {
	return []Option{
		{Value: ExcitationSinusoidal, Label: "정현파 (사인파)"},
		{Value: ExcitationSeismic, Label: "지진파/파도 데이터 (CSV)"},
	}
}

// FluidOptions returns the available fluid options.
func FluidOptions() []Option {
	return []Option{
		{Value: FluidWater, Label: "물 (1000 kg/m³)"},
		{Value: FluidOil, Label: "기름 (850 kg/m³)"},
		{Value: FluidLNG, Label: "LNG (450 kg/m³)"},
	}
}

// BoundaryOptions returns the available boundary method options.
func BoundaryOptions() []Option {
	return []Option{
		{Value: BoundaryDBC, Label: "DBC (기본, 빠름)"},
		{Value: BoundaryMDBC, Label: "mDBC (고정밀)"},
	}
}
