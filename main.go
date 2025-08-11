package main

import (
	filesystem "LiveBuilder/Filesystem"
	filelistwidgets "LiveBuilder/frontend/FileListWidgets"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	plc := filelistwidgets.NewFileListContainer(filesystem.PACKAGE_DIR_ID)

	myApp := app.New()
	myWindow := myApp.NewWindow("test window")
	//
	cnt := plc.GetContainer()
	fmt.Println(cnt.Size())
	myWindow.SetContent(cnt)
	//
	windowPadding := fyne.NewSize(30, 60)
	windowSize := fyne.NewSize(
		cnt.Size().Width+windowPadding.Width,
		cnt.Size().Height+windowPadding.Height,
	)

	myWindow.Resize(windowSize)

	myWindow.ShowAndRun()
}
