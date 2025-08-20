package buildmanager

/*
Actually executes lb build --verbose --debug
*/

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

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

func (self *LBBuildManager) Build() error {
	if self.buildPath == "" {
		return fmt.Errorf("buildPath Not set")
	}

	self.updateChannel <- LogUpdate{
		Append:  false,
		Message: "Running lb build",
	}

	build_command := self.makeBuildCommand()

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
	err := executeCommand(build_command, cmdOutChan)
	close(cmdOutChan)
	wg.Wait()

	self.updateChannel <- LogUpdate{
		Append:  true,
		Message: "lb build finished!",
	}

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

func (self *LBBuildManager) makeBuildCommand() *exec.Cmd {
	//cmd := exec.Command("lb", []string{"build", "--verbose", "--debug"}...)
	cmd := exec.Command("lb", []string{"build"}...)
	cmd.Dir = self.buildPath
	return cmd
}
