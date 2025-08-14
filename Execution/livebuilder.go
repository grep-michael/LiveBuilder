package execution

import (
	appstate "LiveBuilder/AppState"
	filesystem "LiveBuilder/Filesystem"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type LogUpdate struct {
	Message string
	Append  bool // true to append, false to replace
}

type LiveBuilder struct {
	state       *appstate.State
	workingDir  string
	subscribers []chan LogUpdate
	subMutex    sync.RWMutex
}

func NewLiveBuilder() *LiveBuilder {
	return &LiveBuilder{
		state:       appstate.GetGlobalState(),
		subscribers: make([]chan LogUpdate, 0),
	}
}

// GetSubscriber returns a new channel that will receive log updates
func (self *LiveBuilder) GetSubscriber() <-chan LogUpdate {
	self.subMutex.Lock()
	defer self.subMutex.Unlock()

	subscriber := make(chan LogUpdate, 100)
	self.subscribers = append(self.subscribers, subscriber)
	return subscriber
}

func (self *LiveBuilder) publishLog(message string, append bool) {
	update := LogUpdate{
		Message: message,
		Append:  append,
	}

	//self.subMutex.RLock()
	//defer self.subMutex.RUnlock()
	self.subMutex.Lock()
	defer self.subMutex.Unlock()

	for _, subscriber := range self.subscribers {
		select {
		case subscriber <- update:
		default:
			// Channel is full, skip this subscriber to prevent blocking
		}
	}
}

func (self *LiveBuilder) logReplace(message string) {
	self.publishLog(message, false)
}

func (self *LiveBuilder) logAppend(message string) {
	self.publishLog(message, true)
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

func (self *LiveBuilder) BuildLB() {
	cmd := "lb"
	args := []string{"build", "--verbose", "--debug"}
	err := self.executeCommand(cmd, args)
	if err != nil {
		log.Printf("BuildLB cmd.wait returned error: %v\n", err)
	} else {
		log.Printf("BuildLB cmd.wait return no error, finished!")
	}
	self.logReplace("LB Build Ended!")
}

func (self *LiveBuilder) DropPackages() {
	self.logReplace("Dropping Packages\n")
	packageMap := self.state.GetDirectoryEntryMap(filesystem.PACKAGE_DIR_ID)

	outfile_path := filepath.Join(self.GetBuildDir(), "config/package-lists/live.list.chroot")
	outFile, _ := os.OpenFile(outfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer outFile.Close()

	for _, value := range packageMap {
		inFile, _ := os.Open(value.FullPath())
		defer inFile.Close()

		io.Copy(outFile, inFile)
		outFile.WriteString("\n")
		msg := fmt.Sprintf("Added %s package to config/package-lists/live.list.chroot\n", value.Name())
		self.logAppend(msg)
	}
}

func (self *LiveBuilder) DropSplashImages() {
	self.logReplace("Dropping Splash images\n")
	log.Println("Dropping Splash images")
	splashMap := filesystem.GetFileManager().GetFileSystem(filesystem.SPLASH_SCREENS_ID)
	for _, value := range splashMap {
		inFile, _ := os.Open(value.FullPath())
		defer inFile.Close()

		outfile_path := filepath.Join(self.GetBuildDir(), "config/includes.binary/isolinux", value.Name())
		os.MkdirAll(filepath.Dir(outfile_path), 0777)
		outFile, _ := os.OpenFile(outfile_path, os.O_CREATE|os.O_WRONLY, 0644)
		defer outFile.Close()

		io.Copy(outFile, inFile)
		msg := fmt.Sprintf("Added %s splash to %s\n", value.Name(), outfile_path)
		log.Println(msg)
		self.logAppend(msg)
	}
}

func (self *LiveBuilder) DropCustomFiles() {
	self.logReplace("Dropping custom files\n")
	customFileMap := self.state.GetDirectoryEntryMap(filesystem.CUSTOMFILES_DIR_ID)

	for _, value := range customFileMap {
		inFile, _ := os.Open(value.FullPath())
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		var outfileIdentifier string
		if scanner.Scan() {
			outfileIdentifier = scanner.Text()
		}

		outfile_path := filepath.Join(self.GetBuildDir(), outfileIdentifier)
		os.MkdirAll(filepath.Dir(outfile_path), 0777)
		outFile, _ := os.OpenFile(outfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer outFile.Close()

		for scanner.Scan() {
			outFile.WriteString(scanner.Text() + "\n")
		}

		msg := fmt.Sprintf("Added %s file to %s\n", value.Name(), outfile_path)
		self.logAppend(msg)
	}
}

func (self *LiveBuilder) NukeBuild() {
	err := os.RemoveAll(self.GetBuildDir())
	if err != nil {
		self.logAppend(fmt.Sprintf("Error nuking build: %v\n", err))
	}
}

func (self *LiveBuilder) executeCommand(command string, args []string) error {
	log.Println("executing command")
	log.Println(command)
	log.Println(args)

	cmd := exec.Command(command, args...)
	cmd.Dir = self.GetBuildDir()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error occured getting stdout pipe: %v\n", err)
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error occured getting stderr pipe: %v\n", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return err
	}

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			self.logAppend(scanner.Text() + "\n")
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			self.logAppend("CMD Error: " + scanner.Text() + "\n")
		}
	}()

	return cmd.Wait()
}

func (self *LiveBuilder) parseLBCommand() (string, []string, error) {
	cleaned := strings.ReplaceAll(self.state.LBConfigCMD, "\\\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\\\r\n", " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	tokens := self.parseShellCommand(cleaned)

	if len(tokens) == 0 {
		log.Println("empty command")
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

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
