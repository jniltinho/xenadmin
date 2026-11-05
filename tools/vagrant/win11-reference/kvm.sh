#!/bin/sh
# Windows 11 with the official XCP-ng Center (UI reference) on KVM via dockur/windows, on the same Docker network as
# the KVM lab (tools/vagrant/xcp-ng-lab, XCP_LAB_BACKEND=kvm), so both run together.
#   ./kvm.sh up          # first run downloads and installs Windows 11 (~30 min); web console http://127.0.0.1:8006
#   ./kvm.sh shot NAME   # screenshot to _output/reference/NAME.png
#   ./kvm.sh ps 'cmd'    # PowerShell in the logged-on desktop session (via ssh + scheduled task)
#   ./kvm.sh key esc ret # keys through the QEMU monitor
#   ./kvm.sh down
set -e
cd "$(dirname "$0")"
ROOT=$(git rev-parse --show-toplevel)
WORK=$ROOT/_output/win11-kvm                    # git-ignored: disk, ISO, the generated password
OUT=$ROOT/_output/reference
IMAGE=docker.io/dockurr/windows:6.06
mkdir -p "$WORK" "$OUT"
[ -f "$WORK/password" ] || (umask 077; head -c 18 /dev/urandom | base64 | tr -d '/+=' > "$WORK/password")
PASS=$(cat "$WORK/password")
ssh_win() { SSHPASS=$PASS sshpass -e ssh -p 2222 -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR xcpng@127.0.0.1 "$@"; }

case "${1:-}" in
  up)
    docker network inspect xcp-kvm-lab >/dev/null 2>&1 || { echo "start the KVM lab first (XCP_LAB_BACKEND=kvm lab.sh up)" >&2; exit 1; }
    docker rm -f xcpng-win11 >/dev/null 2>&1 || true
    docker run -d --name xcpng-win11 --network xcp-kvm-lab --device /dev/kvm --device /dev/net/tun --cap-add NET_ADMIN \
      -e VERSION=11 -e RAM_SIZE=8G -e CPU_CORES=4 -e DISK_SIZE=64G -e USERNAME=xcpng -e PASSWORD="$PASS" \
      -p 127.0.0.1:8006:8006 -p 127.0.0.1:2222:22 -v "$WORK":/storage -v "$PWD/oem":/oem --stop-timeout 120 "$IMAGE" >/dev/null
    echo "installing; follow it at http://127.0.0.1:8006" ;;
  shot)
    docker exec xcpng-win11 sh -c "echo 'screendump /storage/shot.ppm' | nc -q 1 -U /run/shm/monitor.sock >/dev/null"
    sleep 1; convert "$WORK/shot.ppm" "$OUT/$2.png" && echo "$OUT/$2.png" ;;
  ps) # the script runs as a scheduled task in the logged-on session (GUI apps need the interactive desktop)
    shift
    enc() { printf '%s' "$1" | iconv -t UTF-16LE | base64 -w0; }
    inner=$(enc "$*")
    outer=$(enc "\$ProgressPreference = 'SilentlyContinue'; \$a = New-ScheduledTaskAction -Execute powershell.exe -Argument '-NoProfile -WindowStyle Hidden -ExecutionPolicy Bypass -EncodedCommand $inner'; \$p = New-ScheduledTaskPrincipal -UserId xcpng -LogonType Interactive; Register-ScheduledTask -Force -TaskName xcpngui -Action \$a -Principal \$p | Out-Null; Start-ScheduledTask -TaskName xcpngui; Start-Sleep -Milliseconds 800; while ((Get-ScheduledTask -TaskName xcpngui).State -eq 'Running') { Start-Sleep -Milliseconds 300 }; Unregister-ScheduledTask -TaskName xcpngui -Confirm:\$false")
    ssh_win "powershell -NoProfile -EncodedCommand $outer" ;;
  key) shift; for k in "$@"; do docker exec xcpng-win11 sh -c "echo 'sendkey $k' | nc -q 1 -U /run/shm/monitor.sock >/dev/null"; done ;;  # QEMU key names: esc, ret, tab, alt-f4
  down) docker stop xcpng-win11 >/dev/null 2>&1 || true ;;
  *) echo "usage: $0 up | shot NAME | ps 'command' | key KEY... | down" >&2; exit 2 ;;
esac
