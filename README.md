# LiveBuilder

A powerful Golang tool for creating custom Debian live ISOs using the debian-live/live-build framework. LiveBuilder simplifies the process of building personalized Debian-based live systems with automated configuration management.

## Prerequisites

Before using LiveBuilder, ensure you have the following installed:

- **Go 1.19+**: [Download and install Go](https://golang.org/dl/)
- **live-build**: Debian's live-build package
- **debootstrap**: For creating base Debian systems
- **Root privileges**: Required for some live-build operations

### Installing Prerequisites on Debian/Ubuntu

```bash
sudo apt update
sudo apt install live-build debootstrap squashfs-tools xorriso isolinux syslinux-efi grub-pc-bin grub-efi-amd64-bin mtools
```

## 🔧 Installation

### From Source

```bash
git clone https://github.com/grep-michael/LiveBuilder.git
cd LiveBuilder
go build -o livebuilder
sudo mv livebuilder /usr/local/bin/
```


## Customizing

Custom files are in ~/.config/LiveBuilder/CustomFiles
Each custom file should include a .meta.json to specify where it should be position in the live-build structure and to tag it

## 🙏 Acknowledgments

- [Debian Live Project](https://www.debian.org/devel/debian-live/) for the live-build framework
- [fyne](https://fyne.io/) Fyne gui
- [Go Community](https://golang.org/) for the excellent programming language
- All contributors who have helped improve this project

## 📚 Additional Resources

- [Debian Live Manual](https://live-team.pages.debian.net/live-manual/html/live-manual/index.en.html)
- [live-build Documentation](https://manpages.debian.org/testing/live-build/live-build.7.en.html)
- [Go Documentation](https://golang.org/doc/)
- [Debian Package Management](https://www.debian.org/doc/manuals/debian-reference/ch02.en.html)

---

**Keywords**: debian, live-build, iso, golang, custom-debian, live-system, debian-live, iso-builder, linux-distribution, bootable-iso, debian-customization, live-cd, live-usb
