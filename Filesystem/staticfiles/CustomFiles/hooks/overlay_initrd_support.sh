#!/bin/sh
PREREQ=""
prereqs() { echo "$PREREQ"; }
case $1 in
    prereqs) prereqs; exit 0;;
esac

. /scripts/live-functions

CMDLINE=$(cat /proc/cmdline)
OVERLAY=$(echo "$CMDLINE" | sed -n 's/.*overlay=\([^ ]*\).*/\1/p')

if [ -n "$OVERLAY" ]; then
    echo "live-boot: overlay list = $OVERLAY"

    mkdir -p /run/live/overlay
    idx=0
    OLDIFS=$IFS
    IFS=:

    for f in $OVERLAY; do
        echo "Mounting layer: $f"
        mkdir -p /run/live/overlay/$idx
        mount -t squashfs -o loop "/run/live/$f" "/run/live/overlay/$idx" 2>/dev/null || \
        mount -t squashfs -o loop "/run/live/findiso/live/$f" "/run/live/overlay/$idx"
        idx=$((idx+1))
    done
    IFS=$OLDIFS

    mkdir -p /run/live/rootfs/{upper,work}
    LOWERDIRS=$(printf "/run/live/overlay/%s:" $(seq 0 $((idx-1))))
    LOWERDIRS=${LOWERDIRS%:}

    mount -t overlay overlay \
        -o lowerdir=$LOWERDIRS,upperdir=/run/live/rootfs/upper,workdir=/run/live/rootfs/work \
        /root
fi
echo "OVERLAY PARSER FINISHED"
sleep 3