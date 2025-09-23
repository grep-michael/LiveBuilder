#!/bin/sh
#
# /usr/lib/live/boot/9991-overlayparser
# Runs after live-boot's 9990-overlay. Attempts to inject an extra desktop
# squashfs into the overlay by building a new overlay with the extra layer first,
# then moving it into place.
#
# Heavy debug output and final 'read' so you can inspect what happened.
#

# helper: bail to shell with message
drop_shell() {
    echo "[overlayparser] Dropping to shell for manual debugging..."
    /bin/sh || true
    exit 0
}

echo "[overlayparser] ===== START 9991-overlayparser ====="
sleep 1

# show cmdline
CMDLINE="$(cat /proc/cmdline 2>/dev/null || true)"
echo "[overlayparser] cmdline: $CMDLINE"
sleep 1

# pick DE param
DE_BASENAME=""
case "$CMDLINE" in
    *desktop=openbox*)
        DE_BASENAME="openbox.squashfs"
        DE_SUBPATH="minimal/openbox.squashfs"
        ;;
    *desktop=lxqt*)
        DE_BASENAME="lxqt.squashfs"
        DE_SUBPATH="medium/lxqt.squashfs"
        ;;
    *desktop=maximum*)
        DE_BASENAME="maximum.squashfs"
        DE_SUBPATH="maximum/maximum.squashfs"
        ;;
    *)
        echo "[overlayparser] No desktop= kernel parameter; nothing to inject."
        echo "[overlayparser] ===== END ====="
        read dummy
        exit 0
        ;;
esac

echo "[overlayparser] Requested extra: $DE_SUBPATH ($DE_BASENAME)"
sleep 1

# Attempt to locate squashfs under the mounted medium
DE_PATH=""
# Prefer scanning under run/live/medium first, then fallback to whole FS search
echo "[overlayparser] Scanning /run/live/medium for $DE_BASENAME ..."
FOUND="$(find /run/live/medium -type f -name "$DE_BASENAME" 2>/dev/null | head -n1 || true)"
if [ -n "$FOUND" ]; then
    DE_PATH="$FOUND"
    echo "[overlayparser] Found at: $DE_PATH"
else
    echo "[overlayparser] Not found under /run/live/medium; doing a wider search (this may be slow)..."
    FOUND="$(find / -type f -name "$DE_BASENAME" 2>/dev/null | head -n1 || true)"
    if [ -n "$FOUND" ]; then
        DE_PATH="$FOUND"
        echo "[overlayparser] Found (fallback) at: $DE_PATH"
    fi
fi

if [ -z "$DE_PATH" ]; then
    echo "[overlayparser] ERROR: Could not locate $DE_BASENAME anywhere."
    echo "[overlayparser] /run/live/medium contents (short):"
    ls -l /run/live/medium 2>/dev/null || echo "/run/live/medium not present"
    echo "[overlayparser] All squashfs under /run/live/medium (if any):"
    find /run/live/medium -name '*.squashfs' 2>/dev/null || true
    sleep 2
    echo "[overlayparser] Press ENTER to drop to shell and inspect, or Ctrl+C to continue boot."
    read dummy
    drop_shell
fi

sleep 1

# mount the extra squashfs read-only to a known place
EXTRA_MNT="/run/live/rootfs-extra"
mkdir -p "$EXTRA_MNT"
echo "[overlayparser] Mounting extra squashfs read-only at $EXTRA_MNT"
mount -t squashfs -o loop,ro "$DE_PATH" "$EXTRA_MNT" 2>&1 || {
    echo "[overlayparser] ERROR: mount of extra squashfs failed."
    echo "mount output:"
    mount | tail -n 20
    sleep 2
    drop_shell
}

echo "[overlayparser] Extra mounted. ls $EXTRA_MNT (short):"
ls -l "$EXTRA_MNT" | sed -n '1,80p'
sleep 1

# Locate existing overlay mount (we expect an overlay mount at /root)
echo "[overlayparser] Locating current overlay mountpoint (searching for type overlay and mountpoint /root)..."
# Look for a mount whose target is /root or whose fs type is overlay and may include /root
OLD_MOUNT_LINE="$(grep ' /root ' /proc/self/mounts | grep overlay || true)"
if [ -z "$OLD_MOUNT_LINE" ]; then
    # try searching for any overlay mount and prefer /root fallback
    OLD_MOUNT_LINE="$(grep ' overlay ' /proc/self/mounts || true)"
fi

echo "[overlayparser] mount lines relevant:"
echo "$OLD_MOUNT_LINE"
sleep 1

if [ -z "$OLD_MOUNT_LINE" ]; then
    echo "[overlayparser] ERROR: Could not find an existing overlay mount line in /proc/self/mounts."
    echo "[overlayparser] Current mounts:"
    cat /proc/self/mounts
    sleep 2
    drop_shell
fi

# Parse out lowerdir, upperdir and workdir from the mount options
# The mount line looks like: overlay /root overlay rw,lowerdir=/a:/b,upperdir=/cow,workdir=/cow/.work 0 0
MNT_OPTS="$(echo "$OLD_MOUNT_LINE" | awk '{print $4}')"
echo "[overlayparser] mount options: $MNT_OPTS"
# helper to extract opt
getoptval() {
    echo "$MNT_OPTS" | tr ',' '\n' | sed -n "s/^$1=//p" | head -n1
}

OLD_LOWERDIR="$(getoptval lowerdir)"
OLD_UPPERDIR="$(getoptval upperdir)"
OLD_WORKDIR="$(getoptval workdir)"

echo "[overlayparser] parsed old lowerdir: $OLD_LOWERDIR"
echo "[overlayparser] parsed old upperdir: $OLD_UPPERDIR"
echo "[overlayparser] parsed old workdir: $OLD_WORKDIR"
sleep 1

if [ -z "$OLD_LOWERDIR" ] || [ -z "$OLD_UPPERDIR" ] || [ -z "$OLD_WORKDIR" ]; then
    echo "[overlayparser] ERROR: could not parse required overlay options. Aborting."
    echo "mount line:"
    echo "$OLD_MOUNT_LINE"
    sleep 2
    drop_shell
fi

# Build new lowerdir string with extra first
# NOTE: overlay lowerdir is colon-separated list. Prepend extra mount path.
# The extra must present first so it appears on top.
NEW_LOWERDIR="$EXTRA_MNT:$OLD_LOWERDIR"
echo "[overlayparser] new lowerdir will be: $NEW_LOWERDIR"
sleep 1

# Create a new overlay mountpoint
NEW_ROOT="/run/live/newroot"
mkdir -p "$NEW_ROOT"
echo "[overlayparser] Mounting new overlay at $NEW_ROOT ..."
mount -t overlay overlay -o lowerdir="$NEW_LOWERDIR",upperdir="$OLD_UPPERDIR",workdir="$OLD_WORKDIR" "$NEW_ROOT" 2>&1 || {
    echo "[overlayparser] ERROR: failed to mount new overlay at $NEW_ROOT"
    echo "last mount output:"
    mount | tail -n 40
    sleep 2
    drop_shell
}

echo "[overlayparser] New overlay mounted at $NEW_ROOT. ls (short):"
ls -l "$NEW_ROOT" | sed -n '1,80p'
sleep 1

# Now swap it into place. We try to move the mount over /root.
echo "[overlayparser] Attempting to move new overlay into /root ..."
# The mount --move moves the entire mount tree from NEW_ROOT to /root
mount --move "$NEW_ROOT" /root 2>&1 || {
    echo "[overlayparser] ERROR: mount --move failed."
    echo "mounts:"
    cat /proc/self/mounts | tail -n 40
    sleep 2
    echo "[overlayparser] You can inspect /run/live/newroot and /run/live/rootfs-extra, then press ENTER to drop to shell."
    read dummy
    drop_shell
}

echo "[overlayparser] Successfully moved new overlay into /root."
echo "[overlayparser] Final /root listing (short):"
ls -l /root | sed -n '1,80p'
sleep 1

echo "[overlayparser] ===== FINISHED ====="
echo "Press ENTER to continue boot (or inspect with shell)..."
read dummy

# Optionally give you a shell after the pause for more inspection
echo "[overlayparser] Opening shell for final inspection..."
/bin/sh || true

exit 0
