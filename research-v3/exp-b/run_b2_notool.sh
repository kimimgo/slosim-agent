#!/bin/bash
# =============================================================================
# EXP-B: B2 (−DSPHTool) — 도메인 프롬프트 O, 도구 X
# ollama_generate.py 사용 + domain_prompt.txt
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROMPTS_DIR="$(dirname "$SCRIPT_DIR")/exp-a/prompts"
RESULTS_DIR="${SCRIPT_DIR}/results"
OLLAMA_URL="${LOCAL_ENDPOINT:-http://localhost:11434}"
DOMAIN_PROMPT="${SCRIPT_DIR}/domain_prompt.txt"

SCENARIOS=(S01 S02 S03 S04 S05 S06 S07 S08 S09 S10)
MODELS=("qwen3:32b" "qwen3:latest")

mkdir -p "$RESULTS_DIR"

echo "========================================"
echo "EXP-B B2 (−DSPHTool, Domain Prompt): $(date)"
echo "========================================"

for MODEL in "${MODELS[@]}"; do
    MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')

    for SCENARIO in "${SCENARIOS[@]}"; do
        OUTDIR="${RESULTS_DIR}/B2_${SCENARIO}_${MODEL_SAFE}"

        if [ -d "$OUTDIR" ] && [ -f "$OUTDIR/generated.xml" ] && [ -s "$OUTDIR/generated.xml" ]; then
            echo "[$(date +%H:%M:%S)] SKIP B2 $SCENARIO $MODEL_SAFE (exists)"
            continue
        fi

        PROMPT_FILE="${PROMPTS_DIR}/${SCENARIO}.txt"
        if [ ! -f "$PROMPT_FILE" ]; then
            echo "[$(date +%H:%M:%S)] ERROR: $PROMPT_FILE not found"
            continue
        fi

        echo ""
        echo "[$(date +%H:%M:%S)] ====== B2 $SCENARIO $MODEL_SAFE ======"
        echo "  Prompt: $(head -c 80 "$PROMPT_FILE")..."

        START=$(date +%s)

        python3 "${SCRIPT_DIR}/ollama_generate.py" \
            "$MODEL" \
            "$DOMAIN_PROMPT" \
            "$PROMPT_FILE" \
            "$OUTDIR" \
            "$OLLAMA_URL"

        END=$(date +%s)
        ELAPSED=$((END - START))

        # 메타데이터
        python3 -c "
import json
data = {
    'scenario': '${SCENARIO}',
    'model': '${MODEL}',
    'condition': 'B2_no_tool_domain_prompt',
    'duration_seconds': ${ELAPSED},
    'ollama_direct': True,
    'has_domain_prompt': True,
    'has_tools': False
}
json.dump(data, open('${OUTDIR}/metadata.json', 'w'), indent=2)
"

        echo "  Duration: ${ELAPSED}s"
        echo "[$(date +%H:%M:%S)] Done"
        sleep 5
    done
done

echo ""
echo "========================================"
echo "B2 (−DSPHTool) complete: $(date)"
echo "========================================"
