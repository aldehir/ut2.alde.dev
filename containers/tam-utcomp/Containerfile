FROM ghcr.io/aldehir/ut2/base:3369.3-release-2

RUN dnf install -y bsdtar && \
    mkdir Community && \
    cd Community && \
    curl -sfL https://cdn.alde.dev/ut2k4/packs/ut2004-ufc-community-maps.zip | bsdtar -x --no-same-owner -vf -

ARG UT2U_VERSION=v0.1.4

COPY --chmod=755 install-packages.sh Community/
COPY packages.csv Community/
COPY System .System

RUN curl -sfL -o /usr/bin/ut2u https://github.com/aldehir/ut2u/releases/download/$UT2U_VERSION/ut2u-linux-amd64 && \
    chmod +x /usr/bin/ut2u && \
    rm -f StaticMeshes/DanielsMeshes.usx && \
    cp -v .System/*.ini System/ && \
    (cd Community && ./install-packages.sh) && \
    ut2u package check-deps System/UT2004.ini && \
    echo "Package dependency check passed!"

LABEL org.opencontainers.image.created="2023-09-01T12:00:00Z" \
      org.opencontainers.image.authors="Alde Rojas" \
      org.opencontainers.image.title="UT2004 Dedicated Server - TAM/UTComp Configuration" \
      org.opencontainers.image.description="UT2004 Dedicated Server - TAM/UTComp Configuration" \
      org.opencontainers.image.source="https://github.com/aldehir/ut2004-tam-utcomp-container"

CMD ["DM-Antalus?Game=xGame.xDeathMatch"]
