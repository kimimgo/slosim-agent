package agent

// pipeline.go documents the E2E sloshing simulation pipeline (PIPE-01).
//
// The pipeline is NOT explicitly orchestrated by code, but rather emerges
// from the LLM agent's autonomous tool usage guided by the sloshing system prompt.
//
// The sloshing_coder system prompt (internal/llm/prompt/sloshing_coder.go)
// instructs the agent to follow this sequence:
//
// 1. [User Input] Natural language request (e.g., "1m tank sloshing simulation")
//
// 2. [Inference] Agent interprets conditions:
//    - Tank dimensions (L × W × H)
//    - Fluid height
//    - Excitation frequency & amplitude
//    - Particle spacing (dp)
//    - Simulation time (TimeMax)
//
// 3. [Confirmation] Agent asks user: "Proceed with these conditions?"
//
// 4. [Tool: xml_generator] Generate DualSPHysics XML case file
//    → simulations/{timestamp}_{case_name}/case_Def.xml
//
// 5. [Tool: gencase] Pre-processing (particle geometry generation)
//    → .bi4 binary particle data
//
// 6. [Tool: solver] Launch GPU SPH simulation (background)
//    → Returns job_id, status=RUNNING
//
// 7. [Tool: job_manager] Monitor progress
//    → Poll until status=COMPLETED or FAILED
//
// 8. [Tool: partvtk] Convert results to VTK format
//    → vtk/*.vtk files for visualization
//
// 9. [Tool: measuretool] Extract time-series data (water level, pressure)
//    → csv/*.csv measurement data
//
// 10. [Tool: report] Generate Markdown report
//     → report.md with conditions summary + results
//
// 11. [User Notification] "Simulation complete. Report: simulations/{id}/report.md"
//
// ─────────────────────────────────────────────────────────────────────────────
//
// Implementation Notes:
//
// - The agent (internal/llm/agent/agent.go) automatically handles tool calls
//   via streamAndHandleEvents() → tool.Run() loop
//
// - Each tool returns ToolResponse, which the agent feeds back to the LLM
//
// - The LLM decides the next tool based on:
//   1. System prompt guidelines (sloshing_coder.go)
//   2. Previous tool results
//   3. User input
//
// - No explicit orchestrator code is needed; the prompt is the orchestrator
//
// - Error handling: Each tool returns IsError=true on failure, and the agent
//   propagates this to the LLM for corrective action
//
// - Background execution: solver tool launches a goroutine, returns immediately;
//   job_manager tool tracks completion asynchronously
//
// ─────────────────────────────────────────────────────────────────────────────
//
// v0.1 Constraints:
// - Single job limit (v0.1 does not support concurrent simulations)
// - Rectangular tank only
// - Sinusoidal excitation only
// - Water (1000 kg/m³) single-phase fluid
// - 3D simulation only
//
// Future (v0.2+):
// - Multi-job queue with priority scheduling
// - Real-time progress monitoring with TUI dashboard
// - Parametric studies (automated parameter sweeps)
// - pv-agent MCP server for image/animation rendering (replaced pvpython.go)
//
// ─────────────────────────────────────────────────────────────────────────────
