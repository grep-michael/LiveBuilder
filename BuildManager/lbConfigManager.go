package buildmanager

import (
	appstate "LiveBuilder/AppState"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type LBConfigManager struct {
	buildPath     string
	updateChannel chan LogUpdate
}

func NewLBConfigManager(updateChannel chan LogUpdate) *LBConfigManager {
	return &LBConfigManager{
		updateChannel: updateChannel,
	}
}

func (self *LBConfigManager) SetBuildPath(buildPath string) {
	self.buildPath = buildPath
}

func (self *LBConfigManager) ConfigureLB() error {
	if self.buildPath == "" {
		return fmt.Errorf("buildPath Not set")
	}
	lb_config_command, err := self.parseLBCommand()
	if err != nil {
		return err
	}
	cmdOutChan := make(chan CommandOut, 100)
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
	err = executeCommand(lb_config_command, cmdOutChan)
	close(cmdOutChan)
	wg.Wait()

	return err

}

func (self *LBConfigManager) transformToLogUpdate(cmdout CommandOut) LogUpdate {
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

func (self *LBConfigManager) parseLBCommand() (*exec.Cmd, error) {

	tokens := parseShellCommand(appstate.GetGlobalState().LBConfigCMD)
	if len(tokens) == 0 {
		log.Println("empty command")
		return nil, fmt.Errorf("empty command")
	}
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Dir = self.buildPath

	return cmd, nil
}
