package buildmanager

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	//"path/filepath"
	"sync"
	"text/template"
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

	selectedCommandTemplate := appstate.GetGlobalState().GetDirectoryEntryMap(filesystem.LBCONFIGS_DIR_ID)

	if len(selectedCommandTemplate) != 1 {
		return nil, fmt.Errorf("Incorrect number of lb configs selected, must be only 1")
	}

	var lb_cmd string
	var err error
	for _, val := range selectedCommandTemplate {
		lb_cmd, err = self.loadLBConfigTemplate(val)
		if err != nil {
			return nil, err
		}
	}

	tokens := parseShellCommand(lb_cmd)
	if len(tokens) == 0 {
		log.Println("empty command")
		return nil, fmt.Errorf("empty command")
	}
	cmd := exec.Command(tokens[0], tokens[1:]...)
	cmd.Dir = self.buildPath

	return cmd, nil
}

func (self *LBConfigManager) loadLBConfigTemplate(dir filesystem.DirectoryEntry) (string, error) {
	state := appstate.GetGlobalState()

	content, err := os.ReadFile(dir.FullPath())
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	tmpl, err := template.New("lbconfig").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, state.LBcfg); err != nil {
		return "", err
	}
	str := buf.String()
	log.Printf("Build lb config from template: %s\n", str)
	return str, nil
}
