#!/usr/bin/env python3
"""Generate a triangulated surface VTK for the SPHERIC tank box geometry.

This creates a VTK PolyData file with inward-pointing normals (toward fluid)
for use with DualSPHysics mDBC <norgeometry> <geometryfile>.

Box: 900mm x 62mm x 508mm (5 faces, no top for lateral case)
     or 6 faces (with top for roof impact case)
"""

import sys

def generate_box_vtk(output_path, include_top=False,
                     lx=0.900, ly=0.062, lz=0.508):
    """Generate VTK PolyData of box surface mesh.

    Triangle winding ensures outward normals from solid (= inward toward fluid).
    """
    # 8 vertices of the box
    vertices = [
        (0,   0,   0),      # 0: origin
        (lx,  0,   0),      # 1: right-front-bottom
        (lx,  ly,  0),      # 2: right-back-bottom
        (0,   ly,  0),      # 3: left-back-bottom
        (0,   0,   lz),     # 4: left-front-top
        (lx,  0,   lz),     # 5: right-front-top
        (lx,  ly,  lz),     # 6: right-back-top
        (0,   ly,  lz),     # 7: left-back-top
    ]

    # Triangles with winding order for inward-pointing normals (toward fluid)
    # Each face = 2 triangles, normals point INTO the box interior (fluid domain)

    triangles = []

    # Bottom (z=0): normal +Z (into fluid above)
    # CCW when viewed from above: 0→1→2, 0→2→3
    triangles.append((0, 1, 2))
    triangles.append((0, 2, 3))

    # Left (x=0): normal +X (into fluid to the right)
    # CCW when viewed from +X: 0→3→7, 0→7→4
    triangles.append((0, 3, 7))
    triangles.append((0, 7, 4))

    # Right (x=lx): normal -X (into fluid to the left)
    # CCW when viewed from -X: 1→5→6, 1→6→2
    triangles.append((1, 5, 6))
    triangles.append((1, 6, 2))

    # Front (y=0): normal +Y (into fluid behind)
    # CCW when viewed from +Y: 0→4→5, 0→5→1
    triangles.append((0, 4, 5))
    triangles.append((0, 5, 1))

    # Back (y=ly): normal -Y (into fluid in front)
    # CCW when viewed from -Y: 3→2→6, 3→6→7
    triangles.append((3, 2, 6))
    triangles.append((3, 6, 7))

    if include_top:
        # Top (z=lz): normal -Z (into fluid below)
        # CCW when viewed from below: 4→7→6, 4→6→5
        triangles.append((4, 7, 6))
        triangles.append((4, 6, 5))

    n_verts = len(vertices)
    n_tris = len(triangles)

    with open(output_path, 'w') as f:
        f.write("# vtk DataFile Version 3.0\n")
        f.write(f"SPHERIC Tank Box Surface ({'6' if include_top else '5'} faces)\n")
        f.write("ASCII\n")
        f.write("DATASET POLYDATA\n")
        f.write(f"POINTS {n_verts} float\n")
        for v in vertices:
            f.write(f"{v[0]:.6f} {v[1]:.6f} {v[2]:.6f}\n")
        # Each triangle entry: "3 v0 v1 v2" = 4 ints
        f.write(f"POLYGONS {n_tris} {n_tris * 4}\n")
        for t in triangles:
            f.write(f"3 {t[0]} {t[1]} {t[2]}\n")

    print(f"Generated: {output_path}")
    print(f"  Vertices: {n_verts}, Triangles: {n_tris}")
    print(f"  Box: {lx}x{ly}x{lz}m, Top: {'yes' if include_top else 'no'}")


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("output", help="Output VTK file path")
    parser.add_argument("--top", action="store_true", help="Include top face")
    args = parser.parse_args()
    generate_box_vtk(args.output, include_top=args.top)
