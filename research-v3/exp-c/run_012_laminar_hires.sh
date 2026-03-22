#!/usr/bin/env bash
# Run 012: Oil dp=0.001, Laminar+SPS, ν=4.55e-5 m²/s (hi-res convergence)
# GPU 1 (free after sweep)
set -euo pipefail

DSPH_BIN="/opt/dsph/bin"
SIMDATA="/data/exp-c"
CASES_DIR="$HOME/slosim-agent/research-v3/exp-c/cases"
PROBES_DIR="$HOME/slosim-agent/research-v3/exp-c/probes"

GPU=1
OUT_DIR="$SIMDATA/run_012"
LOG_FILE="$SIMDATA/run_012.log"

echo "=== Run 012: Oil dp=0.001 Laminar+SPS (GPU $GPU) ===" | tee "$LOG_FILE"
echo "Started at $(date)" | tee -a "$LOG_FILE"

mkdir -p "$OUT_DIR/out"
cp "$CASES_DIR/run_012_oil_lat_dp001_laminar.xml" "$OUT_DIR/run_012.xml"
cp "$PROBES_DIR/pressure_probes_oil_dp001.txt" "$OUT_DIR/pressure_probes.txt"

echo "[GPU $GPU] GenCase..." | tee -a "$LOG_FILE"
docker run --rm --gpus "device=$GPU" \
    -v "$OUT_DIR:/data" -w /data --entrypoint '' \
    dsph-agent:latest \
    $DSPH_BIN/GenCase_linux64 /data/run_012 /data/out/run_012 -save:all \
    >> "$LOG_FILE" 2>&1
echo "[GPU $GPU] GenCase done" | tee -a "$LOG_FILE"

echo "[GPU $GPU] Solver starting at $(date)..." | tee -a "$LOG_FILE"
docker run --rm --gpus "device=$GPU" \
    -v "$OUT_DIR:/data" -w /data --entrypoint '' \
    dsph-agent:latest \
    $DSPH_BIN/DualSPHysics5.4_linux64 /data/out/run_012 /data/out -gpu -svres \
    >> "$LOG_FILE" 2>&1
echo "[GPU $GPU] Solver done at $(date)" | tee -a "$LOG_FILE"

echo "[GPU $GPU] MeasureTool..." | tee -a "$LOG_FILE"
docker run --rm --gpus "device=$GPU" \
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
echo "[GPU $GPU] MeasureTool done" | tee -a "$LOG_FILE"

echo ""
echo "=== Run 012 COMPLETE at $(date) ==="
echo "Results: $OUT_DIR/PressLateral_Press.csv"
if [ -f "$OUT_DIR/PressLateral_Press.csv" ]; then
    lines=$(wc -l < "$OUT_DIR/PressLateral_Press.csv")
    echo "  CSV lines: $lines"
fi
