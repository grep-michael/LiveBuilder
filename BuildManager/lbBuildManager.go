package buildmanager

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type CommandSpec struct {
	Name string
	Args []string
}

func (cs CommandSpec) CreateCmd() *exec.Cmd {
	return exec.Command(cs.Name, cs.Args...)
}

type BuildList []CommandSpec

var FullBuildList = BuildList{
	{Name: "lb", Args: []string{"bootstrap"}},
	{Name: "lb", Args: []string{"chroot"}},
	{Name: "lb", Args: []string{"binary"}},
	{Name: "lb", Args: []string{"installer"}},
	{Name: "lb", Args: []string{"source"}},
}

var FilesystemBuild = BuildList{
	{Name: "lb", Args: []string{"bootstrap"}},
	{Name: "lb", Args: []string{"chroot"}},
	{Name: "lb", Args: []string{"binary_chroot"}},
	{Name: "lb", Args: []string{"binary_rootfs"}},
	{Name: "lb", Args: []string{"binary_linux-image"}},
	{Name: "lb", Args: []string{"binary_manifest"}},
}
var JustConfig = BuildList{}

var DefaultBuildList = "FullBuild"

var BuildListMap = map[string]BuildList{
	"FullBuild":       FullBuildList,
	"FilesystemBuild": FilesystemBuild,
	"JustConfigure":   JustConfig,
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
	for _, commandSpec := range commands {
		command := commandSpec.CreateCmd()
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
