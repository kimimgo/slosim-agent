# EXP-D: Autonomous Baffle Optimization

## NL Prompt (Input to slosim-agent)

fuel_tank.stl 파일의 자동차 연료탱크에 대해 슬로싱 시뮬레이션을 설정하세요.
조건: 50% 수위, 급정거 시나리오 (x방향 0.5Hz, 50mm 진폭).
baffle 없이 baseline 시뮬레이션을 먼저 실행하고,
결과를 분석하여 슬로싱을 최소화하는 baffle 위치를 자율적으로 결정하세요.
3회 이내 반복으로 최적 baffle 배치를 찾아주세요.

## Expected Agent Workflow

1. `stl_import` — fuel_tank.stl 로드, BBox 확인
2. `xml_generator` or modified XML — baseline 케이스 생성 (no baffle)
3. `gencase` → `solver` → `measuretool` — baseline 시뮬레이션
4. `analysis` — baseline SWL 측정, 슬로싱 패턴 분석
5. `baffle_generator` — 결과 기반 baffle 배치 결정
6. 반복 2-3회 → 최적 baffle 도출

## Evaluation Metrics (M-D)

| Metric | Definition | Method |
|--------|-----------|--------|
| M-D1 | STL loading success | GenCase no error |
| M-D2 | Baseline simulation complete | Solver normal exit |
| M-D3 | SWL reduction rate | (baseline_SWL - optimized_SWL) / baseline_SWL |
| M-D4 | Autonomous iteration count | Agent self-decided iterations |
| M-D5 | Total wall-clock time | NL input → final result |

## Tank Specifications (Frosina2018)

- Dimensions: ~500mm x 350mm x 250mm (beveled box)
- STL file: cases/fuel_tank.stl
- dp = 0.008 m
- Fill level: 50% (~125mm)
- Excitation: sinusoidal x-direction, 0.5Hz, 50mm amplitude
