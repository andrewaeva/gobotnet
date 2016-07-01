package prgCmd

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"os/user"
)

func CmdTest() {
	fmt.Println("test prgCmd")
	fmt.Println(getName())
	fmt.Println(getOS())
	fmt.Println(getUid())
	fmt.Println(getUsername())
	fmt.Println(getHomeDir())
}

func CmdExec(cmd string) ([]byte, error) {
	return exec.Command("cmd", "/C", cmd).Output()
}

func getUsername() string {
	usr, _ := user.Current()
	return usr.Username
}

func getUid() string {
	usr, _ := user.Current()
	return usr.Uid
}

func getHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func getName() string {
	usr, _ := user.Current()
	return usr.Name
}

func getScreenshot() {

}

func getOS() string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows NT\CurrentVersion`, registry.READ)
	if err != nil {
		return ""
	}
	value, _, err := key.GetStringValue("ProductName")
	if err != nil {
		return ""
	}
	return value
}
