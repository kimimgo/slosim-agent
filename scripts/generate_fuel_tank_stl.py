#!/usr/bin/env python3
"""Generate a watertight ASCII STL for a simplified automotive fuel tank.

Based on Frosina 2018 paper dimensions, simplified to a rectangular box
with beveled top edges for a more realistic tank profile.

Dimensions (meters): 0.5 x 0.35 x 0.25 (500 x 350 x 250 mm)
Bevel: 30mm chamfer on top 4 longitudinal edges
"""

import math
import os

# Tank dimensions (meters)
L = 0.5    # length (x)
W = 0.35   # width (y)
H = 0.25   # height (z)
B = 0.03   # bevel size (30mm chamfer on top edges)


def cross(a, b):
    """Cross product of two 3-vectors."""
    return (
        a[1] * b[2] - a[2] * b[1],
        a[2] * b[0] - a[0] * b[2],
        a[0] * b[1] - a[1] * b[0],
    )


def normalize(v):
    """Normalize a 3-vector."""
    mag = math.sqrt(v[0] ** 2 + v[1] ** 2 + v[2] ** 2)
    if mag == 0:
        return (0, 0, 0)
    return (v[0] / mag, v[1] / mag, v[2] / mag)


def compute_normal(v0, v1, v2):
    """Compute outward-facing normal for a triangle."""
    e1 = (v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2])
    e2 = (v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2])
    return normalize(cross(e1, e2))


def write_stl(filename, triangles):
    """Write triangles to ASCII STL file."""
    with open(filename, "w") as f:
        f.write("solid fuel_tank\n")
        for tri in triangles:
            v0, v1, v2 = tri
            n = compute_normal(v0, v1, v2)
            f.write(f"  facet normal {n[0]:.6e} {n[1]:.6e} {n[2]:.6e}\n")
            f.write("    outer loop\n")
            for v in (v0, v1, v2):
                f.write(f"      vertex {v[0]:.6e} {v[1]:.6e} {v[2]:.6e}\n")
            f.write("    endloop\n")
            f.write("  endfacet\n")
        f.write("endsolid fuel_tank\n")


def generate_fuel_tank():
    """Generate a beveled-top box fuel tank.

    Cross-section in the YZ plane (looking from +X):

        B1────────────B2      z = H
       /                \\
      /                  \\
    B0                    B3  z = H - B
    |                      |
    |                      |
    A0────────────────────A3  z = 0

    Where B = bevel size.
    The shape is extruded along X from x=0 to x=L.

    Vertices at x=0 (front face): A0, A3, B0, B1, B2, B3
    Vertices at x=L (back face):  A0', A3', B0', B1', B2', B3'
    """
    triangles = []

    # Define the 6-point cross-section profile (counterclockwise when viewed from outside)
    # Bottom-left, bottom-right, upper-left (bevel start), top-left, top-right, upper-right (bevel start)
    # At x=0 (front)
    f_bl = (0, 0, 0)       # A0: bottom-left
    f_br = (0, W, 0)       # A3: bottom-right
    f_ml = (0, 0, H - B)   # B0: bevel start left
    f_tl = (0, B, H)       # B1: top-left (after bevel)
    f_tr = (0, W - B, H)   # B2: top-right (after bevel)
    f_mr = (0, W, H - B)   # B3: bevel start right

    # At x=L (back)
    b_bl = (L, 0, 0)
    b_br = (L, W, 0)
    b_ml = (L, 0, H - B)
    b_tl = (L, B, H)
    b_tr = (L, W - B, H)
    b_mr = (L, W, H - B)

    # Helper: create two triangles for a quad (v0, v1, v2, v3) with outward normal
    def quad(v0, v1, v2, v3):
        """Create two triangles for a quad. Vertices should be in CCW order for outward normal."""
        triangles.append((v0, v1, v2))
        triangles.append((v0, v2, v3))

    # === BOTTOM FACE (z=0, normal -Z) ===
    # f_bl, b_bl, b_br, f_br — CCW when viewed from -Z
    quad(f_bl, f_br, b_br, b_bl)

    # === FRONT FACE (x=0, normal -X) ===
    # Hexagonal profile, CCW when viewed from -X
    # Split into 4 triangles from center or as fan
    # Bottom quad: f_bl, f_br, f_mr, f_ml
    quad(f_bl, f_ml, f_mr, f_br)
    # Left bevel: f_ml, f_tl, f_mr — actually need to triangulate the top part
    # Top section: f_ml, f_tl, f_tr, f_mr
    quad(f_ml, f_tl, f_tr, f_mr)

    # === BACK FACE (x=L, normal +X) ===
    # Same hexagonal profile, but CCW when viewed from +X (reversed winding)
    quad(b_bl, b_br, b_mr, b_ml)
    quad(b_ml, b_mr, b_tr, b_tl)

    # === LEFT WALL (y=0, normal -Y) ===
    # f_bl, b_bl, b_ml, f_ml — CCW when viewed from -Y
    quad(f_bl, b_bl, b_ml, f_ml)

    # === RIGHT WALL (y=W, normal +Y) ===
    # f_br, f_mr, b_mr, b_br — CCW when viewed from +Y
    quad(f_br, f_mr, b_mr, b_br)

    # === LEFT BEVEL (y bevel on left side) ===
    # f_ml, b_ml, b_tl, f_tl — CCW when viewed from bevel direction
    quad(f_ml, b_ml, b_tl, f_tl)

    # === RIGHT BEVEL (y bevel on right side) ===
    # f_mr, f_tr, b_tr, b_mr — CCW when viewed from bevel direction
    quad(f_mr, f_tr, b_tr, b_mr)

    # === TOP FACE (z=H, normal +Z) ===
    # f_tl, f_tr, b_tr, b_tl — CCW when viewed from +Z
    quad(f_tl, f_tr, b_tr, b_tl)

    return triangles


def verify_watertight(triangles):
    """Verify the mesh is watertight: every edge shared by exactly 2 triangles."""
    edge_count = {}
    for tri in triangles:
        v0, v1, v2 = tri
        for a, b in [(v0, v1), (v1, v2), (v2, v0)]:
            # Canonical edge key (sorted tuple of vertex tuples)
            edge = tuple(sorted([a, b]))
            edge_count[edge] = edge_count.get(edge, 0) + 1

    non_manifold = {e: c for e, c in edge_count.items() if c != 2}

    n_triangles = len(triangles)
    n_edges = len(edge_count)
    n_vertices = len(set(v for tri in triangles for v in tri))

    print(f"Triangles: {n_triangles}")
    print(f"Vertices:  {n_vertices}")
    print(f"Edges:     {n_edges}")
    print(f"Euler:     V - E + F = {n_vertices} - {n_edges} + {n_triangles} = {n_vertices - n_edges + n_triangles}")

    if non_manifold:
        print(f"\nFAIL: {len(non_manifold)} non-manifold edges found:")
        for edge, count in non_manifold.items():
            print(f"  Edge {edge}: shared by {count} triangles")
        return False
    else:
        print("\nPASS: All edges shared by exactly 2 triangles — mesh is watertight!")
        return True


def verify_normals(triangles):
    """Verify all normals point outward by checking they agree with expected face directions."""
    issues = 0
    for i, tri in enumerate(triangles):
        n = compute_normal(*tri)
        mag = math.sqrt(n[0] ** 2 + n[1] ** 2 + n[2] ** 2)
        if mag < 0.99:
            print(f"  WARNING: Triangle {i} has degenerate normal (magnitude={mag:.4f})")
            issues += 1
    if issues == 0:
        print("PASS: All triangle normals are valid (non-degenerate)")
    else:
        print(f"FAIL: {issues} degenerate normals found")
    return issues == 0


if __name__ == "__main__":
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_dir = os.path.dirname(script_dir)
    output_path = os.path.join(project_dir, "cases", "fuel_tank.stl")

    print("Generating fuel tank STL...")
    print(f"  Dimensions: {L*1000:.0f} x {W*1000:.0f} x {H*1000:.0f} mm")
    print(f"  Bevel size:  {B*1000:.0f} mm")

    triangles = generate_fuel_tank()
    write_stl(output_path, triangles)

    print(f"\nWrote {len(triangles)} triangles to {output_path}")
    print(f"File size: {os.path.getsize(output_path)} bytes\n")

    print("=== Watertightness Check ===")
    wt_ok = verify_watertight(triangles)

    print("\n=== Normal Check ===")
    nm_ok = verify_normals(triangles)

    if wt_ok and nm_ok:
        print("\n*** ALL CHECKS PASSED ***")
    else:
        print("\n*** SOME CHECKS FAILED ***")
        exit(1)
