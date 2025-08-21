package preflightchecks

import (
	"fmt"
	"testing"
)

func TestCommanCheck(t *testing.T) {
	err := CheckCommands()
	fmt.Println(err)
}
func TestLBVersion(t *testing.T) {
	err := CheckLBversion()
	fmt.Println(err)
}
