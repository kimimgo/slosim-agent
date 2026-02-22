---
name: pv-render
description: pv MCP 서버 호출로 후처리 시각화. 시뮬레이션 결과를 렌더링, 슬라이스, 등고선, 애니메이션으로 변환. "렌더링", "시각화", "애니메이션", "render", "visualize" 시 자동 발동.
autoTrigger: true
triggerPattern: "렌더링|시각화|애니메이션|render|visualize|postprocess"
---

# pv-render: ParaView MCP 후처리 시각화

> 시뮬레이션 완료 후 pv MCP 서버를 호출하여 시각화 수행

## 워크플로우

```
1. inspect_data  → 파일 메타데이터 확인 (필드, 타임스텝, 범위)
2. render/slice  → 단일 뷰 시각화 → PNG
3. contour/clip  → 등고선/클리핑 시각화
4. animate       → 시계열 애니메이션 → GIF
5. split_animate → 멀티 패널 동기화 애니메이션
```

## pv MCP 도구 매핑

| 시각화 요청 | MCP 도구 | 출력 |
|------------|---------|------|
| "결과 확인" | `inspect_data` | 필드 목록, 범위 |
| "렌더링" | `render` | PNG |
| "슬라이스" | `slice` | PNG |
| "등고선" | `contour` | PNG |
| "애니메이션" | `animate` | GIF |
| "비교 애니메이션" | `split_animate` | GIF (2-4 pane) |
| "유선" | `streamlines` | PNG |
| "통계" | `extract_stats` | JSON |

## 사용 예시

```
사용자: "시뮬레이션 결과 시각화해줘"
→ inspect_data로 VTK 파일 확인
→ render로 초기 뷰 생성
→ animate로 시계열 GIF 생성

사용자: "압력 분포 슬라이스"
→ slice (field="Pressure", normal=[0,1,0])
→ contour (field="Pressure", levels=10)
```

## 컬러맵 가이드

| 물리량 | 권장 컬러맵 |
|--------|-----------|
| Velocity | coolwarm, jet |
| Pressure | viridis, plasma |
| VonMises | hot, inferno |
| Temperature | coolwarm |

## 규칙

1. 항상 `inspect_data`로 시작하여 가용 필드 확인
2. 렌더링 전 카메라 각도를 데이터 범위에 맞게 설정
3. 애니메이션은 0.2배속 슬로우모션 기본
4. 결과 파일 경로: `./simulations/<job_id>/` 하위
