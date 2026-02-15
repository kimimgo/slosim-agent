//go:build e2e

package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// E2E Prompt Quality Tests for Qwen3 SloshingCoder
//
// Tests whether Qwen3:32b generates correct tool calls given the sloshing
// domain system prompt and tool definitions.
//
// Run with: go test -tags e2e -timeout 600s ./internal/llm/tools/ -run TestPromptQuality -v
//
// Prerequisites:
//   - Ollama running at localhost:11434 with qwen3:32b model
// =============================================================================

const (
	ollamaEndpoint = "http://localhost:11434/v1/chat/completions"
	ollamaModel    = "qwen3:32b"
)

// --- OpenAI-compatible API types ---

type chatRequest struct {
	Model     string        `json:"model"`
	Messages  []chatMessage `json:"messages"`
	Tools     []chatTool    `json:"tools,omitempty"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatTool struct {
	Type     string       `json:"type"`
	Function chatFunction `json:"function"`
}

type chatFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type chatChoice struct {
	Message      chatChoiceMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type chatChoiceMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []toolCall `json:"tool_calls,omitempty"`
}

type toolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// --- System prompt (from sloshing_coder.go) ---

const sloshingTestPrompt = `당신은 슬로싱(Sloshing) 해석 전문 AI 어시스턴트입니다.
비전문가도 자연어로 시뮬레이션을 요청할 수 있도록 돕습니다.

# 절대 규칙
1. 해석 요청 시 반드시 xml_generator를 첫 번째로 호출하세요
2. 기존 XML 파일이 있으면 gencase부터 시작합니다
3. 누락된 파라미터를 사용자에게 묻지 말고 아래 규칙으로 자동 채워서 tool을 호출하세요
4. error_recovery는 시뮬레이션 실행 중 에러가 발생한 경우에만 사용합니다

# Tool 호출 순서 (반드시 이 순서를 따르세요)
1. xml_generator → XML 케이스 파일 생성 (첫 번째로 호출)
2. gencase → 해석 준비 (파티클 생성)
3. solver → 시뮬레이션 실행
4. partvtk → 결과 변환
5. measuretool → 수위/압력 측정
6. report → 리포트 생성

# 파라미터 자동 결정 규칙 (누락 시 이 값을 사용)
1. dp = min(L,W,H)/50 (최소 0.005m, 최대 0.05m)
2. 시뮬레이션 시간(time_max) = 5/freq (초)
3. 유체 높이 미지정 시: 탱크 높이의 50%
4. 진폭 미지정 시: 탱크 길이의 5%
5. out_path 미지정 시: simulations/sloshing_case`

// --- DualSPHysics tool definitions (subset: core pipeline only) ---

func sloshingToolDefs() []chatTool {
	return []chatTool{
		{
			Type: "function",
			Function: chatFunction{
				Name:        "xml_generator",
				Description: "DualSPHysics XML 케이스 자동 생성 — 탱크 치수, 유체 높이, 주파수 등 해석 조건을 XML 케이스 파일로 변환합니다. 시뮬레이션을 시작하려면 이 도구를 먼저 호출하세요.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"tank_length":  map[string]any{"type": "number", "description": "탱크 길이 (m)"},
						"tank_width":   map[string]any{"type": "number", "description": "탱크 너비 (m)"},
						"tank_height":  map[string]any{"type": "number", "description": "탱크 높이 (m)"},
						"fluid_height": map[string]any{"type": "number", "description": "유체 높이 (m)"},
						"freq":         map[string]any{"type": "number", "description": "가진 주파수 (Hz)"},
						"amplitude":    map[string]any{"type": "number", "description": "가진 진폭 (m)"},
						"dp":           map[string]any{"type": "number", "description": "파티클 간격 (m)"},
						"time_max":     map[string]any{"type": "number", "description": "시뮬레이션 시간 (s)"},
						"out_path":     map[string]any{"type": "string", "description": "출력 경로"},
					},
					"required": []string{"tank_length", "tank_width", "tank_height", "fluid_height", "freq", "amplitude", "dp", "time_max", "out_path"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "gencase",
				Description: "DualSPHysics GenCase — XML 케이스 파일에서 파티클 지오메트리를 생성합니다. xml_generator 이후에 호출하세요.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"case_path": map[string]any{"type": "string", "description": "XML 케이스 파일 경로"},
						"save_path": map[string]any{"type": "string", "description": "출력 디렉토리 경로"},
					},
					"required": []string{"case_path", "save_path"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "solver",
				Description: "DualSPHysics GPU 솔버 — 시뮬레이션을 백그라운드에서 실행합니다. gencase 이후에 호출하세요.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"case_name": map[string]any{"type": "string", "description": "케이스 이름"},
						"data_dir":  map[string]any{"type": "string", "description": "데이터 디렉토리"},
						"out_dir":   map[string]any{"type": "string", "description": "출력 디렉토리"},
						"gpu":       map[string]any{"type": "boolean", "description": "GPU 사용 여부"},
					},
					"required": []string{"case_name", "data_dir", "out_dir"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "partvtk",
				Description: "DualSPHysics PartVTK — 시뮬레이션 결과를 VTK 형식으로 변환합니다.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"data_dir": map[string]any{"type": "string", "description": "입력 데이터 디렉토리"},
						"save_vtk": map[string]any{"type": "string", "description": "VTK 출력 경로"},
					},
					"required": []string{"data_dir", "save_vtk"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "measuretool",
				Description: "DualSPHysics MeasureTool — 특정 위치에서 수위, 압력, 속도 데이터를 추출합니다.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"data_dir":    map[string]any{"type": "string", "description": "입력 데이터 디렉토리"},
						"save_csv":    map[string]any{"type": "string", "description": "CSV 출력 경로"},
						"points_file": map[string]any{"type": "string", "description": "측정 포인트 파일 경로"},
					},
					"required": []string{"data_dir", "save_csv", "points_file"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "report",
				Description: "슬로싱 해석 리포트 생성 — 시뮬레이션 결과를 Markdown 리포트로 정리합니다.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"sim_dir": map[string]any{"type": "string", "description": "시뮬레이션 디렉토리"},
					},
					"required": []string{"sim_dir"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "error_recovery",
				Description: "에러 복구 도구 — 시뮬레이션 실행 중 에러가 발생한 경우에만 사용합니다. 새 시뮬레이션 시작 시에는 사용하지 마세요.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"job_id":     map[string]any{"type": "string", "description": "Job ID"},
						"output_dir": map[string]any{"type": "string", "description": "출력 디렉토리"},
					},
					"required": []string{"output_dir"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "monitor",
				Description: "시뮬레이션 모니터링 — Run.csv를 파싱하여 진행률과 안정성을 확인합니다. solver 실행 후 사용합니다.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"sim_dir": map[string]any{"type": "string", "description": "시뮬레이션 디렉토리"},
					},
					"required": []string{"sim_dir"},
				},
			},
		},
	}
}

// --- Scenario definitions ---

type promptScenario struct {
	Name            string   // Test case name
	UserPrompt      string   // User's natural language input
	ExpectedTools   []string // First tool call should be one of these
	ForbiddenTools  []string // These tools should NOT be called first
	ValidateParams  bool     // Whether to validate tool call parameters
	ExpectedParamKV map[string]any
}

func testScenarios() []promptScenario {
	return []promptScenario{
		{
			Name:           "Act1_BasicSloshing",
			UserPrompt:     "1m x 0.5m x 0.6m 탱크에서 0.5Hz 1cm 진폭 슬로싱 해석해줘",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery", "monitor", "solver", "gencase"},
			ValidateParams: true,
			ExpectedParamKV: map[string]any{
				"tank_length": 1.0,
				"tank_width":  0.5,
				"tank_height": 0.6,
				"freq":        0.5,
				"amplitude":   0.01,
			},
		},
		{
			Name:           "Act1_LNGTank",
			UserPrompt:     "LNG 화물창 슬로싱 해석 좀 해줘. 탱크 크기는 가로 40m 세로 30m 높이 27m이고 물 60% 채워져 있어. 배가 좌우로 0.12Hz로 흔들려.",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery", "monitor"},
			ValidateParams: true,
			ExpectedParamKV: map[string]any{
				"tank_length": 40.0,
				"tank_width":  30.0,
				"tank_height": 27.0,
				"freq":        0.12,
			},
		},
		{
			Name:           "Act1_SmallTank_NoFreq",
			UserPrompt:     "소형 탱크로 슬로싱 해석해줘. 주파수는 1Hz.",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery", "solver"},
			ValidateParams: true,
			ExpectedParamKV: map[string]any{
				"freq": 1.0,
			},
		},
		{
			Name:           "Act3_RunExisting",
			UserPrompt:     "cases/SloshingTank_Def.xml 파일로 시뮬레이션 실행해줘",
			ExpectedTools:  []string{"gencase", "xml_generator"},
			ForbiddenTools: []string{"error_recovery"},
		},
		{
			Name:           "Act7_Parametric",
			UserPrompt:     "충진율 30%, 50%, 70%, 90%로 비교해줘. 탱크 크기는 1m x 0.5m x 0.6m, 주파수 0.5Hz, 진폭 1cm",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery"},
		},
		// AGENT-E2E: v0.3 추가 시나리오
		{
			Name:           "MDBC01_BoundaryMethod",
			UserPrompt:     "1m x 0.5m x 0.6m 탱크에서 mDBC 경계조건으로 슬로싱 해석해줘. 주파수 0.5Hz, 진폭 1cm",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery", "monitor", "solver"},
		},
		{
			Name:           "GEO02_LShapedTank",
			UserPrompt:     "L자 형태 탱크로 슬로싱 해석해줘. 주파수 0.3Hz",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery"},
		},
		{
			Name:           "EXC01_SeismicInput",
			UserPrompt:     "지진파 데이터 파일 earthquake.csv 로 슬로싱 해석해줘. 탱크 크기 2m x 1m x 1.5m",
			ExpectedTools:  []string{"xml_generator"},
			ForbiddenTools: []string{"error_recovery", "monitor"},
		},
		{
			Name:           "ErrorRecovery_NotCalledFirst",
			UserPrompt:     "해석이 불안정해졌어. 타임스텝 줄여서 다시 해봐",
			ExpectedTools:  []string{"error_recovery", "monitor"},
			ForbiddenTools: []string{"xml_generator"},
		},
	}
}

// --- Helper: call Ollama API ---

func callOllama(t *testing.T, systemPrompt, userPrompt string, tools []chatTool) (*chatResponse, time.Duration) {
	t.Helper()

	// Append /no_think suffix to suppress Qwen3 extended thinking
	promptWithNoThink := userPrompt + " /no_think"

	req := chatRequest{
		Model: ollamaModel,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: promptWithNoThink},
		},
		Tools:     tools,
		MaxTokens: 1024,
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	// Use HTTP client with 3-minute timeout (Qwen3:32b with CPU offload is slow)
	client := &http.Client{Timeout: 3 * time.Minute}

	start := time.Now()
	resp, err := client.Post(ollamaEndpoint, "application/json", bytes.NewReader(body))
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Ollama request failed after %s: %v", elapsed, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read Ollama response")

	if resp.StatusCode != 200 {
		t.Fatalf("Ollama returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	err = json.Unmarshal(respBody, &chatResp)
	require.NoError(t, err, "Failed to parse Ollama response: %s", string(respBody))

	if chatResp.Error != nil {
		t.Fatalf("Ollama API error: %s", chatResp.Error.Message)
	}

	return &chatResp, elapsed
}

// --- Test: Individual scenario prompt quality ---

func TestPromptQuality_Scenarios(t *testing.T) {
	skipIfNoOllama(t)

	tools := sloshingToolDefs()
	scenarios := testScenarios()

	var passCount, failCount int

	for _, sc := range scenarios {
		t.Run(sc.Name, func(t *testing.T) {
			t.Logf("Prompt: %s", sc.UserPrompt)

			resp, elapsed := callOllama(t, sloshingTestPrompt, sc.UserPrompt, tools)
			t.Logf("Response time: %s", elapsed)

			require.NotEmpty(t, resp.Choices, "No choices in response")

			choice := resp.Choices[0]

			// Log the full response for debugging
			if choice.Message.Content != "" {
				t.Logf("Text response: %s", truncateE2EOutput(choice.Message.Content, 500))
			}

			// Check if model made tool calls
			if len(choice.Message.ToolCalls) == 0 {
				t.Logf("WARN: No tool calls — model responded with text only")
				t.Logf("Finish reason: %s", choice.FinishReason)

				// Check if the text response at least mentions the expected tool
				for _, expected := range sc.ExpectedTools {
					if strings.Contains(strings.ToLower(choice.Message.Content), expected) {
						t.Logf("Text mentions '%s' — model may intend to call it", expected)
					}
				}
				failCount++
				t.Errorf("Expected tool call to one of %v, got text-only response", sc.ExpectedTools)
				return
			}

			// Analyze first tool call
			firstTool := choice.Message.ToolCalls[0].Function.Name
			t.Logf("First tool call: %s", firstTool)
			t.Logf("All tool calls: %v", toolCallNames(choice.Message.ToolCalls))

			// Check forbidden tools
			for _, forbidden := range sc.ForbiddenTools {
				if firstTool == forbidden {
					failCount++
					t.Errorf("FORBIDDEN: First tool call was '%s' (should not be called first)", forbidden)
					return
				}
			}

			// Check expected tools
			found := false
			for _, expected := range sc.ExpectedTools {
				if firstTool == expected {
					found = true
					break
				}
			}
			if !found {
				failCount++
				t.Errorf("Expected first tool to be one of %v, got '%s'", sc.ExpectedTools, firstTool)
				return
			}

			t.Logf("PASS: Correct first tool call '%s'", firstTool)
			passCount++

			// Validate parameters if requested
			if sc.ValidateParams && firstTool == "xml_generator" {
				args := choice.Message.ToolCalls[0].Function.Arguments
				var params map[string]any
				err := json.Unmarshal([]byte(args), &params)
				if err != nil {
					t.Logf("WARN: Could not parse tool arguments: %s", err)
					return
				}

				t.Logf("Tool arguments: %s", args)

				for key, expectedVal := range sc.ExpectedParamKV {
					actualVal, ok := params[key]
					if !ok {
						t.Logf("WARN: Missing parameter '%s' (expected %v)", key, expectedVal)
						continue
					}
					// Compare as float64 (JSON numbers are float64)
					expectedFloat, eOk := toFloat64(expectedVal)
					actualFloat, aOk := toFloat64(actualVal)
					if eOk && aOk {
						if !floatClose(expectedFloat, actualFloat, 0.1) {
							t.Logf("WARN: Parameter '%s' = %v (expected ~%v)", key, actualVal, expectedVal)
						} else {
							t.Logf("OK: Parameter '%s' = %v", key, actualVal)
						}
					}
				}
			}
		})
	}

	t.Logf("\n=== Prompt Quality Summary ===")
	t.Logf("Pass: %d / %d", passCount, passCount+failCount)
	t.Logf("Fail: %d / %d", failCount, passCount+failCount)
	t.Logf("==============================")
}

// --- Test: All tools (full tool set like production) ---

func TestPromptQuality_FullToolSet(t *testing.T) {
	skipIfNoOllama(t)

	// Test with ALL tools including built-in ones (view, bash, edit, etc.)
	// to reproduce the actual production environment
	allTools := append(sloshingToolDefs(), builtinToolDefs()...)

	prompt := "1m x 0.5m x 0.6m 탱크에서 0.5Hz 1cm 진폭 슬로싱 해석해줘"

	t.Logf("Testing with %d tools (production-like)", len(allTools))
	t.Logf("Prompt: %s", prompt)

	resp, elapsed := callOllama(t, sloshingTestPrompt, prompt, allTools)
	t.Logf("Response time: %s", elapsed)

	require.NotEmpty(t, resp.Choices)
	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) == 0 {
		t.Logf("Text response: %s", truncateE2EOutput(choice.Message.Content, 500))
		t.Error("No tool calls — model responded with text only (full tool set)")
		return
	}

	firstTool := choice.Message.ToolCalls[0].Function.Name
	t.Logf("First tool call: %s", firstTool)
	t.Logf("All tool calls: %v", toolCallNames(choice.Message.ToolCalls))

	assert.Equal(t, "xml_generator", firstTool,
		"With full tool set, first tool should be xml_generator")
}

// --- Test: Minimal prompt (diagnose if system prompt is the issue) ---

func TestPromptQuality_MinimalPrompt(t *testing.T) {
	skipIfNoOllama(t)

	// Ultra-minimal system prompt to test if Qwen3 can do tool calls at all
	minimalPrompt := `You are a sloshing simulation assistant.
When the user asks for a simulation, ALWAYS call xml_generator first.
Never call error_recovery unless there is an actual error.`

	tools := sloshingToolDefs()
	prompt := "1m x 0.5m x 0.6m 탱크에서 0.5Hz 1cm 진폭 슬로싱 해석해줘"

	t.Logf("Testing with minimal system prompt (%d bytes)", len(minimalPrompt))

	resp, elapsed := callOllama(t, minimalPrompt, prompt, tools)
	t.Logf("Response time: %s", elapsed)

	require.NotEmpty(t, resp.Choices)
	choice := resp.Choices[0]

	if len(choice.Message.ToolCalls) == 0 {
		t.Logf("Text response: %s", truncateE2EOutput(choice.Message.Content, 500))
		t.Error("No tool calls even with minimal prompt")
		return
	}

	firstTool := choice.Message.ToolCalls[0].Function.Name
	t.Logf("First tool call: %s (minimal prompt)", firstTool)

	assert.Equal(t, "xml_generator", firstTool,
		"Even with minimal prompt, first tool should be xml_generator")
}

// --- Built-in tool definitions (to simulate production environment) ---

func builtinToolDefs() []chatTool {
	return []chatTool{
		{
			Type: "function",
			Function: chatFunction{
				Name:        "bash",
				Description: "Execute a bash command",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"command": map[string]any{"type": "string", "description": "The bash command to execute"},
					},
					"required": []string{"command"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "view",
				Description: "View file contents with line numbers",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path": map[string]any{"type": "string", "description": "File path to read"},
					},
					"required": []string{"file_path"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "edit",
				Description: "Edit a file by replacing text",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path":  map[string]any{"type": "string"},
						"old_string": map[string]any{"type": "string"},
						"new_string": map[string]any{"type": "string"},
					},
					"required": []string{"file_path", "old_string", "new_string"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "glob",
				Description: "Find files matching a pattern",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"pattern": map[string]any{"type": "string"},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "grep",
				Description: "Search file contents for a pattern",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"pattern": map[string]any{"type": "string"},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "write",
				Description: "Write content to a file",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path": map[string]any{"type": "string"},
						"content":   map[string]any{"type": "string"},
					},
					"required": []string{"file_path", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: chatFunction{
				Name:        "ls",
				Description: "List directory contents",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path": map[string]any{"type": "string"},
					},
					"required": []string{"path"},
				},
			},
		},
	}
}

// --- Helper functions ---

func toolCallNames(calls []toolCall) []string {
	names := make([]string, len(calls))
	for i, c := range calls {
		names[i] = c.Function.Name
	}
	return names
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func floatClose(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// checkOllamaModel verifies a specific model is available
func checkOllamaModel(t *testing.T, model string) {
	t.Helper()
	resp, err := http.Get("http://localhost:11434/v1/models")
	if err != nil {
		t.Skipf("Ollama not reachable: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), model) {
		// Try without tag
		base := strings.Split(model, ":")[0]
		if !strings.Contains(string(body), base) {
			t.Skipf("Model %s not available in Ollama", model)
		}
	}
}

func init() {
	// Ensure OLLAMA_HOST can be overridden via env
	if host := os.Getenv("OLLAMA_HOST"); host != "" {
		// Note: ollamaEndpoint is const, would need to make it var for this
		_ = host
	}
}
