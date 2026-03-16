#!/usr/bin/env python3
"""SPHERIC T10 Oil ρ=920 α=1.0 — IsoSurface rendering + pressure plot synced video
Upper: 3D isosurface with sensor marker
Lower: pressure-time plot with moving cursor
Output: MP4 video"""
import pyvista as pv
import numpy as np
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
from matplotlib.gridspec import GridSpec
from PIL import Image
import os
import subprocess
import sys

# ── Paths ──
script_dir = os.path.dirname(os.path.abspath(__file__))
iso_dir = os.path.join(script_dir, '../../simulations/oil_rho920_a1_0/iso')
press_csv = os.path.join(script_dir, '../../simulations/oil_rho920_a1_0/out/measure_lateral_Press.csv')  # z=0.093 (free surface, matches experiment)
exp_file = os.path.join(script_dir, '../../paper-pof/data/spheric/lateral_oil_1x.txt')
frame_dir = os.path.join(script_dir, '../figures/render_a1_0_frames')
out_dir = os.path.join(script_dir, '../figures/organized/expc_spheric')
os.makedirs(frame_dir, exist_ok=True)
os.makedirs(out_dir, exist_ok=True)

# ── Load pressure data ──
def load_press(path):
    t, p = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i < 4: continue
            parts = line.strip().split(';')
            if len(parts) >= 3:
                t.append(float(parts[1]))
                p.append(float(parts[2]))
    return np.array(t), np.array(p)

def load_exp(path):
    t, p = [], []
    with open(path) as f:
        for i, line in enumerate(f):
            if i == 0: continue
            parts = line.strip().split('\t')
            if len(parts) >= 2:
                try:
                    t.append(float(parts[0]))
                    p.append(float(parts[1]) * 100.0)
                except: continue
    return np.array(t), np.array(p)

t_sph, p_sph = load_press(press_csv)
t_exp, p_exp = load_exp(exp_file)

# ── Sensor position (at rest, attached to tank wall) ──
SENSOR_X0, SENSOR_Y0, SENSOR_Z0 = 0.005, 0.031, 0.093  # lateral wall sensor

# ── Tank geometry ──
TANK_L, TANK_W, TANK_H = 0.900, 0.062, 0.508

# ── Pitch motion parameters (from XML mvrotsinu) ──
PITCH_CX, PITCH_CZ = 0.45, 0.0        # rotation center (bottom center)
PITCH_FREQ = 0.651                      # Hz
PITCH_AMPL = 4.0                        # degrees

def pitch_angle(t):
    """Current pitch angle in radians at time t.
    DualSPHysics axis p1→p2 = (0,-Y,0). Right-hand rule around -Y:
    positive rotation tilts left wall (x=0) upward in XZ plane.
    Negate so standard 2D rotation matrix matches solver convention."""
    return -np.radians(PITCH_AMPL) * np.sin(2 * np.pi * PITCH_FREQ * t)

def rotate_point_pitch(x, z, t):
    """Rotate a point (x,z) around pitch center by current angle"""
    theta = pitch_angle(t)
    dx, dz = x - PITCH_CX, z - PITCH_CZ
    rx = PITCH_CX + dx * np.cos(theta) - dz * np.sin(theta)
    rz = PITCH_CZ + dx * np.sin(theta) + dz * np.cos(theta)
    return rx, rz

def make_tank_wireframe_rotated(t):
    """Tank wireframe rotated by pitch angle at time t"""
    pts_rest = np.array([
        [0,0,0],[TANK_L,0,0],[TANK_L,TANK_W,0],[0,TANK_W,0],
        [0,0,TANK_H],[TANK_L,0,TANK_H],[TANK_L,TANK_W,TANK_H],[0,TANK_W,TANK_H]
    ])
    pts = pts_rest.copy()
    for i in range(len(pts)):
        rx, rz = rotate_point_pitch(pts[i, 0], pts[i, 2], t)
        pts[i, 0] = rx
        pts[i, 2] = rz
    lines = [
        [0,1],[1,2],[2,3],[3,0],  # bottom
        [4,5],[5,6],[6,7],[7,4],  # top
        [0,4],[1,5],[2,6],[3,7]   # verticals
    ]
    all_pts = []
    all_lines = []
    for l in lines:
        idx = len(all_pts)
        all_pts.extend([pts[l[0]], pts[l[1]]])
        all_lines.append([2, idx, idx+1])
    edges = pv.PolyData(np.array(all_pts))
    edges.lines = np.hstack(all_lines)
    return edges

# ── Frames: every 10th frame = ~70 frames for 7s at 0.01s interval ──
iso_files = sorted([f for f in os.listdir(iso_dir) if f.endswith('.vtk')])
total_frames = len(iso_files)
# Sample every 7 frames → ~100 frames → 10s video at 10fps
step = 7
frame_indices = list(range(0, total_frames, step))
dt_per_frame = 7.0 / total_frames  # time per iso frame

print(f"Total iso files: {total_frames}")
print(f"Rendering {len(frame_indices)} frames (step={step})")

# ── pyvista offscreen setup ──
pv.OFF_SCREEN = True

# Tank wireframe is now generated per-frame (rotated) — see make_tank_wireframe_rotated()

# ── Render one composite frame ──
def render_frame(frame_idx, out_path):
    iso_file = iso_files[frame_idx]
    time_s = frame_idx * dt_per_frame

    # 1) Render 3D isosurface with pyvista
    mesh = pv.read(os.path.join(iso_dir, iso_file))

    plotter = pv.Plotter(off_screen=True, window_size=[1600, 700])
    plotter.set_background('white')

    # Isosurface (fluid)
    plotter.add_mesh(mesh, color='#3498db', opacity=0.75, smooth_shading=True)

    # Tank wireframe — rotated by pitch angle
    tank_wire = make_tank_wireframe_rotated(time_s)
    plotter.add_mesh(tank_wire, color='black', line_width=2, style='wireframe')

    # Sensor marker — rotated with tank (attached to wall)
    sx, sz = rotate_point_pitch(SENSOR_X0, SENSOR_Z0, time_s)
    sensor = pv.Sphere(radius=0.008, center=(sx, SENSOR_Y0, sz))
    plotter.add_mesh(sensor, color='red')
    plotter.add_point_labels(
        np.array([[sx + 0.03, SENSOR_Y0, sz + 0.02]]),
        ['Sensor P'], font_size=14, text_color='red', bold=True,
        shape=None, render_points_as_spheres=False
    )

    # Pitch center marker (small green sphere)
    plotter.add_mesh(pv.Sphere(radius=0.005, center=(PITCH_CX, SENSOR_Y0, PITCH_CZ)),
                     color='green', opacity=0.7)
    ang_deg = np.degrees(pitch_angle(time_s))
    plotter.add_text(f'pitch = {ang_deg:+.1f}°', position='upper_edge',
                     font_size=11, color='darkgreen')

    # Camera: zoomed out side view — tank is 0.9m×0.062m×0.508m
    plotter.camera_position = [
        (0.45, -1.4, 0.5),   # camera far back (Y=-1.4) and elevated
        (0.45, 0.031, 0.20),  # focal point (tank center, slightly above mid)
        (0, 0, 1)             # up vector
    ]

    # Time annotation
    plotter.add_text(f't = {time_s:.2f} s', position='upper_left',
                     font_size=14, color='black')
    plotter.add_text('SPHERIC T10 Oil ρ=920 α=1.0', position='upper_right',
                     font_size=11, color='gray')

    img_3d = plotter.screenshot(return_img=True)
    plotter.close()

    # 2) Create pressure plot with matplotlib
    fig_plot, ax_plot = plt.subplots(figsize=(16, 3.5), dpi=100)
    ax_plot.plot(t_exp, p_exp/1000, 'k-', lw=1.5, alpha=0.7, label='Experiment')
    ax_plot.plot(t_sph, p_sph/1000, '-', color='#3498db', lw=1.5, alpha=0.8, label='SPH α=1.0')
    # Current time cursor
    ax_plot.axvline(x=time_s, color='red', lw=2, alpha=0.9, ls='--')
    # Current pressure value
    idx_t = np.argmin(np.abs(t_sph - time_s))
    if idx_t < len(p_sph):
        ax_plot.plot(t_sph[idx_t], p_sph[idx_t]/1000, 'ro', markersize=10, zorder=20)
        ax_plot.annotate(f'{p_sph[idx_t]/1000:.2f} kPa', xy=(t_sph[idx_t], p_sph[idx_t]/1000),
                        xytext=(15, 15), textcoords='offset points', fontsize=10,
                        color='red', fontweight='bold',
                        arrowprops=dict(arrowstyle='->', color='red', lw=1.5))
    ax_plot.set_xlabel('Time (s)', fontsize=12)
    ax_plot.set_ylabel('Pressure (kPa)', fontsize=12)
    ax_plot.set_xlim(0, 7)
    ax_plot.set_ylim(-0.5, max(max(p_exp)/1000, max(p_sph)/1000) * 1.2)
    ax_plot.legend(fontsize=10, loc='upper right')
    ax_plot.grid(alpha=0.2)
    # Sensor info
    ax_plot.text(0.01, 0.95, f'Sensor: ({SENSOR_X0}, {SENSOR_Y0}, {SENSOR_Z0}) m — rotates with tank',
                transform=ax_plot.transAxes, fontsize=9, va='top', color='red')
    plt.tight_layout()

    import tempfile
    with tempfile.NamedTemporaryFile(suffix='.png', delete=False) as tf:
        tmp_path = tf.name
    fig_plot.savefig(tmp_path, dpi=100, bbox_inches='tight')
    plt.close(fig_plot)

    # 3) Composite: stack vertically
    img_plot = np.array(Image.open(tmp_path))
    os.unlink(tmp_path)

    # Resize 3D to match plot width
    h3d, w3d = img_3d.shape[:2]
    hp, wp = img_plot.shape[:2]
    target_w = max(w3d, wp)

    img_3d_pil = Image.fromarray(img_3d).resize((target_w, int(h3d * target_w / w3d)), Image.LANCZOS)
    img_plot_pil = Image.fromarray(img_plot).resize((target_w, int(hp * target_w / wp)), Image.LANCZOS)

    composite = Image.new('RGB', (target_w, img_3d_pil.height + img_plot_pil.height), (255, 255, 255))
    composite.paste(img_3d_pil, (0, 0))
    composite.paste(img_plot_pil, (0, img_3d_pil.height))
    composite.save(out_path)

# ── Render all frames ──
print("Rendering frames...")
for i, fi in enumerate(frame_indices):
    out_path = os.path.join(frame_dir, f'frame_{i:04d}.png')
    if os.path.exists(out_path):
        continue
    render_frame(fi, out_path)
    if (i+1) % 10 == 0 or i == 0:
        time_s = fi * dt_per_frame
        print(f"  [{i+1}/{len(frame_indices)}] t={time_s:.2f}s")

print("Encoding video...")
# ── ffmpeg encode ──
video_path = os.path.join(out_dir, 'oil_rho920_a1_0_isosurface.mp4')
subprocess.run([
    'ffmpeg', '-y', '-framerate', '10',
    '-i', os.path.join(frame_dir, 'frame_%04d.png'),
    '-c:v', 'libx264', '-pix_fmt', 'yuv420p',
    '-crf', '18', '-preset', 'medium',
    video_path
], check=True, capture_output=True)

print(f"Done: {video_path}")

# ── Also save one representative frame (peak moment) ──
# Find peak time in SPH data
peak_idx = np.argmax(p_sph)
peak_time = t_sph[peak_idx]
peak_frame = int(peak_time / dt_per_frame)
# Find closest rendered frame
closest = min(frame_indices, key=lambda x: abs(x - peak_frame))
closest_i = frame_indices.index(closest)
snapshot_src = os.path.join(frame_dir, f'frame_{closest_i:04d}.png')
snapshot_dst = os.path.join(out_dir, 'oil_rho920_a1_0_peak_snapshot.png')
if os.path.exists(snapshot_src):
    Image.open(snapshot_src).save(snapshot_dst)
    print(f"Snapshot at peak (t={closest * dt_per_frame:.2f}s): {snapshot_dst}")
