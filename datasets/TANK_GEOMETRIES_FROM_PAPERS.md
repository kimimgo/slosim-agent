# Tank Geometries & Full Parametric Specification — Extracted from Research Papers

논문에서 추출한 **실제 실험/시뮬레이션 탱크의 전체 파라미터**. 파라메트릭 STL 생성기의 프리셋으로 사용.
Part B (docs/STL_DATASETS_AND_PARAMETRICS.md)의 심층 파라메트릭 연구와 결합.

> **JSON 프리셋 스키마 = Part D 전체 스키마 + 논문 고유 실험 조건**

---

## 1. SPHERIC Benchmark Test 10 (2011) — Impact Pressure

**Source**: Souto-Iglesias, Botia-Vera et al. (UPM Madrid)
**Status**: SPH 솔버 검증 골드 스탠다드 (DualSPHysics, SPHinXsys, GPUSPH 등 모두 사용)
**DOI**: 10.1016/j.oceaneng.2015.05.013

### 탱크 형상 (Rectangular 2D)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Length | `L` | 900 | mm | 여진 방향 |
| Breadth | `B` | 31 / 62 / 124 | mm | 3가지 변형 (0.5x, 1x, 2x) |
| Height | `H` | 508 | mm | |
| Material | — | Plexiglas (PMMA) | — | |
| Wall thickness | `t_wall` | ~5 (추정) | mm | 논문 명시 없음 |
| Corner radius | `r_corner` | 0 | mm | Sharp corner |
| Aspect ratio L/H | — | 1.772 | — | |
| Depth ratio H/L | — | 0.564 | — | Intermediate depth |

### 유체 물성 (3종)

| Fluid | Density ρ (kg/m³) | Dynamic Viscosity μ (Pa·s) | Surface Tension σ (N/m) | Kinematic Viscosity ν (m²/s) |
|-------|-------------------|---------------------------|------------------------|------------------------------|
| **Water** | 998 | 8.94e-4 | 0.0728 | 8.96e-7 |
| **Sunflower Oil** | 990 | 0.045 | 0.033 | 4.55e-5 |
| **Glycerin** | 1261 | 0.934 | 0.064 | 7.41e-4 |

Temperature: 19°C, Atmospheric pressure

### 충전율 & 자연주파수

| Fill Condition | Fill Height h (mm) | Fill Ratio h/H | Sloshing Mode | Natural Period T₁ (s) | Natural Freq ω₁ (rad/s) |
|----------------|-------------------|----------------|---------------|----------------------|--------------------------|
| **Low (shallow)** | 93 | 0.183 | Lateral wall impact (Sensor 1) | 1.919 | 3.274 |
| **High (deep)** | 355.6 | 0.700 | Roof/ceiling impact (Sensor 3) | 1.167 | 5.383 |

Natural frequency formula (shallow water dispersion relation):
```
ω₁ = √(g · k₁ · tanh(k₁ · h))
k₁ = π / L = π / 0.9 = 3.491 rad/m
```

### 여진 조건 (Excitation)

| Test Case | Fluid | Fill | T/T₁ | Period T (s) | Freq f (Hz) | Type | Center |
|-----------|-------|------|-------|-------------|-------------|------|--------|
| Water/Low | Water | 18% | 0.85 | 1.630 | 0.613 | Rotation (pitch) | Bottom center [0.45, 0, 0] |
| Oil/Low | Oil | 18% | 0.80 | 1.535 | 0.651 | Rotation (pitch) | Bottom center |
| Water/High | Water | 70% | 1.00 | 1.168 | 0.856 | Rotation (pitch, **resonance**) | Bottom center |
| Oil/High | Oil | 70% | 1.00 | 1.168 | 0.856 | Rotation (pitch, **resonance**) | Bottom center |

Rotation axis: Center of bottom line (x=L/2, y=0, z=0)

### 센서 배치

| Sensor | Position | Target Measurement |
|--------|----------|-------------------|
| Sensor 1 | Lateral wall (low position) | Low-fill lateral impact pressure |
| Sensor 2 | Lateral wall (mid) | Intermediate pressure |
| Sensor 3 | Ceiling (top wall) | High-fill roof impact pressure |

### 실험 통계

- **100 independent repeats** per condition (통계적 유의성 확보)
- 3 min rest between repeats (잔류 진동 소멸)
- Data channels: pressure [mbar], angle [deg], velocity [deg/s], acceleration [deg/s²]
- Acquisition frequency: unspecified (ISOPE benchmark에서 15-50 kHz 범위)

### 바플 (Baffle) — 없음

원본 SPHERIC Test 10은 baffle 없음 (bare tank). Baffle 추가는 후속 연구에서.

---

## 2. Chen & Xue (2018) — OpenFOAM Fill Level Parametric Study

**Source**: Chen, Xue — Hohai University, China
**DOI**: 10.3390/w10121752
**Status**: OpenFOAM InterDyMFoam 검증, 6개 fill level 파라메트릭

### 탱크 형상 (Rectangular 3D)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Length | `L` | 600 | mm | 여진 방향 |
| Width | `W` | 300 | mm | |
| Height | `H` | 650 | mm | |
| Material | — | Plexiglas | — | 8mm 두께 |
| Wall thickness | `t_wall` | 8 | mm | 명시됨 |
| Aspect ratio L/W | — | 2.0 | — | |
| Depth ratio H/L | — | 1.083 | — | Deep tank |

### 유체 물성

| Fluid | Density ρ (kg/m³) | Dynamic Viscosity μ (Pa·s) |
|-------|-------------------|---------------------------|
| **Water** | 998 | 1.003e-3 |
| **Air** | 1.225 | 1.789e-5 |

### 충전율 & 자연주파수

| Fill Level h (mm) | Fill Ratio h/H | h/L Ratio | Natural Freq ω₁ (rad/s) | Category |
|-------------------|----------------|-----------|--------------------------|----------|
| **83 (13.8%)** | 0.128 | 0.138 | **4.749** | Shallow |
| **185 (30.8%)** | 0.285 | 0.308 | **6.333** | Near-critical |
| 90 | 0.138 | 0.150 | — | Shallow (3D probe test) |
| 200 | 0.308 | 0.333 | — | Near-critical (3D probe test) |

Natural frequency formula:
```
ωₙ = √(g · (nπ/L) · tanh(nπh/L)),  n = 1,2,3...
```

**Key finding**: h/L ≈ 0.24가 soft-spring → hard-spring 전이 임계값

### 여진 조건

| Parameter | Value | Notes |
|-----------|-------|-------|
| Motion type | Horizontal sinusoidal `x = -A sin(ωt)` | Sway (lateral) |
| Amplitude `A` | 7 mm | Fixed for most tests |
| Frequency range | Wide sweep around ω₁ | 0.5ω₁ ~ 2.0ω₁ |
| Platform | 6-DOF hexapod | Hohai University |
| Motion accuracy | < 0.5 mm displacement error | |

### 센서 배치

- **5 pressure sensors** embedded in left wall (Figure 4 참조)
- Vertical array: varying heights from bottom to top
- **Camera** (front-mounted): free surface profile recording
- **3D probes** (Figure 11): transverse direction pressure distribution at same height, different Y positions

### 바플 — 없음

Bare tank study. 원본 연구는 fill level 영향에 집중.

### CFD 설정 (OpenFOAM)

| Parameter | Value |
|-----------|-------|
| Solver | InterDyMFoam (VOF, two-phase NS) |
| Turbulence | None (laminar) |
| Time discretization | 1st order implicit Euler |
| Gradient scheme | Gauss linear |
| Divergence scheme | van Leer |
| Pressure-velocity coupling | PIMPLE (PISO+SIMPLE) |
| Surface tension | Disabled (C=0) |

---

## 3. Liu et al. (2024) — Pitch Excitation Parametric Study

**Source**: Liu, Li, Peng, Zhou, Gao — Jiangsu University of Science and Technology + NGI
**DOI**: 10.3390/w16111551
**Status**: 대형 탱크 피치 여진 실험, 3개 충전율 × 5개 주파수 × 3개 진폭

### 탱크 형상 (Rectangular 3D — Large Scale)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Length | `L` | 1000 | mm | 여진 방향 |
| Width | `B` | 500 | mm | |
| Height | `H` | 1000 | mm | |
| Material | — | Plexiglas | — | 20mm 두께 |
| Wall thickness | `t_wall` | 20 | mm | 명시됨 |
| Aspect ratio L/B | — | 2.0 | — | |
| Depth ratio H/L | — | 1.0 | — | Deep/square |

### 유체 물성

| Fluid | Notes |
|-------|-------|
| **Water** + Rhodamine B | Red coloring (ρ, μ 변화 없음), free surface 시각화용 |

### 충전율 & 자연주파수

| Fill Condition | Fill Ratio h/H | Theoretical f₁ (Hz) | Measured f₁ (Hz) | Error Δ (%) | Category |
|----------------|----------------|---------------------|-------------------|-------------|----------|
| **Shallow** | 0.20 (20%H = 200mm) | 0.659 | 0.657 | 0.3% | Shallow |
| **Near-critical** | 0.30 (30%H = 300mm) | 0.758 | 0.754 | 0.5% | Near-critical |
| **High** | 0.70 (70%H = 700mm) | 0.872 | 0.869 | 0.3% | Deep |

Natural frequency formula (same as standard rectangular):
```
fₙ = (1/2π) · √(g · (nπ/L) · tanh(nπh/L))
```

### 여진 조건 (Test Matrix)

**Equal-amplitude pitch:**

| Fill Rate | Amplitude A (°) | Frequencies f (Hz) | Count |
|-----------|-----------------|-------------------|-------|
| 20%H | 2 | 0.53, 0.59, **0.66 (f₁)**, 0.73, 0.79 | 5 |
| 30%H | 2 | 0.61, 0.68, **0.76 (f₁)**, 0.84, 0.91 | 5 |
| 70%H | 2 | 0.70, 0.78, **0.87 (f₁)**, 0.96, 1.04 | 5 |

**Variable-amplitude pitch:**

| Fill Rate | Amplitudes A (°) | Frequency f (Hz) |
|-----------|-----------------|------------------|
| 20%H | 1, 2, 3 | 0.66 (resonance) |
| 30%H | 1, 2, 3 | 0.76 (resonance) |
| 70%H | 1, 2, 3 | 0.87 (resonance) |

Frequency range: **0.8f₁ ~ 1.2f₁** (resonance neighborhood sweep)
Sampling time: **400 s** per test case
Pitch motion: `θ(t) = A·sin(2πft)`

### 실험 장비 (Experimental Apparatus)

| Equipment | Specification |
|-----------|---------------|
| Motion platform | 3-DOF shaking table |
| Roll range | ±30° |
| Pitch range | ±15° |
| Heave range | ±0.35 m |
| Rotational accuracy | 0.05° |
| Translational accuracy | 0.05 mm |
| Platform size | 2 m × 4 m |
| Max bearing capacity | 7 tons |
| Pressure transducer | CY302 miniature intelligent, 0-200 kPa, 1 kHz, 0.1 ms response, 0.1% FS |
| High-speed camera | FASTCAM SA-Z, 1024×1024, max 200,000 fps, min exposure 159 ns |

### 센서 배치

- Pressure monitoring points: right wall + top wall (ceiling)
- Multiple heights: various liquid-carrying heights
- **P1 ~ P_n**: wall-mounted pressure transducers
- Camera: front-facing, free surface tracking

### 바플 — 없음

Bare tank study. 파형 분류 (4 distinct waveforms), 3D vortex wave 발견.

### Key Findings (파라메트릭 영향)

- **20%H resonance**: Violent sloshing, hydraulic jump
- **70%H resonance**: Frequency shift (soft-spring nonlinearity), 3D vortex waves
- **30%H**: Near-critical depth, transitional behavior
- Power spectrum: dominant = excitation frequency + harmonics (2f, 3f, ...)

---

## 4. ISOPE Sloshing Benchmark (2012) — GTT Mark III Scale Model

**Source**: Loysel, Chollet, Gervaise, Brosset (GTT), ISOPE 2012
**9개 기관 참여**: ECM, ECN-BV, GTT x3, Marintek, PNU, UDE, UPM, UR
**Status**: 가장 대규모 슬로싱 벤치마크. LNG 멤브레인 탱크 단면 축소 모형.

### 탱크 형상 (Rectangular 2D — LNG Cross-Section)

| Parameter | Symbol | Value (Main) | Value (UR Scale) | Unit | Notes |
|-----------|--------|-------------|------------------|------|-------|
| Length | `L` | 946 | 700 | mm | 여진 방향 |
| Breadth | `W` | 118 | 87 | mm | |
| Height | `H` | 670 | 496 | mm | |
| Material | — | Plexiglas (PMMA) | — | — | |
| Tolerance | — | < 1 mm | — | mm | 정밀 가공 |
| Scale ratio | — | 1:1 (원본) | 1:1.35 (Froude) | — | |

### 유체 물성

| Fluid | Notes |
|-------|-------|
| **Water + Air** | Standard sloshing fluid |

### 충전율

- **High fill levels** (ceiling impact 조건)
- 14 test conditions at high fill levels
- 7 conditions analyzed in detail

### 여진 조건

| Parameter | Value | Notes |
|-----------|-------|-------|
| Motion type | Combined Tx + Ry + Tz | Hexapod 3-DOF |
| Motion system | Hexapod (6-DOF) | |
| Conditions | 14 different combinations | High fill scenarios |
| Acquisition freq | 15-50 kHz | Institution-dependent |

### 센서 배치

| System | Sensors | Type | Location |
|--------|---------|------|----------|
| Full | 72 pressure sensors | PCB / Kulite / Kistler | Ceiling (roof) |
| Reduced | 16 sensors | Same types | Key impact locations |

### LNG 멤브레인 탱크 파라메트릭 컨텍스트 (Part B 연동)

이 벤치마크는 실물 LNG 멤브레인 탱크 (Mark III)의 **2D 단면 축소 모형**이다.
실물 파라미터:

| Parameter | Value | Notes |
|-----------|-------|-------|
| Total capacity | 120,000 - 270,000 m³ | Full LNGC |
| Cross-section | Octagonal (8-sided) | |
| Upper chamfer angle | 135° | |
| Lower chamfer angle | 135° | |
| Upper chamfer height | ~8.63 m (typical) | |
| Lower chamfer height | ~3.77 m (typical) | |
| Tank count | 4-6 per ship | Centerline arrangement |
| Membrane system | Mark III: SS 304L corrugated 1.2mm / NO96: Invar 0.7mm | |
| Primary insulation | PUF 160mm (Mark III) / Plywood+Perlite (NO96) | |
| Design standard | IGC Code | Full secondary barrier |
| Design life | 10⁸ wave encounters (20 yr, North Atlantic) | |

---

## 5. English et al. (2021) — mDBC Sloshing Validation

**Source**: English, Dominguez, Vacondio, Crespo et al.
**DOI**: 10.1007/s40571-021-00403-3
**Status**: DualSPHysics mDBC (Modified Dynamic Boundary Condition) 검증

### 탱크 형상 (SPHERIC Test 10 재현)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Length | `L` | 900 | mm | Same as SPHERIC |
| Height | `H` | 508 | mm | Same as SPHERIC |
| Width | `B` | 62 | mm | 1x variant |

### 충전 & 여진

| Parameter | Value |
|-----------|-------|
| Fill height h | 93 mm (18% fill, low) |
| Sensor | Sensor 1 (lateral wall) |
| Excitation | Same as SPHERIC Test 10 |

### SPH 수치 파라미터 (DualSPHysics 전용)

| Parameter | Value 1 | Value 2 | Notes |
|-----------|---------|---------|-------|
| Particle spacing `dp` | 0.004 m | 0.002 m | 해상도 테스트 |
| Total particles (dp=0.002) | 26,791 | — | |
| Physical simulation time | 7.0 s | — | |
| Boundary condition | mDBC (modified dynamic BC) | — | vs. standard DBC |
| Kernel | Wendland C2 (추정) | — | DualSPHysics default |
| Viscosity model | Artificial viscosity (추정) | — | |

### 추가 검증 탱크 (Still Water)

| Parameter | Value | Notes |
|-----------|-------|-------|
| Length | 2400 mm | Rectangular |
| Height | 1200 mm | |
| Water height | 500 mm | |
| Bottom obstacle | Trigonal wedge, h=240mm | Tank center bottom |
| Purpose | Hydrostatic pressure validation | |

---

## 6. NASA Anti-Slosh Baffle (2023) — Cylindrical Ring Baffle

**Source**: Yang, Sansone, Brodnick, Williams (NASA MSFC)
**NTRS**: 20230005181
**Status**: 원통형 탱크 링 바플 압력 모델. 우주 추진제 탱크 직접 적용.

### 탱크 형상 (Cylindrical Upright)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Diameter | `D` | 2.84 | m | = 284 cm (Scholl et al. 실험 탱크) |
| Radius | `R` | 1.42 | m | |
| Height | `H` | — | m | (논문 미명시, 충분한 높이) |
| Wall thickness | `t_wall` | — | mm | |

### 바플 파라메트릭 (Ring Baffle) — **핵심 파라미터**

| Parameter | Symbol | Tested Values | Unit | Notes |
|-----------|--------|--------------|------|-------|
| **Baffle depth ratio** | `d/R` | 0.00, 0.10, 0.25 | — | Surface부터 깊이 |
| Baffle width | `W` | — | m | |
| Baffle type | — | Rigid ring | — | Annular plate |
| Baffle count | `n` | 1 | — | Single ring |

### 자연주파수 (Cylindrical, neglecting bottom effect)

```
ω² = (g · λ₁ / R) · tanh(λ₁ · h / R)
λ₁ = 1.841  (first root of J'₁, Bessel derivative)
```

### 압력 모델 (Key Equations)

| Equation | Formula | Notes |
|----------|---------|-------|
| Max pressure | `p_max = √2 · ρ · g · η` | η = slosh wave height |
| Nondimensional pressure | `p̄ = p / (ρ·g·η)` | Max = √2 ≈ 1.414 |
| Phase shift (d/R=0.25) | `Δθ = 6.95°` | Below baffle vs above |
| Phase shift (d/R=0.10) | `Δθ = 8.76°` | |
| Phase shift (d/R=0.00) | `Δθ = 88.4°` | Near-surface, **max damping** |

### 바플 설계 가이드라인 (Part B 연동)

Cylindrical tank ring baffle 최적 설계 원칙:
- **표면 가까울수록 감쇄 효과 극대** (d/R → 0이 최적)
- d/R = 0 → Δθ ≈ 90° (거의 완전 위상 분리)
- 너무 표면 가까우면 구조 하중 과다 → **d/R = 0.05 ~ 0.10** 실용적 범위
- Ring width: 0.05R ~ 0.15R (typical)
- Porosity: 0% (solid ring, NASA 사용) ~ 20% (perforated, 해양 탱크)

### 우주 추진제 탱크 추가 파라미터 (Part B 연동)

| PMD Component | Key Parameters | Notes |
|---------------|---------------|-------|
| Vanes | Count: 2-8 (typical 4-6), angle: 45°-60° to wall | Surface tension 이용 |
| Gallery Arms | Count: 3-6, length: 0.5R-0.8R, mesh: 50-200 μm | |
| Anti-Slosh Baffle Rings | Diameter: 0.7D-0.9D, spacing: 0.2H-0.4H | This paper focus |
| Diaphragm | Flexible membrane separating fuel/gas | Alternative to baffles |
| Bladder | Balloon-like polymer fuel containment | |

---

## 7. Zhao et al. (2024) — Elastic LNG Tank with Baffles (SPHinXsys)

**Source**: Zhao, Wu, Yu, Haidn, Hu — TU Munich
**arXiv**: 2409.16226
**Status**: SPH FSI framework (SPHinXsys), multi-phase (water+air), elastic wall + baffle

### 탱크 형상 (Cylindrical Horizontal — Grotle & Æsøy Validation)

| Parameter | Symbol | Value | Unit | Notes |
|-----------|--------|-------|------|-------|
| Type | — | Horizontal cylindrical | — | Semi-elliptical heads (2:1 ratio) |
| Wall thickness | `t_wall` | 18 | mm | 0.018 m |
| Fill ratio | `hw/D` | 0.255 (validation) / 0.50 (LNG study) | — | |
| Sensor position | — | 0.122 m from left straight section start | — | Free surface height |

### Validation Excitation

| Parameter | Value | Notes |
|-----------|-------|-------|
| Motion type | Pitching (rotation) | |
| Amplitude `A` | 3° | |
| Frequency `f` | 0.5496 Hz | |
| Period `T` | 1.819 s | |
| Rotation axis | Tank bottom center plane intersection | |
| Motion function | `θ(t) = A·sin(2πft)` | |

### 탱크 재료 물성 (Steel — Elastic Wall)

| Property | Value | Unit | Notes |
|----------|-------|------|-------|
| Density | 7890 | kg/m³ | Steel |
| Poisson ratio | 0.27 | — | |
| Young's modulus | 135 | GPa | |

### 유체 물성 (Multi-phase)

| Phase | Density ρ (kg/m³) | Notes |
|-------|-------------------|-------|
| Water | 1000 | Heavy phase |
| Air | 1.226 | Light phase |

### 바플 설정 — **핵심 파라메트릭**

**Baffle Type 1: Ring Baffles**

| Parameter | Value | Notes |
|-----------|-------|-------|
| Count | 2 | Equally spaced |
| Type | Annular ring (ring baffle) | |
| Material | Same as tank wall (rigid) or elastic | |

**Baffle Type 2: Vertical Baffles**

| Parameter | Value | Notes |
|-----------|-------|-------|
| Count | 2 | Equally spaced |
| Type | Flat plate (vertical) | |
| Material | Same as tank wall (rigid) or elastic | |

**Elastic Baffle Study (Section 5):**

| Configuration | hw/D | Wall thickness | Baffle Type | Baffle E (kPa) | Baffle ρ (kg/m³) | Baffle ν |
|---------------|------|---------------|-------------|-----------------|-------------------|----------|
| Rigid baffle | 0.50 | 30 mm | Vertical, center | ∞ (rigid) | — | — |
| Elastic E=500 | 0.50 | 30 mm | Vertical, center | 500 | 2500 | 0.47 |
| Elastic E=50 | 0.50 | 30 mm | Vertical, center | 50 | 2500 | 0.47 |
| Elastic E=5 | 0.50 | 30 mm | Vertical, center | 5 | 2500 | 0.47 |

**Key Finding**: 바플이 rigid에 가까울수록 슬로싱 억제 효과 증가

### SPH 수치 파라미터 (SPHinXsys)

| Parameter | Value | Notes |
|-----------|-------|-------|
| Particle spacing Δx | 0.006 m (production) / 0.0045 m (refined) | |
| Method | WCSPH with multi-phase Riemann solver | |
| Fluid-Structure | Single-framework FSI (SPH for both) | No data transfer error |
| Damping | Artificial-viscosity-based, first 1.0 s | Preloading stabilization |
| Time stepping | Dual-criteria CFL | |
| EOS | `p = c²(ρ - ρ₀)`, c = 10·Vmax | Weakly compressible |
| Vmax estimation | `Vmax = √(2g·hw)` | |

### Cylindrical Horizontal Tank 추가 파라미터 (Part B 연동)

| Parameter | Symbol | Range | Notes |
|-----------|--------|-------|-------|
| L/D ratio | — | 2 - 6 | Optimal: 3-4 |
| **Head type** | — | Semi-elliptical 2:1 (this paper) | Also: flat, torispherical, hemispherical |
| **Saddle supports** (if applicable) | | | |
| Position ratio A/R | — | 0.25 | Min stress (Zick method) |
| Saddle angle θ | — | 120° - 150° | Usually 120° |
| Saddle width | — | 0.1L - 0.2L | |
| **Swash plates** | — | at L/3, 2L/3 | Opening area 30-50% |
| **Transverse baffles** | — | spacing 0.3L-0.5L | + limber holes + air holes |

---

## 8. Frosina et al. (2018) — Automotive Fuel Tank (FCA)

**Source**: Frosina, Senatore, Andreozzi, Fortunato, Giliberti
**Fiat Chrysler Automobiles + University of Naples**
**DOI**: 10.3390/en11030682
**Status**: 실제 자동차 연료탱크 CAD → CFD 워크플로우 시연

### 탱크 형상 (Automotive — Complex 3D)

| Parameter | Value | Notes |
|-----------|-------|-------|
| Type | Real automotive fuel tank | Complex 3D geometry from CAD |
| Material | Plexiglas faces (시각화용) | |
| Exact dimensions | **FCA 기밀 — 비공개** | |
| Capacity | 40-80 L (추정, 일반 승용차 범위) | |
| Shape | Irregular, multi-chamber | |
| CFD source | CAD geometry → CFD mesh (3D VOF) | |
| CFD accuracy | < 3% discrepancy (free surface + CoG) | |

### 유체

| Fluid | Notes |
|-------|-------|
| Water | Tinted with dark blue food colorant (시각화용) |

### 충전율

- 3 different fill levels tested (specific values not reported)

### 실험 장비 (Hexapod)

| Equipment | Specification |
|-----------|---------------|
| Platform | Moog Inc. Hexapod 8-DOF |
| Actuators | 6 linear + 2 tilt |
| Total stroke | 600 mm |
| Max excursion | 400 mm (all directions) |
| Peak acceleration | ~1g (10 m/s²) |
| Pitch/Roll range | > 50° |
| Min frequency | ~0.7 Hz |
| Camera | 1280 × 720, 25 fps |

### 자동차 연료탱크 파라메트릭 (Part B 연동)

| Parameter | Range | Notes |
|-----------|-------|-------|
| Capacity | 40 - 80 L | Passenger car |
| Material | HDPE or Steel | HDPE: 복잡 형상 가능 |
| **Baffle hole diameter** | 11 - 14 mm (7/16" - 9/16") | |
| **Holes per baffle** | 1 - 2 | |
| **Limber holes** | Bottom | 연료 유동 |
| **Air holes** | Top | 공기 배출 |
| **Support** | Saddle/strap, 4 risers | Max 150L |
| Wall thickness | 3-6 mm (HDPE), 0.8-1.5 mm (steel) | |

---

## Full Parametric JSON Presets (STL Generator용)

> 아래 프리셋은 Part B (8개 탱크타입 전체 파라미터) + Part D (JSON 스키마)를 논문 실험 조건과 결합한 **완전 사양**.

### preset_spheric_low — SPHERIC Test 10 Low Fill (Lateral Impact)

```json
{
  "name": "SPHERIC Test 10 - Low Fill (Lateral Impact)",
  "source": {
    "paper": "Souto-Iglesias & Botia-Vera (2011), SPHERIC Benchmark Test 10",
    "doi": "10.1016/j.oceaneng.2015.05.013",
    "dataset": "http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/",
    "repeats": 100
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.900,
    "width": 0.062,
    "height": 0.508,
    "wall_thickness": 0.005,
    "corner_radius": 0.0
  },
  "head": {
    "type": "flat",
    "note": "Sharp rectangular corners, no head curvature"
  },
  "baffles": [],
  "supports": {
    "type": "none"
  },
  "fill_level": 0.183,
  "fill_height_m": 0.093,
  "fluid": {
    "name": "water",
    "density": 998,
    "dynamic_viscosity": 8.94e-4,
    "kinematic_viscosity": 8.96e-7,
    "surface_tension": 0.0728,
    "temperature_C": 19
  },
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "center": [0.45, 0.0, 0.0],
    "axis": "z",
    "period_s": 1.630,
    "frequency_hz": 0.613,
    "frequency_ratio_T_T1": 0.85,
    "amplitude_deg": null,
    "motion_function": "theta(t) = A * sin(2*pi*t/T)"
  },
  "natural_frequency": {
    "mode": 1,
    "period_s": 1.919,
    "frequency_hz": 0.521,
    "omega_rad_s": 3.274,
    "formula": "omega = sqrt(g * (pi/L) * tanh(pi*h/L))",
    "k1_rad_m": 3.491,
    "regime": "shallow"
  },
  "sensors": [
    {"id": "S1", "type": "pressure", "location": "lateral_wall_low", "target": "impact_pressure"}
  ],
  "sph_reference": {
    "solver": "DualSPHysics",
    "dp_options_m": [0.004, 0.002],
    "simulation_time_s": 7.0,
    "boundary_condition": "mDBC",
    "total_particles_dp002": 26791
  },
  "alternative_fluids": [
    {"name": "sunflower_oil", "density": 990, "dynamic_viscosity": 0.045, "surface_tension": 0.033},
    {"name": "glycerin", "density": 1261, "dynamic_viscosity": 0.934, "surface_tension": 0.064}
  ],
  "breadth_variants_m": [0.031, 0.062, 0.124]
}
```

### preset_spheric_high — SPHERIC Test 10 High Fill (Roof Impact)

```json
{
  "name": "SPHERIC Test 10 - High Fill (Roof Impact, Resonance)",
  "source": {
    "paper": "Souto-Iglesias & Botia-Vera (2011), SPHERIC Benchmark Test 10",
    "doi": "10.1016/j.oceaneng.2015.05.013",
    "dataset": "http://canal.etsin.upm.es/ftp/SPHERIC_BENCHMARKS/"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.900,
    "width": 0.062,
    "height": 0.508,
    "wall_thickness": 0.005,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.700,
  "fill_height_m": 0.3556,
  "fluid": {
    "name": "water",
    "density": 998,
    "dynamic_viscosity": 8.94e-4,
    "kinematic_viscosity": 8.96e-7,
    "surface_tension": 0.0728,
    "temperature_C": 19
  },
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "center": [0.45, 0.0, 0.0],
    "axis": "z",
    "period_s": 1.168,
    "frequency_hz": 0.856,
    "frequency_ratio_T_T1": 1.00,
    "note": "RESONANCE condition — maximum sloshing severity"
  },
  "natural_frequency": {
    "mode": 1,
    "period_s": 1.167,
    "frequency_hz": 0.857,
    "omega_rad_s": 5.383,
    "formula": "omega = sqrt(g * (pi/L) * tanh(pi*h/L))",
    "k1_rad_m": 3.491,
    "regime": "intermediate_to_deep"
  },
  "sensors": [
    {"id": "S3", "type": "pressure", "location": "ceiling_top", "target": "roof_impact_pressure"}
  ]
}
```

### preset_chen2018_shallow — Chen & Xue OpenFOAM Low Fill

```json
{
  "name": "Chen & Xue (2018) - OpenFOAM Shallow Fill (13.8%)",
  "source": {
    "paper": "Chen & Xue (2018), Hohai University",
    "doi": "10.3390/w10121752",
    "solver": "OpenFOAM InterDyMFoam"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.600,
    "width": 0.300,
    "height": 0.650,
    "wall_thickness": 0.008,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.138,
  "fill_height_m": 0.083,
  "fluid": {
    "name": "water",
    "density": 998,
    "dynamic_viscosity": 1.003e-3,
    "temperature_C": 20
  },
  "excitation": {
    "type": "translation",
    "subtype": "sway_horizontal",
    "amplitude_m": 0.007,
    "frequency_sweep": true,
    "frequency_range": "0.5*omega1 to 2.0*omega1",
    "motion_function": "x(t) = -A * sin(omega*t)"
  },
  "natural_frequency": {
    "mode": 1,
    "omega_rad_s": 4.749,
    "frequency_hz": 0.756,
    "formula": "omega = sqrt(g * (n*pi/L) * tanh(n*pi*h/L))",
    "regime": "shallow",
    "h_over_L": 0.138
  },
  "sensors": [
    {"id": "P1-P5", "type": "pressure", "location": "left_wall_vertical_array", "count": 5}
  ],
  "cfd_settings": {
    "solver": "InterDyMFoam",
    "method": "VOF",
    "turbulence": "laminar",
    "time_scheme": "implicit_euler_1st",
    "pressure_velocity": "PIMPLE",
    "surface_tension": false
  },
  "parametric_study": {
    "variable": "fill_level",
    "values_percent": [13.8, 30.8],
    "critical_h_over_L": 0.24,
    "transition": "soft-spring → hard-spring at h/L ≈ 0.24"
  }
}
```

### preset_chen2018_nearcritical — Chen & Xue OpenFOAM Near-Critical Fill

```json
{
  "name": "Chen & Xue (2018) - OpenFOAM Near-Critical Fill (30.8%)",
  "source": {
    "paper": "Chen & Xue (2018), Hohai University",
    "doi": "10.3390/w10121752"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.600,
    "width": 0.300,
    "height": 0.650,
    "wall_thickness": 0.008,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.308,
  "fill_height_m": 0.185,
  "fluid": {
    "name": "water",
    "density": 998,
    "dynamic_viscosity": 1.003e-3
  },
  "excitation": {
    "type": "translation",
    "subtype": "sway_horizontal",
    "amplitude_m": 0.007,
    "frequency_sweep": true,
    "frequency_range": "0.5*omega1 to 2.0*omega1"
  },
  "natural_frequency": {
    "mode": 1,
    "omega_rad_s": 6.333,
    "frequency_hz": 1.008,
    "regime": "near_critical",
    "h_over_L": 0.308
  }
}
```

### preset_liu2024_shallow_resonance — Liu et al. Large Tank Shallow Resonance

```json
{
  "name": "Liu et al. (2024) - Large Tank 20% Fill Resonance Pitch",
  "source": {
    "paper": "Liu, Li, Peng, Zhou, Gao (2024), JUST + NGI",
    "doi": "10.3390/w16111551"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 1.000,
    "width": 0.500,
    "height": 1.000,
    "wall_thickness": 0.020,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.20,
  "fill_height_m": 0.200,
  "fluid": {
    "name": "water_rhodamineB",
    "density": 998,
    "dynamic_viscosity": 1.003e-3,
    "note": "Rhodamine B dye, no property change"
  },
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "amplitude_deg": 2.0,
    "frequency_hz": 0.66,
    "frequency_ratio_f_f1": 1.0,
    "motion_function": "theta(t) = A * sin(2*pi*f*t)",
    "note": "RESONANCE — violent sloshing, hydraulic jump observed"
  },
  "natural_frequency": {
    "mode": 1,
    "theoretical_hz": 0.659,
    "measured_hz": 0.657,
    "error_percent": 0.3,
    "regime": "shallow"
  },
  "sensors": [
    {"type": "pressure", "location": "right_wall + ceiling", "transducer": "CY302, 0-200kPa, 1kHz"},
    {"type": "camera", "model": "FASTCAM SA-Z", "resolution": "1024x1024", "max_fps": 200000}
  ],
  "test_matrix": {
    "equal_amplitude": {
      "fill_rates": ["20%H", "30%H", "70%H"],
      "amplitude_deg": 2,
      "frequency_ratios_f_f1": [0.8, 0.9, 1.0, 1.1, 1.2]
    },
    "variable_amplitude": {
      "fill_rates": ["20%H", "30%H", "70%H"],
      "amplitudes_deg": [1, 2, 3],
      "frequency": "f1 (resonance)"
    },
    "sampling_time_s": 400
  },
  "apparatus": {
    "platform": "3-DOF shaking table",
    "roll_range_deg": 30,
    "pitch_range_deg": 15,
    "heave_range_m": 0.35,
    "rotation_accuracy_deg": 0.05,
    "translation_accuracy_mm": 0.05,
    "platform_size_m": [2.0, 4.0],
    "max_load_tons": 7
  }
}
```

### preset_liu2024_high_resonance — Liu et al. Large Tank High Fill Resonance

```json
{
  "name": "Liu et al. (2024) - Large Tank 70% Fill Resonance Pitch",
  "source": {
    "paper": "Liu, Li, Peng, Zhou, Gao (2024)",
    "doi": "10.3390/w16111551"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 1.000,
    "width": 0.500,
    "height": 1.000,
    "wall_thickness": 0.020,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.70,
  "fill_height_m": 0.700,
  "fluid": {"name": "water", "density": 998, "dynamic_viscosity": 1.003e-3},
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "amplitude_deg": 2.0,
    "frequency_hz": 0.87,
    "frequency_ratio_f_f1": 1.0,
    "note": "Soft-spring frequency shift observed; 3D vortex waves at resonance"
  },
  "natural_frequency": {
    "mode": 1,
    "theoretical_hz": 0.872,
    "measured_hz": 0.869,
    "error_percent": 0.3,
    "regime": "deep"
  }
}
```

### preset_isope_lng — ISOPE LNG Benchmark Mark III Scale

```json
{
  "name": "ISOPE LNG Sloshing Benchmark - Mark III 2D Section",
  "source": {
    "paper": "Loysel, Chollet, Gervaise, Brosset (GTT), ISOPE 2012",
    "institutions": ["ECM", "ECN-BV", "GTT_x3", "Marintek", "PNU", "UDE", "UPM", "UR"],
    "note": "9-institution benchmark, largest sloshing comparison study"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.946,
    "width": 0.118,
    "height": 0.670,
    "wall_thickness": 0.005,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.80,
  "fluid": {"name": "water", "density": 998, "dynamic_viscosity": 8.94e-4},
  "excitation": {
    "type": "combined",
    "dof": ["Tx", "Ry", "Tz"],
    "platform": "hexapod_6dof",
    "conditions": 14,
    "analyzed_conditions": 7
  },
  "sensors": {
    "full_array": {"count": 72, "location": "ceiling", "types": ["PCB", "Kulite", "Kistler"]},
    "reduced_array": {"count": 16, "location": "key_impact_positions"}
  },
  "acquisition_frequency_khz": [15, 50],
  "scaled_variant": {
    "institution": "UR",
    "dimensions": {"length": 0.700, "width": 0.087, "height": 0.496},
    "scale_ratio": "1:1.35",
    "scaling_law": "Froude"
  },
  "full_scale_lng_context": {
    "tank_type": "lng_membrane",
    "cross_section": "octagonal",
    "membrane_systems": [
      {"name": "Mark_III", "primary": "SS_304L_corrugated", "thickness_mm": 1.2, "insulation": "PUF_160mm"},
      {"name": "NO96", "primary": "Invar_36Ni", "thickness_mm": 0.7, "insulation": "plywood_perlite"}
    ],
    "chamfer": {
      "upper_angle_deg": 135,
      "lower_angle_deg": 135,
      "upper_height_m": 8.63,
      "lower_height_m": 3.77
    },
    "capacity_range_m3": [120000, 270000],
    "tank_count": [4, 6],
    "design_standard": "IGC_Code",
    "design_life_wave_encounters": 1e8
  }
}
```

### preset_nasa_cylindrical_ring_baffle — NASA Ring Baffle Study

```json
{
  "name": "NASA MSFC Cylindrical Tank - Ring Baffle Parametric",
  "source": {
    "paper": "Yang, Sansone, Brodnick, Williams (2023), NASA MSFC",
    "ntrs": "20230005181",
    "experimental_ref": "Scholl et al. (cylindrical tank)"
  },
  "tank_type": "cylindrical_upright",
  "dimensions": {
    "diameter": 2.84,
    "radius": 1.42,
    "height": null,
    "wall_thickness": null,
    "note": "Height not explicitly stated; sufficient for sloshing tests"
  },
  "head": {
    "type": "flat",
    "bottom": "flat",
    "note": "Straight cylinder, bottom effect neglected in frequency formula"
  },
  "baffles": [
    {
      "type": "ring",
      "material": "rigid",
      "depth_ratio_d_R": [0.00, 0.10, 0.25],
      "width_ratio_W_R": null,
      "thickness": null,
      "porosity": 0.0,
      "count": 1,
      "design_guide": {
        "optimal_depth": "d/R → 0 (near surface) maximizes damping",
        "practical_range": "d/R = 0.05 ~ 0.10",
        "phase_shift_deg": {"d_R_0.25": 6.95, "d_R_0.10": 8.76, "d_R_0.00": 88.4},
        "note": "d/R=0 → ~90° phase separation, near-total wave blocking"
      }
    }
  ],
  "supports": {"type": "none"},
  "fill_level": 0.50,
  "fluid": {"name": "water", "density": 998, "dynamic_viscosity": 8.94e-4},
  "natural_frequency": {
    "formula": "omega^2 = (g * lambda1 / R) * tanh(lambda1 * h / R)",
    "lambda1": 1.841,
    "note": "First root of J'_1 (Bessel derivative), cylindrical 1st mode"
  },
  "pressure_model": {
    "max_pressure": "p_max = sqrt(2) * rho * g * eta",
    "nondimensional": "p_bar = p / (rho * g * eta), max = sqrt(2) ≈ 1.414",
    "note": "Analytical pressure load model for baffle design"
  },
  "aerospace_pmd_context": {
    "vanes": {"count": [2, 8], "typical": [4, 6], "angle_to_wall_deg": [45, 60]},
    "gallery_arms": {"count": [3, 6], "length_ratio_R": [0.5, 0.8], "mesh_um": [50, 200]},
    "baffle_rings": {"diameter_ratio_D": [0.7, 0.9], "spacing_ratio_H": [0.2, 0.4]},
    "diaphragm": "flexible membrane separating fuel/pressurant gas",
    "bladder": "balloon-like polymer fuel containment"
  }
}
```

### preset_zhao2024_lng_horizontal_ring — Zhao et al. LNG Horizontal with Ring Baffles

```json
{
  "name": "Zhao et al. (2024) - Elastic LNG Horizontal Tank with Ring Baffles",
  "source": {
    "paper": "Zhao, Wu, Yu, Haidn, Hu (2024), TU Munich",
    "arxiv": "2409.16226",
    "solver": "SPHinXsys",
    "validation_ref": "Grotle & Æsøy (2017)"
  },
  "tank_type": "cylindrical_horizontal",
  "dimensions": {
    "diameter": null,
    "length": null,
    "wall_thickness": 0.018,
    "note": "Grotle & Æsøy experimental tank (exact D from their paper)"
  },
  "head": {
    "type": "ellipsoidal_2to1",
    "note": "Semi-elliptical heads at 2:1 ratio"
  },
  "baffles": [
    {
      "type": "ring",
      "count": 2,
      "material": "rigid_or_elastic",
      "position": "equally_spaced"
    }
  ],
  "supports": {
    "type": "saddle",
    "note": "Horizontal tank standard; Zick method: A/R=0.25, θ=120°-150°"
  },
  "fill_level_hw_D": 0.255,
  "tank_material": {
    "name": "steel",
    "density": 7890,
    "poisson_ratio": 0.27,
    "youngs_modulus_GPa": 135
  },
  "fluid": {
    "water": {"density": 1000},
    "air": {"density": 1.226}
  },
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "amplitude_deg": 3.0,
    "frequency_hz": 0.5496,
    "period_s": 1.819,
    "rotation_axis": "tank_bottom_center_plane_intersection",
    "motion_function": "theta(t) = A * sin(2*pi*f*t)"
  },
  "sph_settings": {
    "solver": "SPHinXsys",
    "method": "WCSPH_multiphase_Riemann",
    "particle_spacing_m": [0.006, 0.0045],
    "fsi": "single_framework_SPH",
    "damping": "artificial_viscosity_first_1s",
    "eos": "p = c^2 * (rho - rho0), c = 10*Vmax",
    "vmax": "sqrt(2*g*hw)"
  },
  "elastic_baffle_study": {
    "hw_D": 0.50,
    "wall_thickness_m": 0.030,
    "baffle_type": "vertical_center",
    "variants": [
      {"name": "rigid", "E": "infinity"},
      {"name": "elastic_500kPa", "E_kPa": 500, "density": 2500, "poisson": 0.47},
      {"name": "elastic_50kPa", "E_kPa": 50, "density": 2500, "poisson": 0.47},
      {"name": "elastic_5kPa", "E_kPa": 5, "density": 2500, "poisson": 0.47}
    ],
    "finding": "More rigid → more effective sloshing suppression"
  },
  "horizontal_tank_parametric": {
    "typical_L_D_ratio": [2, 6],
    "optimal_L_D": [3, 4],
    "swash_plates": {"position": "L/3 and 2L/3", "opening_area_percent": [30, 50]},
    "transverse_baffles": {"spacing_ratio_L": [0.3, 0.5], "features": ["limber_holes_bottom", "air_holes_top"]},
    "saddle_design_zick": {
      "position_A_R": 0.25,
      "saddle_angle_deg": [120, 150],
      "saddle_width_ratio_L": [0.1, 0.2],
      "count": 2
    },
    "head_types": ["flat", "torispherical", "ellipsoidal_2to1", "hemispherical"]
  }
}
```

### preset_mdbc_validation — DualSPHysics mDBC Validation

```json
{
  "name": "DualSPHysics mDBC Validation Tank (English et al. 2021)",
  "source": {
    "paper": "English, Dominguez, Vacondio, Crespo et al. (2021)",
    "doi": "10.1007/s40571-021-00403-3"
  },
  "tank_type": "rectangular",
  "dimensions": {
    "length": 0.900,
    "width": 0.062,
    "height": 0.508,
    "wall_thickness": 0.005,
    "corner_radius": 0.0
  },
  "head": {"type": "flat"},
  "baffles": [],
  "supports": {"type": "none"},
  "fill_level": 0.183,
  "fill_height_m": 0.093,
  "fluid": {"name": "water", "density": 998, "dynamic_viscosity": 8.94e-4},
  "excitation": {
    "type": "rotation",
    "subtype": "pitch",
    "center": [0.45, 0.0, 0.0],
    "note": "Same as SPHERIC Test 10"
  },
  "sph_settings": {
    "solver": "DualSPHysics",
    "dp_options_m": [0.004, 0.002],
    "simulation_time_s": 7.0,
    "boundary_condition": "mDBC",
    "comparison_bc": "standard_DBC",
    "total_particles_dp002": 26791,
    "kernel": "Wendland_C2",
    "note": "mDBC shows superior pressure accuracy vs standard DBC"
  },
  "still_water_validation": {
    "tank_length_m": 2.400,
    "tank_height_m": 1.200,
    "water_height_m": 0.500,
    "bottom_obstacle": {
      "type": "trigonal_wedge",
      "height_m": 0.240,
      "position": "center_bottom"
    }
  }
}
```

---

## Cross-Reference: Parametric Study Variables Matrix

> Part C (docs/STL_DATASETS_AND_PARAMETRICS.md) 연동 — 논문별 sweep 변수 정리

| Paper | Swept Variable | Range | Fixed Variables |
|-------|---------------|-------|-----------------|
| **SPHERIC** | Fluid type (water/oil/glycerin) | 3 fluids | Tank, fill, excitation |
| **SPHERIC** | Fill level | 18%, 70% | Tank, fluid |
| **SPHERIC** | Tank breadth | 31, 62, 124 mm | L, H, fill |
| **Chen 2018** | Fill level | 13.8%, 30.8% | Tank, amplitude |
| **Chen 2018** | Excitation frequency | 0.5ω₁ ~ 2.0ω₁ | Tank, fill, amplitude |
| **Liu 2024** | Fill level | 20%, 30%, 70% | Tank, amplitude |
| **Liu 2024** | Excitation frequency | 0.8f₁ ~ 1.2f₁ | Tank, fill |
| **Liu 2024** | Pitch amplitude | 1°, 2°, 3° | Tank, fill, at resonance |
| **NASA 2023** | Baffle depth d/R | 0.00, 0.10, 0.25 | Tank, fill |
| **Zhao 2024** | Baffle type | ring vs. vertical | Tank, fill |
| **Zhao 2024** | Baffle elasticity (E) | rigid, 500, 50, 5 kPa | Tank, fill, baffle type |
| **ISOPE** | Excitation condition | 14 combinations (Tx+Ry+Tz) | Tank, fill |

### Response Variables (모든 논문 공통)

| Variable | Unit | Measured By |
|----------|------|-------------|
| Peak wall pressure | Pa, mbar | Pressure transducers |
| Free surface elevation η | m, mm | Wave gauge / camera |
| Impact force on wall | N | Force transducer / integrated pressure |
| Frequency response (PSD) | Pa²/Hz | FFT of pressure time series |
| Phase shift across baffle | deg | Pressure time correlation |
| Stress/strain distribution | Pa, — | FEM / SPH solid (Zhao 2024) |

---

## Natural Frequency Formula Summary (Quick Reference)

### Rectangular Tank

```
ω₁ = √(g · (π/L) · tanh(π·h/L))

Shallow (h/L < 0.1): ω₁ ≈ √(g·π·h/L)     ← linear in h
Deep    (h/L > 1.0): ω₁ ≈ √(g·π/L)         ← independent of h
Critical h/L:        ~0.24 (soft→hard spring transition, Chen 2018)
```

### Cylindrical Upright Tank

```
ω₁ = √(g · λ₁/R · tanh(λ₁·h/R))

λ₁ = 1.841  (1st root of J'₁, Bessel derivative)
R = D/2
```

### Cylindrical Horizontal Tank (approximate)

```
Effective length L_eff depends on fill level and cross-section geometry.
For shallow fill: treat as rectangular with L_eff ≈ chord length at fill height.
```

---

## Design Standards Reference (Part B 연동)

| Standard | Scope | Key Tanks |
|----------|-------|-----------|
| **ASME BPVC VIII Div.1** | Pressure vessel design | Cylindrical (upright/horizontal) |
| **API 650** | Welded tanks for oil storage | Cylindrical upright (≤45mm shell) |
| **API 620** | Large low-pressure storage | Cylindrical upright |
| **IGC Code** | Ships carrying liquefied gases | LNG membrane, Moss sphere, SPB |
| **DNV CN 31.12** | LNG Type B prismatic tank | SPB (IHI) |
| **NFPA 30** | Flammable liquids storage | All tank types |

### Head Type Quick Reference (ASME)

| Head Type | Crown Radius L | Knuckle Radius r | Depth | Use Case |
|-----------|----------------|-------------------|-------|----------|
| Flat | N/A | N/A | 0 | Cheapest, rectangular/small cylindrical |
| Torispherical (F&D) | 1.0D | 0.06D | ~0.194D | Most common pressure vessel |
| 2:1 Ellipsoidal | ≈0.9D | ≈0.17D | D/4 | Good volume/pressure compromise |
| Hemispherical | D/2 | N/A | D/2 | Highest pressure, highest cost |
| Conical | N/A | N/A | varies | Bottom drain, cone angle 30°-60° |

### Baffle Type Quick Reference

| Baffle Type | Key Parameter | Optimal Value | Best For |
|-------------|--------------|---------------|----------|
| Vertical (flat plate) | h_b/h_w | ≥ 0.8 | General sloshing reduction |
| Horizontal (shelf) | y_b/H | 0.3-0.7 | Near free surface placement |
| Perforated | Porosity ε | 20% | Energy dissipation |
| Ring (cylindrical) | d/R | 0.05-0.10 | Cylindrical/aerospace tanks |
| Swash plate (horizontal) | Opening area | 30-50% | Horizontal vessels |
| T-shaped | Combined | Varies | Enhanced lateral+vertical control |
