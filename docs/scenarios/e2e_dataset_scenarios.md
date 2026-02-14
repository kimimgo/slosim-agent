# E2E 유저 시나리오: 논문 데이터셋 기반 (10개)

> 각 시나리오는 **실제 사용자가 slosim-agent에 입력할 자연어** + **검증 기준**으로 구성.
> SPHERIC Test 10 (Low/High Fill Water)는 이미 검증 완료. 나머지 데이터셋 커버.

---

## Scenario 1: SPHERIC Test 10 — Oil Low Fill (점성 유체)

**데이터셋**: SPHERIC Benchmark, Sunflower Oil, 18.3% fill
**난이도**: 중 (유체 물성 변경)

```
유저: SPHERIC 벤치마크 탱크(900x62x508mm)에 해바라기 오일을 18% 채워서
      슬로싱 해석해줘. 오일 밀도 990, 점도 0.045 Pa·s.
      주기 1.535초로 바닥 중심 기준 피치 회전 4도.
```

### 검증 기준
- [ ] XML: `rhop0 value="990"`, `Visco` 값이 오일 점도 반영 (인공점도 α ≈ 0.1-0.5)
- [ ] 여진: `mvrotsinu` freq=0.651 Hz, ampl=4°, axisp = [0.45, ±0.031, 0]
- [ ] 솔버 완주 (7초), 파티클 이탈 < 5%
- [ ] 물 대비 자유표면 진동 감쇠 확인 (오일의 높은 점도 효과)
- [ ] 압력 CSV: Sensor 1 위치 (좌벽 하단) 시계열 추출

### 핵심 파라미터
| 항목 | 값 |
|------|-----|
| dp | 0.004 m |
| 파티클 수 | ~136K |
| TimeMax | 7.0 s |
| 예상 GPU 시간 | ~3분 (RTX 4090) |

---

## Scenario 2: Chen 2018 — Shallow Fill Sway Excitation (얕은 수심 수평 가진)

**데이터셋**: Chen & Xue (2018), OpenFOAM 검증, h=83mm (13.8%)
**난이도**: 중 (수평 병진 운동)

```
유저: 600x300x650mm 직사각형 탱크에 물 83mm 높이로 채워줘.
      수평 방향(x축)으로 진폭 7mm, 주파수 0.756Hz로 좌우 흔들어.
      10초 동안 해석.
```

### 검증 기준
- [ ] XML: tank 600x300x650mm, fluid height 83mm
- [ ] 여진: `mvrectsinu` (sway), ampl x=0.007, freq x=0.756 Hz
- [ ] `pointmin`/`simulationdomain` x 방향 ±7mm 이상 마진
- [ ] 솔버 완주, shallow water wave 패턴 확인
- [ ] h/L = 0.138 → shallow regime → 비선형 파형 (hydraulic bore/breaking)

### 핵심 파라미터
| 항목 | 값 |
|------|-----|
| dp | 0.005 m |
| 파티클 수 | ~120K |
| TimeMax | 10.0 s |
| 자연주파수 ω₁ | 4.749 rad/s (f₁ = 0.756 Hz) |
| 가진/자연주파수 비 | 1.0 (공진) |

---

## Scenario 3: Chen 2018 — Near-Critical Fill (전이 수심)

**데이터셋**: Chen & Xue (2018), h=185mm (30.8%), soft→hard spring 전이
**난이도**: 상 (비선형 전이 구간)

```
유저: 같은 탱크(600x300x650mm)인데 이번엔 물 185mm로 채워줘.
      수평 진동 7mm, 자연주파수(1.008Hz)로 가진.
      이 수위가 critical depth 근처라서 비선형 효과가 클 거야.
```

### 검증 기준
- [ ] XML: fluid height 185mm, freq = ω₁/(2π) = 1.008 Hz
- [ ] h/L = 0.308 → near-critical (h/L ≈ 0.24 전이 기준)
- [ ] 자유표면 변동 > shallow case 대비 증가 확인
- [ ] 압력 시계열에서 고차 하모닉 (2f, 3f) 존재 확인 가능
- [ ] 솔버 완주 (10초)

---

## Scenario 4: Liu 2024 — Large Tank Shallow Pitch Resonance (대형 탱크)

**데이터셋**: Liu et al. (2024), 1000x500x1000mm, 20%H fill, pitch resonance
**난이도**: 상 (대형 탱크 + 공진)

```
유저: 1미터 정사각 탱크 (1000x500x1000mm)에 물 200mm 채워서
      피치 공진 해석해줘. 공진주파수가 0.66Hz이고 진폭 2도야.
      회전축은 탱크 바닥 중심.
```

### 검증 기준
- [ ] XML: tank 1000x500x1000mm, fluid 200mm (20%)
- [ ] 여진: `mvrotsinu` freq=0.66, ampl=2°, axis = bottom center
- [ ] 공진 조건 → violent sloshing, hydraulic jump 패턴
- [ ] `simulationdomain` z에 충분한 마진 (회전으로 인한 하단 이탈 방지)
- [ ] 솔버 완주 (최소 30초, ~20주기)

### 핵심 파라미터
| 항목 | 값 |
|------|-----|
| dp | 0.008 m (1m 탱크이므로 dp 크게) |
| 파티클 수 | ~200K |
| TimeMax | 30.0 s |
| 이론 f₁ | 0.659 Hz (실측 0.657 Hz) |

---

## Scenario 5: Liu 2024 — High Fill Variable Amplitude (진폭 파라메트릭)

**데이터셋**: Liu et al. (2024), 70%H fill, 3 amplitudes at resonance
**난이도**: 상 (파라메트릭 스터디)

```
유저: 1미터 탱크에 물 70% 채우고 공진주파수(0.87Hz) 피치 회전인데,
      진폭을 1도, 2도, 3도로 바꿔가며 파라메트릭 해석 해줘.
      각각 20초씩.
```

### 검증 기준
- [ ] `parametric_study` 도구 호출: amplitude = [1, 2, 3], 3개 케이스
- [ ] 각 케이스: mvrotsinu freq=0.87, ampl={{amplitude}}
- [ ] GenCase × 3 → Solver × 3 → PartVTK × 3 전체 파이프라인
- [ ] 진폭 증가 → 슬로싱 강도 비례 증가 확인
- [ ] result_store에 3개 결과 저장 + compare 호출

---

## Scenario 6: ISOPE 2012 — LNG Mark III 고충전 천장 충격 (LNG 벤치마크)

**데이터셋**: ISOPE Sloshing Benchmark (2012), GTT Mark III, 946x118x670mm
**난이도**: 상 (복합 운동 + 고충전)

```
유저: ISOPE 벤치마크 LNG 탱크 (946x118x670mm) 슬로싱 해석.
      물 90% 채우고 수평 진동 + 피치 회전 동시 가진.
      수평 f=0.6Hz 진폭 20mm, 피치 f=0.6Hz 진폭 3도.
      천장 압력 충격을 보고 싶어.
```

### 검증 기준
- [ ] XML: tank 946x118x670mm, fluid height ~603mm (90%)
- [ ] 복합 여진: `mvrectsinu` (sway) + `mvrotsinu` (pitch) 동시
- [ ] `simulationdomain` 상부 마진 (z + 80% 이상, roof impact 대비)
- [ ] 천장 측정점 포함 (z ≈ 0.668m)
- [ ] 솔버 완주, 천장 충격 압력 피크 > 2000 Pa 확인

### 핵심 파라미터
| 항목 | 값 |
|------|-----|
| dp | 0.004 m |
| 파티클 수 | ~500K |
| TimeMax | 15.0 s |
| 예상 GPU 시간 | ~15분 (RTX 4090) |

---

## Scenario 7: NASA 2023 — Cylindrical Tank Ring Baffle (원통형 + 바플)

**데이터셋**: NASA Anti-Slosh Baffle (2023), D=2.84m, ring baffle
**난이도**: 최상 (원통형 geometry + STL baffle)

```
유저: 직경 2.84m, 높이 3m 원통형 탱크에 물 50% 채워줘.
      링 바플을 수면 아래 0.1R 깊이에 설치하고
      수평 방향 0.5Hz 사인파 10mm 진폭으로 흔들어.
      바플 있을 때와 없을 때 비교해줘.
```

### 검증 기준
- [ ] `geometry` 도구: CylindricalGeometry(R=1.42, H=3.0, fluidH=1.42)
- [ ] 바플 유무 2가지 케이스 비교 (parametric_study 또는 수동 2회)
- [ ] 원통형 자연주파수: ω² = g·λ₁/R·tanh(λ₁h/R), λ₁=1.841
- [ ] 바플 존재 시 자유표면 진폭 감소 확인
- [ ] 측정점: 벽면 + 바플 상하 위치

---

## Scenario 8: English 2021 — mDBC vs DBC 비교 (수치기법 검증)

**데이터셋**: English et al. (2021), DualSPHysics mDBC 검증
**난이도**: 중 (같은 형상, BC 옵션만 변경)

```
유저: SPHERIC 벤치마크 탱크(900x62x508mm) 물 18% 해석인데,
      경계조건을 DBC와 mDBC로 각각 돌려서 비교해줘.
      dp=0.002m로 고해상도로.
      센서1 위치 압력 시계열 비교.
```

### 검증 기준
- [ ] 2개 케이스: `Boundary="DBC"` vs `Boundary="mDBC"` (XML parameter 변경)
- [ ] dp=0.002 → 파티클 수 ~540K (고해상도)
- [ ] mDBC: `<parameter key="BoundaryMethod" value="2" />` (mDBC 옵션)
- [ ] 좌벽 하단 압력 시계열에서 mDBC가 DBC보다 스무스한 결과
- [ ] GPU 시간: dp=0.002 → ~20분 예상

### 핵심 파라미터
| 항목 | 값 |
|------|-----|
| dp | 0.002 m (production quality) |
| 파티클 수 | ~540K |
| TimeMax | 7.0 s |

---

## Scenario 9: Zhao 2024 — Horizontal Cylindrical Tank Pitch (수평 원통형)

**데이터셋**: Zhao et al. (2024), TU Munich, horizontal cylindrical + validation
**난이도**: 최상 (수평 원통형 = STL import 필요)

```
유저: 수평 원통형 탱크 (직경 1m, 길이 3m)에 물 25.5% 채워서
      피치 3도, 주파수 0.55Hz로 슬로싱 해석해줘.
      STL 파일은 datasets/horizontal_cylinder.stl 에 있어.
```

### 검증 기준
- [ ] `stl_import` 도구 호출: STL 파일 로딩 + watertight 검증
- [ ] STL 바운딩 박스 확인: ~(3.0 × 1.0 × 1.0)
- [ ] 여진: `mvrotsinu` freq=0.5496, ampl=3°
- [ ] 유체 높이: 0.255 × D (수평 원통 단면 기준 계산)
- [ ] 솔버 완주, 자유표면 높이 시계열 추출

### 비고
- STL 파일은 사전 생성 필요 (FreeCAD/OpenSCAD로 수평 원통 생성)
- DualSPHysics `drawfilestl` 사용하여 경계 입력

---

## Scenario 10: Frosina 2018 — Automotive Fuel Tank (자동차 연료탱크)

**데이터셋**: Frosina et al. (2018), FCA (Fiat), 실제 자동차 연료탱크
**난이도**: 최상 (복잡 3D CAD 형상 + STL)

```
유저: 자동차 연료탱크 STL 파일(datasets/fuel_tank.stl)로 슬로싱 해석해줘.
      용량 60리터, 물 절반 채우고, 급제동 시나리오:
      x축 방향으로 0.5g (4.9 m/s²) 감속 1초 + 자유 진동 4초.
```

### 검증 기준
- [ ] `stl_import` 도구: watertight 검증 (복잡 형상)
- [ ] 급제동 여진: `mvrectfile` 또는 `mvrectunif` 가속도 프로파일
- [ ] 초기 0-1초: 감속 가속도 → 1-5초: 자유 진동 (댐핑)
- [ ] XML: gravity + external acceleration 조합
- [ ] 유체 중심(CoG) 이동 추적 가능
- [ ] 솔버 완주 (5초)

### 비고
- 실제 FCA 탱크 CAD는 비공개 → 유사 형상 STL 생성 필요
- 복잡 형상: dp=0.005m 이상으로 파티클 수 제한

---

## 시나리오 실행 우선순위

| 순위 | 시나리오 | 데이터셋 | 난이도 | 사전 준비 |
|:---:|:---:|:---:|:---:|:---:|
| 1 | #1 Oil Low Fill | SPHERIC | 중 | XML 수정만 |
| 2 | #2 Shallow Sway | Chen 2018 | 중 | 새 XML 작성 |
| 3 | #3 Near-Critical | Chen 2018 | 상 | #2 변형 |
| 4 | #8 mDBC vs DBC | English 2021 | 중 | SPHERIC 변형 |
| 5 | #4 Large Pitch | Liu 2024 | 상 | 새 XML |
| 6 | #6 LNG Roof Impact | ISOPE 2012 | 상 | 새 XML + 복합운동 |
| 7 | #5 Amplitude Param | Liu 2024 | 상 | parametric_study |
| 8 | #7 Cylindrical Baffle | NASA 2023 | 최상 | geometry + baffle |
| 9 | #9 Horizontal Cyl | Zhao 2024 | 최상 | STL 생성 필요 |
| 10 | #10 Fuel Tank | Frosina 2018 | 최상 | STL 생성 필요 |

---

## 데이터셋 ↔ 시나리오 매핑

| 데이터셋 | 시나리오 수 | 번호 |
|---------|:---:|:---:|
| SPHERIC Test 10 (2011) | 1 | #1 |
| Chen & Xue (2018) | 2 | #2, #3 |
| Liu et al. (2024) | 2 | #4, #5 |
| ISOPE Benchmark (2012) | 1 | #6 |
| NASA Baffle (2023) | 1 | #7 |
| English mDBC (2021) | 1 | #8 |
| Zhao LNG FSI (2024) | 1 | #9 |
| Frosina Automotive (2018) | 1 | #10 |
