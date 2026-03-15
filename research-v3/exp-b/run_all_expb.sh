#!/bin/bash
# =============================================================================
# EXP-B 전체 실행: B4 → B2 → B1
# 2×2 Factorial: (Domain Prompt × Tools) ablation
# =============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "╔════════════════════════════════════════╗"
echo "║  EXP-B Full Ablation: $(date +%H:%M)        ║"
echo "╚════════════════════════════════════════╝"
echo ""
echo "B0 (Full) = EXP-A results (already done)"
echo "B1 (−DomainPrompt) = Generic prompt + Tools"
echo "B2 (−DSPHTool) = Domain prompt + No tools"
echo "B4 (Bare LLM) = Generic prompt + No tools"
echo ""

START_ALL=$(date +%s)

# --- Phase 1: B4 (가장 빠름, Ollama API만) ---
echo "═══ Phase 1: B4 (Bare LLM) ═══"
bash "${SCRIPT_DIR}/run_b4_bare.sh"

# --- Phase 2: B2 (Ollama API + 도메인 프롬프트) ---
echo ""
echo "═══ Phase 2: B2 (−DSPHTool) ═══"
bash "${SCRIPT_DIR}/run_b2_notool.sh"

# --- Phase 3: B1 (slosim 바이너리, 가장 느림) ---
echo ""
echo "═══ Phase 3: B1 (−DomainPrompt) ═══"
bash "${SCRIPT_DIR}/run_b1_batch.sh"

END_ALL=$(date +%s)
TOTAL=$((END_ALL - START_ALL))

echo ""
echo "╔════════════════════════════════════════╗"
echo "║  EXP-B Complete: ${TOTAL}s total      ║"
echo "╚════════════════════════════════════════╝"

# Summary table
echo ""
echo "=== Results Summary ==="
for COND in B1 B2 B4; do
    echo ""
    echo "--- $COND ---"
    for D in "${SCRIPT_DIR}/results/${COND}_"*; do
        if [ -f "$D/metadata.json" ]; then
            NAME=$(basename "$D")
            DUR=$(python3 -c "import json; print(json.load(open('$D/metadata.json'))['duration_seconds'])" 2>/dev/null || echo "?")
            XML_SIZE=0
            if [ -f "$D/generated.xml" ]; then
                XML_SIZE=$(wc -c < "$D/generated.xml" 2>/dev/null || echo 0)
            elif [ -f "$D/simulations" ]; then
                XML_COUNT=$(find "$D/simulations" -name "*.xml" 2>/dev/null | wc -l)
                XML_SIZE="(${XML_COUNT} XMLs)"
            fi
            echo "  $NAME: ${DUR}s, XML=${XML_SIZE}"
        fi
    done
done
