#!/bin/sh
# Build/test the net10.0 (Linux) targets inside the .NET SDK container.
# usage: tools/linux/container.sh build|test|build-win [extra dotnet args]
# build-win compiles the core for net481 (Windows app) to catch regressions.
# Integration tests read XCP_HOST/XCP_USER/XCP_PASS from ./.env (never committed).
set -e
cd "$(dirname "$0")/../.."
cmd=${1:-build}; [ $# -gt 0 ] && shift
envfile=""; [ -f .env ] && envfile="--env-file .env"
# forward shell overrides only when set: a bare "-e VAR" would unset the .env value
for v in XCP_HOST XCP_USER XCP_PASS XCP_RO_HOST XCP_ALLOW_WRITE XCP_LAB_HOSTS; do eval "[ -n \"\${$v+x}\" ]" && envfile="$envfile -e $v"; done
case $cmd in
  build) run="dotnet build tests/XenAdmin.Desktop.Tests -f net10.0 $*" ;;
  test)  run="dotnet test tests/XenAdmin.Desktop.Tests --logger 'console;verbosity=normal' $*" ;;
  build-win)
    refs=/src/_output/nuget/microsoft.netframework.referenceassemblies.net481/1.0.3/build/.NETFramework/v4.8.1
    run="[ -d $refs ] || (mkdir -p /tmp/r && cd /tmp/r && dotnet new classlib -o . >/dev/null && dotnet add package Microsoft.NETFramework.ReferenceAssemblies.net481 --version 1.0.3 >/dev/null); dotnet build XenModel/XenModel.csproj -f net481 -p:FrameworkPathOverride=$refs $*" ;;
  *) echo "usage: $0 build|test|build-win" >&2; exit 2 ;;
esac
docker run --rm $envfile \
  -e NUGET_PACKAGES=/src/_output/nuget -v "$PWD":/src -w /src \
  mcr.microsoft.com/dotnet/sdk:10.0 sh -c "$run; rc=\$?; chown -R $(id -u):$(id -g) /src; exit \$rc"
