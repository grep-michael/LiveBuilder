package usbimager

import (
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
)

type FileObject struct {
	Type FileType
	Path string
	Info *SystemFileInfo
}

func (self *FileObject) Resize(size int64) error {
	if self.Type != TypeRegularFile {
		//might not be an error
		log.Println("Cant make outfile larger, not of type Regular File")
		return nil
	}
	//if self.Info.Size > 1 {
	//	log.Println("attemping to resize file that has a predefined size, skipping resize")
	//	return nil
	//}
	file, err := os.OpenFile(self.Path, os.O_RDWR, 0644)
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
func (self *FileObject) UmountPartitions() error {

	_, stderr, err := Utils.Run_command("sudo", "umount", self.Path+"*")

	if err != nil {
		//erroring doesnt necessarily mean the umount failed, could also mean it wasnt mounted in the first place, fuck it we ball
		log.Printf("umount command failed: %+v\n", err)
		log.Printf("%s\n\n", stderr)
	}
	return nil
}
func (self *FileObject) WipeFS() error {

	stdout, stderr, err := Utils.Run_command("sudo", "wipefs", "-af", "-b", self.Path)

	if err != nil {
		log.Printf("wipefs command failed: %+v\n", err)
	}
	log.Printf("WipeFs out\nstdout (%s)\nstderr (%s)\n", stdout, stderr)
	return err

}
