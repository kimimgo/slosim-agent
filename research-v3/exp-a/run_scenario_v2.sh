#!/bin/bash
# =============================================================================
# EXP-A v2: Cross-model 단일 시나리오 실행 (h100 대응)
# Usage: ./run_scenario_v2.sh <scenario> <model> <trial> [ollama_port]
# Example: ./run_scenario_v2.sh S01 llama3.3:70b-instruct-q4_K_M 1 11434
# =============================================================================
set -euo pipefail

SCENARIO="${1:?Usage: $0 <scenario> <model> <trial> [ollama_port]}"
MODEL="${2:?Missing model}"
TRIAL="${3:?Missing trial number}"
OLLAMA_PORT="${4:-11434}"

export LOCAL_ENDPOINT="http://localhost:${OLLAMA_PORT}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"
AGENT_DIR="$(dirname "$BASE_DIR")"
PROMPT_FILE="${SCRIPT_DIR}/prompts/${SCENARIO}.txt"

# 모델 이름 → config 파일 매핑
MODEL_SAFE=$(echo "$MODEL" | tr ':/' '_')
CONFIG_DIR="${SCRIPT_DIR}/configs"
CONFIG_SRC=""

case "$MODEL" in
    llama3.3:70b*) CONFIG_SRC="${CONFIG_DIR}/config_llama33_70b.json" ;;
    llama3.1:8b*)  CONFIG_SRC="${CONFIG_DIR}/config_llama31_8b.json" ;;
    gemma3:27b*)   CONFIG_SRC="${CONFIG_DIR}/config_gemma3_27b.json" ;;
    qwen3:14b*)    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_14b.json" ;;
    qwen3:32b*)    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_32b.json" ;;
    qwen3*|*latest*|*8b*) CONFIG_SRC="${CONFIG_DIR}/config_qwen3_8b.json" ;;
    *) echo "WARN: No config for model ${MODEL}, using default" ;;
esac

OUTDIR="${SCRIPT_DIR}/results/${SCENARIO}_${MODEL_SAFE}_trial${TRIAL}"

if [ ! -f "$PROMPT_FILE" ]; then
    echo "ERROR: Prompt file not found: $PROMPT_FILE"
    exit 1
fi

PROMPT=$(cat "$PROMPT_FILE")
echo "========================================"
echo "EXP-A v2 Cross-Model Runner"
echo "========================================"
echo "Scenario:  $SCENARIO"
echo "Model:     $MODEL"
echo "Trial:     $TRIAL"
echo "Ollama:    localhost:${OLLAMA_PORT}"
echo "Output:    $OUTDIR"
echo "========================================"

mkdir -p "$OUTDIR"

# config 적용
if [ -n "${CONFIG_SRC}" ] && [ -f "$CONFIG_SRC" ]; then
    cp "$CONFIG_SRC" "${OUTDIR}/.opencode.json"
    echo "[$(date +%H:%M:%S)] Applied config: $(basename $CONFIG_SRC)"
fi

# Ollama 모델 사전 확인
echo "[$(date +%H:%M:%S)] Checking model on port ${OLLAMA_PORT}..."
curl -s "http://localhost:${OLLAMA_PORT}/api/tags" | python3 -c "
import json, sys
data = json.load(sys.stdin)
models = [m['name'] for m in data.get('models', [])]
target = '${MODEL}'
if target not in models:
    print(f'WARNING: {target} not found. Available: {models[:5]}...')
else:
    print(f'OK: {target} found')
" 2>/dev/null || echo "WARN: Could not verify model"

# 이전 시뮬레이션 격리
SIMDIR="${AGENT_DIR}/simulations_${MODEL_SAFE}_${SCENARIO}_t${TRIAL}"
mkdir -p "$SIMDIR"

START_TIME=$(date +%s)
echo "[$(date +%H:%M:%S)] Starting ${SCENARIO} with ${MODEL}..."

# slosim-agent 실행 (config를 출력 디렉토리에서 읽도록)
timeout 3600 "${AGENT_DIR}/slosim" \
    -c "${OUTDIR}" \
    -p "$PROMPT" \
    -q \
    -f json \
    > "${OUTDIR}/agent_output.json" 2> "${OUTDIR}/agent_stderr.log"

EXIT_CODE=$?
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo "[$(date +%H:%M:%S)] Completed in ${DURATION}s (exit: ${EXIT_CODE})"

# 메타데이터 저장
cat > "${OUTDIR}/metadata.json" << METAEOF
{
    "scenario": "${SCENARIO}",
    "model": "${MODEL}",
    "trial": ${TRIAL},
    "ollama_port": ${OLLAMA_PORT},
    "exit_code": ${EXIT_CODE},
    "duration_seconds": ${DURATION},
    "start_time": "$(date -d @${START_TIME} +%Y-%m-%dT%H:%M:%S)",
    "end_time": "$(date -d @${END_TIME} +%Y-%m-%dT%H:%M:%S)",
    "hostname": "$(hostname)",
    "prompt_file": "${PROMPT_FILE}"
}
METAEOF

# 생성된 XML 복사
find "${OUTDIR}" -name "*.xml" -exec cp {} "${OUTDIR}/" \; 2>/dev/null || true
# simulations 디렉토리에서도 복사
if [ -d "${AGENT_DIR}/simulations" ]; then
    find "${AGENT_DIR}/simulations" -name "*.xml" -exec cp {} "${OUTDIR}/" \; 2>/dev/null || true
fi

# Ollama 모델 언로드
curl -s "http://localhost:${OLLAMA_PORT}/api/generate" \
    -d "{\"model\": \"${MODEL}\", \"keep_alive\": 0}" > /dev/null 2>&1 || true

echo "[$(date +%H:%M:%S)] Done: ${OUTDIR}/"
