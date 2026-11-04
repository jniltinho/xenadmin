# Runs inside the interactive session (via a scheduled task): click X Y | dclick X Y | rclick X Y | keys SENDKEYS
param([string]$Cmd, [string]$A, [string]$B)
Add-Type -AssemblyName System.Windows.Forms
Add-Type @"
using System; using System.Runtime.InteropServices;
public static class U {
  [DllImport("user32.dll")] public static extern bool SetCursorPos(int x, int y);
  [DllImport("user32.dll")] public static extern void mouse_event(uint f, uint x, uint y, uint d, IntPtr e);
}
"@
function Click([int]$x, [int]$y, [uint32]$down = 2, [uint32]$up = 4) {
  [U]::SetCursorPos($x, $y) | Out-Null; Start-Sleep -Milliseconds 150
  [U]::mouse_event($down, 0, 0, 0, [IntPtr]::Zero); [U]::mouse_event($up, 0, 0, 0, [IntPtr]::Zero)
}
switch ($Cmd) {
  "click"  { Click $A $B }
  "dclick" { Click $A $B; Start-Sleep -Milliseconds 80; Click $A $B }
  "rclick" { Click $A $B 8 16 }
  "keys"   { [System.Windows.Forms.SendKeys]::SendWait($A) }
}
