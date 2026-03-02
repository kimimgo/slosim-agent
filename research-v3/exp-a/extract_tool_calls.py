#!/usr/bin/env python3
"""
Extract tool call history from OpenCode SQLite DB.
Analyzes M-A1 (Tool Selection) and pipeline completeness.
Usage: python3 extract_tool_calls.py <db_path> [--scenario S01]
"""
import sqlite3
import json
import sys
import os
import re
from pathlib import Path


# Expected tool call sequences for each scenario
EXPECTED_TOOLS = {
    "Easy": ["xml_generator", "gencase", "solver"],
    "Medium": ["xml_generator", "gencase", "solver"],
    "Hard": ["xml_generator|geometry", "gencase", "solver"],  # Hard may use geometry tool
}

# Full pipeline tools (M-A5)
FULL_PIPELINE = ["xml_generator", "gencase", "solver", "partvtk|measuretool"]


def extract_messages(db_path):
    """Extract messages from OpenCode DB."""
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    try:
        cursor.execute("""
            SELECT m.id, m.role, m.parts, m.created_at, s.id as session_id
            FROM messages m
            JOIN sessions s ON m.session_id = s.id
            ORDER BY m.created_at ASC
        """)
        rows = cursor.fetchall()
    except sqlite3.OperationalError as e:
        print(f"DB Error: {e}")
        return []
    finally:
        conn.close()

    return rows


def parse_tool_calls(messages):
    """Parse tool calls from message parts."""
    tool_calls = []

    for msg_id, role, parts_json, created_at, session_id in messages:
        try:
            parts = json.loads(parts_json) if parts_json else []
        except json.JSONDecodeError:
            continue

        for part in parts:
            if isinstance(part, dict):
                # Tool use (assistant calling a tool)
                if part.get("type") == "tool-invocation":
                    tool_calls.append({
                        "timestamp": created_at,
                        "role": role,
                        "tool_name": part.get("toolName", "unknown"),
                        "input": part.get("args", {}),
                        "session_id": session_id,
                    })
                # Tool result
                elif part.get("type") == "tool-result":
                    # Match with previous tool call
                    pass

    return tool_calls


def analyze_tool_sequence(tool_calls, scenario_tier="Easy"):
    """Analyze tool call sequence against expected pattern."""
    tool_names = [tc["tool_name"] for tc in tool_calls]

    analysis = {
        "total_tool_calls": len(tool_calls),
        "unique_tools": list(set(tool_names)),
        "tool_sequence": tool_names,
    }

    # M-A1: Did the agent select the right tools?
    expected = EXPECTED_TOOLS.get(scenario_tier, EXPECTED_TOOLS["Easy"])
    for expected_tool in expected:
        alternatives = expected_tool.split("|")
        found = any(alt in tool_names for alt in alternatives)
        analysis[f"has_{alternatives[0]}"] = found

    # Check for xml_generator specifically
    analysis["ma1_xml_generator"] = "xml_generator" in tool_names
    analysis["ma1_gencase"] = "gencase" in tool_names
    analysis["ma1_solver"] = any(t in tool_names for t in ["solver", "dualsphysics"])

    # M-A1 score: proportion of expected tools used
    expected_found = sum(1 for et in expected
                         if any(alt in tool_names for alt in et.split("|")))
    analysis["ma1_score"] = expected_found / len(expected) if expected else 0

    # M-A5: Full pipeline
    pipeline_found = sum(1 for pt in FULL_PIPELINE
                          if any(alt in tool_names for alt in pt.split("|")))
    analysis["ma5_score"] = pipeline_found / len(FULL_PIPELINE) if FULL_PIPELINE else 0

    return analysis


def analyze_from_artifacts(result_dir):
    """Infer tool usage from simulation artifacts (fallback when DB is unavailable)."""
    result_path = Path(result_dir)
    sim_dir = result_path / "simulations"

    artifacts = {
        "xml_generator": False,
        "gencase": False,
        "solver": False,
        "partvtk": False,
        "measuretool": False,
    }

    # Check for XML
    xmls = list(result_path.glob("*.xml")) + list(sim_dir.glob("*.xml")) if sim_dir.exists() else list(result_path.glob("*.xml"))
    artifacts["xml_generator"] = len(xmls) > 0

    # Check for GenCase outputs (.bi4)
    bi4s = list(sim_dir.glob("**/*.bi4")) if sim_dir.exists() else []
    artifacts["gencase"] = len(bi4s) > 0

    # Check for solver output (Run.out or .out)
    outs = list(sim_dir.glob("**/Run.out")) if sim_dir.exists() else []
    if not outs:
        outs = list(result_path.glob("**/Run.out"))
    artifacts["solver"] = len(outs) > 0

    # Check for VTK files
    vtks = list(sim_dir.glob("**/*PartFluid*.vtk")) if sim_dir.exists() else []
    artifacts["partvtk"] = len(vtks) > 0

    # Check for measurement CSVs
    csvs = list(sim_dir.glob("**/*.csv")) if sim_dir.exists() else []
    artifacts["measuretool"] = len(csvs) > 0

    # Scores
    core_tools = ["xml_generator", "gencase", "solver"]
    ma1_score = sum(1 for t in core_tools if artifacts[t]) / len(core_tools)

    all_tools = ["xml_generator", "gencase", "solver", "partvtk"]
    ma5_score = sum(1 for t in all_tools if artifacts[t]) / len(all_tools)

    return {
        "method": "artifact_inference",
        "artifacts": artifacts,
        "ma1_score": ma1_score,
        "ma5_score": ma5_score,
    }


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 extract_tool_calls.py <db_path|result_dir>")
        sys.exit(1)

    target = sys.argv[1]

    if target.endswith(".db"):
        # DB-based analysis
        messages = extract_messages(target)
        tool_calls = parse_tool_calls(messages)

        print(f"Total messages: {len(messages)}")
        print(f"Total tool calls: {len(tool_calls)}")

        for tc in tool_calls:
            print(f"  [{tc['timestamp']}] {tc['tool_name']}")

        analysis = analyze_tool_sequence(tool_calls)
        print(f"\nM-A1 Score: {analysis['ma1_score']:.0%}")
        print(f"M-A5 Score: {analysis['ma5_score']:.0%}")
    else:
        # Artifact-based analysis (fallback)
        if os.path.isdir(target):
            # Check if it's a single result or batch directory
            subdirs = sorted(Path(target).glob("S*_*/"))
            if subdirs:
                print("=" * 60)
                print("M-A1/M-A5 Tool Usage Analysis (artifact-based)")
                print("=" * 60)
                for subdir in subdirs:
                    result = analyze_from_artifacts(str(subdir))
                    scenario = subdir.name.split("_")[0]
                    print(f"\n{scenario}: M-A1={result['ma1_score']:.0%}, M-A5={result['ma5_score']:.0%}")
                    for tool, found in result["artifacts"].items():
                        icon = "✓" if found else "✗"
                        print(f"  {icon} {tool}")
            else:
                result = analyze_from_artifacts(target)
                print(f"M-A1: {result['ma1_score']:.0%}")
                print(f"M-A5: {result['ma5_score']:.0%}")
                for tool, found in result["artifacts"].items():
                    icon = "✓" if found else "✗"
                    print(f"  {icon} {tool}")


if __name__ == "__main__":
    main()
