#!/bin/bash
# Chain: wait for run_009b → extract pressure → run EXP-D solvers → analyze
set -e
cd /home/imgyu/workspace/02_active/sim/slosim-agent

echo "=== Chain: run_009b → MeasureTool → EXP-D ==="

# Step 1: Extract pressure data from run_009b (10ms PART resolution)
echo "[1/4] Extracting run_009b pressure data..."
docker exec dsph-gpu /opt/dsph/bin/MeasureTool_linux64 \
    -dirin /data/exp1/run_009b/run_009b \
    -points /data/exp1/run_009b/oil_lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_009b/pressure_009b 2>&1 | tail -5
echo "  MeasureTool done."

# Step 2: Copy pressure CSV to host
echo "[2/4] Copying run_009b pressure CSV..."
mkdir -p research-v3/exp-c/analysis/run_009b
# MeasureTool creates *_Fluid.csv
docker cp dsph-gpu:/data/exp1/run_009b/pressure_009b_Fluid.csv \
    research-v3/exp-c/analysis/run_009b/pressure_009b.csv 2>/dev/null && \
    echo "  Copied pressure_009b.csv" || \
echo "  WARNING: trying alternate names..."; \
docker cp dsph-gpu:/data/exp1/run_009b/pressure_009b.csv \
    research-v3/exp-c/analysis/run_009b/pressure_009b.csv 2>/dev/null || true

# Step 3: Run EXP-D solvers
echo "[3/4] Running EXP-D solvers..."
bash research-v3/exp-d/run_solvers.sh

# Step 4: Analyze EXP-D
echo "[4/4] Analyzing EXP-D results..."
python3 research-v3/exp-d/analysis/analyze_expd.py

echo ""
echo "=== Chain complete: $(date) ==="
