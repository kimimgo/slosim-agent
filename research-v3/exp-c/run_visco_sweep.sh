#!/usr/bin/env bash
# Oil Visco parametric sweep: α = [0.5, 1.0, 2.0, 5.0]
# dp=0.002, TimeMax=3s (short), GPU 1 & 3 parallel
# Baseline α=0.1 is already running as run_011
set -euo pipefail

DSPH_BIN="/opt/dsph/bin"
SIMDATA="/data/exp-c"
CASES_DIR="$HOME/slosim-agent/research-v3/exp-c/cases"
PROBES_DIR="$HOME/slosim-agent/research-v3/exp-c/probes"
SWEEP_DIR="$SIMDATA/visco_sweep"

# Visco values to test (baseline 0.1 already running)
ALPHAS=(0.5 1.0 2.0 5.0)

# Create base XML for sweep (dp=0.002, TimeMax=3s)
create_xml() {
    local alpha=$1
    local outdir=$2
    local xmlfile="$outdir/sweep.xml"

    cat > "$xmlfile" << 'XMLEOF'
<?xml version="1.0" encoding="UTF-8" ?>
<!-- Oil Visco Sweep: dp=0.002, alpha=ALPHA_PLACEHOLDER
     SPHERIC Test 10, Sunflower Oil (rho=990, mu=0.045)
     Tank: 900x62x508mm, Fill: h=93mm, Pitch f=0.651Hz A=4deg -->
<case>
    <casedef>
        <constantsdef>
            <gravity x="0" y="0" z="-9.81" />
            <rhop0 value="990" />
            <rhopgradient value="2" />
            <hswl value="0" auto="true" />
            <gamma value="7" />
            <speedsystem value="0" auto="true" />
            <coefsound value="20" />
            <speedsound value="0" auto="true" />
            <coefh value="1.2" />
            <cflnumber value="0.2" />
        </constantsdef>
        <mkconfig boundcount="240" fluidcount="9" />
        <geometry>
            <definition dp="0.002" units_comment="metres (m)">
                <pointmin x="-0.05" y="-0.01" z="-0.05" />
                <pointmax x="0.95" y="0.08" z="0.6" />
            </definition>
            <commands>
                <mainlist>
                    <setshapemode>dp | bound</setshapemode>
                    <setdrawmode mode="full" />
                    <setmkfluid mk="0" />
                    <drawbox>
                        <boxfill>solid</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="0.900" y="0.062" z="0.093" />
                    </drawbox>
                    <setmkbound mk="0" />
                    <drawbox>
                        <boxfill>bottom | left | right | front | back</boxfill>
                        <point x="0" y="0" z="0" />
                        <size x="0.900" y="0.062" z="0.508" />
                    </drawbox>
                    <shapeout file="SPHERIC_Tank" />
                </mainlist>
            </commands>
        </geometry>
        <motion>
            <objreal ref="0">
                <begin mov="1" start="0" />
                <mvrotsinu id="1" duration="3" anglesunits="degrees">
                    <axisp1 x="0.45" y="0.031" z="0" />
                    <axisp2 x="0.45" y="-0.031" z="0" />
                    <freq v="0.651" units_comment="Hz" />
                    <ampl v="4.0" units_comment="degrees" />
                    <phase v="0" units_comment="degrees" />
                </mvrotsinu>
            </objreal>
        </motion>
    </casedef>
    <execution>
        <special>
            <gauges>
                <default>
                    <savevtkpart value="false" />
                    <computedt value="0.001" units_comment="1kHz gauge" />
                    <computetime start="0" end="3" />
                    <output value="true" />
                    <outputdt value="0.001" units_comment="1kHz CSV" />
                </default>
                <swl name="SWL_Left">
                    <point0 x="0.007" y="0.031" z="0.0" />
                    <point2 x="0.007" y="0.031" z="0.508" />
                    <pointdp coefdp="0.5" />
                    <masslimit value="0.2" />
                </swl>
            </gauges>
        </special>
        <parameters>
            <parameter key="SavePosDouble" value="0" />
            <parameter key="StepAlgorithm" value="2" comment="2:Symplectic" />
            <parameter key="Kernel" value="2" comment="2:Wendland" />
            <parameter key="ViscoTreatment" value="1" comment="1:Artificial" />
            <parameter key="Visco" value="ALPHA_PLACEHOLDER" />
            <parameter key="ViscoBoundFactor" value="1" />
            <parameter key="DensityDT" value="2" comment="2:Fourtakas" />
            <parameter key="DensityDTvalue" value="0.1" />
            <parameter key="Shifting" value="0" />
            <parameter key="RigidAlgorithm" value="1" />
            <parameter key="DtIni" value="0.0001" />
            <parameter key="DtMin" value="0.00001" />
            <parameter key="CoefDtMin" value="0.05" />
            <parameter key="TimeMax" value="3.0" units_comment="seconds (short sweep)" />
            <parameter key="TimeOut" value="0.1" units_comment="seconds (31 PART files)" />
            <parameter key="MinFluidStop" value="0" />
            <parameter key="RhopOutMin" value="700" />
            <parameter key="RhopOutMax" value="1300" />
            <simulationdomain>
                <posmin x="default - 15%" y="default" z="default - 10%" />
                <posmax x="default + 15%" y="default" z="default + 80%" />
            </simulationdomain>
        </parameters>
    </execution>
</case>
XMLEOF
    # Replace placeholder with actual alpha
    sed -i "s/ALPHA_PLACEHOLDER/$alpha/g" "$xmlfile"
}

# Run single case on specified GPU
run_case() {
    local alpha=$1
    local gpu=$2
    local name="alpha_${alpha}"
    local outdir="$SWEEP_DIR/$name"

    mkdir -p "$outdir/out"
    create_xml "$alpha" "$outdir"
    cp "$PROBES_DIR/pressure_probes_oil_dp001.txt" "$outdir/pressure_probes.txt"

    echo "[GPU $gpu] Starting $name at $(date)"

    # GenCase
    docker run --rm --gpus "device=$gpu" \
        -v "$outdir:/data" -w /data --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/GenCase_linux64 /data/sweep /data/out/sweep -save:all \
        > "$outdir/run.log" 2>&1

    # Solver
    docker run --rm --gpus "device=$gpu" \
        -v "$outdir:/data" -w /data --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/DualSPHysics5.4_linux64 /data/out/sweep /data/out -gpu -svres \
        >> "$outdir/run.log" 2>&1

    # MeasureTool
    docker run --rm --gpus "device=$gpu" \
        -v "$outdir:/data" -w /data --entrypoint '' \
        dsph-agent:latest \
        $DSPH_BIN/MeasureTool_linux64 \
            -dirin /data/out \
            -points /data/pressure_probes.txt \
            -onlytype:-all,+fluid \
            -vars:+press,+vel \
            -savecsv /data/PressResult \
        >> "$outdir/run.log" 2>&1

    echo "[GPU $gpu] $name DONE at $(date)"
    if [ -f "$outdir/PressResult_Press.csv" ]; then
        lines=$(wc -l < "$outdir/PressResult_Press.csv")
        echo "  CSV lines: $lines"
    else
        echo "  WARNING: No CSV output!"
    fi
}

echo "=== Oil Visco Parametric Sweep ==="
echo "Alpha values: ${ALPHAS[*]}"
echo "GPU 1 & 3 parallel"
echo ""

# Phase 1: GPU 1 = α=0.5, GPU 3 = α=1.0 (parallel)
run_case "${ALPHAS[0]}" 1 &
PID1=$!
run_case "${ALPHAS[1]}" 3 &
PID2=$!
wait $PID1 $PID2
echo ""

# Phase 2: GPU 1 = α=2.0, GPU 3 = α=5.0 (parallel)
run_case "${ALPHAS[2]}" 1 &
PID3=$!
run_case "${ALPHAS[3]}" 3 &
PID4=$!
wait $PID3 $PID4

echo ""
echo "=== ALL SWEEP DONE at $(date) ==="
echo "Results in: $SWEEP_DIR/alpha_*/PressResult_Press.csv"
