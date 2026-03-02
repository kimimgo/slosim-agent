# DualSPHysics Case Design

description: DualSPHysics 슬로싱 케이스 설계 전문 스킬. dp(입자 간격) 선택 의사결정, 경계 조건(DBC/mDBC) 선택, 게이지/프로브 설정, mDBC normals 설정, simulation domain 설계를 가이드. "케이스 설계", "case design", "dp 선택", "수렴 연구 설계", "XML 설계" 시 발동.

---

## 1. dp 선택 의사결정 트리

### 1.1 기본 규칙

```
수심 입자층 수 = fill_height / dp
최소 23층 (dp=0.004) → 조잡, 정성적 확인용
권장 46층 (dp=0.002) → 정량 검증 가능
고해상도 93층 (dp=0.001) → 수렴 확인용, GPU 메모리/디스크 주의
```

### 1.2 목적별 dp 추천

| 목적 | dp [m] | 입자 수 (3D, 0.9×0.062×0.093m) | 예상 시간 (RTX 4090) | 디스크 |
|------|--------|-------------------------------|---------------------|--------|
| 정성적 확인 | 0.004 | ~80K | ~5분 | ~1GB |
| 정량 검증 | 0.002 | ~650K | ~30분 | ~10GB |
| 수렴 확인 | 0.001 | ~5M | ~4시간 | ~100GB+ |

### 1.3 수렴 연구 설계

반드시 **3개 해상도** (2:1 비율):
- coarse: dp = 0.004 (빠른 검증)
- medium: dp = 0.002 (기준)
- fine: dp = 0.001 (수렴 확인)

이유: Richardson extrapolation + GCI 계산에 최소 3개 해상도 필요 (2개는 가정된 수렴 차수 사용 → 보수적).

**dp=0.001 주의사항**:
- `TimeOut=0.001`이면 7000개 × 187MB = **1.3TB** → `TimeOut=0.01` 사용
- MeasureTool 출력은 TimeOut에 의존 → 100Hz 제한됨
- 게이지 CSV는 `computedt=0.0001` (10kHz)로 별도 확보

### 1.4 SPH 문헌 dp 비교

| 논문 | dp | 입자 수 | 비고 |
|------|-----|---------|------|
| English et al. (2022) | 0.004, 0.002 | ~80K, ~650K | mDBC 검증, 정성적 비교만 |
| Green & Peiró (2018) | GPU 해상도 | — | 에너지 소산 분석 |
| Flow3D (격자 기반) | — | — | RMS 4.79%, 격자 수렴 연구 |
| **우리 (EXP-1)** | **0.004/0.002/0.001** | **80K/650K/5M** | **3-level 수렴, M1-M8 정량적** |

## 2. 경계 조건 선택 매트릭스

### 2.1 DBC vs mDBC

```
DBC (Boundary=1): 안정적, ~5mm 인공 감쇠층. lateral impact에 적합.
mDBC (Boundary=2): 감쇠 적음, 정확. normals 설정 필수, 불안정 가능.
```

### 2.2 의사결정 매트릭스

| 유동 유형 | 정확도 요구 | 추천 BC | 이유 |
|-----------|-----------|---------|------|
| Lateral water impact | 중 | **DBC** | 안정적, 피크 캡처 충분 |
| Lateral oil impact | 중 | **DBC** | 작은 피크에서 안정성 우선 |
| Roof impact | 중-고 | **DBC** | mDBC 공명 불안정 위험 |
| 정밀 벽면 압력 | 고 | **mDBC** | 감쇠 없는 측정 |
| 공명 운동 | — | **DBC** | mDBC 발산 위험 |

### 2.3 mDBC normals 설정 (검증된 dp/2 오프셋 기법)

```xml
<normals active="true">
    <!-- 핵심: point를 dp/2만큼 유체쪽으로 오프셋 -->
    <!-- Bottom face (z=0): point z = dp/2 = 0.001 -->
    <norplane mkbound="0">
        <onlypos>
            <posmin x="-0.01" y="-0.01" z="-0.01" />
            <posmax x="0.910" y="0.072" z="0.005" />
        </onlypos>
        <point x="0.45" y="0.031" z="0.001" />
        <normal x="0" y="0" z="1" />
        <maxdisth v="2" />
        <clear v="false" />
    </norplane>
    <!-- 각 면(Left/Right/Front/Back)에 동일 패턴 적용 -->
</normals>
```

**주의사항**:
- `setfrdrawmode auto="true"` 는 normals와 비호환 → 사용 금지
- `<point>`는 boundary interface 위치 (입자 위치 아님). 입자 위에 point를 놓으면 distance=0 → normal=(0,0,0)
- `maxdisth v="0"` = 거리 제한 없음 (포럼 확인)
- 공명 운동 + mDBC = 불안정 가능 (Free-slip, No-slip 모두)

## 3. 게이지/프로브 설정

### 3.1 computedt 선택

```xml
<default>
    <savevtkpart value="false" />
    <computedt value="0.0001" units_comment="10kHz gauge output" />
    <computetime start="0" end="7" />
    <output value="true" />
    <outputdt value="0.0001" units_comment="10kHz CSV output" />
</default>
```

**핵심**: `computedt=0.0001` (10kHz). 실험이 20kHz이면 절반은 확보해야 임팩트 피크를 캡처.
v1에서 `computedt=0.005` (200Hz) → 피크 놓침 → v2에서 10kHz로 수정하여 해결.

### 3.2 TimeOut vs computedt 분리 전략

```
TimeOut  = PART 파일 출력 간격 (바이너리 파티클, 수백MB/개)
computedt = Gauge CSV 출력 간격 (텍스트, 경량)
```

**함정**: dp=0.001에서 `TimeOut=0.001`이면 7000개 × 187MB = 1.3TB.
**해결**: `TimeOut=0.01` (100Hz)로 디스크 절약, `computedt=0.0001`로 게이지는 10kHz 유지.

### 3.3 프로브 배치 (1.5h rule)

```
프로브 위치 = 벽면 + 1.5 × h (smoothing length)
h ≈ coefh × sqrt(3) × dp = 1.2 × 1.732 × dp ≈ 2.078 × dp
1.5h ≈ 3.12 × dp
```

dp=0.002이면 벽에서 최소 6.24mm (실질적으로 프로브 x ≥ 0.005).

**프로브 파일 예시** (`spheric_probes.txt`):
```
POINTS
0.005 0.031 0.046
```

### 3.4 프로브 파일 생성 로직

```python
def generate_probe_file(tank_length, tank_width, fill_height, dp, probe_type='lateral'):
    h = 1.2 * (3**0.5) * dp
    offset = 1.5 * h

    if probe_type == 'lateral':
        # 벽면 x=0 에서 offset만큼 떨어진 위치, 수면 아래
        x = round(offset, 4)
        y = round(tank_width / 2, 4)
        z = round(fill_height / 2, 4)
    elif probe_type == 'roof':
        # 천장 z=tank_height 에서 offset만큼 아래
        x = round(tank_length / 2, 4)
        y = round(tank_width / 2, 4)
        z = round(fill_height - offset, 4)

    return f"POINTS\n{x} {y} {z}"
```

## 4. Simulation Domain

```xml
<simulationdomain>
    <posmin x="default - 15%" y="default - 10%" z="default - 10%" />
    <posmax x="default + 15%" y="default + 10%" z="default + 80%" />
</simulationdomain>
```

- z 방향 +80%: sloshing wave가 튀어올라도 simulation domain 이탈 방지
- y 방향 ±10%: mDBC 사용 시 경계 입자 이탈 방지
- x 방향 ±15%: 격렬한 수평 운동 대비

## 5. Go 이식 인터페이스

```go
// internal/llm/tools/case_design.go

type DPConfig struct {
    DP            float64 `json:"dp"`
    Layers        int     `json:"layers"`
    EstParticles  int     `json:"est_particles"`
    EstTimeMin    int     `json:"est_time_min"`
    EstDiskGB     float64 `json:"est_disk_gb"`
}

type CaseConfig struct {
    DP        float64 `json:"dp"`
    BC        string  `json:"bc"`         // "DBC" or "mDBC"
    TimeOut   float64 `json:"time_out"`
    ComputeDT float64 `json:"compute_dt"`
    ProbeFile string  `json:"probe_file"`
}

type CaseDesigner struct {
    TankLength  float64
    TankWidth   float64
    FillHeight  float64
}

func (d *CaseDesigner) SelectDP(purpose string) DPConfig {
    // purpose: "qualitative", "quantitative", "convergence"
    // Returns recommended dp and estimated resources
}

func (d *CaseDesigner) DesignConvergenceMatrix(baseDP float64) []CaseConfig {
    // Returns 3 configs: coarse (2×baseDP), medium (baseDP), fine (baseDP/2)
    // Each with appropriate TimeOut/computedt settings
}

func (d *CaseDesigner) SelectBC(flowType, accuracy string) string {
    // flowType: "lateral_water", "lateral_oil", "roof", "resonance"
    // accuracy: "medium", "high"
    // Returns "DBC" or "mDBC"
}

func (d *CaseDesigner) GenerateProbeFile(dp float64, probeType string) string {
    // probeType: "lateral", "roof"
    // Returns probe file content with 1.5h rule
}
```
