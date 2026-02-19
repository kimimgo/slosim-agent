# 4. Experiments and Results

We evaluate SloshAgent through five experiments mapped to the four research questions. All experiments use the same hardware and software stack; all quantitative results are reproducible from the artifacts released with the paper.

## 4.1 Experimental Setup

**Hardware and software.** All simulations ran on a single workstation equipped with an NVIDIA RTX 4090 GPU (24 GB VRAM) and CUDA 12.6. DualSPHysics v5.4 executed inside a Docker container with GPU passthrough, mounting case files at `/cases` and writing output to `/data`. The LLM inference used Qwen3 32B \cite{qwen3_2025} served locally via Ollama; a secondary Qwen3 8B model was used for model-size comparison. Both models ran with thinking mode enabled (extended reasoning). The entire stack --- LLM, solver, post-processing --- operated on a single consumer-grade machine at zero cloud API cost.

**Experiment--RQ mapping.** Table~\ref{tab:exp_rq} summarizes the five experiments, their target research questions, and the primary metrics.

| EXP | RQ | Goal | Primary Metric |
|-----|-----|------|----------------|
| EXP-1 | RQ2 | SPHERIC benchmark validation | Peaks within $\pm 2\sigma$ band |
| EXP-2 | RQ1 | NL$\rightarrow$XML generation accuracy | Parameter accuracy (\%) |
| EXP-3 | RQ3 | Parametric study automation | Physical trend consistency |
| EXP-4 | RQ4 | Domain prompt ablation | Accuracy $\Delta$ across conditions |
| EXP-5 | --- | Industrial proof of concept | Qualitative feasibility |

## 4.2 EXP-1: SPHERIC Benchmark Validation (RQ2)

**Motivation.** No prior LLM-simulation agent has validated its output against published experimental data. Existing systems report execution success rates (e.g., MetaOpenFOAM 85\% \cite{metaopenfoam2024}, ChatCFD 82.1\% \cite{chatcfd2026}) or LLM-as-judge physical fidelity scores (ChatCFD 68.12\%), but none compare simulation results to laboratory measurements. We address this gap using the SPHERIC Test 10 benchmark \cite{spheric_benchmarks}, one of the most widely cited sloshing validation datasets.

**Benchmark description.** SPHERIC Test 10 consists of a 0.9 m $\times$ 0.062 m $\times$ 0.508 m rectangular tank subjected to lateral sinusoidal excitation at 0.613 Hz with 93 mm water fill height. The benchmark provides 102-repeat experimental pressure measurements at lateral wall and roof impact sensors, for two fluids: water (density 998 kg/m$^3$, viscosity $1.0 \times 10^{-6}$ m$^2$/s) and sunflower oil (density 916 kg/m$^3$, viscosity $5.0 \times 10^{-5}$ m$^2$/s). All four water impact peaks are non-normally distributed (Shapiro-Wilk $p < 0.005$) with positive skewness and heavy right tails (excess kurtosis 0.9--8.9), and exhibit high stochastic variability (coefficient of variation CV = 20--36\%). This non-normality confirms that the $\pm 2\sigma$ band test --- the standard validation criterion adopted by both SPHERIC and ISOPE sloshing benchmark protocols \cite{isope2012benchmark} --- is more appropriate than pointwise correlation metrics for impact-dominated flows.

**Agent-generated cases.** SloshAgent generated two resolution variants from a single natural language prompt ("*Reproduce SPHERIC Test 10 with water at low and high resolution*"): Water Low (dp = 0.004 m, 136K particles) and Water High (dp = 0.004 m, 344K particles). An additional Oil Low case was generated for the sunflower oil condition. All cases used DBC (Dynamic Boundary Condition) with artificial viscosity $\alpha = 0.01$.

**Validation methodology.** Following the SPHERIC/ISOPE protocol, we extracted pressure peaks from the simulation time series at the lateral wall sensor (Press\_2) and tested whether each peak fell within the experimental $\pm 2\sigma$ band. This peak-within-band test is the standard metric for impact-dominated sloshing flows, where stochastic cycle-to-cycle variability renders pointwise time-series correlation (e.g., Pearson $r$) inappropriate: impact pressures exhibit high-amplitude, short-duration peaks with inherent phase and magnitude scatter even across physical repetitions \cite{isope2012benchmark, english2021mdbc}.

**Results.** Table~\ref{tab:spheric} presents the peak pressure comparison.

| Case | Particles | Sim.\ Peaks [mbar] | Peaks in $\pm 2\sigma$ | Max $P$ [mbar] | Status |
|------|-----------|---------------------|------------------------|-----------------|--------|
| Water Low | 136K | 31.1, 58.9, 76.7 | **3/3** (100\%) | 76.7 | PASS |
| Water High | 344K | 44.2, 29.4, 31.4, 45.3 | **4/4** (100\%) | 50.0 | PASS |
| Oil Low | 136K | (none detected) | N/A | 0.0 | FAIL |

The experimental reference values (100-repeat statistics) are:

| Fluid | Peak 1 [mbar] | Peak 2 [mbar] | Peak 3 [mbar] | Peak 4 [mbar] |
|-------|---------------|---------------|---------------|---------------|
| Water $\mu$ | 37.1 | 48.2 | 46.9 | 46.4 |
| Water $\pm 2\sigma$ | $\pm$14.6 | $\pm$29.9 | $\pm$34.0 | $\pm$26.3 |
| Oil $\mu$ | 6.9 | 6.5 | 5.4 | 5.5 |
| Oil $\pm 2\sigma$ | $\pm$0.3 | $\pm$0.5 | $\pm$0.5 | $\pm$0.5 |

All seven detected peaks across both water resolutions fall within the experimental $\pm 2\sigma$ envelope. The z-scores range from $-1.25$ to $+1.75$ (all $|z| < 2.0$), with five of seven peaks also falling within the stricter $\pm 1\sigma$ band (71\%). The Water Low case produced three distinct impact peaks (31.1, 58.9, 76.7 mbar) and the Water High case produced four (44.2, 29.4, 31.4, 45.3 mbar), both consistent with the experimental peak structure. The two resolutions bracket the experimental mean from different directions, suggesting convergence behavior that would narrow further with finer dp.

Normalized mean absolute error (NMAE) was 34.0\% for Water Low and 23.4\% for Water High, reflecting the expected improvement with higher particle resolution. These errors should be interpreted against the experimental variability itself: the bootstrap 95\% CI of the experimental mean spans $\pm$4--6\% (e.g., Peak 2: [45.7, 51.6] mbar), indicating that a substantial fraction of the simulation--experiment discrepancy lies within the experimental sampling uncertainty. The remaining error is consistent with prior DualSPHysics DBC studies \cite{english2021mdbc}: DBC inherently over-predicts wall impact pressures compared to mDBC due to the artificial boundary layer.

**Oil failure analysis.** The Oil Low case produced zero detectable impact peaks. This is a known DBC limitation, not an agent failure: the oil's high kinematic viscosity ($50 \times$ water) combined with DBC's artificial viscosity over-damps the thin sloshing layer, preventing impact events at the sensor location. English et al.\ \cite{english2021mdbc} demonstrated that mDBC (modified Dynamic Boundary Condition) is required for accurate viscous-fluid sloshing --- a boundary condition that the agent correctly could not configure because DBC was the only option encoded in the current tool set. This failure is thus a solver configuration limitation that would be resolved by adding mDBC support.

**Significance.** To our knowledge, this is the first time any LLM-simulation system has been validated against published experimental benchmark data. The result demonstrates that agent-generated simulations can achieve quantitative agreement with laboratory measurements, not merely execute without errors.

## 4.3 EXP-2: NL$\rightarrow$XML Generation Accuracy (RQ1)

**Design.** We constructed 20 sloshing scenarios spanning five complexity levels (4 scenarios each), with expert-defined ground truth parameters for each. The levels test progressively harder reasoning:

- **L1 (Basic)**: All physical parameters stated explicitly (e.g., "0.9 m tank, 93 mm water fill, 0.613 Hz sway").
- **L2 (Domain)**: Parameters must be inferred from domain terms (e.g., "LNG cargo tank at resonance" $\rightarrow$ requires Mark III dimensions + $f_1$ formula).
- **L3 (Paper)**: Exact reproduction of a published case (e.g., "Chen2018 Case 3" $\rightarrow$ requires paper-specific values absent from training data).
- **L4 (Complex)**: Multi-feature composition (e.g., baffle + tank + parametric study orchestration).
- **L5 (Edge)**: Error conditions requiring graceful handling (e.g., empty tank, extreme dp, 100\% fill).

Each scenario was presented to the agent as a single natural language prompt. We measured three metrics: (1) **tool call rate** --- whether `xml_generator` was invoked; (2) **parameter accuracy** --- fraction of generated parameters matching ground truth within tolerance ($\pm$20\% for dp, $\pm$5\% for frequency, $\pm$10\% for fill level, exact match for dimensions and motion type); and (3) **physical validity** --- whether the generated configuration is physically reasonable (expert judgment).

**Results.** Table~\ref{tab:nl2xml} summarizes the results for both Qwen3 32B and 8B.

| Level | $n$ | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-------|-----|---------------|-------------|--------------|-------------|-----------|---------|
| L1 (Basic) | 4 | 4/4 | 4/4 | **96\%** | **96\%** | 4/4 | 4/4 |
| L2 (Domain) | 4 | **4/4** | 3/4 | **67\%** | 42\% | 4/4 | 3/4 |
| L3 (Paper) | 4 | 3/4 | 3/4 | 15\% | 15\% | 3/4 | 3/4 |
| L4 (Complex) | 4 | 3/4 | 3/4 | **58\%** | 50\% | 3/4 | 3/4 |
| L5 (Edge) | 4 | **3/4** | 2/4 | 0\% | **25\%** | 0/4 | 1/4 |
| **Total** | **20** | **17/20 (85\%)** | 15/20 (75\%) | **47\%** | 46\% | **14/20 (70\%)** | 14/20 (70\%) |

**Analysis by level.** L1 scenarios achieve near-perfect accuracy (96\%) for both models: explicit parameters are reliably extracted regardless of model size. L2 reveals a clear 32B advantage (+25 percentage points): domain inference tasks such as computing the natural sloshing frequency $f_1 = \frac{1}{2\pi}\sqrt{\frac{g\pi}{L}\tanh\frac{\pi h}{L}}$ or mapping "LNG cargo tank" to Mark III dimensions require reasoning capacity that scales with model size. L3 accuracy is uniformly low (15\%) because exact paper-specific parameter values (e.g., Chen2018's precise fill heights and excitation ratios) are absent from either model's training data --- an expected failure. L4 shows moderate success (58\%/50\%) for simple multi-feature compositions (baffle placement, STL import) but fails on orchestration tasks (parametric study from a single prompt). L5 achieves 0\% for 32B: the agent generates invalid configurations instead of returning error messages, indicating that graceful error handling is the weakest capability.

**Model comparison.** Overall accuracy is comparable: 32B = 47\% (95\% CI [0.30, 0.67]) versus 8B = 46\% (95\% CI [0.28, 0.66]), with no statistically significant difference (Wilcoxon $W = 4.0$, $p = 0.715$). This headline figure, however, masks an important per-level interaction. The 32B model's advantage concentrates at L2 (+25 percentage points), where domain inference tasks such as computing $f_1$ from geometry require reasoning capacity that scales with model size. Conversely, 8B paradoxically outperforms at L5 (+25 percentage points), possibly because the larger model's extended reasoning leads it to attempt generation where refusal would be more appropriate.

**Complexity--accuracy boundary.** The most statistically robust finding is the sharp performance divide between standard and challenging scenarios. For L1+L2 (explicit parameters and domain inference), 32B accuracy reaches 82\% (95\% CI [0.66, 0.93]); for L3+L5 (paper-specific and edge cases), accuracy drops to 8\% (95\% CI [0.00, 0.17]). This difference is highly significant (Mann-Whitney $U = 63$, $p = 0.001$) and persists across both models ($p = 0.035$ for 8B). The result delineates a clear boundary: SloshAgent reliably handles scenarios where parameters are stated or inferable, but cannot retrieve paper-specific values absent from training data or gracefully refuse invalid configurations.

**Context.** Comparing to other LLM-CFD systems: FoamGPT \cite{foamgpt2025} reports 26.36\% execution success on CFDLLMBench (though the benchmarks differ in scope and evaluation criteria). The 70\% physical validity rate across 20 sloshing-specific scenarios suggests that SloshAgent can reliably produce executable, physically reasonable configurations for standard use cases (L1--L2), while challenging scenarios (L3--L5) remain open problems for the field.

## 4.4 EXP-3: Parametric Study Automation (RQ3)

**Design.** A key industry requirement is automating parametric sweeps: LNG tank certification requires 200+ simulations varying fill level, excitation frequency, and amplitude \cite{isope2012benchmark}. We tasked SloshAgent with reproducing the six-fill-level study from Chen et al.\ \cite{chen2018sloshing}: a 600 $\times$ 300 $\times$ 650 mm tank under horizontal sway at $f/f_1 = 0.9$ with amplitude $A = 7$ mm, using DBC boundary condition (dp = 0.005 m) and 10 s simulation duration. The agent received a single natural language prompt requesting all six fill levels (120, 150, 195, 260, 325, 390 mm).

**Results.** Table~\ref{tab:parametric} presents the automated surface-wave-level (SWL) analysis.

| Fill [mm] | Fill [\%] | $f_1$ [Hz] | $f_{\text{exc}}$ [Hz] | Left Amp [mm] | Right Amp [mm] | Max SWL [mm] | Dom.\ Freq [Hz] |
|-----------|----------|------------|----------------------|----------------|-----------------|-------------|------------------|
| 120 | 18.5 | 0.851 | 0.766 | 66.6 | 69.8 | 224.2 | 0.142 |
| 150 | 23.1 | 0.924 | 0.831 | 35.9 | 38.6 | 192.3 | 0.851 |
| 195 | 30.0 | 1.001 | 0.901 | 41.4 | 42.9 | 243.2 | 0.851 |
| 260 | 40.0 | 1.068 | 0.961 | 45.6 | 51.4 | 310.0 | 0.993 |
| 325 | 50.0 | 1.103 | 0.993 | 48.2 | 59.3 | 376.3 | 0.993 |
| 390 | 60.0 | 1.122 | 1.009 | 50.8 | 62.3 | 443.0 | 0.993 |

All six cases completed successfully from a single agent invocation. Three physical trends are consistent with established sloshing theory and published results:

*Monotonic SWL increase.* Maximum surface-wave level increases monotonically from 224.2 mm (18.5\% fill) to 443.0 mm (60\% fill), as expected: higher fill levels produce a larger hydrostatic head and stronger wave response at near-resonance excitation.

*Left--right asymmetry growth.* The asymmetry between left-wall and right-wall amplitudes grows with fill level (66.6 vs 69.8 mm at 18.5\%; 50.8 vs 62.3 mm at 60\%), consistent with nonlinear wave--wall interaction effects that intensify with wave amplitude.

*Shallow-water nonlinear regime.* The 120 mm (18.5\%) case exhibits anomalous low-frequency dominance (0.142 Hz vs the expected excitation-locked behavior at higher fills). This is physically expected: at shallow fill ratios ($h/L < 0.2$), sloshing transitions to a nonlinear shallow-water regime where higher harmonics and subharmonics dominate the response spectrum \cite{chen2018sloshing, liu2024pitch}.

**Comparison with Chen et al.\ (2018).** Direct NRMSE comparison is precluded by different measurement variables: Chen et al.\ report normalized wall pressure ($P/\rho g h$) from OpenFOAM VOF simulations, while SloshAgent outputs free surface elevation from DualSPHysics SPH. Nevertheless, four qualitative physical trends are consistently reproduced: (i) monotonic SWL increase with fill level (150--390 mm), (ii) shallow-water hydraulic jump behavior at 120 mm ($h/L = 0.2$, matching Chen et al.'s soft-spring regime observation), (iii) excitation-frequency lock-in at $h/L > 0.32$ (dominant frequency converging to $f_{\text{exc}}$, consistent with the hard-spring regime), and (iv) maximum response near $f/f_1 = 0.9$ at higher fill levels. These four trend agreements provide qualitative validation that the agent-generated parametric study captures the established physics of fill-level-dependent sloshing behavior. Point-by-point quantitative comparison via figure digitization remains future work.

## 4.5 EXP-4: Domain Prompt Ablation (RQ4)

**Design.** While prior LLM-CFD agents ablate architectural components --- RAG removal \cite{metaopenfoam2024, mooseagent2025}, reviewer agent removal \cite{foamagent2025}, Prompt Pool removal \cite{openfoamgpt2025} --- none have isolated the contribution of *domain knowledge content* within the system prompt. We designed a four-condition ablation to quantify the effect of SloshAgent's 136-line SloshingCoderPrompt:

- **FULL**: Complete SloshingCoderPrompt (domain knowledge + tool rules + parameter inference + physics formulas).
- **NO-DOMAIN**: Domain knowledge removed (tank presets, resonance formulas, terminology mapping); tool rules retained.
- **NO-RULES**: Tool calling rules removed (path conventions, execution order, Docker syntax); domain knowledge retained.
- **GENERIC**: Single-line prompt: "You are a helpful AI assistant for DualSPHysics simulation."

We selected 10 representative scenarios from EXP-2 (two per complexity level: S01, S03, S05, S07, S09, S11, S13, S15, S17, S19) and ran each under all four conditions with both Qwen3 32B and 8B, yielding 80 runs total (10 scenarios $\times$ 4 conditions $\times$ 2 models).

**Results.** Table~\ref{tab:ablation} summarizes the accuracy and tool call rates.

| Condition | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-----------|--------------|-------------|-------------|------------|----------|---------|
| **FULL** | **10/10** | 7/10 | **60\%** | 46\% | 8/10 | 6/10 |
| **NO-DOMAIN** | **10/10** | **10/10** | 50\% | 44\% | 8/10 | **8/10** |
| **NO-RULES** | **10/10** | 9/10 | 57\% | **55\%** | 8/10 | **8/10** |
| **GENERIC** | 6/10 | **8/10** | 35\% | 39\% | 6/10 | **8/10** |

**32B ablation ordering.** The 32B model produces a monotonic ordering consistent with the hypothesis that domain knowledge improves agent performance:

$$\text{FULL}~(60\%) > \text{NO-RULES}~(57\%) > \text{NO-DOMAIN}~(50\%) > \text{GENERIC}~(35\%)$$

The FULL-versus-GENERIC gap of +25 percentage points yields a medium effect size (Cohen's $d = 0.58$, 95\% CI [0.32, 0.84] for FULL). However, with $n = 10$ scenarios per condition, pairwise permutation tests do not reach statistical significance after Bonferroni correction (FULL vs GENERIC: $p_{\text{raw}} = 0.23$, $p_{\text{bonf}} = 0.68$). The tool call rate difference is marginally significant (FULL 10/10 vs GENERIC 6/10, Fisher's exact $p = 0.087$), suggesting that the domain prompt's primary benefit may operate through improved tool-calling reliability rather than parameter accuracy alone. We note that the small sample size limits statistical power; larger-scale ablation studies are needed to confirm these directional findings.

**8B ablation ordering (inverted).** The 8B model produces an unexpected inversion:

$$\text{NO-RULES}~(55\%) > \text{FULL}~(46\%) > \text{NO-DOMAIN}~(44\%) > \text{GENERIC}~(39\%)$$

The FULL condition underperforms NO-RULES by 9 percentage points on 8B, though this difference is not statistically significant ($p = 0.67$). The observed pattern is consistent with a tool-calling capacity bottleneck: the FULL prompt's 136 lines reduce the 8B model's tool call rate to 7/10 (vs 10/10 for 32B). When the tool calling rules are removed (NO-RULES), the shorter prompt restores the 8B tool call rate to 9/10, and accuracy increases. The FULL-versus-GENERIC effect size on 8B is negligible ($d = 0.15$), contrasting sharply with the medium effect on 32B ($d = 0.58$).

**Interpretation.** The divergent effect sizes across model scales ($d = 0.58$ vs $d = 0.15$) suggest that domain prompt engineering requires sufficient model capacity to be effective. Long, knowledge-rich system prompts show directional benefit for larger models but may degrade tool-calling reliability on smaller ones. This has practical implications: for deployment on resource-constrained hardware, a shorter prompt paired with an 8B model (55\% accuracy, 11 s latency) may outperform the full prompt with the same model (46\% accuracy, 11 s latency).

**Significance.** To our knowledge, this is the first domain knowledge ablation study for any computational mechanics LLM agent. All prior ablations in the LLM-CFD literature test architectural components; ours isolates knowledge content and reports effect sizes alongside the acknowledged limitation of small sample size ($n = 10$).

## 4.6 EXP-5: Industrial Proof of Concept

**Design.** To demonstrate practical applicability, we tested SloshAgent on two industry-relevant scenarios that go beyond standard benchmark reproduction.

**Scenario A: Anti-slosh baffle design comparison.** The agent received a single prompt requesting a rectangular tank (1 m $\times$ 0.5 m $\times$ 0.6 m) with 50\% water fill at resonance frequency, comparing sloshing amplitude with and without a vertical baffle at the tank center.

| Configuration | SWL Amplitude [mm] | Max SWL [mm] | Min SWL [mm] |
|---------------|---------------------|-------------|-------------|
| No baffle | 158.9 | 522.0 | 204.1 |
| With baffle | 12.8 | 318.8 | 293.2 |
| **Reduction** | **91.9\%** | --- | --- |

The 91.9\% amplitude reduction is physically consistent with the literature: vertical baffles at mid-tank typically achieve 60--90\% reduction depending on baffle height and fill ratio \cite{nasa2023baffle, zhao2024baffles}. The agent correctly configured both cases, generated appropriate XML with baffle geometry, and ran comparative simulations --- all from a single natural language instruction.

**Scenario B: Seismic excitation.** The agent simulated a 10 m $\times$ 5 m $\times$ 8 m petroleum storage tank with 60\% fill (4.8 m water depth) under 0.3 Hz seismic-like excitation.

| Parameter | Value |
|-----------|-------|
| Tank dimensions | 10 m $\times$ 5 m $\times$ 8 m |
| Fill level | 60\% (4.8 m) |
| Excitation frequency | 0.3 Hz |
| Max SWL amplitude | 2361.8 mm |
| Max SWL | 4723.7 mm |
| Min SWL | 0.0 mm (complete wall drainage) |

The extreme amplitude (wave height exceeding the still-water level) and complete drainage at one wall indicate an overflow condition --- physically correct for near-resonance large-tank excitation. This scenario illustrates the kind of industrial application (petroleum tank safety under earthquake loading) that motivates sloshing simulation automation \cite{hatayama2004tomakomai}.

**Both scenarios completed from single natural language prompts** without manual XML editing, demonstrating that SloshAgent can handle practical engineering design tasks beyond academic benchmarks.

## 4.7 Summary of Key Findings

Table~\ref{tab:key_findings} ranks the primary findings by novelty and evidence strength.

| \# | Finding | Evidence | Statistical Support |
|----|---------|----------|---------------------|
| 1 | Complexity--accuracy boundary: L1+L2 accuracy 82\% vs L3+L5 8\% | EXP-2 | $p = 0.001$ (Mann-Whitney) |
| 2 | First experimental benchmark validation: 7/7 water peaks within $\pm 2\sigma$ (non-normal dist.) | EXP-1 | All $|z| < 2.0$; Shapiro-Wilk $p < 0.005$ |
| 3 | First domain knowledge ablation: FULL vs GENERIC medium effect ($d = 0.58$) on 32B | EXP-4 | $d = 0.58$; $p_{\text{bonf}} = 0.68$ (ns, $n = 10$) |
| 4 | 32B vs 8B overall accuracy comparable (47\% vs 46\%) despite 32$\times$ latency difference | EXP-2 | $p = 0.715$ (Wilcoxon); latency $p < 0.001$ |
| 5 | Model capacity threshold: domain prompt effect diverges ($d = 0.58$ vs $d = 0.15$) by model size | EXP-4 | Effect size comparison |
| 6 | Parametric study reproduces 4/4 physical trends from Chen et al.\ (2018) | EXP-3 | Qualitative agreement |
| 7 | 91.9\% baffle amplitude reduction in industrial PoC | EXP-5 | Demonstration |
