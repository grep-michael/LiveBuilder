package main

import (
	filesystem "LiveBuilder/Filesystem"
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
	log.SetOutput(LOGFILE)
	log.SetPrefix("LiveBuilder: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	defer LOGFILE.Close()
	configureLogging()
	log.Println("App Start")

	mainWindow := frontend.NewMainWindow("Live Builder")
	mainWindow.ShowAndRun()

}
