#!/usr/bin/env python3
"""EXP-B M-A3 채점기 — B0/B1/B2/B4 2×2 Factorial Ablation"""
import json
import re
import xml.etree.ElementTree as ET
from pathlib import Path
import sys

GT_PATH = Path(__file__).parent.parent / "exp-a" / "ground_truth.json"
RESULTS_DIR = Path(__file__).parent / "results"

with open(GT_PATH) as f:
    GROUND_TRUTH = json.load(f)


def parse_xml_params(xml_path):
    """DualSPHysics XML에서 주요 파라미터 추출"""
    try:
        tree = ET.parse(xml_path)
        root = tree.getroot()
    except Exception as e:
        return None, str(e)

    params = {}

    # dp
    for dp_el in root.iter('dp'):
        params['dp'] = float(dp_el.get('value', 0))
    if not params.get('dp'):
        defn = root.find('.//definition')
        if defn is not None:
            params['dp'] = float(defn.get('dp', 0))

    # rhop0
    for rho in root.iter('rhop0'):
        params['rhop0'] = float(rho.get('value', 0))

    # Execution parameters
    for p in root.iter('parameter'):
        if p.get('key') == 'TimeMax':
            params['timemax'] = float(p.get('value', 0))
        if p.get('key') == 'BoundaryMethod':
            params['boundary_method'] = int(p.get('value', 0))
        if p.get('key') == 'Visco':
            params['visco'] = float(p.get('value', 0))

    # Geometry
    boxes = []
    for db in root.iter('drawbox'):
        size_el = db.find('size')
        if size_el is not None:
            boxes.append({
                'x': float(size_el.get('x', 0)),
                'y': float(size_el.get('y', 0)),
                'z': float(size_el.get('z', 0))
            })
    if boxes:
        if len(boxes) >= 2:
            params['tank'] = boxes[1]
            params['fluid_box'] = boxes[0]
        elif len(boxes) == 1:
            params['tank'] = boxes[0]

    # Motion
    for motion_tag in ['mvrectsinu', 'mvrotsinu', 'mvrectunif']:
        for mot in root.iter(motion_tag):
            params['motion_type'] = motion_tag
            freq_el = mot.find('freq')
            ampl_el = mot.find('ampl')
            if freq_el is not None:
                params['freq_x'] = float(freq_el.get('x', 0))
                params['freq_y'] = float(freq_el.get('y', 0))
            if ampl_el is not None:
                params['ampl_x'] = float(ampl_el.get('x', 0))
                params['ampl_y'] = float(ampl_el.get('y', 0))
                params['ampl_z'] = float(ampl_el.get('z', 0))
            params['duration'] = float(mot.get('duration', 0))
            break

    for cyl in root.iter('drawcylinder'):
        params['has_cylinder'] = True

    # STL import detection
    for stl_tag in ['drawfilestl', 'drawstl']:
        for el in root.iter(stl_tag):
            params['has_stl'] = True
            params['stl_file'] = el.get('file', '')
            break

    return params, None


def score_scenario(scenario_id, params):
    """Ground truth와 비교하여 M-A3 점수 계산"""
    gt = GROUND_TRUTH.get(scenario_id, {})
    if not gt:
        return [], 0, 0

    checks = []

    # Tank dimensions
    if 'x' in gt.get('tank', {}):
        tank_gt = gt['tank']
        tank_act = params.get('tank', {})
        for dim in ['x', 'y', 'z']:
            if dim in tank_gt:
                expected = tank_gt[dim]
                actual = tank_act.get(dim, 0)
                tol = expected * 0.10  # 10% tolerance for bare LLM
                passed = abs(actual - expected) <= tol
                checks.append({
                    'name': f'tank_{dim}',
                    'expected': expected,
                    'actual': actual,
                    'pass': passed
                })
    elif gt.get('tank', {}).get('shape') == 'stl':
        has_stl = params.get('has_stl', False)
        checks.append({
            'name': 'stl_import_used',
            'expected': True,
            'actual': has_stl,
            'pass': has_stl
        })
        stl_file = params.get('stl_file', '')
        checks.append({
            'name': 'stl_file_referenced',
            'expected': 'stl file present',
            'actual': stl_file if stl_file else 'none',
            'pass': bool(stl_file)
        })
    elif gt.get('tank', {}).get('shape') in ('cylinder', 'horizontal_cylinder'):
        has_cyl = params.get('has_cylinder', False)
        checks.append({
            'name': 'geometry_type',
            'expected': 'cylinder',
            'actual': 'cylinder' if has_cyl else 'box',
            'pass': has_cyl
        })

    # Fill height
    fill_h = gt.get('fluid', {}).get('fill_height')
    if fill_h:
        fluid_box = params.get('fluid_box', {})
        actual_fill = fluid_box.get('z', 0) if fluid_box else 0
        tol = fill_h * 0.15
        passed = abs(actual_fill - fill_h) <= tol
        checks.append({
            'name': 'fill_height',
            'expected': fill_h,
            'actual': actual_fill,
            'pass': passed
        })

    # Motion type
    motion_gt = gt.get('motion', {})
    if motion_gt.get('xml_tag'):
        actual_motion = params.get('motion_type', 'none')
        passed = actual_motion == motion_gt['xml_tag']
        checks.append({
            'name': 'motion_type',
            'expected': motion_gt['xml_tag'],
            'actual': actual_motion,
            'pass': passed
        })

    # Frequency
    if motion_gt.get('freq_hz'):
        actual_freq = params.get('freq_x', 0) or params.get('freq_y', 0)
        expected_freq = motion_gt['freq_hz']
        tol = expected_freq * 0.10
        passed = abs(actual_freq - expected_freq) <= tol
        checks.append({
            'name': 'frequency',
            'expected': expected_freq,
            'actual': actual_freq,
            'pass': passed
        })

    # Amplitude
    if motion_gt.get('amplitude_m'):
        actual_ampl = params.get('ampl_x', 0) or params.get('ampl_y', 0)
        expected_ampl = motion_gt['amplitude_m']
        tol = expected_ampl * 0.15
        passed = abs(actual_ampl - expected_ampl) <= tol
        checks.append({
            'name': 'amplitude',
            'expected': expected_ampl,
            'actual': actual_ampl,
            'pass': passed
        })
    elif motion_gt.get('amplitude_deg'):
        actual_motion = params.get('motion_type', 'none')
        if actual_motion == 'mvrotsinu':
            actual_ampl = params.get('ampl_x', 0) or params.get('ampl_y', 0)
            exp_deg = motion_gt['amplitude_deg']
            if isinstance(exp_deg, list):
                exp_deg = exp_deg[0]
            tol = exp_deg * 0.15
            passed = abs(actual_ampl - exp_deg) <= tol
        else:
            passed = False
            actual_ampl = params.get('ampl_x', 0)
        checks.append({
            'name': 'amplitude',
            'expected': motion_gt['amplitude_deg'],
            'actual': actual_ampl,
            'pass': passed
        })

    # TimeMax
    if gt.get('timemax'):
        actual_tm = params.get('timemax', 0)
        expected_tm = gt['timemax']
        tol = expected_tm * 0.10
        passed = abs(actual_tm - expected_tm) <= tol
        checks.append({
            'name': 'timemax',
            'expected': expected_tm,
            'actual': actual_tm,
            'pass': passed
        })

    passed_count = sum(1 for c in checks if c['pass'])
    total = len(checks)
    return checks, passed_count, total


def check_xml_validity(xml_path):
    """M-A2: XML이 파서블한지 + <case> 루트 존재 확인"""
    try:
        tree = ET.parse(xml_path)
        root = tree.getroot()
        has_case = root.tag == 'case' or root.find('.//case') is not None
        has_casedef = root.find('.//casedef') is not None
        has_execution = root.find('.//execution') is not None
        return {
            'parseable': True,
            'has_case': has_case,
            'has_casedef': has_casedef,
            'has_execution': has_execution,
            'score': sum([has_case, has_casedef, has_execution]) / 3 * 100
        }
    except Exception as e:
        return {
            'parseable': False,
            'error': str(e),
            'score': 0
        }


def main():
    conditions = ['B0', 'B1', 'B2', 'B4']
    scenarios = ['S01', 'S02', 'S03', 'S04', 'S05', 'S06', 'S07', 'S08', 'S09', 'S10']
    models = ['qwen3_32b', 'qwen3_latest']

    print("=" * 70)
    print("  EXP-B Ablation Study — 2×2 Factorial Results")
    print("=" * 70)
    print()
    print("  | Condition | Domain Prompt | Tools |")
    print("  |-----------|:---:|:---:|")
    print("  | B0 Full   |  ✓  |  ✓  |")
    print("  | B1 -Prompt|  ✗  |  ✓  |")
    print("  | B2 -Tool  |  ✓  |  ✗  |")
    print("  | B4 Bare   |  ✗  |  ✗  |")
    print()

    all_results = {}

    for model in models:
        print(f"\n{'─' * 60}")
        print(f"  Model: {model}")
        print(f"{'─' * 60}")

        for cond in conditions:
            for scenario in scenarios:
                # B0 = EXP-A results (look in exp-a results)
                if cond == 'B0':
                    result_dir = Path(__file__).parent.parent / "exp-a" / "results"
                    xml_path = result_dir / f"{scenario}_{model}_trial1" / "simulations" / "sloshing_case.xml"
                    if not xml_path.exists():
                        # Try alternative names in simulations/
                        alt_names = ['parametric_case.xml', 'ISOPE_LNG_Benchmark.xml', 'spheric_benchmark.xml',
                                     'fuel_tank.xml', 'fuel_tank_analysis.xml', 'case.xml', 'base.xml']
                        for alt in alt_names:
                            alt_path = result_dir / f"{scenario}_{model}_trial1" / "simulations" / alt
                            if alt_path.exists():
                                xml_path = alt_path
                                break
                    if not xml_path.exists():
                        # Fallback: check root level (e.g. S05_8B puts XML outside simulations/)
                        root_dir = result_dir / f"{scenario}_{model}_trial1"
                        for alt in ['generated.xml', 'sloshing_case.xml', 'case.xml', 'base.xml',
                                    'fuel_tank.xml', 'parametric_case.xml']:
                            alt_path = root_dir / alt
                            if alt_path.exists():
                                xml_path = alt_path
                                break
                else:
                    result_dir = RESULTS_DIR / f"{cond}_{scenario}_{model}"
                    xml_path = result_dir / "generated.xml"
                    if not xml_path.exists():
                        # Try simulations directory for B1
                        xml_path = result_dir / "simulations" / "sloshing_case.xml"

                key = f"{cond}_{scenario}_{model}"

                if not xml_path.exists():
                    print(f"\n  {cond} {scenario}: NO XML")
                    all_results[key] = {'ma2': 0, 'ma3': 0, 'status': 'NO_XML'}
                    continue

                # M-A2: XML Validity
                validity = check_xml_validity(str(xml_path))

                # M-A3: Parameter Fidelity
                if validity['parseable']:
                    params, err = parse_xml_params(str(xml_path))
                    if params:
                        checks, passed, total = score_scenario(scenario, params)
                        pct = (passed / total * 100) if total > 0 else 0
                    else:
                        checks, pct = [], 0
                else:
                    checks, pct = [], 0

                # Read metadata for duration
                meta_path = (result_dir if cond != 'B0'
                             else result_dir / f"{scenario}_{model}_trial1")
                meta_file = meta_path / "metadata.json"
                duration = "?"
                if meta_file.exists():
                    try:
                        meta = json.load(open(meta_file))
                        duration = meta.get('duration_seconds', '?')
                    except:
                        pass

                status = '✓' if pct >= 80 else '△' if pct >= 50 else '✗'
                print(f"\n  {cond} {scenario}: {status} M-A2={validity['score']:.0f}% M-A3={pct:.0f}% ({duration}s)")
                for c in checks:
                    mark = 'PASS' if c['pass'] else 'FAIL'
                    print(f"    [{mark}] {c['name']}: expected={c['expected']}, actual={c['actual']}")

                all_results[key] = {
                    'ma2': validity['score'],
                    'ma3': pct,
                    'duration': duration,
                    'status': 'OK'
                }

    # Summary table
    print(f"\n\n{'=' * 70}")
    print("  SUMMARY TABLE — M-A3 Parameter Fidelity (%)")
    print(f"{'=' * 70}")

    # Tier labels for scenarios
    tier_labels = {s: GROUND_TRUTH.get(s, {}).get('tier', '?') for s in scenarios}
    tier_abbr = {'Easy': 'E', 'Medium': 'M', 'Hard': 'H'}

    for model in models:
        print(f"\n  {model}:")
        header = f"  {'Cond':<6}"
        for s in scenarios:
            t = tier_abbr.get(tier_labels[s], '?')
            header += f" {s}({t})"
        header += f" {'Mean':>7}"
        print(header)
        print(f"  {'─' * (6 + 7 * len(scenarios) + 8)}")

        for cond in conditions:
            scores = []
            row = f"  {cond:<6}"
            for s in scenarios:
                key = f"{cond}_{s}_{model}"
                r = all_results.get(key, {})
                if r.get('status') in ('OK', 'NO_XML'):
                    row += f" {r['ma3']:>5.0f}%"
                    scores.append(r['ma3'])
                else:
                    row += f" {'N/A':>6}"
            if scores:
                mean = sum(scores) / len(scores)
                row += f" {mean:>6.1f}%"
            else:
                row += f" {'—':>7}"
            print(row)

    # 2×2 Factorial Analysis
    print(f"\n\n{'=' * 70}")
    print("  2×2 FACTORIAL ANALYSIS")
    print(f"{'=' * 70}")

    for model in models:
        print(f"\n  {model}:")
        vals = {}
        for cond in conditions:
            scores = []
            for s in scenarios:
                key = f"{cond}_{s}_{model}"
                r = all_results.get(key, {})
                if r.get('status') in ('OK', 'NO_XML'):
                    scores.append(r['ma3'])
            if scores:
                vals[cond] = sum(scores) / len(scores)

        if all(c in vals for c in conditions):
            print(f"  B0 (Prompt+Tools) = {vals['B0']:.1f}%")
            print(f"  B1 (Tools only)   = {vals['B1']:.1f}%")
            print(f"  B2 (Prompt only)  = {vals['B2']:.1f}%")
            print(f"  B4 (Neither)      = {vals['B4']:.1f}%")
            print()
            prompt_effect = (vals['B0'] + vals['B2']) / 2 - (vals['B1'] + vals['B4']) / 2
            tool_effect = (vals['B0'] + vals['B1']) / 2 - (vals['B2'] + vals['B4']) / 2
            interaction = vals['B0'] - vals['B1'] - vals['B2'] + vals['B4']
            print(f"  Main effect of Prompt: {prompt_effect:+.1f}%")
            print(f"  Main effect of Tools:  {tool_effect:+.1f}%")
            print(f"  Interaction:           {interaction:+.1f}%")
        else:
            available = {k: v for k, v in vals.items()}
            print(f"  Available: {available}")
            print(f"  (Insufficient data for factorial analysis)")


if __name__ == '__main__':
    main()
