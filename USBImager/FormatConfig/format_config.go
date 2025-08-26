package formatconfig

import (
	"encoding/json"
	"fmt"
)

// Grub represents GRUB bootloader configuration
type Grub struct {
	Redirect string `json:"Redirect"`
	BIOS     bool   `json:"BIOS,omitempty"`
	UEFI     bool   `json:"UEFI,omitempty"`
}

// Partition represents a disk partition with its properties
type Partition struct {
	Label     string `json:"label"`
	StartAt   int64  `json:"StartAt,omitempty"` // omitempty for optional fields
	Size      int64  `json:"Size,omitempty"`    // omitempty for optional fields
	Type      string `json:"Type"`
	Bootable  bool   `json:"Bootable,omitempty"` // omitempty for optional fields
	LoopPath  string `json:"LoopPath"`
	MountPath string `json:"MountPath"`
	Grub      *Grub  `json:"Grub,omitempty"` // pointer for optional nested struct
}

// DiskImage represents the root disk image structure
type DiskImage struct {
	Label      string       `json:"label"`
	LabelID    string       `json:"label_id"`
	Units      string       `json:"units"`
	InFile     string       `json:"inFile"`
	OutFile    string       `json:"outFile"`
	Partitions []*Partition `json:"partitions"`
}

func NewDiskImageFromJSON(jsonString string) (*DiskImage, error) {
	var diskImage DiskImage
	err := json.Unmarshal([]byte(jsonString), &diskImage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &diskImage, nil
}

func NewDiskImage(label, labelID, units, inFile, outFile string) *DiskImage {
	return &DiskImage{
		Label:      label,
		LabelID:    labelID,
		Units:      units,
		InFile:     inFile,
		OutFile:    outFile,
		Partitions: make([]*Partition, 0),
	}
}

func (di *DiskImage) AddPartition(partition *Partition) {
	di.Partitions = append(di.Partitions, partition)
}

func NewPartition(label, partType, loopPath, mountPath string) *Partition {
	return &Partition{
		Label:     label,
		Type:      partType,
		LoopPath:  loopPath,
		MountPath: mountPath,
	}
}

func main() {
	// Example JSON data
	jsonData := `{
		"label":"dos",
		"label_id":"0x12345678",
		"units":"sectors",
		"inFile":"/path",
		"outFile":"/path",
		"partitions":[
			{
				"label":"BOOT",
				"StartAt":2048,
				"Size":100000,
				"Type":"W95 FAT32 (LBA)",
				"Bootable":true,
				"LoopPath":"",
				"MountPath":"",
				"Grub":{
					"Redirect":"SYSTEM"
				}
			},
			{
				"label":"SYSTEM",
				"Type":"Linux",
				"LoopPath":"",
				"MountPath":""
			}
		]
	}`

	// Create DiskImage from JSON string
	diskImage, err := NewDiskImageFromJSON(jsonData)
	if err != nil {
		fmt.Printf("Error creating DiskImage: %v\n", err)
		return
	}

	// Display parsed data
	fmt.Printf("Disk Image: %s (%s)\n", diskImage.Label, diskImage.LabelID)
	fmt.Printf("Files: %s -> %s\n", diskImage.InFile, diskImage.OutFile)
	fmt.Printf("Units: %s\n", diskImage.Units)
	fmt.Printf("Partitions: %d\n", len(diskImage.Partitions))

	for i, partition := range diskImage.Partitions {
		fmt.Printf("  [%d] %s (%s)\n", i, partition.Label, partition.Type)
		if partition.StartAt > 0 {
			fmt.Printf("      Start: %d, Size: %d\n", partition.StartAt, partition.Size)
		}
		if partition.Bootable {
			fmt.Printf("      Bootable: true\n")
		}
		if partition.Grub != nil {
			fmt.Printf("      GRUB Redirect: %s\n", partition.Grub.Redirect)
		}
	}

	// Convert back to JSON
	output, err := json.MarshalIndent(*diskImage, "", "  ")
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return
	}

	fmt.Printf("\nJSON Output:\n%s\n", string(output))
}
