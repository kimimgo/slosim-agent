#!/bin/bash
# =============================================================================
# EXP-A→C Bridge: Agent-generated S02 XML → DualSPHysics Solver → Physics
# =============================================================================
# Purpose: Prove that agent-generated XML (EXP-A, M-A3=100%) produces
#          physically valid simulation results (EXP-C comparison)
# Usage:   bash research-v3/exp-c/run_agent_bridge.sh [gpu_id]
# =============================================================================
set -euo pipefail

GPU_ID="${1:-3}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BRIDGE_DIR="${SCRIPT_DIR}/agent-bridge"
AGENT_XML="${SCRIPT_DIR}/../exp-a/results/S02_qwen3_32b_trial1/sloshing_case.xml"
CASE_NAME="AgentS02"

echo "========================================"
echo "EXP-A→C Bridge Experiment"
echo "========================================"
echo "Source: S02 agent-generated XML (M-A3=100%)"
echo "GPU:    ${GPU_ID}"
echo "Output: ${BRIDGE_DIR}"
echo "========================================"

# Setup
mkdir -p "${BRIDGE_DIR}"
cp "${AGENT_XML}" "${BRIDGE_DIR}/${CASE_NAME}_Def.xml"

# Add SWL gauge probes (same positions as EXP-C SPHERIC)
# Left wall: x=0.01, Center: x=0.3, Right wall: x=0.59
cat > "${BRIDGE_DIR}/probe_gauges.txt" << 'EOF'
0.01 0.15 0.0 0.01 0.15 0.10
0.30 0.15 0.0 0.30 0.15 0.10
0.59 0.15 0.0 0.59 0.15 0.10
EOF

echo "[$(date +%H:%M:%S)] Running GenCase..."
docker run --rm \
    --entrypoint '' \
    -v "${BRIDGE_DIR}:/data" \
    -v "$(dirname ${SCRIPT_DIR})/cases:/cases:ro" \
    dsph-agent:latest \
    /opt/dsph/bin/GenCase_linux64 /data/${CASE_NAME}_Def /data/${CASE_NAME}

echo "[$(date +%H:%M:%S)] GenCase complete. Running Solver (GPU ${GPU_ID})..."
docker run --rm \
    --gpus "\"device=${GPU_ID}\"" \
    --entrypoint '' \
    -v "${BRIDGE_DIR}:/data" \
    dsph-agent:latest \
    /opt/dsph/bin/DualSPHysics5.4_linux64 /data/${CASE_NAME} /data/${CASE_NAME} -gpu -svres

echo "[$(date +%H:%M:%S)] Solver complete. Running MeasureTool..."

# Extract SWL elevation at probe locations
docker run --rm \
    --entrypoint '' \
    -v "${BRIDGE_DIR}:/data" \
    dsph-agent:latest \
    /opt/dsph/bin/MeasureTool_linux64 \
    -dirin /data/${CASE_NAME} \
    -pointsdef file:gauges:/data/probe_gauges.txt \
    -onlytype:-all,+fluid \
    -vars:-all,+vel,+press,+mass \
    -savevtk /data/gauges \
    -savecsv /data/gauges

echo "[$(date +%H:%M:%S)] MeasureTool complete."

# Extract PartVTK for visualization
docker run --rm \
    --entrypoint '' \
    -v "${BRIDGE_DIR}:/data" \
    dsph-agent:latest \
    /opt/dsph/bin/PartVTK_linux64 \
    -dirin /data/${CASE_NAME} \
    -savevtk /data/PartFluid \
    -onlytype:-all,+fluid

echo "[$(date +%H:%M:%S)] All complete. Results in ${BRIDGE_DIR}/"
echo ""
echo "Next steps:"
echo "  1. Check gauges_*.csv for SWL time series"
echo "  2. Compare with Chen 2018 experimental data"
echo "  3. Compute correlation (r) and M2 metric"
