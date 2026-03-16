#!/bin/bash
# SPHERIC T10 Oil ρ=920 — Artificial Viscosity α Parametric Study
# α = 1.0, 1.1, 1.2, ..., 3.0 (21 cases)
# Local RTX 4090
set -euo pipefail

PROJ_DIR="/home/imgyu/workspace/02_active/sim/slosim-agent"
SIM_BASE="/mnt/simdata/dualsphysics"
TEMPLATE="$PROJ_DIR/cases/SPHERIC_Test10_Oil_Low_Def.xml"
PROBE_FILE="$SIM_BASE/spheric_probes.txt"
DOCKER_COMPOSE="docker compose -f $PROJ_DIR/docker-compose.yml"

# Alphas: 1.0 to 3.0, step 0.1
ALPHAS=()
for i in $(seq 10 30); do
    a=$(echo "scale=1; $i/10" | bc)
    ALPHAS+=("$a")
done

echo "============================================="
echo "SPHERIC T10 Oil ρ=920 — α Parametric Study"
echo "α values: ${ALPHAS[*]}"
echo "Total cases: ${#ALPHAS[@]}"
echo "============================================="

generate_xml() {
    local alpha="$1"
    local tag="${alpha//./_}"
    local case_name="oil_rho920_a${tag}"
    local case_dir="$SIM_BASE/$case_name"
    local xml_file="$case_dir/${case_name}.xml"

    mkdir -p "$case_dir"

    # Copy and modify template
    sed -e 's|<rhop0 value="990"|<rhop0 value="920"|' \
        -e "s|<parameter key=\"Visco\" value=\"0.1\".*|<parameter key=\"Visco\" value=\"${alpha}\" comment=\"Art viscosity alpha=${alpha}\" />|" \
        -e "s|<parameter key=\"RhopOutMin\" value=\"700\"|<parameter key=\"RhopOutMin\" value=\"600\"|" \
        -e "s|<parameter key=\"RhopOutMax\" value=\"1300\"|<parameter key=\"RhopOutMax\" value=\"1200\"|" \
        "$TEMPLATE" > "$xml_file"

    echo "$xml_file"
}

run_case() {
    local alpha="$1"
    local tag="${alpha//./_}"
    local case_name="oil_rho920_a${tag}"
    local case_dir="$SIM_BASE/$case_name"
    local xml_file="$case_dir/${case_name}.xml"
    local out_dir="$case_dir/out"
    local measure_dir="$case_dir/measure"

    # Skip if already complete
    if [[ -f "$measure_dir/pressure_Press.csv" ]]; then
        echo "[SKIP] α=$alpha — already has pressure data"
        return 0
    fi

    echo ""
    echo "================================================"
    echo "[START] α=$alpha (case: $case_name)"
    echo "================================================"
    local start_time=$(date +%s)

    # 1) GenCase
    echo "[1/3] GenCase..."
    $DOCKER_COMPOSE run --rm --entrypoint '' -w /data dsph \
        /opt/dsph/bin/GenCase_linux64 \
        "/data/${case_name}/${case_name}" \
        "/data/${case_name}/out/${case_name}" \
        -save:all 2>&1 | tail -3

    # 2) Solver
    echo "[2/3] DualSPHysics GPU solver..."
    $DOCKER_COMPOSE run --rm --entrypoint '' dsph \
        /opt/dsph/bin/DualSPHysics5.4_linux64 \
        "/data/${case_name}/out/${case_name}" \
        "/data/${case_name}/out/${case_name}" \
        -gpu -svres 2>&1 | tail -5

    # 3) MeasureTool
    echo "[3/3] MeasureTool (pressure probes)..."
    mkdir -p "$measure_dir"
    cp "$PROBE_FILE" "$measure_dir/probes.txt"
    $DOCKER_COMPOSE run --rm --entrypoint '' dsph \
        /opt/dsph/bin/MeasureTool_linux64 \
        -dirin "/data/${case_name}/out/${case_name}" \
        -pointsdef "/data/${case_name}/measure/probes.txt" \
        -savecsv "/data/${case_name}/measure/pressure" \
        -vars:press,vel 2>&1 | tail -3

    local end_time=$(date +%s)
    local elapsed=$((end_time - start_time))
    echo "[DONE] α=$alpha — ${elapsed}s"
}

# Phase 1: Generate all XMLs
echo ""
echo "=== Phase 1: Generating ${#ALPHAS[@]} XML cases ==="
for alpha in "${ALPHAS[@]}"; do
    xml=$(generate_xml "$alpha")
    echo "  Created: $xml"
done

# Phase 2: Run all cases sequentially
echo ""
echo "=== Phase 2: Running ${#ALPHAS[@]} simulations ==="
TOTAL_START=$(date +%s)
completed=0
for alpha in "${ALPHAS[@]}"; do
    run_case "$alpha"
    completed=$((completed + 1))
    echo "  Progress: $completed/${#ALPHAS[@]}"
done
TOTAL_END=$(date +%s)
TOTAL_ELAPSED=$(( (TOTAL_END - TOTAL_START) / 60 ))

echo ""
echo "============================================="
echo "ALL DONE — ${#ALPHAS[@]} cases in ${TOTAL_ELAPSED} minutes"
echo "============================================="

# Phase 3: Verify outputs
echo ""
echo "=== Verification ==="
for alpha in "${ALPHAS[@]}"; do
    tag="${alpha//./_}"
    f="$SIM_BASE/oil_rho920_a${tag}/measure/pressure_Press.csv"
    if [[ -f "$f" ]]; then
        lines=$(wc -l < "$f")
        echo "  α=$alpha: OK ($lines lines)"
    else
        echo "  α=$alpha: MISSING"
    fi
done
