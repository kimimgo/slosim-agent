#!/bin/bash
# Run 006 v2 — Post-processing with MOVING probes (1.5h offset, rotating)
# Prerequisites: solver must have completed (t=0 → 7s, no partbegin)
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
cd "$ROOT"

LOG="research-v2/exp1_spheric/scripts/postprocess_run006_v2.log"
log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG"; }

log "=== Run 006 v2 Post-Processing ==="

# Verify solver completed
FINAL_T=$(tail -1 simulations/exp1/run_006/RunPARTs.csv 2>/dev/null | cut -d';' -f2 | head -c 5)
log "Final simulation time: ${FINAL_T:-unknown}"

# Step 1: Generate moving probe positions (Python)
log "=== Step 1: Generate moving probe CSV ==="
python3 << 'PYEOF'
import numpy as np

pivot_x, pivot_z = 0.45, 0.0
freq, ampl_deg = 0.856, 4.0
ampl_rad = np.radians(ampl_deg)
y_probe = 0.031

dp = 0.002
coefh = 1.2
h = coefh * np.sqrt(3) * dp
offset = 1.5 * h
roof_z = 0.508
z_tank = roof_z - offset  # 0.50176

x_probes = [0.05, 0.225, 0.45, 0.675, 0.85]

dt_part = 0.001
n_parts = 7001
times = np.arange(n_parts) * dt_part

with open("/tmp/roof_moving_probes_v2.csv", "w") as f:
    for t in times:
        theta = ampl_rad * np.sin(2 * np.pi * freq * t)
        cos_t, sin_t = np.cos(theta), np.sin(theta)
        parts = []
        for x_tank in x_probes:
            dx = x_tank - pivot_x
            x_g = pivot_x + dx * cos_t - z_tank * sin_t
            z_g = dx * sin_t + z_tank * cos_t
            parts.extend([f"{x_g:.6f}", f"{y_probe:.6f}", f"{z_g:.6f}"])
        f.write(";".join(parts) + "\n")
print(f"Generated /tmp/roof_moving_probes_v2.csv: {n_parts} lines, {len(x_probes)} probes")
PYEOF
log "Moving probes CSV generated"

# Step 2: MeasureTool — Roof moving probes
log "=== Step 2: MeasureTool (Roof Moving Probes) ==="
docker compose run --rm \
    -v /tmp/roof_moving_probes_v2.csv:/cases/roof_moving_probes.csv \
    dsph measuretool \
    -dirin /data/exp1/run_006 \
    -pointspos /cases/roof_moving_probes.csv \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_006/PressRoofMoving \
    >> "$LOG" 2>&1
log "MeasureTool Roof Moving done"

# Step 3: MeasureTool — Lateral probes (consistent with Run 002)
log "=== Step 3: MeasureTool (Lateral, fixed) ==="
cat > /tmp/lateral_probe_v2.txt << 'PROBES'
POINTS
0.050 0.031 0.093
PROBES

docker compose run --rm \
    -v /tmp/lateral_probe_v2.txt:/cases/lateral_probe.txt \
    dsph measuretool \
    -dirin /data/exp1/run_006 \
    -points /cases/lateral_probe.txt \
    -onlytype:-all,+fluid \
    -vars:press,vel \
    -savecsv /data/exp1/run_006/PressConsistent \
    >> "$LOG" 2>&1
log "MeasureTool Lateral done"

# Step 4: Quick analysis
log "=== Step 4: Quick Analysis ==="
python3 -c "
import numpy as np

# Roof moving probes
d = np.genfromtxt('simulations/exp1/run_006/PressRoofMoving_Press.csv', delimiter=';', skip_header=4)
print('=== Roof Moving Probes ===')
xs = [0.05, 0.225, 0.45, 0.675, 0.85]
for i, x in enumerate(xs):
    p = d[:, i+2] / 100
    print(f'  x={x}: max={p.max():.1f} mbar, nonzero={np.sum(np.abs(p)>0.1)}')

# Experiment comparison
exp = np.genfromtxt('datasets/spheric/case_1/roof_water_1x.txt', delimiter='\t', skip_header=1)
print(f'  Experiment max: {exp[:,1].max():.1f} mbar')
" >> "$LOG" 2>&1
log "Analysis done"

# Step 5: Regenerate paper figures
log "=== Step 5: Paper Figures ==="
python3 research-v2/exp1_spheric/analysis/paper_figures.py >> "$LOG" 2>&1
log "Figures regenerated"

log "=== ALL DONE ==="
