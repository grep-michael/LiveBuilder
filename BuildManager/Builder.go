package buildmanager

/*
Master object for doing all the backend building
1. imports all the custom files
2. runs the nessacary live-build commands
3. formats resulting files our specific desired output
*/

import (
	filesystem "LiveBuilder/Filesystem"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type UpdateType string

const (
	START  UpdateType = "start"
	UPDATE UpdateType = "start"
	END    UpdateType = "start"
)

type LogUpdate struct {
	UpdateType
	Message string
	Append  bool // true to append, false to replace
}

type BuildManager struct {
	updateChannel   chan LogUpdate
	stepFinished    chan bool
	importer        *Importer
	lbconfigManager *LBConfigManager
	lbBuildManager  *LBBuildManager
	buildPath       string
	subscribers     []chan LogUpdate
	subMutex        sync.RWMutex
}

func NewBuilder() *BuildManager {
	builder := &BuildManager{
		updateChannel: make(chan LogUpdate, 100),
		stepFinished:  make(chan bool),
	}
	builder.importer = NewImporter(builder.updateChannel)
	builder.lbconfigManager = NewLBConfigManager(builder.updateChannel)
	builder.lbBuildManager = NewLBBuildManager(builder.updateChannel)
	go builder.listenForUpdates()
	return builder
}

func (self *BuildManager) GetDefaultBuildPath() string {
	path, err := os.MkdirTemp("", "LiveBuilder-*")
	if err != nil {
		log.Println(err)
		appdata, _ := filesystem.GetAppDataDir()
		buildpath := filepath.Join(appdata, "build")
		return buildpath
	}
	return path

}

func (self *BuildManager) Build(buildPath string) {

	if err := self.InitializeBuildPath(buildPath); err != nil {
		self.updateChannel <- LogUpdate{
			Append:  false,
			Message: fmt.Sprintf("Error occured initializing build path: %v\n", err),
		}
		return
	}
	log.Printf("Building to path: %s\n", self.buildPath)
	self.importer.SetBuildPath(self.buildPath)
	self.lbconfigManager.SetBuildPath(self.buildPath)
	self.lbBuildManager.SetBuildPath(self.buildPath)

	if err := self.lbconfigManager.ConfigureLB(); err != nil {
		self.updateChannel <- LogUpdate{
			Append:  true,
			Message: fmt.Sprintf("Error occured in configuring LB: %v\n", err),
		}
		return
	}
	if err := self.importer.ImportAll(); err != nil {
		self.updateChannel <- LogUpdate{
			Append:  true,
			Message: fmt.Sprintf("Error occured in importer: %v\n", err),
		}
		return
	}
	if err := self.lbBuildManager.Build(); err != nil {
		self.updateChannel <- LogUpdate{
			Append:  true,
			Message: fmt.Sprintf("Error occured in lb build: %v\n", err),
		}
		return
	}
	if err := self.copyISO(); err != nil {
		self.updateChannel <- LogUpdate{
			Append:  true,
			Message: fmt.Sprintf("Error occured copying iso file: %v\n", err),
		}
		return
	}

}

func (self *BuildManager) copyISO() error {
	//create folder for iso to be copied to
	appdata, _ := filesystem.GetAppDataDir()
	iso_path := filepath.Join(appdata, filesystem.ISO_DIR_ID)
	if err := os.MkdirAll(iso_path, 0777); err != nil {
		return err
	}
	//get iso filepath
	var iso_files []string
	filepath.WalkDir(self.buildPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".iso" {
			iso_files = append(iso_files, s)
		}
		return nil
	})
	for _, iso_file := range iso_files {
		dest := filepath.Join(iso_path, filepath.Base(iso_file))
		self.updateChannel <- LogUpdate{
			Append:  true,
			Message: fmt.Sprintf("Copying:%s -> %s\n", iso_file, dest),
		}

		err := copyFile(iso_file, dest)
		if err != nil {
			log.Printf("Error copying iso file %s\n", err)
			return err
		}
	}
	return nil
}

func (self *BuildManager) GetSubscriber() <-chan LogUpdate {
	self.subMutex.Lock()
	defer self.subMutex.Unlock()

	subscriber := make(chan LogUpdate, 100)
	self.subscribers = append(self.subscribers, subscriber)
	return subscriber
}

func (self *BuildManager) listenForUpdates() {
	for update := range self.updateChannel {
		self.subMutex.RLock()
		for _, subscriber := range self.subscribers {
			select {
			case subscriber <- update:
			default:
			}
		}
		self.subMutex.RUnlock()
	}
}

func (self *BuildManager) InitializeBuildPath(buildPath string) error {
	if buildPath == "" {
		self.buildPath = self.GetDefaultBuildPath()
		log.Printf("Using default build directory: %s\n", self.buildPath)
	} else {
		self.buildPath = buildPath
	}

	if err := self.NukeBuild(); err != nil {
		return err
	}
	if err := os.MkdirAll(self.buildPath, 0777); err != nil {
		return err
	}
	return nil
}

func (self *BuildManager) NukeBuild() error {
	if err := os.RemoveAll(self.buildPath); err != nil {
		return err
	}
	return nil
	cmd := exec.Command("lb", "clean")
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}
	return cmd.Wait()
}

func copyFile(src, dst string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close() // Ensure the source file is closed

	// Create the destination file for writing (creates if not exists, truncates if exists)
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		// Close the destination file and handle potential errors during close
		cerr := destinationFile.Close()
		if err == nil { // Only update err if no previous error occurred
			err = cerr
		}
	}()

	// Copy the contents from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Ensure all data is written to disk
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
