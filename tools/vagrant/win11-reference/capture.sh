#!/bin/sh
# Launches the official XCP-ng Center inside the Windows 11 VM's interactive desktop,
# connects to a host and saves screenshots to _output/reference/ (never committed).
# usage: XCP_PASS=... ./capture.sh <host> [user]      step: ./capture.sh shot <name> | key <scancodes> | type <text>
set -e
cd "$(dirname "$0")"
VM=xcpng-center-win11-reference
OUT=../../../_output/reference
mkdir -p "$OUT"

shot() { VBoxManage controlvm "$VM" screenshotpng "$OUT/$1.png" && echo "$OUT/$1.png"; }

case "$1" in
  shot) shot "$2"; exit ;;
  key)  shift; VBoxManage controlvm "$VM" keyboardputscancode "$@"; exit ;;
  type) VBoxManage controlvm "$VM" keyboardputstring "$2"; exit ;;
  ui)   # ui click X Y | ui keys '{ENTER}'  (runs ui.ps1 in the logged-on desktop)
    shift; ARGS=$(printf ' "%s"' "$@" | sed "s/'/''/g")
    vagrant powershell -e -c "
Copy-Item -Force C:\vagrant\ui.ps1 C:\xcpng\ui.ps1
\$a = New-ScheduledTaskAction -Execute powershell.exe -Argument ('-NoProfile -WindowStyle Hidden -ExecutionPolicy Bypass -File C:\xcpng\ui.ps1' + '$ARGS')
\$p = New-ScheduledTaskPrincipal -UserId vagrant -LogonType Interactive
Register-ScheduledTask -Force -TaskName xcpngui -Action \$a -Principal \$p | Out-Null
Start-ScheduledTask -TaskName xcpngui; Start-Sleep -Milliseconds 500
while ((Get-ScheduledTask -TaskName xcpngui).State -eq 'Running') { Start-Sleep -Milliseconds 300 }
Unregister-ScheduledTask -TaskName xcpngui -Confirm:\$false" >/dev/null 2>&1
    exit ;;
esac

HOST=${1:?host}; XUSER=${2:-root}; : "${XCP_PASS:?set XCP_PASS}"
PW=$(printf %s "$XCP_PASS" | sed "s/'/''/g")   # PowerShell single-quote escaping
VBoxManage controlvm "$VM" setvideomodehint 1440 900 32 || true
# a scheduled task is the simplest way to start a GUI app in the logged-on session from WinRM
vagrant powershell -e -c "
\$pw = '$PW'
\$exe = (Get-ChildItem -Recurse C:\\xcpng\\app -Filter 'XCP-ng Center.exe' | Select-Object -First 1).FullName
Get-Process 'XCP-ng Center' -ErrorAction SilentlyContinue | Stop-Process -Force
\$a = New-ScheduledTaskAction -Execute \$exe -Argument ('connect $HOST $XUSER ' + [char]34 + \$pw + [char]34)
\$p = New-ScheduledTaskPrincipal -UserId vagrant -LogonType Interactive
Register-ScheduledTask -Force -TaskName xcpng -Action \$a -Principal \$p | Out-Null
Start-ScheduledTask -TaskName xcpng; Start-Sleep 3; Unregister-ScheduledTask -TaskName xcpng -Confirm:\$false
" 2>&1 | grep -v '^$' || true
sleep 25
shot 01-after-connect
