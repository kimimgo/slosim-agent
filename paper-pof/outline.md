# Paper 2 Outline: Physical Validation (PoF)

## Abstract (~250 words)
- Problem: AI 에이전트가 생성한 SPH 시뮬레이션을 신뢰할 수 있는가?
- Gap: SPH 슬로싱 정량 검증 관행 부재 (SPHERIC T10 Oil 정량 보고 0편)
- Solution: SlosimAgent 생성 케이스를 SPHERIC T10 + Rafiee 2011로 검증
- Results: Oil 26.4%, P2 4.8%, Bridge 5/5, Baffle 28.5%
- Significance: 최초의 AI-generated SPH sloshing validation

## 1. Introduction
### 1.1 Background
- AI 에이전트 기반 시뮬레이션의 부상
- "생성된 결과를 신뢰할 수 있는가?"라는 질문
- SPH sloshing validation 문헌 현황

### 1.2 Research Gaps (동일 3 Gaps, Validation 관점)
- G1: 에이전트 생성 결과의 물리적 정확도 미검증
- G2: SPH sloshing agent의 실험 대비 검증 = 0
- G3: 도메인 특성(인공점성, 커널, 경계) 영향 미분석

### 1.3 Contributions
1. SPHERIC T10 Oil의 최초 SPH 정량 에러 보고
2. 에이전트 생성 케이스의 다중 벤치마크 검증
3. SPH 도메인 특성이 AI 생성 결과에 미치는 영향 분석
4. 배플 자동 설계의 물리적 유효성 검증

## 2. Related Work
### 2.1 SPHERIC Test 10 History
- Botia-Vera et al. (2010): 원 실험
- English et al. (2021): mDBC, 정량 메트릭 0개 (263인용)
- Delorme et al. (2009): 초기 SPH 검증
### 2.2 SPH Sloshing Validation Literature
- 1,176편 서베이 → Oil 정량 보고 0편
### 2.3 AI-Generated Simulation Validation

## 3. SlosimAgent Overview (Brief — ref to Paper 1)
- 아키텍처 요약 (1차 논문 인용)
- 본 논문에서 사용한 도구 체인

## 4. Benchmark Cases
### 4.1 SPHERIC Test 10
- 실험 조건: 900mm×508mm×62mm, Water+Oil, Lateral sway
- 기존 SPH 연구와의 차이점
### 4.2 Rafiee 2011
- 1.3m×0.1m×0.9m, near-resonance sway, f≈f_natural
### 4.3 Bridge Experiment
- 5가지 물리 기준 (정수압, 역위상, 주파수, 대칭, 에너지)

## 5. Experiments
### 5.1 SPHERIC T10 Pressure Validation
- Water lateral: 시계열 비교
- Oil lateral: 26.4% avg error (Laminar+SPS, 10ms)
- Peak 분석: Peaks 2-4 = 6.5-16.5%, Peak 1 = 68.1% (물리적 설명)
- Artificial vs Laminar+SPS 비교

### 5.2 Rafiee 2011 Near-Resonance
- dp sweep: dp004=4.8% vs dp002=43.4%
- 위상 지연 현상: SPH 인공점성 → 에너지 축적 지연
- dp004>dp002: 큰 커널의 벽면 압력 감지 이점

### 5.3 Bridge Multi-Physics Validation
- 정수압: 3.0% 오차
- 역위상: r=-0.720
- 주파수: 4.4% 오차
- 5/5 PASS

### 5.4 Baffle Optimization (EXP-D)
- 28.5% SWL 저감 (좌우 대칭)
- 에이전트 자동 생성 배플의 공학적 유효성

### 5.5 SPH Domain Characteristics
- 인공점성 모델 비교 (Artificial vs Laminar+SPS)
- 출력 해상도 영향 (0.1s vs 10ms)
- 경계조건 (DBC vs mDBC)

## 6. Discussion
### 6.1 SPH 정량 검증 관행 부재 문제
### 6.2 AI 생성 시뮬레이션의 신뢰 경계
### 6.3 SPH 고유 현상이 결과에 미치는 영향
### 6.4 Limitations and Future Work

## 7. Conclusion

## References
