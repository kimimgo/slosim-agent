# Paper 1: CS Journal — Verification

## Title
SlosimAgent: Natural-Language-Driven AI Agent for Autonomous SPH Sloshing Simulation

## Target Journal
EAAI / Engineering with Computers / Expert Systems with Applications

## Framing
V-model 좌측 — "설계대로 만들었는가?" (Verification)

## Gap × Evidence

| Gap | Evidence (Verification) |
|-----|------------------------|
| G1 비전문가 접근성 | M-A3 0%→82% + DualSPHysics 전체 튜토리얼 커버 |
| G2 SPH solver agent | 15종 도구 아키텍처 + 2×2 factorial ablation |
| G3 Local SLM 충분 | 8B>70B cross-model + 8B≡32B 물리 동등성 + eval-harness |

## Planned Sections
1. Introduction (4-gap + contribution)
2. Related Work (7 CFD agents + SPH literature)
3. System Architecture (agent loop, tool design, prompt engineering)
4. Evaluation Framework (M-A3, eval-harness, DualSPHysics tutorial coverage)
5. Experiments
   - EXP-A: M-A3 parameter fidelity (10 scenarios × 3 trials × 2 models)
   - EXP-B: 2×2 factorial ablation (prompt × tool)
   - Cross-Model: 5 models (Qwen3 32B/14B/8B, LLaMA 70B/8B)
   - 8B vs 32B physics equivalence
6. Discussion (tool design as bottleneck, prompt > size)
7. Conclusion + Limitation ("실험 검증은 별도 연구")

## Key Numbers
- EXP-A v4: 32B=82.2%, 8B=78.5%, Δ=3.8pp
- EXP-B: B0=67%, B1=0%, B2=46.1%, B4=0%
- Cross-Model: Qwen3:8B(78.5%) > LLaMA:70B(51.4%)
- Factorial: Prompt=+56.5%pp, Tool=+10.5%pp, Interaction=+21.0%pp

## TODO
- [ ] DualSPHysics 전체 튜토리얼 시나리오 매핑 (현재 10개 → 확장)
- [ ] eval-harness 벤치마크 프레임워크 설계
- [ ] 8B vs 32B 생성 케이스 물리 동등성 실험
- [ ] LaTeX 원고 작성 (Overleaf 연동)
- [ ] 2차 논문 예고 limitation 문구
