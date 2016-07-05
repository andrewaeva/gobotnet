package gobotnet

import (
	"errors"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
)

func CheckFileExist(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0777)
}

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
