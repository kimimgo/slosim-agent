# pv-agent Improvement Request

> **From**: slosim-agent (DualSPHysics sloshing simulation post-processing)
> **Date**: 2026-02-16 (updated 2026-02-17)
> **Context**: 10-case production post-processing (72K~669K SPH particles, VTK Legacy format)
> **pv-agent version**: commit 744d954 (main)

---

## Executive Summary

slosim-agent에서 pv-agent MCP 서버를 10개 슬로싱 시뮬레이션 케이스 후처리에 실전 투입한 결과, **7개 개선 요청**을 식별했다. 핵심은 VTK 파일 시리즈 지원 부재와 비디오 출력 미지원으로 인해 호출 에이전트 쪽에서 500줄+ 커스텀 스크립트를 작성해야 했다는 점이다.

### Priority

| Priority | 이슈 | Impact |
|----------|------|--------|
| **P0-CRITICAL** | Docker 컨테이너 kill 미흡 → CPU 90% 무한 행 | **시스템 장애** (2026-02-17 발견) |
| **P0-CRITICAL** | ffmpeg 타임아웃 미처리 → MCP 서버 크래시 | 서버 불안정 |
| P0 | VTK 파일 시리즈 미지원 | animate 도구 사용 불가 → 커스텀 스크립트 필수 |
| P0 | 비디오(MP4) 출력 미지원 | ffmpeg 수동 호출 필수 |
| P1 | 기본 타임아웃 120초 부족 | 대형 애니메이션 반복 실패 |
| P1 | speed_factor 중복 프레임 렌더링 | 100s 물리시간 × 0.2x = 7500 프레임 계산, 101개 VTK만 있음 |
| P1 | 출력 파일명 커스터마이징 불가 | `render.png` 하드코딩 → str.replace 해킹 필요 |
| P2 | startup 시 고아 컨테이너 정리 없음 | 재시작 후 잔여 컨테이너 누적 |
| P2 | 진행률(progress) 콜백 없음 | 100+ 프레임 렌더링 시 진행 상황 파악 불가 |
| P2 | Base 템플릿 scalar bar 색상 자동 감지 없음 | split_animate는 있으나 base.py.j2는 없음 |
| P3 | 배치 렌더링 도구 없음 | 동일 설정으로 N개 파일 일괄 스냅샷 불가 |

---

## P0-CRITICAL-1: Docker 컨테이너 Kill 미흡 (2026-02-17)

### 증상

- pv-agent MCP 서버 Python 프로세스가 CPU 90.8% 소비, 954분 CPU 시간 (무한 행)
- kitware/paraview-for-ci:v5.13.0 Docker 컨테이너 4개가 2일간 고아 상태 실행
- MCP inspect_data 호출 시 50분간 응답 없음 (블로킹)

### Root Cause

위치: `src/pv_agent/core/runner.py` - `_run_docker()` 타임아웃 핸들링

현재 코드에서 `proc.kill()`은 `docker run` CLI 프로세스만 SIGKILL합니다.
실제 Docker 컨테이너는 containerd가 별도 관리하므로 CLI가 죽어도 컨테이너는 계속 실행됩니다.
`docker run --rm`은 정상 종료 시에만 auto-remove하므로, CLI가 강제 종료되면 컨테이너가 영구적으로 고아화됩니다.

연쇄 효과:
1. 애니메이션 렌더링 요청 -> 120초 타임아웃
2. CLI kill -> 컨테이너는 여전히 pvpython 168프레임 풀 렌더링 중
3. 다음 MCP 호출 -> 또 Docker 컨테이너 생성 -> 또 타임아웃 -> 또 고아화
4. 4개 컨테이너가 동시에 GPU/CPU 점유 -> 시스템 전체 90% CPU

### 수정 제안

`_run_docker()`에서 `--name` 파라미터로 컨테이너에 고유 이름을 부여하고,
타임아웃 시 `docker stop <name>` + `docker rm -f <name>`으로 실제 컨테이너도 정리해야 합니다.

---

## P0-CRITICAL-2: ffmpeg 타임아웃 미처리

### Root Cause

위치: `src/pv_agent/pipeline/engine.py` - `compile_video()` 부근

`asyncio.wait_for(proc.communicate(), timeout=300)` 에서
`asyncio.TimeoutError`가 처리 없이 전파되어 MCP 서버 크래시를 유발합니다.

### 수정 제안

try/except로 TimeoutError를 잡아서 proc.kill() 후 에러 메시지를 반환해야 합니다.

---

## P1-NEW: 타임아웃 기본값 상향 + 도구별 전달

위치: `src/pv_agent/config.py` + `src/pv_agent/pipeline/engine.py`

- PV_TIMEOUT 기본값: 120초 -> 600초로 상향 필요
- engine.py의 execute 호출 시 timeout을 명시적으로 전달해야 합니다.

---

## P2-NEW: 서버 시작 시 고아 컨테이너 정리

서버 시작 시 이전 세션의 pv_agent_ 접두사 컨테이너를 자동 정리하는 startup hook이 필요합니다.

---

## P0-1: VTK File Series Support (`SourceDef.files`)

### Problem

DualSPHysics (및 대부분의 particle-based CFD)는 타임스텝별 **개별 VTK 파일**을 출력한다:

```
PartFluid_0000.vtk
PartFluid_0001.vtk
...
PartFluid_0100.vtk
```

현재 `SourceDef.file`은 단일 파일만 받으므로, animate 도구에 전체 시리즈를 전달할 수 없다. 실제 워크어라운드:

```python
# 현재: 컴파일된 스크립트를 문자열 치환으로 패칭 (fragile!)
script = compiler.compile(pipeline)
script = script.replace(
    f'FileNames=["{vtk_files[0]}"]',
    f"FileNames={json.dumps(vtk_files)}",
)
```

### Proposed Solution

```python
class SourceDef(BaseModel):
    file: str | None = None           # Single file (existing)
    files: list[str] | None = None    # File series (NEW)
    file_pattern: str | None = None   # Glob pattern (NEW), e.g. "PartFluid_*.vtk"
    timestep: float | str | None = None
    blocks: list[str] | None = None
```

`base.py.j2` 템플릿에서:

```jinja2
{% if source_files %}
reader = LegacyVTKReader(FileNames={{ source_files | tojson }})
{% else %}
reader = {{ reader_class }}(FileName="{{ source_file }}")
{% endif %}
```

### Impact

- animate, split_animate 모든 애니메이션 도구에서 VTK 시리즈 직접 지원
- 호출 에이전트의 str.replace 해킹 제거
- file_pattern 지원 시 LLM이 glob 패턴만 전달하면 됨 (더 자연스러운 인터페이스)

---

## P0-2: Video Output (MP4/WebM)

### Problem

animate와 split_animate 도구 모두 **PNG 프레임만 출력**한다. 실제 사용에서는 거의 항상 MP4 비디오가 필요하며, 호출 에이전트가 직접 ffmpeg를 호출해야 한다:

```python
# 현재: 호출 에이전트가 직접 구현해야 하는 코드
def _compile_ffmpeg(frames_dir, output_mp4, fps, text):
    cmd = [
        "ffmpeg", "-y",
        "-framerate", str(fps),
        "-i", os.path.join(frames_dir, "frame_%06d.png"),
        "-vf", f"drawtext=text='{text}':fontsize=28:fontcolor=white...",
        "-c:v", "libx264", "-pix_fmt", "yuv420p",
        "-preset", "fast", "-crf", "23",
        output_mp4,
    ]
    subprocess.run(cmd, ...)
```

### Proposed Solution

AnimationDef에 비디오 출력 옵션 추가:

```python
class AnimationDef(BaseModel):
    # ... existing fields ...
    output_format: Literal["frames", "mp4", "webm", "gif"] = "frames"
    video_codec: str = "libx264"     # h264 (mp4), libvpx-vp9 (webm)
    video_quality: int = 23          # CRF value (lower=better, 18-28)
    text_overlay: str | None = None  # ffmpeg drawtext
```

Runner에서 프레임 렌더링 후 자동 ffmpeg 호출:

```python
# runner.py
async def _compile_video(self, frames_dir, output_path, fps, format, ...):
    """Compile PNG frames to video using ffmpeg."""
    ...
```

### Impact

- animate/split_animate에서 `output_format="mp4"` 한 줄로 비디오 생성
- 호출 에이전트의 ffmpeg 코드 제거 (~30줄)
- ffmpeg 유무 자동 감지 + 미설치 시 PNG fallback

---

## P1-1: Frame Deduplication in Animation Sampling

### Problem

`speed_factor=0.2` (5배 슬로우모션) + 물리시간 100초인 경우:

```
animation_duration = 100 / 0.2 = 500초
target_frames = 500 * 15fps = 7,500 프레임
실제 VTK 파일 수 = 101개
```

7,500개 프레임 중 **동일 VTK를 최대 74회** 중복 렌더링한다. 각 프레임이 5~20초 걸리면 수 시간 낭비.

### Workaround Applied

```python
# slosim-agent postprocess script에서 적용한 수정
target_frames = max(1, min(int(round(anim_dur * fps)), len(all_timesteps)))

sampled = []
seen = set()
for fi in range(target_frames):
    t_target = ...
    best_idx = min(range(len(all_timesteps)), key=lambda k: abs(all_timesteps[k] - t_target))
    if best_idx not in seen:
        seen.add(best_idx)
        sampled.append(all_timesteps[best_idx])
```

### Proposed Solution

`animate.py.j2`와 `split_animate.py.j2` 템플릿 모두에서:

```python
# Cap target_frames at available timestep count to prevent duplicate renders
_target_frames = max(1, min(int(round(_anim_duration * _fps)), len(_all_timesteps)))

# Deduplicate sampled timesteps
_seen_indices = set()
_sampled = []
for fi in range(_target_frames):
    _t_target = _all_timesteps[0] + (fi / max(1, _target_frames - 1)) * _physics_duration
    _best_idx = min(range(len(_all_timesteps)), key=lambda k: abs(_all_timesteps[k] - _t_target))
    if _best_idx not in _seen_indices:
        _seen_indices.add(_best_idx)
        _sampled.append(_all_timesteps[_best_idx])
```

결과 JSON에 `effective_fps` 추가 (비디오 컴파일에 필요):

```python
_effective_fps = len(_sampled) / _anim_duration if _anim_duration > 0 else _fps
result["effective_fps"] = _effective_fps
```

### Impact

- 101 VTK → 101 렌더링 (7,500 → 101, **74배 속도 향상**)
- `effective_fps` 반환으로 ffmpeg/MP4 컴파일 시 올바른 재생 속도 보장

---

## P1-2: Customizable Output Filename

### Problem

render 도구가 `render.png`으로 하드코딩:

```python
# base.py.j2
SaveScreenshot(os.path.join(OUTPUT_DIR, "render.png"), ...)
```

동일 케이스에서 Press 스냅샷과 Vel 스냅샷을 구분하려면 파일명 변경이 필요한데, 현재는 str.replace 해킹만 가능:

```python
script = compiler.compile(pipeline)
script = script.replace(
    'os.path.join(OUTPUT_DIR, "render.png")',
    f'os.path.join(OUTPUT_DIR, "snapshot_press.png")',
)
```

### Proposed Solution

OutputDef에 filename 필드 추가:

```python
class OutputDef(BaseModel):
    # ... existing fields ...
    filename: str | None = None  # None = auto ("render.png", "frame_%06d.png", etc.)
```

render tool 레벨에서:

```python
@mcp.tool()
async def render(
    file_path: str,
    field_name: str,
    # ... existing params ...
    output_filename: str = "render.png",  # NEW
) -> Image:
```

### Impact

- str.replace 해킹 제거
- 동일 디렉토리에 여러 필드 스냅샷 저장 가능
- animate에서도 `frame_prefix` 파라미터로 확장 가능

---

## P2-1: Progress Reporting

### Problem

100+ 프레임 애니메이션 렌더링 시 (실측 5~20분), 호출 에이전트가 진행 상황을 전혀 파악할 수 없다. MCP 호출은 완료까지 블로킹되므로, LLM이 사용자에게 "진행 중"이라고 말할 수도 없다.

### Proposed Solution

MCP notifications 또는 result.json 중간 업데이트:

**Option A**: Progress 파일 (간단)

```python
# animate.py.j2 프레임 루프 내부
for _frame_i, _t in enumerate(_sampled):
    # ... render frame ...
    # Write progress
    _progress = {"frame": _frame_i + 1, "total": len(_sampled), "pct": round((_frame_i + 1) / len(_sampled) * 100)}
    Path(os.path.join(OUTPUT_DIR, "progress.json")).write_text(json.dumps(_progress))
```

**Option B**: FastMCP notifications (MCP 2025.11+ spec)

```python
@mcp.tool()
async def animate(...):
    async for progress in animate_impl_streaming(...):
        await mcp.notify_progress(progress)
```

### Impact

- 호출 에이전트가 사용자에게 "3/101 프레임 렌더링 중..." 표시 가능
- 타임아웃 판단 근거 제공 (진행이 멈춘 건지 느린 건지 구분)

---

## P2-2: Base Template Scalar Bar Color Auto-Detection

### Problem

`split_animate.py.j2`는 배경색 밝기를 감지하여 scalar bar 텍스트 색상을 자동 조정:

```python
# split_animate.py.j2 (잘 되어 있음)
_bg = {{ pane.background }}
_tc = [0,0,0] if (0.299*_bg[0] + 0.587*_bg[1] + 0.114*_bg[0]) > 0.5 else [1,1,1]
_sb.TitleColor = _tc
_sb.LabelColor = _tc
```

그러나 `base.py.j2` (render/animate 공통 템플릿)에는 이 로직이 **없다**. 어두운 배경(`[0.1, 0.1, 0.15]`)에서 scalar bar 텍스트가 검은색으로 렌더링되어 안 보임.

### Proposed Solution

`base.py.j2`의 scalar bar 섹션에 동일 로직 추가.

---

## P3-1: Batch Rendering Tool

### Problem

10개 케이스 × 2 필드 = 20 스냅샷을 렌더링하려면 render 도구를 20번 개별 호출해야 한다. 각 호출마다 pvpython 프로세스 시작/종료 오버헤드 발생 (~3초).

### Proposed Solution

```python
@mcp.tool()
async def batch_render(
    renders: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    """Render multiple snapshots in a single pvpython session.

    Each item in renders should contain:
        file_path, field_name, and optionally colormap, camera, output_filename, etc.

    Returns list of {filename, size_bytes, status}.
    """
```

단일 pvpython 세션에서 여러 파일을 순차 렌더링. 프로세스 시작 오버헤드 1회만 발생.

---

## Appendix: Production Run Data

### Test Environment

- **GPU**: NVIDIA RTX 4090 (24GB VRAM, 14.4GB occupied by uvicorn service)
- **ParaView**: 6.0.1, EGL headless (`--force-offscreen-rendering`)
- **OS**: Ubuntu 22.04, Linux 5.15
- **pv-agent render_backend**: gpu (EGL)

### 10-Case Results

| Case | Particles | VTK Files | Snapshot (2x) | Animation | Total Time* |
|------|----------|-----------|---------------|-----------|-------------|
| e2e_test | 145K | 101 | 247+231 KB | 0.3MB, 101fr | ~9 min |
| chen_nearcrit | 260K | 101 | 136+119 KB | 0.4MB, 101fr | ~12 min |
| chen_shallow | 119K | 101 | 121+116 KB | 0.1MB, 101fr | ~7 min |
| english_dbc | 633K | 93 | 1123+1123 KB | 32MB, 49fr | ~18 min |
| isope_lng | 301K | 101 | 149+265 KB | 0.1MB, 101fr | ~13 min |
| liu_amp1 | 669K | 101 | 1489+1489 KB | 30MB, 49fr | ~20 min |
| liu_large | 192K | 101 | 248+166 KB | 0.3MB, 101fr | ~10 min |
| nasa_cylinder | 143K | 101 | 627+757 KB | 0.1MB, 101fr | ~9 min |
| spheric_low | 72K | 101 | 265+131 KB | 0.3MB, 101fr | ~6 min |
| zhao_horizcyl | 529K | 101 | 670+622 KB | 0.0MB, 101fr | ~16 min |

*Total time = 2 snapshots + animation frames + ffmpeg + stats extraction (approximate)

### Workaround Code Volume

| Workaround | Lines | For Issue |
|-----------|-------|-----------|
| `_build_animation_pvscript()` | 95 | P0-1 (VTK series) + P1-1 (dedup) |
| `_compile_ffmpeg()` | 25 | P0-2 (video output) |
| `str.replace` filename hack | 4 | P1-2 (output filename) |
| Dynamic timeout calculation | 3 | P2-1 (no progress info) |
| **Total** | **~130 lines** | Ideally 0 with fixes |

### Observations

1. **english_dbc (633K) / liu_amp1 (669K)**: 고해상도 케이스에서 프레임당 렌더링 시간이 ~15-20초. 중복 프레임 제거(P1-1) 없이는 수 시간 소요했을 것.
2. **MP4 크기 이상**: 일부 케이스의 MP4가 0.0~0.1MB로 매우 작음. 이는 ffmpeg `drawtext` 필터의 특수문자 이스케이프 문제로 빈 프레임이 생성되었을 가능성. 비디오 출력이 pv-agent 내부에 통합되면 이런 문제 감소.
3. **Point Gaussian representation**: `base.py.j2`에서 `ShaderPreset = 'Sphere'`가 설정되는지 확인 필요. `split_animate.py.j2`에는 있으나 base에는 미확인.
