package usbimager

import (
	loopmanager "LiveBuilder/USBImager/LoopManager"
	"LiveBuilder/Utils"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type DiskPartitionare struct {
	Device         string
	PartitionTable *PartitionTabelBuilder
	loopManager    *loopmanager.LoopManager
}

func NewDiskPartionare(device string) *DiskPartitionare {
	return &DiskPartitionare{
		Device:      device,
		loopManager: loopmanager.NewLoopManager(device),
	}
}
func (self *DiskPartitionare) SetPartitionTable(table *PartitionTabelBuilder) {
	self.PartitionTable = table
}
func (self *DiskPartitionare) PartitionDisk() error {
	if self.PartitionTable == nil {
		return fmt.Errorf("No FileObject supplied")
	}

	log.Printf("PARTITION TABLE\n%s\n", self.PartitionTable.ToSfdisk())

	partitionReader := strings.NewReader(self.PartitionTable.ToSfdisk())

	var out bytes.Buffer
	cmd := exec.Command("sudo", "sfdisk", "-f", self.Device)
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
	defer self.loopManager.CloseLoop()
	if err := self.loopManager.OpenLoop(); err != nil {
		fmt.Println(err)
		return err
	}
	for i, partition := range self.PartitionTable.partitions {
		device := fmt.Sprintf("%sp%d", self.loopManager.LoopPath, i+1)
		cmd, _ := getFormatCommandForDeivce(partition.partType, device, partition.volumeName)
		stdout, stderr, err := Utils.Run_command("sudo", strings.Split(cmd, " ")...)
		log.Printf("MKFS COMMAND: %s\n", cmd)
		log.Printf("stdout (%s)\nstderr (%s)\n", stdout, stderr)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
