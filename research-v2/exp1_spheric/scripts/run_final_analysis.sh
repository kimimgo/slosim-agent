#!/bin/bash
# Run all analysis after pipeline completion
# Usage: bash research-v2/exp1_spheric/scripts/run_final_analysis.sh
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

echo "=== EXP-1 Final Analysis ==="
echo ""

# Check data availability
check_file() {
    if [ ! -f "$1" ]; then
        echo "MISSING: $1"
        return 1
    fi
    echo "  OK: $1"
    return 0
}

echo "Checking data files..."
READY=true
check_file simulations/exp1/run_001/PressConsistent_Press.csv || READY=false
check_file simulations/exp1/run_002/PressConsistent_Press.csv || READY=false
check_file simulations/exp1/run_003/PressConsistent_Press.csv || READY=false
check_file simulations/exp1/run_005/PressConsistent_Press.csv || READY=false
check_file simulations/exp1/run_006/PressRoof_Press.csv || READY=false

if [ "$READY" = false ]; then
    echo ""
    echo "Some data files missing. Run the pipeline first."
    echo "Continuing with available data..."
fi

echo ""
echo "=== Step 1: Convergence Analysis (3-level GCI) ==="
python3 research-v2/exp1_spheric/analysis/convergence_analysis.py 2>&1

echo ""
echo "=== Step 2: Oil Lateral + Water Roof Analysis ==="
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py 2>&1

echo ""
echo "=== Step 3: Publication Figures ==="
python3 research-v2/exp1_spheric/analysis/paper_figures.py 2>&1

echo ""
echo "=========================================="
echo "=== ALL ANALYSIS COMPLETE ==="
echo "=========================================="
echo "Figures: research-v2/exp1_spheric/figures/"
echo "Metrics: research-v2/exp1_spheric/analysis/metrics.json"
echo "Oil/Roof: research-v2/exp1_spheric/analysis/oil_roof_metrics.json"
ls -la research-v2/exp1_spheric/figures/*.png 2>/dev/null
