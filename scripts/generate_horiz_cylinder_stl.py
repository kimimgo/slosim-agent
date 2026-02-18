"""Generate a watertight horizontal cylinder STL (ASCII format).
D=1.0m, L=3.0m, axis along X, center at y=0, z=0.5.
36 segments for circular cross-section.
"""
import math
import sys

def vec_sub(a, b):
    return (a[0]-b[0], a[1]-b[1], a[2]-b[2])

def cross(a, b):
    return (a[1]*b[2]-a[2]*b[1], a[2]*b[0]-a[0]*b[2], a[0]*b[1]-a[1]*b[0])

def normalize(v):
    length = math.sqrt(v[0]**2 + v[1]**2 + v[2]**2)
    if length < 1e-12:
        return (0, 0, 1)
    return (v[0]/length, v[1]/length, v[2]/length)

def write_triangle(f, v1, v2, v3, normal=None):
    if normal is None:
        e1 = vec_sub(v2, v1)
        e2 = vec_sub(v3, v1)
        normal = normalize(cross(e1, e2))
    f.write(f"  facet normal {normal[0]:.6e} {normal[1]:.6e} {normal[2]:.6e}\n")
    f.write("    outer loop\n")
    f.write(f"      vertex {v1[0]:.6e} {v1[1]:.6e} {v1[2]:.6e}\n")
    f.write(f"      vertex {v2[0]:.6e} {v2[1]:.6e} {v2[2]:.6e}\n")
    f.write(f"      vertex {v3[0]:.6e} {v3[1]:.6e} {v3[2]:.6e}\n")
    f.write("    endloop\n")
    f.write("  endfacet\n")

def generate_cylinder(filename, radius=0.5, length=3.0, cy=0.0, cz=0.5, segments=36):
    """Generate horizontal cylinder along X axis."""
    angles = [2 * math.pi * i / segments for i in range(segments)]

    # Circle vertices at x=0 (left cap) and x=length (right cap)
    left_circle = [(0.0, cy + radius * math.cos(a), cz + radius * math.sin(a)) for a in angles]
    right_circle = [(length, cy + radius * math.cos(a), cz + radius * math.sin(a)) for a in angles]
    left_center = (0.0, cy, cz)
    right_center = (length, cy, cz)

    with open(filename, 'w') as f:
        f.write("solid horizontal_cylinder\n")

        for i in range(segments):
            j = (i + 1) % segments

            # Side wall: 2 triangles per segment (outward normals)
            # Quad: left_circle[i], left_circle[j], right_circle[j], right_circle[i]
            write_triangle(f, left_circle[i], right_circle[i], right_circle[j])
            write_triangle(f, left_circle[i], right_circle[j], left_circle[j])

            # Left cap: triangle fan pointing -X (inward = left)
            write_triangle(f, left_center, left_circle[j], left_circle[i],
                          normal=(-1.0, 0.0, 0.0))

            # Right cap: triangle fan pointing +X (outward = right)
            write_triangle(f, right_center, right_circle[i], right_circle[j],
                          normal=(1.0, 0.0, 0.0))

        f.write("endsolid horizontal_cylinder\n")

    # Verify watertight
    total_tris = segments * 4  # 2 side + 1 left cap + 1 right cap per segment
    total_edges = segments * 6  # Each triangle has 3 edges, each edge shared by 2 triangles
    print(f"Generated: {filename}")
    print(f"  Triangles: {total_tris}")
    print(f"  Segments: {segments}")
    print(f"  Radius: {radius}m, Length: {length}m")
    print(f"  Center: y={cy}, z={cz}")

    # Edge-sharing verification
    edges = {}
    tris = []
    # Re-read for verification
    with open(filename, 'r') as f:
        verts = []
        for line in f:
            line = line.strip()
            if line.startswith("vertex"):
                parts = line.split()
                verts.append((float(parts[1]), float(parts[2]), float(parts[3])))
                if len(verts) == 3:
                    tris.append(tuple(verts))
                    for k in range(3):
                        v1 = verts[k]
                        v2 = verts[(k+1)%3]
                        # Normalize edge direction
                        edge = tuple(sorted([
                            (round(v1[0],8), round(v1[1],8), round(v1[2],8)),
                            (round(v2[0],8), round(v2[1],8), round(v2[2],8))
                        ]))
                        edges[edge] = edges.get(edge, 0) + 1
                    verts = []

    non_manifold = {e: c for e, c in edges.items() if c != 2}
    if non_manifold:
        print(f"  WARNING: {len(non_manifold)} non-manifold edges!")
        for e, c in list(non_manifold.items())[:5]:
            print(f"    Edge {e}: {c} faces")
    else:
        print(f"  Watertight: YES ({len(edges)} edges, all shared by exactly 2 faces)")

if __name__ == "__main__":
    outfile = sys.argv[1] if len(sys.argv) > 1 else "cases/horiz_cylinder.stl"
    generate_cylinder(outfile)
