#!/bin/bash
# Send EXP-1 results summary + plots to Telegram
# Usage: bash research-v2/exp1_spheric/scripts/notify_telegram.sh
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

TG="/home/imgyu/.claude/scripts/telegram-send.sh"
FIG_DIR="research-v2/exp1_spheric/figures"

# Run final analysis first
echo "Running final analysis..."
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py > /dev/null 2>&1 || true
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py > /dev/null 2>&1 || true
python3 research-v2/exp1_spheric/analysis/paper_figures.py > /dev/null 2>&1 || true

# Get verdict
VERDICT=$(python3 research-v2/exp1_spheric/scripts/final_verdict.py 2>&1)
echo "$VERDICT"

# Send text summary
bash "$TG" send --project slosim-agent "$(cat <<'MSG'
📊 EXP-1 SPHERIC Test 10 — 결과 요약

MSG
)
$VERDICT"

# Send key figures
for fig in convergence_study.png fig_timeseries.png fig_oil_lateral.png fig_water_roof.png; do
    if [ -f "$FIG_DIR/$fig" ]; then
        echo "Sending $fig..."
        bash "$TG" photo "$FIG_DIR/$fig" --caption "$fig"
        sleep 1
    fi
done

echo "Telegram notification sent!"
