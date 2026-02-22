# Table 5: EXP-4 Domain Prompt Ablation Results

| Condition | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |
|-----------|---|-----------|---------------|---------------|-------------|
| FULL | 10 | 10/10 | 60% | 8/10 | 361.0s |
| NO-DOMAIN | 10 | 10/10 | 50% | 8/10 | 194.1s |
| NO-RULES | 10 | 10/10 | 57% | 8/10 | 356.9s |
| GENERIC | 10 | 6/10 | 35% | 6/10 | 267.6s |

## Per-Scenario Breakdown

| Scenario | Level | FULL | NO-DOMAIN | NO-RULES | GENERIC |
|----------|-------|------|-----------|----------|--------|
| S01 | L1 | 100% V | 100% V | 100% V | 0% X |
| S03 | L1 | 86% V | 100% V | 86% V | 100% V |
| S05 | L2 | 86% V | 14% V | 86% V | 0% X |
| S07 | L2 | 86% V | 100% V | 57% V | 86% V |
| S09 | L3 | 17% V | 33% V | 17% V | 17% V |
| S11 | L3 | 29% V | 29% V | 29% V | 29% V |
| S13 | L4 | 100% V | 20% V | 100% V | 20% V |
| S15 | L4 | 100% V | 100% V | 100% V | 0% X |
| S17 | L5 | 0% X | 0% X | 0% X | 100% V |
| S19 | L5 | 0% X | 0% X | 0% X | 0% X |

## Model: qwen3:32b-64k
## Date: 2026-02-19 14:48
