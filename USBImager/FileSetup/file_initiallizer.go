package filesetup

import (
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
)

func SetUpFiles(config *formatconfig.DiskImage) error {
	infile := newFileObject(config.InFile, false)
	outfile := newFileObject(config.OutFile, true)
	if outfile == nil || infile == nil {
		return fmt.Errorf("One or both files doesnt exist")
	}
	if infile.isNotAllowedFile() {
		log.Printf("Infile not supported type: %s\n", infile.fileType.String())
		return fmt.Errorf("Infile not supported type: %s\n", infile.fileType.String())
	}
	if outfile.isNotAllowedFile() {
		log.Printf("Outfile not supported type: %s\n", outfile.fileType.String())
		return fmt.Errorf("Outfile not supported type: %s\n", outfile.fileType.String())
	}

	unmount(outfile)

	if err := wipeFS(outfile); err != nil {
		return err
	}
	if infile.size <= 1 {
		return fmt.Errorf("Infile size is to small, this will break resizing the outFile if the outfile is a file and not a device")
	}
	megabyteBuffer := int64(1000000) //50MB buffer for boot partitions, can be lowered tbh
	if err := resize(infile.size+megabyteBuffer, outfile); err != nil {
		return err
	}

	return nil
}

func wipeFS(file *fileObject) error {
	stdout, stderr, err := Utils.Run_command("sudo", "wipefs", "-af", file.path)
	if err != nil {
		log.Printf("wipefs command failed: %+v\n", err)
	}
	log.Printf("WipeFs out\nstdout (%s)\nstderr (%s)\n", stdout, stderr)
	return err
}

func resize(size int64, fileobject *fileObject) error {
	if fileobject.fileType != TypeRegularFile {
		//might not be an error
		log.Println("Cant make outfile larger, not of type Regular File")
		return nil
	}
	file, err := os.OpenFile(fileobject.path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("Error opening file: %v\n", err)
	}
	err = file.Truncate(size)
	file.Close()
	if err != nil {
		return fmt.Errorf("Error truncating file: %v\n", err)
	}
	return err
}
func unmount(fileobject *fileObject) error {

	_, stderr, err := Utils.Run_command("sudo", "umount", fileobject.path+"*")

	if err != nil {
		//erroring doesnt necessarily mean the umount failed, could also mean it wasnt mounted in the first place, fuck it we ball
		log.Printf("umount command failed: %+v\n", err)
		log.Printf("%s\n\n", stderr)
	}
	return nil
}
