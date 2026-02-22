# Phase 4-5 Plan: Analysis + Paper Writing

## GAP vs Experiment Alignment Summary

| GAP | RQ | EXP | Status | Evidence |
|-----|----|----|--------|----------|
| GAP-1 (Domain: No sloshing AI) | RQ1 | EXP-2 | **COVERED** | 8B: 46%, 32B: 47% accuracy, L1=96% |
| GAP-2 (Knowledge: No domain ablation) | RQ4 | EXP-4 | **COVERED** | 4-condition ablation, 8B+32B comparison |
| GAP-3 (Validation: No exp benchmark) | RQ2 | EXP-1 | **PARTIAL** | Peaks in +/-2sigma, but r/NRMSE not computed |
| GAP-4 (Industry: $4.13B zero AI) | RQ3 | EXP-3+5 | **COVERED** | 6 fill levels, baffle 91.9% reduction |
| GAP-5 (Method: Particle+LLM void) | RQ3 | System | **COVERED** | 14 tools, DualSPHysics pipeline |

## Critical Fix: EXP-1 Quantitative Metrics

Paper skeleton promises "Pearson r > 0.9, NRMSE". Current EXP-1 only reports:
- Peak pressure in +/-2sigma band (qualitative)
- Need: time-series correlation (r), NRMSE for continuous validation

Action: Recompute from existing data or clearly state peak-based validation.

---

## Phase 4: Integrated Analysis (Agent Team Task)

### Task 4.1: EXP-1 Metric Strengthening
- Verify if time-series r/NRMSE can be computed from existing probe data
- If not, clearly frame as peak-based validation (still valid, ISOPE standard)
- Update Table 2 with final metrics

### Task 4.2: Unified Results Table
- Create master results table linking GAP→RQ→EXP→Key Finding
- Cross-reference all 5 experiments

### Task 4.3: Final Figure Set
- Fig 1: System architecture diagram (text-based, for LaTeX)
- Fig 2: EXP-1 SPHERIC pressure validation (existing: fig2a/b/c)
- Fig 3: EXP-3 parametric SWL (existing: fig3/fig3b)
- Fig 4: EXP-4 ablation 8B vs 32B (existing: fig4_comparison.png)
- Fig 5: EXP-5 baffle comparison (existing: fig5_baffle_comparison.png)
- Fig 6: EXP-5 seismic (existing: fig6_seismic.png)
- Table 1: Competitor comparison (from skeleton, verified)
- Table 2: SPHERIC validation metrics
- Table 3: NL->XML accuracy by level (8B+32B)
- Table 4: Chen2018 parametric results
- Table 5: Ablation results (8B+32B)
- Table 6: Industrial PoC summary

### Task 4.4: Statistical Significance
- Report 95% CI for EXP-2/4 where applicable
- Note sample sizes (n=20 for EXP-2, n=10 for EXP-4)

---

## Phase 5: Paper Writing (Agent Team Task)

### Target: 4-8 page AI4Science workshop paper (LaTeX)

### Structure (from paper_skeleton.md)
1. **Abstract** (~250 words)
2. **Introduction** (~1 page) - 1.1 Sloshing problem, 1.2 Bottleneck, 1.3 Related gap, 1.4 Contributions
3. **Related Work** (~1 page) - 2.1 Sloshing SOTA, 2.2 LLM+CFD agents, 2.3 AI4Science
4. **System Design** (~1.5 pages) - 3.1 Architecture, 3.2 Tools, 3.3 Prompt, 3.4 XML pipeline
5. **Experiments** (~2 pages) - 4.1 Setup, 4.2-4.6 EXP-1~5
6. **Discussion** (~0.5 pages) - 5.1 Limitations, 5.2 Generalization, 5.3 Industry impact
7. **Conclusion** (~0.5 pages)

### Writing Guidelines
- Academic English, CS/AI conference style
- Quantitative claims backed by experiment data
- Each claim maps to a GAP
- LaTeX format: use standard article class or workshop template
- BibTeX references from research/references.bib

---

## Agent Team Design

### Lead: cc-slosim-1 (Opus) - Orchestrator
- GAP alignment analysis
- Quality control, cross-referencing
- Git push/commit
- Final review

### Agent 1: "analyst" (Sonnet) - Data Analyst
- Task 4.1: Recompute/verify EXP-1 metrics
- Task 4.2: Unified results table
- Task 4.4: Statistical analysis
- Reads experiment JSONs, produces analysis markdown

### Agent 2: "writer" (Sonnet) - Paper Writer
- Phase 5: Write paper sections (Markdown first, LaTeX conversion last)
- Uses paper_skeleton.md as blueprint
- Uses gap_refinement_cycle3.md for narrative
- Uses experiment results for evidence

### Agent 3: "figure-gen" (Sonnet) - Figure Generator
- Task 4.3: Generate/refine final figures
- Architecture diagram (Fig 1)
- Verify existing figures meet publication quality
- Python matplotlib scripts

### Workflow
1. analyst produces integrated analysis -> writer uses it
2. figure-gen works in parallel on figures
3. writer produces draft sections -> lead reviews
4. lead commits + pushes
