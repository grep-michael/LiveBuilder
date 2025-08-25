package Utils

import (
	"bytes"
	"os/exec"
)

type stderr string
type stdout string

func Run_command(cmd string, args ...string) (stdout, stderr, error) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = &out
	c.Stderr = &err
	if err := c.Run(); err != nil {
		return "", "", err
	}
	return stdout(out.String()), stderr(err.String()), nil
}
