# NFS (nfs-ganesha, userspace: no kernel nfsd on the host), SMB (samba) and iSCSI (tgt, userspace: no LIO modules
# on the host) servers for the lab's shared SRs.
FROM debian:12
RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends \
      nfs-ganesha nfs-ganesha-vfs rpcbind samba tgt && rm -rf /var/lib/apt/lists/*
COPY ganesha.conf /etc/ganesha/ganesha.conf
COPY smb.conf /etc/samba/smb.conf
COPY entry.sh /entry.sh
ENTRYPOINT ["/bin/sh", "/entry.sh"]
