# Table 5: EXP-4 Domain Prompt Ablation Results

| Condition | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |
|-----------|---|-----------|---------------|---------------|-------------|
| FULL | 10 | 7/10 | 46% | 6/10 | 12.6s |
| NO-DOMAIN | 10 | 10/10 | 44% | 8/10 | 8.8s |
| NO-RULES | 10 | 9/10 | 55% | 8/10 | 12.2s |
| GENERIC | 10 | 8/10 | 39% | 8/10 | 8.7s |

## Per-Scenario Breakdown

| Scenario | Level | FULL | NO-DOMAIN | NO-RULES | GENERIC |
|----------|-------|------|-----------|----------|--------|
| S01 | L1 | 100% V | 100% V | 100% V | 0% X |
| S03 | L1 | 86% V | 86% V | 86% V | 71% V |
| S05 | L2 | 0% X | 14% V | 86% V | 14% V |
| S07 | L2 | 57% V | 86% V | 29% V | 57% V |
| S09 | L3 | 17% V | 17% V | 17% V | 17% V |
| S11 | L3 | 0% X | 14% V | 29% V | 14% V |
| S13 | L4 | 100% V | 20% V | 100% V | 20% V |
| S15 | L4 | 100% V | 100% V | 100% V | 100% V |
| S17 | L5 | 0% X | 0% X | 0% X | 100% V |
| S19 | L5 | 0% X | 0% X | 0% X | 0% X |

## Model: qwen3:latest
## Date: 2026-02-19 08:40
