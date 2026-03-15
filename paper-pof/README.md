# Paper 2: Physics of Fluids — Validation

## Title
Physical Validation of AI-Generated SPH Sloshing Simulations: SPHERIC Test 10 Benchmark Study

## Target Journal
Physics of Fluids (Kim & Park 2023 게재 이력)

## Framing
V-model 우측 — "현실 실험과 일치하는가?" (Validation)

## Gap × Evidence

| Gap | Evidence (Validation) |
|-----|----------------------|
| G1 비전문가 접근성 | SPHERIC T10 실험 대비 정량 검증 (실제로 물리적으로 맞는 결과) |
| G2 SPH solver agent | Bridge 5/5 물리 PASS + 배플 28.5% SWL 저감 |
| G3 Local SLM 충분 | 도메인 특성 반영 실험 비교 (SPH 인공점성, 커널, 경계조건) |

## Key Novelty
- **SPHERIC Test 10 Oil**: SPH 기반 최초 정량 에러 보고 (1,176편 서베이 확인)
- "맞추기 어려운" 실험 검증에 집중 (near-resonance, 2유체, violent sloshing)

## Planned Sections
1. Introduction (동일 Gap, validation framing)
2. Related Work (SPH sloshing validation 문헌, SPHERIC T10 history)
3. Agent System Overview (1차 논문 ref, 간략 요약)
4. Benchmark Cases
   - SPHERIC Test 10 (Water + Oil, 2유체)
   - Rafiee 2011 (near-resonance sway)
   - Bridge experiment (다중 물리 기준)
5. Experiments
   - EXP-C: SPHERIC T10 압력 비교 (Oil 26.4%, Water lateral)
   - Rafiee 2011: dp 수렴성, P2 피크 4.8%
   - Bridge: 정수압(3.0%), 역위상(r=-0.72), 주파수(4.4%)
   - EXP-D: 배플 최적화 (28.5% SWL 저감)
   - SPH 도메인 특성 분석 (인공점성, Laminar+SPS, 커널 크기)
6. Discussion (SPH 검증 관행 부재, 정량 메트릭 필요성)
7. Conclusion

## Key Numbers
- SPHERIC Oil: 26.4% avg error (Laminar+SPS, 10ms)
- Rafiee 2011 P2: 4.8% peak error (dp=0.004)
- Bridge: 5/5 physics PASS
- Baffle: 28.5% SWL reduction (좌우 대칭)
- Literature gap: 1,176편 중 SPH Oil SPHERIC T10 정량 보고 = 0편

## TODO
- [ ] 1차 논문 완성 후 착수
- [ ] 추가 벤치마크 확장 (배플 파라메트릭: 형상/위치/높이)
- [ ] SPH 수렴 분석 심화 (dp sweep, 경계조건 비교)
- [ ] 도메인 특성 반영 실험 비교 설계
- [ ] LaTeX 원고 작성
