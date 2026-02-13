# slosim-agent v0.1 기능 요구사항 정의서 (FRD)

*v0.1 — 2026-02-13*
*Status: Draft (리뷰 대기)*

---

## 1. v0.1 범위 정의

### 1.1 목표

**v0.1은 MVP**: 자연어 → DualSPHysics 시뮬레이션 → 결과 확인까지 최소 End-to-End 파이프라인이 동작하는 것.

### 1.2 v0.1 성공 기준 (Definition of Done)

다음 E2E 시나리오가 완전히 동작해야 v0.1이다:

```
사용자: "1m × 0.5m 직사각형 탱크에 물 50% 채우고 0.5Hz로 흔들어줘"
  ↓
AI가 해석 조건을 추론하고 사용자에게 확인 요청
  ↓
사용자: "좋아, 실행해"
  ↓
GenCase → DualSPHysics GPU 실행 (백그라운드) → 완료 대기
  ↓
PartVTK로 VTK 추출 + MeasureTool로 수위 CSV 추출
  ↓
Markdown 리포트 생성 (조건 요약 + 수위 시계열 데이터)
  ↓
사용자에게 "해석 완료. 리포트: simulations/<id>/report.md" 안내
```

### 1.3 PRD 대비 v0.1 스코프

| PRD 기능 | v0.1 (P0) | v0.2 (P1) | 향후 (P2) | 비고 |
|----------|:---------:|:---------:|:---------:|------|
| **IN-01** 자연어 입력 → 조건 추론 | **O** | | | 직사각형 탱크 한정 |
| **IN-02** STL 파일 입력 | | | **O** | 메시 품질 검증 복잡 |
| **IN-03** 단순 입력 → AI 제안 | **O** | | | "LNG 탱크 해석" → 표준 치수 |
| **SIM-01** DualSPHysics GPU 통합 | **O** | | | Docker 기반 |
| **SIM-02** GenCase XML 자동 생성 | **O** | | | 직사각형 + 정현파 |
| **SIM-03** AI 파라미터 자율 판단 | **O** | | | dp, TimeMax, freq 등 |
| **SIM-04** 백그라운드 Job 실행 | **O** | | | 단일 Job만 |
| **SIM-05** 실시간 모니터링 | | **O** | | v0.1은 완료/실패만 감지 |
| **POST-01** pvpython 후처리 | | **O** | | v0.1은 PartVTK + MeasureTool만 |
| **POST-02** 정지 이미지 스냅샷 | | **O** | | pvpython 의존 |
| **POST-03** 애니메이션 생성 | | | **O** | pvpython + ffmpeg |
| **POST-04** 시각화 항목 (압력, 속도 등) | | **O** | | v0.1은 VTK 파일만 |
| **RPT-01** Markdown 리포트 | **O** | | | 텍스트 + 수치 데이터 |
| **RPT-02** 해석 조건 요약 | **O** | | | |
| **RPT-03** 결과 이미지 포함 | | **O** | | pvpython 의존 |
| **RPT-04** AI 물리적 해석 | | **O** | | Qwen3 성능 검증 후 |
| **PARA-01~03** 파라메트릭 스터디 | | | **O** | 멀티 Job 관리 필요 |
| **DOM-01** 탱크 형상 다양화 | 직사각형만 | 원통형 | L형 등 | |
| **DOM-02** 유체: 물 | **O** | | | 단일 유체 |
| **DOM-03** 가진 조건 다양화 | 정현파만 | 지진파 | 임의 가진 | |
| **DOM-04** 2D/3D | 3D만 | 2D | | |

---

## 2. 기능 상세 명세

### TOOL-01: GenCase Tool (XML → 파티클 전처리)

| 항목 | 내용 |
|------|------|
| **ID** | TOOL-01 |
| **기능명** | GenCase Tool |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | 없음 (독립) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/gencase.go`, `gencase_test.go` |

**사용자 스토리:**
> AI 에이전트로서, XML 케이스 파일에서 파티클 지오메트리를 생성할 수 있어야 한다.
> 이를 통해 사용자의 자연어 요청을 실제 시뮬레이션 입력으로 변환할 수 있다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `case_path` | string | Y | XML 케이스 파일 경로 (`.xml` 확장자 제외) |
| `save_path` | string | Y | 출력 디렉토리 |
| `dp` | float | N | 파티클 간격 오버라이드 (m) |

**출력:**
- `{save_path}/{case_name}.bi4` — 바이너리 파티클 데이터
- `{save_path}/{case_name}.xml` — 처리된 XML 복사본
- `{save_path}/{case_name}_Bound.vtk` — 경계 파티클 VTK (선택)
- `{save_path}/{case_name}_Fluid.vtk` — 유체 파티클 VTK (선택)

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | GenCase가 유효한 XML에서 `.bi4` 파일 생성 성공 | `cases/SloshingTank_Def.xml` 입력 → `.bi4` 존재 확인 |
| AC-2 | 잘못된 XML 입력 시 한국어 에러 메시지 반환 | 빈 XML 입력 → `ToolResponse.IsError=true`, 한국어 메시지 |
| AC-3 | 경로에 `.xml` 포함 시 자동 제거 | `"Tank_Def.xml"` 입력 → `.xml` strip 후 실행 |
| AC-4 | Docker 컨테이너 내 실행 | `docker compose run --rm dsph GenCase ...` 호출 확인 |
| AC-5 | dp 오버라이드 파라미터 적용 | dp=0.02 지정 → GenCase `-dp:0.02` 인자 포함 |

**기술 명세:**
```
실행 명령: docker compose run --rm dsph GenCase {case_path} -save:{save_path}
GenCase는 .xml 자동 추가: case_path에 .xml 포함하지 않아야 함
dp 오버라이드: -dp:{value} 옵션 추가
```

---

### TOOL-02: Solver Tool (DualSPHysics GPU 실행)

| 항목 | 내용 |
|------|------|
| **ID** | TOOL-02 |
| **기능명** | DualSPHysics Solver Tool |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | TOOL-01 (GenCase 출력 필요) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/solver.go`, `solver_test.go` |

**사용자 스토리:**
> AI 에이전트로서, GenCase가 생성한 파티클 데이터에 대해 GPU SPH 시뮬레이션을 실행할 수 있어야 한다.
> 시뮬레이션은 백그라운드에서 실행되어 TUI가 블로킹되지 않아야 한다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `case_name` | string | Y | 케이스 이름 (GenCase 출력 기준) |
| `data_dir` | string | Y | 입력 데이터 디렉토리 (GenCase 출력 경로) |
| `out_dir` | string | Y | 시뮬레이션 출력 디렉토리 |
| `gpu` | bool | N | GPU 사용 여부 (기본: true) |

**출력:**
- `{out_dir}/Part_XXXX.bi4` — 타임스텝별 파티클 데이터
- `{out_dir}/Run.csv` — 실행 통계 (타임스텝, 시간, 파티클 수)
- `{out_dir}/Run.out` — 실행 로그
- 반환값: Job ID (백그라운드 추적용)

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | GPU 모드로 시뮬레이션 완료 (Part_*.bi4 생성) | `SloshingTank_Def` (0.1초) 실행 → Part 파일 존재 |
| AC-2 | 백그라운드 실행: 즉시 Job ID 반환 | `Run()` 호출 → 즉시 ToolResponse 반환, Job ID 포함 |
| AC-3 | 시뮬레이션 완료 감지 | Run.out 파일의 "Finished" 문자열 파싱 |
| AC-4 | 시뮬레이션 실패 감지 | RhopOut 초과 시 에러 상태 전환 |
| AC-5 | 실행 중 TUI 응답 유지 | 시뮬레이션 중 채팅 입력 가능 |

**기술 명세:**
```
실행 명령: docker compose run --rm dsph DualSPHysics5.4_linux64 {case_name} {out_dir} -gpu
백그라운드: goroutine으로 실행, Job ID로 추적
완료 감지: Run.out 파일의 "Finished" 문자열 또는 프로세스 종료
실패 감지: stderr에 "Exception", "Error", RhopOut 관련 메시지
```

---

### TOOL-03: PartVTK Tool (VTK 내보내기)

| 항목 | 내용 |
|------|------|
| **ID** | TOOL-03 |
| **기능명** | PartVTK Tool |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | TOOL-02 (Solver 완료 필요) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/partvtk.go`, `partvtk_test.go` |

**사용자 스토리:**
> AI 에이전트로서, 시뮬레이션 결과를 VTK 형식으로 내보낼 수 있어야 한다.
> 이를 통해 후처리 도구(ParaView 등)에서 시각화가 가능하다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `data_dir` | string | Y | 시뮬레이션 출력 디렉토리 (Part_*.bi4 위치) |
| `out_file` | string | Y | VTK 출력 파일 경로 (접미사 자동 추가) |
| `only_fluid` | bool | N | 유체 파티클만 (기본: true) |
| `vars` | string[] | N | 출력 변수 (기본: `["vel","rhop","press"]`) |
| `first` | int | N | 시작 타임스텝 |
| `last` | int | N | 종료 타임스텝 |

**출력:**
- `{out_file}_XXXX.vtk` — 타임스텝별 VTK 파일

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | Part_*.bi4 → VTK 변환 성공 | Solver 완료 후 실행 → `.vtk` 파일 생성 확인 |
| AC-2 | 유체만 추출 옵션 동작 | `only_fluid=true` → 경계 파티클 미포함 확인 |
| AC-3 | 변수 선택 동작 | `vars=["vel","press"]` → VTK에 해당 필드만 포함 |
| AC-4 | 타임스텝 범위 필터링 | `first=10, last=20` → 해당 범위만 변환 |

**기술 명세:**
```
실행 명령: docker compose run --rm dsph PartVTK -dirdata {data_dir} -filexml AUTO -savevtk {out_file} -onlytype:-bound -vars:{vars}
```

---

### TOOL-04: MeasureTool (수위/압력 측정)

| 항목 | 내용 |
|------|------|
| **ID** | TOOL-04 |
| **기능명** | MeasureTool |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | TOOL-02 (Solver 완료 필요) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/measuretool.go`, `measuretool_test.go` |

**사용자 스토리:**
> AI 에이전트로서, 특정 위치에서 수위, 압력, 속도 등의 시계열 데이터를 추출할 수 있어야 한다.
> 이를 통해 사용자에게 정량적 결과를 제공할 수 있다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `data_dir` | string | Y | 시뮬레이션 출력 디렉토리 |
| `points_file` | string | Y | 측정 포인트 파일 (x y z 좌표) |
| `out_csv` | string | Y | CSV 출력 파일 경로 |
| `vars` | string[] | N | 측정 변수 (기본: `["vel","rhop","press"]`) |
| `elevation` | bool | N | 수위 높이 계산 모드 |

**출력:**
- `{out_csv}_Vel.csv` — 속도 시계열
- `{out_csv}_Press.csv` — 압력 시계열
- `{out_csv}_Elevation.csv` — 수위 시계열 (elevation=true 시)

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | probe_points.txt 기반 측정 성공 | Solver 완료 후 → CSV 파일 생성, 데이터 행 > 0 |
| AC-2 | 수위 높이(elevation) 모드 동작 | `elevation=true` → `_Elevation.csv` 파일 생성 |
| AC-3 | CSV 파일 포맷 정합성 | 헤더 행 + 숫자 데이터, 파싱 가능 확인 |

**기술 명세:**
```
실행 명령: docker compose run --rm dsph MeasureTool -dirdata {data_dir} -filexml AUTO -points {points_file} -vars:{vars} -savecsv {out_csv}
수위 모드: -elevation 옵션 추가
기존 probe_points.txt: 좌벽(x=0.05), 중앙(x=0.50), 우벽(x=0.95) 각 6개 높이
```

---

### AGENT-01: Sloshing Coder 시스템 프롬프트

| 항목 | 내용 |
|------|------|
| **ID** | AGENT-01 |
| **기능명** | 슬로싱 전문 시스템 프롬프트 |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | 없음 (독립, 모든 Tool보다 선행) |
| **담당** | agent-eng |
| **파일** | `internal/llm/prompt/sloshing_coder.go` |

**사용자 스토리:**
> CFD 비전문가로서, "LNG 탱크 슬로싱 해석해줘"처럼 자연어로 요청하면
> AI가 탱크 치수, 유체 높이, 가진 주파수 등을 알아서 결정하고 확인해주길 원한다.

**프롬프트 요구사항:**

| 항목 | 내용 |
|------|------|
| 모델 | Qwen3 32B (Ollama 로컬) |
| 토큰 제한 | 시스템 프롬프트 ≤ 8K 토큰 |
| 언어 | 사용자 대면 = 한국어, Tool 호출 = JSON |
| 역할 | 슬로싱 해석 전문 AI 어시스턴트 |

**프롬프트에 포함해야 할 도메인 지식:**

```
1. 직사각형 탱크 1차 공진 주파수 공식:
   f₁ = (1/2π) × √(g × π/L × tanh(π/L × h))

2. 파라미터 자동 결정 규칙:
   - dp: 탱크 최소 치수 / 50 (최소 0.005, 최대 0.05)
   - TimeMax: 최소 5/freq 주기 (공진 근처면 10/freq)
   - TimeOut: TimeMax / 50 ~ 100
   - Visco: 0.01 (인공 점성, 기본값)

3. 슬로싱 시나리오 분류:
   - Normal: freq/f₁ < 0.8
   - Near-Resonance: 0.8 ≤ freq/f₁ < 0.95
   - Resonance: 0.95 ≤ freq/f₁ ≤ 1.05

4. 사용자 대면 용어 변환:
   - dp → "입자 간격 (해석 정밀도)"
   - GenCase → "해석 준비"
   - DualSPHysics → "시뮬레이션 실행"
   - PartVTK → "결과 변환"
```

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | "1m 탱크 슬로싱" 입력 → 치수/조건 JSON 추론 | 프롬프트 테스트 (Ollama API) |
| AC-2 | 추론 결과에 CFD 전문 용어 없음 | 응답에 "dp", "SPH", "CFL" 직접 노출 없음 |
| AC-3 | 공진 주파수 자동 계산 | L=1m, h=0.3m → f₁≈0.76Hz 근사값 산출 |
| AC-4 | Tool 호출 JSON 정합성 | gencase tool 호출 시 올바른 JSON 스키마 |
| AC-5 | 시스템 프롬프트 ≤ 8K 토큰 | 토큰 카운트 검증 |

---

### AGENT-02: XML 케이스 자동 생성

| 항목 | 내용 |
|------|------|
| **ID** | AGENT-02 |
| **기능명** | XML 케이스 자동 생성 |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | AGENT-01 (프롬프트에서 파라미터 결정) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/xml_generator.go`, `xml_generator_test.go` |

**사용자 스토리:**
> AI 에이전트로서, 추론된 해석 조건(치수, 유체 높이, 주파수, dp 등)을
> DualSPHysics 호환 XML 케이스 파일로 자동 변환할 수 있어야 한다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `tank_length` | float | Y | 탱크 길이 (m) |
| `tank_width` | float | Y | 탱크 너비 (m) |
| `tank_height` | float | Y | 탱크 높이 (m) |
| `fluid_height` | float | Y | 유체 높이 (m) |
| `freq` | float | Y | 가진 주파수 (Hz) |
| `amplitude` | float | Y | 가진 진폭 (m) |
| `dp` | float | Y | 파티클 간격 (m) |
| `time_max` | float | Y | 시뮬레이션 시간 (s) |
| `out_path` | string | Y | XML 출력 경로 |

**출력:**
- `{out_path}.xml` — DualSPHysics 호환 XML 케이스 파일
- `{out_path}_probe_points.txt` — 자동 생성된 측정 포인트 파일

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | 생성된 XML이 GenCase에서 성공 실행 | XML 생성 → GenCase 실행 → `.bi4` 생성 |
| AC-2 | XML은 attribute-only 문법 | `xmllint` + 커스텀 검증 (텍스트 노드 없음) |
| AC-3 | 필수 섹션 모두 포함 | `constantsdef`, `geometry`, `motion`, `execution/parameters` |
| AC-4 | SWL gauge 3개 자동 포함 | 좌벽, 중앙, 우벽 수위 게이지 |
| AC-5 | probe_points.txt 자동 생성 | 탱크 치수 기반 측정 포인트 좌표 |
| AC-6 | 기존 템플릿과 구조 호환 | `cases/Sloshing_Normal_Def.xml`과 동일 구조 |

**기술 명세 — XML 템플릿 구조:**

```xml
<?xml version="1.0" encoding="UTF-8" ?>
<case>
  <casedef>
    <constantsdef>
      <gravity x="0" y="0" z="-9.81" />
      <rhop0 value="1000" />              <!-- 물 밀도 (고정) -->
      <rhopgradient value="2" />
      <hswl value="0" auto="true" />
      <gamma value="7" />
      <speedsystem value="0" auto="true" />
      <coefsound value="20" />
      <speedsound value="0" auto="true" />
      <coefh value="1.2" />
      <cflnumber value="0.2" />
    </constantsdef>
    <mkconfig boundcount="240" fluidcount="9" />
    <geometry>
      <definition dp="{dp}">
        <pointmin x="-{margin}" y="-{margin}" z="-{margin}" />
        <pointmax x="{L+margin}" y="{W+margin}" z="{H+margin}" />
      </definition>
      <commands><mainlist>
        <setshapemode>dp | bound</setshapemode>
        <setdrawmode mode="full" />
        <setmkfluid mk="0" />
        <drawbox>
          <boxfill>solid</boxfill>
          <point x="0" y="0" z="0" />
          <size x="{L}" y="{W}" z="{fluid_height}" />
        </drawbox>
        <setmkbound mk="0" />
        <drawbox>
          <boxfill>bottom | left | right | front | back</boxfill>
          <point x="0" y="0" z="0" />
          <size x="{L}" y="{W}" z="{H}" />
        </drawbox>
        <shapeout file="Tank" />
      </mainlist></commands>
    </geometry>
    <motion>
      <objreal ref="0">
        <begin mov="1" start="0" />
        <mvrectsinu id="1" duration="{time_max}" anglesunits="degrees">
          <freq x="{freq}" y="0" z="0" />
          <ampl x="{amplitude}" y="0" z="0" />
          <phase x="0" y="0" z="0" />
        </mvrectsinu>
      </objreal>
    </motion>
  </casedef>
  <execution>
    <special><gauges>
      <!-- 좌벽, 중앙, 우벽 SWL 게이지 -->
    </gauges></special>
    <parameters>
      <!-- Symplectic, Wendland, Artificial viscosity 등 고정 파라미터 -->
      <parameter key="TimeMax" value="{time_max}" />
      <parameter key="TimeOut" value="{time_out}" />
    </parameters>
  </execution>
</case>
```

---

### JOB-01: Job Manager (백그라운드 실행 관리)

| 항목 | 내용 |
|------|------|
| **ID** | JOB-01 |
| **기능명** | Job Manager |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | TOOL-02 (Solver) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/job_manager.go`, `job_manager_test.go` |

**사용자 스토리:**
> 사용자로서, 시뮬레이션이 실행 중일 때도 채팅을 계속할 수 있어야 한다.
> "진행 상황 알려줘" 라고 물으면 현재 상태를 알려줘야 한다.

**Job 상태 머신:**

```
PENDING → RUNNING → COMPLETED
                  → FAILED
                  → CANCELLED
```

**입력/출력 — Job 제출:**

| 파라미터 | 타입 | 설명 |
|---------|------|------|
| `command` | string[] | Docker 실행 명령어 |
| `work_dir` | string | 시뮬레이션 작업 디렉토리 |

반환: `job_id` (UUID)

**입력/출력 — Job 조회:**

| 파라미터 | 타입 | 설명 |
|---------|------|------|
| `job_id` | string | Job UUID |

반환: `{ status, progress_pct, elapsed_time, error_message }`

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | Job 제출 시 즉시 반환 (non-blocking) | Submit() → 즉시 Job ID 반환, goroutine 실행 |
| AC-2 | Job 상태 조회 가능 | Status(job_id) → RUNNING/COMPLETED/FAILED |
| AC-3 | 완료/실패 시 pubsub 이벤트 발행 | JobCompletedMsg / JobFailedMsg 이벤트 확인 |
| AC-4 | 진행률 추정 (Part 파일 카운팅) | Part_*.bi4 파일 수 / 예상 총 파일 수 |
| AC-5 | v0.1: 단일 Job만 관리 | 동시 2개 제출 시 거절 |

**기술 명세:**
```go
// 진행률 추정
expectedParts = TimeMax / TimeOut  // 예: 5.0/0.1 = 50
currentParts = len(glob("Part_*.bi4"))
progressPct = currentParts / expectedParts * 100

// 완료 감지: Run.out에 "Finished" 또는 프로세스 exit code 0
// 실패 감지: stderr에 "Exception", "Error" 또는 exit code != 0
```

---

### RPT-01: Markdown 리포트 생성

| 항목 | 내용 |
|------|------|
| **ID** | RPT-01 |
| **기능명** | Markdown 리포트 생성 |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | TOOL-03 (PartVTK), TOOL-04 (MeasureTool) |
| **담당** | agent-eng |
| **파일** | `internal/llm/tools/report.go`, `report_test.go` |

**사용자 스토리:**
> 사용자로서, 시뮬레이션 완료 후 결과를 이해하기 쉬운 리포트로 받고 싶다.
> 어떤 조건으로 해석했는지, 수위 변화가 어땠는지 요약이 필요하다.

**입력:**

| 파라미터 | 타입 | 필수 | 설명 |
|---------|------|:----:|------|
| `sim_dir` | string | Y | 시뮬레이션 결과 디렉토리 |
| `case_config` | object | Y | 해석 조건 (탱크 치수, 주파수 등) |
| `csv_files` | string[] | N | MeasureTool CSV 파일 경로 |

**출력:**
- `{sim_dir}/report.md` — Markdown 리포트

**v0.1 리포트 구조:**

```markdown
# 슬로싱 해석 리포트

## 1. 해석 조건
- 탱크: {L}m × {W}m × {H}m (직사각형)
- 유체: 물 (1000 kg/m³), 높이 {h}m ({fill_ratio}% 충진)
- 가진: 정현파 {freq}Hz, 진폭 {amp}m
- 입자 간격: {dp}m (총 파티클 약 {N}개)
- 시뮬레이션 시간: {time_max}초

## 2. 실행 정보
- GPU: NVIDIA RTX 4090
- 실행 시간: {elapsed}
- 타임스텝 수: {steps}

## 3. 결과 요약
- VTK 파일: {vtk_count}개 생성 ({vtk_dir})
- 측정 데이터: {csv_files}

## 4. 수위 시계열 데이터
| 시간(s) | 좌벽(m) | 중앙(m) | 우벽(m) |
|---------|---------|---------|---------|
| ...     | ...     | ...     | ...     |
(MeasureTool CSV에서 추출)
```

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | 리포트 Markdown 파일 생성 | 시뮬레이션 완료 후 → `report.md` 존재 |
| AC-2 | 해석 조건 섹션에 모든 입력 파라미터 포함 | 탱크 치수, 주파수, dp, TimeMax 확인 |
| AC-3 | 수위 테이블 데이터 포함 (CSV가 있는 경우) | MeasureTool CSV → Markdown 테이블 변환 |
| AC-4 | CFD 전문 용어 없음 (사용자 대면 텍스트) | "dp" → "입자 간격", "SPH" 미사용 등 |

---

### TUI-01: 시뮬레이션 상태 표시

| 항목 | 내용 |
|------|------|
| **ID** | TUI-01 |
| **기능명** | 시뮬레이션 상태 메시지 |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | JOB-01 (Job Manager) |
| **담당** | tui-dev |
| **파일** | `internal/tui/components/core/status.go` (기존 수정) |

**사용자 스토리:**
> 사용자로서, 시뮬레이션 진행 상태를 채팅 화면에서 확인하고 싶다.
> "실행 중 (35%) — 약 2분 남음" 같은 안내가 보이면 좋겠다.

**v0.1 범위:** 기존 OpenCode 상태 바를 활용한 최소 구현. 별도 Sim Dashboard는 v0.2.

**표시 형태:**
```
[상태바] 🔄 시뮬레이션 실행 중 (35%) — 약 2분 남음
[상태바] ✅ 시뮬레이션 완료 — 결과: simulations/20260213_143022/
[상태바] ❌ 시뮬레이션 실패 — 시뮬레이션이 불안정해졌습니다
```

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | Job 실행 중 상태 바에 진행률 표시 | JobUpdateMsg 수신 → 상태 바 업데이트 |
| AC-2 | 완료 시 결과 경로 안내 | JobCompletedMsg → 결과 디렉토리 표시 |
| AC-3 | 실패 시 한국어 에러 안내 | JobFailedMsg → 원인 + 해결 방법 표시 |

---

### PIPE-01: E2E 파이프라인 오케스트레이션

| 항목 | 내용 |
|------|------|
| **ID** | PIPE-01 |
| **기능명** | 전체 파이프라인 자동 실행 |
| **우선순위** | P0 (v0.1 필수) |
| **의존성** | AGENT-01, AGENT-02, TOOL-01~04, JOB-01, RPT-01, TUI-01 (전체) |
| **담당** | agent-eng (Agent Loop에서 Tool 순차 호출) |
| **파일** | `internal/llm/agent/agent.go` (기존 수정) |

**사용자 스토리:**
> 비전문가 사용자로서, 한 번의 자연어 요청으로 전체 과정이 자동으로 진행되길 원한다.
> 중간에 확인이 필요한 부분만 물어보고, 나머지는 알아서 처리해야 한다.

**v0.1 파이프라인:**

```
1. [사용자 입력]  "1m 탱크 슬로싱 해석"
2. [AGENT-01]     조건 추론 → 사용자 확인 요청
3. [사용자 확인]  "좋아, 실행해"
4. [AGENT-02]     XML 케이스 생성 → simulations/{id}/case_Def.xml
5. [TOOL-01]      GenCase 실행 → 파티클 생성
6. [TOOL-02]      Solver 실행 (백그라운드) → Job ID 반환
7. [JOB-01]       진행 모니터링 → TUI 상태 업데이트
8. [대기]          완료/실패 이벤트 대기
9. [TOOL-03]      PartVTK → VTK 파일 생성
10.[TOOL-04]      MeasureTool → CSV 측정 데이터 추출
11.[RPT-01]       Markdown 리포트 생성
12.[사용자 안내]  "해석 완료. 리포트: simulations/{id}/report.md"
```

**수락 기준:**

| AC | 기준 | 검증 방법 |
|----|------|----------|
| AC-1 | 자연어 → 시뮬레이션 완료까지 자동 진행 | E2E 시나리오 수동 테스트 |
| AC-2 | 중간 확인 단계 존재 (AI 추론 결과 확인) | 조건 추론 후 반드시 사용자에게 질문 |
| AC-3 | 전 과정에서 에러 시 한국어 안내 | GenCase 실패, Solver 실패 각각 테스트 |
| AC-4 | 결과 디렉토리에 모든 산출물 | XML, bi4, VTK, CSV, report.md 확인 |

**결과 디렉토리 구조:**
```
simulations/{timestamp}_{case_name}/
├── case_Def.xml           ← AGENT-02
├── case.bi4               ← TOOL-01 (GenCase)
├── case.xml               ← TOOL-01 (GenCase)
├── Part_0000.bi4 ...      ← TOOL-02 (Solver)
├── Run.csv                ← TOOL-02 (Solver)
├── Run.out                ← TOOL-02 (Solver)
├── vtk/
│   └── fluid_0000.vtk ... ← TOOL-03 (PartVTK)
├── csv/
│   ├── measure_Vel.csv    ← TOOL-04 (MeasureTool)
│   └── measure_Press.csv  ← TOOL-04 (MeasureTool)
├── probe_points.txt       ← AGENT-02
└── report.md              ← RPT-01
```

---

## 3. 의존성 그래프

```
AGENT-01 (시스템 프롬프트)         ← 선행 조건 없음, 최우선
    │
    ▼
AGENT-02 (XML 생성)               ← AGENT-01의 파라미터 결정 로직 필요
    │
    ▼
TOOL-01 (GenCase)                  ← XML 파일 필요
    │
    ▼
TOOL-02 (Solver) + JOB-01 (Job Manager)  ← GenCase .bi4 필요, 동시 개발
    │
    ├──▶ TOOL-03 (PartVTK)        ← Solver 완료 후
    ├──▶ TOOL-04 (MeasureTool)    ← Solver 완료 후
    │
    ▼
RPT-01 (리포트)                    ← VTK + CSV 결과 필요
    │
    ▼
TUI-01 (상태 표시)                 ← JOB-01 pubsub 이벤트 필요
    │
    ▼
PIPE-01 (E2E 오케스트레이션)       ← 전체 통합
```

**구현 순서 (Critical Path):**

```
Week 1: AGENT-01 → AGENT-02 → TOOL-01
Week 2: TOOL-02 + JOB-01 (병렬)
Week 3: TOOL-03 + TOOL-04 (병렬) → RPT-01
Week 4: TUI-01 → PIPE-01 (통합 + QA)
```

---

## 4. v0.1 제외 항목 (명시적 Out of Scope)

| 제외 항목 | 이유 | 대상 버전 |
|----------|------|----------|
| STL 파일 입력 (IN-02) | 메시 품질 검증 복잡, 별도 Tool 필요 | v0.3+ |
| 실시간 모니터링 (SIM-05) | v0.1은 완료/실패 감지만, Part 카운팅으로 진행률 추정 | v0.2 |
| pvpython 렌더링 (POST-01~04) | ParaView Docker 통합 미결, VTK만 생성 | v0.2 |
| 이미지/애니메이션 (POST-02, POST-03) | pvpython 의존 | v0.2 |
| AI 물리적 해석 코멘트 (RPT-04) | Qwen3 도메인 성능 검증 선행 필요 | v0.2 |
| 파라메트릭 스터디 (PARA-01~03) | 멀티 Job 관리 + 비교 UI 필요 | v0.3+ |
| 원통형/L형 탱크 (DOM-01) | XML 지오메트리 생성 복잡도 | v0.2/v0.3 |
| 지진파/임의 가진 (DOM-03) | 외부 데이터 파싱 필요 | v0.2 |
| 2D 시뮬레이션 (DOM-04) | XML 구조 차이, 별도 템플릿 | v0.2 |
| Sim Dashboard TUI (전용 뷰) | v0.1은 상태 바 활용, 전용 뷰는 v0.2 | v0.2 |
| Case Wizard TUI | v0.1은 채팅으로 대화형 입력 | v0.2 |
| Result Viewer TUI | v0.1은 리포트 파일 경로 안내 | v0.2 |
| 동시 다중 Job | v0.1은 단일 Job만 | v0.2 |

---

## 5. 기술 제약 및 전제 조건

| 항목 | 전제 |
|------|------|
| Docker | DualSPHysics Docker 이미지 빌드 완료 (`docker compose build` 성공) |
| GPU | NVIDIA RTX 4090 + CUDA 12.6 + nvidia-container-toolkit 설치됨 |
| Ollama | Qwen3 32B 모델 로딩됨 (`ollama run qwen3:32b` 응답) |
| 디스크 | `./simulations/` 경로에 충분한 공간 (케이스당 ~500MB) |
| 네트워크 | Ollama 로컬 API (localhost:11434) 접근 가능 |

---

## 6. DualSPHysics 파라미터 레퍼런스 (v0.1 고정값)

v0.1에서 AI가 결정하는 파라미터와 고정 파라미터를 구분한다.

### AI 결정 파라미터

| 파라미터 | 범위 | 결정 기준 |
|---------|------|----------|
| `dp` | 0.005 ~ 0.05 m | tank_min_dim / 50 |
| `tank L/W/H` | 0.1 ~ 10.0 m | 사용자 입력 또는 표준 치수 |
| `fluid_height` | 0.1×H ~ 0.9×H | 충진율 (기본 50%) |
| `freq` | 0.1 ~ 5.0 Hz | 사용자 입력 또는 공진 근처 |
| `amplitude` | 0.01 ~ 0.5 m | 탱크 길이의 1~10% |
| `TimeMax` | 1.0 ~ 30.0 s | 최소 5/freq 주기 |
| `TimeOut` | 0.01 ~ 1.0 s | TimeMax / 50~100 |

### 고정 파라미터 (v0.1에서 변경 불가)

| 파라미터 | 값 | 비고 |
|---------|-----|------|
| gravity | (0, 0, -9.81) | 지구 중력 |
| rhop0 | 1000 kg/m³ | 물 밀도 |
| gamma | 7 | 물 상태방정식 |
| coefsound | 20 | 음속 계수 |
| coefh | 1.2 | Smoothing length |
| cflnumber | 0.2 | CFL 조건 |
| StepAlgorithm | 2 (Symplectic) | 시간 적분 |
| Kernel | 2 (Wendland) | SPH 커널 |
| ViscoTreatment | 1 (Artificial) | 점성 모델 |
| Visco | 0.01 | 점성 계수 |
| DensityDT | 2 (Fourtakas) | 밀도 확산 |
| DensityDTvalue | 0.1 | DDT 계수 |
| Shifting | 0 (None) | 파티클 시프팅 |
| RhopOutMin/Max | 700 / 1300 | 밀도 범위 |

---

## 7. 용어 사전 (사용자 대면 vs 내부)

v0.1에서 사용자에게 노출되는 모든 텍스트는 왼쪽 열을 사용한다.

| 사용자 대면 (한국어) | 내부 기술 용어 |
|---------------------|---------------|
| 입자 간격 (해석 정밀도) | dp (particle spacing) |
| 해석 준비 | GenCase (pre-processing) |
| 시뮬레이션 실행 | DualSPHysics solver |
| 결과 변환 | PartVTK / MeasureTool |
| 수위 변화 | SWL (Still Water Level) gauge |
| 벽면 힘 | Force gauge |
| 해석이 불안정해졌습니다 | Particle divergence / RhopOut error |
| 해석 정밀도를 높이세요 | dp 감소 권장 |
| 시뮬레이션 시간 | TimeMax |
| 출력 간격 | TimeOut |
| 가진 주파수 | Excitation frequency |
| 공진 주파수 | Natural frequency (f₁) |
