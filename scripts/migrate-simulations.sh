#!/usr/bin/env bash
# migrate-simulations.sh — 시뮬레이션 결과물 폴더 구조 체계화
#
# 구조: simulations/{case_name}/
#   ├── vtk/        ← PartFluid_*.vtk
#   ├── measure/    ← CSV, probe_points.txt
#   ├── viz/        ← PNG, MP4 (최종 시각화)
#   └── (bi4, Run.csv, RunPARTs.csv 등은 루트 유지)
#
# Usage:
#   bash scripts/migrate-simulations.sh --dry-run    # 미리보기
#   bash scripts/migrate-simulations.sh              # 실행

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SIM_DIR="$PROJECT_DIR/simulations"
RESULT_STORE="$PROJECT_DIR/result_store"

DRY_RUN=false
if [[ "${1:-}" == "--dry-run" ]]; then
    DRY_RUN=true
    echo "[DRY-RUN] 변경 사항만 출력합니다."
    echo ""
fi

# Helper: move file (or print in dry-run mode)
move_file() {
    local src="$1"
    local dst_dir="$2"
    if [[ ! -e "$src" ]]; then return; fi
    if $DRY_RUN; then
        echo "  MOVE: $(basename "$src") -> $dst_dir/"
    else
        mkdir -p "$dst_dir"
        mv -n "$src" "$dst_dir/" 2>/dev/null || true
    fi
}

# Helper: move directory contents
move_dir() {
    local src_dir="$1"
    local dst_dir="$2"
    if [[ ! -d "$src_dir" ]]; then return; fi
    local count
    count=$(find "$src_dir" -maxdepth 1 -type f | wc -l)
    if [[ "$count" -eq 0 ]]; then return; fi
    if $DRY_RUN; then
        echo "  MOVE DIR: $src_dir/ ($count files) -> $dst_dir/"
    else
        mkdir -p "$dst_dir"
        find "$src_dir" -maxdepth 1 -type f -exec mv -n {} "$dst_dir/" \; 2>/dev/null || true
    fi
}

# Helper: fix root-owned files in a directory
fix_root_ownership() {
    local dir="$1"
    local root_count
    root_count=$(find "$dir" -maxdepth 1 -user root 2>/dev/null | wc -l || echo 0)
    if [[ "$root_count" -gt 0 ]]; then
        if $DRY_RUN; then
            echo "  FIX OWNERSHIP: $root_count root-owned files in $(basename "$dir")"
        else
            docker run --rm -v "$dir:/fix" alpine chown -R "$(id -u):$(id -g)" /fix 2>/dev/null || true
        fi
    fi
}

echo "=== Step 1: 케이스 디렉토리 내 파일 정리 ==="
echo ""

for case_dir in "$SIM_DIR"/*/; do
    # Skip non-case directories
    case_name="$(basename "$case_dir")"
    if [[ "$case_name" == "output" || "$case_name" == "data" || "$case_name" == "_legacy" ]]; then
        continue
    fi

    echo "[$case_name]"

    # Fix root ownership first
    fix_root_ownership "$case_dir"

    # 1a. PartFluid_*.vtk → vtk/
    vtk_count=$(find "$case_dir" -maxdepth 1 -name 'PartFluid_*.vtk' 2>/dev/null | wc -l)
    if [[ "$vtk_count" -gt 0 ]]; then
        if $DRY_RUN; then
            echo "  MOVE: PartFluid_*.vtk ($vtk_count files) -> vtk/"
        else
            mkdir -p "${case_dir}vtk"
            find "$case_dir" -maxdepth 1 -name 'PartFluid_*.vtk' -exec mv -n {} "${case_dir}vtk/" \;
        fi
    fi

    # 1b. CSV (Run.csv, RunPARTs.csv 제외) → measure/
    for csv in "$case_dir"*.csv; do
        [[ ! -e "$csv" ]] && continue
        base="$(basename "$csv")"
        if [[ "$base" != "Run.csv" && "$base" != "RunPARTs.csv" ]]; then
            move_file "$csv" "${case_dir}measure"
        fi
    done

    # 1c. probe_points.txt, *_probe_points.txt → measure/
    for probe in "$case_dir"*probe_points*.txt "$case_dir"*probes*.txt; do
        [[ ! -e "$probe" ]] && continue
        move_file "$probe" "${case_dir}measure"
    done

    # 1d. standalone *.png, *.mp4 → viz/
    for media in "$case_dir"*.png "$case_dir"*.mp4; do
        [[ ! -e "$media" ]] && continue
        move_file "$media" "${case_dir}viz"
    done

    # 1e. frames/ directory → viz/frames/
    if [[ -d "${case_dir}frames" ]]; then
        if $DRY_RUN; then
            frame_count=$(find "${case_dir}frames" -name '*.png' | wc -l)
            echo "  MOVE DIR: frames/ ($frame_count PNGs) -> viz/frames/"
        else
            mkdir -p "${case_dir}viz"
            mv -n "${case_dir}frames" "${case_dir}viz/frames" 2>/dev/null || true
        fi
    fi

    echo ""
done

echo "=== Step 2: simulations/output/ MP4 → 각 케이스 viz/ ==="
echo ""

OUTPUT_DIR="$SIM_DIR/output"
if [[ -d "$OUTPUT_DIR" ]]; then
    for mp4 in "$OUTPUT_DIR"/*.mp4; do
        [[ ! -e "$mp4" ]] && continue
        mp4_name="$(basename "$mp4" .mp4)"

        # Try to match MP4 name to case directory (case-insensitive, underscore-flexible)
        matched=false
        for case_dir in "$SIM_DIR"/*/; do
            case_name="$(basename "$case_dir")"
            [[ "$case_name" == "output" || "$case_name" == "data" || "$case_name" == "_legacy" ]] && continue

            # Normalize names for comparison (lowercase, strip underscores)
            norm_mp4=$(echo "$mp4_name" | tr '[:upper:]' '[:lower:]' | tr -d '_')
            norm_case=$(echo "$case_name" | tr '[:upper:]' '[:lower:]' | tr -d '_')

            if [[ "$norm_mp4" == *"$norm_case"* || "$norm_case" == *"$norm_mp4"* ]]; then
                move_file "$mp4" "${case_dir}viz"
                matched=true
                break
            fi
        done

        if ! $matched; then
            echo "  SKIP: $mp4_name.mp4 (매칭 케이스 없음 → _legacy로 이동)"
            move_file "$mp4" "$SIM_DIR/_legacy"
        fi
    done

    # Remove output/ if empty
    if ! $DRY_RUN; then
        rmdir "$OUTPUT_DIR" 2>/dev/null || true
    fi
else
    echo "  (output/ 디렉토리 없음)"
fi
echo ""

echo "=== Step 3: result_store/postprocess/ viz 파일 → 각 케이스 viz/ ==="
echo ""

PP_DIR="$RESULT_STORE/postprocess"
if [[ -d "$PP_DIR" ]]; then
    for pp_case in "$PP_DIR"/*/; do
        [[ ! -d "$pp_case" ]] && continue
        pp_name="$(basename "$pp_case")"
        target_case="$SIM_DIR/$pp_name"

        if [[ -d "$target_case" ]]; then
            echo "[$pp_name]"
            for media in "$pp_case"*.png "$pp_case"*.mp4; do
                [[ ! -e "$media" ]] && continue
                move_file "$media" "${target_case}/viz"
            done
        else
            echo "  SKIP: $pp_name (매칭 케이스 디렉토리 없음)"
        fi
    done

    # Remove postprocess/ if empty (recursively remove empty subdirs first)
    if ! $DRY_RUN; then
        find "$PP_DIR" -type d -empty -delete 2>/dev/null || true
    fi
else
    echo "  (postprocess/ 디렉토리 없음)"
fi
echo ""

echo "=== Step 4: result_store/postprocess_v2/ 삭제 ==="
PP_V2="$RESULT_STORE/postprocess_v2"
if [[ -d "$PP_V2" ]]; then
    if $DRY_RUN; then
        echo "  DELETE: $PP_V2/"
    else
        rm -rf "$PP_V2"
    fi
else
    echo "  (postprocess_v2/ 없음)"
fi
echo ""

echo "=== Step 5: simulations/ 루트 정체불명 파일 → _legacy/ ==="
echo ""

# Files at simulations/ root that don't belong to a case
for file in "$SIM_DIR"/*; do
    [[ -d "$file" ]] && continue  # Skip directories
    base="$(basename "$file")"
    move_file "$file" "$SIM_DIR/_legacy"
done
echo ""

echo "=== Step 6: 빈 디렉토리 정리 ==="
if ! $DRY_RUN; then
    find "$SIM_DIR" -type d -empty -delete 2>/dev/null || true
    echo "  완료"
else
    echo "  (dry-run: 건너뜀)"
fi

echo ""
echo "=== 마이그레이션 완료 ==="
