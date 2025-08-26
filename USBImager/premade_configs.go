package usbimager

import (
	formatconfig "LiveBuilder/USBImager/FormatConfig"
	"LiveBuilder/Utils"
	"fmt"
)

func NewFat32BootAndExt4System(inFile, outFile string) *formatconfig.DiskImage {
	cfgJson := GenerateConfigJsonFromTemplate(inFile, outFile, `{
		"label":"dos",
		"label_id":"0x12345678",
		"units":"sectors",
		"inFile":"{{.InFile}}",
		"outFile":"{{.OutFile}}",
		"partitions":[
			{
				"label":"BOOT",
				"StartAt":2048,
				"Size":100000,
				"Type":"W95 FAT32 (LBA)",
				"Bootable":true,
				"LoopPath":"",
				"MountPath":"",
				"Grub":{
					"BIOS":true,
                  	"UEFI":true,
					"Redirect":"SYSTEM"
				}
			},
			{
				"label":"SYSTEM",
				"Type":"Linux",
				"LoopPath":"",
				"MountPath":""
			}
		]
	}`)

	cfg, err := formatconfig.NewDiskImageFromJSON(cfgJson)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cfg

}

func GenerateConfigJsonFromTemplate(inFile, outFile, jsonTemplate string) string {
	data := struct {
		InFile  string
		OutFile string
	}{
		InFile:  inFile,
		OutFile: outFile,
	}
	out := Utils.BuildTemplate(data, jsonTemplate)
	return out
}
