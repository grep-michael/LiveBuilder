package usbimager

import (
	"fmt"
	"testing"
)

func TestImageUSB(t *testing.T) {
	imager := NewUSBImager()
	err := imager.ImageUSB("/tmp/Fake.iso", "/tmp/usb.img")
	fmt.Println(err)
}
func TestGidName(t *testing.T) {
	disk := getGroupNameForGID("6")
	fmt.Println(disk)
}

func TestUtil(t *testing.T) {
	stat_block, err := NewFileObject("/tmp/file", true)
	fmt.Println(stat_block)
	fmt.Println(err)
}
