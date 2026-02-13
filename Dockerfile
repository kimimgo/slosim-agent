# =============================================================================
# DualSPHysics v5.4 GPU Docker Image
# Optimized for NVIDIA RTX 4090 (sm_89, Ada Lovelace) + CUDA 12.6
# =============================================================================

# ---------------------------------------------------------------------------
# Stage 1: Build DualSPHysics GPU binary from source
# ---------------------------------------------------------------------------
FROM nvidia/cuda:12.6.3-devel-ubuntu22.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    g++ \
    make \
    && rm -rf /var/lib/apt/lists/*

# Copy entire DualSPHysics tree to preserve relative paths
COPY DualSPHysics/ /opt/DualSPHysics/

# Build MoorDynPlus library
WORKDIR /opt/DualSPHysics/src/source/Source_DSphMoorDynPlus
RUN make -f Makefile_MoorDynPlus -j$(nproc) && \
    echo "=== MoorDynPlus built ==="

# Build DualSPHysics GPU binary
COPY build-dsph.sh /tmp/build-dsph.sh
RUN chmod +x /tmp/build-dsph.sh && /tmp/build-dsph.sh

# ---------------------------------------------------------------------------
# Stage 2: Lean runtime image
# ---------------------------------------------------------------------------
FROM nvidia/cuda:12.6.3-runtime-ubuntu22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgomp1 \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /opt/dsph/bin /opt/dsph/examples /data

# Copy GPU binary built from source (output goes to bin/linux/ per Makefile)
COPY --from=builder /opt/DualSPHysics/bin/linux/DualSPHysics5.4_linux64 /opt/dsph/bin/

# Copy prebuilt tools
COPY --from=builder /opt/DualSPHysics/bin/linux/GenCase_linux64 /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/PartVTK_linux64 /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/PartVTKOut_linux64 /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/MeasureTool_linux64 /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/IsoSurface_linux64 /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/DsphConfig.xml /opt/dsph/bin/

# Copy shared libraries
COPY --from=builder /opt/DualSPHysics/bin/linux/libChronoEngine.so /opt/dsph/bin/
COPY --from=builder /opt/DualSPHysics/bin/linux/libdsphchrono.so /opt/dsph/bin/

# Copy examples
COPY --from=builder /opt/DualSPHysics/examples/ /opt/dsph/examples/

# Copy XML format documentation (for advanced skill reference)
COPY --from=builder /opt/DualSPHysics/doc/xml_format/ /opt/dsph/doc/xml_format/

# Environment
ENV DSPH_HOME=/opt/dsph
ENV PATH="${DSPH_HOME}/bin:${PATH}"
ENV LD_LIBRARY_PATH="${DSPH_HOME}/bin:/usr/local/cuda/lib64:${LD_LIBRARY_PATH}"

# Symlinks for convenience
RUN cd /opt/dsph/bin && \
    ln -sf GenCase_linux64 gencase && \
    ln -sf DualSPHysics5.4_linux64 dualsphysics && \
    ln -sf PartVTK_linux64 partvtk && \
    ln -sf PartVTKOut_linux64 partvtkout && \
    ln -sf MeasureTool_linux64 measuretool && \
    ln -sf IsoSurface_linux64 isosurface && \
    chmod +x /opt/dsph/bin/*

WORKDIR /data

CMD ["/bin/bash"]
