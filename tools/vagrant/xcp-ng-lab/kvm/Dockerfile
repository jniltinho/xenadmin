# QEMU/KVM for the nested XCP-ng lab: KVM exposes NPT to nested Xen (VirtualBox on AMD does not), so guests can run.
FROM ubuntu:24.04
RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get install -y -qq --no-install-recommends \
      qemu-system-x86 qemu-utils iproute2 && rm -rf /var/lib/apt/lists/*
COPY run.sh /run.sh
ENTRYPOINT ["/bin/sh", "/run.sh"]
