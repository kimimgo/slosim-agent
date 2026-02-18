# Sloshing Simulation Workflow Pain Points
## Evidence-Based Analysis for Paper Introduction

**Date**: 2026-02-18
**Purpose**: Quantify "before SloshAgent" pain points for paper Introduction section
**Research Branch**: research/paper

---

## 1. DualSPHysics-Specific Difficulties

### 1.1 XML Configuration Complexity

DualSPHysics requires hand-crafted XML files following a specialized schema (XML Guide v4.2, a 200+ page document). Common XML mistakes documented in the DualSPHysics forums include:

- **Hydrostatic pressure initialization**: Simulations show ~0.5 kPa at t=0 while experiments show 0 kPa. Root cause: hydrostatic pressure from initial water column. Fix requires running ~1 second without motion for particle stabilization, then subtracting the hydrostatic baseline — a non-obvious procedure for newcomers. [Source: [Tank Sloshing Validation forum thread](https://forums.dual.sphysics.org/discussion/2616/tank-sloshing-validation)]

- **Motion definition errors**: Coupled rotational + translational motion (required for combined sloshing) fails silently if `mov` reference values differ. The motion runs only the first component, with no clear error message. The fix (use identical `ref='0'` for all motions) is undocumented in tutorials. [Source: [GitHub Discussion #216](https://github.com/DualSPHysics/DualSPHysics/discussions/216)]

- **fillbox seed point placement**: In 2D simulations, the `fillbox` seed point must be at the exact correct y-position or particle creation fails. An alternative (`hswl auto=false`) exists but is not mentioned in tutorials. [Source: [DualSPHysics FAQ](https://dual.sphysics.org/faq/)]

- **Constant 'b' errors**: Occur when fluid height (`hswl`) is zero or when gravity is zero. The error message does not indicate which parameter is the cause. [Source: DualSPHysics forums]

- **mvrotsinu vs. mvrectsinu attribute syntax**: Rotational sinusoidal motion uses scalar `v` attribute (`<freq v="0.613" />`), while translational sinusoidal uses vector attributes (`x/y/z`). Mixing these produces incorrect motion or silent failure.

### 1.2 DesignSPHysics GUI Limitations (the only graphical frontend)

DesignSPHysics is a FreeCAD plugin intended to lower the XML authoring barrier. However:

- **Beta quality**: Released v0.7.0 in September 2023. Official status: "early Beta stage and is not meant to be used in a stable environment." [Source: [GitHub Releases](https://github.com/DualSPHysics/DesignSPHysics/releases)]

- **mDBC not supported via GUI** (as of January 2024): mDBC boundary conditions — which are needed for accurate pressure prediction on boundaries (crucial for sloshing) — cannot be configured through DesignSPHysics. Users must fall back to hand-editing XML. [Source: [Issue #171](https://github.com/DualSPHysics/DesignSPHysics/issues/171)]

- **FreeCAD version dependency**: Requires FreeCAD 0.18+; version mismatches cause tool failures. An issue from 2020 showed users stuck on outdated v5.0.3 because upgrade paths were unclear. [Source: [Issue #93](https://github.com/DualSPHysics/DesignSPHysics/issues/93)]

- **Non-intuitive save workflow**: Users must choose "Save Case and Run GenCase" from a dropdown — not simply "Save Case." This is a critical but poorly documented distinction that causes tutorial failures. [Source: DualSPHysics forum discussion]

- **MK reference errors**: Incorrect references to Material Key (MK) values — pointing to MKBound when MKFluid is needed — cause motion definitions to silently fail. [Source: DesignSPHysics GitHub issues]

### 1.3 Boundary Condition: DBC vs. mDBC Decision

A fundamental parameter decision with significant accuracy implications:

- **Standard DBC** creates a non-physical gap between fluid and boundary particles, degrading pressure accuracy near walls. This is the default but produces noisy pressure fields in sloshing scenarios.

- **mDBC** (Modified Dynamic Boundary Conditions, introduced in v5.0) corrects this but introduces new complexity:
  - Not yet implemented for floating bodies
  - Requires normal vector computation, which is non-trivial for imported STL geometries
  - Documentation only covers basic cases; no examples for multi-body or complex geometries
  - [Source: [mDBC paper (Computational Particle Mechanics, 2021)](https://link.springer.com/article/10.1007/s40571-021-00403-3)]

### 1.4 Post-Processing Complexity

After GPU simulation, a separate multi-tool post-processing chain is required:

1. `PartVTK` — convert binary particle data to VTK format
2. `MeasureTool` — extract time series at measurement points (requires correctly formatted `POINTS` header in probes file)
3. `ParaView` or `pvpython` — visualization (requires separate installation, scripting knowledge)
4. Custom scripts — comparison with experimental data

**Documented errors**:
- "Number of particles is zero" exception in `JForces2::PrepareData` when modified solver code is used. [Source: [PostProcessing error thread](https://forums.dual.sphysics.org/discussion/2680/postprocessing-error)]
- PartVTK and MeasureTool command-line flags differ in behavior and are not interchangeable with the documented `-save:` syntax (positional args required).

### 1.5 Particle Resolution (dp) Sensitivity

The particle spacing parameter `dp` controls both accuracy and computational cost:

- Finer `dp` → higher accuracy but particle count scales as O(1/dp³), making costs cubic
- Typical convergence studies test 5 resolution levels (e.g., dp = l/75, l/93.75, l/125, l/187.5, l/375)
- Artificial viscosity coefficient `α` requires trial-and-error tuning (values in literature range 0.01–0.5)
- `coefsound` parameter (speed of sound multiplier) has significant influence on pressure fields; value of 60 is common but derived through numerical trials, not formula
- [Source: [MDPI Journal of Marine Science and Engineering (2019)](https://www.mdpi.com/2077-1312/7/8/247)]

---

## 2. OpenFOAM / General Mesh-Based CFD Difficulties

### 2.1 Learning Curve: Weeks to Months

- OpenFOAM learning curves are explicitly described as requiring "a few weeks/months before getting something sensible out of it"
- The official "3 Weeks Series" tutorial is the minimum recommended entry path for serious work
- Components each requiring separate expertise: fluid dynamics, geometry/meshing, numerical methods, data analysis, computing/programming
- [Source: [curiosityFluids blog](https://curiosityfluids.com/2019/04/23/tips-for-tackling-the-openfoam-learning-curve/)]

### 2.2 interDyMFoam Sloshing Setup Complexity

Setting up a sloshing simulation in OpenFOAM (`interDyMFoam`) requires:

1. **blockMesh** — structured mesh generation (separate syntax/tool)
2. **dynamicMeshDict** — mesh motion configuration (SDA motion functions, 6-DOF parameters)
3. **transportProperties** — fluid properties for air and water phases
4. **VOF model configuration** — alpha field initialization, interface compression schemes
5. **User-Defined Functions (UDF)** — oscillating motion via frame motion UDF
6. **Version dependency**: `interDyMFoam` was deprecated in OpenFOAM v7 (2019); cases written for older versions fail on newer installations requiring solver migration

Typical mesh sizes: 144,000 cells for a 3D sloshing tank case.

### 2.3 Violent Sloshing: Known Unsolved Challenges

For near-resonance / violent sloshing:

> "There are still many uncertainties associated with wave breaking and splashing, formation of air pocket and air bubbles, behaviours of impact pressure and dynamic interaction between wave impact and structural response during violent sloshing."

- [Source: ISOPE proceedings, inter-laboratory benchmark comparisons]

---

## 3. Validation Challenges

### 3.1 Pressure Peak Stochasticity

The most important metric for sloshing validation — impact pressure peaks — is inherently stochastic:

- Only **9 of 16 peak pressures** from CFD fall within ±2 standard deviations of experimental results in benchmark studies
- Standard deviations are "quite large compared to the peaks themselves"
- Even nominally identical repeated experiments produce scatter: "residual sub-surface velocities undoubtedly play some role in the variability of the subsequent wave impact pressures"
- The most difficult validation cases are "those involving large free-surface deformations and violent impact phenomena" — precisely the resonance cases most relevant industrially
- [Source: ISOPE benchmark papers, [Tandfonline: Comparison of experimental and numerical sloshing loads](https://www.tandfonline.com/doi/full/10.1080/17445302.2010.522372)]

### 3.2 Initialization Artifact: The "Cold Start" Problem

A fundamental SPH-specific challenge:

- Initial particle arrangement is artificial (structured grid)
- Physical reorganization takes ~1 second of simulation time
- During this period, pressure data is dominated by numerical noise
- **This is not documented** in tutorials; users discover it through validation failures
- Correct procedure: pre-run with motion disabled → subtract hydrostatic baseline → begin measurement

### 3.3 Natural Frequency Calculation Errors

A critical conceptual trap for non-experts:

- Engineers often assume resonance occurs at the natural sloshing frequency
- In reality, due to **nonlinear effects**, resonance frequency differs from natural frequency
- The transition from "soft-spring" to "hard-spring" response (fill level dependent) is not captured in linear analytical formulae
- Fill level sensitivity: natural period changes non-linearly with fill depth
- For wave breaking and hydraulic jump conditions, "nonlinearities... may invalid the analytical solutions"
- [Source: [ResearchGate: Fill level and frequency studies](https://www.researchgate.net/publication/241075215)]

### 3.4 Inter-Laboratory Measurement Variability (ISOPE Benchmark)

The 2012 ISOPE sloshing benchmark (the premier experimental reference) found:

- Primary sources of inter-lab variability: tank motion accuracy, tank alignment, filling precision
- Piezoelectric pressure sensors show thermal sensitivity — temperature differential between fluid and sensor corrupts measurements
- Second benchmark required high-speed video cameras and independent motion measurement systems to achieve repeatability

---

## 4. Cost and Time Quantification

### 4.1 CFD Consulting Market Rates

From [CadCrowd CFD rate analysis](https://www.cadcrowd.com/blog/what-are-cfd-engineering-rates-consulting-services-costs-company-pricing/):

| Service Level | Cost Range |
|--------------|-----------|
| Simple 2D analysis | ~$2,000 |
| Complex transient analysis (automotive) | ~$20,000 |
| Hourly consulting rate | $80–$120/hour |
| Typical project timeline | 2–5 weeks |
| Follow-up simulations | A few days |

**Note**: Sloshing simulation falls in the "complex transient" category due to free-surface dynamics, motion coupling, and validation requirements. A single parametric study (e.g., 5 fill levels × 3 frequencies) could reach the $20,000 range.

### 4.2 Computational Runtime

GPU vs CPU comparison for DualSPHysics (1M+ particle cases):

| Hardware | Runtime (representative large case) |
|----------|-------------------------------------|
| Intel Core i7 CPU (single core) | 40+ hours |
| GTX 480 GPU (2011-era) | ~42 minutes |
| Speedup factor | 56×–100× |

Modern RTX 4090 (our platform) provides further ~3-5× improvement over GTX 480-era GPUs. CPU-only sloshing simulations are impractical for iterative workflows.

### 4.3 Workflow Automation Impact (AI4Science Analogies)

From analogous domains where LLM agents have been measured:

- **Automotive CFD (Siemens Simcenter)**: "Model generation from ~two months to ~one week" with automation tools (>80% reduction)
- **Coscientist (chemistry)**: Nobel-level experiment designed and executed in <4 minutes vs. days/weeks for human chemists
- **OpenFOAMGPT 2.0**: 100% case reproduction rate with multi-agent framework; simple cases: 2 interaction rounds / ~6,600 tokens; complex cases: 10 rounds / ~96,000 tokens
- **Foam-Agent**: 83.6% success rate with Claude 3.5 Sonnet vs. 37.3% baseline (OpenFOAM-GPT)
- **NL2FOAM fine-tuned model**: 82.6% first-attempt XML generation success rate

---

## 5. The Expert–Beginner Gap

### 5.1 Quantified Barriers

- OpenFOAM: "weeks to months" before any sensible output for newcomers
- DualSPHysics XML Guide: 200+ pages of documentation for a single input format
- Minimum prerequisites: fluid dynamics theory, SPH numerics, Linux command line, Python scripting (post-processing), ParaView scripting, Docker/HPC environment

### 5.2 Documented Beginner Failure Patterns

1. **Tutorial following fails**: Motion does not execute despite following official tutorial step-by-step (tank does not move) — caused by subtle axis definition
2. **Version mismatch cascade**: Outdated DesignSPHysics (v5.0.3 from 2019) incompatible with DualSPHysics v5.x
3. **Silent parameter errors**: Wrong `mv` reference → sequential instead of simultaneous motions, with no error message
4. **Post-processing chain breaks**: GenCase → solver → PartVTK → MeasureTool pipeline fails at any step; errors in one stage produce cryptic failures in downstream stages
5. **Overly ambitious first projects**: Violent resonance cases attempted before mastering basic setup — produces numerically unstable or physically incorrect results

### 5.3 Institutional Knowledge Problem

From Rescale's analysis of simulation workflows:
> "Institutional knowledge becomes trapped in rigid systems. When a senior engineer changes teams or leaves the company, their accumulated expertise often leaves with them."

This is acute for sloshing simulation where expert knowledge includes:
- Which dp values are stable for which geometries
- How to set coefsound and viscosity for specific liquids
- How to handle the cold-start initialization artifact
- Which boundary condition (DBC vs. mDBC) to use for each case type

---

## 6. Comparison with Existing LLM-CFD Tools (and Their Limitations)

| System | Solver | Domain | SPH | Sloshing | Local LLM |
|--------|--------|--------|-----|----------|-----------|
| OpenFOAMGPT v1/v2 | OpenFOAM | General CFD | No | No | No |
| FoamGPT | OpenFOAM | General CFD | No | No | Partial |
| Foam-Agent | OpenFOAM | General CFD | No | No | No |
| ChatCFD | OpenFOAM | General CFD | No | No | No |
| NL2FOAM | OpenFOAM | General CFD | No | No | No |
| **SloshAgent (ours)** | **DualSPHysics** | **Sloshing** | **Yes** | **Yes** | **Yes (Qwen3)** |

All existing LLM-CFD systems target OpenFOAM (mesh-based, general). None address:
- SPH (particle-based) solvers
- Domain-specific sloshing scenarios
- Local/private LLM deployment
- GPU-accelerated SPH pipeline automation

---

## 7. Key Quotes for Paper Use

**On CFD expertise barriers** (Foam-Agent, NeurIPS ML4PS 2025):
> "CFD simulation requires a wide range of expertise in fluid mechanics, numerical methods, geometric reasoning, and high-performance computing... OpenFOAM's text-based configuration system requires familiarity with specific file formats, syntax, and interdependencies, creating steep learning curve barriers especially for newcomers and interdisciplinary researchers."

**On simulation cost** (Siemens Simcenter blog):
> "Some papers stated such simulations were too troublesome and that the tools were not efficient enough to provide results in an industrial timeline."

**On validation difficulty** (ISOPE benchmark results):
> "Only nine of 16 peak pressures fell within two standard deviations of experimental results, and standard deviations were quite large compared to the peaks themselves."

**On knowledge loss** (Rescale agentic engineering):
> "When a senior engineer changes teams or leaves the company, their accumulated expertise often leaves with them."

**On automation potential** (Siemens workflow analysis):
> "Automated workflow solutions have proven to decrease model generation from around two months to approximately one week, a savings of over 80% in project/engineering time."

---

## 8. Summary Table for Introduction Paragraph

| Pain Point | Evidence | Impact |
|-----------|----------|--------|
| XML specification | 200+ page guide, 5+ documented error types | Hours–days per case for experts |
| GUI (DesignSPHysics) | Beta quality, mDBC unsupported, FreeCAD dependency | Limits to basic DBC cases only |
| Learning curve | "Weeks to months" before useful output | Excludes non-CFD engineers entirely |
| Consulting cost | $80–120/hr, $2K–$20K/project, 2–5 weeks | Prohibitive for design iteration |
| Validation uncertainty | 9/16 pressure peaks within 2σ | Expert judgment required to interpret |
| Parameter tuning | dp, α, coefsound all require trial-and-error | 3–5 resolution studies standard |
| Post-processing chain | 4+ tools, multiple formats, custom scripts | Additional hours per simulation |
| Motion coupling | Silent errors, undocumented `ref` syntax | Tutorial-following fails |
| Cold-start artifact | ~1s initialization, hydrostatic subtraction | Non-obvious, undocumented |
| Expert knowledge loss | Institutional knowledge exits with engineers | Every project restarts from scratch |

---

## Sources

- [Tank Sloshing Validation — DualSPHysics Forums](https://forums.dual.sphysics.org/discussion/2616/tank-sloshing-validation)
- [Sloshing tank tutorial problem — DualSPHysics Forums](https://forums.dual.sphysics.org/discussion/2397/sloshing-tank-tutorial-problem)
- [A beginner's discussion of the problems encountered — DualSPHysics Forums](https://forums.dual.sphysics.org/discussion/2634/a-beginners-discussion-of-the-problems-encountered)
- [PostProcessing error — DualSPHysics Forums](https://forums.dual.sphysics.org/discussion/2680/postprocessing-error)
- [Predefined coupled motion — GitHub Discussion #216](https://github.com/DualSPHysics/DualSPHysics/discussions/216)
- [Troubles running rotational motion — DesignSPHysics Issue #93](https://github.com/DualSPHysics/DesignSPHysics/issues/93)
- [mDBC with FreeCAD — DesignSPHysics Issue #171](https://github.com/DualSPHysics/DesignSPHysics/issues/171)
- [DesignSPHysics Releases (v0.7.0 Beta)](https://github.com/DualSPHysics/DesignSPHysics/releases)
- [mDBC paper: English et al. 2021 (Computational Particle Mechanics)](https://link.springer.com/article/10.1007/s40571-021-00403-3)
- [Experimental Validation SPH Sloshing (MDPI J. Marine Sci. Eng. 2019)](https://www.mdpi.com/2077-1312/7/8/247)
- [OpenFOAM learning curve — curiosityFluids 2019](https://curiosityfluids.com/2019/04/23/tips-for-tackling-the-openfoam-learning-curve/)
- [Tank Sloshing Simulation "Slosh My World" — Siemens Simcenter Blog](https://blogs.sw.siemens.com/simcenter/tank-sloshing-simulation-slosh-my-world/)
- [CFD consulting rates — CadCrowd](https://www.cadcrowd.com/blog/what-are-cfd-engineering-rates-consulting-services-costs-company-pricing/)
- [Agentic Engineering for M&S Workflows — Rescale](https://rescale.com/blog/agentic-engineering-for-modeling-and-simulation-workflows/)
- [Foam-Agent: Multi-Agent CFD Framework (NeurIPS ML4PS 2025)](https://arxiv.org/html/2505.04997)
- [OpenFOAMGPT 2.0 — arxiv 2504.19338](https://arxiv.org/html/2504.19338v1)
- [OpenFOAMGPT 1.0 — arxiv 2501.06327](https://arxiv.org/html/2501.06327v1)
- [Comparison of experimental and numerical sloshing loads (Tandfonline)](https://www.tandfonline.com/doi/full/10.1080/17445302.2010.522372)
- [DualSPHysics GPU performance — PLOS One](https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0020685)
- [ISOPE 2012 sloshing benchmark — OnePetro](https://onepetro.org/ISOPEIOPEC/proceedings-abstract/ISOPE13/ISOPE13/ISOPE-I-13-515/13505)
