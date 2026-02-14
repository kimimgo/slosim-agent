# Sloshing Simulation STL Datasets & Tank Parametric Specification

## Part A: Downloadable Tank STL/CAD Datasets

### 1. Rectangular / Storage Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| GrabCAD Water Tank Collection | [Link](https://grabcad.com/library?time=all_time&sort=recent&query=water+tank) | STEP, IGES, SLDPRT | Yes (account) | 수백 개 직사각형/원통형 물탱크 | ★★★★★ 내부 캐비티 명확 |
| GrabCAD Modular Water Tank | [Link](https://grabcad.com/library/modular-water-tank-1) | SLDPRT, STEP | Yes | 모듈형 직사각형 물탱크 | ★★★★ 칸막이 구조 |
| GrabCAD STEP/IGES Tank Collection | [Link](https://grabcad.com/library?page=5&softwares=step-slash-iges&sort=most_downloaded&tags=tank) | STEP, IGES | Yes | CAD 네이티브 탱크 모음 | ★★★★★ 정확한 형상 |
| TraceParts Industrial Tanks | [Link](https://www.traceparts.com/en/search/traceparts-classification-material-handling-and-lifting-equipment-material-handling-conveyors-systems-packaging-storage-equipment-tanks) | STEP, STL, IGES, DWG 등 | Yes (account) | 산업용 탱크 카탈로그 (실물 규격) | ★★★★★ 실제 치수 |

### 2. Cylindrical Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| Cults3D Horizontal Storage Tank | [Link](https://cults3d.com/en/3d-model/home/horizontal-storage-tank-industrial-3d-model) | STL | Paid | 수평 원통형 저장탱크 | ★★★★ 공명 주파수 연구 |
| CGTrader Horizontal Vessel | [Link](https://www.cgtrader.com/free-3d-models/industrial/other/horizontal-vessel) | SLDPRT, SLDASM | Yes | 수평 압력용기 (디쉬 엔드) | ★★★★ 벤치마크 |
| CGTrader Pressure Vessel | [Link](https://www.cgtrader.com/free-3d-models/industrial/other/pressure-vessel) | SLDPRT, SLDASM | Yes | 수직/수평 압력용기 | ★★★★ 내부 캐비티 명확 |
| CGTrader Free Vessel STL (39개) | [Link](https://www.cgtrader.com/free-3d-print-models/vessel) | STL, OBJ | Yes | 다양한 vessel 모델 | ★★★ 내부 형상 확인 필요 |

### 3. LNG Carrier Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| Free3D LNG Moss Sphere | [Link](https://free3d.com/3d-model/lng-cargo-tank-k-moss-sphere-3815.html) | LWO, BLEND, OBJ | $20 | Moss 구형 LNG 탱크 (cutaway) | ★★★★★ 공명 모드 검증 |
| GrabCAD LNG Sphere Tank | [Link](https://grabcad.com/library/lng-sphere-tank-1) | STEP, SLDPRT | Yes | LNG 구형 저장탱크 | ★★★★★ 무료 CAD |
| GrabCAD LNG Tag Collection | [Link](https://grabcad.com/library/tag/lng) | 다양 | Yes | LNG 운반선 관련 모델 | ★★★ 탱크 캐비티 추출 필요 |

> **Note**: Mark III / NO96 membrane tank의 상세 3D 모델은 GTT 독점 기술로 공개 모델 없음. 논문 치수 기반 파라메트릭 생성 필요.

### 4. Ship / Ballast Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| GrabCAD Fore Peak Ballast Tank | [Link](https://grabcad.com/library/fore-peak-tank-forward-ballast-tank-1) | STEP, IGES | Yes | 선박 전방 밸러스트 탱크 | ★★★★★ bulkhead 구조 |
| GrabCAD Fuel Tank FSAE | [Link](https://grabcad.com/library/fuel-tank-fsae-1) | SLDPRT, STEP | Yes | FSAE 레이싱카 연료탱크 | ★★★★ 비정형 형상 |

### 5. Chemical / Industrial Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| GrabCAD IBC Bulk Container | [Link](https://grabcad.com/library/ibc-bulk-container-1) | STEP, IGES, SLDPRT | Yes | 1000L IBC 컨테이너 | ★★★★ 운송 슬로싱 |
| GrabCAD Stirred Tank Reactor | [Link](https://grabcad.com/library/continuous-stirred-tank-reactor) | STEP, IGES | Yes | 교반 탱크 (baffle + impeller) | ★★★★★ anti-slosh 연구 |

### 6. Aerospace / Propellant Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| NASA 3D Resources | [Link](https://www.nasa.gov/3d-resources/) / [GitHub](https://github.com/nasa/NASA-3D-Resources) | OBJ, STL, 3DS | Yes (public domain) | 우주선/로켓 모델 (PMD 탱크 포함) | ★★★★★ 검증 데이터 풍부 |
| NASA NTRS Propellant Sloshing | [Link](https://ntrs.nasa.gov/api/citations/20170000667/downloads/20170000667.pdf) | PDF (치수만) | Yes | PMD 설계 상세 | ★★★★★ 재구성 가능 |

### 7. Automotive Fuel Tanks

| Source | URL | Format | Free | Description | Sloshing Suitability |
|--------|-----|--------|------|-------------|---------------------|
| Cults3D Fuel Tank | [Link](https://cults3d.com/en/tags/fuel+tank) | STL | Mixed | 자동차/오토바이 연료탱크 | ★★★★ baffle 포함 가능 |
| CGTrader Fuel Tank STL (2,494개) | [Link](https://www.cgtrader.com/3d-print-models/fuel-tank) | STL, OBJ | Mixed | 대규모 연료탱크 컬렉션 | ★★★★ 선별 필요 |

### DualSPHysics STL 사용 시 주의사항

1. **Watertight mesh 필수** — GenCase는 닫힌 메시만 처리 가능 (Blender/MeshLab로 검증)
2. **단위 확인** — DualSPHysics는 SI(m) 사용, 3D 프린팅 모델은 mm인 경우 많음
3. **내부 캐비티** — 시각화용 모델은 외피만 있을 수 있으므로 Boolean 연산으로 내부 추출
4. **STEP > STL** — 가능하면 STEP/IGES 선호 (정확한 곡면, 해상도 손실 없음)
5. **`-loadstl` + `fill`** — GenCase XML에서 STL 로딩 후 파티클 채우기 필수

---

## Part B: Tank Parametric Design Specification

### 1. Rectangular Tank

```
┌─────────────────────────────┐
│           W                 │
│  ┌─────────────────────┐    │
│  │                     │ H  │  t_wall
│  │    ┃ baffle         │    │
│  │    ┃ h_b            │    │
│  │    ┃                │    │
│  └─────────────────────┘    │
│           L                 │
└─────────────────────────────┘
```

#### Basic Dimensions

| Parameter | Symbol | Range | Unit | Notes |
|-----------|--------|-------|------|-------|
| Length | `L` | 0.5 - 5.0 | m | 여기(excitation) 방향 |
| Width | `W` | 0.3 - 3.0 | m | 횡방향 |
| Height | `H` | 0.4 - 4.0 | m | |
| Wall thickness | `t_wall` | 3 - 15 | mm | |
| Corner radius | `r_corner` | 0 - 20 | mm | CFD 메시 품질용 |
| Aspect ratio L/W | | 1.0 - 3.0 | - | 일반적 1.5 - 2.0 |
| Depth ratio H/L | | 0.3 - 1.5 | - | shallow < 0.5, deep > 1.0 |

#### Fill Level (Critical for Sloshing)

| Fill Ratio h/H | Sloshing Severity | Research Priority |
|----------------|-------------------|-------------------|
| 10% | Low amplitude, high frequency | Medium |
| 30% | Moderate | High |
| **40-80%** | **Worst-case range** | **Critical** |
| 90% | Roof impact risk | Medium |

#### Natural Frequency (1st mode)

```
ω₁ = √(g · k₁ · tanh(k₁ · h))
k₁ = π / L

Shallow (h/L < 0.1): ω₁ ≈ √(g·π·h/L)
Deep    (h/L > 1.0): ω₁ ≈ √(g·π/L)
```

#### Baffle Parameters

**Vertical Baffle:**

| Parameter | Symbol | Optimal | Notes |
|-----------|--------|---------|-------|
| Height ratio | `h_b/h_w` | 0.5 - 0.9 | ≥0.8 most effective |
| Position | `x_b/L` | 0.5 | Center is optimal |
| Thickness | `t_b` | 2 - 10 mm | |
| Count | `n_b` | 1 - 3 | |

**Horizontal Baffle:**

| Parameter | Symbol | Optimal | Notes |
|-----------|--------|---------|-------|
| Height position | `y_b/H` | 0.3 - 0.7 | Just below liquid surface |
| Length ratio | `l_b/L` | 0.6 - 0.9 | |

**Perforated Baffle:**

| Parameter | Symbol | Optimal | Notes |
|-----------|--------|---------|-------|
| Porosity | `ε` | 10 - 20% | **20% optimal** |
| Hole diameter | `d_hole` | 10 - 50 mm | |
| Hole spacing | `s_hole` | 2d - 5d | |
| Pattern | | Uniform | Better than non-uniform |

---

### 2. Cylindrical Tank — Upright

```
       ╭──────────╮  ← Roof (cone/dome/floating)
      │            │
      │            │ H
      │            │
      │     D      │
      └────────────┘  ← Bottom (flat/conical/dished)
```

#### Basic Dimensions (API 650)

| Parameter | Symbol | Range | Unit | Notes |
|-----------|--------|-------|------|-------|
| Diameter | `D` | 3 - 80 | m | |
| Height | `H` | 5 - 25 | m | |
| H/D ratio | | 0.5 - 1.5 | - | Optimal: 0.8 |
| Shell thickness | `t_shell` | 6 - 45 | mm | API 650 max: 45mm |
| Corrosion allowance | `CA` | 1.5 - 3 | mm | |

#### Bottom Head Types

| Type | Key Parameter | Formula |
|------|--------------|---------|
| Flat | `t_min` = 4.76 mm | API 650 |
| Conical | cone angle θ = 30°-60° | `h_cone = (D/2)·tan(θ)` |
| Torispherical (ASME F&D) | L = 1.0D, r = 0.06D | `M = 0.25·(3+√(L/r))` |
| 2:1 Ellipsoidal | depth = D/4 | L ≈ 0.9D, r ≈ 0.17D |
| Hemispherical | depth = D/2 | R = D/2 |

#### Roof Types

| Type | Thickness | Notes |
|------|-----------|-------|
| Cone Roof | min 5 mm | Slope 9.5° - 37° |
| Dome Roof | min 5 mm | R = 0.8D ~ 1.2D |
| Floating Roof | — | Minimizes evaporation |

#### Internal Components

- **Stiffener Rings**: max diameter 52m, wind girder at top
- **Agitator Baffles**: 4ea at 90° spacing, width = 0.1D

#### Natural Frequency (cylindrical)

```
k₁ = 1.841 / R    (1st Bessel root / radius)
ω₁ = √(g · k₁ · tanh(k₁ · h))
```

---

### 3. Cylindrical Tank — Horizontal

```
    ┌──╮          ╭──┐
    │  │          │  │
    │  │   L      │  │ D
    │  │          │  │
    └──╯          ╰──┘
   Head          Head
  ▽ Saddle     ▽ Saddle
```

#### Basic Dimensions

| Parameter | Symbol | Range | Notes |
|-----------|--------|-------|-------|
| Diameter | `D` | 0.5 - 4.0 m | |
| Length | `L` | 2 - 15 m | |
| L/D ratio | | 2 - 6 | Optimal: 3 - 4 |
| Wall thickness | `t` | 3 - 20 mm | ASME VIII |

#### Head Types

Same as upright: Flat, 2:1 Ellipsoidal (most common), ASME F&D, Hemispherical

#### Saddle Supports (Zick Method)

| Parameter | Symbol | Value | Notes |
|-----------|--------|-------|-------|
| Position ratio | `A/R` | 0.25 | Min stress point |
| Saddle angle | `θ` | 120° - 150° | Usually 120° |
| Saddle width | `b` | 0.1L - 0.2L | |
| Count | | 2 | Standard |

#### Internal Baffles

- **Swash Plates**: at L/3, 2L/3; opening area 30-50% of cross-section
- **Transverse Baffles**: spacing 0.3L - 0.5L; limber holes (bottom) + air holes (top)

---

### 4. LNG Membrane Tank (Mark III / NO96)

```
  ╱‾‾‾‾‾‾‾‾‾‾‾‾‾╲    ← Upper chamfer (135°)
 │                 │
 │                 │ H
 │                 │
  ╲_______________╱    ← Lower chamfer (135°)
         L
   Octagonal cross-section
```

#### Tank Geometry

| Parameter | Range | Notes |
|-----------|-------|-------|
| Total capacity | 120,000 - 270,000 m³ | Full LNGC |
| Common capacity | 125,000 - 135,000 m³ | Most frequent |
| Cross-section | Octagonal | 8-sided |
| Upper chamfer height | ~8.63 m | Example value |
| Lower chamfer height | ~3.77 m | Example value |
| Chamfer angle | 135° | Upper & lower |
| Tank count | 4 - 6 | Centerline arrangement |

#### Membrane Systems

| System | Primary | Thickness | Secondary | Insulation |
|--------|---------|-----------|-----------|------------|
| Mark III | SS 304L corrugated | 1.2 mm | Triplex | PUF 160mm |
| NO96 | Invar (36% Ni) | 0.7 mm | Invar 0.7mm | Plywood + Perlite |

#### Design Standards
- IGC Code: Membrane type, full secondary barrier (< -10°C)
- Design life: 10⁸ wave encounters (20 years, North Atlantic)

---

### 5. LNG Moss Spherical Tank

```
       ╭───╮
      ╱     ╲
     │       │ D = 36-42m
      ╲     ╱
       ╰─┬─╯
         │ Skirt (equator attachment)
    ─────┴─────  Hull
```

| Parameter | Symbol | Range | Notes |
|-----------|--------|-------|-------|
| Sphere diameter | `D` | 36 - 42 m | 125k m³ ship: 36.6m |
| Shell thickness | `t` | up to 50 mm | Aluminium alloy |
| Tank count | `n` | 4 - 5 | Usually 5 |
| Material | | Al 5083/5454 or 9% Ni steel | |
| Classification | | IMO Type B | Partial secondary barrier |
| Support | | Cylindrical skirt at equator | |

---

### 6. IMO Type B Prismatic Tank (SPB — IHI)

| Parameter | Value | Notes |
|-----------|-------|-------|
| Material | Al alloy or 9% Ni steel | |
| Plate thickness | 15 - 25 mm | |
| Insulation | PUF blocks | |
| Cross-section | Octagonal with chamfers | |
| Centerline bulkhead | Liquid-tight, longitudinal | |
| Swash bulkhead | Mid-length, transverse | Splits into 4 compartments |
| Pump sump | Aft centerline | |

---

### 7. Automotive Fuel Tank

| Parameter | Range | Notes |
|-----------|-------|-------|
| Capacity | 40 - 80 L | Passenger car |
| Material | HDPE or Steel | HDPE: complex shapes |
| Baffle hole diameter | 11 - 14 mm | 7/16" - 9/16" |
| Holes per baffle | 1 - 2 | |
| Limber holes | Bottom | Fuel flow |
| Air holes | Top | Air evacuation |
| Support | Saddle/strap, 4 risers | Max 150L |

---

### 8. Aerospace Propellant Tank

#### PMD (Propellant Management Device)

| Component | Key Parameters |
|-----------|---------------|
| **Vanes** | Count: 2-8 (typical 4-6), angle: 45°-60° to wall |
| **Gallery Arms** | Count: 3-6, length: 0.5R-0.8R, mesh: 50-200 μm |
| **Anti-Slosh Baffle Rings** | Diameter: 0.7D-0.9D, spacing: 0.2H-0.4H |
| **Diaphragm** | Flexible membrane separating fuel/gas |
| **Bladder** | Balloon-like polymer fuel containment |

Design principle: Baffles placed slightly below free surface maximize damping.

---

## Part C: Common Parametric Study Variables

### Variables Commonly Swept in Research

| Variable Group | Parameter | Study Range |
|---------------|-----------|-------------|
| **Geometry** | L, W, H, D | ±20-50% |
| | H/L or H/D ratio | 0.3 - 1.5 |
| **Fill level** | h/H | 0.1 - 0.9 (step 0.1) |
| **Baffle** | Height ratio h_b/h | 0.3 - 0.9 |
| | Position x_b/L | 0.2 - 0.8 |
| | Porosity ε | 0 - 0.3 |
| **Excitation** | Frequency ω/ω_n | 0.5 - 2.0 |
| | Amplitude A/L | 0.01 - 0.1 |
| **Fluid** | Density ρ | 800-1000 kg/m³ |
| | Viscosity μ | 0.001 - 0.1 Pa·s |

### Response Variables (Outputs)

- Max free surface elevation
- Peak pressure on walls / roof
- Sloshing force & moment
- Energy dissipation rate
- Natural frequency shift (with baffles)

### Frequency Ratio Effects

| ω_exc / ω_n | Effect |
|-------------|--------|
| **1.0** | **Resonance** — max sloshing, roof impact, max loads |
| 2.0 | Destructive interference, reduced amplitude |
| ≠ 1, 2 | General sloshing |

---

## Part D: Parametric JSON Schema (for STL Generator)

```json
{
  "tank_type": "rectangular | cylindrical_upright | cylindrical_horizontal | lng_membrane | lng_moss | spb | automotive | aerospace",
  "dimensions": {
    "length": 1.0,
    "width": 0.6,
    "height": 0.8,
    "diameter": null,
    "wall_thickness": 0.005,
    "corner_radius": 0.01
  },
  "head": {
    "type": "flat | torispherical | ellipsoidal_2to1 | hemispherical | conical",
    "crown_radius_ratio": 1.0,
    "knuckle_radius_ratio": 0.06,
    "cone_angle_deg": 45
  },
  "baffles": [
    {
      "type": "vertical | horizontal | perforated",
      "position_ratio": 0.5,
      "height_ratio": 0.8,
      "thickness": 0.003,
      "porosity": 0.15,
      "hole_diameter": 0.02,
      "hole_spacing": 0.05
    }
  ],
  "fill_level": 0.5,
  "supports": {
    "type": "saddle | skirt | legs | none",
    "count": 2,
    "position_ratio": 0.25,
    "angle_deg": 120
  },
  "lng_specific": {
    "upper_chamfer_height": 8.63,
    "lower_chamfer_height": 3.77,
    "chamfer_angle_deg": 135,
    "centerline_bulkhead": true,
    "swash_bulkhead": true
  }
}
```

---

## Part E: Research Papers with Datasets

> Semantic Scholar API + SerpAPI (Google Scholar) 검색 결과. 2026-02-14 기준.

### Tier 1: Immediately Downloadable Datasets

| # | Resource | Tank Type | Data Type | Link |
|---|----------|-----------|-----------|------|
| 1 | **SPHERIC Benchmark Test 10** | Rectangular 2D | Pressure, motion, video (100 repeats) | [FTP](http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/) |
| 2 | **DualSPHysics 05_SloshingTank** | Rectangular 2D | XML case, motion files | [GitHub](https://github.com/DualSPHysics/DualSPHysics/wiki/7.-Testcases) |
| 3 | **beatrizmoya/sloshingfluids** | Various | Abaqus FEM state variables | [GitHub](https://github.com/beatrizmoya/sloshingfluids) |
| 4 | **CCP-WSI Data Repository** (Blind Test 5) | Circular | Experimental benchmark | [CCP-WSI](https://ccp-wsi.ac.uk/data_repository/) |
| 5 | **xevious54/Sloshing-Tank** | Rectangular | OpenFOAM VOF case files | [GitHub](https://github.com/xevious54/Sloshing-Tank) |
| 6 | **LGST-LAB/slosh_ml** | Propellant | Contour files, MATLAB GUI | [GitHub](https://github.com/LGST-LAB/slosh_ml) |
| 7 | **SPHERIC Validation Tests** (9,10,14,15,19) | Various | SPH benchmarks full set | [SPHERIC](https://www.spheric-sph.org/validation-tests) |

### Tier 2: Open Access Papers (Geometry Reproducible)

#### Rectangular Tank

| Paper | Year | DOI / Link | Key Data | Cited |
|-------|------|-----------|----------|-------|
| Souto-Iglesias et al. "Canonical problems in sloshing" | 2011-2015 | [10.1016/j.oceaneng.2015.05.013](https://doi.org/10.1016/j.oceaneng.2015.05.013) | 900x508mm tank, pressure series, SPHERIC FTP | — |
| Chen & Xue "Liquid Sloshing with Different Filling Levels Using OpenFOAM" | 2018 | [10.3390/W10121752](https://www.mdpi.com/2073-4441/10/12/1752/pdf) | 6 experimental wave height datasets, fill 20-70% | 84 |
| Liu et al. "Sloshing under Pitch Excitations" | 2024 | [10.3390/w16111551](https://www.mdpi.com/2073-4441/16/11/1551/pdf) | Pressure, free surface, power spectrum, fill 20/50/70% | — |
| Luo et al. "Stratified Sloshing" | 2021 | [10.1155/2021/6639223](https://doi.org/10.1155/2021/6639223) | Two-layer fluid, OpenFOAM, velocity field | — |
| JMSE "Perforated and Imperforate Baffles" | 2022 | [10.3390/jmse10101335](https://www.mdpi.com/2077-1312/10/10/1335/pdf) | Sloshing pressure per baffle type | — |
| Frosina et al. "Sloshing in a Fuel Tank" | 2018 | [MDPI Energies 11/3/682](https://www.mdpi.com/1996-1073/11/3/682) | Real automotive CAD geometry, CFD mesh | 31 |

#### Cylindrical Tank

| Paper | Year | DOI / Link | Key Data | Cited |
|-------|------|-----------|----------|-------|
| CCP-WSI Blind Test 5 "Circular Tank" | 2025 | [OnePetro](https://onepetro.org/ISOPEIOPEC/proceedings-abstract/ISOPE25/ISOPE25/ISOPE-I-25-332/713515) | CCP-WSI Data Repository | — |
| Ma et al. "Multi-Phase Sloshing Using DualSPHysics" | 2025 | [ISOPE25](https://onepetro.org/ISOPEIOPEC/proceedings-abstract/ISOPE25/ISOPE25/713553) | SolidWorks CAD + DualSPHysics + ParaView | — |
| "Sloshing in Real-Scale Water Storage Tank" | 2020 | [10.3390/w12082098](https://www.mdpi.com/2073-4441/12/8/2098/pdf) | Seismic loading, large cylindrical | — |
| "Liquid Sloshing in 3D Flexible Cylindrical Tank" | 2024 | [10.1063/5.0235933](https://doi.org/10.1063/5.0235933) | OpenFOAM+FEniCS+preCICE, seismic data | — |
| "Liquid Hydrogen Horizontal Tank" | 2023 | [10.3390/en16041851](https://www.mdpi.com/1996-1073/16/4/1851/pdf) | LH2 horizontal cylindrical, 3D model | — |

#### LNG / Prismatic Tank

| Paper | Year | DOI / Link | Key Data | Cited |
|-------|------|-----------|----------|-------|
| Jiao et al. "LNG ship motions + sloshing by DualSPHysics" | 2024 | [10.1016/j.oceaneng.2024.119148](https://doi.org/10.1016/j.oceaneng.2024.119148) | DualSPHysics XML, experimental validation | 12 |
| Jiao et al. "Two side-by-side LNG ships" | 2024 | [10.1016/j.oceaneng.2024.117022](https://doi.org/10.1016/j.oceaneng.2024.117022) | DualSPHysics + experiment | 30 |
| Zhao et al. "Elastic LNG Tank with Baffles (SPHinXsys)" | 2024 | [arXiv:2409.16226](https://arxiv.org/abs/2409.16226) | Ring/vertical baffles, FSI | — |
| ISOPE "Sloshing Model Test Benchmark 1st/2nd" | 2012/2013 | [GTT PDF](https://gtt.fr/sites/default/files/2012-isope-results-the-first-sloshing-model-test-benchmark-2012-05-23.pdf) | Mark III type, multi-institution pressure data | 63/39 |
| Grotle & AEsoy "Sloshing + Thermodynamic Response" | 2017 | [MDPI Energies 10/9/1338](https://www.mdpi.com/1996-1073/10/9/1338) | **CAD-to-STL workflow explicit** (Siemens NX) | 46 |

#### Baffles / Special

| Paper | Year | DOI / Link | Key Data | Cited |
|-------|------|-----------|----------|-------|
| Ye et al. "Sloshing Suppression DRL + SPH" | 2025 | [arXiv:2505.02354](https://arxiv.org/abs/2505.02354) | Elastic baffle + Deep RL control | — |
| Wang et al. "Anti-slosh Baffle in Tank Trucks" | 2024 | [10.1088/1361-6501/ad6629](https://doi.org/10.1088/1361-6501/ad6629) | "All data included within article" | — |
| NASA "Anti-Slosh Baffle Pressure Load Model" | 2023 | [NASA NTRS](https://ntrs.nasa.gov/citations/20230005181) | Cylindrical ring baffles, analytical + CFD | — |
| "Stepped-Base Rectangular Tank" | 2021 | [10.1063/5.0044682](https://doi.org/10.1063/5.0044682) / [UPC PDF](https://upcommons.upc.edu/bitstream/2117/351929/1/31985279.pdf) | Experimental + numerical, non-standard geometry | — |
| "Sloshing Tank with Flexible Floating Body" (Zenodo) | 2025 | [ScienceDirect](https://www.sciencedirect.com/science/article/pii/S235234092500839X) | Velocity vectors, wall pressure → Zenodo raw data | — |

### Tier 3: DualSPHysics Core Papers (Solver Reference)

| Paper | Year | DOI | Notes |
|-------|------|-----|-------|
| Buruchenko & Crespo "Validation for Liquid Sloshing" | 2014 | [10.1115/ICONE22-30968](https://doi.org/10.1115/ICONE22-30968) | DualSPHysics sloshing validation origin |
| English et al. "mDBC for Tank Sloshing" | 2021 | [10.1007/s40571-021-00403-3](https://link.springer.com/content/pdf/10.1007/s40571-021-00403-3.pdf) | Modified Dynamic BC, open access |
| Dominguez et al. "DualSPHysics: Multiphysics" | 2021 | [10.1007/s40571-021-00404-2](https://link.springer.com/content/pdf/10.1007/s40571-021-00404-2.pdf) | Full capability overview |
| Zhan et al. "DualSPHysics+ Enhanced Accuracy" | 2024 | [10.1016/j.cpc.2024.109389](https://doi.org/10.1016/j.cpc.2024.109389) | Accuracy/energy conservation improvements |

---

## References

### Standards
- ASME BPVC Section VIII Division 1 — Pressure vessel design
- API 650 — Welded tanks for oil storage
- API 620 — Large low-pressure storage tanks
- IGC Code — Ships carrying liquefied gases
- DNV CN 31.12 — LNG Type B prismatic tank

### CAD/STL Sources
- [GrabCAD Library](https://grabcad.com/library)
- [CGTrader Free Models](https://www.cgtrader.com/free-3d-models)
- [NASA 3D Resources](https://www.nasa.gov/3d-resources/)
- [TraceParts](https://www.traceparts.com/)
- [Free3D](https://free3d.com/)

### Open Datasets & Repositories
- [SPHERIC Benchmark FTP](http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/)
- [SPHERIC Validation Tests](https://www.spheric-sph.org/validation-tests)
- [DualSPHysics Test Cases](https://github.com/DualSPHysics/DualSPHysics/wiki/7.-Testcases)
- [CCP-WSI Data Repository](https://ccp-wsi.ac.uk/data_repository/)
- [beatrizmoya/sloshingfluids](https://github.com/beatrizmoya/sloshingfluids) — Abaqus FEM dataset
- [xevious54/Sloshing-Tank](https://github.com/xevious54/Sloshing-Tank) — OpenFOAM VOF
- [LGST-LAB/slosh_ml](https://github.com/LGST-LAB/slosh_ml) — Propellant sloshing MATLAB

### Key Academic Papers
- Souto-Iglesias et al. (2011-2015) — SPHERIC canonical sloshing benchmark
- Chen & Xue (2018) — OpenFOAM sloshing with different fill levels
- English et al. (2021) — mDBC for tank sloshing in DualSPHysics
- Jiao et al. (2024) — LNG ship motions + DualSPHysics sloshing
- Grotle & AEsoy (2017) — CAD-to-STL sloshing workflow
- ISOPE Sloshing Model Test Benchmark (2012/2013) — Mark III multi-institution
- Buruchenko & Crespo (2014) — DualSPHysics sloshing validation
- NASA NTRS (2023) — Anti-slosh baffle pressure load model
- Akyildiz & Unal (2006) — Rectangular tank sloshing benchmark
- L.P. Zick (1951) — Horizontal vessel saddle design
