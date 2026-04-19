# slosim-agent

Natural-language front-end for DualSPHysics GPU sloshing simulations. The
agent turns a Korean or English prompt into a full pipeline run — case
XML generation, GenCase preprocessing, GPU solve, post-processing, and
physics-aware analysis — and reads the residuals back to decide what to
do next.

Built on a fork of [OpenCode](https://github.com/opencode-ai/opencode)
(Go + BubbleTea) integrated with
[DualSPHysics v5.4](https://dual.sphysics.org/) and local Qwen3 via
Ollama.

---

## What it actually does

```
> "0.6m x 0.3m 탱크에 물 50% 채우고 0.5Hz로 흔들어"

[agent] xml_generator      → cases/sloshing_case.xml
[agent] gencase            → /data/out/sloshing_case (Docker, CUDA 12.6)
[agent] solver (GPU)       → RTX 4090, t = 5/f s
[agent] partvtk            → simulations/sloshing_case/vtk/
[agent] analysis           → below resonance → "일반 출렁임"
```

The agent is not a shell wrapper. It picks numerical parameters, reads
solver residuals, and replans on divergence. The next section shows
exactly where.

---

## Where the LLM makes decisions

Three surfaces where the model's output changes the simulation outcome
— not just its logging.

| Surface | What the model decides | Source |
|---|---|---|
| **Parameter inference** | `dp`, `time_max`, fluid height, excitation amplitude when the user omits them; standard tank dimensions by colloquial name ("LNG 탱크", "실험 탱크"); resonance-frequency-aware excitation setup. 5 numeric default rules + 4 named tank presets | `sloshing_coder.go:77-88`, consumed by `xml_generator.go` |
| **Reflection on residuals** | Classifies the run as Normal / Near-Resonance / Resonance / Overshoot via `f/f₁`, which the agent loop consumes to decide whether to re-run with refined `dp`, invoke `measuretool`, or stop. Agent loop re-invokes the model after each tool result; loop cap 25 iterations | `analysis.go:92-125`, `monitor.go`, `agent.go:284-344` |
| **Failure-aware planning** | Detects divergence (`RhopOut`, Inf velocity) from solver logs; emits 6 concrete fix-action types (TimeStep, viscosity, dp, domain size, GPU memory, NaN recovery) with retry policy | `error_recovery.go:139-162` |

What the model does **not** decide: SPH kernel choice, MPI/GPU
partitioning, or the physics laws encoded in DualSPHysics itself. Those
are delegated.

---

## Architecture

```
cmd/root.go          Cobra CLI (interactive TUI / -p non-interactive)
internal/app/        Service orchestrator
internal/llm/
  agent/             Agent loop: coder, task, title, summarizer
  provider/          OpenAI / Gemini / GROQ / OpenRouter / XAI / Local (Ollama)
  prompt/            System prompts (sloshing-domain specialisation)
  tools/             15 DualSPHysics tools + 11 generic tools
internal/tui/        BubbleTea components (chat, simulation, results,
                     parametric, logs)
internal/db/         SQLite + sqlc + goose migrations
cases/               Pre-defined XML scenarios
```

Each tool implements the same interface (`BaseTool.Info() + Run(ctx,
call)`), so built-in tools and MCP tools are interchangeable. Providers
share one `SendMessages()` / `StreamResponse()` contract.

---

## DualSPHysics tool surface

Fifteen domain tools exposed to the agent. Several carry a feature ID
tied to the FRD (`docs/FRD_v0.1.md`).

| Tool | Feature | Responsibility |
|---|---|---|
| `xml_generator` | — | Tank dimensions → DualSPHysics XML |
| `stl_import` | STL-01 | Frosina-pattern STL boundary import |
| `gencase` | — | XML → binary particles |
| `solver` | — | GPU SPH solve |
| `partvtk` | — | VTK extraction |
| `measuretool` | — | Probe-point time series |
| `baffle_generator` | — | Internal-baffle STL generation |
| `monitor` | MON-01 | Live `Run.csv` parsing |
| `error_recovery` | NFR-01 | Divergence detection + fix actions |
| `job_manager` | JOB-02 | Background job lifecycle |
| `seismic_input` | EXC-01 | Earthquake time-series ingestion |
| `parametric_study` | PARA-01 | Sweep orchestration |
| `result_store` | STORE-01 | Result persistence + retrieval |
| `analysis` | RPT-03 | Physics-aware interpretation (Qwen3) |
| `report` | — | Markdown report synthesis |

---

## Pre-defined cases

Sixty-five reference XML scenarios under `cases/`, covering benchmark
papers and industrial geometries:

- `SloshingTank_Def.xml` — baseline benchmark
- `Sloshing_{Normal,NearRes,Res}[_Guard]_Def.xml` — frequency regimes
  × baffle on/off
- `SPHERIC_Test10_{High,Low,Oil_Low}_Def.xml` — SPHERIC validation
- `Chen2018_*`, `English2021_*`, `Liu2024_*` — paper reproductions
- `Frosina2018_FuelTank_Def.xml` — reference for STL import pattern
- `NASA2023_Cylinder_Def.xml`, `Zhao2024_HorizCyl_Def.xml` — cylindrical
- `ISOPE2012_LNG_Def.xml` — large-scale LNG

Cases span three axes: geometry (rectangular, cylindrical, L-shaped),
fluid (water, oil), and excitation (sinusoidal sway/pitch, seismic).

---

## Research and paper trails

Development has proceeded alongside experiment bookkeeping and manuscript
preparation. These live in-tree as evidence of iterative work rather
than end-of-line artefacts.

- `paper-cs/` — manuscript targeting **EAAI / Engineering with
  Computers / Expert Systems with Applications**. Agent architecture +
  ablation study. Key results: parameter fidelity 82.2% (EXP-A),
  Qwen3:8B (78.5%) > LLaMA:70B (51.4%) on cross-model comparison,
  prompt contribution +56.5 pp in 2x2 factorial ablation.
  <!-- TODO: submission status -->
- `paper-pof/` — manuscript targeting **Physics of Fluids**. Physical
  validation against experiments. Key results: SPHERIC Test 10 Water
  lateral 19.5% mean error / r = 0.655; first quantitative SPH error
  report on SPHERIC Test 10 Oil (verified via 1,176-paper literature
  survey); baffle-induced 28.5% SWL reduction.
  <!-- TODO: submission status -->
- `paper/` — earlier integrated draft (legacy reference).
- `research/`, `research-v2/`, `research-v3/` — versioned experiment
  logs (EXP-A parameter fidelity, EXP-B ablation 2x2, EXP-C oil
  parametric, SPHERIC validation). `research-v2/experiment_registry.json`
  catalogues runs.
- Failure-to-learning loop documented in `research-v2/POSTMORTEM_V1.md`:
  peak amplitude error +63.5% → 19.5%, correlation r = -0.087 → 0.655,
  Oil peak detection 0/4 → 4/4 across the v1 → v2 cycle.

---

## Roadmap

- **Now.** Broaden validation coverage against published sloshing
  datasets; tighten divergence-recovery policy; stabilise parametric
  sweep UI.
- **Next.** Multi-GPU job scheduling; richer reflection loop that reads
  probe-point spectra, not just scalar residuals; uncertainty reporting
  over parameter sweeps.
- **Domain generalisation.** The decomposition (NL → parametric case
  generator → GPU solver → residual-driven reflection) is not specific
  to sloshing. It applies to any tank- or reactor-shaped problem where
  (a) a parametric case generator exists, (b) a GPU solver is available,
  and (c) the residuals can be post-processed. Extension slots already
  present in code: `stl_import.go` for arbitrary CAD geometries,
  `seismic_input.go` for non-sinusoidal excitation time-series,
  cylindrical and L-shaped tank types defined in
  `sloshing_coder.go:130-142` and `geometry.go`. Natural adjacents:
  multi-phase flows, chemical process tank dynamics, general CFD
  pre/post workflows.

---

## Testing

46 test files (~12 k lines) across `internal/`. Tool-level tests
sit next to their implementations (`internal/llm/tools/*_test.go`);
integration tests under the `*_integration_test.go` suffix exercise
multi-step flows (error recovery, parametric study, result store).

---

## Build & run

```bash
# Build
go build -o slosim-agent ./main.go

# Test
go test ./...

# Interactive TUI
./slosim-agent

# Non-interactive (single prompt)
./slosim-agent -c . -p "0.6m x 0.3m 탱크에 물 50% 채우고 0.5Hz로 흔들어" -q -f json

# DualSPHysics container
docker compose build
```

Requirements: NVIDIA GPU (tested on RTX 4090, sm_89), Docker with
NVIDIA Container Toolkit, Ollama serving a Qwen3 model
(`qwen3:32b` or `qwen3:latest`).

Config lives at `.opencode/config.json` or `~/.opencode.json`. Set
`LOCAL_ENDPOINT` for the Ollama URL.

---

## License & acknowledgements

MIT © 2026 Imgyu Kim. Forked from
[OpenCode](https://github.com/opencode-ai/opencode) (© 2025 Kujtim
Hoxha). SPH solver: [DualSPHysics](https://dual.sphysics.org/) (v5.4,
GPU). LLM weights: [Qwen3](https://qwenlm.github.io/) (Alibaba).
