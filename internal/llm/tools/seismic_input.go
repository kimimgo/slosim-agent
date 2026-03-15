package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// SeismicInputParams defines parameters for seismic wave file parsing (EXC-01).
type SeismicInputParams struct {
	FilePath       string  `json:"file_path"`
	TimeColumn     int     `json:"time_column,omitempty"`
	AccelColumn    int     `json:"accel_column,omitempty"`
	SkipRows       int     `json:"skip_rows,omitempty"`
	Delimiter      string  `json:"delimiter,omitempty"`
	ScaleFactor    float64 `json:"scale_factor,omitempty"`
	ValidateOnly   bool    `json:"validate_only,omitempty"`
	OutputFormat   string  `json:"output_format,omitempty"` // "dsph" (DualSPHysics format) or "raw"
}

type seismicInputTool struct{}

const (
	SeismicInputToolName    = "seismic_input"
	seismicInputDescription = `지진파 시계열 데이터 파일 파싱 및 검증 도구 (EXC-01).

사용법:
- file_path: 지진파 데이터 파일 경로 (.txt, .csv, .dat 등)
- time_column: 시간 열 인덱스 (기본: 0, 0부터 시작)
- accel_column: 가속도 열 인덱스 (기본: 1)
- skip_rows: 헤더 건너뛸 행 수 (기본: 0)
- delimiter: 구분자 (기본: 공백/탭 자동 감지)
- scale_factor: 가속도 스케일링 계수 (기본: 1.0, 단위 변환용)
- validate_only: true면 검증만 수행, false면 DualSPHysics 포맷으로 변환
- output_format: "dsph" (DualSPHysics 2열 포맷) 또는 "raw" (원본)

출력:
- 검증 결과 (샘플 수, 지속 시간, 최대/최소 가속도)
- validate_only=false 시 변환된 파일 경로 반환`
)

func NewSeismicInputTool() BaseTool {
	return &seismicInputTool{}
}

func (s *seismicInputTool) Info() ToolInfo {
	return ToolInfo{
		Name:        SeismicInputToolName,
		Description: seismicInputDescription,
		Parameters: map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "지진파 데이터 파일 경로",
			},
			"time_column": map[string]any{
				"type":        "integer",
				"description": "시간 열 인덱스 (0부터 시작)",
			},
			"accel_column": map[string]any{
				"type":        "integer",
				"description": "가속도 열 인덱스",
			},
			"skip_rows": map[string]any{
				"type":        "integer",
				"description": "헤더 건너뛸 행 수",
			},
			"delimiter": map[string]any{
				"type":        "string",
				"description": "구분자 (공백/탭/쉼표)",
			},
			"scale_factor": map[string]any{
				"type":        "number",
				"description": "가속도 스케일링 계수",
			},
			"validate_only": map[string]any{
				"type":        "boolean",
				"description": "검증만 수행 (파일 생성 안 함)",
			},
			"output_format": map[string]any{
				"type":        "string",
				"description": "출력 포맷 (dsph 또는 raw)",
			},
		},
		Required: []string{"file_path"},
	}
}

func (s *seismicInputTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params SeismicInputParams
	if err := UnmarshalToolInput(call.Input, &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.FilePath == "" {
		return NewTextErrorResponse("file_path를 지정해주세요"), nil
	}

	// Default values
	if params.TimeColumn == 0 && params.AccelColumn == 0 {
		params.AccelColumn = 1
	}
	if params.ScaleFactor == 0 {
		params.ScaleFactor = 1.0
	}
	if params.OutputFormat == "" {
		params.OutputFormat = "dsph"
	}

	// Validate file existence
	absPath := params.FilePath
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(getWorkingDirectory(), absPath)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("파일을 찾을 수 없습니다: %s", absPath)), nil
	}

	// Parse seismic data
	timeData, accelData, err := parseSeismicFile(absPath, params)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("지진파 파일 파싱 실패: %s", err)), nil
	}

	// Validation
	if len(timeData) == 0 {
		return NewTextErrorResponse("데이터가 비어 있습니다"), nil
	}

	duration := timeData[len(timeData)-1] - timeData[0]
	minAccel, maxAccel := findMinMax(accelData)
	dt := 0.0
	if len(timeData) > 1 {
		dt = timeData[1] - timeData[0]
	}

	validationMsg := fmt.Sprintf(`지진파 데이터 검증 완료:
- 샘플 수: %d
- 지속 시간: %.3f초
- 시간 간격 (dt): %.4f초 (평균)
- 가속도 범위: [%.4f, %.4f] m/s²
- 스케일 계수: %.2f
`,
		len(timeData),
		duration,
		dt,
		minAccel*params.ScaleFactor,
		maxAccel*params.ScaleFactor,
		params.ScaleFactor,
	)

	if params.ValidateOnly {
		return NewTextResponse(validationMsg), nil
	}

	// Convert to DualSPHysics format
	outputPath := generateOutputPath(absPath, params.OutputFormat)
	if err := writeSeismicOutput(outputPath, timeData, accelData, params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("출력 파일 생성 실패: %s", err)), nil
	}

	resultMsg := validationMsg + fmt.Sprintf("\n변환된 파일 경로: %s\n", outputPath)
	return NewTextResponse(resultMsg), nil
}

func parseSeismicFile(filePath string, params SeismicInputParams) ([]float64, []float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var timeData, accelData []float64
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum <= params.SkipRows {
			continue
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Detect delimiter
		delimiter := params.Delimiter
		if delimiter == "" {
			if strings.Contains(line, "\t") {
				delimiter = "\t"
			} else if strings.Contains(line, ",") {
				delimiter = ","
			} else {
				delimiter = " "
			}
		}

		fields := strings.Fields(line)
		if delimiter != " " {
			fields = strings.Split(line, delimiter)
		}

		// Clean empty fields
		cleanFields := []string{}
		for _, f := range fields {
			if strings.TrimSpace(f) != "" {
				cleanFields = append(cleanFields, strings.TrimSpace(f))
			}
		}
		fields = cleanFields

		if len(fields) <= params.TimeColumn || len(fields) <= params.AccelColumn {
			continue
		}

		timeVal, err := strconv.ParseFloat(fields[params.TimeColumn], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: 시간 값 파싱 실패: %s", lineNum, err)
		}

		accelVal, err := strconv.ParseFloat(fields[params.AccelColumn], 64)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: 가속도 값 파싱 실패: %s", lineNum, err)
		}

		timeData = append(timeData, timeVal)
		accelData = append(accelData, accelVal)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return timeData, accelData, nil
}

func writeSeismicOutput(outputPath string, timeData, accelData []float64, params SeismicInputParams) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Write header comment
	fmt.Fprintf(writer, "# DualSPHysics seismic wave input\n")
	fmt.Fprintf(writer, "# Columns: time(s) acceleration(m/s^2)\n")
	fmt.Fprintf(writer, "# Scale factor: %.4f\n", params.ScaleFactor)

	for i := 0; i < len(timeData); i++ {
		fmt.Fprintf(writer, "%.6f %.8f\n", timeData[i], accelData[i]*params.ScaleFactor)
	}

	return nil
}

func generateOutputPath(inputPath, format string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	suffix := "_dsph"
	if format == "raw" {
		suffix = "_converted"
	}

	return filepath.Join(dir, name+suffix+".dat")
}

func findMinMax(data []float64) (float64, float64) {
	if len(data) == 0 {
		return 0, 0
	}

	min, max := data[0], data[0]
	for _, v := range data[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
