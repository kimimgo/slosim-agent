"""SPHERIC Test 10 visualization - Mesa-compatible pvpython script.
Individual file loading to avoid Mesa memory corruption with large file lists.
"""
import os
import sys
import glob

from paraview.simple import *

BASE_DIR = "/home/imgyu/workspace/02_active/slosim-agent/simulations"
RESOLUTION = [1280, 720]

def render_case(vtk_dir, case_label):
    vtk_pattern = os.path.join(vtk_dir, "PartFluid_*.vtk")
    vtk_files = sorted(glob.glob(vtk_pattern))
    if not vtk_files:
        print(f"ERROR: No VTK files at {vtk_pattern}")
        return False

    frames_dir = os.path.join(vtk_dir, "frames")
    os.makedirs(frames_dir, exist_ok=True)

    print(f"[{case_label}] Rendering {len(vtk_files)} frames...")

    paraview.simple._DisableFirstRenderCameraReset()

    view = CreateRenderView()
    view.ViewSize = RESOLUTION
    view.Background = [0.08, 0.08, 0.15]

    # Camera: front view (quasi-2D tank, 900x62x508mm)
    view.CameraPosition = [0.45, -1.2, 0.25]
    view.CameraFocalPoint = [0.45, 0.031, 0.25]
    view.CameraViewUp = [0, 0, 1]
    view.CameraParallelScale = 0.4

    for i, vtk_file in enumerate(vtk_files):
        reader = LegacyVTKReader(FileNames=[vtk_file])
        display = Show(reader, view)
        display.Representation = 'Points'
        display.PointSize = 3.0

        ColorBy(display, ('POINTS', 'Press'))
        pressLUT = GetColorTransferFunction('Press')
        pressLUT.ApplyPreset('Blue to Red Rainbow', True)
        pressLUT.RescaleTransferFunction(0, 4000)

        Render()
        frame_path = os.path.join(frames_dir, f"frame_{i:04d}.png")
        SaveScreenshot(frame_path, view, ImageResolution=RESOLUTION)

        Delete(display)
        Delete(reader)

        if (i + 1) % 25 == 0 or i == 0:
            print(f"  [{case_label}] Frame {i+1}/{len(vtk_files)}")

    Delete(view)
    print(f"[{case_label}] Done: {len(vtk_files)} frames in {frames_dir}")
    return True


def render_case_3d(vtk_dir, case_label, cam_pos, cam_focus, cam_scale=0.5):
    """Render 3D tank case with custom camera."""
    vtk_pattern = os.path.join(vtk_dir, "PartFluid_*.vtk")
    vtk_files = sorted(glob.glob(vtk_pattern))
    if not vtk_files:
        print(f"ERROR: No VTK files at {vtk_pattern}")
        return False

    frames_dir = os.path.join(vtk_dir, "frames")
    os.makedirs(frames_dir, exist_ok=True)

    print(f"[{case_label}] Rendering {len(vtk_files)} frames...")

    paraview.simple._DisableFirstRenderCameraReset()

    view = CreateRenderView()
    view.ViewSize = RESOLUTION
    view.Background = [0.08, 0.08, 0.15]
    view.CameraPosition = cam_pos
    view.CameraFocalPoint = cam_focus
    view.CameraViewUp = [0, 0, 1]
    view.CameraParallelScale = cam_scale

    for i, vtk_file in enumerate(vtk_files):
        reader = LegacyVTKReader(FileNames=[vtk_file])
        display = Show(reader, view)
        display.Representation = 'Points'
        display.PointSize = 3.0

        ColorBy(display, ('POINTS', 'Press'))
        pressLUT = GetColorTransferFunction('Press')
        pressLUT.ApplyPreset('Blue to Red Rainbow', True)
        pressLUT.RescaleTransferFunction(0, 4000)

        Render()
        frame_path = os.path.join(frames_dir, f"frame_{i:04d}.png")
        SaveScreenshot(frame_path, view, ImageResolution=RESOLUTION)

        Delete(display)
        Delete(reader)

        if (i + 1) % 25 == 0 or i == 0:
            print(f"  [{case_label}] Frame {i+1}/{len(vtk_files)}")

    Delete(view)
    print(f"[{case_label}] Done: {len(vtk_files)} frames in {frames_dir}")
    return True


if __name__ == "__main__":
    # Scenario 9: Zhao 2024 HorizCyl (D=1m, L=3m, ~90% fill, pitch 3deg)
    # Primitives approach: drawcylinder solid boundary + fluid overwrite
    render_case_3d(
        f"{BASE_DIR}/zhao_horizcyl", "Zhao_HorizCyl",
        cam_pos=[1.5, -3.0, 0.5],
        cam_focus=[1.5, 0.0, 0.5],
        cam_scale=1.0
    )
    print("Zhao HorizCyl rendering complete.")
