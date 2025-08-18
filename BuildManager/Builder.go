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
	"log"
	"os"
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
	appdata, _ := filesystem.GetAppDataDir()
	buildpath := filepath.Join(appdata, "build")
	return buildpath
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
}
