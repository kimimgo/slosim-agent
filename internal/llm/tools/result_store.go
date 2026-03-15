package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ResultStoreParams defines parameters for result store operations (STORE-01).
type ResultStoreParams struct {
	Action     string                 `json:"action"`
	ResultID   string                 `json:"result_id,omitempty"`
	ResultData *SimulationResult      `json:"result_data,omitempty"`
	Query      *ResultQuery           `json:"query,omitempty"`
}

// SimulationResult holds metadata and paths for a simulation result.
type SimulationResult struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Timestamp   string                 `json:"timestamp"`
	CaseFile    string                 `json:"case_file"`
	OutputDir   string                 `json:"output_dir"`
	Parameters  map[string]any         `json:"parameters,omitempty"`
	Status      string                 `json:"status"`
	Duration    float64                `json:"duration,omitempty"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// ResultQuery defines search criteria for results.
type ResultQuery struct {
	Name       string   `json:"name,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	FromDate   string   `json:"from_date,omitempty"`
	ToDate     string   `json:"to_date,omitempty"`
	Status     string   `json:"status,omitempty"`
	Limit      int      `json:"limit,omitempty"`
}

type resultStoreTool struct {
	storeDir string
}

const (
	ResultStoreToolName    = "result_store"
	resultStoreDescription = `시뮬레이션 결과 저장소 도구 (STORE-01) — 과거 시뮬레이션 결과를 관리합니다.

사용 가능한 액션:
- save: 새 결과 저장 (result_data 필수)
- get: 특정 결과 조회 (result_id 필수)
- list: 모든 결과 목록 조회
- search: 조건에 맞는 결과 검색 (query 필수)
- delete: 결과 삭제 (result_id 필수)
- compare: 여러 결과 비교 (result_id에 쉼표로 구분)

하이브리드 구조:
- SQLite 메타데이터 (ID, 이름, 타임스탬프, 파라미터, 상태)
- 파일시스템 원본 데이터 (VTK, CSV, 리포트)

예시:
{
  "action": "search",
  "query": {"tags": ["sloshing", "benchmark"], "limit": 10}
}`
)

func NewResultStoreTool() BaseTool {
	return &resultStoreTool{
		storeDir: filepath.Join(getWorkingDirectory(), "result_store"),
	}
}

// NewResultStoreToolWithDir creates a result store with a custom directory (for testing).
func NewResultStoreToolWithDir(storeDir string) BaseTool {
	return &resultStoreTool{
		storeDir: storeDir,
	}
}

func (r *resultStoreTool) Info() ToolInfo {
	return ToolInfo{
		Name:        ResultStoreToolName,
		Description: resultStoreDescription,
		Parameters: map[string]any{
			"action": map[string]any{
				"type":        "string",
				"description": "수행할 액션 (save, get, list, search, delete, compare)",
			},
			"result_id": map[string]any{
				"type":        "string",
				"description": "결과 ID (get, delete 시 필수)",
			},
			"result_data": map[string]any{
				"type":        "object",
				"description": "저장할 결과 데이터 (save 시 필수)",
			},
			"query": map[string]any{
				"type":        "object",
				"description": "검색 조건 (search 시 필수)",
			},
		},
		Required: []string{"action"},
	}
}

func (r *resultStoreTool) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var params ResultStoreParams
	if err := UnmarshalToolInput(call.Input, &params); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("파라미터 파싱 오류: %s", err)), nil
	}

	// Ensure store directory exists
	if err := os.MkdirAll(r.storeDir, 0755); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("저장소 디렉토리 생성 실패: %s", err)), nil
	}

	switch params.Action {
	case "save":
		return r.handleSave(params)
	case "get":
		return r.handleGet(params)
	case "list":
		return r.handleList()
	case "search":
		return r.handleSearch(params)
	case "delete":
		return r.handleDelete(params)
	case "compare":
		return r.handleCompare(params)
	default:
		return NewTextErrorResponse(fmt.Sprintf("알 수 없는 액션입니다: %s", params.Action)), nil
	}
}

func (r *resultStoreTool) handleSave(params ResultStoreParams) (ToolResponse, error) {
	if params.ResultData == nil {
		return NewTextErrorResponse("result_data를 지정해주세요"), nil
	}

	result := params.ResultData

	// Generate ID if not provided
	if result.ID == "" {
		result.ID = fmt.Sprintf("sim_%d", time.Now().Unix())
	}

	// Set timestamp if not provided
	if result.Timestamp == "" {
		result.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Save metadata to JSON file
	metadataPath := filepath.Join(r.storeDir, result.ID+".json")
	metadataJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("메타데이터 직렬화 실패: %s", err)), nil
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("메타데이터 저장 실패: %s", err)), nil
	}

	resultMsg := fmt.Sprintf(`결과 저장 완료:
- ID: %s
- 이름: %s
- 출력 디렉토리: %s
- 상태: %s
- 메타데이터 경로: %s
`, result.ID, result.Name, result.OutputDir, result.Status, metadataPath)

	return NewTextResponse(resultMsg), nil
}

func (r *resultStoreTool) handleGet(params ResultStoreParams) (ToolResponse, error) {
	if params.ResultID == "" {
		return NewTextErrorResponse("result_id를 지정해주세요"), nil
	}

	metadataPath := filepath.Join(r.storeDir, params.ResultID+".json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("결과를 찾을 수 없습니다: %s", params.ResultID)), nil
	}

	var result SimulationResult
	if err := json.Unmarshal(data, &result); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("메타데이터 파싱 실패: %s", err)), nil
	}

	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return NewTextResponse(string(resultJSON)), nil
}

func (r *resultStoreTool) handleList() (ToolResponse, error) {
	files, err := os.ReadDir(r.storeDir)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("저장소 읽기 실패: %s", err)), nil
	}

	results := []SimulationResult{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(r.storeDir, file.Name()))
			if err != nil {
				continue
			}

			var result SimulationResult
			if err := json.Unmarshal(data, &result); err != nil {
				continue
			}

			results = append(results, result)
		}
	}

	resultJSON, _ := json.MarshalIndent(results, "", "  ")
	return NewTextResponse(fmt.Sprintf("총 %d개의 결과가 저장되어 있습니다:\n%s", len(results), string(resultJSON))), nil
}

func (r *resultStoreTool) handleSearch(params ResultStoreParams) (ToolResponse, error) {
	if params.Query == nil {
		return NewTextErrorResponse("query를 지정해주세요"), nil
	}

	query := params.Query
	if query.Limit == 0 {
		query.Limit = 100
	}

	files, err := os.ReadDir(r.storeDir)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("저장소 읽기 실패: %s", err)), nil
	}

	results := []SimulationResult{}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(r.storeDir, file.Name()))
			if err != nil {
				continue
			}

			var result SimulationResult
			if err := json.Unmarshal(data, &result); err != nil {
				continue
			}

			// Apply filters
			if query.Name != "" && result.Name != query.Name {
				continue
			}
			if query.Status != "" && result.Status != query.Status {
				continue
			}
			if len(query.Tags) > 0 && !containsAnyTag(result.Tags, query.Tags) {
				continue
			}

			results = append(results, result)

			if len(results) >= query.Limit {
				break
			}
		}
	}

	resultJSON, _ := json.MarshalIndent(results, "", "  ")
	return NewTextResponse(fmt.Sprintf("검색 결과 %d개:\n%s", len(results), string(resultJSON))), nil
}

func (r *resultStoreTool) handleDelete(params ResultStoreParams) (ToolResponse, error) {
	if params.ResultID == "" {
		return NewTextErrorResponse("result_id를 지정해주세요"), nil
	}

	metadataPath := filepath.Join(r.storeDir, params.ResultID+".json")
	if err := os.Remove(metadataPath); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("결과 삭제 실패: %s", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("결과 %s 삭제 완료", params.ResultID)), nil
}

func (r *resultStoreTool) handleCompare(params ResultStoreParams) (ToolResponse, error) {
	if params.ResultID == "" {
		return NewTextErrorResponse("비교할 result_id 목록을 쉼표로 구분하여 지정해주세요"), nil
	}

	// Split result IDs by comma
	rawIDs := strings.Split(params.ResultID, ",")
	ids := make([]string, 0, len(rawIDs))
	for _, id := range rawIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			ids = append(ids, trimmed)
		}
	}

	if len(ids) < 2 {
		return NewTextErrorResponse("비교하려면 최소 2개의 result_id가 필요합니다"), nil
	}

	// Load actual result data for comparison
	results := make([]SimulationResult, 0, len(ids))
	for _, id := range ids {
		metadataPath := filepath.Join(r.storeDir, id+".json")
		data, err := os.ReadFile(metadataPath)
		if err != nil {
			return NewTextErrorResponse(fmt.Sprintf("결과를 찾을 수 없습니다: %s", id)), nil
		}

		var result SimulationResult
		if err := json.Unmarshal(data, &result); err != nil {
			return NewTextErrorResponse(fmt.Sprintf("메타데이터 파싱 실패: %s", id)), nil
		}
		results = append(results, result)
	}

	compareMsg := fmt.Sprintf("결과 비교 (PARA-03):\n")
	compareMsg += fmt.Sprintf("비교 대상 ID 수: %d\n\n", len(ids))

	for _, result := range results {
		compareMsg += fmt.Sprintf("- %s (%s): 상태=%s, 소요시간=%.1fs\n",
			result.Name, result.ID, result.Status, result.Duration)
		if len(result.Parameters) > 0 {
			paramsJSON, _ := json.Marshal(result.Parameters)
			compareMsg += fmt.Sprintf("  파라미터: %s\n", string(paramsJSON))
		}
	}

	return NewTextResponse(compareMsg), nil
}

func containsAnyTag(resultTags, queryTags []string) bool {
	for _, qt := range queryTags {
		for _, rt := range resultTags {
			if rt == qt {
				return true
			}
		}
	}
	return false
}
