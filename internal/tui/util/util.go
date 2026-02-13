package util

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func ReportError(err error) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeError,
		Msg:  err.Error(),
	})
}

type InfoType int

const (
	InfoTypeInfo InfoType = iota
	InfoTypeWarn
	InfoTypeError
)

func ReportInfo(info string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeInfo,
		Msg:  info,
	})
}

func ReportWarn(warn string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeWarn,
		Msg:  warn,
	})
}

type (
	InfoMsg struct {
		Type InfoType
		Msg  string
		TTL  time.Duration
	}
	ClearStatusMsg struct{}

	// Job status messages for DualSPHysics simulation pipeline
	JobUpdateMsg struct {
		JobID       string
		ProgressPct int
		ETA         string // e.g. "약 2분 남음"
	}
	JobCompletedMsg struct {
		JobID     string
		ResultDir string // e.g. "simulations/20260213_143022/"
	}
	JobFailedMsg struct {
		JobID   string
		Message string // 한국어 에러 메시지
	}
	ClearJobStatusMsg struct{}
)

func Clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
