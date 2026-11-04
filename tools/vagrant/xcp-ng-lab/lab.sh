#!/usr/bin/env bash
# Two nested XCP-ng 8.3 hosts in VirtualBox for pool and write tests (Windows and desktop apps).
# XCP-ng has no Vagrant box, so the hosts are installed unattended from a remastered ISO:
# each ISO boots straight into the installer with an answerfile served from the host.
#   XCP_LAB_PASS=... ./lab.sh up        # creates and installs xcp-lab-1/2 (192.168.56.11/.12)
#   ./lab.sh status | halt | start | destroy
set -euo pipefail
cd "$(dirname "$0")"
ROOT=$(git rev-parse --show-toplevel)
WORK=$ROOT/_output/xcp-lab                     # git-ignored: ISOs, answerfiles (password hash), disks, logs
ISO_URL=${XCP_LAB_ISO_URL:-https://mirrors.xcp-ng.org/isos/8.3/xcp-ng-8.3.0-20260806.iso}
ISO=$WORK/$(basename "$ISO_URL")
NET=${XCP_LAB_NET:-vboxnet0}                    # host-only network, host is 192.168.56.1
PORT=8099
HOSTS=(1 2)
ip_of() { echo "192.168.56.1$1"; }
vm_of() { echo "xcp-lab-$1"; }
mkdir -p "$WORK"

remaster() { # $1 = host number; boot menu default -> unattended install from http answerfile
  local n=$1 iso=$WORK/xcp-lab-$1.iso cfg=$WORK/isolinux-$1.cfg
  [ -f "$iso" ] && return
  xorriso -osirrox on -indev "$ISO" -extract /boot/isolinux/isolinux.cfg "$cfg" 2>/dev/null
  chmod u+w "$cfg"
  sed -i -e 's/^TIMEOUT .*/TIMEOUT 10/' \
    -e "0,/--- \/boot\/vmlinuz console=hvc0 console=tty0 --- \/install.img/s##--- /boot/vmlinuz console=hvc0 console=tty0 install answerfile=http://192.168.56.1:$PORT/lab-$n.xml answerfile_device=eth0 --- /install.img#" "$cfg"
  grep -q "lab-$n.xml" "$cfg" || { echo "isolinux.cfg layout changed" >&2; exit 1; }
  xorriso -indev "$ISO" -outdev "$iso" -map "$cfg" /boot/isolinux/isolinux.cfg -boot_image any replay 2>/dev/null
}

answerfile() { # $1 = host number
  cat > "$WORK/www/lab-$1.xml" <<XML
<?xml version="1.0"?>
<installation srtype="ext">
  <primary-disk>sda</primary-disk>
  <keymap>us</keymap>
  <root-password type="hash">$HASH</root-password>
  <source type="local"/>
  <admin-interface name="eth0" proto="static">
    <ipaddr>$(ip_of "$1")</ipaddr>
    <subnet>255.255.255.0</subnet>
    <gateway>192.168.56.1</gateway>
  </admin-interface>
  <nameserver>192.168.56.1</nameserver>
  <timezone>America/Sao_Paulo</timezone>
  <hostname>$(vm_of "$1")</hostname>
</installation>
XML
}

create() { # $1 = host number
  local vm; vm=$(vm_of "$1")
  VBoxManage showvminfo "$vm" >/dev/null 2>&1 && return
  VBoxManage createvm --name "$vm" --ostype Linux_64 --register --basefolder "$WORK/vms"
  VBoxManage modifyvm "$vm" --memory 8192 --cpus 4 --nested-hw-virt on --paravirtprovider none \
    --firmware bios --boot1 disk --boot2 dvd --boot3 none --boot4 none --ioapic on \
    --nic1 hostonly --hostonlyadapter1 "$NET" --nictype1 82540EM --macaddress1 0800270A0B0"$1" \
    --uart1 0x3F8 4 --uartmode1 file "$WORK/$vm-serial.log" --audio-enabled off --graphicscontroller vmsvga
  VBoxManage storagectl "$vm" --name SATA --add sata --portcount 2
  VBoxManage createmedium disk --filename "$WORK/vms/$vm/$vm.vdi" --size 102400
  VBoxManage storageattach "$vm" --storagectl SATA --port 0 --type hdd --medium "$WORK/vms/$vm/$vm.vdi"
  VBoxManage storageattach "$vm" --storagectl SATA --port 1 --type dvddrive --medium "$WORK/xcp-lab-$1.iso"
}

wait_xapi() { # $1 = ip; XAPI answers on 443 once the installed host has booted
  for _ in $(seq 1 120); do
    curl -sk -o /dev/null "https://$1/" && return 0
    sleep 10
  done
  return 1
}

case "${1:-}" in
  up)
    : "${XCP_LAB_PASS:?set XCP_LAB_PASS (root password of the lab hosts)}"
    [ -f "$ISO" ] || curl -fL -o "$ISO" "$ISO_URL"
    mkdir -p "$WORK/www"
    HASH=$(openssl passwd -6 -stdin <<<"$XCP_LAB_PASS")
    for n in "${HOSTS[@]}"; do remaster "$n"; answerfile "$n"; create "$n"; done
    for n in "${HOSTS[@]}"; do VBoxManage startvm "$(vm_of "$n")" --type headless || true; done
    for _ in $(seq 1 30); do ip -4 addr show "$NET" | grep -q 192.168.56.1/ && break; sleep 2; done
    python3 -m http.server "$PORT" --bind 192.168.56.1 --directory "$WORK/www" >"$WORK/http.log" 2>&1 &
    trap 'kill $! 2>/dev/null; rm -rf "$WORK/www"' EXIT   # answerfiles only live during the install
    for n in "${HOSTS[@]}"; do
      wait_xapi "$(ip_of "$n")" && echo "$(vm_of "$n") ready at https://$(ip_of "$n")" \
        || { echo "$(vm_of "$n") did not come up; see $WORK/$(vm_of "$n")-serial.log" >&2; exit 1; }
    done
    for n in "${HOSTS[@]}"; do VBoxManage storageattach "$(vm_of "$n")" --storagectl SATA --port 1 --type dvddrive --medium emptydrive; done
    ;;
  start)   for n in "${HOSTS[@]}"; do VBoxManage startvm "$(vm_of "$n")" --type headless || true; done ;;
  halt)    for n in "${HOSTS[@]}"; do VBoxManage controlvm "$(vm_of "$n")" acpipowerbutton || true; done ;;
  status)  for n in "${HOSTS[@]}"; do printf '%s %s ' "$(vm_of "$n")" "$(ip_of "$n")"; VBoxManage showvminfo "$(vm_of "$n")" --machinereadable 2>/dev/null | grep '^VMState=' || echo absent; done ;;
  destroy) for n in "${HOSTS[@]}"; do VBoxManage controlvm "$(vm_of "$n")" poweroff 2>/dev/null || true; VBoxManage unregistervm "$(vm_of "$n")" --delete 2>/dev/null || true; done ;;
  *) echo "usage: XCP_LAB_PASS=... $0 up | start | halt | status | destroy" >&2; exit 2 ;;
esac
