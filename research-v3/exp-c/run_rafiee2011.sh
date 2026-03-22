#!/bin/bash
# Rafiee & Thiagarajan (2011) Sloshing Benchmark — DualSPHysics GPU
# Tank: 1.3m x 0.1m x 0.9m, Fill: 0.18m, Sway a=0.1m, w=3.1165 rad/s
# GPU 0: dp=0.004 (coarse), GPU 1: dp=0.002 (fine)
# Usage: ./run_rafiee2011.sh [from pajulab ~/slosim-agent/]

set -euo pipefail

COMPOSE="docker compose"
CASES_DIR="/cases"
DATA_DIR="/data"
PROBES_DIR="/data/exp-c/probes"
RAFIEE_DIR="exp-c/rafiee2011"

echo "=== Rafiee 2011 Sloshing Benchmark ==="
echo "  Tank: 1.3m x 0.1m x 0.9m"
echo "  Fill: 0.18m, Sway a=0.1m, w=3.1165 rad/s"
echo ""

# --- Setup ---
mkdir -p simulations/${RAFIEE_DIR}/{dp004,dp002}

# Copy probe files to simulations for Docker access
mkdir -p simulations/exp-c/probes
cp research-v3/exp-c/probes/rafiee2011_probes_dp004.txt simulations/exp-c/probes/
cp research-v3/exp-c/probes/rafiee2011_probes_dp002.txt simulations/exp-c/probes/

# Copy XML cases to cases/ for Docker access
cp research-v3/exp-c/cases/rafiee2011_dp004_dbc.xml cases/rafiee2011_dp004.xml
cp research-v3/exp-c/cases/rafiee2011_dp002_dbc.xml cases/rafiee2011_dp002.xml

run_case() {
    local DP=$1
    local GPU=$2
    local OUTDIR="${RAFIEE_DIR}/dp${DP}"
    local CASE="rafiee2011_dp${DP}"
    local PROBE_FILE="rafiee2011_probes_dp${DP}.txt"

    echo ""
    echo "=== [GPU ${GPU}] Running dp=${DP} ==="
    echo "  Output: simulations/${OUTDIR}/"

    # 1) GenCase
    echo "  [1/4] GenCase..."
    ${COMPOSE} run --rm -e NVIDIA_VISIBLE_DEVICES=${GPU} dsph \
        gencase ${CASES_DIR}/${CASE} ${DATA_DIR}/${OUTDIR}/${CASE} -createdirs:1 \
        2>&1 | tail -5

    # 2) Solver
    echo "  [2/4] DualSPHysics GPU ${GPU}..."
    ${COMPOSE} run --rm -e NVIDIA_VISIBLE_DEVICES=${GPU} dsph \
        dualsphysics ${DATA_DIR}/${OUTDIR}/${CASE} ${DATA_DIR}/${OUTDIR} -gpu -svres \
        2>&1 | tail -10

    # 3) MeasureTool — pressure probes
    echo "  [3/4] MeasureTool (pressure probes)..."
    ${COMPOSE} run --rm dsph measuretool \
        -dirin ${DATA_DIR}/${OUTDIR} \
        -points ${DATA_DIR}/exp-c/probes/${PROBE_FILE} \
        -onlytype:-all,+fluid \
        -vars:-all,+press \
        -savecsv ${DATA_DIR}/${OUTDIR}/Pressure \
        2>&1 | tail -3

    # 4) PartVTK — fluid visualization
    echo "  [4/4] PartVTK (fluid VTK)..."
    mkdir -p simulations/${OUTDIR}/vtk
    ${COMPOSE} run --rm dsph partvtk \
        -dirin ${DATA_DIR}/${OUTDIR} \
        -savevtk ${DATA_DIR}/${OUTDIR}/vtk/PartFluid \
        -onlytype:-all,+fluid \
        -vars:+vel,+rhop,+press \
        2>&1 | tail -3

    echo "  [GPU ${GPU}] dp=${DP} DONE"
    echo "  Pressure CSV: simulations/${OUTDIR}/Pressure_Press.csv"
}

# --- Run in parallel on different GPUs ---
echo "Starting parallel runs on GPU 0 (dp004) and GPU 1 (dp002)..."
echo ""

run_case "004" "0" &
PID_004=$!

run_case "002" "1" &
PID_002=$!

# Wait for both
echo "Waiting for GPU 0 (dp004)..."
wait $PID_004 && echo "GPU 0 (dp004) COMPLETE" || echo "GPU 0 (dp004) FAILED"

echo "Waiting for GPU 1 (dp002)..."
wait $PID_002 && echo "GPU 1 (dp002) COMPLETE" || echo "GPU 1 (dp002) FAILED"

echo ""
echo "=== ALL RUNS COMPLETE ==="
echo "Results:"
echo "  dp=0.004: simulations/${RAFIEE_DIR}/dp004/"
echo "  dp=0.002: simulations/${RAFIEE_DIR}/dp002/"
echo ""
echo "Next: Run analysis script to compare with experimental data"
