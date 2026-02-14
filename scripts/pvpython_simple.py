"""
Minimal ParaView offscreen test: render single VTK frame to PNG.
Usage: /opt/paraview/bin/pvpython --force-offscreen-rendering --mesa scripts/pvpython_simple.py
"""
import os
import glob
from paraview.simple import *

DATA_DIR = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                        "simulations", "e2e_test")
vtk_files = sorted(glob.glob(os.path.join(DATA_DIR, "PartFluid_*.vtk")))
print(f"VTK files: {len(vtk_files)}")

# Load single frame
reader = LegacyVTKReader(FileNames=[vtk_files[50]])  # mid-simulation

view = CreateRenderView()
view.ViewSize = [1280, 720]
view.Background = [0.15, 0.15, 0.2]

display = Show(reader, view)
display.Representation = 'Points'
display.PointSize = 3.0

ColorBy(display, ('POINTS', 'Press'))
lut = GetColorTransferFunction('Press')
lut.ApplyPreset('Blue to Red Rainbow', True)
lut.RescaleTransferFunction(0, 5000)

view.CameraPosition = [0.5, -1.0, 0.5]
view.CameraFocalPoint = [0.5, 0.25, 0.3]
view.CameraViewUp = [0, 0, 1]

Render()

out = os.path.join(DATA_DIR, "test_snapshot.png")
SaveScreenshot(out, view, ImageResolution=[1280, 720])
print(f"Saved: {out}")
print(f"Size: {os.path.getsize(out)} bytes")
