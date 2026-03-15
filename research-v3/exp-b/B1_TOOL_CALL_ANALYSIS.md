# EXP-B B1 Tool Call Analysis

## B1 (−DomainPrompt): 도구 호출 패턴

Generic CoderPrompt + 18개 DualSPHysics 도구 available.
결과: **DualSPHysics 도구 0회 호출**

### 세션별 상세

| Session | Scenario | Model | Tool Calls | 행동 설명 |
|---------|----------|-------|------------|----------|
| 1 | S01 | 32B | `glob` | SPHERIC_Test10*.xml 검색 → 기존 파일 발견 → end_turn |
| 2 | S04 | 32B | `glob`, `glob` | *pitch*resonance*.xml 검색 → 미발견 → end_turn |
| 3 | S07 | 32B | `glob` | cases/*.xml 검색 → 기존 파일 목록 확인 → end_turn |
| 4 | S01 | 8B | `bash`, `glob`, `grep`, `edit` | echo로 설명 출력 → 파일 검색 → oil grep → 편집 시도 → read first 에러 |
| 5 | S04 | 8B | `glob`, `glob` | SloshingTank_Def.xml 검색 → 발견 → end_turn |
| 6 | S07 | 8B | `edit` | SloshingTank_Def.xml 직접 편집 시도 → read first 에러 → end_turn |

### 도구 호출 통계

| 도구 유형 | 도구명 | 호출 횟수 |
|----------|--------|:---:|
| **범용 (Generic)** | glob | 8 |
| | edit | 2 |
| | bash | 1 |
| | grep | 1 |
| **DualSPHysics** | xml_generator | 0 |
| | gencase | 0 |
| | solver | 0 |
| | (전부) | 0 |

### 대조: B0 (Full System) 도구 호출

EXP-A 세션에서 domain prompt 포함 시:
- `xml_generator`, `gencase`, `solver`, `partvtk`, `measuretool`, `analysis`, `job_manager`, `monitor`, `error_recovery`, `stl_import`
- **DualSPHysics 도구만 사용** — generic tools 미사용

### 분석

1. **Generic CoderPrompt**는 "코딩 작업을 도와달라"는 지시만 포함
2. 모델은 슬로싱 시뮬레이션 요청을 "기존 XML 파일 찾아서 수정" 문제로 해석
3. `xml_generator` 도구가 available하지만, 모델이 이를 "슬로싱 시뮬레이션에 필요한 도구"로 인식하지 못함
4. 도메인 프롬프트가 도구의 **사용 컨텍스트**를 제공하는 핵심 역할

### 논문 implications

- Domain prompt ≠ just "prompt engineering" — 도구 discovery & activation의 핵심
- LLM은 도구 목록만으로는 어떤 도구를 어떤 상황에 써야 하는지 모름
- 도메인 프롬프트는 "이 도구들은 이 순서로 이렇게 써라"는 orchestration instruction
