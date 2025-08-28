package buildmanager

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type BuildList []*exec.Cmd

var FullBuildList = BuildList{
	exec.Command("lb", []string{"bootstrap"}...),
	exec.Command("lb", []string{"chroot"}...),
	exec.Command("lb", []string{"binary"}...),
	exec.Command("lb", []string{"installer"}...),
	exec.Command("lb", []string{"source"}...),
}
var FilesystemBuild = BuildList{
	exec.Command("lb", []string{"bootstrap"}...),
	exec.Command("lb", []string{"chroot"}...),
	exec.Command("lb", []string{"binary_chroot"}...),
	exec.Command("lb", []string{"binary_rootfs"}...),
	//exec.Command("lb", []string{"binary_linux-image"}...),
}

var DefaultBuildList = "FullBuild"

var BuildListMap = map[string]BuildList{
	"FullBuild":       FullBuildList,
	"FilesystemBuild": FilesystemBuild,
}

type LBBuildManager struct {
	buildPath     string
	updateChannel chan LogUpdate
}

func NewLBBuildManager(updateChannel chan LogUpdate) *LBBuildManager {
	return &LBBuildManager{
		updateChannel: updateChannel,
	}
}

func (self *LBBuildManager) SetBuildPath(buildPath string) {
	self.buildPath = buildPath
}

func (self *LBBuildManager) BuildConditional(commands BuildList) error {

	if self.buildPath == "" {
		return fmt.Errorf("buildPath Not set")
	}

	self.updateChannel <- LogUpdate{
		Append:  false,
		Message: "Running lb build",
	}

	cmdOutChan := make(chan CommandOut, 20)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for cmdOut := range cmdOutChan {
			logUpdate := self.transformToLogUpdate(cmdOut)
			select {
			case self.updateChannel <- logUpdate:
			default:
				log.Println("Warning: GUI update channel is full, dropping message")
			}
		}
	}()
	for _, command := range commands {
		command.Dir = self.buildPath
		err := executeCommand(command, cmdOutChan)
		if err != nil {
			close(cmdOutChan)
			wg.Wait()
			return err
		}
	}

	self.updateChannel <- LogUpdate{
		Append:  true,
		Message: "lb build finished!",
	}

	close(cmdOutChan)
	wg.Wait()
	return nil

}

func (self *LBBuildManager) Build() error {
	err := self.BuildConditional(BuildListMap[DefaultBuildList])
	return err
}

func (self *LBBuildManager) transformToLogUpdate(cmdout CommandOut) LogUpdate {
	var msg string
	if cmdout.OutType == STDERR {
		msg = fmt.Sprintf("STD Error: %s\n", cmdout.OutPut)
	} else {
		msg = fmt.Sprintf("%s\n", cmdout.OutPut)
	}
	return LogUpdate{
		Append:  true,
		Message: msg,
	}
}
