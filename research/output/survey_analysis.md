# Survey Analysis

## Collection Statistics
- Total papers: 2175
- Total rounds: 1
- Queries executed: 18
- Citation chains: 0
- Datasets found: 0

## Database Coverage
- arxiv: 18 queries
- dblp: 18 queries
- crossref: 18 queries
- openalex: 18 queries
- semantic_scholar: 18 queries
- core: 18 queries
- serpapi: 18 queries

## Top Keywords (TF-IDF)
- sloshing: 36.5447
- tank: 30.2582
- neural: 27.9854
- simulation: 26.267
- network: 23.3227
- automation: 23.0071
- multi-agent: 22.9926
- llm: 21.8561
- model: 20.8243
- agent: 20.5288
- analysis: 20.5281
- language: 20.5272
- graph: 20.5057
- seismic: 20.2241
- particle: 20.1564

## Temporal Trend
- 1995: 1 papers
- 1998: 1 papers
- 1999: 1 papers
- 2001: 1 papers
- 2003: 1 papers
- 2005: 2 papers
- 2006: 1 papers
- 2007: 2 papers
- 2008: 2 papers
- 2009: 4 papers
- 2010: 3 papers
- 2011: 4 papers
- 2012: 7 papers
- 2013: 3 papers
- 2014: 8 papers
- 2015: 5 papers
- 2016: 12 papers
- 2017: 5 papers
- 2018: 65 papers
- 2019: 78 papers
- 2020: 79 papers
- 2021: 79 papers
- 2022: 168 papers
- 2023: 285 papers
- 2024: 374 papers
- 2025: 681 papers
- 2026: 78 papers

## Venue Distribution
- Unknown: 1027
- arXiv.org: 69
- IEEE Access: 15
- Engineering Applications of Computational Fluid Mechanics: 14
- NTNU: 14
- Applied Sciences: 13
- Lecture Notes in Computer Science: 13
- Studies in Computational Intelligence: 12
- SpringerBriefs in Applied Sciences and Technology: 12
- Ocean Engineering: 12

## Coverage Gaps
- Topic keywords missing from corpus: sph-specific

## Gap Evolution by Round
### Round 1
Summary: 교전 라운드 1 (Engagement Round 1) 완료. 2,175편 수집 (18 queries × 7 databases). 슬로싱 도메인 4대 산업 응용(LNG/자동차/원자력/우주) + LLM+CFD 자동화 경쟁자(AutoCFD/OpenFOAMGPT/FoamGPT) + SPH+LLM 블루오션 확인. 핵심 발견: (1) LLM+CFD 자동화는 전부 mesh-based(OpenFOAM)이며 particle-based SPH용은 전무, (2) Pasimodo+RAG가 가장 가까운 경쟁자이나 RAG 기반 보조 도구에 불과, (3) 슬로싱 도메인의 산업적 중요성(Falcon 1 손실, LNG 모형시험 20,000+hr, 원자력 규정)이 강력한 motivation 제공, (4) DesignSPHysics GUI의 한계점(설치 복잡, 시뮬레이션 불일치, 파라메트릭 불가)이 우리 시스템의 차별점 근거.
- SPH + LLM Agent 조합 연구 전무 (severity: critical)
- Particle-Based Solver용 Tool Interface 설계 패턴 부재 (severity: high)
- 전산역학/슬로싱 도메인 특화 프롬프트 체계적 평가 부재 (severity: high)
- NL→벤치마크 검증 End-to-End 파이프라인 부재 (severity: medium-high)
- 슬로싱 산업 PoC 및 현장 기대효과 검증 부재 (severity: medium-high)
- 로컬/오픈웨이트 LLM 시뮬레이션 적용 (severity: medium)
