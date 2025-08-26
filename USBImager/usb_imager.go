package usbimager

import (
	filesetup "LiveBuilder/USBImager/FileSetup"
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	grubinstaller "LiveBuilder/USBImager/GrubInstaller"
	partitioning "LiveBuilder/USBImager/Partitioning"
	"fmt"
)

type USBImager struct {
	config *formatconfig.DiskImage
}

func NewUSBImager() *USBImager {
	return &USBImager{}
}

func (self *USBImager) ImageUSB(config *formatconfig.DiskImage) error {
	self.config = config

	fmt.Printf("%+v\n", self.config)

	if err := filesetup.SetUpFiles(self.config); err != nil {
		fmt.Println(err)
		return nil
	}
	diskpart := partitioning.CreatePartitionareForConfig(config)
	diskpart.PartitionDisk()
	diskpart.WriteFileSystems()

	grubin := grubinstaller.NewGrubInstaller(self.config.OutFile)
	grubin.Init()
	defer grubin.Finish()
	grubin.InstallFromConfig(self.config)
	return nil
}
