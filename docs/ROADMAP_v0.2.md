# slosim-agent v0.2 Roadmap — ✅ COMPLETE

> 2026-02-14 수립, 2026-02-15 완료

## 최종 상태

- **19/19 GitHub Issues 전부 종료** (alpha 13 + beta 6)
- E2E GPU 파이프라인 **10/10 시나리오 검증 완료** (RTX 4090)
- STL watertight 실구현 + Frosina 시나리오 검증
- Agent Teams 병렬 개발 성공 (lead + 2 sonnet agents)
- `go build/test/vet` 전부 PASS

---

## v0.2-alpha (~ 02/28) — 수정 반영 + 핵심 통합 검증

> 목표: 실제로 돌아가는 상태 만들기

| # | 이슈 | 유형 | 우선순위 |
|---|------|------|---------|
| #18 | E2E 버그 수정 커밋 | chore | P0 |
| #1 | Docker 바이너리명 불일치 ✅수정됨 | fix | P0 |
| #2 | GenCase -save: 파라미터 ✅수정됨 | fix | P0 |
| #3 | SimulationDomain 마진 ✅수정됨 | fix | P0 |
| #4 | MeasureTool POINTS 헤더 ✅수정됨 | fix | P0 |
| #6 | **slosim-agent 바이너리 E2E** | test | P0 |
| #7 | **SloshingCoderPrompt + Qwen3 품질** | test | P0 |
| #5 | ParaView Mesa 호환성 workaround | fix | P1 |
| #9 | pvpython/animation.go Mesa 반영 | feat | P1 |
| #8 | monitor + error_recovery E2E | test | P1 |
| #11 | xml_generator POINTS 헤더 자동 | feat | P1 |
| #19 | Go 테스트 integration tag 분리 | test | P1 |

### alpha 완료 기준 ✅
- [x] 버그 수정 커밋 (#18)
- [x] `./slosim-agent`로 전체 파이프라인 1회 자동 실행 (#6)
- [x] Qwen3 + SloshingCoderPrompt가 올바른 도구 호출 (#7)
- [x] `go test ./...` CI에서 PASS (#19)

---

## v0.2-beta (~ 03/15) — 전체 도구 E2E + TUI

> 목표: 시나리오 18기능 전부 검증

| # | 이슈 | 유형 | 우선순위 |
|---|------|------|---------|
| #10 | xml_generator 모션 기반 마진 자동 | feat | P1 |
| #12 | TUI 5개 컴포넌트 검증 | test | P2 |
| #13 | geometry + seismic_input E2E | test | P2 |
| #14 | analysis + report E2E | test | P2 |
| #16 | monitor.go TimeMax 동적 읽기 | feat | P2 |

### beta 완료 기준 ✅
- [x] 시나리오 Act 1~6, 8 전부 E2E 통과
- [x] TUI Dashboard 기본 구현
- [x] 원통형 탱크 시뮬레이션 성공 (NASA + Zhao HorizCyl)
- [ ] CSV 파도 데이터 임포트 → v0.3으로 이동

---

## v0.2-rc (~ 03/31) — 고급 기능 + 최적화

> 목표: 시나리오 전체 통과 + 사용자 피드백

| # | 이슈 | 유형 | 우선순위 |
|---|------|------|---------|
| #15 | parametric_study + result_store + stl_import E2E | test | P2 |
| #17 | stl_import.go watertight 실제 구현 | feat | P2 |

### rc 완료 기준 ✅
- [x] GPU E2E 10개 논문 시나리오 검증 (8 PASS, 2 PARTIAL)
- [x] STL 임포트 → 시뮬레이션 성공 (Frosina Fuel Tank)
- [ ] 파라메트릭 스터디 → v1.0으로 이동

---

## 이슈 의존 관계

```
#18 커밋 ──→ #6 바이너리 E2E ──→ #7 LLM 품질
                  │
                  ├──→ #8 monitor/error_recovery
                  ├──→ #9 Mesa 스크립트 반영
                  └──→ #19 test tag 분리

#11 POINTS 헤더 ──→ #13 geometry/seismic E2E
#10 마진 자동 ─────→ #13 geometry/seismic E2E

#6 바이너리 E2E ──→ #12 TUI 검증
                  ──→ #14 analysis/report E2E
                  ──→ #15 parametric/store/stl E2E
```

---

## 핵심 원칙

1. **실제 실행 우선**: mock/dry-run이 아니라 실제 Docker + GPU + LLM으로 검증
2. **시나리오 기반**: "김 대리" 시나리오의 각 Act를 acceptance criteria로 사용
3. **버그는 즉시 커밋**: 발견 → 수정 → 커밋 → 다음 단계 (미커밋 누적 금지)
4. **인간 판단 영역 분리**: 물리 정확성, UX 품질, 시각화 품질은 사용자 검토
