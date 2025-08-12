package execution

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	"bufio"
	"fmt"
	"fyne.io/fyne/v2/widget"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

type LiveBuilder struct {
	outputLabel *widget.Label
	state       appstate.State
	workingDir  string
}

func NewLiveBuilder(label *widget.Label) *LiveBuilder {
	return &LiveBuilder{
		outputLabel: label,
		state:       *appstate.GetGlobalState(),
	}
}

func (self *LiveBuilder) SetWorkingDir(dir string) {
	self.workingDir = dir
}

func (self *LiveBuilder) setDefaultDir() {
	appdata, _ := filesystem.GetAppDataDir()
	buildpath := filepath.Join(appdata, "build")
	self.workingDir = buildpath
}
func (self *LiveBuilder) GetBuildDir() string {
	if self.workingDir == "" {
		self.setDefaultDir()
	}
	os.MkdirAll(self.workingDir, 0777)
	return self.workingDir
}

func (self *LiveBuilder) ConfigureLB() error {
	cmd, args, _ := self.parseLBCommand()
	return self.executeCommand(cmd, args)
}

func (self *LiveBuilder) executeCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = self.GetBuildDir()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Channel to collect output
	outputChan := make(chan string, 100)
	doneChan := make(chan bool)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outputChan <- scanner.Text()
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outputChan <- "ERROR: " + scanner.Text()
		}
	}()

	// Update label with output
	go func() {
		var allOutput []string
		for line := range outputChan {
			allOutput = append(allOutput, line)
			finalText := strings.Join(allOutput, "\n")
			self.outputLabel.SetText(finalText)
			self.outputLabel.Refresh()
		}
		doneChan <- true
	}()

	err = cmd.Wait()
	close(outputChan)
	<-doneChan
	return err
}

func (self *LiveBuilder) parseLBCommand() (string, []string, error) {
	// Remove line continuations (backslash followed by newline)
	cleaned := strings.ReplaceAll(self.state.LBConfigCMD, "\\\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\\\r\n", " ") // Windows line endings

	// Normalize whitespace
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	// Parse into tokens
	tokens := self.parseShellCommand(cleaned)

	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}

	return tokens[0], tokens[1:], nil
}

func (self *LiveBuilder) parseShellCommand(input string) []string {
	var tokens []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	for _, char := range input {
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteRune(char)
			}

		case unicode.IsSpace(char) && !inQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(char)
		}
	}

	// Add the last token if there is one
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
