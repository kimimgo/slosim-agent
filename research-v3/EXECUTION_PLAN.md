# Research-v3: 실험 실행 계획

## 인프라 현황 (2026-03-02)

### oliveeelab (즉시 사용 가능)

| 항목 | 상태 |
|------|------|
| GPU | 1× RTX 4090 (24GB, 21GB free) |
| 디스크 | 1.2TB free |
| Docker | v28.5 + `dsph-agent:latest` (2.92GB) |
| Ollama | qwen3:32b (20.2GB), qwen3:latest=8b (5.2GB) |
| slosim-agent | 빌드 가능 (Go 소스 있음) |
| 상태 | **READY** |

### pajulab (셋업 필요)

| 항목 | 상태 | 조치 |
|------|------|------|
| GPU | 4× RTX 4090 (24GB × 4, 전부 idle) | — |
| 디스크 / | **100% FULL** (4.2GB free / 859GB) | Docker 정리 (~226GB 회수) |
| 디스크 /data | 4.6TB free (38%) | 작업 디렉토리로 사용 |
| Docker | v27.5, DSPH 이미지 없음 | 이미지 빌드 또는 전송 |
| Ollama | Qwen3 모델 없음 | `ollama pull qwen3:32b` + `qwen3` |
| slosim-agent | 없음 | 바이너리 전송 또는 빌드 |
| 상태 | **SETUP REQUIRED** |

---

## Phase 0: pajulab 셋업

### Step 0-1: 디스크 확보

```bash
# pajulab에서 실행
# Docker 미사용 이미지/볼륨/컨테이너 정리 (예상 ~226GB 회수)
docker system prune -a --volumes
# 또는 선택적으로:
docker container prune -f        # 중지된 컨테이너 (~4.8GB)
docker volume prune -f           # 미사용 볼륨 (~84.5GB)
docker image prune -a -f         # 미사용 이미지 (~137.6GB)
```

### Step 0-2: Qwen3 모델 설치

```bash
# pajulab에서 실행
ollama pull qwen3:32b           # ~20GB
ollama pull qwen3               # ~5GB (8B)
```

### Step 0-3: DualSPHysics Docker 이미지

옵션 A: oliveeelab에서 전송 (빠름)
```bash
# oliveeelab
docker save dsph-agent:latest | ssh pajulab 'docker load'
```

옵션 B: pajulab에서 직접 빌드
```bash
# pajulab
git clone ... slosim-agent
cd slosim-agent
docker compose build             # CUDA 12.6 + DualSPHysics GPU 빌드
```

### Step 0-4: slosim-agent 바이너리

```bash
# oliveeelab에서 빌드 후 전송
cd ~/workspace/02_active/sim/slosim-agent
CGO_ENABLED=0 go build -o slosim-agent ./main.go
scp slosim-agent pajulab:/data/slosim-agent/
scp -r cases/ pajulab:/data/slosim-agent/cases/
scp -r .opencode/ pajulab:/data/slosim-agent/.opencode/
```

---

## Phase 1: EXP-A 실행 배분

### VRAM 관리 전략

하나의 GPU에서 Ollama(LLM)와 DualSPHysics(Solver)는 시간 분할:
1. Ollama qwen3:32b 로드 (~20GB VRAM) → tool call 생성 → XML 생성
2. `keep_alive=0`으로 모델 언로드 → VRAM 해제
3. DualSPHysics GPU solver 실행 (VRAM 사용량은 파티클 수에 비례)
4. Solver 완료 → Ollama 재로드 → 후처리/분석

→ **1 GPU = 1 시나리오 직렬 실행** (병렬 불가: LLM↔Solver VRAM 교대)

### 실행 배분

| 서버 | GPU 수 | 병렬도 | 담당 |
|------|:---:|:---:|------|
| oliveeelab | 1 | 1 | Qwen3-32B: S1~S5 (Easy+Medium 전반) |
| pajulab GPU0 | 1 | 1 | Qwen3-32B: S6~S10 (Medium 후반+Hard) |
| pajulab GPU1 | 1 | 1 | Qwen3-8B: S1~S10 전체 |
| pajulab GPU2-3 | 2 | 2 | 반복 실행 (32B 2nd/3rd trial) |

### 실행 스크립트

```bash
#!/bin/bash
# run_exp_a.sh — EXP-A 단일 시나리오 실행
SCENARIO=$1    # S1, S2, ...
MODEL=$2       # qwen3:32b or qwen3
TRIAL=$3       # 1, 2, 3
GPU_ID=$4      # 0, 1, 2, 3

export CUDA_VISIBLE_DEVICES=$GPU_ID
export OLLAMA_HOST=http://localhost:11434

OUTDIR="exp-a/${SCENARIO}_${MODEL}_trial${TRIAL}"
mkdir -p "$OUTDIR"

# NL 프롬프트 파일에서 읽기
PROMPT=$(cat "exp-a/prompts/${SCENARIO}.txt")

# 실행 (타임아웃 1시간)
timeout 3600 ./slosim-agent -p "$PROMPT" -q -f json \
  2>&1 | tee "${OUTDIR}/log.json"

# 결과 수집
echo "exit_code=$?" >> "${OUTDIR}/result.txt"
echo "scenario=$SCENARIO" >> "${OUTDIR}/result.txt"
echo "model=$MODEL" >> "${OUTDIR}/result.txt"
echo "trial=$TRIAL" >> "${OUTDIR}/result.txt"
echo "gpu=$GPU_ID" >> "${OUTDIR}/result.txt"
```

### 시간 추정

| 단계 | oliveeelab (1GPU) | pajulab (3-4GPU) | 합계 |
|------|:-:|:-:|:-:|
| 32B × 10 시나리오 × 1st trial | S1-S5 (~3h) | S6-S10 (~3h) | ~3h (병렬) |
| 32B × 10 시나리오 × 2nd trial | — | GPU2-3 (~3h) | ~3h |
| 32B × 10 시나리오 × 3rd trial | S1-S5 (~3h) | S6-S10 (~3h) | ~3h |
| 8B × 10 시나리오 × 1회 | — | GPU1 (~1h, 대부분 조기종료) | ~1h |
| **합계** | ~6h | ~7h | **~7h (wall clock)** |

---

## Phase 2: EXP-B 실행 배분

EXP-A의 B0(Full) 결과를 재사용하므로, B1~B4만 새로 실행.

### 조건별 구현

| 조건 | 구현 방법 | 실행 서버 |
|------|----------|----------|
| B0 Full | EXP-A 결과 재사용 | — |
| B1 −DomainPrompt | 시스템 프롬프트 교체 (코드 분기) | oliveeelab |
| B2 −DSPHTool | DualSPHysics 도구 비활성화 (코드 분기) | oliveeelab |
| B3 −PostProcess | 후처리 도구 비활성화 | pajulab GPU0 |
| B4 Bare LLM | 모든 도구 비활성화 | pajulab GPU1 |

```
4 조건 × 3 시나리오 = 12 실행
예상: ~2h (B2/B4는 GenCase 실패 시 조기종료)
```

---

## Phase 3: EXP-C 이관

v2 데이터 복사만 수행. GPU 불필요.

```bash
# research-v2에서 필요한 데이터 복사
git checkout research-v2 -- \
  research-v2/exp1_spheric/analysis/ \
  research-v2/exp1_spheric/figures/ \
  research-v2/exp1_spheric/scripts/
```

---

## 전체 타임라인

```
Day 1 AM:  Phase 0 — pajulab 셋업 (디스크 정리, 모델 pull, Docker 전송)
Day 1 PM:  Phase 1 — EXP-A 1st trial (oliveeelab S1-S5 + pajulab S6-S10)
Day 1 EVE: Phase 1 — EXP-A 8B + 2nd trial (pajulab 병렬)
Day 2 AM:  Phase 1 — EXP-A 3rd trial + 결과 정리
Day 2 PM:  Phase 2 — EXP-B ablation (코드 분기 구현 + 실행)
Day 2 EVE: Phase 3 — EXP-C 이관 + 전체 결과 Table/Figure 생성
Day 3:     논문 작성 시작
```

## 사전 준비 체크리스트

- [ ] pajulab 디스크 정리 (docker system prune)
- [ ] pajulab Qwen3 모델 pull (32b + 8b)
- [ ] pajulab DualSPHysics Docker 이미지 전송
- [ ] pajulab slosim-agent 바이너리 + cases 전송
- [ ] oliveeelab slosim-agent 최신 빌드
- [ ] EXP-A NL 프롬프트 10개 파일 생성
- [ ] EXP-A 실행 스크립트 작성
- [ ] EXP-B 코드 분기 구현 (B1~B4 조건별)
- [ ] 결과 수집/분석 스크립트 작성
