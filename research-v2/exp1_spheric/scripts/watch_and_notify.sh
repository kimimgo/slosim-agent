#!/bin/bash
# Watch for pipeline completion, then run analysis + Telegram notification
# Usage: nohup bash research-v2/exp1_spheric/scripts/watch_and_notify.sh &
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

LOG="research-v2/exp1_spheric/scripts/watch.log"
log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG"; }

MONITOR_PID=60551

log "Watching monitor PID $MONITOR_PID and pipeline..."

# Wait for monitor to finish OR all data to appear
while true; do
    # Check if all required CSVs exist
    ALL_READY=true
    for f in \
        simulations/exp1/run_003/PressConsistent_Press.csv \
        simulations/exp1/run_005/PressConsistent_Press.csv \
        simulations/exp1/run_006/PressRoof_Press.csv; do
        [ -f "$f" ] || ALL_READY=false
    done

    if [ "$ALL_READY" = true ]; then
        log "All data files detected!"
        break
    fi

    # Check if monitor died
    if ! ps -p $MONITOR_PID > /dev/null 2>&1; then
        log "Monitor PID $MONITOR_PID exited"
        # Re-check data availability
        sleep 5
        for f in \
            simulations/exp1/run_003/PressConsistent_Press.csv \
            simulations/exp1/run_005/PressConsistent_Press.csv \
            simulations/exp1/run_006/PressRoof_Press.csv; do
            [ -f "$f" ] || ALL_READY=false
        done
        if [ "$ALL_READY" = false ]; then
            log "WARNING: Monitor exited but not all data ready"
        fi
        break
    fi

    sleep 120  # Check every 2 minutes
done

# Run analysis + notification
log "=== Running final analysis ==="
bash research-v2/exp1_spheric/scripts/run_final_analysis.sh >> "$LOG" 2>&1

log "=== Sending Telegram notification ==="
bash research-v2/exp1_spheric/scripts/notify_telegram.sh >> "$LOG" 2>&1

log "=== DONE ==="
