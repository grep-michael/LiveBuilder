package usbimager

import (
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	"strconv"
)

func CreatePartitionareForConfig(config *formatconfig.DiskImage) *DiskPartitionare {
	partitionTable := buildPartitionTable(config)

	diskpart := NewDiskPartionare(config.OutFile)
	diskpart.SetPartitionTable(partitionTable)

	return diskpart
}

func buildPartitionTable(config *formatconfig.DiskImage) *PartitionTabelBuilder {
	partitionTable := NewPartitionTable(TableTypesnameToCode[config.Label])

	for i, cfgPart := range config.Partitions {
		partition := buildPartition(
			config.OutFile+strconv.FormatInt(int64(i)+1, 10),
			cfgPart,
		)
		partitionTable.WithPartitionDefinition(partition)
	}
	return partitionTable
}

func buildPartition(label string, cfgPart *formatconfig.Partition) *PartitionDefinitionBuilder {
	partitionBuilder := NewPartitionBuilder(label).
		WithName(cfgPart.Label).
		OfType(PartitionNameToCode[cfgPart.Type])

	return partitionBuilder.
		StartAtIf(cfgPart.StartAt).
		WithSizeIf(cfgPart.Size).
		SetBootableIf(cfgPart.Bootable)
}
