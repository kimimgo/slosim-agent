---
name: sim-monitor
description: DualSPHysics 시뮬레이션 실시간 모니터링. Run.csv 파싱으로 진행률, 에너지, 시간 추적. "모니터링", "진행률", "monitor", "progress" 시 자동 발동.
autoTrigger: true
triggerPattern: "모니터링|진행률|monitor|progress|상태"
---

# sim-monitor: 실시간 시뮬레이션 모니터링

> Run.csv를 주기적으로 파싱하여 시뮬레이션 진행 상태를 추적

## Run.csv 포맷

```csv
Run;TimeStep;Time;Steps;FinalTime;Runtime;Particles;...
1;0;0.000000;0;2.000000;0.00;12345;...
2;100;0.001234;100;2.000000;1.23;12345;...
```

## 핵심 메트릭

| 메트릭 | Run.csv 컬럼 | 의미 |
|--------|-------------|------|
| 진행률 | `Time / FinalTime` | 시뮬레이션 완료 비율 |
| ETA | `(FinalTime - Time) * (Runtime / Time)` | 예상 잔여 시간 |
| 파티클 수 | `Particles` | 활성 파티클 (메모리 지표) |
| 타임스텝 | `TimeStep` | 현재 타임스텝 번호 |
| 런타임 | `Runtime` | 누적 실행 시간 (초) |

## 모니터링 워크플로우

```
1. Job ID로 시뮬레이션 디렉토리 특정
2. Run.csv 존재 확인 (솔버 시작 후 생성)
3. 마지막 행 파싱 → 진행률 계산
4. TUI SimDashboard로 업데이트 전송 (pubsub)
5. 완료 시 최종 요약 리포트
```

## TUI 연동

```go
// pubsub 이벤트로 Dashboard 업데이트
pubsub.Publish(SimProgressEvent{
    JobID:    jobID,
    Progress: progress,
    ETA:      eta,
    Step:     currentStep,
})
```

## 에러 감지

| 패턴 | 의미 | 대응 |
|------|------|------|
| Run.csv 미생성 (30초 이상) | GenCase 실패 또는 GPU 할당 실패 | 로그 확인 제안 |
| 파티클 수 급감 | 파티클 유출 (경계 문제) | dp 조정 또는 경계 확인 제안 |
| 런타임/스텝 비율 급증 | CFL 조건 위반 | dt 조정 제안 |
| 진행률 정체 | GPU 행 또는 디스크 I/O | 프로세스 상태 확인 |

## 규칙

1. 폴링 간격: 5초 (기본) — 부하 최소화
2. Run.csv 파일 락 대비: 읽기 실패 시 재시도
3. 완료 판정: `Time >= FinalTime * 0.999`
4. 모니터링 결과는 TUI와 로그 양쪽에 기록
