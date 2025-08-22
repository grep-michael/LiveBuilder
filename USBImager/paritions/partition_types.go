package paritions

import ()

// PartitionType represents a partition type code
type PartitionType string

// Partition type constants
const (
	FAT12             PartitionType = "01"
	FAT16Small        PartitionType = "04"
	Extended          PartitionType = "05"
	FAT16             PartitionType = "06"
	HPFS_NTFS_exFAT   PartitionType = "07"
	W95_FAT32         PartitionType = "0B"
	W95_FAT32_LBA     PartitionType = "0C"
	W95_Extended_LBA  PartitionType = "0F"
	HiddenFAT12       PartitionType = "11"
	HiddenFAT16       PartitionType = "16"
	HiddenHPFS_NTFS   PartitionType = "17"
	LinuxSwap_Solaris PartitionType = "82"
	Linux             PartitionType = "83"
	LinuxExtended     PartitionType = "85"
	LinuxLVM          PartitionType = "8E"
	FreeBSD           PartitionType = "A5"
	OpenBSD           PartitionType = "A6"
	NeXTSTEP          PartitionType = "A7"
	UFS_Darwin        PartitionType = "A8"
	DarwinBoot        PartitionType = "AB"
	EFI_FAT           PartitionType = "EF"
	LinuxRAID         PartitionType = "FD"
)

// PartitionNameToCode maps partition type names to their codes
var PartitionNameToCode = map[string]PartitionType{
	"FAT12":                 FAT12,
	"FAT16 (<32M)":          FAT16Small,
	"Extended":              Extended,
	"FAT16":                 FAT16,
	"HPFS/NTFS/exFAT":       HPFS_NTFS_exFAT,
	"W95 FAT32":             W95_FAT32,
	"W95 FAT32 (LBA)":       W95_FAT32_LBA,
	"W95 Extended (LBA)":    W95_Extended_LBA,
	"Hidden FAT12":          HiddenFAT12,
	"Hidden FAT16":          HiddenFAT16,
	"Hidden HPFS/NTFS":      HiddenHPFS_NTFS,
	"Linux swap / Solaris":  LinuxSwap_Solaris,
	"Linux":                 Linux,
	"Linux extended":        LinuxExtended,
	"Linux LVM":             LinuxLVM,
	"FreeBSD":               FreeBSD,
	"OpenBSD":               OpenBSD,
	"NeXTSTEP":              NeXTSTEP,
	"UFS Darwin":            UFS_Darwin,
	"Darwin boot":           DarwinBoot,
	"EFI (FAT-12/16/32)":    EFI_FAT,
	"Linux raid autodetect": LinuxRAID,
}

// PartitionsCodeToName maps codes back to their descriptive names
var PartitionsCodeToName = map[PartitionType]string{
	FAT12:             "FAT12",
	FAT16Small:        "FAT16 (<32M)",
	Extended:          "Extended (contains logical partitions)",
	FAT16:             "FAT16",
	HPFS_NTFS_exFAT:   "HPFS/NTFS/exFAT",
	W95_FAT32:         "W95 FAT32",
	W95_FAT32_LBA:     "W95 FAT32 (LBA)",
	W95_Extended_LBA:  "W95 Extended (LBA)",
	HiddenFAT12:       "Hidden FAT12",
	HiddenFAT16:       "Hidden FAT16",
	HiddenHPFS_NTFS:   "Hidden HPFS/NTFS",
	LinuxSwap_Solaris: "Linux swap / Solaris",
	Linux:             "Linux",
	LinuxExtended:     "Linux extended",
	LinuxLVM:          "Linux LVM",
	FreeBSD:           "FreeBSD",
	OpenBSD:           "OpenBSD",
	NeXTSTEP:          "NeXTSTEP",
	UFS_Darwin:        "UFS Darwin",
	DarwinBoot:        "Darwin boot",
	EFI_FAT:           "EFI (FAT-12/16/32)",
	LinuxRAID:         "Linux raid autodetect",
}
