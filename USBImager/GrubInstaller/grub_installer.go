package grubinstaller

import (
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	loopmanager "LiveBuilder/USBImager/LoopManager"
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
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

	return grubin
}

func (self *GrubInstaller) Mount() {
	self.LoopManager.OpenLoop()
	self.LoopManager.MountPartitions()
}
func (self *GrubInstaller) Unmount() {
	self.LoopManager.UnmountPartitions()
	self.LoopManager.CloseLoop()
}

func (self *GrubInstaller) installUEFIToPartition(parition *loopmanager.PartitionInfo) {
	if len(parition.MountPoints) <= 0 {
		log.Println("attemped to install to unmounted parition")
		return
	}
	mount_pount := parition.GetMountPointsString()
	stdout, stderr, _ := Utils.Run_command("sudo", "grub-install", "--target=i386-pc", "--boot-directory="+mount_pount+"/boot", self.DevicePath)
	fmt.Println(stdout, stderr)

	self.swapToSignedBinaries(mount_pount)
}

func (self *GrubInstaller) installBioToPartition(parition *loopmanager.PartitionInfo) {
	if len(parition.MountPoints) <= 0 {
		log.Println("attemped to install to unmounted parition")
		return
	}
	mount_pount := parition.GetMountPointsString()

	stdout, stderr, _ := Utils.Run_command("sudo", "grub-install", "--target=x86_64-efi", "--efi-directory="+mount_pount,
		"--boot-directory="+mount_pount+"/boot", "--removable")
	fmt.Println(stdout, stderr)
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

func (self *GrubInstaller) dropRedirectStub(parition *loopmanager.PartitionInfo, grubcfg *formatconfig.Grub) {

	if len(parition.MountPoints) <= 0 {
		log.Println("attemping to stub to unmounted filesystem")
		return
	}
	mount_point := parition.GetMountPointsString()

	data := struct {
		MainSystem string
	}{
		MainSystem: grubcfg.Redirect,
	}

	stub := Utils.BuildTemplate(data, `search --no-floppy --set=root --label {{.MainSystem}}
configfile /boot/grub/grub.cfg
`)
	log.Printf("Dropping grub redirect stub to \n--%s\n%s\n", mount_point, stub)
	Utils.WriteFile(filepath.Join(mount_point, "boot", "grub", "grub.cfg"), stub, 0644)
	time.Sleep(1 * time.Second)
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

func (self *GrubInstaller) InstallFromConfig(config *formatconfig.DiskImage) error {
	self.Mount()
	defer self.Unmount()
	for _, cfgPart := range config.Partitions {

		if cfgPart.Grub != nil {
			part := self.getPartitionByLabel(cfgPart.Label)
			if cfgPart.Grub.BIOS {
				self.installBioToPartition(part)
			}
			if cfgPart.Grub.UEFI {
				self.installUEFIToPartition(part)
			}
			if cfgPart.Grub.Redirect != "" {
				fmt.Println("dropping stub")
				fmt.Println(cfgPart.Grub.Redirect)
				self.dropRedirectStub(part, cfgPart.Grub)
			}
		}
	}
	return nil
}
