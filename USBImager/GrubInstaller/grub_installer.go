package grubinstaller

import (
	loopmanager "LiveBuilder/USBImager/LoopManager"
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type GrubInstaller struct {
	LoopManager *loopmanager.LoopManager
	DevicePath  string
}

func NewGrubInstaller(devicePath string) *GrubInstaller {
	grubin := &GrubInstaller{
		DevicePath:  devicePath,
		LoopManager: loopmanager.NewLoopManager(devicePath),
	}

	grubin.Init()

	return grubin
}

func (self *GrubInstaller) Init() {
	self.LoopManager.OpenLoop()
	self.LoopManager.MountPartitions()
}
func (self *GrubInstaller) Finish() {
	self.LoopManager.UnmountPartitions()
	self.LoopManager.CloseLoop()
}

func (self *GrubInstaller) InstallToBootPartition(parition loopmanager.PartitionInfo) {

	if len(parition.MountPoints) <= 0 {
		log.Println("attemped to install to unmounted parition")
		return
	}

	mount_pount := parition.MountPoints[0]

	stdout, stderr, _ := Utils.Run_command("sudo", "grub-install", "--target=i386-pc", "--boot-directory="+mount_pount+"/boot", self.DevicePath)
	fmt.Println(stdout, stderr)

	stdout, stderr, _ = Utils.Run_command("sudo", "grub-install", "--target=x86_64-efi", "--efi-directory="+mount_pount,
		"--boot-directory="+mount_pount+"/boot", "--removable")
	fmt.Println(stdout, stderr)

	self.swapToSignedBinaries(mount_pount)
	self.dropRedirectStub(mount_pount)

}

func (self *GrubInstaller) swapToSignedBinaries(mount_point string) error {
	//signed boot.efi
	if _, err := os.Stat("/usr/lib/shim/shimx64.efi.signed"); err == nil {
		_, _, err = Utils.Run_command("cp", "/usr/lib/shim/shimx64.efi.signed", filepath.Join(mount_point, "EFI", "BOOT", "BOOTX64.EFI"))
		if err != nil {
			return err
		}
	} else {
		log.Printf("shim-signed not found; install shim-signed for Secure Boot.")
		return nil
	}
	//signed grub.efi
	if _, err := os.Stat("/usr/lib/grub/x86_64-efi-signed/grubx64.efi.signed"); err == nil {
		_, _, err := Utils.Run_command("cp", "/usr/lib/grub/x86_64-efi-signed/grubx64.efi.signed", filepath.Join(mount_point, "EFI", "BOOT", "grubx64.efi"))
		if err != nil {
			return err
		}
	} else {
		log.Printf("Signed GRUB not found; install grub-efi-amd64-signed for Secure Boot.")
		return nil
	}
	return nil
}

func (self *GrubInstaller) InstallToBootPartitionByLabel(label string) error {
	partition := self.getPartitionByLabel(label)
	if partition == nil {
		log.Printf("Failed to find any partition with name %s\n", label)
		return fmt.Errorf("Failed to find any partition with name %s\n", label)
	}
	self.InstallToBootPartition(*partition)
	return nil
}

func (self *GrubInstaller) InstallConfigToPartitionByLabel(label string, config string) error {
	partition := self.getPartitionByLabel(label)
	if partition == nil {
		log.Printf("Failed to find any partition with name %s\n", label)
		return fmt.Errorf("Failed to find any partition with name %s\n", label)
	}
	self.InstallConfigToPartition(*partition, config)
	return nil
}

func (self *GrubInstaller) InstallConfigToPartition(partition loopmanager.PartitionInfo, config string) {

}

func (self *GrubInstaller) dropRedirectStub(mount_point string) {
	stub := `search --no-floppy --set=root --label SYSTEM
configfile /boot/grub/grub.cfg
`
	writeFile(filepath.Join(mount_point, "boot", "grub", "grub.cfg"), stub, 0644)
}

func (self *GrubInstaller) getPartitionByLabel(label string) *loopmanager.PartitionInfo {
	partitions := self.LoopManager.GetPartitions()
	for _, partition := range partitions {
		if partition.Label == label {
			return &partition
		}
	}
	return nil
}

func writeFile(path, content string, perm os.FileMode) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}
