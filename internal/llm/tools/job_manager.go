package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type JobManagerParams struct {
	Action  string   `json:"action"`
	JobID   string   `json:"job_id,omitempty"`
	Command []string `json:"command,omitempty"`
	WorkDir string   `json:"work_dir,omitempty"`
}

// JobStatus represents the state of a managed job.
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
	JobStatusCancelled JobStatus = "CANCELLED"
)

type managedJob struct {
	ID        string    `json:"id"`
	Status    JobStatus `json:"status"`
	WorkDir   string    `json:"work_dir"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Error     string    `json:"error,omitempty"`
	cancel    context.CancelFunc
}

type jobManagerTool struct {
	mu   sync.Mutex
	jobs map[string]*managedJob // JOB-02: Support multiple jobs
	maxConcurrent int
}

const (
	JobManagerToolName    = "job_manager"
	jobManagerDescription = `백그라운드 Job 관리 도구 — 시뮬레이션 등 장시간 작업을 백그라운드에서 관리합니다.

사용 가능한 액션:
- submit: 새 Job 제출 (command, work_dir 필수). 즉시 Job ID 반환.
- status: Job 상태 조회 (job_id 필수). 상태, 진행률, 경과 시간 반환.
- list: 모든 Job 목록 조회.
- cancel: 실행 중인 Job 취소 (job_id 필수).

v0.3: 동시 최대 3개 Job 실행 가능합니다.`
)

func NewJobManagerTool() BaseTool {
	return &jobManagerTool{
		jobs:          make(map[string]*managedJob),
		maxConcurrent: 3, // v0.3: Allow 3 concurrent jobs
	}
}

func (j *jobManagerTool) Info() ToolInfo {
	return ToolInfo{
		Name:        JobManagerToolName,
		Description: jobManagerDescription,
		Parameters: map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "수행할 액션 (submit, status, cancel)",
			},
			"job_id": map[string]any{
				"type":        "string",
				"description": "Job ID (status, cancel 시 필수)",
			},
			"command": map[string]any{
				"type":        "array",
				"description": "실행할 명령어 (submit 시 필수)",
				"items": map[string]any{
					"type": "string",
				},
			},
			"work_dir": map[string]any{
				"type":        "string",
				"description": "작업 디렉토리 (submit 시 필수)",
			},
		},
		Required: []string{"action"},
	}
}

func (j *jobManagerTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params JobManagerParams
	if err := json.Unmarshal([]byte(call.Input), &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	switch params.Action {
	case "submit":
		return j.handleSubmit(ctx, params)
	case "status":
		return j.handleStatus(params)
	case "list":
		return j.handleList()
	case "cancel":
		return j.handleCancel(params)
	default:
		return NewTextErrorResponse(fmt.Sprintf("알 수 없는 액션입니다: %s (사용 가능: submit, status, list, cancel)", params.Action)), nil
	}
}

func (j *jobManagerTool) handleSubmit(_ context.Context, params JobManagerParams) (ToolResponse, error) {
	if len(params.Command) == 0 {
		return NewTextErrorResponse("실행할 명령어(command)를 지정해주세요"), nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	// v0.3: Check concurrent job limit
	runningCount := 0
	for _, job := range j.jobs {
		if job.Status == JobStatusRunning || job.Status == JobStatusPending {
			runningCount++
		}
	}

	if runningCount >= j.maxConcurrent {
		return NewTextErrorResponse(fmt.Sprintf("동시 실행 가능한 Job 수(%d)를 초과했습니다. 일부 Job이 완료될 때까지 기다려주세요", j.maxConcurrent)), nil
	}

	jobID := generateJobID()
	jobCtx, cancel := context.WithCancel(context.Background())

	job := &managedJob{
		ID:        jobID,
		Status:    JobStatusRunning,
		WorkDir:   params.WorkDir,
		StartTime: time.Now(),
		cancel:    cancel,
	}
	j.jobs[jobID] = job

	// Launch in background goroutine
	go func() {
		cmd := exec.CommandContext(jobCtx, params.Command[0], params.Command[1:]...)
		if params.WorkDir != "" {
			cmd.Dir = params.WorkDir
		}

		err := cmd.Run()

		j.mu.Lock()
		defer j.mu.Unlock()

		job.EndTime = time.Now()
		if err != nil {
			if jobCtx.Err() == context.Canceled {
				job.Status = JobStatusCancelled
			} else {
				job.Status = JobStatusFailed
				job.Error = err.Error()
			}
		} else {
			job.Status = JobStatusCompleted
		}
	}()

	result, _ := json.Marshal(map[string]string{
		"job_id": jobID,
		"status": string(JobStatusRunning),
	})

	return NewTextResponse(string(result)), nil
}

func (j *jobManagerTool) handleStatus(params JobManagerParams) (ToolResponse, error) {
	if params.JobID == "" {
		return NewTextErrorResponse("Job ID를 지정해주세요"), nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	job, exists := j.jobs[params.JobID]
	if !exists {
		return NewTextErrorResponse(fmt.Sprintf("Job을 찾을 수 없습니다: %s", params.JobID)), nil
	}

	elapsed := time.Since(job.StartTime).Round(time.Second)
	result, _ := json.Marshal(map[string]string{
		"job_id":       job.ID,
		"status":       string(job.Status),
		"elapsed_time": elapsed.String(),
		"error":        job.Error,
	})

	return NewTextResponse(string(result)), nil
}

func (j *jobManagerTool) handleList() (ToolResponse, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	type jobSummary struct {
		JobID   string `json:"job_id"`
		Status  string `json:"status"`
		Elapsed string `json:"elapsed_time"`
	}

	summaries := []jobSummary{}
	for _, job := range j.jobs {
		elapsed := time.Since(job.StartTime).Round(time.Second)
		if job.Status == JobStatusCompleted || job.Status == JobStatusFailed {
			elapsed = job.EndTime.Sub(job.StartTime).Round(time.Second)
		}
		summaries = append(summaries, jobSummary{
			JobID:   job.ID,
			Status:  string(job.Status),
			Elapsed: elapsed.String(),
		})
	}

	result, _ := json.MarshalIndent(summaries, "", "  ")
	return NewTextResponse(fmt.Sprintf("총 %d개의 Job이 등록되어 있습니다:\n%s", len(summaries), string(result))), nil
}

func (j *jobManagerTool) handleCancel(params JobManagerParams) (ToolResponse, error) {
	if params.JobID == "" {
		return NewTextErrorResponse("Job ID를 지정해주세요"), nil
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	job, exists := j.jobs[params.JobID]
	if !exists {
		return NewTextErrorResponse(fmt.Sprintf("Job을 찾을 수 없습니다: %s", params.JobID)), nil
	}

	if job.Status != JobStatusRunning {
		return NewTextErrorResponse(fmt.Sprintf("실행 중인 Job만 취소할 수 있습니다 (현재 상태: %s)", job.Status)), nil
	}

	job.cancel()
	return NewTextResponse(fmt.Sprintf("Job %s 취소가 요청되었습니다", params.JobID)), nil
}
