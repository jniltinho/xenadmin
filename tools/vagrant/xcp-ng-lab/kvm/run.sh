#!/bin/sh
# Runs the lab hosts on a bridge that holds the container's own interface, so they sit on the Docker network
# itself and the host reaches them directly. usage: run.sh <gateway> <host numbers...>; files in /work.
set -e
GW=$1; shift
IP=$(ip -4 -o addr show eth0 | awk '{print $4}')
ip link add br0 type bridge
ip link set br0 up
ip addr flush dev eth0
ip link set eth0 master br0
ip addr add "$IP" dev br0
ip route add default via "$GW"
for n in "$@"; do
  ip tuntap add "tap$n" mode tap
  ip link set "tap$n" master br0 up
  qemu-img create -q -f qcow2 "/work/xcp-lab-kvm-$n.qcow2" 100G 2>/dev/null || true   # keeps an existing disk
  qemu-system-x86_64 -name "xcp-lab-$n" -accel kvm -cpu host -smp 4 -m 8192 -display none -vnc ":$n" \
    -device ahci,id=ahci \
    -drive "file=/work/xcp-lab-kvm-$n.qcow2,if=none,id=disk,cache=unsafe" -device ide-hd,drive=disk,bus=ahci.0 \
    -drive "file=/work/xcp-lab-kvm-$n.iso,if=none,id=cd,media=cdrom,readonly=on" -device ide-cd,drive=cd,bus=ahci.1 \
    -boot order=cd \
    -netdev "tap,id=net,ifname=tap$n,script=no,downscript=no" -device "e1000,netdev=net,mac=52:54:00:0a:0b:0$n" \
    -serial "file:/work/xcp-lab-kvm-$n-serial.log" &
done
wait
