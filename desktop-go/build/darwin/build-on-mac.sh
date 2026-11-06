#!/bin/sh
# Native macOS build of the desktop app, on a Mac (macOS 14+, Go 1.27, Node 24, the Xcode command line tools):
# the frontend, the app with Wails' native (CGO) backend for this Mac's architecture or ARCH=amd64|arm64, the
# "XCP-ng Center (Go).app" bundle and its zip in dist/. With SIGN_IDENTITY ("Developer ID Application: ...") the
# app is signed with the hardened runtime; with NOTARY_PROFILE too (xcrun notarytool store-credentials) the zip is
# notarized and the ticket stapled. `make macos` on Linux makes the same bundle, cross-built and unsigned (amd64).
set -eu
cd "$(dirname "$0")/../.."
VERSION=${VERSION:-$(git describe --tags --match 'go-v*' 2>/dev/null | sed 's/^go-v//' || true)}
VERSION=${VERSION:-dev}
case "${ARCH:-$(uname -m)}" in arm64 | aarch64) ARCH=arm64 ;; x86_64 | amd64) ARCH=amd64 ;; *) echo "unknown ARCH" >&2; exit 1 ;; esac

make icons
(cd frontend && npm ci --no-audit --no-fund && npm run build)

out=dist/macos-$ARCH
app="$out/XCP-ng Center (Go).app"
rm -rf "$out" && mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"
# the minimum macOS for clang's objects and for Go's external link (else ld links for 10.13 and warns on each object)
min=-mmacosx-version-min=14.0
CGO_ENABLED=1 GOOS=darwin GOARCH=$ARCH MACOSX_DEPLOYMENT_TARGET=14.0 CGO_CFLAGS="${CGO_CFLAGS:-} $min" \
  CGO_LDFLAGS="${CGO_LDFLAGS:-} $min" go build -tags production -trimpath \
  -ldflags "-s -w -X xcpng-center/internal/buildinfo.Version=$VERSION -X xcpng-center/internal/buildinfo.Commit=$(git rev-parse --short HEAD)" \
  -o "$app/Contents/MacOS/xcp-ng-center-go" .
cp ../XenAdmin.Desktop/packaging/macos/xcp-ng-center.icns "$app/Contents/Resources/"
sed -e "s/@VERSION@/$VERSION/" -e "s/@SHORT_VERSION@/${VERSION%%-*}/" build/darwin/Info.plist > "$app/Contents/Info.plist"

if [ -n "${SIGN_IDENTITY:-}" ]; then
  codesign --force --deep --options runtime --timestamp --sign "$SIGN_IDENTITY" "$app"
  codesign --verify --strict --verbose=2 "$app"
fi
zip="dist/xcp-ng-center-go-${VERSION}_macos_$ARCH.zip"
rm -f "$zip" && ditto -c -k --keepParent "$app" "$zip"
if [ -n "${SIGN_IDENTITY:-}" ] && [ -n "${NOTARY_PROFILE:-}" ]; then
  xcrun notarytool submit "$zip" --keychain-profile "$NOTARY_PROFILE" --wait
  xcrun stapler staple "$app" && rm -f "$zip" && ditto -c -k --keepParent "$app" "$zip"
fi
echo "$zip"
