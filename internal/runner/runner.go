package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(cmd []string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	
	fmt.Printf("🔄 Running: %s\n", joinCmd(cmd))
	return c.Run()
}

func joinCmd(cmd []string) string {
	return "'" + strings.Join(cmd, "' '") + "'"
}
