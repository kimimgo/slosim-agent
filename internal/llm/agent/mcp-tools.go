package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/permission"
	"github.com/opencode-ai/opencode/internal/version"
)

// mcpSessionPool manages persistent MCP server sessions.
// Each server gets one long-lived ClientSession, reused across tool calls.
type mcpSessionPool struct {
	mu       sync.RWMutex
	client   *mcp.Client
	sessions map[string]*mcp.ClientSession
}

func (p *mcpSessionPool) getSession(ctx context.Context, name string, cfg config.MCPServer) (*mcp.ClientSession, error) {
	p.mu.RLock()
	if s, ok := p.sessions[name]; ok {
		p.mu.RUnlock()
		return s, nil
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	// Double-check after acquiring write lock
	if s, ok := p.sessions[name]; ok {
		return s, nil
	}

	cmd := exec.Command(cfg.Command, cfg.Args...)
	cmd.Env = cfg.EnvSlice()
	transport := &mcp.CommandTransport{Command: cmd}
	session, err := p.client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("MCP connect %q: %w", name, err)
	}
	p.sessions[name] = session
	logging.Info("MCP session established", "server", name)
	return session, nil
}

func (p *mcpSessionPool) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, s := range p.sessions {
		if err := s.Close(); err != nil {
			logging.Error("MCP session close error", "server", name, "error", err)
		}
	}
	p.sessions = make(map[string]*mcp.ClientSession)
}

// pool is the package-level singleton for MCP session management.
var pool = &mcpSessionPool{
	client: mcp.NewClient(&mcp.Implementation{
		Name:    "slosim-agent",
		Version: version.Version,
	}, nil),
	sessions: make(map[string]*mcp.ClientSession),
}

// CloseAllMCPSessions gracefully closes all MCP server sessions.
// Call this during application shutdown.
func CloseAllMCPSessions() {
	pool.closeAll()
}

// mcpTool wraps an MCP server tool as a BaseTool.
type mcpTool struct {
	mcpName     string
	tool        *mcp.Tool
	mcpConfig   config.MCPServer
	permissions permission.Service
}

func (b *mcpTool) Info() tools.ToolInfo {
	properties, required := extractInputSchema(b.tool.InputSchema)
	return tools.ToolInfo{
		Name:        fmt.Sprintf("%s_%s", b.mcpName, b.tool.Name),
		Description: b.tool.Description,
		Parameters:  properties,
		Required:    required,
	}
}

func (b *mcpTool) Run(ctx context.Context, params tools.ToolCall) (tools.ToolResponse, error) {
	sessionID, messageID := tools.GetContextValues(ctx)
	if sessionID == "" || messageID == "" {
		return tools.ToolResponse{}, fmt.Errorf("session ID and message ID are required")
	}
	permissionDescription := fmt.Sprintf("execute %s with the following parameters: %s", b.Info().Name, params.Input)
	p := b.permissions.Request(
		permission.CreatePermissionRequest{
			SessionID:   sessionID,
			Path:        config.WorkingDirectory(),
			ToolName:    b.Info().Name,
			Action:      "execute",
			Description: permissionDescription,
			Params:      params.Input,
		},
	)
	if !p {
		return tools.NewTextErrorResponse("permission denied"), nil
	}

	session, err := pool.getSession(ctx, b.mcpName, b.mcpConfig)
	if err != nil {
		return tools.NewTextErrorResponse(fmt.Sprintf("MCP session error: %s", err)), nil
	}
	return runTool(ctx, session, b.tool.Name, params.Input)
}

func runTool(ctx context.Context, session *mcp.ClientSession, toolName string, input string) (tools.ToolResponse, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return tools.NewTextErrorResponse(fmt.Sprintf("error parsing parameters: %s", err)), nil
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return tools.NewTextErrorResponse(err.Error()), nil
	}
	if result.IsError {
		var errMsg string
		for _, c := range result.Content {
			if tc, ok := c.(*mcp.TextContent); ok {
				errMsg += tc.Text
			}
		}
		return tools.NewTextErrorResponse(errMsg), nil
	}

	var textParts []string
	for _, c := range result.Content {
		switch v := c.(type) {
		case *mcp.TextContent:
			textParts = append(textParts, v.Text)
		case *mcp.ImageContent:
			textParts = append(textParts, fmt.Sprintf("[image: %s, %d bytes]", v.MIMEType, len(v.Data)))
		default:
			textParts = append(textParts, fmt.Sprintf("%v", v))
		}
	}
	return tools.NewTextResponse(strings.Join(textParts, "\n")), nil
}

// NewMcpTool creates a BaseTool wrapping an MCP server tool.
func NewMcpTool(name string, tool *mcp.Tool, permissions permission.Service, mcpConfig config.MCPServer) tools.BaseTool {
	return &mcpTool{
		mcpName:     name,
		tool:        tool,
		mcpConfig:   mcpConfig,
		permissions: permissions,
	}
}

var (
	mcpToolsMu sync.RWMutex
	mcpTools   []tools.BaseTool
)

// GetMcpTools returns all tools from configured MCP servers.
// Sessions are established on first call and reused for subsequent tool invocations.
func GetMcpTools(ctx context.Context, permissions permission.Service) []tools.BaseTool {
	mcpToolsMu.RLock()
	if len(mcpTools) > 0 {
		mcpToolsMu.RUnlock()
		return mcpTools
	}
	mcpToolsMu.RUnlock()

	mcpToolsMu.Lock()
	defer mcpToolsMu.Unlock()
	if len(mcpTools) > 0 {
		return mcpTools
	}

	for name, m := range config.Get().MCPServers {
		if m.Type != "" && m.Type != config.MCPStdio {
			logging.Warn("MCP server type not supported, skipping", "server", name, "type", m.Type)
			continue
		}
		session, err := pool.getSession(ctx, name, m)
		if err != nil {
			logging.Error("MCP connect failed", "server", name, "error", err)
			continue
		}
		result, err := session.ListTools(ctx, &mcp.ListToolsParams{})
		if err != nil {
			logging.Error("MCP ListTools failed", "server", name, "error", err)
			continue
		}
		for _, t := range result.Tools {
			mcpTools = append(mcpTools, NewMcpTool(name, t, permissions, m))
		}
		logging.Info("MCP tools loaded", "server", name, "count", len(result.Tools))
	}
	return mcpTools
}

// extractInputSchema converts the go-sdk Tool.InputSchema (any) to
// Properties map and Required slice for ToolInfo.
func extractInputSchema(schema any) (map[string]any, []string) {
	m, ok := schema.(map[string]any)
	if !ok {
		return make(map[string]any), nil
	}
	props, _ := m["properties"].(map[string]any)
	if props == nil {
		props = make(map[string]any)
	}
	var required []string
	if req, ok := m["required"].([]any); ok {
		for _, r := range req {
			if s, ok := r.(string); ok {
				required = append(required, s)
			}
		}
	}
	return props, required
}
