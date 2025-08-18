package buildmanager

/*
Imports (Drops/copies) selected files/customizations into the config directory of live-build
*/

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Importer struct {
	buildPath     string
	updateChannel chan LogUpdate
}

func NewImporter(updateChan chan LogUpdate) *Importer {
	return &Importer{
		updateChannel: updateChan,
	}
}

func (self *Importer) SetBuildPath(buildPath string) {
	self.buildPath = buildPath
}

func (self *Importer) ImportAll() error {
	if self.buildPath == "" {
		return fmt.Errorf("Build Path Not set")
	}
	self.updateChannel <- LogUpdate{
		Append:     false,
		Message:    "Starting file import\n\n",
		UpdateType: START,
	}
	// fancy shit
	operations := []struct {
		fn   func() error
		name string
	}{
		{self.DropCustomFiles, "DropCustomFiles"},
		{self.DropPackages, "DropPackages"},
		{self.DropSplashImages, "DropSplashImages"},
	}

	for _, op := range operations {
		if err := op.fn(); err != nil {
			return fmt.Errorf("%s error: %v", op.name, err)
		}
	}
	self.updateChannel <- LogUpdate{
		Append:     true,
		Message:    "Importing finished\n",
		UpdateType: END,
	}
	return nil
}

func (self *Importer) DropPackages() error {
	self.updateChannel <- LogUpdate{
		Append:     true,
		Message:    "Dropping Packages\n",
		UpdateType: UPDATE,
	}
	packageMap := appstate.GetGlobalState().GetDirectoryEntryMap(filesystem.PACKAGE_DIR_ID)

	for _, value := range packageMap {
		self.dropFileFromDirectoryEntry(value)
	}
	return nil
}

func (self *Importer) DropSplashImages() error {
	self.updateChannel <- LogUpdate{
		Append:     true,
		Message:    "Dropping Splash images\n",
		UpdateType: UPDATE,
	}
	log.Println("Dropping Splash images")
	splashMap := filesystem.GetFileManager().GetFileSystem(filesystem.SPLASH_SCREENS_ID)
	for _, value := range splashMap {
		self.dropFileFromDirectoryEntry(value)
	}
	return nil
}

func (self *Importer) DropCustomFiles() error {
	self.updateChannel <- LogUpdate{
		Append:     true,
		Message:    "Dropping Custom Files\n",
		UpdateType: UPDATE,
	}
	customFileMap := appstate.GetGlobalState().GetDirectoryEntryMap(filesystem.CUSTOMFILES_DIR_ID)

	for _, value := range customFileMap {
		self.dropFileFromDirectoryEntry(value)
	}
	return nil
}

func (self *Importer) dropFileFromDirectoryEntry(file filesystem.DirectoryEntry) error {
	var inFile, outFile *os.File
	var err error

	outFilePath := filepath.Join(self.buildPath, file.MetaData.InstallPath)
	inFIlePath := file.FullPath()

	if inFile, err = os.Open(inFIlePath); err != nil {
		return err
	}
	defer inFile.Close()
	if err := os.MkdirAll(filepath.Dir(outFilePath), 0777); err != nil {
		return err
	}
	if outFile, err = os.OpenFile(outFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, inFile); err != nil {
		return err
	}
	outFile.WriteString("\n")

	msg := fmt.Sprintf("Added %s file to %s\n", file.Name(), file.MetaData.InstallPath)
	self.updateChannel <- LogUpdate{
		Append:     true,
		Message:    msg,
		UpdateType: UPDATE,
	}

	return nil
}
