# EXP-B Error Pattern Analysis: B0 vs B2

## Parameter-Level Comparison (S01, S04)

### S01: SPHERIC benchmark (pitch rotation, oil)

| Parameter | Ground Truth | B0 (Prompt+Tools) | B2 (Prompt only) | Who's better? |
|-----------|-------------|-------------------|-------------------|:---:|
| tank_x | 0.9 | ✅ 0.9 | ✅ 0.9 | = |
| tank_y | 0.062 | ✅ 0.062 | ✅ 0.062 | = |
| tank_z | 0.508 | ✅ 0.508 | ❌ 0.091 | B0 |
| fill_height | 0.091 | ✅ 0.091 | ❌ 0.508 | B0 |
| motion_type | mvrotsinu | ❌ mvrectsinu | ✅ mvrotsinu | **B2** |
| frequency | 0.6515 | ✅ 0.6515 | ✅ 0.6515 | = |
| amplitude | 4° | ❌ 0.045 | ❌ 0.0698 | = (both fail) |
| timemax | 7.68 | ✅ 7.675 | ✅ 7.675 | = |

### S04: Pitch resonance (water)

| Parameter | Ground Truth | B0 (Prompt+Tools) | B2 (Prompt only) | Who's better? |
|-----------|-------------|-------------------|-------------------|:---:|
| tank_x | 1.0 | ✅ 1.0 | ✅ 1.0 | = |
| tank_y | 0.5 | ✅ 0.5 | ✅ 0.5 | = |
| tank_z | 1.0 | ✅ 1.0 | ❌ 0.2 | B0 |
| fill_height | 0.2 | ✅ 0.2 | ❌ 1.0 | B0 |
| motion_type | mvrotsinu | ❌ mvrectsinu | ✅ mvrotsinu | **B2** |
| frequency | 0.66 | ✅ 0.66 | ✅ 0.66 | = |
| amplitude | 2° | ❌ 0.035 | ❌ 0.035 | = (both fail) |
| timemax | 30 | ✅ 30.0 | ✅ 30.0 | = |

## Error Analysis

### 1. B0 Advantage: Structured Input Prevents Dimension Swaps
- `xml_generator` 도구는 `tank_height`와 `fluid_height`를 별도 파라미터로 받음
- LLM이 올바른 값을 파라미터에 매핑 → geometry 정확
- B2에서는 LLM이 XML을 직접 생성하면서 `<drawbox>` size의 z값과 fluid box z값을 혼동

### 2. B2 Advantage: Domain Prompt Teaches Motion Semantics
- 도메인 프롬프트에 "피치 회전 → mvrotsinu" 지시가 있음
- LLM이 프롬프트 텍스트에서 직접 motion type을 결정 → 올바른 mvrotsinu
- B0의 `xml_generator` 도구에는 pitch→mvrotsinu 매핑 로직에 버그 (P1 패턴)
- 도구가 모든 motion을 mvrectsinu로 생성하는 경향

### 3. Common Failures Across Both
- **Amplitude**: 양쪽 모두 degrees를 radians로 변환하거나 다른 단위로 출력
  - B0: 0.045 (maybe sin(4°)×tank_dim?), B2: 0.0698 (≈π×4/180)
  - Ground truth expects raw degrees (4) — DualSPHysics XML은 실제로 degrees 사용
- **Cylindrical geometry**: 양쪽 모두 `<drawbox>`로 박스 생성 (drawcylinder 미지원)

## Implications for Paper Discussion

1. **Prompt vs Tool Trade-off**: 프롬프트는 semantic understanding(무엇을 해야 하는지),
   도구는 structural correctness(어떻게 표현하는지)를 담당
2. **Tool Bugs Propagate**: xml_generator의 P1 버그가 모든 B0 pitch 시나리오에 영향
   → 도구 품질이 시스템 전체 정확도의 상한선
3. **Synergy (+13.9%)**: Prompt가 도구에 올바른 의도를 전달하고,
   도구가 구조화된 XML로 변환 → 단순 합(+6.9% + 65.3%)보다 높은 성능
4. **Model Size Irrelevant**: 아키텍처 결정(프롬프트 설계, 도구 구현)이
   모델 스케일링보다 훨씬 중요 — 비용 효율적 배포에 시사점
