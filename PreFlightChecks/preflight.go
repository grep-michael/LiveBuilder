package preflightchecks

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CheckAll(exitOnError bool) {
	funcs := []func() error{
		CheckLBversion,
		CheckCommands,
	}

	for _, fun := range funcs {
		err := fun()
		if err != nil {
			fmt.Println(err)
			log.Println(err)
			if exitOnError {
				os.Exit(1)
			}
		}
	}

}

func CheckCommands() error {
	commands := []string{
		"lb",
		"grub-installer",
		"mkfs.vfat",
		"mkfs.ext4",
		"parted",
		"sfdisk",
	}
	var missingCommands []string
	for _, command := range commands {
		cmd := exec.Command("which", command)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			missingCommands = append(missingCommands, command)
		}
	}
	if len(missingCommands) > 0 {
		return fmt.Errorf("Missing required packages (may require root to use): %v\n", strings.Join(missingCommands, ", "))
	}
	return nil
}

const (
	LB_VERSION = "20250505"
)

func CheckLBversion() error {
	var outbuf bytes.Buffer
	cmd := exec.Command("lb", "--version")
	cmd.Stdout = &outbuf
	if err := cmd.Run(); err != nil {
		return err
	}

	vers := strings.TrimSpace(outbuf.String())
	if vers != LB_VERSION {
		log.Printf("Untested lb version\ntested version: %s, installed version %s\n", LB_VERSION, vers)
		return fmt.Errorf("Untested lb version\ntested version: %s, installed version %s\n", LB_VERSION, vers)
	}
	return nil
}
