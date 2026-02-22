# 3. System Design

This section describes SloshAgent's architecture, tool interface design, domain-specialized prompt, and XML generation pipeline. Figure~\ref{fig:architecture} provides an end-to-end overview.

## 3.1 Architecture Overview

SloshAgent automates the full sloshing simulation lifecycle: a user issues a natural language request (e.g., "simulate a 1 m rectangular tank with 50\% fill at the first natural frequency"); the system converts this into a DualSPHysics XML case, executes the GPU-accelerated SPH solver, post-processes the results, and delivers a physics-informed analysis report. The pipeline comprises four layers.

The **natural language interface** uses a local Qwen3 32B model \cite{qwen3_2025} via Ollama (zero cloud API cost). The model receives the SloshingCoderPrompt (Section 3.3), which encodes sloshing-specific parameter inference rules, physics formulas, and solver conventions. This prompt acts as the sole orchestrator: unlike MetaOpenFOAM's four fixed-role agents \cite{metaopenfoam2024} or Foam-Agent 2.0's six-agent hierarchy \cite{foamagent2025}, SloshAgent uses a single agent with a rich tool set --- an architecture comparable to ChemCrow's 18-tool ReAct loop \cite{bran2024chemcrow} and MDCrow's 40+ tools \cite{mdcrow2025}.

The **ReAct agent loop** \cite{yao2023react} streams the LLM response, collects tool calls, executes them sequentially, and feeds results back until the model produces a final text response with no further tool invocations. This loop is implemented in 120 lines of Go: the `processGeneration` method iterates until the finish reason changes from `tool_use` to `end_turn`.

The **tool layer** exposes 14 built-in DualSPHysics tools (3,866 lines of Go; Section 3.2) plus 12 MCP tools from a pv-agent server for ParaView post-processing. All tools implement a uniform `BaseTool` interface: `Info()` returns the tool schema; `Run(ctx, ToolCall)` returns a `ToolResponse`. The LLM invokes all 26 tools through the same mechanism.

The **solver backend** runs DualSPHysics v5.4 GPU inside a Docker container (CUDA 12.6, NVIDIA RTX 4090). The container mounts `./cases:/cases` for XML input and `./simulations:/data` for output, isolating the CUDA toolchain for reproducibility.

The single-agent design is deliberate. For sloshing's largely sequential pipeline (XML $\rightarrow$ pre-processing $\rightarrow$ solving $\rightarrow$ post-processing $\rightarrow$ analysis), a single agent with domain-aware tool selection achieves comparable automation with lower complexity and token cost than multi-agent alternatives.

## 3.2 Tool Interface Design for Sloshing SPH

Table~\ref{tab:tools} summarizes the 14 built-in tools, each addressing a specific sloshing pain point.

| Tool | Lines | Pain Point Addressed |
|------|-------|----------------------|
| `xml_generator` | 336 | Eliminates hand-crafting XML (200+ page manual) |
| `gencase` | 136 | Automates pre-processing (75--80\% of CFD time \cite{cadence2024}) |
| `solver` | 108 | GPU-accelerated SPH with background launch |
| `job_manager` | 251 | Async lifecycle for hours-long GPU simulations |
| `monitor` | 233 | Run.csv parsing with divergence detection |
| `error_recovery` | 364 | Diagnosis of 5 silent error types |
| `partvtk` | 133 | VTK conversion for visualization |
| `measuretool` | 124 | Pressure and free-surface extraction |
| `analysis` | 218 | AI physics interpretation (resonance, stability) |
| `report` | 241 | Structured Markdown reporting |
| `parametric_study` | 273 | Multi-case sweeps (200+ cases/tank) |
| `stl_import` | 435 | CAD mesh import with watertight validation |
| `seismic_input` | 295 | Earthquake/wave time-series parsing |
| `result_store` | 349 | Persistent result storage and retrieval |

Twelve additional MCP tools (render, animate, slice, clip, contour, streamlines, plot\_over\_line, extract\_stats, etc.) are provided by a pv-agent server wrapping ParaView's `pvpython`, communicating via the Model Context Protocol \cite{mcp2024} over stdio --- analogous to Foam-Agent 2.0's MCP architecture \cite{foamagent2025}.

Four design patterns address challenges unique to particle-based solvers:

**IsError pattern for LLM self-correction.** When a tool detects an error --- GenCase failure, NaN in solver output, or density violations --- it returns `ToolResponse{IsError: true}` with a diagnostic message. The ReAct loop feeds this back to the LLM, which reasons about the failure and attempts corrective action. This implements self-reflection \cite{shinn2023reflexion} at the tool-response level, addressing DualSPHysics's five silent error types without requiring a separate reflection agent.

**Run.csv divergence detection.** SPH instability manifests as exponential kinetic energy growth. The `monitor` and `error_recovery` tools parse the solver's semicolon-delimited Run.csv and flag divergence when the energy ratio exceeds 1.2$\times$ for five consecutive steps --- an SPH-specific heuristic distinct from mesh-based residual convergence.

**Asynchronous GPU execution.** The `solver` tool launches DualSPHysics in a background goroutine via `context.Background()`, returning a job ID immediately. The `job_manager` tracks up to three concurrent jobs with mutex-protected state, supporting submit, status, list, and cancel operations. Cancellation propagates through Go's `context.WithCancel`. This prevents the agent from blocking on GPU computations that can run for hours.

**pv-agent MCP server.** Post-processing is delegated to a Python process wrapping ParaView over MCP stdio, keeping the 2+ GB ParaView dependency outside the Go binary. The server uses Mesa software rendering for headless environments.

## 3.3 Domain-Specialized Sloshing Prompt

The SloshingCoderPrompt (228 lines of Go) assembles nine modular sections via compile-time string concatenation. Four ablation modes --- `full`, `no-domain`, `no-rules`, and `generic` --- are selectable via an environment variable, enabling the controlled ablation experiments of Section 4.5.

The prompt encodes five categories of sloshing knowledge:

**(1) Parameter inference rules.** Deterministic defaults when the user omits values: $\mathit{dp} = \min(L,W,H)/50$ (clamped to $[0.005, 0.05]$ m), $t_{\max} = 5/f$, fluid height = 50\% of tank height, amplitude = 5\% of tank length. These encode the expertise that beginners lack --- overly fine $\mathit{dp}$ causes prohibitive compute; overly coarse values degrade free-surface accuracy.

**(2) Tank presets.** Four geometries (LNG: $40{\times}40{\times}27$ m, ship: $20{\times}10{\times}8$ m, small: $1{\times}0.5{\times}0.6$ m, experimental: $0.6{\times}0.3{\times}0.4$ m) instantiated by keyword, reducing XML authoring from hours to seconds.

**(3) Physics formulas.** The first-mode natural frequency $f_1 = \frac{1}{2\pi}\sqrt{\frac{g\pi}{L}\tanh\!\left(\frac{\pi h}{L}\right)}$ enables autonomous resonance detection and safety warnings.

**(4) Terminology mapping.** A bilingual Korean--English table for seven CFD terms (e.g., "공진 주파수" $\rightarrow$ natural frequency $f_1$), supporting the Korean shipyards that produce ~70\% of LNG carriers \cite{lloyd2024}.

**(5) Docker path conventions.** Explicit mount-point mappings (`/cases/` for XML, `/data/` for output, lowercase binaries) address a recurring source of silent failures on DualSPHysics forums.

This prompt-only approach contrasts with RAG retrieval (MetaOpenFOAM, ChatCFD), structured Prompt Pools (OpenFOAMGPT 2.0), and fine-tuning on 28K+ pairs (AutoCFD, FoamGPT). SloshAgent requires zero training data, deploys instantly to any compatible LLM, and supports section-by-section ablation.

## 3.4 XML Generation Pipeline

The `xml_generator` tool (336 lines of Go) is the entry point of every workflow. It accepts nine parameters (tank dimensions, fluid height, frequency, amplitude, $\mathit{dp}$, simulation time, output path) and produces a complete DualSPHysics XML case plus a probe points file.

The tool encodes four expert-level auto-configuration rules:

**Geometry construction.** Two nested `drawbox` commands define the fluid domain (solid fill) and tank boundary (bottom + four walls). The simulation domain extends by $5{\times}\mathit{dp}$ in each direction, with additional $x$-axis expansion proportional to excitation amplitude to prevent particle ejection.

**Boundary method selection.** Both DBC (`BoundaryMethod=1`) and mDBC (`BoundaryMethod=2`) \cite{english2021} are supported. Since mDBC is unsupported by DesignSPHysics \cite{designsphysics2023}, SloshAgent provides the only automated path to mDBC sloshing simulations.

**SWL gauge placement.** Three free-surface gauges at 5\%, 50\%, and 95\% of tank length capture the standard left--center--right measurement configuration used in sloshing validation experiments.

**Motion specification.** Horizontal sinusoidal motion uses the correct `mvrectsinu` vector syntax (`freq x="..." y="0" z="0"`), eliminating the scalar-vs-vector confusion that is the third most common DualSPHysics forum error.

The key architectural insight is the separation of concerns: the LLM infers *parameters* (what dimensions, frequency, fill level) from natural language via reasoning; the tool handles *syntax* (how to express those parameters in valid XML) via deterministic code. Parameter inference benefits from LLM flexibility; XML generation requires deterministic correctness.
