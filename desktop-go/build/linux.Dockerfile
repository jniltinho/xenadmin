# Builds the desktop app for Linux: Go plus the GTK 4 and WebKitGTK 6.0 headers Wails v3 links against (CGO), the
# AppImage tooling (pinned, as tools/linux/Dockerfile) and zip for the macOS app.
FROM ubuntu:24.04
ARG APPIMAGETOOL_VERSION=1.9.1
ARG APPIMAGE_RUNTIME_VERSION=20251108
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      ca-certificates curl zip file desktop-file-utils build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev libx11-dev \
    && rm -rf /var/lib/apt/lists/* \
    && curl -fsSL -o /usr/local/bin/appimagetool \
         https://github.com/AppImage/appimagetool/releases/download/${APPIMAGETOOL_VERSION}/appimagetool-x86_64.AppImage \
    && curl -fsSL -o /usr/local/share/appimage-runtime-x86_64 \
         https://github.com/AppImage/type2-runtime/releases/download/${APPIMAGE_RUNTIME_VERSION}/runtime-x86_64 \
    && chmod +x /usr/local/bin/appimagetool
ARG GO_VERSION=1.27.3
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz
# no FUSE inside containers
ENV PATH=/usr/local/go/bin:$PATH APPIMAGE_EXTRACT_AND_RUN=1
