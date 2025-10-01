#!/bin/bash
set -e

# Ensure storage drivers are in initramfs
cat >> /etc/initramfs-tools/modules << 'MODEOF'
ahci
libahci
raid_class
dm-mod
MODEOF