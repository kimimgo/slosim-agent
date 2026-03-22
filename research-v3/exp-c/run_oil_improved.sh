#!/usr/bin/env bash
# EXP-C: Improved Oil Lateral simulations
# Run 010: dp=0.002, corrected probe (~22 min)
# Run 011: dp=0.001, high-res (~4-8 hours)
#
# Usage: ./run_oil_improved.sh [010|011|all]
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CASES_DIR="$SCRIPT_DIR/cases"
PROBES_DIR="$SCRIPT_DIR/probes"
SIMDATA="/mnt/simdata/dualsphysics/exp-c"

DOCKER_IMAGE="dsph-agent:latest"
DSPH_BIN="/opt/dsph/bin"

run_simulation() {
    local RUN_NAME="$1"
    local XML_FILE="$2"
    local PROBE_FILE="$3"
    local DP="$4"

    local OUT_DIR="$SIMDATA/$RUN_NAME"
    mkdir -p "$OUT_DIR"

    echo "=============================================="
    echo "  $RUN_NAME: $(basename "$XML_FILE")"
    echo "  dp=$DP, Output: $OUT_DIR"
    echo "=============================================="

    # Copy case file
    cp "$XML_FILE" "$OUT_DIR/${RUN_NAME}.xml"
    cp "$PROBE_FILE" "$OUT_DIR/pressure_probes.txt"

    echo "[1/4] GenCase..."
    docker run --rm --gpus all \
        -v "$OUT_DIR:/data" \
        -v "$CASES_DIR:/cases:ro" \
        -w /data \
        --entrypoint '' \
        "$DOCKER_IMAGE" \
        $DSPH_BIN/GenCase_linux64 /data/${RUN_NAME} /data/${RUN_NAME}_out -save:all
    echo "  GenCase done."

    echo "[2/4] DualSPHysics GPU solver..."
    docker run --rm --gpus all \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        "$DOCKER_IMAGE" \
        $DSPH_BIN/DualSPHysics5.4_linux64 /data/${RUN_NAME}_out/${RUN_NAME} /data/${RUN_NAME}_out -gpu -svres
    echo "  Solver done."

    echo "[3/4] MeasureTool (pressure)..."
    docker run --rm --gpus all \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        "$DOCKER_IMAGE" \
        $DSPH_BIN/MeasureTool_linux64 \
            -dirin /data/${RUN_NAME}_out \
            -points /data/pressure_probes.txt \
            -onlytype:-all,+fluid \
            -vars:+press,+vel \
            -savevtk /data/PressPoints \
            -savecsv /data/PressLateral
    echo "  MeasureTool done."

    echo "[4/4] PartVTK (fluid particles for viz)..."
    docker run --rm --gpus all \
        -v "$OUT_DIR:/data" \
        -w /data \
        --entrypoint '' \
        "$DOCKER_IMAGE" \
        $DSPH_BIN/PartVTK_linux64 \
            -dirin /data/${RUN_NAME}_out \
            -savevtk /data/vtk/Fluid \
            -onlytype:-all,+fluid \
            -vars:+vel,+press
    echo "  PartVTK done."

    echo ""
    echo "  $RUN_NAME COMPLETE"
    echo "  Pressure CSV: $OUT_DIR/PressLateral_Press.csv"
    echo "  Run.csv: $OUT_DIR/${RUN_NAME}_out/Run.csv"
    echo ""
}

# Parse argument
TARGET="${1:-all}"

case "$TARGET" in
    010)
        run_simulation "run_010" \
            "$CASES_DIR/run_010_oil_lat_dp002_probe_fix.xml" \
            "$PROBES_DIR/pressure_probes_oil_dp002.txt" \
            "0.002"
        ;;
    011)
        run_simulation "run_011" \
            "$CASES_DIR/run_011_oil_lat_dp001_hires.xml" \
            "$PROBES_DIR/pressure_probes_oil_dp001.txt" \
            "0.001"
        ;;
    all)
        echo "Running both simulations sequentially..."
        echo ""
        run_simulation "run_010" \
            "$CASES_DIR/run_010_oil_lat_dp002_probe_fix.xml" \
            "$PROBES_DIR/pressure_probes_oil_dp002.txt" \
            "0.002"
        run_simulation "run_011" \
            "$CASES_DIR/run_011_oil_lat_dp001_hires.xml" \
            "$PROBES_DIR/pressure_probes_oil_dp001.txt" \
            "0.001"
        ;;
    *)
        echo "Usage: $0 [010|011|all]"
        exit 1
        ;;
esac

echo "========================================"
echo "  All requested simulations complete!"
echo "========================================"
