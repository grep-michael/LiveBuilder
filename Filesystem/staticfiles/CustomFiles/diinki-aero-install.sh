
#!/bin/bash

# download rustup
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
. "$HOME/.cargo/env"

#donwload eww
git clone https://github.com/elkowar/eww.git
cd eww

cargo build --release

sudo cargo install

#hyprland install user
useradd -m -s /bin/bash installer
echo "installer ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/installer
chmod 440 /etc/sudoers.d/installer

# hyprland install 
#bash -c "yes | $(curl -sSL https://raw.githubusercontent.com/JaKooLit/Debian-Hyprland/main/auto-install.sh)"
sudo -u installer bash << 'EOF'
cd ~
git clone https://github.com/JaKooLit/Debian-Hyprland.git
cd Debian-Hyprland
chmod +x install.sh
yes | ./install.sh
EOF