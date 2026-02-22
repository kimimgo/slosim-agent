# Table 3: EXP-2 NL→XML Generation Results

| Level | Name | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |
|-------|------|---|-----------|---------------|---------------|-------------|
| L1 | Basic | 4 | 4/4 | 96% | 4/4 | 397.0s |
| L2 | Domain | 4 | 4/4 | 67% | 4/4 | 730.4s |
| L3 | Paper | 4 | 3/4 | 15% | 3/4 | 518.1s |
| L4 | Complex | 4 | 3/4 | 58% | 3/4 | 447.5s |
| L5 | Edge | 4 | 3/4 | 0% | 0/4 | 441.5s |
| **Total** | | **20** | **17/20** | **47%** | **14/20** | **506.9s** |

## Detailed Results

| ID | Level | Tool | Accuracy | Valid | Key Params | Latency |
|----|-------|------|----------|-------|------------|--------|
| S01 | L1 | xml_generator | 100% | Y | L=0.9, h=0.093, f=0.613 | 427.6s |
| S02 | L1 | xml_generator | 100% | Y | L=0.6, h=0.15, f=1.2 | 345.3s |
| S03 | L1 | xml_generator | 86% | Y | L=1, h=0.3, f=0.77 | 326.4s |
| S04 | L1 | xml_generator | 100% | Y | L=0.6, h=0.2, f=0.5 | 488.7s |
| S05 | L2 | xml_generator | 86% | Y | L=40, h=8.1, f=0.104 | 557.7s |
| S06 | L2 | xml_generator | 25% | Y | L=1, h=0.3, f=0.75 | 782.5s |
| S07 | L2 | xml_generator | 57% | Y | L=40, h=10, f=0.3 | 795.4s |
| S08 | L2 | xml_generator | 100% | Y | L=1, h=0.42, f=1 | 786.0s |
| S09 | L3 | xml_generator | 17% | Y | L=40, h=10.8, f=0.1053 | 665.2s |
| S10 | L3 | none | 0% | N |  | 826.6s |
| S11 | L3 | xml_generator | 29% | Y | L=0.6, h=0.2, f=1.02 | 301.9s |
| S12 | L3 | xml_generator | 17% | Y | L=1, h=0.3, f=0.756 | 278.8s |
| S13 | L4 | xml_generator | 100% | Y | L=40, h=13.5, f=0.14 | 395.9s |
| S14 | L4 | xml_generator | 33% | Y | L=40, h=13.5, f=0.125 | 416.9s |
| S15 | L4 | xml_generator | 100% | Y | L=1, h=0.3, f=0.5 | 254.5s |
| S16 | L4 | none | 0% | N |  | 722.7s |
| S17 | L5 | xml_generator | 0% | N | L=40, h=13.5, f=0.125 | 332.8s |
| S18 | L5 | xml_generator | 0% | N | L=1, h=0.3, f=1 | 296.8s |
| S19 | L5 | xml_generator | 0% | N | L=40, h=27, f=0.14 | 398.4s |
| S20 | L5 | none | 0% | N |  | 737.8s |

## Model: qwen3:32b-64k
## Date: 2026-02-19 11:30
