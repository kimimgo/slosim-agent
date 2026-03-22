#!/usr/bin/env bash
# Run after run_003v2 completes: cleanup bi4 → start run_011 (Oil dp=0.001)
set -euo pipefail

DSPH_BIN="/opt/dsph/bin"
SIMDATA="$HOME/slosim-agent/simulations/exp-c"
CASES_DIR="$HOME/slosim-agent/research-v3/exp-c/cases"
PROBES_DIR="$HOME/slosim-agent/research-v3/exp-c/probes"

echo "=== Phase 1: Cleanup run_003v2 bi4 files ==="
echo "Disk before cleanup:"
df -h / | tail -1

docker run --rm -v "$SIMDATA/run_003v2:/data" alpine sh -c \
    'rm -rf /data/out/data/Part_*.bi4 /data/out/data/PartMotionRef.ibi4 /data/out/data/PartInfo.ibi4 /data/out/data/PartOut_*.obi4'

echo "Disk after cleanup:"
df -h / | tail -1

echo ""
echo "=== Phase 2: Run 011 (Oil dp=0.001, mDBC, GPU 2) ==="
OUT_DIR="$SIMDATA/run_011"
LOG_FILE="$SIMDATA/run_011.log"
mkdir -p "$OUT_DIR"

cp "$CASES_DIR/run_011_oil_lat_dp001_hires.xml" "$OUT_DIR/run_011.xml"
cp "$PROBES_DIR/pressure_probes_oil_dp001.txt" "$OUT_DIR/pressure_probes.txt"

echo "[GPU 2] GenCase..." | tee "$LOG_FILE"
mkdir -p "$OUT_DIR/out"
docker run --rm --gpus "device=2" \
    -v "$OUT_DIR:/data" -w /data --entrypoint '' \
    dsph-agent:latest \
    $DSPH_BIN/GenCase_linux64 /data/run_011 /data/out/run_011 -save:all \
    >> "$LOG_FILE" 2>&1
echo "[GPU 2] GenCase done" | tee -a "$LOG_FILE"

echo "[GPU 2] Solver starting at $(date)..." | tee -a "$LOG_FILE"
docker run --rm --gpus "device=2" \
    -v "$OUT_DIR:/data" -w /data --entrypoint '' \
    dsph-agent:latest \
    $DSPH_BIN/DualSPHysics5.4_linux64 /data/out/run_011 /data/out -gpu -svres \
    >> "$LOG_FILE" 2>&1
echo "[GPU 2] Solver done at $(date)" | tee -a "$LOG_FILE"

echo "[GPU 2] MeasureTool..." | tee -a "$LOG_FILE"
docker run --rm --gpus "device=2" \
    -v "$OUT_DIR:/data" -w /data --entrypoint '' \
    dsph-agent:latest \
    $DSPH_BIN/MeasureTool_linux64 \
        -dirin /data/out \
        -points /data/pressure_probes.txt \
        -onlytype:-all,+fluid \
        -vars:+press,+vel \
        -savevtk /data/PressPoints \
        -savecsv /data/PressLateral \
    >> "$LOG_FILE" 2>&1
echo "[GPU 2] MeasureTool done" | tee -a "$LOG_FILE"

echo ""
echo "=== Run 011 COMPLETE at $(date) ==="
echo "Results: $OUT_DIR/PressLateral_Press.csv"
if [ -f "$OUT_DIR/PressLateral_Press.csv" ]; then
    lines=$(wc -l < "$OUT_DIR/PressLateral_Press.csv")
    echo "  CSV lines: $lines"
fi
