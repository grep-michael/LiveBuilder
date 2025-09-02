package main

import (
	filesystem "LiveBuilder/Filesystem"
	preflightchecks "LiveBuilder/PreFlightChecks"
	frontend "LiveBuilder/frontend"
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
	guiMain()
}

func guiMain() {
	//# Regenerate X11 authorization
	//xhost +local:
	preflightchecks.CheckAll(false)

	mainWindow := frontend.NewMainWindow("Live Builder")
	mainWindow.ShowAndRun()
}
