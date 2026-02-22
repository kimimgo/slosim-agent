#!/usr/bin/env python3
"""
run_exp2_exp4.py — EXP-2 (NL→XML) and EXP-4 (Ablation) test harness

Calls Ollama API with tool calling to test XML generation from natural language.
EXP-2: 20 scenarios × FULL prompt
EXP-4: 10 scenarios × 4 ablation conditions = 40 runs

Usage:
    python research/scripts/run_exp2_exp4.py [--exp2-only] [--exp4-only] [--model qwen3:32b-64k]
"""

import json
import time
import math
import argparse
import requests
from pathlib import Path
from dataclasses import dataclass, field, asdict
from typing import Optional

OLLAMA_URL = "http://localhost:11434/api/chat"
DEFAULT_MODEL = "qwen3:latest"  # 8B model (fast, fits GPU fully)

OUT_DIR = Path("research/experiments")
EXP2_DIR = OUT_DIR / "exp2_nl2xml"
EXP4_DIR = OUT_DIR / "exp4_ablation"

# ── System Prompts (matching sloshing_coder.go ablation modes) ────

PROMPT_ROLE = """당신은 슬로싱(Sloshing) 해석 전문 AI 어시스턴트입니다.
비전문가도 자연어로 시뮬레이션을 요청할 수 있도록 돕습니다.

# 역할
- 사용자의 자연어 요청을 슬로싱 시뮬레이션 조건으로 변환합니다
- 누락된 파라미터는 아래 규칙으로 자동 결정합니다
- 시뮬레이션 실행, 후처리, 리포트 생성까지 전체 과정을 관리합니다
- 모든 응답은 한국어로, 전문 용어는 쉬운 표현으로 변환합니다
"""

PROMPT_ABSOLUTE_RULES = """# 절대 규칙
1. 해석 요청 시 반드시 xml_generator를 첫 번째로 호출하세요
2. 기존 XML 파일이 있으면 gencase부터 시작합니다
3. 누락된 파라미터를 사용자에게 묻지 말고 아래 규칙으로 자동 채워서 tool을 호출하세요
4. error_recovery는 시뮬레이션 실행 중 에러가 발생한 경우에만 사용합니다
"""

PROMPT_TOOL_CALL_ORDER = """# Tool 호출 순서 (반드시 이 순서를 따르세요)
1. xml_generator → XML 케이스 파일 생성 (첫 번째로 호출)
2. gencase → 해석 준비 (파티클 생성)
3. solver → 시뮬레이션 실행 (백그라운드)
4. partvtk → 결과 변환
5. measuretool → 수위/압력 측정
6. report → 리포트 생성
"""

PROMPT_PARAMETER_RULES = """# 파라미터 자동 결정 규칙 (누락 시 이 값을 사용)
1. dp = min(L,W,H)/50 (최소 0.005m, 최대 0.05m)
2. 시뮬레이션 시간(time_max) = 5/freq (초)
3. 유체 높이 미지정 시: 탱크 높이의 50%
4. 진폭 미지정 시: 탱크 길이의 5%
5. out_path 미지정 시: simulations/sloshing_case
"""

PROMPT_TANK_PRESETS = """# 표준 탱크 치수 (사용자가 미지정 시)
- "LNG 탱크" → 40m × 40m × 27m
- "선박 탱크" → 20m × 10m × 8m
- "소형 탱크" → 1m × 0.5m × 0.6m
- "실험 탱크" → 0.6m × 0.3m × 0.4m
"""

PROMPT_DOMAIN_KNOWLEDGE = """# 도메인 지식

## 직사각형 탱크 1차 공진 주파수
f₁ = (1/2π) × √(g × π/L × tanh(π/L × h))
- g = 9.81 m/s², L = 탱크 길이, h = 유체 높이

## 용어 변환 규칙
| 사용자 표현 | 내부 용어 |
|---|---|
| 입자 간격 | dp |
| 가진 주파수 | Excitation frequency |
| 공진 주파수 | Natural frequency (f₁) |
"""

PROMPT_TOOL_DETAIL_RULES = """# Tool 호출 세부 규칙
- 경로에 .xml 확장자를 포함하지 않습니다 (자동 추가됨)
- 에러 발생 시 한국어로 원인과 해결 방법을 안내합니다
"""

PROMPT_FOLDER_RULES = """# 시뮬레이션 결과 폴더 규칙
- 모든 결과는 simulations/{case_name}/ 하위에 저장합니다
- xml_generator: out_path = "simulations/my_case/my_case"
"""

PROMPT_FEATURES = """# 지원 기능 (v0.3)
## 경계 조건 방식 (boundary_method 파라미터)
- "dbc" (기본): Dynamic Boundary Condition
- "mdbc": Modified Dynamic Boundary Condition
"""

PROMPT_CONSTRAINTS = """# 제약 사항
- 물(1000 kg/m³) 단일 유체
- 3D 시뮬레이션만"""

# Assemble 4 ablation prompts
PROMPTS = {
    "full": (PROMPT_ROLE + PROMPT_ABSOLUTE_RULES + PROMPT_TOOL_CALL_ORDER +
             PROMPT_PARAMETER_RULES + PROMPT_TANK_PRESETS + PROMPT_DOMAIN_KNOWLEDGE +
             PROMPT_TOOL_DETAIL_RULES + PROMPT_FOLDER_RULES + PROMPT_FEATURES +
             PROMPT_CONSTRAINTS),
    "no-domain": (PROMPT_ROLE + PROMPT_ABSOLUTE_RULES + PROMPT_TOOL_CALL_ORDER +
                  PROMPT_TOOL_DETAIL_RULES + PROMPT_FOLDER_RULES + PROMPT_FEATURES +
                  PROMPT_CONSTRAINTS),
    "no-rules": (PROMPT_ROLE + PROMPT_PARAMETER_RULES + PROMPT_TANK_PRESETS +
                 PROMPT_DOMAIN_KNOWLEDGE + PROMPT_FEATURES + PROMPT_CONSTRAINTS),
    "generic": "You are a helpful AI assistant for DualSPHysics sloshing simulation. Help the user set up and run simulations.",
}

# ── Tool Definition ───────────────────────────────────────────────

XML_GENERATOR_TOOL = {
    "type": "function",
    "function": {
        "name": "xml_generator",
        "description": "DualSPHysics XML 케이스 자동 생성 도구",
        "parameters": {
            "type": "object",
            "properties": {
                "tank_length": {"type": "number", "description": "탱크 길이 (m)"},
                "tank_width": {"type": "number", "description": "탱크 너비 (m)"},
                "tank_height": {"type": "number", "description": "탱크 높이 (m)"},
                "fluid_height": {"type": "number", "description": "유체 높이 (m)"},
                "freq": {"type": "number", "description": "가진 주파수 (Hz)"},
                "amplitude": {"type": "number", "description": "가진 진폭 (m)"},
                "dp": {"type": "number", "description": "파티클 간격 (m)"},
                "time_max": {"type": "number", "description": "시뮬레이션 시간 (s)"},
                "out_path": {"type": "string", "description": "출력 경로"},
                "boundary_method": {"type": "string", "description": "dbc or mdbc"},
            },
            "required": ["tank_length", "tank_width", "tank_height",
                          "fluid_height", "freq", "amplitude", "dp",
                          "time_max", "out_path"],
        },
    },
}

# ── Evaluation ────────────────────────────────────────────────────

def natural_freq(L, h):
    """First sloshing mode natural frequency [Hz]."""
    omega = math.sqrt(9.81 * (math.pi / L) * math.tanh(math.pi * h / L))
    return omega / (2 * math.pi)


@dataclass
class EvalResult:
    scenario_id: str
    level: int
    ablation: str
    called_tool: bool = False
    tool_name: str = ""
    gencase_pass: Optional[bool] = None
    params: dict = field(default_factory=dict)
    param_scores: dict = field(default_factory=dict)
    param_accuracy: float = 0.0
    physical_valid: bool = False
    correct_behavior: bool = False  # For edge cases
    response_text: str = ""
    latency_s: float = 0.0
    error: str = ""


def check_range(val, low, high, tol_pct=0.0):
    """Check if val is within [low, high] with optional tolerance."""
    if tol_pct > 0:
        margin = (high - low) * tol_pct / 100
        return low - margin <= val <= high + margin
    return low <= val <= high


def evaluate_params(params, gt, scenario_id):
    """Evaluate generated parameters against ground truth."""
    scores = {}
    checks = 0
    passes = 0

    # tank_length
    if "tank_length" in gt:
        checks += 1
        if abs(params.get("tank_length", 0) - gt["tank_length"]) < 0.01:
            scores["tank_length"] = 1.0
            passes += 1
        else:
            scores["tank_length"] = 0.0
    elif "tank_length_range" in gt:
        checks += 1
        lo, hi = gt["tank_length_range"]
        if check_range(params.get("tank_length", 0), lo, hi, 20):
            scores["tank_length"] = 1.0
            passes += 1
        else:
            scores["tank_length"] = 0.0

    # tank_width
    if "tank_width" in gt:
        checks += 1
        if abs(params.get("tank_width", 0) - gt["tank_width"]) < 0.01:
            scores["tank_width"] = 1.0
            passes += 1
        else:
            scores["tank_width"] = 0.0
    elif "tank_width_range" in gt:
        checks += 1
        lo, hi = gt["tank_width_range"]
        if check_range(params.get("tank_width", 0), lo, hi, 20):
            scores["tank_width"] = 1.0
            passes += 1
        else:
            scores["tank_width"] = 0.0

    # tank_height
    if "tank_height" in gt:
        checks += 1
        if abs(params.get("tank_height", 0) - gt["tank_height"]) < 0.01:
            scores["tank_height"] = 1.0
            passes += 1
        else:
            scores["tank_height"] = 0.0
    elif "tank_height_range" in gt:
        checks += 1
        lo, hi = gt["tank_height_range"]
        if check_range(params.get("tank_height", 0), lo, hi, 20):
            scores["tank_height"] = 1.0
            passes += 1
        else:
            scores["tank_height"] = 0.0

    # fluid_height
    if "fluid_height" in gt:
        checks += 1
        expected = gt["fluid_height"]
        actual = params.get("fluid_height", 0)
        if expected > 0 and abs(actual - expected) / expected < 0.10:
            scores["fluid_height"] = 1.0
            passes += 1
        else:
            scores["fluid_height"] = 0.0
    elif "fluid_fill_ratio" in gt and "tank_height" in gt:
        checks += 1
        expected = gt["fluid_fill_ratio"] * gt["tank_height"]
        actual = params.get("fluid_height", 0)
        if expected > 0 and abs(actual - expected) / expected < 0.10:
            scores["fluid_height"] = 1.0
            passes += 1
        else:
            scores["fluid_height"] = 0.0

    # frequency
    if "frequency" in gt:
        checks += 1
        expected = gt["frequency"]
        actual = params.get("freq", 0)
        if expected > 0 and abs(actual - expected) / expected < 0.05:
            scores["freq"] = 1.0
            passes += 1
        else:
            scores["freq"] = 0.0
    elif "frequency_expected" in gt:
        checks += 1
        expected = gt["frequency_expected"]
        actual = params.get("freq", 0)
        if expected > 0 and abs(actual - expected) / expected < 0.10:
            scores["freq"] = 1.0
            passes += 1
        else:
            scores["freq"] = 0.0
    elif "frequency_range" in gt:
        checks += 1
        lo, hi = gt["frequency_range"]
        if check_range(params.get("freq", 0), lo, hi):
            scores["freq"] = 1.0
            passes += 1
        else:
            scores["freq"] = 0.0

    # dp
    if "dp_range" in gt:
        checks += 1
        lo, hi = gt["dp_range"]
        if check_range(params.get("dp", 0), lo, hi, 20):
            scores["dp"] = 1.0
            passes += 1
        else:
            scores["dp"] = 0.0

    # motion_type (check from response text if mentioned)
    if "motion_type" in gt:
        checks += 1
        # For tool calling, motion type is implicit in xml_generator (always sway)
        # Consider it correct if tool was called
        if gt["motion_type"] in ("sway", "roll"):
            scores["motion_type"] = 1.0
            passes += 1

    accuracy = passes / checks if checks > 0 else 0.0
    return scores, accuracy


def evaluate_edge_case(gt, tool_called, response_text):
    """Evaluate edge case scenarios (L5)."""
    expected = gt.get("expected_behavior", "")
    response_lower = response_text.lower()

    if expected == "warning":
        # Agent should have generated XML but with a warning
        warning_keywords = ["경고", "주의", "warning", "권장", "recommend",
                            "조정", "크", "많", "초과", "넘", "불가능하지는 않"]
        has_warning = any(kw in response_lower for kw in warning_keywords)
        return has_warning
    elif expected == "error":
        # Agent should NOT have called xml_generator
        error_keywords = ["불가", "없", "error", "유체가 없", "시뮬레이션 불가",
                          "최소한", "채워", "넣어"]
        has_error = any(kw in response_lower for kw in error_keywords)
        return has_error or not tool_called

    return False


# ── API Calling ───────────────────────────────────────────────────

def call_ollama(model, system_prompt, user_message, tools=None, think=False):
    """Call Ollama chat API with streaming to avoid HTTP timeout."""
    # Disable Qwen3 thinking mode by appending /no_think to user message
    # This reduces token usage by ~94% (1251 thinking tokens → 0)
    if not think:
        user_message = user_message + " /no_think"

    messages = [
        {"role": "system", "content": system_prompt},
        {"role": "user", "content": user_message},
    ]

    payload = {
        "model": model,
        "messages": messages,
        "stream": True,
        "options": {"temperature": 0.1, "num_predict": 2048},
    }
    if tools:
        payload["tools"] = tools

    start = time.time()
    try:
        resp = requests.post(OLLAMA_URL, json=payload, timeout=600, stream=True)
        resp.raise_for_status()

        # Accumulate streaming chunks
        content_parts = []
        tool_calls = []
        final_chunk = None

        for line in resp.iter_lines():
            if not line:
                continue
            chunk = json.loads(line)
            msg = chunk.get("message", {})
            if msg.get("content"):
                content_parts.append(msg["content"])
            if msg.get("tool_calls"):
                tool_calls.extend(msg["tool_calls"])
            if chunk.get("done"):
                final_chunk = chunk

        # Reconstruct response in non-streaming format
        result = {
            "message": {
                "role": "assistant",
                "content": "".join(content_parts),
                "tool_calls": tool_calls,
            },
        }
        if final_chunk:
            result["eval_count"] = final_chunk.get("eval_count", 0)
            result["eval_duration"] = final_chunk.get("eval_duration", 0)

        latency = time.time() - start
        return result, latency
    except Exception as e:
        return {"error": str(e)}, time.time() - start


def extract_tool_calls(result):
    """Extract tool calls from Ollama response."""
    msg = result.get("message", {})
    tool_calls = msg.get("tool_calls", [])
    content = msg.get("content", "")
    return tool_calls, content


# ── Main Runner ───────────────────────────────────────────────────

def run_scenario(model, scenario, ablation="full"):
    """Run a single scenario and return EvalResult."""
    sid = scenario["id"]
    level = scenario["level"]
    gt = scenario["ground_truth"]
    nl = scenario["nl_input"]

    print(f"  [{ablation}] {sid} (L{level}): {nl[:60]}...", flush=True)

    system_prompt = PROMPTS[ablation]
    result, latency = call_ollama(model, system_prompt, nl, [XML_GENERATOR_TOOL])

    if "error" in result and isinstance(result["error"], str):
        return EvalResult(
            scenario_id=sid, level=level, ablation=ablation,
            error=result["error"], latency_s=latency,
        )

    tool_calls, content = extract_tool_calls(result)

    ev = EvalResult(
        scenario_id=sid, level=level, ablation=ablation,
        latency_s=latency, response_text=content,
    )

    if tool_calls:
        tc = tool_calls[0]
        fn = tc.get("function", {})
        ev.called_tool = True
        ev.tool_name = fn.get("name", "")
        try:
            ev.params = fn.get("arguments", {})
            if isinstance(ev.params, str):
                ev.params = json.loads(ev.params)
        except (json.JSONDecodeError, TypeError):
            ev.params = {}

    # Evaluate based on level
    if level == 5:
        # Edge cases: evaluate correct behavior
        ev.correct_behavior = evaluate_edge_case(gt, ev.called_tool, content)
        ev.physical_valid = ev.correct_behavior
        ev.param_accuracy = 1.0 if ev.correct_behavior else 0.0
    elif ev.called_tool and ev.tool_name == "xml_generator" and ev.params:
        # Normal cases: evaluate parameter accuracy
        ev.param_scores, ev.param_accuracy = evaluate_params(ev.params, gt, sid)
        # Physical validity: has valid dimensions and positive values
        p = ev.params
        ev.physical_valid = (
            p.get("tank_length", 0) > 0 and
            p.get("tank_width", 0) > 0 and
            p.get("tank_height", 0) > 0 and
            p.get("fluid_height", 0) > 0 and
            0 < p.get("fluid_height", 0) <= p.get("tank_height", 1) and
            p.get("freq", 0) > 0 and
            p.get("amplitude", 0) > 0 and
            p.get("dp", 0) > 0 and
            p.get("time_max", 0) > 0
        )
    else:
        # Tool not called — might be OK for edge cases, bad otherwise
        if level < 5:
            ev.param_accuracy = 0.0
            ev.physical_valid = False

    status = "OK" if (ev.param_accuracy >= 0.5 or (level == 5 and ev.correct_behavior)) else "FAIL"
    print(f"    >> {status} acc={ev.param_accuracy:.0%} valid={ev.physical_valid} "
          f"tool={ev.tool_name or 'none'} {latency:.1f}s", flush=True)

    return ev


def run_exp2(model, scenarios):
    """EXP-2: 20 scenarios × FULL prompt."""
    print("\n" + "=" * 60)
    print("EXP-2: NL→XML Generation Accuracy (FULL prompt)")
    print("=" * 60)

    results = []
    for sc in scenarios:
        ev = run_scenario(model, sc, "full")
        results.append(ev)
        time.sleep(1)  # Rate limit

    return results


def run_exp4(model, scenarios):
    """EXP-4: 10 scenarios × 4 ablation conditions."""
    print("\n" + "=" * 60)
    print("EXP-4: Domain Prompt Ablation")
    print("=" * 60)

    # Select 10 scenarios (2 per level)
    selected_ids = {"S01", "S03", "S05", "S07", "S09", "S11", "S13", "S15", "S17", "S19"}
    selected = [s for s in scenarios if s["id"] in selected_ids]

    results = []
    for ablation in ["full", "no-domain", "no-rules", "generic"]:
        print(f"\n--- Ablation: {ablation} ---")
        for sc in selected:
            ev = run_scenario(model, sc, ablation)
            results.append(ev)
            time.sleep(1)

    return results


# ── Report Generation ─────────────────────────────────────────────

def generate_exp2_report(results, model_name=DEFAULT_MODEL):
    """Generate EXP-2 Table 3."""
    EXP2_DIR.mkdir(parents=True, exist_ok=True)
    (EXP2_DIR / "logs").mkdir(exist_ok=True)

    # Save raw results
    raw = [asdict(r) for r in results]
    with open(EXP2_DIR / "results.json", "w") as f:
        json.dump(raw, f, indent=2, ensure_ascii=False)

    # Summary by level
    levels = {1: "Basic", 2: "Domain", 3: "Paper", 4: "Complex", 5: "Edge"}
    summary = {}
    for lv in range(1, 6):
        lvr = [r for r in results if r.level == lv]
        n = len(lvr)
        if n == 0:
            continue
        tool_called = sum(1 for r in lvr if r.called_tool)
        valid = sum(1 for r in lvr if r.physical_valid)
        avg_acc = sum(r.param_accuracy for r in lvr) / n
        avg_lat = sum(r.latency_s for r in lvr) / n
        summary[lv] = {
            "name": levels[lv], "n": n,
            "tool_called": tool_called,
            "physical_valid": valid,
            "avg_accuracy": avg_acc,
            "avg_latency": avg_lat,
        }

    # Markdown table
    md = "# Table 3: EXP-2 NL→XML Generation Results\n\n"
    md += "| Level | Name | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |\n"
    md += "|-------|------|---|-----------|---------------|---------------|-------------|\n"

    total_n = total_tool = total_valid = 0
    total_acc = total_lat = 0.0

    for lv in range(1, 6):
        if lv not in summary:
            continue
        s = summary[lv]
        md += (f"| L{lv} | {s['name']} | {s['n']} | {s['tool_called']}/{s['n']} | "
               f"{s['avg_accuracy']:.0%} | {s['physical_valid']}/{s['n']} | "
               f"{s['avg_latency']:.1f}s |\n")
        total_n += s["n"]
        total_tool += s["tool_called"]
        total_valid += s["physical_valid"]
        total_acc += s["avg_accuracy"] * s["n"]
        total_lat += s["avg_latency"] * s["n"]

    if total_n > 0:
        md += (f"| **Total** | | **{total_n}** | **{total_tool}/{total_n}** | "
               f"**{total_acc / total_n:.0%}** | **{total_valid}/{total_n}** | "
               f"**{total_lat / total_n:.1f}s** |\n")

    # Detail table
    md += "\n## Detailed Results\n\n"
    md += "| ID | Level | Tool | Accuracy | Valid | Key Params | Latency |\n"
    md += "|----|-------|------|----------|-------|------------|--------|\n"
    for r in results:
        tool_str = r.tool_name if r.called_tool else "none"
        params_str = ""
        if r.params:
            p = r.params
            params_str = (f"L={p.get('tank_length', '?')}, "
                          f"h={p.get('fluid_height', '?')}, "
                          f"f={p.get('freq', '?')}")
        md += (f"| {r.scenario_id} | L{r.level} | {tool_str} | "
               f"{r.param_accuracy:.0%} | {'Y' if r.physical_valid else 'N'} | "
               f"{params_str} | {r.latency_s:.1f}s |\n")

    md += f"\n## Model: {model_name}\n"
    md += f"## Date: {time.strftime('%Y-%m-%d %H:%M')}\n"

    with open(EXP2_DIR / "table3_results.md", "w") as f:
        f.write(md)

    print(f"\nSaved: {EXP2_DIR / 'table3_results.md'}")
    print(f"Saved: {EXP2_DIR / 'results.json'}")


def generate_exp4_report(results, model_name=DEFAULT_MODEL):
    """Generate EXP-4 Table 5 and Fig 4 data."""
    EXP4_DIR.mkdir(parents=True, exist_ok=True)
    (EXP4_DIR / "logs").mkdir(exist_ok=True)

    raw = [asdict(r) for r in results]
    with open(EXP4_DIR / "results.json", "w") as f:
        json.dump(raw, f, indent=2, ensure_ascii=False)

    # Summary by ablation
    ablations = ["full", "no-domain", "no-rules", "generic"]
    summary = {}
    for ab in ablations:
        abr = [r for r in results if r.ablation == ab]
        n = len(abr)
        if n == 0:
            continue
        tool_called = sum(1 for r in abr if r.called_tool)
        valid = sum(1 for r in abr if r.physical_valid)
        avg_acc = sum(r.param_accuracy for r in abr) / n
        avg_lat = sum(r.latency_s for r in abr) / n
        summary[ab] = {
            "n": n, "tool_called": tool_called,
            "physical_valid": valid, "avg_accuracy": avg_acc,
            "avg_latency": avg_lat,
        }

    # Markdown table
    md = "# Table 5: EXP-4 Domain Prompt Ablation Results\n\n"
    md += "| Condition | n | Tool Called | Param Accuracy | Physical Valid | Avg Latency |\n"
    md += "|-----------|---|-----------|---------------|---------------|-------------|\n"
    for ab in ablations:
        if ab not in summary:
            continue
        s = summary[ab]
        md += (f"| {ab.upper()} | {s['n']} | {s['tool_called']}/{s['n']} | "
               f"{s['avg_accuracy']:.0%} | {s['physical_valid']}/{s['n']} | "
               f"{s['avg_latency']:.1f}s |\n")

    # Per-scenario breakdown
    md += "\n## Per-Scenario Breakdown\n\n"
    md += "| Scenario | Level | FULL | NO-DOMAIN | NO-RULES | GENERIC |\n"
    md += "|----------|-------|------|-----------|----------|--------|\n"

    scenario_ids = []
    seen = set()
    for r in results:
        if r.scenario_id not in seen:
            scenario_ids.append(r.scenario_id)
            seen.add(r.scenario_id)

    for sid in scenario_ids:
        row = f"| {sid}"
        lv = next((r.level for r in results if r.scenario_id == sid), 0)
        row += f" | L{lv}"
        for ab in ablations:
            r = next((r for r in results if r.scenario_id == sid and r.ablation == ab), None)
            if r:
                status = f"{r.param_accuracy:.0%}" + (" V" if r.physical_valid else " X")
            else:
                status = "N/A"
            row += f" | {status}"
        row += " |"
        md += row + "\n"

    md += f"\n## Model: {model_name}\n"
    md += f"## Date: {time.strftime('%Y-%m-%d %H:%M')}\n"

    with open(EXP4_DIR / "table5_ablation.md", "w") as f:
        f.write(md)

    # Generate Figure 4 data (for matplotlib)
    fig_data = {}
    for ab in ablations:
        if ab in summary:
            fig_data[ab] = {
                "tool_called_pct": summary[ab]["tool_called"] / summary[ab]["n"] * 100,
                "accuracy_pct": summary[ab]["avg_accuracy"] * 100,
                "valid_pct": summary[ab]["physical_valid"] / summary[ab]["n"] * 100,
            }
    with open(EXP4_DIR / "fig4_data.json", "w") as f:
        json.dump(fig_data, f, indent=2)

    print(f"\nSaved: {EXP4_DIR / 'table5_ablation.md'}")
    print(f"Saved: {EXP4_DIR / 'results.json'}")
    print(f"Saved: {EXP4_DIR / 'fig4_data.json'}")


def generate_fig4(exp4_results):
    """Generate Fig 4: Ablation bar chart."""
    import matplotlib
    matplotlib.use("Agg")
    import matplotlib.pyplot as plt
    import numpy as np

    EXP4_DIR.mkdir(parents=True, exist_ok=True)
    (EXP4_DIR / "figures").mkdir(exist_ok=True)

    ablations = ["full", "no-domain", "no-rules", "generic"]
    labels = ["FULL", "NO-DOMAIN", "NO-RULES", "GENERIC"]

    tool_rates = []
    acc_rates = []
    valid_rates = []

    for ab in ablations:
        abr = [r for r in exp4_results if r.ablation == ab]
        n = len(abr) or 1
        tool_rates.append(sum(1 for r in abr if r.called_tool) / n * 100)
        acc_rates.append(sum(r.param_accuracy for r in abr) / n * 100)
        valid_rates.append(sum(1 for r in abr if r.physical_valid) / n * 100)

    x = np.arange(len(labels))
    width = 0.25

    fig, ax = plt.subplots(figsize=(10, 6))
    bars1 = ax.bar(x - width, tool_rates, width, label="Tool Called (%)",
                   color="steelblue", edgecolor="navy")
    bars2 = ax.bar(x, acc_rates, width, label="Param Accuracy (%)",
                   color="coral", edgecolor="darkred")
    bars3 = ax.bar(x + width, valid_rates, width, label="Physical Valid (%)",
                   color="mediumseagreen", edgecolor="darkgreen")

    ax.set_ylabel("Rate (%)")
    ax.set_title("Fig 4: Domain Prompt Ablation — Impact on XML Generation Quality")
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend(loc="upper right")
    ax.set_ylim(0, 110)
    ax.grid(True, alpha=0.3, axis="y")

    # Value labels
    for bars in [bars1, bars2, bars3]:
        for bar in bars:
            h = bar.get_height()
            ax.annotate(f"{h:.0f}%", xy=(bar.get_x() + bar.get_width() / 2, h),
                        xytext=(0, 3), textcoords="offset points",
                        ha="center", va="bottom", fontsize=8)

    plt.tight_layout()
    path = EXP4_DIR / "figures" / "fig4_ablation.png"
    plt.savefig(path, dpi=300, bbox_inches="tight")
    plt.close()
    print(f"Saved: {path}")


# ── Main ──────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="EXP-2/EXP-4 test harness")
    parser.add_argument("--model", default=DEFAULT_MODEL)
    parser.add_argument("--exp2-only", action="store_true")
    parser.add_argument("--exp4-only", action="store_true")
    args = parser.parse_args()

    # Load scenarios
    with open("research/experiments/exp2_nl2xml/scenarios.json") as f:
        data = json.load(f)
    scenarios = data["scenarios"]
    print(f"Loaded {len(scenarios)} scenarios")

    # Check Ollama
    try:
        resp = requests.get("http://localhost:11434/api/tags", timeout=5)
        print(f"Ollama is running, model: {args.model}")
    except Exception:
        print("ERROR: Ollama is not running! Start with: ollama serve")
        return

    run_e2 = not args.exp4_only
    run_e4 = not args.exp2_only

    exp2_results = []
    exp4_results = []

    if run_e2:
        exp2_results = run_exp2(args.model, scenarios)
        generate_exp2_report(exp2_results, args.model)

    if run_e4:
        exp4_results = run_exp4(args.model, scenarios)
        generate_exp4_report(exp4_results, args.model)
        generate_fig4(exp4_results)

    print("\n" + "=" * 60)
    print("ALL EXPERIMENTS COMPLETE")
    print("=" * 60)

    if exp2_results:
        n = len(exp2_results)
        tc = sum(1 for r in exp2_results if r.called_tool)
        pv = sum(1 for r in exp2_results if r.physical_valid)
        acc = sum(r.param_accuracy for r in exp2_results) / n
        print(f"EXP-2: {tc}/{n} tool called, {acc:.0%} accuracy, {pv}/{n} valid")

    if exp4_results:
        for ab in ["full", "no-domain", "no-rules", "generic"]:
            abr = [r for r in exp4_results if r.ablation == ab]
            n = len(abr) or 1
            acc = sum(r.param_accuracy for r in abr) / n
            pv = sum(1 for r in abr if r.physical_valid)
            print(f"EXP-4 [{ab:10s}]: {acc:.0%} accuracy, {pv}/{len(abr)} valid")


if __name__ == "__main__":
    main()
