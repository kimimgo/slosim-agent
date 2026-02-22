# Table 3: EXP-2 NL→XML Generation Results

| Level | Name | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |
|-------|------|---|-----------|---------------|---------------|-------------|
| L1 | Basic | 4 | 4/4 | 96% | 4/4 | 8.3s |
| L2 | Domain | 4 | 3/4 | 42% | 3/4 | 13.0s |
| L3 | Paper | 4 | 3/4 | 15% | 3/4 | 13.1s |
| L4 | Complex | 4 | 3/4 | 50% | 3/4 | 12.6s |
| L5 | Edge | 4 | 2/4 | 25% | 1/4 | 15.7s |
| **Total** | | **20** | **15/20** | **46%** | **14/20** | **12.6s** |

## Detailed Results

| ID | Level | Tool | Accuracy | Valid | Key Params | Latency |
|----|-------|------|----------|-------|------------|--------|
| S01 | L1 | xml_generator | 100% | Y | L=0.9, h=0.093, f=0.613 | 7.3s |
| S02 | L1 | xml_generator | 100% | Y | L=0.6, h=0.15, f=1.2 | 12.3s |
| S03 | L1 | xml_generator | 86% | Y | L=1, h=0.3, f=0.433 | 9.1s |
| S04 | L1 | xml_generator | 100% | Y | L=0.6, h=0.2, f=0.5 | 4.6s |
| S05 | L2 | none | 0% | N |  | 16.8s |
| S06 | L2 | xml_generator | 25% | Y | L=1, h=0.3, f=0.77 | 9.9s |
| S07 | L2 | xml_generator | 43% | Y | L=40, h=5, f=0.3 | 16.5s |
| S08 | L2 | xml_generator | 100% | Y | L=1, h=0.42, f=1 | 9.0s |
| S09 | L3 | xml_generator | 17% | Y | L=40, h=10.8, f=0.374 | 10.2s |
| S10 | L3 | none | 0% | N |  | 16.7s |
| S11 | L3 | xml_generator | 29% | Y | L=0.6, h=0.2, f=1 | 10.4s |
| S12 | L3 | xml_generator | 17% | Y | L=40, h=13.5, f=0.127 | 14.9s |
| S13 | L4 | xml_generator | 100% | Y | L=40, h=13.5, f=0.221 | 10.5s |
| S14 | L4 | none | 0% | N |  | 16.9s |
| S15 | L4 | xml_generator | 100% | Y | L=40, h=13.5, f=0.5 | 15.8s |
| S16 | L4 | xml_generator | 0% | Y | L=40, h=13.5, f=1 | 7.5s |
| S17 | L5 | xml_generator | 0% | N | L=40, h=13.5, f=0.126 | 16.1s |
| S18 | L5 | none | 100% | Y |  | 16.7s |
| S19 | L5 | none | 0% | N |  | 16.8s |
| S20 | L5 | xml_generator | 0% | N | L=10, h=5, f=0.267 | 13.4s |

## Model: qwen3:latest
## Date: 2026-02-19 08:32
