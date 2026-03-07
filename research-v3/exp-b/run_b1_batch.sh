#!/bin/bash
# =============================================================================
# EXP-B B1 배치 — Generic Prompt + Tools
# 3 시나리오 × 2 모델 = 6 실행
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCENARIOS=(S01 S02 S03 S04 S05 S06 S07 S08 S09 S10)
MODELS=("qwen3:32b" "qwen3:latest")

echo "========================================"
echo "EXP-B B1 Batch (−DomainPrompt): $(date)"
echo "========================================"

for MODEL in "${MODELS[@]}"; do
    MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
    for S in "${SCENARIOS[@]}"; do
        OUTDIR="${SCRIPT_DIR}/results/B1_${S}_${MODEL_SAFE}"
        if [ -d "$OUTDIR" ] && [ -f "$OUTDIR/metadata.json" ]; then
            echo "[$(date +%H:%M:%S)] SKIP B1 $S $MODEL_SAFE (exists)"
            continue
        fi
        echo ""
        echo "[$(date +%H:%M:%S)] ====== B1 $S $MODEL ======"
        bash "${SCRIPT_DIR}/run_b1_noprompt.sh" "$S" "$MODEL" 0
        echo "[$(date +%H:%M:%S)] Done B1 $S $MODEL"
        sleep 10
    done
done

echo ""
echo "========================================"
echo "B1 Batch complete: $(date)"
echo "========================================"
