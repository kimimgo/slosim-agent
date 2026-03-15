#!/bin/bash
# Post-process Run 007 (Water Roof, Visco=0.001)
# 1. Quick test: PartVTK → check max particle z (need >505mm for roof impact)
# 2. If success: MeasureTool with roof probes → analysis pipeline
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

LOG="research-v2/exp1_spheric/scripts/postprocess_run007.log"
log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG"; }

RUN_DIR="simulations/exp1/run_007"

# Step 1: Check max particle height via PartVTK at key timesteps
log "=== Step 1: Check max fluid particle z ==="

# Extract fluid particles at t=1.0, 1.5, 2.0, 2.5 to find peak wave height
MAX_Z=0
PEAK_T=0
for PART in 1000 1500 1800 2000 2200 2500 2800 3000; do
    if [ ! -f "$RUN_DIR/Part_${PART}.bi4" ] 2>/dev/null; then
        # Try zero-padded name
        PADDED=$(printf "%04d" $PART)
        if [ ! -f "$RUN_DIR/Part_${PADDED}.bi4" ] 2>/dev/null; then
            log "  Part $PART not found, skipping"
            continue
        fi
    fi

    docker compose run --rm dsph partvtk \
        -dirin /data/exp1/run_007 \
        -savevtk /data/exp1/run_007/check_height \
        -onlytype:-all,+fluid \
        -first:$PART -last:$PART \
        > /dev/null 2>&1 || continue

    # Parse VTK for max z coordinate
    VTK_FILE="$RUN_DIR/check_height_$(printf '%04d' $PART).vtk"
    if [ -f "$VTK_FILE" ]; then
        Z_MAX=$(python3 -c "
import struct
with open('$VTK_FILE', 'rb') as f:
    content = f.read()
# Find POINTS header
idx = content.find(b'POINTS')
if idx < 0:
    print('0')
else:
    line_end = content.index(b'\n', idx)
    header = content[idx:line_end].decode('ascii')
    n_points = int(header.split()[1])
    data_start = line_end + 1
    max_z = 0
    for i in range(n_points):
        offset = data_start + i * 12  # 3 floats * 4 bytes
        x, y, z = struct.unpack('>fff', content[offset:offset+12])
        if z > max_z:
            max_z = z
    print(f'{max_z*1000:.1f}')
" 2>/dev/null || echo "0")

        T_SEC=$(echo "scale=3; $PART / 1000" | bc)
        log "  Part $PART (t=${T_SEC}s): max z = ${Z_MAX} mm"

        # Track peak
        if [ "$(echo "$Z_MAX > $MAX_Z" | bc -l)" = "1" ]; then
            MAX_Z=$Z_MAX
            PEAK_T=$T_SEC
        fi

        # Clean up temp VTK
        rm -f "$VTK_FILE"
    fi
done

log "  Peak height: ${MAX_Z} mm at t=${PEAK_T}s (roof=508mm)"

# Step 2: Decision based on max z
if [ "$(echo "$MAX_Z > 505" | bc -l)" = "1" ]; then
    log "SUCCESS: Wave reaches roof! Proceeding to MeasureTool..."
elif [ "$(echo "$MAX_Z > 500" | bc -l)" = "1" ]; then
    log "MARGINAL: Wave close to roof (${MAX_Z}mm). Proceeding with MeasureTool anyway..."
else
    log "FAIL: Wave still too short (${MAX_Z}mm). Need alternative approach."
    log "  Consider: Visco=0.005, or dp=0.001, or mDBC, or accept PARTIAL"
    exit 1
fi

# Step 3: MeasureTool — roof probes at z=0.503 (5mm below roof)
log "=== Step 3: MeasureTool (Roof Probes) ==="
cat > /tmp/roof_probe_007.txt << 'PROBES'
POINTS
0.050 0.031 0.503
0.150 0.031 0.503
0.225 0.031 0.503
0.450 0.031 0.503
0.675 0.031 0.503
0.750 0.031 0.503
0.850 0.031 0.503
PROBES

docker compose run --rm \
    -v /tmp/roof_probe_007.txt:/cases/roof_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_007 \
    -points /cases/roof_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_007/PressRoof \
    >> "$LOG" 2>&1
log "MeasureTool Roof done"

# Step 4: MeasureTool — lateral probes (same as Run 002 for consistency)
log "=== Step 4: MeasureTool (Lateral Probes) ==="
docker compose run --rm dsph measuretool \
    -dirin /data/exp1/run_007 \
    -points /cases/spheric_probes.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_007/PressConsistent \
    >> "$LOG" 2>&1
log "MeasureTool Lateral done"

# Step 5: Analysis
log "=== Step 5: Running analysis ==="
python3 research-v2/exp1_spheric/analysis/oil_roof_analysis.py >> "$LOG" 2>&1
python3 research-v2/exp1_spheric/analysis/paper_figures.py >> "$LOG" 2>&1
log "Analysis done"

# Step 6: Final Verdict
log "=== Step 6: Final Verdict ==="
python3 research-v2/exp1_spheric/scripts/final_verdict.py | tee -a "$LOG"

log "=== POSTPROCESS RUN 007 COMPLETE ==="
