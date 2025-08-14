package buildmanager

/*
Imports (Drops/copies) selected files/customizations into the config directory of live-build
*/

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	"bufio"
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

	outfile_path := filepath.Join(self.buildPath, "config/package-lists/live.list.chroot")
	outFile, err := os.OpenFile(outfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer outFile.Close()
	if err != nil {
		return err
	}

	for _, value := range packageMap {
		inFile, err := os.Open(value.FullPath())
		defer inFile.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
		_, err = outFile.WriteString("\n")
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Added %s package to config/package-lists/live.list.chroot\n", value.Name())
		self.updateChannel <- LogUpdate{
			Append:     true,
			Message:    msg,
			UpdateType: UPDATE,
		}
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
		inFile, err := os.Open(value.FullPath())
		defer inFile.Close()
		if err != nil {
			return err
		}

		outfile_path := filepath.Join(self.buildPath, "config/includes.binary/isolinux", value.Name())
		os.MkdirAll(filepath.Dir(outfile_path), 0777)
		outFile, err := os.OpenFile(outfile_path, os.O_CREATE|os.O_WRONLY, 0644)
		defer outFile.Close()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Added %s splash to %s\n", value.Name(), outfile_path)
		log.Println(msg)
		self.updateChannel <- LogUpdate{
			Append:     true,
			Message:    msg,
			UpdateType: UPDATE,
		}
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
		inFile, err := os.Open(value.FullPath())
		defer inFile.Close()
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(inFile)
		var outfileIdentifier string
		if scanner.Scan() {
			outfileIdentifier = scanner.Text()
		}

		outfile_path := filepath.Join(self.buildPath, outfileIdentifier)
		os.MkdirAll(filepath.Dir(outfile_path), 0777)
		outFile, err := os.OpenFile(outfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer outFile.Close()
		if err != nil {
			return err
		}

		for scanner.Scan() {
			outFile.WriteString(scanner.Text() + "\n")
		}

		msg := fmt.Sprintf("Added %s file to %s\n", value.Name(), outfile_path)
		self.updateChannel <- LogUpdate{
			Append:     true,
			Message:    msg,
			UpdateType: UPDATE,
		}
	}
	return nil
}
