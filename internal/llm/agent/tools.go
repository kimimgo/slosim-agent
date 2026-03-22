package agent

import (
	"context"

	"github.com/opencode-ai/opencode/internal/history"
	"github.com/opencode-ai/opencode/internal/llm/tools"
	"github.com/opencode-ai/opencode/internal/lsp"
	"github.com/opencode-ai/opencode/internal/message"
	"github.com/opencode-ai/opencode/internal/permission"
	"github.com/opencode-ai/opencode/internal/session"
)

func CoderAgentTools(
	permissions permission.Service,
	sessions session.Service,
	messages message.Service,
	history history.Service,
	lspClients map[string]*lsp.Client,
) []tools.BaseTool {
	ctx := context.Background()
	otherTools := GetMcpTools(ctx, permissions)
	if len(lspClients) > 0 {
		otherTools = append(otherTools, tools.NewDiagnosticsTool(lspClients))
	}
	return append(
		[]tools.BaseTool{
			// DualSPHysics core pipeline (highest priority for sloshing agent)
			tools.NewXMLGeneratorTool(),
			tools.NewSTLImportTool(),
			tools.NewGenCaseTool(),
			tools.NewSolverTool(),
			tools.NewPartVTKTool(),
			tools.NewMeasureToolTool(),
			tools.NewBaffleGeneratorTool(),
			tools.NewReportTool(),
			tools.NewAnalysisTool(),
			tools.NewMonitorTool(),
			tools.NewErrorRecoveryTool(),
			tools.NewJobManagerTool(),
			tools.NewSeismicInputTool(),
			tools.NewParametricStudyTool(),
			tools.NewResultStoreTool(),
			// General-purpose tools
			tools.NewBashTool(permissions),
			tools.NewViewTool(lspClients),
			tools.NewGlobTool(),
			tools.NewGrepTool(),
			tools.NewLsTool(),
			tools.NewWriteTool(lspClients, permissions, history),
			tools.NewEditTool(lspClients, permissions, history),
			tools.NewFetchTool(permissions),
			tools.NewSourcegraphTool(),
			tools.NewPatchTool(lspClients, permissions, history),
			NewAgentTool(sessions, messages, lspClients),
		}, otherTools...,
	)
}

func TaskAgentTools(lspClients map[string]*lsp.Client) []tools.BaseTool {
	return []tools.BaseTool{
		tools.NewGlobTool(),
		tools.NewGrepTool(),
		tools.NewLsTool(),
		tools.NewSourcegraphTool(),
		tools.NewViewTool(lspClients),
	}
}
