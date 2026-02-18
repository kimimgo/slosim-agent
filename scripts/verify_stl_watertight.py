#!/usr/bin/env python3
"""Verify an ASCII STL file is watertight.

Reads an STL file, counts triangles/vertices/edges,
and checks that every edge is shared by exactly 2 triangles.

Usage: python3 verify_stl_watertight.py <stl_file>
"""

import re
import sys


def parse_ascii_stl(filepath):
    """Parse an ASCII STL file and return list of triangles [(v0,v1,v2), ...]."""
    triangles = []
    current_verts = []

    with open(filepath, "r") as f:
        for line in f:
            line = line.strip()
            if line.startswith("vertex"):
                parts = line.split()
                v = (float(parts[1]), float(parts[2]), float(parts[3]))
                current_verts.append(v)
                if len(current_verts) == 3:
                    triangles.append(tuple(current_verts))
                    current_verts = []

    return triangles


def round_vertex(v, decimals=8):
    """Round vertex coordinates to avoid floating-point comparison issues."""
    return tuple(round(c, decimals) for c in v)


def verify(triangles):
    """Verify watertightness and return (pass, stats_dict)."""
    edge_count = {}
    vertex_set = set()

    for tri in triangles:
        verts = [round_vertex(v) for v in tri]
        vertex_set.update(verts)
        for i in range(3):
            a, b = verts[i], verts[(i + 1) % 3]
            edge = tuple(sorted([a, b]))
            edge_count[edge] = edge_count.get(edge, 0) + 1

    n_tri = len(triangles)
    n_vert = len(vertex_set)
    n_edge = len(edge_count)
    euler = n_vert - n_edge + n_tri

    non_manifold = {e: c for e, c in edge_count.items() if c != 2}

    stats = {
        "triangles": n_tri,
        "vertices": n_vert,
        "edges": n_edge,
        "euler": euler,
        "non_manifold_edges": len(non_manifold),
        "non_manifold_details": non_manifold,
    }

    return len(non_manifold) == 0, stats


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 verify_stl_watertight.py <stl_file>")
        sys.exit(1)

    filepath = sys.argv[1]
    print(f"Reading STL: {filepath}")

    triangles = parse_ascii_stl(filepath)
    if not triangles:
        print("ERROR: No triangles found in file")
        sys.exit(1)

    passed, stats = verify(triangles)

    print(f"\nTriangles: {stats['triangles']}")
    print(f"Vertices:  {stats['vertices']}")
    print(f"Edges:     {stats['edges']}")
    print(f"Euler:     V - E + F = {stats['euler']} (expected 2 for closed surface)")

    if passed:
        print("\nPASS: Mesh is watertight (all edges shared by exactly 2 triangles)")
    else:
        print(f"\nFAIL: {stats['non_manifold_edges']} non-manifold edges:")
        for edge, count in stats["non_manifold_details"].items():
            print(f"  {edge} -> {count} triangles")
        sys.exit(1)


if __name__ == "__main__":
    main()
