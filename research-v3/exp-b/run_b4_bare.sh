#!/bin/bash
# =============================================================================
# EXP-B: B4 Bare LLM — 도구 없이 Ollama 직접 호출
# ollama_generate.py 사용 + generic system prompt
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROMPTS_DIR="$(dirname "$SCRIPT_DIR")/exp-a/prompts"
RESULTS_DIR="${SCRIPT_DIR}/results"
OLLAMA_URL="${LOCAL_ENDPOINT:-http://localhost:11434}"

SCENARIOS=(S01 S02 S03 S04 S05 S06 S07 S08 S09 S10)
MODELS=("qwen3:32b" "qwen3:latest")

SYSTEM_PROMPT="You are a helpful assistant. The user wants a fluid dynamics simulation. Generate an XML configuration file for DualSPHysics based on their request. Output ONLY the raw XML content, no explanations."

mkdir -p "$RESULTS_DIR"

echo "========================================"
echo "EXP-B B4 (Bare LLM): $(date)"
echo "========================================"

for MODEL in "${MODELS[@]}"; do
    MODEL_SAFE=$(echo "$MODEL" | tr ':' '_')

    for SCENARIO in "${SCENARIOS[@]}"; do
        OUTDIR="${RESULTS_DIR}/B4_${SCENARIO}_${MODEL_SAFE}"

        if [ -d "$OUTDIR" ] && [ -f "$OUTDIR/generated.xml" ] && [ -s "$OUTDIR/generated.xml" ]; then
            echo "[$(date +%H:%M:%S)] SKIP B4 $SCENARIO $MODEL_SAFE (exists)"
            continue
        fi

        PROMPT_FILE="${PROMPTS_DIR}/${SCENARIO}.txt"
        if [ ! -f "$PROMPT_FILE" ]; then
            echo "[$(date +%H:%M:%S)] ERROR: $PROMPT_FILE not found"
            continue
        fi

        echo ""
        echo "[$(date +%H:%M:%S)] ====== B4 $SCENARIO $MODEL_SAFE ======"
        echo "  Prompt: $(head -c 80 "$PROMPT_FILE")..."

        START=$(date +%s)

        python3 "${SCRIPT_DIR}/ollama_generate.py" \
            "$MODEL" \
            "$SYSTEM_PROMPT" \
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
    'condition': 'B4_bare_llm',
    'duration_seconds': ${ELAPSED},
    'ollama_direct': True,
    'has_domain_prompt': False,
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
echo "B4 Bare LLM complete: $(date)"
echo "========================================"

# Summary
for MODEL_SAFE in qwen3_32b qwen3_latest; do
    echo ""
    echo "=== $MODEL_SAFE ==="
    for S in "${SCENARIOS[@]}"; do
        D="${RESULTS_DIR}/B4_${S}_${MODEL_SAFE}"
        if [ -f "$D/metadata.json" ]; then
            DUR=$(python3 -c "import json; print(json.load(open('$D/metadata.json'))['duration_seconds'])" 2>/dev/null)
            XML_SIZE=$(wc -c < "$D/generated.xml" 2>/dev/null || echo 0)
            echo "  $S: ${DUR}s, XML=${XML_SIZE}B"
        fi
    done
done
