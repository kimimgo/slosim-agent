#!/usr/bin/env bash
# EXP-C: Parallel simulations on pajulab (4x RTX 4090)
# GPU 0: run_003v2 (Water dp=0.001, ~2hr, 6.1M particles)
# GPU 1: run_010   (Oil dp=0.002, probe fix, ~22min, 891K particles)
# GPU 2: run_011   (Oil dp=0.001, high-res, ~4-8hr, ~7M particles)
#
# Usage: scp to pajulab, then: bash run_pajulab_parallel.sh
set -euo pipefail

DSPH_BIN="/opt/dsph/bin"
BASE_DIR="$HOME/slosim-agent/research-v3/exp-c"
CASES_DIR="$BASE_DIR/cases"
PROBES_DIR="$BASE_DIR/probes"
SIMDATA="$HOME/slosim-agent/simulations/exp-c"

run_full_pipeline() {
    local GPU_ID="$1"
    local RUN_NAME="$2"
    local XML_FILE="$3"
    local PROBE_FILE="$4"
    local LOG_FILE="$SIMDATA/${RUN_NAME}.log"

    local OUT_DIR="$SIMDATA/$RUN_NAME"
    mkdir -p "$OUT_DIR"

    echo "[GPU $GPU_ID] Starting $RUN_NAME at $(date)" | tee "$LOG_FILE"

    # Copy files
    cp "$XML_FILE" "$OUT_DIR/${RUN_NAME}.xml"
    cp "$PROBE_FILE" "$OUT_DIR/pressure_probes.txt"

    # GenCase: input=/data/case.xml, output=/data/out/case (creates case.xml, case.bi4 inside out/)
    echo "[GPU $GPU_ID] GenCase..." | tee -a "$LOG_FILE"
    mkdir -p "$OUT_DIR/out"
    docker run --rm --gpus "device=$GPU_ID" \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/GenCase_linux64 /data/${RUN_NAME} /data/out/${RUN_NAME} -save:all \
        >> "$LOG_FILE" 2>&1
    echo "[GPU $GPU_ID] GenCase done" | tee -a "$LOG_FILE"

    # Solver: reads /data/out/case.xml, outputs PART files to /data/out/
    echo "[GPU $GPU_ID] Solver starting at $(date)..." | tee -a "$LOG_FILE"
    docker run --rm --gpus "device=$GPU_ID" \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/DualSPHysics5.4_linux64 /data/out/${RUN_NAME} /data/out -gpu -svres \
        >> "$LOG_FILE" 2>&1
    echo "[GPU $GPU_ID] Solver done at $(date)" | tee -a "$LOG_FILE"

    # MeasureTool: reads PART files from /data/out/
    echo "[GPU $GPU_ID] MeasureTool..." | tee -a "$LOG_FILE"
    docker run --rm --gpus "device=$GPU_ID" \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/MeasureTool_linux64 \
            -dirin /data/out \
            -points /data/pressure_probes.txt \
            -onlytype:-all,+fluid \
            -vars:+press,+vel \
            -savevtk /data/PressPoints \
            -savecsv /data/PressLateral \
        >> "$LOG_FILE" 2>&1
    echo "[GPU $GPU_ID] MeasureTool done" | tee -a "$LOG_FILE"

    echo "[GPU $GPU_ID] $RUN_NAME COMPLETE at $(date)" | tee -a "$LOG_FILE"
    echo "[GPU $GPU_ID] Results: $OUT_DIR/PressLateral_Press.csv" | tee -a "$LOG_FILE"
}

# Setup
mkdir -p "$SIMDATA"
echo "=============================================="
echo "  SPHERIC EXP-C: Parallel Simulation Launch"
echo "  $(date)"
echo "  Output: $SIMDATA"
echo "=============================================="

# Disk constraint: 142GB free, each dp=0.001 run ~70GB
# Strategy: run_010 (small) parallel with run_003v2, then run_011 after cleanup

echo "  Phase 1: run_010 (GPU 1, ~2min) + run_003v2 (GPU 0, ~2hr) in parallel"
echo "  Phase 2: cleanup run_003v2 bi4 → run_011 (GPU 2, ~4-8hr)"
echo ""

# Phase 1: Small Oil + Water in parallel
run_full_pipeline 1 "run_010" \
    "$CASES_DIR/run_010_oil_lat_dp002_probe_fix.xml" \
    "$PROBES_DIR/pressure_probes_oil_dp002.txt" &
PID_010=$!

run_full_pipeline 0 "run_003v2" \
    "$CASES_DIR/run_003v2_water_lat_dp001.xml" \
    "$PROBES_DIR/pressure_probes_water_dp001.txt" &
PID_003=$!

echo "  Phase 1 launched:"
echo "    GPU 0: run_003v2 (Water dp=0.001) PID=$PID_003"
echo "    GPU 1: run_010   (Oil dp=0.002)   PID=$PID_010"
echo ""

wait $PID_010 && echo "==> run_010 finished" || echo "==> run_010 FAILED"
wait $PID_003 && echo "==> run_003v2 finished" || echo "==> run_003v2 FAILED"

# Cleanup run_003v2 PART data to free ~70GB (keep CSV results)
echo ""
echo "  Cleaning run_003v2 bi4 files to free disk..."
docker run --rm -v "$SIMDATA/run_003v2:/data" alpine sh -c "rm -rf /data/out/data/Part_*.bi4 /data/out/data/PartMotionRef.ibi4" 2>/dev/null
echo "  Disk after cleanup:"
df -h / | tail -1

# Phase 2: High-res Oil
echo ""
echo "  Phase 2: run_011 (GPU 2, Oil dp=0.001)"
run_full_pipeline 2 "run_011" \
    "$CASES_DIR/run_011_oil_lat_dp001_hires.xml" \
    "$PROBES_DIR/pressure_probes_oil_dp001.txt"
echo "==> run_011 finished"

echo ""
echo "=============================================="
echo "  ALL SIMULATIONS COMPLETE at $(date)"
echo "=============================================="
echo ""
echo "Results:"
for run in run_003v2 run_010 run_011; do
    csv="$SIMDATA/$run/PressLateral_Press.csv"
    if [ -f "$csv" ]; then
        lines=$(wc -l < "$csv")
        echo "  $run: $csv ($lines lines)"
    else
        echo "  $run: MISSING"
    fi
done
