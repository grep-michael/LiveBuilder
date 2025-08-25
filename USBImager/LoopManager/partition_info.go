package loopmanager

import (
	"fmt"
	"strings"
)

type PartitionInfo struct {
	DeviceName   string   // e.g., "loop0p1"
	DevicePath   string   // e.g., "/dev/loop0p1"
	FSType       string   // e.g., "vfat", "ext4"
	FSVersion    string   // e.g., "FAT32", "1.0"
	Label        string   // e.g., "BOOT", "SYSTEM"
	UUID         string   // e.g., "5D8D-ED10"
	FSAvailable  string   // e.g., "450M" (if mounted)
	FSUse        string   // e.g., "12%" (if mounted)
	MountPoints  []string // e.g., ["/mnt/boot"] (if mounted)
	PartitionNum int      // partition number (1, 2, 3, etc.)
}

func (p *PartitionInfo) IsMounted() bool {
	return len(p.MountPoints) > 0
}

func (p *PartitionInfo) GetMountPointsString() string {
	if len(p.MountPoints) == 0 {
		return ""
	}

	return strings.Join(p.MountPoints, ", ")
}

func (p *PartitionInfo) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Device: %s", p.DevicePath))

	if p.FSType != "" {
		parts = append(parts, fmt.Sprintf("Type: %s", p.FSType))
	}
	if p.FSVersion != "" {
		parts = append(parts, fmt.Sprintf("Version: %s", p.FSVersion))
	}
	if p.Label != "" {
		parts = append(parts, fmt.Sprintf("Label: %s", p.Label))
	}
	if p.UUID != "" {
		parts = append(parts, fmt.Sprintf("UUID: %s", p.UUID))
	}
	if p.FSAvailable != "" {
		parts = append(parts, fmt.Sprintf("Available: %s", p.FSAvailable))
	}
	if p.FSUse != "" {
		parts = append(parts, fmt.Sprintf("Use: %s", p.FSUse))
	}
	if p.IsMounted() {
		parts = append(parts, fmt.Sprintf("Mounted: %s", p.GetMountPointsString()))
	}

	return strings.Join(parts, ", ")
}
