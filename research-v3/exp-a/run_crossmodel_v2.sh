#!/bin/bash
# =============================================================================
# Cross-Model Experiment v2 — JSON compat fix 적용 배치 실행
# 3 non-Qwen3-32b/8b models × 10 scenarios = 30 runs
# Sequential execution (single Ollama instance on port 11434)
#
# Usage: nohup bash run_crossmodel_v2.sh > /tmp/crossmodel_v2.log 2>&1 &
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
AGENT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
SLOSIM="${AGENT_DIR}/slosim"

# Models to test (tool-calling capable, non-duplicated with existing pajulab data)
MODELS=(
    "llama3.3:70b-instruct-q4_K_M"
    "llama3.1:8b-instruct-q8_0"
    "qwen3:14b"
)

SCENARIOS=(S01 S02 S03 S04 S05 S06 S07 S08 S09 S10)
RESULTS_BASE="${SCRIPT_DIR}/results_crossmodel"
OLLAMA_PORT=11434

echo "============================================================"
echo " Cross-Model Experiment v2 (JSON compat fix)"
echo " Models: ${#MODELS[@]}, Scenarios: ${#SCENARIOS[@]}"
echo " Started: $(date)"
echo "============================================================"

# Warm-up: check Ollama is running
if ! curl -s "http://localhost:${OLLAMA_PORT}/api/tags" > /dev/null 2>&1; then
    echo "[ERROR] Ollama not running on port ${OLLAMA_PORT}"
    exit 1
fi

for MODEL in "${MODELS[@]}"; do
    MODEL_SAFE=$(echo "$MODEL" | tr ':/' '_')
    MODEL_DIR="${RESULTS_BASE}/${MODEL_SAFE}"
    mkdir -p "$MODEL_DIR"

    echo ""
    echo "──────────────────────────────────────────────────"
    echo " Model: ${MODEL} (${MODEL_SAFE})"
    echo "──────────────────────────────────────────────────"

    # Pre-load model
    echo "[$(date +%H:%M:%S)] Warming up ${MODEL}..."
    curl -s "http://localhost:${OLLAMA_PORT}/api/generate" \
        -d "{\"model\": \"${MODEL}\", \"prompt\": \"hello\", \"stream\": false}" \
        > /dev/null 2>&1

    for SCENARIO in "${SCENARIOS[@]}"; do
        OUTDIR="${MODEL_DIR}/${SCENARIO}"
        PROMPT_FILE="${SCRIPT_DIR}/prompts/${SCENARIO}.txt"
        CONFIG_FILE="${SCRIPT_DIR}/configs/config_${MODEL_SAFE}.json"

        # Skip if already completed
        if [ -f "${OUTDIR}/metadata.json" ]; then
            DUR=$(python3 -c "import json; d=json.load(open('${OUTDIR}/metadata.json')); print(d.get('duration_seconds',0))" 2>/dev/null || echo "0")
            if [ "$DUR" -gt "5" ]; then
                echo "[$(date +%H:%M:%S)] SKIP ${SCENARIO} — already done (${DUR}s)"
                continue
            fi
        fi

        if [ ! -f "$PROMPT_FILE" ]; then
            echo "[$(date +%H:%M:%S)] SKIP ${SCENARIO} — no prompt file"
            continue
        fi

        mkdir -p "$OUTDIR"

        # Apply config
        if [ -f "$CONFIG_FILE" ]; then
            cp "$CONFIG_FILE" "${OUTDIR}/.opencode.json"
        else
            echo "[WARN] No config for ${MODEL_SAFE}, creating default"
            cat > "${OUTDIR}/.opencode.json" << CFGEOF
{
  "providers": { "local": { "apiKey": "dummy", "disabled": false } },
  "agents": {
    "coder": { "model": "local.${MODEL}", "maxTokens": 0 },
    "summarizer": { "model": "local.${MODEL}", "maxTokens": 0 },
    "task": { "model": "local.${MODEL}", "maxTokens": 0 },
    "title": { "model": "local.${MODEL}", "maxTokens": 0 }
  }
}
CFGEOF
        fi

        PROMPT=$(cat "$PROMPT_FILE")
        START=$(date +%s)
        echo "[$(date +%H:%M:%S)] Running ${SCENARIO} with ${MODEL}..."

        LOCAL_ENDPOINT="http://localhost:${OLLAMA_PORT}" \
            timeout 600 "$SLOSIM" \
            -c "$OUTDIR" \
            -p "$PROMPT" \
            -q -f json \
            > "${OUTDIR}/output.json" 2>"${OUTDIR}/stderr.log" || true

        EXIT_CODE=${PIPESTATUS[0]:-$?}
        END=$(date +%s)
        DURATION=$((END - START))

        # Find generated XML
        XML_FILE=""
        for f in "${OUTDIR}"/*.xml; do
            if [ -f "$f" ]; then
                XML_FILE="$f"
                break
            fi
        done

        # Save metadata
        cat > "${OUTDIR}/metadata.json" << METAEOF
{
    "scenario": "${SCENARIO}",
    "model": "${MODEL}",
    "exit_code": ${EXIT_CODE},
    "duration_seconds": ${DURATION},
    "xml_found": $([ -n "$XML_FILE" ] && echo "true" || echo "false"),
    "timestamp": "$(date -Iseconds)"
}
METAEOF

        STATUS="✓"
        [ -z "$XML_FILE" ] && STATUS="✗"
        echo "[$(date +%H:%M:%S)] ${STATUS} ${SCENARIO} done (${DURATION}s, xml=$([ -n "$XML_FILE" ] && echo "yes" || echo "no"))"

        # Brief cooldown
        sleep 3
    done

    # Unload model
    curl -s "http://localhost:${OLLAMA_PORT}/api/generate" \
        -d "{\"model\": \"${MODEL}\", \"keep_alive\": 0}" > /dev/null 2>&1 || true

    echo "[$(date +%H:%M:%S)] ${MODEL} COMPLETE"
done

echo ""
echo "============================================================"
echo " ALL EXPERIMENTS COMPLETE: $(date)"
echo " Results: ${RESULTS_BASE}/"
echo "============================================================"

# Summary
echo ""
echo "SUMMARY:"
printf "%-6s" ""
for MODEL in "${MODELS[@]}"; do
    printf " %-15s" "$(echo "$MODEL" | cut -d: -f1 | cut -c1-12)"
done
echo ""

for SCENARIO in "${SCENARIOS[@]}"; do
    printf "%-6s" "$SCENARIO"
    for MODEL in "${MODELS[@]}"; do
        MODEL_SAFE=$(echo "$MODEL" | tr ':/' '_')
        META="${RESULTS_BASE}/${MODEL_SAFE}/${SCENARIO}/metadata.json"
        if [ -f "$META" ]; then
            XML_FOUND=$(python3 -c "import json; print(json.load(open('${META}'))['xml_found'])" 2>/dev/null || echo "false")
            DUR=$(python3 -c "import json; print(json.load(open('${META}'))['duration_seconds'])" 2>/dev/null || echo "?")
            if [ "$XML_FOUND" = "True" ] || [ "$XML_FOUND" = "true" ]; then
                printf " %-15s" "✓ ${DUR}s"
            else
                printf " %-15s" "✗ ${DUR}s"
            fi
        else
            printf " %-15s" "—"
        fi
    done
    echo ""
done
