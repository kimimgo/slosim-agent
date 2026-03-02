#!/bin/bash
# Monitor Run 003 completion, then run post-processing + Run 005 + Run 006
# Usage: nohup bash research-v2/exp1_spheric/scripts/monitor_and_continue.sh &
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

LOG="research-v2/exp1_spheric/scripts/pipeline.log"

log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG"; }

log "=== Pipeline monitor started ==="
log "Waiting for Run 003 to complete..."

# Wait for Run 003 solver to finish (container exits)
while docker ps --filter "name=dsph-run" --format "{{.Names}}" 2>/dev/null | grep -q dsph-run; do
    sleep 60
done
log "Run 003 container stopped"

# Verify Run 003 completed successfully
if tail -5 simulations/exp1/run_003/Run.out 2>/dev/null | grep -q "Finished"; then
    log "Run 003 completed successfully"
else
    log "WARNING: Run 003 may not have completed normally"
    tail -5 simulations/exp1/run_003/Run.out >> "$LOG" 2>&1
fi

# Step 1: MeasureTool for Run 003
log "=== MeasureTool for Run 003 ==="
docker compose run --rm \
    -v /tmp/lateral_probe.txt:/cases/lateral_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_003 \
    -points /cases/lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_003/PressConsistent \
    >> "$LOG" 2>&1
log "MeasureTool Run 003 done"

# Step 2: Run 005 solver (Oil, GenCase already done)
log "=== Solver Run 005 (Oil Lat) ==="
docker compose run --rm dsph dualsphysics \
    /data/exp1/run_005/run_005 /data/exp1/run_005 \
    -gpu -svres \
    >> "$LOG" 2>&1
log "Run 005 solver done"

# MeasureTool for Run 005
docker compose run --rm \
    -v /tmp/lateral_probe.txt:/cases/lateral_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_005 \
    -points /cases/lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_005/PressConsistent \
    >> "$LOG" 2>&1
log "MeasureTool Run 005 done"

# Step 3: Run 006 solver (Water Roof, GenCase already done)
log "=== Solver Run 006 (Water Roof) ==="
docker compose run --rm dsph dualsphysics \
    /data/exp1/run_006/run_006 /data/exp1/run_006 \
    -gpu -svres \
    >> "$LOG" 2>&1
log "Run 006 solver done"

# MeasureTool for Run 006 (roof probes)
cat > /tmp/roof_probe.txt << 'PROBES'
POINTS
0.225 0.031 0.503
0.450 0.031 0.503
0.675 0.031 0.503
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

# Step 4: Run convergence analysis (3-level GCI)
log "=== Running convergence analysis ==="
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py >> "$LOG" 2>&1
log "Convergence analysis done"

# Step 5: Oil + Roof analysis
log "=== Running Oil/Roof analysis ==="
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py >> "$LOG" 2>&1
log "Oil/Roof analysis done"

# Step 6: Paper-quality figures
log "=== Generating paper figures ==="
python3 research-v2/exp1_spheric/analysis/paper_figures.py >> "$LOG" 2>&1
log "Paper figures done"

log "=========================================="
log "=== ALL PIPELINE STEPS COMPLETE ==="
log "=========================================="
log "Check figures at: research-v2/exp1_spheric/figures/"
log "Check metrics at: research-v2/exp1_spheric/analysis/metrics.json"
