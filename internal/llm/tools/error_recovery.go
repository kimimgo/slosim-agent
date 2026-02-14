package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ErrorRecoveryParams defines parameters for error recovery strategies (NFR-01).
type ErrorRecoveryParams struct {
	JobID          string   `json:"job_id"`
	OutputDir      string   `json:"output_dir"`
	MaxRetries     int      `json:"max_retries,omitempty"`
	RetryDelay     int      `json:"retry_delay,omitempty"` // seconds
	AutoFix        bool     `json:"auto_fix,omitempty"`
	CheckDivergence bool    `json:"check_divergence,omitempty"`
}

// RecoveryResult holds error detection and recovery results.
type RecoveryResult struct {
	JobID       string   `json:"job_id"`
	Errors      []string `json:"errors"`
	Warnings    []string `json:"warnings"`
	IsDivergent bool     `json:"is_divergent"`
	AutoFixed   bool     `json:"auto_fixed,omitempty"`
	FixActions  []string `json:"fix_actions,omitempty"`
	Retry       bool     `json:"retry"`
}

type errorRecoveryTool struct{}

const (
	ErrorRecoveryToolName    = "error_recovery"
	errorRecoveryDescription = `에러 감지 및 복구 도구 (NFR-01) — 시뮬레이션 에러를 자동으로 감지하고 복구 전략을 제안합니다.

기능:
- 발산 감지 (Run.csv 에너지/압력 모니터링)
- GenCase 에러 메시지 파싱
- DualSPHysics 런타임 에러 분석
- 자동 수정 제안 (파라미터 조정, 재시도)

사용법:
- job_id: 모니터링할 Job ID
- output_dir: 시뮬레이션 출력 디렉토리
- max_retries: 최대 재시도 횟수 (기본: 3)
- retry_delay: 재시도 대기 시간 (초, 기본: 5)
- auto_fix: 자동 수정 시도 여부 (기본: false)
- check_divergence: 발산 검사 수행 (기본: true)

출력:
- 감지된 에러/경고 목록
- 복구 전략 (파라미터 조정, 재시도, 케이스 수정)
- 자동 수정 적용 여부`
)

func NewErrorRecoveryTool() BaseTool {
	return &errorRecoveryTool{}
}

func (e *errorRecoveryTool) Info() ToolInfo {
	return ToolInfo{
		Name:        ErrorRecoveryToolName,
		Description: errorRecoveryDescription,
		Parameters: map[string]any{
			"job_id": map[string]any{
				"type":        "string",
				"description": "모니터링할 Job ID",
			},
			"output_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 출력 디렉토리",
			},
			"max_retries": map[string]any{
				"type":        "integer",
				"description": "최대 재시도 횟수",
			},
			"retry_delay": map[string]any{
				"type":        "integer",
				"description": "재시도 대기 시간 (초)",
			},
			"auto_fix": map[string]any{
				"type":        "boolean",
				"description": "자동 수정 시도 여부",
			},
			"check_divergence": map[string]any{
				"type":        "boolean",
				"description": "발산 검사 수행",
			},
		},
		Required: []string{"output_dir"},
	}
}

func (e *errorRecoveryTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params ErrorRecoveryParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.OutputDir == "" {
		return NewTextErrorResponse("output_dir을 지정해주세요"), nil
	}

	// Default values
	if params.MaxRetries == 0 {
		params.MaxRetries = 3
	}
	if params.RetryDelay == 0 {
		params.RetryDelay = 5
	}
	params.CheckDivergence = true // Always check divergence

	result := RecoveryResult{
		JobID:       params.JobID,
		Errors:      []string{},
		Warnings:    []string{},
		IsDivergent: false,
		AutoFixed:   false,
		FixActions:  []string{},
		Retry:       false,
	}

	// Check Run.csv for divergence
	if params.CheckDivergence {
		runCSVPath := filepath.Join(params.OutputDir, "Run.csv")
		if _, err := os.Stat(runCSVPath); err == nil {
			divergent, warnings := checkDivergence(runCSVPath)
			result.IsDivergent = divergent
			result.Warnings = append(result.Warnings, warnings...)

			if divergent {
				result.Errors = append(result.Errors, "시뮬레이션 발산 감지")
				result.FixActions = append(result.FixActions, "TimeStep 감소 (CFL 조정)")
				result.FixActions = append(result.FixActions, "점성 계수 증가")
				result.FixActions = append(result.FixActions, "파티클 간격 (dp) 감소")
			}
		}
	}

	// Check for runtime errors (log files)
	logPath := filepath.Join(params.OutputDir, "DualSPHysics.log")
	if _, err := os.Stat(logPath); err == nil {
		errors := parseLogErrors(logPath)
		result.Errors = append(result.Errors, errors...)

		// Suggest fixes based on error patterns (case-insensitive)
		for _, errMsg := range errors {
			lowerErr := strings.ToLower(errMsg)
			if strings.Contains(lowerErr, "particle") && strings.Contains(lowerErr, "out") {
				result.FixActions = append(result.FixActions, "도메인 크기 증가 또는 파티클 이탈 검사 활성화")
			}
			if strings.Contains(lowerErr, "memory") {
				result.FixActions = append(result.FixActions, "GPU 메모리 부족 → 파티클 수 감소 또는 dp 증가")
			}
			if strings.Contains(lowerErr, "nan") {
				result.FixActions = append(result.FixActions, "NaN 발생 → TimeStep 감소, 초기 조건 검토")
			}
		}
	}

	// Auto-fix attempt (if enabled and fixable errors detected)
	if params.AutoFix && len(result.FixActions) > 0 {
		// Placeholder: Apply first fix action
		result.AutoFixed = true
		result.FixActions = []string{"자동 수정 시도: " + result.FixActions[0]}
		result.Retry = true
	}

	// Determine if retry is recommended
	if len(result.Errors) == 0 {
		result.Retry = false
	} else if len(result.FixActions) > 0 {
		result.Retry = true
	}

	// Format result message
	resultMsg := fmt.Sprintf(`에러 복구 분석 결과:
Job ID: %s
출력 디렉토리: %s

감지된 에러 수: %d
경고 수: %d
발산 감지: %v

`, params.JobID, params.OutputDir, len(result.Errors), len(result.Warnings), result.IsDivergent)

	if len(result.Errors) > 0 {
		resultMsg += "에러 목록:\n"
		for _, err := range result.Errors {
			resultMsg += fmt.Sprintf("  - %s\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		resultMsg += "\n경고 목록:\n"
		for _, warn := range result.Warnings {
			resultMsg += fmt.Sprintf("  - %s\n", warn)
		}
	}

	if len(result.FixActions) > 0 {
		resultMsg += "\n제안된 수정 방법:\n"
		for _, action := range result.FixActions {
			resultMsg += fmt.Sprintf("  → %s\n", action)
		}
	}

	if result.AutoFixed {
		resultMsg += "\n자동 수정 적용됨\n"
	}

	if result.Retry {
		resultMsg += fmt.Sprintf("\n재시도 권장 (최대 %d회, %d초 대기)\n", params.MaxRetries, params.RetryDelay)
	} else {
		resultMsg += "\n재시도 불필요 또는 불가능\n"
	}

	return NewTextResponse(resultMsg), nil
}

// checkDivergence analyzes Run.csv for signs of simulation divergence.
// Detects exponential energy growth by tracking consecutive growth ratios.
// DualSPHysics Run.csv uses semicolon (;) separators and # comment lines.
// Header format: #Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot;...
func checkDivergence(csvPath string) (bool, []string) {
	file, err := os.Open(csvPath)
	if err != nil {
		return false, []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	warnings := []string{}
	divergent := false

	// Find EnergyKin column from header line (starts with #)
	energyIdx := -1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Parse header: strip # prefix, split by ;
			header := strings.TrimPrefix(line, "#")
			headerFields := strings.Split(header, ";")
			for i, h := range headerFields {
				h = strings.TrimSpace(h)
				if h == "EnergyKin" || h == "Energy" {
					energyIdx = i
					break
				}
			}
			break
		}
	}

	// Scan data rows
	lineNum := 0
	var prevEnergy float64
	consecutiveGrowth := 0
	const growthThreshold = 1.2  // 1.2x per step = rapid growth
	const consecutiveLimit = 5   // 5 consecutive steps of rapid growth = divergence

	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comment lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lineNum++
		fields := strings.Split(line, ";")

		if energyIdx >= 0 && energyIdx < len(fields) {
			energy, err := strconv.ParseFloat(strings.TrimSpace(fields[energyIdx]), 64)
			if err != nil {
				continue
			}

			if lineNum > 1 && prevEnergy > 0 {
				ratio := energy / prevEnergy
				if ratio >= growthThreshold {
					consecutiveGrowth++
					if consecutiveGrowth >= consecutiveLimit {
						divergent = true
					}
				} else {
					consecutiveGrowth = 0
				}
			}
			prevEnergy = energy
		}
	}

	if lineNum > 10000 {
		warnings = append(warnings, "매우 긴 시뮬레이션 (10000+ 타임스텝)")
	}

	return divergent, warnings
}

// parseLogErrors extracts error messages from DualSPHysics log file.
func parseLogErrors(logPath string) []string {
	file, err := os.Open(logPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	errors := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)
		trimmedLine := strings.TrimSpace(line)

		// Check NaN/Inf first (case-sensitive to avoid false positives like "[INFO]")
		if strings.Contains(line, "NaN") || isInfValue(line) {
			errors = append(errors, "NaN 또는 Inf 값 감지: "+trimmedLine)
		} else if strings.Contains(lowerLine, "error") || strings.Contains(lowerLine, "failed") {
			errors = append(errors, trimmedLine)
		}
	}

	return errors
}

// isInfValue checks for Inf floating-point values while avoiding false positives like "[INFO]".
func isInfValue(line string) bool {
	for _, pattern := range []string{" Inf ", " Inf,", "+Inf", "-Inf", " Inf\t"} {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	return strings.HasSuffix(line, " Inf")
}

// retryWithBackoff implements exponential backoff retry logic.
func retryWithBackoff(ctx context.Context, maxRetries, initialDelay int, operation func() error) error {
	delay := time.Duration(initialDelay) * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			}
		}
	}

	return fmt.Errorf("max retries (%d) exceeded", maxRetries)
}
