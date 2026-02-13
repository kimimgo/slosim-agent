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
)

// MonitorParams defines parameters for real-time simulation monitoring (MON-01).
type MonitorParams struct {
	SimDir string `json:"sim_dir"`
}

// MonitorStatus represents the current simulation state.
type MonitorStatus struct {
	CurrentTime    float64 `json:"current_time"`
	ProgressPct    float64 `json:"progress_pct"`
	ParticleCount  int     `json:"particle_count"`
	EnergyKin      float64 `json:"energy_kin,omitempty"`
	EnergyPot      float64 `json:"energy_pot,omitempty"`
	ParticlesOut   int     `json:"particles_out,omitempty"`
	IsUnstable     bool    `json:"is_unstable"`
	UnstableReason string  `json:"unstable_reason,omitempty"`
}

type monitorTool struct{}

const (
	MonitorToolName    = "monitor"
	monitorDescription = `시뮬레이션 실시간 모니터링 도구 (MON-01) — Run.csv를 파싱하여 진행 상황 및 안정성을 확인합니다.

사용법:
- sim_dir: 시뮬레이션 출력 디렉토리 (Run.csv 위치)

반환 정보:
- 현재 시뮬레이션 시간 및 진행률
- 파티클 수
- 운동/위치 에너지
- 발산/불안정 감지

Docker 없이 로컬에서 파일을 직접 읽습니다.`
)

func NewMonitorTool() BaseTool {
	return &monitorTool{}
}

func (m *monitorTool) Info() ToolInfo {
	return ToolInfo{
		Name:        MonitorToolName,
		Description: monitorDescription,
		Parameters: map[string]any{
			"sim_dir": map[string]any{
				"type":        "string",
				"description": "시뮬레이션 출력 디렉토리",
			},
		},
		Required: []string{"sim_dir"},
	}
}

func (m *monitorTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params MonitorParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	if params.SimDir == "" {
		return NewTextErrorResponse("시뮬레이션 디렉토리(sim_dir)를 지정해주세요"), nil
	}

	// Parse Run.csv
	runCSV := filepath.Join(params.SimDir, "Run.csv")
	if _, err := os.Stat(runCSV); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("Run.csv 파일을 찾을 수 없습니다: %s", runCSV)), nil
	}

	status, err := parseRunCSV(runCSV)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Run.csv 파싱 실패: %s", err)), nil
	}

	// Format response
	result, _ := json.MarshalIndent(status, "", "  ")
	msg := fmt.Sprintf("모니터링 결과:\n%s\n", string(result))

	if status.IsUnstable {
		msg += fmt.Sprintf("\n⚠️  해석이 불안정합니다: %s\n", status.UnstableReason)
		return NewTextResponse(msg), nil
	}

	msg += "\n✅ 해석이 정상적으로 진행 중입니다.\n"
	return NewTextResponse(msg), nil
}

func parseRunCSV(csvPath string) (*MonitorStatus, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.HasPrefix(line, "#") {
			lastLine = line
		}
	}

	if lastLine == "" {
		return nil, fmt.Errorf("Run.csv에 데이터가 없습니다")
	}

	// Parse last line (latest timestep)
	// Format: Time;TotalSteps;Nparticles;Nfloat;Nbound;PartOut;EnergyKin;EnergyPot;...
	fields := strings.Split(lastLine, ";")
	if len(fields) < 6 {
		return nil, fmt.Errorf("Run.csv 형식이 올바르지 않습니다")
	}

	currentTime, _ := strconv.ParseFloat(fields[0], 64)
	particleCount, _ := strconv.Atoi(fields[2])
	particlesOut, _ := strconv.Atoi(fields[5])

	var energyKin, energyPot float64
	if len(fields) >= 8 {
		energyKin, _ = strconv.ParseFloat(fields[6], 64)
		energyPot, _ = strconv.ParseFloat(fields[7], 64)
	}

	// Detect instability
	isUnstable := false
	unstableReason := ""

	if particlesOut > particleCount/10 {
		isUnstable = true
		unstableReason = fmt.Sprintf("%d개 파티클이 이탈했습니다 (전체의 %.1f%%)", particlesOut, float64(particlesOut)/float64(particleCount)*100)
	}

	// Calculate progress (assume TimeMax is known from context, here we use a placeholder)
	progressPct := currentTime / 10.0 * 100 // Placeholder: assume TimeMax=10s

	return &MonitorStatus{
		CurrentTime:    currentTime,
		ProgressPct:    progressPct,
		ParticleCount:  particleCount,
		EnergyKin:      energyKin,
		EnergyPot:      energyPot,
		ParticlesOut:   particlesOut,
		IsUnstable:     isUnstable,
		UnstableReason: unstableReason,
	}, nil
}
