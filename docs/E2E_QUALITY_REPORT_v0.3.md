# slosim-agent v0.3 E2E 품질 평가 리포트

**평가일**: 2026-02-15
**평가자**: cc-slosim-1 (Claude Opus 4.6)
**모델**: Qwen3:latest (8B, 6GB VRAM)
**환경**: RTX 4090 (24GB, 6GB available), Ollama, Ubuntu

---

## 1. 평가 개요

### 평가 방법
- Non-interactive 모드(`./slosim-agent -p "..." -q -f json`)로 3개 시나리오 순차 실행
- SQLite DB(`opencode.db`)에서 전체 메시지/도구 호출 추출·분석
- TUI interactive 모드 `script` 명령어로 캡처
- Test code 기반이 아닌 **실제 사용자 관점 직접 평가**

### 평가 시나리오

| # | 프롬프트 | 기대 동작 |
|---|---------|----------|
| 1 | "직사각형 탱크 1m×0.5m×0.6m에 물 50% 채우고 주기 0.8초 진폭 3cm로 5초간 슬로싱 해석해줘" | xml_generator → gencase → solver → partvtk → measuretool → report |
| 2 | "2m×1m×1.5m 탱크에서 충전율 30%, 50%, 70%로 파라메트릭 비교해줘. 가진 주기 1.2초 진폭 5cm." | parametric_study → (케이스별) xml_generator → gencase → solver |
| 3 | "mDBC 경계조건으로 0.5m×0.3m×0.4m 소형 탱크에 LNG 유체 60% 채우고, 주기 0.5초 진폭 1cm로 정밀 해석해줘" | xml_generator(boundary_method=mdbc) → gencase → solver → partvtk → measuretool |

---

## 2. 시나리오별 E2E 평가

### Scenario 1: 기본 슬로싱 해석

**결과: CRITICAL FAIL**

| 항목 | 기대값 | 실제값 | 판정 |
|------|--------|--------|------|
| 첫 번째 도구 | `xml_generator` | `parametric_study` | FAIL |
| 파이프라인 완수 | 6단계 (xml→gen→solver→vtk→measure→report) | 0단계 | FAIL |
| 총 도구 호출 | ~6 | 28 (parametric_study x27 + result_store x1) | FAIL |
| 시뮬레이션 실행 | solver 1회 | 0회 | FAIL |
| 최종 응답 품질 | 해석 결과 요약 | "결과를 확인하고 싶으시면 다음 중..." (191자) | FAIL |
| 루프 탈출 | N/A | 27회 동일 호출 후 result_store로 겨우 탈출 | FAIL |

**상세 분석**:
1. 첫 호출에서 `parametric_study`에 `"tank_dimensions": ["1m x 0.5m x 0.6m"]` (string → float64 type error) 전달
2. 에러 수신 후 `fill_ratio: [0.5, 0.6, 0.7]`로 변경하여 재시도 (사용자가 요청하지 않은 파라메트릭 스터디)
3. 이후 26회 동일한 parametric_study 호출 반복 (무한 루프)
4. xml_generator, gencase, solver는 **한 번도 호출되지 않음**
5. 최종 응답은 "어떤 작업을 수행하시겠습니까?" — 아무것도 수행하지 않았음을 자인

### Scenario 2: 파라메트릭 비교

**결과: FAIL**

| 항목 | 기대값 | 실제값 | 판정 |
|------|--------|--------|------|
| 첫 번째 도구 | `parametric_study` | `parametric_study` | **PASS** |
| 파라미터 파싱 | fill_ratio: [0.3, 0.5, 0.7] | fill_ratio: [0.5, 0.6, 0.7] | FAIL (값 불일치) |
| 케이스 실행 | 3케이스 각각 solver 실행 | 0케이스 실행 | FAIL |
| 총 도구 호출 | ~10-15 | 18 (parametric_study x18) | FAIL |
| 반복 루프 | 없음 | 18회 동일 호출 | FAIL |

**상세 분석**:
1. **유일하게 첫 도구 선택이 맞는 시나리오** (프롬프트에 "파라메트릭 비교" 키워드 포함)
2. 그러나 사용자가 지정한 30/50/70%가 아닌 50/60/70%으로 변환 (프롬프트 추종 실패)
3. parametric_study는 메타데이터(manifest) 생성만 수행 → 실제 시뮬레이션 실행은 별도 도구 필요
4. 동일 parametric_study를 18회 반복 (스터디가 이미 생성되었음을 인식 못함)

### Scenario 3: mDBC 정밀 해석

**결과: CRITICAL FAIL**

| 항목 | 기대값 | 실제값 | 판정 |
|------|--------|--------|------|
| 첫 번째 도구 | `xml_generator` (boundary_method=mdbc) | `parametric_study` | FAIL |
| mDBC 파라미터 반영 | boundary_method="mdbc" | 파라미터 없음 | FAIL |
| LNG 유체 특성 반영 | fluid_height=0.24m (0.4×0.6) | 미반영 | FAIL |
| 총 도구 호출 | ~6 | 59 (parametric_study x59) | FAIL |
| 실행 시간 | ~60초 (solver 포함) | 300초 (timeout) | FAIL |
| 시뮬레이션 완료 | 완료 | 미시작 (timeout) | FAIL |

**상세 분석**:
1. 프롬프트의 "mDBC", "LNG", "정밀 해석" 키워드를 모두 무시
2. parametric_study를 59회 반복하다 5분 timeout으로 강제 종료
3. 가장 심각한 무한 루프 케이스 (exit code 124)

---

## 3. 근본 원인 분석

### 3.1 Qwen3:8B 도구 선택 실패 (P0 - Critical)

**증상**: 시스템 프롬프트가 "반드시 xml_generator를 첫 번째로 호출하세요"라고 명시했으나 3/3 시나리오에서 무시됨

**원인**:
- Qwen3:8B는 18개 도구 중 `parametric_study`를 "가장 관련 있는 도구"로 인식
- 도구 이름(`parametric_study`)이 시뮬레이션 설정 요청과 표면적으로 매칭됨
- 시스템 프롬프트의 "절대 규칙"과 "Tool 호출 순서" 지시를 8B 모델이 추종하지 못함
- 도구 description이 "여러 파라미터 조합을 자동으로 실행합니다"로 범용적이라 8B 모델이 혼동

**영향**: 18개 도구를 구현했으나 올바른 도구가 **단 한 번도** 선택되지 않음

### 3.2 에이전트 루프 보호 미비 (P0 - Critical)

**증상**: 동일 도구를 27~59회 반복 호출, 무한 루프

**원인**: `internal/llm/agent/agent.go:276-310` — `for {}` 무한 루프
```go
for {
    agentMessage, toolResults, err := a.streamAndHandleEvents(...)
    if agentMessage.FinishReason() == message.FinishReasonToolUse {
        msgHistory = append(msgHistory, agentMessage, *toolResults)
        continue  // ← 탈출 조건 없음
    }
    return AgentEvent{Done: true}
}
```

누락된 보호 장치:
1. **최대 반복 횟수** (max_iterations) — 없음
2. **중복 도구 호출 감지** (same tool+params → abort) — 없음
3. **에스컬레이션 로직** (N회 실패 → 사용자에게 위임) — 없음

### 3.3 프롬프트-모델 크기 불일치 (P1 - High)

**증상**: 시스템 프롬프트(87줄)의 복잡한 지시사항을 8B 모델이 추종 불가

**원인**:
- 프롬프트에 "절대 규칙", "순서", "자동 결정 규칙", "도메인 지식" 등 다층 지시가 포함
- 8B 모델의 instruction following 능력 한계
- 특히 "도구 순서" 같은 절차적 지시는 70B+ 모델에서도 어려운 과제

### 3.4 도구 결과 해석 실패 (P1 - High)

**증상**: parametric_study가 "스터디 시작... 3개 케이스 생성됨"을 반환했으나, 모델이 "완료됨"으로 인식하지 못함

**원인**: 도구 결과 메시지에 "다음 단계: solver 실행 필요" 같은 명시적 가이드가 없음

---

## 4. TUI UI/UX 평가

### 4.1 현재 TUI 상태 (Interactive 모드 캡처)

```
┌──────────────────────────────────────────┐
│ ⌬ OpenCode unknown                       │
│ https://github.com/opencode-ai/opencode  │
│                                          │
│ cwd: /home/imgyu/.../slosim-agent        │
│                                          │
│ • gopls ()                               │
│                                          │
│ press enter to send, write \ for newline │
│ ────                                     │
│ >                                        │
│ [ctrl+? help]  No diagnostics  [Qwen3]   │
└──────────────────────────────────────────┘
```

### 4.2 UI 디자인 개선점

| # | 현재 상태 | 개선안 | 우선순위 |
|---|----------|--------|---------|
| UI-1 | "OpenCode unknown" 브랜딩 | "slosim-agent v0.3" 자체 브랜딩 필요 | HIGH |
| UI-2 | opencode-ai/opencode URL 표시 | 프로젝트 URL 또는 슬로싱 관련 배너로 교체 | HIGH |
| UI-3 | Sim Dashboard 미연결 | `dashboard.go` 존재하나 tui.go에 라우팅 없음 | HIGH |
| UI-4 | Results Viewer 미연결 | `viewer.go` 존재하나 tui.go에 라우팅 없음 | HIGH |
| UI-5 | Error Panel 미연결 | `panel.go` SemanticTokens 적용 완료됐으나 통합 안됨 | MEDIUM |
| UI-6 | Case Wizard 미연결 | `wizard.go` + `form.go` 존재하나 접근 불가 | MEDIUM |
| UI-7 | "No diagnostics" 텍스트 | 슬로싱 도메인에선 "시뮬레이션 대기 중"이 적절 | LOW |
| UI-8 | 도구 실행 진행률 미표시 | 28회 반복 호출 중 사용자에게 아무 피드백 없음 | HIGH |

### 4.3 UX 개선점

| # | 현재 상태 | 개선안 | 우선순위 |
|---|----------|--------|---------|
| UX-1 | 무한 도구 루프 피드백 없음 | 반복 호출 감지 시 spinner/경고 표시 | P0 |
| UX-2 | 시뮬레이션 진행 상태 미표시 | Sim Dashboard를 도구 실행 시 자동 활성화 | P0 |
| UX-3 | 에러 발생 시 기술적 메시지만 표시 | 사용자 친화적 에러 메시지 + 재시도 버튼 | P1 |
| UX-4 | 대화 모드와 시뮬레이션 모드 구분 없음 | 탭/뷰 전환: Chat ↔ Simulation ↔ Results | P1 |
| UX-5 | 도구 호출 과정 완전 불투명 | 도구 이름 + 상태를 실시간 로그 스트림으로 표시 | P1 |
| UX-6 | 파이프라인 단계 시각화 없음 | `xml → gen → solver → vtk → measure` 스텝퍼 | P2 |
| UX-7 | 결과 파일 접근 경로 안내 없음 | 완료 후 output_dir 자동 안내 + 파일 목록 | P2 |
| UX-8 | 시뮬레이션 취소 불가 | Ctrl+C로 현재 파이프라인만 취소하는 기능 | P2 |

---

## 5. v0.3 컴포넌트 통합 현황

| 컴포넌트 | 코드 존재 | SemanticTokens | tui.go 라우팅 | 키바인딩 | 실사용 가능 |
|----------|:---------:|:--------------:|:------------:|:--------:|:----------:|
| Chat (기본) | O | O | O | — | O |
| Log Viewer | O | O | O | ctrl+l | O |
| Sim Dashboard | O | O | **X** | — | **X** |
| Results Viewer | O | O | **X** | — | **X** |
| Error Panel | O | O | **X** | — | **X** |
| Case Wizard | O | O | **X** | — | **X** |
| Parametric View | O | O | **X** | — | **X** |

→ **7개 중 2개만 실사용 가능** (28.6%)

---

## 6. 종합 평점

### 기능 동작 (Functional Quality)

| 평가 항목 | 점수 | 근거 |
|----------|:----:|------|
| 도구 선택 정확도 | 1/10 | 3시나리오 중 올바른 첫 도구 선택 0회 (시나리오2는 우연 일치) |
| 파이프라인 완수율 | 0/10 | 0/3 시나리오에서 시뮬레이션 미실행 |
| 프롬프트 추종도 | 1/10 | 시스템 프롬프트 "절대 규칙" 완전 무시 |
| 에러 핸들링 | 2/10 | 에러 후 파라미터 수정 시도는 있으나 올바른 도구로 전환 못함 |
| 루프 안전성 | 0/10 | 무한 루프 보호 장치 완전 부재 |
| 응답 품질 | 1/10 | 최종 응답이 "무엇을 하시겠습니까?" 수준 |
| **소계** | **5/60** | **8.3%** |

### UI/UX 품질 (Visual/Interaction Quality)

| 평가 항목 | 점수 | 근거 |
|----------|:----:|------|
| 기본 대화 UI | 7/10 | BubbleTea 기반 깔끔한 레이아웃, Qwen3 모델 표시 |
| 브랜딩 | 2/10 | "OpenCode" 원본 브랜딩 그대로, slosim 정체성 없음 |
| 시뮬레이션 UI | 0/10 | Sim Dashboard 등 6개 컴포넌트 미통합 |
| 진행률 피드백 | 0/10 | 도구 실행 중 사용자 피드백 전무 |
| 에러 표시 | 3/10 | 기술 메시지 표시되나 친화적이지 않음 |
| **소계** | **12/50** | **24.0%** |

### 종합

| 카테고리 | 점수 | 비중 | 가중 점수 |
|----------|:----:|:----:|:---------:|
| 기능 동작 | 8.3% | 70% | 5.8% |
| UI/UX | 24.0% | 30% | 7.2% |
| **총점** | | | **13.0%** |

---

## 7. 권장 개선 사항 (우선순위순)

### P0 — 즉시 수정 (v0.3 출시 차단)

1. **에이전트 루프 보호** (`agent.go`)
   - 최대 반복 횟수 제한 (기본 25회)
   - 동일 도구+파라미터 연속 3회 감지 시 자동 중단 + 사용자 알림
   - 중단 시 "도구 호출이 반복되고 있습니다. 다른 방법을 시도하겠습니다." 메시지

2. **도구 라우팅 강화** (`sloshing_coder.go` 또는 에이전트 코드)
   - 방법 A: 시스템 프롬프트를 더 단순/강력하게 ("ONLY use xml_generator for new simulations")
   - 방법 B: 에이전트 코드에서 첫 도구 호출을 xml_generator로 강제
   - 방법 C: 도구 필터링 — 첫 턴에서 parametric_study 숨김
   - **권장**: 방법 A + C 조합

3. **도구 결과에 다음 단계 안내 추가**
   - `xml_generator` 결과: "→ 다음: gencase를 호출하세요"
   - `parametric_study` 결과: "→ 다음: 각 케이스에 대해 gencase → solver 실행하세요"

### P1 — v0.3-beta (1주 내)

4. **TUI 컴포넌트 통합** — Sim Dashboard, Results Viewer를 tui.go에 라우팅
5. **브랜딩 커스터마이징** — "OpenCode" → "slosim-agent", URL 교체
6. **도구 실행 로그 스트림** — 채팅 UI에 도구 이름+상태 실시간 표시
7. **Qwen3:32B 전용 최적화** — 8B는 도구 선택 능력 부족, 최소 32B 권장

### P2 — v0.4 (장기)

8. 파이프라인 스텝퍼 UI (xml → gen → solver → vtk → measure 시각화)
9. Case Wizard 연동 (양식 기반 시뮬레이션 설정)
10. 시뮬레이션 취소/재시도 UX

---

## 8. 결론

v0.3은 TUI 디자인 시스템(SemanticTokens, Widget Library) 측면에서 좋은 기반을 마련했으나, **핵심 기능(시뮬레이션 실행)이 전혀 작동하지 않는** 상태입니다. 근본 원인은:

1. **Qwen3:8B의 tool selection 능력 부족** — 18개 도구 중 올바른 도구를 선택하지 못함
2. **에이전트 루프 보호 부재** — 무한 반복을 방지하는 안전장치 없음
3. **TUI 컴포넌트 미통합** — 6개 신규 컴포넌트가 설계만 되고 연결되지 않음

**출시 판정: NOT READY** — P0 이슈 3건 해결 후 재평가 필요.
