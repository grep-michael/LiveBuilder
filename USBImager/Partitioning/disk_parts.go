package usbimager

func StandardLinuxMBRBootPart(filepath string) *DiskPartitionare {
	partition1 := NewPartitionBuilder(filepath + "1").
		WithName("BOOT").
		StartAt("2048").
		WithSize("100000").
		OfType(W95_FAT32_LBA).
		SetBootable(true)

	partition2 := NewPartitionBuilder(filepath + "2").
		WithName("SYSTEM").
		OfType(Linux)

	partitionTable := NewPartitionTable(TABLETYPE_MBR).
		WithPartitionDefinition(partition1).
		WithPartitionDefinition(partition2)

	diskpart := NewDiskPartionare(filepath)
	diskpart.SetPartitionTable(partitionTable)

	return diskpart
}
