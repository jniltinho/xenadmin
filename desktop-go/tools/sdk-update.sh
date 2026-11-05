#!/bin/sh
# Replace sdk/xenapi with a new XenServer Go SDK and re-apply sdk/patches/xcp-ng.patch (see sdk/xenapi/PATCHES.md).
# Nothing in sdk/xenapi changes unless the patch applies cleanly to the new SDK.
# usage: tools/sdk-update.sh XenServer-SDK-<version>.zip
set -eu
zip=${1:?usage: $0 XenServer-SDK-<version>.zip}
cd "$(dirname "$0")/.."
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
unzip -q "$zip" -d "$work/zip"
src=$(find "$work/zip" -type d -path '*XenServerGo/src' | head -1)
[ -n "$src" ] || { echo "no XenServerGo/src in $zip" >&2; exit 1; }
mkdir "$work/stage"
cp "$src"/*.go "$work/stage/"
printf 'module go/xenapi\n\ngo 1.22.2\n' > "$work/stage/go.mod"
if ! patch -p1 -d "$work/stage" --dry-run < sdk/patches/xcp-ng.patch >/dev/null; then
  echo "the XCP-NG patches do not apply to this SDK: port them by hand (PATCHES.md), then regenerate sdk/patches/xcp-ng.patch" >&2
  exit 1
fi
patch -p1 -d "$work/stage" < sdk/patches/xcp-ng.patch >/dev/null
rm -f sdk/xenapi/*.go
cp "$work/stage"/*.go "$work/stage/go.mod" sdk/xenapi/
echo "SDK replaced and patched; update the version in PATCHES.md and run make check"
