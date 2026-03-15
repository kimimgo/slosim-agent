# Paper 1 Outline: SlosimAgent Verification

## Abstract (~250 words)
- Problem: CFD 비전문가의 SPH 시뮬레이션 접근 불가
- Gap: SPH agent 0편, Local SLM 0편
- Solution: SlosimAgent (15 tools, Qwen3, DualSPHysics v5.4)
- Results: M-A3 82.2%, 8B>70B cross-model, factorial ablation
- Significance: tool architecture > model size

## 1. Introduction
### 1.1 Background
- 슬로싱 공학적 중요성 (LNG, 항공, 지진)
- SPH 적합성 (meshless, violent free-surface)
- LLM-CFD agent 동향 (7종, 2024-2025)

### 1.2 Research Gaps
- G1: 비전문가 접근성 = 0
- G2: SPH solver agent = 0
- G3: Local SLM agent = 0

### 1.3 Contributions
1. 최초의 SPH LLM 에이전트
2. DualSPHysics 전체 튜토리얼 커버
3. eval-harness 벤치마크 프레임워크
4. Factorial ablation으로 도구 설계 기여 분리
5. Cross-model 일반화 (prompt > size)

## 2. Related Work
### 2.1 LLM-Based CFD Agents (7종)
### 2.2 SPH Sloshing Simulation
### 2.3 LLM Tool Use and Agent Architecture

## 3. System Architecture
### 3.1 Overview (agent loop, Go + BubbleTea)
### 3.2 Tool Design (15 DualSPHysics tools)
- xml_generator, gencase, solver, partvtk, measuretool, ...
- Tool interface: BaseTool.Info() + Run()
### 3.3 Domain-Specific Prompt Engineering
- SloshingCoderPrompt: 파라미터 자동결정 규칙
- Tool call 순서 강제 (A/B/C pathways)
### 3.4 DualSPHysics Tutorial Coverage
- 전체 튜토리얼 시나리오 매핑

## 4. Evaluation Framework
### 4.1 M-A3 Parameter Fidelity Metric (8 params)
### 4.2 eval-harness Benchmark Design
### 4.3 Scenario Design (10 scenarios, 3 tiers)

## 5. Experiments
### 5.1 EXP-A: Parameter Fidelity
- 10 scenarios × 2 models × 3 trials
- v3→v4 tool fix impact (+12.8pp)
- Table 3: per-scenario scores

### 5.2 EXP-B: Factorial Ablation
- 2×2 (prompt × tool)
- B0=67%, B1=0%, B2=46.1%, B4=0%
- Hierarchical dependency: prompt is prerequisite
- Table 4: factorial effects

### 5.3 Cross-Model Generalization
- 5 models (Qwen3 32B/14B/8B, LLaMA 70B/8B)
- 8B(78.5%) > 70B(51.4%)
- JSON compat layer discovery
- Table 6: cross-model comparison
- Figure: size vs score scatter

### 5.4 8B vs 32B Physics Equivalence
- 동일 시나리오 생성 XML → GPU 시뮬레이션 → 결과 비교
- 물리 결과 차이 ≈ 0 (동일 XML → 동일 결과)

## 6. Discussion
### 6.1 Tool Design as Bottleneck
### 6.2 Prompt Optimization > Model Size
### 6.3 Threshold Effect (LLaMA 8B = 0%)
### 6.4 Limitations
- "실험 대비 물리 검증은 companion paper에서 다룸"

## 7. Conclusion

## References
