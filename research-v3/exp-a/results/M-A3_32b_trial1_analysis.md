# EXP-A M-A3 Parameter Fidelity Results — Qwen3-32B, Trial 1

**Date**: 2026-03-02
**Infrastructure**: pajulab (4× RTX 4090), Ollama qwen3:32b
**Note**: S01 oliveeelab timeout (3600s) — pajulab re-run in progress

## Summary Table

| Scenario | Tier | Duration | M-A3 Score | Key Failures |
|----------|------|----------|------------|--------------|
| S01 | Easy | 157s | **64%** (7/11) | density (1000 vs 990), viscosity, motion_type, amplitude |
| S02 | Easy | 99s | **100%** (9/9) | None — perfect |
| S03 | Medium | 109s | **89%** (8/9) | timemax (4.96 vs 10) |
| S04 | Medium | 120s | **78%** (7/9) | motion_type (rectsinu vs rotsinu), amplitude |
| S05 | Medium | 150s | **67%** (6/9) | tank_z, fill_height, motion_type |
| S06 | Medium | 135s | **89%** (8/9) | timemax (8.33 vs 15) |
| S07 | Hard | 151s | **63%** (5/8) | geometry (box vs cylinder), timemax, dp |
| S08 | Easy | 180s | **50%** (5/10) | motion_type, freq, amplitude, timemax |
| S09 | Hard | 153s | **44%** (4/9) | geometry (box vs cylinder), motion_type, amplitude, timemax |
| S10 | Hard | 213s | **N/A** | STL file not provided — structural limitation |

**Overall (excl S10)**: Mean M-A3 = **71.6%** (9 scenarios, 644/900)
**By Tier**: Easy avg=71.3% (S01:64, S02:100, S08:50), Medium avg=80.8% (S03:89, S04:78, S05:67, S06:89), Hard avg=53.5% (S07:63, S09:44)

## Detailed Per-Scenario Analysis

### S02 — Chen Shallow Sway (Easy) ✅ 100%
- tank (0.6×0.3×0.65): PASS
- fill_height (0.083m): PASS
- density (1000): PASS
- motion_type (mvrectsinu): PASS
- frequency (0.756Hz): PASS
- amplitude (0.007m): PASS
- timemax (10s): PASS
- dp (0.005, range 0.003-0.01): PASS
- gauges: PASS

### S03 — Chen Near-Critical (Medium) ✅ 89%
- tank (0.6×0.3×0.65): PASS
- fill_height (0.185m): PASS
- density (1000): PASS
- motion_type (mvrectsinu): PASS
- frequency (1.008Hz): PASS
- amplitude (0.007m): PASS
- timemax: **FAIL** (4.96 vs 10) — agent halved simulation time
- dp (0.005): PASS
- gauges: PASS

### S04 — Liu Large Pitch 30s (Medium) ⚠️ 78%
- tank (1.0×0.5×1.0): PASS
- fill_height (0.2m): PASS
- density (1000): PASS
- motion_type: **FAIL** (mvrectsinu vs mvrotsinu) — pitch→translation confusion
- frequency (0.66Hz): PASS
- amplitude: **FAIL** (linear amplitude vs 2° rotation)
- timemax (30s): PASS
- dp (0.01): PASS
- gauges: PASS

### S05 — Liu Amplitude Parametric (Medium) ⚠️ 67%
- tank (1.0×0.5): PASS, z: **FAIL** (0.5 vs 1.0)
- fill_height: **FAIL** (0.35 vs 0.7m)
- density (1000): PASS
- motion_type: **FAIL** (mvrectsinu vs mvrotsinu)
- frequency (0.87Hz): PASS
- parametric study: PARTIAL (2 cases vs 3)
- dp (0.01): PASS
- gauges: PASS

### S06 — ISOPE LNG Roof Impact (Medium) ✅ 89%
- tank (0.946×0.118×0.67): PASS
- fill_height (0.603m, 90%): PASS
- density (1000): PASS
- motion_type (mvrectsinu): PASS
- frequency (0.6Hz): PASS
- amplitude (0.02m): PASS
- timemax: **FAIL** (8.33 vs 15) — agent calculated 5 periods
- dp (0.005): PASS
- gauges: PASS

### S07 — NASA Cylinder Baffle (Hard) ⚠️ 63%
- geometry: **FAIL** (drawbox instead of drawcylinder) — core Hard-tier test
- fill_height (1.5m): PASS
- density (1000): PASS
- motion_type (mvrectsinu): PASS (horizontal sway correct)
- frequency (0.5Hz): PASS
- amplitude (0.01m): PASS
- timemax: **FAIL** (10 vs 15)
- dp: **FAIL** (0.05, range 0.01-0.03)

### S08 — English mDBC/DBC (Easy) ⚠️ 50%
- tank (0.9×0.062×0.508): PASS
- fill_height (0.09144m): PASS
- density (1000): PASS
- DBC boundary: PASS
- dp (0.002): PASS
- motion_type: **FAIL** (mvrectsinu vs mvrotsinu) — SPHERIC benchmark pitch
- frequency: **FAIL** (0.519 vs 0.6515 Hz)
- amplitude: **FAIL** (0.045m linear vs 4° rotation)
- timemax: **FAIL** (9.63 vs 7.68)
- Note: Prompt didn't specify motion params — tested domain knowledge of SPHERIC benchmark

### S09 — Horizontal Cylinder Pitch (Hard) ⚠️ 44%
- geometry: **FAIL** (drawbox 3×1×1 instead of horizontal cylinder)
- fill_ratio (25%): PASS (0.25m in 1m height box)
- density (1000): PASS
- motion_type: **FAIL** (mvrectsinu vs mvrotsinu)
- frequency (0.55Hz): PASS
- amplitude: **FAIL** (0.0785m linear vs 3° rotation)
- timemax: **FAIL** (9.09 vs 15)
- dp (0.02): PASS (edge of range 0.008-0.02)

### S10 — STL Fuel Tank (Hard) — N/A
- Agent correctly identified STL file not found
- No XML generated
- Scored as N/A (structural limitation — file not available)

## Systematic Error Patterns

### P1: Pitch Rotation Confusion (affects 5/10 scenarios)
Agent consistently uses `mvrectsinu` (horizontal translation) when `mvrotsinu` (pitch rotation) is required.
- **Affected**: S01, S04, S05, S08, S09
- **Root cause**: System prompt doesn't sufficiently distinguish rotation motion commands

### P2: TimeMax Underestimation (affects 4/10 scenarios)
Agent calculates TimeMax as a small number of periods rather than using the specified value.
- **Affected**: S03, S06, S07, S09
- **Pattern**: Agent computes ~5 oscillation periods instead of user-specified simulation time

### P3: Cylindrical Geometry Not Supported (affects 2/10 scenarios)
Agent defaults to rectangular drawbox when cylindrical geometry is requested.
- **Affected**: S07, S09
- **Root cause**: xml_generator tool only supports rectangular geometry

### P4: SPHERIC Benchmark Domain Knowledge Gap (affects 1 scenario)
When prompt references SPHERIC benchmark without explicit motion params, agent fabricates values.
- **Affected**: S08
- **Root cause**: No embedded knowledge of standard SPHERIC benchmark parameters

### P5: Amplitude Unit Confusion (co-occurs with P1)
When rotation motion is needed, agent converts degree amplitude to a linear displacement.
- **Affected**: S04, S05, S08, S09
