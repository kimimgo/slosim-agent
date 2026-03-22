#!/bin/bash
# =============================================================================
# H100 4-GPU 병렬 Cross-Model 실험 러너
# 4개 모델을 4개 GPU에 동시 할당하여 최대 처리량 달성
#
# Usage: nohup bash run_h100_parallel.sh > /tmp/h100_exp.log 2>&1 &
#
# GPU 할당:
#   GPU 0: llama3.3:70b-instruct-q4_K_M (42.5GB)
#   GPU 1: gemma3:27b (17.4GB)
#   GPU 2: llama3.1:8b-instruct-q8_0 (8.5GB)
#   GPU 3: qwen3:14b (9.3GB)
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
AGENT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

# 4개 모델 정의 (모델명, Ollama 포트, GPU ID)
declare -A MODELS
MODELS[0]="llama3.3:70b-instruct-q4_K_M"
MODELS[1]="gemma3:27b"
MODELS[2]="llama3.1:8b-instruct-q8_0"
MODELS[3]="qwen3:14b"

PORTS=(11434 11435 11436 11437)
SCENARIOS=(S01 S02 S03 S04 S05 S06 S07 S08 S09 S10)
TRIAL=1  # temperature=0 → 1 trial로 충분 (deterministic)

LOG_DIR="/tmp/h100_exp_logs"
mkdir -p "$LOG_DIR"

echo "============================================================"
echo " H100 4-GPU Parallel Cross-Model Experiment"
echo " Started: $(date)"
echo "============================================================"
echo ""

# ──────────────────────────────────────────────────
# Step 1: 4개 Ollama 인스턴스 시작 (GPU별)
# ──────────────────────────────────────────────────
echo "[SETUP] Starting 4 Ollama instances..."

# 기존 Ollama 정리
pkill -f "ollama serve" 2>/dev/null || true
sleep 2

for GPU in 0 1 2 3; do
    PORT=${PORTS[$GPU]}
    echo "[SETUP] GPU ${GPU}: Ollama on port ${PORT} (model: ${MODELS[$GPU]})"

    CUDA_VISIBLE_DEVICES=$GPU OLLAMA_HOST="0.0.0.0:${PORT}" \
        nohup ollama serve > "${LOG_DIR}/ollama_gpu${GPU}.log" 2>&1 &
done

echo "[SETUP] Waiting 10s for Ollama instances to start..."
sleep 10

# Ollama 인스턴스 health check
for GPU in 0 1 2 3; do
    PORT=${PORTS[$GPU]}
    if curl -s "http://localhost:${PORT}/api/tags" > /dev/null 2>&1; then
        echo "[SETUP] GPU ${GPU}: Ollama OK (port ${PORT})"
    else
        echo "[ERROR] GPU ${GPU}: Ollama FAILED on port ${PORT}"
        exit 1
    fi
done

# ──────────────────────────────────────────────────
# Step 2: 각 Ollama에 모델 워밍업 (사전 로드)
# ──────────────────────────────────────────────────
echo ""
echo "[WARMUP] Pre-loading models..."

for GPU in 0 1 2 3; do
    PORT=${PORTS[$GPU]}
    MODEL="${MODELS[$GPU]}"
    echo "[WARMUP] GPU ${GPU}: Loading ${MODEL}..."
    curl -s "http://localhost:${PORT}/api/generate" \
        -d "{\"model\": \"${MODEL}\", \"prompt\": \"hello\", \"stream\": false}" \
        > /dev/null 2>&1 &
done
wait
echo "[WARMUP] All models loaded."

# ──────────────────────────────────────────────────
# Step 3: 4개 모델 × 10 시나리오 병렬 실행
# ──────────────────────────────────────────────────
echo ""
echo "============================================================"
echo "[RUN] Starting experiments: 4 models × 10 scenarios = 40 runs"
echo "============================================================"

# 각 GPU에서 10개 시나리오를 순차 실행 (GPU 간은 병렬)
run_model() {
    local GPU=$1
    local MODEL="${MODELS[$GPU]}"
    local PORT=${PORTS[$GPU]}
    local MODEL_SAFE=$(echo "$MODEL" | tr ':/' '_')
    local LOG="${LOG_DIR}/gpu${GPU}_${MODEL_SAFE}.log"

    echo "[GPU${GPU}] Starting ${MODEL} (10 scenarios)..." | tee -a "$LOG"

    for SCENARIO in "${SCENARIOS[@]}"; do
        RESULT_DIR="${SCRIPT_DIR}/results/${SCENARIO}_${MODEL_SAFE}_trial${TRIAL}"

        # 이미 완료된 시나리오 건너뛰기
        if [ -f "${RESULT_DIR}/metadata.json" ]; then
            DUR=$(python3 -c "import json; d=json.load(open('${RESULT_DIR}/metadata.json')); print(d.get('duration_seconds',0))" 2>/dev/null || echo "0")
            if [ "$DUR" -gt "5" ]; then
                echo "[GPU${GPU}] SKIP ${SCENARIO} — already done (${DUR}s)" | tee -a "$LOG"
                continue
            fi
        fi

        echo "[GPU${GPU}] Running ${SCENARIO}..." | tee -a "$LOG"
        START=$(date +%s)

        bash "${SCRIPT_DIR}/run_scenario_v2.sh" \
            "${SCENARIO}" "${MODEL}" "${TRIAL}" "${PORT}" \
            >> "$LOG" 2>&1

        END=$(date +%s)
        DUR=$((END - START))
        echo "[GPU${GPU}] ${SCENARIO} done (${DUR}s)" | tee -a "$LOG"

        # 시나리오 간 5초 쿨다운
        sleep 5
    done

    echo "[GPU${GPU}] ALL DONE: ${MODEL}" | tee -a "$LOG"
}

# 4개 GPU 동시 실행
for GPU in 0 1 2 3; do
    run_model $GPU &
done

echo "[RUN] 4 parallel streams launched. Waiting for completion..."
wait

# ──────────────────────────────────────────────────
# Step 4: 결과 요약
# ──────────────────────────────────────────────────
echo ""
echo "============================================================"
echo " RESULTS SUMMARY"
echo " Completed: $(date)"
echo "============================================================"

printf "%-6s" ""
for GPU in 0 1 2 3; do
    MODEL_SHORT=$(echo "${MODELS[$GPU]}" | cut -d: -f1)
    printf "  %-15s" "$MODEL_SHORT"
done
echo ""
echo "--------------------------------------------------------------"

for SCENARIO in "${SCENARIOS[@]}"; do
    printf "%-6s" "$SCENARIO"
    for GPU in 0 1 2 3; do
        MODEL_SAFE=$(echo "${MODELS[$GPU]}" | tr ':/' '_')
        META="${SCRIPT_DIR}/results/${SCENARIO}_${MODEL_SAFE}_trial${TRIAL}/metadata.json"
        if [ -f "$META" ]; then
            DUR=$(python3 -c "import json; d=json.load(open('${META}')); print(f\"{d['duration_seconds']}s\")" 2>/dev/null || echo "?")
            EXIT=$(python3 -c "import json; d=json.load(open('${META}')); print(d['exit_code'])" 2>/dev/null || echo "?")
            if [ "$EXIT" = "0" ]; then
                printf "  %-15s" "✓ ${DUR}"
            else
                printf "  %-15s" "✗ (${EXIT})"
            fi
        else
            printf "  %-15s" "—"
        fi
    done
    echo ""
done

echo ""
echo "Logs: ${LOG_DIR}/"
echo "Results: ${SCRIPT_DIR}/results/"

# Ollama 인스턴스 정리
echo "[CLEANUP] Stopping Ollama instances..."
pkill -f "ollama serve" 2>/dev/null || true

echo "DONE."
