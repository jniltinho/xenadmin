#!/bin/sh
# usage: entry.sh nfs | smb   (SMB_PASS sets the password of the "xcp" share user)
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
esac
