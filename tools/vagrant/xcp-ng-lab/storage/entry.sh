#!/bin/sh
# usage: entry.sh nfs | smb | iscsi   (SMB_PASS sets the password of the "xcp" share user)
set -e
case "$1" in
  nfs)
    mkdir -p /export /run/rpcbind /var/run/ganesha /var/lib/nfs/ganesha; rpcbind -w
    exec ganesha.nfsd -F -L /dev/stdout -f /etc/ganesha/ganesha.conf ;;
  smb)
    : "${SMB_PASS:?SMB_PASS required}"
    id xcp >/dev/null 2>&1 || useradd -M -s /usr/sbin/nologin xcp
    printf '%s\n%s\n' "$SMB_PASS" "$SMB_PASS" | smbpasswd -s -a xcp >/dev/null
    exec smbd --foreground --no-process-group ;;
  iscsi)
    # one target, one 4 GiB sparse LUN, open to any initiator (a throwaway lab)
    IQN=${ISCSI_IQN:-iqn.2026-11.org.xcp-ng.lab:sr}
    [ -f /lun/lun1.img ] || truncate -s 4G /lun/lun1.img
    tgtd -f & pid=$!
    for i in $(seq 20); do tgtadm --lld iscsi --mode target --op show >/dev/null 2>&1 && break; sleep 0.5; done
    tgtadm --lld iscsi --mode target --op new --tid 1 --targetname "$IQN"
    tgtadm --lld iscsi --mode logicalunit --op new --tid 1 --lun 1 --backing-store /lun/lun1.img
    tgtadm --lld iscsi --mode target --op bind --tid 1 --initiator-address ALL
    wait $pid ;;
esac
