# Interview Prep Prompt — slosim-agent 프로젝트 전시

Claude Code에 그대로 붙여넣어 실행하는 1회성 작업 프롬프트. 목표는 이 저장소(`slosim-agent`)의 README를 **면접관이 5분 안에 세 가지 오해를 풀고 나가도록** 재구성하는 것이다.

---

## 1. 미션

이 저장소는 채용 면접 후속 평가 대상이다. 면접에서 지원자가 말로 충분히 전달하지 못한 세 가지를 **코드·커밋·폴더 구조가 증거로 스스로 말하게** 만들어라.

세 가지 오해:

| # | 면접관의 의심 | 이 프롬프트가 반박할 근거 |
|---|--------------|--------------------------|
| 1 | "LLM을 단순 wrapper로 쓴 것 아닌가?" | 코드 상 LLM이 실제로 *판단하는* 3개 지점을 표로 적시 |
| 2 | "면접 직전 급조한 포트폴리오 아닌가?" | 커밋 히스토리·브랜치 수·논문 폴더 수를 git에서 직접 뽑아 표면화 |
| 3 | "슬로싱에만 쓰는 좁은 도구 아닌가?" | Roadmap 말미에 도메인 확장 경로를 추상적으로 제시 |

---

## 2. 하드 제약 (절대 어기지 말 것)

### 2.1 수치는 전부 git/파일시스템에서 파생시킨다

README나 새 문서에 들어가는 **모든 정량 수치**는 아래 명령으로 직접 확인한 뒤 삽입한다. 추측·반올림·기억 금지.

```bash
git log --oneline | wc -l              # 커밋 수
git log --format='%ai' | awk '{print $1}' | sort -u | wc -l   # 작업일 수
git log --format='%ai' | sort | head -1                        # 첫 커밋일
git log --format='%ai' | sort | tail -1                        # 최근 커밋일
git tag | wc -l                                                # 릴리스 수
git branch -a | grep -v HEAD | wc -l                           # 브랜치 수
ls -d paper* 2>/dev/null | wc -l                               # 논문 폴더 수
ls -d research* 2>/dev/null | wc -l                            # 리서치 디렉토리 수
find internal/llm/tools -name '*.go' -not -name '*_test.go' | wc -l   # 도구 파일 수
find cases -name '*.xml' | wc -l                               # 사전 정의 케이스 수
```

결과를 README에 넣을 때는 실제 값을 그대로 쓴다. 프롬프트 작성자(대화 상대)가 "289 commits, 5 releases" 같은 숫자를 언급하더라도 **그것은 참고용 기억이며 실제 repo 값이 정답**이다. 불일치 시 repo 값을 따른다.

### 2.2 면접 녹취에서 나온 수치만 예외적으로 인용 가능

아래 두 수치만 "면접 중 지원자가 본인 입으로 언급한 정량 값"으로 간주하고, 인용 시 출처를 `(as described in interview)` 같은 완곡한 형태로 남긴다. 이 값들이 코드/문서에 근거로 남아있지 않다면 `<!-- TODO: locate source in paper/ or research/ -->` 주석과 함께 인용한다.

- `alpha = 980` (수치 확산 계수 관련)
- `13% 오차` (특정 검증 케이스의 상대 오차)

이 외의 물리 수치(주파수, 충전률, 압력 등)는 **반드시 cases/·paper/·research-v*/ 내 실제 파일**에서 찾아 인용하고 파일 경로를 각주 형식으로 남긴다. 찾지 못하면 해당 문단에 `<!-- TODO -->` 표기 후 본인(저장소 오너)이 채우게 남긴다.

### 2.3 CLAUDE.md 는 README 전면에 노출하지 않는다

일부 면접관은 AI 보조 개발을 "전부 AI에 시킨" 것으로 오해한다. CLAUDE.md·`.claude/agents/`·`.claude/commands/`는 README 본문에서 언급하지 않거나, 언급하더라도 "개발 생산성을 위한 내부 문서화" 수준으로 한 줄만 처리한다. 전시 공간의 주인공은 `internal/llm/`, `internal/llm/tools/`, `cases/`, `paper*/`, `research*/`다.

### 2.4 LG·특정 기업명 직접 언급 금지

Roadmap 말미 도메인 확장 항목에서 "화학 프로세스 시뮬레이션", "다상 유동 전반", "산업용 tank/reactor geometry 일반화" 같은 추상적 표현을 사용한다. 특정 랩·회사·인명을 쓰지 않는다. 지원자스러운 인상을 주지 않기 위한 의도된 제약이다.

### 2.5 paper 폴더 3종의 정체는 추측 금지

`paper/`, `paper-cs/`, `paper-pof/` 각각의 투고 대상 저널과 상태는 저장소 오너만 안다. README에는 아래 형태로만 표기한다:

```md
- `paper-cs/` — CS 저널 타깃 원고 <!-- TODO: 저널명/상태 채우기 -->
- `paper-pof/` — 물리 검증 중심 원고 <!-- TODO: 저널명/상태 채우기 -->
- `paper/` — 선행 통합 원고 (legacy) <!-- TODO: 역할 확정 -->
```

CLAUDE.md의 "Git Branch Strategy" 섹션에 힌트가 있으나, 그것은 브랜치에 대한 설명이지 폴더 자체의 현재 상태를 보장하지 않는다. 작성 시 혼동하지 말 것.

---

## 3. 사전 읽기 (수정 전 반드시 확인)

순서대로 읽고 내용을 내재화한다.

1. `README.md` — 현재 빈약한 상태 파악
2. `CLAUDE.md` — 프로젝트 구조·도구 목록·브랜치 전략
3. `ARCHITECTURE.md` — 팀·워크플로우 설계
4. `PRD.md` — 제품 요구사항
5. `AGENTS.md` — 에이전트 시스템 개요
6. `internal/llm/prompt/sloshing_coder.go` — 도메인 특화 프롬프트 (LLM의 "판단" 지점 근거)
7. `internal/llm/tools/` 전 파일 — 도구 구현 (LLM이 무엇을 결정하고 무엇을 위임하는지)
8. `internal/llm/agent/` — agent loop, reflection 로직 위치
9. `cases/*.xml` 파일 목록 — 사전 정의 시나리오 폭
10. `paper*/README.md`, `paper*/outline.md` — 논문 원고 존재 증거
11. `research*/README.md`, `research*/experiment_registry.json` 유사 파일 — 실험 이력 증거
12. `docs/FRD_v0.1.md` — Feature ID 체계

---

## 4. 출력 대상

### 4.1 주 타깃: `README.md` 전면 교체

아래 구조로 작성한다. 각 섹션의 목적·금지사항을 엄수한다.

```md
# slosim-agent

<1~2줄 엘리베이터 피치. "자연어 → DualSPHysics GPU 시뮬레이션 → 물리 해석"
 파이프라인을 하나의 에이전트로 통합한 실험 플랫폼.>

## What this actually does

<5~7줄. 동작 예시 한 줄(자연어 입력 → 결과) + 스크린샷/asciicast 자리 표시>
<!-- TODO: demo asciinema or screenshot -->

## Why the LLM is not a wrapper

<섹션 3.1의 표 삽입. 이것이 README의 핵심 섹션이다.>

## Evidence of sustained work

<섹션 3.2의 증거 블록. git-derived 수치 삽입.>

## Architecture (one screen)

<CLAUDE.md "Core Layers" 를 5~8줄로 압축. 파일 경로 링크 포함.>

## DualSPHysics tool surface

<internal/llm/tools/ 의 15개 DSPH 도구를 표로. 각 행은 Tool / Feature ID / 1줄 책임. CLAUDE.md의 표를 그대로 쓰지 말고 직접 파일에서 파싱해 재구성.>

## Cases & papers

<cases/ 디렉토리의 .xml 개수와 대표 4~5개. paper*/ 3종 폴더의 역할 (섹션 2.5 형식).>

## Research history

<research/, research-v2/, research-v3/ 의 실험 그룹을 1~2줄씩. 실제 존재하는 exp 디렉토리명으로 교차검증.>

## Roadmap

- <현재 진행>
- <단기>
- <장기: 도메인 확장 — 2.4 제약 준수>

## Build & run

<CLAUDE.md Build & Development Commands 에서 "핵심 5개"만 발췌.>

## License / Acknowledgements

<OpenCode fork 명시. DualSPHysics, Qwen3 credit.>
```

### 4.2 부 타깃 (필요 시)

`docs/INTERVIEW_NOTES.md` 는 **만들지 않는다.** 면접관용 가공 문서를 별도로 두면 "이 저장소 자체가 면접용 연출"이라는 인상을 준다. 모든 것은 README에 자연스럽게 녹여야 한다.

---

## 5. 섹션 3.1 — "Why the LLM is not a wrapper" 표 작성 규칙

아래 3개 영역별로 **코드 파일·함수명·줄 번호**를 인용한다. 추상적 주장 금지. 표 예시 골격:

| Decision surface | What the LLM decides | Evidence |
|------------------|----------------------|----------|
| **Parameter inference** | 자연어 입력에서 `dp`, `time_max`, `fill_height`, `amplitude` 등 SPH 실행 파라미터를 도메인 규칙에 따라 선택 | `internal/llm/prompt/sloshing_coder.go:L<실제 라인>` — parameter inference rules; `internal/llm/tools/xml_generator.go` — 결과 주입 지점 |
| **Reflection on residuals** | 시뮬레이션 산출 메트릭(예: 압력 시계열, 에너지 보존 잔차)을 읽고 재실행 여부·후속 도구 선택을 판단 | `internal/llm/agent/<파일>` + `internal/llm/tools/analysis.go`, `monitor.go`, `error_recovery.go` |
| **Failure-aware planning** | GenCase/Solver 에러 로그를 해석해 복구 전략 생성 (재파라미터화, 케이스 수정, 사용자 질의) | `internal/llm/tools/error_recovery.go` + 관련 integration 테스트 |

**작성 절차:**

1. 세 영역 각각에 대해 실제 파일을 `Read`·`Grep`으로 열어 LLM의 판단이 들어가는 호출 지점을 찾는다.
2. 찾은 지점의 파일 경로와 **함수명 또는 줄 번호**를 `file:line` 형식으로 인용한다.
3. 증거가 예상보다 약한 영역이 있다면 과장하지 말고 "(scaffolding in place — recovery policy WIP)" 같은 정직한 주석을 남긴다. 허위 주장이 발각되면 1번 오해를 강화한다.

---

## 6. 섹션 3.2 — "Evidence of sustained work" 증거 블록 작성 규칙

git에서 직접 뽑은 수치로 아래 bullet을 채운다.

- Commits: `<N>` across `<작업일 수>` active days (`<첫 커밋일>` → `<최근 커밋일>`)
- Branches: `main` + research/paper forks (실제 브랜치명 나열; CLAUDE.md 브랜치 전략 섹션 참고하되 실제 `git branch -a`가 정답)
- Pre-defined cases: `<N>` XML scenarios in `cases/` (대표 4~5개 파일명)
- DSPH tool implementations: `<N>` Go files in `internal/llm/tools/` (15개 DSPH 도구 외 11개 built-in)
- Paper drafts: `<N>` target venues (`paper/`, `paper-cs/`, `paper-pof/`) — 2.5 형식
- Experiment trails: `<N>` versioned research directories (`research/`, `research-v2/`, `research-v3/`)

**태그/릴리스가 0개라면** "5 releases" 같은 허위 표기 금지. 대신 "continuous development; no release cut yet" 같은 솔직한 표현을 쓴다. 태그가 있을 경우에만 릴리스 수를 표기한다.

수치 근거를 남기기 위해 README 하단 `<!-- generated: YYYY-MM-DD via scripts/repo_stats.sh -->` 같은 주석을 남기는 것은 허용. 스크립트를 실제로 만들지는 않는다 (생성형 허위 증거 금지).

---

## 7. Roadmap 도메인 확장 항목 작성 규칙

Roadmap 마지막 bullet 예시:

> **Domain generalization.** The agent pipeline (NL → case generator → GPU solver → analysis) is not tied to sloshing physics. The same decomposition applies to any tank/reactor geometry problem where (a) a parametric case generator exists, (b) a GPU-accelerated solver is available, and (c) residual metrics can be post-processed. Near-term candidates: multi-phase flows, chemical process tank dynamics, general CFD pre/post workflows.

- "화학", "reactor", "process" 같은 추상적 단어까지는 허용.
- 특정 기업·랩·인명·제품명 금지.
- "어디에 지원 중" 같은 맥락 금지.

---

## 8. 톤·스타일

- 영어 README 기본. 한국어 섹션을 별도로 두지 않는다.
- 이모지 사용 금지.
- "state-of-the-art", "cutting-edge", "blazing fast" 같은 마케팅 수식어 금지.
- 1인칭·복수형 혼용 금지. "The agent does X." 형식의 건조한 서술.
- 모든 파일 참조는 `` `path/to/file.go` `` 백틱 + 필요 시 `:line` 형식.
- Mermaid 다이어그램은 아키텍처 1개만 허용 (한 화면에 들어가는 수준).

---

## 9. 작업 종료 전 자체 평가 (의무)

README 수정이 끝난 직후 본인에게 아래 3문항을 묻고, 각 문항에 대해 README 내 **구체적 섹션/줄**로 답할 수 있는지 확인한다. 답할 수 없다면 README를 고쳐야 한다.

1. 면접관이 "이 프로젝트에서 LLM이 실제로 뭘 판단하느냐?"라고 물으면 README의 어느 표를 가리킬 것인가?
2. "언제부터 만든 거냐? 면접 보고 급하게 만든 건가?"에 어느 섹션으로 답할 것인가?
3. "슬로싱 말고 다른 도메인에도 쓸 수 있나?"에 어느 bullet으로 답할 것인가?

세 질문 모두 "README의 X 섹션"으로 답할 수 있어야 작업 완료다. 그 후 다음을 보고한다:

- 실제 교체한 파일 목록
- 3개 자체 평가 문항의 답 (섹션명 + 첫 줄 발췌)
- 남겨둔 `<!-- TODO -->` 주석 위치 전수 (저장소 오너가 채울 수 있도록)
- 수치가 프롬프트 작성자의 사전 언급(예: "289 commits")과 달랐다면, 실제 값과 그 차이

---

## 10. 하지 말 것 (요약)

- 수치 지어내기 (특히 커밋·릴리스·케이스 수)
- CLAUDE.md를 README 본문에 강조
- `.claude/agents/` 를 셀링 포인트로 내세우기
- 특정 기업·랩·인명 언급
- paper 폴더 3종의 저널명·상태 추측
- 면접용 별도 문서(`INTERVIEW_NOTES.md` 등) 신설
- 이모지, 마케팅 수식어, "AI-powered" 남발
- LLM의 판단 지점을 실제 코드로 증명하지 않고 추상적으로만 기술
- 커밋 / 푸시 (이 프롬프트는 README 편집까지만 지시한다. 커밋은 저장소 오너가 직접 검토 후 수행)

---

## 11. 참고: 이 프롬프트 자체에 대한 메타 규칙

이 프롬프트 파일(`prompts/interview-prep.md`)을 실행 중 수정하지 말 것. 실행 결과로 드러난 문제점은 별도 이슈로만 남긴다.
