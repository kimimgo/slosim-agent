# EXP-1 Runbook — SPHERIC Test 10 실행·후처리·비교

## 1. DualSPHysics 실행 (GenCase → Solver)

### 1.1 파이프라인 개요

```
XML 케이스 → GenCase → Solver (GPU) → PART 파일 → MeasureTool → CSV → 분석
                                    → Gauge CSV (10kHz SWL/Force)
```

### 1.2 Docker 환경

```bash
# 볼륨 마운트
./cases:/cases        # XML 정의 파일
./simulations:/data   # 시뮬레이션 입출력

# 바이너리 (소문자 symlink 사용)
/opt/dsph/bin/gencase          # GenCase_linux64
/opt/dsph/bin/dualsphysics     # DualSPHysics5.4_linux64
/opt/dsph/bin/partvtk          # PartVTK_linux64
/opt/dsph/bin/measuretool      # MeasureTool_linux64
/opt/dsph/bin/isosurface       # IsoSurface_linux64
```

### 1.3 Run 실행 순서

```bash
# 1) XML을 cases/에 복사 (Docker 마운트 경로)
cp research-v2/exp1_spheric/cases/run_NNN_*.xml cases/RUN_NNN.xml

# 2) GenCase: XML → 바이너리 파티클
docker compose run --rm dsph gencase \
  /cases/RUN_NNN /data/exp1/run_NNN/run_NNN -createdirs:1

# 3) Solver: GPU 실행
docker compose run --rm dsph dualsphysics \
  /data/exp1/run_NNN/run_NNN /data/exp1/run_NNN -gpu -svres

# 4) MeasureTool: 압력 추출 (PART 파일 기반)
#    프로브 파일은 dp별로 선택:
#    Run 001 (dp=0.004): pressure_probes_lateral_dp004.txt
#    Run 002,004,005 (dp=0.002): pressure_probes_lateral_dp002.txt
#    Run 003 (dp=0.001): pressure_probes_lateral_dp001.txt
#    Run 006 (roof): pressure_probes_roof_dp002.txt
docker compose run --rm dsph measuretool \
  -dirin /data/exp1/run_NNN \
  -points /data/exp1/probes/pressure_probes_lateral_dpXXX.txt \
  -onlytype:-all,+fluid \
  -vars:-all,+press \
  -savecsv /data/exp1/run_NNN/PressureLateral

# mDBC 케이스 솔버 플래그 (Run 004, 005, 006):
docker compose run --rm dsph dualsphysics \
  /data/exp1/run_NNN/run_NNN /data/exp1/run_NNN -gpu -svres -mdbc_noslip
```

### 1.4 TimeOut 전략 (압력 시간해상도)

MeasureTool 해상도 = TimeOut 간격. 기본 0.01s (100Hz)는 피크 캡처에 부족.

**Run별 전략:**

| Run | dp | 입자 수 | TimeOut | PART 파일 수 | 예상 디스크 |
|-----|----|---------|---------|-----------|----|
| 001 | 0.004 | 136K | 0.001 (1kHz) | 7,000 | ~7GB |
| 002 | 0.002 | 545K | 0.001 (1kHz) | 7,000 | ~28GB |
| 003 | 0.001 | 2.2M | 0.01 (100Hz) | 700 | ~25GB |
| 004 | 0.002 | 545K | 0.001 (1kHz) | 7,000 | ~28GB |
| 005 | 0.002 | 545K | 0.001 (1kHz) | 7,000 | ~28GB |
| 006 | 0.002 | 545K | 0.001 (1kHz) | 7,000 | ~28GB |

**Run 003 (fine, 2.2M)은 100Hz로 제한** — 디스크 제약. 피크 값은 Force gauge (10kHz)로 보완.

### 1.5 XML 수정 사항 (TimeOut 변경)

`<parameter key="TimeOut" value="0.001" />` 로 변경 (Run 001, 002, 004, 005, 006)
Run 003은 `value="0.01"` 유지.

### 1.6 Gauge 출력 (인라인, 10kHz)

Solver 실행 시 자동 생성:
- `GaugesSWL_SWL_{Left,Center,Right}.csv` — 수면 높이 (현재 0 문제 있음, 미해결)
- `GaugesForce_Force_Wall.csv` — 벽면 전체 힘 (정상 작동, 10kHz)

Force gauge는 **총 벽력**이므로 점 압력 대체 불가하지만, **피크 타이밍 식별**에 유용.

---

## 2. 후처리 — 압력 추출 + 실험 비교

### 2.1 압력 프로브 위치 (핵심 규칙)

> **벽에서 1.5h 떨어진 유체 영역에 프로브 배치** (DualSPHysics 공식 가이드)
> 경계 입자 간극(boundary gap) 안에 프로브를 놓으면 보간 실패.

```
h = coefh × sqrt(3) × dp = 1.2 × 1.732 × dp ≈ 2.078 × dp
1.5h ≈ 3.117 × dp
```

| dp | h | 1.5h | 프로브 x좌표 (좌벽에서) |
|----|---|------|----------------------|
| 0.004 | 0.0083 | 0.0125 | **0.013** |
| 0.002 | 0.0042 | 0.0063 | **0.007** |
| 0.001 | 0.0021 | 0.0031 | **0.004** |

### 2.2 SPHERIC 실험 센서 위치

논문 기준: Sensor 1 (P1) = 좌측 벽, 정수면 높이 (z ≈ 0.093m)
Sensor 3 (P3) = 천장 중앙 (z ≈ 0.508m)

**프로브 파일 (dp별 조정 필요):**

```
# pressure_probes_lateral_dpXXX.txt
POINTS
{1.5h}  0.031  0.093   #P1_main (left wall, at still water level)
{1.5h}  0.031  0.083   #P1_below10mm
{1.5h}  0.031  0.103   #P1_above10mm
{L-1.5h}  0.031  0.093 #P1_right (symmetric check)
```

### 2.3 실험 데이터 형식

**위치**: `datasets/spheric/case_1/`

| 파일 | 내용 | 형식 | 샘플링 |
|------|------|------|--------|
| `lateral_water_1x.txt` | Water Lateral 시계열 | Tab-sep, 6열 | 20kHz (dt=50μs) |
| `lateral_oil_1x.txt` | Oil Lateral 시계열 | Tab-sep, 6열 | 20kHz |
| `roof_water_1x.txt` | Water Roof 시계열 | Tab-sep, 6열 | 20kHz |
| `Water_4first_peak_*.txt` | 피크 통계 (103회 반복) | Tab-sep, 4열 | — |

**시계열 열:**
1. Time [s]
2. **Pressure [mbar]** ← 비교 대상
3. Position_smooth_splines [deg]
4. Velocity [deg/s]
5. Acceleration [deg/s²]
6. Position_original [deg]

**단위 변환**: 실험 mbar → SI Pa : `P_Pa = P_mbar × 100`

### 2.4 비교 절차

```python
# 의사코드
import numpy as np
from scipy.signal import find_peaks, correlate

# 1. 실험 데이터 로드
exp = load_tsv("datasets/spheric/case_1/lateral_water_1x.txt")
exp_t, exp_p = exp[:, 0], exp[:, 1] * 100  # mbar → Pa

# 2. 시뮬레이션 압력 로드
sim = load_csv("simulations/exp1/run_NNN/Pressure_Press.csv")
sim_t, sim_p = sim[:, 1], sim[:, 2]  # Time, Press_0 (P1_main)

# 3. 피크 검출
exp_peaks, _ = find_peaks(exp_p, height=500, distance=15000)  # 최소 0.75s 간격
sim_peaks, _ = find_peaks(sim_p, height=500, distance=750)     # 1kHz 기준

# 4. 첫 4개 피크 비교
for i in range(min(4, len(sim_peaks))):
    sim_val = sim_p[sim_peaks[i]]
    # 실험 피크 통계 (103회)
    exp_peaks_stats = load_tsv("Water_4first_peak_*.txt")
    exp_mean = np.mean(exp_peaks_stats[:, i])
    exp_std = np.std(exp_peaks_stats[:, i])

    # M1: Peak-in-band
    in_band = (exp_mean - 2*exp_std) <= sim_val/100 <= (exp_mean + 2*exp_std)

    # M2: Mean absolute peak error
    error = abs(sim_val/100 - exp_mean) / exp_mean

# 5. Cross-correlation (M5, M6)
# 시뮬레이션을 실험 샘플링에 리샘플링
sim_interp = np.interp(exp_t, sim_t, sim_p)
corr = correlate(sim_interp, exp_p, mode='full')
lag = np.argmax(corr) - len(exp_p) + 1
tau = lag * 0.00005  # 실험 dt
r_max = np.max(corr) / (np.std(sim_interp) * np.std(exp_p) * len(exp_p))
```

### 2.5 수렴 분석 (Run 001/002/003)

```python
# Richardson Extrapolation + GCI
dp = [0.004, 0.002, 0.001]
p_peak = [peak_001, peak_002, peak_003]  # 각 run의 첫 피크

r = dp[0] / dp[1]  # refinement ratio = 2
p_order = np.log(abs((p_peak[0]-p_peak[1])/(p_peak[1]-p_peak[2]))) / np.log(r)

# GCI (ASME V&V 20)
Fs = 1.25  # safety factor (3 grids)
GCI_fine = Fs * abs((p_peak[1]-p_peak[2])/p_peak[2]) / (r**p_order - 1)
# M4: GCI_fine < 10% → PASS
```

---

## 3. ParaView 후처리 (pv MCP 서버)

### 3.1 전제 조건

시뮬레이션 완료 후 PartVTK로 VTK 변환 필요:

```bash
# Fluid VTK 생성
docker compose run --rm dsph partvtk \
  -dirin /data/exp1/run_NNN \
  -savevtk /data/exp1/run_NNN/vtk/PartFluid \
  -onlytype:-all,+fluid \
  -vars:+vel,+rhop,+press

# IsoSurface (자유 표면 메시)
docker compose run --rm dsph isosurface \
  -dirin /data/exp1/run_NNN \
  -saveiso /data/exp1/run_NNN/iso/Surface \
  -vars:+vel,+press,+rhop \
  -onlytype:+fluid
```

### 3.2 pv MCP 호출 (렌더링)

```
# 1) 데이터 검사
pv.inspect_data(file_path="simulations/exp1/run_NNN/vtk/PartFluid_0100.vtk")

# 2) 압력 필드 스냅샷
pv.render(
  file_path="simulations/exp1/run_NNN/vtk/PartFluid_*.vtk",
  color_by="Press",
  colormap="Cool to Warm",
  output="research-v2/exp1_spheric/figures/pressure_t2.png"
)

# 3) 슬라이스 (중앙 Y 평면)
pv.slice(
  file_path="simulations/exp1/run_NNN/iso/Surface_*.vtk",
  normal=[0, 1, 0],
  origin=[0.45, 0.031, 0.25],
  color_by="Press"
)

# 4) 시계열 애니메이션
pv.animate(
  file_path="simulations/exp1/run_NNN/iso/Surface_*.vtk",
  color_by="Press",
  output="research-v2/exp1_spheric/figures/sloshing_pressure.gif"
)

# 5) 멀티패널 (압력 + 속도 + 시계열 그래프)
pv.split_animate(
  file_path="simulations/exp1/run_NNN/iso/Surface_*.vtk",
  layout="2x2",
  panes=[
    {"type": "render", "color_by": "Press"},
    {"type": "render", "color_by": "Vel"},
    {"type": "plot", "csv": "simulations/exp1/run_NNN/Pressure_Press.csv", "x": "Time", "y": "Press_0"},
    {"type": "plot", "csv": "simulations/exp1/run_NNN/GaugesForce_Force_Wall.csv", "x": "time", "y": "forcex"}
  ]
)
```

### 3.3 수렴 비교 시각화

```
# dp=0.004 / 0.002 / 0.001 동일 시점 비교
pv.split_animate(
  layout="1x3",
  panes=[
    {"file": "run_001/iso/Surface_0100.vtk", "color_by": "Press", "title": "dp=0.004"},
    {"file": "run_002/iso/Surface_0100.vtk", "color_by": "Press", "title": "dp=0.002"},
    {"file": "run_003/iso/Surface_0100.vtk", "color_by": "Press", "title": "dp=0.001"}
  ]
)
```

---

## 실행 체크리스트

- [ ] XML TimeOut 수정 (0.01 → 0.001, Run 003 제외)
- [ ] dp별 압력 프로브 파일 생성 (1.5h 규칙)
- [ ] Run 001 실행 + 압력 추출 검증
- [ ] Run 002-006 순차 실행
- [ ] MeasureTool 압력 CSV → 피크 검출
- [ ] 실험 데이터 비교 (M1-M8 메트릭)
- [ ] PartVTK + IsoSurface 생성
- [ ] pv MCP 렌더링 + 애니메이션
- [ ] validation_report.md 작성
