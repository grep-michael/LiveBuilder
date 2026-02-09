#!/bin/bash
set -e

KIOSK_USER="kiosk"
KIOSK_HOME="/home/${KIOSK_USER}"
KIOSK_URL="https://www.mcdonalds.com/us/en-us/full-menu.html"

apt update

apt install -y \
    unclutter \
    xorg \
    chromium \
    openbox \
    lightdm \
    lightdm-gtk-greeter \
    locales \
    network-manager \
    wpasupplicant \
    dbus \
    wireless-tools

# -------------------------
# Create kiosk user
# -------------------------
if ! id "${KIOSK_USER}" >/dev/null 2>&1; then
    useradd -m -s /bin/bash "${KIOSK_USER}"
fi

# Allow LightDM autologin
groupadd -f autologin
usermod -aG autologin "${KIOSK_USER}"

# Password not required for autologin
passwd -d "${KIOSK_USER}" || true

# -------------------------
# Openbox autostart
# -------------------------
mkdir -p "${KIOSK_HOME}/.config/openbox"

cat > "${KIOSK_HOME}/.config/openbox/autostart" << EOF
#!/bin/bash
xset s off
xset s noblank
xset -dpms
unclutter -idle 0.1 -grab -root &

while true; do
  xrandr --auto
  chromium \
    --noerrdialogs \
    --no-memcheck \
    --no-first-run \
    --start-maximized \
    --disable \
    --disable-translate \
    --disable-infobars \
    --disable-suggestions-service \
    --disable-save-password-bubble \
    --disable-session-crashed-bubble \
    --incognito \
    --kiosk "${KIOSK_URL}"
  sleep 5
done &
EOF

chmod +x "${KIOSK_HOME}/.config/openbox/autostart"
chown -R "${KIOSK_USER}:${KIOSK_USER}" "${KIOSK_HOME}/.config"

# -------------------------
# LightDM config
# -------------------------
cat > /etc/lightdm/lightdm.conf << EOF
[LightDM]
greeter-session=lightdm-gtk-greeter

[Seat:*]
autologin-user=${KIOSK_USER}
autologin-user-timeout=0
autologin-session=openbox
xserver-command=/usr/bin/X -nocursor -nolisten tcp
EOF

systemctl enable lightdm

# -------------------------
# NetworkManager
# -------------------------
systemctl enable NetworkManager.service
systemctl disable networking.service 2>/dev/null || true

sed -i '/^iface .* inet/d' /etc/network/interfaces 2>/dev/null || true

mkdir -p /etc/NetworkManager/conf.d
cat > /etc/NetworkManager/conf.d/10-globally-managed.conf <<EOF
[keyfile]
unmanaged-devices=none
EOF

mkdir -p /etc/NetworkManager/system-connections
cat > /etc/NetworkManager/system-connections/wifi.nmconnection <<EOF
[connection]
id=wifi
type=wifi
autoconnect=true

[wifi]
mode=infrastructure
ssid=ITAD_TEST

[wifi-security]
key-mgmt=wpa-psk
psk=p@55word

[ipv4]
method=auto

[ipv6]
method=ignore
EOF

chmod 600 /etc/NetworkManager/system-connections/wifi.nmconnection

mkdir -p /etc/X11/xorg.conf.d/
cat > /etc/X11/xorg.conf.d/10-no-dpms.conf <<EOF
Section "ServerFlags"
    Option "BlankTime" "0"
    Option "StandbyTime" "0"
    Option "SuspendTime" "0"
    Option "OffTime" "0"
EndSection
EOF