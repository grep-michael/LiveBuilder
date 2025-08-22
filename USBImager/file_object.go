package usbimager

import (
	"LiveBuilder/USBImager/paritions"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type FileObject struct {
	Type FileType
	Path string
	Info *SystemFileInfo
}

func (self *FileObject) resize(size int64) error {
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

func (self *FileObject) umountPartitions() error {

	_, stderr, err := run("sudo", "umount", self.Path+"*")

	if err != nil {
		//erroring doesnt necessarily mean the umount failed, could also mean it wasnt mounted in the first place, fuck it we ball
		log.Printf("umount command failed: %+v\n", err)
		log.Printf("%s\n\n", stderr)
	}
	return nil
}

func (self *FileObject) wipeFS() error {

	stdout, stderr, err := run("sudo", "wipefs", "-a", "-b", self.Path)

	if err != nil {
		log.Printf("wipefs command failed: %+v\n", err)
	}
	log.Printf("WipeFs out\nstdout (%s)\nstderr (%s)\n", stdout, stderr)
	return err

}

func (self *FileObject) partition() error {

	partition1 := paritions.NewPartitionBuilder(self.Path + "1").
		StartAt("2048").
		WithSize("512M").
		OfType(paritions.W95_FAT32_LBA).
		SetBootable(true)

	partition2 := paritions.NewPartitionBuilder(self.Path + "2").
		OfType(paritions.Linux)

	partitionTable := paritions.NewPartitionTable(paritions.TABLETYPE_MBR).
		WithPartitionDefinition(partition1).
		WithPartitionDefinition(partition2).
		ToSfdisk()
	log.Println("PARTITION TABLE")
	log.Println(partitionTable)

	partitionMapFile, err := os.Create("/tmp/partition.sfdisk")
	if err != nil {
		return fmt.Errorf("Failed to create partition file %v\n", err)
	}
	partitionMapFile.WriteString(partitionTable)
	if _, err := partitionMapFile.Seek(0, 0); err != nil {
		return err
	}

	var out bytes.Buffer
	cmd := exec.Command("sudo", "sfdisk", self.Path)
	cmd.Stdin = partitionMapFile
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Printf("sfdisk error: %s\n", out.String())
		return err
	}
	log.Printf("sfdisk out: %s\n", out.String())
	partitionMapFile.Close()
	return nil
}
