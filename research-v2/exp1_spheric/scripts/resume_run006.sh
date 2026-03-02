#!/bin/bash
# Resume monitoring for Run 006 after crash+restart
# Run 006 restarted at t=3.446s via -partbegin:3446
# This script waits for solver, runs MeasureTool, then triggers analysis
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

LOG="research-v2/exp1_spheric/scripts/resume_run006.log"
log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG"; }

log "Waiting for Run 006 solver to finish..."

# Wait for Docker container to exit (solver completion)
while docker ps --format '{{.Names}}' 2>/dev/null | grep -q dsph; do
    # Check progress
    LAST_T=$(tail -1 simulations/exp1/run_006/RunPARTs.csv 2>/dev/null | cut -d';' -f2 | head -c 5)
    log "  Run 006 at t=${LAST_T:-?}/7.0"
    sleep 120
done

log "Run 006 solver finished!"

# Check if completed or crashed
FINAL_T=$(tail -1 simulations/exp1/run_006/RunPARTs.csv 2>/dev/null | cut -d';' -f2 | head -c 5)
if (( $(echo "$FINAL_T < 6.5" | bc -l 2>/dev/null || echo 1) )); then
    log "WARNING: Run 006 may have crashed again at t=$FINAL_T"
fi

# Run MeasureTool for Roof — multiple x positions to find best Sensor 3 match
# SPHERIC paper says "Sensor 3" for roof impact but exact x position unclear from PDF
# Probes at z=0.503 (5mm below roof at z=0.508), y=0.031 (tank center breadth)
log "=== MeasureTool Run 006 (Roof) ==="
cat > /tmp/roof_probe.txt << 'PROBES'
POINTS
0.050 0.031 0.503
0.150 0.031 0.503
0.225 0.031 0.503
0.450 0.031 0.503
0.675 0.031 0.503
0.750 0.031 0.503
0.850 0.031 0.503
PROBES

docker compose run --rm \
    -v /tmp/roof_probe.txt:/cases/roof_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_006 \
    -points /cases/roof_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_006/PressRoof \
    >> "$LOG" 2>&1
log "MeasureTool Run 006 done"

# Run analysis
log "=== Running analysis ==="
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py >> "$LOG" 2>&1
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py >> "$LOG" 2>&1
python3 research-v2/exp1_spheric/analysis/paper_figures.py >> "$LOG" 2>&1
log "Analysis done"

# Final verdict
log "=== Final Verdict ==="
python3 research-v2/exp1_spheric/scripts/final_verdict.py >> "$LOG" 2>&1

# Telegram notification
log "=== Sending Telegram ==="
bash research-v2/exp1_spheric/scripts/notify_telegram.sh >> "$LOG" 2>&1

log "=== ALL DONE ==="
