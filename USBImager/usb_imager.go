package usbimager

import (
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
	err = self.prepOutFile(outFile, calculateNeededSizeForISO(inFile))
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (self *USBImager) initalizeFileInfos(iso_file, out_file string) (FileObject, FileObject, error) {
	isoFileInfo, err := NewFileObject(iso_file, false)
	if err != nil {
		return FileObject{}, FileObject{}, err
	}
	outFileInfo, err := NewFileObject(out_file, true)
	if err != nil {
		return FileObject{}, FileObject{}, err
	}
	if isoFileInfo.Type != TypeRegularFile {
		return FileObject{}, FileObject{}, fmt.Errorf("iso file type not support: %s\n", isoFileInfo.Type.String())
	}
	if (1<<outFileInfo.Type)&ALLOWEDTYPES == 0 {
		return FileObject{}, FileObject{}, fmt.Errorf("destinated file type not support: %s\n", outFileInfo.Type.String())
	}
	return isoFileInfo, outFileInfo, nil
}

func (self *USBImager) prepOutFile(outFile FileObject, outFileSize int64) error {
	if outFileSize == 0 {
		outFileSize = 4 * 1024 * 1024 * 1024 //4gb default
	}

	if outFile.Info.Size < outFileSize {
		log.Println("Out file to small, resizing")
		err := outFile.resize(outFileSize)
		if err != nil {
			return err
		}
	}
	//if its a storage device unmount it and wipe it
	if outFile.Type == TypeBlockDevice {
		err := outFile.umountPartitions()
		if err != nil {
			return err
		}
		err = outFile.wipeFS()
		if err != nil {
			return err
		}
	}
	//partition it
	err := outFile.partition()
	if err != nil {
		return err
	}

	return nil
}
