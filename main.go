package main

import (
	filesystem "LiveBuilder/Filesystem"
	preflightchecks "LiveBuilder/PreFlightChecks"
	usbimager "LiveBuilder/USBImager"
	frontend "LiveBuilder/frontend"
	"fmt"
	"log"
	"os"
)

var LOGFILE *os.File

func configureLogging() {
	log_file, err := filesystem.GetAppDataDir()
	log_file += "/app.log"
	LOGFILE, err = os.OpenFile(log_file, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InitLoging("LiveBuilder: ", "	", LOGFILE)

}

func main() {
	defer LOGFILE.Close()
	configureLogging()
	log.Println("App Start")
	testMain()
}

func guiMain() {
	preflightchecks.CheckAll(false)

	mainWindow := frontend.NewMainWindow("Live Builder")
	mainWindow.ShowAndRun()
}

func testMain() {
	imager := usbimager.NewUSBImager()
	err := imager.ImageUSB("/tmp/Fake.iso", "/dev/sdd")
	fmt.Println(err)
}
