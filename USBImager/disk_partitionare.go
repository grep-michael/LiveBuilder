package usbimager

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type DiskPartitionare struct {
	Device         FileObject
	PartitionTable *PartitionTabelBuilder
	loopPath       string
}

func NewDiskPartionare(device FileObject) *DiskPartitionare {
	return &DiskPartitionare{
		Device: device,
	}
}
func (self *DiskPartitionare) SetPartitionTable(table *PartitionTabelBuilder) {
	self.PartitionTable = table
}
func (self *DiskPartitionare) PartitionDisk() error {
	if self.PartitionTable == nil {
		return fmt.Errorf("No FileObject supplied")
	}

	log.Println("PARTITION TABLE")
	log.Println(self.PartitionTable.ToSfdisk())

	partitionReader := strings.NewReader(self.PartitionTable.ToSfdisk())

	var out bytes.Buffer
	cmd := exec.Command("sudo", "sfdisk", "-f", self.Device.Path)
	cmd.Stdin = partitionReader
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		log.Printf("sfdisk error: %s\n", out.String())
		return err
	}
	log.Printf("sfdisk out: %s\n", out.String())

	return nil
}
func (self *DiskPartitionare) WriteFileSystems() error {
	if err := self.openLoop(); err != nil {
		fmt.Println(err)
		return err
	}
	for i, partition := range self.PartitionTable.partitions {
		device := fmt.Sprintf("%sp%d", self.loopPath, i+1)
		cmd, _ := getFormatCommandForDeivce(partition.partType, device, partition.volumeName)
		stdout, stderr, err := run("sudo", strings.Split(cmd, " ")...)
		log.Printf("MKFS COMMAND: %s\n", cmd)
		log.Printf("stdout (%s)\nstderr (%s)\n", stdout, stderr)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	//self.closeLoop()

	return nil
}
func (self *DiskPartitionare) openLoop() error {
	_, _, err := run("sudo", "losetup", "-Pf", self.Device.Path)
	if err != nil {
		log.Printf("Looping setup error: %+v\n", err)
		return err
	}

	// Find loop device
	stdout, _, err := run("sudo", "losetup", "-j", self.Device.Path)

	if err != nil {
		log.Printf("Looping detection error: %+v\n", err)
		return err
	}

	var loopDev string
	fmt.Sscanf(string(stdout), "%s:", &loopDev)

	loopDev = loopDev[:len(loopDev)-1]

	fmt.Println(loopDev)
	self.loopPath = loopDev
	return nil
}
func (self *DiskPartitionare) closeLoop() {
	_, _, err := run("sudo", "losetup", "-d", self.loopPath)
	if err != nil {
		log.Printf("ERROR UNLOOPING: %+v\n", err)
		panic(err)
	}
}
