#!/bin/bash
set -e

cd /opt/DualSPHysics/src/source

# Patch Makefile: fix CUDA toolkit path for container
sed -i 's|DIRTOOLKIT=/exports/opt/NVIDIA/cuda-12.8|DIRTOOLKIT=/usr/local/cuda|' Makefile

# Patch Makefile: add sm_89 (RTX 4090) to CUDA 12 gencode list
sed -i '/compute_86.*compute_86/a\  GENCODE:=$(GENCODE) -gencode=arch=compute_89,code=\\"sm_89,compute_89\\"' Makefile

echo "=== Patched Makefile for RTX 4090 (sm_89) ==="
echo "--- CUDA 12 gencode section ---"
sed -n '/ifeq ($(CUDAVER),12)/,/^endif/p' Makefile

# Build using the original Makefile
echo "=== Starting build with CUDA=12 ==="
make -j$(nproc) CUDA=12 2>&1

echo "=== Build complete ==="
# Check where the binary ended up
find /opt/DualSPHysics -name "DualSPHysics5.4_linux64" -type f 2>/dev/null
