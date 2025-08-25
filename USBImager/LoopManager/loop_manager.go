package loopmanager

import (
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LoopManager struct {
	DevicePath string
	LoopPath   string
	loopPaths  []string
	Partitions []PartitionInfo
}

func NewLoopManager(device string) *LoopManager {
	return &LoopManager{
		DevicePath: device,
	}
}

func (self *LoopManager) OpenLoop() error {
	_, _, err := Utils.Run_command("sudo", "losetup", "-Pf", self.DevicePath)
	if err != nil {

		log.Printf("Looping setup error: %+v\n", err)
		return err
	}
	stdout, _, err := Utils.Run_command("sudo", "losetup", "-j", self.DevicePath)
	if err != nil {
		log.Printf("Looping detection error: %+v\n", err)
		return err
	}

	var loopDev string
	fmt.Sscanf(string(stdout), "%s:", &loopDev)
	loopDev = loopDev[:len(loopDev)-1]
	self.LoopPath = loopDev

	return nil
}

func (self *LoopManager) CloseLoop() {

	_, _, err := Utils.Run_command("sudo", "losetup", "-d", self.LoopPath)
	if err != nil {
		log.Printf("ERROR UNLOOPING: %+v\n", err)
		panic(err)
	}
}
func (self *LoopManager) loadPartitions() {
	time.Sleep(1 * time.Second)
	partitions, err := ParseLsblkFilesystems(self.LoopPath)
	if err != nil {
		log.Println(err)
	}
	self.Partitions = partitions
}

func (self *LoopManager) MountPartitions() error {
	self.loadPartitions()
	for _, partition := range self.Partitions {
		tmp_dir := filepath.Join("/tmp", "loopmanager", partition.Label)
		if err := os.MkdirAll(tmp_dir, 0777); err != nil {
			log.Printf("error making mount directory error: %v\n", err)
			return err
		}
		fmt.Println("sudo", "mount", tmp_dir, partition.DevicePath)
		stdout, stderr, err := Utils.Run_command("sudo", "mount", partition.DevicePath, tmp_dir)
		if err != nil {
			log.Printf("Mounting error: %v - %s\n", err, (string(stdout) + ":" + string(stderr)))
			return fmt.Errorf("Mounting error: %v - %s", err, (string(stdout) + ":" + string(stderr)))
		}
	}
	self.loadPartitions()
	return nil
}

func (self *LoopManager) UnmountPartitions() error {
	self.loadPartitions()
	for _, partition := range self.Partitions {
		for _, mount_point := range partition.MountPoints {
			_, stderr, err := Utils.Run_command("sudo", "umount", mount_point)
			if err != nil {
				log.Printf("unmount error: %v - %s", err, stderr)
			}
		}
	}
	self.loadPartitions()
	return nil
}

func (self *LoopManager) GetPartitions() []PartitionInfo {
	return self.Partitions
}
