@echo off
rem Runs once at the end of the unattended Windows setup (dockur/windows copies this folder to C:\OEM).
rem Installs the official XCP-ng Center v25.04 (UI reference) and OpenSSH, so the host can drive the desktop.
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0;" ^
  "Set-Service sshd -StartupType Automatic; Start-Service sshd;" ^
  "New-NetFirewallRule -Name sshd -DisplayName sshd -Protocol TCP -LocalPort 22 -Action Allow;" ^
  "New-Item -ItemType Directory -Force C:\xcpng | Out-Null;" ^
  "Invoke-WebRequest -UseBasicParsing -Uri https://github.com/xcp-ng/xenadmin/releases/download/v25.04/XCP-ng.Center.zip -OutFile C:\xcpng\XCP-ng.Center.zip;" ^
  "Expand-Archive -Force C:\xcpng\XCP-ng.Center.zip C:\xcpng\app;" ^
  "Copy-Item C:\OEM\ui.ps1 C:\xcpng\ui.ps1" > C:\OEM\install.log 2>&1
