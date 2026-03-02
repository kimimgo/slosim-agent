# DualSPHysics Post-Processing & Analysis

description: DualSPHysics 시뮬레이션 후처리 및 분석 스킬. MeasureTool CSV 파싱, SPH 노이즈 필터링(Median+MovingAvg), 교차상관(Cross-Correlation) 시간 정렬, 피크 검출 및 실험 데이터 비교를 가이드. "후처리", "분석", "필터링", "시계열 비교", "피크 검출", "postprocess", "analysis" 시 발동.

---

## 1. MeasureTool 후처리

### 1.1 실행

```bash
docker compose run --rm dsph measuretool \
    -dirin /data/exp1/run_002 \
    -points /cases/spheric_probes.txt \
    -onlytype:-all,+fluid \
    -vars:press \
    -savecsv /data/exp1/run_002/PressConsistent
```

### 1.2 프로브 파일 규칙 (1.5h rule)

```
프로브 위치 = 벽면 + 1.5 × h (smoothing length)
h ≈ coefh × sqrt(3) × dp = 1.2 × 1.732 × dp ≈ 2.078 × dp
1.5h ≈ 3.12 × dp
```

dp=0.002이면 벽에서 최소 6.24mm. 벽 너무 가까우면 0 또는 노이즈만 나옴.
**중요**: 모든 run에서 동일 프로브 파일 사용 → `PressConsistent` 네이밍 규칙.

### 1.3 CSV 파싱

```python
# MeasureTool CSV: 세미콜론 구분, 헤더 4줄
# Line 1: ;PosX [m]:;0.005;...
# Line 2: ;PosY [m]:;0.031;...
# Line 3: ;PosZ [m]:;0.046;...
# Line 4: Part;Time [s];Press_0 [Pa];...

data = np.genfromtxt(path, delimiter=';', skip_header=4)
t = data[:, 1]          # Time [s]
p = data[:, 2] / 100    # Press_0 [Pa] → [mbar]
```

**단위 변환**: Pa → mbar: `/ 100`

### 1.4 프로브 위치 스냅 문제

**함정**: MeasureTool은 요청한 위치에서 가장 가까운 입자를 찾아 측정.
dp가 다르면 스냅 위치가 달라짐 → 다른 run 간 비교가 불공정.

**해결**: CSV 헤더의 PosX/PosY/PosZ에서 실제 스냅 위치 확인.
모든 run에서 동일 프로브 파일 사용.

## 2. SPH 노이즈 필터링

### 2.1 원리

SPH는 입자법 특성상 압력에 고주파 노이즈가 포함됨.
**적용 순서**: Median filter (스파이크 제거) → Moving average (스무딩)

### 2.2 케이스별 필터 파라미터

```python
# Water (큰 피크 ~30-40 mbar): 강한 필터
FILTER_WATER = {'medfilt_kernel': 5, 'smooth_window': 11}

# Oil (작은 피크 ~5 mbar): 약한 필터
FILTER_OIL = {'medfilt_kernel': 3, 'smooth_window': 7}
```

**핵심**: Oil은 Water보다 약하게 필터링. Oil M1이 ±0.016 mbar 차이로 IN/OUT 갈림.

### 2.3 구현

```python
from scipy.signal import medfilt
from scipy.ndimage import uniform_filter1d

def filter_pressure(p, fluid='water'):
    params = {
        'water': {'medfilt_kernel': 5, 'smooth_window': 11},
        'oil':   {'medfilt_kernel': 3, 'smooth_window': 7},
    }
    fp = params.get(fluid, params['water'])

    # 1) Median filter: isolated spike 제거
    p_filtered = medfilt(p, kernel_size=fp['medfilt_kernel'])
    # 2) Moving average: residual noise smoothing
    p_filtered = uniform_filter1d(p_filtered, size=fp['smooth_window'])

    return p_filtered
```

### 2.4 과필터링 감지

**증상**: 피크가 깎여서 M1 FAIL
**감지**: 필터 전/후 피크 높이 비교 — 20% 이상 감소하면 과필터링 경고
**대응**: kernel/window 크기 축소

```python
def check_overfilter(p_raw, p_filtered, threshold=0.20):
    from scipy.signal import find_peaks
    peaks_raw = find_peaks(p_raw, height=5.0)[0]
    peaks_filt = find_peaks(p_filtered, height=5.0)[0]
    if len(peaks_raw) == 0:
        return False
    max_raw = max(p_raw[peaks_raw])
    max_filt = max(p_filtered[peaks_filt]) if len(peaks_filt) > 0 else 0
    reduction = 1 - max_filt / max_raw
    if reduction > threshold:
        print(f"WARNING: peak reduction {reduction:.1%} > {threshold:.0%}")
        return True
    return False
```

## 3. Cross-Correlation (시간 정렬)

### 3.1 알고리즘

```python
from scipy.signal import correlate
import numpy as np

def cross_correlation(t_sim, p_sim, t_exp, p_exp, max_shift_s=2.0):
    """Compute cross-correlation with optimal time shift.

    Returns: (r_max, tau) where r_max is max correlation, tau is time shift.
    """
    dt_sim = np.median(np.diff(t_sim))
    dt_exp = np.median(np.diff(t_exp))

    # 공통 시간축으로 보간
    dt = max(dt_sim, dt_exp)
    t_common = np.arange(0, min(t_sim[-1], t_exp[-1]), dt)
    p_sim_i = np.interp(t_common, t_sim, p_sim)
    p_exp_i = np.interp(t_common, t_exp, p_exp)

    # 정규화
    p_sim_n = (p_sim_i - np.mean(p_sim_i)) / (np.std(p_sim_i) + 1e-12)
    p_exp_n = (p_exp_i - np.mean(p_exp_i)) / (np.std(p_exp_i) + 1e-12)

    n = len(t_common)
    corr = np.correlate(p_sim_n, p_exp_n, mode='full') / n
    lags = np.arange(-n + 1, n) * dt

    # 물리적 범위 제한 (±2초)
    mask = np.abs(lags) <= max_shift_s
    idx_max = np.argmax(corr[mask])
    return corr[mask][idx_max], lags[mask][idx_max]
```

### 3.2 결과 해석

| r_max | 해석 |
|-------|------|
| > 0.8 | 우수 (파형이 매우 유사) |
| 0.5-0.8 | 양호 (주요 특성 일치) |
| < 0.5 | 불량 (상당한 차이) |

- `tau > 0`: 시뮬레이션이 실험보다 느림 (SPH에서 흔함)
- `tau < 0`: 시뮬레이션이 실험보다 빠름
- `|tau| < 1 period`: M6 PASS 기준

## 4. 피크 검출 및 비교

### 4.1 피크 검출

```python
from scipy.signal import find_peaks

def detect_peaks(t, p, min_height, min_dist_s):
    """Detect pressure peaks with minimum height and distance.

    Args:
        min_height: minimum peak height [mbar]
        min_dist_s: minimum distance between peaks [seconds]
    """
    dt = np.median(np.diff(t))
    min_dist = int(min_dist_s / dt) if dt > 0 else 1000
    peaks, props = find_peaks(p, height=min_height, distance=min_dist)
    return peaks
```

### 4.2 케이스별 파라미터

| 케이스 | min_height [mbar] | min_dist_s [s] | 비고 |
|--------|------------------|----------------|------|
| Water Lateral | 10.0 | 1.2 | 큰 임팩트 피크 |
| Oil Lateral | 2.0 | 1.0 | 작은 피크, 낮은 임계값 |
| Water Roof | 5.0 | 0.8 | 중간 크기 피크 |

### 4.3 실험 피크 통계 로딩

SPHERIC Test 10은 102회 반복 실험 → 각 피크의 평균(μ)과 표준편차(σ) 제공.

```python
def load_exp_peaks(peak_file):
    """Load experimental peak statistics (102 repetitions).

    Returns: (means, stds) arrays of shape (n_peaks,)
    """
    peak_data = np.genfromtxt(peak_file, delimiter='\t', skip_header=2)
    means = np.nanmean(peak_data, axis=0)
    stds = np.nanstd(peak_data, axis=0)
    return means, stds
```

### 4.4 M1: Peak-in-band 분석

```python
def peak_in_band(sim_peaks, exp_means, exp_stds, n_sigma=2):
    """Check if simulation peaks fall within exp mean ± n*sigma."""
    n_peaks = min(len(sim_peaks), len(exp_means))
    in_band = 0
    results = []
    for i in range(n_peaks):
        lo = exp_means[i] - n_sigma * exp_stds[i]
        hi = exp_means[i] + n_sigma * exp_stds[i]
        is_in = lo <= sim_peaks[i] <= hi
        if is_in:
            in_band += 1
        results.append({
            'peak': i + 1,
            'sim': sim_peaks[i],
            'exp_mean': exp_means[i],
            'exp_2sigma': 2 * exp_stds[i],
            'in_band': is_in,
        })
    return in_band, n_peaks, results
```

### 4.5 Roof 특수 처리: 다중 프로브

Roof 측정은 여러 프로브 위치를 시도하여 최적 매칭을 찾아야 함.

```python
def select_best_probe(t_sim, probes_dict, t_exp, p_exp):
    """Select probe with highest cross-correlation to experiment."""
    best_r, best_col = -1, 0
    for col_idx, (p_sim, x_pos) in probes_dict.items():
        r, tau = cross_correlation(t_sim, p_sim, t_exp, p_exp)
        if r > best_r:
            best_r = r
            best_col = col_idx
    return best_col, best_r
```

## 5. 전체 분석 파이프라인 요약

```
1. MeasureTool 실행 (동일 프로브 파일)
2. CSV 로딩 → Pa를 mbar 변환
3. 필터링 (케이스별 Median + MovingAvg)
4. 과필터링 검증
5. Cross-Correlation → 시간 정렬 (M5, M6)
6. 피크 검출 → 실험 비교 (M1, M2)
7. (Roof) 다중 프로브 선택
8. → dsph-metrics 스킬로 정량 검증
```

## 6. 흔한 함정

| 함정 | 증상 | 해결 |
|------|------|------|
| 프로브 벽 너무 가까움 | 0 또는 노이즈만 | 1.5h rule 적용 |
| 프로브 dp 스냅 | run 간 다른 위치 비교 | 동일 프로브 파일 + 헤더 확인 |
| Oil 과필터링 | M1 FAIL (±0.016 mbar) | Med3+Avg7 (Water보다 약하게) |
| numpy bool JSON | TypeError | `default=lambda x: bool(x)` handler 추가 |
| CSV 구분자 | 파싱 오류 | 세미콜론 `;` 확인 |

## 7. Go 이식 인터페이스

```go
// internal/llm/tools/analysis.go

type FilterConfig struct {
    MedfiltKernel int     `json:"medfilt_kernel"`
    SmoothWindow  int     `json:"smooth_window"`
    Fluid         string  `json:"fluid"`  // "water" or "oil"
}

type CorrResult struct {
    RMax float64 `json:"r_max"`
    Tau  float64 `json:"tau_s"`
}

type Peak struct {
    Index    int     `json:"index"`
    Time     float64 `json:"time_s"`
    Pressure float64 `json:"pressure_mbar"`
}

type PeakConfig struct {
    MinHeight float64 `json:"min_height_mbar"`
    MinDistS  float64 `json:"min_dist_s"`
}

type PeakStats struct {
    Means []float64 `json:"means_mbar"`
    Stds  []float64 `json:"stds_mbar"`
}

type Analyzer struct {
    SimDir string
    ExpDir string
}

func (a *Analyzer) LoadGaugeCSV(path string) (times, pressures []float64, err error) {
    // Parse semicolon-delimited CSV, skip 4 header lines
    // Convert Pa → mbar
}

func (a *Analyzer) FilterPressure(data []float64, config FilterConfig) []float64 {
    // Apply median filter → moving average
    // Returns filtered pressure series
}

func (a *Analyzer) CrossCorrelation(sim, exp TimeSeries) CorrResult {
    // Interpolate to common time axis
    // Normalize, compute correlation, find max within ±2s
}

func (a *Analyzer) DetectPeaks(data TimeSeries, config PeakConfig) []Peak {
    // Find peaks with min height and min distance constraints
}

func (a *Analyzer) PeakInBand(simPeaks []Peak, expStats PeakStats, nSigma float64) (inBand, total int) {
    // Check each sim peak against exp mean ± nSigma*std
}
```

## 8. 참조 스크립트

| 스크립트 | 역할 |
|---------|------|
| `research-v2/exp1_spheric/analysis/convergence_analysis.py` | Water Lateral 수렴 분석 |
| `research-v2/exp1_spheric/analysis/oil_roof_analysis.py` | Oil + Roof 분석 |
| `research-v2/exp1_spheric/analysis/paper_figures.py` | 논문용 그래프 생성 |
| `research-v2/exp1_spheric/scripts/final_verdict.py` | 전체 PASS/FAIL 판정 |
