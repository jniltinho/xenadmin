# Build image for the desktop app: .NET SDK, AppImage tooling (pinned) and zip for the macOS app.
FROM mcr.microsoft.com/dotnet/sdk:10.0
ARG APPIMAGETOOL_VERSION=1.9.1
ARG APPIMAGE_RUNTIME_VERSION=20251108
RUN apt-get update \
 && apt-get install -y --no-install-recommends file desktop-file-utils make zip \
 && rm -rf /var/lib/apt/lists/* \
 && curl -fsSL -o /usr/local/bin/appimagetool \
      https://github.com/AppImage/appimagetool/releases/download/${APPIMAGETOOL_VERSION}/appimagetool-x86_64.AppImage \
 && curl -fsSL -o /usr/local/share/appimage-runtime-x86_64 \
      https://github.com/AppImage/type2-runtime/releases/download/${APPIMAGE_RUNTIME_VERSION}/runtime-x86_64 \
 && chmod +x /usr/local/bin/appimagetool
# repo is bind-mounted with the host owner
RUN git config --system --add safe.directory /src
# no FUSE inside containers
ENV APPIMAGE_EXTRACT_AND_RUN=1 DOTNET_CLI_TELEMETRY_OPTOUT=1 DOTNET_NOLOGO=1
