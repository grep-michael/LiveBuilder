package isoinstaller

import (
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	loopmanager "LiveBuilder/USBImager/LoopManager"
	"LiveBuilder/Utils"
	"fmt"
	"log"
	"os"
)

func InstallISOFromConfig(config *formatconfig.DiskImage) {
	//run("mount", "-o", "loop,ro", iso, isoMnt)
	iso_mount, err := os.MkdirTemp("", "iso-*")
	if err != nil {
		log.Println(err)
	}
	stdout, stderr, err := Utils.Run_command("sudo", "mount", "-o", "loop,ro", config.InFile, iso_mount)
	if err != nil {
		log.Println(err)
		log.Println(stdout)
		log.Println(stderr)
	}

	outfile_manager := loopmanager.NewLoopManager(config.OutFile)
	outfile_manager.OpenLoop()
	outfile_manager.MountPartitions()

	defer func() {
		Utils.Run_command("sudo", "umount", iso_mount)
		outfile_manager.UnmountPartitions()
		outfile_manager.CloseLoop()
	}()

	for _, cfgPartition := range config.Partitions {
		if cfgPartition.IsoPartition {
			part_mount := outfile_manager.GetPartitionsMountPointByLabel(cfgPartition.Label)
			fmt.Println("rsync", "-aH", "--info=progress2", iso_mount+"/", part_mount)
			stdout, stderr, err := Utils.Run_command("rsync", "-aH", "--info=progress2", iso_mount+"/", part_mount)
			if err != nil {
				log.Println(err)
				log.Println(stdout)
				log.Println(stderr)
			}
		}
	}

}
