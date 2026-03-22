// ============================================================
// Domain-Specific AI Agent Architecture for Automated
// Sloshing Simulation
// Target: Engineering with Computers
// ============================================================

#set document(
  title: "Domain-Specific AI Agent Architecture for Automated Sloshing Simulation: Enabling Non-Expert CFD Analysis through Tool-Augmented Local SLM",
  author: "Imgyu Kim",
)
#set page(margin: (x: 2.5cm, y: 2.5cm), numbering: "1")
#set text(font: "New Computer Modern", size: 11pt)
#set par(justify: true, leading: 0.65em)
#set heading(numbering: "1.1")
#set math.equation(numbering: "(1)")
#set bibliography(style: "springer-mathphys")

// Title
#align(center)[
  #text(size: 16pt, weight: "bold")[
    Domain-Specific AI Agent Architecture for Automated Sloshing Simulation: Enabling Non-Expert CFD Analysis through Tool-Augmented Local SLM
  ]
  #v(1em)
  #text(size: 11pt)[Imgyu Kim#super[1]]
  #v(0.5em)
  #text(size: 10pt, style: "italic")[
    #super[1]Department of Ocean Systems Engineering, Korea Maritime & Ocean University, Busan, Republic of Korea
  ]
  #v(0.3em)
  #text(size: 10pt)[Corresponding author: kimimgo\@gmail.com]
]

#v(2em)

// Abstract
#align(center)[#text(weight: "bold")[Abstract]]
#block(inset: (x: 2em))[
  Sloshing simulation using smoothed particle hydrodynamics (SPH) solvers requires expert knowledge in particle resolution, boundary conditions, and post-processing workflows, creating substantial barriers for non-expert engineers in maritime and industrial applications. While recent advances in large language model (LLM)-based agents have demonstrated promising results for computational fluid dynamics (CFD) automation, all existing approaches exclusively target finite volume method (FVM) solvers such as OpenFOAM and rely on cloud-based proprietary models. We present SlosimAgent, a domain-specific AI agent architecture that enables non-expert users to configure, execute, and analyze DualSPHysics sloshing simulations through natural language interaction. The system introduces three key innovations: (1) a three-layer architecture integrating a BubbleTea terminal user interface with a ReAct-based reasoning loop that orchestrates 18 specialized tools via the Model Context Protocol (MCP), (2) deployment of a locally-running Qwen3-32B small language model (SLM) through Ollama, eliminating cloud API dependencies while preserving domain reasoning capability, and (3) sloshing-specific domain knowledge encoding including automatic parameter determination rules, resonance frequency calculation, and standard tank configurations. The proposed architecture addresses five identified research gaps: the absence of any LLM agent for particle-based SPH solvers, the lack of local SLM deployment for CFD automation, the nonexistence of sloshing-specific AI agents, the unexplored application of MCP in scientific computing, and the absence of physics fidelity metrics for SPH simulations.
]

#v(1em)
*Keywords:* AI agent; sloshing simulation; DualSPHysics; SPH; small language model; tool augmentation; Model Context Protocol

#v(2em)

= Introduction

Liquid sloshing --- the dynamic motion of fluid within partially filled containers subjected to external excitation --- is a critical engineering concern across multiple industries, including liquefied natural gas (LNG) transportation @IbrahimSloshing2005 @FaltinsenSloshing2009, fuel tank design in aerospace, and seismic assessment of industrial storage facilities. Accurate prediction of sloshing dynamics, particularly peak pressures on tank walls and free surface behavior, is essential for structural integrity assessment and operational safety.

The smoothed particle hydrodynamics (SPH) method has emerged as a particularly suitable numerical approach for sloshing problems due to its Lagrangian, meshless nature, which naturally handles large free-surface deformations, wave breaking, and fluid fragmentation @DualSPHysics2021 @PillotonSPH2022. DualSPHysics, an open-source GPU-accelerated SPH solver, achieves computational speedups exceeding 100$times$ compared to single-core CPU implementations, making it practical for engineering-scale sloshing analyses @DualSPHysics2021. However, configuring a DualSPHysics simulation requires expertise in particle resolution selection, boundary condition specification, time integration parameters, and post-processing tool chains --- knowledge that is typically confined to specialized CFD researchers.

Recent advances in large language model (LLM)-based agents have demonstrated remarkable potential for automating CFD workflows. OpenFOAMGPT 2.0 @OpenFOAMGPT2 employs a multi-agent framework with GPT-4o to achieve 100% success rate across 450+ OpenFOAM simulations. Foam-Agent 2.0 @FoamAgent2 leverages Claude 3.5 Sonnet to attain 88.2% success rate with high-performance computing (HPC) deployment. ChatCFD @ChatCFD introduces the first physics fidelity metric (68.12%) using DeepSeek-R1/V3 across 315 benchmark cases. NL2FOAM @NL2FOAM fine-tunes Qwen2.5-7B on 28,716 CFD training cases, achieving 88.7% accuracy through chain-of-thought reasoning.

Despite these advances, a systematic analysis of the literature reveals five critical research gaps:

+ *No LLM agent for particle-based SPH solvers.* All existing LLM-CFD works --- OpenFOAMGPT @OpenFOAMGPT @OpenFOAMGPT2, Foam-Agent @FoamAgent2, ChatCFD @ChatCFD, NL2FOAM @NL2FOAM, MetaOpenFOAM @MetaOpenFOAM, and Engineering.ai @EngineeringAI --- exclusively target OpenFOAM or other FVM solvers. No publication addresses DualSPHysics or any other SPH solver.

+ *No local SLM deployment for CFD automation.* All approaches rely on cloud-based proprietary models: GPT-4o (OpenFOAMGPT), Claude 3.5 Sonnet (Foam-Agent), DeepSeek-R1/V3 (ChatCFD). This creates data privacy concerns and infrastructure dependencies unsuitable for industrial deployment, particularly for defense and proprietary applications.

+ *No sloshing-specific AI agent.* Existing agents target general-purpose CFD (turbulence, heat transfer, multiphase) without domain-specific knowledge for sloshing: resonance frequency estimation, fill-level dependent behavior, excitation characterization, and standard tank configurations.

+ *No MCP-based tool integration for scientific computing.* The Model Context Protocol @MCP has been widely adopted for IDE and developer tooling but has not been applied to scientific simulation workflows.

+ *No SPH physics fidelity metrics.* ChatCFD pioneered physics fidelity evaluation (68.12%) for FVM solvers @ChatCFD, but no equivalent metrics exist for SPH-based simulations.

This paper presents SlosimAgent, a domain-specific AI agent architecture that addresses all five gaps. Our contributions are:

+ *First LLM agent for DualSPHysics:* A three-layer architecture (TUI $arrow$ Agent Core $arrow$ Tool Layer $arrow$ DualSPHysics) with ReAct-based reasoning that orchestrates the complete SPH simulation workflow.

+ *Local SLM deployment for CFD:* Integration of Qwen3-32B via Ollama for fully local inference, demonstrating that a 32-billion parameter model can effectively drive domain-specific CFD automation without cloud dependencies.

+ *Sloshing-specific domain knowledge encoding:* Automatic parameter determination rules (particle resolution, simulation time, boundary conditions), resonance frequency calculation, standard tank configurations, and terminology translation for non-expert users.

+ *MCP-based scientific tool integration:* The first application of the Model Context Protocol to scientific computing, providing a standardized interface for 18 specialized simulation tools.

The remainder of this paper is organized as follows. Section 2 reviews related work across four domains. Section 3 presents the system architecture. Section 4 details the implementation. Section 5 describes the experimental design. Section 6 concludes with discussion and future directions.


= Related Work

== LLM-Based Agents for CFD

The application of LLMs to CFD automation has rapidly evolved from simple code generation to sophisticated multi-agent architectures. OpenFOAMGPT @OpenFOAMGPT introduced the concept of LLM-assisted OpenFOAM case setup using retrieval-augmented generation (RAG) with GPT-4o. Its successor, OpenFOAMGPT 2.0 @OpenFOAMGPT2, deployed a multi-agent framework comprising pre-processing, prompt generation, simulation, and post-processing agents, achieving a 100% success rate across 450+ simulation configurations.

Foam-Agent 2.0 @FoamAgent2 demonstrated that Claude 3.5 Sonnet could drive OpenFOAM simulations with 88.2% success rate while enabling deployment on HPC clusters. MetaOpenFOAM @MetaOpenFOAM applied the MetaGPT assembly-line paradigm to decompose CFD workflows into specialized sub-agent roles.

ChatCFD @ChatCFD made a significant contribution by introducing the first physics fidelity metric for LLM-generated CFD configurations. Using DeepSeek-R1/V3 with a structured knowledge base of 315 benchmark cases, it achieved 82.1% execution success and 68.12% physics fidelity. This metric --- measuring whether simulation results conform to known physical behavior --- represents an important step beyond mere execution success.

NL2FOAM @NL2FOAM explored fine-tuning smaller models (Qwen2.5-7B) on domain-specific data, achieving 88.7% accuracy on 28,716 training cases. Their ablation study demonstrated that chain-of-thought prompting improved accuracy by 10.5% and pass\@1 by 20.9%, highlighting the importance of structured reasoning for CFD configuration.

Engineering.ai @EngineeringAI achieved 100% success with integrated Gmsh mesh generation, while recent work on lightweight coding agents for CFD @CodingAgentsCFD demonstrated that single-agent architectures can be effective without multi-agent overhead. SimLLM @SimLLM applied multi-stage supervised fine-tuning and direct preference optimization (SFT+DPO) for simulation code generation, improving Qwen2.5-Coder-7B executability from 80.4% to 86.0%.

A critical observation across all these works is their exclusive focus on FVM solvers, particularly OpenFOAM. No prior work addresses particle-based methods such as SPH, which require fundamentally different configuration paradigms (particle resolution instead of mesh quality, boundary particle layers instead of mesh boundary conditions).


== Tool-Augmented Language Model Agents

The ReAct framework @ReAct established the paradigm of interleaving reasoning traces with action execution, enabling LLMs to dynamically interact with external tools. Toolformer @Toolformer demonstrated that language models can learn to invoke APIs autonomously. ChemCrow @ChemCrow showed that augmenting an LLM with 13 domain-specific chemistry tools substantially improved performance on expert-level tasks, providing a template for domain-specific tool augmentation.

AstaBench @AstaBench evaluated 57 AI agents across 22 architectures, providing empirical evidence that tool-augmented agents consistently outperform pure language model approaches on complex, multi-step tasks. Multi-agent frameworks including AutoGen @AutoGen have further demonstrated that agent-to-agent communication can improve collaborative task completion.

Our work extends this line of research by applying tool augmentation to scientific simulation, where tools must interface with compiled GPU solvers, manage long-running background processes, and handle domain-specific file formats.


== DualSPHysics and SPH for Sloshing

DualSPHysics @DualSPHysics2021 is an open-source SPH solver that leverages CUDA-based GPU acceleration to achieve computational speedups exceeding 100$times$ over CPU implementations. The solver has been extensively validated for sloshing problems, including long-time LNG tank simulations @PillotonSPH2022, fluid-structure interaction in elastic tanks @SloshingElasticTank, and prismatic tank validation with less than 1% error compared to experimental data @SPHValidation2019.

The SPH method offers natural advantages for sloshing simulation: free-surface tracking without volume-of-fluid reconstruction, handling of wave breaking and fluid fragmentation, and straightforward implementation of moving boundaries. However, the DualSPHysics workflow involves multiple discrete steps --- XML case definition, GenCase particle generation, GPU solver execution, PartVTK/MeasureTool post-processing --- each requiring specific parameter knowledge.

DesignSPHysics @DesignSPHysics provides a FreeCAD-based graphical interface for DualSPHysics, but requires CAD expertise and does not incorporate AI-assisted parameter selection. No prior work has attempted to automate the DualSPHysics workflow through AI agent technology.


== Small Language Models for Local Deployment

The trend toward smaller, locally-deployable language models has been driven by privacy concerns, latency requirements, and infrastructure costs. The Qwen2.5-Coder report @Qwen25Coder demonstrated that a 7-billion parameter model can achieve competitive coding performance (73.7 on the Aider benchmark, comparable to GPT-4o). Industry analyses project that small language model (SLM) deployments will exceed large language model deployments by a factor of three by 2027.

Ollama @Ollama has emerged as the de facto standard for local LLM deployment, providing quantized model serving with hardware-optimized inference. For our application, local deployment is essential: field engineers working with proprietary tank designs cannot expose simulation parameters to cloud services, and offshore installations may lack reliable internet connectivity.

NL2FOAM @NL2FOAM and SimLLM @SimLLM have demonstrated that fine-tuned small models can match or exceed proprietary models on domain-specific CFD tasks, suggesting that a 32-billion parameter model (Qwen3-32B) with appropriate prompting may be sufficient for sloshing simulation automation.


= System Architecture

SlosimAgent employs a three-layer architecture designed to separate concerns between user interaction, agent reasoning, and domain-specific tool execution. @fig-architecture illustrates the overall system structure.

#figure(
  ```
  ┌──────────────────────────────────────────────┐
  │            Layer 1: TUI (BubbleTea)           │
  │   ┌─────────────┐    ┌───────────────────┐   │
  │   │  Chat View   │    │  Sim Dashboard    │   │
  │   │  (NL input)  │    │  - Job Status     │   │
  │   │              │    │  - Progress Bar   │   │
  │   │              │    │  - Live Results   │   │
  │   └─────────────┘    └───────────────────┘   │
  ├──────────────────────────────────────────────┤
  │          Layer 2: Agent Core (Go)             │
  │   ┌───────────┐    ┌──────────────────┐      │
  │   │ ReAct Loop │    │  Tool Executor   │      │
  │   │ (Qwen3    │    │  (MCP Client)    │      │
  │   │  32B)     │    │                  │      │
  │   └───────────┘    └──────────────────┘      │
  ├──────────────────────────────────────────────┤
  │        Layer 3: Tools (18 MCP Tools)          │
  │   ┌────────┐ ┌──────┐ ┌────────┐ ┌──────┐   │
  │   │XML Gen │ │GenCase│ │ Solver │ │PartVTK│  │
  │   └────────┘ └──────┘ └────────┘ └──────┘   │
  │   ┌────────┐ ┌──────┐ ┌────────┐ ┌──────┐   │
  │   │Measure │ │PV Insp│ │PV Rend │ │PV Anim│  │
  │   └────────┘ └──────┘ └────────┘ └──────┘   │
  │   ┌────────┐ ┌──────┐ ┌────────┐ ┌──────┐   │
  │   │PV Slice│ │PV Clip│ │PV Stats│ │Report │  │
  │   └────────┘ └──────┘ └────────┘ └──────┘   │
  ├──────────────────────────────────────────────┤
  │       DualSPHysics v5.4 (CUDA / GPU)         │
  │       Background Job Manager (pubsub)         │
  └──────────────────────────────────────────────┘
  ```,
  caption: [Three-layer architecture of SlosimAgent. Layer 1 provides the user interface, Layer 2 contains the agent reasoning and tool orchestration logic, and Layer 3 exposes domain-specific tools via MCP.],
) <fig-architecture>

== Layer 1: Terminal User Interface

The TUI layer is built on the BubbleTea framework (Go), forked from OpenCode @OpenCode to provide a domain-specialized interface. It comprises two primary views:

- *Chat View:* Accepts natural language input from the user and displays the agent's responses, tool invocations, and reasoning traces.
- *Simulation Dashboard:* Provides real-time monitoring of submitted simulation jobs, including progress bars, convergence indicators, error detection alerts, and result previews.

The TUI layer communicates with the Agent Core via an event-driven publish-subscribe (pubsub) architecture, ensuring that long-running GPU simulations do not block the user interface.


== Layer 2: Agent Core

The Agent Core implements a ReAct reasoning loop @ReAct that interleaves natural language reasoning with tool execution. Given a user query $q$, the agent iteratively generates thought-action-observation triples:

$ "Thought"_t arrow "Action"_t arrow "Observation"_t $

where each action invokes one of the 18 available tools via the MCP protocol, and observations are the structured tool responses. The loop terminates when the agent determines that the user's request has been fulfilled.

The language model backbone is Qwen3-32B @Qwen25Coder, served locally through Ollama @Ollama with 4-bit quantization (Q4_K_M), requiring approximately 20 GB of VRAM. The model receives a domain-specific system prompt (Section 4.2) that encodes sloshing simulation knowledge and tool usage conventions.

The tool executor implements the MCP client specification @MCP, communicating with tool servers via JSON-RPC 2.0. Each tool call includes structured input parameters and returns typed responses (text, image, or error), enabling the agent to chain multiple tools in the correct execution order.


== Layer 3: Tool Suite

The tool layer exposes 18 specialized tools organized into five functional categories (@tbl-tools). Each tool implements the `BaseTool` interface:

```go
type BaseTool interface {
    Info() ToolInfo
    Run(ctx context.Context, params ToolCall) (ToolResponse, error)
}
```

#figure(
  table(
    columns: (auto, auto, auto),
    inset: 8pt,
    align: (left, left, left),
    [*Category*], [*Tool*], [*Function*],
    [Case Setup], [`xml_generator`], [Generate DualSPHysics XML case file],
    [], [`geometry`], [Create complex tank geometries (cylindrical, L-shaped)],
    [], [`seismic_input`], [Parse external excitation data (CSV/seismic)],
    [Simulation], [`gencase`], [Execute GenCase particle generation],
    [], [`solver`], [Launch DualSPHysics GPU solver (background)],
    [], [`error_recovery`], [Diagnose and recover from solver errors],
    [Post-processing], [`partvtk`], [Convert binary output to VTK format],
    [], [`measuretool`], [Extract pressure/elevation time series],
    [Visualization], [`pv_inspect_data`], [Query VTK metadata (fields, bounds)],
    [], [`pv_render`], [Render field visualization (PNG)],
    [], [`pv_animate`], [Generate time-lapse animation (MP4)],
    [], [`pv_slice`], [Cross-section visualization],
    [], [`pv_contour`], [Isosurface visualization],
    [], [`pv_clip`], [Clipping plane visualization],
    [], [`pv_streamlines`], [Velocity streamline visualization],
    [], [`pv_plot_over_line`], [Line sampling for graph data],
    [], [`pv_extract_stats`], [Field statistics (min/max/mean)],
    [Reporting], [`report`], [Generate Markdown analysis report],
  ),
  caption: [Complete tool suite organized by functional category. All 18 tools are accessible via the Model Context Protocol.],
) <tbl-tools>


= Implementation

== Tool Execution Pipeline

A typical sloshing simulation workflow follows a prescribed tool invocation sequence. Given a natural language request such as _"Simulate sloshing in a small rectangular tank with 50% fill level at 1 Hz excitation"_, the agent executes:

+ `xml_generator` $arrow$ Generate XML case definition with inferred parameters
+ `gencase` $arrow$ Create particle geometry from XML
+ `solver` $arrow$ Launch GPU simulation (background job)
+ `partvtk` $arrow$ Convert results to VTK format
+ `measuretool` $arrow$ Extract pressure and elevation data
+ `pv_inspect_data` $arrow$ Verify output field metadata
+ `pv_render` $arrow$ Generate visualization snapshots
+ `pv_animate` $arrow$ Create time-lapse animation
+ `report` $arrow$ Compile analysis report with results

The agent enforces this ordering through its system prompt while allowing flexibility to skip or repeat steps based on intermediate observations (e.g., re-running with adjusted parameters if solver divergence is detected).


== Domain Knowledge Encoding <domain-knowledge>

The sloshing-specific system prompt encodes three categories of domain knowledge:

*Automatic parameter determination.* When users omit simulation parameters, the agent applies the following rules:
- Particle spacing: $d p = min(L, W, H) \/ 50$, bounded to $[0.005, 0.05]$ meters
- Simulation duration: $t_"max" = 5 \/ f$ seconds, where $f$ is the excitation frequency
- Default fill level: 50% of tank height
- Default excitation amplitude: 5% of tank length

*Resonance frequency calculation.* For rectangular tanks, the first-mode natural frequency is computed as:

$ f_1 = 1/(2 pi) sqrt(g dot pi / L dot tanh(pi / L dot h)) $ <eq-resonance>

where $g = 9.81 "m/s"^2$, $L$ is the tank length, and $h$ is the fluid depth. This enables the agent to warn users about resonance conditions and suggest appropriate simulation durations.

*Standard tank configurations.* The system prompt includes pre-defined configurations for common engineering scenarios:
- LNG tank: $40 times 40 times 27$ m
- Ship tank: $20 times 10 times 8$ m
- Small laboratory tank: $1 times 0.5 times 0.6$ m
- Experimental validation tank: $0.6 times 0.3 times 0.4$ m

*Boundary condition selection.* Two boundary methods are supported:
- Dynamic Boundary Condition (DBC): Default method, computationally efficient
- Modified Dynamic Boundary Condition (mDBC): Higher pressure accuracy near walls, activated when users request "high-precision" or "precise boundary" analysis


== Local SLM Integration

The language model is Qwen3-32B served through Ollama with 4-bit quantization (Q4_K_M format). Key deployment parameters:

#figure(
  table(
    columns: (auto, auto),
    inset: 8pt,
    align: (left, left),
    [*Parameter*], [*Value*],
    [Model], [Qwen3-32B (Q4_K_M)],
    [VRAM requirement], [$approx$ 20 GB],
    [GPU], [NVIDIA RTX 4090 (24 GB)],
    [Inference framework], [Ollama],
    [Context window], [8,192 tokens],
    [System prompt size], [$approx$ 2,000 tokens],
  ),
  caption: [Local SLM deployment configuration. The model shares the GPU with DualSPHysics, with time-multiplexed resource allocation.],
) <tbl-deployment>

A notable design constraint is GPU resource sharing between the language model and the SPH solver. During simulation execution, Ollama releases VRAM to allow DualSPHysics full GPU utilization. The agent pauses LLM inference during solver execution and resumes for post-processing analysis.


== MCP Integration

The Model Context Protocol @MCP provides the standardized interface between the agent core and tool implementations. Each tool is registered as an MCP server, exposing:

- *Tool description:* Name, purpose, and parameter schema (JSON Schema)
- *Input validation:* Required and optional parameters with type checking
- *Structured response:* Typed responses (text, image, error) with optional metadata

This architecture enables tool composability --- new tools (e.g., additional post-processing filters or alternative solver backends) can be added by implementing the `BaseTool` interface and registering with the MCP server, without modifying the agent core.


= Experimental Design

To validate the proposed architecture, we design a comprehensive evaluation framework addressing three dimensions: execution success, physics fidelity, and usability.

== Benchmark Suite

We construct a benchmark of sloshing simulation cases spanning four complexity levels:

#figure(
  table(
    columns: (auto, auto, auto, auto),
    inset: 8pt,
    align: (left, left, left, left),
    [*Level*], [*Cases*], [*Description*], [*Validation*],
    [L1: Basic], [10], [Rectangular tank, sinusoidal excitation], [Analytical solution],
    [L2: Standard], [10], [Variable fill levels, multiple frequencies], [Experimental data],
    [L3: Complex], [5], [Cylindrical/L-shaped tanks, seismic input], [Published SPH results],
    [L4: Expert], [5], [Near-resonance, mDBC, parametric study], [Expert-configured reference],
  ),
  caption: [Sloshing simulation benchmark suite. 30 cases across four complexity levels.],
) <tbl-benchmark>

== Evaluation Metrics

We propose three categories of metrics:

*Execution metrics:*
- *Success rate:* Percentage of cases where the agent produces a completed simulation without human intervention
- *Tool sequence correctness:* Whether the agent invokes tools in the correct order with valid parameters

*Physics fidelity metrics (extending ChatCFD @ChatCFD to SPH):*
- *Free surface accuracy:* Comparison of wave elevation time series against reference data using normalized root-mean-square error (NRMSE)
- *Pressure fidelity:* Wall pressure time history compared to experimental/analytical references
- *Conservation metrics:* Total particle count stability, energy conservation over simulation duration

*Usability metrics:*
- *Natural language comprehension:* Correct parameter extraction from varied user phrasings
- *Iteration count:* Number of agent turns required to complete a simulation request
- *Error recovery rate:* Ability to diagnose and resolve solver divergence or configuration errors

== Ablation Studies

We plan ablation experiments to isolate the contribution of each architectural component:

+ *Model size:* Qwen3-32B vs. Qwen3-8B (impact of parameter count on domain reasoning)
+ *Domain prompt:* With vs. without sloshing-specific system prompt (impact of domain knowledge encoding)
+ *Tool augmentation:* Full 18-tool suite vs. reduced tool set (impact of tool specialization)
+ *Boundary method:* DBC vs. mDBC (impact on physics fidelity)


= Conclusion

We have presented SlosimAgent, a domain-specific AI agent architecture for automated sloshing simulation that addresses five identified research gaps in the LLM-for-CFD literature. The three-layer architecture --- integrating a BubbleTea TUI, ReAct-based agent core with Qwen3-32B local inference, and 18 MCP-exposed tools --- provides the first AI agent capable of driving DualSPHysics SPH simulations through natural language interaction.

Our work makes four primary contributions: (1) the first LLM agent for a particle-based SPH solver, (2) local SLM deployment for CFD automation eliminating cloud dependencies, (3) sloshing-specific domain knowledge encoding enabling non-expert use, and (4) the first application of MCP to scientific computing workflows.

The proposed evaluation framework extends ChatCFD's physics fidelity concept to SPH simulations, introducing particle-specific metrics for free surface accuracy, pressure fidelity, and conservation properties. Planned ablation studies will quantify the contribution of model size, domain prompts, and tool specialization to overall system performance.

Future work will focus on three directions: (1) expanding the domain to multi-phase sloshing (oil-water, LNG with boil-off gas), (2) integrating fine-tuned models using the SFT+DPO approach demonstrated by SimLLM @SimLLM, and (3) developing a standardized SPH physics fidelity benchmark comparable to ChatCFD's 315-case FVM benchmark.

#bibliography("references.bib")
