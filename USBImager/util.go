package usbimager

import (
	fileobject "LiveBuilder/USBImager/FileObject"
	"log"
	"os"
	"path/filepath"
)

func CalculateNeededSizeForISO(iso fileobject.FileObject) int64 {
	bytes := iso.Info.Size
	//const GB = 1024 * 1024 * 1024
	if bytes <= 0 {
		return 0
	}
	//gbNeeded := (bytes + GB - 1) / GB //ceiling division
	gbNeeded := bytes + 52428800 //~50MiB for boot partition
	return gbNeeded              //* GB
}

func WriteFile(path, content string, perm os.FileMode) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		log.Fatalf("write %s: %v", path, err)
	}
}
