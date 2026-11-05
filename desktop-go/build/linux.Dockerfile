# Builds the desktop app for Linux: Go plus the GTK 4 and WebKitGTK 6.0 headers Wails v3 links against (CGO).
FROM ubuntu:24.04
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      ca-certificates curl zip build-essential pkg-config libgtk-4-dev libwebkitgtk-6.0-dev \
    && rm -rf /var/lib/apt/lists/*
ARG GO_VERSION=1.27.3
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH=/usr/local/go/bin:$PATH
