package wizard

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
)

// NewWizardForm creates a 4-step huh form for simulation setup.
// Returns the form and a pointer to the params that will be populated.
func NewWizardForm() (*huh.Form, *WizardParams) {
	params := DefaultWizardParams()

	// String holders for numeric inputs (huh Input returns strings)
	var (
		geometryStr       string = string(params.Geometry)
		fluidTypeStr      string = string(params.FluidType)
		excitationTypeStr string = string(params.ExcitationType)
		boundaryStr       string = string(params.BoundaryMethod)
		tankLengthStr     string = fmt.Sprintf("%.2f", params.TankLength)
		tankWidthStr      string = fmt.Sprintf("%.2f", params.TankWidth)
		tankHeightStr     string = fmt.Sprintf("%.2f", params.TankHeight)
		fluidHeightStr    string = fmt.Sprintf("%.2f", params.FluidHeight)
		freqStr           string = fmt.Sprintf("%.2f", params.Freq)
		amplitudeStr      string = fmt.Sprintf("%.3f", params.Amplitude)
		dpStr             string = fmt.Sprintf("%.3f", params.DP)
		timeMaxStr        string = fmt.Sprintf("%.1f", params.TimeMax)
		outPathStr        string = params.OutPath
	)

	form := huh.NewForm(
		// Step 1: Tank Geometry
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("탱크 형상").
				Description("시뮬레이션할 탱크의 형태를 선택하세요").
				Options(
					huh.NewOption("직사각형 탱크", GeometryRectangular),
					huh.NewOption("원통형 탱크", GeometryCylindrical),
					huh.NewOption("L형 탱크", GeometryLShaped),
					huh.NewOption("STL 파일 (CAD)", GeometrySTL),
				).
				Value(&geometryStr),
			huh.NewInput().
				Title("탱크 길이 (m)").
				Value(&tankLengthStr).
				Validate(validatePositiveFloat),
			huh.NewInput().
				Title("탱크 너비 (m)").
				Value(&tankWidthStr).
				Validate(validatePositiveFloat),
			huh.NewInput().
				Title("탱크 높이 (m)").
				Value(&tankHeightStr).
				Validate(validatePositiveFloat),
		).Title("1단계: 탱크 형상"),

		// Step 2: Fluid Properties
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("유체 종류").
				Description("시뮬레이션할 유체를 선택하세요").
				Options(
					huh.NewOption("물 (1000 kg/m³)", FluidWater),
					huh.NewOption("기름 (850 kg/m³)", FluidOil),
					huh.NewOption("LNG (450 kg/m³)", FluidLNG),
				).
				Value(&fluidTypeStr),
			huh.NewInput().
				Title("유체 높이 (m)").
				Description("탱크 높이 이내로 설정하세요").
				Value(&fluidHeightStr).
				Validate(validatePositiveFloat),
		).Title("2단계: 유체 속성"),

		// Step 3: Excitation
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("가진 방식").
				Description("탱크를 흔드는 방식을 선택하세요").
				Options(
					huh.NewOption("정현파 (사인파)", ExcitationSinusoidal),
					huh.NewOption("지진파/파도 데이터 (CSV)", ExcitationSeismic),
				).
				Value(&excitationTypeStr),
			huh.NewInput().
				Title("가진 주파수 (Hz)").
				Value(&freqStr).
				Validate(validatePositiveFloat),
			huh.NewInput().
				Title("가진 진폭 (m)").
				Value(&amplitudeStr).
				Validate(validatePositiveFloat),
		).Title("3단계: 가진 조건"),

		// Step 4: Simulation Parameters
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("경계 조건 방식").
				Description("DBC는 빠르지만 정밀도 낮음, mDBC는 고정밀").
				Options(
					huh.NewOption("DBC (기본, 빠름)", BoundaryDBC),
					huh.NewOption("mDBC (고정밀)", BoundaryMDBC),
				).
				Value(&boundaryStr),
			huh.NewInput().
				Title("파티클 간격 dp (m)").
				Description("작을수록 정밀, 0이면 자동 계산").
				Value(&dpStr).
				Validate(validateNonNegativeFloat),
			huh.NewInput().
				Title("시뮬레이션 시간 (초)").
				Value(&timeMaxStr).
				Validate(validatePositiveFloat),
			huh.NewInput().
				Title("출력 경로").
				Description(".xml 확장자는 자동 추가됩니다").
				Value(&outPathStr),
		).Title("4단계: 시뮬레이션 파라미터"),
	).WithTheme(huh.ThemeCatppuccin())

	// After form completion, populate params from string values
	// This is done in CollectParams after form.Run()
	_ = &params // Keep reference for later use in CollectParams

	return form, &WizardParams{
		// Pointers to the string holders — CollectParams will parse them
		formState: &formState{
			geometryStr:       &geometryStr,
			fluidTypeStr:      &fluidTypeStr,
			excitationTypeStr: &excitationTypeStr,
			boundaryStr:       &boundaryStr,
			tankLengthStr:     &tankLengthStr,
			tankWidthStr:      &tankWidthStr,
			tankHeightStr:     &tankHeightStr,
			fluidHeightStr:    &fluidHeightStr,
			freqStr:           &freqStr,
			amplitudeStr:      &amplitudeStr,
			dpStr:             &dpStr,
			timeMaxStr:        &timeMaxStr,
			outPathStr:        &outPathStr,
		},
		Geometry:       GeometryRectangular,
		FluidType:      FluidWater,
		ExcitationType: ExcitationSinusoidal,
		BoundaryMethod: BoundaryDBC,
		TankLength:     1.0,
		TankWidth:      0.5,
		TankHeight:     0.6,
		FluidHeight:    0.3,
		Freq:           0.5,
		Amplitude:      0.05,
		DP:             0.02,
		TimeMax:        10.0,
		OutPath:        "simulations/sloshing_case",
	}
}

// formState holds pointers to the string values used by huh form inputs.
type formState struct {
	geometryStr       *string
	fluidTypeStr      *string
	excitationTypeStr *string
	boundaryStr       *string
	tankLengthStr     *string
	tankWidthStr      *string
	tankHeightStr     *string
	fluidHeightStr    *string
	freqStr           *string
	amplitudeStr      *string
	dpStr             *string
	timeMaxStr        *string
	outPathStr        *string
}

// CollectParams parses the form's string values into the WizardParams struct.
// Call this after form.Run() completes successfully.
func (p *WizardParams) CollectParams() error {
	if p.formState == nil {
		return nil // Already has direct values
	}

	s := p.formState

	p.Geometry = *s.geometryStr
	p.FluidType = *s.fluidTypeStr
	p.ExcitationType = *s.excitationTypeStr
	p.BoundaryMethod = *s.boundaryStr
	p.OutPath = *s.outPathStr
	p.Density = FluidDensity(p.FluidType)

	var err error
	if p.TankLength, err = strconv.ParseFloat(*s.tankLengthStr, 64); err != nil {
		return fmt.Errorf("탱크 길이 파싱 실패: %w", err)
	}
	if p.TankWidth, err = strconv.ParseFloat(*s.tankWidthStr, 64); err != nil {
		return fmt.Errorf("탱크 너비 파싱 실패: %w", err)
	}
	if p.TankHeight, err = strconv.ParseFloat(*s.tankHeightStr, 64); err != nil {
		return fmt.Errorf("탱크 높이 파싱 실패: %w", err)
	}
	if p.FluidHeight, err = strconv.ParseFloat(*s.fluidHeightStr, 64); err != nil {
		return fmt.Errorf("유체 높이 파싱 실패: %w", err)
	}
	if p.Freq, err = strconv.ParseFloat(*s.freqStr, 64); err != nil {
		return fmt.Errorf("주파수 파싱 실패: %w", err)
	}
	if p.Amplitude, err = strconv.ParseFloat(*s.amplitudeStr, 64); err != nil {
		return fmt.Errorf("진폭 파싱 실패: %w", err)
	}
	if p.DP, err = strconv.ParseFloat(*s.dpStr, 64); err != nil {
		return fmt.Errorf("dp 파싱 실패: %w", err)
	}
	if p.TimeMax, err = strconv.ParseFloat(*s.timeMaxStr, 64); err != nil {
		return fmt.Errorf("시간 파싱 실패: %w", err)
	}

	return nil
}

// validatePositiveFloat validates that a string is a positive float.
func validatePositiveFloat(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("숫자를 입력하세요")
	}
	if v <= 0 {
		return fmt.Errorf("0보다 큰 값을 입력하세요")
	}
	return nil
}

// validateNonNegativeFloat validates that a string is a non-negative float.
func validateNonNegativeFloat(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("숫자를 입력하세요")
	}
	if v < 0 {
		return fmt.Errorf("0 이상의 값을 입력하세요")
	}
	return nil
}
