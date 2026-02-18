# Sloshing Simulation Datasets

Research-grade datasets for sloshing simulation validation and testing.

## Structure

```
datasets/
├── spheric/           SPHERIC SPH Benchmark (Souto-Iglesias et al.)
│   ├── case_1/        Pressure impact data (oil + water, 100 repeats)
│   ├── case_2/        TLD (Tuned Liquid Damper) data
│   └── case_3/        FSI (Fluid-Structure Interaction) data
├── dualsphysics/      DualSPHysics official examples
│   ├── Examples_main.pdf
│   └── Examples_motion.pdf
├── papers/            Open access research papers
│   ├── ISOPE_2012_sloshing_benchmark.pdf    Mark III LNG multi-institution
│   ├── English_2021_mDBC_sloshing.pdf       DualSPHysics mDBC validation
│   ├── Chen_2018_OpenFOAM_sloshing.pdf      Fill level parametric study
│   ├── Frosina_2018_fuel_tank_sloshing.pdf  Automotive CAD geometry + CFD
│   ├── Liu_2024_pitch_sloshing.pdf          Pitch excitation experimental
│   ├── Zhao_2024_LNG_baffles_SPHinXsys.pdf LNG elastic baffles FSI
│   └── NASA_2023_antislosh_baffle.pdf       Ring baffle pressure model
└── ccp-wsi/           CCP-WSI Blind Test Series (access via website)
```

## Key References (for paper citation)

1. **SPHERIC Benchmark** — Souto-Iglesias, Botia-Vera et al. (2011-2015)
   - FTP: http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/
   - DOI: 10.1016/j.oceaneng.2015.05.013

2. **DualSPHysics SloshingTank** — Official test case 05
   - GitHub: https://github.com/DualSPHysics/DualSPHysics
   - Full examples available in binary distribution

3. **CCP-WSI Blind Test 5** — Circular tank sloshing (2025)
   - Data: https://ccp-wsi.ac.uk/data_repository/

4. **Jiao et al. (2024)** — LNG + DualSPHysics coupling
   - DOI: 10.1016/j.oceaneng.2024.119148
   - DOI: 10.1016/j.oceaneng.2024.117022

5. **ISOPE Sloshing Model Test Benchmark (2012/2013)**
   - Loysel, Chollet, Gervaise, Brosset (GTT)

## Notes

- Large binary files (PDFs, data) are in .gitignore
- SPHERIC case_1 contains ~35MB of experimental pressure time series
- Jiao papers are Elsevier (not open access) — only landing page HTMLs stored
- DualSPHysics full sloshing examples require downloading the binary package from https://dual.sphysics.org/downloads/
