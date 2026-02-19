# Model Comparison: Qwen3 8B vs 32B

## EXP-2: NL→XML Generation Accuracy

| Level | Name | 8B Accuracy | 32B Accuracy | 8B Tool Call | 32B Tool Call |
|-------|------|------------|-------------|-------------|--------------|
| L1 | Basic | **96%** | **96%** | 4/4 | 4/4 |
| L2 | Domain | 42% | **67%** | 3/4 | **4/4** |
| L3 | Paper | **15%** | **15%** | 3/4 | 3/4 |
| L4 | Complex | 50% | **58%** | 3/4 | 3/4 |
| L5 | Edge | **25%** | 0% | 2/4 | **3/4** |
| **Total** | | 46% | **47%** | 15/20 | **17/20** |

### Key Observations
- **L1 (Basic)**: 동일한 96% — 명시적 파라미터는 모델 크기 무관
- **L2 (Domain)**: 32B가 25%p 우위 (67% vs 42%) — 도메인 지식 추론 능력 차이
- **L3 (Paper)**: 동일한 15% — 논문 고유 파라미터는 양쪽 모두 학습 부재
- **L5 (Edge)**: 8B가 우위 (25% vs 0%) — 32B가 edge case에서 오히려 과잉 생성
- **Tool Calling**: 32B (17/20) > 8B (15/20) — 큰 모델이 tool calling 더 안정적

## EXP-4: Domain Prompt Ablation

| Condition | 8B Accuracy | 32B Accuracy | 8B Tool Call | 32B Tool Call |
|-----------|------------|-------------|-------------|--------------|
| FULL | 46% | **60%** | 7/10 | **10/10** |
| NO-DOMAIN | 44% | **50%** | **10/10** | **10/10** |
| NO-RULES | **55%** | 57% | 9/10 | **10/10** |
| GENERIC | 39% | **35%** | **8/10** | 6/10 |

### Key Observations
- **FULL**: 32B (60%) >> 8B (46%) — 큰 모델이 전체 프롬프트를 더 잘 활용
- **Tool Call Rate**: 32B FULL=10/10, 8B FULL=7/10 — 32B가 긴 프롬프트에서도 안정적
- **어블레이션 패턴**:
  - 32B: FULL(60%) > NO-RULES(57%) > NO-DOMAIN(50%) > GENERIC(35%) ← **기대한 순서!**
  - 8B: NO-RULES(55%) > FULL(46%) > NO-DOMAIN(44%) > GENERIC(39%) ← 역전 현상
- **8B FULL 역전 원인**: 긴 시스템 프롬프트가 8B 모델의 tool calling을 방해 (7/10 vs 10/10)

## Latency Comparison

| Model | Avg Latency/call | Total Time (60 calls) | GPU Allocation |
|-------|-----------------|----------------------|----------------|
| 8B | **11.0s** | **~11 min** | 5.8GB (100% GPU) |
| 32B | **350.0s** | **~5.8 hours** | 14.2GB (71% GPU) |

- 32B는 8B 대비 **32배 느림** (thinking 토큰 + CPU offload)
- 실용적 배포 시 8B가 크게 유리 (정확도 차이 1%p vs 속도 32배)

## Research Implications

1. **모델 크기 효과**: 32B > 8B는 주로 L2(도메인) 시나리오에서 발현 (+25%p)
2. **프롬프트 길이 한계**: 8B 모델은 긴 시스템 프롬프트에서 tool calling 실패율 증가
3. **어블레이션 일관성**: 32B만 기대한 순서(FULL>NO-RULES>NO-DOMAIN>GENERIC) 유지
4. **Thinking overhead**: Qwen3의 thinking 모드가 실용성을 크게 제한 (94% 토큰 소비)
5. **실용적 trade-off**: 8B + 짧은 프롬프트(NO-RULES)가 최적 성능/비용 비율

## Date: 2026-02-19
