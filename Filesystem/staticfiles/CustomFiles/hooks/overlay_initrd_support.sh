#!/bin/sh
# /usr/lib/live/boot/9991-extrafs
set +e


LOG_FILE="/run/extrafs-boot.log"
log_msg() {
    echo "$@"
    if command -v tee >/dev/null 2>&1; then
        echo "$@" >> "$LOG_FILE"
    else
        # Manual append if tee not available
        echo "$@" >> "$LOG_FILE" 2>/dev/null || true
    fi
}

log_msg "[extrafs] START"

CMDLINE="$(cat /proc/cmdline)"

# Extract extrafs= parameter (colon-separated list)
EXTRAFS="$(log_msg "$CMDLINE" | sed -n 's/.*extrafs=\([^ ]*\).*/\1/p')"

if [ -z "$EXTRAFS" ]; then
    log_msg "[extrafs] No extrafs= found on kernel cmdline"
    exit 0
fi

log_msg "[extrafs] extrafs= string: $EXTRAFS"

EXTRAFS_SQUASHFS_LAYERS="/run/live/extrafs_squashfs_layers"
EXTRAFS_WORK_DIR="/run/live/extrafs_work"
EXTRAFS_UPPER_DIR="/run/live/extrafs_upper"

mkdir -p "$EXTRAFS_SQUASHFS_LAYERS"
mkdir -p "$EXTRAFS_WORK_DIR"
mkdir -p "$EXTRAFS_UPPER_DIR"

idx=0
OLDIFS=$IFS
IFS=:
EXTRA_FS_LOWER_DIRS=""
for f in $EXTRAFS; do
    log_msg "[extrafs] Candidate squashfs: $f"
    #Find extra squashfs files
    FOUND="$(find /run/live -type f -name "$(basename "$f")" 2>/dev/null | head -n1)"
    if [ -z "$FOUND" ]; then
        log_msg "[extrafs] Not under /run/live, searching whole FS (slow)..."
        FOUND="$(find / -type f -name "$(basename "$f")" 2>/dev/null | head -n1)"
    fi
    #mount found squashfs
    if [ -n "$FOUND" ]; then
        log_msg "[extrafs] Located: $FOUND"
        LAYER_DIR="${EXTRAFS_SQUASHFS_LAYERS}/layer_${idx}"
        mkdir -p "$LAYER_DIR"
        log_msg "[extrafs] Mounting layer $idx: $FOUND"
        mount -t squashfs -o loop,ro "$FOUND" "$LAYER_DIR"

        if [ -z "$EXTRA_FS_LOWER_DIRS" ]; then
            EXTRA_FS_LOWER_DIRS="$LAYER_DIR"
        else
            EXTRA_FS_LOWER_DIRS="$LAYER_DIR:$EXTRA_FS_LOWER_DIRS"
        fi
        
        idx=$((idx + 1))

    else
        log_msg "[extrafs] ERROR: $f not found anywhere."
    fi
done
IFS=$OLDIFS
ls -l "$EXTRAFS_SQUASHFS_LAYERS"
#create our overlay
EXTRA_MNT="/run/live/rootfs-extra"
mkdir -p "$EXTRA_MNT"
mount -t overlay overlay -o lowerdir="$EXTRA_FS_LOWER_DIRS",upperdir="$EXTRAFS_UPPER_DIR",workdir="$EXTRAFS_WORK_DIR" "$EXTRA_MNT" || true
log_msg "[extrafs] finished mounting extrafs_overlay"

OLD_MOUNT_LINE="$(grep ' /root ' /proc/self/mounts | grep overlay || true)"
if [ -z "$OLD_MOUNT_LINE" ]; then
    OLD_MOUNT_LINE="$(grep ' overlay ' /proc/self/mounts || true)"
fi

MNT_OPTS="$(log_msg "$OLD_MOUNT_LINE" | awk '{print $4}')"
log_msg "[extrafs] mount options: $MNT_OPTS"
getoptval() {
    log_msg "$MNT_OPTS" | tr ',' '\n' | sed -n "s/^$1=//p" | head -n1
}

OLD_LOWERDIR="$(getoptval lowerdir)"
OLD_UPPERDIR="$(getoptval upperdir)"
OLD_WORKDIR="$(getoptval workdir)"
log_msg "[extrafs] parsed old lowerdir: $OLD_LOWERDIR"
log_msg "[extrafs] parsed old upperdir: $OLD_UPPERDIR"
log_msg "[extrafs] parsed old workdir: $OLD_WORKDIR"

NEW_LOWERDIR="$EXTRA_MNT:$OLD_LOWERDIR"
NEW_ROOT="/run/live/newroot"
mkdir -p "$NEW_ROOT"
log_msg "[extrafs] Mounting new overlay at $NEW_ROOT ..."

mount -t overlay overlay -o lowerdir="$NEW_LOWERDIR",upperdir="$OLD_UPPERDIR",workdir="$OLD_WORKDIR" "$NEW_ROOT" 2>&1 || {
    log_msg "[extrafs] ERROR: failed to mount new overlay at $NEW_ROOT"
    /bin/sh
    exit 1
}

log_msg "[extrafs] Moving /root to $NEW_ROOT"
mount --move "$NEW_ROOT" /root 2>&1 || {
    log_msg "[extrafs] Failed to move root"
    /bin/sh
    exit 1
}
log_msg "[extrafs] Successfully moved new overlay into /root."
log_msg "[extrafs] Final /root listing (short):"
ls -l /root | sed -n '1,80p'
log_msg "[extrafs] END"
read dummy