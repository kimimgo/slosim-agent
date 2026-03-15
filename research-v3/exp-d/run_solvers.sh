#!/bin/bash
# EXP-D: Run baseline and baffle solvers sequentially in Docker
# Usage: bash run_solvers.sh
# Prerequisites: GenCase already done for both cases
set -e

echo "=== EXP-D Solver Execution ==="
echo "Start: $(date)"

# Check GPU is free
if docker exec dsph-gpu ps aux 2>/dev/null | grep -q DualSPHysics; then
    echo "ERROR: GPU solver still running. Wait for it to finish."
    exit 1
fi

# Step 1: Baseline solver (~2-3 min for dp=0.008, 8s simulation)
echo ""
echo "[1/2] Running BASELINE solver..."
docker exec dsph-gpu /opt/dsph/bin/DualSPHysics5.4_linux64 \
    /data/exp-d/baseline/baseline \
    /data/exp-d/baseline/baseline \
    -gpu -svres
echo "  Baseline done: $(date)"

# Step 2: Baffle solver (~2-3 min)
echo ""
echo "[2/2] Running BAFFLE solver..."
docker exec dsph-gpu /opt/dsph/bin/DualSPHysics5.4_linux64 \
    /data/exp-d/baffle_center/baffle_center \
    /data/exp-d/baffle_center/baffle_center \
    -gpu -svres
echo "  Baffle done: $(date)"

# Step 3: Copy gauge CSVs to host
echo ""
echo "[3/3] Copying SWL gauge data..."
OUTDIR="/home/imgyu/workspace/02_active/sim/slosim-agent/research-v3/exp-d/results"
mkdir -p "$OUTDIR"

for case in baseline baffle_center; do
    for gauge in SWL_Left SWL_Right; do
        src="/data/exp-d/${case}/${case}/GaugesSWL_${gauge}.csv"
        if docker exec dsph-gpu test -f "$src"; then
            label=$(echo $case | sed 's/baffle_center/baffle/')
            docker cp "dsph-gpu:${src}" "${OUTDIR}/${label}_${gauge}.csv"
            echo "  Copied: ${label}_${gauge}.csv"
        else
            echo "  WARNING: ${src} not found"
        fi
    done
done

echo ""
echo "=== EXP-D Complete: $(date) ==="
echo "Results: $OUTDIR"
echo "Next: python3 research-v3/exp-d/analysis/analyze_expd.py"
