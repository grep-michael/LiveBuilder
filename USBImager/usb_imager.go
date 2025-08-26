package usbimager

import (
	filesetup "LiveBuilder/USBImager/FileSetup"
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	grubinstaller "LiveBuilder/USBImager/GrubInstaller"
	isoinstaller "LiveBuilder/USBImager/IsoInstaller"
	partitioning "LiveBuilder/USBImager/Partitioning"
	"fmt"
	"log"
)

type USBImager struct {
	config *formatconfig.DiskImage
}

func NewUSBImager() *USBImager {
	return &USBImager{}
}

func (self *USBImager) ImageUSB(config *formatconfig.DiskImage) error {
	self.config = config

	log.Printf("USB imager using config\n%+v\n", self.config)

	if err := filesetup.SetUpFiles(self.config); err != nil {
		fmt.Println(err)
		return nil
	}

	self.partitionConfig()
	self.installGrub()
	self.copyInfileToOutFile()
	return nil
}

func (self *USBImager) copyInfileToOutFile() {
	isoinstaller.InstallISOFromConfig(self.config)
}

func (self *USBImager) partitionConfig() {
	diskpart := partitioning.CreatePartitionareForConfig(self.config)
	diskpart.PartitionDisk()
	diskpart.WriteFileSystems()
}
func (self *USBImager) installGrub() {
	grubin := grubinstaller.NewGrubInstaller(self.config.OutFile)
	grubin.InstallFromConfig(self.config)
}
