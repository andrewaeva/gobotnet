package cmd

import (
	"os/exec"
)
import logging "../logging"

func CmdExec(cmd string) ([]byte, error) {
	logging.DebugLogging(cmd)
	return exec.Command("cmd", "/C", cmd).Output()
}
