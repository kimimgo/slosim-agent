#!/usr/bin/env python3
"""Ollama API로 LLM 직접 호출 — B2/B4 공용"""
import json
import subprocess
import sys
import re
from pathlib import Path


def main():
    if len(sys.argv) < 6:
        print("Usage: ollama_generate.py <model> <system_prompt_file_or_text> <nl_prompt_file> <outdir> <ollama_url>")
        sys.exit(1)

    model = sys.argv[1]
    system_source = sys.argv[2]
    nl_prompt_file = sys.argv[3]
    outdir = sys.argv[4]
    ollama_url = sys.argv[5]

    # System prompt: file path or inline text
    if Path(system_source).is_file():
        system_prompt = Path(system_source).read_text()
    else:
        system_prompt = system_source

    # NL prompt
    nl_prompt = Path(nl_prompt_file).read_text().strip()
    full_prompt = nl_prompt + "\n\nGenerate a complete DualSPHysics v5.4 XML case file for this simulation. Output ONLY the raw XML."

    payload = json.dumps({
        "model": model,
        "system": system_prompt,
        "prompt": full_prompt,
        "stream": False,
        "options": {
            "temperature": 0,
            "num_predict": 8192
        }
    })

    outdir = Path(outdir)
    outdir.mkdir(parents=True, exist_ok=True)

    try:
        result = subprocess.run(
            ["curl", "-s", f"{ollama_url}/api/generate", "-d", payload],
            capture_output=True, text=True, timeout=600
        )

        (outdir / "response.json").write_text(result.stdout)
        if result.stderr:
            (outdir / "stderr.log").write_text(result.stderr)

        print(f"  Response: {len(result.stdout)} bytes")

        # Extract XML
        data = json.loads(result.stdout)
        resp = data.get("response", "")

        xml_match = re.search(r'(<\?xml.*?</case>)', resp, re.DOTALL)
        if xml_match:
            xml = xml_match.group(1)
        elif '<case' in resp:
            xml_match = re.search(r'(<case.*?</case>)', resp, re.DOTALL)
            xml = xml_match.group(1) if xml_match else resp
        else:
            xml = resp

        (outdir / "generated.xml").write_text(xml)
        print(f"  XML extracted: {len(xml)} chars")

    except Exception as e:
        print(f"  ERROR: {e}", file=sys.stderr)
        (outdir / "stderr.log").write_text(str(e))


if __name__ == "__main__":
    main()
