#!/bin/bash
# =============================================================================
# EXP-A: 단일 시나리오 실행 스크립트
# Usage: ./run_scenario.sh <scenario> <model> <trial> [gpu_id]
# Example: ./run_scenario.sh S01 qwen3:32b 1 0
# =============================================================================
set -euo pipefail

SCENARIO="${1:?Usage: $0 <scenario> <model> <trial> [gpu_id]}"
MODEL="${2:?Missing model (e.g., qwen3:32b or qwen3)}"
TRIAL="${3:?Missing trial number (1, 2, 3)}"
GPU_ID="${4:-0}"

# 경로 설정
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"
AGENT_DIR="$(dirname "$BASE_DIR")"
PROMPT_FILE="${SCRIPT_DIR}/prompts/${SCENARIO}.txt"

# 모델 이름에서 특수문자 제거 (디렉토리용)
MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
OUTDIR="${SCRIPT_DIR}/results/${SCENARIO}_${MODEL_SAFE}_trial${TRIAL}"

# 프롬프트 파일 확인
if [ ! -f "$PROMPT_FILE" ]; then
    echo "ERROR: Prompt file not found: $PROMPT_FILE"
    exit 1
fi

PROMPT=$(cat "$PROMPT_FILE")
echo "========================================"
echo "EXP-A Scenario Runner"
echo "========================================"
echo "Scenario:  $SCENARIO"
echo "Model:     $MODEL"
echo "Trial:     $TRIAL"
echo "GPU:       $GPU_ID"
echo "Prompt:    $PROMPT"
echo "Output:    $OUTDIR"
echo "========================================"

# 출력 디렉토리 생성
mkdir -p "$OUTDIR"

# GPU 설정
export CUDA_VISIBLE_DEVICES="$GPU_ID"

# 모델에 맞는 설정 파일 적용
CONFIG_DIR="${SCRIPT_DIR}/configs"
if [[ "$MODEL" == *"32b"* ]]; then
    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_32b.json"
elif [[ "$MODEL" == *"8b"* ]] || [[ "$MODEL" == "qwen3" ]] || [[ "$MODEL" == *"latest"* ]]; then
    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_8b.json"
fi
if [ -n "${CONFIG_SRC:-}" ] && [ -f "$CONFIG_SRC" ]; then
    cp "$CONFIG_SRC" "$HOME/.opencode.json"
    echo "[$(date +%H:%M:%S)] Applied config: $(basename $CONFIG_SRC)"
fi

# Ollama 모델 사전 로드 확인
echo "[$(date +%H:%M:%S)] Checking Ollama model availability..."
if ! curl -s http://localhost:11434/api/tags | python3 -c "
import json, sys
data = json.load(sys.stdin)
models = [m['name'] for m in data.get('models', [])]
if '${MODEL}' not in models:
    print(f'Model ${MODEL} not found. Available: {models}')
    sys.exit(1)
print(f'Model ${MODEL} found.')
" 2>/dev/null; then
    echo "WARNING: Could not verify model. Proceeding anyway..."
fi

# 시작 시간 기록
START_TIME=$(date +%s)
echo "[$(date +%H:%M:%S)] Starting scenario ${SCENARIO}..."

# slosim-agent 실행 (타임아웃 1시간)
# -c 로 작업 디렉토리 지정, -p 로 프롬프트, -q 조용히, -f json 출력
timeout 3600 "${AGENT_DIR}/slosim" \
    -c "${AGENT_DIR}" \
    -p "$PROMPT" \
    -q \
    -f json \
    > "${OUTDIR}/agent_output.json" 2> "${OUTDIR}/agent_stderr.log" || true

EXIT_CODE=$?
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo "[$(date +%H:%M:%S)] Completed in ${DURATION}s (exit: ${EXIT_CODE})"

# 결과 메타데이터 저장
cat > "${OUTDIR}/metadata.json" << METAEOF
{
    "scenario": "${SCENARIO}",
    "model": "${MODEL}",
    "trial": ${TRIAL},
    "gpu_id": ${GPU_ID},
    "exit_code": ${EXIT_CODE},
    "duration_seconds": ${DURATION},
    "start_time": "$(date -d @${START_TIME} +%Y-%m-%dT%H:%M:%S)",
    "end_time": "$(date -d @${END_TIME} +%Y-%m-%dT%H:%M:%S)",
    "hostname": "$(hostname)",
    "prompt_file": "${PROMPT_FILE}"
}
METAEOF

# Ollama 모델 언로드 (VRAM 해제)
echo "[$(date +%H:%M:%S)] Unloading model from VRAM..."
curl -s http://localhost:11434/api/generate \
    -d "{\"model\": \"${MODEL}\", \"keep_alive\": 0}" > /dev/null 2>&1 || true

echo "[$(date +%H:%M:%S)] Done. Results in ${OUTDIR}/"
