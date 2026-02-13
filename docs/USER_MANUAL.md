# slosim-agent 사용자 매뉴얼

> **AI 기반 슬로싱 시뮬레이션 에이전트**
> 자연어로 슬로싱 분석을 실행하는 터미널 도구
>
> GitHub: https://github.com/kimimgo/slosim-agent

---

## 목차

1. [설치 방법](#1-설치-방법)
2. [빠른 시작 가이드](#2-빠른-시작-가이드)
3. [자연어 입력 예시](#3-자연어-입력-예시)
4. [TUI 사용법](#4-tui-사용법)
5. [Tool별 설명](#5-tool별-설명)
6. [파라메트릭 스터디](#6-파라메트릭-스터디)
7. [결과 비교 및 리포트 해석](#7-결과-비교-및-리포트-해석)
8. [트러블슈팅 FAQ](#8-트러블슈팅-faq)

---

## 개요

slosim-agent는 세 가지 핵심 기술을 결합한 슬로싱 시뮬레이션 도구입니다:

| 구성요소 | 역할 |
|----------|------|
| **OpenCode TUI** | Go/BubbleTea 기반 터미널 UI (포크) |
| **Qwen3 SLM** | 로컬 Ollama에서 구동되는 AI 모델 — 자연어 → 시뮬레이션 파라미터 변환 |
| **DualSPHysics v5.4** | GPU 기반 SPH(Smoothed Particle Hydrodynamics) 솔버 |

**시뮬레이션 파이프라인:**

```
자연어 입력 → Qwen3 파라미터 추론 → XML 생성 → GenCase → DualSPHysics GPU
    → PartVTK + MeasureTool → pvpython 시각화 → Markdown 리포트
```

---

## 1. 설치 방법

### 1.1 사전 요구사항

- NVIDIA GPU (CUDA 지원, compute capability ≥ 6.0)
- NVIDIA Driver ≥ 525 + CUDA Toolkit ≥ 12.0
- Docker ≥ 24.0 + NVIDIA Container Toolkit
- Go ≥ 1.22
- Git

### 1.2 Docker (DualSPHysics GPU 솔버)

NVIDIA Container Toolkit이 설치되어 있어야 합니다:

```bash
# NVIDIA Container Toolkit 설치 (Ubuntu)
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | \
  sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg
curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
  sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
  sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

DualSPHysics Docker 이미지 빌드/풀:

```bash
# 프로젝트 클론
git clone https://github.com/kimimgo/slosim-agent.git
cd slosim-agent

# DualSPHysics GPU Docker 이미지 빌드
docker build -t dsph-gpu:5.4 -f docker/Dockerfile.dsph .

# GPU 접근 확인
docker run --rm --gpus all dsph-gpu:5.4 nvidia-smi
```

### 1.3 Go 빌드 (slosim-agent TUI)

```bash
cd slosim-agent

# 빌드
go build -o slosim-agent ./cmd/slosim-agent

# PATH에 추가 (선택)
sudo mv slosim-agent /usr/local/bin/
```

### 1.4 Ollama + Qwen3 설정

```bash
# Ollama 설치
curl -fsSL https://ollama.ai/install.sh | sh

# Qwen3 모델 다운로드
ollama pull qwen3

# Ollama 서버 실행 확인
ollama serve &
curl http://localhost:11434/api/tags  # 모델 목록 확인
```

### 1.5 설정 파일

`~/.config/slosim-agent/config.yaml`:

```yaml
llm:
  provider: ollama
  model: qwen3
  endpoint: http://localhost:11434

solver:
  docker_image: dsph-gpu:5.4
  gpu_device: 0           # GPU 인덱스
  max_particles: 5000000  # 최대 입자 수

output:
  base_dir: ~/slosim-results
  keep_vtk: true          # VTK 파일 보존 여부
```

---

## 2. 빠른 시작 가이드

### 2.1 첫 시뮬레이션 실행

```bash
# slosim-agent 시작
slosim-agent
```

TUI가 열리면 프롬프트에 자연어로 입력합니다:

```
가로 2m, 세로 1m, 높이 1.5m 직사각형 탱크에 물 50% 채우고,
주기 1.2초 진폭 5도로 롤링 운동 10초간 시뮬레이션 해줘
```

### 2.2 실행 흐름

에이전트가 자동으로 다음 단계를 수행합니다:

1. **파라미터 추론** — Qwen3가 자연어에서 시뮬레이션 파라미터를 추출
2. **XML 생성** — DualSPHysics 입력 파일(`.xml`) 자동 생성
3. **GenCase** — XML로부터 입자 배치 생성 (`.bi4`)
4. **DualSPHysics GPU** — SPH 솔버 실행 (Docker 내부, GPU 가속)
5. **후처리** — PartVTK(VTK 변환) + MeasureTool(수위/압력 측정)
6. **시각화** — pvpython으로 스냅샷 및 그래프 생성
7. **리포트** — Markdown 형식 결과 리포트 자동 생성

### 2.3 결과 확인

시뮬레이션 완료 후 TUI에서 바로 리포트를 확인하거나, 출력 디렉토리에서 파일을 확인합니다:

```
~/slosim-results/
└── sim_20260214_070500/
    ├── report.md          # 결과 리포트
    ├── input.xml          # DualSPHysics 입력 파일
    ├── snapshots/         # 시각화 이미지
    ├── data/              # 수위·압력 CSV
    └── vtk/               # VTK 파일 (선택 보존)
```

---

## 3. 자연어 입력 예시

### 3.1 LNG 탱크

```
LNG 멤브레인 탱크 (40m x 40m x 27m), 충전율 70%, 선박 횡동요 주기 8초 진폭 10도, 
20초 시뮬레이션. 벽면 압력 측정 포함.
```

```
Mark III 타입 LNG 탱크에서 충전율 30%일 때 슬로싱 현상 분석해줘.
배의 롤링 주기는 12초, 진폭 15도.
```

### 3.2 직사각형 탱크

```
가로 1m 세로 0.5m 높이 0.8m 사각 탱크, 물 60% 채움.
좌우 흔들림 주기 0.8초 진폭 3cm, 5초간 시뮬레이션.
```

```
2m x 1m x 1.5m 직사각형 수조에서 물이 절반 찼을 때
수평 가진 주파수를 0.5Hz에서 2Hz까지 변화시키며 공진 찾아줘.
```

### 3.3 원통형 탱크

```
직경 3m, 높이 2m 원통형 탱크에 물 40% 채우고
수평 가진 주기 1.5초, 진폭 2cm로 10초 시뮬레이션.
```

### 3.4 지진파 입력

```
직사각형 수조 (4m x 2m x 3m, 물 50%)에 
엘센트로 지진파(El Centro 1940) 적용해서 슬로싱 분석해줘.
```

```
고베 지진 가속도 이력을 사용해서 원통형 저장탱크(직경 10m, 높이 8m, 충전율 60%)의
슬로싱 응답을 분석해줘.
```

### 3.5 파라메트릭 스터디

```
2m x 1m x 1.5m 탱크에서 충전율 20%, 40%, 60%, 80%에 대해
롤링 주기 1.2초 진폭 5도 조건으로 비교 시뮬레이션 해줘.
```

```
사각 탱크 슬로싱에서 가진 주기를 0.8초부터 2.0초까지 0.2초 간격으로 변화시키면서
최대 벽면 압력을 비교해줘.
```

---

## 4. TUI 사용법

### 4.1 화면 구성

```
┌─────────────────────────────────────────────────┐
│  slosim-agent v0.1.0          [Qwen3 ●] [GPU ●] │  ← 상태바
├─────────────────────────────────────────────────┤
│                                                   │
│  [대화 영역]                                      │  ← 자연어 입출력
│  User: 직사각형 탱크 2m x 1m...                   │
│  Agent: 파라미터를 추출했습니다:                   │
│    - 탱크: 2.0 x 1.0 x 1.5m (rectangular)        │
│    - 충전율: 50%                                  │
│    - 운동: Roll, T=1.2s, A=5°                     │
│    ...                                            │
│                                                   │
├─────────────────────────────────────────────────┤
│  ❯ 입력 프롬프트                                  │  ← 자연어 입력
├─────────────────────────────────────────────────┤
│  [Sim Dashboard]  Step: 45000/100000  ETA: 2m30s │  ← 시뮬레이션 현황
│  ████████████████░░░░░░░░  45%  Particles: 250K  │
└─────────────────────────────────────────────────┘
```

### 4.2 키보드 단축키

| 키 | 기능 |
|----|------|
| `Enter` | 입력 전송 |
| `Ctrl+C` | 현재 시뮬레이션 중단 / 프로그램 종료 |
| `Ctrl+L` | 화면 클리어 |
| `Tab` | 자동완성 (명령어/파라미터) |
| `↑` / `↓` | 입력 히스토리 탐색 |
| `Ctrl+D` | Sim Dashboard 토글 |
| `Ctrl+R` | 결과 리포트 보기 |
| `Ctrl+S` | 현재 설정 저장 |
| `Ctrl+O` | 출력 디렉토리 열기 |
| `Esc` | 현재 패널 닫기 |
| `F1` | 도움말 |

### 4.3 상태바 표시

| 아이콘 | 의미 |
|--------|------|
| `[Qwen3 ●]` (초록) | Ollama 연결 정상 |
| `[Qwen3 ○]` (빨강) | Ollama 연결 실패 |
| `[GPU ●]` (초록) | GPU 사용 가능 |
| `[GPU ○]` (빨강) | GPU 미감지 / Docker 오류 |
| `[SIM ▶]` | 시뮬레이션 실행 중 |
| `[SIM ✓]` | 시뮬레이션 완료 |
| `[SIM ✗]` | 시뮬레이션 오류 |

### 4.4 Sim Dashboard

시뮬레이션 실행 중 하단에 표시되는 대시보드:

- **Progress bar** — 진행률 및 현재 스텝
- **ETA** — 예상 잔여 시간
- **Particles** — 총 입자 수
- **dt** — 현재 타임스텝 크기
- **GPU Memory** — VRAM 사용량
- **Errors** — 발산/오류 감지 시 경고

---

## 5. Tool별 설명

slosim-agent는 내부적으로 다음 Tool들을 조합하여 시뮬레이션 파이프라인을 실행합니다.

### 5.1 XML Generator

자연어에서 추출된 파라미터를 DualSPHysics XML 입력 파일로 변환합니다.

- 탱크 형상 (직사각형, 원통형, LNG 멤브레인 등)
- 유체 속성 (밀도, 점성)
- 운동 조건 (롤링, 피칭, 서징, 지진파)
- 수치 파라미터 (입자 간격, 시간 스텝, 커널)
- 측정점 위치 (압력 센서, 수위 게이지)

### 5.2 GenCase

DualSPHysics의 전처리기. XML 파일을 읽어 초기 입자 배치를 생성합니다.

- **입력:** `input.xml`
- **출력:** `input_Actual.xml`, `input.bi4` (바이너리 입자 데이터)
- 탱크 벽면 경계 입자 + 유체 입자 생성

### 5.3 DualSPHysics Solver

GPU 가속 SPH 솔버. Docker 컨테이너 내부에서 `--gpus all`로 실행됩니다.

- **입력:** `input.bi4`, `input_Actual.xml`
- **출력:** `Part_0000.bi4` ~ `Part_NNNN.bi4` (시간별 입자 데이터)
- Symplectic 시간 적분, Wendland 커널, Delta-SPH 보정
- 자동 타임스텝 (CFL 조건 기반)

### 5.4 PartVTK

바이너리 입자 데이터(`.bi4`)를 VTK 형식으로 변환합니다.

- **입력:** `Part_*.bi4`
- **출력:** `PartFluid_*.vtk`, `PartBound_*.vtk`
- 속도, 압력, 밀도 등 물리량 포함

### 5.5 MeasureTool

지정된 위치에서 물리량을 측정합니다.

- **수위 측정** — 특정 x, y 좌표에서의 자유수면 높이 시계열
- **압력 측정** — 벽면 특정 지점에서의 압력 시계열
- **출력:** CSV 파일 (시간, 측정값)

### 5.6 pvpython (ParaView Python)

ParaView의 Python 인터페이스를 사용한 자동 시각화.

- 입자 분포 스냅샷 (시간별)
- 압력 컨투어 이미지
- 수위/압력 시계열 그래프 (matplotlib)
- 애니메이션 프레임 생성

### 5.7 Report

시뮬레이션 결과를 Markdown 리포트로 자동 생성합니다.

- 입력 파라미터 요약
- 주요 결과 (최대 수위, 최대 압력, 공진 주파수 등)
- 시각화 이미지 인라인 포함
- 수치 데이터 테이블

### 5.8 Job Manager

시뮬레이션 작업의 생명주기를 관리합니다.

- 작업 큐잉 및 스케줄링
- 진행 상태 모니터링 (Sim Dashboard에 반영)
- 병렬 실행 관리 (파라메트릭 스터디)
- 작업 이력 관리

### 5.9 Analysis

후처리 데이터에 대한 분석을 수행합니다.

- FFT 분석 (주파수 응답 추출)
- 최대/평균/RMS 값 계산
- 공진 주파수 탐색
- 파라메트릭 스터디 결과 비교

### 5.10 SeismicInput

지진파 가속도 이력 데이터를 처리합니다.

- 내장 지진파: El Centro (1940), Kobe (1995), Northridge (1994) 등
- 사용자 정의 가속도 시계열 CSV 입력 지원
- 필터링, 스케일링, 기준선 보정
- DualSPHysics 외부 운동 입력 형식으로 변환

### 5.11 Animation

시뮬레이션 결과를 동영상으로 생성합니다.

- VTK 프레임 → MP4/GIF 변환
- pvpython 기반 렌더링
- 카메라 앵글 자동 설정
- 시간 스탬프 오버레이

### 5.12 Geometry

탱크 형상을 정의하고 관리합니다.

- 기본 형상: 직사각형, 원통형, 프리즘형
- LNG 탱크 프리셋 (Mark III, NO96, Type-B 등)
- 배플(Baffle) 추가 지원
- 치수 검증 및 입자 수 추정

### 5.13 STL Import

외부 3D 모델(STL 파일)을 탱크 형상으로 가져옵니다.

- STL 파일 읽기 및 검증 (수밀성 확인)
- 좌표계 변환 및 스케일링
- DualSPHysics XML drawfilestl 요소로 변환

### 5.14 ParametricStudy

여러 조건을 자동으로 변화시키며 반복 시뮬레이션을 수행합니다.

- 변수 범위 지정 (충전율, 주기, 진폭 등)
- 조합 생성 (그리드 또는 사용자 정의)
- 병렬/순차 실행 관리
- 결과 자동 비교 테이블 및 그래프

### 5.15 ResultStore

시뮬레이션 결과를 구조화하여 저장·관리합니다.

- 결과 디렉토리 자동 생성 및 정리
- 메타데이터 인덱싱 (검색 가능)
- 이전 결과와 비교 기능
- 디스크 사용량 관리 (VTK 정리 옵션)

### 5.16 ErrorRecovery

시뮬레이션 오류를 자동으로 감지하고 복구를 시도합니다.

- **입자 폭발 (particle explosion)** → 입자 간격 축소 후 재시도
- **타임스텝 발산** → CFL 계수 조정
- **메모리 부족** → 입자 수 감소 제안
- **Docker 오류** → 컨테이너 재시작
- 복구 불가 시 원인 분석 리포트 제공

---

## 6. 파라메트릭 스터디

### 6.1 사용법

자연어로 변수 범위를 지정하면 자동으로 파라메트릭 스터디가 실행됩니다:

```
충전율 20%에서 80%까지 10% 간격으로 변화시키면서
각각 롤링 주기 1.2초, 진폭 5도 조건으로 시뮬레이션 비교해줘.
탱크는 2m x 1m x 1.5m 직사각형.
```

### 6.2 다중 변수

두 개 이상의 변수를 동시에 변화시킬 수 있습니다:

```
충전율 [30%, 50%, 70%] x 가진 주기 [0.8초, 1.0초, 1.2초, 1.5초]
모든 조합에 대해 직사각형 탱크 (2m x 1m x 1.5m) 슬로싱 비교 분석해줘.
```

→ 총 12개 케이스 자동 생성 및 실행

### 6.3 결과 비교

파라메트릭 스터디 완료 후 자동 생성되는 비교 자료:

- **비교 테이블** — 각 케이스별 최대 수위, 최대 압력
- **트렌드 그래프** — 변수 변화에 따른 응답 변화
- **최적/최악 케이스** 자동 식별
- **공진 영역** 하이라이팅

---

## 7. 결과 비교 및 리포트 해석

### 7.1 Markdown 리포트 구조

자동 생성되는 `report.md`의 구조:

```markdown
# 슬로싱 시뮬레이션 리포트
## 1. 시뮬레이션 개요
- 날짜, 소요 시간, 입자 수

## 2. 입력 조건
- 탱크 형상 및 치수
- 유체 속성
- 운동 조건 (주기, 진폭, 유형)
- 수치 파라미터 (dp, 커널, 시간 적분)

## 3. 결과 요약
- 최대 자유수면 높이
- 최대 벽면 압력
- 주요 주파수 성분 (FFT)

## 4. 시각화
- 시간별 입자 분포 스냅샷
- 수위 시계열 그래프
- 압력 시계열 그래프
- 압력 컨투어

## 5. 데이터 테이블
- 수위 측정 데이터 (시간, 높이)
- 압력 측정 데이터 (시간, 압력)

## 6. 분석 및 코멘트
- AI 기반 결과 해석
- 공진 여부 판단
- 안전성 코멘트
```

### 7.2 수위 데이터

`data/elevation.csv` 형식:

```csv
time(s),probe1_h(m),probe2_h(m),probe3_h(m)
0.000,0.750,0.750,0.750
0.010,0.751,0.749,0.750
...
```

- 각 프로브 위치에서의 자유수면 높이 시계열
- 그래프에서 최대값, 공진 응답을 시각적으로 확인 가능

### 7.3 압력 데이터

`data/pressure.csv` 형식:

```csv
time(s),sensor1_P(Pa),sensor2_P(Pa),sensor3_P(Pa)
0.000,7357.5,7357.5,3678.7
0.010,7360.2,7355.1,3680.1
...
```

- 벽면 압력 센서 위치에서의 압력 시계열
- 충격 압력(impact pressure) 피크 식별

### 7.4 이미지

`snapshots/` 디렉토리:

| 파일 | 내용 |
|------|------|
| `snapshot_t0.0s.png` ~ `snapshot_tN.Ns.png` | 시간별 입자 분포 |
| `elevation_plot.png` | 수위 시계열 그래프 |
| `pressure_plot.png` | 압력 시계열 그래프 |
| `fft_plot.png` | 주파수 응답 스펙트럼 |
| `max_pressure_contour.png` | 최대 압력 분포 |

---

## 8. 트러블슈팅 FAQ

### Docker 관련

**Q: `docker: Error response from daemon: could not select device driver`**

Docker NVIDIA 런타임이 설정되지 않았습니다.

```bash
sudo nvidia-ctk runtime configure --runtime=docker
sudo systemctl restart docker
```

**Q: `docker: Error response from daemon: OCI runtime create failed`**

NVIDIA 드라이버와 CUDA 버전 불일치. 드라이버를 업데이트하세요:

```bash
nvidia-smi  # 현재 드라이버 버전 확인
sudo apt install nvidia-driver-535  # 최신 드라이버 설치
```

**Q: Docker 이미지를 찾을 수 없다는 오류**

```bash
docker images | grep dsph-gpu  # 이미지 존재 여부 확인
# 없으면 빌드
cd slosim-agent && docker build -t dsph-gpu:5.4 -f docker/Dockerfile.dsph .
```

### GPU 인식

**Q: `[GPU ○]` 상태바에 GPU가 빨간색으로 표시됨**

```bash
# 1. 호스트에서 GPU 확인
nvidia-smi

# 2. Docker 내부에서 GPU 확인
docker run --rm --gpus all nvidia/cuda:12.0-base nvidia-smi

# 3. NVIDIA Container Toolkit 재설치
sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker
```

**Q: `CUDA error: out of memory`**

입자 수가 GPU VRAM을 초과했습니다. 설정에서 조정하세요:

```yaml
# config.yaml
solver:
  max_particles: 2000000  # 값을 줄이기
```

또는 자연어로: `입자 간격을 0.01m로 크게 해서 다시 실행해줘`

### Ollama 연결

**Q: `[Qwen3 ○]` Ollama 연결 실패**

```bash
# Ollama 프로세스 확인
pgrep -f ollama

# 재시작
ollama serve &

# 모델 확인
ollama list  # qwen3이 목록에 있는지 확인

# 포트 확인
curl http://localhost:11434/api/tags
```

**Q: Qwen3 응답이 너무 느림**

- GPU에서 Ollama가 실행 중인지 확인: `OLLAMA_NUM_GPU=1 ollama serve`
- 더 작은 모델 사용 가능: `ollama pull qwen3:1.7b` (정확도 감소 가능)

### 시뮬레이션 발산

**Q: `Particles out! ... excluded particles`**

입자가 도메인을 벗어났습니다. 원인 및 해결:

- 운동 진폭이 너무 큼 → 진폭을 줄이세요
- 입자 간격(dp)이 너무 큼 → `dp를 절반으로 줄여서 다시 해줘`
- 시뮬레이션 도메인이 작음 → 자동으로 ErrorRecovery가 조정 시도

**Q: `dt=0 or very small timestep detected`**

수치 발산 징후. ErrorRecovery가 자동으로:
1. CFL 계수를 0.2로 축소
2. 인공 점성(artificial viscosity) 증가
3. 재실행 시도

수동 해결: `입자 간격을 더 작게, 인공 점성을 높여서 다시 실행해줘`

**Q: 시뮬레이션은 완료되었지만 결과가 비물리적**

- 입자 간격이 너무 큰 경우 (해상도 부족)
- 시뮬레이션 시간이 너무 짧아 정상 상태 미도달
- 자연어로: `입자 간격 0.005m로 줄이고 시뮬레이션 시간 20초로 늘려서 다시 해줘`

### 기타

**Q: 리포트가 생성되지 않음**

pvpython이 설치되어 있는지 확인:

```bash
docker run --rm dsph-gpu:5.4 pvpython --version
```

**Q: 이전 시뮬레이션 결과를 다시 보고 싶음**

```
이전 시뮬레이션 목록 보여줘
```

또는:

```
sim_20260214_070500 결과 다시 보여줘
```

**Q: STL 파일로 커스텀 탱크를 사용하고 싶음**

```
/path/to/tank.stl 파일을 탱크로 사용해서 물 50% 채우고
롤링 주기 1.5초 진폭 8도로 시뮬레이션 해줘
```

---

## 부록: DualSPHysics 주요 수치 파라미터

| 파라미터 | 기본값 | 설명 |
|----------|--------|------|
| dp (입자 간격) | 자동 계산 | 탱크 크기 기반 자동 설정 (약 최소 치수/50) |
| Kernel | Wendland | SPH 커널 함수 |
| Time Integration | Symplectic | 2차 정확도 시간 적분 |
| Viscosity | Artificial (α=0.01) | 인공 점성 처리 |
| CFL | 0.2 | 쿠랑 수 (타임스텝 제어) |
| Delta-SPH | 0.1 | 밀도 확산 보정 계수 |

대부분의 경우 기본값으로 충분하며, 필요 시 자연어로 조정 가능합니다:
`입자 간격 0.005m, CFL 0.1로 설정해서 실행해줘`

---

*본 매뉴얼은 slosim-agent v0.1.0 기준으로 작성되었습니다.*
*최신 정보는 [GitHub](https://github.com/kimimgo/slosim-agent)를 참고하세요.*
