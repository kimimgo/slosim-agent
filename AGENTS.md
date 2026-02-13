# slosim-agent

AI agent for sloshing simulation using DualSPHysics SPH solver.

## Architecture
- Base: OpenCode (Go + BubbleTea TUI)
- Solver: DualSPHysics v5.4 (GPU, CUDA 12.6)
- LLM: Qwen3 (local via Ollama)

## GenCase Rule
GenCase auto-appends .xml. Do NOT include .xml in path.

## XML Format  
DualSPHysics XML uses attributes only: <gravity x="0" y="0" z="-9.81" />
