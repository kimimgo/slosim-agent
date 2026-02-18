# Table 2: SPHERIC Test 10 — Peak Pressure Validation

| Case | Particles | dp [m] | Output [Hz] | Sim Peaks | Peaks in ±2σ | Max P [mbar] |
|------|-----------|--------|-------------|-----------|-------------|-------------|
| Water Low | 136K | 0.004 | 200 | 31.1, 58.9, 76.7 | 3/3 | 76.7 |
| Water High | 344K | 0.004 | 100 | 44.2, 29.4, 31.4, 45.3 | 4/4 | 50.0 |
| Oil Low | 136K | 0.004 | 200 | — | N/A (no peaks) | 0.0 |

## Experimental Reference (100-repeat statistics)

| | Peak 1 | Peak 2 | Peak 3 | Peak 4 |
|-----|--------|--------|--------|--------|
| Water μ [mbar] | 37.1 | 48.2 | 46.9 | 46.4 |
| Water ±2σ | ±14.6 | ±29.9 | ±34.0 | ±26.3 |
| Oil μ [mbar] | 6.9 | 6.5 | 5.4 | 5.5 |
| Oil ±2σ | ±0.3 | ±0.5 | ±0.5 | ±0.5 |

## Key Findings

- **Water Low/High**: All detected peaks fall within experimental ±2σ band
- **Oil Low**: No impact peaks detected at H=93mm sensor location — DBC + artificial viscosity over-damps oil sloshing
- **Implication**: DBC boundary condition adequate for water, mDBC recommended for viscous fluids (cf. English et al., 2021)
