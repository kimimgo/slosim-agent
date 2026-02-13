---
description: Write a BDD scenario for slosim-agent from a non-expert user perspective
argument-hint: "<feature-name>"
---

Write a BDD scenario file for the feature `$ARGUMENTS`.

## Output

Create: `docs/scenarios/$ARGUMENTS.feature`

## Persona

The user is a **Korean office worker** with:
- Zero CFD/engineering knowledge
- No idea what DualSPHysics, SPH, GenCase, or particle methods are
- Comfortable with natural language Korean input
- Expects clear, jargon-free guidance from the AI

## Format

```gherkin
# docs/scenarios/<feature>.feature
Feature: <한국어 기능명>

  Background:
    Given 사용자가 slosim-agent를 실행한다

  Scenario: <정상 흐름 - 한국어>
    When <사용자 행동>
    Then <기대 결과>
    And <추가 검증>

  Scenario: <에러 케이스 - 한국어>
    When <잘못된 입력/상황>
    Then <에러 메시지는 쉬운 한국어>
    And <해결 방법 안내 포함>

  Scenario: <경계 케이스>
    When <극단적 입력>
    Then <적절한 처리>
```

## Required Scenarios Per Feature

Each feature file must include at minimum:
1. **Happy path** — 정상 흐름
2. **Invalid input** — 잘못된 입력 처리
3. **Simulation failure** — 시뮬레이션 실패 시 사용자 안내
4. **Cancellation** — 사용자가 중간에 취소

## Validation Rules

- All user-facing text must be in plain Korean
- If a technical term appears, it must be followed by (쉬운 설명)
- Error messages must suggest what the user can do next
- Never show raw error codes, stack traces, or internal paths to the user
