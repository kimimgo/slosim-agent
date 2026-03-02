#!/bin/bash
# Post-Run003 Pipeline: MeasureTool → Run 005 → Run 006
# 실행: bash research-v2/exp1_spheric/scripts/post_run003_pipeline.sh
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

RESEARCH_CASES="research-v2/exp1_spheric/cases"

echo "=== Step 1: MeasureTool for Run 003 ==="
docker compose run --rm \
    -v /tmp/lateral_probe.txt:/cases/lateral_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_003 \
    -points /cases/lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_003/PressConsistent

echo "=== Step 1 done. Running convergence analysis ==="
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py

echo ""
echo "=== Step 2: Run 005 (Oil Lat, dp=0.002, DBC) ==="
# Copy XML to Docker-accessible location, create output dir
cp "$RESEARCH_CASES/run_005_oil_lat_dp002_dbc.xml" cases/run_005.xml
mkdir -p simulations/exp1/run_005

docker compose run --rm dsph gencase \
    /cases/run_005 \
    /data/exp1/run_005/run_005

docker compose run --rm dsph dualsphysics \
    /data/exp1/run_005/run_005 /data/exp1/run_005 \
    -gpu -svres

echo "=== Run 005 solver done. MeasureTool ==="
docker compose run --rm \
    -v /tmp/lateral_probe.txt:/cases/lateral_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_005 \
    -points /cases/lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_005/PressConsistent

echo ""
echo "=== Step 3: Run 006 (Water Roof, dp=0.002, DBC) ==="
cp "$RESEARCH_CASES/run_006_water_roof_dp002_dbc.xml" cases/run_006.xml
mkdir -p simulations/exp1/run_006

docker compose run --rm dsph gencase \
    /cases/run_006 \
    /data/exp1/run_006/run_006

docker compose run --rm dsph dualsphysics \
    /data/exp1/run_006/run_006 /data/exp1/run_006 \
    -gpu -svres

echo "=== Run 006 solver done. MeasureTool ==="
# Roof probes need different probe file (z near ceiling at z=0.503)
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
    -savecsv /data/exp1/run_006/PressRoof

echo ""
echo "=========================================="
echo "=== ALL RUNS COMPLETE ==="
echo "=========================================="
echo "Next: python3 research-v2/exp1_spheric/analysis/convergence_analysis.py"
echo "Next: Create validation_report.md"
