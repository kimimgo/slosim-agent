#!/bin/bash
# =============================================================================
# pajulab 순차 배치: S04-S10 (S02,S03은 이미 실행됨/실행중)
# Usage: nohup bash run_pajulab_batch.sh > /tmp/pajulab_batch.log 2>&1 &
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
MODEL="qwen3:32b"
TRIAL=1
GPU=0
LOG="/tmp/pajulab_batch.log"

echo "========================================"
echo "pajulab Batch Runner: S04-S10 (32B trial1)"
echo "Started: $(date)"
echo "========================================"

for SCENARIO in S04 S05 S06 S07 S08 S09 S10; do
    RESULT_DIR="${SCRIPT_DIR}/results/${SCENARIO}_qwen3_32b_trial${TRIAL}"

    # 이미 완료된 시나리오 건너뛰기
    if [ -f "${RESULT_DIR}/metadata.json" ]; then
        DUR=$(python3 -c "import json; d=json.load(open('${RESULT_DIR}/metadata.json')); print(d.get('duration_seconds',0))" 2>/dev/null || echo "0")
        if [ "$DUR" -gt "5" ]; then
            echo "[$(date +%H:%M:%S)] SKIP ${SCENARIO} — already completed (${DUR}s)"
            continue
        fi
    fi

    echo ""
    echo "[$(date +%H:%M:%S)] ====== Starting ${SCENARIO} ======"
    bash "${SCRIPT_DIR}/run_scenario.sh" "${SCENARIO}" "${MODEL}" "${TRIAL}" "${GPU}" 2>&1 | tee -a "${LOG}.${SCENARIO}"

    # 결과 확인
    if [ -f "${RESULT_DIR}/metadata.json" ]; then
        DUR=$(python3 -c "import json; d=json.load(open('${RESULT_DIR}/metadata.json')); print(d.get('duration_seconds',0))" 2>/dev/null || echo "?")
        echo "[$(date +%H:%M:%S)] ${SCENARIO} completed in ${DUR}s"
    else
        echo "[$(date +%H:%M:%S)] ${SCENARIO} — no metadata (may have failed)"
    fi

    # 시나리오 간 10초 쿨다운 (Ollama 안정화)
    echo "[$(date +%H:%M:%S)] Cooling down 10s..."
    sleep 10
done

echo ""
echo "========================================"
echo "Batch completed: $(date)"
echo "========================================"

# 전체 결과 요약
echo ""
echo "=== Results Summary ==="
for SCENARIO in S04 S05 S06 S07 S08 S09 S10; do
    META="${SCRIPT_DIR}/results/${SCENARIO}_qwen3_32b_trial${TRIAL}/metadata.json"
    if [ -f "$META" ]; then
        DUR=$(python3 -c "import json; d=json.load(open('${META}')); print(d['duration_seconds'])" 2>/dev/null || echo "?")
        EXIT=$(python3 -c "import json; d=json.load(open('${META}')); print(d['exit_code'])" 2>/dev/null || echo "?")
        echo "${SCENARIO}: ${DUR}s (exit=${EXIT})"
    else
        echo "${SCENARIO}: NOT RUN"
    fi
done
