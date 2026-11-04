# Linux desktop app (XenAdmin.Linux). Every target runs inside the build container;
# only docker and make are needed on the host. Outputs go to dist/.
#   make dist        self-contained single-file build + .tar.gz
#   make appimage    AppImage from the dist build
#   make test        unit tests (+ integration tests when .env defines XCP_*)
#   make smoke       run the AppImage under Xvfb in a clean Ubuntu 22.04 container

VERSION   ?= $(shell git describe --tags --match 'linux-v*' 2>/dev/null | sed 's/^linux-v//' | grep . || echo 0.1.0-dev)
RID       ?= linux-x64
IMAGE     ?= xcp-ng-center-build
NAME      := xcp-ng-center-$(VERSION)-$(RID)
DIST      := dist
UID_GID   := $(shell id -u):$(shell id -g)
RUN       := docker run --rm -e NUGET_PACKAGES=/src/_output/nuget -v "$(CURDIR)":/src -w /src $(IMAGE) sh -ec

.PHONY: all image build test test-win dist appimage smoke clean
all: appimage

image:
	docker build -q -t $(IMAGE) tools/linux

build: image
	$(RUN) 'dotnet build XenAdmin.Linux/XenAdmin.Linux.csproj; chown -R $(UID_GID) /src'

test:
	tools/linux/container.sh test

test-win:
	tools/linux/container.sh build-win

dist: image
	$(RUN) '\
	  rm -rf $(DIST)/$(NAME); \
	  dotnet publish XenAdmin.Linux/XenAdmin.Linux.csproj -c Release -r $(RID) --self-contained \
	    -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true \
	    -p:Version=$(VERSION) -p:DebugType=none -o $(DIST)/$(NAME); \
	  cp XenAdmin.Linux/packaging/xcp-ng-center.desktop XenAdmin.Linux/packaging/xcp-ng-center.png $(DIST)/$(NAME)/; \
	  tar -czf $(DIST)/$(NAME).tar.gz -C $(DIST) $(NAME); \
	  chown -R $(UID_GID) /src'
	@echo "-> $(DIST)/$(NAME)/xcp-ng-center"

appimage: dist
	$(RUN) '\
	  app=_output/AppDir; rm -rf $$app; mkdir -p $$app/usr/bin $$app/usr/share/icons/hicolor/48x48/apps; \
	  cp $(DIST)/$(NAME)/xcp-ng-center $$app/usr/bin/; \
	  cp XenAdmin.Linux/packaging/AppRun XenAdmin.Linux/packaging/xcp-ng-center.desktop XenAdmin.Linux/packaging/xcp-ng-center.png $$app/; \
	  cp XenAdmin.Linux/packaging/xcp-ng-center.png $$app/usr/share/icons/hicolor/48x48/apps/; \
	  desktop-file-validate $$app/xcp-ng-center.desktop; \
	  ARCH=x86_64 appimagetool --no-appstream --runtime-file /usr/local/share/appimage-runtime-x86_64 \
	    $$app $(DIST)/xcp-ng-center-$(VERSION)-x86_64.AppImage; \
	  chown -R $(UID_GID) /src'
	@echo "-> $(DIST)/xcp-ng-center-$(VERSION)-x86_64.AppImage"

smoke:
	docker run --rm -v "$(CURDIR)/$(DIST)":/dist:ro -v "$(CURDIR)/_output":/out ubuntu:22.04 sh -ec '\
	  apt-get update -qq >/dev/null; \
	  DEBIAN_FRONTEND=noninteractive apt-get install -y -qq xvfb x11-apps imagemagick libfontconfig1 libice6 libsm6 >/dev/null; \
	  export APPIMAGE_EXTRACT_AND_RUN=1 HOME=/tmp; \
	  /dist/xcp-ng-center-$(VERSION)-x86_64.AppImage --version; \
	  Xvfb :99 -screen 0 1280x800x24 & sleep 2; export DISPLAY=:99; \
	  /dist/xcp-ng-center-$(VERSION)-x86_64.AppImage & pid=$$!; sleep 15; \
	  kill -0 $$pid; xwd -root -silent | convert xwd:- /out/smoke.png; kill $$pid; \
	  chown $(UID_GID) /out/smoke.png; echo "smoke OK -> _output/smoke.png"'

clean:
	rm -rf $(DIST) _output/AppDir
