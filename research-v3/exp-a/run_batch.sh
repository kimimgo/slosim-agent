#!/bin/bash
# =============================================================================
# EXP-A: 배치 실행 스크립트
# Usage: ./run_batch.sh <model> <scenarios> <trial> [gpu_id]
# Examples:
#   ./run_batch.sh qwen3:32b "S01 S02 S03 S04 S05" 1 0   # oliveeelab
#   ./run_batch.sh qwen3:32b "S06 S07 S08 S09 S10" 1 0   # pajulab GPU0
#   ./run_batch.sh qwen3 "S01 S02 S03 S04 S05 S06 S07 S08 S09 S10" 1 1  # pajulab GPU1 (8B)
# =============================================================================
set -euo pipefail

MODEL="${1:?Usage: $0 <model> <scenarios> <trial> [gpu_id]}"
SCENARIOS="${2:?Missing scenarios (e.g., 'S01 S02 S03')}"
TRIAL="${3:?Missing trial number}"
GPU_ID="${4:-0}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TOTAL=$(echo "$SCENARIOS" | wc -w)
COUNT=0
PASS=0
FAIL=0

echo "================================================================"
echo "EXP-A Batch Run"
echo "================================================================"
echo "Model:      $MODEL"
echo "Scenarios:  $SCENARIOS"
echo "Trial:      $TRIAL"
echo "GPU:        $GPU_ID"
echo "Total:      $TOTAL scenarios"
echo "Started:    $(date)"
echo "================================================================"

BATCH_START=$(date +%s)

for S in $SCENARIOS; do
    COUNT=$((COUNT + 1))
    echo ""
    echo "[$COUNT/$TOTAL] Running $S..."
    echo "────────────────────────────────────────"

    if bash "${SCRIPT_DIR}/run_scenario.sh" "$S" "$MODEL" "$TRIAL" "$GPU_ID"; then
        PASS=$((PASS + 1))
        echo "[$COUNT/$TOTAL] $S: DONE"
    else
        FAIL=$((FAIL + 1))
        echo "[$COUNT/$TOTAL] $S: FAILED (continuing...)"
    fi
done

BATCH_END=$(date +%s)
BATCH_DURATION=$((BATCH_END - BATCH_START))

echo ""
echo "================================================================"
echo "Batch Complete"
echo "================================================================"
echo "Total: $TOTAL | Pass: $PASS | Fail: $FAIL"
echo "Duration: ${BATCH_DURATION}s ($(( BATCH_DURATION / 60 ))m)"
echo "Finished: $(date)"
echo "================================================================"

# 배치 요약 저장
MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
cat > "${SCRIPT_DIR}/results/batch_${MODEL_SAFE}_trial${TRIAL}_$(hostname).json" << EOF
{
    "model": "$MODEL",
    "trial": $TRIAL,
    "gpu_id": $GPU_ID,
    "hostname": "$(hostname)",
    "total": $TOTAL,
    "pass": $PASS,
    "fail": $FAIL,
    "duration_seconds": $BATCH_DURATION,
    "scenarios": "$(echo $SCENARIOS)"
}
EOF
