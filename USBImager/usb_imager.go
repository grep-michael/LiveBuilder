package usbimager

import (
	fileobject "LiveBuilder/USBImager/FileObject"
	grubinstaller "LiveBuilder/USBImager/GrubInstaller"
	partitioning "LiveBuilder/USBImager/Partitioning"
	"fmt"
	"log"
)

type USBImager struct {
}

func NewUSBImager() *USBImager {
	return &USBImager{}
}

func (self *USBImager) ImageUSB(iso_file, out_file string) error {
	inFile, outFile, err := self.initalizeFileInfos(iso_file, out_file)
	if err != nil {
		return err
	}
	err = self.prepOutFile(outFile, CalculateNeededSizeForISO(inFile))
	if err != nil {
		log.Println(err)
	}

	grubin := grubinstaller.NewGrubInstaller("/dev/sdd")
	defer grubin.Finish()
	err = grubin.InstallToBootPartitionByLabel("BOOT")
	fmt.Println(err)

	return nil

}

func (self *USBImager) initalizeFileInfos(iso_file, out_file string) (fileobject.FileObject, fileobject.FileObject, error) {
	isoFileInfo, err := fileobject.NewFileObject(iso_file, false)
	if err != nil {
		return errorFileObject(err)
	}
	outFileInfo, err := fileobject.NewFileObject(out_file, true)
	if err != nil {
		return errorFileObject(err)
	}
	if isoFileInfo.Type != fileobject.TypeRegularFile {
		return errorFileObject(fmt.Errorf("iso file type not support: %s\n", isoFileInfo.Type.String()))
	}
	if (1<<outFileInfo.Type)&fileobject.ALLOWEDTYPES == 0 {
		return errorFileObject(fmt.Errorf("destinated file type not support: %s\n", outFileInfo.Type.String()))
	}
	return isoFileInfo, outFileInfo, nil
}

func (self *USBImager) prepOutFile(outFile fileobject.FileObject, outFileSize int64) error {
	if outFileSize == 0 {
		outFileSize = 2 * 1024 * 1024 * 1024 //4gb default
	}

	if outFile.Info.Size < outFileSize {
		log.Println("Out file to small, resizing")
		err := outFile.Resize(outFileSize)
		if err != nil {
			return err
		}
	}
	//if its a storage device unmount it and wipe it
	if outFile.Type == fileobject.TypeBlockDevice {
		err := outFile.UmountPartitions()
		if err != nil {
			return err
		}
		err = outFile.WipeFS()
		if err != nil {
			return err
		}
	}

	diskpart := partitioning.StandardLinuxMBRBootPart(outFile.Path)
	diskpart.PartitionDisk()
	diskpart.WriteFileSystems()

	return nil
}

func errorFileObject(err error) (fileobject.FileObject, fileobject.FileObject, error) {
	return fileobject.FileObject{}, fileobject.FileObject{}, err
}
