package wizard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TUI-04: Case Wizard Tests

func TestWizardParams_DefaultValues(t *testing.T) {
	params := DefaultWizardParams()

	assert.Equal(t, GeometryRectangular, params.Geometry)
	assert.Equal(t, FluidWater, params.FluidType)
	assert.Equal(t, ExcitationSinusoidal, params.ExcitationType)
	assert.Equal(t, BoundaryDBC, params.BoundaryMethod)
	assert.Greater(t, params.TankLength, 0.0)
	assert.Greater(t, params.TankWidth, 0.0)
	assert.Greater(t, params.TankHeight, 0.0)
}

func TestWizardParams_ToXMLGeneratorParams(t *testing.T) {
	t.Run("rectangular tank with sinusoidal excitation", func(t *testing.T) {
		params := DefaultWizardParams()
		params.TankLength = 1.0
		params.TankWidth = 0.5
		params.TankHeight = 0.6
		params.FluidHeight = 0.3
		params.Freq = 0.5
		params.Amplitude = 0.05
		params.DP = 0.02
		params.TimeMax = 10.0
		params.OutPath = "simulations/test_case"

		xmlParams := params.ToXMLGeneratorParams()

		assert.Equal(t, 1.0, xmlParams.TankLength)
		assert.Equal(t, 0.5, xmlParams.TankWidth)
		assert.Equal(t, 0.6, xmlParams.TankHeight)
		assert.Equal(t, 0.3, xmlParams.FluidHeight)
		assert.Equal(t, 0.5, xmlParams.Freq)
		assert.Equal(t, 0.05, xmlParams.Amplitude)
		assert.Equal(t, 0.02, xmlParams.DP)
		assert.Equal(t, 10.0, xmlParams.TimeMax)
		assert.Equal(t, "simulations/test_case", xmlParams.OutPath)
	})

	t.Run("mdbc boundary method", func(t *testing.T) {
		params := DefaultWizardParams()
		params.BoundaryMethod = BoundaryMDBC

		xmlParams := params.ToXMLGeneratorParams()
		assert.Equal(t, "mdbc", xmlParams.BoundaryMethod)
	})

	t.Run("dbc boundary method", func(t *testing.T) {
		params := DefaultWizardParams()
		params.BoundaryMethod = BoundaryDBC

		xmlParams := params.ToXMLGeneratorParams()
		assert.Equal(t, "dbc", xmlParams.BoundaryMethod)
	})
}

func TestWizardParams_AutoDP(t *testing.T) {
	params := DefaultWizardParams()
	params.TankLength = 1.0
	params.TankWidth = 0.5
	params.TankHeight = 0.6
	params.DP = 0 // Not set — should auto-calculate

	dp := params.AutoDP()
	// dp = min(L,W,H) / 50 = 0.5/50 = 0.01
	assert.InDelta(t, 0.01, dp, 0.001)
}

func TestWizardParams_AutoDP_Bounds(t *testing.T) {
	t.Run("minimum dp is 0.005", func(t *testing.T) {
		params := DefaultWizardParams()
		params.TankLength = 0.1 // Very small tank
		params.TankWidth = 0.1
		params.TankHeight = 0.1

		dp := params.AutoDP()
		// min(0.1,0.1,0.1)/50 = 0.002 → clamped to 0.005
		assert.InDelta(t, 0.005, dp, 0.001)
	})

	t.Run("maximum dp is 0.05", func(t *testing.T) {
		params := DefaultWizardParams()
		params.TankLength = 40.0 // Very large tank
		params.TankWidth = 30.0
		params.TankHeight = 27.0

		dp := params.AutoDP()
		// min(40,30,27)/50 = 27/50 = 0.54 → clamped to 0.05
		assert.InDelta(t, 0.05, dp, 0.001)
	})
}

func TestWizardParams_AutoTimeMax(t *testing.T) {
	params := DefaultWizardParams()
	params.Freq = 0.5

	tmax := params.AutoTimeMax()
	// time_max = 5/freq = 5/0.5 = 10.0
	assert.InDelta(t, 10.0, tmax, 0.1)
}

func TestWizardParams_Validate(t *testing.T) {
	t.Run("valid params pass", func(t *testing.T) {
		params := DefaultWizardParams()
		err := params.Validate()
		require.NoError(t, err)
	})

	t.Run("zero tank length fails", func(t *testing.T) {
		params := DefaultWizardParams()
		params.TankLength = 0
		err := params.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "탱크 길이")
	})

	t.Run("fluid higher than tank fails", func(t *testing.T) {
		params := DefaultWizardParams()
		params.FluidHeight = params.TankHeight + 0.1
		err := params.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "유체 높이")
	})

	t.Run("zero freq fails", func(t *testing.T) {
		params := DefaultWizardParams()
		params.Freq = 0
		err := params.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "주파수")
	})
}

func TestFluidDensity(t *testing.T) {
	assert.Equal(t, 1000.0, FluidDensity(FluidWater))
	assert.Equal(t, 850.0, FluidDensity(FluidOil))
	assert.Equal(t, 450.0, FluidDensity(FluidLNG))
	assert.Equal(t, 1000.0, FluidDensity("unknown")) // default
}

func TestGeometryOptions(t *testing.T) {
	opts := GeometryOptions()
	assert.Len(t, opts, 4) // rectangular, cylindrical, l_shaped, stl
	assert.Equal(t, GeometryRectangular, opts[0].Value)
}

func TestExcitationOptions(t *testing.T) {
	opts := ExcitationOptions()
	assert.Len(t, opts, 2) // sinusoidal, seismic
	assert.Equal(t, ExcitationSinusoidal, opts[0].Value)
}
