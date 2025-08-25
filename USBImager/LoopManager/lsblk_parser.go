package loopmanager

import (
	"LiveBuilder/Utils"
	"encoding/json"
	"fmt"

	//"log"
	"regexp"
	"strconv"
)

// lsblkBlockDevice represents the JSON structure from lsblk --fs --json
type lsblkBlockDevice struct {
	Name        string             `json:"name"`
	FSType      *string            `json:"fstype"`
	FSVer       *string            `json:"fsver"`
	Label       *string            `json:"label"`
	UUID        *string            `json:"uuid"`
	FSAvail     *string            `json:"fsavail"`
	FSUse       *string            `json:"fsuse%"`
	MountPoints []string           `json:"mountpoints"`
	Children    []lsblkBlockDevice `json:"children,omitempty"`
}

// lsblkOutput represents the root JSON structure from lsblk
type lsblkOutput struct {
	BlockDevices []lsblkBlockDevice `json:"blockdevices"`
}

func ParseLsblkFilesystems(loopDevice string) ([]PartitionInfo, error) {
	stdout, stderr, err := Utils.Run_command("lsblk", "--fs", "--json", loopDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %v - %s", err, stderr)
	}
	//log.Printf("PARTITIONS DETECTED\n%s\n", stdout)
	return parseLsblkJSON([]byte(stdout))
}

// parseLsblkJSON parses the lsblk JSON output into PartitionInfo structs
func parseLsblkJSON(jsonData []byte) ([]PartitionInfo, error) {
	var lsblkData lsblkOutput

	if err := json.Unmarshal(jsonData, &lsblkData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var partitions []PartitionInfo

	// Process each block device
	for _, device := range lsblkData.BlockDevices {
		// Process children (partitions) of the device
		for _, child := range device.Children {
			partition := convertToPartitionInfo(child)
			partitions = append(partitions, partition)
		}
	}

	return partitions, nil
}

// convertToPartitionInfo converts lsblkBlockDevice to PartitionInfo
func convertToPartitionInfo(device lsblkBlockDevice) PartitionInfo {
	partition := PartitionInfo{
		DeviceName: device.Name,
		DevicePath: "/dev/" + device.Name,
	}

	// Extract partition number from device name (e.g., "loop0p1" -> 1)
	re := regexp.MustCompile(`p(\d+)$`)
	matches := re.FindStringSubmatch(partition.DeviceName)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			partition.PartitionNum = num
		}
	}

	// Convert pointer fields to strings, handling nil values
	if device.FSType != nil {
		partition.FSType = *device.FSType
	}
	if device.FSVer != nil {
		partition.FSVersion = *device.FSVer
	}
	if device.Label != nil {
		partition.Label = *device.Label
	}
	if device.UUID != nil {
		partition.UUID = *device.UUID
	}
	if device.FSAvail != nil {
		partition.FSAvailable = *device.FSAvail
	}
	if device.FSUse != nil {
		partition.FSUse = *device.FSUse
	}

	// Handle mount points array (filter out null entries)
	for _, mp := range device.MountPoints {
		if mp != "" {
			partition.MountPoints = append(partition.MountPoints, mp)
		}
	}

	return partition
}
