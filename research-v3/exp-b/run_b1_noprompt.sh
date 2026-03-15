#!/bin/bash
# =============================================================================
# EXP-B: B1 (−DomainPrompt) — 도메인 프롬프트 X, 도구 O
# slosim-b1 바이너리 사용 (SLOSIM_GENERIC_PROMPT=1)
# Generic CoderPrompt + 도구 18개 그대로 제공
# =============================================================================
set -euo pipefail

SCENARIO="${1:?Usage: $0 <scenario> <model> [gpu_id]}"
MODEL="${2:?Missing model (e.g., qwen3:32b or qwen3)}"
GPU_ID="${3:-0}"

export LOCAL_ENDPOINT="${LOCAL_ENDPOINT:-http://localhost:11434}"
export SLOSIM_GENERIC_PROMPT=1  # ← B1 핵심: 도메인 프롬프트 비활성화

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BASE_DIR="$(dirname "$SCRIPT_DIR")"
AGENT_DIR="$(dirname "$BASE_DIR")"
PROMPT_FILE="${BASE_DIR}/exp-a/prompts/${SCENARIO}.txt"

MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')
OUTDIR="${SCRIPT_DIR}/results/B1_${SCENARIO}_${MODEL_SAFE}"

if [ ! -f "$PROMPT_FILE" ]; then
    echo "ERROR: Prompt file not found: $PROMPT_FILE"
    exit 1
fi

PROMPT=$(cat "$PROMPT_FILE")
echo "========================================"
echo "EXP-B B1 (−DomainPrompt)"
echo "========================================"
echo "Scenario:  $SCENARIO"
echo "Model:     $MODEL"
echo "GPU:       $GPU_ID"
echo "Generic:   SLOSIM_GENERIC_PROMPT=$SLOSIM_GENERIC_PROMPT"
echo "Prompt:    $PROMPT"
echo "Output:    $OUTDIR"
echo "========================================"

mkdir -p "$OUTDIR"
export CUDA_VISIBLE_DEVICES="$GPU_ID"

# 모델 config 적용
CONFIG_DIR="${BASE_DIR}/exp-a/configs"
if [[ "$MODEL" == *"32b"* ]]; then
    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_32b.json"
else
    CONFIG_SRC="${CONFIG_DIR}/config_qwen3_8b.json"
fi
if [ -n "${CONFIG_SRC:-}" ] && [ -f "$CONFIG_SRC" ]; then
    cp "$CONFIG_SRC" "$HOME/.opencode.json"
    echo "[$(date +%H:%M:%S)] Applied config: $(basename $CONFIG_SRC)"
fi

# Ollama 확인
echo "[$(date +%H:%M:%S)] Checking Ollama model availability..."
curl -s http://localhost:11434/api/tags | python3 -c "
import json, sys
data = json.load(sys.stdin)
models = [m['name'] for m in data.get('models', [])]
if '${MODEL}' not in models:
    print(f'WARNING: Model ${MODEL} not found.')
else:
    print(f'Model ${MODEL} found.')
" 2>/dev/null || echo "WARNING: Could not verify model."

START_TIME=$(date +%s)
echo "[$(date +%H:%M:%S)] Starting B1 scenario ${SCENARIO}..."

# 시뮬레이션 디렉토리 격리
SIMDIR="${AGENT_DIR}/simulations"
if [ -d "$SIMDIR" ] && [ "$(ls -A "$SIMDIR" 2>/dev/null)" ]; then
    BACKUP="${SIMDIR}_backup_$(date +%s)"
    mv "$SIMDIR" "$BACKUP"
    echo "[$(date +%H:%M:%S)] Moved previous simulations to $(basename $BACKUP)"
fi
mkdir -p "$SIMDIR"

# slosim-b1 바이너리 실행 (generic prompt)
SLOSIM_BIN="${AGENT_DIR}/slosim-b1"
if [ ! -f "$SLOSIM_BIN" ]; then
    SLOSIM_BIN="${AGENT_DIR}/slosim"
    echo "[$(date +%H:%M:%S)] WARNING: slosim-b1 not found, using slosim"
fi

timeout 3600 "$SLOSIM_BIN" \
    -c "${AGENT_DIR}" \
    -p "$PROMPT" \
    -q \
    -f json \
    > "${OUTDIR}/agent_output.json" 2> "${OUTDIR}/agent_stderr.log"

EXIT_CODE=$?
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo "[$(date +%H:%M:%S)] Completed in ${DURATION}s (exit: ${EXIT_CODE})"

# 메타데이터
cat > "${OUTDIR}/metadata.json" << METAEOF
{
    "scenario": "${SCENARIO}",
    "model": "${MODEL}",
    "condition": "B1_no_domain_prompt",
    "exit_code": ${EXIT_CODE},
    "duration_seconds": ${DURATION},
    "start_time": "$(date -d @${START_TIME} +%Y-%m-%dT%H:%M:%S)",
    "end_time": "$(date -d @${END_TIME} +%Y-%m-%dT%H:%M:%S)",
    "hostname": "$(hostname)",
    "generic_prompt": true,
    "has_tools": true
}
METAEOF

# 산출물 보존
if [ -d "$SIMDIR" ] && [ "$(ls -A "$SIMDIR" 2>/dev/null)" ]; then
    cp -r "$SIMDIR" "${OUTDIR}/simulations"
    echo "[$(date +%H:%M:%S)] Saved simulation artifacts"
fi
find "$SIMDIR" -name "*.xml" -exec cp {} "${OUTDIR}/" \; 2>/dev/null || true

# VRAM 해제
curl -s http://localhost:11434/api/generate \
    -d "{\"model\": \"${MODEL}\", \"keep_alive\": 0}" > /dev/null 2>&1 || true

echo "[$(date +%H:%M:%S)] Done. Results in ${OUTDIR}/"
