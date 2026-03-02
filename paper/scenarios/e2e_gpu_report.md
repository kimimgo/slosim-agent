# E2E 유저 시나리오 테스트 리포트

**날짜**: 2026-02-22
**시나리오**: 비전문가 자연어 슬로싱 시뮬레이션 요청
**바이너리**: `./slosim` (godotenv + context fix + tool order fix)

## 테스트 입력

```
가로 0.5m 세로 0.3m 높이 0.4m 직사각형 탱크에 물 50% 채우고,
좌우 흔들림 주기 0.5초 진폭 2cm로 1초간 시뮬레이션 해줘
```

## 발견 및 수정한 버그 (3건)

### Bug 1: LOCAL_ENDPOINT 환경변수 미로딩
- **증상**: `unsupported model configured, reverting to default` → gemini-2.5 폴백
- **원인**: `.env` 파일 로딩 기능 없음 (`godotenv` 미사용)
- **수정**: `local.go` init()에서 `godotenv.Load()` 호출
- **파일**: `internal/llm/models/local.go`

### Bug 2: context_window 4096 하드코딩
- **증상**: `Context: ⚠(253%)` — 컨텍스트 초과로 tool call 생성 불가
- **원인**: Ollama `/v1/models` API가 `loaded_context_length=0` 반환 → `cmp.Or(0, 4096)` = 4096
- **수정**: 폴백값을 32768로 변경, `MaxContextLength` 우선 사용
- **파일**: `internal/llm/models/local.go` `convertLocalModel()`

### Bug 3: maxTools=5로 DualSPHysics 도구 미전송
- **증상**: Qwen3가 thinking만 출력, tool call 생성 안 됨
- **원인**: `maxTools: 5` 제한으로 도구 배열의 처음 5개(Bash, Edit, Fetch, Glob, Grep)만 전송. `xml_generator`(인덱스 12)가 잘림
- **수정**: (1) 도구 순서를 DualSPHysics 우선으로 재배치, (2) maxTools를 5→20으로 변경
- **파일**: `internal/llm/agent/tools.go`, `internal/llm/provider/provider.go`

## 파이프라인 진행 결과

| 단계 | 도구 | 결과 | 캡처 |
|------|------|------|------|
| 0. TUI 기동 | — | ✅ 성공 | `run4_step0.txt` |
| 1. 프롬프트 전송 | — | ✅ 성공 | `run4_step1_180s.txt` |
| 2. Qwen3 파라미터 추론 | — | ✅ 정확 (dp=0.006, freq=2Hz, fluid_h=0.2m) | `run4_step2_360s.txt` |
| 3. XML 생성 | `xml_generator` | ✅ 성공 (`sloshing_case.xml` 검증 완료) | `run4_step3_480s.txt` |
| 4. 해석 준비 | `gencase` | ❌ Docker 경로 에러 (CUDA 관련) | `run4_step5_780s.txt` |
| 5. 에러 복구 | — | ⏳ Qwen3 응답 지연 (thinking 루프) | `run4_step7_1140s.txt` |

## 생성된 XML 검증

`simulations/sloshing_case.xml`:
- 탱크 치수: 0.5 x 0.3 x 0.4m ✅
- 유체 높이: 0.2m (50%) ✅
- dp: 0.006m ✅
- 가진: freq=2Hz, amplitude=0.02m, 수평(x축) ✅
- 시뮬레이션 시간: 1초 ✅
- SWL 게이지: 3개 (Left, Center, Right) ✅
- Symplectic + Wendland 커널 ✅

## 캡처 파일 목록 (28개)

총 4회 실행(run1~run4), 단계별 텍스트 캡처 저장.

## 잔여 이슈

1. **GenCase Docker 경로**: `case_path=simulations/sloshing_case`가 Docker 내부 경로(`/cases/...`)와 불일치
2. **Qwen3 에러 복구 지연**: 에러 후 thinking이 과도하게 길어짐 (>5분 응답 없음)
3. **Qwen3 thinking 표시**: `<think>` 태그 내용이 사용자에게 노출됨 — 필터링 필요
