#!/usr/bin/env bash
# EXP-D: Autonomous Baffle Optimization Experiment Runner
# Usage: bash run_expd.sh [--dry-run]
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
SLOSIM_BINARY="$PROJECT_ROOT/slosim-agent"

# Ollama / LLM settings
export LOCAL_ENDPOINT="${LOCAL_ENDPOINT:-http://localhost:11434}"
export SLOSIM_GENERIC_PROMPT="${SLOSIM_GENERIC_PROMPT:-}"

# Colors
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

log() { echo -e "${GREEN}[EXP-D]${NC} $1"; }
warn() { echo -e "${YELLOW}[EXP-D]${NC} $1"; }
err() { echo -e "${RED}[EXP-D]${NC} $1" >&2; }

# Check prerequisites
check_prereqs() {
    log "Checking prerequisites..."

    if [[ ! -f "$SLOSIM_BINARY" ]]; then
        log "Building slosim-agent..."
        cd "$PROJECT_ROOT"
        go build -o slosim-agent ./main.go
    fi

    if [[ ! -f "$PROJECT_ROOT/cases/fuel_tank.stl" ]]; then
        err "fuel_tank.stl not found in cases/"
        exit 1
    fi

    if ! docker compose -f "$PROJECT_ROOT/docker-compose.yml" ps --services 2>/dev/null | grep -q dsph; then
        warn "DualSPHysics Docker service not running. Starting..."
        docker compose -f "$PROJECT_ROOT/docker-compose.yml" up -d
    fi

    # Check GPU
    if ! nvidia-smi &>/dev/null; then
        err "NVIDIA GPU not available"
        exit 1
    fi

    log "All prerequisites OK"
}

# Create results directories
setup_dirs() {
    mkdir -p "$RESULTS_DIR"/{baseline,iter_1,iter_2,iter_3}
    mkdir -p "$SCRIPT_DIR/analysis"
    log "Results directories created"
}

# NL prompt for the agent
NL_PROMPT="fuel_tank.stl 파일의 자동차 연료탱크에 대해 슬로싱 시뮬레이션을 설정하세요. 조건: 50% 수위, 급정거 시나리오 (x방향 0.5Hz, 50mm 진폭). baffle 없이 baseline 시뮬레이션을 먼저 실행하고, 결과를 분석하여 슬로싱을 최소화하는 baffle 위치를 자율적으로 결정하세요. 3회 이내 반복으로 최적 baffle 배치를 찾아주세요. /no_think"

# Run experiment
run_experiment() {
    local start_time=$(date +%s)
    log "Starting EXP-D experiment at $(date)"
    log "NL Prompt: $NL_PROMPT"

    if [[ "${1:-}" == "--dry-run" ]]; then
        log "[DRY-RUN] Would execute: $SLOSIM_BINARY -p \"$NL_PROMPT\""
        log "[DRY-RUN] Results would be saved to $RESULTS_DIR/"
        return 0
    fi

    # Run slosim-agent in non-interactive mode
    cd "$PROJECT_ROOT"
    "$SLOSIM_BINARY" -p "$NL_PROMPT" 2>&1 | tee "$RESULTS_DIR/agent_log.txt"

    local end_time=$(date +%s)
    local elapsed=$((end_time - start_time))
    log "Experiment completed in ${elapsed}s ($(( elapsed / 60 ))m $(( elapsed % 60 ))s)"
    echo "$elapsed" > "$RESULTS_DIR/wall_clock_seconds.txt"
}

# Post-experiment: collect results
collect_results() {
    log "Collecting results..."

    local sim_dir="/mnt/simdata/dualsphysics"

    # Copy relevant CSVs and XMLs to results/
    for dir in "$sim_dir"/*/; do
        local case_name=$(basename "$dir")
        if [[ "$case_name" == *fuel_tank* ]] || [[ "$case_name" == *baffle* ]]; then
            local target_dir="$RESULTS_DIR"
            # Determine which iteration
            if [[ "$case_name" == *baseline* ]] || [[ "$case_name" == *no_baffle* ]]; then
                target_dir="$RESULTS_DIR/baseline"
            elif [[ "$case_name" == *iter1* ]] || [[ "$case_name" == *iter_1* ]]; then
                target_dir="$RESULTS_DIR/iter_1"
            elif [[ "$case_name" == *iter2* ]] || [[ "$case_name" == *iter_2* ]]; then
                target_dir="$RESULTS_DIR/iter_2"
            elif [[ "$case_name" == *iter3* ]] || [[ "$case_name" == *iter_3* ]]; then
                target_dir="$RESULTS_DIR/iter_3"
            fi

            # Copy CSV/XML files (not VTK/bi4 binaries)
            find "$dir" -maxdepth 1 \( -name "*.csv" -o -name "*.xml" -o -name "Run.csv" \) \
                -exec cp {} "$target_dir/" \; 2>/dev/null || true
            log "Copied results from $case_name to $target_dir"
        fi
    done
}

# Score the experiment
score() {
    log "Running scoring..."
    if [[ -f "$SCRIPT_DIR/score_expd.py" ]]; then
        python3 "$SCRIPT_DIR/score_expd.py" "$RESULTS_DIR"
    else
        warn "score_expd.py not found, skipping scoring"
    fi
}

# Main
main() {
    log "=== EXP-D: Autonomous Baffle Optimization ==="
    check_prereqs
    setup_dirs
    run_experiment "${1:-}"
    collect_results
    score
    log "=== EXP-D Complete ==="
}

main "$@"
