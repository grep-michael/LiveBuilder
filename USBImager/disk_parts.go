package usbimager

func StandardLinuxMBRBootPart(fileobject FileObject) *DiskPartitionare {
	partition1 := NewPartitionBuilder(fileobject.Path + "1").
		WithName("BOOT").
		StartAt("2048").
		WithSize("512M").
		OfType(W95_FAT32_LBA).
		SetBootable(true)

	partition2 := NewPartitionBuilder(fileobject.Path + "2").
		WithName("SYSTEM").
		OfType(Linux)

	partitionTable := NewPartitionTable(TABLETYPE_MBR).
		WithPartitionDefinition(partition1).
		WithPartitionDefinition(partition2)

	diskpart := NewDiskPartionare(fileobject)
	diskpart.SetPartitionTable(partitionTable)

	return diskpart
}
