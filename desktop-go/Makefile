# XCP-ng Center (Go). Everything runs in a container (only docker and make needed), as the C# app's Makefile.
# The desktop window is Wails v3 (GTK 4 + WebKitGTK 6.0 on Linux); `cli` builds the same binary without it.
GO_IMAGE ?= golang:1.27.3
VERSION  ?= $(or $(shell git describe --tags --match 'go-v*' 2>/dev/null | sed 's/^go-v//'),dev)
LDFLAGS  := -s -w -X xcpng-center/internal/buildinfo.Version=$(VERSION) -X xcpng-center/internal/buildinfo.Commit=$(shell git rev-parse --short HEAD 2>/dev/null)
NODE_IMAGE ?= node:24
DOCKER   := docker run --rm --user $(shell id -u):$(shell id -g) -e HOME=/tmp
GO_ENV   := -e GOFLAGS=-buildvcs=false -e GOCACHE=/src/.cache/go-build -e GOMODCACHE=/src/.cache/mod
RUN      := $(DOCKER) -v $(CURDIR):/src -w /src $(GO_ENV) $(GO_IMAGE)
# the tree icons of the C# app (XenAdmin/Images, branding-xcp-ng), copied into the frontend at build time
ICONS    := ../branding-xcp-ng/Images/Logo.png $(addprefix ../XenAdmin/Images/,000_ServerDisconnected_h32bit_16.png \
            000_PoolConnected_h32bit_16.png 000_TreeConnected_h32bit_16.png 000_StartVM_h32bit_16.png \
            000_StoppedVM_h32bit_16.png 000_SuspendVM_h32bit_16.png 000_Storage_h32bit_16.png 000_StorageBroken_h32bit_16.png 000_StorageDefault_h32bit_16.png 000_StorageDisabled_h32bit_16.png \
            000_AddApplicationServer_h32bit_16.png \
            000_ServerMaintenance_h32bit_16.png 000_CPU_h32bit_16.png 000_ServerHome_h32bit_16.png 000_RebootVM_h32bit_16.png 000_ExportVirtualAppliance_h32bit_16.png 000_ImportVirtualAppliance_h32bit_16.png padlock.png 075_TickRound_h32bit_16.png 000_Search_h32bit_16.png infra_view_16_textured.png objects_16.png \
            $(foreach i,001_Back 001_Forward 000_AddApplicationServer 000_PoolNew 000_NewStorage 000_CreateVM,$(i)_h32bit_24.png) 001_CreateVM_h32bit_16.png 000_Paused_h32bit_24_green.png 000_Resumed_h32bit_24_green.png 000_paused_h32bit_16.png objects_24.png org_view_24.png saved_searches_24.png 000_Info3_h32bit_16.png 000_WarningAlert_h32bit_32.png 000_PoolNew_h32bit_16.png 002_Configure_h32bit_16.png 000_Tag_h32bit_16.png 000_Folder_open_h32bit_16.png 000_Fields_h32bit_16.png 000_VirtualAppliance_h32bit_16.png 000_defaultSpyglass_h32bit_16.png 000_VMTemplate_h32bit_16.png console_16.png 000_VMSnapShotDiskOnly_h32bit_16.png 000_GetMemoryInfo_h32bit_16.png 000_CPU_h32bit_16.png usb_16.png 000_Tick_h32bit_16.png 000_Abort_h32bit_16.png 000_VirtualStorage_h32bit_16.png 000_Network_h32bit_16.png 000_VMSnapshotDiskMemory_h32bit_16.png tree_plus.png tree_minus.png 001_PerformanceGraph_h32bit_16.png $(foreach i,alert1 alert2 alert3 alert4 alert5 alert6 infra_view notif_none notif_alerts,$(i)_16.png) infra_view_24.png notif_none_24.png \
            $(foreach i,001_PowerOn 001_ShutDown 001_Reboot 000_Paused 000_Resumed 001_ForceShutDown 001_ForceReboot,$(i)_h32bit_24.png) \
            $(foreach i,001_PowerOn 001_ShutDown 001_Reboot 000_paused 000_Resumed 001_ForceShutDown 001_ForceReboot,$(i)_h32bit_16.png) PDChevronUp.png PDChevronDown.png)

.PHONY: test check cli icons frontend linux appimage windows macos dist clean

test:            ## unit tests of the app and the vendored SDK
	$(RUN) sh -c 'go test -tags nogui ./... && cd sdk/xenapi && go test ./...'

check:           ## gofmt, go vet (both modules) and the tests
	$(RUN) sh -c 'test -z "$$(gofmt -l cmd internal main.go)" || { gofmt -l cmd internal main.go; exit 1; }; go vet -tags nogui ./... && (cd sdk/xenapi && go vet ./...) && go test -tags nogui ./...'

cli:             ## static CLI/MCP binaries (no desktop window) for Linux, macOS and Windows, amd64, in dist/
	$(RUN) sh -c 'for t in linux/amd64 darwin/amd64 windows/amd64; do \
	  os=$${t%/*}; arch=$${t#*/}; ext=; [ $$os = windows ] && ext=.exe; \
	  CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -tags nogui -trimpath -ldflags "$(LDFLAGS)" -o dist/xcp-ng-center-go-$(VERSION)-$$os-$$arch$$ext . || exit 1; done'

icons:           ## the C# app's icons into frontend/public/icons (no container: also used on a Mac)
	mkdir -p frontend/public/icons && cp $(ICONS) frontend/public/icons/ && cp ../branding-xcp-ng/Images/xcp-ng-logo1glossy-1.png frontend/public/icons/home-logo.png && cp ../XenAdmin.Desktop/packaging/xcp-ng-center.png frontend/public/icons/app.png

frontend: icons  ## the window's Vue/Tailwind assets into internal/gui/dist (embedded by the GUI build)
	$(DOCKER) -v $(CURDIR):/src -w /src/frontend $(NODE_IMAGE) sh -c 'npm ci --no-audit --no-fund && npm run build'

linux: frontend  ## the desktop app for Linux (amd64) in dist/
	docker build -q -t xcpng-center-go-linux -f build/linux.Dockerfile build
	$(DOCKER) -v $(CURDIR):/src -w /src $(GO_ENV) xcpng-center-go-linux \
	  go build -tags production -trimpath -ldflags "$(LDFLAGS)" -o dist/xcp-ng-center-go-$(VERSION)-linux-amd64-desktop .

# The AppImage holds the binary, its .desktop and icon; GTK 4 and WebKitGTK 6.0 come from the system (Ubuntu 24.04+,
# Fedora 39+, Debian 13+), as for any WebKitGTK app.
# shortcut: the GTK/WebKit libraries are not bundled; bundle them (linuxdeploy's GTK plugin) for older distributions.
appimage: linux ## the desktop app for Linux as an x86_64 AppImage in dist/
	$(DOCKER) -v $(CURDIR):/src -v $(CURDIR)/../XenAdmin.Desktop/packaging:/pkg:ro -w /src xcpng-center-go-linux sh -c '\
	  app=.cache/AppDir; rm -rf $$app; mkdir -p $$app/usr/bin $$app/usr/share/icons/hicolor/256x256/apps; \
	  cp dist/xcp-ng-center-go-$(VERSION)-linux-amd64-desktop $$app/usr/bin/xcp-ng-center-go; \
	  cp build/linux/AppRun build/linux/xcp-ng-center-go.desktop $$app/; \
	  cp /pkg/xcp-ng-center.png $$app/xcp-ng-center-go.png; cp /pkg/xcp-ng-center.png $$app/usr/share/icons/hicolor/256x256/apps/xcp-ng-center-go.png; \
	  desktop-file-validate $$app/xcp-ng-center-go.desktop && \
	  ARCH=x86_64 appimagetool --no-appstream --runtime-file /usr/local/share/appimage-runtime-x86_64 \
	    $$app dist/xcp-ng-center-go-$(VERSION)-x86_64.AppImage'

windows: frontend ## the desktop app for Windows (amd64; WebView2, no CGO) in dist/
	$(RUN) sh -c 'CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags production -trimpath -ldflags "$(LDFLAGS) -H windowsgui" \
	  -o dist/xcp-ng-center-go-$(VERSION)-windows-amd64-desktop.exe .'

# "XCP-ng Center (Go).app" for Intel Macs (amd64; Apple Silicon runs it with Rosetta 2).
# Every desktop build has -tags production: no developer tools, and never a FRONTEND_DEVSERVER_URL in place of the
# embedded window. Cross-built from Linux with Wails' CGO-free backend (-tags purego, experimental in Wails 3.1).
# shortcut: unverified on a Mac and unsigned (`xattr -dr com.apple.quarantine` after download); on a Mac, build
# natively with CGO (`CGO_ENABLED=1 go build`) and sign/notarize before a release that claims macOS support.
macos: frontend ## the desktop app for macOS (amd64) as a zip in dist/
	docker build -q -t xcpng-center-go-linux -f build/linux.Dockerfile build
	$(DOCKER) -v $(CURDIR):/src -v $(CURDIR)/../XenAdmin.Desktop/packaging/macos:/macos:ro -w /src $(GO_ENV) xcpng-center-go-linux sh -c '\
	  out=dist/macos; app="$$out/XCP-ng Center (Go).app"; rm -rf $$out; mkdir -p "$$app/Contents/Resources"; \
	  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -tags production,purego -trimpath \
	    -ldflags "$(LDFLAGS)" -o "$$app/Contents/MacOS/xcp-ng-center-go" . || exit 1; \
	  cp /macos/xcp-ng-center.icns "$$app/Contents/Resources/"; \
	  v=$(VERSION); sed -e "s/@VERSION@/$$v/" -e "s/@SHORT_VERSION@/$${v%%-*}/" build/darwin/Info.plist > "$$app/Contents/Info.plist"; \
	  rm -f dist/xcp-ng-center-go-$(VERSION)_macos_amd64.zip; \
	  (cd $$out && zip -qry ../xcp-ng-center-go-$(VERSION)_macos_amd64.zip "XCP-ng Center (Go).app")'

dist: check cli appimage windows macos

clean:
	rm -rf dist .cache
