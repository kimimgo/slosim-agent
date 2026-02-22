# /research — Tactical Research Planning

Research Sniper v2 대화형 연구 기획. 5단계 전술 파이프라인으로 연구 계획 수립.

## Usage

```
/research              # 새 연구 시작 (SCOPE부터)
/research resume       # 기존 세션 이어서
/research status       # 현재 진행 상황
```

## Tactical Stages

```
0 목표선정(SCOPE)    — 6-phase 연구 인터뷰         표적 식별
1 정찰(RECON)        — Iterative 논문 수집/분석     정보 수집
2 영점조정(ZERO)     — Claude gap reasoning         gap 교정
3 조준(AIM)          — 연구 결론 도출               제원 산출
4 발사(SEND)         — BibTeX + 연구 계획 출력      최종 사격
```

## New Session — SCOPE (목표선정)

1. `get_progress(output_dir="./output")` 호출
2. 세션 없으면 → 6-phase Interview:
   - Phase 1 (Topic): topic, motivation
   - Phase 2 (Scope): hypothesis, research_questions
   - Phase 3 (Capabilities): available_data, compute_resources
   - Phase 4 (Constraints): target_venue, deadline
   - Phase 5 (Related Work): existing_approaches, known_gaps, key_references
   - Phase 6 (Expected Contribution): expected_contributions, novelty_claim
3. `create_profile(...)` 호출

### Gate 1: SCOPE 확인

Profile 요약 표시 → **"이 목표로 정찰(RECON)을 시작할까요?"** → 승인 필수.
**승인 없이 survey_start 호출 NEVER.**

## RECON-ZERO Loop (정찰-영점조정)

4. `survey_start(topic, ...)` → 추천 쿼리 획득
5. `survey_search(session_id, query, ...)` x N회 (병렬 가능)
6. `survey_analyze(session_id, ...)` → `claude_analysis_data` 획득
7. **Gap Reasoning** (Claude가 수행):
   - 수집된 논문 분석
   - 연구 gap 식별 (title, severity, evidence, perspective)

### Gate 2: ZERO 승인

Gap 목록 + severity 표시 → 4지선다:
- **사격 허가(Approve)** → `survey_round_commit(session_id, round_summary, approved_gaps)`
- **영점 수정(Refine)** → gap 수정 후 재제시
- **추가 정찰(More)** → 추가 검색 (citation chase, new queries)
- **최종 발사(Finish)** → 최종 출력으로 이동

**사용자 선택 없이 survey_round_commit 호출 NEVER.**

9. 더 조사 필요 시 → 5번으로 복귀 (교전 라운드 반복)

## SEND (발사)

### Gate 3: SEND 확인

최종 gap 요약 + 출력 파일 목록 표시 → **"발사(SEND) 승인"** → 승인 필수.
**승인 없이 survey_export_plan 호출 NEVER.**

10. `survey_export_plan(session_id, final_gaps)` 호출
11. 4개 파일 생성:
    - `output/research_plan.md` — gap별 experiment TODO
    - `output/paper_skeleton.md` — 논문 뼈대
    - `output/survey_analysis.md` — survey 통계
    - `output/references.bib` — BibTeX

## Dashboard

```
Research Sniper v2 — Aim. Analyze. Publish.
Topic: {topic}
{progress_bar}
교전 라운드: {current_round} | 수집 논문: {total} | 승인 Gap: {approved_count}
```

## Deep Research Mode

복잡한 주제 시 survey-orchestrator 에이전트 스폰:
```
Task(subagent_type="survey-orchestrator", prompt="...")
```
Orchestrator가 worker 에이전트를 병렬 스폰하여 빠르게 수집.
