#!/bin/bash
# =============================================================================
# EXP-A P1+P2 Retest: Pitch 시나리오만 재실행
# Scenarios: S01, S04, S05, S08, S09 (motion_type=pitch)
# Models: qwen3:32b, qwen3:latest (8B)
# Trial: 1 only (temperature=0 → deterministic)
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCENARIOS=(S01 S04 S05 S08 S09)
MODELS=("qwen3:32b" "qwen3:latest")
TRIAL=1

echo "========================================"
echo "EXP-A P1+P2 Pitch Retest"
echo "Scenarios: ${SCENARIOS[*]}"
echo "Models: ${MODELS[*]}"
echo "========================================"

TOTAL=$((${#SCENARIOS[@]} * ${#MODELS[@]}))
COUNT=0

for MODEL in "${MODELS[@]}"; do
    for SCENARIO in "${SCENARIOS[@]}"; do
        COUNT=$((COUNT + 1))
        MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
        OUTDIR="${SCRIPT_DIR}/results_v4/${SCENARIO}_${MODEL_SAFE}"

        # Skip if already done
        if [ -f "${OUTDIR}/metadata.json" ]; then
            echo "[${COUNT}/${TOTAL}] SKIP ${SCENARIO} ${MODEL} (already done)"
            continue
        fi

        echo ""
        echo "[${COUNT}/${TOTAL}] Running ${SCENARIO} with ${MODEL}..."
        echo "========================================"

        "${SCRIPT_DIR}/run_scenario.sh" "${SCENARIO}" "${MODEL}" "${TRIAL}" 0 || {
            echo "WARN: ${SCENARIO} ${MODEL} failed (exit $?), continuing..."
        }

        # Move results to v4 directory
        MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
        SRC="${SCRIPT_DIR}/results/${SCENARIO}_${MODEL_SAFE}_trial${TRIAL}"
        if [ -d "$SRC" ]; then
            mkdir -p "${SCRIPT_DIR}/results_v4"
            mv "$SRC" "$OUTDIR"
            echo "Moved results to ${OUTDIR}"
        fi

        # Wait between runs to let Ollama clean up
        sleep 5
    done
done

echo ""
echo "========================================"
echo "All ${TOTAL} runs complete!"
echo "Results in: ${SCRIPT_DIR}/results_v4/"
echo "========================================"
