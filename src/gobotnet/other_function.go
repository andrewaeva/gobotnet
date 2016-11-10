package gobotnet

import (
	"encoding/json"
	"fmt"
	"github.com/aglyzov/charmap"
	"math/rand"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	letterRunes      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	debugMode   bool = true
)

func IsDebugModeEnable() bool {
	return debugMode
}

func SetDebugMode(mode bool) {
	debugMode = mode
}

func CmdExec(cmd string) ([]byte, error) {

	cmd_li := exec.Command("cmd", "/C")
	cmd_li.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //Это необходимо для того что бы CMD запускалось в скрытом режиме
	output, err := cmd_li.Output()

	if output != nil && len(output) > 0 {
		return charmap.CP866_to_UTF8(output), nil
	} else {
		return output, err
	}
}

func CmdExecOrig(cmd string) ([]byte, error) {
	cmd_split := strings.Split(cmd, "@")
	params := make([]string, len(cmd_split)+2)
	params[0] = "/Q"
	params[1] = "/C"
	copy(params[2:], cmd_split[:])

	cmd_li := exec.Command("cmd", params...)
	cmd_li.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //Это необходимо для того что бы CMD запускалось в скрытом режиме
	output, err := cmd_li.Output()

	if output != nil && len(output) > 0 {
		return charmap.CP866_to_UTF8(output), nil
	} else {
		return output, err
	}
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

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Append(slice []byte, elements ...byte) []byte {
	n := len(slice)
	total := len(slice) + len(elements)
	if total > cap(slice) {
		// Reallocate. Grow to 1.5 times the new size, so we can still grow.
		newSize := total*3/2 + 1
		newSlice := make([]byte, total, newSize)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[:total]
	copy(slice[n:], elements)
	return slice
}

func ParseJsonResponse(jsonStruct interface{}, parseStr []byte) (bool, error) {
	err := json.Unmarshal(parseStr, jsonStruct)
	if err == nil {
		return true, err
	} else {
		return false, err
	}
}

func JoinToString(data []string) string {
	var info string
	for index := range data {
		info += data[index]
	}
	return info
}
